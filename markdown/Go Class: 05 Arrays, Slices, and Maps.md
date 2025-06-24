# Go Class: 05 Arrays, Slices, and Maps

## Summary
This video introduces composite data types in Go: arrays, slices, and maps. It explains their fundamental differences, how they are declared and initialized, and how they behave when passed as arguments or assigned. The video emphasizes the unique properties of slices and maps, including their variable length, pass-by-reference semantics, and the utility of nil values in Go.

## Key Points

*   **Composite Types Overview**
    *   Go's composite types are containers for other data types.
    *   **Strings:** Immutable sequences of bytes.
    *   **Arrays:** Fixed-size, contiguous sequences of elements of the same type.
    *   **Slices:** Variable-length, flexible views into underlying arrays.
    *   **Maps:** Unordered collections of key-value pairs, where keys are unique.

*   **Arrays**
    *   Arrays in Go are typed by both their element type and their fixed size, determined at compile time.
    *   **Declaration:**
        ```go
        var a [3]int // Declares an array 'a' of 3 integers, initialized to zero values.
        var b [3]int{0, 0, 0} // Declares and initializes 'b' with specific values.
        var c [...]int{0, 0, 0} // Compiler infers size from initializer (here, 3).
        ```
    *   **Pass-by-Value:** When an array is assigned or passed to a function, its entire contents are copied.
    *   **Type Mismatch:** Arrays with different sizes, even if they hold the same element type, are considered different types.
    *   Arrays are less commonly used directly in Go due to their fixed size and copying behavior, especially for large data sets.

*   **Slices**
    *   Slices provide a dynamic, variable-length view into an underlying array. They are more flexible than arrays.
    *   **Slice Descriptor:** A slice is a small data structure containing a pointer to the underlying array, its length, and its capacity (the maximum length it can grow to without reallocation).
    *   **Pass-by-Reference:** When a slice is assigned or passed to a function, its descriptor is copied, but both the original and the copy point to the same underlying array. Changes to elements through one slice are reflected in the other.
    *   **Declaration & Initialization:**
        ```go
        var a []int // Declares a nil slice (no storage, length 0, capacity 0).
        var b = []int{1, 2} // Declares and initializes a slice with an underlying array.
        ```
    *   **`append` Function:** Used to add elements to a slice. It returns a *new* slice, which might point to a new, larger underlying array if the original capacity is exceeded. It's crucial to reassign the result.
        ```go
        a = append(a, 1) // Appends 1 to 'a'. If 'a' was nil, it's now []int{1}.
        b = append(b, 3) // Appends 3 to 'b'. 'b' is now []int{1, 2, 3}.
        ```
    *   **`make` Function:** Used to create slices with a specified length and optional capacity.
        ```go
        d := make([]int, 5) // Creates a slice of 5 zero-valued integers.
        ```
    *   **Slicing Operator:** Used to create new slices from existing arrays or slices. The syntax `[low:high]` creates a slice from `low` (inclusive) up to `high` (exclusive).
        ```go
        // Example: t[2] gets the element at index 2.
        // Example: t[:2] gets elements from index 0 up to (but not including) 2.
        // Example: t[2:] gets elements from index 2 up to the end.
        // Example: t[3:5] gets elements from index 3 up to (but not including) 5.
        ```
    *   **Off-by-one Bug:** The half-open interval `[low:high)` is consistent with `for` loop ranges (`for i := low; i < high; i++`).

*   **Slices vs Arrays (Comparison)**
    | Feature             | Slice                                | Array                                     |
    | :------------------ | :----------------------------------- | :---------------------------------------- |
    | Length              | Variable length                      | Fixed at compile time                     |
    | Passing             | Passed by reference (descriptor copy)| Passed by value (full copy)               |
    | Comparability (`==`)| Not comparable                       | Comparable (if elements are comparable)   |
    | Map Key             | Cannot be used as map key            | Can be used as map key                    |
    | Helpers             | `copy`, `append` helpers             | No built-in helpers for resizing          |
    | Use Case            | Useful as function parameters        | Useful as "pseudo" constants (fixed data) |

*   **Arrays as Pseudo-Constants**
    *   Arrays can be useful for representing fixed-size tables of values that are treated as constant data within an algorithm.
    *   Example: A permutation table in a cryptographic algorithm (like DES) can be stored as a fixed-size array.
    *   Since arrays are passed by value, copying these small, fixed-size tables is efficient and ensures immutability within function scopes.

