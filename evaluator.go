package rule

import (
	"strconv"
	"strings"
	"time"
)

// EvalResult represents a typed evaluation result to avoid interface boxing.
type EvalResult struct {
	Type    ValueType
	Bool    bool
	Num     float64
	Str     string
	Arr     []Value
	IsValid bool
	// OriginalValue stores the original any value for complex operations like IN with []any
	OriginalValue any
	// IntValue stores the original int64 value to preserve precision for large integers
	IntValue int64
	// IsInt indicates if this numeric value should be treated as an integer
	IsInt bool
}

// Evaluator is an optimized evaluator that avoids allocations during evaluation.
type Evaluator struct {
	context D
	// Pre-allocated result to reuse
	result EvalResult
}

func NewEvaluator() *Evaluator {
	return &Evaluator{}
}

func (e *Evaluator) Evaluate(node *ASTNode, context D) (bool, error) {
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

	case NodeArray:
		return ErrInvalidNode // Arrays are not directly evaluatable

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
	case ValueIdentifier:
		result.IsValid = false
		return ErrInvalidLiteral // Identifiers should not be in literals
	default:
		result.IsValid = false
		return ErrInvalidLiteral
	}

	return nil
}

func (e *Evaluator) evaluateIdentifier(node *ASTNode, result *EvalResult) error {
	value, exists := e.context[node.Value.StrValue]
	if !exists {
		// For missing attributes, return a special "missing" result
		result.IsValid = false
		result.Type = ValueString // Default type for missing

		return nil
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
			// For missing nested attributes, return invalid result
			result.IsValid = false
			result.Type = ValueString // Default type for missing

			return nil
		}

		// If this is the last segment, return the value
		if child == node.Children[len(node.Children)-1] {
			result.IsValid = true
			e.setResultFromAny(result, currentMap)

			return nil
		}

		// Otherwise, continue navigation
		if nextMap, isMap := currentMap.(map[string]any); isMap {
			current = nextMap
		} else {
			// For invalid nested access, return invalid result
			result.IsValid = false
			result.Type = ValueString // Default type for missing

			return nil
		}
	}

	// Should not reach here, but handle gracefully
	result.IsValid = false
	result.Type = ValueString

	return nil
}

func (e *Evaluator) evaluateUnaryOp(node *ASTNode, result *EvalResult) error {
	switch node.Operator {
	case NOT:
		return e.evaluateNotOperator(node, result)
	case PR:
		return e.evaluatePresenceOperator(node, result)
	case EOF,
		IDENTIFIER,
		STRING,
		NUMBER,
		BOOLEAN,
		ARRAY_START,
		ARRAY_END,
		PAREN_OPEN,
		PAREN_CLOSE,
		DOT,
		COMMA,
		EQ,
		NE,
		LT,
		GT,
		LE,
		GE,
		CO,
		SW,
		EW,
		IN,
		NOT_IN,
		DQ,
		DN,
		BE,
		BQ,
		AF,
		AQ,
		AND,
		OR,
		EQUALS,
		NOT_EQUALS:
		return ErrInvalidOperator // These are not unary operators
	default:
		return ErrInvalidOperator
	}
}

func (e *Evaluator) evaluateBinaryOp(node *ASTNode, result *EvalResult) error {
	switch node.Operator {
	case AND:
		return e.evaluateLogicalAnd(node, result)
	case OR:
		return e.evaluateLogicalOr(node, result)
	case EQ, NE, LT, GT, LE, GE, CO, SW, EW, IN, NOT_IN, EQUALS, NOT_EQUALS, DQ, DN, BE, BQ, AF, AQ:
		return e.evaluateComparisonOperator(node, result)
	case EOF,
		IDENTIFIER,
		STRING,
		NUMBER,
		BOOLEAN,
		ARRAY_START,
		ARRAY_END,
		PAREN_OPEN,
		PAREN_CLOSE,
		DOT,
		COMMA,
		PR,
		NOT:
		result.IsValid = false
		return ErrInvalidOperator // These are not binary operators
	default:
		result.IsValid = false
		return ErrInvalidOperator
	}
}

// evaluateNotOperator handles the NOT unary operator.
func (e *Evaluator) evaluateNotOperator(node *ASTNode, result *EvalResult) error {
	var operandResult EvalResult

	err := e.evaluateNode(node.Left, &operandResult)
	if err != nil {
		return err
	}

	result.Type = ValueBoolean
	result.Bool = !e.toBool(&operandResult)
	result.IsValid = true

	return nil
}

