# Go Class: 33 Reflection

## Summary
This video provides a comprehensive introduction to reflection in Go, focusing on its practical applications. It covers type assertion (downcasting), the concept of runtime type information, and how the `reflect` package enables dynamic operations like deep equality checks, type-based switching, and custom JSON decoding. The instructor demonstrates these concepts with clear code examples, highlighting how reflection allows programs to inspect and manipulate their own structure and data at runtime.

## Key Points

*   **Type Assertion (a.k.a. Downcasting)**
    *   The empty interface (`interface{}`) says nothing as it has no methods, allowing it to hold any type in Go.
    *   To extract the "real" underlying type from an interface, a type assertion is used.
    *   **Single-result form:** `value.(T)`
        *   If the underlying type of `value` is `T`, the assertion succeeds, and `value` is converted to type `T`.
        *   If the underlying type is not `T`, the program will `panic`.
        ```go
        var w io.Writer = os.Stdout
        f := w.(*os.File) // Success: f is now *os.File
        // c := w.(*bytes.Buffer) // Panics at runtime
        ```
    *   **Two-result form (safe assertion):** `value, ok := w.(T)`
        *   `ok` is a boolean indicating whether the assertion succeeded.
        *   If `ok` is `false`, `value` will be the zero value for type `T` (e.g., `nil` for pointers/interfaces), and no panic occurs.
        ```go
        var w io.Writer = os.Stdout
        f, ok := w.(*os.File) // ok is true, f is *os.File
        if ok {
            // Use f as *os.File
        }

        b, ok := w.(*bytes.Buffer) // ok is false, b is nil
        if !ok {
            // Handle the failure gracefully
        }
        ```

*   **What is Reflection?**
    *   Reflection is the ability of a program to examine and modify its own structure and behavior at runtime.
    *   Unlike some older languages (e.g., C) where type information is discarded after compilation, Go compilers embed type information directly into the executable binary.
    *   This embedded type information allows Go programs to query types, inspect fields of structs, call methods, and perform other dynamic operations.

*   **Deep Equality**
    *   The `reflect` package provides utilities for runtime type inspection.
    *   `reflect.DeepEqual(x, y)` can be used to check if two values are deeply equal, even if they contain non-comparable types like slices or maps.
    *   This is particularly useful in unit tests where direct `==` comparison might not work for complex structs.
    ```go
    import "reflect"

    type myStruct struct {
        A string
        B []int // Slices are not comparable with ==
    }

    want := myStruct{A: "hello", B: []int{1, 2, 3}}
    got := myStruct{A: "hello", B: []int{1, 2, 3}}

    if !reflect.DeepEqual(got, want) {
        // Handle mismatch, e.g., t.Errorf("Mismatch: got %#v, want %#v", got, want)
    }
    ```

*   **Switching on Type (Type Switch)**
    *   Go's `switch` statement can operate on the *type* of an interface value.
    *   **Syntax:** `switch v := i.(type)` where `i` is an interface.
    *   Inside each `case` block, the variable `v` will have the specific type declared in that case, allowing type-specific operations without further assertions.
    *   This is commonly used in functions that accept `interface{}` to handle different data types gracefully (e.g., `fmt.Println`).
    ```go
    import "fmt"

    type Stringer interface {
        String() string
    }

    func PrintAnything(arg interface{}) {
        switch v := arg.(type) {
        case string:
            fmt.Printf("It's a string: %s\n", v) // v is of type string here
        case int:
            fmt.Printf("It's an int: %d\n", v) // v is of type int here
        case Stringer:
            fmt.Printf("It's a Stringer: %s\n", v.String()) // v is of type Stringer here
        default:
            fmt.Printf("Unknown type: %T\n", v)
        }
    }
    ```

*   **Custom JSON Decoding with Reflection**
    *   For complex JSON structures where keys might depend on other values (e.g., a "type" field determines the structure of a subsequent field), custom unmarshaling is often required.
    *   By implementing the `UnmarshalJSON([]byte) error` method on a type, you can control how JSON is decoded.
    *   A common pattern is to unmarshal the initial JSON into a generic `map[string]interface{}`.
    *   Then, use type assertions and type switches on values within this map to dynamically extract and assign data to specific fields of your Go struct.
    *   A wrapper struct can be used to prevent infinite recursion if the custom `UnmarshalJSON` method needs to call the standard `json.Unmarshal` on a portion of the data that matches the original type.
    ```go
    import (
        "encoding/json"
        "fmt"
        "strings"
    )

    type Response struct {
        Item   string
        Album  string
        Title  string
        Artist string
    }

    // respWrapper is needed to avoid infinite recursion when calling json.Unmarshal
    type respWrapper struct {
        Response
    }

    // UnmarshalJSON implements the json.Unmarshaler interface for respWrapper
    func (r *respWrapper) UnmarshalJSON(b []byte) error {
        // First, unmarshal into a generic map to inspect the "item" field
        var rawMap map[string]interface{}
        if err := json.Unmarshal(b, &rawMap); err != nil {
            return err
        }

        // Extract the "item" field
        item, ok := rawMap["item"].(string)
        if !ok {
            return fmt.Errorf("item field not found or not a string")
        }
        r.Item = item // Set the Item field in the actual Response struct

        // Use a type switch on the item to determine how to unmarshal the rest
        switch item {
        case "album":
            if albumData, ok := rawMap["album"].(map[string]interface{}); ok {
                if title, ok := albumData["title"].(string); ok {
                    r.Album = title // Assuming Album field stores the title for albums
                }
            }
        case "song":
            if songData, ok := rawMap["song"].(map[string]interface{}); ok {
                if title, ok := songData["title"].(string); ok {
                    r.Title = title
                }
                if artist, ok := songData["artist"].(string); ok {
                    r.Artist = artist
                }
            }
        default:
            // Handle unknown item types or ignore
        }
        return nil
    }

    // Example usage (from video's playground)
    // var j1 = []byte(`{"item": "album", "album": {"title": "Dark Side of the Moon"}}`)
    // var j2 = []byte(`{"item": "song", "song": {"title": "Bella Donna", "artist": "Stevie Nicks"}}`)

    // func main() {
    //     var resp1 respWrapper
    //     if err := json.Unmarshal(j1, &resp1); err != nil {
    //         log.Fatal(err)
    //     }
    //     fmt.Printf("Album Response: %#v\n", resp1.Response)

    //     var resp2 respWrapper
    //     if err := json.Unmarshal(j2, &resp2); err != nil {
    //         log.Fatal(err)
    //     }
    //     fmt.Printf("Song Response: %#v\n", resp2.Response)
    // }
    ```

