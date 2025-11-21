package rule

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

// RaceConditionError represents an error detected during race condition testing.
type RaceConditionError struct {
	GoroutineID int
	Iteration   int
	Query       string
	Context     D
	Expected    bool
	Actual      bool
}

func (e *RaceConditionError) Error() string {
	return fmt.Sprintf("Race condition: Goroutine %d, Iteration %d, Query: %s, Expected: %v, Actual: %v, Context: %+v",
		e.GoroutineID, e.Iteration, e.Query, e.Expected, e.Actual, e.Context)
}

// TestEngine_RaceCondition_BasicQuery tests for race conditions
// when evaluating the same query with different contexts concurrently.
func TestEngine_RaceCondition_BasicQuery(t *testing.T) {
	engine := NewEngine()

	// Pre-create the query we'll test
	query := `skin eq "betnacional"`

	// Pre-compile the query
	err := engine.AddQuery(query)
	require.NoError(t, err)

	// Number of concurrent goroutines
	numGoroutines := 1000
	iterationsPerGoroutine := 100
	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines*iterationsPerGoroutine)
	falseResults := make(chan int, numGoroutines*iterationsPerGoroutine)

	// Launch concurrent goroutines
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < iterationsPerGoroutine; j++ {
				// Each goroutine uses a unique context with different customerid
				// but same skin="betnacional" - should always return true
				context := D{
					"skin":       "betnacional",
					"customerid": goroutineID*10000 + j, // Unique customerid per iteration
				}

				result, err := engine.Evaluate(query, context)
				if err != nil {
					errors <- fmt.Errorf("evaluation error in goroutine %d, iteration %d: %w", goroutineID, j, err)
					continue
				}

				// The query should always return true since skin="betnacional"
				if !result {
					falseResults <- goroutineID*10000 + j
					errors <- &RaceConditionError{
						GoroutineID: goroutineID,
						Iteration:   j,
						Context:     context,
						Expected:    true,
						Actual:      result,
					}
				}
			}
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(errors)
	close(falseResults)

	// Check for race conditions
	falseCount := len(falseResults)
	if falseCount > 0 {
		t.Errorf("❌ RACE CONDITION DETECTED: %d evaluations returned false when they should return true", falseCount)
		t.Errorf("   This indicates the Engine.Evaluate is not thread-safe")

		// Print first few errors for debugging
		errorCount := 0
		for err := range errors {
			if errorCount < 10 {
				t.Errorf("   Error %d: %v", errorCount+1, err)
				errorCount++
			}
		}
	} else {
		t.Logf("✅ No race conditions detected in %d concurrent evaluations", numGoroutines*iterationsPerGoroutine)
	}

	require.Equal(t, 0, falseCount, "All evaluations should return true - race condition detected if any return false")
}

// TestEngine_RaceCondition_DifferentQueries tests for race conditions
// when evaluating different queries concurrently with different contexts.
func TestEngine_RaceCondition_DifferentQueries(t *testing.T) {
	engine := NewEngine()

	// Multiple different queries that should be evaluated concurrently
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
		require.NoError(t, err)
	}

	numGoroutines := 1000
	iterationsPerGoroutine := 5000
	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines*iterationsPerGoroutine)

	// Launch concurrent goroutines evaluating different queries
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < iterationsPerGoroutine; j++ {
				// Each goroutine evaluates a different query with different context
				queryIndex := (goroutineID + j) % len(queries)
				query := queries[queryIndex]

				// Create context based on query
				context := getContextForQuery(query, goroutineID, j)

				result, err := engine.Evaluate(query, context)
				if err != nil {
					errors <- fmt.Errorf("evaluation error in goroutine %d, iteration %d: %w", goroutineID, j, err)
					continue
				}

				expected := getExpectedResult(query, context)
				if result != expected {
					errors <- &RaceConditionError{
						GoroutineID: goroutineID,
						Iteration:   j,
						Query:       query,
						Context:     context,
						Expected:    expected,
						Actual:      result,
					}
				}
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	errorCount := 0
	for range errors {
		errorCount++
	}

	if errorCount > 0 {
		t.Errorf("❌ RACE CONDITION DETECTED: %d evaluations returned incorrect results", errorCount)
	} else {
		t.Logf("✅ No race conditions detected in %d concurrent mixed query evaluations", numGoroutines*iterationsPerGoroutine)
	}

	require.Equal(t, 0, errorCount, "All evaluations should return correct results")
}

