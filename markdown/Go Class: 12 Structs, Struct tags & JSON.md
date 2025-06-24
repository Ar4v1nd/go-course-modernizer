# Go Class: 12 Structs, Struct tags & JSON

## Summary
This video introduces Go structs as a way to group together data of different types, similar to records in a database. It covers various ways to declare and initialize structs, including using dot notation and struct literals. The concept of pointers to structs and embedding structs within other structs is explained. A significant portion is dedicated to understanding struct compatibility (structural vs. named typing) and the implications of using structs as values in maps, highlighting a common "gotcha" related to taking addresses of map elements. Finally, the video delves into struct tags, particularly their use with the `encoding/json` package for marshaling and unmarshaling JSON data, and touches upon other applications like SQL queries.

## Key Points

### Structs: Definition and Basic Initialization
*   **Definition:** Structs are user-defined composite types that group together zero or more named fields, each having a name and a type.
*   **Declaration:** Use the `type` keyword followed by the struct name and the `struct` keyword, enclosing fields in curly braces. Field names start with an uppercase letter for exportability.
    ```go
    type Employee struct {
        Name   string
        Number int
        Boss   *Employee // Pointer to self-referential type
        Hired  time.Time
    }
    ```
*   **Zero Value:** The zero value of a struct is a struct with all its fields initialized to their respective zero values (e.g., empty string for `string`, 0 for `int`, `nil` for pointers, and a default time for `time.Time`).
    ```go
    var e Employee
    // e will be {Name:"", Number:0, Boss:nil, Hired:0001-01-01 00:00:00 +0000 UTC}
    ```
*   **Printing Structs:**
    *   `%v`: Prints the struct values.
    *   `%+v`: Prints the struct values along with their field names, which is often more readable.
    *   `%#v`: Prints the Go-syntax representation of the value.

### Initializing and Accessing Struct Fields
*   **Dot Notation:** Individual fields can be accessed and assigned values using dot notation.
    ```go
    var e Employee
    e.Name = "Matt"
    e.Number = 1
    e.Hired = time.Now()
    ```
*   **Struct Literals:** Structs can be initialized directly using a struct literal.
    *   **Positional Initialization:** Values are assigned in the order of field declaration. All fields must be provided.
        ```go
        e := Employee{"Matt", 1, nil, time.Now()}
        ```
    *   **Named Field Initialization:** Values are assigned by specifying the field name. This allows for partial initialization (unspecified fields get their zero value) and arbitrary order. This is generally preferred for readability and maintainability.
        ```go
        e := Employee{
            Name:   "Matt",
            Number: 1,
            Hired:  time.Now(),
        }
        ```
*   **Pointers to Structs:** A field can be a pointer to another struct, including a pointer to the same struct type, allowing for recursive data structures like employee hierarchies.
    ```go
    type Employee struct {
        Name   string
        Number int
        Boss   *Employee // Pointer to Employee
        Hired  time.Time
    }
    ```
    When accessing fields of a struct pointer, Go automatically dereferences the pointer, so `p.Field` works directly without needing `(*p).Field`.

### Structs and Maps
*   **Storing Structs in Maps:** Maps can store structs as values.
    ```go
    c := make(map[string]Employee)
    c["Matt"] = Employee{Name: "Matt", Number: 1, Hired: time.Now()}
    ```
*   **The "Gotcha" with Map Value Addresses:** Go does not allow taking the address of a struct value directly from a map (`&c["Lamine"]`) or modifying its fields directly (`c["Lamine"].Number++`). This is because map values are not addressable; they might be moved in memory during map operations (like resizing a hash table).
    *   **Solution:** Store pointers to structs in maps instead of the structs themselves.
        ```go
        c := make(map[string]*Employee) // Map stores pointers
        c["Lamine"] = &Employee{Name: "Lamine", Number: 2, Hired: time.Now()}
        // Now you can modify fields through the pointer
        c["Lamine"].Number++
        ```

### Struct Compatibility
*   **Anonymous Structs:** Two anonymous struct types are compatible if:
    *   They have the same field names.
    *   They have the same field types.
    *   The fields are in the same order.
    *   They have the same struct tags (this rule changed in Go 1.8, see below).
    ```go
    // Compatible anonymous structs
    v1 := struct { X int `json:"foo"` } {1}
    v2 := struct { X int `json:"foo"` } {2}
    v1 = v2 // This assignment is valid
    ```
*   **Named Structs:** Named struct types are never compatible with other named struct types, even if their underlying structure is identical. This is called "named typing."
    ```go
    type T1 struct { X int `json:"foo"` }
    type T2 struct { X int `json:"foo"` }
    v1 := T1{1}
    v2 := T2{2}
    // v1 = v2 // This would result in a "type mismatch" error
    ```
*   **Type Conversion for Named Structs:** Named struct types can be *converted* to other named struct types if they are structurally compatible (fields, types, order, and tags match).
    ```go
    type T1 struct { X int `json:"foo"` }
    type T2 struct { X int `json:"foo"` }
    v1 := T1{1}
    v2 := T2{2}
    v1 = T1(v2) // Type conversion is valid
    ```
