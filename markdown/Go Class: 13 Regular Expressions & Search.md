# Go Class: 13 Regular Expressions & Search

## Summary
This video provides an introduction to string searching and regular expressions in Go. It covers basic string manipulation using the `strings` package and then delves into more powerful pattern matching and replacement capabilities offered by the `regexp` package. The instructor emphasizes the importance of careful use of regular expressions due to their complexity and potential for performance issues, while also highlighting Go's specific `regexp` implementation (RE2) designed to mitigate some common pitfalls.

## Key Points

*   **Searching in Strings Overview**
    *   Use the `strings` package for simple, exact string searches.
    *   Use the `regexp` package for complex pattern matching and validation.
    *   Be cautious with regular expressions; complex patterns can be difficult to understand, debug, and test.
    *   Go's `regexp` syntax is a subset of what some other languages offer (based on RE2) to avoid performance issues like catastrophic backtracking.

*   **Simple String Searches**
    *   The `strings` package provides functions for common string operations:
        *   `strings.HasPrefix(s, substr)`: Checks if `s` starts with `substr`.
        *   `strings.HasSuffix(s, substr)`: Checks if `s` ends with `substr`.
        *   `strings.Contains(s, substr)`: Checks if `s` contains `substr`.
        *   `strings.LastIndex(s, substr)`: Returns the last index of `substr` in `s`.
        *   `strings.LastIndexByte(s, char)`: Returns the last index of `char` (as a byte) in `s`.
        *   `strings.Replace(s, old, new, n)`: Replaces the first `n` occurrences of `old` with `new`.
        *   `strings.ReplaceAll(s, old, new)`: Replaces all occurrences of `old` with `new`.
    *   Example using `strings.ReplaceAll`:
        ```go
        package main

        import (
        	"fmt"
        	"strings"
        )

        func main() {
        	test := "Here is $1 which is $2!"
        	test = strings.ReplaceAll(test, "$1", "honey")
        	test = strings.ReplaceAll(test, "$2", "tasty")
        	fmt.Println(test)
        }
        // Output: Here is honey which is tasty!
        ```
    *   Example using `strings.LastIndexByte` to extract a filename:
        ```go
        package main

        import (
        	"fmt"
        	"runtime"
        	"strconv"
        	"strings"
        )

        func B() string {
        	_, file, line, _ := runtime.Caller(1)
        	idx := strings.LastIndexByte(file, '/')
        	return file[idx+1:] + ":" + strconv.Itoa(line)
        }

        func A() string {
        	return B()
        }

        func main() {
        	fmt.Println(A())
        }
        // Output: prog.go:19 (on Go Playground)
        ```

*   **Regular Expressions - Location and Replacement**
    *   Regular expressions allow matching patterns with variable numbers of characters.
    *   `regexp.MustCompile(pattern)`: Compiles a regex pattern into a `Regexp` object. Panics if the pattern is invalid.
    *   `re.FindAllString(s, n)`: Finds all non-overlapping matches of the regex in `s` and returns them as a slice of strings. `n=-1` finds all.
    *   `re.FindAllStringIndex(s, n)`: Finds all non-overlapping matches and returns their start/end byte indices as a slice of `[]int` (e.g., `[[start1 end1] [start2 end2]]`).
    *   Example using `b+` (one or more 'b's):
        ```go
        package main

        import (
        	"fmt"
        	"regexp"
        )

        func main() {
        	te := "aba abba abbba"
        	re := regexp.MustCompile("b+")
        	mm := re.FindAllString(te, -1)
        	id := re.FindAllStringIndex(te, -1)

        	fmt.Println(mm) // Output: [b bb bbb]
        	fmt.Println(id) // Output: [[1 2] [5 7] [10 13]]

        	for _, d := range id {
        		fmt.Println(te[d[0]:d[1]]) // Output: b, bb, bbb (each on a new line)
        	}
        }
        ```
    *   `re.ReplaceAllStringFunc(s, replFunc)`: Replaces all matches in `s` with the result of applying `replFunc` to each match.
    *   Example using `b+` with `strings.ToUpper`:
        ```go
        package main

        import (
        	"fmt"
        	"regexp"
        	"strings"
        )

        func main() {
        	te := "aba abba abbba"
        	re := regexp.MustCompile("b+")
        	up := re.ReplaceAllStringFunc(te, strings.ToUpper)
        	fmt.Println(up) // Output: aBA aBBa aBBBba
        }
        ```

