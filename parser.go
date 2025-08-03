package rule

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Parser struct {
	tokens   []Token
	current  int
	curToken Token
}

func NewParser(tokens []Token) *Parser {
	p := &Parser{
		tokens:  tokens,
		current: 0,
	}
	p.curToken = p.tokens[0]

	return p
}

func (p *Parser) Parse() (*ASTNode, error) {
	ast, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	// Check for trailing tokens after a complete expression
	if p.curToken.Type != EOF {
		return nil, ErrTrailingTokens
	}

	return ast, nil
}

func (p *Parser) advance() {
	if p.current < len(p.tokens)-1 {
		p.current++
		p.curToken = p.tokens[p.current]
	}
}

func (p *Parser) peek() Token {
	if p.current+1 < len(p.tokens) {
		return p.tokens[p.current+1]
	}

	return Token{Type: EOF}
}

// isLargeIntegerString checks if a string represents a large integer.
func (p *Parser) isLargeIntegerString(s string) (int64, bool) {
	// Check if it's a valid integer
	if val, err := strconv.ParseInt(s, 10, 64); err == nil {
		// Check if it would lose precision when converted to float64
		if val > maxSafeInteger || val < minSafeInteger {
			return val, true
		}
	}

	return 0, false
}

func (p *Parser) expect(tokenType TokenType) error {
	if p.curToken.Type != tokenType {
		return fmt.Errorf(
			"expected %s, got %s at position %d",
			tokenType,
			p.curToken.Type,
			p.current,
		)
	}

	p.advance()

	return nil
}

func (p *Parser) parseExpression() (*ASTNode, error) {
	return p.parseOrExpression()
}

func (p *Parser) parseOrExpression() (*ASTNode, error) {
	left, err := p.parseAndExpression()
	if err != nil {
		return nil, err
	}

	for p.curToken.Type == OR {
		op := p.curToken.Type
		p.advance()

		right, parseErr := p.parseAndExpression()
		if parseErr != nil {
			return nil, parseErr
		}

		left = NewBinaryOpNode(op, left, right)
	}

	return left, nil
}

func (p *Parser) parseAndExpression() (*ASTNode, error) {
	left, err := p.parseNotExpression()
	if err != nil {
		return nil, err
	}

	for p.curToken.Type == AND {
		op := p.curToken.Type
		p.advance()

		right, parseErr := p.parseNotExpression()
		if parseErr != nil {
			return nil, parseErr
		}

		left = NewBinaryOpNode(op, left, right)
	}

	return left, nil
}

func (p *Parser) parseNotExpression() (*ASTNode, error) {
	if p.curToken.Type == NOT {
		p.advance()

		operand, err := p.parseNotExpression()
		if err != nil {
			return nil, err
		}

		return NewUnaryOpNode(NOT, operand), nil
	}

	return p.parseComparisonExpression()
}

func (p *Parser) parseComparisonExpression() (*ASTNode, error) {
	left, err := p.parsePrimaryExpression()
	if err != nil {
		return nil, err
	}

	if p.isComparisonOperator(p.curToken.Type) {
		op := p.curToken.Type
		p.advance()

		if op == PR {
			return NewUnaryOpNode(PR, left), nil
		}

		right, parseErr := p.parsePrimaryExpression()
		if parseErr != nil {
			return nil, parseErr
		}

		return NewBinaryOpNode(op, left, right), nil
	}

	// Check for missing operator - if we have another value without an operator, that's an error
	if p.isValue(p.curToken.Type) {
		return nil, ErrMissingOperator
	}

	return left, nil
}

func (p *Parser) parsePrimaryExpression() (*ASTNode, error) {
	switch p.curToken.Type {
	case PAREN_OPEN:
		p.advance()

		// Check for empty parentheses
		if p.curToken.Type == PAREN_CLOSE {
			return nil, ErrEmptyParentheses
		}

		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		if expectErr := p.expect(PAREN_CLOSE); expectErr != nil {
			return nil, expectErr
		}

		return expr, nil

	case STRING:
		value := p.curToken.Value
		p.advance()

		// Check if this string represents a large integer
		if intVal, isLargeInt := p.isLargeIntegerString(value); isLargeInt {
			return NewLargeIntegerLiteralNode(intVal), nil
		}

		return NewStringLiteralNode(value), nil

	case NUMBER:
		value := p.curToken.NumValue
		p.advance()

		return NewNumberLiteralNode(value), nil

	case BOOLEAN:
		value := p.curToken.BoolValue
		p.advance()

		return NewBooleanLiteralNode(value), nil

	case ARRAY_START:
		return p.parseArray()

	case IDENTIFIER:
		return p.parseIdentifierOrProperty()

	case EOF,
		ARRAY_END,
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
		PR,
		DQ,
		DN,
		BE,
		BQ,
		AF,
		AQ,
		DL,
		DG,
		AND,
		OR,
		NOT,
		EQUALS,
		NOT_EQUALS:
		return nil, fmt.Errorf("unexpected token %s at position %d", p.curToken.Type, p.current)

	default:
		return nil, fmt.Errorf("unexpected token %s at position %d", p.curToken.Type, p.current)
	}
}

