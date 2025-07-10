package rule

import (
	"fmt"
	"strings"
	"testing"
)

// FuzzRuleExecution tests rule execution with random strings to find parsing crashes.
func FuzzRuleExecution(f *testing.F) {
	// Seed corpus with known good rules
	f.Add("age gt 18")
	f.Add("name eq \"John\"")
	f.Add("score in [100, 200, 300]")
	f.Add("user.active eq true")
	f.Add("created_at af \"2024-01-01T00:00:00Z\"")
	f.Add("(a and b) or (c and d)")
	f.Add("not (status eq \"inactive\")")
	f.Add("price lt 99.99")
	f.Add("tags co \"important\"")
	f.Add("email ew \".com\"")

	// Add some invalid rules to test error handling
	f.Add("age gt")
	f.Add("name eq")
	f.Add("((((")
	f.Add("and or not")
	f.Add("")

	engine := NewEngine()
	context := D{
		"age":    25,
		"name":   "John",
		"score":  150,
		"price":  50.0,
		"tags":   "important work",
		"email":  "test@example.com",
		"status": "active",
		"user": D{
			"active": true,
		},
		"created_at": "2024-06-01T10:00:00Z",
		"a":          true,
		"b":          false,
		"c":          true,
		"d":          false,
	}

	f.Fuzz(func(t *testing.T, rule string) {
		// Fuzz target should never panic, regardless of input
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Rule engine panicked with rule %q: %v", rule, r)
			}
		}()

		// Test evaluation - should either succeed or return an error gracefully
		result, err := engine.Evaluate(rule, context)
		if err != nil {
			// Error is acceptable, but we should be able to convert it to string
			_ = err.Error()
		} else {
			// Result should be a valid boolean
			_ = result
		}

		// Test AddQuery - should either succeed or return an error gracefully
		addErr := engine.AddQuery(rule)
		if addErr != nil {
			// Error is acceptable, but we should be able to convert it to string
			_ = addErr.Error()
		}
	})
}

// FuzzStringOperations tests string operations with random string inputs.
func FuzzStringOperations(f *testing.F) {
	// Seed corpus with various string patterns
	f.Add("hello world", "hello")
	f.Add("TEST STRING", "test")
	f.Add("unicode: café", "café")
	f.Add("line1\nline2", "\n")
	f.Add("tab\there", "\t")
	f.Add("", "")
	f.Add("a", "b")
	f.Add("null\x00byte", "\x00")
	f.Add("emoji 🚀 test", "🚀")
	f.Add("very long string with many words and characters", "many")

	engine := NewEngine()

	f.Fuzz(func(t *testing.T, str1, str2 string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("String operations panicked with inputs %q, %q: %v", str1, str2, r)
			}
		}()

		context := D{
			"str1": str1,
			"str2": str2,
		}

		// Test all string operations
		stringOps := []string{
			"str1 eq str2",
			"str1 ne str2",
			"str1 co str2",
			"str1 sw str2",
			"str1 ew str2",
			"str2 co str1",
			"str2 sw str1", 
			"str2 ew str1",
		}

		for _, rule := range stringOps {
			_, err := engine.Evaluate(rule, context)
			if err != nil {
				// Errors are acceptable for malformed rules
				_ = err.Error()
			}
		}
	})
}

// FuzzNumericOperations tests numeric operations with random numeric inputs.
func FuzzNumericOperations(f *testing.F) {
	// Seed corpus with various numeric edge cases
	f.Add(int64(0), float64(0))
	f.Add(int64(1), float64(1))
	f.Add(int64(-1), float64(-1))
	f.Add(int64(9223372036854775807), float64(9223372036854775807)) // max int64
	f.Add(int64(-9223372036854775808), float64(-9223372036854775808)) // min int64
	f.Add(int64(42), float64(42.0))
	f.Add(int64(100), float64(99.999999))

	engine := NewEngine()

	f.Fuzz(func(t *testing.T, intVal int64, floatVal float64) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Numeric operations panicked with inputs %d, %f: %v", intVal, floatVal, r)
			}
		}()

		context := D{
			"intVal":   intVal,
			"floatVal": floatVal,
		}

		// Test all numeric operations
		numericOps := []string{
			"intVal eq floatVal",
			"intVal ne floatVal",
			"intVal lt floatVal",
			"intVal gt floatVal", 
			"intVal le floatVal",
			"intVal ge floatVal",
			"floatVal eq intVal",
			"floatVal lt intVal",
			"floatVal gt intVal",
		}

		for _, rule := range numericOps {
			_, err := engine.Evaluate(rule, context)
			if err != nil {
				// Errors should be rare for numeric operations but could happen
				_ = err.Error()
			}
		}
	})
}

