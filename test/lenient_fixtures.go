package test

import "github.com/NSXBet/rule"

/* ---------- Lenient Mode (SQL-ish null semantics) ----------
 *
 * These fixtures run under NewEngineWithOptions(WithLenientMode()).
 * "missing" means the key is absent from the context map.
 *
 * SQL-ish null semantics (2-valued):
 *   null eq  <value>      -> false      null ne  <value>      -> true
 *   null lt/gt/le/ge value -> false      null not in [...]     -> true
 *   null co/sw/ew <value> -> false      null eq null          -> true
 *   datetime ops vs null  -> false      null ne null          -> false
 *
 * Presence (pr), logical (and/or/not) and quantifier operators are
 * unaffected by lenient mode and keep their strict behavior.
 */

//nolint:gochecknoglobals // Test data
var LenientEqualityTests = []Case{
	// eq with missing left
	{"lenient_eq_missing_left_value", "x eq 10", rule.D{}, false},
	{"lenient_eq_missing_left_string", `x eq "abc"`, rule.D{}, false},
	{"lenient_eq_missing_left_bool", "x eq true", rule.D{}, false},
	// eq with missing right (property on the right)
	{"lenient_eq_missing_right_value", "10 eq y", rule.D{}, false},
	{"lenient_eq_missing_right_string", `"abc" eq y`, rule.D{}, false},
	// eq with both missing -> true (null eq null)
	{"lenient_eq_both_missing", "x eq y", rule.D{}, true},
	{"lenient_eq_both_missing_nested", "a.b eq c.d", rule.D{}, true},
	// eq with present values still works
	{"lenient_eq_present_true", "x eq 10", rule.D{"x": 10}, true},
	{"lenient_eq_present_false", "x eq 10", rule.D{"x": 11}, false},
	// == alias
	{"lenient_equals_alias_missing", "x == y", rule.D{}, true},
	{"lenient_equals_alias_missing_vs_value", "x == 10", rule.D{}, false},
}

//nolint:gochecknoglobals // Test data
var LenientInequalityTests = []Case{
	// ne with missing left -> true (null ne value)
	{"lenient_ne_missing_left_value", "x ne 10", rule.D{}, true},
	{"lenient_ne_missing_left_string", `x ne "abc"`, rule.D{}, true},
	{"lenient_ne_missing_left_bool", "x ne true", rule.D{}, true},
	// ne with missing right -> true
	{"lenient_ne_missing_right_value", "10 ne y", rule.D{}, true},
	// ne with both missing -> false (null ne null)
	{"lenient_ne_both_missing", "x ne y", rule.D{}, false},
	// ne with present values still works
	{"lenient_ne_present_true", "x ne 10", rule.D{"x": 11}, true},
	{"lenient_ne_present_false", "x ne 10", rule.D{"x": 10}, false},
	// != alias
	{"lenient_not_equals_alias_missing", "x != 10", rule.D{}, true},
	{"lenient_not_equals_alias_both_missing", "x != y", rule.D{}, false},
}

//nolint:gochecknoglobals // Test data
var LenientRelationalTests = []Case{
	// ordering ops against null -> false
	{"lenient_lt_missing", "x lt 10", rule.D{}, false},
	{"lenient_gt_missing", "x gt 10", rule.D{}, false},
	{"lenient_le_missing", "x le 10", rule.D{}, false},
	{"lenient_ge_missing", "x ge 10", rule.D{}, false},
	// both missing -> still false (not null-comparable)
	{"lenient_lt_both_missing", "x lt y", rule.D{}, false},
	{"lenient_gt_both_missing", "x gt y", rule.D{}, false},
	// missing on the right
	{"lenient_lt_missing_right", "10 lt y", rule.D{}, false},
	// present still works
	{"lenient_lt_present_true", "x lt 10", rule.D{"x": 5}, true},
	{"lenient_lt_present_false", "x lt 10", rule.D{"x": 15}, false},
}

//nolint:gochecknoglobals // Test data
var LenientStringOpTests = []Case{
	{"lenient_co_missing", `x co "abc"`, rule.D{}, false},
	{"lenient_sw_missing", `x sw "abc"`, rule.D{}, false},
	{"lenient_ew_missing", `x ew "abc"`, rule.D{}, false},
	{"lenient_co_both_missing", `x co y`, rule.D{}, false},
	// present still works
	{"lenient_co_present_true", `x co "York"`, rule.D{"x": "New York"}, true},
	{"lenient_co_present_false", `x co "York"`, rule.D{"x": "Boston"}, false},
}

