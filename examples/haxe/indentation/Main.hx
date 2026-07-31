import Parser.IndentationSupport;

class Main {
	static final defaultProgram = "a = 1\nif a:\n    b = 2\n    if b:\n        a = 40\n    b = a + b\nreturn b\n";

	static function run(input:String):Int {
		return IndentationSupport.execute(new Parser("indentation", input, []).parse(null));
	}

	static function checkBaseline():Void {
		if (run("x = 40\nif x:\n    x = x + 2\nreturn x\n") != 42) throw "unexpected true-if result";
		if (run("x = 0\nif x:\n    x = 5\nreturn x\n") != 0) throw "unexpected false-if result";
		var rejected = false;
		try {
			run("x = 1\nif x:\n  x = 2\nreturn x\n");
		} catch (_:Dynamic) {
			rejected = true;
		}
		if (!rejected) throw "unexpectedly accepted bad indentation";
	}

	static function main() {
		checkBaseline();
		var argv:Array<String> = cast js.Lib.global.process.argv;
		var args = argv.slice(2);
		var input:String = args.length == 0
			? defaultProgram
			: cast js.Lib.global.require("fs").readFileSync(args[0], "utf8");
		js.Lib.global.console.log(Std.string(run(input)));
	}
}
