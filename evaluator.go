package rule

import (
	"reflect"
	"strconv"
	"strings"
)

type Evaluator struct {
	context map[string]any
}

func NewEvaluator() *Evaluator {
	return &Evaluator{}
}

func (e *Evaluator) Evaluate(node *ASTNode, context map[string]any) (bool, error) {
	e.context = context
	result, err := e.evaluateNode(node)
	if err != nil {
		return false, err
	}
	return e.toBool(result), nil
}

func (e *Evaluator) evaluateNode(node *ASTNode) (any, error) {
	switch node.Type {
	case NodeLiteral:
		return e.evaluateLiteral(node)

	case NodeIdentifier:
		return e.evaluateIdentifier(node)

	case NodeProperty:
		return e.evaluateProperty(node)

	case NodeUnaryOp:
		return e.evaluateUnaryOp(node)

	case NodeBinaryOp:
		return e.evaluateBinaryOp(node)

	default:
		return nil, ErrInvalidNode
	}
}

func (e *Evaluator) evaluateLiteral(node *ASTNode) (any, error) {
	switch node.Value.Type {
	case ValueString:
		return node.Value.StrValue, nil
	case ValueNumber:
		return node.Value.NumValue, nil
	case ValueBoolean:
		return node.Value.BoolValue, nil
	case ValueArray:
		return node.Value.ArrValue, nil
	default:
		return nil, ErrInvalidLiteral
	}
}

func (e *Evaluator) evaluateIdentifier(node *ASTNode) (any, error) {
	value, exists := e.context[node.Value.StrValue]
	if !exists {
		return nil, ErrAttributeNotFound
	}
	return value, nil
}

func (e *Evaluator) evaluateProperty(node *ASTNode) (any, error) {
	current := e.context

	for i, child := range node.Children {
		key := child.Value.StrValue

		if i == len(node.Children)-1 {
			if value, exists := current[key]; exists {
				return value, nil
			}
			return nil, ErrAttributeNotFound
		}

		next, exists := current[key]
		if !exists {
			return nil, ErrAttributeNotFound
		}

		if nextMap, ok := next.(map[string]any); ok {
			current = nextMap
		} else {
			return nil, ErrInvalidNestedAttribute
		}
	}

	return nil, ErrAttributeNotFound
}

func (e *Evaluator) evaluateUnaryOp(node *ASTNode) (any, error) {
	switch node.Operator {
	case NOT:
		operand, err := e.evaluateNode(node.Left)
		if err != nil {
			return false, err
		}
		return !e.toBool(operand), nil

	case PR:
		_, err := e.evaluateNode(node.Left)
		return err == nil, nil

	default:
		return nil, ErrInvalidOperator
	}
}

func (e *Evaluator) evaluateBinaryOp(node *ASTNode) (any, error) {
	left, err := e.evaluateNode(node.Left)
	if err != nil {
		// For missing attributes, treat as nil for comparison
		if err == ErrAttributeNotFound || err == ErrInvalidNestedAttribute {
			left = nil
		} else {
			return false, err
		}
	}

	right, err := e.evaluateNode(node.Right)
	if err != nil {
		// For missing attributes, treat as nil for comparison
		if err == ErrAttributeNotFound || err == ErrInvalidNestedAttribute {
			right = nil
		} else {
			return false, err
		}
	}

	switch node.Operator {
	case AND:
		return e.toBool(left) && e.toBool(right), nil
	case OR:
		return e.toBool(left) || e.toBool(right), nil
	case EQ, EQUALS:
		return e.compareEqual(left, right), nil
	case NE, NOT_EQUALS:
		return !e.compareEqual(left, right), nil
	case LT:
		return e.compareNumbers(left, right, func(a, b float64) bool { return a < b }), nil
	case GT:
		return e.compareNumbers(left, right, func(a, b float64) bool { return a > b }), nil
	case LE:
		return e.compareNumbers(left, right, func(a, b float64) bool { return a <= b }), nil
	case GE:
		return e.compareNumbers(left, right, func(a, b float64) bool { return a >= b }), nil
	case CO:
		return e.stringContains(left, right), nil
	case SW:
		return e.stringStartsWith(left, right), nil
	case EW:
		return e.stringEndsWith(left, right), nil
	case IN:
		return e.membershipCheck(left, right), nil
	default:
		return nil, ErrInvalidOperator
	}
}

