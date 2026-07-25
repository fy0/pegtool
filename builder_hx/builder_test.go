package builder

import (
	"bytes"
	"strings"
	"testing"

	"github.com/fy0/pigeon/bootstrap"
	base "github.com/fy0/pigeon/builder"
)

func TestBuildParserIncludesNoMatchFormatterHook(t *testing.T) {
	p := bootstrap.NewParser()
	g, err := p.Parse("", strings.NewReader(`start = "a"`))
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		name string
		opts []base.Option
	}{
		{name: "standard"},
		{name: "optimized", opts: []base.Option{base.Optimize(true)}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer
			if err := BuildParser(&out, g, tc.opts...); err != nil {
				t.Fatal(err)
			}

			generated := out.String()
			for _, snippet := range []string{
				"typedef NoMatchErrorFormatter = (position:Position, input:String, expected:Array<String>) -> Exception;",
				"function noMatchErrorFormatter(fn:NoMatchErrorFormatter):Option",
				"public var noMatchErrorFormatter:NoMatchErrorFormatter;",
				"function buildNoMatchError(pos:Position, expected:Array<String>):Exception",
				"this.addErrAt(this.buildNoMatchError(this.maxFailPos, expected), this.maxFailPos, expected);",
				"var _choiceNoMatch = -1;",
			} {
				if !strings.Contains(generated, snippet) {
					t.Fatalf("generated parser missing snippet %q", snippet)
				}
			}
			if strings.Contains(generated, "var choiceNoMatch = -1;") {
				t.Fatal("generated parser declares a sentinel shadowed by Parser.choiceNoMatch")
			}
			if tc.name == "standard" {
				for _, snippet := range []string{
					"if (altI == _choiceNoMatch)",
					"this.incChoiceAltCnt(ch, _choiceNoMatch);",
				} {
					if !strings.Contains(generated, snippet) {
						t.Fatalf("generated parser missing no-match sentinel use %q", snippet)
					}
				}
			}
		})
	}
}

func TestBuildParserOptimizesRuleReferencesByIndex(t *testing.T) {
	p := bootstrap.NewParser()
	g, err := p.Parse("", strings.NewReader("start = other\nother = \"a\""))
	if err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	if err := BuildParser(&out, g, base.OptimizeRefExprByIndex(true)); err != nil {
		t.Fatal(err)
	}

	generated := out.String()
	for _, want := range []string{
		"ret = this.parseRuleIRefExpr(expr);",
		"function parseRuleIRefExpr(ref:RuleIRefExpr):RetValOk {",
		"return this.parseRuleWrap(this.rulesArray[ref.index]);",
	} {
		if !strings.Contains(generated, want) {
			t.Fatalf("generated parser missing %q", want)
		}
	}
	for _, invalid := range []string{
		"ret = p.parseRuleIRefExpr(expr)",
		"function parseSeqExpr(ref:RuleIRefExpr):RetValOk {",
		"return this.parseRuleWrap(this.rulesArray[ref.index])\n",
	} {
		if strings.Contains(generated, invalid) {
			t.Fatalf("generated parser contains invalid Haxe code %q", invalid)
		}
	}
}

func TestBuildParserIndentsNestedGrammar(t *testing.T) {
	p := bootstrap.NewParser()
	g, err := p.Parse("", strings.NewReader("start = item\nitem = value:(\"a\" / \"b\")+ !."))
	if err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	if err := BuildParser(
		&out,
		g,
		base.GrammarOnly(true),
		base.GrammarName("g"),
		base.OptimizeRefExprByIndex(true),
	); err != nil {
		t.Fatal(err)
	}

	want := `var g:Grammar = {
	rules: [
		{
			name: "start",
			expr: RuleIRefExpr({index: 1 /* item */}),
		},
		{
			name: "item",
			varExists: true,
			expr: SeqExpr({
				exprs: [
					LabeledExpr({
						label: "value",
						expr: OneOrMoreExpr({
							expr: ChoiceExpr({
								alternatives: [
									LitMatcher({val: "a", want: "\"a\"",}),
									LitMatcher({val: "b", want: "\"b\"",}),
								],
							}),
						}),
						textCapture: false,
					}),
					NotExpr({
						expr: AnyMatcher({}),
					}),
				],
			}),
		},
	],
}
`
	if got := out.String(); got != want {
		t.Fatalf("unexpected generated grammar:\nwant:\n%s\ngot:\n%s", want, got)
	}
}

func TestBuildParserFormatsCallbacks(t *testing.T) {
	p := bootstrap.NewParser()
	g, err := p.Parse("", strings.NewReader(`{
class ParserCustomData {
	public function new() {}
}
}
start = value:"a" { return value; }
`))
	if err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	if err := BuildParser(&out, g, base.GrammarName("g")); err != nil {
		t.Fatal(err)
	}

	generated := out.String()
	for _, want := range []string{
		"// Code is generated; DO NOT EDIT.\nimport haxe.Exception;",
		"import haxe.iterators.StringIteratorUnicode;\n\nclass ParserCustomData",
		"private function _call_onstart_1(p:Parser):Any {\n\tvar stack = p.vstack.first();\n\tvar c = p.cur;",
		"return (function(c:Current, value:Any):Any {",
		"})(p.cur, stack[\"value\"]);",
	} {
		if !strings.Contains(generated, want) {
			t.Fatalf("generated parser missing %q", want)
		}
	}

	for lineNo, line := range strings.Split(generated, "\n") {
		if strings.TrimRight(line, " \t") != line {
			t.Fatalf("line %d has trailing whitespace: %q", lineNo+1, line)
		}
	}
}

func TestBuildParserWrapsLongUnicodeClassArrays(t *testing.T) {
	p := bootstrap.NewParser()
	g, err := p.Parse("", strings.NewReader(`start = [\p{L}\p{Other_ID_Start}\p{Nl}\p{Mn}\p{Mc}\p{Nd}\p{Pc}\p{Other_ID_Continue}]`))
	if err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	if err := BuildParser(&out, g, base.GrammarOnly(true), base.GrammarName("g")); err != nil {
		t.Fatal(err)
	}

	want := `classes: [
					Unicode.L,
					Unicode.Other_ID_Start,
					Unicode.Nl,`
	if !strings.Contains(out.String(), want) {
		t.Fatalf("generated parser missing wrapped Unicode classes:\n%s", out.String())
	}
}
