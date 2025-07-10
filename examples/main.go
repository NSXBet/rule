// Package main demonstrates all rule engine examples
package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/NSXBet/rule"
	ruleslib "github.com/nikunjy/rules"
)

func main() {
	fmt.Println("🚀 Rule Engine Examples")
	fmt.Println("=======================")
	fmt.Println("This demo showcases the NSXBet Rule Engine capabilities")
	fmt.Println("including 100% compatibility with nikunjy/rules and our extensions.")
	fmt.Println("")

	// Check if specific example was requested
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "basic":
			examples.BasicUsageExample()
			examples.TypeSafetyExample()
			examples.PerformanceExample()
		case "ecommerce":
			examples.EcommerceEligibilityExample()
			examples.DynamicPricingExample()
		case "datetime":
			examples.DateTimeOperationsExample()
			examples.SchedulingExample()
		case "migration":
			examples.MigrationExample()
			examples.BeforeAfterExample()
			examples.PerformanceComparisonExample()
		default:
			fmt.Printf("❌ Unknown example: %s\n", os.Args[1])
			fmt.Println("Available examples: basic, ecommerce, datetime, migration")
			os.Exit(1)
		}
		return
	}

	// Run all examples
	fmt.Println("Running all examples...")
	fmt.Println("Press Enter to continue between sections...")

	// Basic Usage Examples
	fmt.Println("\n" + strings.Repeat("=", 60))
	examples.BasicUsageExample()
	examples.TypeSafetyExample()
	examples.PerformanceExample()
	waitForEnter()

	// E-commerce Examples
	fmt.Println("\n" + strings.Repeat("=", 60))
	examples.EcommerceEligibilityExample()
	examples.DynamicPricingExample()
	waitForEnter()

	// DateTime Examples (Our Extension)
	fmt.Println("\n" + strings.Repeat("=", 60))
	examples.DateTimeOperationsExample()
	examples.SchedulingExample()
	waitForEnter()

	// Migration Examples
	fmt.Println("\n" + strings.Repeat("=", 60))
	examples.MigrationExample()
	examples.BeforeAfterExample()
	examples.PerformanceComparisonExample()

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🎉 All examples completed!")
	fmt.Println("")
	fmt.Println("Key takeaways:")
	fmt.Println("• 100% compatibility with nikunjy/rules")
	fmt.Println("• 25-144x performance improvement")
	fmt.Println("• Zero allocations during evaluation")
	fmt.Println("• Native datetime operators")
	fmt.Println("• Thread-safe query caching")
	fmt.Println("• Clean API with rule.D type alias")
	fmt.Println("")
	fmt.Println("To run specific examples:")
	fmt.Println("  go run main.go basic      # Basic usage and type safety")
	fmt.Println("  go run main.go ecommerce  # E-commerce business rules")
	fmt.Println("  go run main.go datetime   # DateTime operations")
	fmt.Println("  go run main.go migration  # Migration from nikunjy/rules")
}

func waitForEnter() {
	fmt.Print("\nPress Enter to continue...")
	fmt.Scanln()
}
