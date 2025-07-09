package test

import (
	"testing"

	"github.com/NSXBet/rule-engine"
	"github.com/stretchr/testify/require"
)

func TestRulesRound1(t *testing.T) {
	all := [][]TestCase{
		EqualTests,
		RelationalTests,
		StringOpTests,
		InTests,
		PresenceTests,
		LogicalTests,
		PropCompareTests,
		NestedPropTests,
	}

	for _, group := range all {
		for _, tc := range group {
			t.Run(tc.Name, func(t *testing.T) {
				got, err := rule.Evaluate(tc.Query, tc.Ctx)
				require.NoError(t, err, "query=%q", tc.Query)
				require.Equal(t, tc.Result, got, "query=%q", tc.Query)
			})
		}
	}
}