// TestEngine_RaceCondition_SameQueryDifferentContexts tests the most common
// race condition scenario: same query, different contexts concurrently.
func TestEngine_RaceCondition_SameQueryDifferentContexts(t *testing.T) {
	engine := NewEngine()

	query := `skin eq "betnacional" and customerid eq 100`

	// Pre-compile the query
	err := engine.AddQuery(query)
	require.NoError(t, err)

	numGoroutines := 500
	iterationsPerGoroutine := 200
	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines*iterationsPerGoroutine)

	// Launch concurrent goroutines with different contexts
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < iterationsPerGoroutine; j++ {
				// Each goroutine uses a different customerid
				// Only customerid=100 should return true
				customerid := goroutineID*1000 + j
				context := D{
					"skin":       "betnacional",
					"customerid": customerid,
				}

				result, err := engine.Evaluate(query, context)
				if err != nil {
					errors <- fmt.Errorf("evaluation error in goroutine %d, iteration %d: %w", goroutineID, j, err)
					continue
				}

				expected := customerid == 100
				if result != expected {
					errors <- &RaceConditionError{
						GoroutineID: goroutineID,
						Iteration:   j,
						Query:       query,
						Context:     context,
						Expected:    expected,
						Actual:      result,
					}
				}
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	errorCount := 0
	for range errors {
		errorCount++
	}

	if errorCount > 0 {
		t.Errorf("❌ RACE CONDITION DETECTED: %d evaluations returned incorrect results", errorCount)
	} else {
		t.Logf("✅ No race conditions detected in %d concurrent evaluations with different contexts", numGoroutines*iterationsPerGoroutine)
	}

	require.Equal(t, 0, errorCount, "All evaluations should return correct results")
}

// Helper functions

func getContextForQuery(query string, goroutineID, iteration int) D {
	context := D{}

	// Determine skin based on query
	if strings.Contains(query, "betnacional") {
		context["skin"] = "betnacional"
	} else if strings.Contains(query, "betdev") {
		context["skin"] = "betdev"
	} else if strings.Contains(query, "webdev") {
		context["skin"] = "webdev"
	} else if strings.Contains(query, "sandbox") {
		context["skin"] = "sandbox"
	}

	// Determine customerid based on query
	if strings.Contains(query, "customerid in [0]") {
		context["customerid"] = 0
	} else if strings.Contains(query, "customerid in [171, 273, 612, 179, 504]") {
		// Return one of the valid IDs
		validIDs := []int{171, 273, 612, 179, 504}
		context["customerid"] = validIDs[(goroutineID+iteration)%len(validIDs)]
	} else {
		// For other queries, return a unique ID
		context["customerid"] = goroutineID*10000 + iteration
	}

	return context
}

//nolint:gocognit
func getExpectedResult(query string, context D) bool {
	skin, ok := context["skin"].(string)
	if !ok {
		return false
	}

	customerid, ok := context["customerid"].(int)
	if !ok {
		return false
	}

	// Evaluate expected result based on query
	if query == `skin eq "betnacional"` {
		return skin == "betnacional"
	}

	if query == `skin eq "betnacional" and customerid in [0]` {
		return skin == "betnacional" && customerid == 0
	}

	if query == `skin eq "betdev" and customerid in [171, 273, 612, 179, 504]` {
		validIDs := []int{171, 273, 612, 179, 504}
		valid := false

		for _, id := range validIDs {
			if customerid == id {
				valid = true
				break
			}
		}

		return skin == "betdev" && valid
	}

	if query == `skin eq "webdev" and customerid in [171, 273, 612, 179, 504]` {
		validIDs := []int{171, 273, 612, 179, 504}
		valid := false
		for _, id := range validIDs {
			if customerid == id {
				valid = true
				break
			}
		}
		return skin == "webdev" && valid
	}

	if query == `skin eq "sandbox" and customerid in [0]` {
		return skin == "sandbox" && customerid == 0
	}

	return false
}
