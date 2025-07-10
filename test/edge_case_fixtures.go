package test

import (
	"math"
	"strings"

	"github.com/NSXBet/rule"
)

/* ---------- Edge Cases & Advanced Scenarios ---------- */

//nolint:gochecknoglobals // Test data
var EdgeCaseTests = []Case{
	// Deep nesting (3+ levels)
	{"deep_nested_eq", "a.b.c eq d.e.f", rule.D{
		"a": rule.D{"b": rule.D{"c": 42}},
		"d": rule.D{"e": rule.D{"f": 42}},
	}, true},
	{"deep_nested_ne", "a.b.c.d ne e.f.g.h", rule.D{
		"a": rule.D{"b": rule.D{"c": rule.D{"d": "test"}}},
		"e": rule.D{"f": rule.D{"g": rule.D{"h": "different"}}},
	}, true},

	// Mixed types comparison
	{"mixed_string_number", `version eq "2024"`, rule.D{
		"version": "2024",
	}, true},
	{"mixed_float_int", "price eq discount", rule.D{
		"price":    99.0,
		"discount": 99,
	}, true},

	// Zero values
	{"zero_int_eq", "x eq y", rule.D{"x": 0, "y": 0}, true},
	{"zero_float_eq", "x eq y", rule.D{"x": 0.0, "y": 0}, true},
	{"empty_string_eq", "x eq y", rule.D{"x": "", "y": ""}, true},
	{"zero_vs_empty", `x eq ""`, rule.D{"x": 0}, false},

	// Array edge cases
	{"empty_array_membership", "x in y", rule.D{
		"x": 1,
		"y": []any{},
	}, false},
	{"nested_array_membership", "user.id in allowed.users", rule.D{
		"user":    rule.D{"id": 123},
		"allowed": rule.D{"users": []any{123, 456, 789}},
	}, true},

	// Complex logical combinations
	{"complex_and_or", "(a.x eq b.y) and (c.z gt d.w) or (e eq f)", rule.D{
		"a": rule.D{"x": 10},
		"b": rule.D{"y": 10},
		"c": rule.D{"z": 20},
		"d": rule.D{"w": 15},
		"e": 999,
		"f": 999,
	}, true},
	{"nested_not_operations", "not (x.y eq z.w and a.b ne c.d)", rule.D{
		"x": rule.D{"y": 10},
		"z": rule.D{"w": 10},
		"a": rule.D{"b": "same"},
		"c": rule.D{"d": "same"},
	}, true},
}

/* ---------- String Operation Edge Cases ---------- */

//nolint:gochecknoglobals // Test data
var StringEdgeCaseTests = []Case{
	// String operations between properties
	{"prop_contains_nested", "user.email co domain.suffix", rule.D{
		"user":   rule.D{"email": "john@example.com"},
		"domain": rule.D{"suffix": "@example.com"},
	}, true},
	{"prop_startswith_nested", "file.name sw prefix.value", rule.D{
		"file":   rule.D{"name": "user_profile.jpg"},
		"prefix": rule.D{"value": "user_"},
	}, true},
	{"prop_endswith_nested", "document.path ew extension.type", rule.D{
		"document":  rule.D{"path": "report.pdf"},
		"extension": rule.D{"type": ".pdf"},
	}, true},

	// Empty string operations
	{"empty_string_contains", `x co ""`, rule.D{"x": "hello"}, true},
	{"empty_string_startswith", `x sw ""`, rule.D{"x": "hello"}, true},
	{"empty_string_endswith", `x ew ""`, rule.D{"x": "hello"}, true},
	{"contains_empty_string", `"" co "test"`, rule.D{}, false},

	// Case insensitive (compatible with nikunjy/rules)
	{"case_insensitive_eq", `x eq "test"`, rule.D{"x": "Test"}, true},
	{"case_insensitive_contains", `x co "hello"`, rule.D{"x": "Hello World"}, true},
}

/* ---------- Numeric Edge Cases ---------- */

