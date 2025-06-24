# Go Class: 06 Control Statements; Declarations & Types

## Summary
This video provides a comprehensive overview of fundamental Go programming concepts, including control structures (if-else, for loops, switch statements), package management (visibility, imports, initialization), and variable declarations and type compatibility. It highlights Go's unique syntax and rules, such as mandatory braces, the single `for` loop, structural vs. named typing, and strict import requirements.

## Key Points

*   **Control Structures: If-then-else**
    *   All `if-then` statements require braces `{}`.
    *   Go enforces the "One True Brace Style" (1TBS).
    *   `if` statements can start with a short declaration or statement before the condition, separated by a semicolon.
        ```go
        if a == b {
            fmt.Println("a equals b")
        } else {
            fmt.Println("a is not equal to b")
        }

        if err := doSomething(); err != nil {
            return err
        }
        ```

*   **Control Structures: For Loops**
    *   Go has only one loop construct: `for`. There are no `do` or `while` loops.
    *   **Explicit control with an index variable**: Similar to C-style for loops, with optional initialization, condition, and increment parts.
        ```go
        for i := 0; i < 10; i++ {
            fmt.Printf("%d, %d\n", i, i*i)
        }
        ```
    *   **Implicit control through the `range` operator for arrays & slices**:
        *   `for i := range myArray` iterates over indices (0, 1, 2, ...).
        *   `for i, v := range myArray` iterates over indices and values. The value `v` is a copy.
        *   The loop ends when the range is exhausted.
        ```go
        for i := range myArray {
            fmt.Println(i, myArray[i])
        }

        for i, v := range myArray {
            fmt.Println(i, v)
        }
        ```
    *   **Implicit control through the `range` operator for maps**:
        *   `for k := range myMap` iterates over keys.
        *   `for k, v := range myMap` iterates over keys and values. The value `v` is a copy.
        *   Maps have no guaranteed order; iteration order is random.
        ```go
        for k := range myMap {
            fmt.Println(k, myMap[k])
        }

        for k, v := range myMap {
            fmt.Println(k, v)
        }
        ```
    *   **Infinite loop with explicit `break`**:
        *   A `for {}` loop creates an infinite loop.
        *   `break` statement exits the innermost loop.
        *   `continue` statement skips to the next iteration of the innermost loop.
        *   Labels can be used with `break` or `continue` to target specific outer loops in nested structures.
        ```go
        for { // Infinite loop
            // ... do something ...
            if condition {
                break // Exit the loop
            }
        }

        outer:
        for k := range testItemsMap { // keys
            for _, v := range itemsList { // values in list
                if k == v.ID {
                    continue outer // Found it! Continue outer loop
                }
            }
            t.Errorf("key not found: %s", k)
        }
        ```
    *   **Common mistake with `range` and blank identifier**:
        *   If you only need the values from a `range` loop (e.g., for a slice), you must use the blank identifier `_` for the index to avoid a compile error for an unused variable.
        ```go
        for _, v := range myArray {
            fmt.Println(v)
        }
        ```

*   **Control Structures: Switch**
    *   A `switch` statement is a shortcut for a series of `if-then-else` statements.
    *   Cases do not "fall through" by default; `break` statements are not required.
    *   Alternatives (cases) may be empty.
    *   The `default` case is always evaluated last, regardless of its position in the code.
    *   `switch` statements can start with a short declaration/statement.
        ```go
        switch a := f.Get(); a {
        case 0, 1, 2:
            fmt.Println("underflow possible")
        case 3, 4, 5, 6, 7, 8:
            // no-op
        default:
            fmt.Println("warning: overload")
        }
        ```
    *   **Switch on true**: A `switch` statement without an argument (e.g., `switch { ... }`) allows arbitrary boolean comparisons in its `case` clauses. These cases are evaluated in order from top to bottom.
        ```go
        switch {
        case a <= 2:
            fmt.Println("underflow possible")
        case a <= 8:
            // evaluated in order
        default:
            fmt.Println("warning: overload")
        }
        ```

*   **Packages**
    *   **Everything lives in a package**: Every Go source file must begin with a `package` declaration. There is no truly "global" scope outside of a package.
    *   **Visibility**: Names starting with a capital letter are "exported" (visible to other packages). Names starting with a lowercase letter are internal (private to the package).
    *   **Imports**: Each source file must explicitly `import` every package it uses. Unused imports are a compile error.
    *   **No cycles**: Packages cannot have cyclic dependencies (e.g., package A imports B, and B imports A). This ensures fast compilation and avoids initialization issues.
    *   **Initialization**: Items declared at package scope (constants, variables) are initialized before the `main()` function runs. Packages can also define an `init()` function (lowercase `i`), which the Go runtime automatically calls before `main()`. `init()` functions are private to the package and are executed in dependency order.

