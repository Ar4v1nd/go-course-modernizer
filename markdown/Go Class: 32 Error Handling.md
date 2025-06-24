# Go Class: 32 Error Handling

## Summary
This video delves into error handling in Go, moving beyond simple string-based errors to custom error types and the advanced features introduced in Go 1.13 like wrapped errors, `errors.Is`, and `errors.As`. It also discusses the philosophy behind Go's error handling approach, differentiating between "normal" (expected) and "abnormal" (logic bug) errors, and advocating for explicit error handling over exceptions and the "fail hard, fail fast" principle for logic bugs.

## Key Points

*   **Simple Errors:**
    *   Most errors in Go are simple strings created using `fmt.Errorf`.
    *   While useful for basic debugging, they lack sophistication for programmatic error inspection.
    ```go
    func (h HAL9000) OpenPodBayDoors() error {
        if h.kill {
            return fmt.Errorf("I'm sorry %s, I can't do that", h.victim)
        }
        // ...
    }
    ```

*   **Error Types:**
    *   Errors in Go are fundamentally objects that satisfy the built-in `error` interface, which requires an `Error() string` method.
    *   Any concrete type can implement this interface to represent a custom error.
    ```go
    type error interface {
        Error() string
    }

    type Fizgig struct{}
    func (f Fizgig) Error() string {
        return "Your fizgig is bent"
    }
    ```

*   **Custom Error Type Definition:**
    *   A custom error type (`WaveError`) can be defined as a struct to hold additional context, such as an `errKind` (an enumerated type for specific error categories), a `value` (e.g., position), and an `err` field for wrapping underlying errors.
    *   The `Error()` method for a custom type can use a `switch` statement on the `errKind` to provide different, context-rich error messages.
    ```go
    type errKind int
    const (
        _ errKind = iota // so we don't start at 0
        noHeader
        cantReadHeader
        invalidHdrType
        invalidChkLength
        // ...
    )

    type WaveError struct {
        kind errKind
        value int
        err error
    }

    func (e WaveError) Error() string {
        switch e.kind {
        case noHeader:
            return "no header (file too short?)"
        case cantReadHeader:
            return fmt.Sprintf("can't read header[%d]: %s", e.value, e.err.Error())
        case invalidHdrType:
            return "invalid header type"
        case invalidChkLength:
            return fmt.Sprintf("invalid chunk length: %d", e.value)
        }
        return "unknown wave error" // Default case
    }
    ```

*   **Helper Methods for Custom Errors:**
    *   Methods like `with(val int)` and `from(pos int, err error)` can be added to the custom error type to create new error instances with specific contextual data, promoting immutability by returning a copy.
    ```go
    func (e WaveError) with(val int) WaveError {
        el := e
        el.value = val
        return el
    }

    func (e WaveError) from(pos int, err error) WaveError {
        el := e
        el.value = pos
        el.err = err
        return el
    }
    ```

*   **Prototype Errors:**
    *   Exported error variables can be declared as pre-initialized instances of the custom error type, serving as prototypes for common error conditions.
    ```go
    var (
        HeaderMissing      = WaveError{kind: noHeader}
        HeaderReadFailed   = WaveError{kind: cantReadHeader}
        InvalidHeaderType  = WaveError{kind: invalidHdrType}
        InvalidChunkLength = WaveError{kind: invalidChkLength}
        InvalidDataLength  = WaveError{kind: invalidLength}
    )
    ```

*   **Custom Error Type in Use:**
    *   Functions can return these prototype errors directly or use helper methods to add specific details before returning.
    ```go
    func DecodeHeader(b []byte) (*Header, []byte, error) {
        var err error
        var pos int
        // ...
        if len(b) < HeaderSize {
            return &header, nil, HeaderMissing
        }
        // ...
        if err = binary.Read(buf, binary.BigEndian, &header.riff); err != nil {
            return &header, nil, HeaderReadFailed.from(pos, err)
        }
        // ...
    }
    ```

*   **Wrapped Errors (Go 1.13+):**
    *   Go 1.13 introduced the ability to wrap one error within another, creating a chain of errors.
    *   The `%w` format verb with `fmt.Errorf()` is the easiest way to wrap errors.
    ```go
    func (h HAL9000) OpenPodBayDoors() error {
        if h.err != nil {
            return fmt.Errorf("I'm sorry %s, I can't: %w", h.victim, h.err)
        }
        // ...
    }
    ```

*   **Unwrapping Errors:**
    *   Custom error types can implement an `Unwrap() error` method to expose the underlying wrapped error, allowing traversal of the error chain.
    ```go
    func (w *WaveError) Unwrap() error {
        return w.err
    }
    ```

