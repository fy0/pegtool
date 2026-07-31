# Overview

pegtool is a PEG (Parsing Expression Grammar) parser generator written in Go. It reads a `.peg` grammar and generates source code in the target language during the build. The deployed parser does not need pegtool itself.

The project forked from Pigeon in 2024, from a baseline that included Pigeon's 2023 `v1.2.x` release line. The two projects retain broadly similar PEG syntax, but pegtool's grammar, value semantics, generated APIs, and runtimes have evolved independently. Performance was the primary reason for the fork, followed by extensive work on the execution efficiency of generated parsers. Eliminating costs such as implicit value arrays, state copies, and dynamic lookups required some changes to grammar and value semantics. The original analysis and measurements are recorded in [Pigeon #151](https://github.com/mna/pigeon/issues/151).

## Why this fork exists

I have an online TRPG project with a dice expression interpreter written in Go. The interpreter originally used [`pointlander/peg`](https://github.com/pointlander/peg). It was fast, but generated code contained many `goto` statements, and even a small grammar change could produce a large diff. Go packages also commonly commit generated `.go` files, making the history noisy and the parser difficult to review or customize.

The project later moved to Pigeon to make generated code easier to review and customize. Pigeon's generation structure was simpler, clearer, and easier to extend. After the migration, however, its parser was noticeably slower on the project's real grammar. In an early 2024 measurement, upstream Pigeon took about `4.5s` on the same grammar while `pointlander/peg` took about `0.3s`. By removing expensive state paths, reducing intermediate sequence values, and tightening rule lookup, the early fork reduced the Pigeon path to roughly `0.5s`. This work began in 2024, and further optimizations and changes have been made since.

The fork grew out of that second migration: preserve Pigeon's reviewable, portable generation structure while optimizing the hot paths that mattered. The tradeoff is that some performance changes cannot retain upstream Pigeon's implicit value semantics exactly. This is why migrations must regenerate the parser and run tests.

## Similar tools

The Go ecosystem has several PEG parser generators. They generate code differently and optimize for different needs:

| Tool | Generation model and characteristics | Best suited for |
| --- | --- | --- |
| [`pointlander/peg`](https://github.com/pointlander/peg) | Generates highly expanded Go parsers with many `goto` statements and generally good execution speed; grammar changes can produce large generated diffs | Throughput-first projects where generated files do not need sustained review and the grammar is relatively stable |
| [Upstream Pigeon](https://github.com/mna/pigeon) | Generates straightforward code with an understandable, extensible grammar and runtime while retaining upstream syntax and APIs | Projects that prioritize upstream compatibility and extensibility and can accept its performance on their workload |
| pegtool | Retains a reviewable generation structure, optimizes parser hot paths, and supports multiple targets; some value semantics differ from upstream Pigeon | Projects that need stable generated diffs, custom runtimes or diagnostics, or multiple output languages |

No tool is best for every grammar. Compare performance on real input, error diagnostics, generated diffs, target languages, and whether generated files must be committed and reviewed.

## From Haxe to multiple targets

In 2024, the author began learning Haxe and initially planned to use the Heaps.io engine. Haxe's JavaScript transpilation also looked particularly capable. As an aside, GopherJS appeared to have fewer maintenance resources and less funding at the time, evolved more slowly, and lacked dead code elimination (DCE), which made its output too large for this use case. Haxe itself was interesting, but practical issues with Heaps.io ended that application path.

That experiment still produced one of the project's most important structural changes. Inspired by Haxe's multi-target approach, the Pigeon builder was extended to select an output language, beginning with a Haxe target. The project then used it to implement a narrative scripting engine, demonstrating that one PEG core could generate parsers for different host languages.

Development slowed for nearly two years while the author focused on other work. In 2026, more target languages and their runtimes were added, producing today's pegtool with support for Go, Haxe, TypeScript, C#, C99, and Rust.

## Supported target languages

| Target | Option | Generated output |
| --- | --- | --- |
| Go | `-t go` | Go source file |
| Haxe | `-t hx` | Haxe source that can be compiled to JavaScript, HashLink, and other backends |
| TypeScript | `-t ts` | TypeScript module |
| C# | `-t cs` | `partial` parser class |
| C99 | `-t c` | C99 source file |
| Rust | `-t rust` | Rust module |

All targets share the PEG matching structure, while initializers, actions, predicates, and public APIs use the relevant host language. See the [target overview](/en/targets/) for details.

## Basic workflow

```text
grammar.peg
    -> validated and generated by pegtool
    -> parser.go / Parser.hx / parser.ts / ...
    -> target-language compiler
    -> application
```

The basic generation command is:

```powershell
pegtool -t go -o parser.go grammar.peg
```

Generated output depends only on the target language's standard library by default. The exceptions are dependencies explicitly referenced by grammar actions and Haxe's optional external `hxUnicode` mode.

## Important current semantics

- Capture text explicitly with `label:<expr>`.
- Ordinary terminals and sequences do not create implicit semantic value arrays.
- Repetition collects only non-null values explicitly returned by its child expression.
- Unicode classes use Unicode data from generation time.
- Left recursion is rejected during generation and must be rewritten as a head followed by repeated tails.
- Memoization and rule index optimization should be chosen for real workloads, not enabled by default.

## Continue reading

- New users can generate a Go parser with the [Quick Start](/en/guide/getting-started).
- For grammar rules, read [PEG Syntax](/en/guide/grammar) and [Actions and Values](/en/guide/actions-and-values).
- To migrate from upstream Pigeon, use the [automated migration guide and prompt](/en/guide/fork-semantics).
- Before choosing a production configuration, read [Choosing Performance Options](/en/guide/performance-options).
- For every command-line option, see the [CLI Reference](/en/reference/cli).
