package rule

import (
	"testing"
)

// Test lexer peek function
func TestLexerPeek(t *testing.T) {
	lexer := NewLexer("ab")
	
	// Current should be 'a'
	if lexer.current != 'a' {
		t.Error("Expected current to be 'a'")
	}
	
	// Peek should be 'b'
	if lexer.peekChar() != 'b' {
		t.Error("Expected peek to be 'b'")
	}
	
	// Advance to 'b'
	lexer.readChar()
	if lexer.current != 'b' {
		t.Error("Expected current to be 'b'")
	}
	
	// Peek should be 0 (end of input)
	if lexer.peekChar() != 0 {
		t.Error("Expected peek to be 0")
	}
}

// Test lexer with various input types
func TestLexerTokenization(t *testing.T) {
	tests := []struct {
		input    string
		expected []TokenType
	}{
		{"x", []TokenType{IDENTIFIER, EOF}},
		{"123", []TokenType{NUMBER, EOF}},
		{"true", []TokenType{BOOLEAN, EOF}},
		{"false", []TokenType{BOOLEAN, EOF}},
		{`"hello"`, []TokenType{STRING, EOF}},
		{"[1,2,3]", []TokenType{ARRAY_START, NUMBER, COMMA, NUMBER, COMMA, NUMBER, ARRAY_END, EOF}},
		{"()", []TokenType{PAREN_OPEN, PAREN_CLOSE, EOF}},
		{"x.y", []TokenType{IDENTIFIER, DOT, IDENTIFIER, EOF}},
		{"x eq y", []TokenType{IDENTIFIER, EQ, IDENTIFIER, EOF}},
		{"x == y", []TokenType{IDENTIFIER, EQUALS, IDENTIFIER, EOF}},
		{"x != y", []TokenType{IDENTIFIER, NOT_EQUALS, IDENTIFIER, EOF}},
		{"x and y", []TokenType{IDENTIFIER, AND, IDENTIFIER, EOF}},
		{"x or y", []TokenType{IDENTIFIER, OR, IDENTIFIER, EOF}},
		{"not x", []TokenType{NOT, IDENTIFIER, EOF}},
	}
	
	for _, test := range tests {
		lexer := NewLexer(test.input)
		tokens := lexer.Tokenize()
		
		if len(tokens) != len(test.expected) {
			t.Errorf("Input %q: expected %d tokens, got %d", test.input, len(test.expected), len(tokens))
			continue
		}
		
		for i, expectedType := range test.expected {
			if tokens[i].Type != expectedType {
				t.Errorf("Input %q: token %d expected %v, got %v", test.input, i, expectedType, tokens[i].Type)
			}
		}
	}
}

// Test lexer with whitespace handling
func TestLexerWhitespace(t *testing.T) {
	lexer := NewLexer("  x    eq   y  ")
	tokens := lexer.Tokenize()
	
	expected := []TokenType{IDENTIFIER, EQ, IDENTIFIER, EOF}
	
	if len(tokens) != len(expected) {
		t.Errorf("Expected %d tokens, got %d", len(expected), len(tokens))
	}
	
	for i, expectedType := range expected {
		if tokens[i].Type != expectedType {
			t.Errorf("Token %d: expected %v, got %v", i, expectedType, tokens[i].Type)
		}
	}
}

// Test lexer with negative numbers
func TestLexerNegativeNumbers(t *testing.T) {
	lexer := NewLexer("x eq -123")
	tokens := lexer.Tokenize()
	
	expected := []TokenType{IDENTIFIER, EQ, NUMBER, EOF}
	
	if len(tokens) != len(expected) {
		t.Errorf("Expected %d tokens, got %d", len(expected), len(tokens))
	}
	
	// Check that the negative number is properly parsed
	if tokens[2].Type != NUMBER {
		t.Error("Expected third token to be NUMBER")
	}
	if tokens[2].NumValue != -123 {
		t.Errorf("Expected -123, got %f", tokens[2].NumValue)
	}
}

// Test lexer with large integers
func TestLexerLargeIntegers(t *testing.T) {
	lexer := NewLexer("x eq 9223372036854775807")
	tokens := lexer.Tokenize()
	
	expected := []TokenType{IDENTIFIER, EQ, STRING, EOF} // Large integers become strings
	
	if len(tokens) != len(expected) {
		t.Errorf("Expected %d tokens, got %d", len(expected), len(tokens))
	}
	
	// Check that the large integer is stored as string
	if tokens[2].Type != STRING {
		t.Error("Expected large integer to be stored as STRING")
	}
	if tokens[2].Value != "9223372036854775807" {
		t.Errorf("Expected '9223372036854775807', got '%s'", tokens[2].Value)
	}
}