*   **Testing JSON Fragments (Checking for Sub-fragments with Reflection)**
    *   Reflection can be used to check if a smaller, "known" JSON fragment is contained within a larger, "unknown" JSON piece.
    *   This involves unmarshaling both the known and unknown JSON into `map[string]interface{}`.
    *   A recursive `contains` function can then iterate through the known map, using type switches and assertions to compare values and recursively call itself for nested maps.
    *   Helper functions like `matchNum` and `matchString` simplify the comparison of primitive types, handling type assertions and value comparisons.
    ```go
    // Helper to match a numeric key-value pair
    func matchNum(key string, exp float64, data map[string]interface{}) bool {
        if v, ok := data[key]; ok {
            if val, ok := v.(float64); ok && val == exp { // Assert to float64 as JSON numbers decode to float64
                return true
            }
        }
        return false
    }

    // Helper to match a string key-value pair
    func matchString(key string, exp string, data map[string]interface{}) bool {
        if v, ok := data[key]; ok {
            if val, ok := v.(string); ok && strings.EqualFold(val, exp) { // Case-insensitive string comparison
                return true
            }
        }
        return false
    }

    // Contains checks if all key-value pairs in 'exp' are present and match in 'got'
    func contains(exp, got map[string]interface{}) error {
        for k, v := range exp {
            // Use a type switch on the expected value's type
            switch x := v.(type) {
            case float64:
                if !matchNum(k, x, got) {
                    return fmt.Errorf("numeric mismatch for key '%s'", k)
                }
            case string:
                if !matchString(k, x, got) {
                    return fmt.Errorf("string mismatch for key '%s'", k)
                }
            case map[string]interface{}: // Recursive case for nested objects
                if gotVal, ok := got[k].(map[string]interface{}); ok {
                    if err := contains(x, gotVal); err != nil { // Recursively check nested map
                        return err
                    }
                } else {
                    return fmt.Errorf("missing or wrong type for nested object key '%s'", k)
                }
            // Add cases for other types (bool, slice, etc.) as needed
            default:
                return fmt.Errorf("unhandled type %T for key '%s'", x, k)
            }
        }
        return nil // All expected fragments found
    }

    // Example usage (from video's playground)
    // func TestContains(t *testing.T) {
    //     var unknown = []byte(`{"id": 1, "name": "bob", "addr": {"street": "Lazy Lane", "city": "Exit", "zip": "99999"}, "extra": 21.1}`)
    //     var known = []string{
    //         `{"id": 1}`,
    //         `{"extra": 21.1}`,
    //         `{"name": "bob"}`,
    //         `{"addr": {"street": "Lazy Lane", "city": "Exit"}}`, // Partial nested match
    //     }

    //     for _, k := range known {
    //         var w, g map[string]interface{}
    //         if err := json.Unmarshal([]byte(k), &w); err != nil { t.Fatal(err) }
    //         if err := json.Unmarshal(unknown, &g); err != nil { t.Fatal(err) }

    //         if err := contains(w, g); err != nil {
    //             t.Errorf("Test failed for known: %s, error: %v", k, err)
    //         }
    //     }
    // }
    ```

## What's New
*   In Go 1.18, the `reflect.Ptr` and `reflect.PtrTo` functions were renamed to `reflect.Pointer` and `reflect.PointerTo` respectively, for consistency within the `reflect` package. The old names continue to work but are deprecated. [3]
*   In Go 1.21, the `reflect.SliceHeader` and `reflect.StringHeader` types were deprecated. For direct manipulation of slice and string values, the `unsafe.Slice`, `unsafe.SliceData`, `unsafe.String`, or `unsafe.StringData` functions from the `unsafe` package are now preferred. [6]
*   In Go 1.23, a bug in `reflect.DeepEqual` was fixed where it could incorrectly return `true` when comparing `netip.Addr` values representing an IPv4 address and its IPv4-mapped IPv6 form, even though they are different. This change makes `reflect.DeepEqual` more accurate for these specific comparisons. [8]

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