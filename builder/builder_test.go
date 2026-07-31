package builder

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fy0/pigeon/ast"
	"github.com/fy0/pigeon/bootstrap"
)

var grammar = `
{
var test = "some string"

func init() {
	fmt.Println("this is inside the init")
}
}

start = additive eof
additive = left:multiplicative "+" space right:additive {
	fmt.Println(left, right)
} / mul:multiplicative { fmt.Println(mul) }
multiplicative = left:primary op:"*" space right:multiplicative { fmt.Println(left, right, op) } / primary
primary = integer / "(" space additive:additive ")" space { fmt.Println(additive) }
integer "integer" = digits:[0123456789]+ space { fmt.Println(digits) }
space = ' '*
eof = !. { fmt.Println("eof") }
`

func TestBuildParser(t *testing.T) {
	p := bootstrap.NewParser()
	g, err := p.Parse("", strings.NewReader(grammar))
	if err != nil {
		t.Fatal(err)
	}
	if err := BuildParser(io.Discard, g); err != nil {
		t.Fatal(err)
	}
}

func TestBuildParserIncludesNoMatchFormatterHook(t *testing.T) {
	p := bootstrap.NewParser()
	g, err := p.Parse("", strings.NewReader(grammar))
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		name string
		opts []Option
	}{
		{name: "standard"},
		{name: "optimized", opts: []Option{Optimize(true)}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer
			if err := BuildParser(&out, g, tc.opts...); err != nil {
				t.Fatal(err)
			}

			generated := out.String()
			for _, snippet := range []string{
				"func noMatchErrorFormatter(fn func(position, []byte, []string) error) option",
				"noMatchErrorFormatter func(position, []byte, []string) error",
				"func (p *parser) buildNoMatchError(pos position, expected []string) error",
				"p.addErrAt(p.buildNoMatchError(p.maxFailPos, expected), p.maxFailPos, expected)",
			} {
				if !strings.Contains(generated, snippet) {
					t.Fatalf("generated parser missing snippet %q", snippet)
				}
			}
			if strings.HasSuffix(generated, "\n\n") {
				t.Fatal("generated parser has extra trailing blank lines")
			}
			for lineNo, line := range strings.Split(generated, "\n") {
				if strings.TrimRight(line, " \t") != line {
					t.Fatalf("line %d has trailing whitespace: %q", lineNo+1, line)
				}
			}
		})
	}
}

func TestBuildParserMemoizesNamedRulesWithAllocationEfficientTable(t *testing.T) {
	p := bootstrap.NewParser()
	g, err := p.Parse("", strings.NewReader(`start = "a"`))
	if err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	if err := BuildParser(&out, g, Optimize(true)); err != nil {
		t.Fatal(err)
	}

	generated := out.String()
	for _, snippet := range []string{
		"memo1       []map[*rule]*resultTuple",
		"memoBlocks  []*[memoBlockSize]resultTuple",
		"func (p *parser) newMemoizedResult(val any, ok bool) *resultTuple",
		"p.memo1 = make([]map[*rule]*resultTuple, len(p.data)+1)",
		"recoveryKey = recoveryMemoKey{p.pt.offset, r}",
		"cached = memo[p.pt.offset][r]",
		"entries = make(map[*rule]*resultTuple)",
		"A rule cache hit still represents its root expression.",
	} {
		if !strings.Contains(generated, snippet) {
			t.Fatalf("generated parser missing snippet %q", snippet)
		}
	}
	for _, snippet := range []string{
		"[]map[any]*resultTuple",
		"memo[p.pt.offset][expr]",
		"setMemoized := func",
	} {
		if strings.Contains(generated, snippet) {
			t.Fatalf("generated parser still contains expression memoization %q", snippet)
		}
	}
}

