package test

import "github.com/NSXBet/rule"

/* ---------- Attribute-to-attribute (flat) ---------- */

//nolint:gochecknoglobals // Test data
var PropCompareTests = []Case{
	// eq / ne on numbers
	{"prop_eq_int_true", "x eq y", rule.D{"x": 5, "y": 5}, true},
	{"prop_eq_int_false", "x eq y", rule.D{"x": 5, "y": 6}, false},
	{"prop_ne_int_true", "x ne y", rule.D{"x": 1, "y": 2}, true},

	// relational ops on numbers
	{"prop_gt_true", "a gt b", rule.D{"a": 10, "b": 3}, true},
	{"prop_lt_false", "a lt b", rule.D{"a": 10, "b": 3}, false},
	{"prop_le_true", "a le b", rule.D{"a": 7, "b": 7}, true},
	{"prop_ge_false", "a ge b", rule.D{"a": 6, "b": 7}, false},

	// eq / ne on strings
	{"prop_eq_str_true", `first eq last`, rule.D{"first": "go", "last": "go"}, true},
	{"prop_ne_str_true", `first ne last`, rule.D{"first": "go", "last": "rust"}, true},

	// string ops (co / sw / ew)
	{"prop_co_true", `email co domain`, rule.D{
		"email":  "foo@example.com",
		"domain": "@example.com",
	}, true},
	{"prop_sw_false", `login sw prefix`, rule.D{
		"login":  "guest123",
		"prefix": "user_",
	}, false},
	{"prop_ew_true", `file ew ext`, rule.D{
		"file": "report.pdf",
		"ext":  ".pdf",
	}, true},

	// membership: an attribute that is a list
	{"prop_in_list_true", `color in allowed`, rule.D{
		"color":   "red",
		"allowed": []any{"red", "green"},
	}, true},
	{"prop_in_list_false", `color in allowed`, rule.D{
		"color":   "blue",
		"allowed": []any{"red", "green"},
	}, false},
}

/* ---------- Nested-attribute comparisons ---------- */

//nolint:gochecknoglobals // Test data
var NestedPropTests = []Case{
	// nested vs nested equality
	{"nested_eq_true", "x.y eq z.w", rule.D{
		"x": rule.D{"y": 10},
		"z": rule.D{"w": 10},
	}, true},
	{"nested_eq_false", "x.y eq z.w", rule.D{
		"x": rule.D{"y": 10},
		"z": rule.D{"w": 11},
	}, false},

	// nested vs top-level
	{"mixed_depth_eq_true", "x.y eq z", rule.D{
		"x": rule.D{"y": "abc"},
		"z": "abc",
	}, true},

	// relational on nested numbers
	{"nested_gt_true", "m.n gt p.q", rule.D{
		"m": rule.D{"n": 20},
		"p": rule.D{"q": 15},
	}, true},
	{"nested_le_false", "m.n le p.q", rule.D{
		"m": rule.D{"n": 21},
		"p": rule.D{"q": 20},
	}, false},
}
