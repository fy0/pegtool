package builder

import (
	"fmt"
	"unicode"
)

// UnicodeRange is a compact, target-neutral representation of a Unicode
// range table entry. Stride is the distance between matching code points.
type UnicodeRange struct {
	Lo     rune
	Hi     rune
	Stride rune
}

// UnicodeRanges expands category, property, and script names understood by
// Go's unicode package. Target builders use this to stay self-contained.
func UnicodeRanges(names []string) ([]UnicodeRange, error) {
	var ranges []UnicodeRange
	for _, name := range names {
		table := unicode.Categories[name]
		if table == nil {
			table = unicode.Properties[name]
		}
		if table == nil {
			table = unicode.Scripts[name]
		}
		if table == nil {
			return nil, fmt.Errorf("unknown Unicode class %q", name)
		}

		for _, item := range table.R16 {
			ranges = append(ranges, UnicodeRange{
				Lo:     rune(item.Lo),
				Hi:     rune(item.Hi),
				Stride: rune(item.Stride),
			})
		}
		for _, item := range table.R32 {
			ranges = append(ranges, UnicodeRange{
				Lo:     rune(item.Lo),
				Hi:     rune(item.Hi),
				Stride: rune(item.Stride),
			})
		}
	}
	return ranges, nil
}
