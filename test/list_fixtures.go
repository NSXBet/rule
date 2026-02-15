package test

import "github.com/NSXBet/rule"

/* ---------- Array length ---------- */

//nolint:gochecknoglobals // Test data
var ArrayLengthTests = []Case{
	// Basic length checks
	{
		"length_eq_true",
		"items.length eq 3",
		rule.D{"items": []any{1, 2, 3}},
		true,
	},
	{
		"length_eq_false",
		"items.length eq 5",
		rule.D{"items": []any{1, 2, 3}},
		false,
	},
	{
		"length_gt_true",
		"items.length gt 2",
		rule.D{"items": []any{"a", "b", "c"}},
		true,
	},
	{
		"length_gt_false",
		"items.length gt 5",
		rule.D{"items": []any{"a", "b", "c"}},
		false,
	},
	{
		"length_lt_true",
		"items.length lt 5",
		rule.D{"items": []any{1, 2}},
		true,
	},
	{
		"length_ge_true",
		"items.length ge 3",
		rule.D{"items": []any{1, 2, 3}},
		true,
	},
	{
		"length_le_true",
		"items.length le 3",
		rule.D{"items": []any{1, 2, 3}},
		true,
	},
	{
		"length_ne_true",
		"items.length ne 0",
		rule.D{"items": []any{1}},
		true,
	},
	// Empty array
	{
		"length_empty_eq_zero",
		"items.length eq 0",
		rule.D{"items": []any{}},
		true,
	},
	{
		"length_empty_gt_zero",
		"items.length gt 0",
		rule.D{"items": []any{}},
		false,
	},
	// Nested property path with length
	{
		"length_nested_property",
		"data.selections.length gt 1",
		rule.D{
			"data": rule.D{
				"selections": []any{
					rule.D{"id": "1"},
					rule.D{"id": "2"},
					rule.D{"id": "3"},
				},
			},
		},
		true,
	},
	{
		"length_deep_nested",
		"a.b.c.length eq 2",
		rule.D{
			"a": rule.D{
				"b": rule.D{
					"c": []any{10, 20},
				},
			},
		},
		true,
	},
	// Combined with other conditions
	{
		"length_combined_and",
		`items.length gt 2 and status eq "active"`,
		rule.D{
			"items":  []any{1, 2, 3},
			"status": "active",
		},
		true,
	},
	{
		"length_combined_and_false",
		`items.length gt 5 and status eq "active"`,
		rule.D{
			"items":  []any{1, 2, 3},
			"status": "active",
		},
		false,
	},
	{
		"length_combined_or",
		`items.length eq 0 or fallback eq true`,
		rule.D{
			"items":    []any{1},
			"fallback": true,
		},
		true,
	},
	// Missing array property
	{
		"length_missing_property",
		"missing.length gt 0",
		rule.D{},
		false,
	},
	// Non-array property with length
	{
		"length_on_map_with_length_key",
		"obj.length eq 42",
		rule.D{
			"obj": rule.D{"length": 42},
		},
		true,
	},
	// Array of objects
	{
		"length_array_of_objects",
		"selections.length eq 3",
		rule.D{
			"selections": []any{
				rule.D{"event_status": "EVENT_STATUS_IN_PROGRESS"},
				rule.D{"event_status": "EVENT_STATUS_FINISHED"},
				rule.D{"event_status": "EVENT_STATUS_CANCELLED"},
			},
		},
		true,
	},
}

/* ---------- ANY quantifier ---------- */

