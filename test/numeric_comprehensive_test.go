package test

import (
	"math"
	"testing"

	"github.com/NSXBet/rule"
)

// TestAllNumericScenarios comprehensively tests ALL numeric comparison scenarios.
// Our library must NEVER fail ANY numeric comparison.
func TestAllNumericScenarios(t *testing.T) {
	t.Log("🔢 COMPREHENSIVE NUMERIC TESTING")
	t.Log("================================")

	engine := rule.NewEngine()

	tests := []struct {
		name      string
		context   rule.D
		rule      string
		expected  bool
		shouldErr bool
	}{
		// SMALL NUMBERS
		{
			name:     "Very small positive float",
			context:  rule.D{"value": 0.000000001},
			rule:     "value gt 0",
			expected: true,
		},
		{
			name:     "Very small negative float",
			context:  rule.D{"value": -0.000000001},
			rule:     "value lt 0",
			expected: true,
		},
		{
			name:     "Smallest normal float64",
			context:  rule.D{"value": math.SmallestNonzeroFloat64},
			rule:     "value gt 0",
			expected: true,
		},

		// LARGE NUMBERS
		{
			name:     "Max int64",
			context:  rule.D{"value": int64(math.MaxInt64)},
			rule:     "value gt 0",
			expected: true,
		},
		{
			name:     "Min int64",
			context:  rule.D{"value": int64(math.MinInt64)},
			rule:     "value lt 0",
			expected: true,
		},
		{
			name:     "Max float64",
			context:  rule.D{"value": math.MaxFloat64},
			rule:     "value gt 0",
			expected: true,
		},
		{
			name:     "Large integer precision edge",
			context:  rule.D{"value": int64(9007199254740992)}, // 2^53
			rule:     "value eq 9007199254740992",
			expected: true,
		},

		// ZERO VALUES
		{
			name:     "Zero int",
			context:  rule.D{"value": 0},
			rule:     "value eq 0",
			expected: true,
		},
		{
			name:     "Zero float",
			context:  rule.D{"value": 0.0},
			rule:     "value eq 0",
			expected: true,
		},
		{
			name:     "Negative zero",
			context:  rule.D{"value": math.Copysign(0, -1)},
			rule:     "value eq 0",
			expected: true,
		},

		// CROSS-TYPE COMPARISONS
		{
			name:     "Int vs float equality",
			context:  rule.D{"int_val": 42, "float_val": 42.0},
			rule:     "int_val eq float_val",
			expected: true,
		},
		{
			name:     "Int32 vs int64",
			context:  rule.D{"val32": int32(100), "val64": int64(100)},
			rule:     "val32 eq val64",
			expected: true,
		},
		{
			name:     "Uint vs int",
			context:  rule.D{"uint_val": uint(50), "int_val": int(50)},
			rule:     "uint_val eq int_val",
			expected: true,
		},
		{
			name:     "Float32 vs float64",
			context:  rule.D{"f32": float32(3.14), "f64": float64(3.14)},
			rule:     "f32 eq f64",
			expected: false, // Precision difference
		},

		// RELATIONAL OPERATIONS
		{
			name:     "Large number greater than",
			context:  rule.D{"big": int64(1000000000000), "small": int64(999999999999)},
			rule:     "big gt small",
			expected: true,
		},
		{
			name:     "Negative number comparisons",
			context:  rule.D{"neg1": -100, "neg2": -50},
			rule:     "neg1 lt neg2",
			expected: true,
		},
		{
			name:     "Mixed sign comparisons",
			context:  rule.D{"pos": 1, "neg": -1},
			rule:     "pos gt neg",
			expected: true,
		},

		// EDGE CASES
		{
			name:     "Infinity handling",
			context:  rule.D{"inf": math.Inf(1)},
			rule:     "inf gt 0",
			expected: true,
		},
		{
			name:     "Negative infinity",
			context:  rule.D{"neginf": math.Inf(-1)},
			rule:     "neginf lt 0",
			expected: true,
		},
		{
			name:     "NaN comparison should be false",
			context:  rule.D{"nan": math.NaN()},
			rule:     "nan eq nan",
			expected: false, // NaN != NaN by IEEE 754
		},

		// PRECISION TESTS
		{
			name:     "Very close floats",
			context:  rule.D{"a": 1.0000000000000001, "b": 1.0000000000000002},
			rule:     "a eq b",
			expected: false, // Should distinguish tiny differences
		},
		{
			name:     "Float precision boundary",
			context:  rule.D{"value": 0.1 + 0.2},
			rule:     "value eq 0.3",
			expected: true, // Go's == operator handles this correctly
		},

		// VARIOUS INTEGER SIZES
		{
			name:     "Int8 boundaries",
			context:  rule.D{"max8": int8(127), "min8": int8(-128)},
			rule:     "max8 gt min8",
			expected: true,
		},
		{
			name:     "Uint64 max",
			context:  rule.D{"max_uint": uint64(math.MaxUint64)},
			rule:     "max_uint gt 0",
			expected: true,
		},

		// PROPERTY-TO-PROPERTY NUMERIC COMPARISONS
		{
			name:     "Nested numeric property comparison",
			context:  rule.D{"user": rule.D{"age": 25}, "limit": rule.D{"min": 18}},
			rule:     "user.age gt limit.min",
			expected: true,
		},
		{
			name: "Deep nested numeric comparison",
			context: rule.D{
				"stats":      rule.D{"player": rule.D{"score": 1500}},
				"thresholds": rule.D{"expert": rule.D{"score": 1000}},
			},
			rule:     "stats.player.score gt thresholds.expert.score",
			expected: true,
		},
	}

	t.Logf("Running %d comprehensive numeric tests...", len(tests))
	t.Logf("=========================================")

	failedTests := 0

	for i, test := range tests {
		t.Logf("\n%d. %s", i+1, test.name)
		t.Logf("Rule: %s", test.rule)

		result, err := engine.Evaluate(test.rule, test.context)

		if test.shouldErr {
			if err == nil {
				t.Errorf("❌ Expected error but got none")

				failedTests++
			} else {
				t.Logf("✅ Expected error: %v", err)
			}
		} else {
			switch {
			case err != nil:
				t.Errorf("❌ CRITICAL: Numeric comparison failed with error: %v", err)
				t.Errorf("   Context: %+v", test.context)

				failedTests++
			case result != test.expected:
				t.Errorf("❌ CRITICAL: Wrong result - got %t, expected %t", result, test.expected)
				t.Errorf("   Context: %+v", test.context)

				failedTests++
			default:
				t.Logf("✅ Correct result: %t", result)
			}
		}
	}

	t.Logf("\n📊 Numeric Testing Summary:")
	t.Logf("==========================")
	t.Logf("Total tests: %d", len(tests))
	t.Logf("Passed: %d", len(tests)-failedTests)
	t.Logf("Failed: %d", failedTests)

	if failedTests > 0 {
		t.Errorf("CRITICAL: %d numeric tests failed. Our library must NEVER fail numeric comparisons!", failedTests)
	} else {
		t.Logf("🎉 ALL numeric scenarios passed! Library handles all numeric comparisons correctly.")
	}
}