// Test lexer with floating point numbers
func TestLexerFloatingPoint(t *testing.T) {
	lexer := NewLexer("x eq 3.14159")
	tokens := lexer.Tokenize()
	
	if len(tokens) != 4 { // IDENTIFIER, EQ, NUMBER, EOF
		t.Errorf("Expected 4 tokens, got %d", len(tokens))
	}
	
	if tokens[2].Type != NUMBER {
		t.Error("Expected third token to be NUMBER")
	}
	if tokens[2].NumValue != 3.14159 {
		t.Errorf("Expected 3.14159, got %f", tokens[2].NumValue)
	}
}

// Test lexer with string literals
func TestLexerStringLiterals(t *testing.T) {
	lexer := NewLexer(`name eq "John Doe"`)
	tokens := lexer.Tokenize()
	
	if len(tokens) != 4 { // IDENTIFIER, EQ, STRING, EOF
		t.Errorf("Expected 4 tokens, got %d", len(tokens))
	}
	
	if tokens[2].Type != STRING {
		t.Error("Expected third token to be STRING")
	}
	if tokens[2].Value != "John Doe" {
		t.Errorf("Expected 'John Doe', got '%s'", tokens[2].Value)
	}
}

// Test lexer with complex expressions
func TestLexerComplexExpression(t *testing.T) {
	lexer := NewLexer(`(x eq 10 and y gt 5) or (status in ["active", "pending"])`)
	tokens := lexer.Tokenize()
	
	// Should have proper tokenization without errors
	if len(tokens) < 10 {
		t.Error("Expected complex expression to produce many tokens")
	}
	
	// Check that the last token is EOF
	if tokens[len(tokens)-1].Type != EOF {
		t.Error("Expected last token to be EOF")
	}
}

// Test lexer with unsupported characters
func TestLexerUnsupportedCharacters(t *testing.T) {
	// Test with character that gets skipped
	lexer := NewLexer("x @ y")
	tokens := lexer.Tokenize()
	
	// Should have x, y, and EOF (@ gets skipped)
	if len(tokens) != 3 {
		t.Errorf("Expected 3 tokens, got %d", len(tokens))
	}
	
	expected := []TokenType{IDENTIFIER, IDENTIFIER, EOF}
	for i, expectedType := range expected {
		if tokens[i].Type != expectedType {
			t.Errorf("Token %d: expected %v, got %v", i, expectedType, tokens[i].Type)
		}
	}
}

// Test lexer isLargeInteger function
func TestLexerIsLargeInteger(t *testing.T) {
	lexer := NewLexer("")
	
	// Test with large integer
	if !lexer.isLargeInteger("9223372036854775807") {
		t.Error("Expected large integer to be detected")
	}
	
	// Test with small integer
	if lexer.isLargeInteger("123") {
		t.Error("Expected small integer to not be detected as large")
	}
	
	// Test with float
	if lexer.isLargeInteger("123.45") {
		t.Error("Expected float to not be detected as large integer")
	}
	
	// Test with invalid number
	if lexer.isLargeInteger("abc") {
		t.Error("Expected invalid number to not be detected as large integer")
	}
}

// Test lexer readString function
func TestLexerReadString(t *testing.T) {
	lexer := NewLexer(`"hello world"`)
	// Start at opening quote, readString() should handle consuming it
	
	result := lexer.readString()
	if result != "hello world" {
		t.Errorf("Expected 'hello world', got '%s'", result)
	}
}

// Test lexer readIdentifier function  
func TestLexerReadIdentifier(t *testing.T) {
	lexer := NewLexer("myVariable123")
	
	result := lexer.readIdentifier()
	if result != "myVariable123" {
		t.Errorf("Expected 'myVariable123', got '%s'", result)
	}
}

// Test lexer readNumber function
func TestLexerReadNumber(t *testing.T) {
	lexer := NewLexer("123.45")
	
	str, num, isLarge := lexer.readNumber()
	if str != "123.45" {
		t.Errorf("Expected string '123.45', got '%s'", str)
	}
	if num != 123.45 {
		t.Errorf("Expected number 123.45, got %f", num)
	}
	if isLarge {
		t.Error("Expected 123.45 to not be large integer")
	}
}