//nolint:gochecknoglobals // Test data
var NumericEdgeCaseTests = []Case{
	// Large numbers
	{"large_int_eq", "x eq y", rule.D{
		"x": 9223372036854775807,
		"y": 9223372036854775807,
	}, true},
	{"large_float_eq", "x eq y", rule.D{
		"x": 1.7976931348623157e+308,
		"y": 1.7976931348623157e+308,
	}, true},

	// Negative numbers
	{"negative_comparison", "x lt y", rule.D{"x": -10, "y": -5}, true},
	{"negative_zero", "x eq y", rule.D{"x": math.Copysign(0, -1), "y": 0.0}, true},

	// Mixed int/float comparisons
	{"int_float_eq", "x eq y", rule.D{"x": 42, "y": 42.0}, true},
	{"int_float_gt", "x gt y", rule.D{"x": 43, "y": 42.5}, true},
	{"int_float_lt", "x lt y", rule.D{"x": 42, "y": 42.1}, true},
}

/* ---------- Array/Membership Edge Cases ---------- */

//nolint:gochecknoglobals // Test data
var ArrayEdgeCaseTests = []Case{
	// Mixed type arrays
	{"mixed_array_string", `x in ["1", "42", "3.14"]`, rule.D{
		"x": "42",
	}, true},
	{"mixed_array_number", "x in [1, 42, 3.14]", rule.D{
		"x": 42,
	}, true},
	{"mixed_array_no_match", `x in ["test", "42", "false"]`, rule.D{
		"x": 42, // int 42 != string "42"
	}, false},

	// Nested array membership
	{"nested_in_comparison", `user.role in permissions.allowed`, rule.D{
		"user":        rule.D{"role": "admin"},
		"permissions": rule.D{"allowed": []any{"admin", "moderator"}},
	}, true},
	{"deep_nested_in", "a.b.c in d.e.f", rule.D{
		"a": rule.D{"b": rule.D{"c": "value"}},
		"d": rule.D{"e": rule.D{"f": []any{"value", "other"}}},
	}, true},

	// Array with complex values
	{"array_with_floats", "x in [1.0, 2.5, 3.14, 4.8]", rule.D{
		"x": 3.14,
	}, true},
	{"array_with_versions", `x in ["1.0.0", "1.2.3", "2.0.0"]`, rule.D{
		"x": "1.2.3",
	}, true},
}

/* ---------- Presence (pr) Edge Cases ---------- */

//nolint:gochecknoglobals // Test data
var PresenceEdgeCaseTests = []Case{
	// Nested presence
	{"nested_presence_true", "user.profile pr", rule.D{
		"user": rule.D{"profile": rule.D{"name": "John"}},
	}, true},
	{"nested_presence_false", "user.profile pr", rule.D{
		"user": rule.D{"other": "value"},
	}, false},
	{"deep_nested_presence", "a.b.c.d pr", rule.D{
		"a": rule.D{
			"b": rule.D{
				"c": rule.D{"d": "exists"},
			},
		},
	}, true},

	// Presence with zero values
	{"presence_zero_int", "x pr", rule.D{"x": 0}, true},
	{"presence_empty_string", "x pr", rule.D{"x": ""}, true},
	{"presence_false_bool", "x pr", rule.D{"x": false}, true},
	{"presence_empty_array", "x pr", rule.D{"x": []any{}}, true},
}

/* ---------- Extreme Values ---------- */

