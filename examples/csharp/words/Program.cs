var expected = new[] { "alpha", "\u4e16\u754c", "\U00010400" };
var actual = (List<string>)PTParser.Parse(string.Join(",", expected))!;

if (!actual.SequenceEqual(expected))
{
    throw new Exception("unexpected parse result: " + string.Join(",", actual));
}

Console.WriteLine(string.Join("|", actual));
