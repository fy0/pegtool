# Rust

Rust target 生成独立模块，使用 `Value` 枚举表达 action、标签和最终结果。

## 最小可运行示例

先在空目录中创建 binary crate：

```powershell
cargo init --bin --name identifier --edition 2021 .
```

新建 `identifier.peg`，并用下面内容替换 `src/main.rs`。`src/parser.rs` 由 pegtool 生成。

### `identifier.peg`

```peg
Start = value:<Identifier> !. { return value; }
Identifier = [\pL_] [\pL\p{Nd}_]*
```

### `src/main.rs`

```rust
mod parser;

use parser::{PTParser, Value};

fn main() {
    let value = PTParser::parse("name_世界").expect("valid identifier");
    match value {
        Value::String(text) => println!("{text}"),
        other => panic!("unexpected parser result: {other:?}"),
    }
}
```

生成并运行：

```powershell
pegtool -t rust -o src\parser.rs identifier.peg
cargo run --quiet
```

输出：

```text
name_世界
```

## `Value`

生成模块提供可克隆的 `Value`：

- `Null`
- `String`
- `Int`
- `Float`
- `Bool`
- `Array`
- `Object`
- `Custom`

常用读取方法包括 `as_str`、`as_array` 和 `as_i64`。`Value::custom` 可以保存实现 `Any` 的项目类型。

## Action

```peg
Integer = '-'? [0-9]+ {
    match c.text.parse::<i64>() {
        Ok(value) => return Value::Int(value),
        Err(error) => {
            p.add_err(error);
            return Value::Null;
        }
    }
}
```

initializer 在模块作用域输出，可以定义 helper 函数和类型。标签本身也是 `Value`，需要按 variant 解构。

## Options

```rust
let options = ParserOptions {
    filename: "input.txt".to_owned(),
    entrypoint: Some("Start".to_owned()),
    max_expressions: 1_000_000,
    custom_data: None,
};

let value = PTParser::parse_with_options(source, options)?;
```

失败返回 `ParseError`。Rust target 当前没有 memo 开关；`-optimize-parser` 没有 Rust 专用效果。规则索引参数仍应根据真实基准和生成 diff 要求决定。
