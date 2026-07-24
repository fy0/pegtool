package builder

import (
	"bytes"
	"io"
	"strings"
	"testing"

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
		})
	}
}
