package test

import (
	"testing"
	"time"

	"github.com/NSXBet/rule"
	ruleslib "github.com/nikunjy/rules"
)

// DeepCompatibilityTest for finding actual incompatibilities.
type DeepCompatibilityTest struct {
	Name        string
	Rule        string
	Context     map[string]interface{}
	OurContext  rule.D
	Description string
	ExpectDiff  bool   // true if we expect different results
	Reason      string // why we expect different results
}

func TestDeepCompatibility(t *testing.T) {
	tests := []DeepCompatibilityTest{
		// TIME.TIME EDGE CASES
		{
			Name: "time.Time lt RFC3339 string",
			Rule: `timestamp lt "2024-07-09T22:12:00Z"`,
			Context: map[string]interface{}{
				"timestamp": time.Date(2024, 7, 9, 22, 11, 59, 0, time.UTC),
			},
			OurContext: rule.D{
				"timestamp": time.Date(2024, 7, 9, 22, 11, 59, 0, time.UTC),
			},
			Description: "time.Time lexicographic comparison vs RFC3339",
			ExpectDiff:  false,
			Reason:      "Should be same string comparison in both",
		},
		{
			Name: "time.Time with nanoseconds",
			Rule: `timestamp lt "2024-07-09T22:12:00Z"`,
			Context: map[string]interface{}{
				"timestamp": time.Date(2024, 7, 9, 22, 11, 59, 123456789, time.UTC),
			},
			OurContext: rule.D{
				"timestamp": time.Date(2024, 7, 9, 22, 11, 59, 123456789, time.UTC),
			},
			Description: "time.Time with nanoseconds vs RFC3339",
			ExpectDiff:  false,
			Reason:      "Both should convert to same string format",
		},
		{
			Name: "time.Time eq exact RFC3339",
			Rule: `timestamp eq "2024-07-09 22:12:00 +0000 UTC"`,
			Context: map[string]interface{}{
				"timestamp": time.Date(2024, 7, 9, 22, 12, 0, 0, time.UTC),
			},
			OurContext: rule.D{
				"timestamp": time.Date(2024, 7, 9, 22, 12, 0, 0, time.UTC),
			},
			Description: "time.Time exact string match",
			ExpectDiff:  false,
			Reason:      "Should match time.String() format exactly",
		},
		{
			Name: "time.Time with timezone",
			Rule: `timestamp co "+0000"`,
			Context: map[string]interface{}{
				"timestamp": time.Date(2024, 7, 9, 22, 12, 0, 0, time.UTC),
			},
			OurContext: rule.D{
				"timestamp": time.Date(2024, 7, 9, 22, 12, 0, 0, time.UTC),
			},
			Description: "time.Time timezone in string",
			ExpectDiff:  false,
			Reason:      "Both should have timezone in string representation",
		},

		// ARRAY EDGE CASES
		{
			Name: "Empty array in context",
			Rule: `"test" in items`,
			Context: map[string]interface{}{
				"items": []interface{}{},
			},
			OurContext: rule.D{
				"items": []any{},
			},
			Description: "Empty array membership",
			ExpectDiff:  true,
			Reason:      "nikunjy/rules panics on array operations",
		},
		{
			Name: "Mixed type array",
			Rule: `42 in items`,
			Context: map[string]interface{}{
				"items": []interface{}{1, "2", 3.0, true, 42},
			},
			OurContext: rule.D{
				"items": []any{1, "2", 3.0, true, 42},
			},
			Description: "Mixed type array membership",
			ExpectDiff:  true,
			Reason:      "nikunjy/rules panics on array operations",
		},
		{
			Name: "Array with nil values",
			Rule: `null in items`,
			Context: map[string]interface{}{
				"items": []interface{}{1, nil, "test"},
			},
			OurContext: rule.D{
				"items": []any{1, nil, "test"},
			},
			Description: "Array with nil values",
			ExpectDiff:  true,
			Reason:      "null literal vs nil handling might differ",
		},

		// NUMERIC PRECISION EDGE CASES
		{
			Name: "Large integer precision",
			Rule: `big_number eq 9007199254740992`, // 2^53
			Context: map[string]interface{}{
				"big_number": int64(9007199254740992),
			},
			OurContext: rule.D{
				"big_number": int64(9007199254740992),
			},
			Description: "Large integer at float64 precision limit",
			ExpectDiff:  false,
			Reason:      "Should handle large integers correctly",
		},
		{
			Name: "Very large integer",
			Rule: `big_number eq 9223372036854775807`, // max int64
			Context: map[string]interface{}{
				"big_number": int64(9223372036854775807),
			},
			OurContext: rule.D{
				"big_number": int64(9223372036854775807),
			},
			Description: "Max int64 value",
			ExpectDiff:  false,
			Reason:      "Both handle large integers correctly",
		},
		{
			Name: "Float precision edge case",
			Rule: `value eq 0.1`,
			Context: map[string]interface{}{
				"value": float32(0.1),
			},
			OurContext: rule.D{
				"value": float32(0.1),
			},
			Description: "Float32 vs float64 precision",
			ExpectDiff:  false,
			Reason:      "Both reject imprecise float comparisons",
		},

		// STRING CASE SENSITIVITY
		{
			Name: "String case sensitivity eq",
			Rule: `name eq "JOHN"`,
			Context: map[string]interface{}{
				"name": "john",
			},
			OurContext: rule.D{
				"name": "john",
			},
			Description: "Case sensitive string equality",
			ExpectDiff:  false,
			Reason:      "Both libraries are case insensitive",
		},
		{
			Name: "String case sensitivity co",
			Rule: `name co "OHN"`,
			Context: map[string]interface{}{
				"name": "john",
			},
			OurContext: rule.D{
				"name": "john",
			},
			Description: "Case sensitive string contains",
			ExpectDiff:  false,
			Reason:      "Both libraries are case insensitive",
		},

		// ERROR HANDLING DIFFERENCES
		{
			Name: "Invalid nested access",
			Rule: `user.profile.age gt 18`,
			Context: map[string]interface{}{
				"user": "not_an_object",
			},
			OurContext: rule.D{
				"user": "not_an_object",
			},
			Description: "Invalid nested property access",
			ExpectDiff:  true,
			Reason:      "Error handling might differ between libraries",
		},
		{
			Name: "Missing property comparison",
			Rule: `nonexistent eq "test"`,
			Context: map[string]interface{}{
				"other": "value",
			},
			OurContext: rule.D{
				"other": "value",
			},
			Description: "Missing property in comparison",
			ExpectDiff:  false,
			Reason:      "Both should handle missing properties similarly",
		},

		// TYPE COERCION EDGE CASES
		{
			Name: "String number eq",
			Rule: `value eq "42"`,
			Context: map[string]interface{}{
				"value": 42,
			},
			OurContext: rule.D{
				"value": 42,
			},
			Description: "Number to string coercion",
			ExpectDiff:  false,
			Reason:      "Both should reject cross-type comparison",
		},
		{
			Name: "Boolean to string",
			Rule: `flag eq "true"`,
			Context: map[string]interface{}{
				"flag": true,
			},
			OurContext: rule.D{
				"flag": true,
			},
			Description: "Boolean to string coercion",
			ExpectDiff:  false,
			Reason:      "Both should reject cross-type comparison",
		},
		{
			Name: "Zero vs empty string",
			Rule: `value eq ""`,
			Context: map[string]interface{}{
				"value": 0,
			},
			OurContext: rule.D{
				"value": 0,
			},
			Description: "Zero number vs empty string",
			ExpectDiff:  false,
			Reason:      "Both should reject cross-type comparison",
		},

		// UNICODE AND SPECIAL CHARACTERS
		{
			Name: "Unicode string comparison",
			Rule: `name eq "José"`,
			Context: map[string]interface{}{
				"name": "José",
			},
			OurContext: rule.D{
				"name": "José",
			},
			Description: "Unicode string handling",
			ExpectDiff:  false,
			Reason:      "Both should handle Unicode correctly",
		},
		{
			Name: "Special characters in strings",
			Rule: `text co "\n"`,
			Context: map[string]interface{}{
				"text": "line1\nline2",
			},
			OurContext: rule.D{
				"text": "line1\nline2",
			},
			Description: "Special characters in string operations",
			ExpectDiff:  true,
			Reason:      "Different handling of special characters in contains",
		},

		// BOUNDARY VALUES
		{
			Name: "Empty string operations",
			Rule: `text co ""`,
			Context: map[string]interface{}{
				"text": "hello",
			},
			OurContext: rule.D{
				"text": "hello",
			},
			Description: "Empty string contains",
			ExpectDiff:  false,
			Reason:      "Both should handle empty string contains",
		},
		{
			Name: "Zero in numeric operations",
			Rule: `value gt 0`,
			Context: map[string]interface{}{
				"value": 0.0,
			},
			OurContext: rule.D{
				"value": 0.0,
			},
			Description: "Zero in greater than comparison",
			ExpectDiff:  false,
			Reason:      "Both should handle zero correctly",
		},

		// COMPLEX NESTED SCENARIOS
		{
			Name: "Deep nested with missing middle",
			Rule: `a.b.c.d eq "test"`,
			Context: map[string]interface{}{
				"a": map[string]interface{}{
					"x": "wrong_path",
				},
			},
			OurContext: rule.D{
				"a": rule.D{
					"x": "wrong_path",
				},
			},
			Description: "Deep nested access with missing intermediate",
			ExpectDiff:  false,
			Reason:      "Both return false for missing nested properties",
		},
	}

	totalTests := len(tests)
	compatibleTests := 0
	expectedDifferences := 0
	unexpectedDifferences := 0

	t.Logf("Running %d deep compatibility tests...", totalTests)
	t.Logf("==========================================")

	for _, test := range tests {
		t.Logf("\n🔍 %s", test.Name)
		t.Logf("Rule: %s", test.Rule)
		t.Logf("Description: %s", test.Description)

		// Test nikunjy/rules
		rulesResult, rulesErr := ruleslib.Evaluate(test.Rule, test.Context)

		// Test our library
		ourEngine := rule.NewEngine()
		ourResult, ourErr := ourEngine.Evaluate(test.Rule, test.OurContext)

		// Analyze results
		errorsMatch := (rulesErr == nil) == (ourErr == nil)
		resultsMatch := rulesResult == ourResult
		overallMatch := errorsMatch && (rulesErr != nil || resultsMatch)

		if overallMatch {
			compatibleTests++

			if test.ExpectDiff {
				t.Errorf("❌ UNEXPECTED COMPATIBILITY: %s", test.Name)
				t.Errorf("   Expected difference because: %s", test.Reason)
				t.Errorf("   But both returned: %v (err: %v)", rulesResult, rulesErr)

				unexpectedDifferences++
			} else {
				t.Logf("✅ COMPATIBLE: both returned %v", rulesResult)
			}
		} else {
			if test.ExpectDiff {
				t.Logf("✅ EXPECTED DIFFERENCE: %s", test.Reason)

				expectedDifferences++
			} else {
				t.Errorf("❌ UNEXPECTED INCOMPATIBILITY: %s", test.Name)

				unexpectedDifferences++
			}

			t.Logf("   nikunjy/rules: %v (err: %v)", rulesResult, rulesErr)
			t.Logf("   Our library: %v (err: %v)", ourResult, ourErr)
		}
	}

	t.Logf("\n📊 Deep Compatibility Analysis:")
	t.Logf("===============================")
	t.Logf("Total tests: %d", totalTests)
	t.Logf("Compatible: %d", compatibleTests)
	t.Logf("Expected differences: %d", expectedDifferences)
	t.Logf("Unexpected differences: %d", unexpectedDifferences)
	t.Logf("Actual compatibility rate: %.1f%%",
		float64(compatibleTests)/float64(totalTests)*100)

	if unexpectedDifferences > 0 {
		t.Errorf("Found %d unexpected incompatibilities", unexpectedDifferences)
	}
}

