package test

import "strings"

/* ---------- Edge Cases & Advanced Scenarios ---------- */

var EdgeCaseTests = []TestCase{
	// Deep nesting (3+ levels)
	{"deep_nested_eq", "a.b.c eq d.e.f", map[string]any{
		"a": map[string]any{"b": map[string]any{"c": 42}},
		"d": map[string]any{"e": map[string]any{"f": 42}},
	}, true},
	{"deep_nested_ne", "a.b.c.d ne e.f.g.h", map[string]any{
		"a": map[string]any{"b": map[string]any{"c": map[string]any{"d": "test"}}},
		"e": map[string]any{"f": map[string]any{"g": map[string]any{"h": "different"}}},
	}, true},

	// Mixed types comparison
	{"mixed_string_number", `version eq "2024"`, map[string]any{
		"version": "2024",
	}, true},
	{"mixed_float_int", "price eq discount", map[string]any{
		"price":    99.0,
		"discount": 99,
	}, true},

	// Zero values
	{"zero_int_eq", "x eq y", map[string]any{"x": 0, "y": 0}, true},
	{"zero_float_eq", "x eq y", map[string]any{"x": 0.0, "y": 0}, true},
	{"empty_string_eq", "x eq y", map[string]any{"x": "", "y": ""}, true},
	{"zero_vs_empty", `x eq ""`, map[string]any{"x": 0}, false},

	// Array edge cases
	{"empty_array_membership", "x in y", map[string]any{
		"x": 1,
		"y": []any{},
	}, false},
	{"nested_array_membership", "user.id in allowed.users", map[string]any{
		"user":    map[string]any{"id": 123},
		"allowed": map[string]any{"users": []any{123, 456, 789}},
	}, true},

	// Complex logical combinations
	{"complex_and_or", "(a.x eq b.y) and (c.z gt d.w) or (e eq f)", map[string]any{
		"a": map[string]any{"x": 10},
		"b": map[string]any{"y": 10},
		"c": map[string]any{"z": 20},
		"d": map[string]any{"w": 15},
		"e": 999,
		"f": 999,
	}, true},
	{"nested_not_operations", "not (x.y eq z.w and a.b ne c.d)", map[string]any{
		"x": map[string]any{"y": 10},
		"z": map[string]any{"w": 10},
		"a": map[string]any{"b": "same"},
		"c": map[string]any{"d": "same"},
	}, true},
}

/* ---------- String Operation Edge Cases ---------- */

var StringEdgeCaseTests = []TestCase{
	// String operations between properties
	{"prop_contains_nested", "user.email co domain.suffix", map[string]any{
		"user":   map[string]any{"email": "john@example.com"},
		"domain": map[string]any{"suffix": "@example.com"},
	}, true},
	{"prop_startswith_nested", "file.name sw prefix.value", map[string]any{
		"file":   map[string]any{"name": "user_profile.jpg"},
		"prefix": map[string]any{"value": "user_"},
	}, true},
	{"prop_endswith_nested", "document.path ew extension.type", map[string]any{
		"document":  map[string]any{"path": "report.pdf"},
		"extension": map[string]any{"type": ".pdf"},
	}, true},

	// Empty string operations
	{"empty_string_contains", `x co ""`, map[string]any{"x": "hello"}, true},
	{"empty_string_startswith", `x sw ""`, map[string]any{"x": "hello"}, true},
	{"empty_string_endswith", `x ew ""`, map[string]any{"x": "hello"}, true},
	{"contains_empty_string", `"" co "test"`, map[string]any{}, false},

	// Case sensitivity
	{"case_sensitive_eq", `x eq "test"`, map[string]any{"x": "Test"}, false},
	{"case_sensitive_contains", `x co "hello"`, map[string]any{"x": "Hello World"}, false},
}

/* ---------- Numeric Edge Cases ---------- */

var NumericEdgeCaseTests = []TestCase{
	// Large numbers
	{"large_int_eq", "x eq y", map[string]any{
		"x": 9223372036854775807,
		"y": 9223372036854775807,
	}, true},
	{"large_float_eq", "x eq y", map[string]any{
		"x": 1.7976931348623157e+308,
		"y": 1.7976931348623157e+308,
	}, true},

	// Negative numbers
	{"negative_comparison", "x lt y", map[string]any{"x": -10, "y": -5}, true},
	{"negative_zero", "x eq y", map[string]any{"x": -0.0, "y": 0.0}, true},

	// Mixed int/float comparisons
	{"int_float_eq", "x eq y", map[string]any{"x": 42, "y": 42.0}, true},
	{"int_float_gt", "x gt y", map[string]any{"x": 43, "y": 42.5}, true},
	{"int_float_lt", "x lt y", map[string]any{"x": 42, "y": 42.1}, true},
}

