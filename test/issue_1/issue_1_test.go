package issue1

import "testing"

func TestOptionalLabelScope(t *testing.T) {
	for _, tc := range []struct {
		input string
		want  string
	}{
		{input: "foo", want: "foo"},
		{input: "db.foo", want: "db.foo"},
	} {
		t.Run(tc.input, func(t *testing.T) {
			p := newParser("", []byte(tc.input))
			got, err := p.parse(nil)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.want {
				t.Errorf("want %q, got %#v", tc.want, got)
			}
		})
	}
}
