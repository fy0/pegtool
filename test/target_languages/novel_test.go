package target_languages_test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestNovelGrammarTargets(t *testing.T) {
	t.Run("typescript", testNovelTypeScript)
	t.Run("csharp", testNovelCSharp)
	t.Run("rust", testNovelRust)
}

func testNovelTypeScript(t *testing.T) {
	tsc := requireTool(t, "tsc")
	node := requireTool(t, "node")
	dir, generated := generateNovelParser(t, "ts")

	writeFile(t, filepath.Join(dir, "parser.ts"), generated)
	writeFile(t, filepath.Join(dir, "runner.ts"), []byte(`
import { parse } from "./parser";
const source = `+quoteJavaScript(readStory(t))+`;
const sections = parse(source) as string[];
if (sections.length !== 5) throw new Error("section count: " + sections.length);
const first = sections[0].split("\u001f");
const branch = sections[1].split("\u001f");
const last = sections[4].split("\u001f");
if (first[0] !== "\u5f00\u573a") throw new Error("first section: " + first[0]);
if (!first[3].includes("<codeInText>player")) throw new Error("missing inline code");
if (!first[3].includes("<codeBlock> val = select")) throw new Error("missing code block");
if (branch[1] !== "val == 1" || branch[2] !== "end") throw new Error("branch metadata");
if (last[0] !== "end" || !last[3].includes("\u7ed3\u675f")) throw new Error("last section");
parse("a", { entrypoint: "predicateAnd" });
parse("b", { entrypoint: "predicateNot" });
parse("", { entrypoint: "codeProbe" });
parse("x", { entrypoint: "recoveryProbe" });
console.log("novel-ok");
`))

	outDir := filepath.Join(dir, "out")
	cmd := exec.Command(tsc,
		"--strict", "--target", "ES2020", "--module", "commonjs",
		"--outDir", outDir,
		filepath.Join(dir, "parser.ts"), filepath.Join(dir, "runner.ts"),
	)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("tsc failed: %v\n%s", err, output)
	}
	assertOutput(t, exec.Command(node, filepath.Join(outDir, "runner.js")), "novel-ok")
}

func testNovelCSharp(t *testing.T) {
	dotnet := requireTool(t, "dotnet")
	dir, generated := generateNovelParser(t, "cs")

	writeFile(t, filepath.Join(dir, "Parser.cs"), generated)
	writeFile(t, filepath.Join(dir, "story.ns"), readStory(t))
	writeFile(t, filepath.Join(dir, "Program.cs"), []byte(`
var source = File.ReadAllText("story.ns");
var sections = ((List<object?>)PTParser.Parse(source)!).Select(value => (string)value!).ToList();
if (sections.Count != 5) throw new Exception("section count: " + sections.Count);
var first = sections[0].Split('\u001f');
var branch = sections[1].Split('\u001f');
var last = sections[4].Split('\u001f');
if (first[0] != "\u5f00\u573a") throw new Exception("first section: " + first[0]);
if (!first[3].Contains("<codeInText>player")) throw new Exception("missing inline code");
if (!first[3].Contains("<codeBlock> val = select")) throw new Exception("missing code block");
if (branch[1] != "val == 1" || branch[2] != "end") throw new Exception("branch metadata");
if (last[0] != "end" || !last[3].Contains("\u7ed3\u675f")) throw new Exception("last section");
PTParser.Parse("a", new PTParser.ParserOptions { Entrypoint = "predicateAnd" });
PTParser.Parse("b", new PTParser.ParserOptions { Entrypoint = "predicateNot" });
PTParser.Parse("", new PTParser.ParserOptions { Entrypoint = "codeProbe" });
PTParser.Parse("x", new PTParser.ParserOptions { Entrypoint = "recoveryProbe" });
Console.WriteLine("novel-ok");
`))
	writeFile(t, filepath.Join(dir, "NovelTargetTest.csproj"), []byte(`
<Project Sdk="Microsoft.NET.Sdk">
  <PropertyGroup>
    <OutputType>Exe</OutputType>
    <TargetFramework>net8.0</TargetFramework>
    <ImplicitUsings>enable</ImplicitUsings>
    <Nullable>enable</Nullable>
  </PropertyGroup>
</Project>
`))

	project := filepath.Join(dir, "NovelTargetTest.csproj")
	cmd := exec.Command(dotnet, "build", project, "--configuration", "Release", "--nologo", "--verbosity", "quiet")
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"DOTNET_CLI_TELEMETRY_OPTOUT=1",
		"DOTNET_SKIP_FIRST_TIME_EXPERIENCE=1",
		"DOTNET_CLI_WORKLOAD_UPDATE_NOTIFY_DISABLE=true",
	)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("dotnet build failed: %v\n%s", err, output)
	}
	cmd = exec.Command(dotnet, filepath.Join(dir, "bin", "Release", "net8.0", "NovelTargetTest.dll"))
	cmd.Dir = dir
	assertOutput(t, cmd, "novel-ok")
}

