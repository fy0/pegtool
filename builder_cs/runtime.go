package builder_cs

const runtimeCode = `
public readonly struct Position
{
    public Position(int line, int col, int offset)
    {
        Line = line;
        Col = col;
        Offset = offset;
    }

    public int Line { get; }
    public int Col { get; }
    public int Offset { get; }
}

public sealed class Current
{
    internal Current(object data)
    {
        Data = data;
        Pos = new Position(1, 1, 0);
        Text = string.Empty;
    }

    public Position Pos { get; internal set; }
    public string Text { get; internal set; }
    public object Data { get; }
}

public sealed class ParserOptions
{
    public object? CustomData { get; set; }
    public string? Entrypoint { get; set; }
    public string? Filename { get; set; }
    public long MaxExpressions { get; set; } = int.MaxValue;
}

public sealed class PigeonParseException : Exception
{
    internal PigeonParseException(string message, Position position, IReadOnlyList<string>? expected = null)
        : base(message)
    {
        Position = position;
        Expected = expected ?? Array.Empty<string>();
    }

    public Position Position { get; }
    public IReadOnlyList<string> Expected { get; }
}

private enum ExprKind
{
    Action,
    AndCode,
    And,
    AndLogical,
    Any,
    CharClass,
    Choice,
    Code,
    Labeled,
    Literal,
    NotCode,
    Not,
    NotLogical,
    OneOrMore,
    Recovery,
    RuleRef,
    RuleIndexRef,
    Sequence,
    Throw,
    ZeroOrMore,
    ZeroOrOne,
}

private sealed class Grammar
{
    public List<Rule> Rules = new List<Rule>();
}

private sealed class Rule
{
    public string Name = string.Empty;
    public string? DisplayName;
    public Expr Expr = null!;
    public bool VarExists;
    public Position Pos;
}

private readonly struct UnicodeRange
{
    public UnicodeRange(int lo, int hi, int stride)
    {
        Lo = lo;
        Hi = hi;
        Stride = stride;
    }

    public int Lo { get; }
    public int Hi { get; }
    public int Stride { get; }
}

private sealed class Expr
{
    public ExprKind Kind;
    public Expr? Inner;
    public List<Expr>? Items;
    public Func<PigeonParser, object?>? ActionRun;
    public Func<PigeonParser, bool>? PredicateRun;
    public Expr? RecoverExpr;
    public List<string>? FailureLabels;
    public string? Label;
    public bool TextCapture;
    public string? Name;
    public int Index;
    public string? Value;
    public string? Want;
    public int[]? Chars;
    public int[]? Ranges;
    public UnicodeRange[]? Classes;
    public bool IgnoreCase;
    public bool Inverted;
    public bool NotSkip;
    public Position Pos;
}

private struct Savepoint
{
    public int Line;
    public int Col;
    public int Offset;
    public int Index;
}

private struct Match
{
    public Match(object? value, bool ok)
    {
        Value = value;
        Ok = ok;
    }

    public object? Value;
    public bool Ok;
}

private readonly string filename;
private readonly string input;
private readonly string entrypoint;
private readonly long maxExpressions;
private Savepoint pt;
private readonly Dictionary<string, Rule> rules = new Dictionary<string, Rule>();
private List<Rule> rulesArray = new List<Rule>();
private readonly List<Rule> rstack = new List<Rule>();
private readonly List<Exception> errors = new List<Exception>();
private readonly List<Dictionary<string, Expr>> recoveryStack = new List<Dictionary<string, Expr>>();
private int skipDepth;
private bool invertFailures;
private Savepoint maxFailPos;
private readonly List<string> maxFailExpected = new List<string>();
private long expressionCount;
private readonly List<Dictionary<string, object?>> vstack = new List<Dictionary<string, object?>>();

public PigeonParser(string filename, string input, ParserOptions? options = null)
{
    options ??= new ParserOptions();
    this.filename = filename;
    this.input = input;
    entrypoint = options.Entrypoint ?? "{{ .Entrypoint }}";
    maxExpressions = options.MaxExpressions <= 0 ? int.MaxValue : options.MaxExpressions;
    pt = new Savepoint { Line = 1, Col = 1, Offset = 0, Index = 0 };
    maxFailPos = pt;
    Cur = new Current(options.CustomData ?? new object());
}

public Current Cur { get; }

public static object? Parse(string input, ParserOptions? options = null)
{
    options ??= new ParserOptions();
    return new PigeonParser(options.Filename ?? string.Empty, input, options).Parse();
}

public object? Parse()
{
    return ParseGrammar({{ .GrammarVarName }});
}

public void addErr(Exception error)
{
    errors.Add(error);
}

private object? ParseGrammar(Grammar grammar)
{
    if (grammar.Rules.Count == 0)
    {
        throw new PigeonParseException("grammar has no rule", PublicPosition(pt));
    }

    rulesArray = grammar.Rules;
    foreach (var rule in grammar.Rules)
    {
        rules[rule.Name] = rule;
    }

    if (!rules.TryGetValue(entrypoint, out var start))
    {
        throw new PigeonParseException("invalid entrypoint", PublicPosition(pt));
    }

    Match result;
    try
    {
        result = ParseRule(start);
    }
    catch (Exception error)
    {
        addErr(error);
        throw CombinedError();
    }

    if (!result.Ok && errors.Count == 0)
    {
        var expected = maxFailExpected.Distinct().OrderBy(value => value, StringComparer.Ordinal).ToList();
        if (expected.Remove("!."))
        {
            expected.Add("EOF");
        }
        var position = PublicPosition(maxFailPos);
        var prefix = filename.Length == 0 ? string.Empty : filename + ":";
        errors.Add(new PigeonParseException(
            prefix + position.Line + ":" + position.Col + ": no match found, expected: " + string.Join(", ", expected),
            position,
            expected));
    }

    if (errors.Count > 0)
    {
        throw CombinedError();
    }
    return result.Value;
}

private Exception CombinedError()
{
    return errors.Count == 1 ? errors[0] : new AggregateException(errors);
}

private static Position PublicPosition(Savepoint point)
{
    return new Position(point.Line, point.Col, point.Offset);
}

private int? CurrentCodePoint()
{
    if (pt.Index >= input.Length)
    {
        return null;
    }
    return Rune.GetRuneAt(input, pt.Index).Value;
}

private void Advance()
{
    if (pt.Index >= input.Length)
    {
        return;
    }
    var rune = Rune.GetRuneAt(input, pt.Index);
    pt.Index += rune.Utf16SequenceLength;
    pt.Offset += rune.Utf8SequenceLength;
    if (rune.Value == 10)
    {
        pt.Line++;
        pt.Col = 1;
    }
    else
    {
        pt.Col++;
    }
}

private void Restore(Savepoint point)
{
    pt = point;
}

private string SliceFrom(Savepoint point)
{
    return input.Substring(point.Index, pt.Index - point.Index);
}

private bool Skipping()
{
    return skipDepth > 0;
}

private void FailAt(bool matched, Savepoint point, string want)
{
    if (matched != invertFailures || point.Offset < maxFailPos.Offset)
    {
        return;
    }
    if (point.Offset > maxFailPos.Offset)
    {
        maxFailPos = point;
        maxFailExpected.Clear();
    }
    maxFailExpected.Add(invertFailures ? "!" + want : want);
}

private Match ParseRule(Rule rule)
{
    rstack.Add(rule);
    var hasVariables = rule.VarExists && !Skipping();
    if (hasVariables)
    {
        vstack.Add(new Dictionary<string, object?>());
    }
    try
    {
        return ParseExpr(rule.Expr);
    }
    finally
    {
        if (hasVariables)
        {
            vstack.RemoveAt(vstack.Count - 1);
        }
        rstack.RemoveAt(rstack.Count - 1);
    }
}

private Match ParseExpr(Expr expr)
{
    expressionCount++;
    if (expressionCount > maxExpressions)
    {
        throw new InvalidOperationException("max number of expressions parsed");
    }

    switch (expr.Kind)
    {
        case ExprKind.Action: return ParseAction(expr);
        case ExprKind.AndCode: return new Match(null, expr.PredicateRun!(this));
        case ExprKind.And: return ParseAnd(expr, false);
        case ExprKind.AndLogical: return ParseAnd(expr, true);
        case ExprKind.Any: return ParseAny();
        case ExprKind.CharClass: return ParseCharClass(expr);
        case ExprKind.Choice: return ParseChoice(expr);
        case ExprKind.Code: return ParseCode(expr);
        case ExprKind.Labeled: return ParseLabeled(expr);
        case ExprKind.Literal: return ParseLiteral(expr);
        case ExprKind.NotCode: return new Match(null, !expr.PredicateRun!(this));
        case ExprKind.Not: return ParseNot(expr, false);
        case ExprKind.NotLogical: return ParseNot(expr, true);
        case ExprKind.OneOrMore: return ParseOneOrMore(expr);
        case ExprKind.Recovery: return ParseRecovery(expr);
        case ExprKind.RuleRef: return ParseRuleRef(expr);
        case ExprKind.RuleIndexRef: return ParseRuleIndexRef(expr);
        case ExprKind.Sequence: return ParseSequence(expr);
        case ExprKind.Throw: return ParseThrow(expr);
        case ExprKind.ZeroOrMore: return ParseZeroOrMore(expr);
        case ExprKind.ZeroOrOne: return ParseZeroOrOne(expr);
        default: throw new InvalidOperationException("unknown expression type");
    }
}

private Match ParseAction(Expr expr)
{
    if (Skipping())
    {
        var skipped = ParseExpr(expr.Inner!);
        return new Match(null, skipped.Ok);
    }
    var start = pt;
    var result = ParseExpr(expr.Inner!);
    if (result.Ok)
    {
        Cur.Pos = PublicPosition(start);
        Cur.Text = SliceFrom(start);
        result.Value = expr.ActionRun!(this);
    }
    return result;
}

private Match ParseAnd(Expr expr, bool logical)
{
    var start = pt;
    skipDepth++;
    Match result;
    try
    {
        result = ParseExpr(expr.Inner!);
    }
    finally
    {
        skipDepth--;
    }
    var consumed = pt.Offset != start.Offset;
    Restore(start);
    return new Match(null, logical ? result.Ok && consumed : result.Ok);
}

private Match ParseNot(Expr expr, bool logical)
{
    var start = pt;
    invertFailures = !invertFailures;
    skipDepth++;
    Match result;
    try
    {
        result = ParseExpr(expr.Inner!);
    }
    finally
    {
        skipDepth--;
        invertFailures = !invertFailures;
    }
    var consumed = pt.Offset != start.Offset;
    Restore(start);
    return new Match(null, logical ? !result.Ok && consumed : !result.Ok);
}

private Match ParseAny()
{
    if (CurrentCodePoint() == null)
    {
        FailAt(false, pt, ".");
        return new Match(null, false);
    }
    FailAt(true, pt, ".");
    Advance();
    return new Match(null, true);
}

private static int LowerCodePoint(int value)
{
    return Rune.ToLowerInvariant(new Rune(value)).Value;
}

private Match ParseCharClass(Expr expr)
{
    var currentValue = CurrentCodePoint();
    if (currentValue == null)
    {
        FailAt(false, pt, expr.Value!);
        return new Match(null, false);
    }
    var current = expr.IgnoreCase ? LowerCodePoint(currentValue.Value) : currentValue.Value;
    var found = expr.Chars != null && Array.IndexOf(expr.Chars, current) >= 0;

    if (!found && expr.Ranges != null)
    {
        for (var index = 0; index < expr.Ranges.Length; index += 2)
        {
            if (current >= expr.Ranges[index] && current <= expr.Ranges[index + 1])
            {
                found = true;
                break;
            }
        }
    }
    if (!found && expr.Classes != null)
    {
        foreach (var range in expr.Classes)
        {
            if (current >= range.Lo && current <= range.Hi && (current - range.Lo) % range.Stride == 0)
            {
                found = true;
                break;
            }
        }
    }

    var matched = expr.Inverted ? !found : found;
    FailAt(matched, pt, expr.Value!);
    if (!matched)
    {
        return new Match(null, false);
    }
    Advance();
    return new Match(null, true);
}

private Match ParseChoice(Expr expr)
{
    foreach (var alternative in expr.Items!)
    {
        var result = ParseExpr(alternative);
        if (result.Ok)
        {
            return result;
        }
    }
    return new Match(null, false);
}

private Match ParseCode(Expr expr)
{
    if (!expr.NotSkip && Skipping())
    {
        return new Match(null, true);
    }
    return new Match(expr.ActionRun!(this), true);
}

private Match ParseLabeled(Expr expr)
{
    var start = pt;
    var result = ParseExpr(expr.Inner!);
    if (result.Ok && !string.IsNullOrEmpty(expr.Label) && !Skipping())
    {
        vstack[^1][expr.Label!] = expr.TextCapture ? SliceFrom(start) : result.Value;
    }
    return result;
}

private Match ParseLiteral(Expr expr)
{
    var start = pt;
    foreach (var wantedRune in expr.Value!.EnumerateRunes())
    {
        var wanted = expr.IgnoreCase ? Rune.ToLowerInvariant(wantedRune).Value : wantedRune.Value;
        var currentValue = CurrentCodePoint();
        var current = currentValue == null
            ? (int?)null
            : expr.IgnoreCase ? LowerCodePoint(currentValue.Value) : currentValue.Value;
        if (current != wanted)
        {
            FailAt(false, start, expr.Want!);
            Restore(start);
            return new Match(null, false);
        }
        Advance();
    }
    FailAt(true, start, expr.Want!);
    return new Match(null, true);
}

private Match ParseOneOrMore(Expr expr)
{
    var values = new List<object?>();
    var matched = false;
    while (true)
    {
        var before = pt.Offset;
        var result = ParseExpr(expr.Inner!);
        if (!result.Ok)
        {
            return new Match(values.Count == 0 ? null : values, matched);
        }
        if (pt.Offset == before)
        {
            throw new InvalidOperationException("repetition matched empty input");
        }
        matched = true;
        if (result.Value != null)
        {
            values.Add(result.Value);
        }
    }
}

private Match ParseZeroOrMore(Expr expr)
{
    var values = new List<object?>();
    while (true)
    {
        var before = pt.Offset;
        var result = ParseExpr(expr.Inner!);
        if (!result.Ok)
        {
            return new Match(values.Count == 0 ? null : values, true);
        }
        if (pt.Offset == before)
        {
            throw new InvalidOperationException("repetition matched empty input");
        }
        if (result.Value != null)
        {
            values.Add(result.Value);
        }
    }
}

private Match ParseZeroOrOne(Expr expr)
{
    var result = ParseExpr(expr.Inner!);
    return new Match(result.Value, true);
}

private Match ParseSequence(Expr expr)
{
    var start = pt;
    var values = new List<object?>();
    foreach (var item in expr.Items!)
    {
        var result = ParseExpr(item);
        if (!result.Ok)
        {
            Restore(start);
            return new Match(null, false);
        }
        if (!Skipping() && result.Value != null)
        {
            values.Add(result.Value);
        }
    }
    return new Match(values.Count == 0 ? null : values, true);
}

private Match ParseRuleRef(Expr expr)
{
    if (!rules.TryGetValue(expr.Name!, out var rule))
    {
        addErr(new InvalidOperationException("undefined rule: " + expr.Name));
        return new Match(null, false);
    }
    return ParseRule(rule);
}

private Match ParseRuleIndexRef(Expr expr)
{
    if (expr.Index < 0 || expr.Index >= rulesArray.Count)
    {
        addErr(new InvalidOperationException("undefined rule index: " + expr.Index));
        return new Match(null, false);
    }
    return ParseRule(rulesArray[expr.Index]);
}

private Match ParseRecovery(Expr expr)
{
    var recovery = new Dictionary<string, Expr>();
    foreach (var label in expr.FailureLabels!)
    {
        recovery[label] = expr.RecoverExpr!;
    }
    recoveryStack.Add(recovery);
    try
    {
        return ParseExpr(expr.Inner!);
    }
    finally
    {
        recoveryStack.RemoveAt(recoveryStack.Count - 1);
    }
}

private Match ParseThrow(Expr expr)
{
    for (var index = recoveryStack.Count - 1; index >= 0; index--)
    {
        if (recoveryStack[index].TryGetValue(expr.Label!, out var recovery))
        {
            var result = ParseExpr(recovery);
            if (result.Ok)
            {
                return result;
            }
        }
    }
    return new Match(null, false);
}
`
