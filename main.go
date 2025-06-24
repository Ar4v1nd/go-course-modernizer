package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/genai"
)

type VideoItem struct {
	Snippet struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Thumbnails  struct {
			Standard struct {
				URL string `json:"url"`
			} `json:"standard"`
		} `json:"thumbnails"`
		Position int `json:"position"`
	} `json:"snippet"`
	ContentDetails struct {
		VideoId          string `json:"videoId"`
		VideoPublishedAt string `json:"videoPublishedAt"`
	} `json:"contentDetails"`
}

type PlaylistItems struct {
	Items         []VideoItem `json:"items"`
	NextPageToken string      `json:"nextPageToken"`
}

var baseUrl = "https://www.googleapis.com/youtube/v3/playlistItems"
var part = []string{"snippet", "contentDetails"} // Add more if required
var playlistId = "PLoILbKo9rG3skRCj37Kn5Zj803hhiuRK6"
var releaseNotes []*genai.File

func uploadReleaseNotes(ctx context.Context, logger *slog.Logger, client *genai.Client) error {
	logger.Info("Uploading release notes to Gemini")

	releaseNotesDir := "./releasenote"
	// Walk the release notes directory
	err := filepath.Walk(releaseNotesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Error("Error accessing path", "path", path, "error", err)
			return fmt.Errorf("Error accessing path %q: %v", path, err)
		}
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".pdf") {
			return nil // Skip directories and non-pdf files
		}
		// Upload the release notes to Gemini
		f, err := client.Files.UploadFromPath(ctx, path, &genai.UploadFileConfig{
			MIMEType: "application/pdf",
		})
		if err != nil {
			logger.Error("Error uploading release notes file to Gemini", "file", path, "error", err)
			return fmt.Errorf("Error uploading release notes file from %q to Gemini: %v", path, err)
		}
		logger.Info("Successfully uploaded release notes file", "file", path)
		releaseNotes = append(releaseNotes, f)
		return nil
	})
	if err != nil {
		logger.Error("Failed to upload release notes", "error", err)
		return fmt.Errorf("Failed to upload release notes: %w", err)
	}
	logger.Info("Successfully uploaded all release notes", "count", len(releaseNotes))
	return nil
}

