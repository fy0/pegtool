import { parse } from "./parser";

if (parse("+") !== "+" || parse("%") !== "%") {
  throw new Error("configured operator failed");
}

let rejected = false;
try {
  parse("-");
} catch {
  rejected = true;
}
if (!rejected) throw new Error("unexpectedly accepted -");

console.log(parse("%"));
