package rule

import (
	"testing"
)

// Test Error interface.
func TestErrorInterface(t *testing.T) {
	err := ErrInvalidNode
	if err.Error() != "Invalid AST node type" {
		t.Errorf("Expected 'Invalid AST node type', got '%s'", err.Error())
	}

	// Test all error types
	errors := []struct {
		err      *EngineError
		expected string
	}{
		{ErrInvalidNode, "Invalid AST node type"},
		{ErrInvalidLiteral, "Invalid literal value"},
		{ErrInvalidOperator, "Invalid operator"},
		{ErrAttributeNotFound, "Attribute not found in context"},
		{ErrInvalidNestedAttribute, "Invalid nested attribute access"},
		{ErrParseError, "Failed to parse rule"},
		{ErrEvaluationError, "Failed to evaluate rule"},
		{ErrRuleNotFound, "Rule not found - use AddQuery to pre-compile rule"},
		{ErrUnterminatedString, "Unterminated string literal"},
		{ErrMissingOperator, "Missing operator between operands"},
		{ErrInvalidSyntax, "Invalid query syntax"},
		{ErrInvalidInOperand, "IN operator requires an array operand"},
		{ErrInvalidStringOp, "String operators (co/sw/ew) can only be used with string operands"},
		{ErrInvalidPresenceOp, "Presence operator (pr) can only be used with identifiers or properties"},
		{ErrEmptyQuery, "Query cannot be empty"},
		{ErrEmptyParentheses, "Empty parentheses are not allowed"},
		{ErrUnbalancedParens, "Unbalanced parentheses"},
		{ErrTrailingTokens, "Unexpected tokens after complete expression"},
	}

	for _, test := range errors {
		if test.err.Error() != test.expected {
			t.Errorf("Expected '%s', got '%s'", test.expected, test.err.Error())
		}
	}
}

// Test EngineError struct.
func TestEngineError(t *testing.T) {
	customErr := &EngineError{
		Code:    "CUSTOM_ERROR",
		Message: "This is a custom error message",
	}
	if customErr.Code != "CUSTOM_ERROR" {
		t.Error("Custom error code incorrect")
	}

	if customErr.Message != "This is a custom error message" {
		t.Error("Custom error message incorrect")
	}

	if customErr.Error() != "This is a custom error message" {
		t.Error("Custom error Error() method incorrect")
	}
}
