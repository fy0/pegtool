# 目标语言概览

pegtool 共享同一套 PEG 结构语义，但 initializer、action、predicate 和公开调用方式属于目标语言。生成前先确定部署目标，再编写对应的宿主代码。

## 支持矩阵

| Target | CLI 值 | 主要生成 API | Action 值 | 运行时 memo |
| --- | --- | --- | --- | --- |
| Go | `go` | 由 initializer 包装包内 `parse` | `any` | 支持 |
| Haxe | `hx`、`haxe` | `new Parser(...).parse(...)` | `Any` | 支持，但与 `-optimize-parser` 互斥 |
| TypeScript | `ts`、`typescript` | `parse`、`Parser` | `any` | 不支持 |
| C# | `cs`、`csharp`、`c#` | `PTParser.Parse` | `object?` | 不支持 |
| C99 | `c`、`c99` | `pegtool_parse*` | `PegtoolValue` | 不支持 |
| Rust | `rust` | `PTParser::parse*` | `Value` | 不支持 |

每个完整生成物默认只使用目标语言标准库。Haxe 只有在显式使用 `-haxe-use-hxunicode` 且文法包含 Unicode class 时需要 `hxUnicode`。

## 在一份文法中区分目标

initializer 和 action 会先经过 Go `text/template` 渲染。可用的稳定模板变量是 `.Target` 和 `.GrammarOnly`：

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

模板中使用规范值：`go`、`hx`、`ts`、`cs`、`c`、`rust`，不要使用命令别名。

## 如何选 target

- 已有 Go 服务或需要高度定制 wrapper：选 [Go](/targets/go)。
- 需要同时编译到 JavaScript 与 HashLink：选 [Haxe](/targets/haxe)。
- Web/Node 项目希望生成物最小且直接严格编译：选 [TypeScript](/targets/typescript)。
- .NET 项目：选 [C#](/targets/csharp)。
- 需要无 VM 的可移植 C ABI：选 [C99](/targets/c)。
- Rust 项目希望使用显式 `Value` 枚举：选 [Rust](/targets/rust)。

Target 排名依赖文法形状、action 和输入，不能只根据一个微基准决定。Haxe/JS 在长线性输入上可能胜过直接 TypeScript，而 TypeScript 通常有更小的固定成本和生成体积；高回溯 Haxe 文法还可以选择命名规则 memo。