func (e *Evaluator) compareEqual(left, right any) bool {
	if left == nil || right == nil {
		return false
	}

	// Check for exact type match first
	if reflect.TypeOf(left) == reflect.TypeOf(right) {
		return left == right
	}

	// Special handling for large integers to avoid float64 precision loss
	if e.isLargeInteger(left) || e.isLargeInteger(right) {
		return e.compareLargeNumbers(left, right, func(a, b float64) bool { return a == b })
	}

	// Allow numeric cross-comparison only
	leftNum, leftIsNum := e.toNumber(left)
	rightNum, rightIsNum := e.toNumber(right)

	if leftIsNum && rightIsNum {
		return leftNum == rightNum
	}

	// For different types, no match
	return false
}

func (e *Evaluator) compareNumbers(left, right any, op func(float64, float64) bool) bool {
	// Special handling for large integers to avoid float64 precision loss
	if e.isLargeInteger(left) || e.isLargeInteger(right) {
		return e.compareLargeNumbers(left, right, op)
	}

	leftNum, leftOk := e.toNumber(left)
	rightNum, rightOk := e.toNumber(right)

	if leftOk && rightOk {
		return op(leftNum, rightNum)
	}

	// If not both numbers, try string comparison
	leftStr := e.toString(left)
	rightStr := e.toString(right)

	// For string comparison, use lexicographic ordering
	switch {
	case op(1, 0): // GT operation
		return leftStr > rightStr
	case op(0, 1): // LT operation
		return leftStr < rightStr
	case op(1, 1): // GE operation
		return leftStr >= rightStr
	case op(0, 0): // LE operation
		return leftStr <= rightStr
	default:
		return false
	}
}

func (e *Evaluator) stringContains(left, right any) bool {
	leftStr := e.toString(left)
	rightStr := e.toString(right)
	return strings.Contains(leftStr, rightStr)
}

func (e *Evaluator) stringStartsWith(left, right any) bool {
	leftStr := e.toString(left)
	rightStr := e.toString(right)
	return strings.HasPrefix(leftStr, rightStr)
}

func (e *Evaluator) stringEndsWith(left, right any) bool {
	leftStr := e.toString(left)
	rightStr := e.toString(right)
	return strings.HasSuffix(leftStr, rightStr)
}

func (e *Evaluator) membershipCheck(left, right any) bool {
	if arr, ok := right.([]Value); ok {
		for _, item := range arr {
			var itemValue any
			switch item.Type {
			case ValueString:
				itemValue = item.StrValue
			case ValueNumber:
				itemValue = item.NumValue
			case ValueBoolean:
				itemValue = item.BoolValue
			}
			if e.compareEqualStrict(left, itemValue) {
				return true
			}
		}
		return false
	}

	if arr, ok := right.([]any); ok {
		for _, item := range arr {
			if e.compareEqualStrict(left, item) {
				return true
			}
		}
		return false
	}

	return false
}

func (e *Evaluator) compareEqualStrict(left, right any) bool {
	if left == nil && right == nil {
		return true
	}
	if left == nil || right == nil {
		return false
	}
	
	// Check if both are numeric types
	leftIsInt := false
	leftIsFloat := false
	rightIsInt := false
	rightIsFloat := false
	
	switch left.(type) {
	case int, int32, int64:
		leftIsInt = true
	case float64:
		leftIsFloat = true
	}
	
	switch right.(type) {
	case int, int32, int64:
		rightIsInt = true
	case float64:
		rightIsFloat = true
	}
	
	// Allow int/float cross-comparison only
	if (leftIsInt || leftIsFloat) && (rightIsInt || rightIsFloat) {
		leftNum, _ := e.toNumber(left)
		rightNum, _ := e.toNumber(right)
		return leftNum == rightNum
	}
	
	// For all other types, require exact match
	return left == right
}

