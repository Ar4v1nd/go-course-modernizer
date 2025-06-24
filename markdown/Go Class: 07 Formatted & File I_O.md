# Go Class: 07 Formatted & File I/O

## Summary
This video covers formatted input/output using Go's `fmt` package, including various print functions and format verbs. It then delves into file I/O, demonstrating how to interact with files and directories using packages like `os`, `io`, `bufio`, `io/ioutil`, and `strings`. A key emphasis is placed on proper error handling in Go for I/O operations.

## Key Points

*   **Standard I/O Streams:**
    *   Unix systems have three standard I/O streams: Standard Input, Standard Output, and Standard Error.
    *   These are open by default in every program.
    *   Go exposes these through the `os` package: `os.Stdin`, `os.Stdout`, `os.Stderr`.

*   **Formatted I/O with `fmt` Package:**
    *   The `fmt` package provides functions for formatted I/O, often using reflection to print various data types.
    *   **`Println` and `Printf` family:**
        *   `fmt.Println(a ...interface{}) (n int, err error)`: Prints arguments to `os.Stdout` with spaces between them and a newline at the end.
        *   `fmt.Printf(format string, a ...interface{}) (n int, err error)`: Prints formatted output to `os.Stdout` based on a format string.
        *   `fmt.Fprintln(w io.Writer, a ...interface{}) (n int, err error)`: Same as `Println` but prints to a specified `io.Writer` (e.g., `os.Stderr`).
        *   `fmt.Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error)`: Same as `Printf` but prints to a specified `io.Writer`.
        *   `fmt.Sprintln(a ...interface{}) string`: Returns the formatted string instead of printing it.
        *   `fmt.Sprintf(format string, a ...interface{}) string`: Returns the formatted string based on a format string.

*   **Common Format Verbs (from `fmt.Printf`):**
    *   `%s`: the uninterpreted bytes of the string or slice.
    *   `%q`: a double-quoted string, safely escaped with Go syntax.
    *   `%c`: the character represented by the corresponding Unicode code point.
    *   `%d`: base 10 integer.
    *   `%x` / `%X`: base 16 (hexadecimal) with lowercase/uppercase letters for a-f.
    *   `%f`: decimal point float (e.g., 123.456).
    *   `%t`: the word `true` or `false` for booleans.
    *   `%v`: the value in a default format. When printing structs, the plus flag (`%+v`) adds field names.
    *   `%#v`: a Go-syntax representation of the value.
    *   `%T`: a Go-syntax representation of the type of the value.
    *   `%%`: a literal percent sign (consumes no value, escape).
    *   **Width and Precision:** Numbers can be formatted with width (e.g., `%6d` for 6 characters wide) and precision (e.g., `%.2f` for 2 decimal places).
    *   **Flags:**
        *   `0`: zero-pads numbers (e.g., `%06d`).
        *   `-`: left-justifies within the field (e.g., `%-6d`).
        *   `#`: alternate format (e.g., `%#x` for `0x` prefix, `%#v` for Go-syntax representation).

*   **File I/O Packages:**
    *   `os` package: Provides functions to open or create files, list directories, and hosts the `os.File` type.
    *   `io` package: Offers utilities for reading and writing data streams.
    *   `bufio` package: Provides buffered I/O operations, including `bufio.Scanner` for efficient line-by-line reading.
    *   `io/ioutil` package: Contains extra utilities like `ReadAll` (reads entire file into memory) and `WriteFile` (writes data all at once).
    *   `strconv` package: Has utilities to convert to/from string representations (e.g., `ParseInt`, `ParseFloat`, `ParseBool`).

*   **Reading a File - Best Practices:**
    *   **Error Handling:** Always check the error returned by I/O functions. Go functions often return a value and a possible `error`.
        ```go
        file, err := os.Open(fname)
        if err != nil {
            fmt.Fprintln(os.Stderr, "bad file:", err)
            return // or continue, depending on desired behavior
        }
        // Always close the file when done
        defer file.Close()
        ```
    *   **`io.Copy` for efficient copying:**
        ```go
        // Copies content from 'file' to 'os.Stdout' efficiently
        _, err = io.Copy(os.Stdout, file)
        if err != nil {
            fmt.Fprintln(os.Stderr, "copy error:", err)
        }
        ```
    *   **`bufio.Scanner` for line-by-line processing:**
        ```go
        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
            line := scanner.Text()
            // Process each line
            fmt.Println(line)
        }
        if err := scanner.Err(); err != nil {
            fmt.Fprintln(os.Stderr, "scanning error:", err)
        }
        ```
    *   **`strings.Fields` for word counting:**
        ```go
        words := strings.Fields(line) // Splits line by whitespace into a slice of words
        wordCount += len(words)
        ```
    *   **`io/ioutil.ReadAll` for small files:**
        ```go
        data, err := ioutil.ReadAll(file) // Reads entire file into a []byte slice
        if err != nil {
            fmt.Fprintln(os.Stderr, "readall error:", err)
        }
        // Process data as a byte slice
        fmt.Println("File has", len(data), "bytes")
        ```
        *Note: `ReadAll` is convenient but can consume large amounts of memory for big files, potentially crashing the program.*

## What's New

*   The `io/ioutil` package has been deprecated since Go 1.16. Its functionalities, including `ReadAll` and `WriteFile`, have been moved to the `io` and `os` packages. [1]

## Updated Code Snippets

For `io/ioutil.ReadAll`:

```go
data, err := io.ReadAll(file) // Reads entire file into a []byte slice
if err != nil {
    fmt.Fprintln(os.Stderr, "readall error:", err)
}
// Process data as a byte slice
fmt.Println("File has", len(data), "bytes")
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