//nolint:gochecknoglobals // Test data
var QuantifierAnyTests = []Case{
	// Basic any - match found
	{
		"any_basic_match",
		`items any (status eq "active")`,
		rule.D{
			"items": []any{
				rule.D{"status": "inactive"},
				rule.D{"status": "active"},
				rule.D{"status": "inactive"},
			},
		},
		true,
	},
	// Basic any - no match
	{
		"any_basic_no_match",
		`items any (status eq "active")`,
		rule.D{
			"items": []any{
				rule.D{"status": "inactive"},
				rule.D{"status": "disabled"},
			},
		},
		false,
	},
	// Any with empty array
	{
		"any_empty_array",
		`items any (status eq "active")`,
		rule.D{"items": []any{}},
		false,
	},
	// Any with numeric comparison
	{
		"any_numeric",
		"scores any (value gt 90)",
		rule.D{
			"scores": []any{
				rule.D{"value": 50},
				rule.D{"value": 75},
				rule.D{"value": 95},
			},
		},
		true,
	},
	{
		"any_numeric_no_match",
		"scores any (value gt 100)",
		rule.D{
			"scores": []any{
				rule.D{"value": 50},
				rule.D{"value": 75},
			},
		},
		false,
	},
	// Any with boolean
	{
		"any_boolean",
		"items any (is_active eq true)",
		rule.D{
			"items": []any{
				rule.D{"is_active": false},
				rule.D{"is_active": true},
			},
		},
		true,
	},
	// Any with compound sub-expression
	{
		"any_compound_and",
		`items any (status eq "active" and priority gt 5)`,
		rule.D{
			"items": []any{
				rule.D{"status": "active", "priority": 3},
				rule.D{"status": "active", "priority": 8},
				rule.D{"status": "inactive", "priority": 9},
			},
		},
		true,
	},
	{
		"any_compound_and_no_match",
		`items any (status eq "active" and priority gt 10)`,
		rule.D{
			"items": []any{
				rule.D{"status": "active", "priority": 3},
				rule.D{"status": "active", "priority": 8},
			},
		},
		false,
	},
	// Any with OR sub-expression
	{
		"any_compound_or",
		`items any (status eq "cancelled" or status eq "refunded")`,
		rule.D{
			"items": []any{
				rule.D{"status": "active"},
				rule.D{"status": "refunded"},
			},
		},
		true,
	},
	// Any with string operations
	{
		"any_string_contains",
		`items any (name co "John")`,
		rule.D{
			"items": []any{
				rule.D{"name": "Alice Smith"},
				rule.D{"name": "John Doe"},
			},
		},
		true,
	},
	{
		"any_string_starts_with",
		`items any (code sw "PRE_")`,
		rule.D{
			"items": []any{
				rule.D{"code": "POST_123"},
				rule.D{"code": "PRE_456"},
			},
		},
		true,
	},
	// Any with IN operator
	{
		"any_in_operator",
		`items any (color in ["red", "blue"])`,
		rule.D{
			"items": []any{
				rule.D{"color": "green"},
				rule.D{"color": "blue"},
			},
		},
		true,
	},
	// Any with presence operator
	{
		"any_presence",
		"items any (optional_field pr)",
		rule.D{
			"items": []any{
				rule.D{"name": "a"},
				rule.D{"name": "b", "optional_field": "value"},
			},
		},
		true,
	},
	// Any combined with outer conditions
	{
		"any_combined_outer",
		`type eq "order" and items any (status eq "shipped")`,
		rule.D{
			"type": "order",
			"items": []any{
				rule.D{"status": "pending"},
				rule.D{"status": "shipped"},
			},
		},
		true,
	},
	{
		"any_combined_outer_false",
		`type eq "invoice" and items any (status eq "shipped")`,
		rule.D{
			"type": "order",
			"items": []any{
				rule.D{"status": "shipped"},
			},
		},
		false,
	},
	// Any on nested property
	{
		"any_nested_property",
		`data.items any (value gt 0)`,
		rule.D{
			"data": rule.D{
				"items": []any{
					rule.D{"value": -1},
					rule.D{"value": 5},
				},
			},
		},
		true,
	},
	// Any with missing property in elements
	{
		"any_missing_property_in_element",
		`items any (missing_field eq "value")`,
		rule.D{
			"items": []any{
				rule.D{"name": "a"},
				rule.D{"name": "b"},
			},
		},
		false,
	},
	// Any with NOT sub-expression
	{
		"any_not_sub_expression",
		`items any (not (status eq "cancelled"))`,
		rule.D{
			"items": []any{
				rule.D{"status": "cancelled"},
				rule.D{"status": "active"},
			},
		},
		true,
	},
	// Any on non-existent array
	{
		"any_missing_array",
		`missing_items any (status eq "active")`,
		rule.D{},
		false,
	},
	// Any on non-array value
	{
		"any_non_array_value",
		`name any (status eq "active")`,
		rule.D{"name": "not_an_array"},
		false,
	},
}

