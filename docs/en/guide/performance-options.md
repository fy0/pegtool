# Choosing Performance Options

Choose performance options according to grammar behavior instead of enabling them all. Pegtool's three most important tradeoffs are a fixed optimized runtime, named-rule memoization, and indexed rule references. They solve different problems and carry different maintenance costs.

## Recommended starting points

When generated files participate in Git review, start with these stable configurations:

| Scenario | Recommended generation options | Runtime |
| --- | --- | --- |
| Typical production Go grammar | `-t go -optimize-parser` | No memoization by default |
| Go grammar with substantial repeated backtracking | `-t go -optimize-parser` | Use `memoized(true)` after benchmarking |
| Typical, near-linear Haxe grammar | `-t hx -optimize-parser` | This output removes memoization |
| Haxe grammar with substantial repeated backtracking | `-t hx` | Regular output can set `parser.memoize = true` |
| TypeScript | `-t ts` | No memoization currently |
| C#, C99, Rust | Select the relevant target only | No memoization currently |

Complete generation commands:

```powershell
# Go: use the same generation command for ordinary and backtracking-heavy
# grammars; enable memoization at runtime for the latter.
pegtool -t go -optimize-parser -o parser.go grammar.peg

# Haxe: ordinary, near-linear grammar; output contains no memoization.
pegtool -t hx -optimize-parser -o Parser.hx grammar.peg

# Haxe: backtracking-heavy grammar; retain memoization support and set
# parser.memoize = true before parsing.
pegtool -t hx -o Parser.hx grammar.peg

# TypeScript
pegtool -t ts -o parser.ts grammar.peg

# C#
pegtool -t cs -o Parser.cs grammar.peg

# C99
pegtool -t c -o parser.c grammar.peg

# Rust
pegtool -t rust -o parser.rs grammar.peg
```

Haxe has two distinct output forms; one generated file does not simultaneously have and omit memoization:

- With `-optimize-parser`, output is smaller and the memo fields, functions, and cache branches are removed completely.
- Without `-optimize-parser`, output retains optional memoization, which is off by default and activates only after setting `parser.memoize = true`.

Do not add `-optimize-ref-expr-by-index` by default. It may improve throughput, but inserting or reordering a rule can produce a large generated diff.

## Decision order

1. Validate semantics and error output with the default configuration.
2. For a fixed production grammar, use `-optimize-parser` where the target meaningfully supports it.
3. Determine whether the same named rule is evaluated repeatedly at the same input position.
4. Test memoization only when step 3 is true.
5. Test rule indexing only when you still need the last increment of throughput and can accept large diffs.
6. Rerun tests on real input in interleaved and reversed order; do not decide from one absolute timing.

## `-optimize-parser`

This option is not one uniform switch across targets:

| Target | Current effect |
| --- | --- |
| Go | Removes Debug and Statistics hot paths and reduces unnecessary variable-stack operations; named-rule memoization remains available |
| Haxe | Removes Debug, Statistics, memo fields, and cache branches, and flattens some wrappers |
| TypeScript | Currently ignores the option and produces byte-for-byte identical output |
| C#, C99, Rust | No target-specific runtime branch currently |

Go can therefore combine optimized generation and memoization:

```powershell
pegtool -t go -optimize-parser -o parser.go grammar.peg
```

```go
value, err := Parse("input.txt", input, Memoize(true))
```

Haxe requires choosing one form:

```powershell
# Fixed, near-linear grammar
pegtool -t hx -optimize-parser -o Parser.hx grammar.peg

# Backtracking-heavy grammar that needs memoization
pegtool -t hx -o Parser.hx grammar.peg
```

## Named-rule memoization

Memoization caches successful and failed results by named rule, input position, and parse mode. It prevents repeated scans during backtracking, but every first visit requires a lookup and a stored result.

Good candidates include:

- Several ordered alternatives share an expensive prefix rule.
- The same failing rule scans a long input repeatedly from the same offset.
- Predicate input configuration remains stable during a parse.

Poor candidates include:

