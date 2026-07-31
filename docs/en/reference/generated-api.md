# Generated API

Targets share concepts such as entry points, maximum expression counts, custom data, and positions, while their APIs follow target-language conventions.

## API comparison

| Target | Invocation | Options | Errors |
| --- | --- | --- | --- |
| Go | Initializer wraps `parse` | Expose `option` through a project wrapper | `(any, error)` |
| Haxe | `new Parser(filename, input, opts).parse(null)` | `Array<Option>` or public fields | Throws an aggregate exception |
| TypeScript | `parse(input, options)` | `ParserOptions` | Throws `PTParseError` / `Error` |
| C# | `PTParser.Parse(input, options)` | `PTParser.ParserOptions` | Throws `PTParseException` |
| C99 | `pegtool_parse*` | `PegtoolParserOptions` | `PegtoolResult.ok/error` |
| Rust | `PTParser::parse*` | `ParserOptions` | `Result<Value, ParseError>` |

## Common options

- `filename`: used for diagnostics only; it does not read a file.
- `entrypoint`: selects a rule name; defaults to the grammar's first rule.
- `maxExpressions`, or its target-specific spelling: limits expression executions to prevent pathological input from consuming resources indefinitely.
- `customData`, or its target-specific spelling: exposes caller state to actions and predicates.

## TypeScript

```ts
const value = parse(input, {
  filename: "input.txt",
  entrypoint: "Start",
  maxExpressions: 1_000_000,
  customData: context,
});
```

`PTParseError` provides `position` and `expected`. `Position.offset` is a UTF-8 byte offset; `line` and `col` start at 1.

## C#

```csharp
var value = PTParser.Parse(input, new PTParser.ParserOptions
{
    Filename = "input.txt",
    Entrypoint = "Start",
    MaxExpressions = 1_000_000,
    CustomData = context,
});
```

`PTParseException` provides `Position` and `Expected`.

## C99

```c
PegtoolParserOptions options = {0};
options.filename = "input.txt";
options.entrypoint = "Start";
options.max_expressions = 1000000;

PegtoolResult result = pegtool_parse_with_options(input, &options);
/* Consume result.value or result.error. */
pegtool_result_free(&result);
```

`PegtoolResult` owns its value or error. The caller is always responsible for freeing it.

## Rust

```rust
let options = ParserOptions {
    filename: "input.txt".to_owned(),
    entrypoint: Some("Start".to_owned()),
    max_expressions: 1_000_000,
    custom_data: None,
};

let value = PTParser::parse_with_options(input, options)?;
```

Results use the `Value` enum rather than relying on `Any` to infer primitive types.

## Go and Haxe memoization

A Go initializer can expose the memoization option through a wrapper:

```go
func Memoize(enabled bool) Option { return memoized(enabled) }
```

Haxe output generated without `-optimize-parser` can set it directly:

```haxe
parser.memoize = true;
```

TypeScript, C#, C99, and Rust currently provide no runtime memoization.
