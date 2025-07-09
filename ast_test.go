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
