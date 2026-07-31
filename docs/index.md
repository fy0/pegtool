---
layout: home

hero:
  name: Pegtool
  text: 使用 Go 编写的 PEG 解析器生成器
  tagline: 可为 Go、Haxe、TypeScript、C#、C99 与 Rust 生成可审阅、可部署的解析器源码。
  actions:
    - theme: brand
      text: 快速开始
      link: /guide/getting-started
    - theme: alt
      text: 从 Pigeon 迁移
      link: /guide/fork-semantics

features:
  - title: 大量性能优化
    details: 针对真实文法持续优化隐式值、状态、规则查找、缓存和 Unicode 热路径。
  - title: 生成代码易于理解
    details: 日常使用无需阅读生成物；清晰的结构让自定义 runtime、诊断接口或新 target 更容易。
  - title: 默认仅依赖标准库
    details: 生成的 parser 默认只依赖目标语言标准库。
  - title: 多种目标语言
    details: 可按项目的部署环境选择 Go、Haxe、TypeScript、C#、C99 或 Rust target。
---

## Pegtool 是什么

Pegtool 是一个使用 Go 编写的 PEG（Parsing Expression Grammar）解析器生成器，可以从同一套文法生成 Go、Haxe、TypeScript、C#、C99 和 Rust 源码。项目在 2024 年从 Pigeon 分叉，基线包含 Pigeon 2023 年 `v1.2.x` 发布线；此后文法和运行时已经独立演进。这个 fork 的首要原因是改善原版 Pigeon 在实际文法中的性能，围绕这一目标进行了大量优化；部分优化不得不调整文法和值语义，早期分析见 [Pigeon #151](https://github.com/mna/pigeon/issues/151)。

```powershell
pegtool -t go   -o parser.go grammar.peg
pegtool -t hx   -o Parser.hx grammar.peg
pegtool -t ts   -o parser.ts grammar.peg
pegtool -t cs   -o Parser.cs grammar.peg
pegtool -t c    -o parser.c grammar.peg
pegtool -t rust -o parser.rs grammar.peg
```

## 从哪里开始

- 第一次使用：阅读[快速开始](/guide/getting-started)。
- 从原版 Pigeon 迁移：使用[自动迁移指南与 Prompt](/guide/fork-semantics)。
- 准备上线：按负载选择[性能参数](/guide/performance-options)。
- 编写跨语言文法：查看[目标语言概览](/targets/)。

::: tip 先保持默认
除非基准测试证明同一命名规则在同一输入位置被反复求值，否则先不要启用 memo。除非生成文件不参与代码审阅，否则先不要启用规则索引优化。
:::
