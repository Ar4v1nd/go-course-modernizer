# Go Class: 23 CSP, Goroutines, and Channels

## Summary
This video introduces concurrency in Go, focusing on the Communicating Sequential Processes (CSP) model. It explains core Go features like channels and goroutines, demonstrating their use through practical examples such as parallel HTTP requests, a safe web server counter, and a dynamic prime sieve. The video emphasizes how Go's concurrency primitives simplify complex parallel programming by promoting communication over shared memory.

## Key Points

*   **Concurrency in Go (CSP)**
    *   Go's concurrency model is based on Communicating Sequential Processes (CSP), an idea from the 1970s.
    *   CSP simplifies concurrent programming by breaking down problems into independent, communicating parts.
    *   Go provides built-in features: **channels** and **goroutines**.

*   **Channels**
    *   A channel is a one-way communication pipe, similar to a Unix pipe.
    *   Data sent into one end comes out the other in the same order.
    *   Channels remain open until explicitly closed by the sender.
    *   Crucially, Go channels are safe for multiple goroutines to read from and write to concurrently, preventing data races.
    *   Channels serve as both communication and synchronization mechanisms.
    *   The "happens-before" relationship: A send operation on a channel always happens before the corresponding receive operation.
    *   Channels facilitate transferring ownership of data between goroutines, avoiding shared memory issues.

*   **Sequential Processes**
    *   A sequential process is an independent part of a program that executes instructions in a defined sequence (e.g., read, process, write in a loop).
    *   In Go, multiple sequential processes can run concurrently (or in parallel on multi-core hardware) by communicating through channels.

*   **Goroutines**
    *   A goroutine is a lightweight unit of independent execution, often referred to as a coroutine.
    *   Starting a goroutine is simple: prefix a function call with the `go` keyword.
        ```go
        go myFunction()
        ```
    *   Go's runtime efficiently schedules many goroutines onto a smaller number of operating system threads. A Go program can manage tens of thousands of goroutines.
    *   It's essential to ensure goroutines terminate correctly to avoid "goroutine leaks" (orphaned goroutines consuming resources). Termination can be achieved via:
        *   Well-defined loop terminating conditions.
        *   Signaling completion through channels or contexts.
        *   Allowing them to run until the main program exits.

*   **Concurrency Example 1: Parallel HTTP GET**
    *   Demonstrates fetching multiple URLs concurrently to improve performance.
    *   A `get` function is defined to fetch a single URL and send its result (URL, error, latency) back on a channel.
    *   The `main` function launches a separate goroutine for each URL using `go get(...)`.
    *   It then reads results from a shared channel, ensuring all concurrent operations complete and their results are collected.
    *   This approach significantly reduces total execution time compared to sequential fetching.

*   **Concurrency Example 2: Web Server with Counter (Data Race & Solution)**
    *   Illustrates a common concurrency problem: a data race when multiple concurrent requests try to increment a shared counter variable directly.
    *   The unsafe operation (`nextID++`) is a read-modify-write operation, prone to race conditions in a concurrent environment.
    *   Solution: Use a channel to safely manage the counter.
        *   A `counter` goroutine continuously generates incrementing numbers and sends them to a channel.
        *   The HTTP `handler` function reads the next available number from this channel.
        *   This ensures that only one goroutine "owns" and processes the number at any given time, eliminating the race condition.
    *   Channels act as a synchronization point, blocking writers until a reader is ready, and vice-versa.

*   **Concurrency Example 3: Prime Sieve**
    *   A classic example of dynamic concurrency using a pipeline of goroutines and channels.
    *   A `generator` goroutine produces a sequence of numbers.
    *   A `sieve` function iteratively creates `filter` goroutines for each new prime number found.
    *   Each `filter` goroutine receives numbers from an input channel, filters out multiples of its prime, and sends the remaining numbers to an output channel.
    *   When an input channel closes, the `filter` goroutine closes its output channel, creating a "domino effect" that propagates closure through the pipeline.
    *   This example showcases the elegance of Go's concurrency model for building complex, dynamic pipelines, even if this specific implementation has communication overhead for very large numbers.

## What's New

*   **Goroutines**: The statement "It's essential to ensure goroutines terminate correctly to avoid 'goroutine leaks' (orphaned goroutines consuming resources). Termination can be achieved via: Well-defined loop terminating conditions." remains valid. However, a significant language change in Go 1.22 made `for` loop variables per-iteration by default, which helps prevent a common class of accidental sharing bugs (and thus goroutine leaks) when closures capture loop variables. The `vet` tool also no longer reports these specific loop variable capture issues for code targeting Go 1.22 or newer, as the underlying language behavior has made them safe by default. [7]
*   **Channels**: While channels generally maintain their properties, a specific change in Go 1.23 made `time.Timer` and `time.Ticker` channels unbuffered. This means that `len` and `cap` on these specific channels will now return 0, which could affect programs that previously relied on them returning 1. [8]
*   **Concurrency Example 2: Web Server with Counter (Data Race & Solution)**: The channel-based solution for safely managing a counter is still correct. However, for simple integer counters, the `sync/atomic` package introduced new types like `atomic.Int64` in Go 1.19. These provide a more idiomatic and often more efficient way to handle atomic operations without explicit channel communication. [4]

## Citations
*   [1] Go version 1.16
*   [2] Go version 1.17
*   [3] Go version 1.18
*   [4] Go version 1.19
*   [5] Go version 1.20
*   [6] Go version 1.21
*   [7] Go version 1.22
*   [8] Go version 1.23
*   [9] Go version 1.24