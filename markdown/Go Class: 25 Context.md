# Go Class: 25 Context

## Summary
This video introduces Go's `context` package, a fundamental tool for managing the lifecycle of requests and related operations, especially in concurrent applications. It covers explicit and implicit cancellation, timeouts, and how contexts can carry request-specific values. The video emphasizes that contexts form an immutable tree structure, allowing cancellation and values to propagate downwards. It also highlights common pitfalls like goroutine leaks and best practices for using context keys.

## Key Points

*   **Introduction to Context:**
    *   The `context` package provides a common mechanism to manage work in progress, particularly for cancellation.
    *   It was introduced in Go 1.7 to standardize request lifecycle management.
    *   While not exclusively for concurrency, it integrates well with Go's concurrency primitives like `select`.

*   **Cancellation and Timeouts:**
    *   `context` supports two primary forms of cancellation:
        *   **Explicit Cancellation:** Triggered manually by calling a `cancel` function.
        *   **Implicit Cancellation:** Automatically triggered based on a `timeout` (duration) or `deadline` (specific time).
    *   A `context` can also carry request-specific values, such as a trace ID for distributed tracing.
    *   Many standard library functions (e.g., `net/http`, database drivers) accept a `context` parameter for managing their operations.
    *   A `context` provides two key controls:
        *   `ctx.Done()`: Returns a channel that is closed when the context is canceled or times out. This channel can be used in `select` statements to react to cancellation.
        *   `ctx.Err()`: Returns an error value indicating the reason for cancellation (e.g., `context.Canceled` or `context.DeadlineExceeded`).

*   **Context as an Immutable Tree Structure:**
    *   Contexts are organized as an immutable tree. When you add a timeout, deadline, or value to an existing context, you create a *new* child context that points to its parent. The parent context remains unchanged.
    *   This immutability ensures goroutine safety, as contexts can be safely passed across goroutine boundaries without fear of concurrent modification issues.
    *   Cancellation or timeout applied to a context propagates downwards to its entire subtree (all derived child contexts).
    *   Values stored in a context are also inherited by its children and can be retrieved by walking up the tree from any child node.

*   **Context Example (Cancellation and Leaks):**
    *   To create a context with a timeout:
        ```go
        ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
        defer cancel() // Always defer the cancel function to release resources
        ```
    *   To use the context with an HTTP request:
        ```go
        req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
        resp, err := http.DefaultClient.Do(req)
        ```
    *   If the HTTP request exceeds the context's timeout, `http.DefaultClient.Do()` will return an error (e.g., "context deadline exceeded").
    *   **Goroutine Leaks:** If goroutines are started that send results to an *unbuffered* channel, and the receiver stops listening (e.g., a `first` function returns after receiving one result), the sending goroutines will block indefinitely, leading to a goroutine leak.
    *   **Solution to Leaks:** Buffer the channel with a size equal to the number of expected senders (`make(chan result, len(urls))`). This allows all goroutines to send their results without blocking, even if not all results are consumed. The `defer cancel()` call ensures that any remaining active goroutines (e.g., those still performing network requests) are signaled to stop.

*   **Values in Context:**
    *   Context values should be used for data specific to a request's lifecycle, such as:
        *   Trace IDs (for distributed tracing).
        *   Start times (for latency calculations).
        *   Security or authorization data.
    *   **Avoid** using context to carry "optional" function parameters or structural components of your application (e.g., a logger instance). This can lead to hidden dependencies and make code harder to reason about.
    *   To prevent key collisions when storing values, use a **package-specific, private context key type** instead of a string:
        ```go
        type contextKey int
        const TraceKey contextKey = 1
        ```
    *   Storing a value: `ctx = context.WithValue(ctx, TraceKey, traceID)`
    *   Retrieving a value: `traceID, ok := ctx.Value(TraceKey).(string)`
        *   `ctx.Value()` returns an `interface{}`, so a type assertion is necessary. The `ok` boolean indicates if the value was found and successfully asserted to the given type.
    *   Context values are typically used in HTTP middleware to inject request-specific data that can then be accessed by downstream handlers in the request chain.

## What's New

*   The `context` package now provides `WithDeadlineCause` and `WithTimeoutCause` functions, which allow setting a specific cancellation cause that can be retrieved using the `Cause` method. This enhances the detail available when a context is canceled or times out. [6]
*   Go 1.22 changed the behavior of `for` loops to create new variables for each iteration. This addresses a common source of goroutine leaks where closures launched in loops capture the same loop variable across iterations. [7]
*   The `context` package introduced new functions: `WithoutCancel` (returns a copy of a context that is not canceled when the original context is canceled) and `AfterFunc` (registers a function to run after a context has been canceled). [6]

## Citations
- [1] Go 1.16 Release Notes
- [2] Go 1.17 Release Notes
- [3] Go 1.18 Release Notes
- [4] Go 1.19 Release Notes
- [5] Go 1.20 Release Notes
- [6] Go 1.21 Release Notes
- [7] Go 1.22 Release Notes
- [8] Go 1.23 Release Notes
- [9] Go 1.24 Release Notes