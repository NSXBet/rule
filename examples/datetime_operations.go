// Package examples demonstrates datetime operations (our extension beyond nikunjy/rules)
package examples

import (
	"fmt"
	"log"
	"time"

	"github.com/NSXBet/rule"
)

// DateTimeOperationsExample demonstrates our datetime extensions
func DateTimeOperationsExample() {
	fmt.Println("📅 DateTime Operations (Our Extension)")
	fmt.Println("=====================================")

	engine := rule.NewEngine()

	// Create test timestamps
	now := time.Now().UTC()
	oneHourAgo := now.Add(-1 * time.Hour)
	oneHourLater := now.Add(1 * time.Hour)
	oneDayAgo := now.Add(-24 * time.Hour)
	oneWeekLater := now.Add(7 * 24 * time.Hour)

	// Context with multiple datetime formats
	context := rule.D{
		"event": rule.D{
			"created_at":    now,                             // time.Time
			"updated_at":    now.Format(time.RFC3339),        // RFC3339 string
			"published_at":  now.Unix(),                      // Unix timestamp
			"started_at":    oneHourAgo.Format(time.RFC3339), // RFC3339 string
			"deadline":      oneWeekLater.Unix(),             // Unix timestamp
			"last_modified": oneDayAgo,                       // time.Time
		},
		"user": rule.D{
			"joined_at":     "2023-01-15T10:30:00Z",                        // RFC3339
			"last_login":    time.Date(2024, 7, 9, 15, 30, 0, 0, time.UTC), // time.Time
			"trial_expires": int64(1735689600),                             // 2025-01-01 00:00:00 UTC
		},
		"subscription": rule.D{
			"started": "2024-01-01T00:00:00Z",
			"expires": "2024-12-31T23:59:59Z",
			"renewed": int64(1704067200), // 2024-01-01 00:00:00 UTC
		},
	}

	// DateTime comparison examples
	datetimeRules := []struct {
		name        string
		rule        string
		description string
	}{
		{
			"Event Creation Time",
			`event.created_at dq event.updated_at`,
			"Check if created and updated times are the same",
		},
		{
			"Event Started Before Now",
			`event.started_at be event.created_at`,
			"Check if event started before it was created",
		},
		{
			"Future Deadline",
			`event.deadline af event.created_at`,
			"Check if deadline is after creation time",
		},
		{
			"Recent Modification",
			fmt.Sprintf(`event.last_modified af "%s"`, oneDayAgo.Add(-1*time.Hour).Format(time.RFC3339)),
			"Check if last modification was within the last day",
		},
		{
			"User Trial Active",
			fmt.Sprintf(`user.trial_expires af "%s"`, now.Format(time.RFC3339)),
			"Check if user trial is still active",
		},
		{
			"Long-time User",
			`user.joined_at be "2024-01-01T00:00:00Z"`,
			"Check if user joined before 2024",
		},
		{
			"Recent Login",
			fmt.Sprintf(`user.last_login af "%s"`, now.Add(-7*24*time.Hour).Format(time.RFC3339)),
			"Check if user logged in within the last week",
		},
		{
			"Subscription Active",
			fmt.Sprintf(`subscription.started be "%s" and subscription.expires af "%s"`,
				now.Format(time.RFC3339), now.Format(time.RFC3339)),
			"Check if subscription is currently active",
		},
		{
			"Recent Renewal",
			`subscription.renewed aq "2024-01-01T00:00:00Z"`,
			"Check if subscription was renewed on or after 2024",
		},
		{
			"Event Timeline Validation",
			`event.started_at be event.created_at and event.deadline af event.created_at`,
			"Validate event timeline (started before created, deadline after created)",
		},
	}

	fmt.Println("DateTime operator examples:")
	fmt.Println("----------------------------")
	fmt.Printf("Current time: %s\n", now.Format(time.RFC3339))
	fmt.Printf("One hour ago: %s\n", oneHourAgo.Format(time.RFC3339))
	fmt.Printf("One week later: %s\n\n", oneWeekLater.Format(time.RFC3339))

	for _, dr := range datetimeRules {
		result, err := engine.Evaluate(dr.rule, context)
		if err != nil {
			log.Printf("❌ Error evaluating '%s': %v", dr.name, err)
			continue
		}

		status := "❌"
		if result {
			status = "✅"
		}

		fmt.Printf("%s %s: %t\n", status, dr.name, result)
		fmt.Printf("   Rule: %s\n", dr.rule)
		fmt.Printf("   Description: %s\n\n", dr.description)
	}

	fmt.Println("📋 DateTime Operator Reference:")
	fmt.Println("--------------------------------")
	fmt.Println("dq - DateTime Equal")
	fmt.Println("dn - DateTime Not Equal")
	fmt.Println("be - Before")
	fmt.Println("bq - Before or Equal")
	fmt.Println("af - After")
	fmt.Println("aq - After or Equal")
	fmt.Println("")
	fmt.Println("💡 Supported formats:")
	fmt.Println("   • time.Time{} (from Go code)")
	fmt.Println("   • RFC3339 strings: '2024-07-09T22:12:00Z'")
	fmt.Println("   • Unix timestamps: 1720563120")
	fmt.Println("")
	fmt.Println("✨ DateTime operations completed!")
}

