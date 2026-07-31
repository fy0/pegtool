import { parse } from "./parser";

const expected = ["alpha", "\u4e16\u754c", "\u{10400}"];
const actual = parse(expected.join(",")) as string[];

if (actual.length !== expected.length || actual.some((value, index) => value !== expected[index])) {
  throw new Error("unexpected parse result: " + JSON.stringify(actual));
}

console.log(actual.join("|"));
