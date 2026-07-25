package builder

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

func TestIndentWriter(t *testing.T) {
	var dst bytes.Buffer
	w := newIndentWriter(&dst, "\t")

	_, _ = io.WriteString(w, "root\n")
	w.indent()
	_, _ = io.WriteString(w, "field: ")
	_, _ = io.WriteString(w, "value\n\nnext\n")
	w.indent()
	_, _ = io.WriteString(w, "nested\n")
	if err := w.dedent(); err != nil {
		t.Fatal(err)
	}
	if err := w.dedent(); err != nil {
		t.Fatal(err)
	}
	_, _ = io.WriteString(w, "done\n")

	want := "root\n\tfield: value\n\n\tnext\n\t\tnested\ndone\n"
	if got := dst.String(); got != want {
		t.Fatalf("unexpected indented output:\nwant: %q\n got: %q", want, got)
	}
}

func TestIndentWriterRejectsUnderflow(t *testing.T) {
	w := newIndentWriter(io.Discard, "\t")
	if err := w.dedent(); !errors.Is(err, errIndentUnderflow) {
		t.Fatalf("want %v, got %v", errIndentUnderflow, err)
	}
}

func TestBuilderWritesFieldsAcrossPartialLines(t *testing.T) {
	var dst bytes.Buffer
	b := Builder{W: &dst, IndentUnit: "\t"}
	b.initIndentWriter()

	b.Writef("Node({")
	b.Indent()
	b.WriteFieldf("first: 1,")
	b.WriteFieldf("second: 2,")
	b.Dedent()
	b.WriteSpace()
	b.Writeln("})")
	b.Indent()
	b.WriteFieldf("nested: ")
	b.Writef("value\n")
	b.Dedent()

	want := "Node({ first: 1, second: 2, })\n\tnested: value\n"
	if got := dst.String(); got != want {
		t.Fatalf("unexpected field output:\nwant: %q\n got: %q", want, got)
	}
}
