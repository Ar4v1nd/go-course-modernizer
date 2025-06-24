# Go Class: 41 Building Go Programs

## Summary
This video provides a comprehensive guide on building and distributing Go applications, moving beyond basic development commands to production-ready deployments. It covers essential Go build tools, techniques for creating pure Go executables with static linking, cross-compilation for different platforms, recommended project layouts, the importance of documentation and Makefiles, and advanced Docker multi-stage builds for creating small, secure, and reproducible container images.

## Key Points

### Go Build Tools
*   **`go run`**: Used for quickly compiling and running Go programs during development.
*   **`go test`**: Used for running tests during development.
*   **`go build`**: Compiles Go source code into an executable binary. By default, it produces a dynamically linked binary that may depend on system libraries (like `libc`).
*   **`go install`**: Similar to `go build`, but it also copies the compiled binary to `$GOPATH/bin`, making it available in your system's PATH. This is commonly used for installing command-line tools.

### Pure Go Programs
*   **Concept**: A "pure Go" program is a Go executable that has no external runtime dependencies on system libraries (e.g., `libc`). All necessary components are statically linked into the single binary.
*   **Benefits**:
    *   **Smaller Container Images**: Eliminates the need for a full operating system base image, leading to extremely small container sizes (e.g., a few MBs vs. GBs for Java/Python).
    *   **Improved Security**: Reduces the attack surface by removing dependencies on potentially vulnerable system libraries.
    *   **Portability**: The single binary can be easily moved and run on any compatible system without needing to install external libraries.
*   **Building Pure Go**: To build a pure Go executable, you need to explicitly tell the Go compiler to use pure Go implementations for certain standard library packages (like `net` and `os/user`) and to statically link any remaining external C libraries.
    ```bash
    CGO_ENABLED=0 go build -a -tags netgo,osusergo -ldflags "-extldflags '-static' -s -w" -o myprogram .
    ```
    *   `CGO_ENABLED=0`: Disables CGo, preventing Go from linking against C libraries.
    *   `-a`: Forces rebuilding of all packages, including standard library packages, with the specified build tags.
    *   `-tags netgo,osusergo`: Instructs the `net` and `os/user` packages to use their pure Go implementations.
    *   `-ldflags "-extldflags '-static' -s -w"`: Passes flags to the linker:
        *   `-extldflags '-static'`: Forces static linking of external C libraries.
        *   `-s`: Strips the symbol table from the binary.
        *   `-w`: Omits DWARF symbol table information (for debugging).
*   **Verification**: Use `ldd` on Linux to check if the binary is dynamically linked. A pure Go binary will show "not a dynamic executable".

### Go Cross-Compilation
*   Go's build tool is a cross-compiler out-of-the-box, meaning you can compile a Go program on one operating system/architecture to run on another.
*   **Environment Variables**:
    *   `GOARCH`: Specifies the target architecture (e.g., `amd64`, `arm`, `arm64`).
    *   `GOOS`: Specifies the target operating system (e.g., `linux`, `darwin`, `windows`).
    *   `GOARM`: (For ARM architectures) Specifies the ARM version (e.g., `v7`).
*   **Example (Building for Raspberry Pi)**:
    ```bash
    GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 go build -a -tags netgo,osusergo -ldflags "-extldflags '-static' -s -w" -o mainPi ./main.go
    ```
    The resulting `mainPi` binary can be directly copied to a Raspberry Pi and executed.

### Project Layout
*   A well-structured project improves maintainability and onboarding for new developers.
*   **Root Directory**: Contains top-level files like `README.md`, `Makefile`, `go.mod`, `go.sum`.
*   **`cmd/`**: Contains main packages for executable programs. Each program typically resides in its own subdirectory (e.g., `cmd/myprogram/main.go`).
*   **`pkg/`**: Contains reusable libraries or packages that are not executable programs.
*   **`build/`**: Stores files related to building artifacts, such as `Dockerfile`s.
*   **`deploy/`**: Stores deployment-specific files, like Kubernetes YAML configurations.
*   **`scripts/`**: Contains miscellaneous scripts for development, testing, or operations.
*   **`vendor/`**: (Optional) Created by `go mod vendor` to store vendored copies of module dependencies.
*   **Principle**: Aim for the least necessary directory structure to avoid excessive nesting and complexity.