// SchedulingExample demonstrates datetime rules for scheduling systems
func SchedulingExample() {
	fmt.Println("\n🗓️  Scheduling System Example")
	fmt.Println("============================")

	engine := rule.NewEngine()

	// Business hours and scheduling context
	now := time.Date(2024, 7, 10, 14, 30, 0, 0, time.UTC)         // Wednesday 2:30 PM UTC
	businessStart := time.Date(2024, 7, 10, 9, 0, 0, 0, time.UTC) // 9 AM
	businessEnd := time.Date(2024, 7, 10, 17, 0, 0, 0, time.UTC)  // 5 PM

	context := rule.D{
		"meeting": rule.D{
			"requested_start": "2024-07-10T15:00:00Z", // 3 PM
			"requested_end":   "2024-07-10T16:00:00Z", // 4 PM
			"duration":        60,                     // minutes
		},
		"business_hours": rule.D{
			"start": businessStart.Unix(),
			"end":   businessEnd.Unix(),
		},
		"employee": rule.D{
			"timezone":        "UTC",
			"available_from":  "2024-07-10T14:00:00Z", // 2 PM
			"available_until": "2024-07-10T18:00:00Z", // 6 PM
			"last_meeting":    "2024-07-10T13:30:00Z", // 1:30 PM
		},
		"current_time": now.Unix(),
	}

	// Scheduling validation rules
	schedulingRules := []struct {
		name        string
		rule        string
		description string
	}{
		{
			"Within Business Hours Start",
			`meeting.requested_start aq business_hours.start and meeting.requested_start be business_hours.end`,
			"Meeting start time must be within business hours",
		},
		{
			"Within Business Hours End",
			`meeting.requested_end aq business_hours.start and meeting.requested_end bq business_hours.end`,
			"Meeting end time must be within business hours",
		},
		{
			"Future Meeting",
			`meeting.requested_start af current_time`,
			"Meeting must be scheduled for the future",
		},
		{
			"Employee Available Start",
			`meeting.requested_start aq employee.available_from`,
			"Meeting start must be when employee is available",
		},
		{
			"Employee Available End",
			`meeting.requested_end bq employee.available_until`,
			"Meeting end must be before employee unavailable",
		},
		{
			"Buffer After Last Meeting",
			`meeting.requested_start af employee.last_meeting`,
			"New meeting must be after previous meeting ended",
		},
		{
			"Logical Meeting Duration",
			`meeting.requested_end af meeting.requested_start`,
			"Meeting end must be after meeting start",
		},
		{
			"Not Too Far Future",
			`meeting.requested_start be "2024-08-10T00:00:00Z"`,
			"Meeting must be within reasonable future timeframe",
		},
	}

	fmt.Println("Validating meeting schedule:")
	fmt.Println("----------------------------")
	fmt.Printf("Current time: %s\n", now.Format(time.RFC3339))
	fmt.Printf("Requested meeting: %s - %s\n",
		context["meeting"].(rule.D)["requested_start"],
		context["meeting"].(rule.D)["requested_end"])
	fmt.Printf("Business hours: %s - %s\n\n",
		businessStart.Format("15:04"), businessEnd.Format("15:04"))

	validationsPassed := 0
	totalValidations := len(schedulingRules)

	for _, sr := range schedulingRules {
		result, err := engine.Evaluate(sr.rule, context)
		if err != nil {
			log.Printf("❌ Error: %v", err)
			continue
		}

		status := "❌"
		if result {
			status = "✅"
			validationsPassed++
		}

		fmt.Printf("%s %s\n", status, sr.name)
		fmt.Printf("   Rule: %s\n", sr.rule)
		fmt.Printf("   Description: %s\n", sr.description)
		fmt.Printf("   Valid: %t\n\n", result)
	}

	fmt.Printf("📊 Scheduling Validation Summary:\n")
	fmt.Printf("   Passed: %d/%d validations (%.1f%%)\n",
		validationsPassed, totalValidations,
		float64(validationsPassed)/float64(totalValidations)*100)

	if validationsPassed == totalValidations {
		fmt.Println("🎉 Meeting can be scheduled!")
	} else {
		fmt.Println("⚠️  Meeting schedule has conflicts")
	}

	fmt.Println("\n✨ Scheduling validation completed!")
}
