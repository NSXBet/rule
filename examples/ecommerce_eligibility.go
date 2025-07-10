// Package examples demonstrates e-commerce business rule evaluation
package examples

import (
	"fmt"
	"log"

	"github.com/NSXBet/rule"
)

// Customer represents an e-commerce customer
type Customer struct {
	ID       string
	Age      int
	Country  string
	Status   string
	Tier     string
	Email    string
	Verified bool
}

// Order represents a customer order
type Order struct {
	Total    float64
	Items    int
	Category string
	Shipping string
}

// Account represents customer account details
type Account struct {
	Balance       float64
	CreditLimit   float64
	PaymentMethod string
	History       AccountHistory
}

// AccountHistory represents account history
type AccountHistory struct {
	TotalOrders   int
	TotalSpent    float64
	LastOrderDays int
	Disputes      int
}

// EcommerceEligibilityExample demonstrates complex business rules for e-commerce
func EcommerceEligibilityExample() {
	fmt.Println("🛒 E-commerce Eligibility Rules")
	fmt.Println("===============================")

	engine := rule.NewEngine()

	// Sample customer data
	customer := Customer{
		ID:       "CUST123",
		Age:      28,
		Country:  "US",
		Status:   "active",
		Tier:     "gold",
		Email:    "customer@example.com",
		Verified: true,
	}

	order := Order{
		Total:    250.00,
		Items:    3,
		Category: "electronics",
		Shipping: "express",
	}

	account := Account{
		Balance:       1500.00,
		CreditLimit:   5000.00,
		PaymentMethod: "credit_card",
		History: AccountHistory{
			TotalOrders:   15,
			TotalSpent:    3200.00,
			LastOrderDays: 5,
			Disputes:      0,
		},
	}

	// Convert to rule.D context
	context := rule.D{
		"customer": rule.D{
			"id":       customer.ID,
			"age":      customer.Age,
			"country":  customer.Country,
			"status":   customer.Status,
			"tier":     customer.Tier,
			"email":    customer.Email,
			"verified": customer.Verified,
		},
		"order": rule.D{
			"total":    order.Total,
			"items":    order.Items,
			"category": order.Category,
			"shipping": order.Shipping,
		},
		"account": rule.D{
			"balance":        account.Balance,
			"credit_limit":   account.CreditLimit,
			"payment_method": account.PaymentMethod,
			"history": rule.D{
				"total_orders":    account.History.TotalOrders,
				"total_spent":     account.History.TotalSpent,
				"last_order_days": account.History.LastOrderDays,
				"disputes":        account.History.Disputes,
			},
		},
	}

	// Business rules for various eligibility checks
	businessRules := []struct {
		name        string
		rule        string
		description string
	}{
		{
			"Basic Eligibility",
			`customer.age ge 18 and customer.status eq "active" and customer.verified eq true`,
			"Customer must be adult, active, and verified",
		},
		{
			"Order Size Limit",
			`order.total le 10000 or (customer.tier in ["gold", "platinum"] and account.credit_limit ge order.total)`,
			"Large orders require premium tier or sufficient credit",
		},
		{
			"Express Shipping Eligibility",
			`order.shipping eq "express" and (customer.tier in ["gold", "platinum"] or order.total ge 100)`,
			"Express shipping for premium customers or large orders",
		},
		{
			"Electronics Category Rules",
			`order.category eq "electronics" and customer.country in ["US", "CA", "UK"] and customer.age ge 21`,
			"Electronics sales restricted by location and age",
		},
		{
			"Credit Eligibility",
			`account.payment_method eq "credit_card" and account.history.disputes eq 0 and account.history.total_orders ge 5`,
			"Credit purchases require clean history and experience",
		},
		{
			"VIP Treatment",
			`customer.tier eq "platinum" or (account.history.total_spent ge 2000 and account.history.disputes eq 0)`,
			"VIP treatment for platinum tier or high-value customers",
		},
		{
			"Loyalty Discount",
			`account.history.total_orders ge 10 and account.history.last_order_days le 30 and customer.status eq "active"`,
			"Loyalty discount for frequent, recent customers",
		},
		{
			"Risk Assessment",
			`account.history.disputes eq 0 and account.balance ge 0 and customer.verified eq true`,
			"Low-risk customer assessment",
		},
		{
			"Bulk Order Processing",
			`order.items ge 5 and (account.balance ge order.total or account.credit_limit ge order.total)`,
			"Bulk orders require sufficient funds or credit",
		},
		{
			"Premium Support Access",
			`customer.tier in ["gold", "platinum"] and account.history.total_spent ge 1000`,
			"Premium support for valuable customers",
		},
	}

	fmt.Printf("Evaluating eligibility for Customer %s:\n", customer.ID)
	fmt.Println("---------------------------------------")

	eligibleCount := 0
	totalRules := len(businessRules)

	for _, br := range businessRules {
		result, err := engine.Evaluate(br.rule, context)
		if err != nil {
			log.Printf("❌ Error evaluating '%s': %v", br.name, err)
			continue
		}

		status := "❌"
		if result {
			status = "✅"
			eligibleCount++
		}

		fmt.Printf("%s %s\n", status, br.name)
		fmt.Printf("   Rule: %s\n", br.rule)
		fmt.Printf("   Description: %s\n", br.description)
		fmt.Printf("   Result: %t\n\n", result)
	}

	fmt.Printf("📊 Eligibility Summary:\n")
	fmt.Printf("   Passed: %d/%d rules (%.1f%%)\n",
		eligibleCount, totalRules, float64(eligibleCount)/float64(totalRules)*100)

	if eligibleCount == totalRules {
		fmt.Printf("🎉 Customer %s is eligible for all features!\n", customer.ID)
	} else if float64(eligibleCount)/float64(totalRules) >= 0.8 {
		fmt.Printf("⭐ Customer %s has high eligibility (80%+ rules passed)\n", customer.ID)
	} else {
		fmt.Printf("⚠️  Customer %s has limited eligibility\n", customer.ID)
	}

	fmt.Println("\n✨ E-commerce eligibility evaluation completed!")
}

