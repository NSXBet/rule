package test

import (
	"testing"
	"time"

	"github.com/NSXBet/rule"
	ruleslib "github.com/nikunjy/rules"
)

// CompatibilityTest defines a test case for comparing our library with nikunjy/rules
type CompatibilityTest struct {
	Name     string
	Rule     string
	Context  rule.D
	Expected bool
	Notes    string
}

// runCompatibilityTest runs a single compatibility test against both libraries
func runCompatibilityTest(t *testing.T, test CompatibilityTest) bool {
	t.Helper()

	// Test our library
	ourEngine := rule.NewEngine()
	ourResult, ourErr := ourEngine.Evaluate(test.Rule, test.Context)

	// Convert our context to map[string]interface{} for nikunjy/rules
	rulesContext := make(map[string]interface{})
	for k, v := range test.Context {
		rulesContext[k] = v
	}

	// Test nikunjy/rules library
	rulesResult, rulesErr := ruleslib.Evaluate(test.Rule, rulesContext)

	// Compare results
	success := true
	if ourErr != nil && rulesErr != nil {
		// Both errored - consider compatible if both fail
		t.Logf("✅ %s: Both libraries errored (compatible)", test.Name)
	} else if ourErr != nil || rulesErr != nil {
		// Only one errored
		t.Errorf("❌ %s: Error mismatch - Our: %v, Rules: %v", test.Name, ourErr, rulesErr)
		success = false
	} else if ourResult != rulesResult {
		// Different results
		t.Errorf("❌ %s: Result mismatch - Our: %v, Rules: %v", test.Name, ourResult, rulesResult)
		success = false
	} else {
		// Same results
		t.Logf("✅ %s: Compatible - both returned %v", test.Name, ourResult)
	}

	// Also check against expected value
	if ourErr == nil && ourResult != test.Expected {
		t.Errorf("⚠️  %s: Our result %v doesn't match expected %v", test.Name, ourResult, test.Expected)
		success = false
	}

	if test.Notes != "" {
		t.Logf("   Note: %s", test.Notes)
	}

	return success
}

func TestBasicCompatibility(t *testing.T) {
	tests := []CompatibilityTest{
		// Equality operators
		{
			Name:     "String equality",
			Rule:     `name eq "John"`,
			Context:  rule.D{"name": "John"},
			Expected: true,
			Notes:    "Basic string comparison",
		},
		{
			Name:     "Number equality",
			Rule:     `age eq 25`,
			Context:  rule.D{"age": 25},
			Expected: true,
			Notes:    "Basic numeric comparison",
		},
		{
			Name:     "Boolean equality",
			Rule:     `active eq true`,
			Context:  rule.D{"active": true},
			Expected: true,
			Notes:    "Basic boolean comparison",
		},

		// Cross-type comparisons (should be false in both libraries)
		{
			Name:     "String vs Number",
			Rule:     `value eq 42`,
			Context:  rule.D{"value": "42"},
			Expected: false,
			Notes:    "String '42' should not equal number 42",
		},
		{
			Name:     "Number vs String",
			Rule:     `value eq "42"`,
			Context:  rule.D{"value": 42},
			Expected: false,
			Notes:    "Number 42 should not equal string '42'",
		},
		{
			Name:     "Boolean vs Number",
			Rule:     `flag eq 1`,
			Context:  rule.D{"flag": true},
			Expected: false,
			Notes:    "Boolean true should not equal number 1",
		},

		// Numeric cross-type comparisons (int/float should work)
		{
			Name:     "Int vs Float equality",
			Rule:     `value eq 42.0`,
			Context:  rule.D{"value": 42},
			Expected: true,
			Notes:    "Integer 42 should equal float 42.0",
		},
		{
			Name:     "Float vs Int equality",
			Rule:     `value eq 42`,
			Context:  rule.D{"value": 42.0},
			Expected: true,
			Notes:    "Float 42.0 should equal integer 42",
		},

		// Relational operators
		{
			Name:     "Number less than",
			Rule:     `age lt 30`,
			Context:  rule.D{"age": 25},
			Expected: true,
			Notes:    "Numeric less than comparison",
		},
		{
			Name:     "String lexicographic comparison",
			Rule:     `name lt "Zoo"`,
			Context:  rule.D{"name": "Apple"},
			Expected: true,
			Notes:    "Lexicographic string comparison",
		},

		// String operators
		{
			Name:     "String contains",
			Rule:     `email co "@example"`,
			Context:  rule.D{"email": "user@example.com"},
			Expected: true,
			Notes:    "String contains operation",
		},
		{
			Name:     "String starts with",
			Rule:     `name sw "Mr"`,
			Context:  rule.D{"name": "Mr. Smith"},
			Expected: true,
			Notes:    "String starts with operation",
		},
		{
			Name:     "String ends with",
			Rule:     `file ew ".txt"`,
			Context:  rule.D{"file": "document.txt"},
			Expected: true,
			Notes:    "String ends with operation",
		},

		// Array membership
		{
			Name:     "String in array literal",
			Rule:     `role in ["admin", "user"]`,
			Context:  rule.D{"role": "admin"},
			Expected: true,
			Notes:    "String membership in array literal",
		},
		{
			Name:     "Number in array literal",
			Rule:     `score in [85, 90, 95]`,
			Context:  rule.D{"score": 90},
			Expected: true,
			Notes:    "Number membership in array literal",
		},

		// Logical operators
		{
			Name:     "Logical AND both true",
			Rule:     `age gt 18 and status eq "active"`,
			Context:  rule.D{"age": 25, "status": "active"},
			Expected: true,
			Notes:    "Logical AND with both conditions true",
		},
		{
			Name:     "Logical OR one true",
			Rule:     `role eq "admin" or role eq "moderator"`,
			Context:  rule.D{"role": "admin"},
			Expected: true,
			Notes:    "Logical OR with one condition true",
		},
		{
			Name:     "Logical NOT",
			Rule:     `not (age lt 18)`,
			Context:  rule.D{"age": 25},
			Expected: true,
			Notes:    "Logical NOT operation",
		},

		// Presence operator
		{
			Name:     "Property present",
			Rule:     `email pr`,
			Context:  rule.D{"email": "user@example.com"},
			Expected: true,
			Notes:    "Property presence check",
		},
		{
			Name:     "Property absent",
			Rule:     `phone pr`,
			Context:  rule.D{"email": "user@example.com"},
			Expected: false,
			Notes:    "Property absence check",
		},

		// Nested properties
		{
			Name:     "Nested property access",
			Rule:     `user.profile.age gt 18`,
			Context:  rule.D{"user": rule.D{"profile": rule.D{"age": 25}}},
			Expected: true,
			Notes:    "Nested object property access",
		},
		{
			Name:     "Deep nested property",
			Rule:     `app.settings.theme.mode eq "dark"`,
			Context:  rule.D{"app": rule.D{"settings": rule.D{"theme": rule.D{"mode": "dark"}}}},
			Expected: true,
			Notes:    "Deep nested property access",
		},
	}

	totalTests := len(tests)
	passedTests := 0

	t.Logf("Running %d compatibility tests...", totalTests)
	t.Logf("=======================================")

	for _, test := range tests {
		if runCompatibilityTest(t, test) {
			passedTests++
		}
	}

	t.Logf("\nCompatibility Summary:")
	t.Logf("=====================")
	t.Logf("Total tests: %d", totalTests)
	t.Logf("Passed: %d", passedTests)
	t.Logf("Failed: %d", totalTests-passedTests)
	t.Logf("Compatibility rate: %.1f%%", float64(passedTests)/float64(totalTests)*100)

	if passedTests != totalTests {
		t.Errorf("Not all compatibility tests passed")
	}
}

