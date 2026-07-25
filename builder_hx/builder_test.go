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
