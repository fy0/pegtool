# TypeScript Unicode Range Runtime Experiment

Status: matcher-local segmented tables with binary lookup selected.

Baseline commit: `c09ac57`
Flat-table commit: `6ed9a1d`
Date: 2026-07-30

## Scope

This experiment measures the runtime of the generated TypeScript parser's
`parse()` function. It does not measure Pigeon grammar parsing, TypeScript code
generation, `tsc`, esbuild, or JavaScript module loading.

The first shared-table prototype appeared to make the Unicode-heavy workload
36% faster. That run changed object layout, loop structure, and table ownership
at the same time, and later testing found large V8 tiering effects. The purpose
of this follow-up is to separate those variables and include non-JIT runtimes.

## Compared representations

1. **Object baseline**: one `{ lo, hi, stride }` object per range and
   `Array.prototype.some()` in the hot path.
2. **Shared per-class tables**: global flat arrays, matcher references, and two
   linear loops.
3. **Matcher-local combined flat table**: the committed `6ed9a1d` design, with
   one linear loop over a matcher-owned number array.
4. **Matcher-local segmented linear tables**: one inline flat array per Unicode
   class, still scanned linearly. This removes global sharing while retaining
   the shared prototype's runtime shape.
5. **Inflated objects**: compact flat source data expanded into range objects at
   module initialization, then matched with the old `.some()` path.
6. **Matcher-local segmented binary tables**: one inline flat array per Unicode
   class, searched with binary lookup. This is the selected design.

The inflated-object candidate restores baseline parse behavior but recreates all
1,801 range objects at startup. It cannot improve the old hot path and gives up
the allocation benefit of flat storage, so it was rejected once binary lookup
proved both smaller and faster.

## Method

- All variants were generated from the same grammar and compiled with the same
  TypeScript 5.9.3 settings: `--strict --target ES2020 --module commonjs`.
- Parser results were serialized and compared before timing.
- Benchmarks were pinned to logical CPU 2 on an i5-13600KF.
- Each sample alternated variants in short blocks. Three-way tests rotated all
  six execution orders; pairwise tests alternated AB and BA.
- V8 JIT runs warmed each parser through at least 15 million code points.
- JITless and QuickJS runs used time-based warmup because there is no tiering.
- The main V8 and Node comparisons used 21 samples with 48 blocks per sample.
- tjs used 15 samples. Its Unicode test used 12 balanced blocks and its real
  workload used 18 balanced blocks.

Runtime versions:

- Node 24.4.0, V8 13.6.233.10-node.17
- tjs 26.4.0, QuickJS 0.13.0
- Windows 11 build 22631

## What the pooling result means

With adequate V8 warmup, shared per-class pooling did reproduce a real-world
benefit but not the original 36% Unicode result:

| V8 JIT workload | Object baseline | Shared tables | Paired ratio | Result |
| --- | ---: | ---: | ---: | --- |
| NovelScript `test-js.ns` | 1.0796 ms | 0.9345 ms | 0.8661 | 13.4% faster |
| Unicode-heavy | 5.3083 ms | 5.1562 ms | 0.9741 | 2.6% faster |

The Unicode paired-ratio IQR was `0.9495-0.9881`; the original 0.641 ratio did
not reproduce. The NovelScript result closely matches the first run's 0.859
ratio and is a real V8-specific signal.

The same shared representation was not generally faster without a JIT:

| Runtime | Workload | Shared / object | Interpretation |
| --- | --- | ---: | --- |
| Node `--jitless` | NovelScript | 1.0406 | approximately neutral |
| Node `--jitless` | Unicode-heavy | 1.4949 | 49% slower |
| tjs / QuickJS | NovelScript | 1.0660 | 6.6% slower |
| tjs / QuickJS | Unicode-heavy | 1.6955 | 70% slower |

QuickJS also measured matcher-local segmented linear tables at 60.48 ms for the
Unicode workload, versus 34.75 ms for object `.some()`. Global sharing is not
the cause of the interpreter behavior. The native array builtin drives the old
loop, while every iteration of the flat linear loop is interpreted JavaScript.