/* ---------- ALL quantifier ---------- */

//nolint:gochecknoglobals // Test data
var QuantifierAllTests = []Case{
	// All match
	{
		"all_basic_match",
		`items all (status eq "active")`,
		rule.D{
			"items": []any{
				rule.D{"status": "active"},
				rule.D{"status": "active"},
				rule.D{"status": "active"},
			},
		},
		true,
	},
	// Not all match
	{
		"all_basic_no_match",
		`items all (status eq "active")`,
		rule.D{
			"items": []any{
				rule.D{"status": "active"},
				rule.D{"status": "inactive"},
			},
		},
		false,
	},
	// All with empty array (vacuous truth)
	{
		"all_empty_array",
		`items all (status eq "active")`,
		rule.D{"items": []any{}},
		true,
	},
	// All with numeric
	{
		"all_numeric",
		"scores all (value ge 50)",
		rule.D{
			"scores": []any{
				rule.D{"value": 60},
				rule.D{"value": 75},
				rule.D{"value": 90},
			},
		},
		true,
	},
	{
		"all_numeric_false",
		"scores all (value ge 50)",
		rule.D{
			"scores": []any{
				rule.D{"value": 30},
				rule.D{"value": 75},
			},
		},
		false,
	},
	// All with compound expression
	{
		"all_compound",
		`items all (status eq "active" and enabled eq true)`,
		rule.D{
			"items": []any{
				rule.D{"status": "active", "enabled": true},
				rule.D{"status": "active", "enabled": true},
			},
		},
		true,
	},
	{
		"all_compound_false",
		`items all (status eq "active" and enabled eq true)`,
		rule.D{
			"items": []any{
				rule.D{"status": "active", "enabled": true},
				rule.D{"status": "active", "enabled": false},
			},
		},
		false,
	},
	// All combined with outer conditions
	{
		"all_combined_outer",
		`category eq "premium" and items all (quality ge 8)`,
		rule.D{
			"category": "premium",
			"items": []any{
				rule.D{"quality": 9},
				rule.D{"quality": 8},
			},
		},
		true,
	},
	// All on nested property
	{
		"all_nested_property",
		`data.items all (is_valid eq true)`,
		rule.D{
			"data": rule.D{
				"items": []any{
					rule.D{"is_valid": true},
					rule.D{"is_valid": true},
				},
			},
		},
		true,
	},
	// All on missing array
	{
		"all_missing_array",
		`missing all (status eq "active")`,
		rule.D{},
		false,
	},
	// All on non-array
	{
		"all_non_array",
		`name all (status eq "active")`,
		rule.D{"name": "string_value"},
		false,
	},
	// Single element all
	{
		"all_single_element",
		`items all (value gt 0)`,
		rule.D{
			"items": []any{
				rule.D{"value": 5},
			},
		},
		true,
	},
}

/* ---------- NONE quantifier ---------- */

