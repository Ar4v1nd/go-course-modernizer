# Go Class: 21 Homework #4

## Summary
This video covers Homework #4, which involves building a simple web server in Go to manage an in-memory database of items and prices. The exercise demonstrates the use of Go's `net/http` package, specifically focusing on how method values can serve as HTTP handlers, and implements basic CRUD (Create, Read, Update, Delete) operations for the database. It also highlights the importance of input validation and proper HTTP error responses. A crucial point is made about the current implementation's lack of concurrency safety, which will be addressed in a future lesson.

## Key Points

### Introduction to HTTP Interfaces in Go
*   Go's `net/http` package defines the `Handler` interface, which requires a single method: `ServeHTTP(ResponseWriter, *Request)`.
*   The `HandlerFunc` type is a function signature that matches the `ServeHTTP` method.
*   A `ServeHTTP` method is defined on the `HandlerFunc` type itself, allowing any function with the matching signature to implicitly satisfy the `Handler` interface. This enables passing plain functions directly as handlers to the HTTP server.
    ```go
    type Handler interface {
        ServeHTTP(w http.ResponseWriter, r *http.Request)
    }

    type HandlerFunc func(w http.ResponseWriter, r *http.Request)

    func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
        f(w, r)
    }
    ```

### Homework Problem: Web Front-End for a Database
*   The goal is to create a web server that acts as a front-end for a simple in-memory database.
*   The database is represented as a `map[string]dollars`, where `dollars` is a custom type wrapping `float32`.
*   The server should support operations to `list` all items, `create` new items, `update` existing item prices, `read` a specific item, and `delete` items.
*   For simplicity, race conditions are ignored in this exercise but are noted as a critical concern for real-world applications.

### Database and Helper Type Definitions
*   A `dollars` type is defined as `float32` with a `String()` method to format currency output.
    ```go
    type dollars float32

    func (d dollars) String() string {
        return fmt.Sprintf("$%.2f", d)
    }
    ```
*   The `database` type is defined as a `map[string]dollars`.
    ```go
    type database map[string]dollars
    ```

### Implementing the `list` Handler
*   The `list` method is defined on the `database` type, taking `http.ResponseWriter` and `*http.Request` as parameters.
*   It iterates through the `db` map and writes each item and its price to the `ResponseWriter`.
*   The `db.list` method value is registered with `http.HandleFunc("/list", db.list)`. Go's type system automatically converts `db.list` (a method value with the correct signature) into an `http.HandlerFunc`, which then satisfies the `http.Handler` interface.
    ```go
    func (db database) list(w http.ResponseWriter, req *http.Request) {
        for item, price := range db {
            fmt.Fprintf(w, "%s: %s\n", item, price)
        }
    }
    ```

### Implementing the `create` Handler
*   The `create` method extracts `item` and `price` from the URL query parameters using `req.URL.Query().Get()`.
*   **Error Handling (Duplicate Item):** If the `item` already exists in the database, it returns an `http.StatusBadRequest` (400) with a "duplicate item" message.
*   **Error Handling (Invalid Price):** It attempts to parse the `price` string into a `float32` using `strconv.ParseFloat`. If parsing fails, it returns an `http.StatusBadRequest` (400) with an "invalid price" message.
*   If successful, the item is added to the map, and a success message is returned.
    ```go
    func (db database) create(w http.ResponseWriter, req *http.Request) {
        item := req.URL.Query().Get("item")
        priceStr := req.URL.Query().Get("price")

        if _, ok := db[item]; ok {
            msg := fmt.Sprintf("duplicate item: %q", item)
            http.Error(w, msg, http.StatusBadRequest) // 400
            return
        }

        price, err := strconv.ParseFloat(priceStr, 32)
        if err != nil {
            msg := fmt.Sprintf("invalid price: %q", priceStr)
            http.Error(w, msg, http.StatusBadRequest) // 400
            return
        }

        db[item] = dollars(price)
        fmt.Fprintf(w, "added %s with price %s\n", item, db[item])
    }
    ```

