# Go Class: 34 Mechanical Sympathy

## Summary
This video introduces the concept of "Mechanical Sympathy" in programming, drawing parallels from race car driving where understanding the machine leads to better performance. It highlights the historical trend of software becoming less efficient despite hardware advancements, emphasizing the importance of writing code that works *with* the underlying machine architecture. The discussion delves into various aspects of computer performance, including CPU and memory hierarchy, caching, locality, and synchronization costs, providing practical insights and best practices for optimizing Go programs by understanding these mechanical realities.

## Key Points

*   **Introduction to Mechanical Sympathy:**
    *   Understanding how the underlying machine works is crucial for writing faster and more efficient programs.
    *   The term originates from race car driving, where mechanics and drivers who deeply understand the car can extract maximum performance.

*   **Hardware vs. Software Performance:**
    *   Despite significant hardware gains (CPU speed, memory, disk space), perceived software performance has not improved proportionally over decades.
    *   This is attributed to software becoming increasingly inefficient, "cancelling out" hardware advancements.

*   **Performance in the Cloud & Trade-offs:**
    *   Building applications in the cloud inherently involves accepting some performance overhead due to distributed architectures and shared resources.
    *   Performance optimization must be balanced against other critical factors like architecture choice, quality, reliability, scalability, and development/ownership costs.
    *   Simplicity, readability, and maintainability of code remain paramount.

*   **Optimization - Top-down Refinement:**
    *   Optimization efforts should follow a top-down approach for maximum impact:
        *   **Architecture:** Focus on reducing latency and communication costs (e.g., microservices vs. monolith).
        *   **Design:** Optimize algorithms, concurrency models, and abstraction layers.
        *   **Implementation:** Consider programming language choice and memory usage patterns.
    *   Mechanical sympathy primarily applies at the implementation level.
    *   Interpreted languages often incur significant performance penalties (e.g., 10x slower) due to their abstraction layers, making mechanical sympathy difficult.

*   **CPU Performance Trends:**
    *   Around 2005, CPU clock speeds (frequency) hit a physical limit (the "frequency wall").
    *   Hardware performance gains shifted from single-core speed to increasing the number of logical cores.
    *   Memory (DRAM) access speeds have not kept pace with CPU speeds, leading to a widening performance gap.

*   **Unfortunate Realities:**
    *   CPUs are no longer getting faster in terms of clock speed.
    *   The performance gap between CPU and memory is not shrinking.
    *   Software tends to get slower more quickly than CPUs get faster.
    *   Software development costs significantly exceed hardware costs.

*   **Software Bloatation:**
    *   Software often grows in size and complexity ("bloat"), consuming more CPU capacity than necessary.
    *   This bloat can be due to added features or inefficient coding practices, including the use of high-level abstractions that hide underlying machine costs.

*   **The Solution: Simpler Software, Machine Sympathy:**
    *   To address performance and cost challenges, software needs to be simpler and designed to work *with* the machine, not against it.
    *   This means being aware of and leveraging the machine's characteristics.

*   **Memory Hierarchy:**
    *   Modern CPUs utilize a multi-level memory hierarchy (L1, L2, L3 caches, DRAM, SSD, Cloud storage).
    *   Access latency increases dramatically as you move further down the hierarchy (e.g., L1 cache is 1 cycle, DRAM is 400 cycles, Cloud is 240,000,000 cycles).
    *   Efficient memory access is critical for performance.

*   **Memory Caching:**
    *   Caches are designed to reduce memory access time by storing frequently used data closer to the CPU.
    *   Memory is accessed in fixed-size **cache lines** (typically 64 bytes).
    *   **Cache coherency** mechanisms ensure data consistency across multiple CPU cores, but incur overhead.

*   **Locality:**
    *   Caches exploit two types of locality:
        *   **Locality in space:** If one memory location is accessed, nearby locations are likely to be accessed soon. (Encourages contiguous data structures).
        *   **Locality in time:** If a memory location is accessed, it's likely to be accessed again soon. (Encourages keeping frequently used data in cache).
    *   Optimal cache performance is achieved when data is stored contiguously and accessed sequentially.

*   **Cache Efficiency (Factors):**
    *   **Things that make cache less efficient:**
        *   **Synchronization between CPUs:** Locks and mutexes introduce overhead and contention.
        *   **Copying blocks of data around in memory:** Incurs unnecessary memory writes and reads.
        *   **Non-sequential access patterns:** Calling functions (jumps in code), chasing pointers (scattered data) break locality.
    *   **Things that make cache more efficient:**
        *   Keeping code or data in cache longer (high temporal locality).
        *   Keeping data together (so an entire cache line is used, high spatial locality).
        *   Processing memory in sequential order (code or data).
    *   "A little copying is better than a lot of pointer chasing!"

*   **Access Patterns: Slice vs. Linked List:**
    *   **Slices (or arrays) of objects** are generally more cache-efficient than **linked lists with pointers**.
    *   Slices store data contiguously, allowing the CPU to prefetch entire cache lines.
    *   Linked lists scatter data in memory, leading to frequent cache misses and costly pointer chasing.

    ```go
    // Example: Slice of structs (contiguous)
    type MyStruct struct {
        Field1 int
        Field2 string
    }
    mySlice := make([]MyStruct, 100)

    // Example: Linked list (scattered, pointer chasing)
    type Node struct {
        Value MyStruct
        Next  *Node
    }
    head := &Node{}
    // ... nodes are allocated individually, potentially far apart
    ```

