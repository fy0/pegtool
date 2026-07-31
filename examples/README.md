# Examples

Examples are grouped first by generated language and then by grammar:

```text
examples/
  go/<example>/
  rust/<example>/
  haxe/<example>/
  typescript/<example>/
  csharp/<example>/
  c/<example>/
```

The examples use the Go versions as the behavioral baseline. The same four
grammars and scenarios are implemented for every full target:

| Example | Go | Rust | Haxe | TypeScript | C# |
| --- | --- | --- | --- | --- | --- |
| `calculator` | yes | yes | yes | yes | yes |
| `indentation` | yes | yes | yes | yes | yes |
| `json` | yes | yes | yes | yes | yes |
| `predicate_code` | yes | yes | yes | yes | yes |

Generated parser sources are committed so each example can be inspected and
run without regenerating it first. The `words` example is an additional shared
Unicode smoke test for Haxe, TypeScript, C#, and C; it does not replace the
Go-equivalent examples above.

## Go

Regenerate all parsers with `make all`, or generate one parser directly:

```sh
go run . -o examples/go/calculator/calculator.go examples/go/calculator/calculator.peg
```

Run the examples' tests with:

```sh
go test ./examples/go/...
```

## Rust

Generate a parser by selecting the Rust target:

```sh
go run . -t rust -o examples/rust/calculator/src/parser.rs examples/rust/calculator/calculator.peg
```

The Rust examples form one Cargo workspace. Test all of them with:

```sh
cargo test --manifest-path examples/rust/Cargo.toml --workspace
```

## Haxe

Generate, compile to JavaScript, and run the calculator with:

```sh
go run . -t hx -o examples/haxe/calculator/Parser.hx examples/haxe/calculator/calculator.peg
cd examples/haxe/calculator
haxe build.hxml
node main.js "3 + (2 - 5 * 12)"
```

The same `build.hxml` workflow applies to `indentation`, `json`,
`predicate_code`, and `words`.

## TypeScript

Generate, compile, and run the calculator with:

```sh
go run . -t ts -o examples/typescript/calculator/parser.ts examples/typescript/calculator/calculator.peg
cd examples/typescript/calculator
tsc -p tsconfig.json
node dist/main.js "3 + (2 - 5 * 12)"
```

The other TypeScript directories use the same `tsconfig.json` workflow.

## C#

Generate and run the .NET calculator with:

```sh
go run . -t cs -o examples/csharp/calculator/Parser.cs examples/csharp/calculator/calculator.peg
dotnet run --project examples/csharp/calculator/Calculator.csproj -- "3 + (2 - 5 * 12)"
```

Each C# example has its own project file and follows the same workflow.

## C99

The C target currently has the additional `words` smoke test, not the full
Go-equivalent matrix. Generate, compile, and run it with:

```sh
go run . -t c -o examples/c/words/parser.c examples/c/words/words.peg
cd examples/c/words
cc -std=c99 -Wall -Wextra -Werror -pedantic-errors main.c -o words
./words
```