// evaluatePresenceOperator handles the PR (presence) unary operator.
func (e *Evaluator) evaluatePresenceOperator(node *ASTNode, result *EvalResult) error {
	switch node.Left.Type {
	case NodeIdentifier:
		return e.checkIdentifierPresence(node, result)
	case NodeProperty:
		return e.checkPropertyPresence(node, result)
	case NodeBinaryOp, NodeUnaryOp, NodeLiteral, NodeArray:
		return ErrInvalidOperator // Invalid node types for PR operator
	default:
		return ErrInvalidOperator
	}
}

// checkIdentifierPresence checks if a simple identifier exists in the context.
func (e *Evaluator) checkIdentifierPresence(node *ASTNode, result *EvalResult) error {
	_, exists := e.context[node.Left.Value.StrValue]
	result.Type = ValueBoolean
	result.Bool = exists
	result.IsValid = true

	return nil
}

// checkPropertyPresence checks if a nested property exists in the context.
func (e *Evaluator) checkPropertyPresence(node *ASTNode, result *EvalResult) error {
	current := e.context

	for _, child := range node.Left.Children {
		key := child.Value.StrValue

		currentValue, exists := current[key]
		if !exists {
			e.setPresenceResult(result, false)
			return nil
		}

		if child == node.Left.Children[len(node.Left.Children)-1] {
			e.setPresenceResult(result, true)
			return nil
		}

		if nextMap, isMap := currentValue.(map[string]any); isMap {
			current = nextMap
		} else {
			e.setPresenceResult(result, false)
			return nil
		}
	}

	return ErrInvalidOperator // Should not reach here
}

// setPresenceResult sets the result for presence check operations.
func (e *Evaluator) setPresenceResult(result *EvalResult, exists bool) {
	result.Type = ValueBoolean
	result.Bool = exists
	result.IsValid = true
}

// evaluateLogicalAnd handles the AND logical operator with short-circuit evaluation.
func (e *Evaluator) evaluateLogicalAnd(node *ASTNode, result *EvalResult) error {
	var leftResult EvalResult

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

	var rightResult EvalResult

	err = e.evaluateNode(node.Right, &rightResult)
	if err != nil {
		return err
	}

	result.Type = ValueBoolean
	result.Bool = e.toBool(&rightResult)
	result.IsValid = true

	return nil
}

// evaluateLogicalOr handles the OR logical operator with short-circuit evaluation.
func (e *Evaluator) evaluateLogicalOr(node *ASTNode, result *EvalResult) error {
	var leftResult EvalResult

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

	var rightResult EvalResult

	err = e.evaluateNode(node.Right, &rightResult)
	if err != nil {
		return err
	}

	result.Type = ValueBoolean
	result.Bool = e.toBool(&rightResult)
	result.IsValid = true

	return nil
}

// evaluateComparisonOperator handles all comparison operators.
func (e *Evaluator) evaluateComparisonOperator(node *ASTNode, result *EvalResult) error {
	var leftResult, rightResult EvalResult

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

	// If either operand is invalid (missing attribute), comparison is false
	if !leftResult.IsValid || !rightResult.IsValid {
		result.Bool = false
		return nil
	}

	return e.performComparison(node.Operator, &leftResult, &rightResult, result)
}

// performComparison executes the specific comparison operation.
func (e *Evaluator) performComparison(
	operator TokenType,
	left, right *EvalResult,
	result *EvalResult,
) error {
	switch operator {
	case EQ, EQUALS:
		result.Bool = e.compareEqual(left, right)
	case NE, NOT_EQUALS:
		result.Bool = !e.compareEqual(left, right)
	case LT:
		result.Bool = e.compareNumbers(left, right, func(a, b float64) bool { return a < b })
	case GT:
		result.Bool = e.compareNumbers(left, right, func(a, b float64) bool { return a > b })
	case LE:
		result.Bool = e.compareNumbers(left, right, func(a, b float64) bool { return a <= b })
	case GE:
		result.Bool = e.compareNumbers(left, right, func(a, b float64) bool { return a >= b })
	case CO:
		result.Bool = e.stringContains(left, right)
	case SW:
		result.Bool = e.stringStartsWith(left, right)
	case EW:
		result.Bool = e.stringEndsWith(left, right)
	case IN:
		result.Bool = e.membershipCheck(left, right)
	case NOT_IN:
		result.Bool = !e.membershipCheck(left, right)
	case DQ:
		result.Bool = e.compareDateTimes(
			left,
			right,
			func(a, b time.Time) bool { return a.Equal(b) },
		)
	case DN:
		result.Bool = !e.compareDateTimes(
			left,
			right,
			func(a, b time.Time) bool { return a.Equal(b) },
		)
	case BE:
		result.Bool = e.compareDateTimes(
			left,
			right,
			func(a, b time.Time) bool { return a.Before(b) },
		)
	case BQ:
		result.Bool = e.compareDateTimes(
			left,
			right,
			func(a, b time.Time) bool { return a.Before(b) || a.Equal(b) },
		)
	case AF:
		result.Bool = e.compareDateTimes(
			left,
			right,
			func(a, b time.Time) bool { return a.After(b) },
		)
	case AQ:
		result.Bool = e.compareDateTimes(
			left,
			right,
			func(a, b time.Time) bool { return a.After(b) || a.Equal(b) },
		)
	case EOF,
		IDENTIFIER,
		STRING,
		NUMBER,
		BOOLEAN,
		ARRAY_START,
		ARRAY_END,
		PAREN_OPEN,
		PAREN_CLOSE,
		DOT,
		COMMA,
		PR,
		AND,
		OR,
		NOT:
		result.IsValid = false
		return ErrInvalidOperator
	default:
		result.IsValid = false
		return ErrInvalidOperator
	}

	return nil
}

