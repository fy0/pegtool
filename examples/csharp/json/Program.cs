static void CheckBaseline()
{
    var value = (Dictionary<string, object?>)PTParser.Parse("{\"items\":[1,null,true],\"name\":\"pigeon\"}")!;
    var items = (List<object?>)value["items"]!;
    if (!Equals(value["name"], "pigeon") || items.Count != 3 || !Equals(items[0], 1.0) || items[1] != null || !Equals(items[2], true))
    {
        throw new Exception("unexpected object value");
    }

    var escaped = (string)PTParser.Parse("\"line\\n\\uD83D\\uDE00\\/\"")!;
    if (escaped != "line\n\U0001F600/") throw new Exception("unexpected escaped string: " + escaped);

    var rejected = false;
    try
    {
        PTParser.Parse("00");
    }
    catch (PTParser.PTParseException)
    {
        rejected = true;
    }
    if (!rejected) throw new Exception("unexpectedly accepted a leading zero");
}

CheckBaseline();
var input = args.Length == 0
    ? "{\"pigeon\":true,\"targets\":[\"go\",\"rust\",\"haxe\",\"typescript\",\"csharp\"]}"
    : string.Join(" ", args);
Console.WriteLine(System.Text.Json.JsonSerializer.Serialize(PTParser.Parse(input)));
