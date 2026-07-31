<script setup>
import MigrationPrompt from "../.vitepress/theme/components/MigrationPrompt.vue";
</script>

# 从原版 Pigeon 迁移

Pegtool 使用 Go 编写，于 2024 年从 Pigeon 分叉，分叉基线包含 Pigeon 2023 年的 `v1.2.x` 发布线。两者的 PEG 写法总体相似，但 Pegtool 不是 Pigeon 的语法兼容层；其文法、值语义、生成 API 和运行时已经独立演进。性能是这个 fork 的主要起因，此后进行了大量优化；部分优化需要改变原版的隐式值传递语义，背景见 [Pigeon #151](https://github.com/mna/pigeon/issues/151)。

大多数迁移可以交给具备仓库读写和测试能力的 code agent 完成。下面的 Prompt 会自动带上本页的完整链接，可直接交给 agent：

<MigrationPrompt />

## 迁移差异总览

| 项目 | 原版 Pigeon `v1.2.x` | Pegtool 当前行为 |
| --- | --- | --- |
| 终结符值 | 字面量、字符类和 `.` 返回 `[]byte` | 返回目标语言的空值；文本用 `label:<expr>` 捕获 |
| sequence 值 | 返回 `[]any`，每个子表达式占一项 | 返回空值，不创建隐式数组 |
| repetition 值 | 收集每次匹配的结果，包括空值 | 只收集子表达式显式返回的非空值 |
| Go action | 返回 `(value, error)` | 返回单个值；错误用 `p.addErr(err)` 记录 |
| Go predicate | 返回 `(bool, error)` | 返回 `bool`；错误需在代码中显式记录 |
| Go 生成 API | 导出 `Parse`、`ParseFile`、`ParseReader` 和固定 options | 基础符号为包内 API，由 initializer 或普通源码包装 |
| Parser 状态 | `state`、`globalStore` 和 `#{...}` | 使用 `ParserCustomData` / `c.data`，不提供自动回滚的等价状态层 |
| 左递归 | 可用实验性参数开启 | builder 会拒绝左递归 |
| 生成目标 | 主要生成 Go | 支持 Go、Haxe、TypeScript、C#、C99 和 Rust |

这张表是迁移检查表，不是替换规则。尤其是状态、错误处理和公开 API，需要根据项目原有语义逐项改写。

## 文本与值不再隐式生成

原版 Pigeon 的终结符和 sequence 会创建隐式值。例如：

```peg
// 原版 Pigeon
Word <- [a-z]+ {
    return string(c.text), nil
}

Pair <- pair:("a" ":" "b") {
    // pair 的底层类型是 []any
    return pair, nil
}
```

Pegtool 应显式说明需要的是原文还是结构化值：

```peg
// Pegtool
Word = value:<[a-z]+> {
    return value
}

Pair = pair:<("a" ":" "b")> {
    // pair 是匹配到的 string
    return pair
}
```

仅把 `pair:(...)` 原样迁移会得到空值。若需要结构化结果，应让内部 action 显式返回对象，再通过标签接收；不要依赖旧 sequence 的嵌套 `[]any`。

repetition 同样需要检查。原版会为每次匹配保留一项；Pegtool 只收集非空的显式返回值：

```peg
// Pegtool：每个 Name 显式返回 string，因此 Names 能收集这些值
Name  = value:<[a-z]+> { return value }
Names = first:Name rest:(_ "," _ next:Name { return next })* {
    values := []any{first}
    if rest != nil {
        values = append(values, rest.([]any)...)
    }
    return values
}
```

## Go action 与 predicate

原版 Pigeon action 和 predicate 都通过第二个返回值上报错误：

```peg
// 原版 Pigeon
Integer <- [0-9]+ {
    return strconv.Atoi(string(c.text))
}

Allowed <- &{
    return c.globalStore["enabled"].(bool), nil
} .
```

Pegtool 的 Go 代码块返回单个值。解析错误通过当前 parser 记录：

```peg
// Pegtool
Integer = [0-9]+ {
    value, err := strconv.Atoi(string(c.text))
    if err != nil {
        p.addErr(err)
        return nil
    }
    return value
}

Allowed = &{
    return c.data.Enabled
} .
```

迁移 predicate 时，不要简单删除 `error` 返回值。如果旧代码确实可能失败，应先调用 `p.addErr(err)`，再返回合适的布尔结果。

## Go API 需要显式包装

原版生成物固定导出 `Parse`、`ParseFile`、`ParseReader`、`Entrypoint`、`Memoize` 等 API。Pegtool 的 Go runtime 只提供包内的 `parse`、`newParser`、`option` 和 `memoized`；项目可以在 initializer 中保留原有的公开表面：

```peg
{
package parser

type ParserCustomData struct{}
type Option = option

func Memoize(enabled bool) Option {
    return memoized(enabled)
}

func Parse(filename string, input []byte, opts ...Option) (any, error) {
    return parse(filename, input, opts...)
}
}

Start = value:<[a-z]+> !. { return value }
```

`ParseFile` 和 `ParseReader` 若属于项目公开 API，应在普通 Go 源码或 initializer 中重新实现并覆盖测试。即使没有自定义状态，Go 文法也应声明空的 `ParserCustomData`。

## 状态不能机械替换

原版的 `state` 会随 PEG 回溯恢复，`globalStore` 不会，`#{...}` 还具有专门的状态修改语义。Pegtool 的 `c.data` 是调用方提供的 `*ParserCustomData`，修改后不会自动回滚。因此不存在把三个旧概念统一替换成 `c.data` 且保持行为不变的通用改写。

迁移时应先区分每个字段：

1. 只读配置可以放进 `ParserCustomData`。
2. 不需要回滚的累计状态可以显式修改 `c.data`。
3. 需要随 choice 失败而恢复的状态，应由文法保存和恢复，或改为返回不可变结果。
4. memo 命中不会重新执行规则内部的 action 或 predicate；配置若在一次 parse 中变化，应关闭 memo 或重新划分解析生命周期。

Pegtool 另外提供 `label:<expr>` 文本捕获、`&&` / `!!` 消费判断和 `*{...}` lookahead 内代码表达式。这些是当前文法能力，不应反向写进仍由原版 Pigeon 生成的文法。

## 左递归与入口规则

原版可通过实验性参数支持直接或间接左递归；Pegtool 会在生成阶段拒绝它。需要改成首项加尾项重复：

```peg
// 原版可选左递归写法
Expr <- Expr "+" Term / Term

// Pegtool
Expr = first:Term rest:(_ "+" _ next:Term { return next })* {
    return foldAdd(first, rest)
}
```

入口规则也应通过各 target 的生成 API 或项目 wrapper 选择，不要假定原版的 `Entrypoint(...)` option 在所有 target 都存在。

## 验证迁移结果

迁移完成至少检查以下内容：

1. 使用 `pegtool -x grammar.peg` 验证文法。
2. 用项目实际 target 重新生成 parser，不手改生成文件掩盖问题。
3. 编译生成物，并运行原项目完整测试。
4. 补测 sequence、repetition、action 错误、predicate、状态回溯、Unicode 和非法输入。
5. 对比项目原有公开 API；必要时保留兼容 wrapper。
6. 只有真实负载基准证明有收益时才启用 memo 或规则索引优化。

具体语法见 [PEG 语法](/guide/grammar)，值传播规则见[动作与值](/guide/actions-and-values)，各 target 的调用方式见[生成 API](/reference/generated-api)。
