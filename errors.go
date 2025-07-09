package rule

type EngineError struct {
	Code    string
	Message string
}

func (e *EngineError) Error() string {
	return e.Message
}

var (
	ErrInvalidNode       = &EngineError{"INVALID_NODE", "Invalid AST node type"}
	ErrInvalidLiteral    = &EngineError{"INVALID_LITERAL", "Invalid literal value"}
	ErrInvalidOperator   = &EngineError{"INVALID_OPERATOR", "Invalid operator"}
	ErrAttributeNotFound = &EngineError{
		"ATTRIBUTE_NOT_FOUND",
		"Attribute not found in context",
	}
	ErrInvalidNestedAttribute = &EngineError{
		"INVALID_NESTED_ATTRIBUTE",
		"Invalid nested attribute access",
	}
	ErrParseError      = &EngineError{"PARSE_ERROR", "Failed to parse rule"}
	ErrEvaluationError = &EngineError{"EVALUATION_ERROR", "Failed to evaluate rule"}
	ErrRuleNotFound    = &EngineError{"RULE_NOT_FOUND", "Rule not found - use AddQuery to pre-compile rule"}
)
