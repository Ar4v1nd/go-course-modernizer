# Go Class: 30 Concurrency Gotchas

## Summary
This video explores common pitfalls and best practices in Go concurrency. It highlights generic concurrency problems like race conditions and deadlocks, along with Go-specific challenges related to goroutines, channels, mutexes, WaitGroups, and the `select` statement. The core message emphasizes understanding the underlying mechanisms and adopting disciplined coding practices to write robust concurrent applications.

## Key Points

*   **Concurrency Problems Overview**
    *   **Race Conditions:** Occur when multiple goroutines access shared mutable data without proper synchronization, leading to unpredictable results. This typically involves unprotected read and write operations that overlap.
    *   **Deadlock:** A state where two or more goroutines are blocked indefinitely, each waiting for a resource that another blocked goroutine holds. Go's runtime can detect some deadlocks automatically.
    *   **Goroutine Leak:** A goroutine that never terminates, often due to being blocked indefinitely on a channel that will never receive or send a value. This consumes system resources over time.
    *   **Channel Errors:** Common issues include attempting to send on a closed channel, sending or receiving on a `nil` channel, or closing a channel multiple times.

*   **Gotcha 1: Data Race**
    *   A data race occurs when multiple goroutines access the same memory location concurrently, and at least one of the accesses is a write operation, without any synchronization mechanism.
    *   **Example:**
        ```go
        package main

        import (
            "fmt"
            "log"
            "net/http"
        )

        var nextID = 0 // Shared mutable state

        func handler(w http.ResponseWriter, r *http.Request) {
            fmt.Fprintf(w, "<h1>You got %v<h1/>", nextID)
            // unsafe - data race
            nextID++ // Unprotected write, concurrent access by multiple handlers
        }

        func main() {
            http.HandleFunc("/", handler)
            if err := http.ListenAndServe(":8080", nil); err != nil {
                log.Fatal(err)
            }
        }
        ```
    *   Web server handlers run concurrently, making `nextID++` a race condition.

*   **Gotcha 2: Deadlock**
    *   **Channel Deadlock:** Occurs when a goroutine waits to receive from a channel, but no other goroutine will ever send to it. Go's runtime detects this.
        ```go
        package main

        import (
            "fmt"
        )

        func main() {
            ch := make(chan bool) // Unbuffered channel
            fmt.Println("START")
            go func(ok bool) {
                if ok { // 'ok' is false, so this branch is skipped
                    ch <- ok
                }
            }(false)
            <-ch // Main goroutine waits forever
            fmt.Println("DONE")
        }
        ```
    *   **Mutex Deadlock (Unreleased Lock):** A goroutine acquires a mutex but fails to release it, preventing other goroutines from acquiring it.
        ```go
        package main

        import (
            "fmt"
            "sync"
            "time"
        )

        func main() {
            var m sync.Mutex
            done := make(chan bool)
            fmt.Println("START")
            go func() {
                m.Lock() // Lock acquired, but no corresponding Unlock
                // Missing defer m.Unlock()
            }()
            go func() {
                time.Sleep(1 * time.Second)
                m.Lock() // This goroutine will block indefinitely
                defer m.Unlock()
                fmt.Println("SIGNAL")
                done <- true
            }()
            <-done
            fmt.Println("DONE")
        }
        ```
        *   **Best Practice:** Always use `defer m.Unlock()` immediately after `m.Lock()` to guarantee release.
    *   **Dining Philosophers Problem (Circular Wait):** A classic deadlock scenario where multiple goroutines (philosophers) require multiple resources (forks/mutexes) and acquire them in conflicting orders, leading to a circular dependency.
        ```go
        package main

        import (
            "fmt"
            "sync"
            "time"
        )

        func main() {
            var m1, m2 sync.Mutex
            done := make(chan bool)
            fmt.Println("START")
            go func() { // Goroutine 1: Locks m1 then m2
                m1.Lock()
                defer m1.Unlock()
                time.Sleep(1 * time.Second)
                m2.Lock()
                defer m2.Unlock()
                fmt.Println("SIGNAL 1")
                done <- true
            }()
            go func() { // Goroutine 2: Locks m2 then m1 (causes deadlock)
                m2.Lock()
                defer m2.Unlock()
                time.Sleep(1 * time.Second)
                m1.Lock()
                defer m1.Unlock()
                fmt.Println("SIGNAL 2")
                done <- true
            }()
            <-done
            fmt.Println("DONE")
        }
        ```
        *   **Best Practice:** To prevent this, ensure all goroutines acquire multiple mutexes in a consistent, predefined order.

