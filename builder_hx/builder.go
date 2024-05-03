// Package builder generates the parser code for a given grammar. It makes
// no attempt to verify the correctness of the grammar.
package builder

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"unicode"

	"github.com/mna/pigeon/ast"
)

const codeGeneratedComment = "// Code is generated; DO NOT EDIT.\n\n"

func (b *builder) templateRender(text string, trim bool) string {
	t, err := template.New("").Parse(text)

	// 创建一个 bytes.Buffer 来保存渲染后的结果
	var out bytes.Buffer
	// 渲染模板，并将结果写入 tout

	err = t.Execute(&out, map[string]any{
		"target": b.target,
	})

	if err != nil {
		panic(err)
	}

	if trim {
		return strings.TrimSpace(out.String())
	} else {
		return out.String()
	}
}

// generated function templates
var (
	// any, err
	onFuncTemplate = `private function _%s(%s:Current, %s):RetValErr {
%s
}
`
	onPredFuncTemplate = `func (%s *current) %s(%s) (bool, error) {
%s
}
`
	onStateFuncTemplate = `private function _%s(%s:Current,%s):Exception {
%s
}
`
	onStateFuncTemplate2 = `private function _%s(p:Parser, %s:Current, %s):RetValOkErr {
%s
}
`
	callFuncTemplate = `private function _call%s(p:Parser):RetValErr {
	var stack = p.vstack.first();
	var _ = stack;
	return _%[1]s(p.cur, %s);
}
`
	callPredFuncTemplate = `func (p *parser) call%s() (bool, error) {
	stack := p.vstack[len(p.vstack)-1]
	var _ = stack
	return p.cur.%[1]s(%s)
}
`
	callStateFuncTemplate = `private function _call%s(p:Parser):Exception {
	var stack = p.vstack.first();
    var _ = stack;
    return _%[1]s(p.cur, %s);
}
`
	callStateFuncTemplate2 = `private function _call%s(p:Parser):RetValOkErr {
	var stack = p.vstack.first();
    var _ = stack;
    return _%[1]s(p, p.cur, %s);
}
`
)

// Option is a function that can set an option on the builder. It returns
// the previous setting as an Option.
type Option func(*builder) Option

// ReceiverName returns an option that specifies the receiver name to
// use for the current struct (which is the struct on which all code blocks
// except the initializer are generated).
func ReceiverName(nm string) Option {
	return func(b *builder) Option {
		prev := b.recvName
		b.recvName = nm
		return ReceiverName(prev)
	}
}

// Optimize returns an option that specifies the optimize option
// If optimize is true, the Debug and Memoize code is completely
// removed from the resulting parser
func Optimize(optimize bool) Option {
	return func(b *builder) Option {
		prev := b.optimize
		b.optimize = optimize
		return Optimize(prev)
	}
}

// SupportLeftRecursion returns an option that specifies the supportLeftRecursion option.
// If supportLeftRecursion is true, LeftRecursion code is added to the resulting parser.
func SupportLeftRecursion(support bool) Option {
	return func(b *builder) Option {
		prev := b.supportLeftRecursion
		b.supportLeftRecursion = support
		return SupportLeftRecursion(prev)
	}
}

// Nolint returns an option that specifies the nolint option
// If nolint is true, special '// nolint: ...' comments are added
// to the generated parser to suppress warnings by gometalinter or golangci-lint.
func Nolint(nolint bool) Option {
	return func(b *builder) Option {
		prev := b.nolint
		b.nolint = nolint
		return Optimize(prev)
	}
}

// BasicLatinLookupTable returns an option that specifies the basicLatinLookup option
// If basicLatinLookup is true, a lookup slice for the first 128 chars of
// the Unicode table (Basic Latin) is generated for each CharClassMatcher
// to increase the character matching.
func BasicLatinLookupTable(basicLatinLookupTable bool) Option {
	return func(b *builder) Option {
		prev := b.basicLatinLookupTable
		b.basicLatinLookupTable = basicLatinLookupTable
		return BasicLatinLookupTable(prev)
	}
}

// BuildParser builds the PEG parser using the provider grammar. The code is
// written to the specified w.
func BuildParser(w io.Writer, g *ast.Grammar, opts ...Option) error {
	b := &builder{w: w, recvName: "c", target: "hx"}
	b.setOptions(opts)
	return b.buildParser(g)
}