func processVideo(ctx context.Context, wg *sync.WaitGroup, logger *slog.Logger, client *genai.Client, limits chan struct{}, results chan<- map[string]string, item VideoItem) error {
	defer wg.Done()

	limits <- struct{}{}        // Acquire a limit
	defer func() { <-limits }() // Release the limit when done

	summarizerPrompt := `
	You are an expert in Go programming.

	You will be given a YouTube video URL of a Go programming course recorded with Go version 1.15.
	
	Your task is to summarize the video by following these guidelines:
	1. Dissect the video content into distinct chapters based on the topics covered.
	2. For each chapter, summarize the key concepts and best practices as concise bullet points:
		- Include relevant Go code snippets.
		- Do not include video timestamps or references to specific moments in the video.
	3. Return your response in the following strict **Markdown format only**, with no additional text:
	# ` + item.Snippet.Title + `

	## Summary
	(A brief overview of the video content.)

	## Key Points
	(A list of chapters with their summaries in concise bullet points. Include relevant Go code snippets or examples.)
	`

	logger.Info("Sending request to Gemini for video summarization", "videoId", item.ContentDetails.VideoId, "title", item.Snippet.Title)

	parts := []*genai.Part{
		genai.NewPartFromURI(fmt.Sprintf("https://www.youtube.com/watch?v=%s", item.ContentDetails.VideoId), "video/mp4"),
		genai.NewPartFromText(summarizerPrompt),
	}

	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	temperature := float32(0.1)
	thinkingBudget := int32(-1) // Set to -1 for dynamic thinking
	response, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		contents,
		&genai.GenerateContentConfig{
			Temperature:        &temperature,
			ResponseModalities: []string{"TEXT"},
			ThinkingConfig: &genai.ThinkingConfig{
				ThinkingBudget: &thinkingBudget,
			},
		},
	)
	if err != nil {
		logger.Error("Failed to summarize video using Gemini", "videoId", item.ContentDetails.VideoId, "error", err)
		return fmt.Errorf("Failed to summarize video using Gemini: %w", err)
	}

	logger.Info("Received summarization response from Gemini", "videoId", item.ContentDetails.VideoId, "response", len(response.Text()))

	usageMetadata, err := json.MarshalIndent(response.UsageMetadata, "", "  ")
	if err != nil {
		logger.Error("Failed to marshal usage metadata", "videoId", item.ContentDetails.VideoId, "error", err)
	} else {
		// Try to unmarshal usageMetadata into a map to access token counts
		var metaMap map[string]any
		if err := json.Unmarshal(usageMetadata, &metaMap); err == nil {
			if thoughts, ok := metaMap["thoughtsTokenCount"]; ok {
				logger.Info("Thoughts tokens", "count", thoughts)
			}
			if candidates, ok := metaMap["candidatesTokenCount"]; ok {
				logger.Info("Output tokens", "count", candidates)
			}
		} else {
			logger.Info("Usage metadata is not a valid JSON", "usageMetadata", string(usageMetadata))
		}
	}

	// Optional: Sleep for a short duration to avoid hitting API rate limits
	// time.Sleep(1 * time.Minute)

	validatorPrompt := `
	You are a technical content editor who is an expert in Go programming.

	You will be given the summary and key points in Markdown format derived from a Go programming course recorded with Go version 1.15.

	Your task is to evaluate each key point present under the "Key Points" section of the Markdown using the provided Go release notes PDF files (from versions 1.16 to 1.24) by following these guidelines:
	1. Determine if every key point is **still valid and accurate** in the latest Go version (1.24) based on the release notes.
		- **Only** consider the following sections in the release note PDF files while evaluating the key points: "Changes to the language", "Tools" and "Standard library". Ignore any other sections.
		- Do **not** evaluate key points expressing opinions, philosophies, or general design principles.
		- Only focus on factual key points about Go syntax, behavior, deprecation, tooling, etc.
	2. For any key point that is no longer valid or accurate:
		- Briefly explain what has changed in the latest Go version that affects the key point.
		- Cite the **first Go version** where the change was introduced using a numbered format like [1], [2], etc.
		- Do not cite a Go version unless it is directly relevant to the key point. Also, do not cite multiple versions for the same change (choose the most relevant one).
		- Provide updated code snippets if the original code is outdated.
	3. Do not use any prior knowledge about Go. Only base your answers on the provided release note PDFs.
	4. Return your response in the following strict **Markdown format only**, with no additional text:
	# ` + item.Snippet.Title + `

	## Summary
	(Summary passed to you as input, do not change it.)

	## Key Points
	(Key points passed to you as input, do not change them.)

	## What's New
	(A list of changes found in the key points based on the release notes, with each change cited to the relevant Go version in [x] numbered format.)

	## Updated Code Snippets
	(If any code snippets in the key points were outdated, provide the updated versions here. If no updated code snippets are needed, omit this section entirely.)

	## Citations
	(A list of Go version release notes cited in the format [1], [2], etc. For example:
	- [1] Go version 1.16
	- [2] Go version 1.17
	)

	Here are the summary and key points in Markdown format you need to evaluate:
	` + response.Text()

	parts = []*genai.Part{}
	for i, file := range releaseNotes {
		parts = append(parts, genai.NewPartFromText(fmt.Sprintf("[%d] %s", i+1, file.Name)), genai.NewPartFromURI(file.URI, file.MIMEType))
	}
	parts = append(parts, genai.NewPartFromText(validatorPrompt))

	contents = []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	temperature = float32(0.0)    // Set temperature to 0 for validation
	thinkingBudget = int32(24576) // Set a high thinking budget for thorough validation
	response, err = client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		contents,
		&genai.GenerateContentConfig{
			Temperature:        &temperature,
			ResponseModalities: []string{"TEXT"},
			ThinkingConfig: &genai.ThinkingConfig{
				ThinkingBudget: &thinkingBudget,
			},
		},
	)
	if err != nil {
		logger.Error("Failed to validate using Gemini", "videoId", item.ContentDetails.VideoId, "error", err)
		return fmt.Errorf("Failed to validate using Gemini: %w", err)
	}

	logger.Info("Received validation response from Gemini", "videoId", item.ContentDetails.VideoId, "response", len(response.Text()))

	usageMetadata, err = json.MarshalIndent(response.UsageMetadata, "", "  ")
	if err != nil {
		logger.Error("Failed to marshal usage metadata", "videoId", item.ContentDetails.VideoId, "error", err)
	} else {
		// Try to unmarshal usageMetadata into a map to access token counts
		var metaMap map[string]any
		if err := json.Unmarshal(usageMetadata, &metaMap); err == nil {
			if thoughts, ok := metaMap["thoughtsTokenCount"]; ok {
				logger.Info("Thoughts tokens", "count", thoughts)
			}
			if candidates, ok := metaMap["candidatesTokenCount"]; ok {
				logger.Info("Output tokens", "count", candidates)
			}
		} else {
			logger.Info("Usage metadata is not a valid JSON", "usageMetadata", string(usageMetadata))
		}
	}

	results <- map[string]string{
		item.Snippet.Title: response.Text(),
	}

	return nil
}

