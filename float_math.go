package rule

import (
	"encoding/json"
	"math"
	"strconv"
)

func (e *Evaluator) setJSONNumberResult(result *EvalResult, value json.Number) error {
	result.Type = ValueNumber

	if !e.floatMath {
		decimal, err := parseDecimal(string(value))
		if err != nil {
			return err
		}

		result.OriginalValue = decimal
		result.Num, _ = decimal.Float64()

		return nil
	}

	if err := validateDecimalText(string(value)); err != nil {
		return err
	}

	number, err := strconv.ParseFloat(string(value), 64)
	if err != nil || math.IsInf(number, 0) {
		return ErrFloatMathOverflow
	}

	result.Num = number
	result.IsInt = false

	return nil
}

func evaluateFloatMathCall(node *ASTNode, operand, result *EvalResult) error {
	if !isDecimalNumber(operand) {
		return ErrInvalidMathOperand
	}

	value := operand.Num
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return ErrInvalidDecimal
	}

	switch node.Value.StrValue {
	case mathRoundFunction:
		scaled := value * node.floatScale
		if math.IsInf(node.floatScale, 0) || math.IsInf(scaled, 0) {
			return ErrFloatMathOverflow
		}

		value = math.Round(scaled) / node.floatScale
	case mathFloorFunction:
		value = math.Floor(value)
	case mathCeilFunction:
		value = math.Ceil(value)
	default:
		return ErrUnknownMathFunction
	}

	if math.IsNaN(value) || math.IsInf(value, 0) {
		return ErrFloatMathOverflow
	}

	result.Type = ValueNumber
	result.IsValid = true
	result.IsInt = false
	result.Num = value
	result.OriginalValue = node

	return nil
}

func isFloatMathResult(result *EvalResult) bool {
	source, ok := result.OriginalValue.(*ASTNode)
	return ok && source.Type == NodeMathCall
}

func isMathComparison(operator TokenType) bool {
	return operator == EQ || operator == EQUALS || operator == NE || operator == NOT_EQUALS ||
		operator == LT || operator == LE || operator == GT || operator == GE
}

// General comparison dispatch recurses through membership and makes its temporaries escape.
func compareFloatMath(operator TokenType, left, right, result *EvalResult) error {
	if !isMathComparison(operator) {
		return ErrInvalidMathOperand
	}

	if !isDecimalNumber(left) || !isDecimalNumber(right) {
		result.Bool = operator == NE || operator == NOT_EQUALS
		return nil
	}

	if math.IsNaN(left.Num) || math.IsInf(left.Num, 0) || math.IsNaN(right.Num) || math.IsInf(right.Num, 0) {
		return ErrInvalidDecimal
	}

	comparison := 0
	if left.Num < right.Num {
		comparison = -1
	} else if left.Num > right.Num {
		comparison = 1
	}

	result.Bool = compareMathOrdering(operator, comparison)

	return nil
}