type builder struct {
	w   io.Writer
	err error

	// options
	recvName              string
	optimize              bool
	basicLatinLookupTable bool
	globalState           bool
	nolint                bool
	supportLeftRecursion  bool
	haveLeftRecursion     bool

	ruleName  string
	exprIndex int
	argsStack [][]string

	target     string
	rangeTable bool
}

func (b *builder) setOptions(opts []Option) {
	for _, opt := range opts {
		opt(b)
	}
}

func (b *builder) buildParser(grammar *ast.Grammar) error {
	haveLeftRecursion, err := PrepareGrammar(grammar)
	if err != nil {
		return fmt.Errorf("incorrect grammar: %w", err)
	}
	if !b.supportLeftRecursion && haveLeftRecursion {
		return fmt.Errorf("incorrect grammar: %w", ErrHaveLeftRecursion)
	}
	b.haveLeftRecursion = haveLeftRecursion

	b.writeInit(grammar.Init)
	b.writeGrammar(grammar)
	for _, rule := range grammar.Rules {
		b.writeRuleCode(rule)
	}
	b.writeStaticCode()

	return b.err
}

func (b *builder) writeInit(init *ast.CodeBlock) {
	if init == nil {
		return
	}

	// remove opening and closing braces
	val := codeGeneratedComment
	b.writef("%s", val)
	b.writeln(`import haxe.Exception;
import haxe.ds.GenericStack;
import haxe.io.BytesBuffer;
import haxe.io.Bytes;
import haxe.ds.HashMap;
import haxe.iterators.StringIteratorUnicode;
`)
	b.writelnf(b.templateRender(init.Val[1:len(init.Val)-1], false))
}

func (b *builder) writeGrammar(g *ast.Grammar) {
	// transform the ast grammar to the self-contained, no dependency version
	// of the parser-generator grammar.
	b.writelnf("var g:Grammar = {")
	b.writelnf("\trules: [")
	for _, r := range g.Rules {
		b.writeRule(r)
	}
	b.writelnf("\t],")
	b.writelnf("}")
}

func (b *builder) writeRule(r *ast.Rule) {
	if r == nil || r.Name == nil {
		return
	}

	b.exprIndex = 0
	b.ruleName = r.Name.Val

	b.writelnf("{")
	b.writelnf("\tname: %q,", r.Name.Val)
	if r.DisplayName != nil && r.DisplayName.Val != "" {
		b.writelnf("\tdisplayName: %q,", r.DisplayName.Val)
	}
	pos := r.Pos()
	b.writelnf("\tpos: {line: %d, col: %d, offset: %d},", pos.Line, pos.Col, pos.Off)
	b.writef("\texpr: ")
	b.writeExpr(r.Expr)
	if b.haveLeftRecursion {
		b.writelnf("\tleader: %t,", r.Leader)
		b.writelnf("\tleftRecursive: %t,", r.LeftRecursive)
	}
	b.writelnf("},")
}

func (b *builder) writeExpr(expr ast.Expression) {
	b.exprIndex++
	switch expr := expr.(type) {
	case *ast.ActionExpr:
		b.writeActionExpr(expr)
	case *ast.AndCodeExpr:
		b.writeAndCodeExpr(expr)
	case *ast.AndExpr:
		b.writeAndExpr(expr)
	case *ast.AnyMatcher:
		b.writeAnyMatcher(expr)
	case *ast.CharClassMatcher:
		b.writeCharClassMatcher(expr)
	case *ast.ChoiceExpr:
		b.writeChoiceExpr(expr)
	case *ast.LabeledExpr:
		b.writeLabeledExpr(expr)
	case *ast.LitMatcher:
		b.writeLitMatcher(expr)
	case *ast.NotCodeExpr:
		b.writeNotCodeExpr(expr)
	case *ast.NotExpr:
		b.writeNotExpr(expr)
	case *ast.OneOrMoreExpr:
		b.writeOneOrMoreExpr(expr)
	case *ast.RecoveryExpr:
		b.writeRecoveryExpr(expr)
	case *ast.RuleRefExpr:
		b.writeRuleRefExpr(expr)
	case *ast.SeqExpr:
		b.writeSeqExpr(expr)
	case *ast.StateCodeExpr:
		b.writeStateCodeExpr(expr)
	case *ast.CustomParserCodeExpr:
		b.writeCustomParserCodeExpr(expr)
	case *ast.ThrowExpr:
		b.writeThrowExpr(expr)
	case *ast.ZeroOrMoreExpr:
		b.writeZeroOrMoreExpr(expr)
	case *ast.ZeroOrOneExpr:
		b.writeZeroOrOneExpr(expr)
	default:
		b.err = fmt.Errorf("builder: unknown expression type %T", expr)
	}
}

