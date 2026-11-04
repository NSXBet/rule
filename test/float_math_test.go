package test

import (
	"encoding/json"
	"errors"
	"math"
	"strings"
	"testing"

	"github.com/NSXBet/rule"
)

func TestFloatMath(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		context rule.D
		want    bool
	}{
		{"binary midpoint", "round(value, 2) eq 1.00", rule.D{"value": 1.005}, true},
		{"json midpoint", "round(value, 2) eq 1.00", rule.D{"value": json.Number("1.005")}, true},
		{"literal midpoint", "round(1.005, 2) eq 1.00", nil, true},
		{"negative midpoint", "round(-1.005, 2) eq -1.00", nil, true},
		{"positive half", "round(2.5, 0) eq 3", nil, true},
		{"negative half", "round(-2.5, 0) eq -3", nil, true},
		{"floor", "floor(value) eq -2", rule.D{"value": -1.2}, true},
		{"ceil", "ceil(value) eq -1", rule.D{"value": json.Number("-1.2")}, true},
		{"nested calls", "floor(ceil(value)) eq 2", rule.D{"value": 1.2}, true},
		{"right operand", "1.23 eq round(value, 2)", rule.D{"value": 1.234}, true},
		{"membership", "round(value, 2) in [1.23, 1.24]", rule.D{"value": 1.234}, true},
		{
			"json membership",
			"round(value, 2) in options",
			rule.D{"value": 1.234, "options": []any{json.Number("1.23")}},
			true,
		},
		{"nonmembership", "floor(value) not in [1, 2]", rule.D{"value": 3.4}, true},
		{"no string coercion", `round(value, 2) eq "1.23"`, rule.D{"value": 1.234}, false},
		{
			"no quoted integer coercion",
			`floor(value) eq "9007199254740993"`,
			rule.D{"value": uint64(9007199254740993)},
			false,
		},
		{
			"no quoted membership coercion",
			`floor(value) in ["9007199254740993"]`,
			rule.D{"value": uint64(9007199254740993)},
			false,
		},
		{"json numeric policy", "value eq 9007199254740992", rule.D{"value": json.Number("9007199254740993")}, true},
		{"truthiness", "round(value, 2)", rule.D{"value": 1.234}, true},
		{"rounded zero", "round(value, 2)", rule.D{"value": 0.001}, false},
		{"max float unscaled", "round(value, 0) eq value", rule.D{"value": math.MaxFloat64}, true},
		{"largest finite scale", "round(value, 308) gt 0", rule.D{"value": math.SmallestNonzeroFloat64}, false},
		{"quantifier", "items all (round(value, 2) eq 1)", rule.D{"items": []any{rule.D{"value": 1.005}}}, true},
	}
	engine := rule.NewEngineWithOptions(rule.WithFloatMath())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := engine.Evaluate(tt.query, tt.context)
			if err != nil || got != tt.want {
				t.Fatalf("Evaluate(%q) = %v, %v; want %v, nil", tt.query, got, err, tt.want)
			}
		})
	}
}

func TestFloatMathInvalidInputs(t *testing.T) {
	tests := []struct {
		name  string
		query string
		value any
		want  error
	}{
		{"nil", "round(value, 2) ge 0", nil, rule.ErrInvalidMathOperand},
		{"string", "floor(value) ge 0", "1.2", rule.ErrInvalidMathOperand},
		{"nan", "ceil(value) ge 0", math.NaN(), rule.ErrInvalidDecimal},
		{"infinity", "round(value, 2) ge 0", math.Inf(1), rule.ErrInvalidDecimal},
		{"invalid json", "round(value, 2) ge 0", json.Number("broken"), rule.ErrInvalidDecimal},
		{"hex json", "round(value, 2) ge 0", json.Number("0x10"), rule.ErrInvalidDecimal},
		{"rational json", "round(value, 2) ge 0", json.Number("1/2"), rule.ErrInvalidDecimal},
		{"separator json", "round(value, 2) ge 0", json.Number("1_000"), rule.ErrInvalidDecimal},
		{"json nan", "round(value, 2) ge 0", json.Number("NaN"), rule.ErrInvalidDecimal},
		{"bounded text", "round(value, 2) ge 0", json.Number(strings.Repeat("1", 10001)), rule.ErrInvalidDecimal},
		{"bounded exponent", "round(value, 2) ge 0", json.Number("1e10001"), rule.ErrInvalidDecimal},
		{"conversion overflow", "round(value, 0) ge 0", json.Number("1e309"), rule.ErrFloatMathOverflow},
		{"scale overflow", "round(value, 309) ge 0", 1.2, rule.ErrFloatMathOverflow},
		{"product overflow", "round(value, 2) ge 0", math.MaxFloat64, rule.ErrFloatMathOverflow},
		{"literal overflow", "floor(1" + strings.Repeat("0", 309) + ") ge 0", nil, rule.ErrInvalidDecimal},
		{"datetime operator", "round(value, 2) dq 1", 1.2, rule.ErrInvalidMathOperand},
	}
	engine := rule.NewEngineWithOptions(rule.WithFloatMath())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := engine.Evaluate(tt.query, rule.D{"value": tt.value})
			if !errors.Is(err, tt.want) {
				t.Fatalf("Evaluate(%q) error = %v; want %v", tt.query, err, tt.want)
			}
		})
	}
}

