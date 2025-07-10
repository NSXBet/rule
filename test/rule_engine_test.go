package test

import (
	"testing"

	"github.com/NSXBet/rule"
	"github.com/stretchr/testify/require"
)

func TestRulesRound1(t *testing.T) {
	all := [][]Case{
		// Core functionality tests
		EqualTests,
		RelationalTests,
		StringOpTests,
		InTests,
		PresenceTests,
		LogicalTests,
		PropCompareTests,
		NestedPropTests,

		// Comprehensive edge case tests
		EdgeCaseTests,
		ComplexLogicalTests,
		ComplexNestedLogicTests,
		RealWorldTests,
		RealWorldEdgeTests,

		// String and numeric edge cases
		StringEdgeCaseTests,
		NumericEdgeCaseTests,
		SpecialNumericTests,

		// Array and boundary tests
		ArrayEdgeCaseTests,
		BoundaryConditionTests,
		ExtremeValueTests,

		// Performance and stress tests
		PerformanceStressTests,
		TypeCoercionStressTests,

		// Whitespace and formatting tests
		WhitespaceTests,

		// Advanced pattern tests
		ComplexStringPatternTests,
		AdvancedPrecedenceTests,

		// Presence edge cases
		PresenceEdgeCaseTests,

		// Error boundary tests
		ErrorBoundaryTests,

		// Additional edge case tests
		AdditionalEdgeCaseTests,

		// Comprehensive datetime tests
		DateTimeComprehensiveTests,
	}

	for _, group := range all {
		engine := rule.NewEngine() // Fresh engine per test group

		for _, tc := range group {
			t.Run(tc.Name, func(t *testing.T) {
				got, err := engine.Evaluate(tc.Query, tc.Ctx)
				require.NoError(t, err, "query=%q", tc.Query)
				require.Equal(t, tc.Result, got, "query=%q", tc.Query)
			})
		}
	}
}
