# Testing & Quality

**Author:** KleaSCM  
**Email:** KleaSCM@gmail.com

---

## Overview

Steria is built with a strong emphasis on reliability, correctness, and maintainability. The project includes:
- Comprehensive unit tests for all core modules
- Integration tests for end-to-end CLI workflows
- Performance benchmarks for critical operations
- Plans for CI automation and code coverage reporting

## How to Run All Tests

```sh
go test ./... -v
```

## How to Run Benchmarks

```sh
go test ./internal/storage/ -v -bench=.
go test ./internal/security/ -v -bench=.
```

## Interpreting Results
- All tests should pass with no errors.
- Benchmarks report average operation time (ns/op) for key functions.
- See `Tests/README.md` for the latest results and detailed output.

## Coverage & CI
- Code coverage reporting and CI automation are planned for future releases.
- CI will run all tests and benchmarks on every push/PR.

## Where to Find Test Details
- See [`Tests/README.md`](../Tests/README.md) for:
  - Test system structure
  - Latest results
  - How to contribute new tests

---

For more on contributing, see the [Contributing Guide](contributing.md). 