func (p *Parser) parseArray() (*ASTNode, error) {
	if err := p.expect(ARRAY_START); err != nil {
		return nil, err
	}

	var elements []Value

	if p.curToken.Type != ARRAY_END {
		for {
			switch p.curToken.Type {
			case STRING:
				value := p.curToken.Value
				// Check if this string represents a large integer
				if intVal, isLargeInt := p.isLargeIntegerString(value); isLargeInt {
					elements = append(elements, Value{
						Type:     ValueNumber,
						NumValue: float64(intVal),
						IntValue: intVal,
						IsInt:    true,
					})
				} else {
					elements = append(elements, Value{
						Type:     ValueString,
						StrValue: value,
					})
				}

				p.advance()

			case NUMBER:
				elements = append(elements, Value{
					Type:     ValueNumber,
					NumValue: p.curToken.NumValue,
					IsInt:    false,
				})

				p.advance()

			case BOOLEAN:
				elements = append(elements, Value{
					Type:      ValueBoolean,
					BoolValue: p.curToken.BoolValue,
				})

				p.advance()

			case EOF,
				IDENTIFIER,
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
				PR,
				DQ,
				DN,
				BE,
				BQ,
				AF,
				AQ,
				DL,
				DG,
				AND,
				OR,
				NOT,
				EQUALS,
				NOT_EQUALS:
				return nil, fmt.Errorf("unexpected token in array: %s", p.curToken.Type)
			default:
				return nil, fmt.Errorf("unexpected token in array: %s", p.curToken.Type)
			}

			if p.curToken.Type == COMMA {
				p.advance()
			} else {
				break
			}
		}
	}

	if err := p.expect(ARRAY_END); err != nil {
		return nil, err
	}

	return NewArrayLiteralNode(elements), nil
}

func (p *Parser) parseIdentifierOrProperty() (*ASTNode, error) {
	path := []string{p.curToken.Value}
	p.advance()

	for p.curToken.Type == DOT {
		p.advance()

		if p.curToken.Type != IDENTIFIER {
			return nil, fmt.Errorf("expected identifier after dot, got %s", p.curToken.Type)
		}

		path = append(path, p.curToken.Value)
		p.advance()
	}

	if len(path) == 1 {
		return NewIdentifierNode(path[0]), nil
	}

	return NewPropertyNode(path), nil
}

func (p *Parser) isComparisonOperator(tokenType TokenType) bool {
	switch tokenType {
	case EQ, NE, LT, GT, LE, GE, CO, SW, EW, IN, NOT_IN, PR, DQ, DN, BE, BQ, AF, AQ, DL, DG, EQUALS, NOT_EQUALS:
		return true
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
		AND,
		OR,
		NOT:
		return false
	default:
		return false
	}
}

func (p *Parser) isValue(tokenType TokenType) bool {
	switch tokenType {
	case IDENTIFIER, STRING, NUMBER, BOOLEAN, ARRAY_START:
		return true
	case EOF, ARRAY_END, PAREN_OPEN, PAREN_CLOSE, DOT, COMMA,
		EQ, NE, LT, GT, LE, GE, CO, SW, EW, IN, NOT_IN, PR,
		DQ, DN, BE, BQ, AF, AQ, DL, DG, AND, OR, NOT, EQUALS, NOT_EQUALS:
		return false
	default:
		return false
	}
}

func ParseRule(rule string) (*ASTNode, error) {
	// Check for empty query
	if len(strings.TrimSpace(rule)) == 0 {
		return nil, ErrEmptyQuery
	}

	lexer := NewLexer(rule)
	tokens := lexer.Tokenize()

	// Check for lexical errors first
	if lexer.HasErrors() {
		return nil, errors.Join(lexer.GetErrors()...)
	}

	parser := NewParser(tokens)

	ast, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	// Perform semantic validation
	if validationErr := ValidateAST(ast); validationErr != nil {
		return nil, validationErr
	}

	return ast, nil
}