*   **Regular Expressions - Syntax Details**
    *   **Repetition:**
        *   `.`: Matches any character (except newline).
        *   `.*`: Matches zero or more of any character (greedy).
        *   `.+`: Matches one or more of any character (greedy).
        *   `.?`: Matches zero or one of any character (greedy, prefers one).
        *   `a{n}`: Matches exactly `n` repetitions of 'a'.
        *   `a{n,m}`: Matches `n` to `m` repetitions of 'a'.
    *   **Character Classes:**
        *   `[a-z]`: Matches any character in the range 'a' through 'z'.
        *   `[^a-z]`: Matches any character *not* in the range 'a' through 'z'.
    *   **Location (Anchors):**
        *   `^x`: Matches 'x' at the beginning of the string.
        *   `x$`: Matches 'x' at the end of the string.
        *   `^x$`: Matches 'x' if it's the entire string.
        *   `\b`: Matches a word boundary.
        *   `\bx\b`: Matches the word 'x' by itself.
    *   **Capture Groups:**
        *   `(x)`: Creates a capture group for the pattern 'x'. The matched content can be extracted.

*   **Regular Expressions - Built-in Character Classes**
    *   Go's `regexp` package provides convenient built-in character classes:
        *   `\d`: Any decimal digit (0-9).
        *   `\w`: Any "word" character (alphanumeric + underscore: `[0-9A-Za-z_]`).
        *   `\s`: Any whitespace character.
    *   POSIX character classes (within `[[:...:]]`):
        *   `[[:alpha:]]`: Any alphabetic character.
        *   `[[:alnum:]]`: Any alphanumeric character.
        *   `[[:punct:]]`: Any punctuation character.
        *   `[[:print:]]`: Any printable character.
        *   `[[:xdigit:]]`: Any hexadecimal digit.
    *   These built-in classes are fully Unicode-compatible.

*   **UUID Validation Example**
    *   UUIDs have a specific format (e.g., `072665ee-a034-4cc3-a2e8-9f1822c4ebbb`).
    *   Regex can be used for validation, but it needs to be precise.
    *   A simple regex like `[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}` might miss cases like uppercase hex digits or incorrect version/format bytes as per RFC 4122.
    *   RFC 4122 specifies that certain hex characters (version `V` and format marker `W`) have specific requirements.
    *   A more accurate regex would incorporate `[[:xdigit:]]` for case-insensitivity and specific character sets for `V` (1-5) and `W` (8, 9, a, b).
    *   Using anchors (`^` for start, `$` for end) is crucial to ensure the *entire* string matches the pattern, preventing partial matches.
    *   Example demonstrating UUID validation with a more robust regex and testing various valid/invalid UUIDs:
        ```go
        package main

        import (
        	"fmt"
        	"regexp"
        )

        func main() {
        	// Regex for UUID validation (more robust than simple hex-only)
        	// Matches: xxxxxxxx-xxxx-Vxxx-Wxxx-xxxxxxxxxxxx
        // V is version (1-5)
        // W is format marker (10bb - one of 8, 9, a, b)
        	uu := regexp.MustCompile(`^([[:xdigit:]]{8})-([[:xdigit:]]{4})-([1-5][[:xdigit:]]{3})-([89aAbB][[:xdigit:]]{3})-([[:xdigit:]]{12})$`)

        	test := []string{
        		"072665ee-a034-4cc3-a2e8-9f1822c4ebbb", // Valid
        		"072665ee-a034-6cc3-a2e8-9f1822c4ebbb", // Invalid version (6)
        		"072665ee-a034-4cc3-72e8-9f1822c4ebbb", // Invalid type (7)
        		"072665ee-a034-4cc3-a2e8-9f1822c4ebbbc", // Too long
        		"072665ee-a034-3cc3-82e8-9f1822c4ebbb", // Valid (version 3, type 8)
        	}

        	for i, t := range test {
        		match := uu.MatchString(t)
        		fmt.Printf("%d: %s %t\n", i, t, match)
        	}
        }
        /* Output:
        0: 072665ee-a034-4cc3-a2e8-9f1822c4ebbb true
        1: 072665ee-a034-6cc3-a2e8-9f1822c4ebbb false
        2: 072665ee-a034-4cc3-72e8-9f1822c4ebbb false
        3: 072665ee-a034-4cc3-a2e8-9f1822c4ebbbc false
        4: 072665ee-a034-3cc3-82e8-9f1822c4ebbb true
        */
        ```

