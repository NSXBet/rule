package rule

import (
	"testing"
)

// Benchmark optimized engine operations for zero allocations.
func BenchmarkOptimizedEngineSimple(b *testing.B) {
	engine := NewEngine()
	ctx := D{"x": 10}
	rule := "x eq 10"

	// Pre-compile rule
	engine.AddQuery(rule)

	b.ResetTimer()

	for range b.N {
		result, err := engine.Evaluate(rule, ctx)
		if err != nil || !result {
			b.Fatalf("Expected true result, got %v, %v", result, err)
		}
	}
}

func BenchmarkOptimizedEngineComplex(b *testing.B) {
	engine := NewEngine()
	ctx := D{
		"user": D{
			"name": "John",
			"age":  30,
		},
		"status": "active",
	}
	rule := "(user.age gt 18 and status eq \"active\") or user.name co \"Admin\""

	// Pre-compile rule
	engine.AddQuery(rule)

	b.ResetTimer()

	for range b.N {
		result, err := engine.Evaluate(rule, ctx)
		if err != nil || !result {
			b.Fatalf("Expected true result, got %v, %v", result, err)
		}
	}
}

func BenchmarkOptimizedEngineStringOps(b *testing.B) {
	engine := NewEngine()
	ctx := D{"name": "John Doe", "email": "john@example.com"}
	rule := "name co \"John\" and email ew \".com\""

	// Pre-compile rule
	engine.AddQuery(rule)

	b.ResetTimer()

	for range b.N {
		result, err := engine.Evaluate(rule, ctx)
		if err != nil || !result {
			b.Fatalf("Expected true result, got %v, %v", result, err)
		}
	}
}

func BenchmarkOptimizedEngineInOperator(b *testing.B) {
	engine := NewEngine()
	ctx := D{
		"color":   "red",
		"allowed": []any{"red", "green", "blue"},
	}
	rule := "color in allowed"

	// Pre-compile rule
	engine.AddQuery(rule)

	b.ResetTimer()

	for range b.N {
		result, err := engine.Evaluate(rule, ctx)
		if err != nil || !result {
			b.Fatalf("Expected true result, got %v, %v", result, err)
		}
	}
}

func BenchmarkOptimizedEngineNestedProps(b *testing.B) {
	engine := NewEngine()
	ctx := D{
		"user": D{
			"profile": D{
				"settings": D{
					"theme": "dark",
				},
			},
		},
	}
	rule := "user.profile.settings.theme eq \"dark\""

	// Pre-compile rule
	engine.AddQuery(rule)

	b.ResetTimer()

	for range b.N {
		result, err := engine.Evaluate(rule, ctx)
		if err != nil || !result {
			b.Fatalf("Expected true result, got %v, %v", result, err)
		}
	}
}

func BenchmarkOptimizedStandalone(b *testing.B) {
	engine := NewEngine()
	ctx := D{"x": 10}
	rule := "x eq 10"

	b.ResetTimer()

	for range b.N {
		result, err := engine.Evaluate(rule, ctx)
		if err != nil || !result {
			b.Fatalf("Expected true result, got %v, %v", result, err)
		}
	}
}

// Direct evaluator benchmarks.
func BenchmarkZeroAllocEvaluatorDirect(b *testing.B) {
	evaluator := NewEvaluator()
	ast := NewBinaryOpNode(EQ,
		NewIdentifierNode("x"),
		NewNumberLiteralNode(10))
	ctx := D{"x": 10}

	b.ResetTimer()

	for range b.N {
		result, err := evaluator.Evaluate(ast, ctx)
		if err != nil || !result {
			b.Fatalf("Expected true result, got %v, %v", result, err)
		}
	}
}
