package rule

import (
	"testing"
)

// Test AST node helper functions.
func TestASTNodeHelpers(t *testing.T) {
	// Test IsOperator
	binaryNode := NewBinaryOpNode(AND, nil, nil)
	if !binaryNode.IsOperator() {
		t.Error("Binary node should be an operator")
	}

	unaryNode := NewUnaryOpNode(NOT, nil)
	if !unaryNode.IsOperator() {
		t.Error("Unary node should be an operator")
	}

	literalNode := NewStringLiteralNode("test")
	if literalNode.IsOperator() {
		t.Error("Literal node should not be an operator")
	}

	// Test IsLiteral
	if !literalNode.IsLiteral() {
		t.Error("Literal node should be a literal")
	}

	if binaryNode.IsLiteral() {
		t.Error("Binary node should not be a literal")
	}

	// Test IsIdentifier
	identifierNode := NewIdentifierNode("test")
	if !identifierNode.IsIdentifier() {
		t.Error("Identifier node should be an identifier")
	}

	propertyNode := NewPropertyNode([]string{"a", "b"})
	if !propertyNode.IsIdentifier() {
		t.Error("Property node should be an identifier")
	}

	if literalNode.IsIdentifier() {
		t.Error("Literal node should not be an identifier")
	}
}

// Test AST node constructors.
//
//nolint:gocognit,gocyclo,cyclop // Test function complexity is acceptable
func TestASTNodeConstructors(t *testing.T) {
	// Test NewBinaryOpNode
	left := NewNumberLiteralNode(1)
	right := NewNumberLiteralNode(2)
	binNode := NewBinaryOpNode(EQ, left, right)

	if binNode.Type != NodeBinaryOp {
		t.Error("Binary node should have NodeBinaryOp type")
	}

	if binNode.Operator != EQ {
		t.Error("Binary node should have EQ operator")
	}

	if binNode.Left != left {
		t.Error("Binary node left child incorrect")
	}

	if binNode.Right != right {
		t.Error("Binary node right child incorrect")
	}

	// Test NewUnaryOpNode
	operand := NewBooleanLiteralNode(true)
	unaryNode := NewUnaryOpNode(NOT, operand)

	if unaryNode.Type != NodeUnaryOp {
		t.Error("Unary node should have NodeUnaryOp type")
	}

	if unaryNode.Operator != NOT {
		t.Error("Unary node should have NOT operator")
	}

	if unaryNode.Left != operand {
		t.Error("Unary node operand incorrect")
	}

	// Test NewIdentifierNode
	idNode := NewIdentifierNode("myVar")
	if idNode.Type != NodeIdentifier {
		t.Error("Identifier node should have NodeIdentifier type")
	}

	if idNode.Value.Type != ValueIdentifier {
		t.Error("Identifier node value should have ValueIdentifier type")
	}

	if idNode.Value.StrValue != "myVar" {
		t.Error("Identifier node value incorrect")
	}

	// Test NewPropertyNode
	propNode := NewPropertyNode([]string{"user", "profile", "name"})
	if propNode.Type != NodeProperty {
		t.Error("Property node should have NodeProperty type")
	}

	if len(propNode.Children) != 3 {
		t.Error("Property node should have 3 children")
	}

	if propNode.Children[0].Value.StrValue != "user" {
		t.Error("Property node first child incorrect")
	}

	if propNode.Children[1].Value.StrValue != "profile" {
		t.Error("Property node second child incorrect")
	}

	if propNode.Children[2].Value.StrValue != "name" {
		t.Error("Property node third child incorrect")
	}

	// Test NewStringLiteralNode
	strNode := NewStringLiteralNode("hello")
	if strNode.Type != NodeLiteral {
		t.Error("String literal node should have NodeLiteral type")
	}

	if strNode.Value.Type != ValueString {
		t.Error("String literal value should have ValueString type")
	}

	if strNode.Value.StrValue != "hello" {
		t.Error("String literal value incorrect")
	}

	// Test NewNumberLiteralNode
	numNode := NewNumberLiteralNode(42.5)
	if numNode.Type != NodeLiteral {
		t.Error("Number literal node should have NodeLiteral type")
	}

	if numNode.Value.Type != ValueNumber {
		t.Error("Number literal value should have ValueNumber type")
	}

	if numNode.Value.NumValue != 42.5 {
		t.Error("Number literal value incorrect")
	}

	// Test NewBooleanLiteralNode
	boolNode := NewBooleanLiteralNode(true)
	if boolNode.Type != NodeLiteral {
		t.Error("Boolean literal node should have NodeLiteral type")
	}

	if boolNode.Value.Type != ValueBoolean {
		t.Error("Boolean literal value should have ValueBoolean type")
	}

	if boolNode.Value.BoolValue != true {
		t.Error("Boolean literal value incorrect")
	}

	// Test NewArrayLiteralNode
	arrayElements := []Value{
		{Type: ValueString, StrValue: "a"},
		{Type: ValueNumber, NumValue: 1},
		{Type: ValueBoolean, BoolValue: true},
	}

	arrayNode := NewArrayLiteralNode(arrayElements)
	if arrayNode.Type != NodeLiteral {
		t.Error("Array literal node should have NodeLiteral type")
	}

	if arrayNode.Value.Type != ValueArray {
		t.Error("Array literal value should have ValueArray type")
	}

	if len(arrayNode.Value.ArrValue) != 3 {
		t.Error("Array literal should have 3 elements")
	}

	if arrayNode.Value.ArrValue[0].StrValue != "a" {
		t.Error("Array literal first element incorrect")
	}

	if arrayNode.Value.ArrValue[1].NumValue != 1 {
		t.Error("Array literal second element incorrect")
	}

	if arrayNode.Value.ArrValue[2].BoolValue != true {
		t.Error("Array literal third element incorrect")
	}
}