*   **Maps**
    *   Maps are Go's built-in hash tables (dictionaries).
    *   **Declaration:**
        ```go
        var m map[string]int // Declares 'm' as a nil map (no storage allocated).
        ```
    *   **Initialization with `make`:** A nil map cannot be written to; it must be initialized using `make` to allocate underlying storage.
        ```go
        p := make(map[string]int) // Creates an empty, non-nil map.
        ```
    *   **Reading from Maps:**
        *   Reading from a nil map or a non-existent key returns the zero value of the value type (e.g., `0` for `int`).
        *   ```go
            a := p["the"] // If "the" is not in 'p', 'a' will be 0.
            b := m["the"] // If 'm' is nil, 'b' will also be 0.
            ```
    *   **Writing to Maps:**
        *   Writing to a nil map will cause a `panic` (runtime error).
        *   ```go
            // m["and"] = 1 // This would panic if 'm' is nil.
            p["and"] = 1 // This is OK, as 'p' was made.
            ```
    *   **Pass-by-Reference:** Maps are passed by reference, meaning changes within a function affect the original map.
    *   **Key Type Restriction:** The type used for a map's key *must* be comparable (i.e., support `==` and `!=` operators). Slices, maps, and functions are *not* comparable and thus cannot be used as map keys.

*   **Maps - Special Two-Result Lookup**
    *   To differentiate between a missing key and a key with a zero value, maps provide a two-result lookup:
        ```go
        value, ok := myMap[key]
        ```
        `value` holds the element's value, and `ok` is a boolean that is `true` if the key was found, `false` otherwise.
    *   This pattern is very common in Go for safely accessing map elements.
        ```go
        // Example:
        // p is a map[string]int, initialized as non-nil but empty.
        // p["the"]++ // This would insert "the" with value 1.
        // c, ok := p["the"] // c will be 1, ok will be true.
        // w, ok := p["missing"] // w will be 0, ok will be false.

        if w, ok := p["the"]; ok {
            // This block executes only if "the" key exists in 'p'.
            // 'w' holds the value, 'ok' is true.
        }
        ```

*   **Built-in Functions**
    *   Go provides several built-in functions that work across different composite types:
        *   `len(s)`: Returns the length of strings, arrays, slices, and maps.
        *   `cap(a)`: Returns the capacity of arrays and slices. (Arrays' capacity is always equal to their length).
        *   `make(T, x)`: Creates slices (length x, capacity x) or maps (empty, non-nil).
        *   `make(T, x, y)`: Creates slices (length x, capacity y).
        *   `copy(dst, src)`: Copies elements from `src` slice to `dst` slice; returns the number of elements copied (minimum of the two lengths).
        *   `append(s, elems...)`: Appends elements to a slice `s`; returns a *new* slice.
        *   `delete(m, k)`: Deletes key `k` from map `m`. It's safe to delete a non-existent key.

*   **Make Nil Useful**
    *   In Go, `nil` is the zero value for pointer types, slices, maps, channels, and interfaces. It indicates the absence of a value or an uninitialized state.
    *   Go's philosophy is to make `nil` values useful and safe to work with, avoiding common null pointer exceptions found in other languages.
    *   Many built-in functions (like `len`, `cap`, `range`) safely handle `nil` inputs without panicking.
    *   Example: `len(nil_slice)` returns `0`. Iterating over a `nil` slice with `for...range` results in zero iterations.
    *   This design simplifies code by reducing the need for explicit `nil` checks in many common scenarios.

## What's New
*   **Slices - Helpers:** The standard library now includes a new `slices` package, introduced in Go 1.21, which provides many common operations on slices using generic functions. This expands the built-in helpers available for slices beyond `copy` and `append`. [6]
*   **Built-in Functions:** A new built-in function `clear` was added in Go 1.21. It deletes all elements from a map or zeroes all elements of a slice. [6]

## Updated Code Snippets
```go
// New built-in function in Go 1.21:
clear(myMap)   // Deletes all entries from myMap
clear(mySlice) // Zeroes all elements of mySlice
```

## Citations
*   [1] Go 1.16 Release Notes
*   [2] Go 1.17 Release Notes
*   [3] Go 1.18 Release Notes
*   [4] Go 1.19 Release Notes
*   [5] Go 1.20 Release Notes
*   [6] Go 1.21 Release Notes
*   [7] Go 1.22 Release Notes
*   [8] Go 1.23 Release Notes
*   [9] Go 1.24 Release Notes