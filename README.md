# pegtool

[![Go Reference](https://pkg.go.dev/badge/github.com/fy0/pegtool.svg)](https://pkg.go.dev/github.com/fy0/pegtool)
[![Test Status](https://github.com/fy0/pegtool/actions/workflows/ci.yml/badge.svg)](https://github.com/fy0/pegtool/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/fy0/pegtool)](https://goreportcard.com/report/github.com/fy0/pegtool)
[![License](https://img.shields.io/badge/license-BSD--3--Clause-blue.svg)](LICENSE)

pegtool is a PEG (Parsing Expression Grammar) parser generator written in Go.
One grammar can generate reviewable, deployable parser source for Go, Haxe,
TypeScript, C#, C99, or Rust. Generated parsers use only the target language's
standard library by default and do not require pegtool at runtime.

The project forked from [Pigeon](https://github.com/mna/pigeon) in 2024 after
performance work on real grammars required changes to grammar semantics and the
generated runtime. The projects still share a recognizable PEG syntax, but
pegtool is not a compatibility layer for Pigeon. Existing Pigeon grammars should
follow the [migration guide](docs/en/guide/fork-semantics.md) and be regenerated.

[中文文档](https://fy0.github.io/pegtool/) |
[English Documentation](https://fy0.github.io/pegtool/en/)

## Highlights

- Deliver high parser throughput on real grammars by optimizing hot paths and
  reducing intermediate values, state work, and rule lookup overhead.
- Generate Go, Haxe, TypeScript, C#, C99, and Rust from the same PEG structure.
- Keep generated code understandable enough to review, customize, and commit.
- Capture matched text explicitly and avoid unnecessary implicit value arrays.
- Use parser-level diagnostics, alternate entry points, custom caller data, and
  expression limits across supported targets.
- Choose optimization and memoization behavior according to the grammar instead
  of paying for every runtime feature unconditionally.
- Parse UTF-8 input and generate only the Unicode tables required by a grammar.

## Install

```sh
go install github.com/fy0/pegtool@latest
```

This installs the lowercase `pegtool` command into `GOBIN` (or `GOPATH/bin` when
`GOBIN` is unset).

## Quick Start

Create `identifier.peg`:

```peg
{
package parser

type ParserCustomData struct{}

func Parse(input string) (any, error) {
    return parse("", []byte(input))
}
}

Identifier = value:<([A-Za-z_][A-Za-z0-9_]*)> !. {
    return value
}
```

Generate a Go parser:

```sh
pegtool -t go -o parser.go identifier.peg
```

For grammars whose initializer and actions are written for the selected host
language (or guarded by target templates), choose another output target with the
same CLI shape:

```sh
pegtool -t hx   -o Parser.hx grammar.peg
pegtool -t ts   -o parser.ts grammar.peg
pegtool -t cs   -o Parser.cs grammar.peg
pegtool -t c    -o parser.c grammar.peg
pegtool -t rust -o parser.rs grammar.peg
```

Input is read from stdin when no grammar file is provided, and generated source
is written to stdout when `-o` is omitted. Use `pegtool -h` for the complete CLI
reference.

## Generated APIs

Initializer and action code are written in the selected target language. The
generated entry points are therefore target-specific:

| Target | Select with | Primary API |
| --- | --- | --- |
| Go | `-t go` | Define a public wrapper around the generated package-private `parse` function |
| Haxe | `-t hx` | `new Parser(filename, input, options).parse(null)` |
| TypeScript | `-t ts` | `parse(input, options)`; failures use `PTParseError` |
| C# | `-t cs` | `PTParser.Parse(input, options)`; failures use `PTParseException` |
| C99 | `-t c` | `pegtool_parse*`; values and options use the full `Pegtool` prefix |
| Rust | `-t rust` | `PTParser::parse*` returning `Result<Value, ParseError>` |

Aliases `haxe`, `typescript`, `csharp`, `c#`, and `c99` are also accepted.
See the [target overview](docs/en/targets/index.md) and
[generated API reference](docs/en/reference/generated-api.md) before writing
target-specific initializer or action code.

## Grammar Semantics

pegtool deliberately differs from upstream Pigeon in a few important areas:

- Literals, character classes, and sequences do not create implicit semantic
  values. Use `label:<expression>` to capture matched text or return a value from
  an action explicitly.
- Repetitions collect only non-null values explicitly returned by their child
  expression.
- Go actions return one value. Record an action error with `p.addErr(err)` before
  returning the appropriate fallback value.
- `ParserCustomData` / `c.data` carries caller state, but mutations are not
  automatically rolled back during PEG backtracking.
- Direct and indirect left recursion are rejected and must be rewritten.

Read [PEG syntax](docs/en/guide/grammar.md) and
[actions and values](docs/en/guide/actions-and-values.md) for the complete rules.

## Performance Options

Start with the defaults and benchmark real inputs before enabling additional
runtime behavior.

- `-optimize-parser` is a useful production starting point for Go and for
  near-linear Haxe grammars. Its effect differs by target.
- Named-rule memoization helps repeated expensive backtracking, but adds lookup
  and allocation overhead when a grammar is already linear.
- `-optimize-ref-expr-by-index` can reduce rule lookup cost, but rule insertion
  or reordering may cause large generated diffs.
- `-cache` affects how pegtool parses the grammar file; it does not enable
  memoization in the generated parser.

See [Choosing Performance Options](docs/en/guide/performance-options.md) for the
decision guide and current target-specific behavior.

## Documentation

The published documentation is available at
[fy0.github.io/pegtool](https://fy0.github.io/pegtool/).

- [English quick start](docs/en/guide/getting-started.md)
- [中文快速开始](docs/guide/getting-started.md)
- [PEG grammar reference](docs/en/guide/grammar.md)
- [CLI reference](docs/en/reference/cli.md)
- [Migrating from Pigeon](docs/en/guide/fork-semantics.md)
- [Examples by target language](examples/README.md)

## Development

Run the Go test suite:

```sh
go test ./...
```

Build the documentation without starting a development server:

```sh
npm ci
npm run docs:build
```

Generated example sources are committed to the repository. Changes to a builder
or runtime should update the corresponding generated examples and compile them
with the target toolchain. See the
[development guide](docs/en/reference/development.md) for the complete workflow.

## Project History

pegtool is a fork of Pigeon. Pigeon was created by
[@mna](https://github.com/mna) and later maintained by
[@breml](https://github.com/breml). The fork history and migration differences
remain documented so users can evaluate compatibility explicitly.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

pegtool is released under the [BSD 3-Clause License](LICENSE).
