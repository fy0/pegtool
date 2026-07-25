package builder

import (
	"bytes"
	"errors"
	"io"
)

var errIndentUnderflow = errors.New("builder: indentation depth below zero")

type indentWriter struct {
	dst         io.Writer
	unit        []byte
	depth       int
	atLineStart bool
	lastByte    byte
	hasLastByte bool
}

func newIndentWriter(dst io.Writer, unit string) *indentWriter {
	return &indentWriter{
		dst:         dst,
		unit:        []byte(unit),
		atLineStart: true,
	}
}

func (w *indentWriter) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	if w.depth == 0 || len(w.unit) == 0 || (!w.atLineStart && !bytes.ContainsRune(p, '\n')) {
		n, err := w.dst.Write(p)
		if n > 0 {
			w.atLineStart = p[n-1] == '\n'
			w.lastByte = p[n-1]
			w.hasLastByte = true
		}
		return n, err
	}

	var out bytes.Buffer
	out.Grow(len(p) + w.depth*len(w.unit))
	for _, ch := range p {
		if w.atLineStart && ch != '\n' {
			for range w.depth {
				_, _ = out.Write(w.unit)
			}
			w.atLineStart = false
		}

		_ = out.WriteByte(ch)
		w.lastByte = ch
		w.hasLastByte = true
		if ch == '\n' {
			w.atLineStart = true
		}
	}

	n, err := w.dst.Write(out.Bytes())
	if err != nil {
		return 0, err
	}
	if n != out.Len() {
		return 0, io.ErrShortWrite
	}
	return len(p), nil
}

func (w *indentWriter) indent() {
	w.depth++
}

func (w *indentWriter) dedent() error {
	if w.depth == 0 {
		return errIndentUnderflow
	}
	w.depth--
	return nil
}

func (b *Builder) initIndentWriter() {
	if w, ok := b.W.(*indentWriter); ok {
		b.indentWriter = w
		return
	}

	b.indentWriter = newIndentWriter(b.W, b.IndentUnit)
	b.W = b.indentWriter
}

// Indent increases indentation for subsequent lines.
func (b *Builder) Indent() {
	b.indentWriter.indent()
}

// Dedent decreases indentation for subsequent lines.
func (b *Builder) Dedent() {
	if err := b.indentWriter.dedent(); err != nil && b.Err == nil {
		b.Err = err
	}
}

func (b *Builder) writeInlineSeparator() {
	w := b.indentWriter
	if w.atLineStart || !w.hasLastByte {
		return
	}
	if w.lastByte != '{' && w.lastByte != '[' && w.lastByte != ' ' && w.lastByte != '\t' {
		b.Writef(" ")
	}
}

// WriteFieldf writes a field, separating it from preceding content on the
// same line while leaving line-start indentation to the writer.
func (b *Builder) WriteFieldf(f string, args ...any) {
	b.writeInlineSeparator()
	b.Writef(f, args...)
}

// WritelnFieldf writes a field followed by a newline.
func (b *Builder) WritelnFieldf(f string, args ...any) {
	b.writeInlineSeparator()
	b.Writelnf(f, args...)
}

// WriteItemf writes an inline or multiline collection item.
func (b *Builder) WriteItemf(f string, args ...any) {
	b.writeInlineSeparator()
	b.Writef(f, args...)
}

// WritelnItemf writes a collection item followed by a newline.
func (b *Builder) WritelnItemf(f string, args ...any) {
	b.writeInlineSeparator()
	b.Writelnf(f, args...)
}
