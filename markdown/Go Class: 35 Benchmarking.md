# Go Class: 35 Benchmarking

## Summary
This video introduces Go's benchmarking tools, which are integrated with its testing framework. It demonstrates how to write and run benchmarks to measure performance, memory allocation, and the impact of design decisions like recursion vs. iteration, linked lists vs. slices, and dynamic dispatch vs. direct calls. The video also highlights the critical performance issue of false sharing in concurrent programming and how to mitigate it, emphasizing the importance of understanding CPU cache behavior.

## Key Points

*   **Go Benchmarking Basics**
    *   Go provides built-in tools for benchmarking, similar to its unit testing framework.
    *   Benchmark functions reside in `_test.go` files.
    *   They are executed using the `go test -bench` command.
    *   Only functions prefixed with `Benchmark` are run as benchmarks.
    *   The `*testing.B` parameter passed to benchmark functions provides `b.N`, which is the number of iterations the benchmark runs to achieve a stable measurement (typically aiming for 1 second of execution).
    *   `ns/op` (nanoseconds per operation) is the primary metric for performance.

*   **Fibonacci Example: Iteration vs. Recursion**
    *   **Concept**: Compares the performance of a recursive Fibonacci implementation against an iterative one.
    *   **Code (Fibonacci function):**
        ```go
        func Fib(n int, recursive bool) int {
            switch n {
            case 0:
                return 0
            case 1:
                return 1
            default:
                if recursive {
                    return Fib(n-1, true) + Fib(n-2, true)
                }
                a, b := 0, 1
                for i := 1; i < n; i++ {
                    a, b = b, a+b
                }
                return b
            }
        }
        ```
    *   **Code (Benchmark functions):**
        ```go
        func BenchmarkFib20T(b *testing.B) {
            for n := 0; n < b.N; n++ {
                Fib(20, true) // Recursive
            }
        }

        func BenchmarkFib20F(b *testing.B) {
            for n := 0; n < b.N; n++ {
                Fib(20, false) // Iterative
            }
        }
        ```
    *   **Best Practice**: Iterative solutions are generally much faster than naive recursive solutions for problems like Fibonacci due to redundant calculations and function call overhead.

*   **List vs. Slice Example: Mechanical Sympathy**
    *   **Concept**: Compares the performance of summing elements in a linked list (scattered memory) versus a slice (contiguous memory). This demonstrates the impact of CPU cache efficiency (mechanical sympathy).
    *   **Code (Simplified `sumList` and `sumSlice`):**
        ```go
        type node struct {
            v int
            t *node
        }

        func sumList(h *node) (i int) { /* ... iterates through linked nodes ... */ }
        func sumSlice(l []int) (i int) { /* ... iterates through slice elements ... */ }
        ```
    *   **Benchmarking Tools**:
        *   `b.ResetTimer()`: Resets the benchmark timer, allowing measurement of only a specific part of the function (e.g., just the summing, not the list/slice creation).
        *   `go test -benchmem`: Adds memory allocation statistics (`allocs/op` and `B/op`) to the output.
    *   **Best Practice**: Slices (contiguous memory) generally outperform linked lists (scattered memory with pointers) for sequential access due to better CPU cache utilization. Linked lists also incur more memory allocations, which adds overhead.

*   **Forwarding Example: Dynamic Dispatch Overhead**
    *   **Concept**: Illustrates the performance cost of calling functions through interfaces (dynamic dispatch) and multiple layers of indirection compared to direct function calls.
    *   **Code (Interface and forwarding structs):**
        ```go
        type forwarder interface {
            forward(string) int
        }

        type thing1 struct { t forwarder }
        func (t1 *thing1) forward(s string) int { return t1.t.forward(s) }

        type thing2 struct { t forwarder }
        func (t2 *thing2) forward(s string) int { return t2.t.forward(s) }

        type thing3 struct{}
        func (t3 *thing3) forward(s string) int { return len(s) } // Actual work
        ```
    *   **Code (Benchmark functions):**
        ```go
        func BenchmarkDirect(b *testing.B) {
            // ... setup ...
            b.ResetTimer()
            for i := 0; i < b.N; i++ {
                _ = len(someString) // Direct call
            }
        }

        func BenchmarkForward(b *testing.B) {
            // ... setup of t1, t2, t3 chain ...
            b.ResetTimer()
            for i := 0; i < b.N; i++ {
                _ = t1.forward(someString) // Call through interface chain
            }
        }
        ```
    *   **Best Practice**: For performance-critical code, minimize unnecessary layers of abstraction and dynamic dispatch, especially for short, frequently called methods. The `//go:noinline` directive can be used to prevent compiler inlining for specific functions, allowing for more accurate measurement of call overhead.

