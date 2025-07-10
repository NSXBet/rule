package test

import (
	"testing"

	"github.com/NSXBet/rule"
)

// TestPropertyToPropertyComparisons tests all 5 documented property-to-property comparison scenarios.
func TestPropertyToPropertyComparisons(t *testing.T) {
	t.Log("🔗 PROPERTY-TO-PROPERTY COMPARISON TESTS")
	t.Log("========================================")

	engine := rule.NewEngine()

	// Test context matching the README documentation examples
	context := rule.D{
		// Top-level properties
		"age":     25,
		"min_age": 18,
		"score":   1200,
		"minimum": 21,

		// Single-level nested
		"user": rule.D{
			"age": 30,
		},
		"limits": rule.D{
			"maximum": 1500,
			"minimum": 18,
		},

		// Multi-level nested (2 levels)
		"config": rule.D{
			"settings": rule.D{
				"max": 2000,
			},
		},

		// Deep nested (4 levels)
		"system": rule.D{
			"validation": rule.D{
				"rules": rule.D{
					"limits": rule.D{
						"ceiling": 2000,
					},
				},
			},
		},
	}

	tests := []struct {
		name        string
		rule        string
		expected    bool
		description string
	}{
		// 1. Simple to Simple: Compare two top-level properties
		{
			name:        "Simple to Simple - age > min_age",
			rule:        "age gt min_age",
			expected:    true,
			description: "Compare two top-level properties (25 > 18)",
		},
		{
			name:        "Simple to Simple - age == min_age (false)",
			rule:        "age eq min_age",
			expected:    false,
			description: "Compare two top-level properties for equality (25 != 18)",
		},
		{
			name:        "Simple to Simple - score >= minimum",
			rule:        "score ge minimum",
			expected:    true,
			description: "Compare two top-level properties (1200 >= 21)",
		},

		// 2. Nested to Simple: Compare nested property to top-level property
		{
			name:        "Nested to Simple - user.age > minimum",
			rule:        "user.age gt minimum",
			expected:    true,
			description: "Compare 1-level nested to top-level (30 > 21)",
		},
		{
			name:        "Nested to Simple - user.age <= minimum (false)",
			rule:        "user.age le minimum",
			expected:    false,
			description: "Compare 1-level nested to top-level (30 <= 21 is false)",
		},

		// 3. Simple to Nested: Compare top-level property to nested property
		{
			name:        "Simple to Nested - score <= limits.maximum",
			rule:        "score le limits.maximum",
			expected:    true,
			description: "Compare top-level to 1-level nested (1200 <= 1500)",
		},
		{
			name:        "Simple to Nested - minimum == limits.minimum",
			rule:        "minimum gt limits.minimum",
			expected:    true,
			description: "Compare top-level to 1-level nested (21 > 18)",
		},

		// 4. Nested to Nested: Compare two single-level nested properties
		{
			name:        "Nested to Nested - user.age > limits.minimum",
			rule:        "user.age gt limits.minimum",
			expected:    true,
			description: "Compare two 1-level nested properties (30 > 18)",
		},
		{
			name:        "Nested to Nested - user.age < limits.maximum",
			rule:        "user.age lt limits.maximum",
			expected:    true,
			description: "Compare two 1-level nested properties (30 < 1500)",
		},

		// 5. Deep Nested: Compare 2-level vs 4-level nested properties
		{
			name:        "Deep Nested - config.settings.max == system deep nested",
			rule:        "config.settings.max eq system.validation.rules.limits.ceiling",
			expected:    true,
			description: "Compare 2-level vs 4-level nested properties (2000 == 2000)",
		},
		{
			name:        "Deep Nested - config.settings.max >= system deep nested",
			rule:        "config.settings.max ge system.validation.rules.limits.ceiling",
			expected:    true,
			description: "Compare 2-level vs 4-level nested properties (2000 >= 2000)",
		},
	}

	for i, test := range tests {
		t.Logf("\n%d. %s", i+1, test.name)
		t.Logf("   Rule: %s", test.rule)
		t.Logf("   Description: %s", test.description)

		result, err := engine.Evaluate(test.rule, context)
		if err != nil {
			t.Errorf("❌ Error evaluating rule: %v", err)
			continue
		}

		if result != test.expected {
			t.Errorf("❌ Expected %t, got %t", test.expected, result)
		} else {
			t.Logf("✅ Correct result: %t", result)
		}
	}
}

// TestPropertyToPropertyStringComparisons tests string operations in property comparisons.
func TestPropertyToPropertyStringComparisons(t *testing.T) {
	t.Log("\n📝 PROPERTY-TO-PROPERTY STRING COMPARISONS")
	t.Log("==========================================")

	engine := rule.NewEngine()

	context := rule.D{
		// String properties for comparison
		"username": "john_doe",
		"prefix":   "john",
		"suffix":   "_doe",
		"domain":   "example.com",

		"user": rule.D{
			"name":     "john_doe",
			"email":    "john@example.com",
			"category": "premium",
		},
		"profile": rule.D{
			"username": "john_doe",
			"level":    "premium",
		},
		"settings": rule.D{
			"theme": rule.D{
				"name": "dark",
			},
		},
		"defaults": rule.D{
			"ui": rule.D{
				"theme": rule.D{
					"preference": "dark",
				},
			},
		},
	}

	tests := []struct {
		name     string
		rule     string
		expected bool
	}{
		// String equality comparisons
		{"String eq - simple to simple", "username eq user.name", true},
		{"String eq - nested to nested", "user.category eq profile.level", true},
		{"String eq - 2-level vs 3-level nested", "settings.theme.name eq defaults.ui.theme.preference", true},

		// String contains operations
		{"String co - simple to simple", "username co prefix", true},
		{"String co - nested to simple", "user.email co domain", true},

		// String starts with / ends with
		{"String sw - simple to simple", "username sw prefix", true},
		{"String ew - simple to simple", "username ew suffix", true},

		// String inequality
		{"String ne - different strings", "username ne domain", true},
		{"String ne - same strings", "username ne user.name", false},
	}

	for i, test := range tests {
		t.Logf("\n%d. %s", i+1, test.name)
		t.Logf("   Rule: %s", test.rule)

		result, err := engine.Evaluate(test.rule, context)
		if err != nil {
			t.Errorf("❌ Error evaluating rule: %v", err)
			continue
		}

		if result != test.expected {
			t.Errorf("❌ Expected %t, got %t", test.expected, result)
		} else {
			t.Logf("✅ Correct result: %t", result)
		}
	}
}

