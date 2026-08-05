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

// Option configures an Engine created with NewEngineWithOptions.
type Option func(*Engine)

// WithLenientMode enables lenient (neutral) evaluation semantics.
//
// In lenient mode, a comparison predicate involving a missing attribute
// returns true (neutral), imposing no constraint. This is the identity
// element for AND-chains (the dominant pattern in betting lifecycle rules),
// so a missing optional field drops out of the conjunction instead of
// failing it:
//
//	x eq 10  (x absent) -> true	x ne 10  (x absent) -> true
//	x lt 10  (x absent) -> true	x in [...]  (x absent) -> true
//	x co "a" (x absent) -> true	x dl 30  (x absent) -> true
//
// Composition is intentional: neutral only inside an AND-chain. Under
// negation or quantifiers neutrality becomes decisive, e.g.
// `not (x eq 10)` -> false and `selections none (r gt 0)` sees true per
// element where r is absent and returns false. For such rules prefer strict
// mode (the default).
//
// Presence (pr), logical (and/or/not) and quantifier operators are
// unaffected. The default (strict) mode is unchanged. Opt-in only.
func WithLenientMode() Option {
	return func(e *Engine) {
		e.evaluator.lenient = true
	}
}

// NewEngineWithOptions creates a new rule engine configured with the given
// options. NewEngine() is equivalent to NewEngineWithOptions() with no
// options (strict mode).
func NewEngineWithOptions(opts ...Option) *Engine {
	e := NewEngine()

	for _, opt := range opts {
		if opt != nil {
			opt(e)
		}
	}

	return e
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
