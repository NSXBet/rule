package rule

type TokenType uint8

const (
	EOF TokenType = iota
	IDENTIFIER
	STRING
	NUMBER
	BOOLEAN
	ARRAY_START
	ARRAY_END
	PAREN_OPEN
	PAREN_CLOSE
	DOT
	COMMA

	// Operators
	EQ
	NE
	LT
	GT
	LE
	GE
	CO
	SW
	EW
	IN
	PR

	// Datetime Operators
	DQ // datetime equal
	DN // datetime not equal
	BE // before
	BQ // before or equal
	AF // after
	AQ // after or equal

	// Logical
	AND
	OR
	NOT

	// Aliases
	EQUALS     // ==
	NOT_EQUALS // !=
)

type Token struct {
	Type      TokenType
	Value     string
	Start     int
	End       int
	NumValue  float64
	BoolValue bool
}

var keywordMap = map[string]TokenType{
	"eq":    EQ,
	"ne":    NE,
	"lt":    LT,
	"gt":    GT,
	"le":    LE,
	"ge":    GE,
	"co":    CO,
	"sw":    SW,
	"ew":    EW,
	"in":    IN,
	"pr":    PR,
	"dq":    DQ,
	"dn":    DN,
	"be":    BE,
	"bq":    BQ,
	"af":    AF,
	"aq":    AQ,
	"and":   AND,
	"or":    OR,
	"not":   NOT,
	"true":  BOOLEAN,
	"false": BOOLEAN,
}

func (t TokenType) String() string {
	switch t {
	case EOF:
		return "EOF"
	case IDENTIFIER:
		return "IDENTIFIER"
	case STRING:
		return "STRING"
	case NUMBER:
		return "NUMBER"
	case BOOLEAN:
		return "BOOLEAN"
	case ARRAY_START:
		return "["
	case ARRAY_END:
		return "]"
	case PAREN_OPEN:
		return "("
	case PAREN_CLOSE:
		return ")"
	case DOT:
		return "."
	case COMMA:
		return ","
	case EQ:
		return "eq"
	case NE:
		return "ne"
	case LT:
		return "lt"
	case GT:
		return "gt"
	case LE:
		return "le"
	case GE:
		return "ge"
	case CO:
		return "co"
	case SW:
		return "sw"
	case EW:
		return "ew"
	case IN:
		return "in"
	case PR:
		return "pr"
	case DQ:
		return "dq"
	case DN:
		return "dn"
	case BE:
		return "be"
	case BQ:
		return "bq"
	case AF:
		return "af"
	case AQ:
		return "aq"
	case AND:
		return "and"
	case OR:
		return "or"
	case NOT:
		return "not"
	case EQUALS:
		return "=="
	case NOT_EQUALS:
		return "!="
	default:
		return "UNKNOWN"
	}
}

func (t Token) String() string {
	if t.Value != "" {
		return t.Type.String() + "(" + t.Value + ")"
	}
	return t.Type.String()
}