func TestTimeCompatibility(t *testing.T) {
	// Test time.Time compatibility specifically
	testTime := time.Date(2024, 7, 9, 22, 12, 0, 0, time.UTC)

	tests := []CompatibilityTest{
		{
			Name:     "time.Time string conversion",
			Rule:     `created_at eq "2024-07-09 22:12:00 +0000 UTC"`,
			Context:  rule.D{"created_at": testTime},
			Expected: true,
			Notes:    "time.Time should convert to string like nikunjy/rules",
		},
		{
			Name:     "time.Time string contains",
			Rule:     `created_at co "2024-07-09"`,
			Context:  rule.D{"created_at": testTime},
			Expected: true,
			Notes:    "String operations should work on time.Time string representation",
		},
		{
			Name:     "time.Time lexicographic comparison",
			Rule:     `created_at lt "2024-07-10"`,
			Context:  rule.D{"created_at": testTime},
			Expected: true,
			Notes:    "Lexicographic comparison should work on time.Time strings",
		},
	}

	for _, test := range tests {
		runCompatibilityTest(t, test)
	}
}

func TestIncompatibleFeatures(t *testing.T) {
	// Test features that are intentionally different
	ourEngine := rule.NewEngine()

	tests := []struct {
		name  string
		rule  string
		ctx   rule.D
		notes string
	}{
		{
			name:  "DateTime operators (our extension)",
			rule:  `created_at dq "2024-07-09T22:12:00Z"`,
			ctx:   rule.D{"created_at": "2024-07-09T22:12:00Z"},
			notes: "datetime operators are our extension, not in nikunjy/rules",
		},
		{
			name:  "DateTime before operator",
			rule:  `start_time be end_time`,
			ctx:   rule.D{"start_time": "2024-07-09T22:11:00Z", "end_time": "2024-07-09T22:12:00Z"},
			notes: "datetime operators are our extension",
		},
	}

	t.Log("Testing intentionally incompatible features (our extensions):")
	t.Log("===========================================================")

	for _, test := range tests {
		result, err := ourEngine.Evaluate(test.rule, test.ctx)
		if err != nil {
			t.Errorf("❌ %s failed: %v", test.name, err)
		} else {
			t.Logf("✅ %s: %v (%s)", test.name, result, test.notes)
		}
	}
}
