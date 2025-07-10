// Package examples demonstrates basic usage of the rule engine
package examples

import (
	"fmt"
	"log"

	"github.com/NSXBet/rule-engine"
)

// BasicUsageExample demonstrates simple rule evaluation
func BasicUsageExample() {
	fmt.Println("🚀 Basic Rule Engine Usage")
	fmt.Println("==========================")

	// Create a new rule engine
	engine := rule.NewEngine()

	// Define context data using rule.D for cleaner syntax
	context := rule.D{
		"user": rule.D{
			"age":    25,
			"status": "active",
			"name":   "John Doe",
			"email":  "john@example.com",
		},
		"account": rule.D{
			"balance": 1000.50,
			"type":    "premium",
		},
	}

	// Example rules
	rules := []struct {
		name string
		rule string
	}{
		{"Adult user", `user.age ge 18`},
		{"Active account", `user.status eq "active"`},
		{"Premium customer", `account.type eq "premium" and account.balance gt 500`},
		{"Email verification", `user.email co "@" and user.email ew ".com"`},
		{"VIP eligibility", `user.age gt 21 and account.balance ge 1000 and user.status eq "active"`},
	}

	fmt.Println("Evaluating rules:")
	fmt.Println("-----------------")

	for _, r := range rules {
		result, err := engine.Evaluate(r.rule, context)
		if err != nil {
			log.Printf("❌ Error evaluating '%s': %v", r.name, err)
			continue
		}

		status := "❌"
		if result {
			status = "✅"
		}
		fmt.Printf("%s %s: %s -> %t\n", status, r.name, r.rule, result)
	}

	fmt.Println("\n✨ Basic usage completed!")
}

// TypeSafetyExample demonstrates type safety features
func TypeSafetyExample() {
	fmt.Println("\n🛡️  Type Safety Demonstration")
	fmt.Println("============================")

	engine := rule.NewEngine()

	context := rule.D{
		"count":  42,
		"flag":   true,
		"text":   "42",
		"score":  95.5,
		"items":  []any{"apple", "banana", "cherry"},
	}

	// Type safety rules
	rules := []struct {
		name     string
		rule     string
		expected bool
		note     string
	}{
		{"Same type int", `count eq 42`, true, "Same types compare correctly"},
		{"Cross-type int/string", `count eq "42"`, false, "Different types never equal"},
		{"Cross-type bool/int", `flag eq 1`, false, "Boolean doesn't equal number"},
		{"Numeric cross-type", `count eq 42.0`, true, "Int/float comparison allowed"},
		{"String contains", `text co "4"`, true, "String operations work correctly"},
		{"Array membership", `"apple" in items`, true, "Array membership works"},
	}

	fmt.Println("Type safety tests:")
	fmt.Println("------------------")

	for _, r := range rules {
		result, err := engine.Evaluate(r.rule, context)
		if err != nil {
			log.Printf("❌ Error: %v", err)
			continue
		}

		status := "❌"
		if result == r.expected {
			status = "✅"
		}
		fmt.Printf("%s %s: %t (expected %t) - %s\n", 
			status, r.name, result, r.expected, r.note)
	}

	fmt.Println("\n✨ Type safety demonstration completed!")
}

// PerformanceExample demonstrates query caching for performance
func PerformanceExample() {
	fmt.Println("\n⚡ Performance & Caching Example")
	fmt.Println("================================")

	engine := rule.NewEngine()

	// Pre-compile frequently used rules for maximum performance
	frequentRules := []string{
		`user.role eq "admin"`,
		`user.subscription.active eq true`,
		`user.age ge 18 and user.status eq "verified"`,
		`account.balance gt 100 and account.type in ["premium", "enterprise"]`,
	}

	fmt.Println("Pre-compiling frequent rules...")
	for _, rule := range frequentRules {
		if err := engine.AddQuery(rule); err != nil {
			log.Fatalf("❌ Failed to compile rule: %s - %v", rule, err)
		}
	}
	fmt.Println("✅ Rules pre-compiled and cached")

	// Test context
	context := rule.D{
		"user": rule.D{
			"role":   "admin",
			"age":    30,
			"status": "verified",
			"subscription": rule.D{
				"active": true,
			},
		},
		"account": rule.D{
			"balance": 500.0,
			"type":    "premium",
		},
	}

	fmt.Println("\nTesting cached rules (lightning fast!):")
	fmt.Println("--------------------------------------")

	for _, rule := range frequentRules {
		result, err := engine.Evaluate(rule, context)
		if err != nil {
			log.Printf("❌ Error: %v", err)
			continue
		}
		fmt.Printf("✅ %s -> %t\n", rule, result)
	}

	fmt.Println("\n⚡ Performance example completed!")
	fmt.Println("Note: Subsequent evaluations of these rules will be 25-100ns (cached AST)")
}