func TestFloatMathSharedCompiledRule(t *testing.T) {
	t.Parallel()

	exact := rule.NewEngine()
	approximate := rule.NewEngineWithOptions(rule.WithFloatMath())

	compiled, err := exact.CompileRule("round(value, 2) eq 1.01")
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range []struct {
		name   string
		engine *rule.Engine
		want   bool
	}{{"exact", exact, true}, {"float", approximate, false}} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			for range 100 {
				got, evalErr := tt.engine.EvaluateCompiled(compiled, rule.D{"value": json.Number("1.005")})
				if evalErr != nil || got != tt.want {
					t.Fatalf("EvaluateCompiled = %v, %v; want %v, nil", got, evalErr, tt.want)
				}
			}
		})
	}
}

func TestFloatMathLenient(t *testing.T) {
	for _, lenient := range []bool{false, true} {
		options := []rule.Option{rule.WithFloatMath()}
		if lenient {
			options = append(options, rule.WithLenientMode())
		}

		got, err := rule.NewEngineWithOptions(options...).Evaluate("round(value, 2) eq 1", nil)
		if err != nil || got != lenient {
			t.Fatalf("lenient=%v: got %v, %v", lenient, got, err)
		}
	}
}

func TestFloatMathAllocations(t *testing.T) {
	engine := rule.NewEngineWithOptions(rule.WithFloatMath())

	compiled, err := engine.CompileRule("round(value, 2) eq 1.23 and floor(value) eq 1 and ceil(value) eq 2")
	if err != nil {
		t.Fatal(err)
	}

	for _, input := range []any{1.234, json.Number("1.234")} {
		context := rule.D{"value": input}

		allocations := testing.AllocsPerRun(100, func() {
			got, evalErr := engine.EvaluateCompiled(compiled, context)
			if evalErr != nil || !got {
				t.Fatalf("EvaluateCompiled = %v, %v", got, evalErr)
			}
		})
		if allocations != 0 {
			t.Errorf("input %T: %g allocations per evaluation, want zero", input, allocations)
		}
	}
}

func TestFloatMathErrorsInLists(t *testing.T) {
	engine := rule.NewEngineWithOptions(rule.WithFloatMath())
	for _, query := range []string{
		"items any (round(value, 2) gt 0)",
		"items all (round(value, 2) gt 0)",
		"items none (round(value, 2) gt 0)",
		"items where (round(value, 2) gt 0).length gt 0",
	} {
		_, err := engine.Evaluate(query, rule.D{"items": []any{rule.D{"value": math.MaxFloat64}}})
		if !errors.Is(err, rule.ErrFloatMathOverflow) {
			t.Errorf("Evaluate(%q) error = %v, want %v", query, err, rule.ErrFloatMathOverflow)
		}
	}
}

func BenchmarkMathModes(b *testing.B) {
	for _, operation := range []struct {
		name  string
		query string
	}{
		{"Round", "round(value, 2) eq 1.23"},
		{"Floor", "floor(value) eq 1"},
		{"Ceil", "ceil(value) eq 2"},
	} {
		for _, input := range []struct {
			name  string
			value any
		}{{"Native", 1.234}, {"JSON", json.Number("1.234")}} {
			for _, mode := range []struct {
				name   string
				engine *rule.Engine
			}{
				{"Exact", rule.NewEngine()},
				{"Float", rule.NewEngineWithOptions(rule.WithFloatMath())},
			} {
				b.Run(operation.name+"/"+input.name+"/"+mode.name, func(b *testing.B) {
					benchmarkMathMode(b, mode.engine, operation.query, input.value)
				})
			}
		}
	}
}

func benchmarkMathMode(b *testing.B, engine *rule.Engine, query string, value any) {
	b.Helper()

	compiled, err := engine.CompileRule(query)
	if err != nil {
		b.Fatal(err)
	}

	context := rule.D{"value": value}

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		got, evalErr := engine.EvaluateCompiled(compiled, context)
		if evalErr != nil || !got {
			b.Fatalf("EvaluateCompiled = %v, %v", got, evalErr)
		}
	}
}

func FuzzFloatMathBounds(f *testing.F) {
	for _, text := range []string{
		"1.005", "-1.005", "1e308", "1e-10000", "9007199254740993",
		"NaN", "0x10", "1/2", "1e10001",
	} {
		f.Add(text)
	}

	engine := rule.NewEngineWithOptions(rule.WithFloatMath())

	f.Fuzz(func(t *testing.T, text string) {
		result, err := engine.Evaluate(
			"floor(value) le value and value le ceil(value)",
			rule.D{"value": json.Number(text)},
		)
		if err == nil && !result {
			t.Fatalf("accepted value %q lies outside floor/ceil bounds", text)
		}
	})
}
