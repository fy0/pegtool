package builder

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/fy0/pigeon/ast"

	base "github.com/fy0/pigeon/builder"
)

const codeGeneratedComment = "// Code is generated; DO NOT EDIT.\n\n"

// generated function templates
var (
	callCodeFuncTemplate = `private function _call{{.funcName}}(p:Parser):Any {
{{ if .useStack }} var stack = p.vstack.first(); {{ end }}
	var c = p.cur;
	return (function (c:Current, {{.paramsDef}}):Any {
		{{.code}}
		return null;
	})(p.cur, {{.paramsCall}});
}
`
	callPredFuncTemplate = `private function _call{{.funcName}}(p:Parser):Bool {
{{ if .useStack }} var stack = p.vstack.first(); {{ end }}
	return (function (c:Current, {{.paramsDef}}):Bool {
		{{.code}}
	})(p.cur, {{.paramsCall}});
}
`
)

// BuildParser builds the PEG parser using the provider grammar. The code is
// written to the specified w.
func BuildParser(w io.Writer, g *ast.Grammar, opts ...base.Option) error {
	b0 := base.Builder{W: w, RecvName: "c", Target: "haxe", IndentUnit: "\t"}
	b := &BuilderHx{b0}
	b.Init()
	b.SetOptions(opts)
	return b.BuildParser(g)
}

type BuilderHx struct {
	base.Builder
}