//nolint:gochecknoglobals // Test data
var ExtremeValueTests = []Case{
	// Very long strings
	{"long_string_eq", `x eq "` + strings.Repeat("a", 1000) + `"`, rule.D{
		"x": strings.Repeat("a", 1000),
	}, true},
	{"long_string_contains", `x co "needle"`, rule.D{
		"x": strings.Repeat("a", 500) + "needle" + strings.Repeat("b", 500),
	}, true},

	// Very deep nesting
	{
		"extreme_deep_nesting",
		"level1.level2.level3.level4.level5.level6.level7.level8.level9.level10 eq 42",
		rule.D{
			"level1": rule.D{
				"level2": rule.D{
					"level3": rule.D{
						"level4": rule.D{
							"level5": rule.D{
								"level6": rule.D{
									"level7": rule.D{
										"level8": rule.D{
											"level9": rule.D{
												"level10": 42,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		true,
	},

	// Large arrays
	{"large_array_membership", "x in large_array", rule.D{
		"x": 999,
		"large_array": func() []any {
			arr := make([]any, 1000)
			for i := range 1000 {
				arr[i] = i
			}
			return arr
		}(),
	}, true},

	// Unicode and special characters
	{"unicode_strings", `emoji eq "🚀"`, rule.D{"emoji": "🚀"}, true},
	{"unicode_contains", `text co "🎉"`, rule.D{"text": "Hello 🎉 World"}, true},
	{"special_chars", `path eq "/path/to/file@domain.com"`, rule.D{
		"path": "/path/to/file@domain.com",
	}, true},

	// Datetime operators - RFC 3339 format
	{"datetime_equal_rfc3339", `created_at dq "2024-07-09T19:12:00-03:00"`, rule.D{
		"created_at": "2024-07-09T22:12:00Z", // Same time in UTC
	}, true},
	{"datetime_not_equal_rfc3339", `created_at dn "2024-07-09T19:12:00-03:00"`, rule.D{
		"created_at": "2024-07-09T22:13:00Z", // Different time
	}, true},
	{"datetime_before_rfc3339", `created_at be "2024-07-09T19:12:00-03:00"`, rule.D{
		"created_at": "2024-07-09T22:11:00Z", // 1 minute before
	}, true},
	{"datetime_before_or_equal_rfc3339", `created_at bq "2024-07-09T19:12:00-03:00"`, rule.D{
		"created_at": "2024-07-09T22:12:00Z", // Equal time
	}, true},
	{"datetime_after_rfc3339", `created_at af "2024-07-09T19:12:00-03:00"`, rule.D{
		"created_at": "2024-07-09T22:13:00Z", // 1 minute after
	}, true},
	{"datetime_after_or_equal_rfc3339", `created_at aq "2024-07-09T19:12:00-03:00"`, rule.D{
		"created_at": "2024-07-09T22:12:00Z", // Equal time
	}, true},

	// Datetime operators - Unix timestamp format
	{"datetime_equal_unix", `timestamp dq 1720558320`, rule.D{
		"timestamp": int64(1720558320),
	}, true},
	{"datetime_not_equal_unix", `timestamp dn 1720558320`, rule.D{
		"timestamp": int64(1720558321),
	}, true},
	{"datetime_before_unix", `timestamp be 1720558320`, rule.D{
		"timestamp": int64(1720558319),
	}, true},
	{"datetime_before_or_equal_unix", `timestamp bq 1720558320`, rule.D{
		"timestamp": int64(1720558320),
	}, true},
	{"datetime_after_unix", `timestamp af 1720558320`, rule.D{
		"timestamp": int64(1720558321),
	}, true},
	{"datetime_after_or_equal_unix", `timestamp aq 1720558320`, rule.D{
		"timestamp": int64(1720558320),
	}, true},

	// Mixed format comparisons (RFC3339 vs Unix)
	{"datetime_mixed_rfc3339_vs_unix", `created_at af 1720558320`, rule.D{
		"created_at": "2024-07-09T22:12:01Z", // 1 second after the Unix timestamp
	}, true},
	{"datetime_mixed_unix_vs_rfc3339", `timestamp be "2024-07-09T22:12:01Z"`, rule.D{
		"timestamp": int64(1720558320), // 1 second before the RFC3339 time
	}, true},

	// Edge cases
	{"datetime_false_before", `created_at be "2024-07-09T19:12:00-03:00"`, rule.D{
		"created_at": "2024-07-09T22:13:00Z", // After, so before should be false
	}, false},
	{"datetime_false_after", `created_at af "2024-07-09T19:12:00-03:00"`, rule.D{
		"created_at": "2024-07-09T22:11:00Z", // Before, so after should be false
	}, false},
	{"datetime_false_equal", `created_at dq "2024-07-09T19:12:00-03:00"`, rule.D{
		"created_at": "2024-07-09T22:13:00Z", // Different time
	}, false},
}