func getPlaylistItems(ctx context.Context) ([]VideoItem, error) {
	apiKey, ok := os.LookupEnv("YOUTUBE_API_KEY")
	if !ok {
		return nil, fmt.Errorf("YOUTUBE_API_KEY environment variable is not set")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("YOUTUBE_API_KEY environment variable is empty")
	}

	var videoItems []VideoItem
	url := baseUrl + "?part=" + strings.Join(part, ",") + "&playlistId=" + playlistId + "&key=" + apiKey + "&maxResults=50"
	pageToken := ""

	for {
		var playlistItems PlaylistItems

		reqUrl := url
		if pageToken != "" {
			reqUrl += "&pageToken=" + pageToken
		}

		// Create a timeout context
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)

		// Make a GET request to the YouTube API
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqUrl, nil)
		if err != nil {
			cancel()
			return nil, fmt.Errorf("Failed to create request: %w", err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			cancel()
			return nil, fmt.Errorf("Failed to make request to YouTube API: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			cancel()
			return nil, fmt.Errorf("Received non-200 response from YouTube API: %s", resp.Status)
		}

		// Decode the JSON response into the PlaylistItems struct
		if err := json.NewDecoder(resp.Body).Decode(&playlistItems); err != nil {
			cancel()
			return nil, fmt.Errorf("Failed to decode JSON response: %w", err)
		} else {
			cancel()
		}

		// Append the items to the videoItems slice
		videoItems = append(videoItems, playlistItems.Items...)

		if playlistItems.NextPageToken == "" {
			break // No more pages to fetch
		} else {
			pageToken = playlistItems.NextPageToken // Update the page token for the next iteration
		}
	}

	return videoItems, nil
}

func main() {
	// Initialize logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	if err := godotenv.Load(); err != nil {
		logger.Error("Error loading .env file", "error", err)
		os.Exit(-1)
	}

	logger.Info("Getting the list of video URLs from the playlist")

	ctx := context.Background()

	// Fetch playlist items from YouTube API
	playlistItems, err := getPlaylistItems(ctx)
	if err != nil {
		logger.Error("Error fetching playlist items", "error", err)
		os.Exit(-1)
	}
	logger.Info("Successfully fetched playlist items", "count", len(playlistItems))

	apiKey, ok := os.LookupEnv("GEMINI_API_KEY")
	if !ok {
		logger.Error("GEMINI_API_KEY environment variable is not set")
		os.Exit(-1)
	}
	if apiKey == "" {
		logger.Error("GEMINI_API_KEY environment variable is empty")
		os.Exit(-1)
	}

	// Create a Gemini client
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		logger.Error("Failed to create Gemini client", "error", err)
		os.Exit(-1)
	}

	// Upload release notes to Gemini
	if err := uploadReleaseNotes(ctx, logger, client); err != nil {
		logger.Error("Error uploading release notes to Gemini", "error", err)
		os.Exit(-1)
	}

	wg := new(sync.WaitGroup)
	limits := make(chan struct{}, 5) // Limit to 5 concurrent request to handle Gemini API rate limits
	results := make(chan map[string]string, len(playlistItems))

	for _, item := range playlistItems {
		wg.Add(1)
		go processVideo(ctx, wg, logger, client, limits, results, item)
	}
	wg.Wait()
	close(results)

	// Ensure the "markdown" directory exists
	if err := os.MkdirAll("markdown", 0755); err != nil {
		logger.Error("Error creating markdown directory", "error", err)
		os.Exit(-1)
	}

	// Collect results
	for result := range results {
		for title, summary := range result {
			re := regexp.MustCompile(`/`) // Replace slashes in titles
			title = re.ReplaceAllString(title, "_")
			filePath := filepath.Join("markdown", fmt.Sprintf("%s.md", title))
			err := os.WriteFile(filePath, []byte(summary), 0644)
			if err != nil {
				logger.Error("Error writing summary to file", "title", title, "error", err)
			}
		}
	}
}