// FuzzDateTimeOperations tests datetime operations with random datetime strings.
func FuzzDateTimeOperations(f *testing.F) {
	// Seed corpus with various datetime formats
	f.Add("2024-01-01T00:00:00Z")
	f.Add("2024-12-31T23:59:59Z")
	f.Add("2024-06-15T12:30:45-07:00")
	f.Add("2024-02-29T10:15:30+05:30") // leap year
	f.Add("1970-01-01T00:00:00Z")      // Unix epoch
	f.Add("2038-01-19T03:14:07Z")      // Unix timestamp limit
	f.Add("invalid-date")
	f.Add("2024-13-40T25:70:70Z") // invalid components
	f.Add("")
	f.Add("1234567890") // Unix timestamp
	f.Add("0")

	engine := NewEngine()

	f.Fuzz(func(t *testing.T, dateStr string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("DateTime operations panicked with input %q: %v", dateStr, r)
			}
		}()

		context := D{
			"date1": dateStr,
			"date2": "2024-01-01T00:00:00Z", // Fixed reference date
		}

		// Test all datetime operations
		dateTimeOps := []string{
			"date1 dq date2",
			"date1 dn date2", 
			"date1 be date2",
			"date1 bq date2",
			"date1 af date2",
			"date1 aq date2",
		}

		for _, rule := range dateTimeOps {
			_, err := engine.Evaluate(rule, context)
			if err != nil {
				// Errors are expected for invalid datetime formats
				_ = err.Error()
			}
		}
	})
}

// FuzzPropertyAccess tests property access with random property paths.
func FuzzPropertyAccess(f *testing.F) {
	// Seed corpus with various property access patterns
	f.Add("user.name")
	f.Add("config.settings.theme")
	f.Add("a.b.c.d.e.f.g")
	f.Add("user")
	f.Add("")
	f.Add("nonexistent.property")
	f.Add("user..name") // double dots
	f.Add("user.") // trailing dot
	f.Add(".user") // leading dot
	f.Add("user.name.length") // accessing property on non-object
	f.Add("very_long_property_name_with_many_characters")

	engine := NewEngine()
	
	f.Fuzz(func(t *testing.T, propertyPath string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Property access panicked with path %q: %v", propertyPath, r)
			}
		}()

		context := D{
			"user": D{
				"name": "John",
				"profile": D{
					"age": 25,
				},
			},
			"config": D{
				"settings": D{
					"theme": "dark",
				},
			},
			"simple": "value",
		}

		// Test property access in various rule contexts
		if propertyPath != "" && !strings.Contains(propertyPath, "..") {
			rules := []string{
				fmt.Sprintf("%s pr", propertyPath),
				fmt.Sprintf("%s eq \"test\"", propertyPath),
				fmt.Sprintf("%s ne \"test\"", propertyPath),
			}

			for _, rule := range rules {
				_, err := engine.Evaluate(rule, context)
				if err != nil {
					// Errors are acceptable for invalid property paths
					_ = err.Error()
				}
			}
		}
	})
}

// FuzzArrayOperations tests array operations with random array inputs.
func FuzzArrayOperations(f *testing.F) {
	// Seed corpus with various array scenarios
	f.Add("test", `["test", "value"]`)
	f.Add("42", `[1, 2, 42, 3]`)
	f.Add("missing", `["a", "b", "c"]`)
	f.Add("", `[]`)
	f.Add("null", `[null, "test"]`)
	f.Add("true", `[true, false]`)

	engine := NewEngine()

	f.Fuzz(func(t *testing.T, searchValue, arrayStr string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Array operations panicked with inputs %q, %q: %v", searchValue, arrayStr, r)
			}
		}()

		// Create context with the array and search value
		context := D{
			"searchValue": searchValue,
			"arrayValue":  arrayStr,
		}

		// Test various array operation rules
		arrayRules := []string{
			"searchValue in [\"a\", \"b\", \"c\"]",
			"\"fixed\" in arrayValue", // This will likely error, which is fine
			"searchValue in []",
		}

		for _, rule := range arrayRules {
			_, err := engine.Evaluate(rule, context)
			if err != nil {
				// Errors are acceptable for malformed array syntax
				_ = err.Error()
			}
		}
	})
}