func (e *Evaluator) setResultFromAny(result *EvalResult, value any) {
	result.OriginalValue = value

	switch v := value.(type) {
	case bool:
		result.Type = ValueBoolean
		result.Bool = v
	case int, int8, int16, int32, int64:
		e.setIntegerResult(result, v)
	case uint, uint8, uint16, uint32, uint64:
		e.setUintegerResult(result, v)
	case float32, float64:
		e.setFloatResult(result, v)
	case string:
		result.Type = ValueString
		result.Str = v
	case time.Time:
		// Store time.Time values as strings for compatibility with nikunjy/rules
		// The original time.Time is preserved in OriginalValue for datetime operators
		result.Type = ValueString
		result.Str = v.String()
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

func (e *Evaluator) setIntegerResult(result *EvalResult, value any) {
	result.Type = ValueNumber
	result.IsInt = true

	switch v := value.(type) {
	case int:
		result.Num = float64(v)
		result.IntValue = int64(v)
	case int8:
		result.Num = float64(v)
		result.IntValue = int64(v)
	case int16:
		result.Num = float64(v)
		result.IntValue = int64(v)
	case int32:
		result.Num = float64(v)
		result.IntValue = int64(v)
	case int64:
		result.Num = float64(v)
		result.IntValue = v
	}
}

func (e *Evaluator) setUintegerResult(result *EvalResult, value any) {
	result.Type = ValueNumber
	result.IsInt = true

	switch v := value.(type) {
	case uint:
		result.Num = float64(v)
		if v <= uint(maxSafeInteger) {
			result.IntValue = int64(v)
		} else {
			result.IntValue = maxSafeInteger
		}
	case uint8:
		result.Num = float64(v)
		result.IntValue = int64(v)
	case uint16:
		result.Num = float64(v)
		result.IntValue = int64(v)
	case uint32:
		result.Num = float64(v)
		result.IntValue = int64(v)
	case uint64:
		result.Num = float64(v)
		if v <= uint64(maxSafeInteger) {
			result.IntValue = int64(v)
		} else {
			result.IntValue = maxSafeInteger
		}
	}
}

func (e *Evaluator) setFloatResult(result *EvalResult, value any) {
	result.Type = ValueNumber
	result.IsInt = false

	switch v := value.(type) {
	case float32:
		result.Num = float64(v)
	case float64:
		result.Num = v
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
			return e.equalIgnoreCase(left.Str, right.Str)
		case ValueArray, ValueIdentifier:
			return false // Arrays and identifiers cannot be compared
		}
	}

	// Allow numeric cross-type comparisons (int/float) like the original library
	if left.Type == ValueNumber && right.Type == ValueNumber {
		// If both are large integers, compare as int64 to avoid precision loss
		if left.IsInt && right.IsInt {
			return left.IntValue == right.IntValue
		}

		return left.Num == right.Num
	}

	// No other cross-type comparisons are allowed
	return false
}

// compareEqualStrict performs strict type comparison for membership operations.
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
		return e.equalIgnoreCase(left.Str, right.Str)
	case ValueArray, ValueIdentifier:
		return false // Arrays and identifiers cannot be compared in strict mode
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

// compareInt64 performs int64 comparison using the float64 comparison function.
func (e *Evaluator) compareInt64(left, right int64, op func(float64, float64) bool) bool {
	if left < right {
		return op(-1, 0)
	}

	if left > right {
		return op(1, 0)
	}

	return op(0, 0)
}

func (e *Evaluator) stringContains(left, right *EvalResult) bool {
	leftStr := e.resultToString(left)
	rightStr := e.resultToString(right)

	return e.containsIgnoreCase(leftStr, rightStr)
}

func (e *Evaluator) stringStartsWith(left, right *EvalResult) bool {
	leftStr := e.resultToString(left)
	rightStr := e.resultToString(right)

	return e.hasPrefixIgnoreCase(leftStr, rightStr)
}

func (e *Evaluator) stringEndsWith(left, right *EvalResult) bool {
	leftStr := e.resultToString(left)
	rightStr := e.resultToString(right)

	return e.hasSuffixIgnoreCase(leftStr, rightStr)
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
	case ValueIdentifier:
		// Identifiers should not be converted to results
		result.IsValid = false
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
	case ValueIdentifier:
		return false // Identifiers are not boolean
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
			return trueString
		}

		return "false"
	case ValueArray:
		return "" // Arrays cannot be converted to string
	case ValueIdentifier:
		return "" // Identifiers cannot be converted to string
	default:
		return ""
	}
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