//nolint:gochecknoglobals // Test data
var LenientMembershipTests = []Case{
	// in: null in set -> false
	{"lenient_in_missing", "x in [1,2,3]", rule.D{}, false},
	{"lenient_in_missing_string", `x in ["a","b"]`, rule.D{}, false},
	{"lenient_in_present_true", "x in [1,2,3]", rule.D{"x": 2}, true},
	{"lenient_in_present_false", "x in [1,2,3]", rule.D{"x": 4}, false},
	// not in: null not in set -> true
	{"lenient_not_in_missing", "x not in [1,2,3]", rule.D{}, true},
	{"lenient_not_in_missing_string", `x not in ["a","b"]`, rule.D{}, true},
	{"lenient_not_in_present_true", "x not in [1,2,3]", rule.D{"x": 4}, true},
	{"lenient_not_in_present_false", "x not in [1,2,3]", rule.D{"x": 2}, false},
}

//nolint:gochecknoglobals // Test data
var LenientLogicalTests = []Case{
	// and: false && anything -> false (short-circuit, missing not reached)
	{"lenient_and_left_false", "(x eq 10) and (y ne 5)", rule.D{"x": 99}, false},
	// and: true && missing-comparison -> depends on lenient result of right
	{"lenient_and_right_missing_eq", "(x eq 10) and (y eq 5)", rule.D{"x": 10}, false},
	{"lenient_and_right_missing_ne", "(x eq 10) and (y ne 5)", rule.D{"x": 10}, true},
	// or: true short-circuits
	{"lenient_or_left_true", "(x eq 10) or (y eq 5)", rule.D{"x": 10}, true},
	{"lenient_or_left_false_missing_right", "(x eq 10) or (y eq 5)", rule.D{}, false},
	{"lenient_or_left_false_missing_right_ne", "(x eq 10) or (y ne 5)", rule.D{}, true},
	// not: not (missing eq value) -> not false -> true
	{"lenient_not_missing_eq", "not (x eq 10)", rule.D{}, true},
	// not: not (missing ne value) -> not true -> false
	{"lenient_not_missing_ne", "not (x ne 10)", rule.D{}, false},
	// not (both missing eq) -> not true -> false
	{"lenient_not_both_missing_eq", "not (x eq y)", rule.D{}, false},
}

//nolint:gochecknoglobals // Test data
var LenientNestedTests = []Case{
	{"lenient_nested_eq_missing_chain", "a.b.c.d eq 10", rule.D{}, false},
	{"lenient_nested_ne_missing_chain", "a.b.c.d ne 10", rule.D{}, true},
	{"lenient_nested_eq_partial", "a.b eq c.d", rule.D{"a": rule.D{"b": 1}}, false},
	{"lenient_nested_eq_partial_ne", "a.b ne c.d", rule.D{"a": rule.D{"b": 1}}, true},
	{"lenient_nested_both_missing_eq", "a.b eq c.d", rule.D{}, true},
	{"lenient_nested_both_missing_ne", "a.b ne c.d", rule.D{}, false},
}

//nolint:gochecknoglobals // Test data
var LenientDateTimeTests = []Case{
	// datetime ops vs null -> false
	{"lenient_dq_missing", `created_at dq "2024-01-01T00:00:00Z"`, rule.D{}, false},
	{"lenient_dn_missing", `created_at dn "2024-01-01T00:00:00Z"`, rule.D{}, false},
	{"lenient_be_missing", `created_at be "2024-01-01T00:00:00Z"`, rule.D{}, false},
	{"lenient_bq_missing", `created_at bq "2024-01-01T00:00:00Z"`, rule.D{}, false},
	{"lenient_af_missing", `created_at af "2024-01-01T00:00:00Z"`, rule.D{}, false},
	{"lenient_aq_missing", `created_at aq "2024-01-01T00:00:00Z"`, rule.D{}, false},
	{"lenient_dl_missing", "created_at dl 30", rule.D{}, false},
	{"lenient_dg_missing", "created_at dg 30", rule.D{}, false},
	// both missing -> false
	{"lenient_dq_both_missing", "a dq b", rule.D{}, false},
	// present still works
	{"lenient_af_present", `created_at af "2024-01-01T00:00:00Z"`, rule.D{
		"created_at": "2024-06-01T00:00:00Z",
	}, true},
}

//nolint:gochecknoglobals // Test data
var LenientPresenceTests = []Case{
	// pr is unaffected by lenient mode: missing -> false, nil -> true, present -> true
	{"lenient_pr_missing", "x pr", rule.D{}, false},
	{"lenient_pr_present", "x pr", rule.D{"x": 1}, true},
	{"lenient_pr_nil_explicit", "x pr", rule.D{"x": nil}, true},
	{"lenient_pr_nested_missing", "a.b.c pr", rule.D{}, false},
	{"lenient_pr_nested_present", "a.b pr", rule.D{"a": rule.D{"b": 1}}, true},
	// explicit nil (present) vs absent: only the absent operand routes through
	// lenientCompare, so eq -> false and ne -> true.
	{"lenient_nil_vs_absent_eq", "x eq y", rule.D{"x": nil}, false},
	{"lenient_nil_vs_absent_ne", "x ne y", rule.D{"x": nil}, true},
}
