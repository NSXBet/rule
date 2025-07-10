package rule

type TokenType uint8

const (
	EOF TokenType = iota
	IDENTIFIER
	STRING
	NUMBER
	BOOLEAN
	ARRAY_START //nolint:revive,staticcheck // Token constants use ALL_CAPS convention
	ARRAY_END   //nolint:revive,staticcheck // Token constants use ALL_CAPS convention
	PAREN_OPEN  //nolint:revive,staticcheck // Token constants use ALL_CAPS convention
	PAREN_CLOSE //nolint:revive,staticcheck // Token constants use ALL_CAPS convention
	DOT
	COMMA

	// EQ represents the equality operator.
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

	// DQ represents the datetime equality operator.
	DQ // datetime equal
	DN // datetime not equal
	BE // before
	BQ // before or equal
	AF // after
	AQ // after or equal

	// AND represents the logical AND operator.
	AND
	OR
	NOT

	// EQUALS is an alias for the equality operator.
	EQUALS     // ==
	NOT_EQUALS //nolint:revive,staticcheck // Token constants use ALL_CAPS convention
)

type Token struct {
	Type      TokenType
	Value     string
	Start     int
	End       int
	NumValue  float64
	BoolValue bool
}

//nolint:gochecknoglobals // Static keyword lookup table
var keywordMap = map[string]TokenType{
	"eq":       EQ,
	"ne":       NE,
	"lt":       LT,
	"gt":       GT,
	"le":       LE,
	"ge":       GE,
	"co":       CO,
	"sw":       SW,
	"ew":       EW,
	"in":       IN,
	"pr":       PR,
	"dq":       DQ,
	"dn":       DN,
	"be":       BE,
	"bq":       BQ,
	"af":       AF,
	"aq":       AQ,
	"and":      AND,
	"or":       OR,
	"not":      NOT,
	trueString: BOOLEAN,
	"false":    BOOLEAN,
}

//nolint:gochecknoglobals // Static token string lookup table
var tokenStringMap = map[TokenType]string{
	EOF:         "EOF",
	IDENTIFIER:  "IDENTIFIER",
	STRING:      "STRING",
	NUMBER:      "NUMBER",
	BOOLEAN:     "BOOLEAN",
	ARRAY_START: "[",
	ARRAY_END:   "]",
	PAREN_OPEN:  "(",
	PAREN_CLOSE: ")",
	DOT:         ".",
	COMMA:       ",",
	EQ:          "eq",
	NE:          "ne",
	LT:          "lt",
	GT:          "gt",
	LE:          "le",
	GE:          "ge",
	CO:          "co",
	SW:          "sw",
	EW:          "ew",
	IN:          "in",
	PR:          "pr",
	DQ:          "dq",
	DN:          "dn",
	BE:          "be",
	BQ:          "bq",
	AF:          "af",
	AQ:          "aq",
	AND:         "and",
	OR:          "or",
	NOT:         "not",
	EQUALS:      "==",
	NOT_EQUALS:  "!=",
}

func (t TokenType) String() string {
	if str, exists := tokenStringMap[t]; exists {
		return str
	}

	return "UNKNOWN"
}

func (t Token) String() string {
	if t.Value != "" {
		return t.Type.String() + "(" + t.Value + ")"
	}

	return t.Type.String()
}