*   **Declarations & Compatibility**
    *   **Ways to introduce a name**:
        *   `const` for constant declarations.
        *   `type` for type declarations.
        *   `var` for variable declarations (must have type or initial value, sometimes both).
        *   Short, initialized variable declaration `:=` (only inside a function).
        *   `func` for function declarations (methods may only be declared at package level).
        *   Formal parameters and named returns of a function.
    *   **Variable declarations**:
        *   `var a int` (defaults to 0).
        *   `var b int = 1` (explicit type and value).
        *   `var c = 1` (type inferred as `int`).
        *   `var d = 1.0` (type inferred as `float64`).
        *   Multiple variables can be declared in a block:
            ```go
            var (
                x, y int
                z    float64
                s    string
            )
            ```
    *   **Short declarations (`:=`)**:
        *   Cannot be used outside of a function.
        *   Must be used (instead of `var`) in control statements (e.g., `if`, `for`).
        *   Must declare at least one *new* variable.
        *   It won't re-use an existing declaration from an outer scope if it's the *only* variable being declared in the short declaration. This creates a new variable that shadows the outer one.
        *   This can lead to subtle bugs if not understood, as the compiler will only flag an "unused variable" error if the *outer* shadowed variable is not used elsewhere.

*   **Structural Typing**
    *   Go uses structural typing for most types (arrays, slices, maps, structs, functions).
    *   It's the same type if it has the same structure or behavior.
        *   Arrays: same size and base type.
        *   Slices: same base type.
        *   Maps: same key and value types.
        *   Structs: same sequence of field names and types.
        *   Functions: same parameter and return types.
    *   Example: `a := [...]int{1, 2, 3}` and `b := [3]int{}` are compatible because they are both arrays of 3 integers. `c := [4]int{}` is not compatible with `a` or `b` because its size is different.

*   **Named Typing**
    *   Go uses named typing for non-function user-declared types (types defined with the `type` keyword).
    *   If you declare `type X int`, `X` is a *new, distinct* type from `int`, even though its underlying structure is an `int`.
    *   Direct assignment between `X` and `int` (or other named types with `int` as base) is a type mismatch error.
    *   Explicit type conversion is required for assignments between named types and their underlying types.
        ```go
        type X int
        func main() {
            var a X
            b := 12 // b defaults to int
            a = b   // TYPE MISMATCH (compile error)

            a = 12  // OK, 12 is an untyped literal
            a = X(b) // OK, type conversion
        }
        ```

*   **Numeric Literals**
    *   Go keeps "arbitrary" precision for literal values (256 bits or more) at compile time.
    *   Integer literals are untyped; they can be assigned to any integer type (e.g., `int`, `uint`, `int8`, `int64`) without explicit conversion, as long as the value fits.
    *   Float and complex literals are also untyped, picked by syntax (e.g., `2.0`, `2e9`, `2.0i`, `2i3`).
    *   Mathematical constants can be very precise.
    *   Constant arithmetic done at compile time doesn't lose precision.

*   **Basic Operators**
    *   **Arithmetic**: `+`, `-`, `*`, `/`, `%`, `++`, `--`. Operates on numbers only, except `+` which is overloaded for string concatenation.
    *   **Comparison**: `==`, `!=`, `<`, `<=`, `>`, `>=`. Only numbers and strings support order comparisons.
    *   **Boolean**: `!`, `&&`, `||`. Operates only on boolean types, with shortcut evaluation for `&&` and `||`.
    *   **Bitwise**: `&`, `|`, `^`, `<<`, `>>`, `&^`. Operates on integers.
    *   **Assignment**: `=`, `+=`, `-=`, `*=`, `/=`, `%=`, `&=`, `|=`, `^=`, `<<=`, `>>=`, `&^=`. These combine an operation with assignment.

*   **Operator Precedence**
    *   Go has only five levels of operator precedence, simplifying parsing and reducing ambiguity.
    *   Operators within the same precedence level are evaluated left-to-right.
    *   Levels (highest to lowest):
        1.  Multiplication-like: `*`, `/`, `%`, `<<`, `>>`, `&`, `&^`
        2.  Addition-like: `+`, `-`, `|`, `^`
        3.  Comparison: `==`, `!=`, `<`, `<=`, `>`, `>=`
        4.  Logical AND: `&&`
        5.  Logical OR: `||`
    *   Use parentheses `()` to explicitly control evaluation order if needed for clarity or to override default precedence.

## What's New
*   **Control Structures: For Loops**:
    *   In Go 1.22, variables declared by a `for` loop are created anew for each iteration, addressing accidental sharing bugs. [7]
    *   In Go 1.22, `for` loops can now range over integers. [7]
*   **Packages**:
    *   In Go 1.21, package initialization order was specified more precisely to ensure unambiguous definition. [6]
*   **Declarations & Compatibility**:
    *   The statement "Must be used (instead of var) in control statements (e.g., if, for)" is inaccurate. While `:=` is commonly used for short declarations in control statements, `var` can also be used. This behavior has not changed since Go 1.15.
*   **Named Typing**:
    *   Go 1.23 introduced preview support for generic type aliases, which allow the `type` keyword to define parameterized aliases (fully supported in Go 1.24). Unlike traditional named types, these are aliases and do not create new, distinct types. [8], [9]
*   **Numeric Literals**:
    *   Go 1.18 fixed a bug where the compiler now correctly reports an overflow when passing certain rune constant expressions to predeclared functions, making constant arithmetic behavior more consistent. [3]

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