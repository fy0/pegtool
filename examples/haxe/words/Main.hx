class Main {
	static function main() {
		var last = String.fromCharCode(0x10400);
		var parser = new Parser("words.txt", "alpha,\u4e16\u754c," + last, []);
		var values:Array<Dynamic> = cast parser.parse(null);
		if (values.length != 3 || values[0] != "alpha" || values[1] != "\u4e16\u754c" || values[2] != last) {
			throw 'unexpected parse result: ${values}';
		}
		js.Lib.global.console.log(values.map(value -> Std.string(value)).join("|"));
	}
}
