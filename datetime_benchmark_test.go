package rule

import (
	"runtime"
	"testing"
)

// BenchmarkDateTimeOperators tests all datetime operators for performance and allocations.
func BenchmarkDateTimeOperators(b *testing.B) {
	engine := NewEngine()

	// Test cases for different datetime operations
	testCases := []struct {
		name  string
		query string
		ctx   map[string]any
	}{
		// RFC 3339 vs literal
		{
			"RFC3339_After_Literal",
			`created_at af "2024-07-09T19:12:00-03:00"`,
			map[string]any{"created_at": "2024-07-09T22:13:00Z"},
		},
		{
			"RFC3339_Before_Literal",
			`created_at be "2024-07-09T19:12:00-03:00"`,
			map[string]any{"created_at": "2024-07-09T22:11:00Z"},
		},
		{
			"RFC3339_Equal_Literal",
			`created_at dq "2024-07-09T19:12:00-03:00"`,
			map[string]any{"created_at": "2024-07-09T22:12:00Z"},
		},

		// Unix timestamp vs literal
		{
			"Unix_After_Literal",
			`timestamp af 1720558320`,
			map[string]any{"timestamp": int64(1720558321)},
		},
		{
			"Unix_Before_Literal",
			`timestamp be 1720558320`,
			map[string]any{"timestamp": int64(1720558319)},
		},
		{
			"Unix_Equal_Literal",
			`timestamp dq 1720558320`,
			map[string]any{"timestamp": int64(1720558320)},
		},

		// Property vs Property
		{
			"Unix_Prop_vs_Prop",
			`start_time af end_time`,
			map[string]any{
				"start_time": int64(1720558321),
				"end_time":   int64(1720558320),
			},
		},
		{
			"RFC3339_Prop_vs_Prop",
			`created_at be updated_at`,
			map[string]any{
				"created_at": "2024-07-09T22:11:59Z",
				"updated_at": "2024-07-09T22:12:00Z",
			},
		},
		{
			"Mixed_Format_Prop_vs_Prop",
			`created_at af timestamp`,
			map[string]any{
				"created_at": "2024-07-09T22:12:01Z",
				"timestamp":  int64(1720558320),
			},
		},

		// Nested properties
		{
			"Nested_RFC3339_vs_Literal",
			`event.created_at af "2024-07-09T22:12:00Z"`,
			map[string]any{
				"event": map[string]any{
					"created_at": "2024-07-09T22:12:01Z",
				},
			},
		},
		{
			"Nested_Unix_vs_Literal",
			`session.timestamp be 1720558320`,
			map[string]any{
				"session": map[string]any{
					"timestamp": int64(1720558319),
				},
			},
		},
		{
			"Nested_vs_Nested",
			`session.start_time af user.created_at`,
			map[string]any{
				"session": map[string]any{
					"start_time": "2024-07-09T22:12:01Z",
				},
				"user": map[string]any{
					"created_at": "2024-07-09T22:12:00Z",
				},
			},
		},

		// Complex expressions
		{
			"Complex_DateTime_And",
			`start_time af "2024-07-09T22:12:00Z" and end_time be "2024-07-09T22:15:00Z"`,
			map[string]any{
				"start_time": "2024-07-09T22:12:01Z",
				"end_time":   "2024-07-09T22:14:59Z",
			},
		},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			// Verify the query works first
			result, err := engine.Evaluate(tc.query, tc.ctx)
			if err != nil {
				b.Fatalf("Query failed: %v", err)
			}

			if !result {
				b.Fatalf("Expected true result for test case %s", tc.name)
			}

			// Reset timer and measure allocations
			b.ResetTimer()
			b.ReportAllocs()

			var memStats runtime.MemStats

			runtime.ReadMemStats(&memStats)
			allocsBefore := memStats.Mallocs

			for range b.N {
				_, err := engine.Evaluate(tc.query, tc.ctx)
				if err != nil {
					b.Fatal(err)
				}
			}

			runtime.ReadMemStats(&memStats)
			allocsAfter := memStats.Mallocs
			totalAllocs := allocsAfter - allocsBefore

			if totalAllocs > 0 {
				b.Errorf("Expected 0 allocations, got %d total allocations for %d operations (%.2f allocs/op)",
					totalAllocs, b.N, float64(totalAllocs)/float64(b.N))
			}
		})
	}
}

// BenchmarkDateTimeVsRegularOperators compares datetime operators with regular operators.
func BenchmarkDateTimeVsRegularOperators(b *testing.B) {
	engine := NewEngine()

	b.Run("DateTime_After", func(b *testing.B) {
		query := `timestamp af 1720558320`
		ctx := map[string]any{"timestamp": int64(1720558321)}

		b.ResetTimer()
		b.ReportAllocs()

		for range b.N {
			_, err := engine.Evaluate(query, ctx)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Regular_GreaterThan", func(b *testing.B) {
		query := `timestamp gt 1720558320`
		ctx := map[string]any{"timestamp": int64(1720558321)}

		b.ResetTimer()
		b.ReportAllocs()

		for range b.N {
			_, err := engine.Evaluate(query, ctx)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("String_Contains", func(b *testing.B) {
		query := `text co "test"`
		ctx := map[string]any{"text": "this is a test string"}

		b.ResetTimer()
		b.ReportAllocs()

		for range b.N {
			_, err := engine.Evaluate(query, ctx)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkDateTimeStressCases tests edge cases that might be slower.
func BenchmarkDateTimeStressCases(b *testing.B) {
	engine := NewEngine()

	b.Run("RFC3339_Timezone_Conversion", func(b *testing.B) {
		query := `created_at dq "2024-07-09T19:12:00-03:00"`
		ctx := map[string]any{"created_at": "2024-07-09T22:12:00Z"}

		b.ResetTimer()
		b.ReportAllocs()

		for range b.N {
			_, err := engine.Evaluate(query, ctx)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Large_Unix_Timestamp", func(b *testing.B) {
		query := `timestamp dq 9223372036854775807`
		ctx := map[string]any{"timestamp": int64(9223372036854775807)}

		b.ResetTimer()
		b.ReportAllocs()

		for range b.N {
			_, err := engine.Evaluate(query, ctx)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Deep_Nested_Properties", func(b *testing.B) {
		query := `user.profile.settings.timestamps.created_at af session.events.latest.timestamp`
		ctx := map[string]any{
			"user": map[string]any{
				"profile": map[string]any{
					"settings": map[string]any{
						"timestamps": map[string]any{
							"created_at": "2024-07-09T22:12:01Z",
						},
					},
				},
			},
			"session": map[string]any{
				"events": map[string]any{
					"latest": map[string]any{
						"timestamp": int64(1720558320),
					},
				},
			},
		}

		b.ResetTimer()
		b.ReportAllocs()

		for range b.N {
			_, err := engine.Evaluate(query, ctx)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
