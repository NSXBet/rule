package rule

import (
	"testing"

	"github.com/NSXBet/rule-engine/test"
)

func TestEngineRound1(t *testing.T) {
	engine := NewEngine()

	all := [][]test.TestCase{
		test.EqualTests,
		test.RelationalTests,
		test.StringOpTests,
		test.InTests,
		test.PresenceTests,
		test.LogicalTests,
		test.PropCompareTests,
		test.NestedPropTests,
		test.EdgeCaseTests,
		test.StringEdgeCaseTests,
		test.NumericEdgeCaseTests,
		test.ArrayEdgeCaseTests,
		test.PresenceEdgeCaseTests,
		test.ComplexLogicalTests,
		test.RealWorldTests,
		test.ErrorBoundaryTests,
		test.ExtremeValueTests,
		test.ComplexNestedLogicTests,
		test.RealWorldEdgeTests,
		test.WhitespaceTests,
		test.AdvancedPrecedenceTests,
		test.TypeCoercionStressTests,
		test.PerformanceStressTests,
		test.BoundaryConditionTests,
		test.SpecialNumericTests,
		test.ComplexStringPatternTests,
	}

	// Pre-compile all rules (optional with JIT compilation)
	for _, group := range all {
		for _, tc := range group {
			if err := engine.AddQuery(tc.Query); err != nil {
				t.Fatalf("failed to add query %q: %v", tc.Query, err)
			}
		}
	}

	for _, group := range all {
		for _, tc := range group {
			t.Run(tc.Name, func(t *testing.T) {
				got, err := engine.Evaluate(tc.Query, tc.Ctx)
				if err != nil {
					t.Errorf("query=%q, error=%v", tc.Query, err)
					return
				}
				if got != tc.Result {
					t.Errorf("query=%q, expected=%v, got=%v", tc.Query, tc.Result, got)
				}
			})
		}
	}
}

func TestEngineCompileOnce(t *testing.T) {
	engine := NewEngine()

	rule := "x eq 10 and y gt 5"
	ctx := map[string]any{"x": 10, "y": 6}

	compiled, err := engine.CompileRule(rule)
	if err != nil {
		t.Fatalf("failed to compile rule: %v", err)
	}

	for range 1000 {
		got, err := engine.EvaluateCompiled(compiled, ctx)
		if err != nil {
			t.Fatalf("evaluation failed: %v", err)
		}
		if !got {
			t.Fatalf("expected true, got false")
		}
	}
}

func TestEngineJITCompilation(t *testing.T) {
	engine := NewEngine()

	rule := "x eq 10 and y gt 5"
	ctx := map[string]any{"x": 10, "y": 6}

	// First evaluation should compile just-in-time
	got, err := engine.Evaluate(rule, ctx)
	if err != nil {
		t.Fatalf("JIT evaluation failed: %v", err)
	}
	if !got {
		t.Fatalf("expected true, got false")
	}

	// Second evaluation should use compiled version
	got, err = engine.Evaluate(rule, ctx)
	if err != nil {
		t.Fatalf("cached evaluation failed: %v", err)
	}
	if !got {
		t.Fatalf("expected true, got false")
	}
}

func BenchmarkEngineSimple(b *testing.B) {
	engine := NewEngine()
	rule := "x eq 10"
	ctx := map[string]any{"x": 10}

	b.ResetTimer()

	for b.Loop() {
		_, err := engine.Evaluate(rule, ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEngineComplex(b *testing.B) {
	engine := NewEngine()
	rule := "(x eq 10 and y gt 5) or (z co \"test\" and w in [1,2,3])"
	ctx := map[string]any{
		"x": 10,
		"y": 6,
		"z": "testing",
		"w": 2,
	}

	b.ResetTimer()

	for b.Loop() {
		_, err := engine.Evaluate(rule, ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEnginePrecompiled(b *testing.B) {
	engine := NewEngine()
	rule := "(x eq 10 and y gt 5) or (z co \"test\" and w in [1,2,3])"
	ctx := map[string]any{
		"x": 10,
		"y": 6,
		"z": "testing",
		"w": 2,
	}

	compiled, err := engine.CompileRule(rule)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for b.Loop() {
		_, err := engine.EvaluateCompiled(compiled, ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEngineNested(b *testing.B) {
	engine := NewEngine()
	rule := "user.profile.age ge 18 and user.status eq \"active\""
	ctx := map[string]any{
		"user": map[string]any{
			"profile": map[string]any{
				"age": 25,
			},
			"status": "active",
		},
	}

	compiled, err := engine.CompileRule(rule)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for b.Loop() {
		_, err := engine.EvaluateCompiled(compiled, ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEngineLargeNumbers(b *testing.B) {
	engine := NewEngine()
	rule := "x gt 9223372036854775806"
	ctx := map[string]any{"x": int64(9223372036854775807)}

	compiled, err := engine.CompileRule(rule)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for b.Loop() {
		_, err := engine.EvaluateCompiled(compiled, ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEngineStringOps(b *testing.B) {
	engine := NewEngine()
	rule := "email co \"@example.com\" and name sw \"John\""
	ctx := map[string]any{
		"email": "john.doe@example.com",
		"name":  "John Doe",
	}

	compiled, err := engine.CompileRule(rule)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for b.Loop() {
		_, err := engine.EvaluateCompiled(compiled, ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEngineArrayMembership(b *testing.B) {
	engine := NewEngine()
	rule := "status in [\"active\", \"pending\", \"verified\"]"
	ctx := map[string]any{"status": "active"}

	compiled, err := engine.CompileRule(rule)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for b.Loop() {
		_, err := engine.EvaluateCompiled(compiled, ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEngineZeroAlloc(b *testing.B) {
	engine := NewEngine()
	rule := "x gt 5"
	ctx := map[string]any{"x": 10}

	compiled, err := engine.CompileRule(rule)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for b.Loop() {
		result, err := engine.EvaluateCompiled(compiled, ctx)
		if err != nil {
			b.Fatal(err)
		}
		if !result {
			b.Fatal("Expected true")
		}
	}
}
