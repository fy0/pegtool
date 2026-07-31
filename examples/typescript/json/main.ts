import { parse } from "./parser";

declare const process: { argv: string[] };

function checkBaseline(): void {
  const value = parse('{"items":[1,null,true],"name":"pigeon"}') as Record<string, any>;
  const items = value.items as any[];
  if (value.name !== "pigeon" || items.length !== 3 || items[0] !== 1 || items[1] !== null || items[2] !== true) {
    throw new Error("unexpected object value");
  }

  const escaped = parse('"line\\n\\uD83D\\uDE00\\/"');
  const expected = "line\n" + String.fromCodePoint(0x1f600) + "/";
  if (escaped !== expected) throw new Error("unexpected escaped string: " + String(escaped));

  let rejected = false;
  try {
    parse("00");
  } catch {
    rejected = true;
  }
  if (!rejected) throw new Error("unexpectedly accepted a leading zero");
}

checkBaseline();
const input = process.argv.slice(2).join(" ") || '{"pigeon":true,"targets":["go","rust","haxe","typescript","csharp"]}';
console.log(JSON.stringify(parse(input)));
