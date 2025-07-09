package rule

import (
	"fmt"
	"testing"
)

func TestDebugLargeNumbers(t *testing.T) {
	engine := NewEngine()
	
	// Test the actual failing case
	rule := "x gt 9223372036854775806"
	ctx := map[string]any{"x": int64(9223372036854775807)}
	
	// Check if our isLargeInteger function works
	evaluator := NewEvaluator()
	fmt.Printf("isLargeInteger(9223372036854775806): %v\n", evaluator.isLargeInteger(int64(9223372036854775806)))
	fmt.Printf("isLargeInteger(9223372036854775807): %v\n", evaluator.isLargeInteger(int64(9223372036854775807)))
	
	// Test our large number comparison
	left := int64(9223372036854775807)
	right := int64(9223372036854775806)
	
	fmt.Printf("Direct comparison: %d > %d = %v\n", left, right, left > right)
	
	// Test our evaluator comparison
	result := evaluator.compareLargeNumbers(left, right, func(a, b float64) bool { return a > b })
	fmt.Printf("compareLargeNumbers result: %v\n", result)
	
	// Test the actual engine
	if err := engine.AddQuery(rule); err != nil {
		t.Fatalf("failed to add query: %v", err)
	}
	
	got, err := engine.Evaluate(rule, ctx)
	if err != nil {
		t.Fatalf("evaluation failed: %v", err)
	}
	
	fmt.Printf("Engine result: %v (expected: true)\n", got)
}