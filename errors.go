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

	// ErrUnterminatedString indicates an unterminated string literal in the query.
	ErrUnterminatedString = &EngineError{"UNTERMINATED_STRING", "Unterminated string literal"}
	// ErrMissingOperator indicates missing operator between operands.
	ErrMissingOperator = &EngineError{"MISSING_OPERATOR", "Missing operator between operands"}
	// ErrInvalidSyntax indicates invalid query syntax.
	ErrInvalidSyntax = &EngineError{"INVALID_SYNTAX", "Invalid query syntax"}

	// ErrInvalidInOperand indicates IN operator used with non-array operand.
	ErrInvalidInOperand = &EngineError{"INVALID_IN_OPERAND", "IN operator requires an array operand"}
	ErrInvalidStringOp  = &EngineError{
		"INVALID_STRING_OP",
		"String operators (co/sw/ew) can only be used with string operands",
	}
	ErrInvalidPresenceOp = &EngineError{
		"INVALID_PRESENCE_OP",
		"Presence operator (pr) can only be used with identifiers or properties",
	}
	ErrEmptyQuery       = &EngineError{"EMPTY_QUERY", "Query cannot be empty"}
	ErrEmptyParentheses = &EngineError{"EMPTY_PARENTHESES", "Empty parentheses are not allowed"}
	ErrUnbalancedParens = &EngineError{"UNBALANCED_PARENTHESES", "Unbalanced parentheses"}
	ErrTrailingTokens   = &EngineError{"TRAILING_TOKENS", "Unexpected tokens after complete expression"}
)