*   **Search and Replace with Capture Groups (Advanced)**
    *   Capture groups allow extracting specific parts of a matched string.
    *   `re.FindStringSubmatch(s)`: Returns a slice of strings where the first element is the full match, and subsequent elements are the captured groups.
    *   Using `$1`, `$2`, etc., in the replacement string refers to the content of the corresponding capture groups.
    *   Example reformatting a phone number:
        ```go
        package main

        import (
        	"fmt"
        	"regexp"
        )

        func main() {
        	// Regex to match phone numbers like (214) 514-9548
        	// Capture groups for area code, prefix, and line number
        	phre := regexp.MustCompile(`\(([[:digit:]]{3})\)\s([[:digit:]]{3})-([[:digit:]]{4})`)

        	orig := "call me at (214) 514-3232 today"
        	match := phre.FindStringSubmatch(orig)

        	fmt.Printf("%q\n", match) // Output: ["(214) 514-3232" "214" "514" "3232"]

        	// Reformat to international style: +1 AreaCode-Prefix-LineNumber
        	intl := phre.ReplaceAllString(orig, "+1 $1-$2-$3")
        	fmt.Println(intl) // Output: call me at +1 214-514-3232 today
        }
        ```
    *   **Important Note on Greedy Matching:** By default, quantifiers (`*`, `+`, `?`, `{n,m}`) are "greedy," meaning they try to match as much as possible. This can lead to unexpected behavior if not accounted for.
    *   **`regexp.MustCompile` vs. `regexp.Compile`**: `MustCompile` is used for constant regex patterns that are known to be valid at compile time. It panics on error. `Compile` returns an error and should be used for dynamic or user-provided patterns.
    *   **Non-capturing groups**: `(?:...)` creates a group for grouping parts of the regex without capturing the matched text. This is useful when you need to group for repetition or alternation but don't need to extract that specific part.
    *   **`re.FindAllStringSubmatch(s, n)`**: Finds all matches and their submatches. The `n` parameter limits the number of matches.
    *   **Complexity of URL Matching**: Matching complex patterns like URLs accurately with regex can be very challenging due to various optional components (protocol, port, path, query parameters, fragments) and special characters. It's often better to use dedicated URL parsing libraries if available.
    *   **`u.LiteralPrefix()`**: Returns the literal prefix of the regex that must match at the beginning of the string. This can be useful for optimization or understanding the regex.
    *   **"Clear is better than clever"**: A guiding principle for writing readable and maintainable code, especially with complex tools like regular expressions. Avoid overly clever or concise regex that sacrifices readability.

## What's New
*   The `regexp` package now treats each invalid byte of a UTF-8 string as U+FFFD when processing strings. This may alter the results of functions like `re.FindAllString` or `re.ReplaceAllStringFunc` if the input string contains invalid UTF-8. [3]
*   The `regexp` parser now rejects very deeply nested expressions due to a security fix, causing `regexp.MustCompile` to panic for such patterns. Go 1.19 introduced `syntax.ErrNestingDepth` for this, and Go 1.20 introduced `syntax.ErrLarge` for very large expressions. [4]
*   The `regexp` package now implements the `encoding.TextAppender` interface, providing new ways to append textual representations of regular expressions to byte slices. [9]

## Updated Code Snippets

## Citations
- [1] Go 1.16 Release Notes
- [2] Go 1.17 Release Notes
- [3] Go 1.18 Release Notes
- [4] Go 1.19 Release Notes
- [5] Go 1.20 Release Notes
- [6] Go 1.21 Release Notes
- [7] Go 1.22 Release Notes
- [8] Go 1.23 Release Notes
- [9] Go 1.24 Release Notes