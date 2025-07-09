package rule

import (
	"testing"
)

// Test Token struct and String method.
func TestTokenString(t *testing.T) {
	token := Token{Type: IDENTIFIER, Value: "test"}

	str := token.String()
	if str == "" {
		t.Error("Token string representation should not be empty")
	}

	// Test that token has the expected type and value
	if token.Type != IDENTIFIER {
		t.Error("Token type should be IDENTIFIER")
	}

	if token.Value != "test" {
		t.Error("Token value should be 'test'")
	}
}

// Test Token with different types.
func TestTokenTypes(t *testing.T) {
	tests := []struct {
		tokenType TokenType
		value     string
	}{
		{IDENTIFIER, "variable"},
		{STRING, "hello"},
		{NUMBER, "123"},
		{BOOLEAN, "true"},
		{EQ, ""},
		{NE, ""},
		{LT, ""},
		{GT, ""},
		{LE, ""},
		{GE, ""},
		{AND, ""},
		{OR, ""},
		{NOT, ""},
		{PAREN_OPEN, ""},
		{PAREN_CLOSE, ""},
		{ARRAY_START, ""},
		{ARRAY_END, ""},
		{COMMA, ""},
		{DOT, ""},
		{EOF, ""},
	}

	for _, test := range tests {
		token := Token{Type: test.tokenType, Value: test.value}
		if token.Type != test.tokenType {
			t.Errorf("Expected token type %v, got %v", test.tokenType, token.Type)
		}

		if token.Value != test.value {
			t.Errorf("Expected token value %s, got %s", test.value, token.Value)
		}

		// Test that String() method works
		str := token.String()
		if str == "" && test.tokenType != EOF {
			t.Errorf("Token string representation should not be empty for type %v", test.tokenType)
		}
	}
}

// Test Token with numeric values.
func TestTokenNumeric(t *testing.T) {
	token := Token{
		Type:     NUMBER,
		Value:    "123.45",
		NumValue: 123.45,
	}

	if token.NumValue != 123.45 {
		t.Errorf("Expected numeric value 123.45, got %f", token.NumValue)
	}

	str := token.String()
	if str == "" {
		t.Error("Numeric token string representation should not be empty")
	}
}

// Test Token with boolean values.
func TestTokenBoolean(t *testing.T) {
	token := Token{
		Type:      BOOLEAN,
		Value:     "true",
		BoolValue: true,
	}

	if token.BoolValue != true {
		t.Errorf("Expected boolean value true, got %v", token.BoolValue)
	}

	str := token.String()
	if str == "" {
		t.Error("Boolean token string representation should not be empty")
	}
}

// Test Token position tracking.
func TestTokenPosition(t *testing.T) {
	token := Token{
		Type:  IDENTIFIER,
		Value: "test",
		Start: 5,
		End:   9,
	}

	if token.Start != 5 {
		t.Errorf("Expected start position 5, got %d", token.Start)
	}

	if token.End != 9 {
		t.Errorf("Expected end position 9, got %d", token.End)
	}
}
