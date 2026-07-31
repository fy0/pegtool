# 项目开发

本页面向修改 pegtool 生成器本身的贡献者。只使用生成 parser 的项目不需要这些步骤。

## 仓库结构

| 路径 | 内容 |
| --- | --- |
| `grammar/` | pegtool 自身 PEG 文法 |
| `ast/` | 文法 AST 与遍历 |
| `builder/` | 共享 builder 与 Go runtime |
| `builder_hx/` | Haxe target |
| `builder_ts/` | TypeScript target |
| `builder_cs/` | C# target |
| `builder_c/` | C99 target |
| `builder_rust/` | Rust target |
| `examples/` | 各语言可编译示例 |
| `test/` | parser 行为回归测试 |

## 工具链

当前 CI 使用 Go 1.25。不同 target 的完整验证还需要相应工具链，例如 Haxe/HashLink、Node/TypeScript、.NET、C compiler 和 Cargo。

```powershell
go build .
go test ./...
```

修改 Go 或 Haxe 静态 runtime 源后，要同步再生模板镜像：

```powershell
go generate ./builder ./builder_hx
```

不要直接只改 `generated_static_code.go`；源文件和生成镜像必须保持一致。

## 验证 target

按变更范围至少执行对应包测试：

```powershell
go test ./builder ./builder_ts ./builder_hx ./builder_cs ./builder_c ./builder_rust
```

涉及共享 AST、grammar 或 CLI 时执行：

```powershell
go test ./...
```

生成器测试只能证明输出结构。用户可见 target 改动还应实际使用 `tsc --strict`、Haxe、`dotnet build`、C99 compiler 或 `cargo test` 编译并运行生成 parser。

## 基准测试原则

- correctness 先于吞吐。
- 在同一机器、同一 toolchain 下交错比较。
- 对 JIT runtime 分离预热和测量，并反转执行顺序复测。
- 同时记录生成源码大小、执行时间、分配次数和内存。
- 用户文档只发布已经复核的稳定结论，不发布中间实验记录。

## 文档站

安装依赖并构建：

```powershell
npm install
npm run docs:build
```

`package.json` 对 Vite 使用了安全版本覆盖；更新 VitePress 后应重新运行 `npm audit`，确认覆盖仍然必要且兼容。

本地编辑时可使用 VitePress 的文档脚本：

```powershell
npm run docs:dev
```

发布前以 `docs:build` 为准；输出位于 `docs/.vitepress/dist/`，不会提交到 Git。

## 提交前

1. 格式化 Go 代码并运行 goimports。
2. 再生受影响的生成文件。
3. 运行对应 target 测试与 `go test ./...`。
4. 构建 VitePress 文档。
5. 检查生成物是否出现与语法改动无关的大面积 diff。