// TestPropertyToPropertyDateTimeComparisons tests datetime operations in property comparisons.
func TestPropertyToPropertyDateTimeComparisons(t *testing.T) {
	t.Log("\n📅 PROPERTY-TO-PROPERTY DATETIME COMPARISONS")
	t.Log("===========================================")

	engine := rule.NewEngine()

	context := rule.D{
		// Top-level timestamps
		"start_time": "2024-07-10T10:00:00Z",
		"end_time":   "2024-07-10T12:00:00Z",

		// Nested timestamps
		"event": rule.D{
			"start_time": "2024-07-10T10:00:00Z",
			"end_time":   "2024-07-10T12:00:00Z",
		},
		"session": rule.D{
			"created_at": "2024-07-10T09:30:00Z",
			"expires_at": "2024-07-10T11:30:00Z",
		},

		// Deep nested timestamps
		"booking": rule.D{
			"schedule": rule.D{
				"time": rule.D{
					"start": "2024-07-10T08:00:00Z",
					"end":   "2024-07-10T14:00:00Z",
				},
			},
		},
	}

	tests := []struct {
		name     string
		rule     string
		expected bool
	}{
		// Simple to simple datetime comparisons
		{"DateTime simple - start before end", "start_time be end_time", true},
		{"DateTime simple - end after start", "end_time af start_time", true},
		{"DateTime simple - start equals start", "start_time dq start_time", true},

		// Nested to simple
		{"DateTime nested to simple - event start after simple start", "event.start_time aq start_time", true},
		{"DateTime nested to simple - session created before simple end", "session.created_at be end_time", true},

		// Simple to nested
		{"DateTime simple to nested - simple start before event end", "start_time be event.end_time", true},
		{"DateTime simple to nested - simple end after session created", "end_time af session.created_at", true},

		// Nested to nested
		{
			"DateTime nested to nested - event start after session created",
			"event.start_time af session.created_at",
			true,
		},
		{"DateTime nested to nested - session expires before event end", "session.expires_at be event.end_time", true},

		// Deep nested comparisons
		{
			"DateTime deep nested - booking start before event start",
			"booking.schedule.time.start be event.start_time",
			true,
		},
		{"DateTime deep nested - booking end after event end", "booking.schedule.time.end af event.end_time", true},
	}

	for i, test := range tests {
		t.Logf("\n%d. %s", i+1, test.name)
		t.Logf("   Rule: %s", test.rule)

		result, err := engine.Evaluate(test.rule, context)
		if err != nil {
			t.Errorf("❌ Error evaluating rule: %v", err)
			continue
		}

		if result != test.expected {
			t.Errorf("❌ Expected %t, got %t", test.expected, result)
		} else {
			t.Logf("✅ Correct result: %t", result)
		}
	}
}

// TestPropertyToPropertyErrorHandling tests error scenarios in property comparisons.
func TestPropertyToPropertyErrorHandling(t *testing.T) {
	t.Log("\n⚠️  PROPERTY-TO-PROPERTY ERROR HANDLING")
	t.Log("=====================================")

	engine := rule.NewEngine()

	context := rule.D{
		"existing": 42,
		"user": rule.D{
			"age": 25,
		},
	}

	tests := []struct {
		name        string
		rule        string
		expected    bool
		description string
	}{
		{
			name:        "Missing left property",
			rule:        "nonexistent eq existing",
			expected:    false,
			description: "Should return false when left property doesn't exist",
		},
		{
			name:        "Missing right property",
			rule:        "existing eq nonexistent",
			expected:    false,
			description: "Should return false when right property doesn't exist",
		},
		{
			name:        "Both properties missing",
			rule:        "missing1 eq missing2",
			expected:    false, // Both missing properties evaluate to false in our implementation
			description: "Should return false when both properties are missing (defensive behavior)",
		},
		{
			name:        "Missing nested property",
			rule:        "user.nonexistent eq existing",
			expected:    false,
			description: "Should return false when nested property doesn't exist",
		},
		{
			name:        "Deep missing property",
			rule:        "user.profile.missing eq existing",
			expected:    false,
			description: "Should return false when deep nested property doesn't exist",
		},
	}

	for i, test := range tests {
		t.Logf("\n%d. %s", i+1, test.name)
		t.Logf("   Rule: %s", test.rule)
		t.Logf("   Description: %s", test.description)

		result, err := engine.Evaluate(test.rule, context)
		if err != nil {
			t.Errorf("❌ Error evaluating rule (should handle gracefully): %v", err)
			continue
		}

		if result != test.expected {
			t.Errorf("❌ Expected %t, got %t", test.expected, result)
		} else {
			t.Logf("✅ Correct result: %t", result)
		}
	}
}
