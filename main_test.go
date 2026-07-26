package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMain(t *testing.T) {
	args, stdout, stderr := os.Args, os.Stdout, os.Stderr
	os.Stdout, _ = os.Open(os.DevNull)
	os.Stderr, _ = os.Open(os.DevNull)
	defer func() {
		exit = os.Exit
		os.Args = args
		os.Stdout = stdout
		os.Stderr = stderr
	}()
	exit = func(code int) {
		panic(code)
	}
	grammarFile := filepath.Join(t.TempDir(), "grammar.peg")
	if err := os.WriteFile(grammarFile, []byte(`start = "a"`), 0o600); err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		args []string
		code int
	}{
		{args: nil, code: 3},                        // stdin: no match found
		{args: []string{"-h"}, code: 0},             // help
		{args: []string{"FILE1", "FILE2"}, code: 1}, // want only 1 non-flag arg
		{args: []string{"-x"}, code: 3},             // stdin: no match found
		{args: []string{"-t", "hx", "-x", grammarFile}},
		{args: []string{"-target", "hx", "-x", grammarFile}},
		{args: []string{"-t", "ts", "-x", grammarFile}},
		{args: []string{"-t", "cs", "-x", grammarFile}},
		{args: []string{"-t", "rust", "-x", grammarFile}},
		{args: []string{"-t", "unknown", "-x", grammarFile}, code: 1},
	}

	for _, tc := range cases {
		os.Args = append([]string{"pigeon"}, tc.args...)

		got := runMainRecover()
		if got != tc.code {
			t.Errorf("%q: want code %d, got %d", tc.args, tc.code, got)
		}
	}
}

func TestParserBuilderForTarget(t *testing.T) {
	for _, target := range []string{
		"go", "hx", "haxe", "ts", "typescript", "cs", "csharp", "c#", "rust",
	} {
		builder, formatGo, err := parserBuilderForTarget(target)
		if err != nil {
			t.Errorf("%q: unexpected error: %v", target, err)
			continue
		}
		if builder == nil {
			t.Errorf("%q: nil builder", target)
		}
		if formatGo != (target == "go") {
			t.Errorf("%q: formatGo = %t", target, formatGo)
		}
	}

	if _, _, err := parserBuilderForTarget("invalid"); err == nil {
		t.Fatal("expected invalid target error")
	}
}

func TestUsageIncludesTargetFlagAliases(t *testing.T) {
	if want := "-t TARGET, -target TARGET"; !strings.Contains(usagePage, want) {
		t.Fatalf("usage page missing %q", want)
	}
}

func runMainRecover() (code int) {
	defer func() {
		if e := recover(); e != nil {
			if i, ok := e.(int); ok {
				code = i
				return
			}
			panic(e)
		}
	}()
	main()
	return 0
}