*   **Access Patterns: Short Methods via Dynamic Dispatch:**
    *   Calling many short methods, especially through dynamic dispatch (e.g., virtual functions in C++, interface methods in Go), can be very expensive.
    *   Each call might involve multiple pointer dereferences (e.g., v-table lookup) before any actual work is done.
    *   The cost of calling a function should be proportional to the work it performs.
    *   "Forwarding methods" (methods that simply call another method on an internal object) are a design smell that adds unnecessary overhead.

*   **Synchronization Costs: False Sharing:**
    *   **False sharing** occurs when independent variables, accessed by different CPU cores, happen to reside within the same cache line.
    *   Even though there's no logical data race, writes to one variable by one core will invalidate the entire cache line for other cores, forcing them to re-fetch it.
    *   This "bouncing" of cache lines between cores can drastically reduce performance.
    *   Solution often involves padding structs to ensure frequently accessed independent variables reside on different cache lines.

*   **Other Hidden Costs:**
    *   Disk access, garbage collection (GC), virtual memory, and context switching between processes are other sources of overhead.
    *   While many are OS-managed, GC can be influenced by programmer choices.
    *   To optimize GC:
        *   Reduce unnecessary memory allocations.
        *   Reduce embedded pointers in objects (simplifies GC traversal).
        *   Paradoxically, a larger heap can sometimes lead to better performance (less frequent GC cycles).

*   **Optimization in Go Philosophy:**
    *   Go's design philosophy encourages "good design" by providing choices that allow programmers to be aware of and manage underlying machine costs.
    *   Go doesn't hide the costs involved in operations, making it easier to reason about performance.
    *   This includes choices like:
        *   Allocating contiguously (slices/arrays).
        *   Deciding whether to copy data or pass pointers.
        *   Choosing between stack and heap allocation.
        *   Being explicit about synchronous vs. asynchronous operations.
        *   Avoiding unnecessary abstraction layers.
        *   Avoiding overly short/forwarding methods.

*   **Optimization Principles (Knuth & Fromberger):**
    *   **Don Knuth:** "Premature optimization is the root of all evil." Focus on small efficiencies only in critical 3% of the code, as most optimization attempts on non-critical parts waste time and negatively impact maintainability.
    *   **Michael Fromberger:** "There are only three optimizations: 1. Do less. 2. Do it less often. 3. Do it faster." The biggest performance gains come from "doing less" (simplifying the problem/algorithm), followed by "doing it less often," with "doing it faster" (micro-optimizations) being the last resort.

## What's New

*   **Compiler and Runtime Optimizations:** Go's compiler and runtime have received significant performance improvements that automatically apply principles of mechanical sympathy, reducing the need for manual micro-optimizations in some cases.
    *   **Register-based calling conventions:** Go 1.17 introduced a new way of passing function arguments and results using registers instead of the stack, leading to performance improvements and reduced binary size [2]. This was expanded to more architectures in Go 1.18 [3].
    *   **Profile-Guided Optimization (PGO):** Go 1.20 introduced preview support for PGO, allowing the compiler to perform application- and workload-specific optimizations based on runtime profiles [5]. PGO became ready for general use in Go 1.21, further improving performance by devirtualizing interface method calls and enabling more aggressive inlining [6]. Build time overhead for PGO was significantly reduced in Go 1.23 [8].
    *   **General Runtime Performance:** The runtime has seen continuous improvements, including more prompt memory release to the OS (Go 1.17 [1]), a soft memory limit for better GC control (Go 1.19 [4]), and internal GC tuning for reduced tail latency and memory use (Go 1.21 [6], Go 1.22 [7]). Go 1.24 introduced a new builtin map implementation and more efficient small object allocation, further reducing CPU overheads [9].
*   **Loop Variable Semantics:** In Go 1.22, the behavior of `for` loops changed: variables declared by a `for` loop are now created anew for each iteration, rather than being updated once per loop. This change addresses a common class of "accidental sharing bugs" when variables are captured by closures within loops [7].
*   **Enhanced Slice Operations:** The standard library has introduced new functions and types to work with slices more efficiently and powerfully. This includes `unsafe.Slice` for low-level slice creation (Go 1.17 [2]), and the new `slices` package with generic functions for common slice operations (Go 1.21 [6]).

## Updated Code Snippets
The original code snippet for "Access Patterns: Slice vs. Linked List" is still valid for demonstrating the concept. However, to illustrate the significant change in loop variable semantics introduced in Go 1.22, an additional example is provided:

```go
// Old behavior (Go 1.21 and earlier):
// In this example, if goroutines were launched, they would likely
// all print the final value of 'i' (e.g., 9) because 'i' was a single
// variable shared across all loop iterations.
// for i := 0; i < 10; i++ {
//     go func() {
//         fmt.Println(i)
//     }()
// }

// New behavior (Go 1.22 and later):
// Each iteration of the loop now creates a new 'i' variable.
// As a result, each goroutine captures its own unique 'i' value (0 through 9),
// preventing accidental sharing bugs.
package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(i) // In Go 1.22+, this 'i' is unique to each iteration
		}()
	}
	wg.Wait()
}
```

## Citations
- [1] Go version 1.16
- [2] Go version 1.17
- [3] Go version 1.18
- [4] Go version 1.19
- [5] Go version 1.20
- [6] Go version 1.21
- [7] Go version 1.22
- [8] Go version 1.23
- [9] Go version 1.24