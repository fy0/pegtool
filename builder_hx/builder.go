package builder

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/fy0/pegtool/ast"

	base "github.com/fy0/pegtool/builder"
)

const (
	codeGeneratedComment = "// Code is generated; DO NOT EDIT.\n"
	maxInlineArrayWidth  = 80
)

// generated function templates
var (
	callCodeFuncTemplate = `private function _call{{.funcName}}(p:Parser):Any {
{{if .useStack}}	var stack = p.vstack.first();
{{end}}	var c = p.cur;
	return (function(c:Current{{if .hasParams}}, {{.paramsDef}}{{end}}):Any {
		{{.code}}
		return null;
	})(p.cur{{if .hasParams}}, {{.paramsCall}}{{end}});
}
`
	callPredFuncTemplate = `private function _call{{.funcName}}(p:Parser):Bool {
{{if .useStack}}	var stack = p.vstack.first();
{{end}}	return (function(c:Current{{if .hasParams}}, {{.paramsDef}}{{end}}):Bool {
		{{.code}}
	})(p.cur{{if .hasParams}}, {{.paramsCall}}{{end}});
}
`
)

func haxeRuneArrayWidth(values []rune) int {
	width := 2
	for _, rn := range values {
		width += len(fmt.Sprintf("0x%X, ", rn))
	}
	return width
}

func haxeClassArrayWidth(classes []string) int {
	width := 2
	for _, class := range classes {
		width += len("hxUnicode.U.") + len(class) + len(", ")
	}
	return width
}

// UseHxUnicode makes generated Unicode character classes reference the
// hxUnicode haxelib instead of embedding their range data.
func UseHxUnicode(enable bool) base.Option {
	return func(b *base.Builder) base.Option {
		previous := b.HaxeUseHxUnicode
		b.HaxeUseHxUnicode = enable
		return UseHxUnicode(previous)
	}
}

func writeHaxeRuneArray(b *base.Builder, values []rune, ignoreCase bool) {
	multiline := haxeRuneArrayWidth(values) > maxInlineArrayWidth
	b.WriteArray("rune", multiline, func() {
		for _, rn := range values {
			if ignoreCase {
				rn = unicode.ToLower(rn)
			}
			if multiline {
				b.WritelnItemf("0x%X,", rn)
			} else {
				b.WriteItemf("0x%X,", rn)
			}
		}
	})
}

func writeHaxeUnicodeRanges(b *base.Builder, count int, valueAt func(int) (uint32, uint32, uint32)) {
	var line strings.Builder
	for i := range count {
		lo, hi, stride := valueAt(i)
		value := fmt.Sprintf("[%d,%d,%d],", lo, hi, stride)
		if line.Len() > 0 && line.Len()+1+len(value) > maxInlineArrayWidth {
			b.WritelnItemf("%s", line.String())
			line.Reset()
		}
		if line.Len() > 0 {
			line.WriteByte(' ')
		}
		line.WriteString(value)
	}
	if line.Len() > 0 {
		b.WritelnItemf("%s", line.String())
	}
}

func haxeUnicodeTable(name string) (*unicode.RangeTable, error) {
	if table := unicode.Categories[name]; table != nil {
		return table, nil
	}
	if table := unicode.Properties[name]; table != nil {
		return table, nil
	}
	if table := unicode.Scripts[name]; table != nil {
		return table, nil
	}
	return nil, fmt.Errorf("unknown Unicode class %q", name)
}

func writeHaxeUnicodeTable(b *base.Builder, name string) {
	table, err := haxeUnicodeTable(name)
	if err != nil {
		b.Err = err
		return
	}

	b.WritelnItemf("{R16: [")
	b.Indent()
	writeHaxeUnicodeRanges(b, len(table.R16), func(i int) (uint32, uint32, uint32) {
		item := table.R16[i]
		return uint32(item.Lo), uint32(item.Hi), uint32(item.Stride)
	})
	b.Dedent()
	b.Writelnf("], R32: [")
	b.Indent()
	writeHaxeUnicodeRanges(b, len(table.R32), func(i int) (uint32, uint32, uint32) {
		item := table.R32[i]
		return item.Lo, item.Hi, item.Stride
	})
	b.Dedent()
	b.Writelnf("], LatinOffset: %d},", table.LatinOffset)
}

func writeHaxeClassArray(b *base.Builder, classes []string) {
	if b.HaxeUseHxUnicode {
		multiline := haxeClassArrayWidth(classes) > maxInlineArrayWidth
		b.WriteArray("UnicodeRangeTable", multiline, func() {
			for _, class := range classes {
				if multiline {
					b.WritelnItemf("hxUnicode.U.%s,", class)
				} else {
					b.WriteItemf("hxUnicode.U.%s,", class)
				}
			}
		})
		return
	}

	b.WriteArray("UnicodeRangeTable", true, func() {
		for _, class := range classes {
			writeHaxeUnicodeTable(b, class)
		}
	})
}

