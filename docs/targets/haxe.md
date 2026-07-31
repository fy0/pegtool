# Haxe

Haxe target 可以编译到 JavaScript、HashLink 和其他 Haxe 后端。本项目重点验证了 JavaScript 与 HashLink。

## 最小可运行示例

在空目录中手写下面两个文件。`Parser.hx` 由 pegtool 生成。

### `identifier.peg`

```peg
{
class ParserCustomData {
    public function new() {}
}
}

Start = value:<Identifier> !. { return value; }
Identifier = [\pL_] [\pL\p{Nd}_]*
```

### `Main.hx`

```haxe
class Main {
    static function main() {
        var parser = new Parser("input.txt", "name_世界", []);
        var value:String = cast parser.parse(null);
        #if js
        js.Lib.global.console.log(value);
        #else
        Sys.println(value);
        #end
    }
}
```

生成并编译到 JavaScript：

```powershell
pegtool -t hx -optimize-parser -o Parser.hx identifier.peg
haxe -cp . -main Main -js main.js -D analyzer-optimize
node main.js
```

同样两个文件也可以编译到 HashLink：

```powershell
haxe -cp . -main Main -hl main.hl -D analyzer-optimize
hl main.hl
```

两种运行方式都会输出：

```text
name_世界
```

## Parser 与自定义数据

`new Parser(filename, input, options)` 创建 parser，`parse(null)` 使用生成文件中的默认 grammar。选择其他入口时，在 `parse` 前设置 `parser.entrypoint`。需要自定义数据时可设置公开的 `parser.cur.data`，或在生成模块中提供项目 wrapper。

## `-optimize-parser` 与 memo

Haxe 中二者互斥：

- `-optimize-parser` 会删除 debug、statistics、memo 字段和缓存分支，适合固定、近线性的生产文法。
- 高回溯文法需要 memo 时，不加 `-optimize-parser`，解析前设置 `parser.memoize = true`。

```haxe
var parser = new Parser("input.txt", source, []);
parser.memoize = true;
var value = parser.parse(null);
```

在 HashLink 的重复扫描失败负载中，命名规则 memo 相对同版本 no-memo 快 `2.53x-2.61x`；线性负载则慢约 `1%-1.4%`，缓存简单失败规则的探索性负载甚至慢约 `24%`。因此它不是默认开关。

## Unicode 表

默认模式把文法实际使用的 Unicode 表嵌入 `Parser.hx`，无第三方依赖，字符集由生成器版本固定。若 parser 源文件体积更重要，可以改用外部库：

```powershell
haxelib install hxUnicode
pegtool -t hx -haxe-use-hxunicode -o Parser.hx identifier.peg
haxe -cp . -main Main -js main.js -lib hxUnicode -D analyzer-optimize
```

`-haxe-use-hxunicode` 是体积与依赖选择，不是已验证的性能参数。外部库的 Unicode 版本也可能和生成器不同。

不要把 `-optimize-ref-expr-by-index` 作为默认参数。它会在插入或重排规则后改写大量数字引用，只有生成物不参与审阅且真实基准确认收益时才启用。
