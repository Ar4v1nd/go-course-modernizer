# Go Class: 08 Functions, Parameters & Defer

## Summary
This video segment introduces functions in Go, highlighting their "first-class" nature and various capabilities. It delves into function scope, signatures, and a detailed explanation of parameter passing mechanisms (by value vs. by reference), including practical demonstrations with arrays, slices, and maps. Finally, it introduces the `defer` statement, explaining its purpose in ensuring deferred execution and discussing important "gotchas" related to its scope and argument evaluation.

## Key Points

*   **Functions as First-Class Objects:**
    *   Functions can be defined, even inside other functions.
    *   Anonymous function literals can be created.
    *   Functions can be passed as parameters to other functions or returned as values from functions.
    *   They can be stored in variables, slices, and maps (but not as map keys).
    *   They can be sent and received through channels.
    *   Methods can be written against a function type.
    *   Function variables can be compared against `nil`.

*   **Function Scope:**
    *   Almost anything (constants, types, variables, other functions) can be declared inside a function.
    *   Methods, however, cannot be defined inside a function; they must be defined at package scope.

*   **Function Signatures:**
    *   The signature of a function is defined by the order and type of its parameters and return values.
    *   It does *not* depend on the names given to parameters or return values.
    *   Go uses structural typing for functions, meaning two functions have the same type if their signatures match, regardless of their names or parameter names.
    ```go
    // These functions have the same structural type
    func Do(a string, b int) string { /* ... */ }
    func NotDo(x string, y int) string { /* ... */ }
    ```

*   **Parameter Terms:**
    *   **Formal parameters:** The parameters listed in a function's declaration (e.g., `a` and `b` in `func do(a, b int) int`).
    *   **Actual parameters (arguments):** The values passed to a function when it is called (e.g., `1` and `2` in `do(1, 2)`).
    *   **Pass by value:** The function receives a copy of the actual parameter. Changes to this copy inside the function do not affect the original variable in the caller.
    *   **Pass by reference:** The function receives a mechanism (like a pointer) to access and potentially modify the original actual parameter, allowing changes to be visible to the caller.

*   **Parameter Passing in Go (The Ultimate Truth):**
    *   Technically, *all* parameters in Go are passed by value (i.e., by copying something).
    *   However, the *behavior* can differ based on what is copied:
        *   **By Value (behaves like value copy):** Numbers, booleans, arrays, and structs are copied entirely. Modifying them inside the function does not affect the original.
        *   **By Reference (behaves like reference):** Strings, slices, maps, and channels are passed by copying their *descriptor* or *header*. This descriptor contains a pointer to the underlying data (e.g., array for slices, hash table for maps). Modifying the *contents* of the underlying data through this copied descriptor *will* affect the original data, making it *behave* like pass-by-reference.
        *   Explicit pointers (`*T`) are also passed by value (the pointer address is copied), but dereferencing the copied pointer allows modification of the original variable.
    ```go
    // Array: Passed by value (copy of the array)
    func modifyArray(arr [3]int) {
        arr[0] = 0 // Modifies the copy, not the original
    }

    // Slice: Descriptor copied, but points to same underlying array
    func modifySlice(s []int) {
        s[0] = 0 // Modifies the underlying array, visible to caller
    }

    // Map: Descriptor copied, but points to same underlying hash table
    func modifyMap(m map[string]int) {
        m["key"] = 0 // Modifies the underlying hash table, visible to caller
    }

    // Pointer: Pointer address copied, dereferencing modifies original
    func modifyInt(ptr *int) {
        *ptr = 100 // Modifies the original int
    }
    ```

*   **Return Values:**
    *   Go functions can return multiple values.
    *   If a function returns more than one value, the return types must be enclosed in parentheses (e.g., `(int, error)`).
    *   Every `return` statement in a function must specify all return values.

*   **Recursion:**
    *   Go supports recursion, where a function calls itself.
    *   Each recursive call adds a new stack frame to the call stack, storing local variables and parameters for that specific call.
    *   A base case (stopping criteria) is crucial to prevent infinite recursion and stack overflow errors.
    *   While often elegant for problems like tree/graph traversals, recursion can be less efficient than iteration due to the overhead of managing stack frames.

*   **Defer Statement:**
    *   The `defer` statement schedules a function call to be executed just before the surrounding function returns.
    *   It's commonly used for cleanup actions like closing files, unlocking mutexes, or ensuring resources are released.
    *   **Key characteristics:**
        *   Deferred calls are executed in LIFO (Last-In, First-Out) order.
        *   Arguments to the deferred function are evaluated and *copied* at the time the `defer` statement is encountered, not when the deferred function actually runs.
        *   `defer` operates on *function scope*, meaning the deferred call executes when the *enclosing function* exits, not when a specific block (like an `if` statement or `for` loop) ends.
        *   Deferred anonymous functions can modify named return values of the enclosing function, as they share the same scope.
    ```go
    // Example of defer for file closing
    func main() {
        f, err := os.Open("my_file.txt")
        if err != nil {
            // handle error
            return
        }
        defer f.Close() // f.Close() is scheduled to run when main() exits

        // Do something with the file
        // ...
    }

    // Example of defer with argument copying
    func deferDemo() {
        a := 10
        defer fmt.Println("Deferred a:", a) // 'a' is evaluated and copied as 10 here

        a = 20
        fmt.Println("Current a:", a) // Prints 20
        // When deferDemo() exits, the deferred call runs, printing "Deferred a: 10"
    }

    // Example of defer modifying a named return value
    func doIt(a int) (result int) { // 'result' is a named return value
        defer func() {
            result = 2 // This modifies the 'result' variable
        }()
        result = 1 // Initial assignment to result
        return // Returns the value of 'result' after deferred func runs (which is 2)
    }
    ```

## What's New
*   **For Loop Variable Scope Change:** In Go 1.22, the variables declared by a `for` loop are now created anew for each iteration, rather than being created once and updated. This change aims to prevent accidental sharing bugs, particularly when closures (including deferred functions) capture loop variables. While the core principle of `defer` operating on function scope remains, the behavior of closures capturing loop variables has changed to be safer. [7]
*   **`panic(nil)` Behavior Change:** Starting with Go 1.21, calling `panic` with a `nil` interface value (or an untyped `nil`) now causes a run-time panic of type `*runtime.PanicNilError`. This ensures that `recover` (often used within `defer` statements) is guaranteed not to return `nil` when a panic occurs, making panic/recover behavior more predictable. [6]

## Citations
- [6] Go 1.21 Release Notes
- [7] Go 1.22 Release Notes