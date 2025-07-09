package rule

import (
	"testing"
)

// Test Error interface
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
	}
	
	for _, test := range errors {
		if test.err.Error() != test.expected {
			t.Errorf("Expected '%s', got '%s'", test.expected, test.err.Error())
		}
	}
}

// Test EngineError struct
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