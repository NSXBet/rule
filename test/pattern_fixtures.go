package test

/* ---------- Complex String Patterns ---------- */

//nolint:gochecknoglobals // Test data
var ComplexStringPatternTests = []Case{
	// Email-like patterns
	{
		"email_validation",
		`email co "@" and email co "." and email sw "user"`,
		map[string]any{"email": "user@example.com"},
		true,
	},
	{
		"email_validation_2",
		`email co "@" and email co "." and email sw "admin"`,
		map[string]any{"email": "user@example.com"},
		false,
	},

	// URL-like patterns
	{
		"url_validation",
		`url sw "https://" and url co "example.com" and url ew "/path"`,
		map[string]any{"url": "https://example.com/path"},
		true,
	},
	{
		"url_validation_2",
		`url sw "http://" and url co "example.com" and url ew ".html"`,
		map[string]any{"url": "http://example.com/page.html"},
		true,
	},

	// JSON-like string patterns
	{"json_pattern", `data sw "{" and data ew "}" and data co "key"`, map[string]any{"data": `{"key": "value"}`}, true},

	// SQL-like patterns
	{
		"sql_pattern",
		`query sw "SELECT" and query co "FROM" and query co "WHERE"`,
		map[string]any{"query": "SELECT * FROM users WHERE active = true"},
		true,
	},

	// File path patterns
	{
		"file_path_unix",
		`path sw "/" and path co "home" and path ew ".txt"`,
		map[string]any{"path": "/home/user/document.txt"},
		true,
	},
	{
		"file_path_windows",
		`path sw "C:" and path co "Users" and path ew ".doc"`,
		map[string]any{"path": "C:\\Users\\John\\file.doc"},
		true,
	},
}
