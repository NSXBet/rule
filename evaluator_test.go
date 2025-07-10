package rule

import (
	"testing"
)

// Test evaluator with basic operations.
func TestEvaluatorBasicOperations(t *testing.T) {
	evaluator := NewEvaluator()

	// Test EQ operation
	ast := NewBinaryOpNode(EQ,
		NewIdentifierNode("x"),
		NewNumberLiteralNode(10))

	ctx := map[string]any{"x": 10}

	result, err := evaluator.Evaluate(ast, ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !result {
		t.Error("Expected true for x eq 10 when x=10")
	}

	// Test NE operation
	ast = NewBinaryOpNode(NE,
		NewIdentifierNode("x"),
		NewNumberLiteralNode(10))

	result, err = evaluator.Evaluate(ast, ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result {
		t.Error("Expected false for x ne 10 when x=10")
	}
}

// Test evaluator with string operations.
func TestEvaluatorStringOperations(t *testing.T) {
	evaluator := NewEvaluator()

	// Test CO (contains) operation
	ast := NewBinaryOpNode(CO,
		NewIdentifierNode("name"),
		NewStringLiteralNode("John"))

	ctx := map[string]any{"name": "John Doe"}

	result, err := evaluator.Evaluate(ast, ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !result {
		t.Error("Expected true for name co 'John' when name='John Doe'")
	}

	// Test SW (starts with) operation
	ast = NewBinaryOpNode(SW,
		NewIdentifierNode("name"),
		NewStringLiteralNode("John"))

	result, err = evaluator.Evaluate(ast, ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !result {
		t.Error("Expected true for name sw 'John' when name='John Doe'")
	}

	// Test EW (ends with) operation
	ast = NewBinaryOpNode(EW,
		NewIdentifierNode("name"),
		NewStringLiteralNode("Doe"))

	result, err = evaluator.Evaluate(ast, ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !result {
		t.Error("Expected true for name ew 'Doe' when name='John Doe'")
	}
}

// Test evaluator with IN operation.
func TestEvaluatorInOperation(t *testing.T) {
	evaluator := NewEvaluator()

	ast := NewBinaryOpNode(IN,
		NewIdentifierNode("color"),
		NewArrayLiteralNode([]Value{
			{Type: ValueString, StrValue: "red"},
			{Type: ValueString, StrValue: "green"},
			{Type: ValueString, StrValue: "blue"},
		}))

	ctx := map[string]any{"color": "red"}

	result, err := evaluator.Evaluate(ast, ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !result {
		t.Error("Expected true for color in ['red', 'green', 'blue'] when color='red'")
	}

	// Test with value not in array
	ctx["color"] = "yellow"

	result, err = evaluator.Evaluate(ast, ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result {
		t.Error("Expected false for color in ['red', 'green', 'blue'] when color='yellow'")
	}
}

// Test evaluator with PR (presence) operation.
func TestEvaluatorPresenceOperation(t *testing.T) {
	evaluator := NewEvaluator()

	ast := NewUnaryOpNode(PR, NewIdentifierNode("name"))

	// Test with existing attribute
	ctx := map[string]any{"name": "John"}

	result, err := evaluator.Evaluate(ast, ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !result {
		t.Error("Expected true for name pr when name exists")
	}

	// Test with missing attribute
	ctx = map[string]any{}

	result, err = evaluator.Evaluate(ast, ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result {
		t.Error("Expected false for name pr when name is missing")
	}
}

// Test evaluator with logical operations.
func TestEvaluatorLogicalOperations(t *testing.T) {
	evaluator := NewEvaluator()

	// Test AND operation
	ast := NewBinaryOpNode(AND,
		NewBinaryOpNode(EQ,
			NewIdentifierNode("x"),
			NewNumberLiteralNode(10)),
		NewBinaryOpNode(GT,
			NewIdentifierNode("y"),
			NewNumberLiteralNode(5)))

	ctx := map[string]any{"x": 10, "y": 8}

	result, err := evaluator.Evaluate(ast, ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !result {
		t.Error("Expected true for (x eq 10 and y gt 5) when x=10, y=8")
	}

	// Test OR operation
	ast = NewBinaryOpNode(OR,
		NewBinaryOpNode(EQ,
			NewIdentifierNode("x"),
			NewNumberLiteralNode(10)),
		NewBinaryOpNode(GT,
			NewIdentifierNode("y"),
			NewNumberLiteralNode(5)))
	ctx["y"] = 3

	result, err = evaluator.Evaluate(ast, ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !result {
		t.Error("Expected true for (x eq 10 or y gt 5) when x=10, y=3")
	}
}

// Test evaluator with NOT operation.
func TestEvaluatorNotOperation(t *testing.T) {
	evaluator := NewEvaluator()

	ast := NewUnaryOpNode(NOT,
		NewBinaryOpNode(EQ,
			NewIdentifierNode("x"),
			NewNumberLiteralNode(10)))

	ctx := map[string]any{"x": 5}

	result, err := evaluator.Evaluate(ast, ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !result {
		t.Error("Expected true for not (x eq 10) when x=5")
	}
}

// Test evaluator with nested attributes.
func TestEvaluatorNestedAttributes(t *testing.T) {
	evaluator := NewEvaluator()

	ast := NewBinaryOpNode(EQ,
		NewPropertyNode([]string{"user", "name"}),
		NewStringLiteralNode("John"))

	ctx := map[string]any{
		"user": map[string]any{
			"name": "John",
			"age":  30,
		},
	}

	result, err := evaluator.Evaluate(ast, ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !result {
		t.Error("Expected true for user.name eq 'John' when user.name='John'")
	}
}

// Test evaluator with large numbers.
func TestEvaluatorLargeNumbers(t *testing.T) {
	evaluator := NewEvaluator()

	// Test basic large number comparison
	ast := NewBinaryOpNode(EQ,
		NewIdentifierNode("big"),
		NewStringLiteralNode("9223372036854775807"))

	ctx := map[string]any{"big": "9223372036854775807"}

	result, err := evaluator.Evaluate(ast, ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !result {
		t.Errorf("Expected true for large number equality, got %v", result)
	}

	// Test large number inequality
	ast = NewBinaryOpNode(NE,
		NewIdentifierNode("big"),
		NewStringLiteralNode("9223372036854775806"))

	ctx = map[string]any{"big": "9223372036854775807"}

	result, err = evaluator.Evaluate(ast, ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !result {
		t.Errorf(
			"Expected true for large number inequality (9223372036854775807 != 9223372036854775806), got %v",
			result,
		)
	}
}

// Test evaluator with boolean literals.
func TestEvaluatorBooleanLiterals(t *testing.T) {
	evaluator := NewEvaluator()

	ast := NewBinaryOpNode(EQ,
		NewIdentifierNode("active"),
		NewBooleanLiteralNode(true))

	ctx := map[string]any{"active": true}

	result, err := evaluator.Evaluate(ast, ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !result {
		t.Error("Expected true for active eq true when active=true")
	}
}

// Test evaluator with relational operations.
func TestEvaluatorRelationalOperations(t *testing.T) {
	evaluator := NewEvaluator()

	tests := []struct {
		operator TokenType
		left     float64
		right    float64
		expected bool
	}{
		{LT, 5, 10, true},
		{LT, 10, 5, false},
		{GT, 10, 5, true},
		{GT, 5, 10, false},
		{LE, 5, 10, true},
		{LE, 10, 10, true},
		{LE, 10, 5, false},
		{GE, 10, 5, true},
		{GE, 10, 10, true},
		{GE, 5, 10, false},
	}

	for _, test := range tests {
		ast := NewBinaryOpNode(test.operator,
			NewIdentifierNode("x"),
			NewNumberLiteralNode(test.right))

		ctx := map[string]any{"x": test.left}

		result, err := evaluator.Evaluate(ast, ctx)
		if err != nil {
			t.Errorf("Expected no error for %v operator, got %v", test.operator, err)
		}

		if result != test.expected {
			t.Errorf("Expected %v for %v %v %v, got %v", test.expected, test.left, test.operator, test.right, result)
		}
	}
}

// Test evaluator with type conversions.
func TestEvaluatorTypeConversions(t *testing.T) {
	evaluator := NewEvaluator()

	// Test string to number comparison
	ast := NewBinaryOpNode(EQ,
		NewIdentifierNode("x"),
		NewNumberLiteralNode(10))

	ctx := map[string]any{"x": "10"}

	result, err := evaluator.Evaluate(ast, ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result {
		t.Error("Expected false for string '10' eq number 10 (different categories should never compare equally)")
	}
}

// Test evaluator with EQUALS and NOT_EQUALS aliases.
func TestEvaluatorEqualsAliases(t *testing.T) {
	evaluator := NewEvaluator()

	// Test EQUALS (==)
	ast := NewBinaryOpNode(EQUALS,
		NewIdentifierNode("x"),
		NewNumberLiteralNode(10))

	ctx := map[string]any{"x": 10}

	result, err := evaluator.Evaluate(ast, ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !result {
		t.Error("Expected true for x == 10 when x=10")
	}

	// Test NOT_EQUALS (!=)
	ast = NewBinaryOpNode(NOT_EQUALS,
		NewIdentifierNode("x"),
		NewNumberLiteralNode(10))

	result, err = evaluator.Evaluate(ast, ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result {
		t.Error("Expected false for x != 10 when x=10")
	}
}

// Test evaluator error conditions.
func TestEvaluatorErrorConditions(t *testing.T) {
	evaluator := NewEvaluator()

	// Test missing attribute - this should NOT error, just return false
	ast := NewBinaryOpNode(EQ,
		NewIdentifierNode("missing"),
		NewNumberLiteralNode(10))

	result, err := evaluator.Evaluate(ast, map[string]any{})
	if err != nil {
		t.Errorf("Expected no error for missing attribute, got %v", err)
	}

	if result {
		t.Error("Expected false for missing attribute comparison")
	}

	// Test invalid nested attribute - this might error
	ast = NewBinaryOpNode(EQ,
		NewPropertyNode([]string{"missing", "nested"}),
		NewStringLiteralNode("test"))

	result, err = evaluator.Evaluate(ast, map[string]any{})
	if err != nil {
		// This is expected for invalid nested attributes
		t.Logf("Got expected error for invalid nested attribute: %v", err)
	} else if result {
		// Or might return false
		t.Error("Expected false for invalid nested attribute")
	}

	// Test invalid property access on non-map
	ast = NewBinaryOpNode(EQ,
		NewPropertyNode([]string{"x", "y"}),
		NewStringLiteralNode("test"))

	result, err = evaluator.Evaluate(ast, map[string]any{"x": "not a map"})
	if err != nil {
		// This is expected for property access on non-map
		t.Logf("Got expected error for property access on non-map: %v", err)
	} else if result {
		// Or might return false
		t.Error("Expected false for property access on non-map")
	}
}

// Test evaluator with complex nested expressions.
func TestEvaluatorComplexNested(t *testing.T) {
	evaluator := NewEvaluator()

	// Test (x eq 10 and y gt 5) or (z co "test")
	ast := NewBinaryOpNode(OR,
		NewBinaryOpNode(AND,
			NewBinaryOpNode(EQ,
				NewIdentifierNode("x"),
				NewNumberLiteralNode(10)),
			NewBinaryOpNode(GT,
				NewIdentifierNode("y"),
				NewNumberLiteralNode(5))),
		NewBinaryOpNode(CO,
			NewIdentifierNode("z"),
			NewStringLiteralNode("test")))

	ctx := map[string]any{"x": 5, "y": 3, "z": "testing"}

	result, err := evaluator.Evaluate(ast, ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !result {
		t.Error("Expected true for complex nested expression")
	}
}
