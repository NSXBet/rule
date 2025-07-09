package rule

type NodeType uint8

const (
	NodeBinaryOp NodeType = iota
	NodeUnaryOp
	NodeIdentifier
	NodeLiteral
	NodeArray
	NodeProperty
)

type ASTNode struct {
	Type     NodeType
	Operator TokenType
	Left     *ASTNode
	Right    *ASTNode
	Value    Value
	Children []*ASTNode
}

type ValueType uint8

const (
	ValueString ValueType = iota
	ValueNumber
	ValueBoolean
	ValueArray
	ValueIdentifier
)

type Value struct {
	Type      ValueType
	StrValue  string
	NumValue  float64
	BoolValue bool
	ArrValue  []Value
	// IntValue stores large integers to preserve precision
	IntValue int64
	// IsInt indicates if this numeric value should be treated as an integer
	IsInt bool
}

func NewBinaryOpNode(op TokenType, left, right *ASTNode) *ASTNode {
	return &ASTNode{
		Type:     NodeBinaryOp,
		Operator: op,
		Left:     left,
		Right:    right,
	}
}

func NewUnaryOpNode(op TokenType, operand *ASTNode) *ASTNode {
	return &ASTNode{
		Type:     NodeUnaryOp,
		Operator: op,
		Left:     operand,
	}
}

func NewIdentifierNode(name string) *ASTNode {
	return &ASTNode{
		Type: NodeIdentifier,
		Value: Value{
			Type:     ValueIdentifier,
			StrValue: name,
		},
	}
}

func NewPropertyNode(path []string) *ASTNode {
	node := &ASTNode{
		Type:     NodeProperty,
		Children: make([]*ASTNode, len(path)),
	}

	for i, segment := range path {
		node.Children[i] = NewIdentifierNode(segment)
	}

	return node
}

func NewStringLiteralNode(value string) *ASTNode {
	return &ASTNode{
		Type: NodeLiteral,
		Value: Value{
			Type:     ValueString,
			StrValue: value,
		},
	}
}

func NewNumberLiteralNode(value float64) *ASTNode {
	return &ASTNode{
		Type: NodeLiteral,
		Value: Value{
			Type:     ValueNumber,
			NumValue: value,
			IsInt:    false,
		},
	}
}

func NewLargeIntegerLiteralNode(value int64) *ASTNode {
	return &ASTNode{
		Type: NodeLiteral,
		Value: Value{
			Type:     ValueNumber,
			NumValue: float64(value),
			IntValue: value,
			IsInt:    true,
		},
	}
}

func NewBooleanLiteralNode(value bool) *ASTNode {
	return &ASTNode{
		Type: NodeLiteral,
		Value: Value{
			Type:      ValueBoolean,
			BoolValue: value,
		},
	}
}

func NewArrayLiteralNode(elements []Value) *ASTNode {
	return &ASTNode{
		Type: NodeLiteral,
		Value: Value{
			Type:     ValueArray,
			ArrValue: elements,
		},
	}
}

func (n *ASTNode) IsOperator() bool {
	return n.Type == NodeBinaryOp || n.Type == NodeUnaryOp
}

func (n *ASTNode) IsLiteral() bool {
	return n.Type == NodeLiteral
}

func (n *ASTNode) IsIdentifier() bool {
	return n.Type == NodeIdentifier || n.Type == NodeProperty
}
