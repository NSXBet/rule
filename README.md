# Rule Engine 🚀

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Tests](https://img.shields.io/badge/Tests-248%20Passing-brightgreen.svg)](#)
[![Lint](https://img.shields.io/badge/Lint-100%25%20Clean-brightgreen.svg)](#)

A **blazingly fast**, **zero-allocation** rule engine for Go that evaluates logical expressions in under 100 nanoseconds ⚡

## 📖 Table of Contents

1. [🚀 Getting Started](#-getting-started)
2. [📚 API](#-api)
3. [🎯 Context](#-context)
4. [🔤 Rule Language](#-rule-language)
5. [💾 Query Caching](#-query-caching)
6. [⚡ Benchmarks](#-benchmarks)
7. [🤝 Contributing](#-contributing)
8. [📄 License](#-license)

---

## 🚀 Getting Started

### Installation

```bash
go get github.com/NSXBet/rule-engine
```

### Quick Example

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/NSXBet/rule-engine"
)

func main() {
    // Create a new rule engine
    engine := rule.NewEngine()
    
    // Define your context data
    context := rule.D{
        "user": rule.D{
            "age":    25,
            "status": "active",
            "name":   "John Doe",
        },
    }
    
    // Evaluate a rule
    result, err := engine.Evaluate(`user.age gt 18 and user.status eq "active"`, context)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("User can access: %t\n", result) // Output: User can access: true
}
```

That's it! 🎉 You're now evaluating complex business rules in under 100 nanoseconds with zero memory allocations.

---

## 📚 API

The rule engine provides a simple but powerful API:

### Core Methods

#### `NewEngine() *Engine`
Creates a new rule engine instance. Each engine maintains its own query cache for optimal performance.

```go
engine := rule.NewEngine()
```

#### `Evaluate(query string, context rule.D) (bool, error)`
Evaluates a rule expression against the provided context. Returns `true`/`false` and any parsing/evaluation errors.

```go
result, err := engine.Evaluate(`price lt 100`, rule.D{"price": 50})
// result: true, err: nil
```

#### `AddQuery(query string) error` 
Pre-compiles and caches a query for optimal performance. This is optional but recommended for frequently used rules.

```go
// Pre-compile for better performance
err := engine.AddQuery(`user.role eq "admin"`)
if err != nil {
    log.Fatal(err)
}

// Later evaluations will be faster
result, _ := engine.Evaluate(`user.role eq "admin"`, context)
```

### Error Handling

The engine returns descriptive errors for invalid syntax:

```go
result, err := engine.Evaluate("invalid syntax !!!", context)
if err != nil {
    fmt.Printf("Parse error: %v\n", err)
    // Output: Parse error: unexpected token at position 15
}
```

---

## 🎯 Context

The context is a `rule.D` (alias for `map[string]any`) that contains your data. The engine supports arbitrarily nested structures! 🏗️

### Simple Values

```go
context := rule.D{
    "price":    99.99,
    "quantity": 5,
    "active":   true,
    "name":     "Product A",
}

// Use directly in rules
engine.Evaluate(`price lt 100`, context)      // true
engine.Evaluate(`quantity ge 3`, context)     // true  
engine.Evaluate(`active eq true`, context)    // true
engine.Evaluate(`name co "Product"`, context) // true
```

### Nested Objects

Navigate deep object hierarchies with dot notation:

```go
context := rule.D{
    "user": rule.D{
        "profile": rule.D{
            "settings": rule.D{
                "theme":         "dark",
                "notifications": true,
            },
            "preferences": rule.D{
                "language": "en",
                "timezone": "UTC",
            },
        },
        "subscription": rule.D{
            "plan":   "premium",
            "active": true,
        },
    },
}

// Navigate nested structures easily
engine.Evaluate(`user.profile.settings.theme eq "dark"`, context)              // true
engine.Evaluate(`user.subscription.plan eq "premium"`, context)               // true  
engine.Evaluate(`user.profile.preferences.language eq "en"`, context)         // true
```

### Arrays and Membership

Check if values exist in arrays:

```go
context := rule.D{
    "user": rule.D{
        "roles": []any{"admin", "moderator"},
        "tags":  []any{"vip", "beta-tester"},
    },
    "colors": []any{"red", "green", "blue"},
}

engine.Evaluate(`user.roles in ["admin", "user"]`, context)     // true (admin matches)
engine.Evaluate(`"red" in colors`, context)                      // true
engine.Evaluate(`user.tags in ["vip", "premium"]`, context)    // true (vip matches)
```

### Type Safety

The engine handles type mismatches gracefully - different types never compare as equal:

```go
context := rule.D{
    "count":  42,
    "flag":   true,
    "text":   "42",
}

engine.Evaluate(`count eq 42`, context)     // true (int matches int)
engine.Evaluate(`count eq "42"`, context) // false (int != string)
engine.Evaluate(`flag eq 1`, context)       // false (bool != int)
engine.Evaluate(`text eq "42"`, context)  // true (string matches string)
```

---

## 🔤 Rule Language

Our rule language is intuitive and powerful! Here are all the supported operators:

### Equality Operators

| Operator | Description | Example | Result |
|----------|-------------|---------|---------|
| `eq` | Equal to | `age eq 25` | `true` if age is 25 |
| `ne` | Not equal to | `status ne "inactive"` | `true` if status is not "inactive" |
| `==` | Equal to (alias) | `price == 99.99` | Same as `eq` |
| `!=` | Not equal to (alias) | `role != "guest"` | Same as `ne` |

```go
context := rule.D{"age": 25, "status": "active"}

engine.Evaluate(`age eq 25`, context)               // true
engine.Evaluate(`status ne "inactive"`, context)  // true
engine.Evaluate(`age == 25`, context)               // true  
engine.Evaluate(`status != "guest"`, context)     // true
```

### Relational Operators

| Operator | Description | Example | Result |
|----------|-------------|---------|---------|
| `lt` | Less than | `price lt 100` | `true` if price < 100 |
| `gt` | Greater than | `score gt 80` | `true` if score > 80 |
| `le` | Less than or equal | `age le 18` | `true` if age ≤ 18 |
| `ge` | Greater than or equal | `rating ge 4.5` | `true` if rating ≥ 4.5 |

```go
context := rule.D{"price": 50, "score": 95, "age": 16, "rating": 4.8}

engine.Evaluate(`price lt 100`, context)   // true
engine.Evaluate(`score gt 80`, context)    // true
engine.Evaluate(`age le 18`, context)      // true
engine.Evaluate(`rating ge 4.5`, context) // true
```

### String Operators

| Operator | Description | Example | Result |
|----------|-------------|---------|---------|
| `co` | Contains | `name co "John"` | `true` if name contains "John" |
| `sw` | Starts with | `email sw "admin"` | `true` if email starts with "admin" |
| `ew` | Ends with | `domain ew ".com"` | `true` if domain ends with ".com" |

```go
context := rule.D{
    "name":   "John Doe",
    "email":  "admin@company.com", 
    "domain": "example.com",
}

engine.Evaluate(`name co "John"`, context)        // true
engine.Evaluate(`email sw "admin"`, context)      // true  
engine.Evaluate(`domain ew ".com"`, context)      // true
```

### Membership Operator

| Operator | Description | Example | Result |
|----------|-------------|---------|---------|
| `in` | Member of array | `role in ["admin", "mod"]` | `true` if role is "admin" or "mod" |

```go
context := rule.D{
    "role":   "admin",
    "colors": []any{"red", "green"},
}

engine.Evaluate(`role in ["admin", "user"]`, context)        // true
engine.Evaluate(`"red" in colors`, context)                   // true
engine.Evaluate(`"blue" in ["red", "green", "blue"]`, context) // true
```

### Presence Operator

| Operator | Description | Example | Result |
|----------|-------------|---------|---------|
| `pr` | Property present | `user.email pr` | `true` if user.email exists |

```go
context := rule.D{
    "user": rule.D{
        "name":  "John",
        "email": "john@example.com",
    },
}

engine.Evaluate(`user.email pr`, context)      // true (email exists)
engine.Evaluate(`user.phone pr`, context)      // false (phone doesn't exist)
engine.Evaluate(`user.name pr`, context)       // true (name exists)
```

### Logical Operators

| Operator | Description | Example | Result |
|----------|-------------|---------|---------|
| `and` | Logical AND | `age gt 18 and status eq "active"` | `true` if both conditions are true |
| `or` | Logical OR | `role eq "admin" or role eq "mod"` | `true` if either condition is true |
| `not` | Logical NOT | `not (age lt 18)` | `true` if age is NOT less than 18 |

```go
context := rule.D{"age": 25, "status": "active", "role": "admin"}

engine.Evaluate(`age gt 18 and status eq "active"`, context)  // true
engine.Evaluate(`role eq "admin" or role eq "mod"`, context) // true
engine.Evaluate(`not (age lt 18)`, context)                     // true
```

### DateTime Operators 📅

Perfect for time-based rules and scheduling logic:

| Operator | Description | Example | Result |
|----------|-------------|---------|---------|
| `dq` | DateTime equal | `created_at dq "2024-01-01T10:00:00Z"` | `true` if timestamps are equal |
| `dn` | DateTime not equal | `updated_at dn "2024-01-01T10:00:00Z"` | `true` if timestamps differ |
| `be` | Before | `start_time be "2024-12-31T23:59:59Z"` | `true` if start_time is before |
| `bq` | Before or equal | `deadline bq "2024-12-31T23:59:59Z"` | `true` if deadline is before or equal |
| `af` | After | `event_time af "2024-01-01T00:00:00Z"` | `true` if event_time is after |
| `aq` | After or equal | `publish_date aq "2024-01-01T00:00:00Z"` | `true` if publish_date is after or equal |

Supports both **RFC3339** strings and **Unix timestamps**:

```go
context := rule.D{
    "created_at":   "2024-07-09T22:12:00Z",           // RFC3339
    "updated_at":   int64(1720558320),                // Unix timestamp  
    "publish_date": "2024-01-15T10:30:00-03:00",     // RFC3339 with timezone
}

engine.Evaluate(`created_at af "2024-01-01T00:00:00Z"`, context)     // true
engine.Evaluate(`updated_at be 1720558400`, context)                   // true  
engine.Evaluate(`publish_date aq "2024-01-01T00:00:00Z"`, context)   // true
```

### Complex Expressions

Combine operators with parentheses for complex business logic:

```go
context := rule.D{
    "user": rule.D{
        "age":      25,
        "status":   "active", 
        "role":     "premium",
        "country":  "US",
    },
    "feature_flags": rule.D{
        "beta_enabled": true,
    },
}

// Complex eligibility rule
rule := `(user.age ge 18 and user.status eq "active") and 
         (user.role in ["premium", "enterprise"] or user.country eq "US") and
         feature_flags.beta_enabled eq true`

result, _ := engine.Evaluate(rule, context) // true
```

---

## 💾 Query Caching

The rule engine is smart about performance! 🧠 Here's how caching works:

### Automatic Lazy Caching

Every query gets automatically cached after first use:

```go
engine := rule.NewEngine()

// First evaluation: parses + compiles + caches + evaluates
result1, _ := engine.Evaluate(`user.age gt 18`, context) // ~100ns (includes parsing)

// Subsequent evaluations: uses cached AST
result2, _ := engine.Evaluate(`user.age gt 18`, context) // ~25ns (cached!)
result3, _ := engine.Evaluate(`user.age gt 18`, context) // ~25ns (cached!)
```

### Pre-compilation with AddQuery

For maximum performance, pre-compile frequently used rules:

```go
engine := rule.NewEngine()

// Pre-compile critical business rules at startup
criticalRules := []string{
    `user.role eq "admin"`,
    `user.subscription.active eq true`, 
    `user.age ge 18 and user.status eq "verified"`,
}

for _, rule := range criticalRules {
    if err := engine.AddQuery(rule); err != nil {
        log.Fatalf("Invalid rule: %s - %v", rule, err)
    }
}

// Now all evaluations are lightning fast from the start! ⚡
```

### When to Use AddQuery

✅ **Use AddQuery when:**
- Rules are known at application startup
- You want to validate rule syntax early  
- Maximum performance is critical
- Rules are used frequently (>1000 times)

✅ **Skip AddQuery when:**
- Rules are dynamic/user-generated
- One-time or infrequent evaluations
- Prototyping or development

### Memory Management

The cache is bounded and efficient:
- **Thread-safe**: Multiple goroutines can safely share an engine
- **Memory-efficient**: Only stores parsed AST, not string queries  
- **Bounded growth**: Cache size grows with unique queries only

---

## ⚡ Benchmarks

We believe in **transparency over marketing** 📊. Here are objective performance comparisons to help you choose the right tool:

> 💡 **Disclaimer**: We're not trying to discourage anyone from using `nikunjy/rules` or Go templates - they're excellent libraries with different design goals. We're simply offering another option that might benefit your specific use case, especially when ultra-low latency and zero allocations are critical.

### Performance Results

All benchmarks run on: `Intel(R) Core(TM) i9-14900KF, Go 1.21+`

#### Simple Operations (`x eq 10`)

| Engine | Time/op | Allocs/op | Memory/op | Relative Speed |
|--------|---------|-----------|-----------|----------------|
| **Our Engine** | **24.73 ns** | **0 allocs** | **0 B** | **1x (baseline)** ✅ |
| nikunjy/rules | 3,033 ns | 88 allocs | 5,328 B | 123x slower |
| text/template | 551.0 ns | 14 allocs | 424 B | 22x slower |

#### Complex Operations (`(user.age gt 18 and status eq "active") or user.name co "Admin"`)

| Engine | Time/op | Allocs/op | Memory/op | Relative Speed |
|--------|---------|-----------|-----------|----------------|
| **Our Engine** | **66.38 ns** | **0 allocs** | **0 B** | **1x (baseline)** ✅ |
| nikunjy/rules | 9,589 ns | 190 allocs | 12,905 B | 144x slower |
| text/template | 1,246 ns | 28 allocs | 736 B | 19x slower |

#### String Operations (`name co "John" and email ew ".com"`)

| Engine | Time/op | Allocs/op | Memory/op | Relative Speed |
|--------|---------|-----------|-----------|----------------|
| **Our Engine** | **58.05 ns** | **0 allocs** | **0 B** | **1x (baseline)** ✅ |
| nikunjy/rules | 5,573 ns | 128 allocs | 8,120 B | 96x slower |
| text/template | 871.0 ns | 17 allocs | 424 B | 15x slower |

#### In Operator (`color in ["red", "green", "blue"]`)

| Engine | Time/op | Allocs/op | Memory/op | Relative Speed |
|--------|---------|-----------|-----------|----------------|
| **Our Engine** | **31.92 ns** | **0 allocs** | **0 B** | **1x (baseline)** ✅ |
| nikunjy/rules | 4,664 ns | 106 allocs | 6,648 B | 146x slower |
| text/template | 621.8 ns | 16 allocs | 464 B | 19x slower |

#### Nested Properties (`user.profile.settings.theme eq "dark"`)

| Engine | Time/op | Allocs/op | Memory/op | Relative Speed |
|--------|---------|-----------|-----------|----------------|
| **Our Engine** | **43.34 ns** | **0 allocs** | **0 B** | **1x (baseline)** ✅ |
| nikunjy/rules | 4,793 ns | 108 allocs | 6,824 B | 111x slower |
| text/template | 747.4 ns | 21 allocs | 536 B | 17x slower |

#### Many Queries (5 different queries with pre-compilation)

| Engine | Time/op | Allocs/op | Memory/op | Relative Speed |
|--------|---------|-----------|-----------|----------------|
| **Our Engine** | **144.7 ns** | **0 allocs** | **0 B** | **1x (baseline)** ✅ |
| nikunjy/rules | 20,791 ns | 462 allocs | 30,077 B | 144x slower |
| text/template | 2,791 ns | 62 allocs | 1,800 B | 19x slower |

#### DateTime Operations (`created_at af "2024-01-01T00:00:00Z"`)

| Engine | Time/op | Allocs/op | Memory/op | Relative Speed |
|--------|---------|-----------|-----------|----------------|
| **Our Engine** | **118.9 ns** | **0 allocs** | **0 B** | **1x (baseline)** ✅ |
| nikunjy/rules | ❌ Not supported | ❌ No datetime operators | | |
| text/template | 🔶 Complex setup required | 🔶 Custom functions needed | | |

### 🔬 Run Benchmarks Yourself

Want to verify these results? Run the benchmarks on your hardware:

```bash
# Clone the repository
git clone https://github.com/NSXBet/rule-engine
cd rule-engine

# Run all comparison benchmarks  
make bench

# Or run specific benchmark categories
go test -bench=BenchmarkComparison -benchmem .
go test -bench=BenchmarkDateTime -benchmem .
```

### When Each Tool Shines 🌟

**Our Engine is ideal for:**
- 🚀 Ultra-high performance applications (>100k ops/sec)
- 🎯 Zero-allocation requirements
- ⚡ Sub-100ns latency needs
- 📅 DateTime-heavy business rules
- 🔄 Real-time systems and hot paths

**nikunjy/rules is great for:**
- 🛠️ Rapid prototyping and development
- 📖 Excellent documentation and examples
- 🏗️ Less performance-critical applications
- 👥 Large community and ecosystem

**text/template works well for:**
- 📝 Template generation and formatting
- 🔧 Complex custom function needs
- 🎨 String manipulation and rendering
- 📚 Standard library familiarity

---

## 🔄 Compatibility with nikunjy/rules

We've extensively tested compatibility with the popular `nikunjy/rules` library to ensure smooth migration and familiar behavior. Here's our detailed compatibility analysis:

### ✅ 100% Compatible Features

These features work **identically** between both libraries:

| Feature | Our Library | nikunjy/rules | Status | Notes |
|---------|-------------|---------------|---------|-------|
| **Basic Operators** | ✅ | ✅ | 🟢 **Identical** | `eq`, `ne`, `lt`, `gt`, `le`, `ge` |
| **String Operators** | ✅ | ✅ | 🟢 **Identical** | `co`, `sw`, `ew` with same behavior |
| **Logical Operators** | ✅ | ✅ | 🟢 **Identical** | `and`, `or`, `not` with short-circuit evaluation |
| **Membership Operator** | ✅ | ✅ | 🟢 **Identical** | `in` with arrays, strict type matching |
| **Presence Operator** | ✅ | ✅ | 🟢 **Identical** | `pr` for property existence |
| **Nested Properties** | ✅ | ✅ | 🟢 **Identical** | Dot notation: `user.profile.age` |
| **Type Safety** | ✅ | ✅ | 🟢 **Identical** | No cross-type comparisons except numeric |
| **Numeric Cross-Type** | ✅ | ✅ | 🟢 **Identical** | `int`/`float` comparisons: `42 == 42.0` |
| **time.Time Handling** | ✅ | ✅ | 🟢 **Identical** | Converts to string via `.String()` method |
| **Error Handling** | ✅ | ✅ | 🟢 **Identical** | Same behavior for invalid operations |

### 🔧 Our Extensions (Intentionally Different)

Features that we've added beyond nikunjy/rules capabilities:

| Feature | Our Library | nikunjy/rules | Status | Description |
|---------|-------------|---------------|---------|-------------|
| **DateTime Operators** | ✅ `dq`, `dn`, `be`, `bq`, `af`, `aq` | ❌ Not available | 🟡 **Our Extension** | Native datetime comparison with RFC3339 and Unix timestamps |
| **Performance** | ⚡ **25-144x faster** | ✅ Good | 🟡 **Enhanced** | Sub-100ns evaluation, zero allocations |
| **rule.D Type Alias** | ✅ Clean API | ✅ `map[string]interface{}` | 🟡 **Enhanced** | Cleaner syntax: `rule.D{...}` |
| **Memory Usage** | ✅ **0 allocs/op** | ❌ High allocation | 🟡 **Enhanced** | Zero-allocation evaluation |

### 📊 Migration Compatibility

**Drop-in replacement compatibility: 100%** ✅

```go
// nikunjy/rules code
result, err := rules.Evaluate(`user.age gt 18 and status eq "active"`, context)

// Our library - same API!
engine := rule.NewEngine()
result, err := engine.Evaluate(`user.age gt 18 and status eq "active"`, context)
```

### 🔍 Detailed Compatibility Matrix

#### Context Data Types

| Type | Our Library | nikunjy/rules | Compatibility |
|------|-------------|---------------|---------------|
| `string` | ✅ Full support | ✅ Full support | 🟢 **100%** |
| `int`, `int8-64` | ✅ Full support | ✅ Full support | 🟢 **100%** |
| `uint`, `uint8-64` | ✅ Full support | ✅ Full support | 🟢 **100%** |
| `float32`, `float64` | ✅ Full support | ✅ Full support | 🟢 **100%** |
| `bool` | ✅ Full support | ✅ Full support | 🟢 **100%** |
| `[]any` (arrays) | ✅ Full support | ✅ Full support | 🟢 **100%** |
| `time.Time` | ✅ **Enhanced** | ✅ String conversion | 🟢 **100%** + datetime operators |
| `map[string]any` | ✅ Full support | ✅ Full support | 🟢 **100%** |

#### Rule Syntax Support

| Rule Type | Example | Our Library | nikunjy/rules | Compatibility |
|-----------|---------|-------------|---------------|---------------|
| Simple comparison | `age eq 25` | ✅ | ✅ | 🟢 **100%** |
| String operations | `name co "John"` | ✅ | ✅ | 🟢 **100%** |
| Array membership | `role in ["admin", "user"]` | ✅ | ✅ | 🟢 **100%** |
| Nested properties | `user.profile.age gt 18` | ✅ | ✅ | 🟢 **100%** |
| Complex logical | `(a eq 1 and b gt 2) or c pr` | ✅ | ✅ | 🟢 **100%** |
| DateTime (our extension) | `created_at af "2024-01-01T00:00:00Z"` | ✅ | ❌ | 🟡 **Our Extension** |

### 🧪 Verification

Our compatibility is **proven by comprehensive automated tests** that run the same rules against both libraries and compare results:

- **✅ 22/22 basic compatibility tests pass**
- **✅ 3/3 time.Time compatibility tests pass**  
- **✅ 100% compatibility rate achieved**

See `test/compatibility_test.go` for the complete test suite that validates our compatibility claims.

### 🚀 Migration Guide

**Step 1**: Replace imports
```go
// Before
import "github.com/nikunjy/rules"

// After  
import "github.com/NSXBet/rule-engine"
```

**Step 2**: Update API calls
```go
// Before
result, err := rules.Evaluate(rule, context)

// After
engine := rule.NewEngine()
result, err := engine.Evaluate(rule, context)
```

**Step 3**: Optionally use our enhancements
```go
// Use our cleaner type alias
context := rule.D{"user": rule.D{"age": 25}}

// Use datetime operators for time-based rules
result, _ := engine.Evaluate(`created_at af "2024-01-01T00:00:00Z"`, context)

// Pre-compile frequently used rules for maximum performance
engine.AddQuery(`user.role eq "admin"`)
```

---

## 🤝 Contributing

We'd love your help making this engine even better! 🛠️

### Getting Started

1. **Fork the repository**
2. **Clone your fork**: `git clone https://github.com/yourusername/rule-engine`
3. **Create a branch**: `git checkout -b feature/amazing-feature`
4. **Make your changes** 
5. **Run tests**: `make test`
6. **Run lints**: `make lint` (must be 100% clean ✅)
7. **Commit**: `git commit -m "Add amazing feature"`
8. **Push**: `git push origin feature/amazing-feature`
9. **Create a Pull Request**

### Development Setup

```bash
# Install dependencies
go mod download

# Run tests
make test

# Run benchmarks  
make bench

# Run linter (must pass 100%)
make lint

# Format code
make format
```

### What We're Looking For

- 🐛 **Bug fixes** with test cases
- ⚡ **Performance improvements** with benchmarks
- 📚 **Documentation improvements**
- 🧪 **Additional test coverage**
- 🔧 **New operators** (with use cases)
- 🌐 **Language features** that maintain zero-allocation goals

### Code Standards

- ✅ All tests must pass (`make test`)
- ✅ 100% lint compliance (`make lint`) 
- ✅ Zero allocations in hot paths
- ✅ Comprehensive test coverage for new features
- ✅ Benchmark comparisons for performance changes
- ✅ Clear commit messages and PR descriptions

### Performance Requirements

Any changes to core evaluation logic must maintain:
- **Sub-100ns evaluation times** for simple operations
- **Zero allocations** during rule evaluation  
- **Thread safety** for concurrent usage

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

<div align="center">

**Built with ❤️ for the Go community**

⭐ **Star us on GitHub if this helped you!** ⭐

[Report Bug](https://github.com/NSXBet/rule-engine/issues) | [Request Feature](https://github.com/NSXBet/rule-engine/issues) | [Contribute](https://github.com/NSXBet/rule-engine/pulls)

</div>