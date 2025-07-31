package rule

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func testInOperatorValidation(t *testing.T) {
	engine := NewEngine()

	testCases := []struct {
		name        string
		query       string
		expectError bool
		errorType   *EngineError
	}{
		{
			name:        "IN with string literal",
			query:       `field in "string"`,
			expectError: true,
			errorType:   ErrInvalidInOperand,
		},
		{
			name:        "IN with number literal",
			query:       `field in 123`,
			expectError: true,
			errorType:   ErrInvalidInOperand,
		},
		{
			name:        "IN with boolean literal",
			query:       `field in true`,
			expectError: true,
			errorType:   ErrInvalidInOperand,
		},
		{
			name:        "IN with array (valid)",
			query:       `field in [1,2,3]`,
			expectError: false,
		},
		{
			name:        "IN with identifier (valid - runtime check)",
			query:       `field in otherField`,
			expectError: false,
		},
		{
			name:        "IN with property (valid - runtime check)",
			query:       `field in user.roles`,
			expectError: false,
		},
		{
			name:        "IN with empty array (valid)",
			query:       `field in []`,
			expectError: false,
		},
	}

	runSemanticTestCases(t, engine, testCases)
}

func testStringOperatorValidation(t *testing.T) {
	engine := NewEngine()

	testCases := []struct {
		name        string
		query       string
		expectError bool
		errorType   *EngineError
	}{
		{
			name:        "Contains with number left operand",
			query:       `123 co "test"`,
			expectError: true,
			errorType:   ErrInvalidStringOp,
		},
		{
			name:        "Contains with boolean left operand",
			query:       `true co "test"`,
			expectError: true,
			errorType:   ErrInvalidStringOp,
		},
		{
			name:        "Starts with number left operand",
			query:       `123 sw "1"`,
			expectError: true,
			errorType:   ErrInvalidStringOp,
		},
		{
			name:        "Ends with number right operand",
			query:       `"test" ew 123`,
			expectError: true,
			errorType:   ErrInvalidStringOp,
		},
		{
			name:        "Contains with string operands (valid)",
			query:       `"hello" co "ell"`,
			expectError: false,
		},
		{
			name:        "Contains with identifier (valid - runtime check)",
			query:       `name co "test"`,
			expectError: false,
		},
	}

	runSemanticTestCases(t, engine, testCases)
}

func testPresenceOperatorValidation(t *testing.T) {
	engine := NewEngine()

	testCases := []struct {
		name        string
		query       string
		expectError bool
		errorType   *EngineError
	}{
		{
			name:        "Presence with string literal",
			query:       `"string" pr`,
			expectError: true,
			errorType:   ErrInvalidPresenceOp,
		},
		{
			name:        "Presence with number literal",
			query:       `123 pr`,
			expectError: true,
			errorType:   ErrInvalidPresenceOp,
		},
		{
			name:        "Presence with boolean literal",
			query:       `true pr`,
			expectError: true,
			errorType:   ErrInvalidPresenceOp,
		},
		{
			name:        "Presence with identifier (valid)",
			query:       `field pr`,
			expectError: false,
		},
		{
			name:        "Presence with property (valid)",
			query:       `user.email pr`,
			expectError: false,
		},
	}

	runSemanticTestCases(t, engine, testCases)
}

func testEmptyQueryValidation(t *testing.T) {
	engine := NewEngine()

	testCases := []struct {
		name        string
		query       string
		expectError bool
		errorType   *EngineError
	}{
		{
			name:        "Empty query",
			query:       ``,
			expectError: true,
			errorType:   ErrEmptyQuery,
		},
		{
			name:        "Whitespace only query",
			query:       `   `,
			expectError: true,
			errorType:   ErrEmptyQuery,
		},
		{
			name:        "Tab and newline only",
			query:       "\t\n  ",
			expectError: true,
			errorType:   ErrEmptyQuery,
		},
	}

	runSemanticTestCases(t, engine, testCases)
}

func testComplexValidQueries(t *testing.T) {
	engine := NewEngine()

	testCases := []struct {
		name        string
		query       string
		expectError bool
		errorType   *EngineError
	}{
		{
			name:        "Complex valid query with multiple operators",
			query:       `user.age gt 18 and user.name co "John" and user.roles in ["admin", "user"]`,
			expectError: false,
		},
		{
			name:        "Nested parentheses with valid operators",
			query:       `(user.active pr and user.age ge 21) or (user.premium eq true and user.trial_days in [7, 14, 30])`,
			expectError: false,
		},
	}

	runSemanticTestCases(t, engine, testCases)
}

func runSemanticTestCases(t *testing.T, engine *Engine, testCases []struct {
	name        string
	query       string
	expectError bool
	errorType   *EngineError
},
) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := engine.AddQuery(tc.query)
			if !tc.expectError {
				require.NoError(t, err, "query=%q", tc.query)
				return
			}

			require.Error(t, err, "Expected error for query=%q", tc.query)

			// Check for specific error type if provided
			if tc.errorType != nil {
				require.Contains(t, err.Error(), tc.errorType.Message, "query=%q", tc.query)
			}
		})
	}
}
