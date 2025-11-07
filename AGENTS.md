# Agent Guidelines for articulate-parser

## Build/Test Commands

- **Build**: `task build` or `go build -o bin/articulate-parser main.go`
- **Run tests**: `task test` or `go test -race -timeout 5m ./...`
- **Run single test**: `go test -v -race -run ^TestName$ ./path/to/package`
- **Test with coverage**:
  - `task test:coverage` or
  - `go test -race -coverprofile=coverage/coverage.out -covermode=atomic ./...`
- **Lint**: `task lint` (runs vet, fmt check, staticcheck, golangci-lint)
- **Format**: `task fmt` or `gofmt -s -w .`
- **CI checks**: `task ci` (deps, lint, test with coverage, build)

## Code Style Guidelines

### Imports

- Use `goimports` with local prefix: `github.com/kjanat/articulate-parser`
- Order: stdlib, external, internal packages
- Group related imports together

### Formatting

- Use `gofmt -s` (simplify) and `gofumpt` with extra rules
- Function length: max 100 lines, 50 statements
- Cyclomatic complexity: max 15
- Cognitive complexity: max 20

### Types & Naming

- Use interface-based design (see `internal/interfaces/`)
- Export types/functions with clear godoc comments ending with period
- Use descriptive names: `ArticulateParser`, `MarkdownExporter`
- Receiver names: short (1-2 chars), consistent per type

### Error Handling

- Always wrap errors with context: `fmt.Errorf("operation failed: %w", err)`
- Use `%w` verb for error wrapping to preserve error chain
- Check all error returns (enforced by `errcheck`)
- Document error handling rationale in defer blocks when ignoring close errors

### Comments

- All exported types/functions require godoc comments
- End sentences with periods (`godot` linter enforced)
- Mark known issues with TODO/FIXME/HACK/BUG/XXX

### Security

- Use `#nosec` with justification for deliberate security exceptions (G304 for CLI file paths, G306 for export file permissions)
- Run `gosec` and `govulncheck` for security audits

### Testing

- Enable race detection: `-race` flag
- Use table-driven tests where applicable
- Mark test helpers with `t.Helper()`
- Benchmarks in `*_bench_test.go`, examples in `*_example_test.go`

### Dependencies

- Minimal external dependencies (currently: go-docx, golang.org/x/net, golang.org/x/text)
- Run `task deps:tidy` after adding/removing dependencies

---

## Go 1.24 & 1.25 New Features Reference

This project uses Go 1.24+. Below is a comprehensive summary of new features and changes in Go 1.24 and 1.25 that may be relevant for development and maintenance.

### Go 1.24 Major Changes (Released February 2025)

#### Language Features

**Generic Type Aliases**

- Type aliases can now be parameterized with type parameters
- Example: `type List[T any] = []T`
- Can be disabled via `GOEXPERIMENT=noaliastypeparams` (removed in 1.25)

#### Tooling Enhancements

**Module Tool Dependencies**

- New `tool` directive in go.mod tracks executable dependencies
- Use `go get -tool <package>` to add tool dependencies
- Use `go install tool` and `go get tool` to manage them
- Eliminates need for blank imports in `tools.go` files

**Build Output Formatting**

- Both `go build` and `go test` support `-json` flag for structured JSON output
- New action types distinguish build output from test results

**Authentication**

- New `GOAUTH` environment variable provides flexible authentication for private modules

**Automatic Version Tracking**

- `go build` automatically sets main module version in binaries based on VCS tags
- Adds `+dirty` suffix for uncommitted changes

**Cgo Performance Improvements**

- New `#cgo noescape` annotation: Prevents escape analysis overhead for C function calls
- New `#cgo nocallback` annotation: Indicates C function won't call back to Go

**Toolchain Tracing**

- `GODEBUG=toolchaintrace=1` enables tracing of toolchain selection

#### Runtime & Performance

**Performance Improvements**

- **2-3% CPU overhead reduction** across benchmark suites
- New Swiss Tables-based map implementation (faster lookups)
  - Disable via `GOEXPERIMENT=noswissmap`
- More efficient small object allocation
- Redesigned runtime-internal mutexes
  - Disable via `GOEXPERIMENT=nospinbitmutex`

#### Compiler & Linker

**Method Receiver Restrictions**

- Methods on cgo-generated types now prevented (both directly and through aliases)

**Build IDs**

- Linkers generate GNU build IDs (ELF) and UUIDs (macOS) by default
- Disable via `-B none` flag

#### Standard Library Additions

**File System Safety - `os.Root`**

- New `os.Root` type enables directory-limited operations
- Prevents path escape and symlink breakouts
- Essential for sandboxed file operations

