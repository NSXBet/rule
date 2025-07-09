package rule

import (
	"strconv"
	"strings"
)

// EvalResult represents a typed evaluation result to avoid interface boxing
type EvalResult struct {
	Type      ValueType
	Bool      bool
	Num       float64
	Str       string
	Arr       []Value
	IsValid   bool
	// OriginalValue stores the original any value for complex operations like IN with []any
	OriginalValue any
	// IntValue stores the original int64 value to preserve precision for large integers
	IntValue int64
	// IsInt indicates if this numeric value should be treated as an integer
	IsInt bool
}

// Evaluator is an optimized evaluator that avoids allocations during evaluation
type Evaluator struct {
	context map[string]any
	// Pre-allocated result to reuse
	result EvalResult
}

func NewEvaluator() *Evaluator {
	return &Evaluator{}
}

func (e *Evaluator) Evaluate(node *ASTNode, context map[string]any) (bool, error) {
	e.context = context
	err := e.evaluateNode(node, &e.result)
	if err != nil {
		return false, err
	}
	return e.toBool(&e.result), nil
}

func (e *Evaluator) evaluateNode(node *ASTNode, result *EvalResult) error {
	result.IsValid = false
	
	switch node.Type {
	case NodeLiteral:
		return e.evaluateLiteral(node, result)

	case NodeIdentifier:
		return e.evaluateIdentifier(node, result)

	case NodeProperty:
		return e.evaluateProperty(node, result)

	case NodeUnaryOp:
		return e.evaluateUnaryOp(node, result)

	case NodeBinaryOp:
		return e.evaluateBinaryOp(node, result)

	default:
		return ErrInvalidNode
	}
}

func (e *Evaluator) evaluateLiteral(node *ASTNode, result *EvalResult) error {
	result.Type = node.Value.Type
	result.IsValid = true
	
	switch node.Value.Type {
	case ValueString:
		result.Str = node.Value.StrValue
	case ValueNumber:
		result.Num = node.Value.NumValue
		result.IntValue = node.Value.IntValue
		result.IsInt = node.Value.IsInt
	case ValueBoolean:
		result.Bool = node.Value.BoolValue
	case ValueArray:
		result.Arr = node.Value.ArrValue
	default:
		result.IsValid = false
		return ErrInvalidLiteral
	}
	return nil
}

func (e *Evaluator) evaluateIdentifier(node *ASTNode, result *EvalResult) error {
	value, exists := e.context[node.Value.StrValue]
	if !exists {
		return ErrAttributeNotFound
	}
	
	result.IsValid = true
	e.setResultFromAny(result, value)
	return nil
}

func (e *Evaluator) evaluateProperty(node *ASTNode, result *EvalResult) error {
	current := e.context
	
	for _, child := range node.Children {
		key := child.Value.StrValue
		
		// Navigate to the next level
		currentMap, ok := current[key]
		if !ok {
			return ErrAttributeNotFound
		}
		
		// If this is the last segment, return the value
		if child == node.Children[len(node.Children)-1] {
			result.IsValid = true
			e.setResultFromAny(result, currentMap)
			return nil
		}
		
		// Otherwise, continue navigation
		if nextMap, ok := currentMap.(map[string]any); ok {
			current = nextMap
		} else {
			return ErrInvalidNestedAttribute
		}
	}
	
	return ErrInvalidNestedAttribute
}

func (e *Evaluator) evaluateUnaryOp(node *ASTNode, result *EvalResult) error {
	var operandResult EvalResult
	
	switch node.Operator {
	case NOT:
		err := e.evaluateNode(node.Left, &operandResult)
		if err != nil {
			return err
		}
		result.Type = ValueBoolean
		result.Bool = !e.toBool(&operandResult)
		result.IsValid = true
		return nil
		
	case PR:
		// Presence check - special case, don't evaluate the operand
		if node.Left.Type == NodeIdentifier {
			_, exists := e.context[node.Left.Value.StrValue]
			result.Type = ValueBoolean
			result.Bool = exists
			result.IsValid = true
			return nil
		} else if node.Left.Type == NodeProperty {
			// Check nested property existence
			current := e.context
			for _, child := range node.Left.Children {
				key := child.Value.StrValue
				if currentValue, ok := current[key]; ok {
					if child == node.Left.Children[len(node.Left.Children)-1] {
						result.Type = ValueBoolean
						result.Bool = true
						result.IsValid = true
						return nil
					}
					if nextMap, ok := currentValue.(map[string]any); ok {
						current = nextMap
					} else {
						result.Type = ValueBoolean
						result.Bool = false
						result.IsValid = true
						return nil
					}
				} else {
					result.Type = ValueBoolean
					result.Bool = false
					result.IsValid = true
					return nil
				}
			}
		}
		return ErrInvalidOperator
	
	default:
		return ErrInvalidOperator
	}
}

