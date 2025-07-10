package main

import (
	"fmt"
	"log"
	"time"

	rule "github.com/NSXBet/rule-engine"
)

func main() {
	engine := rule.NewEngine()

	// Test time: 2024-07-09 22:12:00 UTC
	testTime := time.Date(2024, 7, 9, 22, 12, 0, 0, time.UTC)

	// Context with three timestamp formats
	context := rule.D{
		"time_obj":     testTime,                         // time.Time{}
		"time_rfc":     "2024-07-09T22:12:00Z",          // RFC 3339 string  
		"time_unix":    int64(1720563120),               // Unix timestamp
		"earlier_time": time.Date(2024, 7, 9, 22, 11, 0, 0, time.UTC),
		"later_time":   time.Date(2024, 7, 9, 22, 13, 0, 0, time.UTC),
	}

	tests := []struct {
		name     string
		rule     string
		expected bool
	}{
		// Regular operators with time.Time (should work like nikunjy/rules via string conversion)
		{"time.Time eq string", `time_obj eq "2024-07-09 22:12:00 +0000 UTC"`, true},
		{"time.Time ne different", `time_obj ne "different"`, true},
		{"time.Time lt lexicographic", `time_obj lt "2024-07-09 23:00:00 +0000 UTC"`, true},
		{"time.Time gt lexicographic", `time_obj gt "2024-07-09 21:00:00 +0000 UTC"`, true},

		// Datetime operators with time.Time (should work with proper datetime comparison)
		{"time.Time dq RFC", `time_obj dq "2024-07-09T22:12:00Z"`, true},
		{"time.Time af earlier", `time_obj af "2024-07-09T22:11:00Z"`, true},
		{"time.Time be later", `time_obj be "2024-07-09T22:13:00Z"`, true},

		// Mixed format comparisons with datetime operators
		{"RFC vs Unix", `time_rfc dq time_unix`, true},
		{"time.Time vs RFC", `time_obj dq time_rfc`, true},
		{"time.Time vs Unix", `time_obj dq time_unix`, true},

		// Property vs Property datetime comparisons
		{"time.Time af earlier time.Time", `time_obj af earlier_time`, true},
		{"time.Time be later time.Time", `time_obj be later_time`, true},
	}

	fmt.Println("Testing time.Time compatibility...")
	fmt.Println("=================================")

	for _, test := range tests {
		result, err := engine.Evaluate(test.rule, context)
		if err != nil {
			log.Printf("ERROR in %s: %v", test.name, err)
			continue
		}

		status := "✅"
		if result != test.expected {
			status = "❌"
		}

		fmt.Printf("%s %s: %s -> %t (expected %t)\n", 
			status, test.name, test.rule, result, test.expected)
	}

	fmt.Println("\n✅ All time.Time compatibility tests completed!")
}