// parseDateTime attempts to parse a value as a datetime, supporting RFC 3339 and Unix timestamps.
func (e *Evaluator) parseDateTime(result *EvalResult) (time.Time, bool) {
	// Check if the original value is a time.Time (for context values)
	if t, ok := result.OriginalValue.(time.Time); ok {
		return t.UTC(), true
	}

	switch result.Type {
	case ValueString:
		// Try parsing as RFC 3339 first
		if t, err := time.Parse(time.RFC3339, result.Str); err == nil {
			return t.UTC(), true
		}
		// Try parsing as Unix timestamp string
		if unix, err := strconv.ParseInt(result.Str, 10, 64); err == nil {
			return time.Unix(unix, 0).UTC(), true
		}

		return time.Time{}, false
	case ValueNumber:
		// Check if it's a large integer first
		if result.IsInt {
			return time.Unix(result.IntValue, 0).UTC(), true
		}
		// Unix timestamp as float (truncate to seconds)
		return time.Unix(int64(result.Num), 0).UTC(), true
	case ValueBoolean, ValueArray, ValueIdentifier:
		return time.Time{}, false // These types cannot be parsed as datetime
	default:
		return time.Time{}, false
	}
}

// compareDateTimes compares two datetime values using the provided comparison function.
func (e *Evaluator) compareDateTimes(
	left, right *EvalResult,
	op func(time.Time, time.Time) bool,
) bool {
	leftTime, leftOk := e.parseDateTime(left)
	if !leftOk {
		return false
	}

	rightTime, rightOk := e.parseDateTime(right)
	if !rightOk {
		return false
	}

	return op(leftTime, rightTime)
}

// Zero-allocation case-insensitive string comparison functions
// These maintain compatibility with nikunjy/rules while avoiding allocations

// equalIgnoreCase compares two strings case-insensitively without allocations.
func (e *Evaluator) equalIgnoreCase(a, b string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range len(a) {
		if toLowerByte(a[i]) != toLowerByte(b[i]) {
			return false
		}
	}

	return true
}

// containsIgnoreCase checks if a contains b case-insensitively without allocations.
func (e *Evaluator) containsIgnoreCase(a, b string) bool {
	if len(b) == 0 {
		return true
	}

	if len(b) > len(a) {
		return false
	}

	for i := 0; i <= len(a)-len(b); i++ {
		if e.equalIgnoreCaseSubstring(a, i, b) {
			return true
		}
	}

	return false
}

// hasPrefixIgnoreCase checks if a starts with b case-insensitively without allocations.
func (e *Evaluator) hasPrefixIgnoreCase(a, b string) bool {
	if len(b) > len(a) {
		return false
	}

	return e.equalIgnoreCaseSubstring(a, 0, b)
}

// hasSuffixIgnoreCase checks if a ends with b case-insensitively without allocations.
func (e *Evaluator) hasSuffixIgnoreCase(a, b string) bool {
	if len(b) > len(a) {
		return false
	}

	return e.equalIgnoreCaseSubstring(a, len(a)-len(b), b)
}

// equalIgnoreCaseSubstring compares substring of a starting at offset with b.
func (e *Evaluator) equalIgnoreCaseSubstring(a string, offset int, b string) bool {
	if offset+len(b) > len(a) {
		return false
	}

	for i := range len(b) {
		if toLowerByte(a[offset+i]) != toLowerByte(b[i]) {
			return false
		}
	}

	return true
}

const asciiLowerOffset = 32 // Offset between ASCII uppercase and lowercase

// toLowerByte converts ASCII uppercase to lowercase without allocation
// Only handles ASCII characters (A-Z), non-ASCII bytes are returned unchanged.
func toLowerByte(b byte) byte {
	if b >= 'A' && b <= 'Z' {
		return b + asciiLowerOffset
	}

	return b
}