//nolint:gochecknoglobals // Test data
var QuantifierNoneTests = []Case{
	// None match (true)
	{
		"none_basic_match",
		`items none (status eq "cancelled")`,
		rule.D{
			"items": []any{
				rule.D{"status": "active"},
				rule.D{"status": "pending"},
			},
		},
		true,
	},
	// Some match (false)
	{
		"none_basic_has_match",
		`items none (status eq "cancelled")`,
		rule.D{
			"items": []any{
				rule.D{"status": "active"},
				rule.D{"status": "cancelled"},
			},
		},
		false,
	},
	// None with empty array (true)
	{
		"none_empty_array",
		`items none (status eq "cancelled")`,
		rule.D{"items": []any{}},
		true,
	},
	// None with numeric
	{
		"none_numeric",
		"scores none (value lt 0)",
		rule.D{
			"scores": []any{
				rule.D{"value": 10},
				rule.D{"value": 20},
			},
		},
		true,
	},
	{
		"none_numeric_false",
		"scores none (value lt 0)",
		rule.D{
			"scores": []any{
				rule.D{"value": 10},
				rule.D{"value": -5},
			},
		},
		false,
	},
	// None with compound expression
	{
		"none_compound",
		`items none (status eq "error" and severity eq "critical")`,
		rule.D{
			"items": []any{
				rule.D{"status": "error", "severity": "low"},
				rule.D{"status": "ok", "severity": "critical"},
			},
		},
		true,
	},
	{
		"none_compound_false",
		`items none (status eq "error" and severity eq "critical")`,
		rule.D{
			"items": []any{
				rule.D{"status": "error", "severity": "critical"},
				rule.D{"status": "ok", "severity": "low"},
			},
		},
		false,
	},
	// None combined with outer conditions
	{
		"none_combined_outer",
		`is_verified eq true and items none (is_fraud eq true)`,
		rule.D{
			"is_verified": true,
			"items": []any{
				rule.D{"is_fraud": false},
				rule.D{"is_fraud": false},
			},
		},
		true,
	},
	// None on nested property
	{
		"none_nested_property",
		`data.items none (status eq "deleted")`,
		rule.D{
			"data": rule.D{
				"items": []any{
					rule.D{"status": "active"},
					rule.D{"status": "archived"},
				},
			},
		},
		true,
	},
	// None on missing array
	{
		"none_missing_array",
		`missing none (status eq "active")`,
		rule.D{},
		false,
	},
	// None on non-array
	{
		"none_non_array",
		`name none (status eq "active")`,
		rule.D{"name": "string_value"},
		false,
	},
}

/* ---------- Real-world betting scenarios ---------- */