// TestNumericOverflowScenarios tests edge cases around numeric overflow.
func TestNumericOverflowScenarios(t *testing.T) {
	t.Log("\n🔀 NUMERIC OVERFLOW SCENARIOS")
	t.Log("=============================")

	engine := rule.NewEngine()

	tests := []struct {
		name     string
		context  rule.D
		rule     string
		expected bool
	}{
		{
			name:     "Uint overflow converted to int64",
			context:  rule.D{"huge": uint64(math.MaxUint64)},
			rule:     "huge gt 0",
			expected: true,
		},
		{
			name:     "Large uint32 to int64",
			context:  rule.D{"val": uint32(math.MaxUint32)},
			rule:     "val eq 4294967295",
			expected: true,
		},
		{
			name:     "Int overflow boundary",
			context:  rule.D{"max_int": int(math.MaxInt), "max_int64": int64(math.MaxInt64)},
			rule:     "max_int eq max_int64",
			expected: true, // On 64-bit systems
		},
	}

	for i, test := range tests {
		t.Logf("\n%d. %s", i+1, test.name)
		t.Logf("Rule: %s", test.rule)

		result, err := engine.Evaluate(test.rule, test.context)
		if err != nil {
			t.Errorf("❌ CRITICAL: Overflow scenario failed: %v", err)
		} else {
			t.Logf("✅ Result: %t (expected: %t)", result, test.expected)

			if result != test.expected {
				t.Errorf("❌ Wrong result for overflow scenario")
			}
		}
	}
}
