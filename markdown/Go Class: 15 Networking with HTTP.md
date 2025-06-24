# Go Class: 15 Networking with HTTP

## Summary
This video provides an introduction to networking with HTTP in Go, demonstrating how to build simple web servers and clients. It covers the extensive support for web services within Go's standard library, including HTTP, JSON, and templating. The lesson highlights Go's concurrent nature for handling web requests and delves into the interface-driven design of its HTTP packages, showcasing how functions and methods interact with `io.Writer` and `io.Reader` interfaces for flexible data handling.

## Key Points

*   **Go's Network Libraries:**
    *   Go is fundamentally designed for building cloud-native software, especially REST-based services.
    *   The Go standard library provides comprehensive packages for web development, including:
        *   Client and server sockets (`net`).
        *   HTTP protocol support (`net/http`).
        *   Route multiplexing.
        *   HTML templating (`html/template`).
        *   JSON and other data formats (`encoding/json`).
        *   Cryptographic security.
        *   SQL database access.
        *   Compression utilities.
        *   Image generation.
    *   Numerous third-party packages offer additional improvements and functionalities.

*   **Basic HTTP Server:**
    *   A minimal Go HTTP server can be created with just a few lines of code using the `net/http` package.
    *   The `http.HandleFunc` function registers a handler for a specific URL path.
    *   The `http.ListenAndServe` function starts the web server on a given address and port.
    *   Go's HTTP server is concurrent by design; each incoming request is handled in its own goroutine.
    *   **Code Example (server.go):**
        ```go
        package main

        import (
            "fmt"
            "log"
            "net/http"
        )

        func handler(w http.ResponseWriter, r *http.Request) {
            // w is the http.ResponseWriter, used to send the response back to the client.
            // r is the *http.Request, containing details about the incoming request.
            fmt.Fprintf(w, "Hello, world! from %s\n", r.URL.Path[1:])
        }

        func main() {
            // Register the handler function for the root path "/"
            http.HandleFunc("/", handler)

            // Start the HTTP server on port 8080.
            // nil uses the default ServeMux (multiplexer) which http.HandleFunc registers with.
            // log.Fatal logs any error and exits the program.
            log.Fatal(http.ListenAndServe(":8080", nil))
        }
        ```
    *   `http.ResponseWriter` implements the `io.Writer` interface, allowing standard Go I/O functions like `fmt.Fprintf` to write directly to the HTTP response stream.
    *   By default, a `200 OK` status is sent if no other status is explicitly set.

*   **Basic HTTP Client:**
    *   Go's `net/http` package also provides simple client functionalities.
    *   The `http.Get` function makes a GET request to a specified URL.
    *   It returns an `*http.Response` object and an error.
    *   It's crucial to `defer resp.Body.Close()` to ensure the response body (and underlying network connection) is closed, preventing resource leaks.
    *   `io/ioutil.ReadAll` can be used to read the entire response body into a byte slice.
    *   **Code Example (client.go):**
        ```go
        package main

        import (
            "fmt"
            "io/ioutil"
            "log"
            "net/http"
            "os"
        )

        func main() {
            // Make a GET request to the local server, appending a command-line argument.
            resp, err := http.Get("http://localhost:8080/" + os.Args[1])
            if err != nil {
                log.Fatal(err)
            }
            // Ensure the response body is closed when the function exits.
            defer resp.Body.Close()

            // Check if the HTTP status code is OK (200).
            if resp.StatusCode == http.StatusOK {
                // Read the entire response body into a byte slice.
                body, err := ioutil.ReadAll(resp.Body)
                if err != nil {
                    log.Fatal(err)
                }
                // Print the response body as a string.
                fmt.Println(string(body))
            } else {
                fmt.Printf("Server returned status: %s\n", resp.Status)
            }
        }
        ```
*   **Go HTTP Handler Design and Interfaces:**
    *   The `net/http` package uses interfaces for flexible design.
    *   The `http.Handler` interface defines the `ServeHTTP` method.
    *   `http.HandlerFunc` is a function type that matches the `ServeHTTP` method's signature.
    *   Go allows methods to be declared on any user-defined type, including function types. This means a function of type `http.HandlerFunc` implicitly satisfies the `http.Handler` interface because it has a `ServeHTTP` method (which simply calls the underlying function).
    *   This enables a clean and idiomatic way to register functions as HTTP handlers.

