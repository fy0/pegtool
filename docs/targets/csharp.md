# C#

C# target 生成一个独立的 `partial` `PTParser` 类，适合 .NET 8 及启用 nullable reference types 的项目。

## 最小可运行示例

先在空目录中创建一个控制台项目：

```powershell
dotnet new console --framework net8.0
```

保留生成的项目文件，新建 `identifier.peg`，并用下面内容替换 `Program.cs`。`Parser.cs` 由 pegtool 生成。

### `identifier.peg`

```peg
Start = value:<Identifier> !. { return value; }
Identifier = [\pL_] [\pL\p{Nd}_]*
```

### `Program.cs`

```csharp
var value = (string)PTParser.Parse("name_世界")!;
Console.WriteLine(value);
```

生成并运行：

```powershell
pegtool -t cs -o Parser.cs identifier.peg
dotnet run
```

输出：

```text
name_世界
```

## 公开 API

`PTParser.Parse` 返回 `object?`，失败时抛出包含 `Position` 和 `Expected` 的 `PTParseException`。完整 options 调用如下：

```csharp
var value = PTParser.Parse(source, new PTParser.ParserOptions
{
    Filename = "input.txt",
    Entrypoint = "Start",
    MaxExpressions = 1_000_000,
    CustomData = context,
});
```

## Initializer 与 action

C# initializer 会放进生成的 `partial PTParser` 类内部，因此应声明嵌套类型、字段和 helper 方法，不要写 namespace 或顶层语句：

```peg
{
private static int ToInteger(string text) => int.Parse(text);
}

Integer = [0-9]+ { return ToInteger(c.Text); }
```

标签类型为 `object?`。当 PEG 标签也是 C# 关键字时，生成 action 中使用 verbatim identifier，例如标签 `operator` 通过 `@operator` 访问。

## 分配行为

当前 runtime 不为普通 sequence 创建 `List<object?>`，repetition 只在子表达式返回非空值时分配结果集合。测试中的 action-free 长重复负载相对旧实现快约 `18%-19%`，且无用集合分配基本不随输入长度增长。这项优化自动生效，不需要额外参数。

`-optimize-parser` 当前没有 C# 专用效果；规则索引开关仍有 diff 稳定性代价，按[参数指南](/guide/performance-options)决定。