func (b *builder) writeActionExpr(act *ast.ActionExpr) {
	if act == nil {
		b.writelnf("nil,")
		return
	}
	if act.FuncIx == 0 {
		act.FuncIx = b.exprIndex
	}
	b.writelnf("ActionExpr({")
	pos := act.Pos()
	b.writelnf("\tpos: {line: %d, col: %d, offset: %d},", pos.Line, pos.Col, pos.Off)
	b.writelnf("\trun: _call%s,", b.funcName(act.FuncIx))
	b.writef("\texpr: ")
	b.writeExpr(act.Expr)
	b.writelnf("}),")
}

func (b *builder) writeAndCodeExpr(and *ast.AndCodeExpr) {
	if and == nil {
		b.writelnf("nil,")
		return
	}
	b.writelnf("AndCodeExpr({")
	pos := and.Pos()
	if and.FuncIx == 0 {
		and.FuncIx = b.exprIndex
	}
	b.writelnf("\tpos: {line: %d, col: %d, offset: %d},", pos.Line, pos.Col, pos.Off)
	b.writelnf("\trun: _call%s,", b.funcName(and.FuncIx))
	b.writelnf("}),")
}

func (b *builder) writeAndExpr(and *ast.AndExpr) {
	if and == nil {
		b.writelnf("nil,")
		return
	}
	b.writelnf("AndExpr({")
	pos := and.Pos()
	b.writelnf("\tpos: {line: %d, col: %d, offset: %d},", pos.Line, pos.Col, pos.Off)
	b.writef("\texpr: ")
	b.writeExpr(and.Expr)
	b.writelnf("}),")
}

func (b *builder) writeAnyMatcher(any *ast.AnyMatcher) {
	if any == nil {
		b.writelnf("nil,")
		return
	}
	b.writelnf("AnyMatcher({")
	pos := any.Pos()
	b.writelnf("\tline: %d, col: %d, offset: %d,", pos.Line, pos.Col, pos.Off)
	b.writelnf("}),")
}

func (b *builder) writeCharClassMatcher(ch *ast.CharClassMatcher) {
	if ch == nil {
		b.writelnf("nil,")
		return
	}
	b.writelnf("CharClassMatcher({")
	pos := ch.Pos()
	b.writelnf("\tpos: {line: %d, col: %d, offset: %d},", pos.Line, pos.Col, pos.Off)
	b.writelnf("\tval: %q,", ch.Val)
	if len(ch.Chars) > 0 {
		b.writef("\tchars: [")
		for _, rn := range ch.Chars {
			if ch.IgnoreCase {
				b.writef("%q.charCodeAt(0),", unicode.ToLower(rn))
			} else {
				b.writef("%q.charCodeAt(0),", rn)
			}
		}
		b.writelnf("],")
	} else {
		b.writef("\tchars: [],")
	}
	if len(ch.Ranges) > 0 {
		b.writef("\tranges: [")
		for _, rn := range ch.Ranges {
			if ch.IgnoreCase {
				b.writef("%q.charCodeAt(0),", unicode.ToLower(rn))
			} else {
				b.writef("%q.charCodeAt(0),", rn)
			}
		}
		b.writelnf("],")
	} else {
		b.writef("\tranges: [],")
	}
	if len(ch.UnicodeClasses) > 0 {
		b.rangeTable = true
		b.writef("\tclasses: [")
		for _, cl := range ch.UnicodeClasses {
			b.writef("Unicode.%s,", cl)
		}
		b.writelnf("],")
	}
	if b.basicLatinLookupTable {
		b.writelnf("\tbasicLatinChars: %#v,", BasicLatinLookup(ch.Chars, ch.Ranges, ch.UnicodeClasses, ch.IgnoreCase))
	}
	b.writelnf("\tignoreCase: %t,", ch.IgnoreCase)
	b.writelnf("\tinverted: %t,", ch.Inverted)
	b.writelnf("}),")
}

