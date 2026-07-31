# CLI Reference

```text
pegtool [options] [GRAMMAR_FILE]
```

Without a file, pegtool reads the grammar from stdin. Without `-o`, it writes generated source to stdout.

## Output and targets

| Option | Description |
| --- | --- |
| `-o FILE` | Write output to the specified file |
| `-t TARGET`, `-target TARGET` | Select `go`, `hx`, `ts`, `cs`, `c`, or `rust`; defaults to `go` |
| `-x` | Parse and validate the grammar without generating a parser |
| `-receiver-name NAME` | Change the default name `c` of the current-match object in actions |
| `-nolint` | Add lint-suppression comments to generated Go output |

Target aliases:

| Canonical value | Aliases |
| --- | --- |
| `hx` | `haxe` |
| `ts` | `typescript` |
| `cs` | `csharp`, `c#` |
| `c` | `c99` |

## Generation performance and runtime form

| Option | Description |
| --- | --- |
| `-optimize-parser` | Generate a fixed production runtime; effects differ by target |
| `-optimize-ref-expr-by-index` | Replace name lookups with rule-array indexes, potentially producing large diffs |
| `-cache` | Cache only the generator's parsing of `.peg`; does not affect the generated parser |
| `-debug` | Print debugging output while the generator parses `.peg` |

Do not confuse `-cache` with memoization in the generated parser. See [Choosing Performance Options](/en/guide/performance-options) for the complete selection criteria.

## Entry points and multiple grammars

| Option | Description |
| --- | --- |
| `-alternate-entrypoints A,B` | Verify that the listed rules exist; currently does not change the generated API or output |
| `-grammar-only` | Generate only the grammar object and actions without duplicating the runtime |
| `-grammar-name NAME` | Set the grammar variable name; defaults to `g` |
| `-run-func-prefix PREFIX` | Prefix generated action callbacks |

The explicitly supported and tested multi-grammar workflow currently uses one Go package. See [Multiple Grammars and Entry Points](/en/guide/multiple-grammars) for a complete example.

## Haxe-specific option

| Option | Description |
| --- | --- |
| `-haxe-use-hxunicode` | Use `hxUnicode` for Unicode classes, reducing embedded source at the cost of a dependency |

Using this option with a non-Haxe target is an error. A grammar without Unicode classes does not gain an external dependency from the option.

## Help

```powershell
pegtool -h
pegtool -help
```

## Common commands

```powershell
# Validate a grammar.
pegtool -x grammar.peg

# Production Go parser with stable diffs.
pegtool -t go -optimize-parser -o parser.go grammar.peg

# Haxe parser with optional memoization.
pegtool -t hx -o Parser.hx grammar.peg

# Strict TypeScript parser.
pegtool -t ts -o parser.ts grammar.peg
```