func (e *Evaluator) toBool(value any) bool {
	if value == nil {
		return false
	}

	switch v := value.(type) {
	case bool:
		return v
	case string:
		return v != ""
	case float64:
		return v != 0
	case int:
		return v != 0
	case int64:
		return v != 0
	case int32:
		return v != 0
	default:
		return true
	}
}

func (e *Evaluator) toNumber(value any) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case int32:
		return float64(v), true
	case string:
		if num, err := strconv.ParseFloat(v, 64); err == nil {
			return num, true
		}
		return 0, false
	default:
		return 0, false
	}
}

func (e *Evaluator) toString(value any) string {
	if value == nil {
		return ""
	}

	switch v := value.(type) {
	case string:
		return v
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case bool:
		return strconv.FormatBool(v)
	default:
		return ""
	}
}

func (e *Evaluator) isLargeInteger(value any) bool {
	switch v := value.(type) {
	case int64:
		// Check if int64 would lose precision when converted to float64
		return v > 9007199254740992 || v < -9007199254740992 // 2^53
	case int:
		return int64(v) > 9007199254740992 || int64(v) < -9007199254740992
	case float64:
		// Check if it's a large integer-like float
		return v > 9007199254740992 || v < -9007199254740992
	default:
		return false
	}
}

func (e *Evaluator) compareLargeNumbers(left, right any, op func(float64, float64) bool) bool {
	// Convert both to string representations for exact comparison
	leftStr := e.formatLargeNumber(left)
	rightStr := e.formatLargeNumber(right)
	
	// Parse as big numbers for comparison
	leftBig, leftOk := e.parseBigNumber(leftStr)
	rightBig, rightOk := e.parseBigNumber(rightStr)
	
	if !leftOk || !rightOk {
		return false
	}
	
	// Compare as strings with proper numeric sorting
	return e.compareNumberStrings(leftBig, rightBig, op)
}

func (e *Evaluator) formatLargeNumber(value any) string {
	switch v := value.(type) {
	case int64:
		return strconv.FormatInt(v, 10)
	case int:
		return strconv.Itoa(v)
	case float64:
		// Format without scientific notation
		return strconv.FormatFloat(v, 'f', 0, 64)
	case string:
		return v
	default:
		return ""
	}
}

func (e *Evaluator) parseBigNumber(s string) (string, bool) {
	// Simple validation - just check if it's a valid number string
	if s == "" {
		return "", false
	}
	
	// Handle negative numbers
	if s[0] == '-' {
		if len(s) == 1 {
			return "", false
		}
		rest, ok := e.parseBigNumber(s[1:])
		if !ok {
			return "", false
		}
		return "-" + rest, true
	}
	
	// Check all characters are digits
	for _, r := range s {
		if r < '0' || r > '9' {
			return "", false
		}
	}
	
	return s, true
}

func (e *Evaluator) compareNumberStrings(left, right string, op func(float64, float64) bool) bool {
	// Handle sign comparison first
	leftNeg := left[0] == '-'
	rightNeg := right[0] == '-'
	
	if leftNeg && !rightNeg {
		return op(-1, 1) // left is negative, right is positive
	}
	if !leftNeg && rightNeg {
		return op(1, -1) // left is positive, right is negative
	}
	
	// Both same sign, remove negative signs for comparison
	leftAbs := left
	rightAbs := right
	if leftNeg {
		leftAbs = left[1:]
		rightAbs = right[1:]
	}
	
	// Compare lengths first
	if len(leftAbs) != len(rightAbs) {
		if leftNeg {
			// For negative numbers, longer means more negative (smaller)
			if len(leftAbs) > len(rightAbs) {
				return op(-1, 1)
			} else {
				return op(1, -1)
			}
		} else {
			// For positive numbers, longer means bigger
			if len(leftAbs) > len(rightAbs) {
				return op(1, -1)
			} else {
				return op(-1, 1)
			}
		}
	}
	
	// Same length, compare lexicographically
	cmp := 0
	if leftAbs < rightAbs {
		cmp = -1
	} else if leftAbs > rightAbs {
		cmp = 1
	}
	
	// For negative numbers, reverse the comparison
	if leftNeg {
		cmp = -cmp
	}
	
	return op(float64(cmp), 0)
}
