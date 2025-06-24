# Go Class: 27 Concurrent File Processing

## Summary
This video explores different approaches to parallelizing a file processing task in Go: finding duplicate files based on their content using secure hashes. Starting from a sequential implementation, the presenter demonstrates two concurrent models using goroutines and channels: a fixed worker pool (map-reduce style) and a parallel tree walk with goroutine limiting. The video highlights the importance of managing shared resources and understanding Amdahl's Law for optimizing concurrent programs.

## Key Points

*   **Problem Statement: Finding Duplicate Files**
    *   The goal is to find duplicate files in a large directory structure (e.g., Dropbox folder).
    *   Duplication is determined by file content, not just name or date, requiring a secure hash (like MD5 for this example, as it's not for internet security).
    *   The sequential approach takes approximately 5 minutes on a modern quad-core machine.

*   **Sequential Approach**
    *   **Declarations:** Defines structs to hold file hash and path, and a map to store results (hash to a slice of file paths).
        ```go
        type pair struct {
            hash, path string
        }
        type fileList []string
        type results map[string]fileList
        ```
    *   **Hashing:** The `hashFile` function opens a file, defers its closing (crucial for many files), creates an MD5 hasher, copies file content to the hasher, and returns the hexadecimal hash string along with the file path.
        ```go
        func hashFile(path string) pair {
            file, err := os.Open(path)
            if err != nil {
                log.Fatal(err) // Handle error appropriately in real code
            }
            defer file.Close() // Ensure file is closed

            hash := md5.New()
            if _, err := io.Copy(hash, file); err != nil {
                log.Fatal(err) // Handle error appropriately
            }
            return pair{fmt.Sprintf("%x", hash.Sum(nil)), path}
        }
        ```
    *   **Searching:** The `searchTree` function uses `filepath.Walk` (a built-in Go library function implementing the visitor pattern) to traverse the directory tree. It filters for regular, non-empty files, hashes them, and stores the hash-path pairs in a map.
        ```go
        func searchTree(dir string) (results, error) {
            hashes := make(results)
            err := filepath.Walk(dir, func(p string, fi os.FileInfo, err error) error {
                if err != nil { // Always check error from filepath.Walk
                    return err
                }
                if fi.Mode().IsRegular() && fi.Size() > 0 {
                    h := hashFile(p)
                    hashes[h.hash] = append(hashes[h.hash], h.path)
                }
                return nil
            })
            return hashes, err
        }
        ```

*   **Concurrent Approach #1: Fixed Worker Pool (Map-Reduce Style)**
    *   **Architecture:** A "Tree walk" goroutine feeds file paths into a `paths` channel. A fixed pool of "Workers" read from `paths`, hash files, and send hash-path `pairs` to a `pairs` channel. A "Collector" goroutine reads from `pairs`, aggregates results into a map, and sends the final `results` to a `results` channel. `done` channels are used for synchronization.
    *   **Worker Function (`processFiles`):** Reads paths from an input channel, hashes the file, and sends the `pair` to an output channel. Signals completion on a `done` channel.
        ```go
        func processFiles(paths <-chan string, pairs chan<- pair, done chan<- bool) {
            for path := range paths { // Loop until paths channel is closed
                pairs <- hashFile(path)
            }
            done <- true // Signal this worker is done
        }
        ```
    *   **Collector Function (`collectHashes`):** Reads `pair` values from its input channel and aggregates them into the `results` map. Once the input channel is closed and drained, it sends the final `results` map to its output channel.
        ```go
        func collectHashes(pairs <-chan pair, results chan<- results) {
            hashes := make(results)
            for p := range pairs { // Loop until pairs channel is closed
                hashes[p.hash] = append(hashes[p.hash], p.path)
            }
            results <- hashes // Send final aggregated results
        }
        ```
    *   **Main Logic:**
        1.  Determine the number of workers (e.g., `2 * runtime.GOMAXPROCS(0)` for I/O bound tasks).
        2.  Create unbuffered channels: `paths` (string), `pairs` (pair), `done` (bool), `results` (results map).
        3.  Start the `collectHashes` goroutine.
        4.  Start `N` `processFiles` worker goroutines.
        5.  The `filepath.Walk` function (now in `main`) sends file paths to the `paths` channel.
        6.  After `filepath.Walk` completes, `close(paths)` to signal workers no more paths are coming.
        7.  Wait for all workers to signal `true` on the `done` channel.
        8.  `close(pairs)` to signal the collector no more pairs are coming.
        9.  Receive the final `results` from the `results` channel.
    *   **Performance (Evaluation #1):**
        *   Sequential: ~5 minutes.
        *   Concurrent (unbuffered channels): 56.11s.
        *   Adding buffers to `pairs` channel: 52.76s (reduces blocking, keeps workers busy).
        *   Increasing workers to 32 (from 16): 51.36s (further improvement due to I/O bound nature).

*   **Concurrent Approach #2: Parallel Tree Walk**
    *   **Concept:** Instead of a single sequential `filepath.Walk`, each subdirectory encountered starts a new goroutine to walk that subdirectory. This parallelizes the directory traversal itself.
    *   **Implementation:**
        *   The `searchTree` function is replaced by a recursive `walkDir` function that takes a `*sync.WaitGroup` as an argument.
        *   When `walkDir` starts, it calls `defer wg.Done()`.
        *   If `walkDir` encounters a subdirectory, it calls `wg.Add(1)` and starts a new `go walkDir(...)` for that subdirectory, then returns `filepath.SkipDir` to prevent the current `filepath.Walk` from descending.
        *   File processing remains the same (sending paths to the worker pool).
    *   **Performance (Evaluation #2):**
        *   Basic version (parallel tree walk, unbuffered channels): 51.14s.
        *   Adding buffers to all channels: 50.03s.
        *   Increasing workers to 32: 48.75s.
        *   This approach offers a slight improvement because the path identification is also parallelized, reducing idle time for workers waiting for paths.

*   **Concurrent Approach #3: Goroutines Galore! (Limiting Active Goroutines)**
    *   **Problem:** Starting a goroutine for *every* file and directory can lead to resource exhaustion (running out of OS threads) if not managed, especially for I/O-bound tasks. `GOMAXPROCS` only limits CPU-bound goroutines, not those blocked on syscalls.
    *   **Solution: Channels as Counting Semaphores:**
        *   A buffered channel (`limits chan bool`) is used to control the number of *active* goroutines (those performing I/O).
        *   The buffer size `N` of this channel determines the maximum number of concurrent I/O operations.
        *   Before performing I/O, a goroutine attempts to `send` a value to the `limits` channel. If the channel is full, it blocks, effectively limiting active I/O.
        *   After completing I/O, the goroutine `receives` a value from the `limits` channel, freeing up a slot.
    *   **Implementation Details:**
        *   The `processFile` function (for individual files) now takes the `limits` channel.
        *   It calls `limits <- true` at the beginning and `defer func() { <-limits }()` at the end.
        *   The `walkDir` function (for directories) also takes the `limits` channel and passes it down.
        *   The `filepath.Walk` visitor function within `walkDir` now starts `go processFile(...)` for files and `go walkDir(...)` for subdirectories, passing the `limits` channel.
    *   **Performance (Evaluation #3):**
        *   Best time achieved: 46.93s using 32 workers.
        *   Increasing the `limits` buffer beyond the optimal point makes the time grow longer due to increased disk contention.
    *   **Amdahl's Law:** The speedup of a program is limited by the portion that cannot be parallelized. Even with infinite processors, a program that is 95% parallel can only achieve a maximum speedup of 20x. In this example, with 8 logical processors, a speedup of ~6.25x was achieved, indicating approximately 96% of the program was parallelized.

*   **Conclusions**
    *   It's not necessary to strictly limit the total number of goroutines created in a Go program.
    *   It *is* crucial to limit contention for shared resources, especially I/O-bound resources like disk or network access.
    *   Go's channels can effectively act as counting semaphores to manage active goroutines and prevent resource contention.

## What's New
*   The `filepath.Walk` function, while still valid, has a more efficient alternative in `filepath.WalkDir` (introduced in Go 1.16), which uses `fs.DirEntry` instead of `os.FileInfo` for directory entries. `os.FileInfo` itself became an alias for `fs.FileInfo` in Go 1.16, ensuring backward compatibility. [1]

## Updated Code Snippets
```go
import (
	"io"
	"io/fs" // New import for fs.DirEntry
	"log"
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	// ... other imports
)

// Assuming 'pair' and 'results' structs are defined elsewhere as in the original
// type pair struct {
//     hash, path string
// }
// type fileList []string
// type results map[string]fileList

func searchTree(dir string) (results, error) {
    hashes := make(results)
    // Use filepath.WalkDir for improved efficiency and modern API usage
    err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
        if err != nil { // Always check error from filepath.WalkDir
            return err
        }
        // Check if it's a regular file and not empty
        if d.Type().IsRegular() {
            info, err := d.Info() // Get FileInfo for size check
            if err != nil {
                // If we can't get file info, it might be a permission error or similar.
                // Depending on requirements, you might log and skip, or return the error.
                // For this example, we'll return the error as in the original Walk behavior.
                return err
            }
            if info.Size() > 0 {
                h := hashFile(p)
                hashes[h.hash] = append(hashes[h.hash], h.path)
            }
        }
        return nil
    })
    return hashes, err
}
```

## Citations
- [1] Go version 1.16