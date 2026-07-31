# Error Handling

By default, a PEG can report only which matchers were expected at the farthest failure position. After a complex grammar passes through many choices and backtracking paths, that list often reflects the parser's internal route more than an error a user can understand. Replacing such messages in upstream Pigeon generally requires predicates or recovery expressions throughout the grammar, or parsing internal errors after a parse, all of which are expensive to maintain.

pegtool preserves position and `expected` data while adding a parser-level no-match formatter for Go and Haxe. Applications can produce user-facing messages without modifying the generated runtime or scattering error text throughout the grammar.

## Go: `noMatchErrorFormatter`

The Go target's formatter signature is:

```go
func(position, []byte, []string) error
```

The arguments are the farthest failure position, the complete input, and a deduplicated, sorted list of expected items. `!.` is converted to `EOF` first. When the formatter returns a non-null error, the parser uses it instead of the default `no match found, expected: ...` message. Returning `nil` preserves the default error.

This initializer installs an English formatter in the generated parser:

```peg
{
package parser

import (
    "fmt"
    "strings"
    "unicode/utf8"
)

type ParserCustomData struct{}

func friendlyNoMatch(pos position, input []byte, expected []string) error {
    got := "EOF"
    if pos.offset < len(input) {
        value, _ := utf8.DecodeRune(input[pos.offset:])
        got = fmt.Sprintf("%q", value)
    }

    wanted := strings.Join(expected, ", ")
    if wanted == "" {
        wanted = "valid input"
    }
    return fmt.Errorf("line %d, column %d: got %s; expected %s", pos.line, pos.col, got, wanted)
}

func Parse(filename string, input []byte) (any, error) {
    return parse(filename, input, noMatchErrorFormatter(friendlyNoMatch))
}
}

Start "identifier" = value:<Identifier> !. { return value }
Identifier = [\pL_] [\pL\p{Nd}_]*
```

`noMatchErrorFormatter(...)` is an option for one parser instance and does not modify a global variable. A wrapper can pass a different formatter for each parse according to the caller's language or error presentation. The interface remains available in Go output generated with `-optimize-parser`.

## Haxe and other targets

Haxe provides the same parser option with a `String` input:

```haxe
var parser = new Parser("input.txt", source, [
    noMatchErrorFormatter((position, input, expected) ->
        new Exception('${position.line}:${position.col}: expected ${expected.join(", ")}')
    )
]);
```

The Haxe formatter also remains available in output generated with `-optimize-parser`. Other targets currently expose structured diagnostic data:

| Target | Available data |
| --- | --- |
| TypeScript | `PTParseError.position`, `expected` |
| C# | `PTParseException.Position`, `Expected` |
| Rust | `ParseError.position`, `expected` |
| C99 | `PegtoolResult.error`, `position` |

Applications can turn these structured errors into localized messages at their boundaries. These targets do not currently provide the same formatter option as Go and Haxe.

## Semantic errors and display names

The no-match formatter handles only the case where the grammar cannot match. Semantic errors such as integer overflow or forbidden keywords should be recorded explicitly in an action or predicate. A Go action can call `p.addErr(err)`:

```peg
Integer "integer" = [0-9]+ {
    value, err := strconv.Atoi(string(c.text))
    if err != nil {
        p.addErr(fmt.Errorf("integer out of range: %w", err))
        return nil
    }
    return value
}
```

The string after the rule name is its diagnostic display name. In the Go and Haxe runtimes, errors recorded inside the rule show `integer` instead of its implementation name, `Integer`. A display name does not replace the no-match formatter.

To consume invalid input after an error and continue constructing a result, use failure labels and recovery as described below. When only the error text needs to change, prefer a formatter or structured error rather than introducing a recovery branch.

## Failure labels and recovery

`%{Label}` throws a failure label, while `//{...}` installs a recovery expression for that label:

```peg
Value = Number / %{ExpectedNumber}
        //{ExpectedNumber} InvalidNumber
```

The behavior is:

1. Match the normal expression on the left side of `Value`.
2. `%{ExpectedNumber}` searches outward for a recovery that handles the label.
3. Once found, run the recovery expression from the throw position.
4. Continue parsing if recovery succeeds; if it fails or no handler exists, continue propagating the failure.

Custom data is not rolled back automatically. A grammar that modifies `c.data` during recovery is responsible for restoring a consistent state. Failure-label recovery and host-language panic/exception recovery are separate mechanisms.

Error handling commonly involves both actions and the public parser API. See [Actions and Values](/en/guide/actions-and-values) for value semantics and [Generated API](/en/reference/generated-api) for each target's error type.
