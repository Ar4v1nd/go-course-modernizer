# Go Class: 39 Code Coverage

## Summary
This video introduces Go's built-in code coverage tools, demonstrating how to measure the percentage of code exercised by unit tests, generate detailed coverage profiles, and visualize these profiles as an interactive HTML report. It also highlights how to use code coverage to identify and improve test cases, particularly for error handling paths.

## Key Points

*   **Introduction to Code Coverage**:
    *   Code coverage tools help identify which parts of your code are executed when running unit tests.
    *   This information is crucial for understanding the effectiveness of your test suite and pinpointing areas that lack testing.

*   **Example Code Overview**:
    *   The demonstration uses a Go program (`main.go`) with functions like `matchNum`, `matchString`, `contains`, and `CheckData`.
    *   The `contains` function recursively checks if a smaller JSON object (represented as `exp` map) is contained within a larger JSON object (`data` map).
    *   `CheckData` unmarshals JSON byte slices into maps and then calls `contains`.
    *   Unit tests (`main_test.go`) are provided for `TestContains` (positive cases) and `TestNotContains` (negative cases).

*   **Measuring Basic Code Coverage**:
    *   To get a simple percentage of statements covered by your tests, use the `go test` command with the `-cover` flag.
    ```bash
    go test ./... -cover
    ```
    *   This command will output the percentage of statements covered, e.g., `coverage: 85.2% of statements`.

*   **Generating a Detailed Coverage Profile**:
    *   To get more granular information about which specific lines and blocks of code are covered, use the `-coverprofile` flag to output the coverage data to a file.
    *   The `-covermode=count` flag (optional but recommended) tracks how many times each statement is executed, providing a "heat map" effect in the visual report.
    ```bash
    go test ./... -coverprofile=c.out -covermode=count
    ```
    *   This creates a file (e.g., `c.out`) containing the coverage profile. The raw content of this file is not easily human-readable.

*   **Visualizing Code Coverage**:
    *   Go provides a powerful tool to visualize the coverage profile in an HTML format, which can be opened in a web browser.
    ```bash
    go tool cover -html=c.out
    ```
    *   The HTML report displays your source code with color-coded lines:
        *   **Red**: Code that was not executed by any test.
        *   **Green**: Code that was executed by tests. Darker shades of green (or higher counts if `covermode=count` is used) indicate more frequent execution.
        *   **Grey**: Code that is not tracked for coverage (e.g., comments, blank lines, or certain declarations).
    *   This visual representation makes it easy to spot untested code paths.

*   **Improving Code Coverage**:
    *   By examining the red (uncovered) sections in the HTML report, you can identify missing test cases.
    *   For example, if an `else` branch or an error return path is red, you need to craft a test case that specifically triggers that condition.
    *   **Example**: To cover the `else` branch where a key exists but its type is incorrect, or where a key is missing, add specific JSON structures to your `TestNotContains` table.
    ```go
    // In main_test.go, inside TestNotContains
    var known = []string{
        // ... existing test cases ...
        `{"city": {}}`, // Triggers "missing in data" error if "city" is not an object in unknown
        `{"name": {}}`, // Triggers "wrong in data" error if "name" is not an object in unknown
    }
    ```
    *   After adding new tests, rerun the `go test` and `go tool cover` commands to see the improved coverage.

*   **Code Coverage Best Practices**:
    *   Code coverage is a **tool to improve testing**, not a goal in itself.
    *   Aiming for 100% coverage might lead to writing complex or artificial tests for obscure error cases that provide little real value.
    *   Focus on covering critical paths and common error scenarios. Use the visual report to guide your testing efforts, prioritizing areas that are genuinely important for the correctness and robustness of your application.

## What's New
*   **Go 1.20: Program Coverage**: Go 1.20 introduced support for collecting code coverage profiles for programs (applications and integration tests), as opposed to just unit tests. This is achieved by building the program with `go build -cover` and then running the resulting binary with the `GOCOVERDIR` environment variable set to an output directory for coverage profiles. [5, Page 4, Cover]
*   **Go 1.22: Coverage Reporting for Packages without Test Files**: The `go test -cover` command now prints coverage summaries for covered packages that do not have their own test files, reporting `0.0% of statements` instead of `[no test files]`. [7, Page 2, Go command]

## Citations
*   [5] Go version 1.20
*   [7] Go version 1.22