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

// Example_daysLessOperator demonstrates the "dl" (days less) operator.
func Example_daysLessOperator() {
	engine := rule.NewEngine()

	now := time.Now().UTC()

	// Sample context with various timestamp formats calculated relative to now
	context := rule.D{
		"user_registered":  now.AddDate(-2, 0, 0).Format(time.RFC3339), // 2 years ago
		"last_login":       now.AddDate(0, 0, -200).Unix(),             // 200 days ago
		"password_changed": now.AddDate(0, 0, -7).Format(time.RFC3339), // 7 days ago
		"account_created":  now.AddDate(-5, 0, 0).Format(time.RFC3339), // 5 years ago
	}

	// Check if events happened within specific time ranges from NOW
	rules := []string{
		`user_registered dl 365`,  // Within last 365 days (about 1 year)
		`last_login dl 400`,       // Within last 400 days (over 1 year)
		`password_changed dl 30`,  // Within last 30 days
		`account_created dl 1000`, // Within last 1000 days (about 3 years)
		`user_registered dl 1.5`,  // Within last 1.5 days (fractional)
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
	// user_registered dl 365 -> false
	// last_login dl 400 -> true
	// password_changed dl 30 -> true
	// account_created dl 1000 -> false
	// user_registered dl 1.5 -> false
}

// Example_daysLessUseCase demonstrates practical use cases for the "dl" operator.
func Example_daysLessUseCase() {
	engine := rule.NewEngine()

	now := time.Now().UTC()

	// User session and security context
	context := rule.D{
		"user": rule.D{
			"last_login":       now.AddDate(0, 0, -10).Format(time.RFC3339), // 10 days ago
			"password_changed": now.AddDate(0, 0, -20).Format(time.RFC3339), // 20 days ago
			"mfa_enabled":      true,
		},
		"session": rule.D{
			"created_at": now.Add(-12 * time.Hour).Format(time.RFC3339), // 12 hours ago (less than 1 day)
			"ip_address": "192.168.1.100",
		},
	}

	// Security and business rules using "dl" operator
	securityRules := []rule.D{
		{
			"name": "Recent login check",
			"rule": `user.last_login dl 30`,
			"desc": "User logged in within last 30 days",
		},
		{
			"name": "Password age check",
			"rule": `user.password_changed dl 90`,
			"desc": "Password changed within last 90 days",
		},
		{
			"name": "Active session check",
			"rule": `session.created_at dl 1 and user.mfa_enabled eq true`,
			"desc": "Session created within 1 day and MFA enabled",
		},
	}

	for _, ruleData := range securityRules {
		result, err := engine.Evaluate(ruleData["rule"].(string), context)
		if err != nil {
			slog.Error("Rule evaluation failed", "error", err)
			continue
		}

		fmt.Printf("%s: %t\n", ruleData["name"], result)
	}
	// Output:
	// Recent login check: true
	// Password age check: true
	// Active session check: true
}

// Example_daysGreaterOperator demonstrates the "dg" (days greater) operator.
func Example_daysGreaterOperator() {
	engine := rule.NewEngine()

	now := time.Now().UTC()

	// Sample context with various timestamp formats (old timestamps)
	context := rule.D{
		"account_created":  now.AddDate(-5, 0, 0).Format(time.RFC3339),  // 5 years ago
		"last_backup":      now.AddDate(-2, 0, 0).Unix(),                // 2 years ago (Unix timestamp)
		"system_installed": now.AddDate(-6, 0, 0).Format(time.RFC3339),  // 6 years ago
		"recent_update":    now.AddDate(0, 0, -5).Format(time.RFC3339),  // 5 days ago
		"maintenance_done": now.AddDate(-1, -6, 0).Format(time.RFC3339), // 1.5 years ago
	}

	// Check if events happened MORE than specific time ranges from NOW
	rules := []string{
		`account_created dg 365`,   // More than 365 days ago (about 1 year)
		`last_backup dg 400`,       // More than 400 days ago
		`system_installed dg 1000`, // More than 1000 days ago (about 3 years)
		`recent_update dg 30`,      // More than 30 days ago
		`maintenance_done dg 365`,  // More than 365 days ago
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
	// account_created dg 365 -> true
	// last_backup dg 400 -> true
	// system_installed dg 1000 -> true
	// recent_update dg 30 -> false
	// maintenance_done dg 365 -> true
}

// Example_dlVsDgComparison demonstrates the opposite behavior of "dl" vs "dg" operators.
func Example_dlVsDgComparison() {
	engine := rule.NewEngine()

	now := time.Now().UTC()

	// Test with the same timestamp for both operators
	context := rule.D{
		"event_timestamp": now.AddDate(-2, 0, 0).Format(time.RFC3339), // 2 years ago
	}

	// DL (days less) checks if timestamp is WITHIN the threshold from NOW
	// DG (days greater) checks if timestamp is BEYOND the threshold from NOW
	comparisonRules := []struct {
		name string
		rule string
		desc string
	}{
		{"DL - Within 365 days", `event_timestamp dl 365`, "Should be false (beyond 365 days)"},
		{"DG - Beyond 365 days", `event_timestamp dg 365`, "Should be true (beyond 365 days)"},
		{"DL - Within 1000 days", `event_timestamp dl 1000`, "Should be true (within 1000 days)"},
		{"DG - Beyond 1000 days", `event_timestamp dg 1000`, "Should be false (within 1000 days)"},
	}

	for _, ruleData := range comparisonRules {
		result, err := engine.Evaluate(ruleData.rule, context)
		if err != nil {
			slog.Error("Rule evaluation failed", "error", err)
			continue
		}

		fmt.Printf("%-21s: %t (%s)\n", ruleData.name, result, ruleData.desc)
	}
	// Output:
	// DL - Within 365 days : false (Should be false (beyond 365 days))
	// DG - Beyond 365 days : true (Should be true (beyond 365 days))
	// DL - Within 1000 days: true (Should be true (within 1000 days))
	// DG - Beyond 1000 days: false (Should be false (within 1000 days))
}

// Example_daysOperatorsUseCase demonstrates practical use cases combining "dl" and "dg" operators.
func Example_daysOperatorsUseCase() {
	engine := rule.NewEngine()

	now := time.Now().UTC()

	// System maintenance and security context
	context := rule.D{
		"system": rule.D{
			"last_security_scan": now.AddDate(0, 0, -3).Format(time.RFC3339), // 3 days ago
			"last_full_backup":   now.AddDate(0, -8, 0).Format(time.RFC3339), // 8 months ago
			"os_install_date":    now.AddDate(-5, 0, 0).Format(time.RFC3339), // 5 years ago
		},
		"user": rule.D{
			"password_changed": now.AddDate(0, 0, -10).Format(time.RFC3339), // 10 days ago
			"account_created":  now.AddDate(-6, 0, 0).Format(time.RFC3339),  // 6 years ago
		},
	}

	// Business rules combining both operators
	systemRules := []rule.D{
		{
			"name": "Security scan up-to-date",
			"rule": `system.last_security_scan dl 7`,
			"desc": "Security scan within last 7 days",
		},
		{
			"name": "Backup overdue",
			"rule": `system.last_full_backup dg 180`,
			"desc": "Last backup more than 180 days ago",
		},
		{
			"name": "Legacy system",
			"rule": `system.os_install_date dg 1825`,
			"desc": "OS installed more than 5 years ago",
		},
		{
			"name": "Password policy compliance",
			"rule": `user.password_changed dl 90`,
			"desc": "Password changed within last 90 days",
		},
		{
			"name": "Established user",
			"rule": `user.account_created dg 365`,
			"desc": "Account created more than 1 year ago",
		},
	}

	for _, ruleData := range systemRules {
		result, err := engine.Evaluate(ruleData["rule"].(string), context)
		if err != nil {
			slog.Error("Rule evaluation failed", "error", err)
			continue
		}

		status := "✅"
		if !result {
			status = "❌"
		}

		fmt.Printf("%s %s: %t\n", status, ruleData["name"], result)
	}
	// Output:
	// ✅ Security scan up-to-date: true
	// ✅ Backup overdue: true
	// ✅ Legacy system: true
	// ✅ Password policy compliance: true
	// ✅ Established user: true
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
