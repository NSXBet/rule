package test

import "github.com/NSXBet/rule"

// AdditionalEdgeCaseTests contains additional edge case tests for comprehensive coverage.
//
//nolint:gochecknoglobals // Test data
var AdditionalEdgeCaseTests = []Case{
	// String edge cases (focusing on supported features)
	{"quote_in_string", `name eq "John \"The Great\" Doe"`, rule.D{"name": `John "The Great" Doe`}, true},
	{"single_quote_in_string", `name eq "John's book"`, rule.D{"name": "John's book"}, true},
	{"string_with_spaces", `text eq "hello world"`, rule.D{"text": "hello world"}, true},
	{"string_with_numbers", `code eq "ABC123"`, rule.D{"code": "ABC123"}, true},
	{"escaped_backslash", `path eq "C:\\Users\\John"`, rule.D{"path": `C:\Users\John`}, true},
	{"escaped_newline", `text eq "Line 1\nLine 2"`, rule.D{"text": "Line 1\nLine 2"}, true},
	{"escaped_tab", `text eq "Column1\tColumn2"`, rule.D{"text": "Column1\tColumn2"}, true},

	// Boolean string representations
	{"bool_true_caps", `active eq "TRUE"`, rule.D{"active": "TRUE"}, true},
	{"bool_false_caps", `active eq "FALSE"`, rule.D{"active": "FALSE"}, true},
	{"bool_mixed_case", `active eq "True"`, rule.D{"active": "True"}, true},

	// Unicode edge cases
	{"unicode_emoji", `message co "🚀"`, rule.D{"message": "Launch 🚀 successful"}, true},
	{"unicode_accents", `name eq "José"`, rule.D{"name": "José"}, true},
	{"unicode_chinese", `text co "你好"`, rule.D{"text": "Hello 你好 World"}, true},

	// Mathematical edge cases (no scientific notation)
	{"very_large_float", `x lt 999999999999999999.0`, rule.D{"x": 123456789.0}, true},
	{"very_small_float", `x gt 0.000000000000001`, rule.D{"x": 0.001}, true},
	{"negative_zero", `x eq -0.0`, rule.D{"x": 0.0}, true},

	// Complex nested array operations
	{"nested_array_deep", `user.permissions.roles in ["admin", "user"]`, rule.D{
		"user": rule.D{
			"permissions": rule.D{
				"roles": "admin",
			},
		},
	}, true},

	// Performance edge cases
	{"very_long_string", `text co "needle"`, rule.D{
		"text": "This is a very long string that contains the word needle somewhere in the middle of all this text that goes on and on and on and on and on and on and on and on and on and on and on and on and on and on and on and on and on and on and on and on and on and on and on and on and on and on and on and on",
	}, true},

	// Empty string edge cases
	{"empty_string_contains_empty", `x co ""`, rule.D{"x": "hello"}, true},
	{"empty_string_starts_with_empty", `x sw ""`, rule.D{"x": "hello"}, true},
	{"empty_string_ends_with_empty", `x ew ""`, rule.D{"x": "hello"}, true},

	// Zero and null edge cases
	{"zero_string_vs_number", `x eq "0"`, rule.D{"x": 0}, false},        // Different types
	{"false_string_vs_bool", `x eq "false"`, rule.D{"x": false}, false}, // Different types
	{"null_string_literal", `x eq "null"`, rule.D{"x": "null"}, true},
}
