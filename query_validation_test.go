package rule

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func testBasicQueryValidation(t *testing.T) {
	engine := NewEngine()

	testCases := []struct {
		name        string
		query       string
		expectError bool
		errorType   *EngineError
	}{
		{
			name:        "Valid simple query",
			query:       `x eq 10`,
			expectError: false,
		},
		{
			name:        "Valid complex query",
			query:       `user.age gt 18 and status eq "active"`,
			expectError: false,
		},
		{
			name:        "Missing operator between operands",
			query:       `something "else"`,
			expectError: true,
			errorType:   ErrMissingOperator,
		},
		{
			name:        "Unterminated string literal",
			query:       `name eq "unclosed`,
			expectError: true,
			errorType:   ErrUnterminatedString,
		},
		{
			name:        "Missing right operand",
			query:       `x eq`,
			expectError: true,
		},
		{
			name:        "Missing left operand",
			query:       `eq "value"`,
			expectError: true,
		},
		{
			name:        "Multiple adjacent values",
			query:       `x y z`,
			expectError: true,
			errorType:   ErrMissingOperator,
		},
		{
			name:        "Number followed by string without operator",
			query:       `123 "test"`,
			expectError: true,
			errorType:   ErrMissingOperator,
		},
	}

	runTestCases(t, engine, testCases)
}

func testPropertyValidation(t *testing.T) {
	engine := NewEngine()

	testCases := []struct {
		name        string
		query       string
		expectError bool
		errorType   *EngineError
	}{
		{
			name:        "Trailing dot in property",
			query:       `something.else. in ["a"]`,
			expectError: true,
		},
		{
			name:        "Double dot in property",
			query:       `something..else eq "test"`,
			expectError: true,
		},
		{
			name:        "Property ending with dot",
			query:       `something. eq "test"`,
			expectError: true,
		},
		{
			name:        "Property starting with dot",
			query:       `.something eq "test"`,
			expectError: true,
		},
		{
			name:        "Trailing dot with presence operator",
			query:       `user.profile. pr`,
			expectError: true,
		},
	}

	runTestCases(t, engine, testCases)
}

func testParenthesesValidation(t *testing.T) {
	engine := NewEngine()

	testCases := []struct {
		name        string
		query       string
		expectError bool
		errorType   *EngineError
	}{
		{
			name:        "Missing closing parenthesis",
			query:       `(something eq "test"`,
			expectError: true,
		},
		{
			name:        "Extra closing parenthesis",
			query:       `something eq "test")`,
			expectError: true,
			errorType:   ErrTrailingTokens,
		},
		{
			name:        "Multiple missing closing parentheses",
			query:       `((something eq "test")`,
			expectError: true,
		},
		{
			name:        "Multiple extra closing parentheses",
			query:       `(something eq "test"))`,
			expectError: true,
			errorType:   ErrTrailingTokens,
		},
		{
			name:        "Parentheses in wrong order",
			query:       `)(something eq "test"(`,
			expectError: true,
		},
		{
			name:        "Empty parentheses",
			query:       `(())`,
			expectError: true,
			errorType:   ErrEmptyParentheses,
		},
	}

	runTestCases(t, engine, testCases)
}

func testLogicalValidation(t *testing.T) {
	engine := NewEngine()

	testCases := []struct {
		name        string
		query       string
		expectError bool
		errorType   *EngineError
	}{
		{
			name:        "Missing right operand after OR",
			query:       `((something eq "12") or )`,
			expectError: true,
		},
		{
			name:        "Missing left operand before OR",
			query:       `( or something eq "test")`,
			expectError: true,
		},
		{
			name:        "Missing right operand after AND",
			query:       `something eq "test" and`,
			expectError: true,
		},
		{
			name:        "Missing left operand before AND",
			query:       `and something eq "test"`,
			expectError: true,
		},
		{
			name:        "Double OR operator",
			query:       `something eq "test" or or`,
			expectError: true,
		},
		{
			name:        "Double AND operator",
			query:       `something and and "test"`,
			expectError: true,
		},
		{
			name:        "Incomplete and expression",
			query:       `x eq "value" and`,
			expectError: true,
		},
	}

	runTestCases(t, engine, testCases)
}

func testEdgeCaseValidation(t *testing.T) {
	engine := NewEngine()

	testCases := []struct {
		name        string
		query       string
		expectError bool
		errorType   *EngineError
	}{
		{
			name:        "Multiple missing closing in nested",
			query:       `(((something eq "test"`,
			expectError: true,
		},
		{
			name:        "Multiple extra closing in nested",
			query:       `something eq "test")))`,
			expectError: true,
			errorType:   ErrTrailingTokens,
		},
		{
			name:        "Incomplete expression in nested parentheses",
			query:       `((something eq "test") and)`,
			expectError: true,
		},
		{
			name:        "Empty nested parentheses",
			query:       `(something eq "test" and ())`,
			expectError: true,
			errorType:   ErrEmptyParentheses,
		},
		{
			name:        "Property dot + missing operand",
			query:       `something. and (missing eq)`,
			expectError: true,
		},
		{
			name:        "Property double dot + missing paren",
			query:       `(something..else eq "test" or`,
			expectError: true,
		},
		{
			name:        "Unterminated string + paren",
			query:       `something eq "unclosed and )`,
			expectError: true,
			errorType:   ErrUnterminatedString,
		},
		{
			name:        "Multiple validation issues",
			query:       `(something in "string") or .`,
			expectError: true,
		},
		{
			name:        "Operator without operands",
			query:       `eq`,
			expectError: true,
		},
		{
			name:        "Only operator",
			query:       `and`,
			expectError: true,
		},
		{
			name:        "Nested empty expressions",
			query:       `(()) and (())`,
			expectError: true,
			errorType:   ErrEmptyParentheses,
		},
		{
			name:        "Array with missing closing",
			query:       `field in [1, 2, 3`,
			expectError: true,
		},
		{
			name:        "Array with extra closing",
			query:       `field in [1, 2, 3]]`,
			expectError: true,
			errorType:   ErrTrailingTokens,
		},
		{
			name:        "Field with hyphen (tokenized as separate)",
			query:       `field-name eq "test"`,
			expectError: true,
		},
		{
			name:        "Multiple operators without operands",
			query:       `eq ne lt gt`,
			expectError: true,
		},
		{
			name:        "Operator sandwich",
			query:       `field eq ne "test"`,
			expectError: true,
		},
		{
			name:        "Excessive whitespace with missing operand",
			query:       `field    eq     `,
			expectError: true,
		},
		{
			name:        "Tabs and spaces mixing with errors",
			query:       "\tfield\t\teq\t\t",
			expectError: true,
		},
	}

	runTestCases(t, engine, testCases)
}

func testValidComplexQueries(t *testing.T) {
	engine := NewEngine()

	testCases := []struct {
		name        string
		query       string
		expectError bool
		errorType   *EngineError
	}{
		{
			name:        "Complex nested valid expression",
			query:       `((user.age gt 18 and user.active eq true) or (user.premium eq true)) and user.email pr`,
			expectError: false,
		},
		{
			name:        "Deep nesting valid",
			query:       `(((field eq "value")))`,
			expectError: false,
		},
		{
			name:        "Complex array operations",
			query:       `user.roles in ["admin", "user", "guest"] and user.permissions in []`,
			expectError: false,
		},
		{
			name:        "Mixed operators valid",
			query:       `name co "John" and age gt 18 and active eq true and email pr`,
			expectError: false,
		},
	}

	runTestCases(t, engine, testCases)
}

func runTestCases(t *testing.T, engine *Engine, testCases []struct {
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