func (e *Evaluator) evaluateBinaryOp(node *ASTNode, result *EvalResult) error {
	var leftResult, rightResult EvalResult
	
	switch node.Operator {
	case AND:
		err := e.evaluateNode(node.Left, &leftResult)
		if err != nil {
			return err
		}
		if !e.toBool(&leftResult) {
			result.Type = ValueBoolean
			result.Bool = false
			result.IsValid = true
			return nil
		}
		err = e.evaluateNode(node.Right, &rightResult)
		if err != nil {
			return err
		}
		result.Type = ValueBoolean
		result.Bool = e.toBool(&rightResult)
		result.IsValid = true
		return nil
		
	case OR:
		err := e.evaluateNode(node.Left, &leftResult)
		if err != nil {
			return err
		}
		if e.toBool(&leftResult) {
			result.Type = ValueBoolean
			result.Bool = true
			result.IsValid = true
			return nil
		}
		err = e.evaluateNode(node.Right, &rightResult)
		if err != nil {
			return err
		}
		result.Type = ValueBoolean
		result.Bool = e.toBool(&rightResult)
		result.IsValid = true
		return nil
		
	default:
		// Comparison operations
		err := e.evaluateNode(node.Left, &leftResult)
		if err != nil {
			return err
		}
		err = e.evaluateNode(node.Right, &rightResult)
		if err != nil {
			return err
		}
		
		result.Type = ValueBoolean
		result.IsValid = true
		
		switch node.Operator {
		case EQ, EQUALS:
			result.Bool = e.compareEqual(&leftResult, &rightResult)
		case NE, NOT_EQUALS:
			result.Bool = !e.compareEqual(&leftResult, &rightResult)
		case LT:
			result.Bool = e.compareNumbers(&leftResult, &rightResult, func(a, b float64) bool { return a < b })
		case GT:
			result.Bool = e.compareNumbers(&leftResult, &rightResult, func(a, b float64) bool { return a > b })
		case LE:
			result.Bool = e.compareNumbers(&leftResult, &rightResult, func(a, b float64) bool { return a <= b })
		case GE:
			result.Bool = e.compareNumbers(&leftResult, &rightResult, func(a, b float64) bool { return a >= b })
		case CO:
			result.Bool = e.stringContains(&leftResult, &rightResult)
		case SW:
			result.Bool = e.stringStartsWith(&leftResult, &rightResult)
		case EW:
			result.Bool = e.stringEndsWith(&leftResult, &rightResult)
		case IN:
			result.Bool = e.membershipCheck(&leftResult, &rightResult)
		default:
			result.IsValid = false
			return ErrInvalidOperator
		}
		
		return nil
	}
}

func (e *Evaluator) setResultFromAny(result *EvalResult, value any) {
	result.OriginalValue = value
	
	switch v := value.(type) {
	case bool:
		result.Type = ValueBoolean
		result.Bool = v
	case int:
		result.Type = ValueNumber
		result.Num = float64(v)
		result.IntValue = int64(v)
		result.IsInt = true
	case int8:
		result.Type = ValueNumber
		result.Num = float64(v)
		result.IntValue = int64(v)
		result.IsInt = true
	case int16:
		result.Type = ValueNumber
		result.Num = float64(v)
		result.IntValue = int64(v)
		result.IsInt = true
	case int32:
		result.Type = ValueNumber
		result.Num = float64(v)
		result.IntValue = int64(v)
		result.IsInt = true
	case int64:
		result.Type = ValueNumber
		result.Num = float64(v)
		result.IntValue = v
		result.IsInt = true
	case uint:
		result.Type = ValueNumber
		result.Num = float64(v)
		result.IntValue = int64(v)
		result.IsInt = true
	case uint8:
		result.Type = ValueNumber
		result.Num = float64(v)
		result.IntValue = int64(v)
		result.IsInt = true
	case uint16:
		result.Type = ValueNumber
		result.Num = float64(v)
		result.IntValue = int64(v)
		result.IsInt = true
	case uint32:
		result.Type = ValueNumber
		result.Num = float64(v)
		result.IntValue = int64(v)
		result.IsInt = true
	case uint64:
		result.Type = ValueNumber
		result.Num = float64(v)
		result.IntValue = int64(v)
		result.IsInt = true
	case float32:
		result.Type = ValueNumber
		result.Num = float64(v)
		result.IsInt = false
	case float64:
		result.Type = ValueNumber
		result.Num = v
		result.IsInt = false
	case string:
		result.Type = ValueString
		result.Str = v
	case []any:
		// Handle []any slices by storing original value
		result.Type = ValueArray
		result.OriginalValue = v
	default:
		// Fallback to string representation
		result.Type = ValueString
		result.Str = e.toString(value)
	}
}

func (e *Evaluator) compareEqual(left, right *EvalResult) bool {
	// Handle same type comparisons
	if left.Type == right.Type {
		switch left.Type {
		case ValueBoolean:
			return left.Bool == right.Bool
		case ValueNumber:
			// If both are large integers, compare as int64 to avoid precision loss
			if left.IsInt && right.IsInt {
				return left.IntValue == right.IntValue
			}
			return left.Num == right.Num
		case ValueString:
			return left.Str == right.Str
		}
	}
	
	// Allow numeric cross-type comparisons (int/float) like the original library
	if (left.Type == ValueNumber && right.Type == ValueNumber) {
		// If both are large integers, compare as int64 to avoid precision loss
		if left.IsInt && right.IsInt {
			return left.IntValue == right.IntValue
		}
		return left.Num == right.Num
	}
	
	// No other cross-type comparisons are allowed
	return false
}

