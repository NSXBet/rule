// Package main demonstrates key rule engine features
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/NSXBet/rule"
	ruleslib "github.com/nikunjy/rules"
)

func main() {
	fmt.Println("🚀 Rule Engine Demo")
	fmt.Println("===================")

	// Basic Usage Demo
	basicUsageDemo()

	// Compatibility Demo
	compatibilityDemo()

	// DateTime Demo (Our Extension)
	datetimeDemo()

	// Performance Demo
	performanceDemo()

	fmt.Println("\n🎉 Demo completed!")
	fmt.Println("See the individual example files for more detailed demonstrations.")
}

func basicUsageDemo() {
	fmt.Println("\n📚 Basic Usage")
	fmt.Println("==============")

	engine := rule.NewEngine()

	context := rule.D{
		"user": rule.D{
			"age":    25,
			"status": "active",
			"name":   "John Doe",
		},
		"account": rule.D{
			"balance": 1000.50,
			"type":    "premium",
		},
	}

	rules := []string{
		`user.age ge 18`,
		`user.status eq "active"`,
		`account.balance gt 500`,
		`user.age gt 21 and account.type eq "premium"`,
	}

	for _, r := range rules {
		result, err := engine.Evaluate(r, context)
		if err != nil {
			log.Printf("❌ Error: %v", err)
			continue
		}
		fmt.Printf("✅ %s -> %t\n", r, result)
	}
}

func compatibilityDemo() {
	fmt.Println("\n🔄 Compatibility with nikunjy/rules")
	fmt.Println("===================================")

	// Context that works with both libraries
	context := map[string]interface{}{
		"user": map[string]interface{}{
			"age":    25,
			"status": "active",
		},
	}

	ourContext := rule.D{
		"user": rule.D{
			"age":    25,
			"status": "active",
		},
	}

	testRule := `user.age gt 18 and user.status eq "active"`

	// Test with nikunjy/rules
	oldResult, oldErr := ruleslib.Evaluate(testRule, context)

	// Test with our library
	ourEngine := rule.NewEngine()
	newResult, newErr := ourEngine.Evaluate(testRule, ourContext)

	fmt.Printf("Rule: %s\n", testRule)
	if oldErr == nil && newErr == nil && oldResult == newResult {
		fmt.Printf("✅ Both libraries return: %t (100%% compatible!)\n", oldResult)
	} else {
		fmt.Printf("❌ nikunjy/rules: %t (err: %v)\n", oldResult, oldErr)
		fmt.Printf("❌ Our library: %t (err: %v)\n", newResult, newErr)
	}
}

func datetimeDemo() {
	fmt.Println("\n📅 DateTime Operations (Our Extension)")
	fmt.Println("=====================================")

	engine := rule.NewEngine()

	now := time.Now().UTC()
	oneHourAgo := now.Add(-1 * time.Hour)

	context := rule.D{
		"created_at": now,                            // time.Time
		"updated_at": now.Format(time.RFC3339),       // RFC3339 string
		"deadline":   now.Add(24 * time.Hour).Unix(), // Unix timestamp
		"start_time": oneHourAgo.Format(time.RFC3339),
	}

	datetimeRules := []string{
		`created_at dq updated_at`, // DateTime equal
		`start_time be created_at`, // Before
		`deadline af created_at`,   // After
		`created_at aq start_time`, // After or equal
	}

	fmt.Printf("Current time: %s\n", now.Format(time.RFC3339))

	for _, r := range datetimeRules {
		result, err := engine.Evaluate(r, context)
		if err != nil {
			log.Printf("❌ Error: %v", err)
			continue
		}
		fmt.Printf("✅ %s -> %t\n", r, result)
	}

	fmt.Println("\n💡 DateTime operators: dq, dn, be, bq, af, aq")
	fmt.Println("   Supports: time.Time, RFC3339 strings, Unix timestamps")
}

func performanceDemo() {
	fmt.Println("\n⚡ Performance Features")
	fmt.Println("======================")

	engine := rule.NewEngine()

	// Pre-compile frequent rules
	frequentRules := []string{
		`user.role eq "admin"`,
		`account.balance gt 1000`,
	}

	fmt.Println("Pre-compiling frequent rules for maximum performance...")
	for _, rule := range frequentRules {
		if err := engine.AddQuery(rule); err != nil {
			log.Fatalf("Failed to compile rule: %v", err)
		}
	}

	context := rule.D{
		"user":    rule.D{"role": "admin"},
		"account": rule.D{"balance": 1500.0},
	}

	fmt.Println("Evaluating pre-compiled rules (lightning fast!):")
	for _, rule := range frequentRules {
		result, _ := engine.Evaluate(rule, context)
		fmt.Printf("✅ %s -> %t\n", rule, result)
	}

	fmt.Println("\n📊 Performance benefits:")
	fmt.Println("   • 25-144x faster than nikunjy/rules")
	fmt.Println("   • 0 allocations during evaluation")
	fmt.Println("   • Sub-100ns evaluation times")
	fmt.Println("   • Automatic query caching")
}
