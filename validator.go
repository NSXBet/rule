package rule

import "errors"

// ValidateAST performs semantic validation on the parsed AST.
func ValidateAST(node *ASTNode) error {
	if node == nil {
		return nil
	}

	switch node.Type {
	case NodeBinaryOp:
		if err := validateBinaryOperation(node); err != nil {
			return err
		}

		// Recursively validate children
		if err := ValidateAST(node.Left); err != nil {
			return err
		}

		if err := ValidateAST(node.Right); err != nil {
			return err
		}

	case NodeUnaryOp:
		if err := validateUnaryOperation(node); err != nil {
			return err
		}

		// Recursively validate the operand
		if err := ValidateAST(node.Left); err != nil {
			return err
		}

	case NodeLiteral, NodeIdentifier, NodeProperty, NodeArray:
		// These are terminal nodes, no further validation needed
		return nil
	}

	return nil
}

func validateBinaryOperation(node *ASTNode) error {
	if node.Left == nil || node.Right == nil {
		return errors.New("binary operation missing operands")
	}

	switch node.Operator {
	case IN, NOT_IN:
		return validateInOperation(node)
	case CO, SW, EW:
		return validateStringOperation(node)
	case EOF, IDENTIFIER, STRING, NUMBER, BOOLEAN, ARRAY_START, ARRAY_END,
		PAREN_OPEN, PAREN_CLOSE, DOT, COMMA, EQ, NE, LT, GT, LE, GE, PR,
		DQ, DN, BE, BQ, AF, AQ, DL, DG, AND, OR, NOT, EQUALS, NOT_EQUALS:
		// Other operators don't need special validation
		return nil
	}

	return nil
}

func validateUnaryOperation(node *ASTNode) error {
	if node.Left == nil {
		return errors.New("unary operation missing operand")
	}

	switch node.Operator {
	case PR:
		return validatePresenceOperation(node)
	case EOF, IDENTIFIER, STRING, NUMBER, BOOLEAN, ARRAY_START, ARRAY_END,
		PAREN_OPEN, PAREN_CLOSE, DOT, COMMA, EQ, NE, LT, GT, LE, GE,
		CO, SW, EW, IN, NOT_IN, DQ, DN, BE, BQ, AF, AQ, DL, DG, AND, OR, NOT, EQUALS, NOT_EQUALS:
		// Other operators don't apply to unary operations
		return nil
	}

	return nil
}

func validateInOperation(node *ASTNode) error {
	// The right operand of IN must be an array
	if node.Right.Type != NodeArray {
		// Check if it's a literal that's not an array
		switch node.Right.Type {
		case NodeLiteral:
			if node.Right.Value.Type != ValueArray {
				return ErrInvalidInOperand
			}
		case NodeIdentifier, NodeProperty:
			// Allow identifiers/properties as they might evaluate to arrays at runtime
			return nil
		case NodeBinaryOp, NodeUnaryOp, NodeArray:
			return ErrInvalidInOperand
		}
	}

	return nil
}

func validateStringOperation(node *ASTNode) error {
	// String operations (co, sw, ew) should typically work on strings
	// But we'll allow identifiers/properties as they might be strings at runtime
	left := node.Left
	right := node.Right

	// Check if we're trying to use string operators on obviously non-string literals
	if left.Type == NodeLiteral {
		if left.Value.Type != ValueString {
			return ErrInvalidStringOp
		}
	}

	if right.Type == NodeLiteral {
		if right.Value.Type != ValueString {
			return ErrInvalidStringOp
		}
	}

	return nil
}

func validatePresenceOperation(node *ASTNode) error {
	// Presence operator should only work on identifiers or properties
	operand := node.Left
	if operand.Type != NodeIdentifier && operand.Type != NodeProperty {
		return ErrInvalidPresenceOp
	}

	return nil
}
