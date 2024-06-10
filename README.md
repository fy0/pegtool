# pigeon - a PEG parser generator for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/mna/pigeon.svg)](https://pkg.go.dev/github.com/mna/pigeon)
[![Test Status](https://github.com/mna/pigeon/workflows/Go%20Matrix/badge.svg)](https://github.com/mna/pigeon/actions?query=workflow%3AGo%20Matrix)
[![GoReportCard](https://goreportcard.com/badge/github.com/mna/pigeon)](https://goreportcard.com/report/github.com/mna/pigeon)
[![Software License](https://img.shields.io/badge/license-BSD-blue.svg)](LICENSE)

The pigeon command generates parsers based on a [parsing expression grammar (PEG)][0]. Its grammar and syntax is inspired by the [PEG.js project][1], while the implementation is loosely based on the [parsing expression grammar for C# 3.0][2] article. It parses Unicode text encoded in UTF-8.

See the [godoc page][3] for detailed usage. Also have a look at the [Pigeon Wiki](https://github.com/mna/pigeon/wiki) for additional information about Pigeon and PEG in general.

## Features for this fork

* Performance tweak, 5-10x faster than original version(with some trade off).
  * `parser.state` is removed, because it's very slow.
  * `parseSeqExpr` only collect not nil values now. Mainly for performance improvement. For example: `e <- Expr __ Plus __ Expr` returns \[expr, '+', expr], original version return \[expr, nil, '+', nil, expr].
  * Memoized can be worked with Optimize and enable by default.
  * Generated parser do less memory allocated.
  * Generated parser uses fewer lines of code.

* Multiple peg files supported.
  1. `pigeon -o script1.peg.go script1.peg` to generate a normal parser.
  2. Run `pigeon -grammar-only -grammar-name=g2 -run-func-prefix="_s2_" -o script2.peg.go script2.peg` to generate grammar only code in same package.
  3. Use it by `newParser("filename", "expr").parse(g2)`


* `actionExpr` is different
  * Only needs to return one value. Moreover, this is not required. If the return statement is not written, it will automatically return nil. Examples:
    * `expr <- [0-9]+ { fmt.Println(expr) }` is ok in this fork.
    * `expr <- "true" { return 1 }` if you want return something.
  * If you want to add an error by manual, do this:
    * `expr <- "if" { p.addErr(errors.New("keyword is not allowed")) }`, equals to `expr <- "if" { return nil, errors.New("keyword is not allowed") }` of original pigeon.

* `andCodeExpr` and `notCodeExpr`:
    * Like `actionExpr`, return a bool instead of return bool and error
    * `expr <- &{ return c.data.AllowNumber } [0-9]+`

* String capture:
  * `expr <- val:<anotherExpr> { fmt.Println(val.(string)) }`
  * `expr <- val:<(A '=' B)> { fmt.Println(val.(string)) }`

* Logical `and` / `or` match:
  * `expr <- &&testExpr testExpr` // if testExpr return ok but matched nothing (e.g. testExpr <- 'A'*), `&&testEpr` returns false.

* Skip "actionExpr" while looking ahead [issue](https://github.com/mna/pigeon/issues/149), branch feat/skip-code-expr-while-looking-ahead
  * See detail in the issue.
  * `*{}` / `&{}` / `!{}` won't skip.

* Remove ParseFile ParseReader, rename Parse and all options to lowercase [issue](https://github.com/mna/pigeon/issues/150), branch feat/rename-exported-api
  * `ParseReader` converts io.Reader to bytes, then invoke `parse`, it don't make sense.
  * Function `Parse` and all options(`MaxExpressions`,`Entrypoint`,`Statistics`,`Debug`,`Memoize`,`AllowInvalidUTF8`,`Recover`,`GlobalStore`,`InitState`) expose to module user. I think expose them is not a good idea.

* ActionExpr refactored [issue](https://github.com/mna/pigeon/issues/150), branch refactor/actionExpr
  * Unlimited ActionExpr(CodeExpr): grammar like `expr <- firstPart:[0-9]+ { fmt.Println(firstPart) }  secondPart:[a-z]+ { fmt.Println(firstPart, secondPart) }` is allowed for this fork.
  * You can access parser in ActionExpr: `expr <- { fmt.Println(p) }`
  * `stateCodeExpr(#{})` was removed.

* Provide a struct(`ParserCustomData`) to embed, to replace the globalStore
  * Must define a struct `ParserCustomData` in your module.
  * Access data by `c.data`, for example: `expr <- { fmt.Println(c.data.MyOption) }`
  * `globalState` is removed.

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

## Releases

* v1.0.0 is the tagged release of the original implementation.
* Work has started on v2.0.0 with some planned breaking changes.

Github user [@mna][6] created the package in April 2015, and [@breml][5] is the package's maintainer as of May 2017.

### Release policy

Starting of June 2023, the backwards compatibility support for `pigeon` is changed to follow the official [Go Security Policy](https://github.com/golang/go/security/policy).

Over time, the Go ecosystem is evolving.

On one hand side, packages like [golang.org/x/tools](https://pkg.go.dev/golang.org/x/tools), which are critical dependencies of `pigeon`, do follow the official Security Policy and with `pigeon` not following the same guidelines, it was no longer possible to include recent versions of these dependencies and with this it was no longer possible to include critical bugfixes.
On the other hand there are changes to what is considered good practice by the greater community (e.g. change from `interface{}` to `any`). For users following (or even enforcing) these good practices, the code generated by `pigeon` does no longer meet the bar of expectations.
Last but not least, following the Go Security Policy over the last years has been a smooth experience and therefore updating Go on a regular bases feels like duty that is reasonable to be put on users of `pigeon`.

This observations lead to the decision to follow the same Security Policy as Go.

## Installation

Provided you have Go correctly installed with the $GOPATH and $GOBIN environment variables set, run:

```
$ go get -u github.com/mna/pigeon
```

This will install or update the package, and the `pigeon` command will be installed in your $GOBIN directory. Neither this package nor the parsers generated by this command require any third-party dependency, unless such a dependency is used in the code blocks of the grammar.

## Basic usage

```
$ pigeon [options] [PEG_GRAMMAR_FILE]
```

By default, the input grammar is read from `stdin` and the generated code is printed to `stdout`. You may save it in a file using the `-o` flag.

## Example

Given the following grammar:

```
{
// part of the initializer code block omitted for brevity
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
        r := restExpr[3].(int)
        op := restExpr[1].(string)
        l = ops[op](l, r)
    }
    return l
}
}


Input <- expr:Expr EOF {
    return expr
}

Expr <- _ first:Term rest:( _ AddOp _ Term )* _ {
    return eval(first, rest)
}

Term <- first:Factor rest:( _ MulOp _ Factor )* {
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
[3]: https://pkg.go.dev/github.com/mna/pigeon
[4]: http://opensource.org/licenses/BSD-3-Clause
[5]: https://github.com/breml
[6]: https://github.com/mna


## TODO
* ~~performance: Create another version of `parseOneOrMoreExpr/parseZeroOrMoreExpr` which not collect results. Choose expr decide by is labeled, A bit faster.~~
* ~~performance: Remove `pushV` and `popV`, a bit faster.~~
* ~~performance: In `parseCharClassMatcher`, variable `start` can be removed in most case. Lot of of small memory pieces allocated.~~
* ~~performance: Remove Wrap function if they are not needed.~~
* performance: Too many any, can we remove `parseExpr`?
