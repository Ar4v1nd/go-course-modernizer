# Go Class: 38 Testing

## Summary
This video provides a comprehensive overview of testing in Go, covering its built-in features, different layers and goals of testing, and practical refactoring techniques for writing robust tests. It also delves into the philosophy and psychology behind effective software testing, emphasizing the importance of developer responsibility and the distinct roles of developers and dedicated testers in ensuring software quality.

## Key Points

### Go test features
*   Go provides standard tools and conventions for testing.
*   Test files must end with `_test.go` and contain functions named `TestXXX`. These files can reside in the same package directory or a separate directory.
*   Tests are executed using the `go test` command.
*   Go caches test results and will not re-run tests if the source code hasn't changed since the last test run.

### Layers of testing
*   Testing can be conceptualized in layers, broadly categorized into "Developer Testing" and "Tester Testing".
*   **Developer Testing** (transparent, fully automated, integrated with CI/CD) includes:
    *   Unit Testing: Testing individual components or functions in isolation.
    *   Integration Testing: Testing interactions between multiple components or services.
    *   End-to-End Testing: Testing the entire system flow from start to finish.
*   **Tester Testing** (opaque, adversarial, exploratory or automated, periodic or tied to release cycles) includes:
    *   Load/Performance Testing: Evaluating system behavior under various loads.
    *   System Testing: Testing the complete and integrated software system.
    *   Chaos Testing: Intentionally introducing failures to test system resilience.

### Goals
*   Developer testing aims to verify specific aspects of the code. Key areas to test for include:
    *   Extreme values (e.g., minimum, maximum inputs).
    *   Input validation (handling invalid or unexpected inputs).
    *   Race conditions (for concurrent code).
    *   Error conditions (how errors are handled and propagated).
    *   Boundary conditions (values at the edges of valid ranges).
    *   Pre- and post-conditions (state before and after function execution).
    *   Randomized data (fuzzing) to uncover unexpected behaviors.
    *   Configuration & deployment (ensuring correct setup in different environments).
    *   Interfaces to other software (for integration tests).
*   Unit tests should be self-contained and not rely on external resources like networks or real databases. Integration tests, however, should interact with actual external services to verify connectivity and correct interaction.

### Test functions
*   Go test functions have the signature `func TestCrypto(t *testing.T)`.
*   The `t *testing.T` parameter is used to report test failures and other information.
*   `t.Errorf("message", args...)` reports an error but allows the test to continue.
*   `t.Fatalf("message", args...)` reports an error and stops the test immediately.

```go
func TestCrypto(t *testing.T) {
    uuid := "650b5cc5-5c0b-4c00-ad97-36b08553c91d"
    key1 := "75abbabc1f9f8d28d55200b43fd9592"
    key2 := "75abbabc1f9f8d28d28d66200b43fd95962" // Slightly different key

    ct, err := secrets.MakeAppKey(key1, uuid)

    if err != nil {
        t.Errorf("make failed: %s", err)
    }
    // ... further assertions
}
```

### Table-driven tests
*   A common pattern in Go is to use table-driven tests, where test cases are defined in a slice of structs.
*   Each struct in the slice represents a single test case with its input values and expected outputs.
*   A loop iterates over this table, running the same test logic for each case.

```go
func TestValueFreeFloat(t *testing.T) {
    table := []struct {
        v float64
        s string
    }{
        {1, "1"},
        {1.1, "1.1"},
    }

    for _, tt := range table {
        v := Value{tt.v, m: &Machine{}} // Example usage of test data
        if s := v.String(); s != tt.s {
            t.Errorf("%%v: wanted %%s, got %%s", tt.v, tt.s, s)
        }
    }
}
```

### Table-driven subtests
*   For more complex table-driven tests, especially when individual test cases need separate reporting or setup/teardown, `t.Run()` can be used to create subtests.
*   Each subtest is run as an independent goroutine, allowing for parallel execution and clearer failure reporting.
*   The `t.Run()` method takes a name and a function (closure) as arguments.

```go
func TestGraphqlResolver(t *testing.T) {
    table := []struct {
        name string
        // ... other test case data
    }{
        {name: "retrieve_offer"},
        // ... more test cases
    }

    for _, st := range table {
        t.Run(st.name, func(t *testing.T) { // closure
            // Test logic for this subtest
            // Access st.name and other data from the outer scope
        })
    }
}
```

### A complex unit test example
*   For very complex test cases, it's beneficial to define a named struct for the test case data and a method on that struct to encapsulate the test logic.
*   This improves readability and maintainability, especially when test cases involve complex input/output structures.
*   The `run` method on the test struct takes `*testing.T` as a parameter.

