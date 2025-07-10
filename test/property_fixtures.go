package test

/* ---------- Attribute-to-attribute (flat) ---------- */

//nolint:gochecknoglobals // Test data
var PropCompareTests = []TestCase{
	// eq / ne on numbers
	{"prop_eq_int_true", "x eq y", map[string]any{"x": 5, "y": 5}, true},
	{"prop_eq_int_false", "x eq y", map[string]any{"x": 5, "y": 6}, false},
	{"prop_ne_int_true", "x ne y", map[string]any{"x": 1, "y": 2}, true},

	// relational ops on numbers
	{"prop_gt_true", "a gt b", map[string]any{"a": 10, "b": 3}, true},
	{"prop_lt_false", "a lt b", map[string]any{"a": 10, "b": 3}, false},
	{"prop_le_true", "a le b", map[string]any{"a": 7, "b": 7}, true},
	{"prop_ge_false", "a ge b", map[string]any{"a": 6, "b": 7}, false},

	// eq / ne on strings
	{"prop_eq_str_true", `first eq last`, map[string]any{"first": "go", "last": "go"}, true},
	{"prop_ne_str_true", `first ne last`, map[string]any{"first": "go", "last": "rust"}, true},

	// string ops (co / sw / ew)
	{"prop_co_true", `email co domain`, map[string]any{
		"email":  "foo@example.com",
		"domain": "@example.com",
	}, true},
	{"prop_sw_false", `login sw prefix`, map[string]any{
		"login":  "guest123",
		"prefix": "user_",
	}, false},
	{"prop_ew_true", `file ew ext`, map[string]any{
		"file": "report.pdf",
		"ext":  ".pdf",
	}, true},

	// membership: an attribute that is a list
	{"prop_in_list_true", `color in allowed`, map[string]any{
		"color":   "red",
		"allowed": []any{"red", "green"},
	}, true},
	{"prop_in_list_false", `color in allowed`, map[string]any{
		"color":   "blue",
		"allowed": []any{"red", "green"},
	}, false},
}

/* ---------- Nested-attribute comparisons ---------- */

//nolint:gochecknoglobals // Test data
var NestedPropTests = []TestCase{
	// nested vs nested equality
	{"nested_eq_true", "x.y eq z.w", map[string]any{
		"x": map[string]any{"y": 10},
		"z": map[string]any{"w": 10},
	}, true},
	{"nested_eq_false", "x.y eq z.w", map[string]any{
		"x": map[string]any{"y": 10},
		"z": map[string]any{"w": 11},
	}, false},

	// nested vs top-level
	{"mixed_depth_eq_true", "x.y eq z", map[string]any{
		"x": map[string]any{"y": "abc"},
		"z": "abc",
	}, true},

	// relational on nested numbers
	{"nested_gt_true", "m.n gt p.q", map[string]any{
		"m": map[string]any{"n": 20},
		"p": map[string]any{"q": 15},
	}, true},
	{"nested_le_false", "m.n le p.q", map[string]any{
		"m": map[string]any{"n": 21},
		"p": map[string]any{"q": 20},
	}, false},
}
