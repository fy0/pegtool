package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"golang.org/x/tools/imports"

	"github.com/fy0/pigeon/ast"
	builderGo "github.com/fy0/pigeon/builder"
	builderCs "github.com/fy0/pigeon/builder_cs"
	builderHx "github.com/fy0/pigeon/builder_hx"
	builderRust "github.com/fy0/pigeon/builder_rust"
	builderTs "github.com/fy0/pigeon/builder_ts"
)

// exit function mockable for tests
var exit = os.Exit

// ruleNamesFlag is a custom flag that parses a comma-separated
// list of rule names. It implements flag.Value.
type ruleNamesFlag []string

func (r *ruleNamesFlag) String() string {
	return fmt.Sprint(*r)
}

func (r *ruleNamesFlag) Set(value string) error {
	names := strings.Split(value, ",")
	*r = append(*r, names...)
	return nil
}

func ParseReader(filename string, r io.Reader, opts ...option) (any, error) { // nolint: deadcode
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return Parse(filename, b, opts...)
}

// Parse parses the data from b using filename as information in the
// error messages.
func Parse(filename string, b []byte, opts ...option) (any, error) {
	return newParser(filename, b, opts...).parse(g)
}

func main() {
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// define command-line flags
	var (
		// noRecoverFlag      = fs.Bool("no-recover", false, "do not recover from panic")
		dbgFlag            = fs.Bool("debug", false, "set debug mode")
		shortHelpFlag      = fs.Bool("h", false, "show help page")
		longHelpFlag       = fs.Bool("help", false, "show help page")
		nolint             = fs.Bool("nolint", false, "add '// nolint: ...' comments to suppress warnings by gometalinter or golangci-lint")
		outputFlag         = fs.String("o", "", "output file, defaults to stdout")
		optimizeParserFlag = fs.Bool("optimize-parser", false, "generate optimized parser without Debug options")
		recvrNmFlag        = fs.String("receiver-name", "c", "receiver name for the generated methods")
		noBuildFlag        = fs.Bool("x", false, "do not build, only parse")

		cacheFlag = fs.Bool("cache", false, "cache parsing results")

		grammarNameFlag        = fs.String("grammar-name", "g", "default is g, `var g = &grammar{ ... }")
		runFuncPrefixFlag      = fs.String("run-func-prefix", "", "set prefix for generated function name: `(*parser).call_onXXX`. For multiple peg files")
		grammarOnlyFlag        = fs.Bool("grammar-only", false, "use it when you have multiple peg files")
		optimizeRefExprByIndex = fs.Bool("optimize-ref-expr-by-index", false, "generate optimized parser grammar find RefExpr by index (~10% increased)")
		targetFlag             = fs.String("t", "go", "build target, default go")

		// optimizeGrammar        = fs.Bool("optimize-grammar", false, "optimize the given grammar (EXPERIMENTAL FEATURE)")

		altEntrypointsFlag ruleNamesFlag
	)
	fs.StringVar(targetFlag, "target", "go", "build target, default go")
	fs.Var(&altEntrypointsFlag, "alternate-entrypoints", "comma-separated list of rule names that may be used as entrypoints")

	fs.Usage = usage
	err := fs.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "args parse error:\n", err)
		exit(6)
	}

	if *shortHelpFlag || *longHelpFlag {
		fs.Usage()
		exit(0)
	}

	if fs.NArg() > 1 {
		argError(1, "expected one argument, got %q", strings.Join(fs.Args(), " "))
	}

	buildParser, formatGo, err := parserBuilderForTarget(*targetFlag)
	if err != nil {
		argError(1, "%v", err)
	}

	// get input source
	infile := ""
	if fs.NArg() == 1 {
		infile = fs.Arg(0)
	}
	nm, rc := input(infile)
	defer func() {
		err = rc.Close()
		if err != nil {
			fmt.Fprintln(os.Stderr, "close file error:\n", err)
		}
		if r := recover(); r != nil {
			panic(r)
		}
		if err != nil {
			exit(7)
		}
	}()

	// parse input
	// , Recover(!*noRecoverFlag)
	g, err := ParseReader(nm, rc, debug(*dbgFlag), memoized(*cacheFlag))
	if err != nil {
		fmt.Fprintln(os.Stderr, "parse error(s):\n", err)
		exit(3)
	}

	// validate alternate entrypoints
	grammar := g.(*ast.Grammar)
	rules := make(map[string]struct{}, len(grammar.Rules))
	for _, rule := range grammar.Rules {
		rules[rule.Name.Val] = struct{}{}
	}
	for _, entrypoint := range altEntrypointsFlag {
		if entrypoint == "" {
			continue
		}

		if _, ok := rules[entrypoint]; !ok {
			fmt.Fprintf(os.Stderr, "argument error:\nunknown rule name %s used as alternate entrypoint\n", entrypoint)
			exit(9)
		}
	}

	if !*noBuildFlag {
		// if *optimizeGrammar {
		// ast.Optimize(grammar, altEntrypointsFlag...)
		// }

		// generate parser
		out := output(*outputFlag)
		defer func() {
			err := out.Close()
			if err != nil {
				fmt.Fprintln(os.Stderr, "close file error:\n", err)
				exit(8)
			}
		}()

		outBuf := bytes.NewBuffer([]byte{})

		curNmOpt := builderGo.ReceiverName(*recvrNmFlag)
		optimizeParser := builderGo.Optimize(*optimizeParserFlag)
		nolintOpt := builderGo.Nolint(*nolint)
		refExprByIndex := builderGo.OptimizeRefExprByIndex(*optimizeRefExprByIndex)
		runFuncPrefix := builderGo.RunFuncPrefix(*runFuncPrefixFlag)
		grammarOnly := builderGo.GrammarOnly(*grammarOnlyFlag)
		grammarName := builderGo.GrammarName(*grammarNameFlag)

		if err := buildParser(
			outBuf, grammar, curNmOpt, optimizeParser,
			runFuncPrefix, grammarOnly, grammarName,
			nolintOpt, refExprByIndex); err != nil {
			fmt.Fprintln(os.Stderr, "build error: ", err)
			exit(5)
		}

		// Defaults from golang.org/x/tools/cmd/goimports
		options := &imports.Options{
			TabWidth:  8,
			TabIndent: true,
			Comments:  true,
			Fragment:  true,
		}

		if formatGo {
			formattedBuf, err := imports.Process("filename", outBuf.Bytes(), options)
			if err != nil {
				if _, err := out.Write(outBuf.Bytes()); err != nil {
					fmt.Fprintln(os.Stderr, "write error: ", err)
					exit(7)
				}
				fmt.Fprintln(os.Stderr, "format error: ", err)
				exit(6)
			}

			if _, err := out.Write(formattedBuf); err != nil {
				fmt.Fprintln(os.Stderr, "write error: ", err)
				exit(7)
			}
		} else {
			if _, err := out.Write(outBuf.Bytes()); err != nil {
				fmt.Fprintln(os.Stderr, "write error: ", err)
				exit(7)
			}
		}
	}
}

