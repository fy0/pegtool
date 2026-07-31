if (!Equals(PigeonParser.Parse("+"), "+") || !Equals(PigeonParser.Parse("%"), "%"))
{
    throw new Exception("configured operator failed");
}

var rejected = false;
try
{
    PigeonParser.Parse("-");
}
catch (PigeonParser.PigeonParseException)
{
    rejected = true;
}
if (!rejected) throw new Exception("unexpectedly accepted -");

Console.WriteLine(PigeonParser.Parse("%"));
