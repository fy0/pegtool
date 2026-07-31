<script setup>
import MigrationPrompt from "../../.vitepress/theme/components/MigrationPrompt.vue";
</script>

# Migrating from Upstream Pigeon

Pegtool is written in Go and forked from Pigeon in 2024, from a baseline that included Pigeon's 2023 `v1.2.x` release line. The two projects retain broadly similar PEG syntax, but Pegtool is not a compatibility layer for Pigeon. Its grammar, value semantics, generated APIs, and runtimes have evolved independently. Performance was the main reason for the fork, followed by extensive optimization. Some of that work required changes to upstream Pigeon's implicit value propagation semantics. See [Pigeon #151](https://github.com/mna/pigeon/issues/151) for background.

Most migrations can be completed by a code agent with repository read/write access and the ability to run tests. The prompt below automatically includes the full URL of this page and can be given directly to an agent:

<MigrationPrompt />

## Migration differences at a glance

| Area | Upstream Pigeon `v1.2.x` | Current Pegtool behavior |
| --- | --- | --- |
| Terminal values | Literals, character classes, and `.` return `[]byte` | Return the target language's null value; capture text with `label:<expr>` |
| Sequence values | Return `[]any` with one item for each child expression | Return null and do not create implicit arrays |
| Repetition values | Collect the result of every match, including null values | Collect only non-null values explicitly returned by the child expression |
| Go actions | Return `(value, error)` | Return one value; record errors with `p.addErr(err)` |
| Go predicates | Return `(bool, error)` | Return `bool`; errors must be recorded explicitly in code |
| Generated Go API | Exports `Parse`, `ParseFile`, `ParseReader`, and fixed options | Base symbols are package-private and wrapped by the initializer or regular source files |
| Parser state | `state`, `globalStore`, and `#{...}` | Uses `ParserCustomData` / `c.data`, without an equivalent automatically rolled-back state layer |
| Left recursion | Available behind an experimental option | Rejected by the builder |
| Generation targets | Primarily Go | Go, Haxe, TypeScript, C#, C99, and Rust |

This table is a migration checklist, not a list of mechanical replacements. State, error handling, and public APIs in particular must be rewritten according to the existing project's semantics.

## Text and values are no longer implicit

Upstream Pigeon creates implicit values for terminals and sequences. For example:

```peg
// Upstream Pigeon
Word <- [a-z]+ {
    return string(c.text), nil
}

Pair <- pair:("a" ":" "b") {
    // The underlying type of pair is []any.
    return pair, nil
}
```

In Pegtool, state explicitly whether you need source text or a structured value:

```peg
// Pegtool
Word = value:<[a-z]+> {
    return value
}

Pair = pair:<("a" ":" "b")> {
    // pair is the matched string.
    return pair
}
```

Migrating `pair:(...)` unchanged produces a null value. For a structured result, return an object explicitly from an inner action and receive it through a label. Do not rely on the old sequence's nested `[]any` values.

Repetitions also need review. Upstream preserves an item for every match; Pegtool collects only non-null values explicitly returned by a child:

```peg
// Pegtool: Name explicitly returns a string, so Names can collect the values.
Name  = value:<[a-z]+> { return value }
Names = first:Name rest:(_ "," _ next:Name { return next })* {
    values := []any{first}
    if rest != nil {
        values = append(values, rest.([]any)...)
    }
    return values
}
```

## Go actions and predicates

Upstream Pigeon actions and predicates both report errors through a second return value:

```peg
// Upstream Pigeon
Integer <- [0-9]+ {
    return strconv.Atoi(string(c.text))
}

Allowed <- &{
    return c.globalStore["enabled"].(bool), nil
} .
```

Pegtool's Go code blocks return one value. Record parse errors through the current parser:

```peg
// Pegtool
Integer = [0-9]+ {
    value, err := strconv.Atoi(string(c.text))
    if err != nil {
        p.addErr(err)
        return nil
    }
    return value
}

Allowed = &{
    return c.data.Enabled
} .
```

Do not simply remove the `error` result while migrating a predicate. If the old code can fail, call `p.addErr(err)` first and then return the appropriate Boolean result.

## The Go API needs explicit wrappers

Upstream generated output exports fixed APIs such as `Parse`, `ParseFile`, `ParseReader`, `Entrypoint`, and `Memoize`. Pegtool's Go runtime provides only the package-private `parse`, `newParser`, `option`, and `memoized` symbols. Preserve the project's public surface in the initializer:

```peg
{
package parser

type ParserCustomData struct{}
type Option = option

func Memoize(enabled bool) Option {
    return memoized(enabled)
}

func Parse(filename string, input []byte, opts ...Option) (any, error) {
    return parse(filename, input, opts...)
}
}

Start = value:<[a-z]+> !. { return value }
```

If `ParseFile` and `ParseReader` are part of the public API, reimplement them in a regular Go source file or the initializer and cover them with tests. A Go grammar must declare an empty `ParserCustomData` even when it has no custom state.

## State cannot be replaced mechanically

Upstream `state` is restored during PEG backtracking, `globalStore` is not, and `#{...}` has dedicated state-mutation semantics. Pegtool's `c.data` is a caller-provided `*ParserCustomData`; modifications are not automatically rolled back. There is no general rewrite that replaces all three old concepts with `c.data` while preserving behavior.

Classify each field before migrating it:

1. Read-only configuration can live in `ParserCustomData`.
2. Accumulated state that does not need rollback can modify `c.data` explicitly.
3. State that must be restored when a choice fails should be saved and restored by the grammar, or replaced with immutable returned results.
4. A memo hit does not rerun actions or predicates inside a rule. Disable memoization or redefine the parse lifecycle if configuration changes during one parse.

Pegtool also provides `label:<expr>` text capture, `&&` / `!!` consuming assertions, and `*{...}` code expressions that run in lookahead. These are features of the current grammar and must not be copied back into grammars still generated by upstream Pigeon.

## Left recursion and entry points

Upstream Pigeon can enable direct or indirect left recursion with an experimental option. Pegtool rejects it at generation time. Rewrite it as a head followed by repeated tails:

```peg
// Optional left-recursive form in upstream Pigeon
Expr <- Expr "+" Term / Term

// Pegtool
Expr = first:Term rest:(_ "+" _ next:Term { return next })* {
    return foldAdd(first, rest)
}
```

Select entry rules through the generated API for each target or through a project wrapper. Do not assume upstream Pigeon's `Entrypoint(...)` option exists for every target.

## Verify the migration

At minimum, complete these checks after migrating:

1. Validate the grammar with `pegtool -x grammar.peg`.
2. Regenerate the parser for the project's actual target. Do not hand-edit generated files to hide problems.
3. Compile the generated output and run the project's complete test suite.
4. Add coverage for sequences, repetitions, action errors, predicates, state rollback, Unicode, and invalid input.
5. Compare the existing public API and retain compatibility wrappers where needed.
6. Enable memoization or rule index optimization only when benchmarks on real workloads show a benefit.

See [PEG Syntax](/en/guide/grammar) for grammar details, [Actions and Values](/en/guide/actions-and-values) for value propagation, and [Generated API](/en/reference/generated-api) for target-specific invocation.
