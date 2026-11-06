package test

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/NSXBet/rule"
	"github.com/stretchr/testify/require"
)

func TestUnwrapExpansion(t *testing.T) {
	t.Parallel()

	context := rule.D{
		"data": rule.D{
			"items": []any{
				rule.D{"id": "root", "score": 10},
				rule.D{
					"id":    "wrapper-a",
					"score": 100,
					"children": []any{
						rule.D{"id": "child-a", "score": 20},
						rule.D{"id": "child-b", "score": 5},
					},
				},
				rule.D{
					"id":       "wrapper-b",
					"score":    100,
					"children": []any{rule.D{"id": "child-c", "score": 30}},
				},
			},
		},
	}
	queries := []struct {
		name  string
		query string
	}{
		{"source has three entries", "data.items.length eq 3"},
		{"expansion has root plus three children", "data.items unwrap children where (true).length eq 4"},
		{"ordinary root remains", `data.items unwrap children where (id eq "root").length eq 1`},
		{
			"both wrappers are excluded",
			`data.items unwrap children where (id in ["wrapper-a","wrapper-b"]).length eq 0`,
		},
		{
			"all immediate children participate",
			`data.items unwrap children where (id in ["child-a","child-b","child-c"]).length eq 3`,
		},
		{"predicate matches root and two children", "data.items unwrap children where (score ge 10).length eq 3"},
		{"wrapper score is not a fallback", "data.items unwrap children where (score ge 50).length eq 0"},
	}

	engine := rule.NewEngine()

	for _, tc := range queries {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			before, err := json.Marshal(context)
			require.NoError(t, err)

			got, err := engine.Evaluate(tc.query, context)
			require.NoError(t, err)
			require.True(t, got)

			after, err := json.Marshal(context)
			require.NoError(t, err)
			require.Equal(t, before, after, "unwrap must not mutate the supplied structure")
		})
	}
}

func TestUnwrapFieldIsolation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		items     []any
		wantCount int
	}{
		{
			name: "children use their own conflicting fields",
			items: []any{
				rule.D{
					"score":  0,
					"status": "inactive",
					"children": []any{
						rule.D{"score": 20, "status": "active"},
						rule.D{"score": 5, "status": "active"},
					},
				},
			},
			wantCount: 1,
		},
		{
			name: "children do not inherit missing fields from wrapper",
			items: []any{
				rule.D{
					"score":  20,
					"status": "active",
					"children": []any{
						rule.D{"status": "active"},
						rule.D{"score": 20},
					},
				},
			},
			wantCount: 0,
		},
		{
			name: "matching fields do not leak to later siblings or roots",
			items: []any{
				rule.D{
					"children": []any{
						rule.D{"score": 20, "status": "active"},
						rule.D{"status": "active"},
						rule.D{"score": 20},
					},
				},
				rule.D{"score": 20},
				rule.D{"score": 10, "status": "active"},
			},
			wantCount: 2,
		},
	}

	engine := rule.NewEngine()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			query := `items unwrap children where (score ge 10 and status eq "active").length eq ` +
				strconv.Itoa(tc.wantCount)
			got, err := engine.Evaluate(query, rule.D{"items": tc.items})
			require.NoError(t, err)
			require.True(t, got)
		})
	}
}

func TestUnwrapIdentityCollisions(t *testing.T) {
	t.Parallel()

	child := rule.D{"id": "same", "score": 20}
	tests := []struct {
		name      string
		items     []any
		wantCount int
	}{
		{
			name: "wrapper and child share identity without counting wrapper twice",
			items: []any{
				rule.D{"id": "same", "score": 100, "children": []any{child}},
			},
			wantCount: 1,
		},
		{
			name: "equal sibling values remain separate occurrences",
			items: []any{
				rule.D{"children": []any{child, rule.D{"id": "same", "score": 20}}},
			},
			wantCount: 2,
		},
		{
			name: "root and child with same identity are not merged",
			items: []any{
				rule.D{"id": "same", "score": 10},
				rule.D{"children": []any{child}},
			},
			wantCount: 2,
		},
		{
			name: "shared child in different wrappers counts per occurrence",
			items: []any{
				rule.D{"children": []any{child}},
				rule.D{"children": []any{child}},
			},
			wantCount: 2,
		},
	}

	engine := rule.NewEngine()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			query := `items unwrap children where (id eq "same" and score ge 10).length eq ` +
				strconv.Itoa(tc.wantCount)
			got, err := engine.Evaluate(query, rule.D{"items": tc.items})
			require.NoError(t, err)
			require.True(t, got)
		})
	}
}

