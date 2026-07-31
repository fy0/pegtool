# 多文件与入口规则

Pegtool 可以让一个 Go parser runtime 承载多份文法。当前经过明确设计和测试的是同一 Go package 内的 `grammar-only` 工作流；其他 target 暂不建议依赖这一能力。

## 生成主文法

```powershell
pegtool -t go -o base.peg.go base.peg
```

主文件包含 parser runtime、默认文法变量 `g` 和 action callbacks。

## 生成附加文法

```powershell
pegtool -t go `
  -grammar-only `
  -grammar-name=gOperators `
  -run-func-prefix=operators_ `
  -o operators.peg.go `
  operators.peg
```

三个参数分别解决不同问题：

| 参数 | 作用 |
| --- | --- |
| `-grammar-only` | 不重复输出 parser runtime |
| `-grammar-name` | 给附加文法一个唯一变量名，避免都叫 `g` |
| `-run-func-prefix` | 给 action callback 加唯一前缀，避免方法重名 |

名称和前缀必须是目标语言可接受的标识符片段。所有生成的 Go 文件必须属于同一个 package，并且 `ParserCustomData`、公开 wrapper 等共享声明只能定义一次。

可以利用模板变量只在主文件输出共享声明：

```peg
{
package parser

{{ if not .GrammarOnly }}
type ParserCustomData struct{}
{{ end }}
}
```

## 选择附加文法与入口

```go
func ParseOperator(filename string, input []byte) (any, error) {
    p := newParser(filename, input)
    p.entrypoint = "Operator"
    return p.parse(gOperators)
}
```

parser 的默认入口来自主文法第一条规则。附加文法首条规则名称不同时，应显式设置 `p.entrypoint`。

## `-alternate-entrypoints` 的现状

当前 CLI 的 `-alternate-entrypoints=A,B` 只验证这些规则存在，不会生成公开 selector，也不会改变生成文件。运行时入口应由各 target 的 API 选择：

- Go：在 initializer wrapper 中设置 `p.entrypoint`。
- Haxe：解析前设置 `parser.entrypoint`。
- TypeScript、C#、C99、Rust：使用对应 `ParserOptions` 的 entrypoint 字段。

## 索引模式的一致性

若确实启用 `-optimize-ref-expr-by-index`，主 runtime 和所有 grammar-only 片段应保持一致。索引引用依赖每份 grammar 自己的规则顺序，也不应在运行时重排规则数组。

::: warning Target 支持边界
TypeScript grammar-only 输出不会自动导入主模块中的私有 runtime helper；Haxe 目前也只有输出格式层面的覆盖。除 Go 同包方案外，应先做独立编译实验再用于生产。
:::
