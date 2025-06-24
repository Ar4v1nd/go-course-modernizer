# Go Class: 01 Hello world!

## Summary
This video introduces the basics of Go programming, starting with the classic "Hello, world!" program. It covers online environments like the Go Playground and repl.it for writing and running Go code, discusses their capabilities and limitations, and provides a brief overview of how to install and run Go programs locally from the command line.

## Key Points

*   **Go Playground**
    *   Simple Go programs can be run directly in a web browser using the Go Playground at `https://play.golang.org`.
    *   It's ideal for quick tests and learning basic syntax.
    *   **Limitations**:
        *   No input or output except writing to `stdout` and `stderr`.
        *   Cannot write to files, open network sockets, or run web servers due to security restrictions.
    *   The Go documentation at `https://golang.org/doc/` includes runnable examples that utilize the Go Playground.

*   **Simplest Go Program Structure**
    *   Every executable Go program must belong to the `main` package.
    *   The program's execution begins in the `main` function.
    *   Packages used in the program must be imported using the `import` keyword.
    *   Functions from imported packages are called using the `package.Function()` notation.

    ```go
    package main

    import (
    	"fmt"
    )

    func main() {
    	fmt.Println("Hello, world!")
    }
    ```

*   **repl.it**
    *   An alternative online environment for Go programming available at `https://replit.it`.
    *   **Advantage**: Supports full input/output operations, including reading/writing files and network interactions, making it suitable for more complex programs or training.

*   **Go Installation**
    *   Detailed installation instructions for various operating systems are available on the official Go downloads page: `https://golang.org/dl`.
    *   **Mac**: Use Homebrew (`brew install go`) or the provided installer package.
    *   **Windows**: Use the MSI installer file and follow the prompts.
    *   **Linux**: Download the archive (e.g., `go1.15.6.linux-amd64.tar.gz`), extract it to `/usr/local` (creating `/usr/local/go`), and remember to add `/usr/local/go/bin` to your `PATH` environment variable.
        ```bash
        $ sudo tar -C /usr/local -xzf go1.15.6.linux-amd64.tar.gz
        ```

*   **Running a Program from the Command Line**
    *   Once Go is installed locally, navigate to the directory containing your Go source file (e.g., `main.go`).
    *   Use the `go run` command followed by a dot (`.`) to compile and execute the program in the current directory. This creates a temporary binary that is removed after execution.

    ```bash
    $ go run .
    Hello, world!
    ```

## What's New
*   The statement "This creates a temporary binary that is removed after execution" for `go run` is no longer entirely accurate. While `go run` still creates a temporary binary, executables created by `go run` are now cached in the Go build cache. This makes repeated executions faster at the expense of making the cache larger, meaning the binary might persist in the cache rather than being immediately removed from the system. [9]

## Citations
- [9] Go version 1.24