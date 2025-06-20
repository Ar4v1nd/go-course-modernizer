package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
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
var once sync.Once
var releaseNotes []*genai.File

func uploadReleaseNotes(ctx context.Context, logger *slog.Logger, client *genai.Client) error {
	logger.Info("Uploading release notes to Gemini")

	releaseNotesDir := "./releasenotes"
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

func summarizeVideo(ctx context.Context, wg *sync.WaitGroup, logger *slog.Logger, client *genai.Client, limits chan struct{}, results chan<- map[string]string, item VideoItem) error {
	defer wg.Done()

	limits <- struct{}{}        // Acquire a limit
	defer func() { <-limits }() // Release the limit when done

	once.Do(func() {
		err := uploadReleaseNotes(ctx, logger, client)
		if err != nil {
			os.Exit(-1)
		}
	})

	prompt := `
	You're a Principal Go developer and an expert technical content editor.

	Given the above YouTube video URL from a Go programming course recorded with Go version 1.15, and the release notes for Go versions 1.16 to 1.24 (as PDF files), your task is to:
	1. Dissect the video content into distinct chapters based on the topics covered.

	2. For each chapter, summarize the key concepts and best practices as concise bullet points:
		- Include relevant Go code snippets.
		- Do not include video timestamps or references to specific moments in the video.

	3. Fact check each bullet point in your summary using *only* the content from the release note PDF files.

	4. For any bullet point that is outdated or inaccurate:
		- Briefly explain what changed.
		- Cite the specific release note PDF **filename (without extension)** using a numbered format like [1], [2], etc.
		- Always cite the **first release note file** where the change was introduced.
		- Do not cite any release note unless it is directly relevant to the concept discussed in the video.
		- Provide updated code snippets if the original code is outdated.

	5. Return your response in the following strict **Markdown format only**, with no additional text:
	# ` + item.Snippet.Title + `

	## Summary
	(A brief overview of the video content.)

	## Key Points
	(A list of chapters with their summaries in bullet points. Include relevant Go code snippets or examples.)

	## What has changed?
	(Only include changes that affect topics covered in the video. Each point must end with a numbered citation in [#] format. If no relevant changes occurred, state: “No significant changes since the video was recorded.”)

	## Updated Code Examples
	(Updated Go code from the video content that reflects modern usage. If nothing changed, say: “No updated code examples.”)

	## Citations
	(A list of **only the release note PDF files that were cited** in the 'What has changed?' section, each preceded by its citation number. For example:  
	- [1] Go 1.16 Release Notes
	- [2] Go 1.17 Release Notes
	If no citations were needed, say: “No citations needed.”)
	`

	parts := []*genai.Part{
		genai.NewPartFromURI(fmt.Sprintf("https://www.youtube.com/watch?v=%s", item.ContentDetails.VideoId), "video/mp4"),
	}
	for _, file := range releaseNotes {
		parts = append(parts, genai.NewPartFromURI(file.URI, file.MIMEType))
	}
	parts = append(parts, genai.NewPartFromText(prompt))

	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	temperature := float32(0.1)
	response, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		contents,
		&genai.GenerateContentConfig{
			Temperature:        &temperature,
			ResponseModalities: []string{"TEXT"},
		},
	)

	if err != nil {
		logger.Error("Failed to generate content using Gemini", "videoId", item.ContentDetails.VideoId, "error", err)
		return fmt.Errorf("Failed to generate content using Gemini: %w", err)
	}

	logger.Info("Received response from Gemini", "videoId", item.ContentDetails.VideoId, "summaryLength", len(response.Text()))

	results <- map[string]string{
		item.ContentDetails.VideoId: response.Text(),
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

	playlistItems, err := getPlaylistItems(ctx)
	if err != nil {
		logger.Error("Error fetching playlist items", "error", err)
		os.Exit(-1)
	}
	logger.Info("Successfully fetched playlist items", "count", len(playlistItems))

	wg := new(sync.WaitGroup)
	limits := make(chan struct{}, 5) // Limit to 5 concurrent requests
	results := make(chan map[string]string, len(playlistItems))

	apiKey, ok := os.LookupEnv("GEMINI_API_KEY")
	if !ok {
		logger.Error("GEMINI_API_KEY environment variable is not set")
		os.Exit(-1)
	}
	if apiKey == "" {
		logger.Error("GEMINI_API_KEY environment variable is empty")
		os.Exit(-1)
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		logger.Error("Failed to create Gemini client", "error", err)
		os.Exit(-1)
	}

	for _, item := range playlistItems[:6] { // Limit to first 6 items for testing
		wg.Add(1)
		go summarizeVideo(ctx, wg, logger, client, limits, results, item)
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
		for videoId, summary := range result {
			filePath := filepath.Join("markdown", fmt.Sprintf("%s.md", videoId))
			err := os.WriteFile(filePath, []byte(summary), 0644)
			if err != nil {
				logger.Error("Error writing summary to file", "videoId", videoId, "error", err)
			}
		}
	}
}
