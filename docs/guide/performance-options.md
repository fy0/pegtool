# 性能参数怎么选

性能参数应按文法行为选择，而不是全部开启。Pegtool 中最重要的三个取舍分别是：固定 runtime 优化、命名规则 memoization、规则引用索引化。它们解决的问题不同，也有不同的可维护性成本。

## 先给结论

如果生成文件参与 Git 审阅，使用下面的稳定配置作为起点：

| 场景 | 推荐生成参数 | 运行时 |
| --- | --- | --- |
| Go 普通生产文法 | `-t go -optimize-parser` | 默认不 memo |
| Go 明显重复回溯 | `-t go -optimize-parser` | 基准确认后 `memoized(true)` |
| Haxe 普通、近线性文法 | `-t hx -optimize-parser` | 该生成物会裁掉 memo |
| Haxe 明显重复回溯 | `-t hx` | 普通生成物可设置 `parser.memoize = true` |
| TypeScript | `-t ts` | 当前无 memo |
| C#、C99、Rust | 只选择对应 target | 当前无 memo |

完整生成命令如下：

```powershell
# Go：普通和高回溯文法使用相同生成命令；高回溯场景在运行时启用 memo
pegtool -t go -optimize-parser -o parser.go grammar.peg

# Haxe：普通、近线性文法，生成物不包含 memo
pegtool -t hx -optimize-parser -o Parser.hx grammar.peg

# Haxe：高回溯文法，保留 memo 支持，运行前设置 parser.memoize = true
pegtool -t hx -o Parser.hx grammar.peg

# TypeScript
pegtool -t ts -o parser.ts grammar.peg

# C#
pegtool -t cs -o Parser.cs grammar.peg

# C99
pegtool -t c -o parser.c grammar.peg

# Rust
pegtool -t rust -o parser.rs grammar.peg
```

Haxe 有两种生成形态，并不是同一份生成物同时“有 memo”又“没有 memo”：

- 加 `-optimize-parser`：生成物更精简，memo 字段、函数和缓存分支会被完全删除。
- 不加 `-optimize-parser`：生成物保留可选 memo，默认关闭；只有设置 `parser.memoize = true` 才会启用。

默认不加 `-optimize-ref-expr-by-index`。它可能带来额外吞吐，但会让规则插入或重排产生大量生成 diff。

## 决策顺序

1. 先用默认配置验证语义和错误输出。
2. 固定生产文法使用该 target 真正支持的 `-optimize-parser`。
3. 观察是否有同一命名规则在同一输入位置被反复求值。
4. 只有第 3 步成立时才测试 memo。
5. 只有仍需要最后一段吞吐、并接受大 diff 时才测试规则索引。
6. 在真实输入上做交错、反序复测，不根据一次绝对时间决定。

## `-optimize-parser`

这个参数不是跨 target 的统一开关：

| Target | 当前效果 |
| --- | --- |
| Go | 删除 Debug、Statistics 热路径并减少不必要的变量栈操作；仍保留命名规则 memo |
| Haxe | 删除 Debug、Statistics、memo 字段与缓存分支，并展平部分 wrapper |
| TypeScript | 当前不消费该参数，输出逐字节相同 |
| C#、C99、Rust | 当前没有对应的 target 专用 runtime 分支 |

因此 Go 可以组合：

```powershell
pegtool -t go -optimize-parser -o parser.go grammar.peg
```

```go
value, err := Parse("input.txt", input, Memoize(true))
```

Haxe 则必须二选一：

```powershell
# 固定、近线性
pegtool -t hx -optimize-parser -o Parser.hx grammar.peg

# 高回溯，需要 memo
pegtool -t hx -o Parser.hx grammar.peg
```

## 命名规则 memoization

memo 以命名规则、输入位置和解析模式为边界缓存成功与失败结果。它避免回溯重复扫描，但每次首次访问都需要查表和存储。

适合：

- 多个有序分支共享昂贵前缀规则。
- 同一失败规则会从同一 offset 多次扫描较长输入。
- predicate 的输入配置在一次 parse 期间稳定。

不适合：

