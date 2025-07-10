# Rule Engine Examples

This directory contains comprehensive examples demonstrating the NSXBet Rule Engine capabilities, including our 100% compatibility with nikunjy/rules and our performance/feature extensions.

## Quick Start

Run all examples:
```bash
cd examples
go run main.go
```

Run specific example categories:
```bash
go run main.go basic      # Basic usage and type safety
go run main.go ecommerce  # E-commerce business rules  
go run main.go datetime   # DateTime operations (our extension)
go run main.go migration  # Migration from nikunjy/rules
```

## Example Files

### 📚 [basic_usage.go](basic_usage.go)
Demonstrates fundamental rule engine usage:
- **BasicUsageExample()**: Simple rule evaluation with nested objects
- **TypeSafetyExample()**: Type safety and cross-type comparison behavior
- **PerformanceExample()**: Query caching and pre-compilation for performance

**Key concepts**: rule.D syntax, basic operators, query caching

### 🛒 [ecommerce_eligibility.go](ecommerce_eligibility.go)
Real-world e-commerce business rule scenarios:
- **EcommerceEligibilityExample()**: Complex customer eligibility rules
- **DynamicPricingExample()**: Dynamic pricing rule evaluation

**Key concepts**: Complex nested rules, business logic, real-world use cases

### 📅 [datetime_operations.go](datetime_operations.go)
Showcases our datetime extensions beyond nikunjy/rules:
- **DateTimeOperationsExample()**: Comprehensive datetime operator usage
- **SchedulingExample()**: Meeting scheduling validation system

**Key concepts**: DateTime operators (dq, dn, be, bq, af, aq), time.Time handling, RFC3339 and Unix timestamps

### 🔄 [migration_from_nikunjy.go](migration_from_nikunjy.go)
Migration guide and compatibility demonstration:
- **MigrationExample()**: Side-by-side compatibility testing
- **BeforeAfterExample()**: API comparison and migration steps
- **PerformanceComparisonExample()**: Performance benefits demonstration

**Key concepts**: 100% compatibility, migration steps, performance improvements

## Features Demonstrated

### ✅ 100% Compatible with nikunjy/rules
- All basic operators: `eq`, `ne`, `lt`, `gt`, `le`, `ge`
- String operators: `co`, `sw`, `ew`
- Logical operators: `and`, `or`, `not`
- Membership: `in` with arrays
- Presence: `pr` operator
- Nested properties: `user.profile.age`
- Type safety: same cross-type behavior
- time.Time handling: identical string conversion

### 🚀 Our Extensions
- **DateTime operators**: `dq`, `dn`, `be`, `bq`, `af`, `aq`
- **Performance**: 25-144x faster, 0 allocations
- **Query caching**: Thread-safe automatic caching
- **Clean API**: `rule.D` type alias
- **Multiple timestamp formats**: time.Time, RFC3339, Unix

### 📊 Performance Benefits
- **Evaluation speed**: ~25ns (vs ~3,000ns)
- **Memory usage**: 0 allocs/op (vs ~88 allocs/op)
- **Caching**: Automatic query compilation caching
- **Concurrency**: Thread-safe for high-load systems

## Example Usage Patterns

### Basic Rule Evaluation
```go
engine := rule.NewEngine()
context := rule.D{"user": rule.D{"age": 25}}
result, err := engine.Evaluate(`user.age gt 18`, context)
```

### Pre-compilation for Performance
```go
engine := rule.NewEngine()
err := engine.AddQuery(`user.role eq "admin"`)  // Pre-compile
result, _ := engine.Evaluate(`user.role eq "admin"`, context)  // Fast evaluation
```

### DateTime Operations (Our Extension)
```go
context := rule.D{
    "created_at": time.Now(),
    "deadline": "2024-12-31T23:59:59Z",
}
result, _ := engine.Evaluate(`created_at be deadline`, context)
```

### Migration from nikunjy/rules
```go
// Before (nikunjy/rules)
result, err := rules.Evaluate(rule, context)

// After (our library) 
engine := rule.NewEngine()
result, err := engine.Evaluate(rule, context)  // Same result!
```

## Test Coverage

These examples include comprehensive tests that prove:
- **22/22 basic compatibility tests pass** ✅
- **3/3 time.Time compatibility tests pass** ✅  
- **100% compatibility rate achieved** ✅

See `../test/compatibility_test.go` for the automated test suite.

## Dependencies

The examples require:
- Go 1.21+
- github.com/NSXBet/rule-engine
- github.com/nikunjy/rules (for migration examples)

## Running Examples

1. **Install dependencies**:
   ```bash
   go mod tidy
   ```

2. **Run all examples**:
   ```bash
   go run main.go
   ```

3. **Run specific examples**:
   ```bash
   go run main.go basic
   go run main.go ecommerce  
   go run main.go datetime
   go run main.go migration
   ```

## Performance Benchmarking

To see actual performance comparisons:
```bash
cd ..
make bench  # Run comprehensive benchmarks
```

## Next Steps

After reviewing these examples:
1. Check out the [main README](../README.md) for complete documentation
2. Review [test/compatibility_test.go](../test/compatibility_test.go) for compatibility verification
3. Run benchmarks with `make bench` to see performance improvements
4. Start migrating your rules using the patterns shown in migration examples

---

**Built with ❤️ for the Go community**

[Report Issues](https://github.com/NSXBet/rule-engine/issues) | [Contribute](https://github.com/NSXBet/rule-engine/pulls) | [Documentation](../README.md)