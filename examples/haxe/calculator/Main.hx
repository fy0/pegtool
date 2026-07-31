class Main {
	static function main() {
		var argv:Array<String> = cast js.Lib.global.process.argv;
		var args = argv.slice(2);
		var input = args.length == 0 ? "18 + 3 - 27 * (-18 / -3)" : args.join(" ");
		var value:Int = cast new Parser("expression", input, []).parse(null);
		if (args.length == 0 && value != -141) throw 'unexpected result: ${value}';
		js.Lib.global.console.log(Std.string(value));
	}
}
