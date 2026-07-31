# TypeScript

TypeScript target 生成可在 `--strict` 下编译的独立模块，适合浏览器、Node.js 和其他 JavaScript runtime。

## 最小可运行示例

在空目录中手写下面两个文件。`parser.ts` 由 Pegtool 生成。

### `identifier.peg`

```peg
Start = value:<Identifier> !. { return value; }
Identifier = [\pL_] [\pL\p{Nd}_]*
```

### `main.ts`

```ts
import { parse } from "./parser.js";

const value = parse("name_世界");
if (typeof value !== "string") {
  throw new Error("unexpected parser result");
}
console.log(value);
```

生成、严格编译并运行：

```powershell
pegtool -t ts -o parser.ts identifier.peg
tsc --strict --target ES2020 --module NodeNext --moduleResolution NodeNext parser.ts main.ts
node main.js
```

输出：

```text
name_世界
```

## 生成 API

生成模块导出 `parse`、`Parser`、`ParserOptions`、`PigeonParseError`、`Position`、`Current`、`Grammar` 和文法变量。完整 options 调用如下：

```ts
const value = parse(source, {
  filename: "input.txt",
  entrypoint: "Start",
  maxExpressions: 1_000_000,
  customData: { featureEnabled: true },
});
```

## Action

```peg
{
function toInteger(text: string): number {
  return Number.parseInt(text, 10);
}
}

Input = value:Integer !. { return value; }
Integer = [0-9]+ { return toInteger(c.text); }
```

标签和 action 返回值为 `any`。`c.text` 是 JavaScript 字符串，`c.data` 来自 `ParserOptions.customData`。

## Unicode

Unicode class 会在生成阶段转换成紧凑的共享 number table。runtime 使用生成器固化的范围进行匹配，不依赖浏览器的 `\p{...}` RegExp 支持或 Unicode 版本。普通字符、范围、忽略大小写、反选和非 BMP code point 走同一字符类语义。

## 优化与 runtime

当前 TypeScript runtime 不消费 `-optimize-parser`；加不加该参数会生成逐字节相同的文件。TypeScript 也没有运行时 memo 开关。

`-optimize-ref-expr-by-index` 可以跳过规则名称查表，但会让规则插入或重排产生大面积 diff。默认关闭，只在实际负载证明收益且可接受生成物变化时开启。

正常 Node.js 会使用 V8 JIT。`node --jitless` 只适合测量无 JIT 的绝对执行路径，不是生产性能参数。tjs/QuickJS 可以验证无 V8 环境的兼容性；target 仍然是直接 TypeScript，生成命令不变。