```go
package oak

import (
    "bytes"
    "reflect"
    "testing"

    "oak/token" // Assuming 'oak/token' is a package with token definitions
)

type scanTest struct {
    name  string
    input string
    want  []token.Token
}

func (st scanTest) run(t *testing.T) {
    b := bytes.NewBufferString(st.input)
    c := ScanConfig{} // Assuming ScanConfig and NewScanner exist
    s := NewScanner(c, st.name, b)

    var got []token.Token
    for tok := s.Next(); tok.Type != token.EOF; tok = s.Next() {
        got = append(got, tok)
    }

    if !reflect.DeepEqual(st.want, got) {
        t.Errorf("line %%q, wanted %%v, got %%v", st.input, st.want, got)
    }
}

var scanTests = []scanTest{
    {
        name:  "simple-add",
        input: "2 -1 + # comment",
        want: []token.Token{
            {Type: token.Number, Line: 1, Text: "2"},
            {Type: token.Number, Line: 1, Text: "-1"},
            {Type: token.Operator, Line: 1, Text: "+"},
            {Type: token.Comment, Line: 1, Text: "# comment"},
        },
    },
    // ... more test cases
}

func TestScanner(t *testing.T) {
    for _, st := range scanTests {
        t.Run(st.name, st.run) // Using the run method as a subtest
    }
}
```

### More refactoring
*   To further parameterize test result checking, an interface can be defined for a "checker".
*   This allows different types of checks (e.g., deep equality, specific error checks, golden file comparisons) to be plugged into the test cases.
*   A `shouldFail` boolean can be added to test structs to indicate expected failures, allowing tests to verify error handling.

```go
type checker interface {
    check(*testing.T, string, string) bool // Example: checks got, want strings
}

type subTest struct {
    name      string
    shouldFail bool
    checker   checker // parameterize how we check results
    // ... other test data
}

// Example concrete checker type
type checkGolden struct {
    // ... fields for golden file logic
}

func (c checkGolden) check(t *testing.T, got, want string) bool {
    // ... logic to compare 'got' with 'want' (e.g., from a golden file)
    return true // or false on failure
}
```

### Mocking or faking
*   For unit tests that interact with external dependencies (like databases or microservices), it's often necessary to "mock" or "fake" these dependencies to isolate the code under test.
*   Go interfaces are crucial for this, as they allow different implementations (real vs. mock) to be swapped in.
*   A mock implementation can simulate various behaviors, including successful operations, specific errors, or even controlled delays.

```go
type DB interface {
    GetThing(string) (thing, error)
    // ... other methods
}

type mockDB struct {
    shouldFail bool
    // ... other mock state
}

var errShouldFail = errors.New("db should fail")

func (m mockDB) GetThing(key string) (thing, error) {
    if m.shouldFail {
        return thing{}, fmt.Errorf("%%s: %%w", key, errShouldFail)
    }
    // ... return a dummy successful thing
    return thing{}, nil
}
```

### Main test functions
*   Go allows defining a special `TestMain(m *testing.M)` function in a `_test.go` file.
*   This function is executed *before* any tests in the package are run, and it's responsible for setting up and tearing down test environments (e.g., starting/stopping emulators, databases).
*   It takes a `*testing.M` parameter, which has a `Run()` method to execute all tests in the package.
*   The `os.Exit()` function is used to return the test result code.

```go
func TestMain(m *testing.M) {
    stop, err := startEmulator() // Function to start external dependency
    if err != nil {
        log.Println("*** FAILED TO START EMULATOR ***")
        os.Exit(-1)
    }

    result := m.Run() // run all tests

    stop() // Function to stop external dependency
    os.Exit(result)
}
```

### Special test-only packages
*   To add test-only code that needs to access unexported (private) identifiers of a package, a separate test package can be created.
*   This package is named `packagename_test` (e.g., `package myfunc_test` for `package myfunc`).
*   Files in this `_test` package are not included in a regular build, only when running tests.
*   Unlike normal `_test.go` files within the same package, `packagename_test` can only access *exported* identifiers of the main package, treating it as an external consumer. This is useful for "opaque" or "black-box" tests.

```go
// file myfunc_test.go
package myfunc_test // This package is not part of package myfunc, so it has no internal access
```

