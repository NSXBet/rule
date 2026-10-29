package test

import (
	"encoding/json"
	"errors"
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/NSXBet/rule"
)

func TestDecimalMath(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		context rule.D
		want    bool
	}{
		{"decimal midpoint", "round(value, 2) eq 1.01", rule.D{"value": 1.005}, true},
		{"negative midpoint", "round(value, 2) eq -1.01", rule.D{"value": -1.005}, true},
		{"exact literal", "round(9007199254740993.125, 2) eq 9007199254740993.13", nil, true},
		{
			"exact context",
			"round(value, 2) eq 9007199254740993.13",
			rule.D{"value": json.Number("9007199254740993.125")},
			true,
		},
		{"negative floor", "floor(value) eq -2", rule.D{"value": -1.2}, true},
		{"negative ceil", "ceil(value) eq -1", rule.D{"value": -1.2}, true},
		{"below midpoint", "round(value, 2) eq 1.00", rule.D{"value": json.Number("1.004999999999999999999999")}, true},
		{"above midpoint", "round(value, 2) eq 1.01", rule.D{"value": json.Number("1.005000000000000000000001")}, true},
		{"carry", "round(9.995, 2) eq 10", nil, true},
		{"negative carry", "round(-9.995, 2) eq -10", nil, true},
		{"positive half", "round(2.5, 0) eq 3", nil, true},
		{"negative half", "round(-2.5, 0) eq -3", nil, true},
		{"trailing zero places", "round(1.005, 2.0) eq 1.01", nil, true},
		{"float32 source precision", "round(value, 2) eq 1.01", rule.D{"value": float32(1.005)}, true},
		{
			"tiny decimal",
			"round(value, 30) eq 0.000000000000000000000000000001",
			rule.D{"value": json.Number("5e-31")},
			true,
		},
		{"large exact integer", "floor(value) eq 18446744073709551615", rule.D{"value": uint64(math.MaxUint64)}, true},
		{"signed integer", "ceil(value) eq -9223372036854775808", rule.D{"value": int64(math.MinInt64)}, true},
		{"big integer carry", "round(999999999999999999999999.5, 0) eq 1000000000000000000000000", nil, true},
		{
			"fraction preserved",
			"round(1.12345678901234567890123456789, 29) eq 1.12345678901234567890123456789",
			nil,
			true,
		},
		{"distinguish rounded decimals", "round(1.12345678901234567891, 20) ne 1.12345678901234567892", nil, true},
		{"nested calls", "round(round(1.2449, 3), 2) eq 1.25", nil, true},
		{
			"rounded property comparison",
			"round(actual, 2) eq round(expected, 2)",
			rule.D{"actual": 1.005, "expected": 1.006},
			true,
		},
		{"right operand", "1.01 eq round(1.005, 2)", nil, true},
		{"strict greater", "round(1.12345678901234567892, 20) gt 1.12345678901234567891", nil, true},
		{"strict less", "1.12345678901234567891 lt round(1.12345678901234567892, 20)", nil, true},
		{"less or equal", "round(1.005, 2) le 1.01", nil, true},
		{"greater or equal", "round(1.005, 2) ge 1.01", nil, true},
		{"unequal precision", "round(1.005, 2) == 1.005", nil, false},
		{"alias inequality", "round(1.005, 2) != 1.005", nil, true},
		{
			"literal membership",
			"round(9007199254740993.125, 2) in [9007199254740993.12, 9007199254740993.13]",
			nil,
			true,
		},
		{"negative membership", "round(1.005, 2) not in [1.00, 1.02]", nil, true},
		{
			"context membership",
			"round(value, 2) in options",
			rule.D{"value": 1.005, "options": []any{"1.01", json.Number("1.01")}},
			true,
		},
		{"no string coercion", `round(1.005, 2) eq "1.01"`, nil, false},
		{"membership type mismatch", `round(1.005, 2) in ["1.01", true]`, nil, false},
		{"quoted large integer comparison", `floor(9007199254740993) eq "9007199254740993"`, nil, false},
		{"quoted large integer membership", `floor(9007199254740993) in ["9007199254740993"]`, nil, false},
		{"exact truthiness", "round(value, 1000)", rule.D{"value": json.Number("1e-999")}, true},
		{"rounded zero truthiness", "round(value, 2)", rule.D{"value": json.Number("0.0001")}, false},
		{"positive floor", "floor(1.99) eq 1", nil, true},
		{"positive ceil", "ceil(1.01) eq 2", nil, true},
		{"integer floor", "floor(-2) eq -2", nil, true},
		{"integer ceil", "ceil(-2) eq -2", nil, true},
		{"zero", "round(-0.0001, 2) eq 0", nil, true},
		{
			"ordinary function field names",
			"round eq 2 and floor eq 3 and ceil eq 4",
			rule.D{"round": 2, "floor": 3, "ceil": 4},
			true,
		},
		{"nested property", "round(user.odd, 2) ge 1.01", rule.D{"user": rule.D{"odd": 1.005}}, true},
		{
			"quantifier",
			"items all (round(odd, 2) ge 1.01)",
			rule.D{"items": []any{rule.D{"odd": 1.005}, rule.D{"odd": 1.015}}},
			true,
		},
		{
			"where predicate",
			"items where (round(odd, 2) ge 1.01).length eq 1",
			rule.D{"items": []any{rule.D{"odd": 1.005}, rule.D{"odd": 1.004}}},
			true,
		},
		{
			"filtered length argument",
			"round(items where (odd ge 1).length, 2) eq 1",
			rule.D{"items": []any{rule.D{"odd": 1.005}, rule.D{"odd": 0.9}}},
			true,
		},
		{"short circuit", "true or round(value, 2) gt 0", rule.D{"value": nil}, true},
	}
	engine := rule.NewEngine()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := engine.Evaluate(tt.query, tt.context)
			if err != nil {
				t.Fatal(err)
			}

			if got != tt.want {
				t.Fatalf("Evaluate(%q) = %v, want %v", tt.query, got, tt.want)
			}
		})
	}
}

