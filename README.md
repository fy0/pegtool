# pigeon - a PEG parser generator for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/mna/pigeon.svg)](https://pkg.go.dev/github.com/fy0/pigeon)
[![Test Status](https://github.com/mna/pigeon/workflows/Go%20Matrix/badge.svg)](https://github.com/fy0/pigeon/actions?query=workflow%3AGo%20Matrix)
[![GoReportCard](https://goreportcard.com/badge/github.com/mna/pigeon)](https://goreportcard.com/report/github.com/fy0/pigeon)
[![Software License](https://img.shields.io/badge/license-BSD-blue.svg)](LICENSE)

The pigeon command generates parsers based on a [parsing expression grammar (PEG)][0]. Its grammar and syntax is inspired by the [PEG.js project][1], while the implementation is loosely based on the [parsing expression grammar for C# 3.0][2] article. It parses Unicode text encoded in UTF-8.

## Features for this fork

* Performance improvements (5-10x faster than the original version):
  * Removed `parser.state`, because it's very slow.
  * `Memoized` can be worked together with `Optimize` option.
  * More memory efficient - reduces memory allocations
  * Generates cleaner and less lines of code

* `parseSeqExpr` only collects values that are explicitly returned, improving performance:
  * Example 1: `e <- items:(__ Integer __)+ EOF { return items }` returns an empty array `[]` since no values are explicitly marked for collection
  * Example 2: `e <- items:(__ i:Integer __ { return i })+ EOF { return items }` returns `[1, 2, 3]` because integers are explicitly collected
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
  1. `pigeon -o script1.peg.go script1.peg` to generate a normal parser.
  2. Run `pigeon -grammar-only -grammar-name=g2 -run-func-prefix="_s2_" -o script2.peg.go script2.peg` to generate grammar only code in same package.
  3. Use it by `newParser("filename", "expr").parse(g2)`

* Simplified `actionExpr`:
  * The original version required two parameters to return (val, error), but errors are rarely used. So this fork simplifies the return values.
  * Examples:
    * `expr <- [0-9]+ { fmt.Println(expr) }` is ok in this fork, returns nothing.
    * `expr <- "true" { return 1 }` if you want return something.
  * Add an error by manual:
    * `expr <- "if" { p.addErr(errors.New("keyword is not allowed")) }`, equals to `expr <- "if" { return nil, errors.New("keyword is not allowed") }` of original pigeon.

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
go install github.com/fy0/pigeon@latest
```

This will install or update the package, and the `pigeon` command will be installed in your $GOBIN directory. Neither this package nor the parsers generated by this command require any third-party dependency, unless such a dependency is used in the code blocks of the grammar.

## Basic usage

```
pigeon [options] [PEG_GRAMMAR_FILE]
```

By default, the input grammar is read from `stdin` and the generated code is printed to `stdout`. You may save it in a file using the `-o` flag.

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

More examples can be found in the `examples/` subdirectory.

See the [package documentation][3] for detailed usage.

## Contributing

See the CONTRIBUTING.md file.

## License

The [BSD 3-Clause license][4]. See the LICENSE file.

[0]: http://en.wikipedia.org/wiki/Parsing_expression_grammar
[1]: http://pegjs.org/
[2]: http://www.codeproject.com/Articles/29713/Parsing-Expression-Grammar-Support-for-C-Part
[3]: https://pkg.go.dev/github.com/fy0/pigeon
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