// BasicLatinLookup calculates the decision results for the first 256 characters of the UTF-8 character
// set for a given set of chars, ranges and unicodeClasses to speedup the CharClassMatcher.
func BasicLatinLookup(chars, ranges []rune, unicodeClasses []string, ignoreCase bool) (basicLatinChars [128]bool) {
	for _, rn := range chars {
		if rn < 128 {
			basicLatinChars[rn] = true
			if ignoreCase {
				if unicode.IsLower(rn) {
					basicLatinChars[unicode.ToUpper(rn)] = true
				} else {
					basicLatinChars[unicode.ToLower(rn)] = true
				}
			}
		}
	}
	for i := 0; i < len(ranges); i += 2 {
		if ranges[i] < 128 {
			for j := ranges[i]; j < 128 && j <= ranges[i+1]; j++ {
				basicLatinChars[j] = true
				if ignoreCase {
					if unicode.IsLower(j) {
						basicLatinChars[unicode.ToUpper(j)] = true
					} else {
						basicLatinChars[unicode.ToLower(j)] = true
					}
				}
			}
		}
	}
	for _, cl := range unicodeClasses {
		rt := rangeTable(cl)
		for r := rune(0); r < 128; r++ {
			if unicode.Is(rt, r) {
				basicLatinChars[r] = true
			}
		}
	}
	return
}

func (b *builder) writeChoiceExpr(ch *ast.ChoiceExpr) {
	if ch == nil {
		b.writelnf("nil,")
		return
	}
	b.writelnf("ChoiceExpr({")
	pos := ch.Pos()
	b.writelnf("\tpos: {line: %d, col: %d, offset: %d},", pos.Line, pos.Col, pos.Off)
	if len(ch.Alternatives) > 0 {
		b.writelnf("\talternatives: [")
		for _, alt := range ch.Alternatives {
			b.writeExpr(alt)
		}
		b.writelnf("\t],")
	}
	b.writelnf("}),")
}

func (b *builder) writeLabeledExpr(lab *ast.LabeledExpr) {
	if lab == nil {
		b.writelnf("nil,")
		return
	}
	b.writelnf("LabeledExpr({")
	pos := lab.Pos()
	b.writelnf("\tpos: {line: %d, col: %d, offset: %d},", pos.Line, pos.Col, pos.Off)
	if lab.Label != nil && lab.Label.Val != "" {
		b.writelnf("\tlabel: %q,", lab.Label.Val)
	}
	b.writef("\texpr: ")
	b.writeExpr(lab.Expr)
	b.writelnf("}),")
}

func (b *builder) writeLitMatcher(lit *ast.LitMatcher) {
	if lit == nil {
		b.writelnf("nil,")
		return
	}
	b.writelnf("LitMatcher({")
	pos := lit.Pos()
	b.writelnf("\tpos: {line: %d, col: %d, offset: %d},", pos.Line, pos.Col, pos.Off)
	if lit.IgnoreCase {
		b.writelnf("\tval: %q,", strings.ToLower(lit.Val))
	} else {
		b.writelnf("\tval: %q,", lit.Val)
	}
	b.writelnf("\tignoreCase: %t,", lit.IgnoreCase)
	ignoreCaseFlag := ""
	if lit.IgnoreCase {
		ignoreCaseFlag = "i"
	}
	b.writelnf("\twant: %q,", strconv.Quote(lit.Val)+ignoreCaseFlag)
	b.writelnf("}),")
}

func (b *builder) writeNotCodeExpr(not *ast.NotCodeExpr) {
	if not == nil {
		b.writelnf("nil,")
		return
	}
	b.writelnf("NotCodeExpr({")
	pos := not.Pos()
	if not.FuncIx == 0 {
		not.FuncIx = b.exprIndex
	}
	b.writelnf("\tpos: {line: %d, col: %d, offset: %d},", pos.Line, pos.Col, pos.Off)
	b.writelnf("\trun: _call%s,", b.funcName(not.FuncIx))
	b.writelnf("}),")
}

func (b *builder) writeNotExpr(not *ast.NotExpr) {
	if not == nil {
		b.writelnf("nil,")
		return
	}
	b.writelnf("NotExpr({")
	pos := not.Pos()
	b.writelnf("\tpos: {line: %d, col: %d, offset: %d},", pos.Line, pos.Col, pos.Off)
	b.writef("\texpr: ")
	b.writeExpr(not.Expr)
	b.writelnf("}),")
}

