# Go Class: 31 Odds & Ends

## Summary
This video covers several miscellaneous but important topics in Go programming, often referred to as "odds and ends." It delves into how to implement enumerated types using Go's constant blocks and `iota`, the functionality and usage of variable argument lists (variadic functions), the specifics of sized integers and their implications for low-level programming, and the application of bitwise operators. The video concludes with a brief discussion on the rare but sometimes useful `goto` statement.

## Key Points

### Enumerated Types
*   Go does not have a dedicated `enum` keyword like some other languages.
*   "Almost-enum" types can be created using a named integer type and a `const` block.
*   The `iota` keyword is a pre-declared identifier that acts as a constant generator, starting at `0` in each `const` block and incrementing by `1` for each subsequent constant declaration on a new line.

    ```go
    type shoe int

    const (
        tennisShoe shoe = iota // tennisShoe = 0
        dress                  // dress = 1
        sandal                 // sandal = 2
        clog                   // clog = 3
    )
    ```
*   `iota` can be used within expressions to generate sequences of values, such as powers of two for bit flags.

    ```go
    type Flags uint

    const (
        FlagUp Flags = 1 << iota // FlagUp = 1 (2^0)
        FlagBroadcast            // FlagBroadcast = 2 (2^1)
        FlagLoopback             // FlagLoopback = 4 (2^2)
        FlagPointToPoint         // FlagPointToPoint = 8 (2^3)
        FlagMulticast            // FlagMulticast = 16 (2^4)
    )
    ```
*   The first value in a `const` block can be ignored using the blank identifier `_`. `iota` will still increment, but the value won't be assigned to a named constant.

    ```go
    type ByteSize int64

    const (
        _ ByteSize = iota // iota = 0, ignored
        KiB ByteSize = 1 << (10 * iota) // iota = 1, KiB = 2^10
        MiB                              // iota = 2, MiB = 2^20
        GiB                              // iota = 3, GiB = 2^30
        TiB                              // iota = 4, TiB = 2^40
        PiB                              // iota = 5, PiB = 2^50
        EiB                              // iota = 6, EiB = 2^60
    )
    ```

### Variable Argument Lists
*   Go supports variadic functions, which can accept a variable number of arguments of a specified type.
*   The syntax for a variadic parameter is `...type`. It must be the *last* parameter in the function signature.

    ```go
    func sum(nums ...int) int {
        total := 0
        for _, num := range nums {
            total += num
        }
        return total
    }
    ```
*   Inside the function, the variadic parameter `nums` behaves like a slice of `int` (`[]int`).
*   When calling a variadic function, you can pass individual arguments or unpack a slice using `...`.

    ```go
    func main() {
        fmt.Println(sum())         // Call with no arguments (nums is an empty slice)
        fmt.Println(sum(1))        // Call with one argument
        fmt.Println(sum(1, 2, 3, 4)) // Call with multiple arguments

        s := []int{1, 2, 3, 4}
        fmt.Println(sum(s...)) // Unpack a slice into individual arguments
    }
    ```

### Sized Integers
*   Go provides fixed-size integer types (`int8`, `uint8`, `int16`, `uint16`, `int32`, `uint32`, `int64`, `uint64`) for scenarios requiring precise bit-width, such as handling low-level network protocols or binary file formats.
*   `int` and `uint` types are platform-dependent (typically 32-bit or 64-bit) and should be preferred for general-purpose integer arithmetic.
*   Unsigned integers (`uintN`) are useful when values are always non-negative, allowing for a larger positive range.

    ```go
    type TCPFields struct {
        SrcPort     uint16
        DstPort     uint16
        SeqNum      uint32
        AckNum      uint32
        DataOffset  uint8
        Flags       uint8
        WindowSize  uint16
        Checksum    uint16
        UrgentPtr   uint16
    }
    ```

### Bitwise Operators
*   Go supports standard bitwise operators:
    *   `&` (AND)
    *   `|` (OR)
    *   `^` (XOR)
    *   `&^` (AND NOT, or bit clear)
    *   `<<` (left shift)
    *   `>>` (right shift)
