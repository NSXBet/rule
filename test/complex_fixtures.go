package test

import "github.com/NSXBet/rule"

/* ---------- Complex Logical Combinations ---------- */

//nolint:gochecknoglobals // Test data
var ComplexLogicalTests = []Case{
	// Multiple levels of nesting
	{"complex_nested_logic", "((a.x eq b.y) and (c.z gt d.w)) or ((e.f co g.h) and (i.j in k.l))", rule.D{
		"a": rule.D{"x": 10},
		"b": rule.D{"y": 10},
		"c": rule.D{"z": 20},
		"d": rule.D{"w": 15},
		"e": rule.D{"f": "hello world"},
		"g": rule.D{"h": "world"},
		"i": rule.D{"j": "admin"},
		"k": rule.D{"l": []any{"admin", "user"}},
	}, true},

	// Chained comparisons
	{"chained_comparisons", "a eq b and b eq c and c eq d", rule.D{
		"a": 42,
		"b": 42,
		"c": 42,
		"d": 42,
	}, true},
	{"chained_comparisons_fail", "a eq b and b eq c and c eq d", rule.D{
		"a": 42,
		"b": 42,
		"c": 42,
		"d": 43,
	}, false},

	// Complex NOT operations
	{"complex_not_nested", "not (a.b eq c.d and e.f gt g.h)", rule.D{
		"a": rule.D{"b": 10},
		"c": rule.D{"d": 10},
		"e": rule.D{"f": 5},
		"g": rule.D{"h": 15},
	}, true},
	{"double_not", "not (not (x eq y))", rule.D{"x": 42, "y": 42}, true},
}

/* ---------- Real-world Scenarios ---------- */

//nolint:gochecknoglobals // Test data
var RealWorldTests = []Case{
	// User permissions
	{
		"user_permissions",
		`user.role in permissions.roles and user.active eq true and user.profile.verified eq true`,
		rule.D{
			"user": rule.D{
				"role":   "admin",
				"active": true,
				"profile": rule.D{
					"verified": true,
				},
			},
			"permissions": rule.D{
				"roles": []any{"admin", "moderator"},
			},
		},
		true,
	},

	// API rate limiting
	{"rate_limiting", `request.count lt limits.max and request.user.tier eq "premium"`, rule.D{
		"request": rule.D{
			"count": 50,
			"user":  rule.D{"tier": "premium"},
		},
		"limits": rule.D{"max": 100},
	}, true},

	// Feature flags
	{
		"feature_flags",
		`user.id in features.beta_users or (user.plan eq "enterprise" and features.enterprise_enabled eq true)`,
		rule.D{
			"user": rule.D{
				"id":   12345,
				"plan": "enterprise",
			},
			"features": rule.D{
				"beta_users":         []any{11111, 22222, 33333},
				"enterprise_enabled": true,
			},
		},
		true,
	},

	// Configuration matching
	{
		"config_matching",
		`env.stage eq "production" and config.debug eq false and config.monitoring.enabled eq true`,
		rule.D{
			"env": rule.D{"stage": "production"},
			"config": rule.D{
				"debug": false,
				"monitoring": rule.D{
					"enabled": true,
				},
			},
		},
		true,
	},

	// Content filtering
	{
		"content_filtering",
		`post.author.reputation gt 100 and not (post.content co "spam") and post.tags.category ne "adult"`,
		rule.D{
			"post": rule.D{
				"author":  rule.D{"reputation": 150},
				"content": "This is a legitimate post",
				"tags":    rule.D{"category": "general"},
			},
		},
		true,
	},
}

/* ---------- Error & Boundary Cases ---------- */

