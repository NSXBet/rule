package main

import (
	"fmt"
	"github.com/nikunjy/rules"
)

func main() {
	// Test missing attribute behavior
	result, err := rules.Evaluate(`missing eq 10`, map[string]any{})
	fmt.Printf("missing eq 10: result=%v, error=%v\n", result, err)
	
	// Test nested missing attribute
	result2, err2 := rules.Evaluate(`missing.nested eq 10`, map[string]any{})
	fmt.Printf("missing.nested eq 10: result=%v, error=%v\n", result2, err2)
	
	// Test partially missing nested attribute
	result3, err3 := rules.Evaluate(`user.missing eq 10`, map[string]any{"user": map[string]any{}})
	fmt.Printf("user.missing eq 10: result=%v, error=%v\n", result3, err3)
}