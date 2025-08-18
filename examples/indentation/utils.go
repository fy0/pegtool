package main

import (
	"io"
	"os"
)

type Option func(*parser) Option

// ParseReader parses the data from r using filename as information in the
// error messages.
func ParseReader(filename string, r io.Reader, opts ...Option) (any, error) { // nolint: deadcode
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return Parse(filename, b, opts...)
}

func Parse(input string, data []byte, opts ...Option) (any, error) {
	p := newParser(input, data)
	for _, opt := range opts {
		opt(p)
	}
	return p.parse(nil)
}

// ParseFile parses the file identified by filename.
func ParseFile(filename string, opts ...Option) (i any, err error) { // nolint: deadcode
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			err = closeErr
		}
	}()
	return ParseReader(filename, f, opts...)
}