// compareEqualStrict performs strict type comparison for membership operations
func (e *Evaluator) compareEqualStrict(left, right *EvalResult) bool {
	// Strict type comparison - no cross-type conversions
	if left.Type != right.Type {
		return false
	}
	
	switch left.Type {
	case ValueBoolean:
		return left.Bool == right.Bool
	case ValueNumber:
		// If both are large integers, compare as int64 to avoid precision loss
		if left.IsInt && right.IsInt {
			return left.IntValue == right.IntValue
		}
		return left.Num == right.Num
	case ValueString:
		return left.Str == right.Str
	}
	
	return false
}

func (e *Evaluator) compareNumbers(left, right *EvalResult, op func(float64, float64) bool) bool {
	// Allow numeric comparisons
	if left.Type == ValueNumber && right.Type == ValueNumber {
		// If both are large integers, compare as int64 to avoid precision loss
		if left.IsInt && right.IsInt {
			// Use int64 comparison to avoid precision loss
			return e.compareInt64(left.IntValue, right.IntValue, op)
		}
		// Otherwise, use float64 comparison
		return op(left.Num, right.Num)
	}
	
	// Allow string comparisons (lexicographic)
	if left.Type == ValueString && right.Type == ValueString {
		cmpResult := strings.Compare(left.Str, right.Str)
		return op(float64(cmpResult), 0)
	}
	
	// No cross-category comparisons
	return false
}

// compareInt64 performs int64 comparison using the float64 comparison function
func (e *Evaluator) compareInt64(left, right int64, op func(float64, float64) bool) bool {
	if left < right {
		return op(-1, 0)
	} else if left > right {
		return op(1, 0)
	} else {
		return op(0, 0)
	}
}

func (e *Evaluator) stringContains(left, right *EvalResult) bool {
	leftStr := e.resultToString(left)
	rightStr := e.resultToString(right)
	return strings.Contains(leftStr, rightStr)
}

func (e *Evaluator) stringStartsWith(left, right *EvalResult) bool {
	leftStr := e.resultToString(left)
	rightStr := e.resultToString(right)
	return strings.HasPrefix(leftStr, rightStr)
}

func (e *Evaluator) stringEndsWith(left, right *EvalResult) bool {
	leftStr := e.resultToString(left)
	rightStr := e.resultToString(right)
	return strings.HasSuffix(leftStr, rightStr)
}

func (e *Evaluator) membershipCheck(left, right *EvalResult) bool {
	// Handle Value array (from literals)
	if right.Type == ValueArray && right.Arr != nil {
		for _, item := range right.Arr {
			var itemResult EvalResult
			e.setResultFromValue(&itemResult, &item)
			if e.compareEqualStrict(left, &itemResult) {
				return true
			}
		}
		return false
	}
	
	// Handle []any array (from context variables)
	if right.Type == ValueArray && right.OriginalValue != nil {
		if arr, ok := right.OriginalValue.([]any); ok {
			for _, item := range arr {
				var itemResult EvalResult
				e.setResultFromAny(&itemResult, item)
				if e.compareEqualStrict(left, &itemResult) {
					return true
				}
			}
			return false
		}
	}
	
	return false
}

func (e *Evaluator) setResultFromValue(result *EvalResult, value *Value) {
	result.Type = value.Type
	result.IsValid = true
	
	switch value.Type {
	case ValueString:
		result.Str = value.StrValue
	case ValueNumber:
		result.Num = value.NumValue
		result.IntValue = value.IntValue
		result.IsInt = value.IsInt
	case ValueBoolean:
		result.Bool = value.BoolValue
	case ValueArray:
		result.Arr = value.ArrValue
	}
}

func (e *Evaluator) toBool(result *EvalResult) bool {
	if !result.IsValid {
		return false
	}
	
	switch result.Type {
	case ValueBoolean:
		return result.Bool
	case ValueNumber:
		return result.Num != 0
	case ValueString:
		return result.Str != ""
	case ValueArray:
		return len(result.Arr) > 0
	default:
		return false
	}
}

func (e *Evaluator) resultToString(result *EvalResult) string {
	if !result.IsValid {
		return ""
	}
	
	switch result.Type {
	case ValueString:
		return result.Str
	case ValueNumber:
		return strconv.FormatFloat(result.Num, 'f', -1, 64)
	case ValueBoolean:
		if result.Bool {
			return "true"
		}
		return "false"
	default:
		return ""
	}
}

func (e *Evaluator) toNumberFromString(s string) (float64, bool) {
	if num, err := strconv.ParseFloat(s, 64); err == nil {
		return num, true
	}
	return 0, false
}

func (e *Evaluator) toString(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	default:
		return ""
	}
}