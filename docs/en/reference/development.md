# Development

This page is for contributors modifying the pegtool generator itself. Projects that only use generated parsers do not need these steps.

## Repository structure

| Path | Contents |
| --- | --- |
| `grammar/` | pegtool's own PEG grammar |
| `ast/` | Grammar AST and traversal |
| `builder/` | Shared builder and Go runtime |
| `builder_hx/` | Haxe target |
| `builder_ts/` | TypeScript target |
| `builder_cs/` | C# target |
| `builder_c/` | C99 target |
| `builder_rust/` | Rust target |
| `examples/` | Compilable examples for each language |
| `test/` | Parser behavior regression tests |

## Toolchains

Current CI uses Go 1.25. Complete validation of the other targets also requires the relevant toolchains, such as Haxe/HashLink, Node/TypeScript, .NET, a C compiler, and Cargo.

```powershell
go build .
go test ./...
```

After modifying static Go or Haxe runtime source, regenerate the template mirrors:

```powershell
go generate ./builder ./builder_hx
```

Do not modify only `generated_static_code.go`; source files and generated mirrors must stay synchronized.

## Validate targets

Run at least the tests for affected packages:

```powershell
go test ./builder ./builder_ts ./builder_hx ./builder_cs ./builder_c ./builder_rust
```

When shared AST, grammar, or CLI code changes, run:

```powershell
go test ./...
```

Generator tests prove only output structure. A user-visible target change should also compile and run a generated parser with `tsc --strict`, Haxe, `dotnet build`, a C99 compiler, or `cargo test`.

## Benchmark principles

- Correctness comes before throughput.
- Interleave comparisons on the same machine and toolchain.
- Separate warm-up from measurement for JIT runtimes and repeat with reversed execution order.
- Record generated source size, execution time, allocation counts, and memory.
- Publish only reviewed, stable conclusions in user documentation, not intermediate experiments.

## Documentation site

Install dependencies and build the site:

```powershell
npm install
npm run docs:build
```

`package.json` applies a secure version override for Vite. After updating VitePress, rerun `npm audit` to confirm that the override remains necessary and compatible.

For local editing, use the VitePress documentation script:

```powershell
npm run docs:dev
```

Treat `docs:build` as the release check. It writes output to `docs/.vitepress/dist/`, which is not committed to Git.

Pull requests that change the documentation run the same production build in
GitHub Actions. A documentation change pushed to `main` is built and deployed to
GitHub Pages automatically by the `Deploy documentation` workflow.

## Before committing

1. Format Go code and run goimports.
2. Regenerate affected generated files.
3. Run the relevant target tests and `go test ./...`.
4. Build the VitePress documentation.
5. Check generated output for large diffs unrelated to the syntax change.
