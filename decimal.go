package rule

import (
	"math/big"
	"strconv"
)

const (
	maxDecimalTextLength = 10000
	maxDecimalExponent   = 10000
)

func parseDecimal(text string) (*big.Rat, error) {
	if err := validateDecimalText(text); err != nil {
		return nil, err
	}

	value, ok := new(big.Rat).SetString(text)
	if !ok {
		return nil, ErrInvalidDecimal
	}

	return value, nil
}

func validateDecimalText(text string) error {
	if len(text) == 0 || len(text) > maxDecimalTextLength {
		return ErrInvalidDecimal
	}

	index := 0
	if text[index] == '+' || text[index] == '-' {
		index++
	}

	digits := 0

	for index < len(text) && isDecimalDigit(text[index]) {
		index++
		digits++
	}

	if index < len(text) && text[index] == '.' {
		index++
		for index < len(text) && isDecimalDigit(text[index]) {
			index++
			digits++
		}
	}

	if digits == 0 {
		return ErrInvalidDecimal
	}

	if index < len(text) && (text[index] == 'e' || text[index] == 'E') {
		var err error

		index, err = scanDecimalExponent(text, index+1)
		if err != nil {
			return err
		}
	}

	if index != len(text) {
		return ErrInvalidDecimal
	}

	return nil
}

func isDecimalDigit(char byte) bool {
	return char >= '0' && char <= '9'
}

func scanDecimalExponent(text string, index int) (int, error) {
	if index < len(text) && (text[index] == '+' || text[index] == '-') {
		index++
	}

	start := index
	exponent := 0

	for index < len(text) && isDecimalDigit(text[index]) {
		exponent = exponent*decimalRadix + int(text[index]-'0')
		if exponent > maxDecimalExponent {
			return index, ErrInvalidDecimal
		}

		index++
	}

	if index == start {
		return index, ErrInvalidDecimal
	}

	return index, nil
}

func decimalFromResult(result *EvalResult) (*big.Rat, error) {
	if result.Type != ValueNumber {
		return nil, ErrInvalidMathOperand
	}

	switch value := result.OriginalValue.(type) {
	case *big.Rat:
		return value, nil
	case *Value:
		if value.quotedNumber {
			return nil, ErrInvalidMathOperand
		}

		if value.numberText != "" {
			return parseDecimal(value.numberText)
		}
	case uint:
		return new(big.Rat).SetUint64(uint64(value)), nil
	case uint64:
		return new(big.Rat).SetUint64(value), nil
	case float32:
		return parseDecimal(strconv.FormatFloat(float64(value), 'g', -1, 32))
	case float64:
		return parseDecimal(strconv.FormatFloat(value, 'g', -1, 64))
	}

	if result.IsInt {
		return new(big.Rat).SetInt64(result.IntValue), nil
	}

	return parseDecimal(strconv.FormatFloat(result.Num, 'g', -1, 64))
}

func (e *Evaluator) evaluateMathCall(node *ASTNode, context D, result *EvalResult) error {
	if node.Left == nil || node.mathScale == nil {
		return ErrInvalidMathArguments
	}

	var operand EvalResult

	if err := e.evaluateNode(node.Left, context, &operand); err != nil {
		return err
	}

	if !operand.IsValid {
		return nil
	}

	if e.floatMath {
		return evaluateFloatMathCall(node, &operand, result)
	}

	value, err := decimalFromResult(&operand)
	if err != nil {
		return err
	}

	scaled := new(big.Int).Mul(value.Num(), node.mathScale)
	quotient, remainder := new(big.Int), new(big.Int)
	quotient.QuoRem(scaled, value.Denom(), remainder)

	if err = adjustMathQuotient(node.Value.StrValue, value, quotient, remainder); err != nil {
		return err
	}

	result.Type = ValueNumber
	result.IsValid = true
	result.OriginalValue = new(big.Rat).SetFrac(quotient, node.mathScale)

	return nil
}

func adjustMathQuotient(name string, value *big.Rat, quotient, remainder *big.Int) error {
	switch name {
	case mathRoundFunction:
		doubled := new(big.Int).Lsh(new(big.Int).Abs(remainder), 1)
		if doubled.Cmp(value.Denom()) >= 0 {
			quotient.Add(quotient, big.NewInt(int64(value.Sign())))
		}
	case mathFloorFunction:
		if remainder.Sign() < 0 {
			quotient.Sub(quotient, big.NewInt(1))
		}
	case mathCeilFunction:
		if remainder.Sign() > 0 {
			quotient.Add(quotient, big.NewInt(1))
		}
	default:
		return ErrUnknownMathFunction
	}

	return nil
}

func isDecimalNumber(result *EvalResult) bool {
	if source, ok := result.OriginalValue.(*Value); ok && source.quotedNumber {
		return false
	}

	return result.Type == ValueNumber
}

func compareDecimal(operator TokenType, left, right, result *EvalResult) error {
	if !isMathComparison(operator) {
		return ErrInvalidMathOperand
	}

	if !isDecimalNumber(left) || !isDecimalNumber(right) {
		result.Bool = operator == NE || operator == NOT_EQUALS
		return nil
	}

	leftDecimal, err := decimalFromResult(left)
	if err != nil {
		return err
	}

	rightDecimal, err := decimalFromResult(right)
	if err != nil {
		return err
	}

	result.Bool = compareMathOrdering(operator, leftDecimal.Cmp(rightDecimal))

	return nil
}

func compareMathOrdering(operator TokenType, comparison int) bool {
	switch comparison {
	case 0:
		return operator == EQ || operator == EQUALS || operator == LE || operator == GE
	case -1:
		return operator == NE || operator == NOT_EQUALS || operator == LT || operator == LE
	default:
		return operator == NE || operator == NOT_EQUALS || operator == GT || operator == GE
	}
}

func (e *Evaluator) equalMembership(left, right *EvalResult) (bool, error) {
	if e.floatMath && (isFloatMathResult(left) || isFloatMathResult(right)) {
		var result EvalResult

		err := compareFloatMath(EQ, left, right, &result)

		return result.Bool, err
	}

	_, leftExact := left.OriginalValue.(*big.Rat)

	_, rightExact := right.OriginalValue.(*big.Rat)
	if !leftExact && !rightExact {
		return e.compareEqualStrict(left, right), nil
	}

	var result EvalResult
	if err := compareDecimal(EQ, left, right, &result); err != nil {
		return false, err
	}

	return result.Bool, nil
}
