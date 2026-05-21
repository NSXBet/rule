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

	b.ReportAllocs()

	if allocs := testing.AllocsPerRun(1, func() {
		_, _ = engine.Evaluate(rule, ctx)
	}); allocs != 0 {
		b.Fatalf("expected 0 allocs/op, got %f", allocs)
	}

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

	b.ReportAllocs()

	if allocs := testing.AllocsPerRun(1, func() {
		_, _ = engine.Evaluate(rule, ctx)
	}); allocs != 0 {
		b.Fatalf("expected 0 allocs/op, got %f", allocs)
	}

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

	b.ReportAllocs()

	if allocs := testing.AllocsPerRun(1, func() {
		_, _ = engine.Evaluate(rule, ctx)
	}); allocs != 0 {
		b.Fatalf("expected 0 allocs/op, got %f", allocs)
	}

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

	b.ReportAllocs()

	// 1 alloc (24 B) from []any interface boxing in Go runtime.
	if allocs := testing.AllocsPerRun(1, func() {
		_, _ = engine.Evaluate(rule, ctx)
	}); allocs > 1 {
		b.Fatalf("expected <= 1 allocs/op, got %f", allocs)
	}

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

	b.ReportAllocs()

	if allocs := testing.AllocsPerRun(1, func() {
		_, _ = engine.Evaluate(rule, ctx)
	}); allocs != 0 {
		b.Fatalf("expected 0 allocs/op, got %f", allocs)
	}

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

	b.ReportAllocs()

	if allocs := testing.AllocsPerRun(1, func() {
		_, _ = engine.Evaluate(rule, ctx)
	}); allocs != 0 {
		b.Fatalf("expected 0 allocs/op, got %f", allocs)
	}

	b.ResetTimer()

	for range b.N {
		result, err := engine.Evaluate(rule, ctx)
		if err != nil || !result {
			b.Fatalf("Expected true result, got %v, %v", result, err)
		}
	}
}

// Benchmark array length evaluation for zero allocations.
func BenchmarkOptimizedEngineArrayLength(b *testing.B) {
	engine := NewEngine()
	ctx := D{
		"selections": []any{
			D{"id": "1"},
			D{"id": "2"},
			D{"id": "3"},
		},
	}
	rule := "selections.length gt 2"

	// Pre-compile rule
	engine.AddQuery(rule)

	b.ReportAllocs()

	if allocs := testing.AllocsPerRun(1, func() {
		_, _ = engine.Evaluate(rule, ctx)
	}); allocs != 0 {
		b.Fatalf("expected 0 allocs/op, got %f", allocs)
	}

	b.ResetTimer()

	for range b.N {
		result, err := engine.Evaluate(rule, ctx)
		if err != nil || !result {
			b.Fatalf("Expected true result, got %v, %v", result, err)
		}
	}
}

// Benchmark any quantifier evaluation (1 alloc, 24 B from []any interface boxing).
func BenchmarkOptimizedEngineQuantifierAny(b *testing.B) {
	engine := NewEngine()
	ctx := D{
		"selections": []any{
			D{"event_status": "EVENT_STATUS_IN_PROGRESS", "is_live": true},
			D{"event_status": "EVENT_STATUS_CANCELLED", "is_live": false},
			D{"event_status": "EVENT_STATUS_FINISHED", "is_live": false},
		},
	}
	rule := `selections any (event_status eq "EVENT_STATUS_CANCELLED")`

	// Pre-compile rule
	engine.AddQuery(rule)

	b.ReportAllocs()

	// 1 alloc (24 B) from []any interface boxing in Go runtime.
	if allocs := testing.AllocsPerRun(1, func() {
		_, _ = engine.Evaluate(rule, ctx)
	}); allocs > 1 {
		b.Fatalf("expected <= 1 allocs/op, got %f", allocs)
	}

	b.ResetTimer()

	for range b.N {
		result, err := engine.Evaluate(rule, ctx)
		if err != nil || !result {
			b.Fatalf("Expected true result, got %v, %v", result, err)
		}
	}
}

// Benchmark all quantifier evaluation (1 alloc, 24 B from []any interface boxing).
func BenchmarkOptimizedEngineQuantifierAll(b *testing.B) {
	engine := NewEngine()
	ctx := D{
		"items": []any{
			D{"is_valid": true},
			D{"is_valid": true},
			D{"is_valid": true},
		},
	}
	rule := "items all (is_valid eq true)"

	// Pre-compile rule
	engine.AddQuery(rule)

	b.ReportAllocs()

	// 1 alloc (24 B) from []any interface boxing in Go runtime.
	if allocs := testing.AllocsPerRun(1, func() {
		_, _ = engine.Evaluate(rule, ctx)
	}); allocs > 1 {
		b.Fatalf("expected <= 1 allocs/op, got %f", allocs)
	}

	b.ResetTimer()

	for range b.N {
		result, err := engine.Evaluate(rule, ctx)
		if err != nil || !result {
			b.Fatalf("Expected true result, got %v, %v", result, err)
		}
	}
}

// Benchmark none quantifier evaluation (1 alloc, 24 B from []any interface boxing).
func BenchmarkOptimizedEngineQuantifierNone(b *testing.B) {
	engine := NewEngine()
	ctx := D{
		"items": []any{
			D{"status": "active"},
			D{"status": "pending"},
			D{"status": "completed"},
		},
	}
	rule := `items none (status eq "error")`

	// Pre-compile rule
	engine.AddQuery(rule)

	b.ReportAllocs()

	// 1 alloc (24 B) from []any interface boxing in Go runtime.
	if allocs := testing.AllocsPerRun(1, func() {
		_, _ = engine.Evaluate(rule, ctx)
	}); allocs > 1 {
		b.Fatalf("expected <= 1 allocs/op, got %f", allocs)
	}

	b.ResetTimer()

	for range b.N {
		result, err := engine.Evaluate(rule, ctx)
		if err != nil || !result {
			b.Fatalf("Expected true result, got %v, %v", result, err)
		}
	}
}