*   **Gotcha 3: Goroutine Leak**
    *   A goroutine can leak if it attempts a blocking send on an unbuffered channel, but the receiver times out or stops listening, leaving the sender permanently blocked.
    *   **Problematic Code:**
        ```go
        package main

        import (
            "fmt"
            "time"
        )

        type obj struct{}

        func fn() obj {
            time.Sleep(2 * time.Second) // Simulate long work
            return obj{}
        }

        func finishReq(timeout time.Duration) *obj {
            ch := make(chan obj) // Unbuffered channel
            go func() {
                ch <- fn() // Blocking send
            }()
            select {
            case rslt := <-ch:
                return &rslt
            case <-time.After(timeout): // Receiver times out, sender goroutine leaks
                return nil
            }
        }

        func main() {
            fmt.Println("Request started")
            result := finishReq(1 * time.Second) // Timeout is shorter than work
            if result == nil {
                fmt.Println("Request timed out")
            } else {
                fmt.Println("Request completed")
            }
            time.Sleep(3 * time.Second) // Allow time to observe leak
            fmt.Println("Main exiting")
        }
        ```
    *   **Best Practice:** Use a buffered channel (`ch := make(chan obj, 1)`) if the sender needs to complete its operation regardless of immediate receiver readiness. This allows the sender to place the value in the buffer and exit.

*   **Gotcha 4: Incorrect use of WaitGroup**
    *   `wg.Add(1)` must be called *before* the `go` statement that launches the goroutine. If called too late, the main goroutine might call `wg.Wait()` before the launched goroutine has incremented the counter, leading to premature return.
    *   **Problematic Code:**
        ```go
        package main

        import (
            "fmt"
            "sync"
            "time"
        )

        func worker(id int, wg *sync.WaitGroup) {
            defer wg.Done()
            fmt.Printf("Worker %d starting\n", id)
            time.Sleep(time.Duration(id) * 100 * time.Millisecond)
            fmt.Printf("Worker %d done\n", id)
        }

        func main() {
            var wg sync.WaitGroup
            for i := 0; i < 3; i++ {
                go worker(i, &wg)
                wg.Add(1) // BIG MISTAKE: Add is called AFTER go
            }
            wg.Wait() // Might return too soon
            fmt.Println("All workers finished (maybe?)")
        }
        ```
    *   **Corrected Code:**
        ```go
        package main

        import (
            "fmt"
            "sync"
            "time"
        )

        func worker(id int, wg *sync.WaitGroup) {
            defer wg.Done()
            fmt.Printf("Worker %d starting\n", id)
            time.Sleep(time.Duration(id) * 100 * time.Millisecond)
            fmt.Printf("Worker %d done\n", id)
        }

        func main() {
            var wg sync.WaitGroup
            for i := 0; i < 3; i++ {
                wg.Add(1) // RIGHT: Call Add before launching goroutine
                go worker(i, &wg)
            }
            wg.Wait()
            fmt.Println("All workers finished!")
        }
        ```

*   **Gotcha 5: Closure Capture**
    *   When a goroutine is created inside a loop and refers to a loop variable, it captures the *variable itself* (by reference), not its value at the time of goroutine creation. This can lead to unexpected behavior as the loop variable mutates.
    *   **Problematic Code:**
        ```go
        package main

        import (
            "fmt"
            "time"
        )

        func main() {
            for i := 0; i < 10; i++ { // WRONG
                go func() {
                    fmt.Println(i) // 'i' is captured by reference
                }()
            }
            time.Sleep(1 * time.Second) // Give goroutines time to run
        }
        ```
        *   Output will likely be `10` printed multiple times, as `i` will be `10` when the goroutines finally execute.
    *   **Corrected Code (Pass as Parameter):**
        ```go
        package main

        import (
            "fmt"
            "time"
        )

        func main() {
            for i := 0; i < 10; i++ { // RIGHT
                go func(i int) { // Pass 'i' as a parameter
                    fmt.Println(i)
                }(i) // Pass the current value of 'i'
            }
            time.Sleep(1 * time.Second)
        }
        ```
    *   **Corrected Code (Local Copy):**
        ```go
        package main

        import (
            "fmt"
            "time"
        )

        func main() {
            for i := 0; i < 10; i++ { // RIGHT
                i := i // Create a new variable 'i' for each iteration
                go func() {
                    fmt.Println(i) // This 'i' is the local copy
                }()
            }
            time.Sleep(1 * time.Second)
        }
        ```

