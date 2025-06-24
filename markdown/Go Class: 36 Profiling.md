# Go Class: 36 Profiling

## Summary
This video introduces Go's built-in profiling tool, `pprof`, and demonstrates its use in identifying and resolving common performance issues and resource leaks in Go applications. It covers finding leaking goroutines and file descriptors, analyzing CPU usage, and optimizing code based on profiling results. The video also briefly touches on using Prometheus for collecting application metrics.

## Key Points

### Profiling Goroutine Leaks
*   **Identifying Leaks:** Goroutine leaks often manifest as a steadily increasing number of active goroutines, even when application traffic is stable or has ceased. This can lead to increased memory consumption and eventual application crashes.
*   **Using `net/http/pprof`:**
    *   Import `_ "net/http/pprof"` in your `main` package to automatically expose profiling endpoints.
    *   Access the `/debug/pprof/` endpoint in your browser to see available profiles (e.g., `goroutine`, `heap`, `profile`).
    *   Clicking on `goroutine` provides a stack trace of all current goroutines. A continuously increasing count after requests indicates a leak.
*   **Common Goroutine Leak Cause:** Not closing the `io.ReadCloser` (response body) of HTTP requests can lead to hung sockets and leaked goroutines.
    ```go
    // Example of a leak:
    func handler(w http.ResponseWriter, r *http.Request) {
        // ...
        resp, err := cli.Do(req)
        if err != nil {
            // handle error
            return
        }
        // Missing: defer resp.Body.Close()
        // ...
    }
    ```
*   **Using Prometheus for Metrics:**
    *   Prometheus can provide real-time metrics about your application's state, including the number of active goroutines and open file descriptors.
    *   Import Prometheus client libraries: `github.com/prometheus/client_golang/prometheus` and `github.com/prometheus/client_golang/prometheus/promhttp`.
    *   Register custom counters (e.g., `queries` for request count) and expose them via an HTTP endpoint (e.g., `/metrics`).
    ```go
    var queries = prometheus.NewCounter(prometheus.CounterOpts{
        Name: "all_queries",
        Help: "How many queries we've received.",
    })

    func init() {
        prometheus.MustRegister(queries)
    }

    func main() {
        // ...
        http.Handle("/metrics", promhttp.Handler())
        // ...
    }

    func handler(w http.ResponseWriter, r *http.Request) {
        // ...
        queries.Inc() // Increment counter on each request
        // ...
    }
    ```
    *   Monitoring `go_goroutines` and `process_open_fds` metrics from Prometheus can help detect leaks. If these numbers don't decrease after traffic subsides, a leak is likely present.

### CPU Profiling
*   **Generating a CPU Profile:**
    *   Build your Go program as a binary using `go build .` (or `go build ./cmd/your_app`).
    *   Run the compiled binary.
    *   Access the `/debug/pprof/profile` endpoint in your browser (e.g., `http://localhost:8080/debug/pprof/profile`). This will download a profile file after a default duration (e.g., 30 seconds).
*   **Analyzing the CPU Profile:**
    *   Use the `go tool pprof` command with the binary and the downloaded profile file: `go tool pprof <binary> <profile-file>`.
    *   **Interactive Mode:** Running `go tool pprof` without options provides an interactive prompt.
    *   **Top Entries (`-top`):** `go tool pprof -top <binary> <profile-file>` shows the functions consuming the most CPU time.
    *   **Browser UI (`-http`):** `go tool pprof -http=":6060" <binary> <profile-file>` opens a web-based UI in your browser, offering powerful visualizations:
        *   **Graph View:** Shows a call graph with nodes representing functions and edges representing calls. Thicker edges and larger nodes indicate more time spent.
        *   **Flame Graph:** A stacked bar chart where the width of a bar represents the percentage of time spent in a function and its children.
        *   **Source View:** Displays the source code of functions, highlighting lines where CPU time is concentrated.
*   **Optimization Example (Sort Animation):**
    *   **Problem:** An animation program drawing squares for a sorting algorithm was found to spend significant CPU time in the `paintSquareSlow` function.
    *   **Initial `paintSquareSlow`:** This function used nested loops to draw each pixel of an 8x8 square, calling `img.SetColorIndex(x, y, color)` for every pixel.
    *   **Analysis:** Profiling revealed that `img.SetColorIndex` performed redundant work (bounds checks, pixel offset calculations) for each pixel, even though the square's pixels are contiguous in memory.
    *   **Optimization Strategy:**
        1.  **Eliminate unnecessary checks:** The outer function already ensures the square is within image bounds, so `SetColorIndex`'s internal bounds check is redundant.
        2.  **Move multiplication out of loops:** Calculate the starting pixel offset for the square once, outside the inner loops.
        3.  **Strength reduction:** Replace repeated multiplications within loops with additions, leveraging the contiguous memory layout.
        4.  **Reorder loops (y then x):** Process pixels row by row, then column by column, to take advantage of cache locality and slice copying.
        5.  **Use `copy` for entire rows:** Instead of setting each pixel individually, create a pre-filled row slice and use `copy(destination_slice, source_slice)` to paint entire rows at once. This is much faster as `copy` is highly optimized.
    *   **Result:** These optimizations significantly reduced the CPU time spent in the drawing function (e.g., from ~20-27% to ~3% of total CPU time), leading to a substantial speedup (e.g., 6-7x faster for the drawing part).

## What's New
*   The Go runtime's internal metrics, including those related to goroutines, are now exposed through a stable and more efficient `runtime/metrics` package, superseding older functions like `runtime.ReadMemStats` and `debug.GCStats`. [1]

## Citations
*   [1] Go version 1.16