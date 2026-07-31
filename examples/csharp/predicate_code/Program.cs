if (!Equals(PTParser.Parse("+"), "+") || !Equals(PTParser.Parse("%"), "%"))
{
    throw new Exception("configured operator failed");
}

var rejected = false;
try
{
    PTParser.Parse("-");
}
catch (PTParser.PTParseException)
{
    rejected = true;
}
if (!rejected) throw new Exception("unexpectedly accepted -");

Console.WriteLine(PTParser.Parse("%"));
