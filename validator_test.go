package rule

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateAST(t *testing.T) {
	t.Run("NilAST", testNilASTValidation)
	t.Run("ValidQueries", testValidASTQueries)
	t.Run("InvalidQueries", testInvalidASTQueries)
}

func testNilASTValidation(t *testing.T) {
	err := ValidateAST(nil)
	require.NoError(t, err, "ValidateAST(nil) should return nil")
}

func testValidASTQueries(t *testing.T) {
	tests := []string{
		`field eq "value"`,
		`field in [1, 2, 3]`,
		`"hello" co "ell"`,
		`field pr`,
	}

	for _, query := range tests {
		t.Run(query, func(t *testing.T) {
			ast, err := ParseRule(query)
			require.NoError(t, err, "query=%q", query)

			validationErr := ValidateAST(ast)
			require.NoError(t, validationErr, "query=%q", query)
		})
	}
}

func testInvalidASTQueries(t *testing.T) {
	tests := []struct {
		query     string
		errorType *EngineError
	}{
		{`field in "string"`, ErrInvalidInOperand},
		{`123 co "test"`, ErrInvalidStringOp},
		{`"string" pr`, ErrInvalidPresenceOp},
	}

	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			_, parseErr := ParseRule(tt.query)
			if parseErr == nil {
				t.Errorf("Expected parsing error for query '%s', but got none", tt.query)
				return
			}

			if !strings.Contains(parseErr.Error(), tt.errorType.Message) {
				t.Errorf("Expected error containing '%s', got '%s'", tt.errorType.Message, parseErr.Error())
			}
		})
	}
}

func TestValidateInOperation(t *testing.T) {
	tests := []struct {
		name        string
		rightType   NodeType
		rightValue  Value
		expectError bool
	}{
		{
			name:        "Array node (valid)",
			rightType:   NodeArray,
			expectError: false,
		},
		{
			name:        "Identifier (valid - runtime check)",
			rightType:   NodeIdentifier,
			expectError: false,
		},
		{
			name:        "Property (valid - runtime check)",
			rightType:   NodeProperty,
			expectError: false,
		},
		{
			name:        "String literal (invalid)",
			rightType:   NodeLiteral,
			rightValue:  Value{Type: ValueString, StrValue: "test"},
			expectError: true,
		},
		{
			name:        "Number literal (invalid)",
			rightType:   NodeLiteral,
			rightValue:  Value{Type: ValueNumber, NumValue: 123},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &ASTNode{
				Type:     NodeBinaryOp,
				Operator: IN,
				Left:     &ASTNode{Type: NodeIdentifier},
				Right: &ASTNode{
					Type:  tt.rightType,
					Value: tt.rightValue,
				},
			}

			err := validateInOperation(node)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			} else if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestValidateStringOperation(t *testing.T) {
	tests := []struct {
		name        string
		leftValue   Value
		rightValue  Value
		expectError bool
	}{
		{
			name:        "Both strings (valid)",
			leftValue:   Value{Type: ValueString, StrValue: "hello"},
			rightValue:  Value{Type: ValueString, StrValue: "ell"},
			expectError: false,
		},
		{
			name:        "Left number (invalid)",
			leftValue:   Value{Type: ValueNumber, NumValue: 123},
			rightValue:  Value{Type: ValueString, StrValue: "test"},
			expectError: true,
		},
		{
			name:        "Right number (invalid)",
			leftValue:   Value{Type: ValueString, StrValue: "test"},
			rightValue:  Value{Type: ValueNumber, NumValue: 123},
			expectError: true,
		},
		{
			name:        "Both numbers (invalid)",
			leftValue:   Value{Type: ValueNumber, NumValue: 123},
			rightValue:  Value{Type: ValueNumber, NumValue: 456},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &ASTNode{
				Type:     NodeBinaryOp,
				Operator: CO, // Using contains operator
				Left: &ASTNode{
					Type:  NodeLiteral,
					Value: tt.leftValue,
				},
				Right: &ASTNode{
					Type:  NodeLiteral,
					Value: tt.rightValue,
				},
			}

			err := validateStringOperation(node)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			} else if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestValidatePresenceOperation(t *testing.T) {
	tests := []struct {
		name        string
		operandType NodeType
		expectError bool
	}{
		{
			name:        "Identifier (valid)",
			operandType: NodeIdentifier,
			expectError: false,
		},
		{
			name:        "Property (valid)",
			operandType: NodeProperty,
			expectError: false,
		},
		{
			name:        "Literal (invalid)",
			operandType: NodeLiteral,
			expectError: true,
		},
		{
			name:        "Array (invalid)",
			operandType: NodeArray,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &ASTNode{
				Type:     NodeUnaryOp,
				Operator: PR,
				Left: &ASTNode{
					Type: tt.operandType,
				},
			}

			err := validatePresenceOperation(node)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			} else if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