//nolint:gochecknoglobals // Test data
var ListRealWorldTests = []Case{
	// Betting: check if any selection has a cancelled event
	{
		"betting_any_cancelled",
		`selections any (event_status eq "EVENT_STATUS_CANCELLED")`,
		rule.D{
			"selections": []any{
				rule.D{"event_status": "EVENT_STATUS_IN_PROGRESS", "sport_name": "Football"},
				rule.D{"event_status": "EVENT_STATUS_CANCELLED", "sport_name": "Basketball"},
				rule.D{"event_status": "EVENT_STATUS_FINISHED", "sport_name": "Tennis"},
			},
		},
		true,
	},
	{
		"betting_no_cancelled",
		`selections any (event_status eq "EVENT_STATUS_CANCELLED")`,
		rule.D{
			"selections": []any{
				rule.D{"event_status": "EVENT_STATUS_IN_PROGRESS"},
				rule.D{"event_status": "EVENT_STATUS_FINISHED"},
			},
		},
		false,
	},
	// Betting: all selections are live
	{
		"betting_all_live",
		"selections all (is_live eq true)",
		rule.D{
			"selections": []any{
				rule.D{"is_live": true, "sport_name": "Football"},
				rule.D{"is_live": true, "sport_name": "Basketball"},
			},
		},
		true,
	},
	{
		"betting_not_all_live",
		"selections all (is_live eq true)",
		rule.D{
			"selections": []any{
				rule.D{"is_live": true},
				rule.D{"is_live": false},
			},
		},
		false,
	},
	// Betting: no selection is a loss
	{
		"betting_none_loss",
		`selections none (status eq "SELECTION_STATUS_LOSS")`,
		rule.D{
			"selections": []any{
				rule.D{"status": "SELECTION_STATUS_WIN"},
				rule.D{"status": "SELECTION_STATUS_IN_PROGRESS"},
			},
		},
		true,
	},
	// Betting: multiple bet with minimum selections
	{
		"betting_multiple_min_selections",
		`bet_type eq "BET_TYPE_MULTIPLE" and selections.length ge 3`,
		rule.D{
			"bet_type": "BET_TYPE_MULTIPLE",
			"selections": []any{
				rule.D{"id": "1"},
				rule.D{"id": "2"},
				rule.D{"id": "3"},
			},
		},
		true,
	},
	{
		"betting_multiple_insufficient_selections",
		`bet_type eq "BET_TYPE_MULTIPLE" and selections.length ge 3`,
		rule.D{
			"bet_type": "BET_TYPE_MULTIPLE",
			"selections": []any{
				rule.D{"id": "1"},
				rule.D{"id": "2"},
			},
		},
		false,
	},
	// Betting: VIP customer with high-odd live selections
	{
		"betting_vip_high_odd_live",
		`customer_data.is_vip eq true and selections any (is_live eq true and odd gt 3.0)`,
		rule.D{
			"customer_data": rule.D{"is_vip": true},
			"selections": []any{
				rule.D{"is_live": false, "odd": 1.5},
				rule.D{"is_live": true, "odd": 4.2},
			},
		},
		true,
	},
	{
		"betting_vip_no_high_odd_live",
		`customer_data.is_vip eq true and selections any (is_live eq true and odd gt 3.0)`,
		rule.D{
			"customer_data": rule.D{"is_vip": true},
			"selections": []any{
				rule.D{"is_live": true, "odd": 1.5},
				rule.D{"is_live": false, "odd": 4.2},
			},
		},
		false,
	},
	// Betting: complex rule - multiple type, enough selections, no cancelled events, from specific provider
	{
		"betting_complex_validation",
		`bet_type eq "BET_TYPE_MULTIPLE" and selections.length ge 2 and selections none (event_status eq "EVENT_STATUS_CANCELLED") and selections any (provider eq "PROVIDER_SPORTRADAR")`,
		rule.D{
			"bet_type": "BET_TYPE_MULTIPLE",
			"selections": []any{
				rule.D{
					"event_status": "EVENT_STATUS_IN_PROGRESS",
					"provider":     "PROVIDER_SPORTRADAR",
					"is_live":      true,
				},
				rule.D{
					"event_status": "EVENT_STATUS_NOT_STARTED",
					"provider":     "PROVIDER_RAMP",
					"is_live":      false,
				},
			},
		},
		true,
	},
	{
		"betting_complex_validation_cancelled",
		`bet_type eq "BET_TYPE_MULTIPLE" and selections.length ge 2 and selections none (event_status eq "EVENT_STATUS_CANCELLED")`,
		rule.D{
			"bet_type": "BET_TYPE_MULTIPLE",
			"selections": []any{
				rule.D{"event_status": "EVENT_STATUS_IN_PROGRESS"},
				rule.D{"event_status": "EVENT_STATUS_CANCELLED"},
			},
		},
		false,
	},
	// Betting: any selection is super odd
	{
		"betting_any_super_odd",
		"selections any (is_super_odd eq true)",
		rule.D{
			"selections": []any{
				rule.D{"is_super_odd": false, "odd": 1.5},
				rule.D{"is_super_odd": true, "odd": 5.0},
			},
		},
		true,
	},
	// Betting: all selections from the same sport
	{
		"betting_all_same_sport",
		`selections all (sport_name eq "Football")`,
		rule.D{
			"selections": []any{
				rule.D{"sport_name": "Football"},
				rule.D{"sport_name": "Football"},
				rule.D{"sport_name": "Football"},
			},
		},
		true,
	},
	{
		"betting_not_all_same_sport",
		`selections all (sport_name eq "Football")`,
		rule.D{
			"selections": []any{
				rule.D{"sport_name": "Football"},
				rule.D{"sport_name": "Basketball"},
			},
		},
		false,
	},
	// Betting: any selection with recommendation
	{
		"betting_any_recommendation",
		"selections any (is_recommendation eq true)",
		rule.D{
			"selections": []any{
				rule.D{"is_recommendation": false},
				rule.D{"is_recommendation": true},
			},
		},
		true,
	},
	// Betting: length combined with quantifier
	{
		"betting_length_and_quantifier",
		`selections.length gt 1 and selections all (odd gt 1.0)`,
		rule.D{
			"selections": []any{
				rule.D{"odd": 1.5},
				rule.D{"odd": 2.3},
			},
		},
		true,
	},
	// Betting: freebet with single selection constraint
	{
		"betting_freebet_single",
		`is_freebet eq true and selections.length eq 1`,
		rule.D{
			"is_freebet": true,
			"selections": []any{
				rule.D{"id": "sel_1"},
			},
		},
		true,
	},
	{
		"betting_freebet_multiple_invalid",
		`is_freebet eq true and selections.length eq 1`,
		rule.D{
			"is_freebet": true,
			"selections": []any{
				rule.D{"id": "sel_1"},
				rule.D{"id": "sel_2"},
			},
		},
		false,
	},
	// Selection with string operation inside quantifier
	{
		"betting_any_event_name_contains",
		`selections any (event_name co "Barcelona")`,
		rule.D{
			"selections": []any{
				rule.D{"event_name": "Real Madrid vs Atletico"},
				rule.D{"event_name": "FC Barcelona vs PSG"},
			},
		},
		true,
	},
	// Quantifier with NOT outer
	{
		"betting_not_any_cancelled",
		`not (selections any (event_status eq "EVENT_STATUS_CANCELLED"))`,
		rule.D{
			"selections": []any{
				rule.D{"event_status": "EVENT_STATUS_IN_PROGRESS"},
				rule.D{"event_status": "EVENT_STATUS_FINISHED"},
			},
		},
		true,
	},
	{
		"betting_not_any_cancelled_false",
		`not (selections any (event_status eq "EVENT_STATUS_CANCELLED"))`,
		rule.D{
			"selections": []any{
				rule.D{"event_status": "EVENT_STATUS_CANCELLED"},
			},
		},
		false,
	},
}