*   **HTTP Client with JSON and Struct Unmarshalling:**
    *   Go's `encoding/json` package allows easy marshalling (Go to JSON) and unmarshalling (JSON to Go).
    *   To unmarshal JSON into a Go struct, struct fields must be exported (start with an uppercase letter).
    *   JSON tags (`json:"key_name"`) can be used to map JSON keys (often lowercase or snake_case) to Go struct field names (often PascalCase).
    *   **Code Example (struct definition with JSON tags):**
        ```go
        type todo struct {
            UserID    int    `json:"userId"`
            ID        int    `json:"id"`
            Title     string `json:"title"`
            Completed bool   `json:"completed"`
        }
        ```
    *   **Efficient JSON Decoding with `json.NewDecoder`:**
        *   Instead of reading the entire response body into memory with `ioutil.ReadAll` and then unmarshalling, `json.NewDecoder` can read directly from an `io.Reader`.
        *   `resp.Body` (from `http.Get`) implements `io.Reader`.
        *   `json.NewDecoder(resp.Body).Decode(&item)` reads and decodes JSON directly from the stream, which is more memory-efficient for large payloads.
    *   **Code Example (JSON decoding in client):**
        ```go
        // ... imports and todo struct ...

        func main() {
            const url = "https://jsonplaceholder.typicode.com"
            // ... http.Get and error handling ...

            if resp.StatusCode == http.StatusOK {
                var item todo
                // Create a new JSON decoder that reads directly from the response body.
                decoder := json.NewDecoder(resp.Body)
                // Decode the JSON into the 'item' struct.
                err = decoder.Decode(&item)
                if err != nil {
                    log.Fatal(err)
                }
                fmt.Printf("%#v\n", item) // Prints the Go struct
            }
            // ... rest of the code ...
        }
        ```

*   **HTTP Server as a Client with HTML Templating:**
    *   A Go web server can also act as an HTTP client to other services.
    *   The server's handler can fetch data from an external API (e.g., JSONPlaceholder), process it, and then render an HTML response using `html/template`.
    *   `html/template` provides a secure way to generate HTML, preventing common vulnerabilities like HTML injection.
    *   Templates use placeholders (`{{.FieldName}}`) and functions (`{{printf "format" .FieldName}}`) to dynamically insert data.
    *   **Code Example (Server acting as client and rendering HTML):**
        ```go
        package main

        import (
            "encoding/json"
            "fmt"
            "html/template"
            "log"
            "net/http"
            "os" // For os.Exit and os.Stderr
        )

        type todo struct {
            UserID    int    `json:"userId"`
            ID        int    `json:"id"`
            Title     string `json:"title"`
            Completed bool   `json:"completed"`
        }

        // HTML template string
        var form = `
        <h1>Todo #{{.ID}}</h1>
        <div>{{printf "User %d" .UserID}}</div>
        <div>{{printf "%s (completed: %t)" .Title .Completed}}</div>
        `

        func handler(w http.ResponseWriter, r *http.Request) {
            const base = "https://jsonplaceholder.typicode.com"
            // Get the todo ID from the incoming request path (e.g., /todos/1 -> 1)
            todoID := r.URL.Path[1:] // Assuming path is like /todos/1

            // Make a GET request to the external JSON API
            resp, err := http.Get(base + todoID)
            if err != nil {
                http.Error(w, err.Error(), http.StatusServiceUnavailable)
                return // Important to return after sending an error response
            }
            defer resp.Body.Close()

            if resp.StatusCode != http.StatusOK {
                http.Error(w, "External service error", http.StatusInternalServerError)
                return
            }

            var item todo
            // Efficiently decode JSON directly from the response body
            decoder := json.NewDecoder(resp.Body)
            err = decoder.Decode(&item)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }

            // Create and parse the HTML template
            tmpl := template.New("myTemplate") // Give the template a name
            tmpl, err = tmpl.Parse(form)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }

            // Execute the template, writing the formatted HTML to the ResponseWriter
            err = tmpl.Execute(w, item)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
        }

        func main() {
            http.HandleFunc("/todos/", handler) // Handle paths like /todos/1, /todos/2 etc.
            log.Fatal(http.ListenAndServe(":8080", nil))
        }
        ```

## What's New

*   The `io/ioutil` package, including `ioutil.ReadAll`, was deprecated in Go 1.16. Its functionality has been moved to the `io` and `os` packages. New code should use `io.ReadAll` instead for reading response bodies. [1]
*   While `http.HandleFunc` remains valid for registering handlers, Go 1.22 introduced enhanced routing patterns for `net/http.ServeMux`. These enhancements allow handlers to be registered with specific HTTP methods (e.g., "POST /items/create") and support wildcards in URL paths (e.g., `/items/{id}`), providing more expressive and specific routing capabilities. [7]

## Updated Code Snippets

```go
// client.go
package main

import (
	"fmt"
	"io" // Changed from "io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	// Make a GET request to the local server, appending a command-line argument.
	resp, err := http.Get("http://localhost:8080/" + os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	// Ensure the response body is closed when the function exits.
	defer resp.Body.Close()

	// Check if the HTTP status code is OK (200).
	if resp.StatusCode == http.StatusOK {
		// Read the entire response body into a byte slice.
		body, err := io.ReadAll(resp.Body) // Changed from ioutil.ReadAll
		if err != nil {
			log.Fatal(err)
		}
		// Print the response body as a string.
		fmt.Println(string(body))
	} else {
		fmt.Printf("Server returned status: %s\n", resp.Status)
	}
}
```

## Citations
- [1] Go 1.16 Release Notes
- [7] Go 1.22 Release Notes