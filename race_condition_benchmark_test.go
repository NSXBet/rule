package rule

import (
	"sync"
	"testing"
)

// BenchmarkEngine_Sequential benchmarks sequential evaluation
// This serves as a baseline for comparison with concurrent benchmarks
func BenchmarkEngine_Sequential(b *testing.B) {
	engine := NewEngine()
	query := `skin eq "betnacional"`

	err := engine.AddQuery(query)
	if err != nil {
		b.Fatalf("Failed to add query: %v", err)
	}

	context := D{
		"skin":       "betnacional",
		"customerid": 100,
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		result, err := engine.Evaluate(query, context)
		if err != nil {
			b.Fatalf("Evaluation error: %v", err)
		}
		if !result {
			b.Fatalf("Expected true, got false")
		}
	}
}

// BenchmarkEngine_Concurrent_SameQuery benchmarks concurrent evaluation
// of the same query with different contexts
// NOTE: This benchmark will show race conditions if run with -race flag
func BenchmarkEngine_Concurrent_SameQuery(b *testing.B) {
	engine := NewEngine()
	query := `skin eq "betnacional"`

	err := engine.AddQuery(query)
	if err != nil {
		b.Fatalf("Failed to add query: %v", err)
	}

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		context := D{
			"skin":       "betnacional",
			"customerid": 100,
		}

		for pb.Next() {
			result, err := engine.Evaluate(query, context)
			if err != nil {
				b.Errorf("Evaluation error: %v", err)
				continue
			}
			// NOTE: Due to race conditions, result might be incorrect
			// This benchmark measures performance, not correctness
			_ = result
		}
	})
}

// BenchmarkEngine_Concurrent_DifferentQueries benchmarks concurrent evaluation
// of different queries with different contexts
// NOTE: This benchmark will show race conditions if run with -race flag
func BenchmarkEngine_Concurrent_DifferentQueries(b *testing.B) {
	engine := NewEngine()

	queries := []string{
		`skin eq "betnacional"`,
		`skin eq "betdev" and customerid in [171, 273, 612, 179, 504]`,
		`skin eq "webdev" and customerid in [171, 273, 612, 179, 504]`,
		`skin eq "sandbox" and customerid in [0]`,
		`skin eq "betnacional" and customerid in [0]`,
	}

	// Pre-compile all queries
	for _, query := range queries {
		err := engine.AddQuery(query)
		if err != nil {
			b.Fatalf("Failed to add query %q: %v", query, err)
		}
	}

	b.ResetTimer()
	b.ReportAllocs()

	var queryIndex int
	var mu sync.Mutex

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.Lock()
			currentIndex := queryIndex
			queryIndex = (queryIndex + 1) % len(queries)
			mu.Unlock()

			query := queries[currentIndex]
			context := getContextForQuery(query, currentIndex, 0)

			result, err := engine.Evaluate(query, context)
			if err != nil {
				b.Errorf("Evaluation error: %v", err)
				continue
			}
			// NOTE: Due to race conditions, result might be incorrect
			// This benchmark measures performance, not correctness
			_ = result
		}
	})
}

// BenchmarkEngine_Concurrent_ComplexQuery benchmarks concurrent evaluation
// of a complex query with different contexts
// NOTE: This benchmark will show race conditions if run with -race flag
func BenchmarkEngine_Concurrent_ComplexQuery(b *testing.B) {
	engine := NewEngine()
	query := `skin eq "betnacional" and customerid eq 100`

	err := engine.AddQuery(query)
	if err != nil {
		b.Fatalf("Failed to add query: %v", err)
	}

	b.ResetTimer()
	b.ReportAllocs()

	var customerID int
	var mu sync.Mutex

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.Lock()
			currentID := customerID
			customerID++
			mu.Unlock()

			context := D{
				"skin":       "betnacional",
				"customerid": currentID,
			}

			result, err := engine.Evaluate(query, context)
			if err != nil {
				b.Errorf("Evaluation error: %v", err)
				continue
			}
			// NOTE: Due to race conditions, result might be incorrect
			// Only customerid=100 should return true, but race conditions may cause wrong results
			_ = result
		}
	})
}

// BenchmarkEngine_Concurrent_MixedOperations benchmarks concurrent evaluation
// with mixed simple and complex queries
// NOTE: This benchmark will show race conditions if run with -race flag
func BenchmarkEngine_Concurrent_MixedOperations(b *testing.B) {
	engine := NewEngine()

	queries := []string{
		`skin eq "betnacional"`,
		`skin eq "betnacional" and customerid eq 100`,
		`skin eq "betdev" and customerid in [171, 273, 612, 179, 504]`,
	}

	// Pre-compile all queries
	for _, query := range queries {
		err := engine.AddQuery(query)
		if err != nil {
			b.Fatalf("Failed to add query %q: %v", query, err)
		}
	}

	b.ResetTimer()
	b.ReportAllocs()

	var queryIndex int
	var mu sync.Mutex

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.Lock()
			currentIndex := queryIndex
			queryIndex = (queryIndex + 1) % len(queries)
			mu.Unlock()

			query := queries[currentIndex]
			context := D{
				"skin":       "betnacional",
				"customerid": 100,
			}

			result, err := engine.Evaluate(query, context)
			if err != nil {
				b.Errorf("Evaluation error: %v", err)
				continue
			}
			// NOTE: Due to race conditions, result might be incorrect
			// This benchmark measures performance, not correctness
			_ = result
		}
	})
}

// BenchmarkEngine_Concurrent_WithWaitGroup benchmarks concurrent evaluation
// using explicit goroutines and WaitGroup (similar to race condition tests)
// NOTE: This benchmark will show race conditions if run with -race flag
func BenchmarkEngine_Concurrent_WithWaitGroup(b *testing.B) {
	engine := NewEngine()
	query := `skin eq "betnacional"`

	err := engine.AddQuery(query)
	if err != nil {
		b.Fatalf("Failed to add query: %v", err)
	}

	numGoroutines := 10
	iterationsPerGoroutine := b.N / numGoroutines
	if iterationsPerGoroutine == 0 {
		iterationsPerGoroutine = 1
	}

	b.ResetTimer()
	b.ReportAllocs()

	var wg sync.WaitGroup
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < iterationsPerGoroutine; j++ {
				context := D{
					"skin":       "betnacional",
					"customerid": goroutineID*10000 + j,
				}

				result, err := engine.Evaluate(query, context)
				if err != nil {
					b.Errorf("Evaluation error: %v", err)
					// Continue to next iteration
				} else {
					// NOTE: Due to race conditions, result might be incorrect
					// This benchmark measures performance, not correctness
					_ = result
				}
			}
		}(i)
	}

	wg.Wait()
}