func (b *builder) writeOneOrMoreExpr(one *ast.OneOrMoreExpr) {
	if one == nil {
		b.writelnf("nil,")
		return
	}
	b.writelnf("OneOrMoreExpr({")
	pos := one.Pos()
	b.writelnf("\tpos: {line: %d, col: %d, offset: %d},", pos.Line, pos.Col, pos.Off)
	b.writef("\texpr: ")
	b.writeExpr(one.Expr)
	b.writelnf("}),")
}

func (b *builder) writeRecoveryExpr(recover *ast.RecoveryExpr) {
	if recover == nil {
		b.writelnf("nil,")
		return
	}
	b.writelnf("RecoveryExpr({")
	pos := recover.Pos()
	b.writelnf("\tpos: {line: %d, col: %d, offset: %d},", pos.Line, pos.Col, pos.Off)

	b.writef("\texpr: ")
	b.writeExpr(recover.Expr)
	b.writef("\trecoverExpr: ")
	b.writeExpr(recover.RecoverExpr)
	b.writelnf("\tfailureLabel: [")
	for _, label := range recover.Labels {
		b.writelnf("%q,", label)
	}
	b.writelnf("\t],")
	b.writelnf("}),")
}

func (b *builder) writeRuleRefExpr(ref *ast.RuleRefExpr) {
	if ref == nil {
		b.writelnf("nil,")
		return
	}
	b.writelnf("RuleRefExpr({")
	pos := ref.Pos()
	b.writelnf("\tpos: {line: %d, col: %d, offset: %d},", pos.Line, pos.Col, pos.Off)
	if ref.Name != nil && ref.Name.Val != "" {
		b.writelnf("\tname: %q,", ref.Name.Val)
	}
	b.writelnf("}),")
}

func (b *builder) writeSeqExpr(seq *ast.SeqExpr) {
	if seq == nil {
		b.writelnf("nil,")
		return
	}
	b.writelnf("SeqExpr({")
	pos := seq.Pos()
	b.writelnf("\tpos: {line: %d, col: %d, offset: %d},", pos.Line, pos.Col, pos.Off)
	if len(seq.Exprs) > 0 {
		b.writelnf("\texprs: [")
		for _, e := range seq.Exprs {
			b.writeExpr(e)
		}
		b.writelnf("\t],")
	}
	b.writelnf("}),")
}

func (b *builder) writeStateCodeExpr(state *ast.StateCodeExpr) {
	if state == nil {
		b.writelnf("nil,")
		return
	}
	b.globalState = true
	b.writelnf("StateCodeExpr({")
	pos := state.Pos()
	if state.FuncIx == 0 {
		state.FuncIx = b.exprIndex
	}
	b.writelnf("\tpos: {line: %d, col: %d, offset: %d},", pos.Line, pos.Col, pos.Off)
	b.writelnf("\trun: _call%s,", b.funcName(state.FuncIx))
	b.writelnf("}),")
}

func (b *builder) writeCustomParserCodeExpr(state *ast.CustomParserCodeExpr) {
	if state == nil {
		b.writelnf("nil,")
		return
	}
	b.globalState = true
	b.writelnf("CustomParserCodeExpr({")
	pos := state.Pos()
	if state.FuncIx == 0 {
		state.FuncIx = b.exprIndex
	}
	b.writelnf("\tpos: {line: %d, col: %d, offset: %d},", pos.Line, pos.Col, pos.Off)
	b.writelnf("\trun: _call%s,", b.funcName(state.FuncIx))
	b.writelnf("}),")
}

func (b *builder) writeThrowExpr(throw *ast.ThrowExpr) {
	if throw == nil {
		b.writelnf("nil,")
		return
	}
	b.writelnf("ThrowExpr({")
	pos := throw.Pos()
	b.writelnf("\tpos: {line: %d, col: %d, offset: %d},", pos.Line, pos.Col, pos.Off)
	b.writelnf("\tlabel: %q,", throw.Label)
	b.writelnf("}),")
}

