package test

/* ---------- Complex Logical Combinations ---------- */

var ComplexLogicalTests = []TestCase{
	// Multiple levels of nesting
	{"complex_nested_logic", "((a.x eq b.y) and (c.z gt d.w)) or ((e.f co g.h) and (i.j in k.l))", map[string]any{
		"a": map[string]any{"x": 10},
		"b": map[string]any{"y": 10},
		"c": map[string]any{"z": 20},
		"d": map[string]any{"w": 15},
		"e": map[string]any{"f": "hello world"},
		"g": map[string]any{"h": "world"},
		"i": map[string]any{"j": "admin"},
		"k": map[string]any{"l": []any{"admin", "user"}},
	}, true},

	// Chained comparisons
	{"chained_comparisons", "a eq b and b eq c and c eq d", map[string]any{
		"a": 42,
		"b": 42,
		"c": 42,
		"d": 42,
	}, true},
	{"chained_comparisons_fail", "a eq b and b eq c and c eq d", map[string]any{
		"a": 42,
		"b": 42,
		"c": 42,
		"d": 43,
	}, false},

	// Complex NOT operations
	{"complex_not_nested", "not (a.b eq c.d and e.f gt g.h)", map[string]any{
		"a": map[string]any{"b": 10},
		"c": map[string]any{"d": 10},
		"e": map[string]any{"f": 5},
		"g": map[string]any{"h": 15},
	}, true},
	{"double_not", "not (not (x eq y))", map[string]any{"x": 42, "y": 42}, true},
}

/* ---------- Real-world Scenarios ---------- */

var RealWorldTests = []TestCase{
	// User permissions
	{"user_permissions", `user.role in permissions.roles and user.active eq true and user.profile.verified eq true`, map[string]any{
		"user": map[string]any{
			"role":   "admin",
			"active": true,
			"profile": map[string]any{
				"verified": true,
			},
		},
		"permissions": map[string]any{
			"roles": []any{"admin", "moderator"},
		},
	}, true},

	// API rate limiting
	{"rate_limiting", `request.count lt limits.max and request.user.tier eq "premium"`, map[string]any{
		"request": map[string]any{
			"count": 50,
			"user":  map[string]any{"tier": "premium"},
		},
		"limits": map[string]any{"max": 100},
	}, true},

	// Feature flags
	{"feature_flags", `user.id in features.beta_users or (user.plan eq "enterprise" and features.enterprise_enabled eq true)`, map[string]any{
		"user": map[string]any{
			"id":   12345,
			"plan": "enterprise",
		},
		"features": map[string]any{
			"beta_users":         []any{11111, 22222, 33333},
			"enterprise_enabled": true,
		},
	}, true},

	// Configuration matching
	{"config_matching", `env.stage eq "production" and config.debug eq false and config.monitoring.enabled eq true`, map[string]any{
		"env": map[string]any{"stage": "production"},
		"config": map[string]any{
			"debug": false,
			"monitoring": map[string]any{
				"enabled": true,
			},
		},
	}, true},

	// Content filtering
	{"content_filtering", `post.author.reputation gt 100 and not (post.content co "spam") and post.tags.category ne "adult"`, map[string]any{
		"post": map[string]any{
			"author": map[string]any{"reputation": 150},
			"content": "This is a legitimate post",
			"tags":    map[string]any{"category": "general"},
		},
	}, true},
}

/* ---------- Error & Boundary Cases ---------- */

