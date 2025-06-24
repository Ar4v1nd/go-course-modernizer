# Go Class: 16 Homework #3

## Summary
This video guides viewers through solving a homework exercise from "The Go Programming Language" book, focusing on fetching and searching data from the xkcd webcomic API. The solution is broken down into two Go programs: a downloader that fetches comic metadata and stores it as a JSON array on disk, and a searcher that reads this local data to find comics matching specified keywords. The video emphasizes practical Go programming concepts such as HTTP requests, JSON encoding/decoding, file I/O, string manipulation, and robust error handling.

## Key Points

### 1. Problem Overview
*   **Objective:** Create two Go programs to interact with the xkcd webcomic API.
    *   **Program 1 (Downloader):** Fetch metadata for all comics and save it to a local JSON file.
    *   **Program 2 (Searcher):** Read the local JSON file and search for comics whose title or transcript matches a list of search terms provided on the command line.
*   **xkcd API:** Provides comic metadata (month, num, link, year, transcript, img, title, day) as JSON objects at URLs like `https://xkcd.com/num/info.0.json`.

### 2. First Program: Downloader (`xkcd-load.go`)
*   **Purpose:** Downloads comic metadata from the xkcd API and stores it as a single JSON array in a file.
*   **Design Choices:**
    *   Reads comics sequentially by number, starting from 1.
    *   Stops when two consecutive HTTP 404 (Not Found) responses are received, indicating no more comics.
    *   Each API request returns a JSON object as a string. These strings are concatenated into a valid JSON array format `[obj1, obj2, ..., objN]`.
    *   The program does *not* decode the JSON during download; it treats the JSON objects as raw strings to be written to the file.
    *   Optionally takes an output filename from the command line; otherwise, it prints to standard output.
*   **Key Go Concepts:**
    *   **HTTP Requests:** Uses `net/http` package to perform GET requests.
        ```go
        resp, err := http.Get(url)
        if err != nil {
            // Handle network error
        }
        defer resp.Body.Close() // Ensure response body is closed
        ```
    *   **Error Handling:** Checks for network errors and non-200 OK HTTP status codes.
        ```go
        if resp.StatusCode != http.StatusOK {
            // Handle non-OK status (e.g., 404)
            return nil // Indicate no data found
        }
        ```
    *   **File I/O:** Uses `os.Create` to open a file for writing and `io.WriteCloser` interface for flexible output (stdout or file).
        ```go
        var output io.WriteCloser = os.Stdout
        if len(os.Args) > 1 {
            output, err = os.Create(os.Args[1])
            if err != nil {
                fmt.Fprintln(os.Stderr, err)
                os.Exit(-1)
            }
        }
        defer output.Close() // Ensures the output file is closed
        ```
    *   **JSON Array Construction:** Manually prints `[` at the beginning, `]` at the end (using `defer`), and `,` between each JSON object.
        ```go
        fmt.Println("[")
        defer fmt.Println("]")
        // ... inside loop ...
        if cnt > 0 { // For all but the first object
            fmt.Print(output, ",")
        }
        io.Copy(output, bytes.NewBuffer(data)) // Write raw JSON string
        ```
    *   **Looping with Condition:** Uses a `for` loop with a condition to stop after consecutive failures (e.g., two 404s).
        ```go
        for i := 1; fails < 2; i++ {
            // ... fetch comic ...
            if data == nil {
                fails++
                continue
            }
            fails = 0 // Reset failures on success
            // ... write data ...
        }
        ```

### 3. Second Program: Searcher (`xkcd-find.go`)
*   **Purpose:** Reads the downloaded JSON data, decodes it into Go objects, and searches for comics matching provided keywords.
*   **Design Choices:**
    *   Requires the input JSON filename and at least one search term from the command line.
    *   Decodes the entire JSON array into a slice of `xkcd` structs in memory.
    *   Performs a quadratic search (nested loops): iterates through each comic and then through each search term.
    *   A comic is considered a match only if *all* provided search terms are found in its title or transcript.
    *   All comparisons are case-insensitive (strings are converted to lowercase before searching).
*   **Key Go Concepts:**
    *   **Struct Definition with JSON Tags:** Defines a Go struct that maps to the JSON structure, using `json:"fieldname"` tags for mapping.
        ```go
        type xkcd struct {
            Num       int    `json:"num"`
            Day       string `json:"day"`
            Month     string `json:"month"`
            Year      string `json:"year"`
            Title     string `json:"title"`
            Transcript string `json:"transcript"`
            // ... other fields if needed ...
        }
        ```
    *   **JSON Decoding:** Uses `encoding/json` package to decode the JSON file directly into a slice of structs.
        ```go
        input, err := os.Open(fn)
        if err != nil {
            // Handle file open error
        }
        defer input.Close()

        err = json.NewDecoder(input).Decode(&items) // Decode into a slice of xkcd structs
        if err != nil {
            fmt.Fprintln(os.Stderr, "bad json:", err)
            os.Exit(-1)
        }
        ```
    *   **Command Line Argument Processing:** Slices `os.Args` to extract the filename and search terms.
        ```go
        fn := os.Args[1] // Input filename
        // ... check for enough arguments ...
        for _, t := range os.Args[2:] { // Iterate over search terms
            terms = append(terms, strings.ToLower(t)) // Convert to lowercase
        }
        ```
    *   **String Manipulation:** Uses `strings.ToLower` for case-insensitive comparison and `strings.Contains` for substring matching.
        ```go
        lowerTitle := strings.ToLower(item.Title)
        lowerTranscript := strings.ToLower(item.Transcript)
        // ... inside inner loop ...
        if !strings.Contains(lowerTitle, term) && !strings.Contains(lowerTranscript, term) {
            continue outer // Skip to next comic if term not found
        }
        ```
    *   **Labeled `continue`:** Used to break out of the inner loop and continue with the next iteration of the outer loop when a search term doesn't match a comic.
        ```go
        outer:
        for _, item := range items {
            // ... convert title/transcript to lowercase ...
            for _, term := range terms {
                if !strings.Contains(lowerTitle, term) && !strings.Contains(lowerTranscript, term) {
                    continue outer // Skip this comic if any term doesn't match
                }
            }
            // ... if all terms matched, print comic details ...
        }
        ```

## What's New
- The behavior of variables declared by `for` loops changed in Go 1.22. Previously, these variables were created once and updated by each iteration. In Go 1.22, each iteration of the loop creates new variables to avoid accidental sharing bugs [7]. This change makes code that captures loop variables in closures or goroutines safer, as each iteration now has its own distinct variable. The provided `for` loop examples remain functionally correct, but their underlying variable semantics have been updated.

## Citations
- [7] Go version 1.22