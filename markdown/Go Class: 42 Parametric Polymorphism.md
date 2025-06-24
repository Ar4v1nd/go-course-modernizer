# Go Class: 42 Parametric Polymorphism

## Summary
This video introduces generics in Go, explaining them as a form of parametric polymorphism. The speaker discusses the rationale behind adding generics to Go, primarily focusing on replacing dynamic typing with static typing for improved compile-time safety. He emphasizes the Go team's cautious approach to integrating this feature, aiming for simplicity and avoiding unnecessary complexity often seen in other languages. The video demonstrates basic generic types, methods, and functions, including the use of type constraints and type inference.

## Key Points

*   **Introduction to Generics:**
    *   "Generics" is shorthand for "parametric polymorphism," which involves having a type parameter on a type or function.
    *   A generic type acts as a template for creating concrete types, e.g., `type MyType[T any] struct { v T; n int }`.
    *   The `any` keyword is a predeclared identifier that serves as a type parameter constraint, indicating that the type parameter can be any valid Go type (equivalent to `interface{}`).

*   **Why Generics in Go:**
    *   The primary reason for introducing generics is to replace dynamic typing (using `interface{}` and type assertions) with static typing. This shifts type checking from runtime to compile time, improving program safety.
    *   Performance improvements from using generics (e.g., avoiding runtime type assertions) should be considered a bonus, not the main motivation.
    *   It's crucial to continue using (non-empty) interfaces wherever possible, as they remain a powerful abstraction mechanism.
    *   Generics are a powerful feature for abstraction but can also be a source of unnecessary abstraction and complexity if not used judiciously. The goal is to cover common use cases without making the language overly complicated.

*   **Generic Types and Functions (Syntax & Usage):**
    *   **Generic Type Alias:** A type can be parameterized, creating a generic type alias.
        ```go
        type Vector[T any] []T
        ```
    *   **Generic Method:** Methods can be defined on generic types. For methods that modify the underlying data structure (like `append` for slices), a pointer receiver is necessary.
        ```go
        func (v *Vector[T]) Push(x T) {
            *v = append(*v, x) // may reallocate
        }
        ```
    *   **Generic Function:** Functions can also have type parameters.
        ```go
        func Map[F, T any](s []F, f func(F) T) []T {
            r := make([]T, len(s))
            for i, v := range s {
                r[i] = f(v)
            }
            return r
        }
        ```
    *   **Type Inference:** The Go compiler can often infer the type parameters when a generic function or type is used, reducing verbosity. For example, when calling `Map` with a `[]int` and a function `func(int) string`, the compiler infers `F` as `int` and `T` as `string`.

*   **Type Constraints:**
    *   Type parameters can be constrained to ensure they implement specific interfaces or have certain methods. This allows generic code to operate on types that guarantee specific behaviors.
    *   Example of a type constraint using `fmt.Stringer`:
        ```go
        // type constraint: T must have String() method
        type StringableVector[T fmt.Stringer] []T

        func (s StringableVector[T]) String() string {
            var sb strings.Builder
            sb.WriteString("<<")
            for i, v := range s {
                if i > 0 {
                    sb.WriteString(", ")
                }
                sb.WriteString(v.String()) // Safe call due to fmt.Stringer constraint
            }
            sb.WriteString(">>")
            return sb.String()
        }
        ```
    *   When using a generic type with a constraint, explicit instantiation of the type parameter might be required if the compiler cannot infer it (e.g., when using a named type like `num` which has `int` as its underlying type, but `int` itself doesn't implement `fmt.Stringer`).
        ```go
        type num int
        // ... (String() method for num) ...
        var s StringableVector[num] = []num{1, 2, 3} // Explicit 'num' required
        fmt.Println(s)
        ```

## What's New
*   **Introduction of Generics:** The core concepts and syntax for generics, including type parameters on types and functions, the `any` predeclared identifier (alias for `interface{}`), and the ability to define type constraints, were introduced in Go 1.18. [3]
*   **Generic Type Aliases:** While generic types and functions were introduced in Go 1.18, full support for generic type aliases (e.g., `type Vector[T any] []T`) was introduced in Go 1.24. Previously, this feature was in preview in Go 1.23 and required a `GOEXPERIMENT` flag. [8], [9]
*   **Type Inference Improvements:** The Go compiler's type inference capabilities, which allow it to automatically deduce type parameters when generic functions or types are used, have seen multiple improvements in power and precision since their introduction, notably in Go 1.21. [6]

## Citations
*   [1] Go 1.16 Release Notes
*   [2] Go 1.17 Release Notes
*   [3] Go 1.18 Release Notes
*   [4] Go 1.19 Release Notes
*   [5] Go 1.20 Release Notes
*   [6] Go 1.21 Release Notes
*   [7] Go 1.22 Release Notes
*   [8] Go 1.23 Release Notes
*   [9] Go 1.24 Release Notes