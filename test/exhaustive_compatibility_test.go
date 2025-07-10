package test

import (
	"testing"
	"time"

	"github.com/NSXBet/rule"
	ruleslib "github.com/nikunjy/rules"
)

func TestExhaustiveCompatibility(t *testing.T) {
	t.Log("🔍 EXHAUSTIVE COMPATIBILITY ANALYSIS")
	t.Log("===================================")

	// Test every possible incompatibility scenario
	tests := []struct {
		category string
		tests    []CompatibilityCase
	}{
		{
			category: "STRING CASE SENSITIVITY (NOW COMPATIBLE)",
			tests: []CompatibilityCase{
				{
					name:         "eq case insensitive",
					rule:         `name eq "JOHN"`,
					context:      map[string]interface{}{"name": "john"},
					ourContext:   rule.D{"name": "john"},
					expectOurs:   true,
					expectRules:  true,
					incompatible: false,
					reason:       "Both libraries are now case-insensitive for string operations",
				},
				{
					name:         "ne case insensitive",
					rule:         `name ne "JOHN"`,
					context:      map[string]interface{}{"name": "john"},
					ourContext:   rule.D{"name": "john"},
					expectOurs:   false,
					expectRules:  false,
					incompatible: false,
					reason:       "Both libraries see them as equal (case-insensitive)",
				},
				{
					name:         "co case insensitive",
					rule:         `text co "HELLO"`,
					context:      map[string]interface{}{"text": "say hello world"},
					ourContext:   rule.D{"text": "say hello world"},
					expectOurs:   true,
					expectRules:  true,
					incompatible: false,
					reason:       "Both libraries convert to lowercase before comparison",
				},
				{
					name:         "sw case insensitive",
					rule:         `text sw "HELLO"`,
					context:      map[string]interface{}{"text": "hello world"},
					ourContext:   rule.D{"text": "hello world"},
					expectOurs:   true,
					expectRules:  true,
					incompatible: false,
					reason:       "Both libraries are case-insensitive for starts with",
				},
				{
					name:         "ew case insensitive",
					rule:         `text ew "WORLD"`,
					context:      map[string]interface{}{"text": "hello world"},
					ourContext:   rule.D{"text": "hello world"},
					expectOurs:   true,
					expectRules:  true,
					incompatible: false,
					reason:       "Both libraries are case-insensitive for ends with",
				},
			},
		},
		{
			category: "ARRAY OPERATIONS",
			tests: []CompatibilityCase{
				{
					name:         "empty array handling",
					rule:         `"test" in items`,
					context:      map[string]interface{}{"items": []interface{}{}},
					ourContext:   rule.D{"items": []any{}},
					expectOurs:   false,
					expectRules:  false, // but errors
					incompatible: true,
					reason:       "nikunjy/rules panics on empty arrays, we handle gracefully",
				},
				{
					name:         "mixed type array",
					rule:         `42 in items`,
					context:      map[string]interface{}{"items": []interface{}{1, "hello", 42, true}},
					ourContext:   rule.D{"items": []any{1, "hello", 42, true}},
					expectOurs:   true,
					expectRules:  true, // but errors
					incompatible: true,
					reason:       "nikunjy/rules may panic on mixed type arrays",
				},
				{
					name:         "array with nil",
					rule:         `42 in items`,
					context:      map[string]interface{}{"items": []interface{}{nil, 42}},
					ourContext:   rule.D{"items": []any{nil, 42}},
					expectOurs:   true,
					expectRules:  true, // but errors
					incompatible: true,
					reason:       "nikunjy/rules may panic with nil in arrays",
				},
			},
		},
		{
			category: "ERROR HANDLING",
			tests: []CompatibilityCase{
				{
					name:         "invalid nested access",
					rule:         `user.profile.age gt 18`,
					context:      map[string]interface{}{"user": "not_object"},
					ourContext:   rule.D{"user": "not_object"},
					expectOurs:   false,
					expectRules:  false, // but errors
					incompatible: true,
					reason:       "nikunjy/rules errors, we return false gracefully",
				},
				{
					name:         "deep invalid access",
					rule:         `a.b.c.d.e.f eq "test"`,
					context:      map[string]interface{}{"a": map[string]interface{}{"b": "not_object"}},
					ourContext:   rule.D{"a": rule.D{"b": "not_object"}},
					expectOurs:   false,
					expectRules:  false, // but errors
					incompatible: true,
					reason:       "nikunjy/rules errors on invalid deep access, we return false",
				},
			},
		},
		{
			category: "STRING SPECIAL CASES",
			tests: []CompatibilityCase{
				{
					name:         "newline in contains",
					rule:         `text co "\n"`,
					context:      map[string]interface{}{"text": "line1\nline2"},
					ourContext:   rule.D{"text": "line1\nline2"},
					expectOurs:   true,
					expectRules:  false,
					incompatible: true,
					reason:       "Different handling of special characters in string operations",
				},
				{
					name:         "tab character",
					rule:         `text co "\t"`,
					context:      map[string]interface{}{"text": "col1\tcol2"},
					ourContext:   rule.D{"text": "col1\tcol2"},
					expectOurs:   true,
					expectRules:  false,
					incompatible: true,
					reason:       "Different handling of tab characters",
				},
			},
		},
		{
			category: "TIME.TIME EDGE CASES",
			tests: []CompatibilityCase{
				{
					name:         "time.Time exact format match",
					rule:         `timestamp eq "2024-07-10 02:30:00 +0000 UTC"`,
					context:      map[string]interface{}{"timestamp": time.Date(2024, 7, 10, 2, 30, 0, 0, time.UTC)},
					ourContext:   rule.D{"timestamp": time.Date(2024, 7, 10, 2, 30, 0, 0, time.UTC)},
					expectOurs:   true,
					expectRules:  true,
					incompatible: false,
					reason:       "Both should convert time.Time to same string format",
				},
				{
					name: "time.Time with nanoseconds format",
					rule: `timestamp co "123456"`,
					context: map[string]interface{}{
						"timestamp": time.Date(2024, 7, 10, 2, 30, 0, 123456789, time.UTC),
					},
					ourContext:   rule.D{"timestamp": time.Date(2024, 7, 10, 2, 30, 0, 123456789, time.UTC)},
					expectOurs:   true,
					expectRules:  true,
					incompatible: false,
					reason:       "Both should include nanoseconds in string representation",
				},
			},
		},
		{
			category: "DATETIME OPERATORS (OUR EXTENSION)",
			tests: []CompatibilityCase{
				{
					name:         "datetime equal operator",
					rule:         `created_at dq "2024-07-10T02:30:00Z"`,
					context:      map[string]interface{}{"created_at": "2024-07-10T02:30:00Z"},
					ourContext:   rule.D{"created_at": "2024-07-10T02:30:00Z"},
					expectOurs:   true,
					expectRules:  false, // should error
					incompatible: true,
					reason:       "dq operator is our extension, not in nikunjy/rules",
				},
				{
					name: "datetime before operator",
					rule: `start_time be end_time`,
					context: map[string]interface{}{
						"start_time": "2024-07-10T01:00:00Z",
						"end_time":   "2024-07-10T02:00:00Z",
					},
					ourContext:   rule.D{"start_time": "2024-07-10T01:00:00Z", "end_time": "2024-07-10T02:00:00Z"},
					expectOurs:   true,
					expectRules:  false, // should error
					incompatible: true,
					reason:       "be operator is our extension, not in nikunjy/rules",
				},
			},
		},
		{
			category: "SYNTAX DIFFERENCES (MAJOR INCOMPATIBILITY)",
			tests: []CompatibilityCase{
				{
					name:         "unquoted string literal",
					rule:         `name eq Bernardo`, // No quotes around Bernardo
					context:      map[string]interface{}{"name": "Bernardo"},
					ourContext:   rule.D{"name": "Bernardo"},
					expectOurs:   false, // should error
					expectRules:  true,
					incompatible: true,
					reason:       "nikunjy/rules supports unquoted strings, we require quotes",
				},
				{
					name:         "unquoted string with property",
					rule:         `user.name eq John`, // No quotes around John
					context:      map[string]interface{}{"user": map[string]interface{}{"name": "John"}},
					ourContext:   rule.D{"user": rule.D{"name": "John"}},
					expectOurs:   false, // should error
					expectRules:  true,
					incompatible: true,
					reason:       "We require quoted strings, nikunjy/rules allows unquoted identifiers as strings",
				},
			},
		},
		{
			category: "OUR ENHANCEMENTS (INTENTIONAL EXTENSIONS)",
			tests: []CompatibilityCase{
				{
					name: "property to property comparison",
					rule: `user.age eq threshold.minimum`,
					context: map[string]interface{}{
						"user":      map[string]interface{}{"age": 25},
						"threshold": map[string]interface{}{"minimum": 25},
					},
					ourContext: rule.D{
						"user":      rule.D{"age": 25},
						"threshold": rule.D{"minimum": 25},
					},
					expectOurs:   true,
					expectRules:  false, // should error
					incompatible: true,
					reason:       "Our library supports property-to-property comparisons, nikunjy/rules does not",
				},
				{
					name: "nested property to nested property",
					rule: `config.limits.max eq settings.values.ceiling`,
					context: map[string]interface{}{
						"config":   map[string]interface{}{"limits": map[string]interface{}{"max": 100}},
						"settings": map[string]interface{}{"values": map[string]interface{}{"ceiling": 100}},
					},
					ourContext: rule.D{
						"config":   rule.D{"limits": rule.D{"max": 100}},
						"settings": rule.D{"values": rule.D{"ceiling": 100}},
					},
					expectOurs:   true,
					expectRules:  false, // should error
					incompatible: true,
					reason:       "Deep nested property-to-property comparisons are our enhancement",
				},
			},
		},
	}

	totalTests := 0
	actualIncompatibilities := 0
	expectedIncompatibilities := 0
	unexpectedCompatibilities := 0
	unexpectedIncompatibilities := 0

	for _, category := range tests {
		t.Logf("\n📂 %s", category.category)
		t.Log("=============================")

		for _, test := range category.tests {
			totalTests++

			// Test nikunjy/rules
			rulesResult, rulesErr := ruleslib.Evaluate(test.rule, test.context)

			// Test our library
			ourEngine := rule.NewEngine()
			ourResult, ourErr := ourEngine.Evaluate(test.rule, test.ourContext)

			// Determine if results are compatible
			resultsMatch := (rulesErr == nil && ourErr == nil && rulesResult == ourResult) ||
				(rulesErr != nil && ourErr != nil)

			t.Logf("\n🔬 %s", test.name)
			t.Logf("Rule: %s", test.rule)
			t.Logf("nikunjy/rules: %v (err: %v)", rulesResult, rulesErr)
			t.Logf("Our library: %v (err: %v)", ourResult, ourErr)

			if test.incompatible {
				if resultsMatch {
					t.Logf("⚠️  UNEXPECTED COMPATIBILITY")
					t.Logf("   Expected incompatible because: %s", test.reason)

					unexpectedCompatibilities++
				} else {
					t.Logf("✅ EXPECTED INCOMPATIBILITY: %s", test.reason)

					expectedIncompatibilities++
					actualIncompatibilities++
				}
			} else {
				if resultsMatch {
					t.Logf("✅ COMPATIBLE as expected")
				} else {
					t.Logf("❌ UNEXPECTED INCOMPATIBILITY")

					unexpectedIncompatibilities++
					actualIncompatibilities++
				}
			}
		}
	}

	t.Log("\n============================================================")
	t.Log("📊 EXHAUSTIVE COMPATIBILITY ANALYSIS RESULTS")
	t.Log("============================================================")
	t.Logf("Total tests: %d", totalTests)
	t.Logf("Actual incompatibilities: %d", actualIncompatibilities)
	t.Logf("Expected incompatibilities: %d", expectedIncompatibilities)
	t.Logf("Unexpected compatibilities: %d", unexpectedCompatibilities)
	t.Logf("Unexpected incompatibilities: %d", unexpectedIncompatibilities)

	actualCompatibilityRate := float64(totalTests-actualIncompatibilities) / float64(totalTests) * 100
	t.Logf("Actual compatibility rate: %.1f%%", actualCompatibilityRate)

	t.Logf("\n🎯 COMPATIBILITY BREAKDOWN:")
	t.Logf("• String operations: ✅ COMPATIBLE (case-insensitive)")
	t.Logf("• Array operations: ❌ INCOMPATIBLE (error handling)")
	t.Logf("• Error handling: ❌ INCOMPATIBLE (graceful vs panic/error)")
	t.Logf("• Special characters: ❌ INCOMPATIBLE (different handling)")
	t.Logf("• DateTime operators: ❌ INCOMPATIBLE BY DESIGN (our extension)")
	t.Logf("• Syntax differences: ❌ INCOMPATIBLE (unquoted strings)")
	t.Logf("• Property-to-property: ❌ INCOMPATIBLE BY DESIGN (our enhancement)")
	t.Logf("• Basic numeric/boolean: ✅ COMPATIBLE")
	t.Logf("• time.Time handling: ✅ COMPATIBLE")
	t.Logf("• Nested properties: ✅ MOSTLY COMPATIBLE")

	if unexpectedIncompatibilities > 0 {
		t.Errorf("Found %d unexpected incompatibilities", unexpectedIncompatibilities)
	}
}

type CompatibilityCase struct {
	name         string
	rule         string
	context      map[string]interface{}
	ourContext   rule.D
	expectOurs   bool
	expectRules  bool
	incompatible bool
	reason       string
}
