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

	// ErrQuantifierRequiresParens indicates a quantifier operator used without parenthesized sub-expression.
	ErrQuantifierRequiresParens = &EngineError{
		"QUANTIFIER_REQUIRES_PARENS",
		"Quantifier operators (any/all/none) require a parenthesized sub-expression",
	}
	// ErrInvalidQuantifierTarget indicates a quantifier used on a non-identifier/property operand.
	ErrInvalidQuantifierTarget = &EngineError{
		"INVALID_QUANTIFIER_TARGET",
		"Quantifier operators (any/all/none) can only be used with identifiers or properties",
	}

	// ErrWhereRequiresParens indicates where used without parenthesized predicate.
	ErrWhereRequiresParens = &EngineError{
		"WHERE_REQUIRES_PARENS",
		"where operator requires a parenthesized predicate",
	}
	// ErrWhereRequiresLengthOrQuantifier indicates where without .length or quantifier after.
	ErrWhereRequiresLengthOrQuantifier = &EngineError{
		"WHERE_REQUIRES_LENGTH_OR_QUANTIFIER",
		"where expression must be followed by .length or a quantifier (any/all/none)",
	}
	// ErrEmptyWherePredicate indicates where with empty predicate.
	ErrEmptyWherePredicate = &EngineError{
		"EMPTY_WHERE_PREDICATE",
		"where predicate cannot be empty",
	}
	// ErrWhereRequiresLength indicates where followed by non-length property.
	ErrWhereRequiresLength = &EngineError{
		"WHERE_REQUIRES_LENGTH",
		"where can only be followed by .length",
	}

	// ErrUnknownMathFunction indicates an unsupported function call.
	ErrUnknownMathFunction = &EngineError{"UNKNOWN_MATH_FUNCTION", "Unknown math function"}
	// ErrInvalidMathArguments indicates incorrect math function arguments.
	ErrInvalidMathArguments = &EngineError{"INVALID_MATH_ARGUMENTS", "Invalid math function arguments"}
	// ErrInvalidMathPrecision indicates round precision outside the supported integer range.
	ErrInvalidMathPrecision = &EngineError{
		"INVALID_MATH_PRECISION",
		"round precision must be an integer literal from 0 to 1000",
	}
	// ErrInvalidMathOperand indicates a nonnumeric math function argument.
	ErrInvalidMathOperand = &EngineError{"INVALID_MATH_OPERAND", "Math functions require a numeric operand"}
	// ErrInvalidDecimal indicates an invalid or unsupported decimal representation.
	ErrInvalidDecimal = &EngineError{"INVALID_DECIMAL", "Invalid decimal value"}
	// ErrFloatMathOverflow indicates a float math value or scale outside float64 range.
	ErrFloatMathOverflow = &EngineError{"FLOAT_MATH_OVERFLOW", "Float math exceeds float64 range"}
)