### Implementing the `update` Handler
*   The `update` method also extracts `item` and `price` from the URL query parameters.
*   **Error Handling (Item Not Found):** If the `item` does not exist in the database, it returns an `http.StatusNotFound` (404) with a "no such item" message.
*   **Error Handling (Invalid Price):** Same price parsing and error handling as the `create` method.
*   If successful, the item's price is updated in the map, and a success message is returned.
    ```go
    func (db database) update(w http.ResponseWriter, req *http.Request) {
        item := req.URL.Query().Get("item")
        priceStr := req.URL.Query().Get("price")

        if _, ok := db[item]; !ok { // Item must exist to update
            msg := fmt.Sprintf("no such item: %q", item)
            http.Error(w, msg, http.StatusNotFound) // 404
            return
        }

        price, err := strconv.ParseFloat(priceStr, 32)
        if err != nil {
            msg := fmt.Sprintf("invalid price: %q", priceStr)
            http.Error(w, msg, http.StatusBadRequest) // 400
            return
        }

        db[item] = dollars(price)
        fmt.Fprintf(w, "new price for %s is %s\n", item, db[item])
    }
    ```

### Implementing the `fetch` (Read Single Item) Handler
*   The `fetch` method extracts the `item` from the URL query parameters.
*   **Error Handling (Item Not Found):** If the `item` is not found in the database, it returns an `http.StatusNotFound` (404).
*   If found, it prints the item and its price.
    ```go
    func (db database) fetch(w http.ResponseWriter, req *http.Request) {
        item := req.URL.Query().Get("item")

        price, ok := db[item]
        if !ok {
            msg := fmt.Sprintf("no such item: %q", item)
            http.Error(w, msg, http.StatusNotFound) // 404
            return
        }
        fmt.Fprintf(w, "%s has price %s\n", item, price)
    }
    ```

### Implementing the `drop` (Delete Item) Handler
*   The `drop` method extracts the `item` from the URL query parameters.
*   **Error Handling (Item Not Found):** If the `item` does not exist in the database, it returns an `http.StatusNotFound` (404).
*   If found, it uses the built-in `delete()` function to remove the item from the map and returns a success message.
    ```go
    func (db database) drop(w http.ResponseWriter, req *http.Request) {
        item := req.URL.Query().Get("item")

        if _, ok := db[item]; !ok { // Item must exist to drop
            msg := fmt.Sprintf("no such item: %q", item)
            http.Error(w, msg, http.StatusNotFound) // 404
            return
        }

        delete(db, item)
        fmt.Fprintf(w, "dropped %s\n", item)
    }
    ```

### Concurrency Considerations
*   The current in-memory `database` (a `map`) is not safe for concurrent access.
*   Go's `net/http` server handles requests concurrently, meaning multiple goroutines could access and modify the `db` map simultaneously.
*   This can lead to race conditions (e.g., data corruption, unexpected behavior).
*   Addressing concurrency issues (e.g., using mutexes or channels) is a separate, advanced topic that will be covered in future lessons.

## What's New
*   The key point "The `create` method extracts `item` and `price` from the URL query parameters using `req.URL.Query().Get()`" is still accurate in its description of the method used. However, the underlying behavior of URL query parsing has changed. As of Go 1.17, the `net/url` and `net/http` packages now reject semicolons (`;`) as a setting separator in URL queries, accepting only ampersands (`&`). Previously, both were accepted. This means that URLs using semicolons for parameter separation will now be parsed differently, potentially leading to `Get()` not finding parameters that it would have found in Go 1.15. [2]

## Citations
*   [1] Go version 1.16
*   [2] Go version 1.17
*   [3] Go version 1.18
*   [4] Go version 1.19
*   [5] Go version 1.20
*   [6] Go version 1.21
*   [7] Go version 1.22
*   [8] Go version 1.23
*   [9] Go version 1.24