type parserBuilder func(io.Writer, *ast.Grammar, ...builderGo.Option) error

func parserBuilderForTarget(target string) (parserBuilder, bool, error) {
	switch strings.ToLower(strings.TrimSpace(target)) {
	case "go":
		return builderGo.BuildParser, true, nil
	case "hx", "haxe":
		return builderHx.BuildParser, false, nil
	case "ts", "typescript":
		return builderTs.BuildParser, false, nil
	case "cs", "csharp", "c#":
		return builderCs.BuildParser, false, nil
	case "rust":
		return builderRust.BuildParser, false, nil
	default:
		return nil, false, fmt.Errorf(
			"unknown build target %q (supported: go, hx, ts, cs, rust)",
			target,
		)
	}
}

var usagePage = `usage: %s [options] [GRAMMAR_FILE]

Pigeon generates a parser based on a PEG grammar.

By default, pigeon reads the grammar from stdin and writes the
generated parser to stdout. If GRAMMAR_FILE is specified, the
grammar is read from this file instead. If the -o flag is set,
the generated code is written to this file instead.

	-cache
		cache parser results to avoid exponential parsing time in
		pathological cases. Can make the parsing slower for typical
		cases and uses more memory.
	-debug
		output debugging information while parsing the grammar.
	-h -help
		display this help message.
	-nolint
		add '// nolint: ...' comments for generated parser to suppress
		warnings by gometalinter (https://github.com/alecthomas/gometalinter) or
		golangci-lint (https://golangci-lint.run/).
	-no-recover
		do not recover from a panic. Useful to access the panic stack
		when debugging, otherwise the panic is converted to an error.
	-o OUTPUT_FILE
		write the generated parser to OUTPUT_FILE. Defaults to stdout.
	-optimize-ref-expr-by-index
		generate optimized parser grammar find RefExpr by index (~10%% performance increased, cause more git line diff)
	-optimize-parser
		generate optimized parser without Debug and Memoize options and
		with some other optimizations applied.
	-receiver-name NAME
		use NAME as for the receiver name of the generated methods
		for the grammar's code blocks. Defaults to "c".
	-x
		do not generate the parser, only parse the grammar.
 	-alternate-entrypoints RULE[,RULE...]
		comma-separated list of rule names that may be used as alternate
		entrypoints for the parser, in addition to the first rule in the
		grammar.
	-grammar-only
		generate grammar part only, used for multiple peg files.
	-grammar-name
		to set variable name of grammar, default is g, "var g = &grammar{ ... }"
	-run-func-prefix
		set prefix for generated function name: "(*parser).call_onXXX". For multiple peg files
	-t TARGET, -target TARGET
		build target: go, hx, ts, cs, or rust. Defaults to go.

See https://godoc.org/github.com/mna/pigeon for more information.
This version is a fork: https://github.com/fy0/pigeon
`

