# 错误处理

PEG 默认只能报告“在最远失败位置期望哪些 matcher”。复杂文法经过多层 choice 和回溯后，这份列表往往更接近 parser 的内部路径，而不是用户能理解的错误原因。原版 Pigeon 若要替换这类默认消息，通常需要在文法各处插入 predicate/recovery，或者在 parse 之后解析内部错误，维护成本较高。

pegtool 保留位置和 `expected` 数据，同时为 Go 和 Haxe 扩展了 parser 级的 no-match formatter；应用可以在不修改生成 runtime、也不把错误文案散落到文法各处的前提下生成面向用户的错误消息。

## Go：`noMatchErrorFormatter`

Go target 的 formatter 签名是：

```go
func(position, []byte, []string) error
```

三个参数依次是最远失败位置、完整输入和已经去重排序的期望项；`!.` 会先被转换为 `EOF`。formatter 返回非空错误时，parser 使用它替代默认的 `no match found, expected: ...`；返回 `nil` 时保留默认错误。

下面的 initializer 为生成 parser 固定安装一个中文 formatter：

```peg
{
package parser

import (
    "fmt"
    "strings"
    "unicode/utf8"
)

type ParserCustomData struct{}

func friendlyNoMatch(pos position, input []byte, expected []string) error {
    got := "EOF"
    if pos.offset < len(input) {
        value, _ := utf8.DecodeRune(input[pos.offset:])
        got = fmt.Sprintf("%q", value)
    }

    wanted := strings.Join(expected, "、")
    if wanted == "" {
        wanted = "合法输入"
    }
    return fmt.Errorf("第 %d 行第 %d 列：遇到 %s，需要 %s", pos.line, pos.col, got, wanted)
}

func Parse(filename string, input []byte) (any, error) {
    return parse(filename, input, noMatchErrorFormatter(friendlyNoMatch))
}
}

Start "标识符" = value:<Identifier> !. { return value }
Identifier = [\pL_] [\pL\p{Nd}_]*
```

`noMatchErrorFormatter(...)` 是单次 parser 的 option，不需要修改全局变量。项目可以根据调用方语言或错误展示形式，在 wrapper 中为每次 parse 传入不同 formatter。该接口在 Go 的 `-optimize-parser` 生成物中仍然保留。

## Haxe 与其他 target

Haxe 提供同样的 parser option，输入类型改为 `String`：

```haxe
var parser = new Parser("input.txt", source, [
    noMatchErrorFormatter((position, input, expected) ->
        new Exception('${position.line}:${position.col}: expected ${expected.join(", ")}')
    )
]);
```

Haxe 的 formatter 同样保留在 `-optimize-parser` 生成物中。其他 target 当前通过结构化错误暴露诊断数据：

| Target | 可用数据 |
| --- | --- |
| TypeScript | `PTParseError.position`、`expected` |
| C# | `PTParseException.Position`、`Expected` |
| Rust | `ParseError.position`、`expected` |
| C99 | `PegtoolResult.error`、`position` |

这些 target 可以在应用边界把结构化错误转换为本地化消息；当前没有与 Go/Haxe 相同的 formatter option。

## 语义错误与显示名称

no-match formatter 只处理“文法无法匹配”的情况。数字溢出、禁止使用的关键字等语义错误，应在 action 或 predicate 中显式记录；Go action 可调用 `p.addErr(err)`：

```peg
Integer "整数" = [0-9]+ {
    value, err := strconv.Atoi(string(c.text))
    if err != nil {
        p.addErr(fmt.Errorf("整数超出范围: %w", err))
        return nil
    }
    return value
}
```

规则名后的字符串是诊断显示名称。在 Go 和 Haxe runtime 中，它会让规则内部记录的错误显示 `整数`，而不是实现名称 `Integer`；它不能替代 no-match formatter。

如果希望在错误后消费无效输入并继续构造结果，应使用下面的 failure label 与恢复语法。单纯改变错误文本时，优先使用 formatter 或结构化错误，不要为此引入恢复分支。

## Failure label 与恢复

`%{Label}` 抛出 failure label，`//{...}` 为指定 label 安装恢复表达式：

```peg
Value = Number / %{ExpectedNumber}
        //{ExpectedNumber} InvalidNumber
```

语义如下：

1. 先匹配 `Value` 左侧的正常表达式。
2. `%{ExpectedNumber}` 从内向外查找能够处理该 label 的 recovery。
3. 找到后，从抛出位置运行 recovery 表达式。
4. recovery 成功则继续解析；失败或找不到 handler 则继续传播失败。

自定义数据不会自动回滚。recovery 修改 `c.data` 时，应由文法负责恢复一致状态。Failure label recovery 与宿主语言的 panic/exception recovery 是两套机制。

错误处理通常会同时涉及 action 和公开 parser API。相关值语义见[动作与值](/guide/actions-and-values)，各 target 的错误类型见[生成 API](/reference/generated-api)。