func TestUnwrapEmptyInputs(t *testing.T) {
	t.Parallel()

	for _, items := range [][]any{
		{},
		{rule.D{"score": 100, "children": []any{}}},
	} {
		got, err := rule.NewEngine().Evaluate("items unwrap children where (true).length eq 0", rule.D{"items": items})
		require.NoError(t, err)
		require.True(t, got)
	}
}

func TestUnwrapStopsAfterOneLevel(t *testing.T) {
	t.Parallel()

	context := rule.D{
		"items": []any{
			rule.D{
				"id": "wrapper",
				"children": []any{
					rule.D{
						"id":       "child",
						"score":    20,
						"children": []any{rule.D{"id": "grandchild", "score": 100}},
					},
				},
			},
		},
	}

	got, err := rule.NewEngine().Evaluate(
		`items unwrap children where (id eq "child").length eq 1 and
		items unwrap children where (score ge 50).length eq 0`,
		context,
	)
	require.NoError(t, err)
	require.True(t, got)
}

func TestUnwrapSyntaxErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		query   string
		wantErr error
	}{
		{"items unwrap where (true).length eq 0", rule.ErrInvalidUnwrapField},
		{"items unwrap 123 where (true).length eq 0", rule.ErrInvalidUnwrapField},
		{"items unwrap children", rule.ErrWhereRequiresLengthOrQuantifier},
		{"items unwrap children where ().length eq 0", rule.ErrEmptyWherePredicate},
		{"items unwrap children where true", rule.ErrWhereRequiresParens},
		{"items unwrap children where (true).other eq 0", rule.ErrWhereRequiresLength},
		{"items unwrap children where (true) any (true)", rule.ErrWhereRequiresLengthOrQuantifier},
		{"items unwrap children where (true) all (true)", rule.ErrWhereRequiresLengthOrQuantifier},
		{"items unwrap children where (true) none (true)", rule.ErrWhereRequiresLengthOrQuantifier},
	}

	engine := rule.NewEngine()

	for _, tc := range tests {
		t.Run(tc.query, func(t *testing.T) {
			t.Parallel()

			_, err := engine.Evaluate(tc.query, rule.D{})
			require.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestUnwrapPropagatesChildPredicateErrors(t *testing.T) {
	t.Parallel()

	got, err := rule.NewEngine().Evaluate(
		"items unwrap children where (round(score, 2) ge 10).length eq 0",
		rule.D{
			"items": []any{
				rule.D{"children": []any{rule.D{"score": "invalid"}}},
			},
		},
	)
	require.ErrorIs(t, err, rule.ErrInvalidMathOperand)
	require.False(t, got)
}

func TestUnwrapKeywordCutover(t *testing.T) {
	t.Parallel()

	context := rule.D{
		"items": []any{
			rule.D{"score": 10},
			rule.D{"children": []any{rule.D{"score": 20}}},
		},
		"flatten": true,
		"payload": rule.D{"flatten": true},
	}
	tests := []struct {
		name      string
		query     string
		wantError bool
	}{
		{"unwrap replaces wrappers", "items unwrap children where (score ge 10).length eq 2", false},
		{"flatten operator is retired", "items flatten children where (true).length eq 2", true},
		{"flatten is an ordinary field name", "flatten eq true and payload.flatten eq true", false},
	}

	engine := rule.NewEngine()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := engine.Evaluate(tc.query, context)
			if tc.wantError {
				require.Error(t, err)
				require.False(t, got)

				return
			}

			require.NoError(t, err)
			require.True(t, got)
		})
	}
}
