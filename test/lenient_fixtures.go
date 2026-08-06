package test

import "github.com/NSXBet/rule"

/* ---------- Lenient Mode (neutral semantics) ----------
 *
 * These fixtures run under NewEngineWithOptions(WithLenientMode()).
 * "missing" means the key is absent from the context map.
 *
 * Neutral semantics: a comparison predicate involving a missing attribute
 * returns true, so it imposes no constraint. This is the identity element
 * for AND-chains (the dominant pattern in betting lifecycle rules).
 *
 *   x eq 10   (x absent) -> true   x lt 10   (x absent) -> true
 *   x ne 10   (x absent) -> true   x co "a"  (x absent) -> true
 *   x in [...]  (x absent) -> true   x not in [...] -> true
 *   datetime ops with missing operand -> true
 *   x eq y  (both absent) -> true
 *
 * Composition consequences (documented, intentional):
 *   not (x eq 10)         -> not true -> false
 *   x eq 10 or y eq 20    -> OR short-circuits to true
 *   selections none (r gt 0)  -> sees true per element where r is absent -> false
 *
 * Presence (pr), logical (and/or/not) and quantifier operators themselves
 * are unaffected by lenient mode and keep their strict behavior.
 */

//nolint:gochecknoglobals // Test data
var LenientEqualityTests = []Case{
	// eq with missing left -> neutral true
	{"lenient_eq_missing_left_value", "x eq 10", rule.D{}, true},
	{"lenient_eq_missing_left_string", `x eq "abc"`, rule.D{}, true},
	{"lenient_eq_missing_left_bool", "x eq true", rule.D{}, true},
	// eq with missing right (property on the right) -> true
	{"lenient_eq_missing_right_value", "10 eq y", rule.D{}, true},
	{"lenient_eq_missing_right_string", `"abc" eq y`, rule.D{}, true},
	// eq with both missing -> true
	{"lenient_eq_both_missing", "x eq y", rule.D{}, true},
	{"lenient_eq_both_missing_nested", "a.b eq c.d", rule.D{}, true},
	// eq with present values still works (regular path)
	{"lenient_eq_present_true", "x eq 10", rule.D{"x": 10}, true},
	{"lenient_eq_present_false", "x eq 10", rule.D{"x": 11}, false},
	// == alias
	{"lenient_equals_alias_missing", "x == y", rule.D{}, true},
	{"lenient_equals_alias_missing_vs_value", "x == 10", rule.D{}, true},
}

//nolint:gochecknoglobals // Test data
var LenientInequalityTests = []Case{
	// ne with missing left -> neutral true
	{"lenient_ne_missing_left_value", "x ne 10", rule.D{}, true},
	{"lenient_ne_missing_left_string", `x ne "abc"`, rule.D{}, true},
	{"lenient_ne_missing_left_bool", "x ne true", rule.D{}, true},
	// ne with missing right -> true
	{"lenient_ne_missing_right_value", "10 ne y", rule.D{}, true},
	// ne with both missing -> true
	{"lenient_ne_both_missing", "x ne y", rule.D{}, true},
	// ne with present values still works (regular path)
	{"lenient_ne_present_true", "x ne 10", rule.D{"x": 11}, true},
	{"lenient_ne_present_false", "x ne 10", rule.D{"x": 10}, false},
	// != alias
	{"lenient_not_equals_alias_missing", "x != 10", rule.D{}, true},
	{"lenient_not_equals_alias_both_missing", "x != y", rule.D{}, true},
}

//nolint:gochecknoglobals // Test data
var LenientRelationalTests = []Case{
	// ordering ops against null -> neutral true
	{"lenient_lt_missing", "x lt 10", rule.D{}, true},
	{"lenient_gt_missing", "x gt 10", rule.D{}, true},
	{"lenient_le_missing", "x le 10", rule.D{}, true},
	{"lenient_ge_missing", "x ge 10", rule.D{}, true},
	// both missing -> true
	{"lenient_lt_both_missing", "x lt y", rule.D{}, true},
	{"lenient_gt_both_missing", "x gt y", rule.D{}, true},
	// missing on the right
	{"lenient_lt_missing_right", "10 lt y", rule.D{}, true},
	// present still works (regular path)
	{"lenient_lt_present_true", "x lt 10", rule.D{"x": 5}, true},
	{"lenient_lt_present_false", "x lt 10", rule.D{"x": 15}, false},
}

//nolint:gochecknoglobals // Test data
var LenientStringOpTests = []Case{
	{"lenient_co_missing", `x co "abc"`, rule.D{}, true},
	{"lenient_sw_missing", `x sw "abc"`, rule.D{}, true},
	{"lenient_ew_missing", `x ew "abc"`, rule.D{}, true},
	{"lenient_co_both_missing", `x co y`, rule.D{}, true},
	// present still works (regular path)
	{"lenient_co_present_true", `x co "York"`, rule.D{"x": "New York"}, true},
	{"lenient_co_present_false", `x co "York"`, rule.D{"x": "Boston"}, false},
}