- The grammar is nearly linear and produces no useful cache hits.
- A repeatedly called rule matches only a simple literal, making recomputation cheaper than a cache lookup.
- An action or predicate depends on `c.data` that changes during backtracking.
- Side effects must occur again on every rule visit.

In tests, memoization made an expensive repeated-scan HashLink workload `2.53x-2.61x` faster. It made a linear workload about `1%-1.4%` slower and an exploratory workload that cached a simple failing rule about `24%` slower. Backtracking alone is not enough; the repeated work must be more expensive than the cache operations.

In a real Go DiceScript test, narrowing the cache boundary from every expression to named rules reduced parser-related execution in the complete root test suite from about `199.8 ms` to `58.3 ms`. Most of the improvement came from fewer cache entries and allocations, not from eliminating `any` type dispatch.

## `-optimize-ref-expr-by-index`

By default, rule references are looked up by name. Index mode changes this to:

```text
rules["Expression"]  ->  rulesArray[42]
```

This avoids a name-map lookup, but the number depends directly on rule order. Inserting a rule at position 10 can change every later rule number and every reference to those rules.

### When to enable it

- Generated files are not committed or manually reviewed.
- Grammar rule order is effectively frozen.
- A benchmark on the target workload shows a benefit large enough to cover the maintenance cost.
- The runtime does not reorder or concatenate `grammar.rules`.

### When to leave it off

- Small grammar changes should produce local diffs.
- Multiple contributors frequently insert or reorder rules.
- Generated parsers require sustained auditing.

The historical estimate of roughly `10%` was an overall figure, not a current guarantee for Go, Haxe, or TypeScript. The current cross-target tests did not isolate this option either, so Pegtool does not document it as a default.

For maximum throughput, append this only after establishing a stable configuration:

```text
-optimize-ref-expr-by-index
```

## Other commonly confused options

### `-cache`

This caches only the process by which **Pegtool reads a `.peg` grammar**. It does not change the generated file or enable memoization in the generated parser. Test it only when the generator is unusually slow at parsing a pathological grammar; ordinary grammars may become slower and consume more memory.

### `-nolint`

This controls lint-suppression comments in generated Go code only. It does not affect parser throughput or generation algorithms. Haxe and TypeScript currently ignore it.

### `-haxe-use-hxunicode`

This is a Haxe source-size and dependency choice, not a performance switch. The default embeds the generator's Unicode tables. External mode requires `hxUnicode` and produces a smaller file, but the character-set version can change with the dependency.

### `-alternate-entrypoints`

This currently verifies rule names during generation only. It neither changes output nor creates a runtime selector. Select an entry point with target-specific `ParserOptions` or a project wrapper.

## JavaScript and HashLink runtime options

- Use the default V8 JIT for production Node.js execution; do not add `--jitless`.
- Use `node --jitless` only to confirm that an optimization is not solely a warm-up artifact. It is not an acceleration option.
- tjs uses the QuickJS path and is useful for compatibility testing without V8. It does not require different generation options.
- HashLink should use normal HL/JIT. There is currently no corresponding no-JIT switch.
- When compiling Haxe to JavaScript or HashLink, use `-D analyzer-optimize`.

In a cross-target backtracking workload, the Haxe memoized version took about `74%-77%` of the time used by direct TypeScript, retaining the advantage under Node JIT, Node jitless, and tjs. A linear workload should prefer Haxe `-optimize-parser` instead of paying for a cache that does not hit.

## Benchmark your own grammar

A reliable comparison should do at least the following:

1. Verify that every variant returns the same values, errors, and consumed positions.
2. Use short, typical, and worst-case backtracking input.
3. Warm up each variant independently and calibrate batch counts automatically.
4. Alternate A/B runs, then use a separate process to reverse them to B/A.
5. Record throughput, allocation counts, and memory, not only one wall-clock time.
6. Include one `--jitless` control for Node results; report QuickJS and HL separately.

Machine load can change absolute timings substantially. Prefer adjacent paired ratios within the same process and stable measurements such as allocation counts.
