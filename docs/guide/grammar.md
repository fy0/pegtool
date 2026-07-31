# PEG 语法

PEG 用识别规则描述输入。与正则表达式不同，规则可以递归引用；与上下文无关文法不同，选择是有顺序的，第一条成功分支立即胜出。

## 文法结构

文法由可选 initializer 和至少一条规则组成：

```peg
{
// 这里是目标语言代码
}

Start "start rule" <- _ value:Expr EOF { return value }
Expr               <- Number / "(" _ Expr _ ")"
Number             <- [0-9]+
_                  <- [ \t\r\n]*
EOF                <- !.
```

第一条规则是默认入口。规则名后可加字符串作为诊断中的显示名称。定义符可以写成 `=`、`<-`、`←` 或 `⟵`。

普通空白会被忽略。支持 `//` 行注释和 `/* ... */` 块注释；`//{` 是 failure recovery 语法的开头，不应拿来写普通注释。

## 表达式

| 写法 | 含义 |
| --- | --- |
| `A B` | 按顺序匹配 A 和 B |
| `A / B` | 先尝试 A，失败后回到起点尝试 B |
| `(A / B)` | 分组 |
| `A?` | 零次或一次 |
| `A*` | 零次或多次，贪婪匹配 |
| `A+` | 一次或多次，贪婪匹配 |
| `&A` | A 成功则成功，但不消费输入 |
| `!A` | A 失败则成功，但不消费输入 |
| `&&A` | A 成功且实际消费输入时成功 |
| `!!A` | `&&A` 的逻辑反向扩展 |
| `label:A` | 将 A 的语义值绑定到 label |
| `label:<A>` | 将 A 匹配到的原文捕获为字符串 |

### 有序选择

顺序会改变语言：

```peg
BadOperator  = "<" / "<="
GoodOperator = "<=" / "<"
```

`BadOperator` 永远不会把 `<=` 作为完整的第二分支匹配，因为第一分支已经接受了 `<`。通常应把更具体的分支放在前面，并在入口末尾使用 `!.` 拒绝未消费的尾部输入。

### 重复必须前进

`*` 或 `+` 的子表达式应当消费至少一个字符，或者失败。避免下面这种可空重复：

```peg
// 错误示例：("a"?)* 可以成功却不前进
Loop = ("a"?)*
```

## 字面量

```peg
Keyword = "select"i
Rune    = 'x'
Raw     = `no escapes here`
```

双引号表示字符串，单引号只表示一个字符，反引号表示原始字符串。结束引号后的 `i` 表示忽略大小写。

## 字符类与 Unicode

```peg
Digit       = [0-9]
Identifier  = [\pL_] [\pL\p{Nd}_]*
GreekLetter = [\p{Greek}]
NotNewline  = [^\r\n]
Hex         = [0-9a-f]i
```

字符类支持单个字符、范围、转义、首字符 `^` 反选，以及 Unicode 类：

- `\pL` 使用单字母通用类别。
- `\p{Nd}` 使用命名类别。
- `\p{Greek}` 使用 Unicode script/property 名称。

可用名称由当前生成器的 Go Unicode 数据决定。默认生成物会保留生成时的类别语义，而不是依赖浏览器的 `\p{...}` 正则实现。Haxe 的外部 `hxUnicode` 模式是例外，详见[Haxe target](/targets/haxe#unicode-表)。

`.` 匹配一个 Unicode 字符但不匹配 EOF，因此惯用的结束规则是：

```peg
EOF = !.
```

## 标签、捕获与 action

`label:expr` 保存表达式的语义值，`label:<expr>` 保存匹配文本。复杂文本捕获需要分组：

```peg
Assignment = text:<(Name _ "=" _ Value)> { return text }
```

具体返回类型和 repetition 的收集规则见[动作与值](/guide/actions-and-values)。

## 语义谓词与代码表达式

```peg
Allowed = &{ return c.data.AllowValue } Value
Blocked = !{ return c.data.Blocked } .
Probe   = *{ c.data.ProbeCount++ }
```

- `&{...}` 在代码返回真值时成功。
- `!{...}` 对代码结果取反。
- 普通 action 在 syntactic lookahead 中会跳过。
- `*{...}` 是即使处于 lookahead 中也执行的代码表达式。

代码内容及返回类型取决于 target。

## 错误处理

规则显示名称、no-match formatter、action 语义错误以及 `%{Label}` / `//{...}` 恢复语法统一见[错误处理](/guide/error-handling)。

## 左递归

当前 Pegtool builder 会拒绝左递归：

```peg
// 不支持
Expr = Expr "+" Term / Term
```

应改写为先匹配首项，再重复匹配尾项：

```peg
Expr = Term (_ "+" _ Term)*
```
