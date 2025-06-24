# Go Class: 40 Go Modules

## Summary
This video provides an introduction to Go Modules, focusing on their purpose, how they solve dependency management problems, and practical day-to-day usage for Go developers. It covers the core concepts of `go.mod` and `go.sum` files, the role of the Go proxy, versioning strategies, and commands for managing dependencies. The video also touches upon the philosophical aspects of software dependencies and the security implications.

## Key Points

*   **Why Go Modules?**
    *   Go Modules were introduced to address challenges in dependency management, such as avoiding the need for `$GOPATH`, grouping packages for versioning, and supporting semantic versioning.
    *   They provide in-project dependency management, ensuring strong dependency security and availability by leveraging proxy servers and checksums.
    *   Modules aim to offer the benefits of vendoring (local copies of dependencies) without requiring developers to commit all third-party code directly into their repositories.
    *   They protect against risks like flaky repositories, disappearing packages (e.g., the "left-pad" incident in JavaScript), and surreptitious changes to public packages.

*   **Import Compatibility Rule**
    *   If an old package and a new package share the same import path, the new package must be backward compatible with the old one.
    *   An incompatible updated package should use a new URL (version) to avoid conflicts.
    *   Example of importing different major versions:
        ```go
        package hello

        import (
            "github.com/x"    // Refers to v1
            x2 "github.com/x/v2" // Refers to v2
        )
        ```

*   **Some Control Files**
    *   `go.mod`: This file defines the module's name and its direct dependency requirements. It also specifies the Go version the module is built with (from Go 1.13 onwards).
        ```
        module hello
        require github.com/x v1.1
        go 1.13
        ```
    *   `go.sum`: This file contains cryptographic checksums for all direct and transitive dependencies, ensuring the integrity and authenticity of downloaded modules.
        ```
        github.com/x v1.1 h1:KqKTd5Bnrg8aKH3J...
        github.com/y v0.2 h1:Qz0iS0pjZuFQy/z7...
        github.com/z v1.5 h1:r8zfno3MHue2Ht5s...
        ```
    *   **Best Practice**: Always commit both `go.mod` and `go.sum` files to your version control system (e.g., Git repository).

*   **Environment Variables**
    *   Go Modules typically use default proxy and sum database settings:
        *   `GOPROXY=https://proxy.golang.org,direct`
        *   `GOSUMDB=sum.golang.org`
    *   For private repositories, specific environment variables need to be set to bypass the public proxy and sum database:
        *   `GOPRIVATE=github.com/xxx,github.com/yyy`
        *   `GONOSUMDB=github.com/xxx,github.com/yyy`
    *   Access to private GitHub repositories still requires proper authentication setup (e.g., SSH keys or personal access tokens).

*   **Module Proxy**
    *   The Go toolchain interacts with a module proxy (like `proxy.golang.org`) to fetch modules.
    *   The proxy caches modules, providing availability even if the original source disappears.
    *   The `go.sum` file's checksums are verified against a secure sum database (like `sum.golang.org`) to prevent tampering.

*   **`go.mod` Details: Pseudo-versions and Replacements**
    *   `go.mod` can record "pseudo-versions" for non-release or trunk versions of packages. These typically include a date and a commit hash.
        ```
        require (
            github.com/gen2brain/malgo v0.0.0-20181117112449-af6b9a0d538d
        )
        ```
    *   The `replace` directive allows developers to substitute a module path with another, useful for local development, testing unreleased versions, or applying temporary fixes.
        ```
        replace github.com/satori/go.uuid v1.2.0 => github.com/satori/go.uuid v1.2.1-0.20181028125025-b2ce2384e17b
        ```

*   **Maintaining Dependencies**
    *   To start a new project with modules:
        *   `$ go mod init <module-name>` (creates `go.mod`)
        *   `$ go build` (downloads dependencies and updates `go.mod` and `go.sum` in Go 1.15 and earlier; in Go 1.16+, `go build` is read-only for `go.mod`, requiring explicit `go get` or `go mod tidy` for updates).
    *   To update all dependencies to their latest compatible versions:
        *   `$ go get -u ./...`
        *   `$ go mod tidy` (removes unused modules and cleans up `go.mod`/`go.sum`).
    *   To list available versions of a dependency:
        *   `$ go list -m -versions rsc.io/sampler`
    *   To update a single dependency to a specific version:
        *   `$ go get github.com/gorilla/mux@latest` (latest tagged release)
        *   `$ go get github.com/gorilla/mux@v1.6.2` (specific tagged version)
        *   `$ go get github.com/gorilla/mux@e3702bed2` (specific commit hash)
        *   `$ go get github.com/gorilla/mux@master` (specific branch)
    *   **Crucial**: After any dependency changes, you *must* commit the updated `go.mod` and `go.sum` files to your repository.

*   **Vendoring and Local Cache**
    *   Vendoring is still supported: `$ go mod vendor` creates a `vendor/` directory containing copies of all dependencies.
    *   In Go 1.13, `go build -mod=vendor` was required to use the vendored modules; this is no longer necessary in Go 1.14+.
    *   Go maintains a local module cache (typically in `$GOPATH/pkg/mod` or a default location) where downloaded modules are stored.
    *   To clear the local module cache: `$ go clean -modcache`.

## What's New

*   **`go.mod` content for Go 1.17+ modules:** `go.mod` files for modules specifying Go 1.17 or higher now include explicit `require` directives for all transitively imported packages, not just direct dependencies. [2]
*   **`godebug` directive in `go.mod`:** `go.mod` and `go.work` files can now declare `godebug` settings. [8]
*   **`go.sum` behavior with `go mod tidy`:** The `go mod tidy` command now retains additional checksums in the `go.sum` file for modules whose source code is needed to verify that each imported package is provided by only one module in the build list. [3]
*   **`go build` and `go test` default behavior:** `go build` and `go test` commands no longer modify `go.mod` and `go.sum` by default; they operate in a read-only mode for these files. [1]
*   **`go get` command's primary role:** The `go get` command no longer builds or installs packages in module-aware mode; it is now dedicated solely to adjusting dependencies in `go.mod`. [3]
*   **Deprecation and removal of `go get -insecure`:** The `go get -insecure` flag was deprecated in Go 1.16 and removed in Go 1.17. Users should now use `GOINSECURE`, `GOPRIVATE`, or `GONOSUMDB` environment variables to manage insecure schemes or bypass sum validation. [2]
*   **`go mod tidy` with `-go` flag:** The `go mod tidy` subcommand now supports the `-go` flag to set or change the Go version in the `go.mod` file. [2]
*   **`go mod tidy` with `-diff` flag:** The `go mod tidy` subcommand now supports the `-diff` flag, which causes it to print necessary changes as a unified diff without modifying files. [8]
*   **Recommended way to install executables:** The recommended way to build and install executables in module mode is now `go install <package>@<version>`, as `go get` no longer performs building or installation. [1]
*   **`go mod vendor` behavior for Go 1.17+ modules:** For modules specifying Go 1.17 or higher, `go mod vendor` now annotates `vendor/modules.txt` with the Go version of each vendored module and omits `go.mod` and `go.sum` files for vendored dependencies. [2]

## Citations
*   [1] Go version 1.16
*   [2] Go version 1.17
*   [3] Go version 1.18
*   [8] Go version 1.23