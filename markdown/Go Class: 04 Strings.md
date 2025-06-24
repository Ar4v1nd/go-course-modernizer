# Go Class: 04 Strings

## Summary
This video provides an in-depth look at strings in Go, highlighting their dual nature as sequences of bytes (physically) and Unicode runes (logically). It covers the `byte` and `rune` types, string immutability, internal memory representation, and common string operations and functions, including a practical demonstration of a simple search-and-replace program.

## Key Points

*   **Strings have two natures in Go:**
    *   **Physical:** A string is a sequence of bytes, specifically UTF-8 encoded bytes.
    *   **Logical:** A string represents a sequence of Unicode characters, known as `runes`.
*   **`byte` and `rune` types:**
    *   `byte`: A synonym for `uint8`, representing a single 8-bit byte.
    *   `rune`: A synonym for `int32`, representing a Unicode code point (character). This is Go's equivalent of a "wide character" in other languages.
*   **UTF-8 Encoding:**
    *   Go strings are always UTF-8 encoded.
    *   UTF-8 is a variable-width encoding, meaning a single Unicode `rune` can be represented by 1 to 4 bytes.
    *   ASCII characters (0-127) are represented by a single byte in UTF-8.
    *   Non-ASCII characters (like `é`) require multiple bytes (e.g., `é` is `233` as a rune, but `195 169` as UTF-8 bytes).
*   **String Length (`len()`):**
    *   The built-in `len()` function returns the number of *bytes* in a string, not the number of logical characters (runes).
    *   Example:
        ```go
        s := "élite" // 5 runes, but 'é' takes 2 bytes in UTF-8
        fmt.Printf("%T %v\n", s, s) // Output: string élite
        fmt.Printf("%T %v\n", []rune(s), []rune(s)) // Output: []int32 [233 108 105 116 101] (5 runes)
        fmt.Printf("%T %v %d\n", []byte(s), []byte(s), len(s)) // Output: []uint8 [195 169 108 105 116 101] 6 (6 bytes)
        ```
*   **String Immutability:**
    *   Strings in Go are **immutable**. Once created, their underlying byte sequence cannot be changed.
    *   Attempting to modify a character by index will result in a compile-time error:
        ```go
        s := "a string"
        // s[5] = 'a' // SYNTAX ERROR: cannot assign to s[5] (strings are immutable)
        ```
*   **Internal String Representation:**
    *   Internally, a string is represented by a "string descriptor," which is a small data structure containing a pointer to the underlying byte array and the length (number of bytes) of the string.
    *   This design allows for efficient substring operations and memory sharing.
*   **Substring (Slicing):**
    *   Slicing a string (e.g., `s[start:end]`) creates a *new string descriptor* that points to a portion of the *original string's underlying byte array*. No new memory is allocated for the characters themselves.
    *   Example:
        ```go
        s := "hello, world"
        hello := s[0:5] // hello points to the first 5 bytes of s
        world := s[7:12] // world points to bytes 7-11 of s
        ```
*   **String Concatenation:**
    *   Concatenating strings (e.g., `s + "es"` or `s += "es"`) creates a *new string* in a *new memory location*. The original string's underlying data is copied to the new location, and the new characters are appended.
    *   The original string remains unchanged in its memory location. If no other variables reference the original string's data, it becomes eligible for garbage collection.
*   **String Functions (from `strings` package):**
    *   The `strings` package provides many useful functions for string manipulation.
    *   These functions generally return *new* strings, adhering to the immutability principle.
    *   Examples:
        ```go
        s := "a string"
        strings.Contains(s, "g")     // returns true
        strings.HasPrefix(s, "a")    // returns true
        strings.Index(s, "string")   // returns 2 (index of first occurrence)
        s = strings.ToUpper(s)       // returns "A STRING" and reassigns s
        ```
*   **Practical Example: Simple Search and Replace:**
    *   A program can read input line by line using `bufio.NewScanner(os.Stdin)`.
    *   `scanner.Scan()` reads a line, and `scanner.Text()` retrieves it.
    *   `strings.Split(line, old)` splits a string into a slice of strings using `old` as the delimiter.
    *   `strings.Join(slice, new)` joins a slice of strings into a single string, inserting `new` between the elements.
    *   This combination effectively replaces all occurrences of `old` with `new` in a line.
    *   Example (conceptual, as `os.Args` and `bufio` are used):
        ```go
        package main

        import (
            "bufio"
            "fmt"
            "os"
            "strings"
        )

        func main() {
            if len(os.Args) < 3 {
                fmt.Fprintln(os.Stderr, "Not enough args")
                os.Exit(-1)
            }

            old := os.Args[1]
            new := os.Args[2]

            scanner := bufio.NewScanner(os.Stdin)
            for scanner.Scan() {
                s := strings.Split(scanner.Text(), old)
                t := strings.Join(s, new)
                fmt.Println(t)
            }
        }
        ```

## What's New

*   **Internal String Representation:** The `StringHeader` type, which is implicitly referred to as the "string descriptor," has been deprecated. New code should prefer using functions from the `unsafe` package like `unsafe.SliceData`, `unsafe.String`, or `unsafe.StringData` for constructing and deconstructing string values without depending on their exact representation. [5], [6]
*   **String Functions (from `strings` package):** The `strings.Title` function has been deprecated. It doesn't handle Unicode punctuation and language-specific capitalization rules, and is superseded by the `golang.org/x/text/cases` package. [3]

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