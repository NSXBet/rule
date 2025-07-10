package rule

import (
	"github.com/puzpuzpuz/xsync/v4"
)

// D is a type alias for map[string]any, providing a cleaner API for context data.
// Usage: rule.D{"user": rule.D{"age": 25, "active": true}}.
type D = map[string]any

type CompiledRule struct {
	AST  *ASTNode
	Hash uint64
}

type Engine struct {
	compiledRules *xsync.Map[string, *CompiledRule]
	evaluator     *Evaluator
}

func NewEngine() *Engine {
	return &Engine{
		compiledRules: xsync.NewMap[string, *CompiledRule](),
		evaluator:     NewEvaluator(),
	}
}

func (e *Engine) AddQuery(rule string) error {
	if _, exists := e.compiledRules.Load(rule); exists {
		return nil // Already compiled
	}

	ast, err := ParseRule(rule)
	if err != nil {
		return err
	}

	compiled := &CompiledRule{
		AST:  ast,
		Hash: hash(rule),
	}

	e.compiledRules.Store(rule, compiled)

	return nil
}

func (e *Engine) Evaluate(rule string, context D) (bool, error) {
	compiled, exists := e.compiledRules.Load(rule)
	if !exists {
		// Compile just-in-time
		if err := e.AddQuery(rule); err != nil {
			return false, err
		}

		compiled, _ = e.compiledRules.Load(rule)
	}

	return e.evaluator.Evaluate(compiled.AST, context)
}

func (e *Engine) CompileRule(rule string) (*CompiledRule, error) {
	if compiled, exists := e.compiledRules.Load(rule); exists {
		return compiled, nil
	}

	ast, err := ParseRule(rule)
	if err != nil {
		return nil, err
	}

	compiled := &CompiledRule{
		AST:  ast,
		Hash: hash(rule),
	}

	e.compiledRules.Store(rule, compiled)

	return compiled, nil
}

func (e *Engine) EvaluateCompiled(compiled *CompiledRule, context D) (bool, error) {
	return e.evaluator.Evaluate(compiled.AST, context)
}

func (e *Engine) ClearCache() {
	e.compiledRules.Clear()
}

func hash(s string) uint64 {
	h := uint64(hashOffsetBasis)
	for i := range len(s) {
		h = ((h << hashPrime) + h) + uint64(s[i])
	}

	return h
}
