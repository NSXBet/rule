# Plano de Implementação: Modo Lenient (SQL-ish null semantics)

## Objetivo

Adicionar um modo "lenient" ao engine onde campos ausentes (ou `nil`) são
tratados como `null` e propagados com semântica SQL-ish, em vez de virarem
`false` indiscriminadamente em qualquer comparação.

## Motivação

No modo strict (atual), quando um campo está ausente do contexto:

| Query               | Resultado atual | Surpresa? |
|---------------------|-----------------|-----------|
| `x eq 10`           | `false`         | não       |
| `x ne 10`           | `false`         | **SIM** (null ≠ 10 deveria ser true) |
| `not (x eq 10)`     | `true`          | coerente  |
| `x lt 10`           | `false`         | coerente  |
| `x in [...]`        | `false`         | coerente  |
| `x not in [...]`    | `false`         | **SIM** (null não está no set → true) |

O modo lenient corrige essas incongruências aplicando semântica null-aware
estilo SQL (2-valued logic simplificada, sem tri-state "unknown"):

| Operador            | null vs valor concreto | null vs null |
|---------------------|------------------------|--------------|
| `eq` / `==`         | `false`                | `true`       |
| `ne` / `!=`         | `true`                 | `false`      |
| `lt`, `gt`, `le`, `ge` | `false`            | `false`      |
| `co`, `sw`, `ew`    | `false`                | `false`      |
| `in`                | `false`                | `false`      |
| `not in`            | `true`                 | `false`      |
| `pr`                | `false` (ausente) / `true` (nil explícito) | — |
| datetime (`dq,dn,be,bq,af,aq,dl,dg`) | `false` | `false` |

Regras para `pr` (presence) NÃO mudam: ausência → `false`, `nil` explícito →
`true` (presença é sobre "a chave está no map", não sobre o valor).

## Design

### API pública (não quebra nada)

```go
// Modo default (strict, inalterado):
engine := rule.NewEngine()

// Modo lenient:
engine := rule.NewEngineWithOptions(rule.WithLenientMode())
```

- `NewEngine()` continua funcionando exatamente igual (modo strict).
- Adiciona `Option` funcional e `WithLenientMode()`.
- `NewEngineWithOptions(opts ...Option) *Engine`.

### Interno

- `Evaluator` ganha um campo `lenient bool` (zero custo de alocação).
- Único ponto de mudança semântica: `evaluateComparisonOperator`.
  Hoje existe o short-circuit:
  ```go
  if !leftResult.IsValid || !rightResult.IsValid {
      result.Bool = false
      return nil
  }
  ```
  Em modo lenient, substitui por:
  ```go
  if !leftResult.IsValid || !rightResult.IsValid {
      if !e.lenient {
          result.Bool = false
          return nil
      }
      result.Bool = e.lenientCompare(node.Operator, &leftResult, &rightResult)
      return nil
  }
  ```
- Nova função `lenientCompare(operator, left, right) bool` centraliza a tabela
  SQL-ish. Não aloca (operadores são constantes, leitura por valor).
- `pr` (presence) não muda — é tratado em `evaluatePresenceOperator`, fora do
  caminho de comparação.
- `not`, `and`, `or` não mudam — usam `toBool`, que já retorna `false` para
  inválido; `not` inverte → coerente com SQL-ish (NOT null = true via
  `not (x eq 10)` → `not false` → `true`).

### Ausência vs nil explícito

- **Ausência** (chave não está no map): `IsValid = false` → em modo lenient
  tratado como `null`.
- **nil explícito** (`x: nil` no context): hoje cai no `default:` de
  `setResultFromAny` → `IsValid = true`, `Type = ValueString`, `Str = ""`.
  Isso significa que `nil` explícito hoje se comporta como string vazia para
  comparações, mas como presente para `pr`.
  - **Decisão**: no escopo deste PR, lenient trata apenas **ausência** como
    null. `nil` explícito mantém comportamento atual (não quebra testes
    existentes como `nil_value_presence` → `true`). Futuro: nullable tracking
    separado, fora do escopo.

## Passos (commits semânticos)

1. **`feat: add lenient mode option to Engine and Evaluator`**
   - `engine.go`: `Option`, `WithLenientMode()`, `NewEngineWithOptions()`.
   - `evaluator.go`: campo `lenient bool` em `Evaluator`, propagado por
     `NewEvaluator` (interno) e setter a partir do Engine.
   - `NewEngine()` inalterado.

2. **`feat: implement SQL-ish null comparison in lenient mode`**
   - `evaluator.go`: `lenientCompare()` + branch em
     `evaluateComparisonOperator`.
   - Sem mudança de comportamento em strict.

3. **`test: add lenient mode null semantics test suite`**
   - Novo arquivo `test/lenient_fixtures.go` com tabela de casos.
   - Integração em `test/rule_engine_test.go` (novo grupo rodando com
     `NewEngineWithOptions(WithLenientMode())`).
   - Casos: eq/ne/lt/gt/le/ge/co/sw/ew/in/not_in com ausência, ausência ambos
     lados, nil vs concreto, and/or/not, nested missing, datetime.

4. **`test: add lenient mode example in example_test.go`**
   - `Example_lenientMode()` com `// Output:` verificável.

5. **`docs: document lenient mode in README`**
   - Nova seção "## 🟢 Lenient Mode (Null-Aware Semantics)".
   - Atualizar tabela de "Exclusive Features".
   - Atualizar Quick Example/TOC conforme necessário.

6. **`docs: update CLAUDE.md with lenient mode spec`**
   - Adicionar bullet em "Type System Compliance" e "List Operations"
     descrevendo o modo lenient.

## Restrições do projeto respeitadas

- **0 allocs/op**: `lenient bool` é campo de struct, lido por valor; nenhum
  novo slice/map/interface.
- **Sub-100ns**: um branch predizível a mais no caminho hot.
- **API-compatible**: `NewEngine()` inalterado; modo é opt-in.
- **TDD**: testes definem a spec; implementação segue a tabela.
- **Errors tipados**: nenhum novo erro necessário (lenient não introduz
  erros — apenas muda bool de resultado).
- **Sem co-author** nos commits.
- **Scripts de teste em pasta própria**: testes ficam em `test/`.