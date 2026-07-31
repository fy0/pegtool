class Main {
	static function parseJson(input:String):Dynamic {
		return new Parser("json", input, []).parse(null);
	}

	static function checkBaseline():Void {
		var value:Dynamic = parseJson('{"items":[1,null,true],"name":"pigeon"}');
		var items:Array<Dynamic> = cast Reflect.field(value, "items");
		if (Reflect.field(value, "name") != "pigeon" || items.length != 3 || items[0] != 1 || items[1] != null || items[2] != true) {
			throw "unexpected object value";
		}

		var escaped:String = cast parseJson('"line\\n\\uD83D\\uDE00\\/"');
		var expected = "line\n" + String.fromCharCode(0xd83d) + String.fromCharCode(0xde00) + "/";
		if (escaped != expected) throw 'unexpected escaped string: ${escaped}';

		var rejected = false;
		try {
			parseJson("00");
		} catch (_:Dynamic) {
			rejected = true;
		}
		if (!rejected) throw "unexpectedly accepted a leading zero";
	}

	static function main() {
		checkBaseline();
		var argv:Array<String> = cast js.Lib.global.process.argv;
		var args = argv.slice(2);
		var input = args.length == 0 ? '{"pigeon":true,"targets":["go","rust","haxe","typescript","csharp"]}' : args.join(" ");
		js.Lib.global.console.log(haxe.Json.stringify(parseJson(input)));
	}
}