**Cryptography Expansion**

- `crypto/mlkem`: ML-KEM-768/1024 post-quantum key exchange (FIPS 203)
- `crypto/hkdf`: HMAC-based Extract-and-Expand KDF (RFC 5869)
- `crypto/pbkdf2`: Password-based key derivation (RFC 8018)
- `crypto/sha3`: SHA-3 and SHAKE functions (FIPS 202)

**FIPS 140-3 Support**

- New `GOFIPS140` environment variable enables FIPS mode
- New `fips140` GODEBUG setting for cryptographic module compliance

**Weak References - `weak` Package**

- New `weak` package provides low-level weak pointers
- Enables memory-efficient structures like weak maps and caches
- Useful for preventing memory leaks in cache implementations

**Testing Improvements**

- `testing.B.Loop()`: Cleaner syntax replacing manual `b.N` iteration
- Prevents compiler from optimizing away benchmarked code
- New `testing/synctest` package (experimental) for testing concurrent code with fake clocks

**Iterator Support**

- Multiple packages now offer iterator-returning variants:
  - `bytes`: Iterator-based functions
  - `strings`: Iterator-based functions
  - `go/types`: Iterator support

#### Security Enhancements

**TLS Post-Quantum Cryptography**

- `X25519MLKEM768` hybrid key exchange enabled by default in TLS
- Provides quantum-resistant security

**Encrypted Client Hello (ECH)**

- TLS servers can enable ECH via `Config.EncryptedClientHelloKeys`
- Protects client identity during TLS handshake

**RSA Key Validation**

- Keys smaller than 1024 bits now rejected by default
- Use `GODEBUG=rsa1024min=0` to revert (testing only)

**Constant-Time Execution**

- New `crypto/subtle.WithDataIndependentTiming()` enables architecture-specific timing guarantees
- Helps prevent timing attacks

#### Deprecations & Removals

- `runtime.GOROOT()`: Deprecated; use system path instead
- `crypto/cipher` OFB/CFB modes: Deprecated (unauthenticated encryption)
- `x509sha1` GODEBUG: Removed; SHA-1 certificates no longer verified
- Experimental `X25519Kyber768Draft00`: Removed

#### Platform Changes

- **Linux**: Now requires kernel 3.2+ (enforced)
- **macOS**: Go 1.24 is final release supporting Big Sur
- **Windows/ARM 32-bit**: Marked broken
- **WebAssembly**:
  - New `go:wasmexport` directive
  - Reactor/library builds supported via `-buildmode=c-shared`

#### Bootstrap Requirements

- Go 1.24 requires Go 1.22.6+ for bootstrapping
- Go 1.26 will require Go 1.24+

---

### Go 1.25 Major Changes (Released August 2025)

#### Language Changes

- No breaking language changes
- "Core types" concept removed from specification (replaced with clearer prose)

#### Tooling Improvements

**Go Command Enhancements**

- `go build -asan`: Now defaults to leak detection at program exit
- New `go.mod ignore` directive: Specify directories for go command to ignore
- `go doc -http`: Starts documentation server and opens in browser
- `go version -m -json`: Prints JSON-encoded BuildInfo structures
- Module path resolution now supports subdirectories using `<meta>` syntax
- New `work` package pattern matches all packages in work/workspace modules
- Removed automatic toolchain line additions when updating `go` version

**Vet Analyzers**

- **"waitgroup"**: Detects misplaced `sync.WaitGroup.Add` calls
- **"hostport"**: Warns against using `fmt.Sprintf` for constructing addresses
  - Recommends `net.JoinHostPort` instead

#### Runtime Enhancements

**Container-Aware GOMAXPROCS**

- Linux now respects cgroup CPU bandwidth limits
- All OSes periodically update GOMAXPROCS if CPU availability changes
- Disable via environment variables or GODEBUG settings
- Critical for containerized applications

**New Garbage Collector - "Green Tea GC"**

- Experimental `GOEXPERIMENT=greenteagc` enables new GC
- **10-40% reduction in garbage collection overhead**
- Significant for GC-sensitive applications

**Trace Flight Recorder**

- New `runtime/trace.FlightRecorder` API
- Captures execution traces in in-memory ring buffer
- Essential for debugging rare events and production issues

**Other Runtime Changes**

- Simplified unhandled panic output
- VMA names on Linux identify memory purpose (debugging aid)
- New `SetDefaultGOMAXPROCS` function resets GOMAXPROCS to defaults

#### Compiler Fixes & Improvements

**Critical Nil Pointer Bug Fix**

- Fixed Go 1.21 regression where nil pointer checks were incorrectly delayed
- ⚠️ **May cause previously passing code to now panic** (correct behavior)
- Review code for assumptions about delayed nil checks

