# Examples

Examples are grouped first by generated language and then by grammar:

```text
examples/
  go/<example>/
  rust/<example>/
```

The Go and Rust directories contain equivalent calculator, indentation, JSON,
and semantic-predicate examples. Generated parser sources are committed so each
example can be inspected and run without regenerating it first.

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
