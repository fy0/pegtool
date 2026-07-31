# PEG Syntax

A PEG describes input with recognition rules. Unlike regular expressions, rules can refer to one another recursively. Unlike context-free grammars, choices are ordered: the first successful alternative wins immediately.

## Grammar structure

A grammar contains an optional initializer and at least one rule:

```peg
{
// Target-language code goes here.
}

Start "start rule" <- _ value:Expr EOF { return value }
Expr               <- Number / "(" _ Expr _ ")"
Number             <- [0-9]+
_                  <- [ \t\r\n]*
EOF                <- !.
```

The first rule is the default entry point. A string after the rule name provides its display name in diagnostics. Rule definitions may use `=`, `<-`, `←`, or `⟵`.

Ordinary whitespace is ignored. Line comments use `//` and block comments use `/* ... */`. Because `//{` starts failure recovery syntax, do not use it as the beginning of an ordinary comment.

## Expressions

| Syntax | Meaning |
| --- | --- |
| `A B` | Match A followed by B |
| `A / B` | Try A first; if it fails, return to the starting position and try B |
| `(A / B)` | Group expressions |
| `A?` | Match zero or one time |
| `A*` | Match zero or more times, greedily |
| `A+` | Match one or more times, greedily |
| `&A` | Succeed if A succeeds without consuming input |
| `!A` | Succeed if A fails without consuming input |
| `&&A` | Succeed if A succeeds and actually consumes input |
| `!!A` | Logical inverse extension of `&&A` |
| `label:A` | Bind A's semantic value to `label` |
| `label:<A>` | Capture A's matched source text as a string |

### Ordered choice

Order changes the accepted language:

```peg
BadOperator  = "<" / "<="
GoodOperator = "<=" / "<"
```

`BadOperator` can never match `<=` through its complete second alternative because the first alternative has already accepted `<`. Put the more specific alternative first, and place `!.` at the end of the entry rule to reject unconsumed trailing input.

### Repetition must advance

The child of `*` or `+` should consume at least one character or fail. Avoid nullable repetitions such as:

```peg
// Incorrect: ("a"?)* can succeed without advancing.
Loop = ("a"?)*
```

## Literals

```peg
Keyword = "select"i
Rune    = 'x'
Raw     = `no escapes here`
```

Double quotes delimit strings, single quotes delimit exactly one character, and backticks delimit raw strings. An `i` after the closing quote enables case-insensitive matching.

## Character classes and Unicode

```peg
Digit       = [0-9]
Identifier  = [\pL_] [\pL\p{Nd}_]*
GreekLetter = [\p{Greek}]
NotNewline  = [^\r\n]
Hex         = [0-9a-f]i
```

Character classes support individual characters, ranges, escapes, negation with a leading `^`, and Unicode classes:

- `\pL` uses a one-letter general category.
- `\p{Nd}` uses a named category.
- `\p{Greek}` uses a Unicode script or property name.

Available names come from the current generator's Go Unicode data. Generated output retains the category semantics from generation time rather than relying on a browser's implementation of `\p{...}` regular expressions. Haxe's external `hxUnicode` mode is the exception; see [Unicode tables for the Haxe target](/en/targets/haxe#unicode-tables).

`.` matches one Unicode character, but not EOF, so the idiomatic end rule is:

```peg
EOF = !.
```

## Labels, captures, and actions

`label:expr` stores the expression's semantic value, while `label:<expr>` stores its matched text. Group a complex expression when capturing all of its text:

```peg
Assignment = text:<(Name _ "=" _ Value)> { return text }
```

See [Actions and Values](/en/guide/actions-and-values) for return types and repetition collection rules.

## Semantic predicates and code expressions

```peg
Allowed = &{ return c.data.AllowValue } Value
Blocked = !{ return c.data.Blocked } .
Probe   = *{ c.data.ProbeCount++ }
```

- `&{...}` succeeds when the code returns true.
- `!{...}` inverts the code result.
- Ordinary actions are skipped inside syntactic lookahead.
- `*{...}` is a code expression that runs even inside lookahead.

Code contents and return types depend on the target.

## Error handling

Rule display names, no-match formatters, action-level semantic errors, and `%{Label}` / `//{...}` recovery syntax are covered in [Error Handling](/en/guide/error-handling).

## Left recursion

The current pegtool builder rejects left recursion:

```peg
// Not supported.
Expr = Expr "+" Term / Term
```

Rewrite it to match a head followed by repeated tails:

```peg
Expr = Term (_ "+" _ Term)*
```
