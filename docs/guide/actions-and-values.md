# 动作与值

Pegtool 让文本、语义值和宿主语言 action 分工明确。文法不再为每个 sequence 自动创建值数组，因此需要的数据必须显式捕获或返回。

## 值传播规则

| 表达式 | 返回值 |
| --- | --- |
| 字面量、字符类、`.` | `nil` / `null` / 目标语言空值 |
| `label:<expr>` | 匹配原文，类型为字符串 |
| 普通 sequence | 空值 |
| choice | 命中分支的值 |
| `expr?` | 子表达式的值或空值 |
| `expr*`、`expr+` | 子表达式显式返回的非空值集合；没有值时为空值 |
| `&expr`、`!expr` | 空值 |
| action | action 显式返回的任意值 |

## 捕获文本

```peg
Name       = value:<[\pL_] [\pL\p{Nd}_]*> { return value }
Assignment = source:<(Name _ "=" _ Name)> { return source }
```

`value:Name` 得到的是 `Name` 的返回值；`value:<Name>` 得到的是 `Name` 消费的原始文本。终结符本身不返回字符或字节。

## 从重复中收集值

错误的期待：

```peg
// [0-9]+ 的终结符不返回值，因此 digits 是空值
Number = digits:[0-9]+
```

捕获完整文本：

```peg
Number = digits:<[0-9]+> { return digits }
```

收集多个结构化结果时，让重复的子 action 返回需要的值：

```peg
Name  = value:<[a-z]+> { return value }
Names = first:Name rest:(_ "," _ next:Name { return next })* {
  values := []any{first}
  if rest != nil {
    values = append(values, rest.([]any)...)
  }
  return values
}
```

这个示例使用 Go action。重复内部的标签不会自动提升到外层作用域。把需要的数据从内部 action 返回，能避免隐式嵌套数组和大量无用的空值。

## Go action

Go action 返回单个 `any`。错误通过 parser 记录：

```peg
Integer = [0-9]+ {
    value, err := strconv.Atoi(string(c.text))
    if err != nil {
        p.addErr(err)
        return nil
    }
    return value
}
```

可用对象：

- `c.text`：action 所附表达式匹配到的完整 `[]byte`。
- `c.pos`：匹配起点的位置。
- `c.data`：`*ParserCustomData`。
- `p`：当前 parser，可调用 `p.addErr` 或项目自行包装的内部能力。

不要使用其他 Pigeon 版本中的 `return value, err` 写法。

## 其他 target 的 action 值

| Target | 标签/返回值模型 | 当前匹配对象 |
| --- | --- | --- |
| Haxe | `Any` | `c.text`、`c.data` |
| TypeScript | `any` | `c.text`、`c.data` |
| C# | `object?` | `c.Text`、`c.Data` |
| C99 | `PigeonValue` 所有权模型 | `PigeonCurrent *c` |
| Rust | `Value` 枚举 | `current.text` 与 custom data |

C99 action 返回拥有所有权的 `PigeonValue`；Rust action 使用可克隆的 `Value` 枚举。具体调用方式见各[目标语言](/targets/)页面。

## Lookahead 中的 action

普通 action 在 `&expr` 和 `!expr` 的探测过程中会跳过，避免探测产生业务副作用。`*{...}` 明确表示即使在 lookahead 中也执行：

```peg
Probe = *{ c.data.ProbeCount++ }
```

不要依赖在 syntactic lookahead 内部创建并向外传播标签；把需要的判断写成语义谓词，或在真正消费输入的分支中捕获。

## Action、副作用与 memo

memo 命中不会重新执行缓存规则内部的 action 和 predicate。以下条件同时成立时才适合启用：

- 规则结果主要由输入位置决定。
- 一次 parse 期间，predicate 读取的配置保持不变。
- action 的副作用不要求在每次回溯访问时重复发生。

若 predicate 使用不断变化的 `c.data` 控制匹配，应关闭 memoization。
