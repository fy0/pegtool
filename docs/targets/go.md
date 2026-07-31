# Go

Go 是默认 target。生成器会运行 goimports 格式化输出，initializer 必须提供合法的 `package` 声明。

## 最小可运行示例

在空目录中手写下面两个文件。`parser.go` 由 Pegtool 生成，不需要手工创建。

### `identifier.peg`

```peg
{
package main

type ParserCustomData struct{}

func ParseIdentifier(input []byte) (string, error) {
    value, err := parse("input.txt", input)
    if err != nil {
        return "", err
    }
    return value.(string), nil
}
}

Start = value:<Identifier> !. { return value }
Identifier = [\pL_] [\pL\p{Nd}_]*
```

### `main.go`

```go
package main

import "fmt"

func main() {
    value, err := ParseIdentifier([]byte("name_世界"))
    if err != nil {
        panic(err)
    }
    fmt.Println(value)
}
```

生成并运行：

```powershell
go mod init example.com/identifier
pegtool -t go -optimize-parser -o parser.go identifier.peg
go run .
```

输出：

```text
name_世界
```

## 公开 API 与自定义数据

生成 runtime 的基础符号是包内的 `parse`、`newParser`、`option` 和 `memoized`。上例在 initializer 中把内部 `parse` 包装成了项目需要的 `ParseIdentifier`。需要 `ParseFile`、`ParseReader` 或其他签名时，也可以在 initializer 或普通 Go 文件中包装。

即使不需要状态，也必须定义空的 `ParserCustomData`。需要注入数据时，可以提供项目级 wrapper：

```go
func ParseWithData(filename string, input []byte, data *ParserCustomData) (any, error) {
    p := newParser(filename, input)
    p.setCustomData(data)
    return p.parse(g)
}
```

action 和 predicate 通过 `c.data` 访问它。此数据不会随 PEG 回溯自动恢复。

## Action 与错误

```peg
Integer = [0-9]+ {
    value, err := strconv.Atoi(string(c.text))
    if err != nil {
        p.addErr(err)
        return nil
    }
    return value
}
```

Go action 返回一个值，不使用旧式的 `(value, error)` 双返回值。面向用户的 no-match 消息可通过[错误处理](/guide/error-handling)中的 formatter 自定义。

## Memoization

Go 的 `-optimize-parser` 当前仍保留命名规则 memoization，因此两者可以同时使用：

```go
type Option = option

func Memoize(enabled bool) Option {
    return memoized(enabled)
}
```

memo 只适合规则结果在一次 parse 中由输入位置稳定决定的情况。依赖可变 `ParserCustomData` 的 predicate 应谨慎使用。完整选择方法见[性能参数选择](/guide/performance-options)。

多份文法共享 runtime 时，使用[多文件与入口规则](/guide/multiple-grammars)中的 Go 同包方案。
