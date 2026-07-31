# Multiple Grammars and Entry Points

pegtool can place multiple grammars on one Go parser runtime. The explicitly designed and tested workflow is `grammar-only` within one Go package. Relying on this capability is not currently recommended for other targets.

## Generate the main grammar

```powershell
pegtool -t go -o base.peg.go base.peg
```

The main file contains the parser runtime, the default grammar variable `g`, and action callbacks.

## Generate an additional grammar

```powershell
pegtool -t go `
  -grammar-only `
  -grammar-name=gOperators `
  -run-func-prefix=operators_ `
  -o operators.peg.go `
  operators.peg
```

The three options solve different problems:

| Option | Purpose |
| --- | --- |
| `-grammar-only` | Omit a duplicate parser runtime |
| `-grammar-name` | Give the additional grammar a unique variable name instead of `g` |
| `-run-func-prefix` | Give action callbacks a unique prefix to avoid method-name collisions |

Names and prefixes must be valid identifier fragments in the target language. All generated Go files must belong to the same package, and shared declarations such as `ParserCustomData` and public wrappers must be defined only once.

Use a template variable to emit shared declarations only in the main file:

```peg
{
package parser

{{ if not .GrammarOnly }}
type ParserCustomData struct{}
{{ end }}
}
```

## Select an additional grammar and entry point

```go
func ParseOperator(filename string, input []byte) (any, error) {
    p := newParser(filename, input)
    p.entrypoint = "Operator"
    return p.parse(gOperators)
}
```

The parser's default entry point comes from the first rule of the main grammar. Set `p.entrypoint` explicitly when the additional grammar's first rule has a different name.

## Current status of `-alternate-entrypoints`

The current CLI option `-alternate-entrypoints=A,B` only verifies that these rules exist. It does not generate a public selector or change the generated file. Select runtime entry points through each target's API:

- Go: set `p.entrypoint` in an initializer wrapper.
- Haxe: set `parser.entrypoint` before parsing.
- TypeScript, C#, C99, and Rust: use the entrypoint field in the relevant `ParserOptions`.

## Index-mode consistency

If you enable `-optimize-ref-expr-by-index`, the main runtime and all grammar-only fragments must use a consistent mode. Indexed references depend on the rule order in each grammar, and rule arrays must not be reordered at runtime.

::: warning Target support boundary
TypeScript grammar-only output does not automatically import private runtime helpers from the main module. Haxe currently has only output-format coverage for this mode. For anything other than the same-package Go workflow, perform an independent compilation experiment before using it in production.
:::
