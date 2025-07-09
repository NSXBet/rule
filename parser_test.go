package rule

import (
	"testing"
)

// Test parser peek function.
func TestParserPeek(t *testing.T) {
	lexer := NewLexer("x eq 1")
	tokens := lexer.Tokenize()
	parser := NewParser(tokens)

	// Current should be IDENTIFIER
	if parser.curToken.Type != IDENTIFIER {
		t.Error("Expected current token to be IDENTIFIER")
	}

	// Peek should be EQ
	if parser.peek().Type != EQ {
		t.Error("Expected peek token to be EQ")
	}

	// Advance to end
	for parser.curToken.Type != EOF {
		parser.advance()
	}

	// Peek should be EOF
	if parser.peek().Type != EOF {
		t.Error("Expected peek token to be EOF")
	}
}

// Test parser expect with error conditions.
func TestParserExpectError(t *testing.T) {
	lexer := NewLexer("x")
	tokens := lexer.Tokenize()
	parser := NewParser(tokens)

	// Expect a different token type than what we have
	err := parser.expect(NUMBER)
	if err == nil {
		t.Error("Expected error when expecting wrong token type")
	}
}

// Test parser error conditions.
func TestParserErrorConditions(t *testing.T) {
	// Test parseOrExpression with error
	lexer := NewLexer("x or @")
	tokens := lexer.Tokenize()
	parser := NewParser(tokens)

	// This should trigger an error path
	_, err := parser.parseOrExpression()
	if err == nil {
		t.Error("Expected error in parseOrExpression")
	}

	// Test parseAndExpression with error
	lexer = NewLexer("x and @")
	tokens = lexer.Tokenize()
	parser = NewParser(tokens)

	_, err = parser.parseAndExpression()
	if err == nil {
		t.Error("Expected error in parseAndExpression")
	}

	// Test parseNotExpression with error
	lexer = NewLexer("not @")
	tokens = lexer.Tokenize()
	parser = NewParser(tokens)

	_, err = parser.parseNotExpression()
	if err == nil {
		t.Error("Expected error in parseNotExpression")
	}
}

// Test parseArray with error conditions.
func TestParseArrayError(t *testing.T) {
	// Test array with invalid token
	lexer := NewLexer("[x]") // Identifier not allowed in array
	tokens := lexer.Tokenize()
	parser := NewParser(tokens)

	_, err := parser.parseArray()
	if err == nil {
		t.Error("Expected error parsing array with invalid token")
	}

	// Test array without closing bracket
	lexer = NewLexer("[1, 2, 3")
	tokens = lexer.Tokenize()
	parser = NewParser(tokens)

	_, err = parser.parseArray()
	if err == nil {
		t.Error("Expected error parsing array without closing bracket")
	}
}

// Test parsePrimaryExpression with error conditions.
func TestParsePrimaryExpressionError(t *testing.T) {
	// Test with unexpected token
	lexer := NewLexer("eq") // Operator without operands
	tokens := lexer.Tokenize()
	parser := NewParser(tokens)

	_, err := parser.parsePrimaryExpression()
	if err == nil {
		t.Error("Expected error parsing unexpected token")
	}

	// Test parenthesized expression without closing paren
	lexer = NewLexer("(x eq 1")
	tokens = lexer.Tokenize()
	parser = NewParser(tokens)

	_, err = parser.parsePrimaryExpression()
	if err == nil {
		t.Error("Expected error parsing expression without closing paren")
	}
}

// Test parseIdentifierOrProperty with error.
func TestParseIdentifierOrPropertyError(t *testing.T) {
	// Test property with invalid identifier after dot
	lexer := NewLexer("x.eq") // eq is a keyword, not identifier
	tokens := lexer.Tokenize()
	parser := NewParser(tokens)

	_, err := parser.parseIdentifierOrProperty()
	if err == nil {
		t.Error("Expected error parsing property with invalid identifier")
	}
}

// Test parser with complex expressions.
func TestParserComplexExpressions(t *testing.T) {
	tests := []string{
		"x eq 1",
		"x eq 1 and y gt 2",
		"x eq 1 or y gt 2",
		"not (x eq 1)",
		"x in [1, 2, 3]",
		"x.y eq z.w",
		"(x eq 1 and y gt 2) or (z lt 5)",
		"name co \"John\" and age ge 18",
		"status in [\"active\", \"pending\"] and verified eq true",
	}

	for _, test := range tests {
		lexer := NewLexer(test)
		tokens := lexer.Tokenize()
		parser := NewParser(tokens)

		_, err := parser.Parse()
		if err != nil {
			t.Errorf("Failed to parse %q: %v", test, err)
		}
	}
}