func testNovelRust(t *testing.T) {
	rustc := requireTool(t, "rustc")
	dir, generated := generateNovelParser(t, "rust")

	var source bytes.Buffer
	source.Write(generated)
	source.WriteString(`
fn main() {
    let sections = match PTParser::parse(include_str!("story.ns")) {
        Ok(Value::Array(values)) => values,
        Ok(value) => panic!("unexpected result: {:?}", value),
        Err(error) => panic!("parse failed: {}", error),
    };
    assert_eq!(sections.len(), 5);
    let section = |index: usize| match &sections[index] {
        Value::String(value) => value.split('\u{1f}').collect::<Vec<_>>(),
        value => panic!("unexpected section: {:?}", value),
    };
    let first = section(0);
    let branch = section(1);
    let last = section(4);
    assert_eq!(first[0], "\u{5f00}\u{573a}");
    assert!(first[3].contains("<codeInText>player"));
    assert!(first[3].contains("<codeBlock> val = select"));
    assert_eq!(branch[1], "val == 1");
    assert_eq!(branch[2], "end");
    assert_eq!(last[0], "end");
    assert!(last[3].contains("\u{7ed3}\u{675f}"));
    for (entrypoint, input) in [
        ("predicateAnd", "a"),
        ("predicateNot", "b"),
        ("codeProbe", ""),
        ("recoveryProbe", "x"),
    ] {
        let options = ParserOptions {
            entrypoint: Some(entrypoint.to_string()),
            ..ParserOptions::default()
        };
        PTParser::parse_with_options(input, options).unwrap();
    }
    println!("novel-ok");
}
`)
	writeFile(t, filepath.Join(dir, "main.rs"), source.Bytes())
	writeFile(t, filepath.Join(dir, "story.ns"), readStory(t))
	binary := filepath.Join(dir, "novel-test.exe")
	cmd := exec.Command(rustc, "--edition=2021", filepath.Join(dir, "main.rs"), "-O", "-o", binary)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("rustc failed: %v\n%s", err, output)
	}
	assertOutput(t, exec.Command(binary), "novel-ok")
}

func generateNovelParser(t *testing.T, target string) (string, []byte) {
	t.Helper()
	goTool := requireTool(t, "go")
	grammarFile, err := filepath.Abs(filepath.Join("testdata", "novel.peg"))
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	outputFile := filepath.Join(dir, "generated.txt")
	cmd := exec.Command(
		goTool,
		"run", "../..",
		"-t", target,
		"-optimize-ref-expr-by-index",
		"-o", outputFile,
		grammarFile,
	)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("pegtool target %s failed: %v\n%s", target, err, output)
	}
	generated, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatal(err)
	}
	return dir, generated
}

func readStory(t *testing.T) []byte {
	t.Helper()
	value, err := os.ReadFile(filepath.Join("testdata", "story.ns"))
	if err != nil {
		t.Fatal(err)
	}
	return value
}

func quoteJavaScript(value []byte) string {
	return fmt.Sprintf("%q", string(value))
}

func requireTool(t *testing.T, name string) string {
	t.Helper()
	path, err := exec.LookPath(name)
	if err != nil {
		t.Skipf("%s is not installed", name)
	}
	return path
}

func writeFile(t *testing.T, path string, value []byte) {
	t.Helper()
	if err := os.WriteFile(path, value, 0o600); err != nil {
		t.Fatal(err)
	}
}

func assertOutput(t *testing.T, cmd *exec.Cmd, expected string) {
	t.Helper()
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\n%s", err, output)
	}
	if actual := strings.TrimSpace(string(output)); actual != expected {
		t.Fatalf("unexpected output %q", actual)
	}
}
