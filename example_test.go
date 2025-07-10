package rule_test

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/NSXBet/rule"
	ruleslib "github.com/nikunjy/rules"
)

// Example demonstrates basic rule engine usage.
func Example() {
	engine := rule.NewEngine()

	context := rule.D{
		"user": rule.D{
			"age":    25,
			"status": "active",
		},
	}

	result, err := engine.Evaluate(`user.age gt 18 and user.status eq "active"`, context)
	if err != nil {
		slog.Error("Rule evaluation failed", "error", err)
		return
	}

	fmt.Printf("User is eligible: %t", result)
	// Output: User is eligible: true
}

// ExampleEngine_Evaluate demonstrates rule evaluation.
func ExampleEngine_Evaluate() {
	engine := rule.NewEngine()

	context := rule.D{
		"score": 85,
		"level": "premium",
	}

	result, _ := engine.Evaluate(`score ge 80 and level eq "premium"`, context)
	fmt.Printf("Eligible for bonus: %t", result)
	// Output: Eligible for bonus: true
}

// Example_basicUsage shows fundamental operations.
func Example_basicUsage() {
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
			slog.Error("Rule evaluation failed", "error", err)
			continue
		}

		fmt.Printf("%s -> %t\n", r, result)
	}
	// Output:
	// user.age ge 18 -> true
	// user.status eq "active" -> true
	// account.balance gt 500 -> true
	// user.age gt 21 and account.type eq "premium" -> true
}

// Example_dateTimeOperations demonstrates datetime capabilities.
func Example_dateTimeOperations() {
	engine := rule.NewEngine()

	now := time.Date(2024, 7, 10, 15, 30, 0, 0, time.UTC)
	oneHourAgo := now.Add(-1 * time.Hour)

	context := rule.D{
		"created_at": now,
		"updated_at": now.Format(time.RFC3339),
		"deadline":   now.Add(24 * time.Hour).Unix(),
		"start_time": oneHourAgo.Format(time.RFC3339),
	}

	rules := []string{
		`created_at dq updated_at`, // DateTime equal
		`start_time be created_at`, // Before
		`deadline af created_at`,   // After
	}

	for _, r := range rules {
		result, err := engine.Evaluate(r, context)
		if err != nil {
			slog.Error("Rule evaluation failed", "error", err)
			continue
		}

		fmt.Printf("%s -> %t\n", r, result)
	}
	// Output:
	// created_at dq updated_at -> true
	// start_time be created_at -> true
	// deadline af created_at -> true
}

// Example_compatibility shows compatibility with nikunjy/rules.
func Example_compatibility() {
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
		fmt.Printf("Both libraries return: %t (100%% compatible!)", oldResult)
	} else {
		fmt.Printf("nikunjy/rules: %t (err: %v)\n", oldResult, oldErr)
		fmt.Printf("Our library: %t (err: %v)", newResult, newErr)
	}
	// Output:
	// Rule: user.age gt 18 and user.status eq "active"
	// Both libraries return: true (100% compatible!)
}

// Example_migration demonstrates migration from nikunjy/rules.
func Example_migration() {
	// BEFORE (nikunjy/rules style)
	oldContext := map[string]interface{}{
		"user": map[string]interface{}{
			"age":  25,
			"role": "admin",
		},
	}
	oldResult, _ := ruleslib.Evaluate(`user.age gt 18`, oldContext)

	// AFTER (our library)
	engine := rule.NewEngine()
	newContext := rule.D{
		"user": rule.D{
			"age":  25,
			"role": "admin",
		},
	}
	newResult, _ := engine.Evaluate(`user.age gt 18`, newContext)

	fmt.Printf("nikunjy/rules result: %t\n", oldResult)
	fmt.Printf("Our library result: %t\n", newResult)
	fmt.Printf("Compatible: %t", oldResult == newResult)
	// Output:
	// nikunjy/rules result: true
	// Our library result: true
	// Compatible: true
}

// Example_performance shows performance-optimized usage.
func Example_performance() {
	engine := rule.NewEngine()

	// Pre-compile frequent rules for maximum performance
	frequentRules := []string{
		`user.role eq "admin"`,
		`account.balance gt 1000`,
	}

	for _, rule := range frequentRules {
		if err := engine.AddQuery(rule); err != nil {
			slog.Error("Failed to compile rule", "error", err)
			return
		}
	}

	context := rule.D{
		"user":    rule.D{"role": "admin"},
		"account": rule.D{"balance": 1500.0},
	}

	// Lightning-fast evaluation of pre-compiled rules
	for _, rule := range frequentRules {
		result, _ := engine.Evaluate(rule, context)
		fmt.Printf("%s -> %t\n", rule, result)
	}
	// Output:
	// user.role eq "admin" -> true
	// account.balance gt 1000 -> true
}

// Example_ecommerce demonstrates e-commerce business rules.
func Example_ecommerce() {
	engine := rule.NewEngine()

	customer := rule.D{
		"age":              28,
		"membership_years": 3,
		"location":         "US",
		"total_spent":      2500.00,
	}

	order := rule.D{
		"total":       150.00,
		"items_count": 3,
		"category":    "electronics",
		"is_weekend":  true,
	}

	context := rule.D{
		"customer": customer,
		"order":    order,
	}

	// Check eligibility for free shipping
	freeShippingRule := `customer.total_spent gt 1000 and order.total gt 100`
	freeShipping, _ := engine.Evaluate(freeShippingRule, context)

	// Check discount eligibility
	discountRule := `customer.membership_years ge 2 and order.is_weekend eq true`
	weekendDiscount, _ := engine.Evaluate(discountRule, context)

	fmt.Printf("Free shipping eligible: %t\n", freeShipping)
	fmt.Printf("Weekend discount eligible: %t", weekendDiscount)
	// Output:
	// Free shipping eligible: true
	// Weekend discount eligible: true
}
