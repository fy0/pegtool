# pegtool - a PEG parser generator

[![Go Reference](https://pkg.go.dev/badge/github.com/fy0/pegtool.svg)](https://pkg.go.dev/github.com/fy0/pegtool)
[![Test Status](https://github.com/fy0/pegtool/actions/workflows/ci.yml/badge.svg)](https://github.com/fy0/pegtool/actions/workflows/ci.yml)
[![GoReportCard](https://goreportcard.com/badge/github.com/fy0/pegtool)](https://goreportcard.com/report/github.com/fy0/pegtool)
[![Software License](https://img.shields.io/badge/license-BSD-blue.svg)](LICENSE)

The pegtool command generates Go, Haxe, TypeScript, C#, and Rust parsers based on a [parsing expression grammar (PEG)][0]. Its grammar and syntax is inspired by the [PEG.js project][1], while the implementation is loosely based on the [parsing expression grammar for C# 3.0][2] article. It parses Unicode text encoded in UTF-8.

## Features for this fork

* Performance improvements (5-10x faster than the original version):
  * Removed `parser.state`, because it's very slow.
  * `Memoized` can be worked together with `Optimize` option.
  * More memory efficient - reduces memory allocations
  * Generates cleaner and less lines of code

* `parseSeqExpr` does not construct a value slice; repetitions collect only non-nil values explicitly returned by nested actions:
  * Example 1: `e <- items:(__ Integer __)+ EOF { return items }` returns `nil` since no values are explicitly returned
  * Example 2: `e <- items:(__ i:Integer __ { return i })+ EOF { return items }` returns `[1, 2, 3]` because the integers are explicitly returned
  * This is more efficient than the original version which would return `[[nil, 1, nil], [nil, 2, nil], [nil, 3, nil]]` with unnecessary nil values

* String capture:
  * `expr <- val:<anotherExpr> { fmt.Println(val.(string)) }` // you got the string value in `val`
  * `expr <- val:<(A '=' B)> { fmt.Println(val.(string)) }` // you got the string value in `val`, e.g. "A=B"

* Use code to control matching behavior (`andCodeExpr` and `notCodeExpr`):
    * `expr <- &{ return c.data.AllowNumber } [0-9]+` // Only matches digits if c.data.AllowNumber is true
    * `expr <- val:<[0-9]+> &{ return val.(string) == "123" } { return val.(string) }` // Only succeeds if the matched string equals "123"
    * Tips: don't enable memoized if you control matching behavior by code.

* Logical `and` / `or` match:
  * `expr <- &&testExpr testExpr` // if testExpr return ok but matched nothing (e.g. testExpr <- 'A'*), `&&testEpr` returns false.

* Multiple peg files supported:
  1. `pegtool -o script1.peg.go script1.peg` to generate a normal parser.
  2. Run `pegtool -grammar-only -grammar-name=g2 -run-func-prefix="_s2_" -o script2.peg.go script2.peg` to generate grammar only code in same package.
  3. Use it by `newParser("filename", "expr").parse(g2)`

* Simplified `actionExpr`:
  * The original version required two parameters to return (val, error), but errors are rarely used. So this fork simplifies the return values.
  * Examples:
    * `expr <- [0-9]+ { fmt.Println(expr) }` is ok in this fork, returns nothing.
    * `expr <- "true" { return 1 }` if you want return something.
  * Add an error by manual:
    * `expr <- "if" { p.addErr(errors.New("keyword is not allowed")) }`, equals to `expr <- "if" { return nil, errors.New("keyword is not allowed") }` of original Pigeon.

* Provide a struct(`ParserCustomData`) to embed, to replace the `globalStore`
  * Must define a struct `ParserCustomData` in your module.
  * Access data by `c.data`, for example: `expr <- { fmt.Println(c.data.MyOption) }`
  * `globalState` is removed.

* Remove ParseFile ParseReader, rename Parse and all options to lowercase [issue](https://github.com/mna/pigeon/issues/150), branch feat/rename-exported-api
  * `ParseReader` converts io.Reader to bytes, then invoke `parse`, it don't make sense.
  * Function `Parse` and all options(`MaxExpressions`,`Entrypoint`,`Statistics`,`Debug`,`Memoize`,`AllowInvalidUTF8`,`Recover`,`GlobalStore`,`InitState`) expose to module user. I think expose them is not a good idea.

* Skip "actionExpr" while looking ahead [issue](https://github.com/mna/pigeon/issues/149), branch feat/skip-code-expr-while-looking-ahead
  * See detail in the issue.
  * `*{}` / `&{}` / `!{}` won't skip.

* ActionExpr refactored [issue](https://github.com/mna/pigeon/issues/150), branch refactor/actionExpr
  * Unlimited ActionExpr(CodeExpr): grammar like `expr <- firstPart:[0-9]+ { fmt.Println(firstPart) }  secondPart:[a-z]+ { fmt.Println(firstPart, secondPart) }` is allowed for this fork.
  * You can access parser in ActionExpr: `expr <- { fmt.Println(p) }`
  * `stateCodeExpr(#{})` was removed.

* `position` of generated code is removed 
  * It produced a lot of different for version control.
  * You can keep it by set `SetRulePos` to true and rebuild.

* Added `-optimize-ref-expr-by-index` option
  * An option to tweak `RefExpr` the most usually used expr in parser.
  * About ~10% faster with this option.

* Removed `-support-left-recursion` option
  * It's not used much, so I removed it to make maintenance easier

* Removed `-optimize-grammar` option
  * There are bugs present and the effects are not significant.

* Removed `-optimize-basic-latin` option
  * Because there is no evidence to suggest that this is an optimization

* `charClassMatcher` / `anyMatcher` / `litMatcher` not return byte anymore, because of performance.
  * Use string capture or `c.text` instead.

## Installation

```
go install github.com/fy0/pegtool@latest
```

This will install or update the package, and the `pegtool` command will be installed in your $GOBIN directory. Neither this package nor the parsers generated by this command require any third-party dependency, unless such a dependency is used in the code blocks of the grammar.

## Basic usage

```
pegtool [options] [PEG_GRAMMAR_FILE]
```

By default, the input grammar is read from `stdin` and the generated code is printed to `stdout`. You may save it in a file using the `-o` flag.

## Language targets

Use `-t` or `-target` to select a generator. Every generated parser is self-contained and uses only its language's standard library.

```sh
pegtool -t go   -o parser.go grammar.peg
pegtool -t hx   -o Parser.hx grammar.peg
pegtool -t ts   -o parser.ts grammar.peg
pegtool -t cs   -o Parser.cs grammar.peg
pegtool -t c    -o parser.c grammar.peg
pegtool -t rust -o parser.rs grammar.peg
```

The aliases `typescript`, `csharp`, `c#`, `haxe`, and `c99` are also accepted. Grammar templates receive the canonical target values `go`, `haxe`, `ts`, `cs`, `c`, and `rust`, so one grammar can contain target-specific initializer and action code:

```text
{{ if eq .Target "ts" }}
// TypeScript code
{{ else if eq .Target "cs" }}
// C# code
{{ else if eq .Target "c" }}
/* C99 code */
{{ else if eq .Target "rust" }}
// Rust code
{{ end }}
```

The generated TypeScript module exports `parse`, `Parser`, and `ParserOptions`. The C# source exposes `PTParser.Parse` and `PTParser.ParserOptions`. Rust exposes `PTParser::parse`, `PTParser::parse_with_options`, and the cloneable `Value` enum used by action labels and return values.

The C target emits portable C99 and exposes `pegtool_parse`, `pegtool_parse_n`, and their `PegtoolParserOptions` variants. Parse calls return an owned `PegtoolResult`; release both successful values and errors with `pegtool_result_free`. C action labels are `const PegtoolValue *`; labels that are C99 keywords or collide with the internal `p` parameter or current-match variable use `_pegtool_label_<name>` in action code. The current match is available as `PegtoolCurrent *c` by default and can be renamed with `-receiver-name`. Actions return an owned `PegtoolValue` built with functions such as `pegtool_value_string`, `pegtool_value_integer`, or `pegtool_value_clone`. Semantic predicates return an `int`. The generated source uses only the C standard library and can be compiled with `cc -std=c99 parser.c`.

Generated Haxe parsers embed only the Unicode tables used by the grammar and have no Haxe library dependency by default. Pass `-haxe-use-hxunicode` to keep parsers with Unicode character classes smaller; in that mode, install the dependency with `haxelib install hxUnicode` and compile with `-lib hxUnicode`. The generated source reports these commands directly if the library is missing.

TypeScript and Rust initializer blocks are emitted at module scope. A C# initializer block is emitted inside the generated partial `PTParser` class, so it should contain nested types, fields, and helper methods rather than namespace or top-level declarations.

TypeScript and C# action labels use `any` and `object?` respectively. C# labels that are language keywords are emitted as verbatim identifiers, for example the PEG label `params` is referenced as `@params` in C# action code. Rust labels use `Value`; its `String`, `Int`, `Float`, `Bool`, `Array`, `Object`, and `Custom` variants cover parser values and user-defined types.

Github user [@mna][6] created the original package in April 2015, and [@breml][5] is the original package's maintainer as of May 2017.

## Example

Given the following grammar:

```
{
//nolint:unreachable
package main

type ParserCustomData struct {
}

var ops = map[string]func(int, int) int {
    "+": func(l, r int) int {
        return l + r
    },
    "-": func(l, r int) int {
        return l - r
    },
    "*": func(l, r int) int {
        return l * r
    },
    "/": func(l, r int) int {
        return l / r
    },
}

func toAnySlice(v any) []any {
    if v == nil {
        return nil
    }
    return v.([]any)
}

func eval(first, rest any) int {
    l := first.(int)
    restSl := toAnySlice(rest)
    for _, v := range restSl {
        restExpr := toAnySlice(v)
        r := restExpr[1].(int)
        op := restExpr[0].(string)
        l = ops[op](l, r)
    }
    return l
}
}

Input <- expr:Expr EOF {
    return expr
}

Expr <- _ first:Term rest:( _ op:AddOp _ r:Term { return []any{op, r} })* _ {
    return eval(first, rest)
}

Term <- first:Factor rest:( _ op:MulOp _ r:Factor { return []any{op, r} })* {
    return eval(first, rest)
}

Factor <- '(' expr:Expr ')' {
    return expr
} / integer:Integer {
    return integer
}

AddOp <- ( '+' / '-' ) {
    return string(c.text)
}

MulOp <- ( '*' / '/' ) {
    return string(c.text)
}

Integer <- '-'? [0-9]+ {
    v, err := strconv.Atoi(string(c.text))
    if err != nil {
        p.addErr(err)
    }
    return v
}

_ "whitespace" <- [ \n\t\r]*

EOF <- !.
```

The generated parser can parse simple arithmetic operations, e.g.:

```
18 + 3 - 27 * (-18 / -3)

=> -141
```

More examples, grouped by generated language, can be found in the
[`examples/`](examples/README.md) subdirectory.

See the [package documentation][3] for detailed usage.

## Contributing

See the CONTRIBUTING.md file.

## License

The [BSD 3-Clause license][4]. See the LICENSE file.

[0]: http://en.wikipedia.org/wiki/Parsing_expression_grammar
[1]: http://pegjs.org/
[2]: http://www.codeproject.com/Articles/29713/Parsing-Expression-Grammar-Support-for-C-Part
[3]: https://pkg.go.dev/github.com/fy0/pegtool
[4]: http://opensource.org/licenses/BSD-3-Clause
[5]: https://github.com/breml
[6]: https://github.com/mna


## TODO
* ~~performance: Create another version of `parseOneOrMoreExpr/parseZeroOrMoreExpr` which not collect results. Choose expr decide by is labeled, A bit faster.~~
* ~~performance: Remove `pushV` and `popV`, a bit faster.~~
* ~~performance: In `parseCharClassMatcher`, variable `start` can be removed in most case. Lot of of small memory pieces allocated.~~
* ~~performance: Remove Wrap function if they are not needed.~~
* performance: Too many any, can we remove `parseExpr`?
* string capture inside predicate expr not work: `&( alist:<("a")*> &{ fmt.Println(alist) } )`
* auto remove `return nil` if unreachable