// Test parser with invalid expressions.
func TestParserInvalidExpressions(t *testing.T) {
	tests := []string{
		"",
		"x eq",
		"eq 1",
		"x eq 1 and",
		"x eq 1 or",
		"not",
		"(x eq 1",
		"x eq 1 and and y gt 2",
		"x in [1, 2",
		"x.eq",
		"x..y",
		"x.123",
	}

	for _, test := range tests {
		lexer := NewLexer(test)
		tokens := lexer.Tokenize()
		parser := NewParser(tokens)

		_, err := parser.Parse()
		if err == nil {
			t.Errorf("Expected error parsing invalid expression %q", test)
		}
	}
}

// Test parser advance function.
func TestParserAdvance(t *testing.T) {
	lexer := NewLexer("x eq 1")
	tokens := lexer.Tokenize()
	parser := NewParser(tokens)

	// Should start at IDENTIFIER
	if parser.curToken.Type != IDENTIFIER {
		t.Error("Expected to start at IDENTIFIER")
	}

	// Advance to EQ
	parser.advance()

	if parser.curToken.Type != EQ {
		t.Error("Expected to advance to EQ")
	}

	// Advance to NUMBER
	parser.advance()

	if parser.curToken.Type != NUMBER {
		t.Error("Expected to advance to NUMBER")
	}

	// Advance to EOF
	parser.advance()

	if parser.curToken.Type != EOF {
		t.Error("Expected to advance to EOF")
	}

	// Advance beyond EOF should stay at EOF
	parser.advance()

	if parser.curToken.Type != EOF {
		t.Error("Expected to stay at EOF")
	}
}

// Test parser expect with success.
func TestParserExpectSuccess(t *testing.T) {
	lexer := NewLexer("x eq 1")
	tokens := lexer.Tokenize()
	parser := NewParser(tokens)

	// Expect IDENTIFIER
	err := parser.expect(IDENTIFIER)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Should now be at EQ
	if parser.curToken.Type != EQ {
		t.Error("Expected to be at EQ after expecting IDENTIFIER")
	}
}

// Test parser isComparisonOperator.
func TestParserIsComparisonOperator(t *testing.T) {
	lexer := NewLexer("")
	tokens := lexer.Tokenize()
	parser := NewParser(tokens)

	comparisonOps := []TokenType{
		EQ, NE, LT, GT, LE, GE, CO, SW, EW, IN, PR, EQUALS, NOT_EQUALS,
	}

	for _, op := range comparisonOps {
		if !parser.isComparisonOperator(op) {
			t.Errorf("Expected %v to be a comparison operator", op)
		}
	}

	nonComparisonOps := []TokenType{
		IDENTIFIER, STRING, NUMBER, BOOLEAN, AND, OR, NOT,
		PAREN_OPEN, PAREN_CLOSE, ARRAY_START, ARRAY_END, DOT, COMMA, EOF,
	}

	for _, op := range nonComparisonOps {
		if parser.isComparisonOperator(op) {
			t.Errorf("Expected %v to not be a comparison operator", op)
		}
	}
}

// Test standalone ParseRule function.
func TestParseRule(t *testing.T) {
	ast, err := ParseRule("x eq 1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if ast == nil {
		t.Error("Expected AST to not be nil")
	}

	// Test with invalid rule
	_, err = ParseRule("eq 1")
	if err == nil {
		t.Error("Expected error for invalid syntax")
	}
}

// Test parser with nested parentheses.
func TestParserNestedParentheses(t *testing.T) {
	tests := []string{
		"(x eq 1)",
		"((x eq 1))",
		"(x eq 1) and (y gt 2)",
		"((x eq 1) and (y gt 2)) or (z lt 5)",
		"not (x eq 1 and y gt 2)",
		"not ((x eq 1) and (y gt 2))",
	}

	for _, test := range tests {
		lexer := NewLexer(test)
		tokens := lexer.Tokenize()
		parser := NewParser(tokens)

		_, err := parser.Parse()
		if err != nil {
			t.Errorf("Failed to parse nested parentheses %q: %v", test, err)
		}
	}
}

// Test parser with presence operator.
func TestParserPresenceOperator(t *testing.T) {
	lexer := NewLexer("x pr")
	tokens := lexer.Tokenize()
	parser := NewParser(tokens)

	ast, err := parser.Parse()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if ast.Type != NodeUnaryOp {
		t.Error("Expected unary operation for presence operator")
	}

	if ast.Operator != PR {
		t.Error("Expected PR operator")
	}
}