- 文法近似线性，没有有效缓存命中。
- 被重复调用的规则只匹配一个简单字面量，重算比查表便宜。
- action 或 predicate 依赖回溯期间变化的 `c.data`。
- 每次规则访问都必须重新发生副作用。

测试中，HashLink 的昂贵重复扫描负载启用 memo 后快 `2.53x-2.61x`；线性负载慢约 `1%-1.4%`，简单失败规则的探索性负载慢约 `24%`。这说明“有回溯”还不够，重复工作本身必须比缓存操作昂贵。

Go 的 DiceScript 真实测试中，把缓存边界从每个表达式收缩到命名规则后，完整根测试套件中 parser 相关执行从约 `199.8 ms` 降到 `58.3 ms`。主要收益来自减少缓存项和分配，并不是消除 `any` 类型分派。

## `-optimize-ref-expr-by-index`

默认规则引用按名称查找。索引模式把它改成：

```text
rules["Expression"]  ->  rulesArray[42]
```

这会跳过名称 map 查找，但序号直接取决于规则顺序。若在第 10 条插入规则，后续规则的编号以及所有指向它们的引用都可能变化。

### 什么时候开启

- 生成文件不提交或不参与人工审阅。
- 文法规则顺序基本冻结。
- 目标负载基准证明收益足够覆盖维护成本。
- 运行时不会重排或拼接 `grammar.rules`。

### 什么时候关闭

- 希望小改文法只产生局部 diff。
- 多人频繁插入、重排规则。
- 生成 parser 需要长期审计。

历史上的约 `10%` 是总体估计，不是当前 Go、Haxe、TypeScript 各自的保证。本轮跨 target 测试也没有隔离该参数，因此 Pegtool 文档不把它列为默认项。

若确实追求极限吞吐，可以在稳定配置后追加：

```text
-optimize-ref-expr-by-index
```

## 其他容易混淆的参数

### `-cache`

只缓存 **Pegtool 读取 `.peg` 文法** 时的结果。它不改变生成文件，也不会给生成 parser 启用 memo。只有生成器解析病态文法本身很慢时才测试它；普通文法可能更慢并消耗更多内存。

### `-nolint`

只控制 Go 生成注释，不影响 parser 吞吐或生成算法。Haxe、TypeScript 当前不消费它。

### `-haxe-use-hxunicode`

这是 Haxe 源码体积与依赖选择，不是性能开关。默认嵌入生成器的 Unicode 表；外部模式需要 `hxUnicode`，文件更小，但字符集版本可能随依赖变化。

### `-alternate-entrypoints`

当前只在生成阶段验证规则名存在，不改变输出，也不创建 runtime selector。入口选择应通过各 target 的 `ParserOptions` 或项目 wrapper 完成。

## JavaScript 与 HashLink 运行参数

- Node 生产运行使用默认 V8 JIT，不加 `--jitless`。
- `node --jitless` 用于确认优化不是纯粹的预热波动，不是加速参数。
- tjs 使用 QuickJS 路径，适合无 V8 环境验证；不需要不同的生成参数。
- HashLink 使用正常 HL/JIT；当前没有对应的无 JIT 开关。
- Haxe 编译 JavaScript 或 HashLink时建议加 `-D analyzer-optimize`。

在跨 target 回溯负载中，Haxe memo 版本的耗时约为直接 TypeScript 的 `74%-77%`，在 Node JIT、Node jitless 和 tjs 三种模式下都保持优势；线性负载则应优先使用 Haxe `-optimize-parser`，而不是为没有命中的缓存付费。

## 如何复测自己的文法

一次可靠比较至少应做到：

1. 对每个变体先验证返回值、错误和消费位置一致。
2. 使用短输入、典型输入和最坏回溯输入。
3. 每个变体独立预热，并自动校准批次数量。
4. 交替执行 A/B，另起进程反转为 B/A。
5. 同时记录吞吐、分配次数和内存，不只看一次墙钟时间。
6. Node 结果同时做一次 `--jitless` 控制；QuickJS/HL 结果单独报告。

机器负载会显著改变绝对时间。优先看同一进程内相邻、成对的比率，以及分配数量这类更稳定的指标。
