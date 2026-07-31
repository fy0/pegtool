import { execute, parse } from "./parser";

declare const process: { argv: string[] };
declare function require(name: string): { readFileSync(path: string, encoding: string): string };

const defaultProgram = "a = 1\nif a:\n    b = 2\n    if b:\n        a = 40\n    b = a + b\nreturn b\n";

function run(input: string): number {
  return execute(parse(input, { customData: { indentation: 0 } }));
}

function checkBaseline(): void {
  if (run("x = 40\nif x:\n    x = x + 2\nreturn x\n") !== 42) throw new Error("unexpected true-if result");
  if (run("x = 0\nif x:\n    x = 5\nreturn x\n") !== 0) throw new Error("unexpected false-if result");
  let rejected = false;
  try {
    run("x = 1\nif x:\n  x = 2\nreturn x\n");
  } catch {
    rejected = true;
  }
  if (!rejected) throw new Error("unexpectedly accepted bad indentation");
}

checkBaseline();
const path = process.argv[2];
const input = path === undefined ? defaultProgram : require("fs").readFileSync(path, "utf8");
console.log(run(input));
