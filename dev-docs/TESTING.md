# Testing Guide

This guide explains how to run and maintain the test suite for the Flashpipe CLI project.

## Table of Contents

- [Quick Start](#quick-start)
- [Running Tests](#running-tests)
- [Test Coverage](#test-coverage)
- [Writing New Tests](#writing-new-tests)
- [Test Organization](#test-organization)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)

---

## Quick Start

### Prerequisites

- Go 1.21 or higher
- Git (for downloading dependencies)

### Install Dependencies

```bash
go mod download
```

### Run All Tests

```bash
go test ./...
```

### Run Tests with Coverage

```bash
go test ./... -cover
```

---

## Running Tests

### Run Tests for Specific Package

```bash
# Partner Directory repository layer
go test ./internal/repo -v

# Configuration loader
go test ./internal/deploy -v

# API layer
go test ./internal/api -v
```

### Run Specific Test Function

```bash
# Run a single test
go test ./internal/repo -run TestParseContentType

# Run tests matching a pattern
go test ./internal/repo -run TestBinary
```

### Run with Verbose Output

```bash
go test ./internal/repo -v
```

### Run with Coverage Report

```bash
# Generate coverage report
go test ./internal/repo -coverprofile=coverage.out

# View coverage in terminal
go tool cover -func=coverage.out

# View coverage in browser
go tool cover -html=coverage.out
```

### Run with Race Detection

```bash
go test ./internal/repo ./internal/deploy -race
```

### Run Only Short Tests (Skip Integration Tests)

```bash
go test ./... -short
```

---

## Test Coverage

### Current Coverage by Package

| Package | Coverage | Status |
|---------|----------|--------|
| `internal/repo` | 74.9% | ✅ Good |
| `internal/deploy` | 82.6% | ✅ Excellent |
| `internal/analytics` | 42.9% | ⚠️ Moderate |
| `internal/str` | 35.0% | ⚠️ Low |
| `internal/file` | 5.3% | 🔴 Needs Work |
| `internal/sync` | 3.4% | 🔴 Needs Work |

### Generate Coverage Report for All Packages

```bash
# Create coverage directory
mkdir -p coverage

# Generate coverage for each package
go test ./internal/repo -coverprofile=coverage/repo.out
go test ./internal/deploy -coverprofile=coverage/deploy.out

# View combined report
go tool cover -html=coverage/repo.out
```

### Coverage Goals

- **Critical paths:** >90% coverage
- **New features:** >80% coverage
- **Overall project:** >70% coverage

---

## Writing New Tests

### Test File Naming

- Test files must end with `_test.go`
- Place test files in the same package as the code being tested
- Example: `partnerdirectory.go` → `partnerdirectory_test.go`

### Test Function Naming

```go
func TestFunctionName(t *testing.T)               // Basic test
func TestFunctionName_Scenario(t *testing.T)      // Specific scenario
func TestFunctionName_EdgeCase(t *testing.T)      // Edge case
```

### Test Structure

```go
package mypackage

import (
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestMyFunction(t *testing.T) {
    // Setup
    input := "test input"
    expected := "expected output"
    
    // Execute
    result := MyFunction(input)
    
    // Assert
    assert.Equal(t, expected, result)
}
```

### Table-Driven Tests

```go
func TestParseContentType(t *testing.T) {
    tests := []struct {
        name        string
        input       string
        wantExt     string
        wantError   bool
    }{
        {
            name:      "simple xml",
            input:     "xml",
            wantExt:   "xml",
            wantError: false,
        },
        {
            name:      "with encoding",
            input:     "xml; encoding=UTF-8",
            wantExt:   "xml",
            wantError: false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ext, err := parseContentType(tt.input)
            
            if tt.wantError {
                assert.Error(t, err)
                return
            }
            
            require.NoError(t, err)
            assert.Equal(t, tt.wantExt, ext)
        })
    }
}
```

### Testing with Temporary Files

```go
func TestFileOperation(t *testing.T) {
    // Create temp directory
    tempDir, err := os.MkdirTemp("", "test-*")
    require.NoError(t, err)
    defer os.RemoveAll(tempDir) // Clean up
    
    // Create test file
    testFile := filepath.Join(tempDir, "test.txt")
    err = os.WriteFile(testFile, []byte("content"), 0644)
    require.NoError(t, err)
    
    // Run test
    result := ProcessFile(testFile)
    
    // Verify
    assert.True(t, result)
}
```

### Using `require` vs `assert`

```go
// Use require for fatal errors (stop test execution)
require.NoError(t, err)
require.NotNil(t, result)
require.Len(t, items, 5)

// Use assert for non-fatal assertions (continue test execution)
assert.Equal(t, expected, actual)
assert.Contains(t, str, substring)
assert.True(t, condition)
```

---

## Test Organization

### Directory Structure

```
internal/
├── repo/
│   ├── partnerdirectory.go
│   └── partnerdirectory_test.go       (708 lines, 25 tests)
├── deploy/
│   ├── config_loader.go
│   ├── config_loader_test.go          (556 lines, 20 tests)
│   ├── utils.go
│   └── utils_test.go                  (562 lines, 18 tests)
└── api/
    ├── partnerdirectory.go
    └── partnerdirectory_test.go
```

### Test Categories

1. **Unit Tests** - Test individual functions in isolation
2. **Integration Tests** - Test interaction between components
3. **End-to-End Tests** - Test complete workflows

---

## Best Practices

### DO ✅

- Write tests for new code before submitting PR
- Use descriptive test names that explain what is being tested
- Test both happy paths and error cases
- Clean up resources (files, connections) with `defer`
- Use table-driven tests for multiple scenarios
- Keep tests independent (no shared state)
- Mock external dependencies (HTTP, database, file system when appropriate)

### DON'T ❌

- Commit tests that require manual intervention
- Write tests that depend on external services (use mocks)
- Write flaky tests (random failures)
- Share state between tests
- Test implementation details (test behavior, not internals)
- Write overly complex tests (keep them simple)

### Code Coverage Guidelines

- Aim for **>80% coverage** for critical code paths
- Don't obsess over 100% coverage
- Focus on testing **important logic** and **edge cases**
- Skip trivial getters/setters
- Document any intentionally uncovered code

### Test Maintenance

- Update tests when changing code behavior
- Remove obsolete tests for removed features
- Refactor tests to reduce duplication
- Keep test code as clean as production code

---

## Troubleshooting

### Tests Fail on Windows but Pass on Linux

**Issue:** Line ending differences (CRLF vs LF)

**Solution:** Tests already handle this by:
```go
// Detect line ending style
lineEnding := "\n"
if strings.Contains(string(data), "\r\n") {
    lineEnding = "\r\n"
}
```

### Tests Are Slow

**Causes:**
- Too many file I/O operations
- Network calls (should be mocked)
- Large test data

**Solutions:**
```bash
# Run only fast tests
go test ./... -short

# Run tests in parallel
go test ./... -parallel 4

# Profile slow tests
go test ./... -cpuprofile=cpu.prof
go tool pprof cpu.prof
```

### Coverage Report Shows Uncovered Lines

**Check:**
1. Are there error paths not tested?
2. Is the code actually reachable?
3. Should this code be tested, or is it trivial?

**Example:**
```go
// Intentionally uncovered - OS-specific error handling
if runtime.GOOS == "windows" {
    // Windows-specific path (hard to test cross-platform)
}
```

### Test Fixtures Are Missing

**Issue:** Test data files not found

**Solution:** Use relative paths from test file location:
```go
testDataPath := filepath.Join("testdata", "config.yml")
```

### Race Conditions Detected

**Issue:** `go test -race` reports data races

**Solution:**
1. Identify shared state
2. Add proper synchronization (mutex, channels)
3. Make tests independent

---

## Continuous Integration

### GitHub Actions Example

```yaml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Run tests
        run: go test ./... -race -coverprofile=coverage.out
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage.out
```

---

## Additional Resources

- [Go Testing Documentation](https://pkg.go.dev/testing)
- [Testify Documentation](https://github.com/stretchr/testify)
- [Table-Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests)
- [Test Coverage Summary](./TEST_COVERAGE_SUMMARY.md)

---

## Getting Help

If you have questions about:
- Writing tests → See "Writing New Tests" section above
- Running tests → See "Running Tests" section above
- Coverage goals → See [TEST_COVERAGE_SUMMARY.md](./TEST_COVERAGE_SUMMARY.md)
- Test failures → Check existing tests for examples

---

**Last Updated:** December 2024  
**Maintainer:** Development Team