### Documentation
*   **Importance**: Good documentation, especially a comprehensive `README.md`, is crucial for reducing technical debt and enabling efficient collaboration.
*   **`README.md` Content**: Should provide essential information for anyone interacting with the project:
    *   **Overview**: Who is it for, what does it do, why does it exist?
    *   **Developer Setup**: Steps to get the development environment ready.
    *   **Project & Directory Structure**: Explanation of the project layout.
    *   **Dependency Management**: How dependencies are managed (e.g., Go modules).
    *   **Build/Install**: Instructions on how to build and/or install the application (including `make` targets).
    *   **Testing**: How to run unit tests, integration tests, end-to-end tests, and load tests.
    *   **Running**: How to run the application locally, in Docker, or in a cloud environment.
    *   **Database & Schema**: Details about database requirements and schema.
    *   **Credentials & Security**: How to handle sensitive information.
    *   **Debugging & Monitoring**: How to debug and monitor the application (metrics, logs).
    *   **CLI Tools**: Usage instructions for any command-line tools.

### Makefiles
*   Even with Go's simple build commands, Makefiles remain useful for orchestrating complex workflows.
*   **Reasons to use Makefiles**:
    *   **Calculate Parameters**: Dynamically generate build parameters (e.g., version strings from Git).
    *   **Other Steps/Dependencies**: Define dependencies between tasks (e.g., lint before build, build Docker image after Go binary).
    *   **Simplify Options**: Encapsulate long and complex command-line options for Go build commands.
    *   **Non-Go Commands**: Integrate commands for Docker, cloud providers, or other tools.
*   **Example (Makefile for Go Project)**:
    ```makefile
    SOURCES := $(wildcard *.go cmd/*/*.go) # Finds all Go source files

    # Calculates version from Git tags and commits
    VERSION=$(shell git describe --tags --long --dirty=2 /dev/null)
    ## We must have tagged the repo at least once for VERSION to work
    ifeq ($(VERSION),)
    VERSION = UNKNOWN
    endif

    .PHONY: sort
    sort: $(SOURCES)
    	go build -ldflags "-X main.version=$(VERSION)" -o ./cmd/sort

    .PHONY: lint
    lint:
    	golangci-lint run

    .PHONY: committed
    committed:
    	@git diff --exit-code > /dev/null || (echo "** COMMIT YOUR CHANGES FIRST **"; exit 1)

    .PHONY: docker
    docker: $(SOURCES) build/Dockerfile
    	docker build -t sort-anim:latest -f build/Dockerfile --build-arg VERSION=$(VERSION) .

    .PHONY: publish
    publish: committed lint docker
    	docker tag sort-anim:latest matthol2/sort-anim:$(VERSION)
    	docker push matthol2/sort-anim:$(VERSION)
    ```

### Versioning the Executable
*   **Importance**: Essential for identifying the exact code running in production, aiding in debugging and issue resolution.
*   **Method**:
    1.  Declare a global string variable in your `main` package (e.g., `var version string`). It must be uninitialized.
    2.  Use the `-X` linker flag during `go build` to set its value at compile time.
    3.  Generate the version string from Git using `git describe --tags --long --dirty=2 /dev/null`. This provides a version like `v1.0.0-2-ge5816a3` (tag, commits since tag, Git hash, and "dirty" if uncommitted changes exist).
*   **Reproducible Builds**: Avoid embedding build timestamps in the binary, as this makes every build unique even with identical source code and tools. The goal is to produce the exact same binary (same hash) from the same source and tools.

### Building in Docker
*   **Multi-stage Builds**: A Dockerfile technique that uses multiple `FROM` instructions to create smaller, more efficient images.
    1.  **Builder Stage**: Uses a larger image with all necessary build tools (e.g., `golang:1.15-alpine`). It compiles the Go application.
    2.  **Runtime Stage**: Uses a minimal base image (e.g., `busybox:musl` or `scratch`) and copies only the compiled binary and essential runtime dependencies (like CA certificates and timezone data) from the builder stage.
*   **Advantages**:
    *   **Minimal Image Size**: The final image contains only the application and its absolute necessities.
    *   **No Host Dependencies**: You don't need Go installed on your CI/CD server or local machine to build the Docker image.
    *   **Consistent Builds**: Ensures the build environment is always the same, improving reproducibility.
