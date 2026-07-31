# 快速开始

Pegtool 使用 Go 编写。下面统一使用命令名 `pegtool`。

## 安装

从 Go 模块安装最新版本：

```powershell
go install github.com/fy0/pigeon@latest
pegtool -h
```

从当前仓库构建：

```powershell
go build -o .\bin\pegtool.exe .
```

## 创建第一份文法

::: tip 使用其他目标语言
本节以 Go 为例。要生成其他语言的 parser，请直接查看 [Haxe](/targets/haxe)、[TypeScript](/targets/typescript)、[C#](/targets/csharp)、[C99](/targets/c) 或 [Rust](/targets/rust) 的最小可运行示例。

不准备运行 Go 示例？可以跳到[切换目标语言](#切换目标语言)，或直接查看[下一步](#下一步)。
:::

新建 `identifier.peg`。下面的 Go 文法识别一个 Unicode 标识符，并返回捕获到的文本：

```peg
{
package main

type ParserCustomData struct{}

func ParseIdentifier(filename string, input []byte) (string, error) {
    value, err := parse(filename, input)
    if err != nil {
        return "", err
    }
    return value.(string), nil
}
}

Start = value:<Identifier> !. {
    return value
}

Identifier = [\pL_] [\pL\p{Nd}_]*
```

这里的关键点是：

- initializer 中的 Go 代码声明生成文件所属的 package 和公开调用函数。
- 即使不使用自定义状态，Go target 也需要声明空的 `ParserCustomData`。
- 第一条规则是默认入口。
- `<Identifier>` 显式捕获匹配文本。
- `!.` 要求输入已经结束，避免只接受合法前缀。
- `\pL` 和 `\p{Nd}` 使用生成器所认识的 Unicode 类。

## 生成 Go parser

```powershell
go mod init example.com/identifier
pegtool -t go -o parser.go identifier.peg
```

新建 `main.go` 调用生成的 parser：

```go
package main

import "fmt"

func main() {
    value, err := ParseIdentifier("input.txt", []byte("name_世界"))
    if err != nil {
        panic(err)
    }
    fmt.Println(value)
}
```

编译并运行：

```powershell
go run .
# 输出：name_世界
```

## 切换目标语言

使用 `-t` 或 `-target` 选择输出：

| 目标 | 参数 | 常用输出名 |
| --- | --- | --- |
| Go | `go` | `parser.go` |
| Haxe | `hx`、`haxe` | `Parser.hx` |
| TypeScript | `ts`、`typescript` | `parser.ts` |
| C# | `cs`、`csharp`、`c#` | `Parser.cs` |
| C99 | `c`、`c99` | `parser.c` |
| Rust | `rust` | `parser.rs` |

文法中的动作代码属于目标语言。需要从同一文件生成多种语言时，使用[目标条件模板](/targets/#在一份文法中区分目标)。

## 下一步

1. 阅读[PEG 语法](/guide/grammar)，理解有序选择与回溯。
2. 阅读[动作与值](/guide/actions-and-values)，避免依赖上游 Pigeon 的序列返回值。
3. 在真实输入上基准测试后，再选择[性能参数](/guide/performance-options)。