### Philosophy of Testing
*   **Testing Culture:** "Your tests are the contract about what your software does and does not do. Unit tests at the package level should lock in the behavior of the package’s API. They describe, in code, what the package promises to do. If there is a unit test for each input permutation, you have defined the contract for what the code will do in code, not documentation." — Dave Cheney
*   "This is a contract you can assert as simply as typing `go test`. At any stage, you can know with a high degree of of confidence, that the behavior people relied on before your change continues to function after your change." — Dave Cheney
*   **Assumption about code:** You should assume your code doesn't work unless:
    *   You have tests (unit, integration, etc.).
    *   They work correctly (i.e., they fail when they should and pass when they should).
    *   You run them.
    *   They pass.
*   Your work isn't done until you've added or updated the tests. This is basic code hygiene: start clean, stay clean.

### Psychology of computer programming
*   "The hardest bugs are those where your mental model of the situation is just wrong, so you can't see the problem at all." — Brian Kernighan. This applies directly to testing.
*   Developers typically test to show that things *are* done and working according to their understanding of the problem and solution.
*   Most difficulties in software development are failures of imagination (e.g., not anticipating edge cases or unexpected inputs).

### Program correctness
*   There are eight levels of program correctness (Gries & Conway), in order of increasing difficulty:
    1.  It compiles (and passes static analysis).
    2.  It has no bugs that can be found just running the program.
    3.  It works for some hand-picked test data.
    4.  It works for typical, reasonable input.
    5.  It works with test data chosen to be difficult.
    6.  It works for all input that follows the specifications.
    7.  It works for all valid inputs and likely error cases.
    8.  It works for all input.
*   "It works" means it produces the desired behavior or fails safely.
*   Achieving higher levels of correctness (especially 7 and 8) often requires more than just empirical testing; it might involve formal methods or extensive fuzzing.

### Developer testing isn't enough
*   You can have 100% code coverage and still be wrong because:
    *   The code may be bug-free, but not match the requirements.
    *   The requirements may not match expectations.
    *   You can't test code that's missing (e.g., missing features due to incomplete requirements).
*   **Testers test to show that things *don't* work.** Their mindset is adversarial, trying to break the system.
*   Testers cannot test a system well if the requirements aren't documented (a major limitation of the Agile method as practiced by some).
*   Code and unit tests are simply not enough documentation for comprehensive testing.

### Testing is not "quality assurance"
*   Confusing "test" and "QA" is a basic mistake.
*   QA is a different discipline in software development; it's about managing the overall quality process, not just testing.
*   Software development is not a manufacturing process; you can't "test in" or "prove" quality by simply running tests.
*   Testing is not about running "acceptance" tests to show that things work.
*   **Testing is about surfacing defects by causing the system to fail (breaking it).**
*   The wrong testing mindset leads to inadequate testing.

### Reality check
*   The engineering notion of "good, fast, cheap - pick any two" applies to software testing. You can't have all three in the real world.
*   Effective and thorough testing is hard and expensive.
*   Software is often annoying because most organizations pick "fast and cheap" over "good" when it comes to testing.

## What's New

*   **Randomized data (fuzzing):** The concept of using randomized data for testing to uncover unexpected behaviors is now explicitly supported by Go's built-in fuzzing feature, introduced in Go 1.18 [3].
*   **`go get` command behavior:** The use of `go get` to build and install packages in module-aware mode was deprecated in Go 1.16 [1]. As of Go 1.18, `go get` no longer builds or installs packages in module-aware mode; it is now dedicated to adjusting dependencies in `go.mod` [3]. `go install example.com/cmd@latest` is the recommended way to install executables [3].
*   **`-i` flag deprecation:** The `-i` flag, accepted by `go build`, `go install`, and `go test`, was deprecated in Go 1.16 [1]. As of Go 1.20, `go build` and `go test` commands no longer accept this flag [5].
*   **`io/ioutil` package deprecation:** The `io/ioutil` package has been deprecated in Go 1.16, with its functionality moved to other packages like `io` and `os`. New code is encouraged to use the new definitions [1].
*   **Build constraint syntax (`//go:build`):** Go 1.17 introduced `//go:build` lines as a more readable way to write build constraints, preferring them over the older `// +build` lines [2]. The `gofmt` tool now automatically synchronizes these two forms, and the `vet` tool warns about mismatches [2].
*   **Loop variable capture in `for` loops:** The behavior of variables declared by a `for` loop changed in Go 1.22. Previously, these variables were created once and updated by each iteration, which could lead to accidental sharing bugs when captured by closures (e.g., in `t.Run` subtests). In Go 1.22, each iteration of the loop creates new variables, resolving these potential bugs by default [7].

## Citations
- [1] Go 1.16 Release Notes
- [2] Go 1.17 Release Notes
- [3] Go 1.18 Release Notes
- [5] Go 1.20 Release Notes
- [7] Go 1.22 Release Notes