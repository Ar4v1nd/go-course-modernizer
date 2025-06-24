# Go Class: 00 Intro and Why Use Go?

## Summary
This video introduces a Go programming course, outlining its scope from fundamental concepts to advanced topics like concurrency, performance optimization, and backend development. The instructor emphasizes the core reasons for choosing Go: its simplicity and readability for improved software engineering, and its efficiency and suitability for cloud-native applications.

## Key Points

*   **Course Introduction**
    *   The course covers Go programming from the basics, assuming prior knowledge of general programming constructs like if statements, for loops, and function calls.
    *   Topics include Go program structure, object-oriented programming in Go, concurrent programming, mechanical sympathy, benchmarking, profiling, the Go toolchain, and best practices.
    *   The course concludes with building backend applications using REST and GraphQL in Go.

*   **Recommended Reading**
    *   The book "The Go Programming Language" by Alan A. A. Donovan and Brian W. Kernighan is highly recommended.
    *   The course structure and some examples/exercises align with the book's content.

*   **Why Go? (Software Engineering)**
    *   Go's language design prioritizes software engineering, focusing on "programming in the large" with multiple developers and long-term maintenance.
    *   The goal is to create reliable and maintainable programs that are easy to understand and modify over time.
    *   Simplicity and clarity are valued over "clever" or complex code, as complexity can become a problem rather than a solution.

*   **Why Go? (Simplicity)**
    *   Programs should be written primarily for human readability, with machine execution being secondary.
    *   A complicated language itself can hinder software development.
    *   Go aims for a small, understandable specification that developers can largely keep in their heads, reducing the need for constant reference or expert consultation.
    *   This simplicity makes Go accessible for new developers and promotes consistent code style (e.g., enforced by `gofmt`).

*   **Why Go? (Design Goals)**
    *   Go was designed to combine the ease of use found in interpreted, dynamically typed languages with the efficiency and safety of statically typed, compiled languages.
    *   Key design goals include simplicity, safety, and readability.
    *   The language emphasizes orthogonality and often provides "one right way to do things" to reduce complexity and improve consistency.
    *   Development efforts have focused on runtime improvements (like garbage collection) and tooling, rather than adding numerous new language features.

*   **Why Go? (Go is Boring)**
    *   The "boring" nature of Go (lack of excessive features or complex paradigms) is considered a strength.
    *   Its simplicity is a powerful feature in itself, contributing to reliability and ease of understanding.

*   **Why Go? (Microprocessor Trends)**
    *   Around 2005, microprocessor trends shifted from increasing single-thread performance and frequency to increasing the number of logical cores per CPU.
    *   Many popular programming languages were designed before this shift, when concurrency was not a primary concern.

*   **Why Go? (New Computing Landscape)**
    *   Modern computing environments are characterized by multicore processors, networked systems, massive clusters, web programming models (like REST), huge programs, large development teams, and long build times.
    *   Older languages often treat concurrency as an afterthought, making it difficult to leverage modern hardware effectively.
    *   To achieve performance in this new landscape, software must either be concurrent or waste fewer resources.

*   **Why Go? (Cloud Cost Efficiency)**
    *   In cloud environments, computing resources are rented, making efficiency directly tied to cost.
    *   Go's performance allows applications to run significantly faster and use fewer resources, leading to substantial cost savings (e.g., reducing server count from 30 to 2).

*   **Why Go? (Built for the Cloud)**
    *   Go has become a preferred language for cloud infrastructure and application development.
    *   Its compiled nature produces self-contained binaries, simplifying deployment into containers (e.g., Docker, Kubernetes).
    *   Go binaries do not require a JVM, interpreter, or external libraries like libc, resulting in smaller, more secure containers with fewer potential vulnerabilities.
    *   Many prominent cloud-native projects are built with Go (e.g., Docker, Kubernetes, Prometheus, CockroachDB).

*   **Why Go? (Bitly's Perspective)**
    *   Go offers the speed and robustness of a compiled, statically-typed language without excessive complexity.
    *   It provides clear ways to express concurrent solutions for parallelizable problems.
    *   It achieves these benefits with minimal sacrifice in functionality.
    *   Tools like `gofmt` ensure consistently styled code, aiding collaboration and readability.

*   **Why Go? (Dennis Ritchie's Wisdom)**
    *   A programming language that is not overly complex and doesn't try to include "everything" can actually be easier to program in.

## What's New
- The statement "Development efforts have focused on runtime improvements (like garbage collection) and tooling, rather than adding numerous new language features" is no longer accurate. Go 1.18 introduced Generics, a significant new language feature, which represents a major change to the language. [3]
- The statement "Go binaries do not require a JVM, interpreter, or external libraries like libc, resulting in smaller, more secure containers with fewer potential vulnerabilities" is not universally accurate. On OpenBSD, Go binaries (especially non-static ones) now make system calls through libc. [1]

## Updated Code Snippets
(No updated code snippets are needed.)

## Citations
- [1] Go version 1.16
- [3] Go version 1.18