//nolint:gochecknoglobals // Test data
var LenientMembershipTests = []Case{
	// in with missing operand -> neutral true
	{"lenient_in_missing", "x in [1,2,3]", rule.D{}, true},
	{"lenient_in_missing_string", `x in ["a","b"]`, rule.D{}, true},
	{"lenient_in_missing_right", "v in y", rule.D{"v": 2}, true},
	{"lenient_in_both_missing", "x in y", rule.D{}, true},
	// present still works (regular path)
	{"lenient_in_present_true", "x in [1,2,3]", rule.D{"x": 2}, true},
	{"lenient_in_present_false", "x in [1,2,3]", rule.D{"x": 4}, false},
	// not in with missing operand -> neutral true
	{"lenient_not_in_missing", "x not in [1,2,3]", rule.D{}, true},
	{"lenient_not_in_missing_string", `x not in ["a","b"]`, rule.D{}, true},
	{"lenient_not_in_missing_right", "v not in y", rule.D{"v": 2}, true},
	{"lenient_not_in_both_missing", "x not in y", rule.D{}, true},
	// present still works (regular path)
	{"lenient_not_in_present_true", "x not in [1,2,3]", rule.D{"x": 4}, true},
	{"lenient_not_in_present_false", "x not in [1,2,3]", rule.D{"x": 2}, false},
}

//nolint:gochecknoglobals // Test data
var LenientLogicalTests = []Case{
	// and: neutral-then-false -> false (real constraint still applies)
	{"lenient_and_neutral_then_false", "(x eq 10) and (y eq 5)", rule.D{"y": 9}, false},
	// and: true-then-neutral -> true
	{"lenient_and_true_then_neutral", "(y eq 5) and (x eq 10)", rule.D{"y": 5}, true},
	// and: false-then-neutral (short-circuit) -> false
	{"lenient_and_false_then_neutral_sc", "(y eq 5) and (x eq 10)", rule.D{"y": 9}, false},
	// or: neutral short-circuits to true
	{"lenient_or_neutral_short_circuit", "(x eq 10) or (y eq 5)", rule.D{"y": 9}, true},
	// or: true-then-irrelevant -> true
	{"lenient_or_true_then_neutral", "(y eq 5) or (x eq 10)", rule.D{"y": 5}, true},
	// not: not (neutral) -> not true -> false
	{"lenient_not_neutral", "not (x eq 10)", rule.D{}, false},
}

//nolint:gochecknoglobals // Test data
var LenientNestedTests = []Case{
	{"lenient_nested_eq_missing_chain", "a.b.c.d eq 10", rule.D{}, true},
	{"lenient_nested_ne_missing_chain", "a.b.c.d ne 10", rule.D{}, true},
	{"lenient_nested_eq_partial", "a.b eq c.d", rule.D{"a": rule.D{"b": 1}}, true},
	{"lenient_nested_eq_partial_ne", "a.b ne c.d", rule.D{"a": rule.D{"b": 1}}, true},
	{"lenient_nested_both_missing_eq", "a.b eq c.d", rule.D{}, true},
	{"lenient_nested_both_missing_ne", "a.b ne c.d", rule.D{}, true},
}

//nolint:gochecknoglobals // Test data
var LenientDateTimeTests = []Case{
	// datetime ops vs null -> neutral true
	{"lenient_dq_missing", `created_at dq "2024-01-01T00:00:00Z"`, rule.D{}, true},
	{"lenient_dn_missing", `created_at dn "2024-01-01T00:00:00Z"`, rule.D{}, true},
	{"lenient_be_missing", `created_at be "2024-01-01T00:00:00Z"`, rule.D{}, true},
	{"lenient_bq_missing", `created_at bq "2024-01-01T00:00:00Z"`, rule.D{}, true},
	{"lenient_af_missing", `created_at af "2024-01-01T00:00:00Z"`, rule.D{}, true},
	{"lenient_aq_missing", `created_at aq "2024-01-01T00:00:00Z"`, rule.D{}, true},
	{"lenient_dl_missing", "created_at dl 30", rule.D{}, true},
	{"lenient_dg_missing", "created_at dg 30", rule.D{}, true},
	// both missing -> true
	{"lenient_dq_both_missing", "a dq b", rule.D{}, true},
	// present still works (regular path)
	{"lenient_af_present", `created_at af "2024-01-01T00:00:00Z"`, rule.D{
		"created_at": "2024-06-01T00:00:00Z",
	}, true},
	{"lenient_bq_present_false", `created_at bq "2024-01-01T00:00:00Z"`, rule.D{
		"created_at": "2024-06-01T00:00:00Z",
	}, false},
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
	// lenientCompare, so the comparison is neutral -> true.
	{"lenient_nil_vs_absent_eq", "x eq y", rule.D{"x": nil}, true},
	{"lenient_nil_vs_absent_ne", "x ne y", rule.D{"x": nil}, true},
}
