# Go Class: 09 Closures

## Summary
This video introduces closures in Go, explaining them as functions that capture and refer to variables from their enclosing scope. It differentiates between variable scope and lifetime, highlighting Go's escape analysis. The core concept of a closure is demonstrated with a Fibonacci generator, showing how captured variables maintain state across function calls. The video also covers the internal representation of closures and a common pitfall related to loop variable capture, providing a clear solution.

## Key Points

*   **What are Closures?**
    *   A closure is a function defined inside another function.
    *   It "closes over" (captures) local variables from its outer (enclosing) function's scope.
    *   The inner function gets a *reference* to these captured variables, allowing it to access and modify them.

*   **Scope vs. Lifetime**
    *   **Scope** is static, determined at compile time, defining where a variable name is visible in the code.
    *   **Lifetime** is dynamic, determined at runtime, defining how long a variable's value exists in memory.
    *   Go's compiler uses **escape analysis** to determine if a local variable's lifetime needs to extend beyond its immediate function's execution. If so, the variable is allocated on the heap instead of the stack, preventing dangling pointers.

*   **Fibonacci Generator Example**
    *   A function can return another function (a closure).
    *   The returned closure retains access to the outer function's local variables, even after the outer function has completed execution.
    *   Each time the outer function is called, it creates a *new, independent set* of captured variables for its returned closure.

    ```go
    func fib() func() int {
        a, b := 0, 1 // Captured by the inner function
        return func() int {
            a, b = b, a+b
            return b
        }
    }

    func main() {
        f := fib() // Creates one Fibonacci generator
        g := fib() // Creates another independent Fibonacci generator

        fmt.Println(f(), f(), f(), f()) // Output: 1 2 3 5
        fmt.Println(g(), g(), g(), g()) // Output: 1 2 3 5 (starts fresh)
    }
    ```

*   **How Closures Work Internally**
    *   A function pointer typically points only to the executable code.
    *   A closure is represented as a data structure containing two pointers:
        *   A pointer to the function's code.
        *   A pointer to an "environment" (or "context") that holds references to the captured variables.
    *   When the closure is invoked, it uses its environment pointer to access the correct instances of the captured variables.

*   **Useful Applications of Closures**
    *   Creating stateful functions (like the Fibonacci generator).
    *   Implementing custom callback functions, where the callback needs access to contextual data without explicit parameters.
    *   Example: Providing a custom comparison function to `sort.Slice`, where the comparison logic depends on data from the outer scope.

    ```go
    type item struct {
        name string
        value int
    }

    func main() {
        items := []item{{"apple", 3}, {"banana", 1}, {"cherry", 2}}

        // The anonymous function is a closure that captures 'items'
        sort.Slice(items, func(i, j int) bool {
            return items[i].value < items[j].value // Sorts by 'value'
        })
        fmt.Println(items) // Output: [{banana 1} {cherry 2} {apple 3}]
    }
    ```

*   **Common Pitfall: Loop Variable Capture**
    *   When creating closures inside a loop, the closure captures a *reference* to the loop variable itself, not its value at each iteration.
    *   If the closures are executed *after* the loop has completed, they will all see the *final* value of the loop variable. This is a common source of bugs, especially with goroutines.

    ```go
    func main() {
        funcs := make([]func(), 4)
        for i := 0; i < 4; i++ {
            funcs[i] = func() {
                fmt.Printf("Value: %d, Address: %p\n", i, &i) // Captures reference to the *same* 'i'
            }
        }
        // Calling funcs later will print '4' for all of them,
        // as 'i' has reached its final value after the loop.
        for _, f := range funcs {
            f()
        }
    }
    ```

*   **Fixing the Loop Variable Capture Problem**
    *   To capture the value of the loop variable for each iteration, create a *new variable* inside the loop's body and assign the loop variable's current value to it.
    *   The closure then captures a reference to this *new, distinct* variable, ensuring it holds the correct value from that specific iteration.

    ```go
    func main() {
        funcs := make([]func(), 4)
        for i := 0; i < 4; i++ {
            // Create a new variable 'val' for each iteration, copying 'i'
            val := i
            funcs[i] = func() {
                fmt.Printf("Value: %d, Address: %p\n", val, &val) // Captures reference to the *new* 'val'
            }
        }
        // Calling funcs now will print 0, 1, 2, 3 respectively,
        // each referring to a distinct memory location for its captured 'val'.
        for _, f := range funcs {
            f()
        }
    }
    ```

## What's New
*   **Loop Variable Capture:** The behavior of `for` loop variables changed in Go 1.22. Previously, variables declared by a `for` loop were created once and updated by each iteration, leading to the "loop variable capture" pitfall where closures captured a reference to the same variable. In Go 1.22, each iteration of the loop creates new variables, avoiding this accidental sharing bug. [7]
*   **Explicit Fix for Loop Variable Capture:** The explicit pattern of creating a new variable inside the loop (e.g., `val := i`) to capture the loop variable's value for each iteration is no longer strictly necessary in Go 1.22 and later, as the language now handles this by default. The code using this pattern will continue to work as expected. [7]

## Updated Code Snippets
```go
// Original "Common Pitfall" code snippet with updated comment for Go 1.22+ behavior.
func main() {
    funcs := make([]func(), 4)
    for i := 0; i < 4; i++ {
        funcs[i] = func() {
            fmt.Printf("Value: %d, Address: %p\n", i, &i) // In Go 1.22+, captures a *new* 'i' for each iteration
        }
    }
    // Calling funcs now will print 0, 1, 2, 3 respectively,
    // each referring to a distinct memory location for its captured 'i'.
    for _, f := range funcs {
        f()
    }
}
```

## Citations
- [7] Go 1.22 Release Notes