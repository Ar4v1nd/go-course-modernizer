# Go Class: 18 Methods and Interfaces

## Summary
This video introduces methods and interfaces in Go, fundamental concepts for object-oriented programming. It explains how methods are functions bound to specific types via receivers and how interfaces define abstract behavior through method sets. The video emphasizes Go's structural typing ("duck typing"), where types implicitly satisfy interfaces by implementing their methods, without explicit declaration. It also covers the distinction between value and pointer receivers, interface composition, and the "has-a" relationship through embedding, contrasting it with traditional inheritance.

## Key Points

*   **Methods as Type-Bound Functions:**
    *   In Go, a method is a function associated with a specific type.
    *   Methods are declared with a *receiver* parameter before the function name.
    *   Example:
        ```go
        type MyType int

        func (m MyType) MyMethod() {
            // Method logic
        }
        ```

*   **Receivers: Value vs. Pointer:**
    *   **Value Receiver:** The method operates on a copy of the receiver. Changes to the receiver inside the method do not affect the original variable. Use for methods that do not modify the receiver's state.
        ```go
        type Point struct { X, Y float64 }
        func (p Point) Offset(dx, dy float64) Point {
            return Point{p.X + dx, p.Y + dy} // Returns a new Point
        }
        ```
    *   **Pointer Receiver:** The method operates on the original receiver via its memory address. Changes inside the method affect the original variable. Use for methods that modify the receiver's state.
        ```go
        func (p *Point) Move(dx, dy float64) {
            p.X += dx
            p.Y += dy // Modifies the original Point
        }
        ```

*   **Interfaces: Abstract Behavior and Structural Typing:**
    *   An interface defines a *method set*, specifying abstract behavior.
    *   A concrete type *satisfies* an interface if it implements *all* methods in the interface's method set. This is Go's *structural typing* (or "duck typing").
    *   There is no explicit `implements` keyword; satisfaction is implicit and checked at compile time.
    *   Example: `fmt.Stringer` interface and its satisfaction by a custom type.
        ```go
        // fmt.Stringer interface (from fmt package)
        // type Stringer interface {
        //     String() string
        // }

        type IntSlice []int
        func (is IntSlice) String() string {
            // Converts IntSlice to a string representation like "[1;2;3]"
            return "..."
        }

        // In main:
        var mySlice IntSlice = []int{1, 2, 3}
        var s fmt.Stringer = mySlice // mySlice (IntSlice) implicitly satisfies fmt.Stringer
        // fmt.Printf("%v", s) will call mySlice.String()
        ```

*   **Benefits of Interfaces:**
    *   **Abstraction:** Allows writing functions that operate on abstract behavior rather than specific concrete types.
    *   **Polymorphism:** Enables different concrete types to be used interchangeably if they satisfy the same interface.
    *   **Decoupling:** Reduces coupling between components, making code more flexible and testable.
    *   Example: `io.Writer` interface for writing data to various destinations.
        ```go
        // io.Writer interface (from io package)
        // type Writer interface {
        //     Write(p []byte) (n int, err error)
        // }

        // A function that can write to any io.Writer
        func WriteData(w io.Writer, data []byte) (int, error) {
            return w.Write(data)
        }
        ```

*   **Interface Composition:**
    *   Interfaces can embed other interfaces, combining their method sets.
    *   This promotes small, focused interfaces that can be composed into larger ones as needed.
    *   Example: `io.Reader`, `io.Writer`, and `io.ReadWriter`.
        ```go
        // io.Reader interface
        // type Reader interface { Read(p []byte) (n int, err error) }

        // io.Writer interface
        // type Writer interface { Write(p []byte) (n int, err error) }

        // io.ReadWriter interface
        type ReadWriter interface {
            Reader // embeds Reader's methods
            Writer // embeds Writer's methods
        }
        ```

*   **Composition vs. Inheritance (Go's Approach):**
    *   Go favors *composition* ("has-a" relationship) over traditional class-based *inheritance* ("is-a" relationship).
    *   Embedding a struct within another struct promotes the embedded struct's fields and methods to the outer struct.
    *   Example: `ColoredPoint` embedding `Point`.
        ```go
        type ColoredPoint struct {
            Point // Embedded field, promotes Point's fields (X, Y) and methods (e.g., Distance)
            Color color.RGBA
        }
        ```
    *   A `ColoredPoint` *has a* `Point`, but it *is not* a `Point` in terms of type substitution. You cannot pass a `ColoredPoint` directly where a `Point` is expected, even if it has all the `Point`'s methods. This strictness ensures type safety.

*   **Interface Satisfiability Rules:**
    *   A type satisfies an interface if it has *all* the methods specified by the interface.
    *   The receiver type of the concrete method must match the interface's expectation (value receiver for value methods, pointer receiver for pointer methods).
    *   If an interface method requires a pointer receiver, a value of the concrete type cannot directly satisfy the interface unless its address is taken.
    *   Methods must be declared in the *same package* as the type they are associated with. This allows the compiler to know all methods for a given type at compile time.

## What's New
*   **Interfaces: Abstract Behavior and Structural Typing:** The definition of an interface in Go has expanded beyond just a method set. With the introduction of generics in Go 1.18, interfaces can now also define a *type set*, allowing them to be used as type constraints. [3]

## Citations
- [3] Go version 1.18