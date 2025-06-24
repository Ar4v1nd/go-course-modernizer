# Go Class: 28 Conventional Synchronization

## Summary
This video introduces conventional synchronization primitives available in Go's `sync` and `sync/atomic` packages. It explains the problem of mutual exclusion and race conditions in concurrent programming, demonstrating how these primitives (Mutex, RWMutex, atomic operations, Once, Pool, WaitGroup) can be used to safely manage shared resources, contrasting them with Go's higher-level CSP model.

## Key Points

*   **Mutual Exclusion:**
    *   The fundamental problem in concurrent programming is a race condition, where multiple goroutines access and modify shared data concurrently, leading to unpredictable results.
    *   Mutual exclusion ensures that only one goroutine can access a "critical section" (a piece of code that manipulates shared data) at any given instant.
    *   This is typically achieved using locks: acquire the lock before accessing data, and release it when done. Any other goroutine attempting to acquire the lock will block until it's released.

*   **`sync.Mutex` (Mutual Exclusion Lock):**
    *   Provides basic mutual exclusion.
    *   Methods: `Lock()` to acquire the lock, `Unlock()` to release it.
    *   Best practice: Use `defer m.Unlock()` immediately after `m.Lock()` to ensure the lock is always released, even if the function returns early or panics.
    *   A `Mutex` has the same cost for both readers and writers, as it grants exclusive access.
    *   Example:
        ```go
        package main

        import (
        	"fmt"
        	"sync"
        )

        func do() int {
        	var m sync.Mutex // Declare a Mutex
        	var n int64
        	var w sync.WaitGroup

        	for i := 0; i < 1000; i++ {
        		w.Add(1)
        		go func() {
        			m.Lock() // Acquire the lock
        			n++      // Critical section: DATA RACE without lock
        			m.Unlock() // Release the lock
        			w.Done()
        		}()
        	}

        	w.Wait()
        	return int(n)
        }

        func main() {
        	fmt.Println(do()) // Will consistently print 1000
        }
        ```

