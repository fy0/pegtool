# 生成 API

各 target 共享入口规则、最大表达式数、自定义数据和位置信息等概念，但 API 采用目标语言习惯。

## API 对照

| Target | 调用入口 | Options | 错误 |
| --- | --- | --- | --- |
| Go | initializer 自行包装 `parse` | 自行暴露 `option` | `(any, error)` |
| Haxe | `new Parser(filename, input, opts).parse(null)` | `Array<Option>` 或公开字段 | 抛出聚合 exception |
| TypeScript | `parse(input, options)` | `ParserOptions` | 抛出 `PigeonParseError` / `Error` |
| C# | `PigeonParser.Parse(input, options)` | `PigeonParser.ParserOptions` | 抛出 `PigeonParseException` |
| C99 | `pigeon_parse*` | `PigeonParserOptions` | `PigeonResult.ok/error` |
| Rust | `PigeonParser::parse*` | `ParserOptions` | `Result<Value, ParseError>` |

## 通用 options

- `filename`：只用于诊断，不读取文件。
- `entrypoint`：选择规则名；默认是文法第一条规则。
- `maxExpressions` / 对应命名：限制表达式执行数，防止病态输入无限消耗。
- `customData` / 对应命名：把调用方状态暴露给 action 和 predicate。

## TypeScript

```ts
const value = parse(input, {
  filename: "input.txt",
  entrypoint: "Start",
  maxExpressions: 1_000_000,
  customData: context,
});
```

`PigeonParseError` 提供 `position` 和 `expected`。`Position.offset` 是 UTF-8 字节 offset，`line`、`col` 从 1 开始。

## C#

```csharp
var value = PigeonParser.Parse(input, new PigeonParser.ParserOptions
{
    Filename = "input.txt",
    Entrypoint = "Start",
    MaxExpressions = 1_000_000,
    CustomData = context,
});
```

`PigeonParseException` 提供 `Position` 与 `Expected`。

## C99

```c
PigeonParserOptions options = {0};
options.filename = "input.txt";
options.entrypoint = "Start";
options.max_expressions = 1000000;

PigeonResult result = pigeon_parse_with_options(input, &options);
/* consume result.value or result.error */
pigeon_result_free(&result);
```

`PigeonResult` 拥有其 value/error；调用方始终负责释放。

## Rust

```rust
let options = ParserOptions {
    filename: "input.txt".to_owned(),
    entrypoint: Some("Start".to_owned()),
    max_expressions: 1_000_000,
    custom_data: None,
};

let value = PigeonParser::parse_with_options(input, options)?;
```

结果使用 `Value` 枚举，不通过 `Any` 猜测基础类型。

## Go 与 Haxe memo

Go 的公开 memo option 由 initializer 包装：

```go
func Memoize(enabled bool) Option { return memoized(enabled) }
```

Haxe 未使用 `-optimize-parser` 时可直接设置：

```haxe
parser.memoize = true;
```

TypeScript、C#、C99 和 Rust 当前不提供运行时 memo。