func grammarUsesUnicodeClasses(g *ast.Grammar) bool {
	found := false
	ast.Inspect(g, func(expr ast.Expression) bool {
		if found || expr == nil {
			return false
		}
		if matcher, ok := expr.(*ast.CharClassMatcher); ok && len(matcher.UnicodeClasses) > 0 {
			found = true
			return false
		}
		return true
	})
	return found
}

func writeHxUnicodeRequirement(b *base.Builder) {
	if !b.HaxeUseHxUnicode || !b.RangeTable {
		return
	}
	b.Writeln(`#if !hxUnicode
#error "This parser uses hxUnicode. Run 'haxelib install hxUnicode' and compile with '-lib hxUnicode'."
#end`)
	b.Writeln("")
}

// BuildParser builds the PEG parser using the provider grammar. The code is
// written to the specified w.
func BuildParser(w io.Writer, g *ast.Grammar, opts ...base.Option) error {
	b0 := base.Builder{W: w, RecvName: "c", Target: "haxe", GrammarName: "g", IndentUnit: "\t"}
	b := &BuilderHx{b0}
	b.Init()
	b.SetOptions(opts)
	b.RangeTable = grammarUsesUnicodeClasses(g)
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
		if b.GrammarOnly && init == nil {
			return
		}
		b.Writef("%s", codeGeneratedComment)
		b.Writeln(`import haxe.Exception;
import haxe.ds.GenericStack;
import haxe.io.BytesBuffer;
import haxe.ds.HashMap;
import haxe.iterators.StringIteratorUnicode;`)

		if init == nil {
			b.Writeln(`
class ParserCustomData {
	public function new() {}
}`)
			b.Writeln("")
			return
		}

		// remove opening and closing braces
		val := strings.TrimSpace(b.TemplateRender(init.Val[1:len(init.Val)-1], false))
		if val != "" {
			b.Writeln("")
			b.Writeln(val)
		}
		b.Writeln("")
	}

	b.Shims.WriteGrammar = func(b *base.Builder, g *ast.Grammar) {
		writeHxUnicodeRequirement(b)
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
		writeHxUnicodeRequirement(b)
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

	b.Shims.WriteAndCodeExpr = func(b *base.Builder, expr *ast.AndCodeExpr) {
		if expr == nil {
			b.WriteNilLine()
			return
		}
		if expr.FuncIx == 0 {
			expr.FuncIx = b.ExprIndex
		}
		b.WriteExprBlock("andCodeExpr", false, func() {
			b.WriteRulePos(expr.Pos())
			b.WriteFieldf("run: _call%s,", b.Shims.FuncName(b, expr.FuncIx))
		})
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

	b.Shims.WriteCodeExpr = func(b *base.Builder, expr *ast.CodeExpr) {
		if expr == nil {
			b.WriteNilLine()
			return
		}
		if expr.FuncIx == 0 {
			expr.FuncIx = b.ExprIndex
		}
		b.WriteExprBlock("codeExpr", false, func() {
			b.WriteRulePos(expr.Pos())
			b.WriteFieldf("run: _call%s,", b.Shims.FuncName(b, expr.FuncIx))
			b.WriteFieldf("notSkip: %v,", expr.NotSkip)
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

		b.Writeln("")
		b.Writef("%s", b.TemplateRenderBase(funcTpl, false, map[string]any{
			"funcName":   b.FuncName(funcIx),
			"paramsDef":  params,
			"code":       val,
			"paramsCall": args.String(),
			"useStack":   len(argsInfo) > 0,
			"hasParams":  params != "",
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
				writeHaxeRuneArray(b, ch.Chars, ch.IgnoreCase)
			}
			if len(ch.Ranges) > 0 {
				b.WriteFieldf("ranges: ")
				writeHaxeRuneArray(b, ch.Ranges, ch.IgnoreCase)
			}
			if len(ch.UnicodeClasses) > 0 {
				b.RangeTable = true
				b.WriteFieldf("classes: ")
				writeHaxeClassArray(b, ch.UnicodeClasses)
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

	b.Shims.WriteNotCodeExpr = func(b *base.Builder, expr *ast.NotCodeExpr) {
		if expr == nil {
			b.WriteNilLine()
			return
		}
		if expr.FuncIx == 0 {
			expr.FuncIx = b.ExprIndex
		}
		b.WriteExprBlock("notCodeExpr", false, func() {
			b.WriteRulePos(expr.Pos())
			b.WriteFieldf("run: _call%s,", b.Shims.FuncName(b, expr.FuncIx))
		})
	}

	b.Shims.WriteRulePos = func(b *base.Builder, pos ast.Pos) {
		if b.SetRulePos {
			b.WritelnFieldf("pos: {line: %d, col: %d, offset: %d},", pos.Line, pos.Col, pos.Off)
		}
	}

	b.Shims.WriteStaticCodeWrap = func(b *base.Builder) {
		b.Writeln("")
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