*   **`sync.RWMutex` (Read-Write Mutex):**
    *   An optimization for scenarios where data is read much more frequently than it's written.
    *   Allows multiple readers to hold a read lock (`RLock()`) concurrently.
    *   A writer's lock (`Lock()`) blocks all readers and other writers.
    *   Readers are preferred over writers (writers might starve if there's a continuous stream of readers).
    *   Example:
        ```go
        type InfoClient struct {
        	mu        sync.RWMutex
        	token     string
        	tokenTime time.Time
        	TTL       time.Duration
        }

        func (i *InfoClient) CheckToken() (string, time.Duration) {
        	i.mu.RLock() // Acquire a read lock
        	defer i.mu.RUnlock() // Release the read lock when done
        	return i.token, i.TTL - time.Since(i.tokenTime)
        }

        func (i *InfoClient) ReplaceToken(ctx context.Context) (string, error) {
        	token, ttl, err := i.getAccessToken(ctx) // Do expensive work outside the lock
        	if err != nil {
        		return "", err
        	}

        	i.mu.Lock() // Acquire a write lock (blocks all readers/writers)
        	defer i.mu.Unlock() // Release the write lock
        	i.token = token
        	i.tokenTime = time.Now()
        	i.TTL = time.Duration(ttl) * time.Second
        	return token, nil
        }
        ```

*   **`sync/atomic` Primitives:**
    *   Provide low-level, hardware-backed atomic operations for scalar values (integers, pointers).
    *   More efficient than mutexes for simple operations like increments, decrements, or swaps, as they often avoid context switching to the OS.
    *   Less flexible than mutexes; only specific operations are supported.
    *   Example:
        ```go
        package main

        import (
        	"fmt"
        	"sync"
        	"sync/atomic" // Import the atomic package
        )

        func do() int {
        	var n int64
        	var w sync.WaitGroup

        	for i := 0; i < 1000; i++ {
        		w.Add(1)
        		go func() {
        			atomic.AddInt64(&n, 1) // Atomically add 1 to n
        			w.Done()
        		}()
        	}

        	w.Wait()
        	return int(n) // Will consistently print 1000
        }

        func main() {
        	fmt.Println(do())
        }
        ```

*   **`sync.Once` (Only-Once Execution):**
    *   Guarantees that a function will be executed exactly once, even if called concurrently by multiple goroutines.
    *   Useful for initializing singletons or performing one-time setup tasks.
    *   Example:
        ```go
        package main

        import (
        	"fmt"
        	"sync"
        	"time"
        )

        var once sync.Once
        var x *Singleton // Our singleton instance

        type Singleton struct {
        	// ... fields
        }

        func NewSingleton() *Singleton {
        	fmt.Println("Initializing Singleton...")
        	time.Sleep(100 * time.Millisecond) // Simulate expensive initialization
        	return &Singleton{}
        }

        func initialize() {
        	x = NewSingleton()
        }

        func main() {
        	var wg sync.WaitGroup
        	for i := 0; i < 5; i++ {
        		wg.Add(1)
        		go func(id int) {
        			defer wg.Done()
        			fmt.Printf("Goroutine %d calling once.Do...\n", id)
        			once.Do(initialize) // initialize() will be called only once
        			fmt.Printf("Goroutine %d finished once.Do. Singleton: %p\n", id, x)
        		}(i)
        	}
        	wg.Wait()
        	fmt.Println("All goroutines finished.")
        }
        ```

*   **`sync.Pool` (Object Pooling):**
    *   Provides a way to efficiently reuse frequently allocated, temporary objects.
    *   Reduces garbage collection pressure and allocation overhead.
    *   Objects are stored as `interface{}`, requiring type assertion when retrieved.
    *   Example:
        ```go
        package main

        import (
        	"bytes"
        	"fmt"
        	"io"
        	"sync"
        )

        var bufPool = sync.Pool{
        	New: func() interface{} {
        		return new(bytes.Buffer) // Function to create new objects if pool is empty
        	},
        }

        func Log(w io.Writer, key, val string) {
        	b := bufPool.Get().(*bytes.Buffer) // Get a buffer from the pool, type assert
        	b.Reset()                           // Reset the buffer for reuse
        	b.WriteString(key)
        	b.WriteString(": ")
        	b.WriteString(val)
        	w.Write(b.Bytes())
        	bufPool.Put(b) // Return the buffer to the pool
        }

        func main() {
        	var wg sync.WaitGroup
        	for i := 0; i < 10; i++ {
        		wg.Add(1)
        		go func(id int) {
        			defer wg.Done()
        			Log(bytes.NewBuffer(nil), fmt.Sprintf("goroutine-%d", id), "some log message")
        		}(i)
        	}
        	wg.Wait()
        	fmt.Println("Logging complete.")
        }
        ```

*   **Other Primitives:**
    *   `sync.Cond`: A condition variable used with a `Mutex` to allow goroutines to wait for a specific condition to be met before proceeding.
    *   `sync.Map`: A concurrent map implementation that is safe for concurrent use by multiple goroutines without explicit locking. It stores `interface{}` values.
    *   `sync.WaitGroup`: Used to wait for a collection of goroutines to finish. (Covered in previous sections).

## What's New
*   The `sync/atomic` package now defines new atomic types `Bool`, `Int32`, `Int64`, `Uint32`, `Uint64`, `Uintptr`, and `Pointer` to make it easier to use atomic values. [4]
*   The `sync` package now includes `OnceFunc`, `OnceValue`, and `OnceValues` functions, which provide a more convenient way to lazily initialize a value on first use. [6]
*   The implementation of `sync.Map` has been improved for performance, especially for map modifications, and it now includes a `Map.Clear` method to delete all entries. [8]

## Updated Code Snippets
(No updated code snippets are needed.)

## Citations
*   [4] Go version 1.19
*   [6] Go version 1.21
*   [8] Go version 1.23