*   **False Sharing Example: Cache Coherency Issues**
    *   **Concept**: Demonstrates how concurrent writes to different variables that happen to reside on the same CPU cache line can lead to performance degradation (false sharing) due to cache coherency protocols.
    *   **Code (Concurrent counting with shared counter slice):**
        ```go
        // Simplified, actual code involves goroutines and channels
        // False sharing occurs if 'cnt' elements are on the same cache line
        func count(cnt *uint64, in <-chan int) {
            var total uint64 // Local accumulator
            for i := range in {
                total += uint64(i)
            }
            *cnt = total // Write final sum to shared counter
            wg.Done()
        }
        ```
    *   **Benchmarking Tool**: `go test -bench=. -benchtime=10s -cpu=2,4,8` allows running benchmarks with varying numbers of CPU cores/threads.
    *   **Results**: With false sharing, increasing cores does not yield proportional speedup; performance may even degrade.
    *   **Solution**: Pad the shared data structure to ensure that variables accessed by different cores reside on different cache lines. This prevents cache line contention.
        *   Example: `cnt := make([]uint64, nworker * 8)` (padding each `uint64` with enough space to ensure it's on a separate cache line).
    *   **Best Practice**: Be mindful of data layout and access patterns in concurrent programs to avoid false sharing, which can severely limit scalability.

*   **General Benchmarking Considerations**
    *   **Environment**: Ensure the benchmarking environment is quiet and consistent (no other demanding processes).
    *   **Cache**: Consider if data/code is in cache.
    *   **Garbage Collection**: Be aware of GC pauses.
    *   **Virtual Memory**: Paging in/out can affect results.
    *   **Branch Prediction**: CPU branch prediction can influence performance.
    *   **Compiler Optimizations**: Be aware of compiler optimizations (e.g., inlining) that might alter the measured code.
    *   **Parallelism**: Understand physical vs. virtual cores and how they impact parallel execution.
    *   **Machine Sharing**: Avoid running benchmarks on machines shared with other processes.

## What's New
*   The `*testing.B` parameter still provides `b.N` for benchmark iterations, but Go 1.24 introduces the `testing.B.Loop` method as a faster and less error-prone way to perform benchmark iterations. This method handles the iteration loop internally, making the benchmark function execute exactly once per `-count` flag, which simplifies expensive setup and cleanup steps, and ensures function call parameters and results are kept alive to prevent the compiler from fully optimizing away the loop body [9].

## Updated Code Snippets
*   **Fibonacci Example: Benchmark functions**
    ```go
    func BenchmarkFib20T(b *testing.B) {
        b.ResetTimer()
        b.Loop(func() {
            Fib(20, true) // Recursive
        })
    }

    func BenchmarkFib20F(b *testing.B) {
        b.ResetTimer()
        b.Loop(func() {
            Fib(20, false) // Iterative
        })
    }
    ```
*   **Forwarding Example: Benchmark functions**
    ```go
    func BenchmarkDirect(b *testing.B) {
        // ... setup ...
        b.ResetTimer()
        b.Loop(func() {
            _ = len(someString) // Direct call
        })
    }

    func BenchmarkForward(b *testing.B) {
        // ... setup of t1, t2, t3 chain ...
        b.ResetTimer()
        b.Loop(func() {
            _ = t1.forward(someString) // Call through interface chain
        })
    }
    ```

## Citations
*   [9] Go version 1.24