var ErrorBoundaryTests = []TestCase{
	// Missing nested attributes should return false (not error)
	{"missing_nested_left", "missing.attr eq 10", map[string]any{}, false},
	{"missing_nested_right", "10 eq missing.attr", map[string]any{}, false},
	{"missing_nested_both", "missing1.attr eq missing2.attr", map[string]any{}, false},
	
	// Deeply nested missing attributes
	{"deep_missing_chain", "a.b.c.d.e.f eq 10", map[string]any{
		"a": map[string]any{"b": map[string]any{"c": map[string]any{}}},
	}, false},
	
	// Type mismatches in nested structures
	{"nested_type_mismatch", "user.profile.age eq 25", map[string]any{
		"user": map[string]any{"profile": "not_an_object"},
	}, false},
	
	// Boolean edge cases
	{"bool_vs_string", `active eq "true"`, map[string]any{"active": true}, false},
	{"bool_vs_number", "active eq 1", map[string]any{"active": true}, false},
	{"bool_false_vs_zero", "active eq 0", map[string]any{"active": false}, false},
	
	// Array membership with wrong types
	{"array_membership_wrong_type", "value in array", map[string]any{
		"value": "not_an_array",
		"array": "also_not_an_array",
	}, false},
	
	// Presence of deeply nested attributes
	{"presence_missing_chain", "user.profile.settings.theme pr", map[string]any{
		"user": map[string]any{"profile": map[string]any{}},
	}, false},
	
	// Complex operator precedence
	{"operator_precedence", "a eq b and c eq d or e eq f", map[string]any{
		"a": 1, "b": 2, "c": 3, "d": 3, "e": 5, "f": 5,
	}, true}, // Should be (a eq b and c eq d) or (e eq f) = (false and true) or true = true
}

/* ---------- Complex Nested Logic ---------- */

var ComplexNestedLogicTests = []TestCase{
	// Deeply nested logical expressions
	{"deeply_nested_logic", "((a eq b) and (c eq d)) or ((e eq f) and (g eq h)) or ((i eq j) and (k eq l))", map[string]any{
		"a": 1, "b": 2, "c": 3, "d": 3, "e": 5, "f": 6, "g": 7, "h": 8, "i": 9, "j": 9, "k": 11, "l": 11,
	}, true},
	
	// Mixed operators with complex nesting
	{"mixed_operators_complex", `(user.age ge 18 and user.verified eq true) and (user.email co "@" and user.domain ne "banned.com") or user.role eq "admin"`, map[string]any{
		"user": map[string]any{
			"age":      25,
			"verified": true,
			"email":    "user@example.com",
			"domain":   "example.com",
			"role":     "user",
		},
	}, true},
	
	// Multiple levels of NOT
	{"multiple_nots", "not (not (not (x eq 1)))", map[string]any{"x": 1}, false},
	{"not_complex_expression", "not ((a eq b and c eq d) or (e eq f and g eq h))", map[string]any{
		"a": 1, "b": 1, "c": 2, "d": 2, "e": 3, "f": 4, "g": 5, "h": 6,
	}, false},
}

/* ---------- Real-world Edge Cases ---------- */

var RealWorldEdgeTests = []TestCase{
	// JSON-like nested structures (no array indexing)
	{"json_like_structure", `response.data.status eq "success" and response.data.count gt 0`, map[string]any{
		"response": map[string]any{
			"data": map[string]any{
				"status": "success",
				"count":  2,
			},
		},
	}, true},
	
	// HTTP-like headers
	{"http_headers", `headers.content_type co "application/json" and headers.authorization sw "Bearer"`, map[string]any{
		"headers": map[string]any{
			"content_type":  "application/json; charset=utf-8",
			"authorization": "Bearer token123",
		},
	}, true},
	
	// Database-like queries
	{"database_query", `user.created_at ge "2023-01-01" and user.status eq "active" and user.subscription.plan in ["premium", "enterprise"]`, map[string]any{
		"user": map[string]any{
			"created_at": "2023-06-15",
			"status":     "active",
			"subscription": map[string]any{
				"plan": "premium",
			},
		},
	}, true}, // String comparison "2023-06-15" ge "2023-01-01" is true (lexicographic)
	
	// Configuration validation
	{"config_validation", `config.database.host sw "localhost" and config.database.port ge 1024 and config.database.ssl eq true`, map[string]any{
		"config": map[string]any{
			"database": map[string]any{
				"host": "localhost:5432",
				"port": 5432,
				"ssl":  true,
			},
		},
	}, true},
}