func TestValueNumberTextParsedLiterals(t *testing.T) {
	t.Parallel()

	for _, text := range []string{
		"0",
		"-0.000",
		"001.2300",
		"9007199254740993.12500",
		"1.12345678901234567890123456789",
		"1.0049999999999999999999",
		"9007199254740993",
		"-9223372036854775808",
		"9223372036854775807",
	} {
		t.Run(text, func(t *testing.T) {
			t.Parallel()

			node, err := ParseRule("value eq " + text)
			if err != nil {
				t.Fatal(err)
			}

			if got := node.Right.Value.NumberText(); got != text {
				t.Errorf("NumberText() = %q; want %q", got, text)
			}
		})
	}
}

func TestValueNumberTextParsedArray(t *testing.T) {
	t.Parallel()

	node, err := ParseRule("value in [9007199254740993.12500, -0.000, 9007199254740993, " +
		"001.2300, 1.0049999999999999999999]")
	if err != nil {
		t.Fatal(err)
	}

	for index, want := range []string{
		"9007199254740993.12500",
		"-0.000",
		"9007199254740993",
		"001.2300",
		"1.0049999999999999999999",
	} {
		if got := node.Right.Value.ArrValue[index].NumberText(); got != want {
			t.Errorf("array[%d].NumberText() = %q; want %q", index, got, want)
		}
	}
}

func TestValueNumberTextConstructedNumbers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value Value
		want  string
	}{
		{"decimal", NewNumberLiteralNode(1.125).Value, "1.125"},
		{"negative", NewNumberLiteralNode(-1.125).Value, "-1.125"},
		{"zero", NewNumberLiteralNode(0).Value, "0"},
		{"float integer", NewNumberLiteralNode(42).Value, "42"},
		{"large integer", NewLargeIntegerLiteralNode(9007199254740993).Value, "9007199254740993"},
		{"min int64", NewLargeIntegerLiteralNode(-9223372036854775808).Value, "-9223372036854775808"},
		{"max int64", NewLargeIntegerLiteralNode(9223372036854775807).Value, "9223372036854775807"},
		{"integer field wins", Value{Type: ValueNumber, NumValue: 1, IsInt: true, IntValue: 42}, "42"},
		{"integer zero wins", Value{Type: ValueNumber, NumValue: 42, IsInt: true}, "0"},
		{"float field wins", Value{Type: ValueNumber, NumValue: 1.125, IntValue: 42}, "1.125"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.value.NumberText(); got != tt.want {
				t.Errorf("NumberText() = %q; want %q", got, tt.want)
			}
		})
	}
}
