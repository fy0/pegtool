class Main {
	static function parseOperator(input:String):String {
		return cast new Parser("operator", input, []).parse(null);
	}

	static function main() {
		if (parseOperator("+") != "+" || parseOperator("%") != "%") throw "configured operator failed";
		var rejected = false;
		try {
			parseOperator("-");
		} catch (_:Dynamic) {
			rejected = true;
		}
		if (!rejected) throw "unexpectedly accepted -";
		js.Lib.global.console.log(parseOperator("%"));
	}
}
