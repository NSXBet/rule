package rule

import (
	"errors"
	"math"
	"math/big"
	"strconv"
)

const (
	maxDecimalPlaces  = 1000
	decimalRadix      = 10
	mathRoundFunction = "round"
	mathFloorFunction = "floor"
	mathCeilFunction  = "ceil"
)

func (p *Parser) parseMathCall() (*ASTNode, error) {
	name := p.curToken.Value
	if name != mathRoundFunction && name != mathFloorFunction && name != mathCeilFunction {
		return nil, ErrUnknownMathFunction
	}

	p.advance()
	p.advance()

	if p.curToken.Type == PAREN_CLOSE || p.curToken.Type == COMMA || p.curToken.Type == EOF {
		return nil, ErrInvalidMathArguments
	}

	operand, err := p.parseExpression()
	if err != nil {
		return nil, errors.Join(ErrInvalidMathArguments, err)
	}

	var places int64

	if name == mathRoundFunction {
		if p.curToken.Type != COMMA {
			return nil, ErrInvalidMathArguments
		}

		p.advance()

		places, err = p.parseMathPrecision()
		if err != nil {
			return nil, err
		}
	}

	if p.curToken.Type != PAREN_CLOSE {
		return nil, ErrInvalidMathArguments
	}

	p.advance()

	node := &ASTNode{
		Type: NodeMathCall,
		Left: operand,
		Value: Value{
			StrValue: name,
			IntValue: places,
		},
	}
	node.mathScale = new(big.Int).Exp(big.NewInt(decimalRadix), big.NewInt(places), nil)
	node.floatScale = math.Pow10(int(places))

	return node, nil
}

func (p *Parser) parseMathPrecision() (int64, error) {
	text := numericTokenText(p.curToken)
	if text == "" {
		return 0, ErrInvalidMathPrecision
	}

	precision, err := parseDecimal(text)
	if err != nil || !precision.IsInt() || !precision.Num().IsInt64() {
		return 0, ErrInvalidMathPrecision
	}

	places := precision.Num().Int64()
	if places < 0 || places > maxDecimalPlaces {
		return 0, ErrInvalidMathPrecision
	}

	p.advance()

	return places, nil
}

func numericTokenText(token Token) string {
	if token.Type == NUMBER {
		return token.Value
	}

	// The lexer also stores unquoted large int64 literals as STRING tokens.
	if token.Type == STRING && token.End-token.Start == len(token.Value) {
		return token.Value
	}

	return ""
}

func validateMathCallNode(node *ASTNode) error {
	switch node.Value.StrValue {
	case mathRoundFunction:
		if node.Value.IntValue < 0 || node.Value.IntValue > maxDecimalPlaces {
			return ErrInvalidMathPrecision
		}
	case mathFloorFunction, mathCeilFunction:
		if node.Value.IntValue != 0 {
			return ErrInvalidMathArguments
		}
	default:
		return ErrUnknownMathFunction
	}

	if err := validateMathOperand(node.Left); err != nil {
		return err
	}

	return ValidateAST(node.Left)
}

func validateMathOperand(node *ASTNode) error {
	if node == nil {
		return ErrInvalidMathArguments
	}

	switch node.Type {
	case NodeIdentifier, NodeProperty, NodeWhereCount, NodeMathCall:
		return nil
	case NodeLiteral:
		if node.Value.Type != ValueNumber || node.Value.quotedNumber {
			return ErrInvalidMathOperand
		}

		text := node.Value.numberText
		if text == "" {
			text = strconv.FormatFloat(node.Value.NumValue, 'g', -1, 64)
			if node.Value.IsInt {
				text = strconv.FormatInt(node.Value.IntValue, decimalRadix)
			}
		}

		_, err := parseDecimal(text)

		return err
	case NodeBinaryOp, NodeUnaryOp, NodeArray:
		return ErrInvalidMathOperand
	}

	return ErrInvalidMathOperand
}
