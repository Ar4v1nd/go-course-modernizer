# Go Class: 02 Simple Example

## Summary
This video introduces a simple Go program that processes command-line arguments, demonstrating basic Go syntax, project structuring with packages, and the fundamental concepts of unit testing using Go's built-in `testing` package. It also highlights the shift from `GOPATH` to Go modules for dependency management.

## Key Points

*   **Command-Line Arguments:**
    *   Go's `main` function does not directly accept command-line arguments as parameters.
    *   Command-line arguments are accessed via the `os.Args` slice of strings from the `os` package.
    *   `os.Args[0]` contains the name of the executable program.
    *   Subsequent elements (`os.Args[1:]`) contain the actual arguments passed by the user.
    *   Slicing (e.g., `os.Args[1:]`) is used to get a sub-slice of arguments, excluding the program name. An empty slice is returned if no arguments are provided, preventing out-of-bounds errors.

*   **Project Structure and Packages:**
    *   Go programs are organized into packages. An executable program must be in `package main`.
    *   It's a common convention to place the main executable in a `cmd` subdirectory (e.g., `cmd/myprogram/main.go`).
    *   Reusable code (libraries) are placed in separate packages (e.g., `package hello` in `hello.go`).
    *   Packages are imported using their module path (e.g., `import "your_module_name/hello"`).

*   **Functions:**
    *   Function declaration syntax: `func FunctionName(parameterName parameterType) returnType { ... }`.
    *   In Go, the type of a parameter or variable is declared *after* its name (e.g., `name string`).
    *   Functions can accept slices as parameters, allowing for a variable number of inputs (e.g., `names []string`).
    *   The `strings.Join` function from the `strings` package can concatenate elements of a string slice with a specified separator.

*   **Unit Testing:**
    *   Go has built-in support for unit testing. Test files must end with `_test.go` (e.g., `hello_test.go`).
    *   Test functions must start with `Test` (capital T) followed by a descriptive name (e.g., `func TestSayHello`).
    *   Test functions take a single parameter: `t *testing.T`. This `t` object provides methods for reporting test status (e.g., `t.Errorf` for failures).
    *   Table-driven tests are a common Go idiom for testing multiple scenarios. A slice of anonymous structs is created, where each struct defines inputs and expected outputs for a specific test case. A `for...range` loop iterates through these test cases.
    *   Tests are run from the terminal using the `go test` command (e.g., `go test ./...` to run all tests in the current module, or `go test ./hello` to test a specific package).
    *   Go's compiler is strict: it will not compile if there are unused imported packages.

*   **Go Modules:**
    *   Go modules are the modern way to manage dependencies in Go, replacing the older `GOPATH` system.
    *   A `go.mod` file defines the module's path and tracks its dependencies.
    *   `go mod init <module_name>` initializes a new module and creates the `go.mod` file.
    *   The `go.mod` file specifies the module name and the Go version it targets (e.g., `module hello`, `go 1.14`).
    *   Go automatically manages third-party dependencies listed in `go.mod`, downloading and caching them as needed.

**Example Code Snippets:**

**`cmd/main.go`**
```go
package main

import (
	"fmt"
	"os"
	"your_module_name/hello" // Replace 'your_module_name' with your actual module name
)

func main() {
	var names []string
	if len(os.Args) > 1 {
		names = os.Args[1:] // Get all arguments after the program name
	} else {
		names = []string{"world"} // Default to "world" if no arguments
	}
	fmt.Println(hello.Say(names))
}
```

**`hello.go`**
```go
package hello

import (
	"fmt"
	"strings" // Import strings package for strings.Join
)

// Say returns a greeting string based on the provided names.
// If no names are provided, it defaults to "world".
func Say(names []string) string {
	if len(names) == 0 {
		names = []string{"world"} // Default to "world" if no names
	}
	return fmt.Sprintf("Hello, %s!", strings.Join(names, ", "))
}
```

**`hello_test.go`**
```go
package hello

import (
	"testing"
)

func TestSayHello(t *testing.T) {
	// Define a struct for test cases
	type testCase struct {
		items  []string // Input names
		result string   // Expected result
	}

	// Create a slice of test cases
	subtests := []testCase{
		{items: []string{}, result: "Hello, world!"},
		{items: []string{"test"}, result: "Hello, test!"},
		{items: []string{"Matt"}, result: "Hello, Matt!"},
		{items: []string{"Matt", "Anne"}, result: "Hello, Matt, Anne!"},
		{items: []string{"Matt", "Anne", "Dorothy"}, result: "Hello, Matt, Anne, Dorothy!"},
	}

	// Iterate over test cases
	for _, st := range subtests {
		got := Say(st.items) // Call the function with test input
		if got != st.result { // Compare actual result with expected result
			t.Errorf("wanted %q, got %q for items %v", st.result, got, st.items)
		}
	}
}
```

## What's New

*   **Unit Testing:** The behavior of `for...range` loops, commonly used in table-driven tests, has changed. In Go 1.22, variables declared by a `for` loop are created anew for each iteration, preventing accidental sharing bugs that could occur if closures captured the loop variable from previous iterations. [7]
*   **Go Modules:**
    *   The `GO111MODULE` environment variable now defaults to `on`, meaning module-aware mode is enabled by default, regardless of the presence of a `go.mod` file in the current or parent directory. [1]
    *   The `go get` command's role has shifted. While Go still automatically manages dependencies, `go get` no longer builds or installs packages in module-aware mode. Its primary function is now to adjust dependencies in `go.mod`. To install executables, `go install <module_path>@<version>` is the recommended command. [1]

## Citations
- [1] Go 1.16 Release Notes
- [7] Go 1.22 Release Notes