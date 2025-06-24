# Go Class: 19 Composition

## Summary
This video delves into object-oriented programming in Go, focusing on struct composition as an alternative to traditional inheritance. It explains how embedding structs (and even non-struct types) promotes their fields and methods to the embedding structure's level, simplifying access. The concept is demonstrated through practical examples, including custom sorting implementations using Go's `sort.Interface` and `sort.Reverse`. Finally, the video highlights how Go's design allows for making `nil` values useful, particularly when calling methods on `nil` receivers, leading to more elegant and robust code for certain data structures.

## Key Points

*   **Struct Composition:**
    *   Go achieves object-oriented design through composition, specifically struct embedding.
    *   When a struct embeds another type (e.g., `type PairWithLength struct { Pair; Length int }`), the fields and methods of the embedded type are "promoted" to the embedding structure's level.
    *   This allows direct access to promoted fields (e.g., `pl.Path` instead of `pl.Pair.Path`).
    *   Methods of the embedded type are also promoted and can be called directly on the embedding struct.
    *   If the embedding struct defines a method with the same name as a promoted method, the embedding struct's method takes precedence (similar to method overriding, but Go avoids the term "override" to distinguish from inheritance).
    *   Composition is *not* inheritance: An embedding struct is *not* a subtype of the embedded type. You cannot pass an instance of `PairWithLength` to a function expecting a `Pair` directly; you must explicitly pass the embedded field (`pl.Pair`).
    *   A struct can embed a pointer to another type (e.g., `*PairWithLength`), and promotion still works the same way.

    ```go
    type Pair struct {
        Path string
        Hash string
    }

    func (p Pair) String() string {
        return fmt.Sprintf("Hash of %s is %s", p.Path, p.Hash)
    }

    type PairWithLength struct {
        Pair // Embedded struct
        Length int
    }

    // Example of promoted field access:
    // pl := PairWithLength{Pair{"/usr", "0xdfdfe"}, 121}
    // fmt.Println(pl.Path, pl.Length) // Accesses promoted Path and direct Length

    // Example of overriding a promoted method:
    func (p PairWithLength) String() string {
        return fmt.Sprintf("Hash of %s is %s; length %d", p.Path, p.Hash, p.Length)
    }
    ```

*   **Sorting with Interfaces:**
    *   Go's `sort` package uses interfaces to provide generic sorting functionality.
    *   The `sort.Interface` interface defines three methods:
        *   `Len() int`: Returns the number of elements in the collection.
        *   `Less(i, j int) bool`: Reports whether the element at index `i` should sort before the element at index `j`.
        *   `Swap(i, j int)`: Swaps the elements at indices `i` and `j`.
    *   The `sort.Sort(data sort.Interface)` function can sort any data type that implements these three methods.
    *   Custom types can implement `sort.Interface` by embedding the slice type and defining the required methods.

    ```go
    type Organ struct {
        Name  string
        Weight int
    }

    type Organs []Organ // A slice of Organ structs

    func (s Organs) Len() int {
        return len(s)
    }

    func (s Organs) Swap(i, j int) {
        s[i], s[j] = s[j], s[i]
    }

    // Custom types for sorting by different criteria
    type ByName struct {
        Organs // Embeds the Organs slice
    }

    func (s ByName) Less(i, j int) bool {
        return s.Organs[i].Name < s.Organs[j].Name // Compares by Name
    }

    type ByWeight struct {
        Organs // Embeds the Organs slice
    }

    func (s ByWeight) Less(i, j int) bool {
        return s.Organs[i].Weight < s.Organs[j].Weight // Compares by Weight
    }

    // Example usage:
    // s := Organs{{"brain", 1340}, {"liver", 1494}, {"spleen", 162}}
    // sort.Sort(ByWeight{s}) // Sorts 's' by weight
    // sort.Sort(ByName{s})   // Sorts 's' by name
    ```

*   **Sorting in Reverse:**
    *   The `sort.Reverse` struct in the `sort` package allows sorting in reverse order without modifying the original `Less` implementation.
    *   `sort.Reverse` embeds `sort.Interface` and redefines its `Less()` method to call the embedded interface's `Less()` method with swapped indices (`r.Interface.Less(j, i)`).

    ```go
    // Defined in sort package:
    // type Reverse struct {
    //     Interface
    // }
    // func (r Reverse) Less(i, j int) bool {
    //     return r.Interface.Less(j, i)
    // }

    // Example usage:
    // entries := []string{"charlie", "able", "dog", "baker"}
    // sort.Sort(sort.Reverse(sort.StringSlice(entries))) // Sorts entries in reverse alphabetical order
    ```

*   **Making Nil Useful:**
    *   In Go, it's safe to call methods on a `nil` receiver (for pointer receivers). The method will execute, and it's the method's responsibility to handle the `nil` case.
    *   This allows for simpler and more elegant code, especially for recursive data structures like linked lists or trees, where `nil` naturally represents the end of the list/branch.
    *   Example: A `Sum()` method for a linked list (`*IntList`) can check if the receiver is `nil` as its base case.

    ```go
    type IntList struct {
        Value int
        Tail  *IntList
    }

    // Sum returns the sum of the list elements.
    func (list *IntList) Sum() int {
        if list == nil { // Safe to check if receiver is nil
            return 0
        }
        return list.Value + list.Tail.Sum() // Recursively call Sum on Tail
    }

    // Example usage:
    // var myList *IntList // myList is nil by default
    // total := myList.Sum() // This call is safe and returns 0
    ```

## What's New
*   While the `sort.Interface` and `sort.Sort` functions remain valid for sorting, Go 1.21 introduced the `slices` package, which provides generic functions like `slices.Sort`, `slices.SortFunc`, and `slices.Reverse`. These new functions offer a more direct and type-safe way to sort slices of any element type, often preferred for new code. [6]

## Citations
*   [6] Go version 1.21