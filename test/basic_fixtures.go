package test

import "github.com/NSXBet/rule"

// Case bundles one rule, its context, and the expected boolean result.
type Case struct {
	Name   string
	Query  string
	Ctx    rule.D
	Result bool
}

/* ---------- Equality & inequality ---------- */

//nolint:gochecknoglobals // Test data
var EqualTests = []Case{
	// ints
	{"eq_int_true", "x eq 1", rule.D{"x": 1}, true},
	{"eq_int_false", "x eq 1", rule.D{"x": 2}, false},
	// floats
	{"eq_float_true", "x eq 1.23", rule.D{"x": 1.23}, true},
	{"eq_float_false", "x eq 1.23", rule.D{"x": 1.24}, false},
	// strings
	{"eq_string_true", `x eq "abc"`, rule.D{"x": "abc"}, true},
	{"eq_string_false", `x eq "abc"`, rule.D{"x": "def"}, false},
	// versions
	{"eq_version_true", `version eq "1.2.3"`, rule.D{"version": "1.2.3"}, true},
	{"eq_version_false", `version eq "1.2.3"`, rule.D{"version": "1.2.4"}, false},
	// nested attr
	{
		"eq_nested_true",
		"x.y eq 2",
		rule.D{"x": rule.D{"y": 2}},
		true,
	},
	{
		"eq_nested_false",
		"x.y eq 2",
		rule.D{"x": rule.D{"y": 3}},
		false,
	},
	// alias ==
	{"eq_alias_true", "x == 99", rule.D{"x": 99}, true},
	{"ne_int_true", "x ne 1", rule.D{"x": 2}, true},
	{"ne_alias_true", "x != 42", rule.D{"x": 1}, true},
}

/* ---------- Relational (<, >, <=, >=) ---------- */

//nolint:gochecknoglobals // Test data
var RelationalTests = []Case{
	{"lt_true", "score lt 10", rule.D{"score": 5}, true},
	{"lt_false", "score lt 10", rule.D{"score": 10}, false},
	{"gt_true", "score gt 10", rule.D{"score": 11}, true},
	{"gt_false", "score gt 10", rule.D{"score": 9}, false},
	{"le_true_equal", "score le 7", rule.D{"score": 7}, true},
	{"le_true_less", "score le 7", rule.D{"score": 6}, true},
	{"ge_true_equal", "score ge 7", rule.D{"score": 7}, true},
	{"ge_true_greater", "score ge 7", rule.D{"score": 8}, true},
}

/* ---------- String operations (co, sw, ew) ---------- */

//nolint:gochecknoglobals // Test data
var StringOpTests = []Case{
	{"co_true", `city co "York"`, rule.D{"city": "New York"}, true},
	{"co_false", `city co "York"`, rule.D{"city": "Boston"}, false},
	{"sw_true", `id sw "user_"`, rule.D{"id": "user_123"}, true},
	{"sw_false", `id sw "user_"`, rule.D{"id": "admin_1"}, false},
	{"ew_true", `file ew ".txt"`, rule.D{"file": "report.txt"}, true},
	{"ew_false", `file ew ".txt"`, rule.D{"file": "image.png"}, false},
}

/* ---------- Membership (in) ---------- */

//nolint:gochecknoglobals // Test data
var InTests = []Case{
	{"in_int_true", "x in [1,2,3]", rule.D{"x": 2}, true},
	{"in_int_false", "x in [1,2,3]", rule.D{"x": 4}, false},
	{
		"in_str_true",
		`color in ["red","green","blue"]`,
		rule.D{"color": "green"},
		true,
	},
	{
		"in_str_false",
		`color in ["red","green","blue"]`,
		rule.D{"color": "yellow"},
		false,
	},
}

/* ---------- Presence (pr) ---------- */

//nolint:gochecknoglobals // Test data
var PresenceTests = []Case{
	{"pr_present", "betaUser pr", rule.D{"betaUser": true}, true},
	{"pr_missing", "betaUser pr", rule.D{}, false},
}

/* ---------- NOT, AND, OR, nesting ---------- */

//nolint:gochecknoglobals // Test data
var LogicalTests = []Case{
	{"not_true", "not (x eq 1)", rule.D{"x": 2}, true},
	{"not_false", "not (x eq 1)", rule.D{"x": 1}, false},
	{"and_true", "(x gt 1) and (y lt 5)", rule.D{"x": 2, "y": 3}, true},
	{"and_false_left", "(x gt 1) and (y lt 5)", rule.D{"x": 1, "y": 3}, false},
	{"or_true_left", "(x lt 0) or (y eq 7)", rule.D{"x": -1, "y": 9}, true},
	{"or_true_right", "(x lt 0) or (y eq 7)", rule.D{"x": 1, "y": 7}, true},
	{"or_false", "(x lt 0) or (y eq 7)", rule.D{"x": 1, "y": 8}, false},
}
