import { parse } from "./parser";

declare const process: { argv: string[] };

const defaultExpression = "18 + 3 - 27 * (-18 / -3)";
const input = process.argv.slice(2).join(" ") || defaultExpression;
const value = parse(input) as number;

if (input === defaultExpression && value !== -141) {
  throw new Error("unexpected result: " + value);
}

console.log(value);