*   **Dockerfile Example (Combined)**:
    ```dockerfile
    # Stage 1: Builder
    FROM golang:1.15-alpine AS builder
    RUN /sbin/apk update && \
        /sbin/apk --no-cache add ca-certificates git tzdata && \
        /usr/sbin/update-ca-certificates
    RUN adduser -D -g '' sort
    WORKDIR /home/sort
    COPY go.mod go.sum /home/sort/
    COPY cmd /home/sort/cmd
    COPY *.go /home/sort
    ARG VERSION
    RUN CGO_ENABLED=0 go build -a -tags netgo,osusergo \
        -ldflags "-extldflags '-static' -s -w' -X main.version=${VERSION}" \
        -o sort ./cmd/sort

    # Stage 2: Runtime
    FROM busybox:musl
    COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
    COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
    COPY --from=builder /etc/passwd /etc/passwd
    COPY --from=builder /home/sort/sort /home/sort
    USER sort
    WORKDIR /home/sort
    EXPOSE 8081
    ENTRYPOINT ["/home/sort"]
    ```
*   **Deployment**: The resulting Docker image can be tagged with the Git version and pushed to a Docker registry (e.g., Docker Hub) for deployment.

## What's New
*   The `go install` command is now the recommended way to install command-line tools, especially with version suffixes (e.g., `go install example.com/cmd@latest`). It no longer writes pre-compiled package archives to `$GOROOT/pkg` and its behavior is more aligned with module-aware mode, installing binaries to `GOBIN` (which defaults to `$GOPATH/bin` or `$HOME/go/bin`). [1, p.3], [3, p.5], [5, p.2]
*   The `CGO_ENABLED=0` flag still explicitly disables CGo, but the `go` command now disables CGo by default on systems without a C toolchain. [5, p.3]
*   The `-a` flag, which forces rebuilding of all packages, is generally no longer necessary for ensuring pure Go builds due to improvements in the Go module system and build cache.
*   The `netgo` and `osusergo` build tags are no longer necessary for building pure Go executables for `net` and `os/user` packages, as these packages have been rewritten to use pure Go implementations when CGo is disabled. [1, p.4]
*   The `-s` and `-w` linker flags now behave more consistently across all platforms, with `-s` suppressing symbol table generation and `-w` omitting DWARF debug information generation. [7, p.4]
*   The default value for `GOARM` when cross-compiling to ARM is now always 7. Additionally, `GOARM` now supports suffixes like `,softfloat` or `,hardfloat` to specify floating-point implementation. [6, p.15], [7, p.13-14]
*   `go mod vendor` now annotates `vendor/modules.txt` with Go versions and omits `go.mod` and `go.sum` files for vendored dependencies when the main module specifies Go 1.17 or higher. [2, p.4]
*   Go 1.18 automatically embeds version control information (revision, commit time, dirty flag) and build information (tags, compiler flags, cgo status) into binaries. This information can be accessed via `go version -m` or `runtime/debug.ReadBuildInfo`, reducing the need for manual `-X` linker flags for basic versioning. [3, p.5]
*   Go 1.19 enhanced support for doc comments, including links, lists, and clearer headings, with `gofmt` now reformatting them for clarity. [4, p.2]
*   The Dockerfile example should be updated to use a more recent Go base image (e.g., `golang:1.24-alpine`) for the builder stage. Additionally, the `go build` command within the Dockerfile can be simplified by removing the `-a` and `-tags netgo,osusergo` flags, as they are no longer necessary for pure Go builds. [1, p.4]

## Updated Code Snippets
```bash
# Go Cross-Compilation Example (Building for Raspberry Pi)
GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 go build -ldflags "-extldflags '-static' -s -w" -o mainPi ./main.go
```

```dockerfile
# Dockerfile Example (Combined)
# Stage 1: Builder
FROM golang:1.24-alpine AS builder # Updated Go version
RUN /sbin/apk update && \
    /sbin/apk --no-cache add ca-certificates git tzdata && \
    /usr/sbin/update-ca-certificates
RUN adduser -D -g '' sort
WORKDIR /home/sort
COPY go.mod go.sum /home/sort/
COPY cmd /home/sort/cmd
COPY *.go /home/sort
ARG VERSION
# Simplified build command: -a and -tags netgo,osusergo are no longer needed
RUN CGO_ENABLED=0 go build \
    -ldflags "-extldflags '-static' -s -w' -X main.version=${VERSION}" \
    -o sort ./cmd/sort

# Stage 2: Runtime
FROM busybox:musl
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /home/sort/sort /home/sort
USER sort
WORKDIR /home/sort
EXPOSE 8081
ENTRYPOINT ["/home/sort"]
```

## Citations
- [1] Go version 1.16
- [2] Go version 1.17
- [3] Go version 1.18
- [4] Go version 1.19
- [5] Go version 1.20
- [6] Go version 1.21
- [7] Go version 1.22