func (b *builder) writeZeroOrMoreExpr(zero *ast.ZeroOrMoreExpr) {
	if zero == nil {
		b.writelnf("nil,")
		return
	}
	b.writelnf("ZeroOrMoreExpr({")
	pos := zero.Pos()
	b.writelnf("\tpos: {line: %d, col: %d, offset: %d},", pos.Line, pos.Col, pos.Off)
	b.writef("\texpr: ")
	b.writeExpr(zero.Expr)
	b.writelnf("}),")
}

func (b *builder) writeZeroOrOneExpr(zero *ast.ZeroOrOneExpr) {
	if zero == nil {
		b.writelnf("nil,")
		return
	}
	b.writelnf("ZeroOrOneExpr({")
	pos := zero.Pos()
	b.writelnf("\tpos: {line: %d, col: %d, offset: %d},", pos.Line, pos.Col, pos.Off)
	b.writef("\texpr: ")
	b.writeExpr(zero.Expr)
	b.writelnf("}),")
}

func (b *builder) writeRuleCode(rule *ast.Rule) {
	if rule == nil || rule.Name == nil {
		return
	}

	// keep trace of the current rule, as the code blocks are created
	// in functions named "on<RuleName><#ExprIndex>".
	b.ruleName = rule.Name.Val
	b.pushArgsSet()
	b.writeExprCode(rule.Expr)
	b.popArgsSet()
}

func (b *builder) pushArgsSet() {
	b.argsStack = append(b.argsStack, nil)
}

func (b *builder) popArgsSet() {
	b.argsStack = b.argsStack[:len(b.argsStack)-1]
}

func (b *builder) addArg(arg *ast.Identifier) {
	if arg == nil {
		return
	}
	ix := len(b.argsStack) - 1
	b.argsStack[ix] = append(b.argsStack[ix], arg.Val)
}

func (b *builder) writeExprCode(expr ast.Expression) {
	switch expr := expr.(type) {
	case *ast.ActionExpr:
		b.writeExprCode(expr.Expr)
		b.writeActionExprCode(expr)

	case *ast.AndCodeExpr:
		b.writeAndCodeExprCode(expr)

	case *ast.LabeledExpr:
		b.addArg(expr.Label)
		b.pushArgsSet()
		b.writeExprCode(expr.Expr)
		b.popArgsSet()

	case *ast.NotCodeExpr:
		b.writeNotCodeExprCode(expr)

	case *ast.AndExpr:
		b.pushArgsSet()
		b.writeExprCode(expr.Expr)
		b.popArgsSet()

	case *ast.ChoiceExpr:
		for _, alt := range expr.Alternatives {
			b.pushArgsSet()
			b.writeExprCode(alt)
			b.popArgsSet()
		}

	case *ast.NotExpr:
		b.pushArgsSet()
		b.writeExprCode(expr.Expr)
		b.popArgsSet()

	case *ast.OneOrMoreExpr:
		b.pushArgsSet()
		b.writeExprCode(expr.Expr)
		b.popArgsSet()

	case *ast.RecoveryExpr:
		b.pushArgsSet()
		b.writeExprCode(expr.Expr)
		b.writeExprCode(expr.RecoverExpr)
		b.popArgsSet()

	case *ast.SeqExpr:
		for _, sub := range expr.Exprs {
			b.writeExprCode(sub)
		}

	case *ast.StateCodeExpr:
		b.writeStateCodeExprCode(expr)

	case *ast.CustomParserCodeExpr:
		b.writeCustomParserCodeExprCode(expr)

	case *ast.ZeroOrMoreExpr:
		b.pushArgsSet()
		b.writeExprCode(expr.Expr)
		b.popArgsSet()

	case *ast.ZeroOrOneExpr:
		b.pushArgsSet()
		b.writeExprCode(expr.Expr)
		b.popArgsSet()
	}
}

func (b *builder) writeActionExprCode(act *ast.ActionExpr) {
	if act == nil {
		return
	}
	if act.FuncIx > 0 {
		b.writeFunc(act.FuncIx, act.Code, callFuncTemplate, onFuncTemplate)
		act.FuncIx = 0 // already rendered, prevent duplicates
	}
}

func (b *builder) writeAndCodeExprCode(and *ast.AndCodeExpr) {
	if and == nil {
		return
	}
	if and.FuncIx > 0 {
		b.writeFunc(and.FuncIx, and.Code, callPredFuncTemplate, onPredFuncTemplate)
		and.FuncIx = 0 // already rendered, prevent duplicates
	}
}

