package rule

import (
	"testing"
)

// TestQueryValidation tests the engine's query validation capabilities.
func TestQueryValidation(t *testing.T) {
	t.Run("BasicValidation", testBasicQueryValidation)
	t.Run("PropertyValidation", testPropertyValidation)
	t.Run("ParenthesesValidation", testParenthesesValidation)
	t.Run("LogicalValidation", testLogicalValidation)
	t.Run("EdgeCaseValidation", testEdgeCaseValidation)
	t.Run("ValidComplexQueries", testValidComplexQueries)
}

// TestSemanticValidation tests the engine's semantic validation capabilities.
func TestSemanticValidation(t *testing.T) {
	t.Run("InOperatorValidation", testInOperatorValidation)
	t.Run("StringOperatorValidation", testStringOperatorValidation)
	t.Run("PresenceOperatorValidation", testPresenceOperatorValidation)
	t.Run("EmptyQueryValidation", testEmptyQueryValidation)
	t.Run("ComplexValidQueries", testComplexValidQueries)
}