*   **Go 1.8 Change for Struct Tags:** Prior to Go 1.8, struct tags were part of structural compatibility. From Go 1.8 onwards, differences in struct tags do not prevent type conversions between structurally identical named structs.
    ```go
    type T1 struct { X int `json:"foo"` }
    type T2 struct { X int `json:"bar"` } // Different tag
    v1 := T1{1}
    v2 := T2{2}
    v1 = T1(v2) // Valid conversion in Go 1.8+
    ```

### Zero Value of Structs
*   The zero value for a struct is the "zero" for each field in turn. This means all fields are initialized to their default zero values (e.g., 0 for numbers, empty string for strings, `nil` for pointers and slices).
*   **Best Practice:** It is usually desirable that the zero value of a struct be a natural or sensible default, making the struct immediately usable without explicit initialization.
    *   Example: `bytes.Buffer` struct's zero value is a ready-to-use empty buffer. Its `buf` (byte slice) is `nil` (which can be appended to), `off` (int) is 0, and `lastReadOp` (custom type) is `opInvalid`.

### Passing Structs to Functions
*   **Pass by Value (Default):** Structs are passed by value to functions. This means a copy of the struct is made, and any modifications inside the function will not affect the original struct.
    ```go
    func soldAnother(a album) { // 'a' is a copy
        a.copies++ // Modifies the copy, not the original
    }
    ```
*   **Pass by Pointer (for Modification):** To modify the original struct within a function, pass a pointer to the struct.
    ```go
    func soldAnother(a *album) { // 'a' is a pointer
        a.copies++ // Modifies the original struct through the pointer
    }
    ```
*   **Dot Notation on Pointers:** Go allows using dot notation directly on struct pointers (`a.copies`) without explicit dereferencing (`(*a).copies`).

### Empty Structs
*   An empty struct (`struct{}`) has no fields and takes up no space in memory.
*   **Use Cases:**
    *   **Set Type:** Can be used as the value type in a map to represent a set, where the presence of a key indicates membership.
        ```go
        var isPresent map[int]struct{}
        ```
    *   **Cheap Channel Type:** Can be used as the element type for channels when only signaling is needed, not data transfer.
        ```go
        done := make(chan struct{})
        ```
*   **Singleton:** All empty structs are considered identical (a singleton), as there's no way to differentiate them. Go keeps a single instance of `struct{}` in memory, and all references to `struct{}` point to this same instance.

### Struct Tags and JSON
*   **Purpose:** Struct tags are raw string literals associated with struct fields. They provide metadata about the field, often used by reflection-based libraries for encoding/decoding data to/from external formats (e.g., JSON, XML, SQL).
*   **Format:** A struct tag is a backtick-quoted string after the field's type. It typically contains key-value pairs separated by spaces. The key is followed by a colon and the value in double quotes.
    ```go
    type Response struct {
        Page  int    `json:"page"`
        Words []string `json:"words,omitempty"`
    }
    ```
*   **`encoding/json` Package:** This standard library package uses struct tags to control how Go structs are marshaled (encoded) into JSON and unmarshaled (decoded) from JSON.
    *   **`json:"name"`:** Specifies the JSON property name for the field. If omitted, the Go field name (converted to lowercase) is used.
    *   **`json:"-"`:** Ignores the field during marshaling/unmarshaling.
    *   **`json:"name,omitempty"`:** Omits the field from the JSON output if its value is the zero value (e.g., `0`, `""`, `nil`, empty slice/map).
*   **Example (Marshaling):**
    ```go
    r := &Response{Page: 1, Words: []string{"up", "in", "out"}}
    j, _ := json.Marshal(r)
    fmt.Println(string(j)) // Output: {"page":1,"words":["up","in","out"]}
    ```
    If `Words` was an empty slice, `omitempty` would prevent it from appearing in the JSON.
*   **Example (Unmarshaling):**
    ```go
    var r2 Response
    json.Unmarshal(j, &r2) // Unmarshal JSON bytes into a struct pointer
    fmt.Printf("%#v\n", r2) // Prints the Go struct representation
    ```
*   **Exported Fields Requirement:** For `encoding/json` (and other reflection-based packages) to access and process struct fields, the field names *must* start with an uppercase letter (i.e., they must be exported). If a field is unexported (starts with a lowercase letter), `encoding/json` will ignore it during marshaling and unmarshaling, even if it has a struct tag.
    *   Go's static analysis tool `go vet` will warn about unexported fields with JSON tags, as this is a common mistake.
*   **Other Uses of Struct Tags:** Struct tags are not limited to JSON. They are widely used in other libraries for various purposes, such as:
    *   **SQL Queries:** Libraries like `sqlx` use `db` tags to map struct fields to database column names for easier query building and result scanning.
        ```go
        type Item struct {
            Name string `db:"name"`
            When string `db:"created"`
        }
        // SQL query example using named parameters from struct tags
        // stmt := "INSERT INTO items (name, created) VALUES (:name, :created)"
        // db.NamedExec(stmt, item)
        ```
    *   **XML, Protocol Buffers (Protobuf), etc.**

This concludes the summary of the video content on Go structs, struct tags, and JSON.

## What's New

*   The `encoding/json` package introduced a new `omitzero` option for struct field tags. This option explicitly omits a field from JSON output if its value is the zero value, and it can use an `IsZero()` bool method if defined by the field's type. This is described as clearer and less error-prone than `omitempty` for the intent to omit zero values, particularly addressing friction with `time.Time` values. If both `omitempty` and `omitzero` are specified, the field is omitted if it's either empty or zero. [8]

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