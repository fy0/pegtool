package builder

import "testing"

func TestUnicodeRanges(t *testing.T) {
	ranges, err := UnicodeRanges([]string{"Latin", "Other_ID_Start"})
	if err != nil {
		t.Fatal(err)
	}
	if len(ranges) == 0 {
		t.Fatal("expected expanded Unicode ranges")
	}

	if _, err := UnicodeRanges([]string{"not-a-unicode-class"}); err == nil {
		t.Fatal("expected an error for an unknown Unicode class")
	}
}
