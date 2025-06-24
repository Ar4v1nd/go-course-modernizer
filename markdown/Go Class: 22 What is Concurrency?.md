# Go Class: 22 What is Concurrency?

## Summary
This video introduces the fundamental concepts of concurrency and parallelism in programming, distinguishing between them and explaining how they relate to program execution. It delves into the definition of a race condition, illustrating it with a practical bank account example, and discusses various strategies to prevent such issues, emphasizing the importance of atomic operations and the trade-offs involved.

## Key Points

*   **What is Concurrency?**
    *   Concurrency refers to the composition of independently executing parts of a program.
    *   Execution happens in some non-deterministic (partial) order.
    *   It implies undefined out-of-order or non-sequential execution of program parts.

*   **Partial Order:**
    *   In a partially ordered system, some operations have a defined sequence (e.g., A must happen before B), while others do not have a strict order relative to each other (e.g., C and D can happen in any order relative to each other, as long as their internal dependencies are met).

*   **Non-determinism:**
    *   Non-determinism means a program can exhibit different execution traces (sequences of operations) on different runs, even with the same input.
    *   This doesn't necessarily mean different *results*, but rather different *internal behaviors*.

*   **Independent Execution (Subroutines vs. Coroutines):**
    *   **Subroutines** (or functions/procedures) are subordinate; when a main program calls a subroutine, the subroutine executes entirely before control returns to the main program.
    *   **Coroutines** (like Go's goroutines) are co-equal; they run alongside the main program or other coroutines, and their execution can be interleaved.
    *   Concurrency can exist even on a single-core processor through interleaving (e.g., operating system interrupt handling).

*   **Concurrency vs. Parallelism:**
    *   **Concurrency** is about *dealing with* multiple things happening out-of-order. It's a way of structuring a program.
    *   **Parallelism** is about multiple things *actually happening at the same time*. It requires multiple processing units (e.g., multi-core processors).
    *   A program needs to be concurrent to achieve parallelism.
    *   Concurrency itself doesn't necessarily make a program faster; parallelism does (by utilizing multiple cores). However, concurrency can improve responsiveness by allowing other tasks to run while one is waiting (e.g., for I/O).

*   **Race Condition:**
    *   A race condition occurs when the system's behavior depends on the non-deterministic sequence or timing of independently executing program parts, and some possible execution orders produce *invalid results*.
    *   A race condition is a *bug*.
    *   **Example (Bank Account Deposit):**
        *   Consider two concurrent deposits to a shared bank account. Each deposit involves a "Read balance", "Modify balance", and "Write balance" sequence.
        *   If these sequences are interleaved (e.g., Read A, Read B, Modify A, Modify B, Write A, Write B), one deposit might overwrite the other's changes, leading to an incorrect final balance.

*   **Solving Race Conditions (Ensuring Consistency):**
    *   Race conditions arise when independent parts of a program *change* shared data.
    *   Solutions aim to ensure operations produce a consistent state for any shared data:
        *   **Don't share anything:** The simplest solution, if feasible.
        *   **Make shared things read-only:** If data is only read and never modified, race conditions on that data are avoided.
        *   **Allow only one writer to shared things:** Restrict modification to a single, controlled entity.
        *   **Make read-modify-write operations atomic:** This is the most common and robust solution for mutable shared state. An atomic operation is indivisible; once it starts, it completes without interruption from other concurrent operations on the same shared resource.
    *   Making operations atomic often involves adding more *sequential* order to parts of the program, which can reduce the overall concurrency and potentially impact performance. This is a necessary trade-off for correctness.

## What's New
*   **Race Condition Detection:** The race detector was improved in Go 1.16 to more precisely follow channel synchronization rules, potentially reporting more races [1].
*   **Atomic Operations:** Go 1.19 introduced new atomic types (e.g., `atomic.Int64`, `atomic.Pointer[T]`) in the `sync/atomic` package, making it easier to use atomic values [4]. Go 1.23 further added `And` and `Or` operators to `sync/atomic` for bitwise atomic operations [8].
*   **Loop Variable Semantics:** In Go 1.22, variables declared by a "for" loop are now created anew for each iteration, addressing a common source of accidental sharing bugs [7]. The `Vet` tool was updated in Go 1.22 to reflect this change and warn about incorrect loop variable capture [7].
*   **Testing Concurrent Code:** Go 1.24 introduced the experimental `testing/synctest` package to provide support for testing concurrent code [9].
*   **Error Handling:** Go 1.20 expanded error wrapping to allow an error to wrap multiple other errors, and `errors.Is` and `errors.As` functions were updated to inspect these multiply wrapped errors [5].
*   **Timer/Ticker Behavior:** In Go 1.23, `time.Timer` and `time.Ticker` channels became unbuffered, guaranteeing no stale values after `Reset` or `Stop` calls. This changes the behavior for programs that poll channel length [8].
*   **Memory Model Revision:** The Go memory model was revised in Go 1.19 to align with other languages (C++, Java, Rust, Swift), focusing on sequentially consistent atomics [4].
*   **Vet Tool Improvements:** The `vet` tool gained new warnings for invalid `testing.T` use in goroutines [1], missing values after `append` [7], deferring `time.Since` [7], and mismatched key-value pairs in `log/slog` calls [7]. Go 1.24 added a `copylock` analyzer for 3-clause for loops [9].

## Updated Code Snippets
No code snippets in the key points were outdated.

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