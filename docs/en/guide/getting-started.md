# Quick Start

Pegtool is written in Go. The examples below use `pegtool` as the command name.

## Installation

Install the latest version from the Go module:

```powershell
go install github.com/fy0/pigeon@latest
pegtool -h
```

Or build it from the current repository:

```powershell
go build -o .\bin\pegtool.exe .
```

## Create your first grammar

::: tip Using another target language
This section uses Go. For a minimal working parser in another language, go directly to [Haxe](/en/targets/haxe), [TypeScript](/en/targets/typescript), [C#](/en/targets/csharp), [C99](/en/targets/c), or [Rust](/en/targets/rust).

If you do not plan to run the Go example, skip to [Switch target languages](#switch-target-languages) or [Next steps](#next-steps).
:::

Create `identifier.peg`. This Go grammar recognizes a Unicode identifier and returns the captured text:

```peg
{
package main

type ParserCustomData struct{}

func ParseIdentifier(filename string, input []byte) (string, error) {
    value, err := parse(filename, input)
    if err != nil {
        return "", err
    }
    return value.(string), nil
}
}

Start = value:<Identifier> !. {
    return value
}

Identifier = [\pL_] [\pL\p{Nd}_]*
```

The important details are:

- The Go code in the initializer declares the package for the generated file and a public entry function.
- The Go target requires an empty `ParserCustomData` declaration even when no custom state is used.
- The first rule is the default entry point.
- `<Identifier>` explicitly captures the matched text.
- `!.` requires the complete input to be consumed instead of accepting only a valid prefix.
- `\pL` and `\p{Nd}` use Unicode classes known to the generator.

## Generate a Go parser

```powershell
go mod init example.com/identifier
pegtool -t go -o parser.go identifier.peg
```

Create `main.go` to call the generated parser:

```go
package main

import "fmt"

func main() {
    value, err := ParseIdentifier("input.txt", []byte("name_世界"))
    if err != nil {
        panic(err)
    }
    fmt.Println(value)
}
```

Build and run it:

```powershell
go run .
# Output: name_世界
```

## Switch target languages

Use `-t` or `-target` to choose the output:

| Target | Value | Typical output name |
| --- | --- | --- |
| Go | `go` | `parser.go` |
| Haxe | `hx`, `haxe` | `Parser.hx` |
| TypeScript | `ts`, `typescript` | `parser.ts` |
| C# | `cs`, `csharp`, `c#` | `Parser.cs` |
| C99 | `c`, `c99` | `parser.c` |
| Rust | `rust` | `parser.rs` |

Action code in a grammar belongs to the target language. To generate multiple languages from one file, use [target-conditional templates](/en/targets/#distinguishing-targets-in-one-grammar).

## Next steps

1. Read [PEG Syntax](/en/guide/grammar) to understand ordered choice and backtracking.
2. Read [Actions and Values](/en/guide/actions-and-values) to avoid relying on upstream Pigeon's sequence return values.
3. Benchmark with real input before choosing [performance options](/en/guide/performance-options).