// FuzzBooleanOperations tests boolean operations with random boolean inputs.
func FuzzBooleanOperations(f *testing.F) {
	// Seed corpus with boolean combinations
	f.Add(true, false)
	f.Add(false, true)
	f.Add(true, true)
	f.Add(false, false)

	engine := NewEngine()

	f.Fuzz(func(t *testing.T, bool1, bool2 bool) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Boolean operations panicked with inputs %t, %t: %v", bool1, bool2, r)
			}
		}()

		context := D{
			"bool1": bool1,
			"bool2": bool2,
		}

		// Test all boolean operations
		booleanOps := []string{
			"bool1 and bool2",
			"bool1 or bool2",
			"not bool1",
			"not bool2",
			"bool1 eq bool2",
			"bool1 ne bool2",
			"(bool1 and bool2) or (not bool1)",
			"not (bool1 or bool2)",
		}

		for _, rule := range booleanOps {
			_, err := engine.Evaluate(rule, context)
			if err != nil {
				// Boolean operations should rarely error
				t.Errorf("Unexpected error in boolean operation %q: %v", rule, err)
			}
		}
	})
}

// FuzzComplexRules tests complex rule combinations with random inputs.
func FuzzComplexRules(f *testing.F) {
	// Seed corpus with complex rule patterns
	f.Add("(age gt 18 and status eq \"active\") or admin eq true")
	f.Add("name co \"test\" and (score in [100, 200] or level gt 5)")
	f.Add("not (expired eq true) and date af \"2024-01-01T00:00:00Z\"")
	f.Add("((a and b) or c) and not d")
	f.Add("user.profile.age ge config.limits.minimum")

	engine := NewEngine()

	f.Fuzz(func(t *testing.T, complexRule string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Complex rule evaluation panicked with rule %q: %v", complexRule, r)
			}
		}()

		context := D{
			"age":     25,
			"status":  "active",
			"admin":   false,
			"name":    "test user",
			"score":   150,
			"level":   3,
			"expired": false,
			"date":    "2024-06-01T00:00:00Z",
			"a":       true,
			"b":       false,
			"c":       true,
			"d":       false,
			"user": D{
				"profile": D{
					"age": 30,
				},
			},
			"config": D{
				"limits": D{
					"minimum": 18,
				},
			},
		}

		// Evaluate the complex rule
		_, err := engine.Evaluate(complexRule, context)
		if err != nil {
			// Complex rules might have syntax errors, which is acceptable
			_ = err.Error()
		}
	})
}

// FuzzMixedTypeComparisons tests comparisons between different data types.
func FuzzMixedTypeComparisons(f *testing.F) {
	// Seed corpus with mixed type scenarios
	f.Add("42", int64(42), 42.0, true)
	f.Add("true", int64(1), 1.0, true)
	f.Add("", int64(0), 0.0, false)
	f.Add("null", int64(-1), -1.0, false)

	engine := NewEngine()

	f.Fuzz(func(t *testing.T, strVal string, intVal int64, floatVal float64, boolVal bool) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Mixed type comparison panicked with inputs %q, %d, %f, %t: %v", 
					strVal, intVal, floatVal, boolVal, r)
			}
		}()

		context := D{
			"strVal":   strVal,
			"intVal":   intVal,
			"floatVal": floatVal,
			"boolVal":  boolVal,
		}

		// Test mixed type comparisons
		mixedTypeRules := []string{
			"strVal eq intVal",    // string vs int
			"strVal eq floatVal",  // string vs float
			"strVal eq boolVal",   // string vs bool
			"intVal eq floatVal",  // int vs float (should work)
			"intVal eq boolVal",   // int vs bool
			"floatVal eq boolVal", // float vs bool
		}

		for _, rule := range mixedTypeRules {
			_, err := engine.Evaluate(rule, context)
			if err != nil {
				// Mixed type comparisons might error or return false
				_ = err.Error()
			}
		}
	})
}