**DWARF5 Support**

- Debug information now uses DWARF version 5
- Reduces binary size and linking time
- Better debugging experience

**Faster Slices**

- Expanded stack allocation for slice backing stores
- Improved slice performance

#### Linker

- New `-funcalign=N` option specifies function entry alignment

#### Standard Library Highlights

**New Packages**

1. **`testing/synctest`** (Promoted from Experimental)
   - Concurrent code testing with virtualized time
   - Control time progression in tests
   - Essential for testing time-dependent concurrent code

2. **`encoding/json/v2`** (Experimental)
   - **Substantially better decoding performance**
   - Improved API design
   - Backward compatible with v1

**Major Package Updates**

| Package | Key Changes |
|---------|------------|
| `crypto` | New `MessageSigner` interface and `SignMessage` function |
| `crypto/ecdsa` | New raw key parsing/serialization functions |
| `crypto/rsa` | **Key generation now 3x faster** |
| `crypto/sha1` | **Hashing 2x faster on amd64 with SHA-NI** |
| `crypto/tls` | New `CurveID` field; SHA-1 algorithms disallowed in TLS 1.2 |
| `net` | Windows now supports file-to-connection conversion; IPv6 multicast improvements |
| `net/http` | **New `CrossOriginProtection` middleware for CSRF defense** |
| `os` | Windows async I/O support; `Root` type expanded with 12 new methods |
| `sync` | **New `WaitGroup.Go` method for convenient goroutine creation** |
| `testing` | New `Attr`, `Output` methods; `AllocsPerRun` panics with parallel tests |
| `unique` | More eager and parallel reclamation of interned values |

#### Performance Notes

**Performance Improvements**

- ECDSA and Ed25519 signing **4x faster** in FIPS 140-3 mode
- SHA3 hashing **2x faster** on Apple M processors
- AMD64 fused multiply-add instructions in v3+ mode
  - ⚠️ **Changes floating-point results** (within IEEE 754 spec)

**Performance Regressions**

- SHA-1, SHA-256, SHA-512 slower without AVX2
  - Most servers post-2015 support AVX2

#### Platform Changes

- **macOS**: Requires version 12 Monterey or later
- **Windows**: 32-bit windows/arm port marked for removal in Go 1.26
- **Loong64**: Race detector now supported
- **RISC-V**:
  - Plugin build mode support
  - New `GORISCV64=rva23u64` environment variable value

#### Deprecations

- `go/ast` functions: `FilterPackage`, `PackageExports`, `MergePackageFiles`
- `go/parser.ParseDir` function
- Old `testing/synctest` API (when `GOEXPERIMENT=synctest` set)

---

### Actionable Recommendations for This Project

#### Immediate Opportunities

1. **Replace `sync.WaitGroup` patterns with `WaitGroup.Go()`** (Go 1.25)

   ```go
   // Old pattern
   wg.Add(1)
   go func() {
       defer wg.Done()
       // work
   }()
   
   // New pattern (Go 1.25)
   wg.Go(func() {
       // work
   })
   ```

2. **Use `testing.B.Loop()` in benchmarks** (Go 1.24)

   ```go
   // Old pattern
   for i := 0; i < b.N; i++ {
       // benchmark code
   }
   
   // New pattern (Go 1.24)
   for b.Loop() {
       // benchmark code
   }
   ```

3. **Consider `os.Root` for file operations** (Go 1.24)
   - Prevents path traversal vulnerabilities
   - Safer for user-provided file paths

4. **Enable Green Tea GC for testing** (Go 1.25)
   - Test with `GOEXPERIMENT=greenteagc`
   - May reduce GC overhead by 10-40%

5. **Leverage container-aware GOMAXPROCS** (Go 1.25)
   - No changes needed; automatic in containers
   - Improves resource utilization

6. **Review floating-point operations** (Go 1.25)
   - AMD64 v3+ uses FMA instructions
   - May change floating-point results (within spec)

7. **Watch nil pointer checks** (Go 1.25)
   - Compiler bug fix may expose latent nil pointer bugs
   - Review crash reports carefully

#### Future Considerations

1. **Evaluate `encoding/json/v2`** when stable
   - Better performance for JSON operations
   - Currently experimental in Go 1.25

2. **Adopt tool directives** in go.mod
   - Cleaner dependency management for build tools
   - Remove `tools.go` if present

3. **Enable FIPS mode if required**
   - Use `GOFIPS140=1` for compliance
   - Performance improvements in Go 1.25 (4x faster signing)

4. **Use `runtime/trace.FlightRecorder`** for production debugging
   - Capture traces of rare events
   - Minimal overhead when not triggered
