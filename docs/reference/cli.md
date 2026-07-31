# 命令行参考

```text
pegtool [options] [GRAMMAR_FILE]
```

不提供文件时从 stdin 读取文法；不提供 `-o` 时把生成源码写到 stdout。

## 输出与 target

| 参数 | 说明 |
| --- | --- |
| `-o FILE` | 写入指定文件 |
| `-t TARGET`、`-target TARGET` | 选择 `go`、`hx`、`ts`、`cs`、`c` 或 `rust`；默认 `go` |
| `-x` | 只解析并验证文法，不生成 parser |
| `-receiver-name NAME` | 修改 action 中当前匹配对象的默认名称 `c` |
| `-nolint` | 为 Go 生成物加入 lint 抑制注释 |

Target 别名：

| 规范值 | 别名 |
| --- | --- |
| `hx` | `haxe` |
| `ts` | `typescript` |
| `cs` | `csharp`、`c#` |
| `c` | `c99` |

## 生成性能与 runtime 形态

| 参数 | 说明 |
| --- | --- |
| `-optimize-parser` | 生成固定生产 runtime；不同 target 的效果不同 |
| `-optimize-ref-expr-by-index` | 用规则数组下标替代名称查找，可能产生大面积 diff |
| `-cache` | 只缓存生成器解析 `.peg` 的过程，不影响生成 parser |
| `-debug` | 输出生成器解析 `.peg` 时的调试信息 |

不要把 `-cache` 当成生成 parser 的 memo 开关。完整选择规则见[性能参数怎么选](/guide/performance-options)。

## 入口与多文法

| 参数 | 说明 |
| --- | --- |
| `-alternate-entrypoints A,B` | 验证所列规则存在；当前不改变生成 API 或输出 |
| `-grammar-only` | 只生成文法对象与 action，不重复 runtime |
| `-grammar-name NAME` | 设置文法变量名，默认 `g` |
| `-run-func-prefix PREFIX` | 给生成 action callback 添加前缀 |

目前明确支持和测试的是 Go 同包多文法工作流。完整示例见[多文件与入口规则](/guide/multiple-grammars)。

## Haxe 专用

| 参数 | 说明 |
| --- | --- |
| `-haxe-use-hxunicode` | Unicode class 改用 `hxUnicode`，减少嵌入源码但增加依赖 |

该参数用于非 Haxe target 时会报错。没有 Unicode class 的文法不会因为它引入外部依赖。

## 帮助

```powershell
pegtool -h
pegtool -help
```

## 常用命令

```powershell
# 验证文法
pegtool -x grammar.peg

# 稳定 diff 的 Go 生产 parser
pegtool -t go -optimize-parser -o parser.go grammar.peg

# 可选择 memo 的 Haxe parser
pegtool -t hx -o Parser.hx grammar.peg

# 严格 TypeScript parser
pegtool -t ts -o parser.ts grammar.peg
```
