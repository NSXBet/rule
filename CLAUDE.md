# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a high-performance rule engine written in Go that evaluates context-based rules in the form of "x eq 10 and y gt 20". The engine is designed for extreme performance with strict constraints:

- Rule evaluation MUST complete in under 1000 nanoseconds
- Rule evaluation MUST NOT allocate memory during runtime
- The evaluator can pre-allocate memory during initialization for runtime optimization

## Architecture

The project uses the `github.com/nikunjy/rules` library as the core rule evaluation engine. The main interface is:

```go
rules.Evaluate(query string, context map[string]any) (bool, error)
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
- Memory allocation during rule evaluation is strictly forbidden
- All pre-computation should happen during initialization
- Target sub-1000 nanosecond evaluation times
- Weight any memory allocation against performance impact

### Dependencies
- Core: `github.com/nikunjy/rules v1.5.0`
- Testing: `github.com/stretchr/testify v1.10.0`
- Go version: 1.24.3

## Testing Strategy

All functionality is validated through the comprehensive test suite in `test/fixtures_test.go`. The test cases define the complete specification - any new features must be validated against these tests. The tests are organized by operation type and include edge cases for nested attributes and complex logical expressions.