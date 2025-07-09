package rule

import (
	"testing"

	"github.com/nikunjy/rules"
	"github.com/stretchr/testify/assert"
)

// TestOriginalLibraryCrossTypeComparison tests the behavior of the original nikunjy/rules library
// for cross-type comparison in array membership operations.
//
// Key findings:
// 1. The original library uses STRICT type matching for "in" operations
// 2. int 42 != string "42" (returns false)
// 3. float 42.0 != int 42 (returns false)  
// 4. bool true != string "true" (returns false)
// 5. Even mixed arrays with matching types return false in some cases
// 6. Only exact type and value matches return true
//
// This means our implementation needs to match this strict behavior.
func TestOriginalLibraryCrossTypeComparison(t *testing.T) {
	tests := []struct {
		name     string
		rule     string
		context  map[string]any
		expected bool
	}{
		{
			name:     "int 42 in string array should return false",
			rule:     "x in [\"test\", \"42\", \"false\"]",
			context:  map[string]any{"x": 42},
			expected: false, // int 42 != string "42"
		},
		{
			name:     "string '42' in string array should return true",
			rule:     "x in [\"test\", \"42\", \"false\"]",
			context:  map[string]any{"x": "42"},
			expected: true,
		},
		{
			name:     "int 42 in mixed array with int should return false",
			rule:     "x in [\"test\", 42, \"false\"]",
			context:  map[string]any{"x": 42},
			expected: false, // Even with int in array, the library returns false
		},
		{
			name:     "string '42' in mixed array with int should return false",
			rule:     "x in [\"test\", 42, \"false\"]",
			context:  map[string]any{"x": "42"},
			expected: false, // string "42" != int 42
		},
		{
			name:     "bool true in string array should return false",
			rule:     "x in [\"test\", \"true\", \"false\"]",
			context:  map[string]any{"x": true},
			expected: false, // bool true != string "true"
		},
		{
			name:     "string 'true' in string array should return true",
			rule:     "x in [\"test\", \"true\", \"false\"]",
			context:  map[string]any{"x": "true"},
			expected: true,
		},
		{
			name:     "float 42.0 in int array should return false",
			rule:     "x in [1, 42, 100]",
			context:  map[string]any{"x": 42.0},
			expected: false, // float 42.0 != int 42
		},
		{
			name:     "int 42 in int array should return true",
			rule:     "x in [1, 42, 100]",
			context:  map[string]any{"x": 42},
			expected: true,
		},
		{
			name:     "simple int in int array test",
			rule:     "x in [42]",
			context:  map[string]any{"x": 42},
			expected: true,
		},
		{
			name:     "SPECIFIC CASE: int 42 in string array with '42' should return false",
			rule:     "x in [\"test\", \"42\", \"false\"]",
			context:  map[string]any{"x": 42},
			expected: false, // This is the specific case from the problem
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := rules.Evaluate(tt.rule, tt.context)
			assert.NoError(t, err, "Rule evaluation should not error")
			assert.Equal(t, tt.expected, result, "Result should match expected value")
		})
	}
}