// TestStringCaseSensitivity specifically tests nikunjy/rules case handling.
func TestStringCaseSensitivity(t *testing.T) {
	t.Logf("🔤 Testing String Case Sensitivity")
	t.Logf("=================================")

	context := map[string]interface{}{
		"name": "john",
	}
	ourContext := rule.D{
		"name": "john",
	}

	caseTests := []struct {
		rule string
		desc string
	}{
		{`name eq "JOHN"`, "Uppercase equality"},
		{`name eq "John"`, "Title case equality"},
		{`name co "OHN"`, "Uppercase contains"},
		{`name sw "JO"`, "Uppercase starts with"},
		{`name ew "HN"`, "Uppercase ends with"},
	}

	for _, test := range caseTests {
		rulesResult, rulesErr := ruleslib.Evaluate(test.rule, context)

		ourEngine := rule.NewEngine()
		ourResult, ourErr := ourEngine.Evaluate(test.rule, ourContext)

		t.Logf("\n%s: %s", test.desc, test.rule)
		t.Logf("nikunjy/rules: %v (err: %v)", rulesResult, rulesErr)
		t.Logf("Our library: %v (err: %v)", ourResult, ourErr)

		if rulesResult != ourResult || (rulesErr == nil) != (ourErr == nil) {
			t.Logf("❌ INCOMPATIBLE")
		} else {
			t.Logf("✅ Compatible")
		}
	}
}