*   Hexadecimal literals start with `0x` (e.g., `0xfff0`).
*   Binary literals start with `0b` (e.g., `0b1111`).
*   Bitwise operators are crucial for manipulating individual bits within an integer, often used with flags.

    ```go
    func main() {
        a := uint16(0xfff7) // 1111111111110111
        b := uint16(0b1111) // 0000000000001111

        // AND NOT: Clears bits in 'a' where corresponding bits in 'b' are set.
        // a &^ b results in 1111111111110000 (0xfff0)
        fmt.Printf("%016b %#04x\n", a&^b, a&^b)

        // AND: Sets bits only where corresponding bits in both 'a' and 'b' are set.
        // a & b results in 0000000000000111 (0x0007)
        fmt.Printf("%016b %#04x\n", a&b, a&b)
    }
    ```
*   Checking for multiple bit flags: Combine flags with `|` (OR) to create a mask, then use `&` (AND) with the target value and compare the result to the mask.

    ```go
    // Example: Check if both TCPFlagSyn and TCPFlagAck are set in tcpHeader.Flags
    // TCPFlagSyn | TCPFlagAck creates a mask with both bits set.
    // (tcpHeader.Flags & (TCPFlagSyn | TCPFlagAck)) isolates those bits in Flags.
    // The comparison checks if *only* those specific bits are set as per the mask.
    synAck := (tcpHeader.Flags & (TCPFlagSyn | TCPFlagAck)) == (TCPFlagSyn | TCPFlagAck)
    ```

### Short Integers and Type Conversion Issues
*   When converting a larger integer type to a smaller one (e.g., `uint32` to `int16`), Go performs truncation, discarding the most significant bits. This can lead to unexpected values, especially with signed integers.
*   Signed integers use the most significant bit to represent the sign (0 for positive, 1 for negative in two's complement). Truncation can change a positive number into a negative one if the discarded bits affect the sign bit of the new, smaller type.
*   Go's strict type system requires explicit type conversions to make these operations clear and prevent silent errors. This forces the programmer to acknowledge potential data loss or sign changes.
*   For general programming, using the platform-dependent `int` or `uint` is recommended to avoid these complexities, as they typically provide sufficient range and handle arithmetic consistently. Fixed-size integers should be reserved for specific low-level requirements.

### `goto` Statement
*   The `goto` statement allows unconditional jumps to a labeled statement within the same function.
*   While generally discouraged in modern programming due to leading to "spaghetti code" and making programs harder to read and maintain, Go does support it.
*   In rare, specific scenarios (e.g., complex error handling or state machine implementations within a single function), `goto` can sometimes simplify code and improve readability compared to deeply nested `if` statements or complex loops.

    ```go
    readFormat: // Label for goto
        err = binary.Read(buf, binary.BigEndian, &header.format)
        if err != nil {
            return &header, nil, HeaderReadFailed.from(pos, err)
        }

        if header.format == junkID {
            // ... find size & consume WAVE junk header
            goto readFormat // Jump back to read the next header
        }

        if header.format != fmtID {
            return &header, nil, InvalidChunkType
        }
    ```

## What's New
*   **Tools: `go build` and `go test` no longer accept the `-i` flag.**
    The `-i` flag, which instructed the `go` command to install packages imported by packages named on the command line, was deprecated in Go 1.16 because it no longer had a significant effect on build times after the introduction of the build cache in Go 1.10. As of Go 1.20, `go build` and `go test` commands explicitly no longer accept this flag. [1], [5]
*   **Standard Library: Error Wrapping Enhancements.**
    Go 1.20 introduced new features for error wrapping. The `fmt.Errorf` function now supports multiple occurrences of the `%w` format verb to wrap multiple errors, and a new function `errors.Join` was added to return an error wrapping a list of errors. While the existing error handling in the `goto` example remains valid, these new features provide more flexible ways to handle and inspect errors. [5]

## Updated Code Snippets
No code snippets in the key points were outdated.

## Citations
*   [1] Go version 1.16
*   [5] Go version 1.20