---
layout: home

hero:
  name: Pegtool
  text: A PEG parser generator written in Go
  tagline: Generate reviewable, deployable parser source for Go, Haxe, TypeScript, C#, C99, and Rust.
  actions:
    - theme: brand
      text: Quick Start
      link: /en/guide/getting-started
    - theme: alt
      text: Migrate from Pigeon
      link: /en/guide/fork-semantics

features:
  - title: Extensive performance work
    details: Continuously optimized for real grammars, including implicit values, state, rule lookup, memoization, and Unicode hot paths.
  - title: Understandable generated code
    details: You rarely need to read generated code, but its clear structure makes custom runtimes, diagnostics, and new targets easier to build.
  - title: Standard library by default
    details: Generated parsers depend only on the target language's standard library by default.
  - title: Multiple target languages
    details: Choose Go, Haxe, TypeScript, C#, C99, or Rust to match your deployment environment.
---

## What is Pegtool?

Pegtool is a PEG (Parsing Expression Grammar) parser generator written in Go. A single grammar can generate source code for Go, Haxe, TypeScript, C#, C99, and Rust. The project forked from Pigeon in 2024, from a baseline that included Pigeon's 2023 `v1.2.x` release line. Its grammar and runtime have evolved independently since then. The main reason for the fork was to improve upstream Pigeon's performance on real grammars. Many optimizations followed, and some required changes to grammar and value semantics. See [Pigeon #151](https://github.com/mna/pigeon/issues/151) for the early analysis.

```powershell
pegtool -t go   -o parser.go grammar.peg
pegtool -t hx   -o Parser.hx grammar.peg
pegtool -t ts   -o parser.ts grammar.peg
pegtool -t cs   -o Parser.cs grammar.peg
pegtool -t c    -o parser.c grammar.peg
pegtool -t rust -o parser.rs grammar.peg
```

## Where to start

- New to Pegtool: read the [Quick Start](/en/guide/getting-started).
- Migrating from upstream Pigeon: use the [automated migration guide and prompt](/en/guide/fork-semantics).
- Preparing for production: choose [performance options](/en/guide/performance-options) according to your workload.
- Writing a cross-language grammar: see the [target overview](/en/targets/).

::: tip Keep the defaults first
Do not enable memoization until benchmarks show that the same named rule is repeatedly evaluated at the same input position. Do not enable rule index optimization unless generated files are excluded from code review.
:::
