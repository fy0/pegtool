const string defaultProgram = "a = 1\nif a:\n    b = 2\n    if b:\n        a = 40\n    b = a + b\nreturn b\n";

static int Run(string input)
{
    var options = new PigeonParser.ParserOptions
    {
        CustomData = new PigeonParser.IndentationData(),
    };
    return PigeonParser.Execute(PigeonParser.Parse(input, options));
}

static void CheckBaseline()
{
    if (Run("x = 40\nif x:\n    x = x + 2\nreturn x\n") != 42) throw new Exception("unexpected true-if result");
    if (Run("x = 0\nif x:\n    x = 5\nreturn x\n") != 0) throw new Exception("unexpected false-if result");
    var rejected = false;
    try
    {
        Run("x = 1\nif x:\n  x = 2\nreturn x\n");
    }
    catch (PigeonParser.PigeonParseException)
    {
        rejected = true;
    }
    if (!rejected) throw new Exception("unexpectedly accepted bad indentation");
}

CheckBaseline();
var input = args.Length == 0 ? defaultProgram : File.ReadAllText(args[0]);
Console.WriteLine(Run(input));
