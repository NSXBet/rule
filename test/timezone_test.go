package test

import (
	"testing"
	"time"

	"github.com/NSXBet/rule"
	ruleslib "github.com/nikunjy/rules"
)

func TestTimezoneHandling(t *testing.T) {
	t.Log("🌍 TIMEZONE HANDLING COMPARISON")
	t.Log("===============================")

	// Same time in different timezones
	utcTime := time.Date(2024, 7, 10, 15, 30, 0, 0, time.UTC)
	estTime := time.Date(2024, 7, 10, 11, 30, 0, 0, time.FixedZone("EST", -4*3600))
	pstTime := time.Date(2024, 7, 10, 8, 30, 0, 0, time.FixedZone("PST", -7*3600))

	t.Logf("UTC time: %s", utcTime.Format(time.RFC3339))
	t.Logf("EST time: %s", estTime.Format(time.RFC3339))
	t.Logf("PST time: %s", pstTime.Format(time.RFC3339))
	t.Logf("Are they equal? UTC==EST: %t, UTC==PST: %t", utcTime.Equal(estTime), utcTime.Equal(pstTime))

	timezoneTests := []struct {
		name        string
		rule        string
		context     map[string]interface{}
		ourContext  rule.D
		description string
	}{
		{
			name: "UTC vs EST time.Time equality",
			rule: `utc_time eq est_time`,
			context: map[string]interface{}{
				"utc_time": utcTime,
				"est_time": estTime,
			},
			ourContext: rule.D{
				"utc_time": utcTime,
				"est_time": estTime,
			},
			description: "Same instant in different timezones should be equal",
		},
		{
			name: "time.Time vs RFC3339 string different timezone",
			rule: `time_obj eq "2024-07-10T11:30:00-04:00"`, // EST equivalent
			context: map[string]interface{}{
				"time_obj": utcTime,
			},
			ourContext: rule.D{
				"time_obj": utcTime,
			},
			description: "time.Time vs equivalent RFC3339 in different timezone",
		},
		{
			name: "datetime operator UTC vs EST",
			rule: `utc_time dq est_time`,
			context: map[string]interface{}{
				"utc_time": utcTime,
				"est_time": estTime,
			},
			ourContext: rule.D{
				"utc_time": utcTime,
				"est_time": estTime,
			},
			description: "Datetime operator with same instant different timezones",
		},
		{
			name: "datetime operator time.Time vs RFC3339 different TZ",
			rule: `time_obj dq "2024-07-10T11:30:00-04:00"`,
			context: map[string]interface{}{
				"time_obj": utcTime,
			},
			ourContext: rule.D{
				"time_obj": utcTime,
			},
			description: "Datetime operator: time.Time vs RFC3339 different timezone",
		},
	}

	for i, test := range timezoneTests {
		t.Logf("\n%d. %s", i+1, test.name)
		t.Logf("Rule: %s", test.rule)
		t.Logf("Description: %s", test.description)

		// Test nikunjy/rules (only for non-datetime operators)
		if !containsDatetimeOperator(test.rule) {
			rulesResult, rulesErr := ruleslib.Evaluate(test.rule, test.context)
			t.Logf("nikunjy/rules: %v (err: %v)", rulesResult, rulesErr)
		} else {
			t.Logf("nikunjy/rules: N/A (doesn't support datetime operators)")
		}

		// Test our library
		ourEngine := rule.NewEngine()
		ourResult, ourErr := ourEngine.Evaluate(test.rule, test.ourContext)
		t.Logf("Our library: %v (err: %v)", ourResult, ourErr)
	}

	t.Log("\n🎯 Timezone handling testing completed!")
}

func containsDatetimeOperator(rule string) bool {
	datetimeOps := []string{"dq", "dn", "be", "bq", "af", "aq"}
	for _, op := range datetimeOps {
		if len(rule) >= len(op) {
			for i := 0; i <= len(rule)-len(op); i++ {
				if i > 0 && rule[i-1] != ' ' {
					continue
				}
				if i+len(op) < len(rule) && rule[i+len(op)] != ' ' {
					continue
				}
				if rule[i:i+len(op)] == op {
					return true
				}
			}
		}
	}
	return false
}
