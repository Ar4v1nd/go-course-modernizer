# Go Class: 10 Slices in Detail

## Summary
This video provides a detailed exploration of Go slices, differentiating between nil and empty slices, explaining the concepts of length and capacity, and demonstrating how these properties behave with slice creation using `make()` and slice operators. It also delves into the internal representation of slices and highlights practical implications, such as JSON encoding differences and potential pitfalls related to slice reallocation during appends.

## Key Points

*   **Nil vs. Empty Slices:**
    *   A nil slice is declared but not initialized (e.g., `var s []int`). Its length and capacity are both 0, and `s == nil` evaluates to `true`. When printed with `%#v`, it appears as `[]int(nil)`. Internally, its address pointer is `nil`.
    *   An empty slice is explicitly initialized as empty (e.g., `t := []int{}`). Its length and capacity are also 0, but `t == nil` evaluates to `false`. When printed with `%#v`, it appears as `[]int{}`. Internally, its address pointer points to a special zero-length internal array, not `nil`.
    *   For practical checks, it's generally best to use `len(s) == 0` to determine if a slice is empty, as this covers both nil and explicitly empty slices.

*   **Creating Slices with `make()`:**
    *   `make([]int, length)`: Creates a slice with the specified `length`. The capacity is implicitly set to the same value as the length. All elements are initialized to their zero value (e.g., `0` for `int`).
        ```go
        u := make([]int, 5) // len(u) = 5, cap(u) = 5, u = {0, 0, 0, 0, 0}
        ```
    *   `make([]int, length, capacity)`: Creates a slice with the specified `length` and `capacity`. The `capacity` must be greater than or equal to `length`. Elements up to `length` are zero-initialized. This is useful for pre-allocating space to avoid frequent reallocations during appends.
        ```go
        v := make([]int, 0, 5) // len(v) = 0, cap(v) = 5, v = {} (but has underlying space for 5 elements)
        ```

*   **Internal Slice Representation (Slice Descriptor):**
    *   A Go slice is internally represented by a "slice descriptor" which is a struct containing three fields:
        *   `len`: The current number of elements in the slice.
        *   `cap`: The maximum number of elements the underlying array can hold, starting from the slice's first element.
        *   `addr`: A pointer to the first element of the underlying array.
    *   Nil slice: `len=0`, `cap=0`, `addr=nil`.
    *   Empty slice (`[]int{}` or `make([]int, 0)`): `len=0`, `cap=0`, `addr` points to a non-nil, zero-length internal array.
    *   `make([]int, 0, 5)`: `len=0`, `cap=5`, `addr` points to an underlying array of 5 zero-initialized integers.

*   **Slice Operators (`[low:high]` and `[low:high:max]`):**
    *   `s[low:high]`: Creates a new slice that references a portion of the original slice's (or array's) underlying array.
        *   The `length` of the new slice is `high - low`.
        *   The `capacity` of the new slice is `capacity_of_original_slice - low`.
        *   This means the new slice might have a capacity greater than its length, allowing for appends without immediate reallocation, but potentially modifying elements beyond its current length in the original underlying array.
        ```go
        a := [3]int{1, 2, 3}
        b := a[0:1] // len(b) = 1, cap(b) = 3 (points to {1, 2, 3})
        ```
    *   `s[low:high:max]`: Provides explicit control over the new slice's capacity.
        *   The `length` of the new slice is `high - low`.
        *   The `capacity` of the new slice is `max - low`.
        *   `max` must be less than or equal to the capacity of the original slice (or array) from `low`.
        *   This is useful to "clip" the capacity of a new slice, ensuring that appends to it will cause a reallocation and thus prevent unintended modifications to the original underlying array.
        ```go
        b := []int{1, 2, 3}
        c := b[0:2:2] // len(c) = 2, cap(c) = 2 (points to {1, 2})
        // Appending to c will now force a reallocation, creating a new underlying array for c.
        ```

*   **Appending to Slices and Reallocation:**
    *   When `append()` is used and the slice's current `length` reaches its `capacity`, Go performs a reallocation.
    *   A new, larger underlying array is allocated (typically doubling the capacity).
    *   All elements from the old underlying array are copied to the new one.
    *   The slice descriptor is updated to point to this new array.
    *   **Important:** If multiple slices share the same underlying array, and one of them triggers a reallocation, that specific slice will then point to a *new* underlying array. The other slices will *still* point to the *original* underlying array. This means modifications to the reallocated slice will no longer affect the original, and vice-versa.
    *   Appending to a nil slice is perfectly valid and will cause Go to allocate an underlying array for it.
        ```go
        var s []int
        s = append(s, 1) // s is now {1}, len=1, cap=1 (or more)
        ```

## What's New
*   **Internal Slice Representation (Slice Descriptor):** The `SliceHeader` type, which is implicitly referred to as the "slice descriptor," has been deprecated. New `unsafe` package functions (`unsafe.SliceData`, `unsafe.StringData`) now provide direct access to the underlying pointer, making the internal representation explicitly accessible without relying on the deprecated `SliceHeader` struct. [5]
*   **Appending to Slices and Reallocation:** The built-in `append()` function's internal formula for determining the new capacity when a reallocation is needed changed in Go 1.17. While the general behavior of reallocation (allocating a new array, copying elements) remains the same, the specific growth factor ("typically doubling the capacity") might be less precise due to the new formula. [2]

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