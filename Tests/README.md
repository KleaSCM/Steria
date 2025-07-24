# Author: KleaSCM
# Email: KleaSCM@gmail.com
# Name of the file: README.md
# Description: Test documentation and instructions for the Steria project.

---

# Steria Project Tests

This directory contains all test cases and testing documentation for the project. 

## Test System Overview
- **Unit tests** are provided for all core modules in `internal/`, `core/`, and utility packages.
- **Integration tests** are provided in `Tests/integration_test.go` to verify end-to-end CLI workflows.
- All tests use Go's standard `testing` package.

## How to Run All Tests

```sh
go test ./... -v
```

## How to Run Storage Benchmarks

```sh
go test ./internal/storage/ -v -bench=.
```

## How to Run Security Benchmarks

```sh
go test ./internal/security/ -v -bench=.
```

## Latest Test Results

### Storage Module
```
=== RUN   TestLoadOrInitRepo_New
--- PASS: TestLoadOrInitRepo_New (0.00s)
=== RUN   TestLoadOrInitRepo_Existing
--- PASS: TestLoadOrInitRepo_Existing (0.00s)
=== RUN   TestCreateCommitAndGetChanges
--- PASS: TestCreateCommitAndGetChanges (0.00s)
=== RUN   TestHasRemote
--- PASS: TestHasRemote (0.00s)
=== RUN   TestSync_NoRemote
--- PASS: TestSync_NoRemote (0.00s)
=== RUN   TestLoadCommit
--- PASS: TestLoadCommit (0.00s)
=== RUN   TestGetCurrentStateAndWorkingState
--- PASS: TestGetCurrentStateAndWorkingState (0.00s)
=== RUN   TestCalculateFileHash
--- PASS: TestCalculateFileHash (0.00s)
goos: linux
goarch: amd64
pkg: steria/internal/storage
cpu: AMD Ryzen 7 8845HS w/ Radeon 780M Graphics     
BenchmarkCreateCommit
BenchmarkCreateCommit-16           16938             68681 ns/op
PASS
ok      steria/internal/storage 1.987s
```

### Security Module
```
=== RUN   TestGenerateKeyPair
--- PASS: TestGenerateKeyPair (0.00s)
=== RUN   TestSignAndVerifyMessage
--- PASS: TestSignAndVerifyMessage (0.00s)
=== RUN   TestVerifySignature_Invalid
--- PASS: TestVerifySignature_Invalid (0.00s)
=== RUN   TestSecureHashAndVerify
--- PASS: TestSecureHashAndVerify (0.00s)
=== RUN   TestVerifySecureHash_Invalid
--- PASS: TestVerifySecureHash_Invalid (0.00s)
goos: linux
goarch: amd64
pkg: steria/internal/security
cpu: AMD Ryzen 7 8845HS w/ Radeon 780M Graphics     
BenchmarkGenerateKeyPair
BenchmarkGenerateKeyPair-16        78223             14537 ns/op
PASS
ok   steria/internal/security 1.297s
```

### Metrics Module
```
=== RUN   TestStartProfiling
--- PASS: TestStartProfiling (0.00s)
PASS
ok   steria/internal/metrics	0.002s
```

### Utils Module
```
=== RUN   TestShouldIgnore
--- PASS: TestShouldIgnore (0.00s)
PASS
ok   steria/internal/utils	0.002s
```

### Core Module
```
=== RUN   TestCreateCommit
--- PASS: TestCreateCommit (0.00s)
PASS
ok   steria/core	0.002s
```

## Summary
- **All unit and integration tests passed** for core, storage, metrics, security, utils, and CLI workflow.
- **Benchmark for CreateCommit:** ~68,681 ns/op (16 threads)
- **Benchmark for GenerateKeyPair:** ~14,537 ns/op (16 threads)
- No test files yet for CLI command packages (`cmd/`).

---

For full coverage, see TODOs in each `_test.go` file. Contributions and additional tests are welcome! 