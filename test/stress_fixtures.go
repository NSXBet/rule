package test

/* ---------- Whitespace and Formatting Edge Cases ---------- */

//nolint:gochecknoglobals // Test data
var WhitespaceTests = []Case{
	// Extra spaces
	{"extra_spaces", "  x    eq    10   ", map[string]any{"x": 10}, true},
	{"tabs_and_spaces", "x\teq\t10", map[string]any{"x": 10}, true},
	{"mixed_whitespace", "  x  \t eq \t 10  ", map[string]any{"x": 10}, true},

	// Parentheses with spaces
	{"spaces_in_parens", "( x eq 10 )", map[string]any{"x": 10}, true},
	{"complex_spacing", "(  x  eq  10  )  and  (  y  gt  5  )", map[string]any{"x": 10, "y": 6}, true},

	// Array formatting
	{"array_spaces", "x in [ 1 , 2 , 3 ]", map[string]any{"x": 2}, true},
	{"array_mixed_spacing", "x in [1, 2,3 ,4]", map[string]any{"x": 3}, true},

	// String with internal spaces
	{"strings_with_spaces", `name eq "John Doe"`, map[string]any{"name": "John Doe"}, true},
	{"property_with_spaces", `user.full_name eq "Jane Smith"`, map[string]any{
		"user": map[string]any{"full_name": "Jane Smith"},
	}, true},
}

/* ---------- Advanced Operator Precedence ---------- */

//nolint:gochecknoglobals // Test data
var AdvancedPrecedenceTests = []Case{
	// Complex precedence without parentheses
	{"precedence_and_or", "a eq 1 and b eq 2 or c eq 3", map[string]any{"a": 1, "b": 2, "c": 4}, true},
	{"precedence_and_or_2", "a eq 1 and b eq 2 or c eq 3", map[string]any{"a": 2, "b": 2, "c": 3}, true},
	{"precedence_and_or_3", "a eq 1 and b eq 2 or c eq 3", map[string]any{"a": 2, "b": 3, "c": 4}, false},

	// NOT with mixed operators
	{"not_with_and_or", "not a eq 1 and b eq 2", map[string]any{"a": 2, "b": 2}, true},
	{"not_with_and_or_2", "not a eq 1 or b eq 2", map[string]any{"a": 1, "b": 3}, false},
	{"not_with_and_or_3", "not a eq 1 or b eq 2", map[string]any{"a": 2, "b": 3}, true},

	// Complex chained operations
	{
		"long_chain",
		"a eq 1 and b eq 2 and c eq 3 or d eq 4 and e eq 5",
		map[string]any{"a": 1, "b": 2, "c": 3, "d": 0, "e": 0},
		true,
	},
	{
		"long_chain_2",
		"a eq 1 and b eq 2 and c eq 3 or d eq 4 and e eq 5",
		map[string]any{"a": 0, "b": 0, "c": 0, "d": 4, "e": 5},
		true,
	},
	{
		"long_chain_3",
		"a eq 1 and b eq 2 and c eq 3 or d eq 4 and e eq 5",
		map[string]any{"a": 0, "b": 0, "c": 0, "d": 0, "e": 0},
		false,
	},
}

/* ---------- Type Coercion Stress Tests ---------- */

//nolint:gochecknoglobals // Test data
var TypeCoercionStressTests = []Case{
	// Numeric precision edge cases
	{"float_precision", "x eq 0.1000000000000001", map[string]any{"x": 0.1000000000000001}, true},
	{"float_precision_2", "x eq 0.1000000000000001", map[string]any{"x": 0.1}, false},

	// Large number comparisons
	{"large_int_comparison", "x gt 9223372036854775806", map[string]any{"x": int64(9223372036854775807)}, true},
	{"large_float_comparison", "x lt 1000000000000000000", map[string]any{"x": 999999999999999999}, true},

	// String numeric comparisons
	{
		"string_numeric_mixed",
		`version ge "1.10.0"`,
		map[string]any{"version": "1.9.0"},
		true,
	}, // lexicographic: "1.9.0" ge "1.10.0" is true
	{
		"string_numeric_mixed_2",
		`version ge "1.10.0"`,
		map[string]any{"version": "1.2.0"},
		true,
	}, // lexicographic: "1.2.0" ge "1.10.0" is true

	// Boolean edge cases
	{"bool_string_strict", `flag eq "true"`, map[string]any{"flag": "true"}, true},
	{"bool_string_strict_2", `flag eq "true"`, map[string]any{"flag": true}, false},
	{"bool_int_strict", "flag eq 1", map[string]any{"flag": 1}, true},
	{"bool_int_strict_2", "flag eq 1", map[string]any{"flag": true}, false},
}