// DynamicPricingExample demonstrates dynamic pricing rules
func DynamicPricingExample() {
	fmt.Println("\n💰 Dynamic Pricing Rules")
	fmt.Println("========================")

	engine := rule.NewEngine()

	// Product and customer context
	context := rule.D{
		"product": rule.D{
			"base_price": 100.0,
			"category":   "electronics",
			"stock":      5,
			"rating":     4.5,
		},
		"customer": rule.D{
			"tier":          "gold",
			"loyalty_years": 3,
			"last_purchase": 15, // days ago
		},
		"market": rule.D{
			"demand":     "high",
			"season":     "holiday",
			"competitor": 95.0,
		},
	}

	// Pricing rules (these would calculate discounts/markups)
	pricingRules := []struct {
		name     string
		rule     string
		modifier string
	}{
		{
			"Loyalty Discount",
			`customer.tier in ["gold", "platinum"] and customer.loyalty_years ge 2`,
			"-10%",
		},
		{
			"Recent Customer Discount",
			`customer.last_purchase le 30`,
			"-5%",
		},
		{
			"Low Stock Premium",
			`product.stock le 10 and market.demand eq "high"`,
			"+15%",
		},
		{
			"Holiday Markup",
			`market.season eq "holiday" and product.category in ["electronics", "toys"]`,
			"+20%",
		},
		{
			"Competitive Pricing",
			`market.competitor lt product.base_price`,
			"Match competitor",
		},
		{
			"High Rating Bonus",
			`product.rating ge 4.0`,
			"+5%",
		},
	}

	fmt.Println("Evaluating pricing rules:")
	fmt.Println("-------------------------")

	basePrice := context["product"].(rule.D)["base_price"].(float64)
	fmt.Printf("Base Price: $%.2f\n\n", basePrice)

	for _, pr := range pricingRules {
		result, err := engine.Evaluate(pr.rule, context)
		if err != nil {
			log.Printf("❌ Error: %v", err)
			continue
		}

		status := "❌"
		if result {
			status = "✅"
		}

		fmt.Printf("%s %s: %s\n", status, pr.name, pr.modifier)
		fmt.Printf("   Rule: %s\n", pr.rule)
		fmt.Printf("   Applied: %t\n\n", result)
	}

	fmt.Println("💡 In a real system, these rules would calculate the final price")
	fmt.Println("   by applying the modifiers to the base price")
	fmt.Println("\n✨ Dynamic pricing evaluation completed!")
}
