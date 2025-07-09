# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a high-performance rule engine written in Go that evaluates context-based rules in the form of "x eq 10 and y gt 20". The engine is designed for extreme performance with strict constraints:

- Rule evaluation MUST complete in under 100 nanoseconds (updated from 1000ns)
- Rule evaluation MUST NOT allocate memory during runtime (0 allocs/op required)
- The evaluator can pre-allocate memory during initialization for runtime optimization
- The engine must be 100+ times faster than the original nikunjy/rules library
- Single unified implementation - no separate optimized/unoptimized versions

## Architecture

The project implements a custom zero-allocation rule engine that is API-compatible with the `github.com/nikunjy/rules` library. The main interface is:

```go
// Create one engine instance and reuse it
engine := rule.NewEngine()
result, err := engine.Evaluate(query string, context map[string]any) (bool, error)
```

### Test-Driven Development

The entire specification is defined through comprehensive test cases in `test/fixtures_test.go`. The rule engine must pass all these tests, which cover:

- **Equality & Inequality**: `eq`, `ne`, `==`, `!=` operations
- **Relational**: `lt`, `gt`, `le`, `ge` operations  
- **String Operations**: `co` (contains), `sw` (starts with), `ew` (ends with)
- **Membership**: `in` operator with arrays
- **Presence**: `pr` operator to check if attribute exists
- **Logical**: `not`, `and`, `or` with proper nesting
- **Attribute Comparisons**: Both flat and nested property comparisons
- **Nested Attributes**: Deep object navigation with dot notation

### Rule Syntax Examples

```
x eq 10                           // equality
score gt 100 and level lt 5       // logical operations
city co "York"                    // string contains
color in ["red","green","blue"]   // membership
user.profile.age ge 18            // nested attributes
not (status eq "inactive")        // negation
```

## Development Commands

### Running Tests
```bash
go test ./test/...
```

### Running Specific Tests
```bash
go test ./test -run TestRulesRound1
```

### Build
```bash
go build ./...
```

### Module Management
```bash
go mod tidy
go mod download
```

## Development Guidelines

### Code Standards
- Use modern Go constructs (any, range patterns, etc.)
- Avoid `fmt.Sprintf` when possible
- Use typed errors only - no `errorf` or `error.New`
- Optimize for speed over convenience
- Replace locks with lock-free structures (xsync v4 or uber's atomic)
- Follow TDD: RED → GREEN → REFACTOR

### Performance Requirements
- Memory allocation during rule evaluation is strictly forbidden (0 allocs/op)
- All pre-computation should happen during initialization
- Target sub-100 nanosecond evaluation times (updated from 1000ns)
- Must be 100+ times faster than nikunjy/rules library
- Weight any memory allocation against performance impact
- Pre-allocated EvalResult structures to avoid interface boxing

### Dependencies
- Benchmarking: `github.com/nikunjy/rules v1.5.0` (for performance comparison)
- Concurrency: `github.com/puzpuzpuz/xsync/v4` (lock-free map for rule caching)
- Go version: 1.24.3

## Testing Strategy

All functionality is validated through the comprehensive test suite in `test/fixtures_test.go`. The test cases define the complete specification - any new features must be validated against these tests. The tests are organized by operation type and include edge cases for nested attributes and complex logical expressions.

### Type System Compliance
- **Strict Type Checking**: Different categories (string/number/boolean) never compare equal
- **Numeric Cross-Type**: int/float comparisons are allowed (42 == 42.0)
- **String Comparisons**: Lexicographic ordering for string relational operations
- **Membership Operations**: Use strict type checking (no cross-type matching)
- **Large Integer Support**: Preserve precision for integers > 2^53 using dual storage

### Zero-Allocation Implementation
- **EvalResult Structure**: Pre-allocated typed result structure to avoid interface boxing
- **Memory Reuse**: Single evaluator instance with reusable result buffer
- **AST Caching**: Pre-compiled rules stored in lock-free concurrent map
- **Allocation Verification**: All benchmarks must show 0 allocs/op

### Performance Benchmarking
- **Baseline Comparison**: Must outperform nikunjy/rules by 100x minimum
- **Sub-100ns Requirement**: All evaluations must complete in under 100 nanoseconds
- **Memory Efficiency**: Zero allocations during evaluation phase
- **Concurrency Safe**: Thread-safe rule compilation and evaluation

### Test Development
- **Package Isolation**: Any custom test scripts MUST be placed in their own separate folder (not in the root directory) to avoid conflicts with the main `rule` package
- **Main Tests**: Use the existing `test/` directory structure for all official functionality tests
- **Debug Scripts**: For debugging, create a separate folder (e.g., `debug/`) with a `main` function and run with `go run ./debug/`. Do NOT create separate modules for each test function.