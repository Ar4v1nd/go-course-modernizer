# Go Class: 03 Basic Types

## Summary
This video introduces fundamental data types in Go, contrasting Go's machine-native approach with interpreted languages. It covers integer, floating-point, boolean, and error types, along with variable declaration methods and the concept of constants. A practical example demonstrates reading numerical input, calculating averages, and handling basic type conversions and errors.

## Key Points

*   **Keywords & Predeclared Identifiers:**
    *   Go has a small set of 25 keywords, indicating its simplicity.
    *   It also includes predeclared identifiers for constants (`true`, `false`, `iota`, `nil`), types (`int`, `uint`, `float64`, `bool`, `string`, `error`, etc.), and built-in functions (`make`, `len`, `append`, etc.). While these can be shadowed, it's generally not recommended.

*   **Machine-Native vs. Interpreted Languages:**
    *   Go is a compiled, machine-native language, meaning variables directly represent values in machine memory (RAM/CPU). This leads to performance advantages.
    *   Interpreted languages (like Python) use an interpreter layer where variables are objects that masquerade as numbers, requiring an extra step to interact with machine hardware.

*   **Integers:**
    *   Go provides signed (`int`, `int8`, `int16`, `int32`, `int64`) and unsigned (`uint`, `uint8`, `uint16`, `uint32`, `uint64`, `uintptr`) integer types.
    *   `int` is the default integer type and its size (32 or 64 bits) depends on the machine's natural word size. For most modern systems (laptops, cloud servers), `int` is 64 bits, offering a very large range.

*   **Floating-Point Numbers:**
    *   Non-integer numbers are represented as floating-point types: `float32` and `float64`.
    *   `float64` is the default floating-point type.
    *   Go also supports complex numbers: `complex64` and `complex128`.
    *   **Best Practice:** Avoid using floating-point numbers for monetary calculations due to potential precision errors (e.g., representing 0.1 in binary can be repeating, leading to rounding issues). Use specialized packages for financial calculations.

*   **Simple Declarations:**
    *   **`var` keyword:** Used for explicit type declaration or grouped declarations with type inference.
        ```go
        var a int // Explicit type
        var (
            b = 2    // Type inferred as int
            f = 2.01 // Type inferred as float64
        )
        ```
    *   **Short declaration operator `:=`:** A concise way to declare and initialize variables, with type inference. It can only be used *inside* functions.
        ```go
        c := 2 // Type inferred as int
        ```
    *   **Type Strictness:** Go is a statically typed language and does not allow implicit type conversions between different types (e.g., `int` and `float64`). Explicit type casting is required.
        ```go
        var myInt int = 5
        var myFloat float64 = 3.14
        myInt = int(myFloat)   // Explicit conversion: myInt becomes 3 (truncation)
        myFloat = float64(myInt) // Explicit conversion: myFloat becomes 5.0
        ```

*   **Special Types:**
    *   **`bool`:** Represents boolean values (`true` or `false`). Unlike some languages (e.g., C), boolean values are not directly convertible to/from integers.
    *   **`error`:** A built-in interface type used for error handling. An `error` variable can be `nil` (indicating no error) or non-`nil` (indicating an error occurred).
    *   **Pointers:** Represent memory addresses. A pointer can be `nil` (pointing to nothing) or non-`nil`. Go restricts direct pointer arithmetic and manipulation to the `unsafe` package for safety.

*   **Initialization:**
    *   All variables in Go are automatically initialized to their "zero value" if not explicitly initialized by the programmer. This prevents common bugs related to uninitialized memory.
    *   Zero values:
        *   Numerical types: `0` (e.g., `0` for `int`, `0.0` for `float64`, `0+0i` for `complex`).
        *   `bool`: `false`.
        *   `string`: `""` (empty string).
        *   Pointers, slices, maps, channels, function variables, interfaces: `nil`.

*   **Constants:**
    *   Declared using the `const` keyword.
    *   Only numbers, strings, and booleans can be declared as constants.
    *   Constants are immutable; their values are fixed at compile time and cannot be changed during program execution.
    *   This immutability is a key feature for ensuring concurrency safety in Go programs.
    *   Constants can be literals or the result of compile-time constant expressions (e.g., `const b = 2 * 1024`).

*   **Practical Example: Calculating Averages:**
    *   Demonstrates reading numbers from standard input (`os.Stdin`) using `fmt.Fscanln`.
    *   Uses the address-of operator (`&`) to pass a pointer to the variable where the scanned value should be stored.
    *   Includes basic error handling for input operations.
    *   Illustrates type casting for arithmetic operations (e.g., `sum / float64(n)`).
    *   Shows how to handle the case of no input values to prevent division by zero.
    *   Uses `os.Stderr` for printing error messages, which is good practice for command-line tools.
    *   Demonstrates running the Go program with input from the console and by redirecting input from a file.
    *   Uses `sum += val` (compound assignment) and `n++` (increment operator). Note that `n++` is a statement in Go and cannot be part of an expression.

## What's New

*   **Keywords & Predeclared Identifiers:**
    *   Go 1.18 introduced `any` and `comparable` as new predeclared identifiers [3].
    *   Go 1.21 introduced `min`, `max`, and `clear` as new built-in functions [6].

*   **Integers:**
    *   Go 1.19 introduced new atomic integer types (`atomic.Int32`, `atomic.Int64`, `atomic.Uint32`, `atomic.Uint64`, `atomic.Uintptr`) in the `sync/atomic` package, which provide atomic operations for integer values [4].

*   **Special Types - `error`:**
    *   Go 1.20 expanded error wrapping to allow an error to wrap multiple other errors, with `errors.Is`, `errors.As`, and `fmt.Errorf` updated to support this. The `errors.Join` function was also added [5].
    *   Go 1.21 changed the behavior of `panic` with a `nil` interface value; it now causes a `*runtime.PanicNilError` instead of returning `nil` from `recover` [6].

*   **Special Types - Pointers:**
    *   Go 1.17 introduced `unsafe.Add` and `unsafe.Slice` functions for more controlled pointer manipulation [2].
    *   Go 1.18's change to passing function arguments and results using registers may affect code that violates `unsafe.Pointer` rules or depends on undocumented behavior of comparing function code pointers [3].
    *   Go 1.20 added `unsafe.SliceData`, `unsafe.String`, and `unsafe.StringData` to provide complete functionality for constructing and deconstructing slice and string values without relying on their exact representation [5].
    *   Go 1.21 introduced the `runtime.Pinner` type, which allows "pinning" Go memory for safer use by non-Go code, including passing Go values that reference pinned memory to C code, which was previously disallowed by cgo pointer passing rules [6].

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