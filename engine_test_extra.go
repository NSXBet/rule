package rule

import (
	"testing"
)

// Test Engine functions beyond basic functionality
func TestEngineClearCache(t *testing.T) {
	engine := NewEngine()
	
	// Add some queries
	err := engine.AddQuery("x eq 1")
	if err != nil {
		t.Fatal(err)
	}
	
	err = engine.AddQuery("y gt 2")
	if err != nil {
		t.Fatal(err)
	}
	
	// Clear cache
	engine.ClearCache()
	
	// Should still work (will recompile)
	result, err := engine.Evaluate("x eq 1", map[string]any{"x": 1})
	if err != nil {
		t.Fatal(err)
	}
	if !result {
		t.Error("Expected true")
	}
}

func TestEngineStandaloneEvaluate(t *testing.T) {
	result, err := Evaluate("x eq 10", map[string]any{"x": 10})
	if err != nil {
		t.Fatal(err)
	}
	if !result {
		t.Error("Expected true")
	}
}

// Test engine AddQuery with parse error
func TestEngineAddQueryParseError(t *testing.T) {
	engine := NewEngine()
	
	// Try to add an invalid query
	err := engine.AddQuery("invalid query syntax [[[")
	if err == nil {
		t.Error("Expected error for invalid query syntax")
	}
}

// Test engine CompileRule with error
func TestEngineCompileRuleError(t *testing.T) {
	engine := NewEngine()
	
	// Try to compile invalid rule
	_, err := engine.CompileRule("invalid syntax [[[")
	if err == nil {
		t.Error("Expected error for invalid rule syntax")
	}
}

// Test engine Evaluate with compilation error
func TestEngineEvaluateCompilationError(t *testing.T) {
	engine := NewEngine()
	
	// Try to evaluate invalid rule
	_, err := engine.Evaluate("invalid syntax [[[", map[string]any{})
	if err == nil {
		t.Error("Expected error for invalid rule syntax")
	}
}

// Test engine with already compiled rule
func TestEngineAddQueryAlreadyCompiled(t *testing.T) {
	engine := NewEngine()
	
	// Add a query
	err := engine.AddQuery("x eq 1")
	if err != nil {
		t.Fatal(err)
	}
	
	// Add the same query again - should not error
	err = engine.AddQuery("x eq 1")
	if err != nil {
		t.Error("Adding same query twice should not error")
	}
	
	// Verify it still works
	result, err := engine.Evaluate("x eq 1", map[string]any{"x": 1})
	if err != nil {
		t.Fatal(err)
	}
	if !result {
		t.Error("Expected true")
	}
}

// Test engine with CompileRule existing rule
func TestEngineCompileRuleExisting(t *testing.T) {
	engine := NewEngine()
	
	// Compile a rule
	compiled1, err := engine.CompileRule("x eq 1")
	if err != nil {
		t.Fatal(err)
	}
	
	// Compile the same rule again - should return same compiled rule
	compiled2, err := engine.CompileRule("x eq 1")
	if err != nil {
		t.Fatal(err)
	}
	
	// Should be the same object
	if compiled1 != compiled2 {
		t.Error("Compiling same rule twice should return same object")
	}
}

// Test engine hash function
func TestEngineHash(t *testing.T) {
	// Test that hash produces different values for different strings
	hash1 := hash("rule1")
	hash2 := hash("rule2")
	
	if hash1 == hash2 {
		t.Error("Different strings should produce different hashes")
	}
	
	// Test that same string produces same hash
	hash3 := hash("rule1")
	if hash1 != hash3 {
		t.Error("Same string should produce same hash")
	}
}