## Binary lookup

Go documents the required invariant for `unicode.RangeTable`: R16 and R32 are
sorted and non-overlapping, and R32 contains only values at or above U+10000.
Pigeon keeps each Unicode class in its own local table, preserving that
invariant. A matcher with several classes tests each table independently.

Within one table, binary lookup finds the only range that could contain the code
point and then applies the original stride check. A stride miss can safely stop
searching that table because its ranges do not overlap.

The test suite exhaustively compares binary lookup with `unicode.Is` for every
code point from U+0000 through U+10FFFF across representative categories,
properties, and scripts, including stride-heavy and non-BMP tables.

## Binary performance

Ratios below 1 are faster. IQR is the interquartile range of paired ratios.

| Runtime | Workload | Object | Binary | Binary / object | IQR |
| --- | --- | ---: | ---: | ---: | ---: |
| V8 JIT | Unicode-heavy | 5.2951 ms | 0.7077 ms | 0.1321 | 0.1308-0.1368 |
| V8 JIT | NovelScript | 1.0749 ms | 0.9628 ms | 0.8824 | 0.8531-0.9123 |
| Node `--jitless` | Unicode-heavy | 6.7784 ms | 0.8448 ms | 0.1254 | 0.1227-0.1309 |
| Node `--jitless` | NovelScript | 12.6710 ms | 11.9247 ms | 0.9457 | 0.9303-0.9627 |
| tjs / QuickJS | Unicode-heavy | 34.7529 ms | 3.5648 ms | 0.1061 | 0.1029-0.1084 |
| tjs / QuickJS | NovelScript | 52.8901 ms | 51.6714 ms | 0.9396 | 0.8998-1.0106 |

The Unicode-heavy improvement is 7.6x on V8, 8.0x in Node JITless, and 9.4x on
tjs. Real NovelScript parsing improves by about 12% on V8 and 5% in Node
JITless. The tjs real-workload interval crosses 1.0, so that result should be
treated as neutral rather than a claimed 6% improvement.

## Generated size

NovelScript was generated with:

```text
-nolint -optimize-ref-expr-by-index -optimize-parser -t ts
```

esbuild 0.27.0 performed minification and zlib level 9 produced gzip sizes.

| Representation | Source | Minified | Minified + gzip |
| --- | ---: | ---: | ---: |
| Object baseline | 144,077 B | 73,282 B | 14,035 B |
| Shared per-class prototype | 83,054 B | 37,923 B | 12,292 B |
| Combined matcher-local flat | 86,551 B | 46,328 B | 12,495 B |
| Segmented matcher-local linear | 86,669 B | 46,347 B | 12,534 B |
| Inflated objects | 86,793 B | 46,369 B | 12,540 B |
| **Segmented matcher-local binary** | **86,925 B** | **46,406 B** | **12,578 B** |

The selected design remains below all acceptance limits: 105 KB source, 50 KB
minified, and 13,963 B gzip. It introduces no pre-scan, global pool, generated
identifier, numbering scheme, or cross-matcher reference.

## Decision and residual risk

Use matcher-local segmented flat tables with binary lookup.

This design keeps grammar-to-output locality and stable diffs, retains the size
benefit of flat numeric data, and removes the interpreter regression caused by
linear JavaScript scanning. It also outperforms the object baseline on
Unicode-heavy workloads across both tested runtime families.

Repeated Unicode classes still duplicate numeric source data across matchers.
That cost is intentional. If a future grammar exceeds the minified-size budget,
sharing should be reconsidered as a separate TypeScript-specific decision rather
than coupled to the lookup algorithm.

The generated tables must remain sorted and non-overlapping within each class.
That is guaranteed by Go's compile-time `unicode.RangeTable` data. Code that
constructs internal grammar expressions manually instead of using generated
tables would need to preserve the same invariant.

Raw benchmark harnesses, generated variants, and JSON results are retained under
`C:\Users\test\AppData\Local\Temp\pigeon-unicode-jitless` and
`C:\Users\test\AppData\Local\Temp\pigeon-unicode-both`.
