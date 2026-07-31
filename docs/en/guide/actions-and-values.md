# Actions and Values

Pegtool gives matched text, semantic values, and host-language actions distinct roles. A grammar no longer creates a value array for every sequence, so any required data must be captured or returned explicitly.

## Value propagation

| Expression | Return value |
| --- | --- |
| Literal, character class, `.` | `nil` / `null` / the target language's null value |
| `label:<expr>` | Matched source text as a string |
| Ordinary sequence | Null value |
| Choice | Value of the selected alternative |
| `expr?` | Child value or null |
| `expr*`, `expr+` | Collection of non-null values explicitly returned by the child; null if there are no values |
| `&expr`, `!expr` | Null value |
| Action | Any value explicitly returned by the action |

## Capturing text

```peg
Name       = value:<[\pL_] [\pL\p{Nd}_]*> { return value }
Assignment = source:<(Name _ "=" _ Name)> { return source }
```

`value:Name` receives the semantic value returned by `Name`; `value:<Name>` receives the original text consumed by `Name`. Terminals do not return characters or bytes themselves.

## Collecting values from repetition

This expectation is incorrect:

```peg
// The terminals in [0-9]+ return no values, so digits is null.
Number = digits:[0-9]+
```

Capture the complete matched text instead:

```peg
Number = digits:<[0-9]+> { return digits }
```

To collect multiple structured results, make an action inside the repetition return each required value:

```peg
Name  = value:<[a-z]+> { return value }
Names = first:Name rest:(_ "," _ next:Name { return next })* {
  values := []any{first}
  if rest != nil {
    values = append(values, rest.([]any)...)
  }
  return values
}
```

This example uses a Go action. Labels inside the repetition are not automatically promoted to the outer scope. Returning required data from the inner action avoids implicit nested arrays and large numbers of useless null values.

## Go actions

A Go action returns one `any` value. Record errors through the parser:

```peg
Integer = [0-9]+ {
    value, err := strconv.Atoi(string(c.text))
    if err != nil {
        p.addErr(err)
        return nil
    }
    return value
}
```

The following objects are available:

- `c.text`: the complete `[]byte` matched by the action's expression.
- `c.pos`: the position where the match started.
- `c.data`: `*ParserCustomData`.
- `p`: the current parser, which exposes `p.addErr` and any internal capabilities wrapped by the project.

Do not use the `return value, err` form from other Pigeon versions.

## Action values in other targets

| Target | Label/return value model | Current-match object |
| --- | --- | --- |
| Haxe | `Any` | `c.text`, `c.data` |
| TypeScript | `any` | `c.text`, `c.data` |
| C# | `object?` | `c.Text`, `c.Data` |
| C99 | Owned `PigeonValue` model | `PigeonCurrent *c` |
| Rust | `Value` enum | `current.text` and custom data |

C99 actions return an owned `PigeonValue`; Rust actions use a cloneable `Value` enum. See each [target language](/en/targets/) page for exact usage.

## Actions inside lookahead

Ordinary actions are skipped while probing through `&expr` or `!expr`, preventing probes from causing application side effects. `*{...}` explicitly runs even inside lookahead:

```peg
Probe = *{ c.data.ProbeCount++ }
```

Do not rely on creating labels inside syntactic lookahead and propagating them outward. Put the required condition in a semantic predicate, or capture the data in a branch that actually consumes input.

## Actions, side effects, and memoization

A memo hit does not rerun actions or predicates inside the cached rule. Enable memoization only when all of the following are true:

- Rule results are determined primarily by the input position.
- Configuration read by predicates stays constant during a parse.
- Action side effects do not need to occur again on every backtracking visit.

Disable memoization when predicates use changing `c.data` to control matching.
