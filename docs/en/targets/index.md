# Target Overview

pegtool uses one shared set of PEG structural semantics, while the initializer, actions, predicates, and public invocation belong to the target language. Choose the deployment target before writing host-language code.

## Support matrix

| Target | CLI values | Primary generated API | Action value | Runtime memoization |
| --- | --- | --- | --- | --- |
| Go | `go` | Initializer wraps package-private `parse` | `any` | Supported |
| Haxe | `hx`, `haxe` | `new Parser(...).parse(...)` | `Any` | Supported, but mutually exclusive with `-optimize-parser` |
| TypeScript | `ts`, `typescript` | `parse`, `Parser` | `any` | Not supported |
| C# | `cs`, `csharp`, `c#` | `PTParser.Parse` | `object?` | Not supported |
| C99 | `c`, `c99` | `pegtool_parse*` | `PegtoolValue` | Not supported |
| Rust | `rust` | `PTParser::parse*` | `Value` | Not supported |

Each complete generated parser depends only on the target language's standard library by default. Haxe requires `hxUnicode` only when `-haxe-use-hxunicode` is explicitly selected and the grammar contains Unicode classes.

## Distinguishing targets in one grammar

Initializers and actions are rendered through Go `text/template` first. The stable template variables are `.Target` and `.GrammarOnly`:

```text
{
{{ if eq .Target "go" }}
package parser
type ParserCustomData struct{}
{{ else if eq .Target "ts" }}
export interface AstNode { kind: string }
{{ else if eq .Target "cs" }}
private sealed record AstNode(string Kind);
{{ else if eq .Target "c" }}
/* C helper declarations */
{{ else if eq .Target "rust" }}
// Rust helper declarations
{{ end }}
}
```

Use the canonical values `go`, `hx`, `ts`, `cs`, `c`, and `rust` in templates, not command-line aliases.

## Choosing a target

- For an existing Go service or highly customized wrappers, choose [Go](/en/targets/go).
- To compile to both JavaScript and HashLink, choose [Haxe](/en/targets/haxe).
- For small generated output that compiles directly in strict Web or Node projects, choose [TypeScript](/en/targets/typescript).
- For a .NET project, choose [C#](/en/targets/csharp).
- For a portable C ABI without a VM, choose [C99](/en/targets/c).
- For an explicit `Value` enum in a Rust project, choose [Rust](/en/targets/rust).

Target rankings depend on grammar shape, actions, and input and cannot be inferred from one microbenchmark. Haxe/JavaScript can outperform direct TypeScript on long linear input, while TypeScript generally has lower fixed cost and smaller generated output. Backtracking-heavy Haxe grammars can also use named-rule memoization.