/* ---------- Performance Stress Tests ---------- */

//nolint:gochecknoglobals // Test data
var PerformanceStressTests = []Case{
	// Long chains of operations
	{
		"long_and_chain",
		"a eq 1 and b eq 2 and c eq 3 and d eq 4 and e eq 5 and f eq 6 and g eq 7 and h eq 8 and i eq 9 and j eq 10",
		map[string]any{
			"a": 1, "b": 2, "c": 3, "d": 4, "e": 5, "f": 6, "g": 7, "h": 8, "i": 9, "j": 10,
		},
		true,
	},

	{
		"long_or_chain",
		"a eq 0 or b eq 0 or c eq 0 or d eq 0 or e eq 0 or f eq 0 or g eq 0 or h eq 0 or i eq 0 or j eq 1",
		map[string]any{
			"a": 1, "b": 2, "c": 3, "d": 4, "e": 5, "f": 6, "g": 7, "h": 8, "i": 9, "j": 1,
		},
		true,
	},

	// Deeply nested property chains
	{"deep_property_chain", "a.b.c.d.e.f.g.h.i.j.k.l.m.n.o.p.q.r.s.t eq 42", map[string]any{
		"a": map[string]any{
			"b": map[string]any{
				"c": map[string]any{
					"d": map[string]any{
						"e": map[string]any{
							"f": map[string]any{
								"g": map[string]any{
									"h": map[string]any{
										"i": map[string]any{
											"j": map[string]any{
												"k": map[string]any{
													"l": map[string]any{
														"m": map[string]any{
															"n": map[string]any{
																"o": map[string]any{
																	"p": map[string]any{
																		"q": map[string]any{
																			"r": map[string]any{
																				"s": map[string]any{
																					"t": 42,
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}, true},
}

/* ---------- Boundary Conditions ---------- */

//nolint:gochecknoglobals // Test data
var BoundaryConditionTests = []Case{
	// Empty contexts
	{"empty_context_missing", "x eq 10", map[string]any{}, false},
	{"empty_context_presence", "x pr", map[string]any{}, false},

	// Null/nil values in context
	{"nil_value_comparison", "x eq null", map[string]any{"x": nil}, false}, // nil != "null"
	{"nil_value_presence", "x pr", map[string]any{"x": nil}, true},

	// Empty string vs missing
	{"empty_vs_missing", `x eq ""`, map[string]any{"x": ""}, true},
	{"empty_vs_missing_2", `x eq ""`, map[string]any{}, false},

	// Zero vs missing
	{"zero_vs_missing", "x eq 0", map[string]any{"x": 0}, true},
	{"zero_vs_missing_2", "x eq 0", map[string]any{}, false},

	// False vs missing
	{"false_vs_missing", "x eq false", map[string]any{"x": false}, true},
	{"false_vs_missing_2", "x eq false", map[string]any{}, false},
}

/* ---------- Special Numeric Values ---------- */

//nolint:gochecknoglobals // Test data
var SpecialNumericTests = []Case{
	// Very small numbers (using regular notation)
	{"very_small_positive", "x gt 0", map[string]any{"x": 1e-100}, true},
	{"very_small_negative", "x lt 0", map[string]any{"x": -1e-100}, true},

	// Boundary values - integers only (no scientific notation)
	{"max_int64", "x eq 9223372036854775807", map[string]any{"x": int64(9223372036854775807)}, true},
	{"min_int64", "x eq -9223372036854775808", map[string]any{"x": int64(-9223372036854775808)}, true},
}