*   **`errors.Is` (Go 1.13+):**
    *   The `errors.Is(err, target)` function checks if `err` (or any error in its chain) is semantically equivalent to `target`.
    *   It compares error *variables* (not types) and traverses the error chain using `Unwrap()`.
    *   Custom error types can implement their own `Is(target error) bool` method for custom comparison logic (e.g., comparing internal `errKind` values).
    ```go
    if audio, err := DecodeWaveFile(fn); err != nil {
        if errors.Is(err, os.ErrPermission) {
            // let's report a security violation
        }
    }

    // Custom Is method for WaveError
    func (w *WaveError) Is(t error) bool {
        e, ok := t.(WaveError) // Using reflection to check type
        if !ok {
            return false
        }
        return e.errKind == w.errKind // Compare based on internal kind
    }
    ```

*   **`errors.As` (Go 1.13+):**
    *   The `errors.As(err, &target)` function checks if `err` (or any error in its chain) can be assigned to `target` (which must be a pointer to an error type).
    *   If a match is found, it copies the matching error from the chain into `target` and returns `true`.
    *   This allows extracting specific error types from a wrapped chain to access their concrete fields.
    ```go
    if audio, err := DecodeWaveFile(fn); err != nil {
        var e os.PathError // a struct
        if errors.As(err, &e) {
            // let's pass back just the underlying file error
            return e
        }
    }
    ```

*   **Philosophy of Error Handling:**
    *   **Normal Errors:** Result from expected input or external conditions (e.g., file not found, network down). Go handles these by returning the `error` type, requiring explicit checking (`if err != nil`). This promotes visibility and deliberate handling.
    *   **Abnormal Errors:** Result from invalid program logic or internal inconsistencies (e.g., nil pointer dereference, off-by-one errors). Go handles these with `panic`.
    ```go
    func (d *digest) checkSum() [Size]byte {
        // finish writing the checksum
        // ...
        if d.nx != 0 { // panic if there's data left over
            panic("d.nx != 0")
        }
        // ...
    }
    ```
    *   **"Fail hard, fail fast":** For abnormal errors (bugs), it's often better for the program to crash immediately. This surfaces bugs quickly during development/testing, provides clear stack traces for debugging, and prevents corrupted states in production (especially critical in distributed systems where silent failures can lead to Byzantine failures).
    *   **Exception Handling (Go's perspective):** While Go has `panic` and `recover` (similar to exceptions), it discourages their general use for error handling. Exceptions introduce invisible control paths, making code harder to analyze and reason about.
    *   **Proactively Prevent Problems:**
        *   Design abstractions to make operations inherently safe (e.g., Go's nil map reads, slice appends).
        *   Break down complex logic into smaller, understandable pieces.
        *   Hide information to reduce corruption chances.
        *   Avoid clever/unsafe code.
        *   Assert invariants (using `panic` for violations).
        *   Never ignore errors.
        *   **Test, test, test:** Comprehensive testing is crucial to catch logic bugs early.
        *   Always validate user/environment input.
    *   **Error Handling Culture:** Go's explicit error handling (`if err != nil`) fosters a culture where programmers actively think about and handle failure cases at the point of writing code, rather than deferring them. This verbosity leads to greater visibility and deliberate handling of every error condition.

## What's New
*   **Wrapped Errors (Go 1.13+):** Go 1.20 expanded error wrapping to support multiple wrapped errors. `fmt.Errorf` now supports multiple `%w` verbs, and the new `errors.Join` function was introduced to create errors that wrap a list of errors [5].
*   **Unwrapping Errors:** The `Unwrap()` method can now return `[]error` to expose multiple underlying wrapped errors, consistent with the multiple error wrapping feature [5].
*   **`errors.Is` (Go 1.13+):** The `vet` tool now warns about `Is` methods (and `As`, `Unwrap`) that have a different signature than the one expected by the `errors` package, helping to catch incorrect implementations [2].
*   **`errors.As` (Go 1.13+):** The `vet` tool now warns about `As` methods (and `Is`, `Unwrap`) that have a different signature than the one expected by the `errors` package, helping to catch incorrect implementations [2].

## Updated Code Snippets
```go
// For Unwrap() to support multiple wrapped errors (Go 1.20+)
func (w *WaveError) Unwrap() []error {
    // Assuming w.err could be a single error or an error that itself wraps multiple.
    // For simplicity, if w.err is a single error, return it in a slice.
    if w.err == nil {
        return nil
    }
    return []error{w.err}
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