func TestDecimalMathMissing(t *testing.T) {
	for _, query := range []string{
		"round(value, 2) eq 1",
		"floor(user.value) ne 1",
		"ceil(round(value, 2)) gt 0",
	} {
		for _, lenient := range []bool{false, true} {
			engine := rule.NewEngine()
			if lenient {
				engine = rule.NewEngineWithOptions(rule.WithLenientMode())
			}

			got, err := engine.Evaluate(query, rule.D{})
			if err != nil || got != lenient {
				t.Errorf("%s (lenient=%v) = %v, %v; want %v, nil", query, lenient, got, err, lenient)
			}
		}
	}
}

func TestDecimalMathInvalidQueries(t *testing.T) {
	tests := []struct {
		query string
		want  error
	}{
		{"sqrt(4) eq 2", rule.ErrUnknownMathFunction},
		{"round(1.2) eq 1", rule.ErrInvalidMathArguments},
		{"round(1, 2, 3) eq 1", rule.ErrInvalidMathArguments},
		{"round() eq 1", rule.ErrInvalidMathArguments},
		{"floor() eq 1", rule.ErrInvalidMathArguments},
		{"ceil(1, 2) eq 1", rule.ErrInvalidMathArguments},
		{"round(1, -1) eq 1", rule.ErrInvalidMathPrecision},
		{"round(1, 1001) eq 1", rule.ErrInvalidMathPrecision},
		{"round(1, 2.0000000000000000000001) eq 1", rule.ErrInvalidMathPrecision},
		{"round(1, places) eq 1", rule.ErrInvalidMathPrecision},
		{`round("1.005", 2) eq 1.01`, rule.ErrInvalidMathOperand},
		{"floor(true) eq 1", rule.ErrInvalidMathOperand},
		{"ceil([1]) eq 1", rule.ErrInvalidMathOperand},
		{"floor(1 eq 1) eq 1", rule.ErrInvalidMathOperand},
		{"round(value, 2) pr", rule.ErrInvalidPresenceOp},
		{"floor(value) any (x eq 1)", rule.ErrInvalidQuantifierTarget},
		{"1 in round(value, 2)", rule.ErrInvalidInOperand},
	}
	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			_, err := rule.NewEngine().CompileRule(tt.query)
			if !errors.Is(err, tt.want) {
				t.Fatalf("CompileRule(%q) error = %v, want %v", tt.query, err, tt.want)
			}
		})
	}
}