//nolint:gochecknoglobals // Test data
var ErrorBoundaryTests = []Case{
	// Missing nested attributes should return false (not error)
	{"missing_nested_left", "missing.attr eq 10", rule.D{}, false},
	{"missing_nested_right", "10 eq missing.attr", rule.D{}, false},
	{"missing_nested_both", "missing1.attr eq missing2.attr", rule.D{}, false},

	// Deeply nested missing attributes
	{"deep_missing_chain", "a.b.c.d.e.f eq 10", rule.D{
		"a": rule.D{"b": rule.D{"c": rule.D{}}},
	}, false},

	// Type mismatches in nested structures
	{"nested_type_mismatch", "user.profile.age eq 25", rule.D{
		"user": rule.D{"profile": "not_an_object"},
	}, false},

	// Boolean edge cases
	{"bool_vs_string", `active eq "true"`, rule.D{"active": true}, false},
	{"bool_vs_number", "active eq 1", rule.D{"active": true}, false},
	{"bool_false_vs_zero", "active eq 0", rule.D{"active": false}, false},

	// Array membership with wrong types
	{"array_membership_wrong_type", "value in array", rule.D{
		"value": "not_an_array",
		"array": "also_not_an_array",
	}, false},

	// Presence of deeply nested attributes
	{"presence_missing_chain", "user.profile.settings.theme pr", rule.D{
		"user": rule.D{"profile": rule.D{}},
	}, false},

	// Complex operator precedence
	{"operator_precedence", "a eq b and c eq d or e eq f", rule.D{
		"a": 1, "b": 2, "c": 3, "d": 3, "e": 5, "f": 5,
	}, true}, // Should be (a eq b and c eq d) or (e eq f) = (false and true) or true = true
}

/* ---------- Complex Nested Logic ---------- */

//nolint:gochecknoglobals // Test data
var ComplexNestedLogicTests = []Case{
	// Deeply nested logical expressions
	{
		"deeply_nested_logic",
		"((a eq b) and (c eq d)) or ((e eq f) and (g eq h)) or ((i eq j) and (k eq l))",
		rule.D{
			"a": 1, "b": 2, "c": 3, "d": 3, "e": 5, "f": 6, "g": 7, "h": 8, "i": 9, "j": 9, "k": 11, "l": 11,
		},
		true,
	},

	// Mixed operators with complex nesting
	{
		"mixed_operators_complex",
		`(user.age ge 18 and user.verified eq true) and (user.email co "@" and user.domain ne "banned.com") or user.role eq "admin"`,
		rule.D{
			"user": rule.D{
				"age":      25,
				"verified": true,
				"email":    "user@example.com",
				"domain":   "example.com",
				"role":     "user",
			},
		},
		true,
	},

	// Multiple levels of NOT
	{"multiple_nots", "not (not (not (x eq 1)))", rule.D{"x": 1}, false},
	{"not_complex_expression", "not ((a eq b and c eq d) or (e eq f and g eq h))", rule.D{
		"a": 1, "b": 1, "c": 2, "d": 2, "e": 3, "f": 4, "g": 5, "h": 6,
	}, false},
}

/* ---------- Real-world Edge Cases ---------- */

//nolint:gochecknoglobals // Test data
var RealWorldEdgeTests = []Case{
	// JSON-like nested structures (no array indexing)
	{"json_like_structure", `response.data.status eq "success" and response.data.count gt 0`, rule.D{
		"response": rule.D{
			"data": rule.D{
				"status": "success",
				"count":  2,
			},
		},
	}, true},

	// HTTP-like headers
	{"http_headers", `headers.content_type co "application/json" and headers.authorization sw "Bearer"`, rule.D{
		"headers": rule.D{
			"content_type":  "application/json; charset=utf-8",
			"authorization": "Bearer token123",
		},
	}, true},

	// Database-like queries
	{
		"database_query",
		`user.created_at ge "2023-01-01" and user.status eq "active" and user.subscription.plan in ["premium", "enterprise"]`,
		rule.D{
			"user": rule.D{
				"created_at": "2023-06-15",
				"status":     "active",
				"subscription": rule.D{
					"plan": "premium",
				},
			},
		},
		true,
	}, // String comparison "2023-06-15" ge "2023-01-01" is true (lexicographic)

	// Configuration validation
	{
		"config_validation",
		`config.database.host sw "localhost" and config.database.port ge 1024 and config.database.ssl eq true`,
		rule.D{
			"config": rule.D{
				"database": rule.D{
					"host": "localhost:5432",
					"port": 5432,
					"ssl":  true,
				},
			},
		},
		true,
	},
}
