const string defaultExpression = "18 + 3 - 27 * (-18 / -3)";
var input = args.Length == 0 ? defaultExpression : string.Join(" ", args);
var value = (int)PTParser.Parse(input)!;

if (input == defaultExpression && value != -141)
{
    throw new Exception("unexpected result: " + value);
}

Console.WriteLine(value);