/* ---------- Nested list quantifiers (list inside list) ---------- */

//nolint:gochecknoglobals // Test data
var NestedListTests = []Case{
	// any inside any - match found
	{
		"nested_any_any_match",
		`groups any (items any (value gt 10))`,
		rule.D{
			"groups": []any{
				rule.D{"items": []any{
					rule.D{"value": 1},
					rule.D{"value": 2},
				}},
				rule.D{"items": []any{
					rule.D{"value": 5},
					rule.D{"value": 15},
				}},
			},
		},
		true,
	},
	// any inside any - no match
	{
		"nested_any_any_no_match",
		`groups any (items any (value gt 100))`,
		rule.D{
			"groups": []any{
				rule.D{"items": []any{
					rule.D{"value": 1},
					rule.D{"value": 2},
				}},
			},
		},
		false,
	},
	// all inside all
	{
		"nested_all_all_match",
		`groups all (items all (is_valid eq true))`,
		rule.D{
			"groups": []any{
				rule.D{"items": []any{
					rule.D{"is_valid": true},
					rule.D{"is_valid": true},
				}},
				rule.D{"items": []any{
					rule.D{"is_valid": true},
				}},
			},
		},
		true,
	},
	{
		"nested_all_all_no_match",
		`groups all (items all (is_valid eq true))`,
		rule.D{
			"groups": []any{
				rule.D{"items": []any{
					rule.D{"is_valid": true},
				}},
				rule.D{"items": []any{
					rule.D{"is_valid": false},
				}},
			},
		},
		false,
	},
	// none inside any
	{
		"nested_any_none",
		`groups any (items none (status eq "error"))`,
		rule.D{
			"groups": []any{
				rule.D{"items": []any{
					rule.D{"status": "error"},
				}},
				rule.D{"items": []any{
					rule.D{"status": "ok"},
					rule.D{"status": "ok"},
				}},
			},
		},
		true,
	},
	// any inside all
	{
		"nested_all_any_match",
		`groups all (items any (priority gt 0))`,
		rule.D{
			"groups": []any{
				rule.D{"items": []any{
					rule.D{"priority": 0},
					rule.D{"priority": 5},
				}},
				rule.D{"items": []any{
					rule.D{"priority": 3},
				}},
			},
		},
		true,
	},
	{
		"nested_all_any_no_match",
		`groups all (items any (priority gt 0))`,
		rule.D{
			"groups": []any{
				rule.D{"items": []any{
					rule.D{"priority": 5},
				}},
				rule.D{"items": []any{
					rule.D{"priority": 0},
					rule.D{"priority": 0},
				}},
			},
		},
		false,
	},
	// Nested length
	{
		"nested_length_inside_any",
		"groups any (items.length gt 2)",
		rule.D{
			"groups": []any{
				rule.D{"items": []any{rule.D{"x": 1}}},
				rule.D{"items": []any{rule.D{"x": 1}, rule.D{"x": 2}, rule.D{"x": 3}}},
			},
		},
		true,
	},
	{
		"nested_length_inside_all",
		"groups all (items.length ge 1)",
		rule.D{
			"groups": []any{
				rule.D{"items": []any{rule.D{"x": 1}}},
				rule.D{"items": []any{rule.D{"x": 1}, rule.D{"x": 2}}},
			},
		},
		true,
	},
	// Condition on parent + nested quantifier
	{
		"nested_parent_condition_and_inner_quantifier",
		`groups any (name eq "vip" and items any (score gt 90))`,
		rule.D{
			"groups": []any{
				rule.D{
					"name":  "regular",
					"items": []any{rule.D{"score": 95}},
				},
				rule.D{
					"name":  "vip",
					"items": []any{rule.D{"score": 50}, rule.D{"score": 99}},
				},
			},
		},
		true,
	},
	{
		"nested_parent_condition_and_inner_quantifier_no_match",
		`groups any (name eq "vip" and items any (score gt 100))`,
		rule.D{
			"groups": []any{
				rule.D{
					"name":  "vip",
					"items": []any{rule.D{"score": 50}, rule.D{"score": 99}},
				},
			},
		},
		false,
	},
	// Empty inner lists
	{
		"nested_any_empty_inner",
		`groups any (items any (x eq 1))`,
		rule.D{
			"groups": []any{
				rule.D{"items": []any{}},
			},
		},
		false,
	},
	{
		"nested_all_empty_inner",
		`groups all (items all (x eq 1))`,
		rule.D{
			"groups": []any{
				rule.D{"items": []any{}},
			},
		},
		true, // vacuous truth
	},
	{
		"nested_none_empty_inner",
		`groups all (items none (x eq 1))`,
		rule.D{
			"groups": []any{
				rule.D{"items": []any{}},
			},
		},
		true, // no elements to match
	},
	// Betting: selections with bet_builder_selections (real-world)
	{
		"betting_nested_bb_any_market",
		`selections any (bet_builder_selections any (market_name eq "Goals"))`,
		rule.D{
			"selections": []any{
				rule.D{
					"sport_name": "Football",
					"bet_builder_selections": []any{
						rule.D{"market_name": "Goals", "odd": 1.5},
						rule.D{"market_name": "Cards", "odd": 2.3},
					},
				},
				rule.D{
					"sport_name": "Basketball",
					"bet_builder_selections": []any{
						rule.D{"market_name": "Points", "odd": 3.1},
					},
				},
			},
		},
		true,
	},
	{
		"betting_nested_bb_no_match",
		`selections any (bet_builder_selections any (market_name eq "Corners"))`,
		rule.D{
			"selections": []any{
				rule.D{
					"bet_builder_selections": []any{
						rule.D{"market_name": "Goals"},
						rule.D{"market_name": "Cards"},
					},
				},
			},
		},
		false,
	},
	// Betting: all bet_builder_selections have odd > 1.0
	{
		"betting_nested_all_bb_odds",
		`selections all (bet_builder_selections all (odd gt 1.0))`,
		rule.D{
			"selections": []any{
				rule.D{
					"bet_builder_selections": []any{
						rule.D{"odd": 1.5},
						rule.D{"odd": 2.0},
					},
				},
				rule.D{
					"bet_builder_selections": []any{
						rule.D{"odd": 1.2},
					},
				},
			},
		},
		true,
	},
	// Betting: sport + bb combined
	{
		"betting_nested_sport_and_bb",
		`selections any (sport_name eq "Football" and bet_builder_selections any (odd gt 2.0))`,
		rule.D{
			"selections": []any{
				rule.D{
					"sport_name": "Football",
					"bet_builder_selections": []any{
						rule.D{"odd": 1.5},
						rule.D{"odd": 2.5},
					},
				},
			},
		},
		true,
	},
	{
		"betting_nested_sport_and_bb_no_match",
		`selections any (sport_name eq "Football" and bet_builder_selections any (odd gt 3.0))`,
		rule.D{
			"selections": []any{
				rule.D{
					"sport_name": "Football",
					"bet_builder_selections": []any{
						rule.D{"odd": 1.5},
						rule.D{"odd": 2.5},
					},
				},
			},
		},
		false,
	},
	// Betting: none of the bet_builders have errors
	{
		"betting_nested_none_bb_errors",
		`selections none (bet_builder_selections any (status eq "error"))`,
		rule.D{
			"selections": []any{
				rule.D{
					"bet_builder_selections": []any{
						rule.D{"status": "ok"},
						rule.D{"status": "ok"},
					},
				},
				rule.D{
					"bet_builder_selections": []any{
						rule.D{"status": "ok"},
					},
				},
			},
		},
		true,
	},
	{
		"betting_nested_none_bb_errors_found",
		`selections none (bet_builder_selections any (status eq "error"))`,
		rule.D{
			"selections": []any{
				rule.D{
					"bet_builder_selections": []any{
						rule.D{"status": "ok"},
					},
				},
				rule.D{
					"bet_builder_selections": []any{
						rule.D{"status": "error"},
					},
				},
			},
		},
		false,
	},
	// Betting: bb length inside quantifier
	{
		"betting_nested_bb_length",
		`selections any (is_bet_builder eq true and bet_builder_selections.length gt 1)`,
		rule.D{
			"selections": []any{
				rule.D{
					"is_bet_builder":         true,
					"bet_builder_selections": []any{rule.D{"a": 1}, rule.D{"a": 2}},
				},
			},
		},
		true,
	},
	// Complex: outer length + inner quantifier
	{
		"nested_complex_length_and_quantifiers",
		`selections.length ge 2 and selections all (bet_builder_selections.length ge 1) and selections any (bet_builder_selections any (odd gt 2.0))`,
		rule.D{
			"selections": []any{
				rule.D{
					"bet_builder_selections": []any{
						rule.D{"odd": 1.5},
					},
				},
				rule.D{
					"bet_builder_selections": []any{
						rule.D{"odd": 2.5},
					},
				},
			},
		},
		true,
	},
}