*   **Select Problems**
    *   `select` statement behavior can be counter-intuitive:
        *   `default` is always active: If no other case is ready, the `default` case is executed immediately without blocking.
        *   A `nil` channel is always ignored: Cases involving `nil` channels are never considered ready.
        *   A full channel (for send) is skipped over: If a send case's channel is full, that case is not ready.
        *   A "done" channel is just another channel: It has no special semantic meaning; it behaves like any other channel.
        *   Available channels are selected at random: If multiple cases are ready, Go non-deterministically picks one.
    *   **Mistake 1: Skipping a full channel to default and losing a message:**
        ```go
        for {
            x := socket.Read()
            select {
            case output <- x: // If output is full, this case is not ready
                // ...
            default: // This default case will be taken immediately
                return // Message 'x' is lost
            }
        }
        ```
        *   The programmer might intend to skip `output` only if it's `nil`, but `select` also skips it if it's full.
    *   **Mistake 2: Reading a "done" channel and aborting when input is backed up:**
        ```go
        for {
            select {
            case x := <-input: // If input has buffered messages
                // ...
            case <-done: // If done channel is ready
                return // No guarantee all input is processed before returning
            }
        }
        ```
        *   If both `input` and `done` are ready, `select` will pick one randomly. If `done` is picked, remaining messages in `input`'s buffer are lost.
        *   **Better Practice:** Use `done` only for error aborts. For normal termination (e.g., EOF), close the input channel, which will cause `x, ok := <-input` to return `ok` as `false`.

*   **Some Thoughts / Best Practices**
    1.  **Don't start a goroutine without knowing how it will stop:** Always plan for the termination of your goroutines to prevent resource leaks.
    2.  **Acquire locks/semaphores as late as possible; release them in the reverse order:** This minimizes the critical section and reduces contention. `defer` is ideal for ensuring locks are released in the correct reverse order.
    3.  **Don't wait for non-parallel work that you could do yourself:** Avoid using goroutines for tasks that are inherently sequential or don't benefit from concurrency, as it adds unnecessary overhead.
        ```go
        func do() int {
            ch := make(chan int)
            go func() { ch <- 1 }() // Unnecessary goroutine
            return <-ch
        }
        ```
        This function is effectively sequential and could simply return `1`.
    4.  **Simplify! Review! Test!**
        *   **Simplify:** Strive for the simplest possible concurrent design.
        *   **Review:** Have other experienced Go programmers review your concurrent code, especially if you are new to Go concurrency.
        *   **Test:** Thoroughly test your concurrent code. Assume your code doesn't work until proven otherwise through rigorous testing.

## What's New
*   **Gotcha 5: Closure Capture:** The behavior of loop variables in `for` loops changed in Go 1.22. Previously, variables declared by a `for` loop were created once and updated per iteration, leading to closure capture issues where goroutines captured the final value of the variable. In Go 1.22, each iteration of the loop creates new variables, avoiding these accidental sharing bugs. [7]

## Updated Code Snippets
*   **Gotcha 5: Closure Capture - Problematic Code:** This code is no longer problematic in Go 1.22+ due to the change in loop variable semantics. The output will now be 0-9, as intended, not multiple 10s.
    ```go
    package main

    import (
        "fmt"
        "time"
    )

    func main() {
        // This code is now correct in Go 1.22+
        for i := 0; i < 10; i++ { 
            go func() {
                fmt.Println(i) // 'i' is now a new variable for each iteration
            }()
        }
        time.Sleep(1 * time.Second) // Give goroutines time to run
    }
    ```

## Citations
*   [7] Go version 1.22