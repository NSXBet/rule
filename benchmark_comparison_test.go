package rule

import (
	"testing"
	"github.com/nikunjy/rules"
)

// Benchmark comparing our optimized engine vs nikunjy/rules library
func BenchmarkComparisonSimple(b *testing.B) {
	ctx := map[string]any{"x": 10}
	rule := "x eq 10"
	
	b.Run("OurEngine", func(b *testing.B) {
		engine := NewEngine()
		engine.AddQuery(rule)
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result, err := engine.Evaluate(rule, ctx)
			if err != nil || !result {
				b.Fatalf("Expected true result, got %v, %v", result, err)
			}
		}
	})
	
	b.Run("NikunjyRules", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result, err := rules.Evaluate(rule, ctx)
			if err != nil || !result {
				b.Fatalf("Expected true result, got %v, %v", result, err)
			}
		}
	})
}

func BenchmarkComparisonComplex(b *testing.B) {
	ctx := map[string]any{
		"user": map[string]any{
			"name": "John",
			"age":  30,
		},
		"status": "active",
	}
	rule := `(user.age gt 18 and status eq "active") or user.name co "Admin"`
	
	b.Run("OurEngine", func(b *testing.B) {
		engine := NewEngine()
		engine.AddQuery(rule)
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result, err := engine.Evaluate(rule, ctx)
			if err != nil || !result {
				b.Fatalf("Expected true result, got %v, %v", result, err)
			}
		}
	})
	
	b.Run("NikunjyRules", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result, err := rules.Evaluate(rule, ctx)
			if err != nil || !result {
				b.Fatalf("Expected true result, got %v, %v", result, err)
			}
		}
	})
}

func BenchmarkComparisonStringOps(b *testing.B) {
	ctx := map[string]any{"name": "John Doe", "email": "john@example.com"}
	rule := `name co "John" and email ew ".com"`
	
	b.Run("OurEngine", func(b *testing.B) {
		engine := NewEngine()
		engine.AddQuery(rule)
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result, err := engine.Evaluate(rule, ctx)
			if err != nil || !result {
				b.Fatalf("Expected true result, got %v, %v", result, err)
			}
		}
	})
	
	b.Run("NikunjyRules", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result, err := rules.Evaluate(rule, ctx)
			if err != nil || !result {
				b.Fatalf("Expected true result, got %v, %v", result, err)
			}
		}
	})
}

func BenchmarkComparisonInOperator(b *testing.B) {
	ctx := map[string]any{
		"color": "red",
		"allowed": []string{"red", "green", "blue"}, // Use []string instead of []any for nikunjy compatibility
	}
	rule := "color in allowed"
	
	b.Run("OurEngine", func(b *testing.B) {
		engine := NewEngine()
		engine.AddQuery(rule)
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result, err := engine.Evaluate(rule, ctx)
			if err != nil || !result {
				b.Fatalf("Expected true result, got %v, %v", result, err)
			}
		}
	})
	
	b.Run("NikunjyRules", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result, err := rules.Evaluate(rule, ctx)
			if err != nil {
				b.Fatalf("Error from nikunjy rules: %v", err)
			}
			if !result {
				b.Fatalf("Expected true result, got %v", result)
			}
		}
	})
}

func BenchmarkComparisonNestedProps(b *testing.B) {
	ctx := map[string]any{
		"user": map[string]any{
			"profile": map[string]any{
				"settings": map[string]any{
					"theme": "dark",
				},
			},
		},
	}
	rule := `user.profile.settings.theme eq "dark"`
	
	b.Run("OurEngine", func(b *testing.B) {
		engine := NewEngine()
		engine.AddQuery(rule)
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result, err := engine.Evaluate(rule, ctx)
			if err != nil || !result {
				b.Fatalf("Expected true result, got %v, %v", result, err)
			}
		}
	})
	
	b.Run("NikunjyRules", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result, err := rules.Evaluate(rule, ctx)
			if err != nil || !result {
				b.Fatalf("Expected true result, got %v, %v", result, err)
			}
		}
	})
}

// Test with different query patterns to show pre-compilation advantage
func BenchmarkComparisonManyQueries(b *testing.B) {
	ctx := map[string]any{"x": 10, "y": 20, "z": 30}
	queries := []string{
		"x eq 10",
		"y gt 15",
		"z lt 50",
		"x lt y",
		"y le z",
	}
	
	b.Run("OurEngine", func(b *testing.B) {
		engine := NewEngine()
		// Pre-compile all queries
		for _, query := range queries {
			engine.AddQuery(query)
		}
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, query := range queries {
				result, err := engine.Evaluate(query, ctx)
				if err != nil || !result {
					b.Fatalf("Expected true result for %s, got %v, %v", query, result, err)
				}
			}
		}
	})
	
	b.Run("NikunjyRules", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, query := range queries {
				result, err := rules.Evaluate(query, ctx)
				if err != nil {
					b.Logf("Error from nikunjy rules for %s: %v", query, err)
					continue
				}
				if !result {
					b.Logf("Expected true result for %s, got %v", query, result)
				}
			}
		}
	})
}