# Go Class: 17 Go does OOP

## Summary
This video introduces Object-Oriented Programming (OOP) concepts and discusses how the Go language implements these principles. It covers abstraction, encapsulation, polymorphism, and inheritance, highlighting Go's unique approach, particularly its use of interfaces and composition instead of traditional class-based inheritance, aligning with a broader definition of OOP.

## Key Points

*   **Object-Oriented Programming (General Concepts)**
    *   For many, the essential elements of OOP are abstraction, encapsulation, polymorphism, and inheritance.
    *   Polymorphism and inheritance are often combined or confused in common understanding.
    *   Go's approach to OOP is similar in principle but different in implementation.

*   **Abstraction**
    *   **Concept:** Decoupling behavior from the implementation details.
    *   **Best Practice:** Design systems where users interact with a simplified view of functionality, without needing to know the complex underlying mechanisms.
    *   **Go Relevance:** Achieved through well-defined functions and interfaces that hide internal complexities.

*   **Encapsulation**
    *   **Concept:** Hiding implementation details from misuse.
    *   **Best Practice:** Protect internal state and logic to maintain the integrity of an abstraction and prevent external dependencies on internal details that might change.
    *   **Go Implementation:** Achieved through package-level visibility. Identifiers (variables, functions, types, fields) starting with an uppercase letter are exported (visible outside the package), while those starting with a lowercase letter are unexported (private to the package).
        ```go
        // In package 'mypackage'
        type MyType struct {
            ExportedField   string // Visible outside mypackage
            unexportedField int    // Hidden within mypackage
        }

        func ExportedFunc() { /* ... */ } // Visible outside mypackage
        func unexportedFunc() { /* ... */ } // Hidden within mypackage
        ```

*   **Polymorphism**
    *   **Concept:** Literally "many shapes" – allowing multiple types to be treated through a single interface, exhibiting different behaviors.
    *   **Traditional Types:** Ad-hoc (function/operator overloading), Parametric (generics), and Subtype (subclasses substituting for superclasses).
    *   **Go Implementation (Protocol-Oriented):** Go achieves polymorphism primarily through explicit interface types. Any type that implements all the methods declared in an interface implicitly satisfies that interface. This separates behavior from implementation.
        ```go
        type Greeter interface {
            Greet() string
        }

        type Person struct {
            Name string
        }
        func (p Person) Greet() string {
            return "Hello, " + p.Name
        }

        type Robot struct {
            Model string
        }
        func (r Robot) Greet() string {
            return "Greetings, human. I am " + r.Model
        }

        // Both Person and Robot can be used where a Greeter is expected.
        ```
    *   Go 1.15 (video version) does not have generics (parametric polymorphism), but it is planned for future versions.

*   **Inheritance**
    *   **Conflicting Meanings:** Often refers to both substitutability (a form of polymorphism) and structural sharing of implementation details.
    *   **Theoretical Ideal:** The Liskov Substitution Principle states that a subclass should be substitutable for its superclass ("is-a" relationship).
    *   **Practical Issues:**
        *   Injects strong dependencies: Subclasses become tightly coupled to the internal implementation of superclasses, making changes difficult and potentially breaking derived classes.
        *   "Leaky abstractions": Inheritance can lead to situations where a subclass doesn't logically fit all behaviors of its superclass (e.g., a `Line` is a `Shape` but doesn't have an `Area`).
        *   Often leads to overly complex and deep class hierarchies, which are hard to manage and maintain.
    *   **Modern Software Design:** There's a growing preference for "Composition over Inheritance" to achieve code reuse and flexibility without the tight coupling issues of inheritance.

*   **Alan Kay's View of OOP**
    *   Alan Kay, who coined the term "Object-Oriented Programming," emphasized messaging between self-contained objects, local retention and protection of state, and late-binding.
    *   His original vision de-emphasized inheritance hierarchies, focusing instead on objects communicating and exhibiting polymorphic behavior.

*   **OOP in Go**
    *   Go supports the core principles of OOP:
        *   **Encapsulation:** Through package-level visibility.
        *   **Abstraction & Polymorphism:** Through interfaces, where substitutability is based purely on abstract behavior.
        *   **Structural Sharing:** Through composition (embedding structs), which allows types to "have-a" relationship without "is-a" inheritance.
            ```go
            type Speaker struct {
                Volume int
            }

            type Phone struct {
                Model string
                Speaker // Embedded Speaker struct
            }

            // A Phone has a Speaker, inheriting its fields and methods.
            // This is composition, not inheritance.
            ```
    *   Go does not have traditional classes or inheritance based on types.
    *   This design choice allows defining methods on *any* user-defined type, not just "classes," and enables any object to implement an interface, promoting a more flexible and less constrained approach to object-oriented design.
    *   Go's approach encourages thinking about software design in terms of behavior and communication, rather than rigid type hierarchies.

## What's New
*   The statement "Go 1.15 (video version) does not have generics (parametric polymorphism), but it is planned for future versions" is no longer accurate. Go 1.18 introduced generics, providing parametric polymorphism through type parameters for functions and types. [3]

## Citations
- [3] Go version 1.18