// usage prints the help page of the command-line tool.
func usage() {
	fmt.Printf(usagePage, os.Args[0])
}

// argError prints an error message to stderr, prints the command usage
// and exits with the specified exit code.
func argError(exitCode int, msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg, args...)
	fmt.Fprintln(os.Stderr)
	usage()
	exit(exitCode)
}

// input gets the name and reader to get input text from.
func input(filename string) (nm string, rc io.ReadCloser) {
	nm = "stdin"
	inf := os.Stdin
	if filename != "" {
		f, err := os.Open(filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			exit(2)
		}
		inf = f
		nm = filename
	}
	r := bufio.NewReader(inf)
	return nm, makeReadCloser(r, inf)
}

// output gets the writer to write the generated parser to.
func output(filename string) io.WriteCloser {
	out := os.Stdout
	if filename != "" {
		f, err := os.Create(filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			exit(4)
		}
		out = f
	}
	return out
}

// create a ReadCloser that reads from r and closes c.
func makeReadCloser(r io.Reader, c io.Closer) io.ReadCloser {
	rc := struct {
		io.Reader
		io.Closer
	}{r, c}
	return io.ReadCloser(rc)
}

// astPos is a helper method for the PEG grammar parser. It returns the
// position of the current match as an ast.Pos.
func (c *current) astPos() ast.Pos {
	return ast.Pos{Line: c.pos.line, Col: c.pos.col, Off: c.pos.offset}
}

// toAnySlice is a helper function for the PEG grammar parser. It converts
// v to a slice of empty interfaces.
func toAnySlice(v any) []any {
	if v == nil {
		return nil
	}
	return v.([]any)
}

// validateUnicodeEscape checks that the provided escape sequence is a
// valid Unicode escape sequence.
func validateUnicodeEscape(escape, errMsg string) (any, error) {
	r, _, _, err := strconv.UnquoteChar("\\"+escape, '"')
	if err != nil {
		return nil, errors.New(errMsg)
	}
	if 0xD800 <= r && r <= 0xDFFF {
		return nil, errors.New(errMsg)
	}
	return nil, nil
}
