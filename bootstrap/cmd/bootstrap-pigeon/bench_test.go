package main

import (
	"os"
	"testing"
)

func Parse(filename string, b []byte, opts ...option) (any, error) {
	return newParser(filename, b, opts...).parse(g)
}

func BenchmarkParsePigeonNoMemo(b *testing.B) {
	d, err := os.ReadFile("../../../grammar/pigeon.peg")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := Parse("", d, memoized(false)); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParsePigeonMemo(b *testing.B) {
	d, err := os.ReadFile("../../../grammar/pigeon.peg")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := Parse("", d, memoized(true)); err != nil {
			b.Fatal(err)
		}
	}
}