// Benchmark complex real-world betting rule with length + quantifiers.
func BenchmarkOptimizedEngineComplexBetting(b *testing.B) {
	engine := NewEngine()
	ctx := D{
		"bet_type": "BET_TYPE_MULTIPLE",
		"customer_data": D{
			"is_vip": true,
		},
		"selections": []any{
			D{
				"event_status": "EVENT_STATUS_IN_PROGRESS",
				"is_live":      true,
				"odd":          2.5,
				"provider":     "PROVIDER_SPORTRADAR",
			},
			D{"event_status": "EVENT_STATUS_NOT_STARTED", "is_live": false, "odd": 1.8, "provider": "PROVIDER_RAMP"},
			D{
				"event_status": "EVENT_STATUS_IN_PROGRESS",
				"is_live":      true,
				"odd":          3.1,
				"provider":     "PROVIDER_SPORTRADAR",
			},
		},
	}
	rule := `bet_type eq "BET_TYPE_MULTIPLE" and selections.length ge 2 and selections none (event_status eq "EVENT_STATUS_CANCELLED") and selections any (is_live eq true and odd gt 2.0)`

	// Pre-compile rule
	engine.AddQuery(rule)

	b.ReportAllocs()

	// Multiple allocs expected from []any interface boxing (one per quantifier).
	if allocs := testing.AllocsPerRun(1, func() {
		_, _ = engine.Evaluate(rule, ctx)
	}); allocs > 2 {
		b.Fatalf("expected <= 2 allocs/op, got %f", allocs)
	}

	b.ResetTimer()

	for range b.N {
		result, err := engine.Evaluate(rule, ctx)
		if err != nil || !result {
			b.Fatalf("Expected true result, got %v, %v", result, err)
		}
	}
}

// Benchmark quantifier with short-circuit (first element matches).
func BenchmarkOptimizedEngineQuantifierAnyShortCircuit(b *testing.B) {
	engine := NewEngine()

	// Create array with 100 elements, first one matches
	items := make([]any, 100)

	items[0] = D{"status": "target"}
	for i := 1; i < 100; i++ {
		items[i] = D{"status": "other"}
	}

	ctx := D{"items": items}
	rule := `items any (status eq "target")`

	// Pre-compile rule
	engine.AddQuery(rule)

	b.ReportAllocs()

	// 1 alloc (24 B) from []any interface boxing in Go runtime.
	if allocs := testing.AllocsPerRun(1, func() {
		_, _ = engine.Evaluate(rule, ctx)
	}); allocs > 1 {
		b.Fatalf("expected <= 1 allocs/op, got %f", allocs)
	}

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

	b.ReportAllocs()

	if allocs := testing.AllocsPerRun(1, func() {
		_, _ = evaluator.Evaluate(ast, ctx)
	}); allocs != 0 {
		b.Fatalf("expected 0 allocs/op, got %f", allocs)
	}

	b.ResetTimer()

	for range b.N {
		result, err := evaluator.Evaluate(ast, ctx)
		if err != nil || !result {
			b.Fatalf("Expected true result, got %v, %v", result, err)
		}
	}
}

// BenchmarkWhereCount validates near-zero allocations for "where (...).length" rules.
// 1 alloc (24 B) from []any interface boxing in Go runtime.
func BenchmarkWhereCount(b *testing.B) {
	engine := NewEngine()
	ctx := D{
		"selections": []any{
			D{"odd": 1.5},
			D{"odd": 2.0},
			D{"odd": 1.4},
			D{"odd": 1.8},
			D{"odd": 1.0},
		},
	}
	query := "selections where (odd ge 1.4).length ge 4"

	if err := engine.AddQuery(query); err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()

	if allocs := testing.AllocsPerRun(1, func() {
		_, _ = engine.Evaluate(query, ctx)
	}); allocs > 1 {
		b.Fatalf("expected <= 1 allocs/op, got %f", allocs)
	}

	b.ResetTimer()

	for range b.N {
		result, err := engine.Evaluate(query, ctx)
		if err != nil || !result {
			b.Fatalf("Expected true result, got %v, %v", result, err)
		}
	}
}

// BenchmarkWhereWithQuantifier validates near-zero allocations for "where (...) any (...)" rules.
// 1 alloc (24 B) from []any interface boxing in Go runtime.
func BenchmarkWhereWithQuantifier(b *testing.B) {
	engine := NewEngine()
	ctx := D{
		"selections": []any{
			D{"odd": 1.5, "provider": "Y"},
			D{"odd": 2.0, "provider": "X"},
			D{"odd": 1.0, "provider": "X"},
		},
	}
	query := `selections where (odd ge 1.4) any (provider eq "X")`

	if err := engine.AddQuery(query); err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()

	if allocs := testing.AllocsPerRun(1, func() {
		_, _ = engine.Evaluate(query, ctx)
	}); allocs > 1 {
		b.Fatalf("expected <= 1 allocs/op, got %f", allocs)
	}

	b.ResetTimer()

	for range b.N {
		result, err := engine.Evaluate(query, ctx)
		if err != nil || !result {
			b.Fatalf("Expected true result, got %v, %v", result, err)
		}
	}
}