func (b *BuilderHx) Init() {
	b.Builder.Init()
	b.CallCodeFuncTemplate = callCodeFuncTemplate
	b.CallPredFuncTemplate = callPredFuncTemplate

	b.Shims.WriteInit = func(b *base.Builder, init *ast.CodeBlock) {
		if init == nil {
			return
		}

		// remove opening and closing braces
		val := codeGeneratedComment
		b.Writef("%s", val)
		b.Writeln(`import haxe.Exception;
import haxe.ds.GenericStack;
import haxe.io.BytesBuffer;
import haxe.io.Bytes;
import haxe.ds.HashMap;
import haxe.iterators.StringIteratorUnicode;
`)
		b.Writelnf(b.TemplateRender(init.Val[1:len(init.Val)-1], false))
	}

	b.Shims.WriteGrammar = func(b *base.Builder, g *ast.Grammar) {
		// transform the ast grammar to the self-contained, no dependency version
		// of the parser-generator grammar.
		m := map[string]*base.ExprInfo{}

		for index, r := range g.Rules {
			info := b.GetExprInfo(r.Expr)
			info.Index = index
			m[r.Name.Val] = info
		}
		b.RuleName2Index = m

		b.Writelnf("var %s:Grammar = {", b.GrammarName)
		b.Indent()
		b.WritelnFieldf("rules: [")
		b.Indent()
		for _, r := range g.Rules {
			b.WriteRule(r)
		}
		b.Dedent()
		b.Writelnf("],")
		b.Dedent()
		b.Writelnf("}")
	}

	b.Shims.WriteGrammar2 = func(b *base.Builder, g *ast.Grammar) {
		// transform the ast grammar to the self-contained, no dependency version
		// of the parser-generator grammar.
		b.Writelnf("var g = {")
		b.Indent()
		for _, r := range g.Rules {
			b.WriteRule(r)
		}
		b.Dedent()
		b.Writelnf("}")
	}

	b.Shims.WriteActionExpr = func(b *base.Builder, act *ast.ActionExpr) {
		if act == nil {
			b.WriteNilLine()
			return
		}
		if act.FuncIx == 0 {
			act.FuncIx = b.ExprIndex
		}
		b.WriteExprBlock("actionExpr", true, func() {
			b.WriteRulePos(act.Pos())
			b.WritelnFieldf("run: _call%s,", b.Shims.FuncName(b, act.FuncIx))
			b.WriteFieldf("expr: ")
			b.WriteExpr(act.Expr)
		})
	}

	b.Shims.WriteFunc = func(b *base.Builder, funcIx int, code *ast.CodeBlock, funcTpl string) {
		if code == nil {
			return
		}
		val := b.TemplateRender(strings.TrimSpace(code.Val)[1:len(code.Val)-1], true)
		if len(val) > 0 && val[0] == '\n' {
			val = val[1:]
		}
		if len(val) > 0 && val[len(val)-1] == '\n' {
			val = val[:len(val)-1]
		}
		var args bytes.Buffer
		ix := len(b.ArgsStack) - 1
		argsInfo := base.StringArrayUniq(b.ArgsStack[ix])
		if ix >= 0 {
			for i, arg := range argsInfo {
				if i > 0 {
					args.WriteString(", ")
				}
				args.WriteString(arg)
			}
		}
		if args.Len() > 0 {
			args.WriteString(":Any")
		}

		params := args.String()
		args.Reset()
		if ix >= 0 {
			for i, arg := range argsInfo {
				if i > 0 {
					args.WriteString(", ")
				}
				fmt.Fprintf(&args, `stack[%q]`, arg)
			}
		}

		b.Writelnf(b.TemplateRenderBase(funcTpl, false, map[string]any{
			"funcName":   b.FuncName(funcIx),
			"paramsDef":  params,
			"code":       val,
			"paramsCall": args.String(),
			"useStack":   len(argsInfo) > 0,
		}))
	}

	b.Shims.WriteCharClassMatcher = func(b *base.Builder, ch *ast.CharClassMatcher) {
		if ch == nil {
			b.WriteNilLine()
			return
		}
		b.WriteExprBlock("charClassMatcher", true, func() {
			pos := ch.Pos()
			b.WriteRulePos(pos)
			b.WritelnFieldf("val: %q,", ch.Val)
			if len(ch.Chars) > 0 {
				b.WriteFieldf("chars: ")
				b.WriteArray("rune", false, func() {
					for _, rn := range ch.Chars {
						if ch.IgnoreCase {
							// b.Writef("%q.charCodeAt(0),", unicode.ToLower(rn))
							b.Writef("0x%X,", unicode.ToLower(rn))
						} else {
							// b.Writef("%q.charCodeAt(0),", rn)
							b.Writef("0x%X,", rn)
						}
					}
				})
			}
			if len(ch.Ranges) > 0 {
				b.WriteFieldf("ranges: ")
				b.WriteArray("rune", false, func() {
					for _, rn := range ch.Ranges {
						if ch.IgnoreCase {
							// b.Writef("%q.charCodeAt(0),", unicode.ToLower(rn))
							b.Writef("0x%X,", unicode.ToLower(rn))
						} else {
							// b.Writef("%q.charCodeAt(0),", rn)
							b.Writef("0x%X,", rn)
						}
					}
				})
			}
			if len(ch.UnicodeClasses) > 0 {
				b.RangeTable = true
				b.WriteFieldf("classes: ")
				b.WriteArray("*unicode.RangeTable", false, func() {
					for _, cl := range ch.UnicodeClasses {
						// b.Writef("RangeTable(%q),", cl)
						b.Writef("Unicode.%s,", cl)
					}
				})
			}
			if ch.IgnoreCase {
				b.WritelnFieldf("ignoreCase: %t,", ch.IgnoreCase)
			}
			if ch.Inverted {
				b.WritelnFieldf("inverted: %t,", ch.Inverted)
			}
		})
	}

	b.Shims.WriteLabeledExpr = func(b *base.Builder, lab *ast.LabeledExpr) {
		if lab == nil {
			b.WriteNilLine()
			return
		}
		b.WriteExprBlock("labeledExpr", true, func() {
			pos := lab.Pos()
			b.WriteRulePos(pos)
			b.WritelnFieldf("label: %q,", lab.Label.Val)
			b.WriteFieldf("expr: ")
			b.WriteExpr(lab.Expr)
			b.WritelnFieldf("textCapture: %v,", lab.TextCapture)
		})
	}

	b.Shims.WriteRulePos = func(b *base.Builder, pos ast.Pos) {
		if b.SetRulePos {
			b.WritelnFieldf("pos: {line: %d, col: %d, offset: %d},", pos.Line, pos.Col, pos.Off)
		}
	}

	b.Shims.WriteStaticCodeWrap = func(b *base.Builder) {
		b.WriteStaticCode(staticCode)
	}

	b.Shims.WriteNilLine = func(b *base.Builder) {
		b.Writelnf("null,")
	}

	b.Shims.WriteExprBlock = func(b *base.Builder, name string, newline bool, inside func()) {
		// &zeroOrMoreExpr{ ... }
		name = strings.ToUpper(name[:1]) + name[1:]
		b.Writef("%s({", name)
		if newline {
			b.Writeln("")
		}
		b.Indent()
		if inside != nil {
			inside()
		}
		b.Dedent()
		if !newline && inside != nil {
			b.WriteSpace()
		}
		b.Writef("}),\n")
	}

	b.Shims.WriteArray = func(b *base.Builder, typeName string, newline bool, inside func()) {
		// Example: [ ... ],
		b.Writef("[")
		if newline {
			b.Writeln("")
		}
		b.Indent()
		if inside != nil {
			inside()
		}
		b.Dedent()
		b.Writeln("],")
	}
}
