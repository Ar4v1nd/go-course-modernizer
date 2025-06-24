# Go Class: 29 Homework #5 (h/w #4 part deux)

## Summary
This video revisits Homework #4, which involved building a simple web server with a fake database. The primary goal of Homework #5 is to identify and fix a concurrency issue (race condition) present in the Homework #4 solution. The video demonstrates how to use Go's built-in race detector to find these issues and then how to resolve them using mutexes.

## Key Points

*   **Revisiting Homework #4's Web Server**
    *   Homework #4 involved creating a REST-based web server that simulated a storefront with items and prices stored in a simple in-memory map.
    *   The previous homework intentionally left a concurrency issue unaddressed.

*   **Creating a Test Driver for Concurrency Issues**
    *   A separate Go program is written to act as a test driver, generating traffic (create, update, delete requests) against the web server.
    *   The driver program uses goroutines to send multiple concurrent requests, increasing the likelihood of exposing race conditions.
    *   The `doQuery` function simplifies sending HTTP GET requests to the server's API endpoints.

    ```go
    func doQuery(cmd, parms string) error {
        resp, err := http.Get("http://localhost:8080/" + cmd + "?" + parms)
        if err != nil {
            fmt.Fprintf(os.Stderr, "err %s = %v\n", parms, err)
            return err
        }
        defer resp.Body.Close()
        fmt.Fprintf(os.Stderr, "got %s = %d (no err)\n", parms, resp.StatusCode)
        return nil
    }

    func runAdds() {
        for _, s := range items {
            if err := doQuery("create", "item="+s.item+"&price="+s.price); err != nil {
                return
            }
        }
    }

    func main() {
        go runAdds()
        go runDeletes()
        go runUpdates()
        time.Sleep(5 * time.Second) // Allow goroutines to run
    }
    ```

*   **Detecting Race Conditions with Go's Race Detector**
    *   Go provides a built-in race detector that can identify concurrent access to shared memory without proper synchronization.
    *   To enable the race detector, compile and run the server program with the `-race` flag: `go run -race .` (from the server's directory).
    *   When the test driver is run against the race-detected server, the race detector outputs warnings, indicating where concurrent reads and writes to the shared map (`db`) are occurring. This confirms the existence of a race condition.

    ```bash
    # Example output from race detector
    WARNING: DATA RACE
    Read at 0x00c00007cc00 by goroutine 8:
      runtime.mapaccess2_faststr()
      /usr/local/Cellar/go/1.14.5/libexec/src/runtime/map_faststr.go:107 +0x0
      main.database.drop()
      /Users/mholiday/Projects/go-material/Go Training/cmd/hw4/main.go:83 +0xac
      net/http.HandlerFunc.ServeHTTP()
      ...
    Previous read at 0x00c00007cc00 by goroutine 7:
      runtime.mapaccess2_faststr()
      /usr/local/Cellar/go/1.14.5/libexec/src/runtime/map_faststr.go:107 +0x0
      main.database.add()
      /Users/mholiday/Projects/go-material/Go Training/cmd/hw4/main.go:28 +0x187
      net/http.HandlerFunc.ServeHTTP()
      ...
    ```

*   **Solving Race Conditions with Mutexes**
    *   To make the database operations thread-safe (goroutine-safe), a `sync.Mutex` is added to the `database` struct.
    *   All methods that access or modify the underlying map (`db`) must acquire a lock on this mutex before accessing the map and release it afterward.
    *   The `defer` keyword is used with `Unlock()` to ensure the mutex is always released, even if the function returns early due to an error.
    *   Crucially, the receiver for all database methods must be a *pointer* to the `database` struct (`*database`) to ensure all goroutines operate on the *same* mutex instance. If it were a value receiver, each method call would receive a copy of the mutex, defeating the purpose of synchronization.

    ```go
    import "sync" // Import the sync package

    type database struct {
        mu sync.Mutex // Add a mutex field
        db map[string]dollars
    }

    // Example for a method that reads from the database
    func (d *database) list(w http.ResponseWriter, req *http.Request) {
        d.mu.Lock() // Acquire lock
        defer d.mu.Unlock() // Release lock when function exits

        for item, price := range d.db { // Access the map via d.db
            fmt.Fprintf(w, "%s: %s\n", item, price)
        }
    }

    // Example for a method that writes to the database
    func (d *database) add(w http.ResponseWriter, req *http.Request) {
        d.mu.Lock() // Acquire lock
        defer d.mu.Unlock() // Release lock when function exits

        item := req.URL.Query().Get("item")
        price := req.URL.Query().Get("price")

        if _, ok := d.db[item]; ok { // Access the map via d.db
            msg := fmt.Sprintf("duplicate item: %q", item)
            http.Error(w, msg, http.StatusBadRequest) // 400
            return
        }
        // ... rest of the add logic ...
        d.db[item] = dollars(p) // Write to the map via d.db
        // ...
    }

    // Initialize the database in main
    func main() {
        db := database{
            db: map[string]dollars{
                "shoes": 50,
                "socks": 5,
            },
        }
        // ... http.HandleFunc calls now pass &db (pointer) ...
    }
    ```

*   **Verification of the Fix**
    *   After applying the mutexes and pointer receivers, running the server with `-race` and then the test driver program no longer produces `DATA RACE` warnings.
    *   The server also runs stably without crashing, demonstrating that the concurrency issues have been successfully resolved.

## What's New
No changes were found that invalidate or make inaccurate the provided key points or code snippets. The core concepts and usage of the Go race detector and `sync.Mutex` remain valid in Go 1.24.

## Updated Code Snippets
No updated code snippets are needed as the original code remains valid.

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