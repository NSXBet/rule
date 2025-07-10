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
}
