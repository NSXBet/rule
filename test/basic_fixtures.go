package test

// Case bundles one rule, its context, and the expected boolean result.
type Case struct {
	Name   string
	Query  string
	Ctx    map[string]any
	Result bool
}

/* ---------- Equality & inequality ---------- */

//nolint:gochecknoglobals // Test data
var EqualTests = []Case{
	// ints
	{"eq_int_true", "x eq 1", map[string]any{"x": 1}, true},
	{"eq_int_false", "x eq 1", map[string]any{"x": 2}, false},
	// floats
	{"eq_float_true", "x eq 1.23", map[string]any{"x": 1.23}, true},
	{"eq_float_false", "x eq 1.23", map[string]any{"x": 1.24}, false},
	// strings
	{"eq_string_true", `x eq "abc"`, map[string]any{"x": "abc"}, true},
	{"eq_string_false", `x eq "abc"`, map[string]any{"x": "def"}, false},
	// versions
	{"eq_version_true", `version eq "1.2.3"`, map[string]any{"version": "1.2.3"}, true},
	{"eq_version_false", `version eq "1.2.3"`, map[string]any{"version": "1.2.4"}, false},
	// nested attr
	{
		"eq_nested_true",
		"x.y eq 2",
		map[string]any{"x": map[string]any{"y": 2}},
		true,
	},
	{
		"eq_nested_false",
		"x.y eq 2",
		map[string]any{"x": map[string]any{"y": 3}},
		false,
	},
	// alias ==
	{"eq_alias_true", "x == 99", map[string]any{"x": 99}, true},
	{"ne_int_true", "x ne 1", map[string]any{"x": 2}, true},
	{"ne_alias_true", "x != 42", map[string]any{"x": 1}, true},
}

/* ---------- Relational (<, >, <=, >=) ---------- */

//nolint:gochecknoglobals // Test data
var RelationalTests = []Case{
	{"lt_true", "score lt 10", map[string]any{"score": 5}, true},
	{"lt_false", "score lt 10", map[string]any{"score": 10}, false},
	{"gt_true", "score gt 10", map[string]any{"score": 11}, true},
	{"gt_false", "score gt 10", map[string]any{"score": 9}, false},
	{"le_true_equal", "score le 7", map[string]any{"score": 7}, true},
	{"le_true_less", "score le 7", map[string]any{"score": 6}, true},
	{"ge_true_equal", "score ge 7", map[string]any{"score": 7}, true},
	{"ge_true_greater", "score ge 7", map[string]any{"score": 8}, true},
}

/* ---------- String operations (co, sw, ew) ---------- */

//nolint:gochecknoglobals // Test data
var StringOpTests = []Case{
	{"co_true", `city co "York"`, map[string]any{"city": "New York"}, true},
	{"co_false", `city co "York"`, map[string]any{"city": "Boston"}, false},
	{"sw_true", `id sw "user_"`, map[string]any{"id": "user_123"}, true},
	{"sw_false", `id sw "user_"`, map[string]any{"id": "admin_1"}, false},
	{"ew_true", `file ew ".txt"`, map[string]any{"file": "report.txt"}, true},
	{"ew_false", `file ew ".txt"`, map[string]any{"file": "image.png"}, false},
}

/* ---------- Membership (in) ---------- */

//nolint:gochecknoglobals // Test data
var InTests = []Case{
	{"in_int_true", "x in [1,2,3]", map[string]any{"x": 2}, true},
	{"in_int_false", "x in [1,2,3]", map[string]any{"x": 4}, false},
	{
		"in_str_true",
		`color in ["red","green","blue"]`,
		map[string]any{"color": "green"},
		true,
	},
	{
		"in_str_false",
		`color in ["red","green","blue"]`,
		map[string]any{"color": "yellow"},
		false,
	},
}

/* ---------- Presence (pr) ---------- */

//nolint:gochecknoglobals // Test data
var PresenceTests = []Case{
	{"pr_present", "betaUser pr", map[string]any{"betaUser": true}, true},
	{"pr_missing", "betaUser pr", map[string]any{}, false},
}

/* ---------- NOT, AND, OR, nesting ---------- */

//nolint:gochecknoglobals // Test data
var LogicalTests = []Case{
	{"not_true", "not (x eq 1)", map[string]any{"x": 2}, true},
	{"not_false", "not (x eq 1)", map[string]any{"x": 1}, false},
	{"and_true", "(x gt 1) and (y lt 5)", map[string]any{"x": 2, "y": 3}, true},
	{"and_false_left", "(x gt 1) and (y lt 5)", map[string]any{"x": 1, "y": 3}, false},
	{"or_true_left", "(x lt 0) or (y eq 7)", map[string]any{"x": -1, "y": 9}, true},
	{"or_true_right", "(x lt 0) or (y eq 7)", map[string]any{"x": 1, "y": 7}, true},
	{"or_false", "(x lt 0) or (y eq 7)", map[string]any{"x": 1, "y": 8}, false},
}