func TestDecimalMathInvalidInputs(t *testing.T) {
	tests := []struct {
		name  string
		value any
		want  error
	}{
		{"nil", nil, rule.ErrInvalidMathOperand},
		{"string", "1.005", rule.ErrInvalidMathOperand},
		{"boolean", true, rule.ErrInvalidMathOperand},
		{"array", []any{1}, rule.ErrInvalidMathOperand},
		{"nan", math.NaN(), rule.ErrInvalidDecimal},
		{"positive infinity", math.Inf(1), rule.ErrInvalidDecimal},
		{"negative infinity", math.Inf(-1), rule.ErrInvalidDecimal},
		{"invalid decimal", json.Number("broken"), rule.ErrInvalidDecimal},
		{"rational", json.Number("1/3"), rule.ErrInvalidDecimal},
		{"hex", json.Number("0x10"), rule.ErrInvalidDecimal},
		{"separator", json.Number("1_000"), rule.ErrInvalidDecimal},
		{"json nan", json.Number("NaN"), rule.ErrInvalidDecimal},
		{"json infinity", json.Number("Inf"), rule.ErrInvalidDecimal},
		{"space", json.Number(" 1.005"), rule.ErrInvalidDecimal},
		{"huge exponent", json.Number("1e10001"), rule.ErrInvalidDecimal},
		{"huge negative exponent", json.Number("1e-10001"), rule.ErrInvalidDecimal},
		{"huge input", json.Number(strings.Repeat("1", 10001)), rule.ErrInvalidDecimal},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, query := range []string{
				"round(value, 2) ge 0",
				"items any (floor(value) ge 0)",
				"items all (ceil(value) ge 0)",
				"items none (round(value, 2) ge 0)",
				"items where (round(value, 2) ge 0).length gt 0",
			} {
				context := rule.D{"value": tt.value, "items": []any{rule.D{"value": tt.value}}}

				_, err := rule.NewEngine().Evaluate(query, context)
				if !errors.Is(err, tt.want) {
					t.Errorf("Evaluate(%q) error = %v, want %v", query, err, tt.want)
				}
			}
		})
	}
}

func TestDecimalMathCompiledConcurrent(t *testing.T) {
	t.Parallel()

	engine := rule.NewEngine()

	compiled, err := engine.CompileRule("round(amount, 2) eq expected")
	if err != nil {
		t.Fatal(err)
	}

	for i := range 16 {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()

			for j := range 100 {
				amount := json.Number("9007199254740993.125")
				expected := json.Number("9007199254740993.13")

				if (i+j)%2 == 0 {
					amount = "-9007199254740993.125"
					expected = "-9007199254740993.13"
				}

				got, evalErr := engine.EvaluateCompiled(compiled, rule.D{"amount": amount, "expected": expected})
				if evalErr != nil || !got {
					t.Fatalf("EvaluateCompiled(amount=%s, expected=%s) = %v, %v", amount, expected, got, evalErr)
				}
			}
		})
	}
}

func FuzzDecimalRounding(f *testing.F) {
	for _, coefficient := range []int64{0, 1005, -1005, 9995, -9995, math.MaxInt64, math.MinInt64} {
		f.Add(coefficient)
	}

	engine := rule.NewEngine()

	f.Fuzz(func(t *testing.T, coefficient int64) {
		rounded := coefficient / 10

		remainder := coefficient % 10
		if remainder >= 5 {
			rounded++
		} else if remainder <= -5 {
			rounded--
		}

		context := rule.D{
			"value":    json.Number(strconv.FormatInt(coefficient, 10) + "e-3"),
			"expected": json.Number(strconv.FormatInt(rounded, 10) + "e-2"),
		}

		got, err := engine.Evaluate("round(value, 2) eq expected", context)
		if err != nil || !got {
			t.Fatalf("coefficient=%d: rounding comparison = %v, %v", coefficient, got, err)
		}
	})
}
