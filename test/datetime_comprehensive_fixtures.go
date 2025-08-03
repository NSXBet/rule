package test

import "github.com/NSXBet/rule"

// DateTimeComprehensiveTests contains exhaustive datetime operator test cases.
//
//nolint:gochecknoglobals // Test data
var DateTimeComprehensiveTests = []Case{
	// Property vs Property - Unix timestamp vs Unix timestamp
	{"datetime_prop_unix_vs_unix_after", `start_time af end_time`, rule.D{
		"start_time": int64(1720558321),
		"end_time":   int64(1720558320),
	}, true},
	{"datetime_prop_unix_vs_unix_before", `start_time be end_time`, rule.D{
		"start_time": int64(1720558319),
		"end_time":   int64(1720558320),
	}, true},
	{"datetime_prop_unix_vs_unix_equal", `start_time dq end_time`, rule.D{
		"start_time": int64(1720558320),
		"end_time":   int64(1720558320),
	}, true},
	{"datetime_prop_unix_vs_unix_not_equal", `start_time dn end_time`, rule.D{
		"start_time": int64(1720558321),
		"end_time":   int64(1720558320),
	}, true},
	{"datetime_prop_unix_vs_unix_before_or_equal", `start_time bq end_time`, rule.D{
		"start_time": int64(1720558320),
		"end_time":   int64(1720558320),
	}, true},
	{"datetime_prop_unix_vs_unix_after_or_equal", `start_time aq end_time`, rule.D{
		"start_time": int64(1720558320),
		"end_time":   int64(1720558320),
	}, true},

	// Property vs Property - RFC 3339 vs RFC 3339
	{"datetime_prop_rfc_vs_rfc_after", `created_at af updated_at`, rule.D{
		"created_at": "2024-07-09T22:12:01Z",
		"updated_at": "2024-07-09T22:12:00Z",
	}, true},
	{"datetime_prop_rfc_vs_rfc_before", `created_at be updated_at`, rule.D{
		"created_at": "2024-07-09T22:11:59Z",
		"updated_at": "2024-07-09T22:12:00Z",
	}, true},
	{"datetime_prop_rfc_vs_rfc_equal", `created_at dq updated_at`, rule.D{
		"created_at": "2024-07-09T22:12:00Z",
		"updated_at": "2024-07-09T22:12:00Z",
	}, true},
	{"datetime_prop_rfc_vs_rfc_timezone_equal", `created_at dq updated_at`, rule.D{
		"created_at": "2024-07-09T19:12:00-03:00", // Same time in different timezone
		"updated_at": "2024-07-09T22:12:00Z",
	}, true},

	// Property vs Property - Mixed formats
	{"datetime_prop_mixed_rfc_vs_unix", `created_at af timestamp`, rule.D{
		"created_at": "2024-07-09T22:12:01Z",
		"timestamp":  int64(1720558320), // 1 second earlier
	}, true},
	{"datetime_prop_mixed_unix_vs_rfc", `timestamp be created_at`, rule.D{
		"timestamp":  int64(1720558319),
		"created_at": "2024-07-09T22:12:00Z",
	}, true},

	// Nested Property vs Literal
	{"datetime_nested_prop_vs_literal_rfc", `event.created_at af "2024-07-09T22:12:00Z"`, rule.D{
		"event": rule.D{
			"created_at": "2024-07-09T22:12:01Z",
		},
	}, true},
	{"datetime_nested_prop_vs_literal_unix", `session.timestamp be 1720558320`, rule.D{
		"session": rule.D{
			"timestamp": int64(1720558319),
		},
	}, true},
	{"datetime_deep_nested_prop", `user.profile.last_login dq "2024-07-09T22:12:00Z"`, rule.D{
		"user": rule.D{
			"profile": rule.D{
				"last_login": "2024-07-09T22:12:00Z",
			},
		},
	}, true},

	// Nested Property vs Nested Property
	{"datetime_nested_vs_nested_same_depth", `event.start_time be event.end_time`, rule.D{
		"event": rule.D{
			"start_time": "2024-07-09T22:11:59Z",
			"end_time":   "2024-07-09T22:12:00Z",
		},
	}, true},
	{"datetime_nested_vs_nested_different_depth", `session.created_at af user.last_login`, rule.D{
		"session": rule.D{
			"created_at": "2024-07-09T22:12:01Z",
		},
		"user": rule.D{
			"last_login": "2024-07-09T22:12:00Z",
		},
	}, true},
	{"datetime_deep_nested_vs_deep_nested", `user.profile.created_at bq admin.profile.last_update`, rule.D{
		"user": rule.D{
			"profile": rule.D{
				"created_at": "2024-07-09T22:12:00Z",
			},
		},
		"admin": rule.D{
			"profile": rule.D{
				"last_update": "2024-07-09T22:12:00Z",
			},
		},
	}, true},

	// Different Numeric Types for Unix Timestamps
	{"datetime_int32_timestamp", `timestamp af 1720558320`, rule.D{
		"timestamp": int32(1720558321),
	}, true},
	{"datetime_float32_timestamp", `timestamp be 1000000`, rule.D{
		"timestamp": float32(999999.0),
	}, true},
	{"datetime_float64_timestamp", `timestamp dq 1720558320`, rule.D{
		"timestamp": float64(1720558320.0),
	}, true},
	{"datetime_uint64_timestamp", `timestamp aq 1720558320`, rule.D{
		"timestamp": uint64(1720558320),
	}, true},

	// String Unix Timestamps
	{"datetime_string_unix_timestamp", `timestamp af "1720558320"`, rule.D{
		"timestamp": "1720558321",
	}, true},
	{"datetime_mixed_string_int_unix", `timestamp_str be timestamp_int`, rule.D{
		"timestamp_str": "1720558319",
		"timestamp_int": int64(1720558320),
	}, true},

	// Boundary Conditions
	{"datetime_epoch_zero", `timestamp af 0`, rule.D{
		"timestamp": int64(1),
	}, true},
	{"datetime_negative_timestamp", `timestamp af -86400`, rule.D{
		"timestamp": int64(-86399), // 1 second after
	}, true},
	{"datetime_year_2038_problem", `timestamp be 2147483648`, rule.D{ // Beyond 32-bit int
		"timestamp": int64(2147483647),
	}, true},

	// Large Integer Precision (> 2^53)
	{"datetime_large_int_precision", `timestamp dq 9223372036854775807`, rule.D{
		"timestamp": int64(9223372036854775807), // Max int64
	}, true},

	// Error Conditions - Invalid Formats (should return false)
	{"datetime_invalid_rfc3339_format", `created_at af "2024-07-09T25:12:00Z"`, rule.D{
		"created_at": "2024-07-09T22:12:00Z",
	}, false}, // Invalid hour in literal
	{"datetime_invalid_property_format", `created_at af "2024-07-09T22:12:00Z"`, rule.D{
		"created_at": "invalid-date-format",
	}, false},
	{"datetime_non_numeric_unix_string", `timestamp af 1720558320`, rule.D{
		"timestamp": "not-a-number",
	}, false},

	// Missing Properties (should return false)
	{"datetime_missing_left_property", `nonexistent af "2024-07-09T22:12:00Z"`, rule.D{
		"other_field": "value",
	}, false},
	{"datetime_missing_right_property", `created_at af nonexistent`, rule.D{
		"created_at": "2024-07-09T22:12:00Z",
	}, false},
	{"datetime_missing_nested_property", `event.nonexistent af "2024-07-09T22:12:00Z"`, rule.D{
		"event": rule.D{
			"other_field": "value",
		},
	}, false},

	// Complex Scenarios with Logical Operators
	{
		"datetime_complex_and_operation",
		`start_time af "2024-07-09T22:12:00Z" and end_time be "2024-07-09T22:15:00Z"`,
		rule.D{
			"start_time": "2024-07-09T22:12:01Z",
			"end_time":   "2024-07-09T22:14:59Z",
		},
		true,
	},
	{
		"datetime_complex_or_operation",
		`created_at be "2024-07-09T22:12:00Z" or updated_at af "2024-07-09T22:15:00Z"`,
		rule.D{
			"created_at": "2024-07-09T22:13:00Z", // After (false)
			"updated_at": "2024-07-09T22:15:01Z", // After (true)
		},
		true,
	},
	{
		"datetime_complex_nested_and_mixed",
		`session.start_time af user.created_at and session.end_time be 1720558400`,
		rule.D{
			"session": rule.D{
				"start_time": "2024-07-09T22:12:01Z",
				"end_time":   int64(1720558399),
			},
			"user": rule.D{
				"created_at": "2024-07-09T22:12:00Z",
			},
		},
		true,
	},

	// DL (Days Less) Operator Tests - comparing timestamps with NOW
	// Note: These tests use timestamps that are safely in the past to ensure deterministic results

	// Basic DL functionality with RFC3339 timestamps
	{"dl_rfc3339_within_threshold", `created_at dl 3650`, rule.D{
		"created_at": "2020-01-01T00:00:00Z", // About 5+ years ago, within 3650 days (10 years)
	}, true},
	{"dl_rfc3339_exactly_threshold", `created_at dl 1`, rule.D{
		"created_at": "2020-01-01T00:00:00Z", // Well beyond 1 day threshold
	}, false},
	{"dl_rfc3339_beyond_threshold", `created_at dl 365`, rule.D{
		"created_at": "2020-01-01T00:00:00Z", // Well beyond 365 days threshold
	}, false},

	// DL with Unix timestamps
	{"dl_unix_within_threshold", `timestamp dl 3650`, rule.D{
		"timestamp": int64(1577836800), // 2020-01-01T00:00:00Z, within 3650 days
	}, true},
	{"dl_unix_beyond_threshold", `timestamp dl 365`, rule.D{
		"timestamp": int64(1577836800), // 2020-01-01T00:00:00Z - well beyond 365 days
	}, false},

	// DL with fractional days
	{"dl_fractional_days_within", `created_at dl 3650.5`, rule.D{
		"created_at": "2020-01-01T00:00:00Z",
	}, true},
	{"dl_fractional_days_beyond", `created_at dl 0.5`, rule.D{
		"created_at": "2020-01-01T00:00:00Z", // Beyond 0.5 days
	}, false},

	// DL with string number as days parameter
	{"dl_string_days_within", `created_at dl "3650"`, rule.D{
		"created_at": "2020-01-01T00:00:00Z",
	}, true},
	{"dl_string_days_beyond", `created_at dl "365"`, rule.D{
		"created_at": "2020-01-01T00:00:00Z",
	}, false},

	// DL with nested properties
	{"dl_nested_property", `event.created_at dl 3650`, rule.D{
		"event": rule.D{
			"created_at": "2020-01-01T00:00:00Z",
		},
	}, true},
	{"dl_deep_nested_property", `user.profile.last_login dl 3650`, rule.D{
		"user": rule.D{
			"profile": rule.D{
				"last_login": "2020-01-01T00:00:00Z",
			},
		},
	}, true},

	// DL error cases
	{"dl_invalid_timestamp", `created_at dl 30`, rule.D{
		"created_at": "invalid-timestamp",
	}, false},
	{"dl_invalid_days_string", `created_at dl "not-a-number"`, rule.D{
		"created_at": "2020-01-01T00:00:00Z",
	}, false},
	{"dl_boolean_days", `created_at dl true`, rule.D{
		"created_at": "2020-01-01T00:00:00Z",
	}, false},
	{"dl_missing_property", `nonexistent dl 30`, rule.D{
		"other_field": "value",
	}, false},

	// DL with zero and negative days
	{"dl_zero_days", `created_at dl 0`, rule.D{
		"created_at": "2020-01-01T00:00:00Z", // Definitely beyond 0 days
	}, false},
	{"dl_negative_days", `created_at dl -1`, rule.D{
		"created_at": "2020-01-01T00:00:00Z", // Should be false for negative threshold
	}, false},

	// DL in complex expressions
	{"dl_with_and_operator", `created_at dl 3650 and status eq "active"`, rule.D{
		"created_at": "2020-01-01T00:00:00Z",
		"status":     "active",
	}, true},
	{"dl_with_or_operator", `created_at dl 365 or updated_at dl 3650`, rule.D{
		"created_at": "2020-01-01T00:00:00Z", // Beyond 365 days (false)
		"updated_at": "2020-01-01T00:00:00Z", // Within 3650 days (true)
	}, true},

	// DG (Days Greater) Operator Tests - comparing timestamps with NOW (opposite of DL)
	// Note: These tests use timestamps that are safely in the past to ensure deterministic results

	// Basic DG functionality with RFC3339 timestamps
	{"dg_rfc3339_beyond_threshold", `created_at dg 365`, rule.D{
		"created_at": "2020-01-01T00:00:00Z", // About 5+ years ago, beyond 365 days (1 year)
	}, true},
	{"dg_rfc3339_within_threshold", `created_at dg 3650`, rule.D{
		"created_at": "2020-01-01T00:00:00Z", // About 5+ years ago, within 3650 days (10 years)
	}, false},
	{"dg_rfc3339_exactly_threshold", `created_at dg 0.5`, rule.D{
		"created_at": "2020-01-01T00:00:00Z", // Well beyond 0.5 days threshold
	}, true},

	// DG with Unix timestamps
	{"dg_unix_beyond_threshold", `timestamp dg 365`, rule.D{
		"timestamp": int64(1577836800), // 2020-01-01T00:00:00Z, beyond 365 days
	}, true},
	{"dg_unix_within_threshold", `timestamp dg 3650`, rule.D{
		"timestamp": int64(1577836800), // 2020-01-01T00:00:00Z, within 3650 days
	}, false},

	// DG with fractional days
	{"dg_fractional_days_beyond", `created_at dg 365.5`, rule.D{
		"created_at": "2020-01-01T00:00:00Z",
	}, true},
	{"dg_fractional_days_within", `created_at dg 3650.5`, rule.D{
		"created_at": "2020-01-01T00:00:00Z", // Within 3650.5 days
	}, false},

	// DG with string number as days parameter
	{"dg_string_days_beyond", `created_at dg "365"`, rule.D{
		"created_at": "2020-01-01T00:00:00Z",
	}, true},
	{"dg_string_days_within", `created_at dg "3650"`, rule.D{
		"created_at": "2020-01-01T00:00:00Z",
	}, false},

	// DG with nested properties
	{"dg_nested_property", `event.created_at dg 365`, rule.D{
		"event": rule.D{
			"created_at": "2020-01-01T00:00:00Z",
		},
	}, true},
	{"dg_deep_nested_property", `user.profile.last_login dg 365`, rule.D{
		"user": rule.D{
			"profile": rule.D{
				"last_login": "2020-01-01T00:00:00Z",
			},
		},
	}, true},

	// DG error cases
	{"dg_invalid_timestamp", `created_at dg 30`, rule.D{
		"created_at": "invalid-timestamp",
	}, false},
	{"dg_invalid_days_string", `created_at dg "not-a-number"`, rule.D{
		"created_at": "2020-01-01T00:00:00Z",
	}, false},
	{"dg_boolean_days", `created_at dg true`, rule.D{
		"created_at": "2020-01-01T00:00:00Z",
	}, false},
	{"dg_missing_property", `nonexistent dg 30`, rule.D{
		"other_field": "value",
	}, false},

	// DG with zero and negative days
	{"dg_zero_days", `created_at dg 0`, rule.D{
		"created_at": "2020-01-01T00:00:00Z", // Definitely beyond 0 days
	}, true},
	{"dg_negative_days", `created_at dg -1`, rule.D{
		"created_at": "2020-01-01T00:00:00Z", // Should be true for negative threshold
	}, true},

	// DG in complex expressions
	{"dg_with_and_operator", `created_at dg 365 and status eq "inactive"`, rule.D{
		"created_at": "2020-01-01T00:00:00Z",
		"status":     "inactive",
	}, true},
	{"dg_with_or_operator", `created_at dg 3650 or updated_at dg 365`, rule.D{
		"created_at": "2020-01-01T00:00:00Z", // Within 3650 days (false)
		"updated_at": "2020-01-01T00:00:00Z", // Beyond 365 days (true)
	}, true},

	// DL vs DG comparison tests (opposite behaviors)
	{"dl_vs_dg_same_input_opposite_results_1", `created_at dl 365`, rule.D{
		"created_at": "2020-01-01T00:00:00Z", // Beyond 365 days - dl should be false
	}, false},
	{"dl_vs_dg_same_input_opposite_results_2", `created_at dg 365`, rule.D{
		"created_at": "2020-01-01T00:00:00Z", // Beyond 365 days - dg should be true
	}, true},
}