/* ---------- Array/Membership Edge Cases ---------- */

var ArrayEdgeCaseTests = []TestCase{
	// Mixed type arrays
	{"mixed_array_string", `x in ["1", "42", "3.14"]`, map[string]any{
		"x": "42",
	}, true},
	{"mixed_array_number", "x in [1, 42, 3.14]", map[string]any{
		"x": 42,
	}, true},
	{"mixed_array_no_match", `x in ["test", "42", "false"]`, map[string]any{
		"x": 42, // int 42 != string "42"
	}, false},

	// Nested array membership
	{"nested_in_comparison", `user.role in permissions.allowed`, map[string]any{
		"user":        map[string]any{"role": "admin"},
		"permissions": map[string]any{"allowed": []any{"admin", "moderator"}},
	}, true},
	{"deep_nested_in", "a.b.c in d.e.f", map[string]any{
		"a": map[string]any{"b": map[string]any{"c": "value"}},
		"d": map[string]any{"e": map[string]any{"f": []any{"value", "other"}}},
	}, true},

	// Array with complex values
	{"array_with_floats", "x in [1.0, 2.5, 3.14, 4.8]", map[string]any{
		"x": 3.14,
	}, true},
	{"array_with_versions", `x in ["1.0.0", "1.2.3", "2.0.0"]`, map[string]any{
		"x": "1.2.3",
	}, true},
}

/* ---------- Presence (pr) Edge Cases ---------- */

var PresenceEdgeCaseTests = []TestCase{
	// Nested presence
	{"nested_presence_true", "user.profile pr", map[string]any{
		"user": map[string]any{"profile": map[string]any{"name": "John"}},
	}, true},
	{"nested_presence_false", "user.profile pr", map[string]any{
		"user": map[string]any{"other": "value"},
	}, false},
	{"deep_nested_presence", "a.b.c.d pr", map[string]any{
		"a": map[string]any{
			"b": map[string]any{
				"c": map[string]any{"d": "exists"},
			},
		},
	}, true},

	// Presence with zero values
	{"presence_zero_int", "x pr", map[string]any{"x": 0}, true},
	{"presence_empty_string", "x pr", map[string]any{"x": ""}, true},
	{"presence_false_bool", "x pr", map[string]any{"x": false}, true},
	{"presence_empty_array", "x pr", map[string]any{"x": []any{}}, true},
}

/* ---------- Extreme Values ---------- */

var ExtremeValueTests = []TestCase{
	// Very long strings
	{"long_string_eq", `x eq "` + strings.Repeat("a", 1000) + `"`, map[string]any{
		"x": strings.Repeat("a", 1000),
	}, true},
	{"long_string_contains", `x co "needle"`, map[string]any{
		"x": strings.Repeat("a", 500) + "needle" + strings.Repeat("b", 500),
	}, true},
	
	// Very deep nesting
	{"extreme_deep_nesting", "level1.level2.level3.level4.level5.level6.level7.level8.level9.level10 eq 42", map[string]any{
		"level1": map[string]any{
			"level2": map[string]any{
				"level3": map[string]any{
					"level4": map[string]any{
						"level5": map[string]any{
							"level6": map[string]any{
								"level7": map[string]any{
									"level8": map[string]any{
										"level9": map[string]any{
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
	}, true},
	
	// Large arrays
	{"large_array_membership", "x in large_array", map[string]any{
		"x": 999,
		"large_array": func() []any {
			arr := make([]any, 1000)
			for i := 0; i < 1000; i++ {
				arr[i] = i
			}
			return arr
		}(),
	}, true},
	
	// Unicode and special characters
	{"unicode_strings", `emoji eq "🚀"`, map[string]any{"emoji": "🚀"}, true},
	{"unicode_contains", `text co "🎉"`, map[string]any{"text": "Hello 🎉 World"}, true},
	{"special_chars", `path eq "/path/to/file@domain.com"`, map[string]any{
		"path": "/path/to/file@domain.com",
	}, true},
}