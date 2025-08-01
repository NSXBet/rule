package rule

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNotInOperator(t *testing.T) {
	engine := NewEngine()

	testCases := []struct {
		name        string
		query       string
		context     map[string]any
		expected    bool
		expectError bool
	}{
		{
			name:     "Simple not in with string not in array",
			query:    `role not in ["admin", "guest"]`,
			context:  map[string]any{"role": "user"},
			expected: true,
		},
		{
			name:     "Simple not in with string in array",
			query:    `role not in ["admin", "guest"]`,
			context:  map[string]any{"role": "admin"},
			expected: false,
		},
		{
			name:     "Number not in array",
			query:    `score not in [100, 200, 300]`,
			context:  map[string]any{"score": 150},
			expected: true,
		},
		{
			name:     "Number in array with not in",
			query:    `score not in [100, 200, 300]`,
			context:  map[string]any{"score": 200},
			expected: false,
		},
		{
			name:     "Property not in array",
			query:    `user.role not in ["banned", "suspended"]`,
			context:  map[string]any{"user": map[string]any{"role": "active"}},
			expected: true,
		},
		{
			name:     "Property in array with not in",
			query:    `user.role not in ["banned", "suspended"]`,
			context:  map[string]any{"user": map[string]any{"role": "banned"}},
			expected: false,
		},
		{
			name:     "Empty array always returns true for not in",
			query:    `anything not in []`,
			context:  map[string]any{"anything": "value"},
			expected: true,
		},
		{
			name:     "Mixed type array with not in",
			query:    `value not in ["string", 42, true]`,
			context:  map[string]any{"value": "different"},
			expected: true,
		},
		{
			name:     "Mixed type array with match",
			query:    `value not in ["string", 42, true]`,
			context:  map[string]any{"value": 42},
			expected: false,
		},
		{
			name:     "Complex expression with not in and and",
			query:    `role not in ["banned"] and status eq "active"`,
			context:  map[string]any{"role": "user", "status": "active"},
			expected: true,
		},
		{
			name:     "Complex expression with not in failing",
			query:    `role not in ["banned"] and status eq "active"`,
			context:  map[string]any{"role": "banned", "status": "active"},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := engine.AddQuery(tc.query)
			require.NoError(t, err, "Failed to add query: %s", tc.query)

			result, err := engine.Evaluate(tc.query, tc.context)
			if tc.expectError {
				require.Error(t, err, "Expected error for query: %s", tc.query)
			} else {
				require.NoError(t, err, "Unexpected error for query: %s", tc.query)
				require.Equal(t, tc.expected, result, "Unexpected result for query: %s", tc.query)
			}
		})
	}
}

func TestNotInVsNotIn(t *testing.T) {
	engine := NewEngine()

	// Test that "not in" and "not (... in ...)" produce the same results
	testCases := []struct {
		name    string
		query1  string
		query2  string
		context map[string]any
	}{
		{
			name:    "String not in array",
			query1:  `role not in ["admin", "guest"]`,
			query2:  `not (role in ["admin", "guest"])`,
			context: map[string]any{"role": "user"},
		},
		{
			name:    "String in array",
			query1:  `role not in ["admin", "guest"]`,
			query2:  `not (role in ["admin", "guest"])`,
			context: map[string]any{"role": "admin"},
		},
		{
			name:    "Number not in array",
			query1:  `score not in [100, 200]`,
			query2:  `not (score in [100, 200])`,
			context: map[string]any{"score": 150},
		},
		{
			name:    "Number in array",
			query1:  `score not in [100, 200]`,
			query2:  `not (score in [100, 200])`,
			context: map[string]any{"score": 100},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err1 := engine.AddQuery(tc.query1)
			require.NoError(t, err1, "Failed to add query1: %s", tc.query1)

			err2 := engine.AddQuery(tc.query2)
			require.NoError(t, err2, "Failed to add query2: %s", tc.query2)

			result1, err1 := engine.Evaluate(tc.query1, tc.context)
			require.NoError(t, err1, "Failed to evaluate query1: %s", tc.query1)

			result2, err2 := engine.Evaluate(tc.query2, tc.context)
			require.NoError(t, err2, "Failed to evaluate query2: %s", tc.query2)

			require.Equal(
				t,
				result1,
				result2,
				"not in vs not (in) should be equivalent\nQuery1: %s\nQuery2: %s\nContext: %v",
				tc.query1,
				tc.query2,
				tc.context,
			)
		})
	}
}

func TestNotInValidation(t *testing.T) {
	engine := NewEngine()

	validQueries := []string{
		`field not in ["a", "b"]`,
		`user.role not in ["admin"]`,
		`score not in [1, 2, 3]`,
		`status not in []`,
	}

	for _, query := range validQueries {
		t.Run("Valid: "+query, func(t *testing.T) {
			err := engine.AddQuery(query)
			require.NoError(t, err, "Valid query should be accepted: %s", query)
		})
	}

	invalidQueries := []struct {
		query     string
		errorType *EngineError
	}{
		{`field not in "string"`, ErrInvalidInOperand},
		{`field not in 123`, ErrInvalidInOperand},
		{`field not in true`, ErrInvalidInOperand},
	}

	for _, tc := range invalidQueries {
		t.Run("Invalid: "+tc.query, func(t *testing.T) {
			err := engine.AddQuery(tc.query)
			require.Error(t, err, "Invalid query should be rejected: %s", tc.query)
			require.Contains(t, err.Error(), tc.errorType.Message, "Error should contain expected message")
		})
	}
}
