package test

import "github.com/NSXBet/rule"

/* ---------- Complex String Patterns ---------- */

//nolint:gochecknoglobals // Test data
var ComplexStringPatternTests = []Case{
	// Email-like patterns
	{
		"email_validation",
		`email co "@" and email co "." and email sw "user"`,
		rule.D{"email": "user@example.com"},
		true,
	},
	{
		"email_validation_2",
		`email co "@" and email co "." and email sw "admin"`,
		rule.D{"email": "user@example.com"},
		false,
	},

	// URL-like patterns
	{
		"url_validation",
		`url sw "https://" and url co "example.com" and url ew "/path"`,
		rule.D{"url": "https://example.com/path"},
		true,
	},
	{
		"url_validation_2",
		`url sw "http://" and url co "example.com" and url ew ".html"`,
		rule.D{"url": "http://example.com/page.html"},
		true,
	},

	// JSON-like string patterns
	{"json_pattern", `data sw "{" and data ew "}" and data co "key"`, rule.D{"data": `{"key": "value"}`}, true},

	// SQL-like patterns
	{
		"sql_pattern",
		`query sw "SELECT" and query co "FROM" and query co "WHERE"`,
		rule.D{"query": "SELECT * FROM users WHERE active = true"},
		true,
	},

	// File path patterns
	{
		"file_path_unix",
		`path sw "/" and path co "home" and path ew ".txt"`,
		rule.D{"path": "/home/user/document.txt"},
		true,
	},
	{
		"file_path_windows",
		`path sw "C:" and path co "Users" and path ew ".doc"`,
		rule.D{"path": "C:\\Users\\John\\file.doc"},
		true,
	},
}
