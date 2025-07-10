package test

import (
	"testing"

	"github.com/NSXBet/rule"
	ruleslib "github.com/nikunjy/rules"
)

func TestAdditionalEdgeCases(t *testing.T) {
	t.Log("🔍 ADDITIONAL EDGE CASE TESTING")
	t.Log("===============================")

	// Additional edge cases to test
	edgeCases := []struct {
		name        string
		rule        string
		context     map[string]interface{}
		ourContext  rule.D
		description string
	}{
		{
			name:        "String case mixed",
			rule:        `name eq "JoHn"`,
			context:     map[string]interface{}{"name": "john"},
			ourContext:  rule.D{"name": "john"},
			description: "Mixed case string equality",
		},
		{
			name:        "Empty string equality",
			rule:        `text eq ""`,
			context:     map[string]interface{}{"text": ""},
			ourContext:  rule.D{"text": ""},
			description: "Empty string comparison",
		},
		{
			name:        "String with spaces",
			rule:        `text co " "`,
			context:     map[string]interface{}{"text": "hello world"},
			ourContext:  rule.D{"text": "hello world"},
			description: "Space character in contains",
		},
		{
			name:        "Number as string key",
			rule:        `"123" eq value`,
			context:     map[string]interface{}{"value": "123"},
			ourContext:  rule.D{"value": "123"},
			description: "Numeric string comparison",
		},
		{
			name:        "Boolean string representation",
			rule:        `flag co "true"`,
			context:     map[string]interface{}{"flag": "the value is true"},
			ourContext:  rule.D{"flag": "the value is true"},
			description: "Boolean keyword in string",
		},
		{
			name:        "Array membership case sensitive",
			rule:        `"ADMIN" in roles`,
			context:     map[string]interface{}{"roles": []interface{}{"admin", "user"}},
			ourContext:  rule.D{"roles": []any{"admin", "user"}},
			description: "Case sensitive array membership",
		},
		{
			name:        "Float vs int equality",
			rule:        `value eq 42`,
			context:     map[string]interface{}{"value": 42.0},
			ourContext:  rule.D{"value": 42.0},
			description: "Float to int comparison",
		},
		{
			name:        "Very small float",
			rule:        `value gt 0`,
			context:     map[string]interface{}{"value": 0.0000001},
			ourContext:  rule.D{"value": 0.0000001},
			description: "Very small positive float",
		},
	}

	for i, test := range edgeCases {
		t.Logf("\n%d. %s", i+1, test.name)
		t.Logf("Rule: %s", test.rule)
		t.Logf("Description: %s", test.description)

		// Test nikunjy/rules
		rulesResult, rulesErr := ruleslib.Evaluate(test.rule, test.context)

		// Test our library
		ourEngine := rule.NewEngine()
		ourResult, ourErr := ourEngine.Evaluate(test.rule, test.ourContext)

		t.Logf("nikunjy/rules: %v (err: %v)", rulesResult, rulesErr)
		t.Logf("Our library: %v (err: %v)", ourResult, ourErr)

		// Check compatibility
		if (rulesErr == nil && ourErr == nil && rulesResult == ourResult) ||
			(rulesErr != nil && ourErr != nil) {
			t.Logf("✅ COMPATIBLE")
		} else {
			t.Logf("❌ INCOMPATIBLE")
		}
	}

	t.Log("\n🎯 Additional edge case testing completed!")
}
