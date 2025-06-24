# Go Class: 11 Homework #2

## Summary
This video focuses on implementing Exercise 5.5 from "The Go Programming Language" book, which involves parsing HTML to count words and images. The instructor demonstrates how to use Go's standard library and an extended library package (`golang.org/x/net/html`) to parse an HTML document represented as a raw string. The core of the solution involves a recursive, depth-first tree traversal function to navigate the HTML document's node structure and count relevant elements.

## Key Points

- **Exercise Objective:** Implement a Go program to count words and images within an HTML document.
- **HTML Source:** Instead of fetching HTML from an HTTP query (which is covered later in the course), the HTML content is provided as a raw string literal within the Go program.
    - Raw strings in Go are enclosed in backticks (`` ` ``) and can span multiple lines, allowing for direct embedding of multi-line text without needing to escape newlines or internal quotation marks.
    ```go
    var raw = `
    <!DOCTYPE html>
    <html>
    <body>
    <h1>My First Heading</h1>
    <p>My first paragraph.</p>
    <p>HTML images are defined with the img tag:</p>
    <img src="xxx.jpg" width="104" height="142">
    </body>
    </html>
    `
    ```
- **HTML Parsing with `golang.org/x/net/html`:**
    - The `html.Parse` function from the `golang.org/x/net/html` package is used to parse the HTML string into a document tree structure.
    - This package is part of Go's extended standard library, indicated by `golang.org/x/`.
    - `html.Parse` expects an `io.Reader` as input. A `bytes.NewReader` is used to convert the raw string into a reader.
    ```go
    import (
        "bytes"
        "fmt"
        "os"
        "strings"
        "golang.org/x/net/html"
    )

    func main() {
        // ... var raw = `...`
        doc, err := html.Parse(bytes.NewReader([]byte(raw)))
        if err != nil {
            fmt.Fprintf(os.Stderr, "parse failed: %s\n", err)
            os.Exit(-1)
        }
        // ...
    }
    ```
- **Counting Words and Images (`countWordsAndImages` function):**
    - A top-level function `countWordsAndImages` is created to encapsulate the counting logic.
    - It takes the parsed HTML document (a pointer to `html.Node`) as input and returns two integers: the total word count and the total image count.
    - This function initializes local variables for `words` and `pics` and then calls a recursive helper function `visit` to traverse the document tree.
    ```go
    func countWordsAndImages(doc *html.Node) (int, int) {
        var words, pics int
        visit(doc, &words, &pics) // Pass pointers to modify counts
        return words, pics
    }
    ```
- **Recursive Tree Traversal (`visit` function):**
    - The `visit` function performs a depth-first traversal of the HTML document tree.
    - It takes the current `html.Node` and pointers to the `words` and `pics` accumulators.
    - **Node Type Checking:**
        - It checks the `n.Type` property of the current node to determine if it's a text node or an element node.
        - If `html.TextNode`: Words are counted by splitting the `n.Data` (text content) using `strings.Fields` and adding the length of the resulting slice to `*pwords`.
        - If `html.ElementNode` and `n.Data == "img"`: The image count `*ppics` is incremented.
    ```go
    func visit(n *html.Node, pwords, ppics *int) {
        // Check if it's a text node and count words
        if n.Type == html.TextNode {
            *pwords += len(strings.Fields(n.Data))
        } else if n.Type == html.ElementNode && n.Data == "img" {
            // Check if it's an image element and count images
            *ppics++
        }

        // Recursively visit children and siblings
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            visit(c, pwords, ppics)
        }
    }
    ```
    - **Traversal Logic:**
        - `n.FirstChild`: Points to the first child of the current node.
        - `c.NextSibling`: Iterates through the siblings of the current child.
        - The recursive call `visit(c, pwords, ppics)` ensures depth-first exploration.
- **Pointer Usage:** Pointers (`*int`) are used for `pwords` and `ppics` in the `visit` function. This allows the recursive calls to modify the same `words` and `pics` variables declared in `countWordsAndImages,` effectively accumulating the counts across the entire tree traversal.
- **Output:** The final counts are printed to the console.
    ```go
    // ... in main()
    words, pics := countWordsAndImages(doc)
    fmt.Printf("%d words and %d images\n", words, pics)
    ```

## What's New
All key points remain valid and accurate in Go version 1.24 based on the provided release notes.

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