func (b *builder) writeNotCodeExprCode(not *ast.NotCodeExpr) {
	if not == nil {
		return
	}
	if not.FuncIx > 0 {
		b.writeFunc(not.FuncIx, not.Code, callPredFuncTemplate, onPredFuncTemplate)
		not.FuncIx = 0 // already rendered, prevent duplicates
	}
}

func (b *builder) writeStateCodeExprCode(state *ast.StateCodeExpr) {
	if state == nil {
		return
	}
	if state.FuncIx > 0 {
		b.writeFunc(state.FuncIx, state.Code, callStateFuncTemplate, onStateFuncTemplate)
		state.FuncIx = 0 // already rendered, prevent duplicates
	}
}

func (b *builder) writeCustomParserCodeExprCode(state *ast.CustomParserCodeExpr) {
	if state == nil {
		return
	}
	if state.FuncIx > 0 {
		b.writeFunc(state.FuncIx, state.Code, callStateFuncTemplate2, onStateFuncTemplate2)
		state.FuncIx = 0 // already rendered, prevent duplicates
	}
}

func (b *builder) writeFunc(funcIx int, code *ast.CodeBlock, callTpl, funcTpl string) {
	if code == nil {
		return
	}
	val := b.templateRender(strings.TrimSpace(code.Val)[1:len(code.Val)-1], true)
	if len(val) > 0 && val[0] == '\n' {
		val = val[1:]
	}
	if len(val) > 0 && val[len(val)-1] == '\n' {
		val = val[:len(val)-1]
	}
	var args bytes.Buffer
	ix := len(b.argsStack) - 1
	if ix >= 0 {
		for i, arg := range b.argsStack[ix] {
			if i > 0 {
				args.WriteString(", ")
			}
			args.WriteString(arg)
		}
	}
	if args.Len() > 0 {
		args.WriteString(":Any")
	}

	fnNm := b.funcName(funcIx)
	b.writelnf(funcTpl, fnNm, b.recvName, args.String(), val)

	args.Reset()
	if ix >= 0 {
		for i, arg := range b.argsStack[ix] {
			if i > 0 {
				args.WriteString(", ")
			}
			args.WriteString(fmt.Sprintf(`stack[%q]`, arg))
		}
	}
	b.writelnf(callTpl, fnNm, args.String())
}

func (b *builder) writeStaticCode() {
	buffer := bytes.NewBufferString("")
	params := struct {
		Optimize              bool
		BasicLatinLookupTable bool
		GlobalState           bool
		LeftRecursion         bool
		Nolint                bool
	}{
		Optimize:              b.optimize,
		BasicLatinLookupTable: b.basicLatinLookupTable,
		GlobalState:           b.globalState,
		LeftRecursion:         b.haveLeftRecursion,
		Nolint:                b.nolint,
	}
	t := template.Must(template.New("static_code").Parse(staticCode))

	err := t.Execute(buffer, params)
	if err != nil {
		// This is very unlikely to ever happen
		panic("executing template: " + err.Error())
	}

	// Clean the ==template== comments from the generated parser
	lines := strings.Split(buffer.String(), "\n")
	buffer.Reset()
	re := regexp.MustCompile(`^\s*//\s*(==template==\s*)+$`)
	reLineEnd := regexp.MustCompile(`//\s*==template==\s*$`)
	for _, line := range lines {
		if !re.MatchString(line) {
			line = reLineEnd.ReplaceAllString(line, "")
			_, err := buffer.WriteString(line + "\n")
			if err != nil {
				// This is very unlikely to ever happen
				panic("unable to write to byte buffer: " + err.Error())
			}
		}
	}

	b.writeln(buffer.String())
	if b.rangeTable {
		b.writeln(rangeTable0)
	}
}

func (b *builder) funcName(ix int) string {
	return "on" + b.ruleName + strconv.Itoa(ix)
}

func (b *builder) writef(f string, args ...any) {
	if b.err == nil {
		_, b.err = fmt.Fprintf(b.w, f, args...)
	}
}

func (b *builder) writelnf(f string, args ...any) {
	b.writef(f+"\n", args...)
}

func (b *builder) writeln(f string) {
	if b.err == nil {
		_, b.err = fmt.Fprint(b.w, f+"\n")
	}
}
