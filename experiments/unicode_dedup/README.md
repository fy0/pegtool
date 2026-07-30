# Unicode Range Size Experiment and Decision

Status: final design selected after comparing shared tables, matcher pooling, and
generic expression pooling.

Baseline commit: `c09ac57`
Working branch: `hx`
Date: 2026-07-30

## Problem

The old TypeScript target emitted one object for every Go `unicode.RangeTable`
range used by every character-class matcher:

```ts
{ lo: 65, hi: 90, stride: 1 }
```

NovelScript therefore contained 1,801 range objects. Of those ranges, 626 were
repeated across Unicode classes: 623 from `L` and 3 from `Other_ID_Start`. The
union of the per-class tables contained 1,175 ranges. Besides duplicated data,
the property names and object syntax were repeated for every range.

## Provisional shared-table design

The first implementation pre-scanned the grammar, assigned stable encoded names
to Unicode classes, and emitted one flat table per class. TypeScript matchers
then referred to those constants. Equivalent target-specific sharing was also
prototyped for Haxe, C#, and Rust; Go was deliberately left unchanged.

NovelScript source sizes for that provisional implementation were:

| Target | Baseline source | Provisional source | Change |
| --- | ---: | ---: | ---: |
| TypeScript | 144,077 B | 83,054 B | -42.4% |
| Haxe | 97,422 B | 85,405 B | -12.3% |
| C# | 200,930 B | 147,104 B | -26.8% |
| Rust | 209,515 B | 146,733 B | -30.0% |

With local esbuild 0.27.0 and zlib level 9, the provisional TypeScript output
measured 83,054 B source, 37,923 B minified, and 12,292 B minified plus gzip.
The locally reproduced old output measured 73,282 B minified and 14,035 B gzip;
the acceptance baseline supplied for this work is 73,281 B and 13,963 B.

Warmup, alternating old/new order, and paired samples were used for performance
checks because workstation load is not stable. The provisional implementation
was faster in the TypeScript Unicode-heavy microbenchmark, while the other
targets were effectively neutral:

| Workload | Baseline median | Provisional median | New / old |
| --- | ---: | ---: | ---: |
| TypeScript NovelScript input | 1.842 ms | 1.604 ms | 0.859 |
| TypeScript Unicode micro | 16.915 ms | 10.850 ms | 0.641 |
| Haxe/JavaScript Unicode micro | 0.272 ms | 0.277 ms | 1.008 |
| C# Unicode micro | 1.804 ms | 1.786 ms | 0.985 |
| Rust Unicode micro | 1.171 ms | 1.183 ms | 1.008 |

Those results established that flat numeric storage and ordinary loops were
sound, but did not by themselves justify global pooling.

## Pooling experiments

The NovelScript grammar has 449 expressions and 77 groups of positionless,
structurally duplicate expressions. Its Unicode data, however, occurs in only
two structurally different character-class matchers. The expensive duplication
is inside their Unicode range data, so a generic expression-interning mechanism
does not naturally target this case.

Pooling exact repeated expression leaves saved only 331 B of source and 1,178 B
after minification, while increasing gzip size by 50 B. This is a poor return for
introducing shared identity and non-local generated code.

Pooling complete matcher tables did meet the size limits. It also made generated
output fragile: a small grammar change produced a `+311/-3` line diff, compared
with `+2/-2` for matcher-local inline tables. Stable sorting does not solve that
problem; adding or changing a pooled item can still rename or move unrelated
generated constants when identifiers are positional, and encoded names add a
different layer of generated machinery.

The shared-table design compressed somewhat better than matcher-local tables,
but its pre-scan, naming, cross-reference, and multi-file collision rules made
the generated parser less directly traceable to the grammar.

## Final TypeScript design

Each `CharClassMatcher` now owns one flat `readonly number[]`, emitted directly
in its `classes` field. Production output keeps the table on one line to avoid
repeating indentation for every range; expanded for explanation, it is:

```ts
classes: [
  lo, hi, stride,
  lo, hi, stride,
]
```

The runtime advances through that table by three numbers at a time and exits on
the first match. It uses an ordinary loop, so it creates no callback closure.
There is no grammar pre-scan, global constant pool, generated identifier,
numbering scheme, or cross-matcher reference. A local grammar edit consequently
changes only the corresponding matcher in normal cases.

This retains some repeated numeric data when multiple matchers use the same
Unicode class. The experiment showed that removing the per-range object syntax
captures most of the useful size reduction without obscuring the grammar-to-code
mapping. The final NovelScript measurements are 86,551 B of source, 46,328 B
minified, and 12,495 B minified plus gzip, all within the requested limits of
105 KB, 50 KB, and 13,963 B.

Compared with the provisional shared-table output, matcher-local ownership costs
3,497 B of source, 8,405 B after minification, and 203 B after gzip. That is the
measured price of keeping generated identifiers, global ordering, and non-local
references out of the final design.

The final direct comparison against `c09ac57` pinned one logical CPU, warmed each
parser through at least 15 million code points, and used 21 paired samples. Each
sample contained 48 short blocks with AB/BA order alternated both within and
between samples. The result comparison ran before timing, and both generated
parsers were compiled with `tsc --strict`.

| Workload | Baseline median | Final median | Paired final / baseline | Ratio IQR |
| --- | ---: | ---: | ---: | ---: |
| NovelScript `test-js.ns` | 1.0717 ms | 1.0833 ms | 0.9948 | 0.9661-1.0489 |
| Unicode-heavy micro | 6.4714 ms | 6.3841 ms | 0.9717 | 0.9203-1.0119 |

Both ratio intervals cross 1.0 and neither run had a Tukey outlier. The result is
no observed performance regression, but not evidence for a reliable speedup.

The generator still expands Go's compile-time `unicode.RangeTable` data. It does
not use JavaScript Unicode property regular expressions, so browser Unicode
versions cannot change parser behavior. Character matching semantics and the
documented TypeScript exports (`parse`, `Parser`, and `ParserOptions`) are
unchanged.

## Why this is not generic deduplication

Automatic pooling of repeated grammar fragments would affect expression
identity, source positions, diagnostics, generated ordering, and optimization
passes across all constructs. The measurements above show little size benefit
for exact expression pooling, and the Unicode duplication occurs below the AST
expression boundary. A general mechanism would therefore add broad semantics
and maintenance risk without addressing the dominant cost as directly as the
flat representation does.

The final change is TypeScript-only. Haxe, C#, and Rust use different data
representations and runtime hot paths; their provisional savings are evidence
for separate, target-owned investigations, not justification for a shared
cross-target abstraction. Go is maintained on a separate branch and is outside
this change. None of those target implementations is modified by the final
patch.

The remaining size risk is linear repetition when a grammar uses the same large
Unicode class in many different matchers. NovelScript has 3,672 B of headroom
against the requested minified limit. If that headroom is exhausted later, the
next experiment should remain TypeScript-specific and quantify diff churn again;
the current data does not justify generic AST interning.

Keeping each table on one line also produces a 17,038-character longest line in
the NovelScript parser. `tsc` and esbuild accept it, but repositories that apply
line-length rules to generated files may need to exclude those files or choose a
stable, measured wrapping policy.

See `representative.annotated.generated.ts` for a complete generated parser with
three short annotations showing matcher-local table ownership and diff locality.
