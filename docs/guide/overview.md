# 项目概览

Pegtool 是一个使用 Go 编写的 PEG（Parsing Expression Grammar）解析器生成器。它读取 `.peg` 文法，在构建阶段生成目标语言源码；部署生成的 parser 时不需要携带 Pegtool 本身。

项目于 2024 年从 Pigeon 分叉，分叉基线包含 Pigeon 2023 年的 `v1.2.x` 发布线。两者的 PEG 写法总体相似，但 Pegtool 的文法、值语义、生成 API 和运行时已经独立演进。这个 fork 的首要原因是性能，此后围绕生成物的执行效率进行了大量优化；为了消除隐式值数组、状态复制和动态查找等成本，部分优化不得不调整文法和值语义。最初的分析和测量记录在 [Pigeon #151](https://github.com/mna/pigeon/issues/151)。

## 为什么建立这个 fork

我有一个 TRPG 网络跑团项目，并用 Go 为它实现了一个骰点表达式解释器。这个解释器最早使用的是 [`pointlander/peg`](https://github.com/pointlander/peg)。它很快，但生成代码中有大量 `goto`，文法稍微一改就会产生大片 diff。Go package 通常又要提交生成的 `.go` 文件，因此提交记录很容易变乱，也不方便审阅和自定义。

为了解决生成物难以审阅和自定义的问题，项目随后迁移到了 Pigeon。Pigeon 的生成结构更简单、清晰，也更适合继续扩展；迁移完成后才发现，它在项目真实文法上的 parser 性能明显不足。2024 年的早期测量中，同一文法使用原版 Pigeon 约需 `4.5s`，此前的 `pointlander/peg` 约需 `0.3s`。通过移除昂贵的 state 路径、减少 sequence 中间值、压缩规则查找等优化，早期分支将 Pigeon 路径降到了约 `0.5s`。这是 2024 年的事情了，后来又陆续做了一些其他优化和修改。

这个 fork 就诞生在第二次迁移之后：保留 Pigeon 较容易阅读、审阅和移植的生成结构，再有针对性地优化热路径。代价是部分性能改动无法完全保持原版的隐式值语义，这也是迁移时必须重新生成并运行测试的原因。

## 类似工具

Go 生态中有多种 PEG parser 生成器。它们的生成方式和设计目标不同，适合的项目也不同：

| 工具 | 生成方式与特点 | 更适合的场景 |
| --- | --- | --- |
| [`pointlander/peg`](https://github.com/pointlander/peg) | 生成高度展开、包含大量 `goto` 的 Go parser，执行速度通常较好；文法改动可能产生大面积生成 diff | 吞吐优先，生成文件不需要长期审阅，文法结构相对稳定 |
| [原版 Pigeon](https://github.com/mna/pigeon) | 生成结构直观，文法和 runtime 容易理解与扩展，保持上游语法和 API | 重视上游兼容性、可扩展性，且现有负载能够接受其性能 |
| Pegtool | 保留可审阅的生成结构，优化 parser 热路径，并支持多个 target；部分值语义与原版 Pigeon 不同 | 需要稳定生成 diff、自定义 runtime 或诊断接口，或者需要生成多种目标语言 |

没有一种工具在所有文法上都更好。选型时应同时比较真实输入上的性能、错误诊断、生成 diff、目标语言，以及生成文件是否需要提交和审阅。

## 从 Haxe 到多语言 target

2024 年，项目作者开始学习 Haxe，最初希望使用 Heaps.io 引擎；另外，Haxe/JavaScript 的转译在作者看来做得非常好。顺便吐槽一下 Go：大概是 GopherJS 的维护人员没有那么多精力，也缺少捐赠支持，当时的演进速度比较慢，并且缺少死代码清除（DCE），导致输出包体过大。不过 Haxe 语言本身很有趣，但使用 Heaps.io 在实际尝试中遇到了一些问题，这条应用路线最终没有继续。

这次尝试仍然留下了项目最重要的结构变化之一：受 Haxe 多 target 思路启发，Pigeon builder 被扩展为可以选择目标语言，并加入了 Haxe target。项目随后用它实现过一套叙事语言脚本引擎，验证了同一 PEG 核心生成不同宿主语言 parser 的可行性。

此后近两年，作者主要忙于其他工作，项目演进一度放缓。到 2026 年又补充了更多目标语言和相应 runtime，形成了现在支持 Go、Haxe、TypeScript、C#、C99 与 Rust 的 Pegtool。

## 支持的目标语言

| Target | 参数 | 生成物 |
| --- | --- | --- |
| Go | `-t go` | Go 源文件 |
| Haxe | `-t hx` | 可继续编译到 JavaScript、HashLink 等后端的 Haxe 源文件 |
| TypeScript | `-t ts` | TypeScript 模块 |
| C# | `-t cs` | `partial` parser 类 |
| C99 | `-t c` | C99 源文件 |
| Rust | `-t rust` | Rust 模块 |

各 target 共享 PEG 的匹配结构，但 initializer、action、predicate 和公开 API 使用对应的宿主语言。具体差异见[目标语言概览](/targets/)。

## 基本工作流

```text
grammar.peg
    -> pegtool 验证并生成
    -> parser.go / Parser.hx / parser.ts / ...
    -> 目标语言编译器
    -> 应用程序
```

生成命令的基本形式是：

```powershell
pegtool -t go -o parser.go grammar.peg
```

生成物默认只依赖目标语言标准库；文法 action 主动引用的依赖，以及 Haxe 显式选择的 `hxUnicode` 模式除外。

## 当前语义重点

- 文本通过 `label:<expr>` 显式捕获。
- 普通终结符和 sequence 不创建隐式语义值数组。
- repetition 只收集子表达式显式返回的非空值。
- Unicode class 由生成时的 Unicode 数据确定。
- 左递归会在生成阶段被拒绝，应改写为首项加尾项重复。
- memoization 和规则索引优化都应按真实负载选择，不是默认全部开启。

## 从哪里继续

- 第一次使用：从[快速开始](/guide/getting-started)生成一个 Go parser。
- 编写规则：阅读 [PEG 语法](/guide/grammar)和[动作与值](/guide/actions-and-values)。
- 从原版 Pigeon 迁移：使用[自动迁移指南与 Prompt](/guide/fork-semantics)。
- 准备生产配置：阅读[性能参数选择](/guide/performance-options)。
- 查找所有参数：查看[命令行参考](/reference/cli)。
