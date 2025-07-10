package rule

import (
	"bytes"
	"strings"
	"testing"
	"text/template"

	"github.com/nikunjy/rules"
)

// Benchmark comparing our optimized engine vs nikunjy/rules vs text/template.
func BenchmarkComparisonSimple(b *testing.B) {
	ctx := map[string]any{"x": 10}
	rule := "x eq 10"
	tmplText := "{{if eq .x 10}}true{{else}}false{{end}}"

	b.Run("OurEngine", func(b *testing.B) {
		engine := NewEngine()
		engine.AddQuery(rule)

		b.ResetTimer()

		for range b.N {
			result, err := engine.Evaluate(rule, ctx)
			if err != nil || !result {
				b.Fatalf("Expected true result, got %v, %v", result, err)
			}
		}
	})

	b.Run("NikunjyRules", func(b *testing.B) {
		b.ResetTimer()

		for range b.N {
			result, err := rules.Evaluate(rule, ctx)
			if err != nil || !result {
				b.Fatalf("Expected true result, got %v, %v", result, err)
			}
		}
	})

	b.Run("TextTemplate", func(b *testing.B) {
		tmpl := template.Must(template.New("test").Parse(tmplText))

		var buf bytes.Buffer

		b.ResetTimer()

		for range b.N {
			buf.Reset()

			err := tmpl.Execute(&buf, ctx)
			if err != nil {
				b.Fatalf("Template execution error: %v", err)
			}

			result := strings.TrimSpace(buf.String())
			if result != "true" {
				b.Fatalf("Expected 'true', got '%s'", result)
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
	tmplText := `{{if or (and (gt .user.age 18) (eq .status "active")) (contains .user.name "Admin")}}true{{else}}false{{end}}`

	b.Run("OurEngine", func(b *testing.B) {
		engine := NewEngine()
		engine.AddQuery(rule)

		b.ResetTimer()

		for range b.N {
			result, err := engine.Evaluate(rule, ctx)
			if err != nil || !result {
				b.Fatalf("Expected true result, got %v, %v", result, err)
			}
		}
	})

	b.Run("NikunjyRules", func(b *testing.B) {
		b.ResetTimer()

		for range b.N {
			result, err := rules.Evaluate(rule, ctx)
			if err != nil || !result {
				b.Fatalf("Expected true result, got %v, %v", result, err)
			}
		}
	})

	b.Run("TextTemplate", func(b *testing.B) {
		funcMap := template.FuncMap{
			"contains": strings.Contains,
		}
		tmpl := template.Must(template.New("test").Funcs(funcMap).Parse(tmplText))

		var buf bytes.Buffer

		b.ResetTimer()

		for range b.N {
			buf.Reset()

			err := tmpl.Execute(&buf, ctx)
			if err != nil {
				b.Fatalf("Template execution error: %v", err)
			}

			result := strings.TrimSpace(buf.String())
			if result != "true" {
				b.Fatalf("Expected 'true', got '%s'", result)
			}
		}
	})
}

func BenchmarkComparisonStringOps(b *testing.B) {
	ctx := map[string]any{"name": "John Doe", "email": "john@example.com"}
	rule := `name co "John" and email ew ".com"`
	tmplText := `{{if and (contains .name "John") (hasSuffix .email ".com")}}true{{else}}false{{end}}`

	b.Run("OurEngine", func(b *testing.B) {
		engine := NewEngine()
		engine.AddQuery(rule)

		b.ResetTimer()

		for range b.N {
			result, err := engine.Evaluate(rule, ctx)
			if err != nil || !result {
				b.Fatalf("Expected true result, got %v, %v", result, err)
			}
		}
	})

	b.Run("NikunjyRules", func(b *testing.B) {
		b.ResetTimer()

		for range b.N {
			result, err := rules.Evaluate(rule, ctx)
			if err != nil || !result {
				b.Fatalf("Expected true result, got %v, %v", result, err)
			}
		}
	})

	b.Run("TextTemplate", func(b *testing.B) {
		funcMap := template.FuncMap{
			"contains":  strings.Contains,
			"hasSuffix": strings.HasSuffix,
		}
		tmpl := template.Must(template.New("test").Funcs(funcMap).Parse(tmplText))

		var buf bytes.Buffer

		b.ResetTimer()

		for range b.N {
			buf.Reset()

			err := tmpl.Execute(&buf, ctx)
			if err != nil {
				b.Fatalf("Template execution error: %v", err)
			}

			result := strings.TrimSpace(buf.String())
			if result != "true" {
				b.Fatalf("Expected 'true', got '%s'", result)
			}
		}
	})
}

//nolint:gocognit // Benchmark function complexity is acceptable
func BenchmarkComparisonInOperator(b *testing.B) {
	ctx := map[string]any{
		"color": "red",
	}
	rule := `color in ["red", "green", "blue"]`
	tmplText := `{{if or (eq .color "red") (eq .color "green") (eq .color "blue")}}true{{else}}false{{end}}`

	b.Run("OurEngine", func(b *testing.B) {
		engine := NewEngine()
		engine.AddQuery(rule)

		b.ResetTimer()

		for range b.N {
			result, err := engine.Evaluate(rule, ctx)
			if err != nil || !result {
				b.Fatalf("Expected true result, got %v, %v", result, err)
			}
		}
	})

	b.Run("NikunjyRules", func(b *testing.B) {
		b.ResetTimer()

		for range b.N {
			result, err := rules.Evaluate(rule, ctx)
			if err != nil {
				b.Fatalf("Error from nikunjy rules: %v", err)
			}

			if !result {
				b.Fatalf("Expected true result, got %v", result)
			}
		}
	})

	b.Run("TextTemplate", func(b *testing.B) {
		tmpl := template.Must(template.New("test").Parse(tmplText))

		var buf bytes.Buffer

		b.ResetTimer()

		for range b.N {
			buf.Reset()

			err := tmpl.Execute(&buf, ctx)
			if err != nil {
				b.Fatalf("Template execution error: %v", err)
			}

			result := strings.TrimSpace(buf.String())
			if result != "true" {
				b.Fatalf("Expected 'true', got '%s'", result)
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
	tmplText := `{{if eq .user.profile.settings.theme "dark"}}true{{else}}false{{end}}`

	b.Run("OurEngine", func(b *testing.B) {
		engine := NewEngine()
		engine.AddQuery(rule)

		b.ResetTimer()

		for range b.N {
			result, err := engine.Evaluate(rule, ctx)
			if err != nil || !result {
				b.Fatalf("Expected true result, got %v, %v", result, err)
			}
		}
	})

	b.Run("NikunjyRules", func(b *testing.B) {
		b.ResetTimer()

		for range b.N {
			result, err := rules.Evaluate(rule, ctx)
			if err != nil || !result {
				b.Fatalf("Expected true result, got %v, %v", result, err)
			}
		}
	})

	b.Run("TextTemplate", func(b *testing.B) {
		tmpl := template.Must(template.New("test").Parse(tmplText))

		var buf bytes.Buffer

		b.ResetTimer()

		for range b.N {
			buf.Reset()

			err := tmpl.Execute(&buf, ctx)
			if err != nil {
				b.Fatalf("Template execution error: %v", err)
			}

			result := strings.TrimSpace(buf.String())
			if result != "true" {
				b.Fatalf("Expected 'true', got '%s'", result)
			}
		}
	})
}

// Test with different query patterns to show pre-compilation advantage.
//
//nolint:gocognit // Benchmark function complexity is acceptable
func BenchmarkComparisonManyQueries(b *testing.B) {
	ctx := map[string]any{"x": 10, "y": 20, "z": 30}
	queries := []string{
		"x eq 10",
		"y gt 15",
		"z lt 50",
		"x lt y",
		"y le z",
	}
	templates := []string{
		"{{if eq .x 10}}true{{else}}false{{end}}",
		"{{if gt .y 15}}true{{else}}false{{end}}",
		"{{if lt .z 50}}true{{else}}false{{end}}",
		"{{if lt .x .y}}true{{else}}false{{end}}",
		"{{if le .y .z}}true{{else}}false{{end}}",
	}

	b.Run("OurEngine", func(b *testing.B) {
		engine := NewEngine()
		// Pre-compile all queries
		for _, query := range queries {
			engine.AddQuery(query)
		}

		b.ResetTimer()

		for range b.N {
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

		for range b.N {
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

	b.Run("TextTemplate", func(b *testing.B) {
		// Pre-compile all templates
		compiledTemplates := make([]*template.Template, len(templates))
		for i, tmplText := range templates {
			compiledTemplates[i] = template.Must(template.New("test").Parse(tmplText))
		}

		var buf bytes.Buffer

		b.ResetTimer()

		for range b.N {
			for _, tmpl := range compiledTemplates {
				buf.Reset()

				err := tmpl.Execute(&buf, ctx)
				if err != nil {
					b.Fatalf("Template execution error: %v", err)
				}

				result := strings.TrimSpace(buf.String())
				if result != "true" {
					b.Fatalf("Expected 'true', got '%s'", result)
				}
			}
		}
	})
}

// Benchmark datetime operations - showcasing our unique advantage.
func BenchmarkComparisonDateTimeOps(b *testing.B) {
	ctx := map[string]any{
		"created_at": "2024-07-09T22:12:01Z",
		"updated_at": "2024-07-09T22:12:00Z",
		"timestamp":  int64(1720558320),
	}
	rule := `created_at af updated_at and timestamp be 1720558400`

	b.Run("OurEngine", func(b *testing.B) {
		engine := NewEngine()
		engine.AddQuery(rule)

		b.ResetTimer()

		for range b.N {
			result, err := engine.Evaluate(rule, ctx)
			if err != nil || !result {
				b.Fatalf("Expected true result, got %v, %v", result, err)
			}
		}
	})
	// Note: nikunjy/rules doesn't support datetime operators like 'af' and 'be'
	// This demonstrates our unique capability and competitive advantage
}