func TestGeneratedParserMemoizesNamedRules(t *testing.T) {
	p := bootstrap.NewParser()
	g, err := p.Parse("memo.peg", strings.NewReader(`{
package generated

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type ParserCustomData struct{}

var actionCalls int
var failureCalls int
}

start = (word "x" / word "y") !. { return actionCalls }
lookahead = &word word !. { return actionCalls }
failure = failed? failed? "z" !.
failed = ("a"? { failureCalls++; return nil }) "x"
collision = first / second
first = "x"
second = "y"
diagnostic = word "x" / word "y"
word = [a-w]+ { actionCalls++; return string(c.text) }
limit = atom "x" / atom "y"
atom = "a"
statistics = (choice "x" / choice "y") !.
choice = "a" / "b"
`))
	if err != nil {
		t.Fatal(err)
	}

	pos := ast.Pos{Filename: "memo.peg", Line: 1, Col: 1}
	ruleRef := func(name string) *ast.RuleRefExpr {
		ref := ast.NewRuleRefExpr(pos)
		ref.Name = ast.NewIdentifier(pos, name)
		return ref
	}
	rule := func(name string, expr ast.Expression) *ast.Rule {
		r := ast.NewRule(pos, ast.NewIdentifier(pos, name))
		r.Expr = expr
		return r
	}
	recovery := func(literal string) *ast.RecoveryExpr {
		expr := ast.NewRecoveryExpr(pos)
		expr.Expr = ruleRef("throwing")
		expr.RecoverExpr = ast.NewLitMatcher(pos, literal)
		expr.Labels = []ast.FailureLabel{"expected"}
		return expr
	}
	throwing := ast.NewThrowExpr(pos)
	throwing.Label = "expected"
	firstAlternative := ast.NewSeqExpr(pos)
	firstAlternative.Exprs = []ast.Expression{ruleRef("recoverA"), ast.NewLitMatcher(pos, "x")}
	recoveryContext := ast.NewChoiceExpr(pos)
	recoveryContext.Alternatives = []ast.Expression{firstAlternative, ruleRef("recoverB")}
	recoveryMemo := ast.NewRecoveryExpr(pos)
	recoveryMemo.Expr = ruleRef("start")
	recoveryMemo.RecoverExpr = ast.NewLitMatcher(pos, "fallback")
	recoveryMemo.Labels = []ast.FailureLabel{"expected"}
	g.Rules = append(g.Rules,
		rule("recoveryContext", recoveryContext),
		rule("recoveryMemo", recoveryMemo),
		rule("recoverA", recovery("a")),
		rule("recoverB", recovery("b")),
		rule("throwing", throwing),
	)

	generatedTest := []byte(`package generated

import (
	"strings"
	"testing"
)

func parseEntry(entrypoint, input string, useMemo bool) (any, *parser, error) {
	actionCalls = 0
	failureCalls = 0
	p := newParser("memo.peg", []byte(input), memoized(useMemo))
	p.entrypoint = entrypoint
	value, err := p.parse(g)
	return value, p, err
}

func TestNamedRuleMemoization(t *testing.T) {
	value, memoParser, err := parseEntry("start", "abcy", true)
	if err != nil {
		t.Fatal(err)
	}
	if value != 1 {
		t.Fatalf("memoized action calls = %v, want 1", value)
	}
	word := memoParser.rules["word"]
	if result := memoParser.memo1[0][word]; result == nil || !result.b || result.end.offset != 3 {
		t.Fatalf("successful word result was not memoized: %#v", result)
	}

	value, _, err = parseEntry("start", "abcy", false)
	if err != nil {
		t.Fatal(err)
	}
	if value != 2 {
		t.Fatalf("plain action calls = %v, want 2", value)
	}

	value, lookaheadParser, err := parseEntry("lookahead", "abc", true)
	if err != nil {
		t.Fatal(err)
	}
	if value != 1 {
		t.Fatalf("lookahead action calls = %v, want 1", value)
	}
	word = lookaheadParser.rules["word"]
	if lookaheadParser.memo1[0][word] == nil || lookaheadParser.memo2[0][word] == nil {
		t.Fatal("normal and lookahead rule results were not memoized separately")
	}

	_, failureParser, err := parseEntry("failure", "z", true)
	if err != nil {
		t.Fatal(err)
	}
	failed := failureParser.memo1[0][failureParser.rules["failed"]]
	if failed == nil || failed.b {
		t.Fatalf("failed rule result was not memoized: %#v", failed)
	}
	if failureCalls != 1 {
		t.Fatalf("memoized failed rule calls = %d, want 1", failureCalls)
	}
	_, _, err = parseEntry("failure", "z", false)
	if err != nil {
		t.Fatal(err)
	}
	if failureCalls != 2 {
		t.Fatalf("plain failed rule calls = %d, want 2", failureCalls)
	}

	if _, _, err := parseEntry("collision", "y", true); err != nil {
		t.Fatal(err)
	}
	_, _, memoErr := parseEntry("diagnostic", "abcz", true)
	_, _, plainErr := parseEntry("diagnostic", "abcz", false)
	if memoErr == nil || plainErr == nil || memoErr.Error() != plainErr.Error() {
		t.Fatalf("memoized diagnostic differs: memo=%v plain=%v", memoErr, plainErr)
	}

	_, counted, err := parseEntry("limit", "ay", true)
	if err != nil {
		t.Fatal(err)
	}
	if counted.ExprCnt != 9 {
		t.Fatalf("memoized expression count = %d, want 9", counted.ExprCnt)
	}
	_, plainCounted, err := parseEntry("limit", "ay", false)
	if err != nil {
		t.Fatal(err)
	}
	if plainCounted.ExprCnt != 9 {
		t.Fatalf("plain expression count = %d, want 9", plainCounted.ExprCnt)
	}
	limited := newParser("memo.peg", []byte("ay"), memoized(true))
	limited.entrypoint = "limit"
	limited.maxExprCnt = 8
	limited.recover = true
	if _, err := limited.parse(g); err == nil || !strings.Contains(err.Error(), "max number of expressions") {
		t.Fatalf("expression limit did not include a rule cache hit: %v", err)
	}

	for _, tc := range []struct {
		input   string
		wantErr bool
	}{
		{input: "b"},
		{input: "a", wantErr: true},
	} {
		_, _, memoErr := parseEntry("recoveryContext", tc.input, true)
		_, _, plainErr := parseEntry("recoveryContext", tc.input, false)
		if (memoErr != nil) != tc.wantErr || (plainErr != nil) != tc.wantErr {
			t.Fatalf("recovery context mismatch for %q: memo=%v plain=%v wantErr=%t", tc.input, memoErr, plainErr, tc.wantErr)
		}
	}

	value, _, err = parseEntry("recoveryMemo", "abcy", true)
	if err != nil {
		t.Fatal(err)
	}
	if value != 1 {
		t.Fatalf("memoized recovery-frame action calls = %v, want 1", value)
	}
	value, _, err = parseEntry("recoveryMemo", "abcy", false)
	if err != nil {
		t.Fatal(err)
	}
	if value != 2 {
		t.Fatalf("plain recovery-frame action calls = %v, want 2", value)
	}
}
`)

	generatedStatisticsTest := []byte(`package generated

import "testing"

func TestNamedRuleMemoizationStatistics(t *testing.T) {
	parse := func(useMemo bool) *Stats {
		stats := &Stats{}
		p := newParser("memo.peg", []byte("ay"), memoized(useMemo), statistics(stats, "no match"))
		p.entrypoint = "statistics"
		if _, err := p.parse(g); err != nil {
			t.Fatal(err)
		}
		return stats
	}

	memoStats := parse(true)
	plainStats := parse(false)
	if got := memoStats.ChoiceAltCnt["choice"]["1"]; got != 1 {
		t.Fatalf("memoized choice count = %d, want 1", got)
	}
	if got := plainStats.ChoiceAltCnt["choice"]["1"]; got != 2 {
		t.Fatalf("plain choice count = %d, want 2", got)
	}
}
`)

	for _, tc := range []struct {
		name string
		opts []Option
	}{
		{name: "standard"},
		{name: "optimized-indexed", opts: []Option{Optimize(true), OptimizeRefExprByIndex(true)}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var first bytes.Buffer
			if err := BuildParser(&first, g, tc.opts...); err != nil {
				t.Fatal(err)
			}
			var second bytes.Buffer
			if err := BuildParser(&second, g, tc.opts...); err != nil {
				t.Fatal(err)
			}
			if first.String() != second.String() {
				t.Fatal("repeated generation is not deterministic")
			}

			dir := t.TempDir()
			files := map[string][]byte{
				"go.mod":         []byte("module generated\n\ngo 1.25.0\n"),
				"parser.go":      first.Bytes(),
				"parser_test.go": generatedTest,
			}
			if tc.name == "standard" {
				files["statistics_test.go"] = generatedStatisticsTest
			}
			for name, data := range files {
				if err := os.WriteFile(filepath.Join(dir, name), data, 0o600); err != nil {
					t.Fatal(err)
				}
			}
			cmd := exec.Command("go", "test", ".", "-count=1")
			cmd.Dir = dir
			cmd.Env = append(os.Environ(), "GOWORK=off")
			if output, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("generated parser failed: %v\n%s\nGenerated parser:\n%s", err, output, first.String())
			}
		})
	}
}

func TestBuildParserIndentsGrammar(t *testing.T) {
	p := bootstrap.NewParser()
	g, err := p.Parse("", strings.NewReader("start = other\nother = \"a\""))
	if err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	if err := BuildParser(&out, g, GrammarOnly(true)); err != nil {
		t.Fatal(err)
	}

	want := `var g = &grammar {
	rules: []*rule{
		{
			name: "start",
			expr: &ruleRefExpr{name: "other",},
		},
		{
			name: "other",
			expr: &litMatcher{val: "a", want: "\"a\"",},
		},
	},
}
`
	if got := out.String(); got != want {
		t.Fatalf("unexpected generated grammar:\nwant:\n%s\ngot:\n%s", want, got)
	}
}

func TestBuildParserKeepsArrayLabelsInline(t *testing.T) {
	p := bootstrap.NewParser()
	g, err := p.Parse("", strings.NewReader(`start = "a" / "b"`))
	if err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	if err := BuildParser(&out, g, GrammarOnly(true)); err != nil {
		t.Fatal(err)
	}

	if want := "alternatives: []any{"; !strings.Contains(out.String(), want) {
		t.Fatalf("generated grammar missing %q", want)
	}
}
