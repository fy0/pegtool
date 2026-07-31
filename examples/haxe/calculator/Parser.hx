// Code is generated; DO NOT EDIT.
import haxe.Exception;
import haxe.ds.GenericStack;
import haxe.io.BytesBuffer;
import haxe.ds.HashMap;
import haxe.iterators.StringIteratorUnicode;

class ParserCustomData {
	public function new() {}
}

class CalculatorSupport {
	public static function evaluate(first:Dynamic, rest:Dynamic):Int {
		var result:Int = cast first;
		var items:Array<Dynamic> = rest == null ? [] : cast rest;
		for (item in items) {
			var pair:Array<Dynamic> = cast item;
			var right:Int = cast pair[1];
			switch (Std.string(pair[0])) {
				case "+": result += right;
				case "-": result -= right;
				case "*": result *= right;
				case "/": result = Std.int(result / right);
				default: throw 'unknown operator ${Std.string(pair[0])}';
			}
		}
		return result;
	}
}

var g:Grammar = {
	rules: [
		{
			name: "Input",
			varExists: true,
			expr: ActionExpr({
				run: _call_onInput_1,
				expr: SeqExpr({
					exprs: [
						LabeledExpr({
							label: "expression",
							expr: RuleRefExpr({name: "Expr",}),
							textCapture: false,
						}),
						RuleRefExpr({name: "EOF",}),
					],
				}),
			}),
		},
		{
			name: "Expr",
			varExists: true,
			expr: ActionExpr({
				run: _call_onExpr_1,
				expr: SeqExpr({
					exprs: [
						RuleRefExpr({name: "_",}),
						LabeledExpr({
							label: "first",
							expr: RuleRefExpr({name: "Term",}),
							textCapture: false,
						}),
						LabeledExpr({
							label: "rest",
							expr: ZeroOrMoreExpr({
								expr: ActionExpr({
									run: _call_onExpr_8,
									expr: SeqExpr({
										exprs: [
											RuleRefExpr({name: "_",}),
											LabeledExpr({
												label: "op",
												expr: RuleRefExpr({name: "AddOp",}),
												textCapture: false,
											}),
											RuleRefExpr({name: "_",}),
											LabeledExpr({
												label: "right",
												expr: RuleRefExpr({name: "Term",}),
												textCapture: false,
											}),
										],
									}),
								}),
							}),
							textCapture: false,
						}),
						RuleRefExpr({name: "_",}),
					],
				}),
			}),
		},
		{
			name: "Term",
			varExists: true,
			expr: ActionExpr({
				run: _call_onTerm_1,
				expr: SeqExpr({
					exprs: [
						LabeledExpr({
							label: "first",
							expr: RuleRefExpr({name: "Factor",}),
							textCapture: false,
						}),
						LabeledExpr({
							label: "rest",
							expr: ZeroOrMoreExpr({
								expr: ActionExpr({
									run: _call_onTerm_7,
									expr: SeqExpr({
										exprs: [
											RuleRefExpr({name: "_",}),
											LabeledExpr({
												label: "op",
												expr: RuleRefExpr({name: "MulOp",}),
												textCapture: false,
											}),
											RuleRefExpr({name: "_",}),
											LabeledExpr({
												label: "right",
												expr: RuleRefExpr({name: "Factor",}),
												textCapture: false,
											}),
										],
									}),
								}),
							}),
							textCapture: false,
						}),
					],
				}),
			}),
		},
		{
			name: "Factor",
			varExists: true,
			expr: ChoiceExpr({
				alternatives: [
					ActionExpr({
						run: _call_onFactor_2,
						expr: SeqExpr({
							exprs: [
								LitMatcher({val: "(", want: "\"(\"",}),
								LabeledExpr({
									label: "expression",
									expr: RuleRefExpr({name: "Expr",}),
									textCapture: false,
								}),
								LitMatcher({val: ")", want: "\")\"",}),
							],
						}),
					}),
					ActionExpr({
						run: _call_onFactor_8,
						expr: LabeledExpr({
							label: "integer",
							expr: RuleRefExpr({name: "Integer",}),
							textCapture: false,
						}),
					}),
				],
			}),
		},
		{
			name: "AddOp",
			expr: ActionExpr({
				run: _call_onAddOp_1,
				expr: ChoiceExpr({
					alternatives: [
						LitMatcher({val: "+", want: "\"+\"",}),
						LitMatcher({val: "-", want: "\"-\"",}),
					],
				}),
			}),
		},
		{
			name: "MulOp",
			expr: ActionExpr({
				run: _call_onMulOp_1,
				expr: ChoiceExpr({
					alternatives: [
						LitMatcher({val: "*", want: "\"*\"",}),
						LitMatcher({val: "/", want: "\"/\"",}),
					],
				}),
			}),
		},
		{
			name: "Integer",
			expr: ActionExpr({
				run: _call_onInteger_1,
				expr: SeqExpr({
					exprs: [
						ZeroOrOneExpr({
							expr: LitMatcher({val: "-", want: "\"-\"",}),
						}),
						OneOrMoreExpr({
							expr: CharClassMatcher({
								val: "[0-9]",
								ranges: [0x30, 0x39,],
							}),
						}),
					],
				}),
			}),
		},
		{
			name: "_",
			displayName: "\"whitespace\"",
			expr: ZeroOrMoreExpr({
				expr: CharClassMatcher({
					val: "[ \\n\\t\\r]",
					chars: [0x20, 0xA, 0x9, 0xD,],
				}),
			}),
		},
		{
			name: "EOF",
			expr: NotExpr({
				expr: AnyMatcher({}),
			}),
		},
	],
}

private function _call_onInput_1(p:Parser):Any {
	var stack = p.vstack.first();
	var c = p.cur;
	return (function(c:Current, expression:Any):Any {
		return expression;
		return null;
	})(p.cur, stack["expression"]);
}

private function _call_onExpr_8(p:Parser):Any {
	var stack = p.vstack.first();
	var c = p.cur;
	return (function(c:Current, op, right:Any):Any {
		return [op, right];
		return null;
	})(p.cur, stack["op"], stack["right"]);
}

private function _call_onExpr_1(p:Parser):Any {
	var stack = p.vstack.first();
	var c = p.cur;
	return (function(c:Current, first, rest:Any):Any {
		return CalculatorSupport.evaluate(first, rest);
		return null;
	})(p.cur, stack["first"], stack["rest"]);
}

private function _call_onTerm_7(p:Parser):Any {
	var stack = p.vstack.first();
	var c = p.cur;
	return (function(c:Current, op, right:Any):Any {
		return [op, right];
		return null;
	})(p.cur, stack["op"], stack["right"]);
}

private function _call_onTerm_1(p:Parser):Any {
	var stack = p.vstack.first();
	var c = p.cur;
	return (function(c:Current, first, rest:Any):Any {
		return CalculatorSupport.evaluate(first, rest);
		return null;
	})(p.cur, stack["first"], stack["rest"]);
}

private function _call_onFactor_2(p:Parser):Any {
	var stack = p.vstack.first();
	var c = p.cur;
	return (function(c:Current, expression:Any):Any {
		return expression;
		return null;
	})(p.cur, stack["expression"]);
}

private function _call_onFactor_8(p:Parser):Any {
	var stack = p.vstack.first();
	var c = p.cur;
	return (function(c:Current, integer:Any):Any {
		return integer;
		return null;
	})(p.cur, stack["integer"]);
}

private function _call_onAddOp_1(p:Parser):Any {
	var c = p.cur;
	return (function(c:Current):Any {
		return c.text;
		return null;
	})(p.cur);
}

private function _call_onMulOp_1(p:Parser):Any {
	var c = p.cur;
	return (function(c:Current):Any {
		return c.text;
		return null;
	})(p.cur);
}

private function _call_onInteger_1(p:Parser):Any {
	var c = p.cur;
	return (function(c:Current):Any {
		var value = Std.parseInt(c.text);
	if (value == null) throw 'invalid integer ${c.text}';
	return value;
		return null;
	})(p.cur);
}

class Errors {
	public static var errNoRule:Exception = new Exception("grammar has no rule");
	public static var errInvalidEntrypoint:Exception = new Exception("invalid entrypoint");
	public static var errInvalidEncoding:Exception = new Exception("invalid encoding");
	public static var errMaxExprCnt:Exception = new Exception("max number of expressions parsed");
}

// alternative of Uint64
typedef Uint64A = UInt;
typedef RuneA = UInt;

typedef RetValOk = {
	var val:Any;
	var ok:Bool;
}

// Option is a function that can set an option on the parser. It returns
// the previous setting as an Option.
// typedef Option = (p:Parser) -> Option; // 注: 这是个函数
typedef Option = (p:Parser) -> Any; // 注: 这是个函数，因static编译需求返回Any

// NoMatchErrorFormatter builds the error reported when parsing cannot match.
// Returning null keeps the default no-match error.
typedef NoMatchErrorFormatter = (position:Position, input:String, expected:Array<String>) -> Exception;

function noMatchErrorFormatter(fn:NoMatchErrorFormatter):Option {
	return function(p:Parser):Option {
		var old = p.noMatchErrorFormatter;
		p.noMatchErrorFormatter = fn;
		return noMatchErrorFormatter(old);
	}
}

// Statistics adds a user provided Stats struct to the parser to allow
// the user to process the results after the parsing has finished.
// Also the key for the "no match" counter is set.
//
// Example usage:
//
//	input := "input"
//	stats := Stats{}
//	_, err := Parse("input-file", []byte(input), Statistics(&stats, "no match"))
//	if err != nil {
//	    log.Panicln(err)
//	}
//	b, err := json.MarshalIndent(stats.ChoiceAltCnt, "", "  ")
//	if err != nil {
//	    log.Panicln(err)
//	}
//	fmt.Println(string(b))
function statistics(stats:Stats, choiceNoMatch :String): Option {
	return function(p :Parser): Option {
		var oldStats = p.stats;
		p.stats = stats;
		var oldChoiceNoMatch = p.choiceNoMatch;
		p.choiceNoMatch = choiceNoMatch;
		if (p.stats.ChoiceAltCnt == null) {
			p.stats.ChoiceAltCnt = new Map<String, Map<String, Int>>();
		}
		return statistics(oldStats, oldChoiceNoMatch);
	}
}

// Debug creates an Option to set the debug flag to b. When set to true,
// debugging information is printed to stdout while parsing.
//
// The default is false.
function debug(b:Bool):Option {
	return function(p:Parser):Option {
		var old = p.debug;
		p.debug = b;
		return debug(old);
	}
}

// Memoize creates an Option to set the memoize flag to b. When set to true,
// the parser caches named rule results by source position and parse mode.
// This avoids repeated named-rule evaluation during backtracking at the
// expense of additional memory.
//
// The default is false.
function memoize(b:Bool):Option {
	return function(p:Parser):Option {
		var old = p.memoize;
		p.memoize = b;
		return memoize(old);
	}
}


// position records a position in the text.
typedef Position = {
    var line:Int;
    var col:Int;
    var offset:Int;
}

function positionToString(p:Position):String {
    return '${p.line}:${p.col} [${p.offset}]';
}

// savepoint stores all state required to go back to this point in the
// parser.
typedef Savepoint = {
    > Position,
    var rn:RuneA;
    var w:Int;
}

function positionClone(x:Position):Position {
	return {
		line: x.line,
		col: x.col,
		offset: x.offset,
	};
}

function savepointClone(x:Savepoint):Savepoint {
	return {
		line: x.line,
		col: x.col,
		offset: x.offset,
		rn: x.rn,
		w: x.w
	};
}

typedef Current = {
	var pos:Position; // start position of the match
	var text:String; // raw text of the match // []byte
	var data:ParserCustomData;
}

enum AllExpr {
	ActionExpr(e:ActionExpr);
	AndCodeExpr(e:AndCodeExpr);
	AndExpr(e:AndExpr);
	AndLogicalExpr(e:AndLogicalExpr);
	AnyMatcher(e:AnyMatcher);

	CharClassMatcher(e:CharClassMatcher);
	ChoiceExpr(e:ChoiceExpr);
	CodeExpr(e:CodeExpr);

	LabeledExpr(e:LabeledExpr);
	LitMatcher(e:LitMatcher);
	NotCodeExpr(e:NotCodeExpr);
	NotExpr(e:NotExpr);
	NotLogicalExpr(e:NotLogicalExpr);
	OneOrMoreExpr(e:OneOrMoreExpr);
	RecoveryExpr(e:RecoveryExpr);
	RuleRefExpr(e:RuleRefExpr);
	SeqExpr(e:SeqExpr);
	ThrowExpr(e:ThrowExpr);
	ZeroOrMoreExpr(e:ZeroOrMoreExpr);
	ZeroOrOneExpr(e:ZeroOrOneExpr);
}

typedef Grammar = {
	var rules:Array<Rule>;
}

typedef Rule = {
	var name:String;
	@:optional var displayName:String;
	var expr:AllExpr;
	@:optional var varExists:Bool;
}

typedef ChoiceExpr = {
	var alternatives:Array<AllExpr>;
}

typedef ActionExpr = {
	var expr:AllExpr;
	var run:(p:Parser) -> Any;
}

typedef RecoveryExpr = {
	var expr:AllExpr;
	var recoverExpr:AllExpr;
	var failureLabel:Array<String>;
}

typedef SeqExpr = {
	var exprs:Array<AllExpr>;
}

typedef ThrowExpr = {
	var label:String;
}

typedef LabeledExpr = {
	var label:String;
	var expr:AllExpr;
	var textCapture:Bool;
}

typedef Expr = {
	var expr:AllExpr;
}

typedef AndExpr = Expr;
typedef NotExpr = Expr;
typedef AndLogicalExpr = Expr;
typedef NotLogicalExpr = Expr;
typedef ZeroOrOneExpr = Expr;
typedef ZeroOrMoreExpr = Expr;
typedef OneOrMoreExpr = Expr;

typedef RuleRefExpr = {
	var name:String;
}


typedef AndCodeExpr = {
	var run:(p:Parser) -> Bool;
}

typedef NotCodeExpr = {
	var run:(p:Parser) -> Bool;
}

typedef LitMatcher = {
	var val:String;
	var want:String;
	@:optional var ignoreCase:Bool;
}

typedef CodeExpr = {
	var run:(p:Parser) -> Any;
	var notSkip:Bool;
}

typedef UnicodeRangeChar = Array<UInt>;

typedef UnicodeRangeTable = {
	var R16:Array<UnicodeRangeChar>;
	var R32:Array<UnicodeRangeChar>;
	var LatinOffset:Int;
}

typedef CharClassMatcher = {
	var val:String;
	@:optional var chars:Array<UInt>;
	@:optional var ranges:Array<UInt>;
	@:optional var classes:Array<UnicodeRangeTable>;
	@:optional var ignoreCase:Bool;
	@:optional var inverted:Bool;
}

typedef AnyMatcher = {}

var unicodeLinearMax = 18;
var unicodeMaxLatin1 = 0x00FF;

private function unicodeIs16(ranges:Array<UnicodeRangeChar>, r:UInt):Bool {
	if (ranges.length <= unicodeLinearMax || r <= unicodeMaxLatin1) {
		for (range in ranges) {
			if (r < range[0]) {
				return false;
			}
			if (r <= range[1]) {
				return range[2] == 1 || (r - range[0]) % range[2] == 0;
			}
		}
		return false;
	}

	var lo = 0;
	var hi = ranges.length;
	while (lo < hi) {
		var middle = lo + Std.int((hi - lo) / 2);
		var range = ranges[middle];
		if (range[0] <= r && r <= range[1]) {
			return range[2] == 1 || (r - range[0]) % range[2] == 0;
		}
		if (r < range[0]) {
			hi = middle;
		} else {
			lo = middle + 1;
		}
	}
	return false;
}

private function unicodeIs32(ranges:Array<UnicodeRangeChar>, r:UInt):Bool {
	if (ranges.length <= unicodeLinearMax) {
		for (range in ranges) {
			if (r < range[0]) {
				return false;
			}
			if (r <= range[1]) {
				return range[2] == 1 || (r - range[0]) % range[2] == 0;
			}
		}
		return false;
	}

	var lo = 0;
	var hi = ranges.length;
	while (lo < hi) {
		var middle = lo + Std.int((hi - lo) / 2);
		var range = ranges[middle];
		if (range[0] <= r && r <= range[1]) {
			return range[2] == 1 || (r - range[0]) % range[2] == 0;
		}
		if (r < range[0]) {
			hi = middle;
		} else {
			lo = middle + 1;
		}
	}
	return false;
}

function unicodeIs(rangeTable:UnicodeRangeTable, r:UInt):Bool {
	var r16 = rangeTable.R16;
	if (r16.length > 0 && r <= r16[r16.length - 1][1]) {
		return unicodeIs16(r16, r);
	}
	var r32 = rangeTable.R32;
	return r32.length > 0 && r >= r32[0][0] && unicodeIs32(r32, r);
}

// errList cumulates the errors found by the parser.
class ErrList {
	var lst:Array<ParserError> = [];

	public function add(err:ParserError) {
		this.lst.push(err);
	}

	function err():ErrList {
		if (this.lst.length == 0) {
			return null;
		}
		this.dedupe();
		return this;
	}

	public var length(get, null):Int;

	function get_length() {
		return this.lst.length;
	}

	function dedupe() {
		// TODO: 去重
		// var cleaned []error
		// set := make(map[string]bool)
		// for _, err := range *e {
		// 	if msg := err.Error(); !set[msg] {
		// 		set[msg] = true
		// 		cleaned = append(cleaned, err)
		// 	}
		// }
		// *e = cleaned
	}

	function error():String {
		switch (this.lst.length) {
			case 0:
				return "";
			case 1:
				return parserErrorToString(this.lst[0]);
			default:
				var arr = [];
				for (i in this.lst) {
					arr.push(parserErrorToString(i));
				}
				return arr.join('\n');
		}
	}

	public function toException():Exception {
		if (this.lst.length == 0) {
			return null;
		}
		return new Exception('Parse Failed:\n' + this.err().error() + '\n');
	}

	public function new() {}
}

// parserError wraps an error with a prefix indicating the rule in which
// the error occurred. The original error is stored in the Inner field.
typedef ParserError = {
	var Inner:Exception;
	var pos:Position;
	var prefix:String;
	var expected:Array<String>;
}

// Error returns the error message.
function parserErrorToString(p:ParserError):String {
	return p.prefix + ": " + p.Inner.toString();
}

typedef ResultTuple = {
	var v:Any;
	var b:Bool;
	var end:Savepoint;
}

var _choiceNoMatch = -1; // const

// Stats stores some statistics, gathered during parsing
typedef Stats = {
	// ExprCnt counts the number of expressions processed during parsing
	// This value is compared to the maximum number of expressions allowed
	// (set by the MaxExpressions option).
	var ExprCnt: Uint64A;

	// ChoiceAltCnt is used to count for each ordered choice expression,
	// which alternative is used how may times.
	// These numbers allow to optimize the order of the ordered choice expression
	// to increase the performance of the parser
	//
	// The outer key of ChoiceAltCnt is composed of the name of the rule as well
	// as the line and the column of the ordered choice.
	// The inner key of ChoiceAltCnt is the number (one-based) of the matching alternative.
	// For each alternative the number of matches are counted. If an ordered choice does not
	// match, a special counter is incremented. The name of this counter is set with
	// the parser option Statistics.
	// For an alternative to be included in ChoiceAltCnt, it has to match at least once.
	var ChoiceAltCnt: Map<String, Map<String, Int>>;
}

class Parser {
	public var filename:String;
	public var pt:Savepoint;
	public var cur:Current;

	public var data:String;
	public var errs = new ErrList(); // x *errList

	public var depth:Int;
	public var recover:Bool;
	public var debug:Bool;

	public var memoize:Bool;
	// Memoization tables indexed directly by source offset. Lookahead parsing
	// uses a separate table because actions are skipped there.
	public var memo1:Array<Map<String, ResultTuple>>;
	public var memo2:Array<Map<String, ResultTuple>>;

	// rules table, maps the rule identifier to the rule node
	public var rules:Map<String, Rule>;
	public var rulesArray:Array<Rule>;
	// variables stack, map of label to value
	public var vstack = new GenericStack<Map<String, Any>>();
	// rule stack, allows identification of the current rule in errors
	public var rstack = new GenericStack<Rule>();

	// parse fail
	public var maxFailPos:Position;
	public var maxFailExpected:Array<String>;
	public var maxFailInvertExpected:Bool = false;
	public var noMatchErrorFormatter:NoMatchErrorFormatter;

	// max number of expressions to be parsed
	public var maxExprCnt:Uint64A;
	// entrypoint for the parser
	public var entrypoint:String;

	public var allowInvalidUTF8:Bool;

	// *Stats
	public var stats:Stats;

	public var choiceNoMatch:String;
	// recovery expression stack, keeps track of the currently available recovery expression, these are traversed in reverse
	public var recoveryStack = new GenericStack<Map<String, Any>>();

	public var _errPos: Position;
	// skip code stack
	public var scStack = new GenericStack<Bool>();
	// save point stack
	// public var spStack = new GenericStack<Savepoint>();

	// newParser creates a parser with the specified input source and options.
	public function new(filename:String, b:String, opts:Array<Option>) {
		var stats:Stats = {
			ExprCnt: 0,
			ChoiceAltCnt: new Map<String, Map<String, Int>>(),
		}

		this.filename = filename;
		this.data = b;
		this.pt = {
			line: 1,
			col: 0,
			offset: 0,
			rn: 0,
			w: 0,
		};
		this.recover = true;
		this.cur = {
			pos: {line: 0, col: 0, offset: 0},
			text: null,
			data: new ParserCustomData(),
		};

		this.maxFailPos = {col: 1, line: 1, offset: 0};
		this.maxFailExpected = []; // make([]string, 0, 20),
		this.stats = stats;
		// start rule is rule [0] unless an alternate entrypoint is specified
		this.entrypoint = "Input";

		// p.spStack.init(5)
		this.setOptions(opts);

		if (this.maxExprCnt == 0) {
			this.maxExprCnt = 2147483647; // p.maxExprCnt = 2147483647;
		}
	}

	// setOptions applies the options to the parser.
	function setOptions(opts:Array<Option>) {
		for (opt in opts) {
			opt(this);
		}
	}

	// setCustomData to the parser.
	function setCustomData(data:ParserCustomData) {
		this.cur.data = data;
	}

	// setCustomData to the parser.
	function checkSkipCode():Bool {
		return this.scStack.first();
	}

	// push a variable set on the vstack.
	function pushV() {
		this.vstack.add(new Map());
	}

	// pop a variable set from the vstack.
	function popV() {
		this.vstack.pop();
	}

	// push a recovery expression with its labels to the recoveryStack
	function pushRecovery(labels:Array<String>, expr:Any) {
		var m = new Map<String, Any>();
		// m := make(map[string]any, len(labels))
		for (fl in labels) {
			m[fl] = expr;
		}
		this.recoveryStack.add(m);
	}

	// pop a recovery expression from the recoveryStack
	function popRecovery() {
		this.recoveryStack.pop();
	}

	function print(prefix:String, s:String):String {
		if (!this.debug) {
			return s;
		}
		trace('${prefix} ${this.pt.line}:${this.pt.col}:${this.pt.offset}: ${s} [${this.pt.rn}]');
		// fmt.Printf("%s %d:%d:%d: %s [%#U]\n", prefix, p.pt.line, p.pt.col, p.pt.offset, s, p.pt.rn)
		return s;
	}

	function printIndent(mark:String, s:String):String {
		// return this.print(strings.Repeat(" ", p.depth)+mark, s)
		return this.print(StringTools.lpad('', ' ', this.depth), s);
	}

	function in_(s:String):String {
		var res = this.printIndent(">", s);
		this.depth++;
		return res;
	}

	function out(s:String):String {
		this.depth--;
		return this.printIndent("<", s);
	}

	public function addErr(err:Exception) {
		// p.addErrAt(err, p.pt.position,[])
		this.addErrAt(err, this.pt, []);
	}

	public function addErrAt(err:Exception, pos:Position, expected:Array<String>) {
		var buf = new BytesBuffer();
		if (this.filename != "") {
			buf.addString(this.filename);
		}
		// if buf.Len() > 0 { 不确定
		if (buf.length > 0) {
			buf.addString(':');
		}
		// buf.WriteString(fmt.Sprintf("%d:%d (%d)", pos.line, pos.col, pos.offset))
		buf.addString('${pos.line}:${pos.col} (${pos.offset})');
		if (this.rstack.first() != null) {
			if (buf.length > 0) {
				buf.addString(": ");
			}
			var rule = this.rstack.first();
			if (rule.displayName != "") {
				buf.addString("rule " + rule.displayName);
			} else {
				buf.addString("rule " + rule.name);
			}
		}
		// pe := &parserError{Inner: err, pos: pos, prefix: buf.String(), expected: expected}
		// p.errs.add(pe)
		this.errs.add({
			Inner: err,
			pos: pos,
			prefix: buf.getBytes().toString(),
			expected: expected,
		});
	}

	function buildNoMatchError(pos:Position, expected:Array<String>):Exception {
		if (this.noMatchErrorFormatter != null) {
			var err = this.noMatchErrorFormatter(pos, this.data, expected);
			if (err != null) {
				return err;
			}
		}

		return new Exception("no match found, expected: " + listJoin(expected, ", ", "or"));
	}

	function failAt(fail:Bool, pos:Position, want:String) {
		// process fail if parsing fails and not inverted or parsing succeeds and invert is set
		if (fail == this.maxFailInvertExpected) {
			if (pos.offset < this.maxFailPos.offset) {
				return;
			}

			if (pos.offset > this.maxFailPos.offset) {
				this.maxFailPos = pos;
				// p.maxFailExpected = p.maxFailExpected[:0]
				this.maxFailExpected = [];
			}

			if (this.maxFailInvertExpected) {
				want = "!" + want;
			}
			this.maxFailExpected.push(want);
		}
	}

	// read advances the parser to the next rune.
	function read() {
		this.pt.offset += this.pt.w;
		var first = this.data.charCodeAt(this.pt.offset);
		if (first == null) {
			this.pt.rn = 0;
			this.pt.w = 0;
		} else {
			this.pt.rn = first;
			this.pt.w = 1;
			#if utf16
			if (first >= 0xD800 && first <= 0xDBFF) {
				var second = this.data.charCodeAt(this.pt.offset + 1);
				if (second != null && second >= 0xDC00 && second <= 0xDFFF) {
					this.pt.rn = ((first - 0xD800) << 10) + second - 0xDC00 + 0x10000;
					this.pt.w = 2;
				}
			}
			#end
		}
		this.pt.col ++;

		if (this.pt.rn == 0x0A) {
			this.pt.line++;
			this.pt.col = 0;
		}

		// TODO 需要重做，rn为空不代表utf8的invalid
        // if (rn == '') {
        //     if (!this.allowInvalidUTF8) {
        //         this.addErr(Errors.errInvalidEncoding);
        //     }
        // }
	}

	// restore parser position to the savepoint pt.
	function restore(pt:Savepoint) {
		if (pt.offset == this.pt.offset) {
			if (this.debug) {
				this.out(this.in_("restore"));
			}
			return;
		}
		this.pt = pt;
		if (this.debug) {
			this.out(this.in_("restore"));
		}
	}

	function sliceFrom(start:Savepoint):String {
		// return p.data[start.position.offset:p.pt.position.offset]
		return this.data.substring(start.offset, this.pt.offset);
	}

	function sliceFromOffset(offset:Int):String {
		return this.data.substring(offset, this.pt.offset);
	}

	function getMemoizedRule(rule:Rule):ResultTuple {
		var memo = this.checkSkipCode() ? this.memo2 : this.memo1;
		if (memo == null) {
			return null;
		}
		var entries = memo[this.pt.offset];
		return entries == null ? null : entries[rule.name];
	}

	function setMemoizedRule(pt:Savepoint, rule:Rule, val:Any, ok:Bool, end:Savepoint) {
		var memo:Array<Map<String, ResultTuple>>;
		if (this.checkSkipCode()) {
			if (this.memo2 == null) {
				this.memo2 = [];
			}
			memo = this.memo2;
		} else {
			if (this.memo1 == null) {
				this.memo1 = [];
			}
			memo = this.memo1;
		}
		var entries = memo[pt.offset];
		if (entries == null) {
			entries = new Map<String, ResultTuple>();
			memo[pt.offset] = entries;
		}
		var result = entries[rule.name];
		if (result == null) {
			entries[rule.name] = {v: val, b: ok, end: savepointClone(end)};
		} else {
			result.v = val;
			result.b = ok;
			result.end = savepointClone(end);
		}
	}

	public function buildRulesTable(g:Grammar):Void {
		this.rules = new Map<String, Rule>();
		for (r in g.rules) {
			this.rules[r.name] = r;
		}
	}

	public function parse(grammar:Grammar):Any {
		if (grammar == null) {
			grammar = g;
		}
		if (grammar.rules.length == 0) {
			this.addErr(Errors.errNoRule);
			throw this.errs.toException();
		}

		this.rulesArray = grammar.rules;
		this.buildRulesTable(grammar);
		// var recoverDefer = () -> {};

		if (this.recover) {
			// panic can be used in action code to stop parsing immediately
			// and return the panic as an error.
		}

		var startRule = this.rules[this.entrypoint];

		if (startRule == null) {
			this.addErr(Errors.errInvalidEntrypoint);
			// return nil, p.errs.err()
			throw this.errs.toException();
		}

		var ret:RetValOk = {val: null, ok: false};

		try {
			this.read(); // advance to first rune
			ret = this.parseRuleWrap(startRule);
		} catch (e:haxe.Exception) {
			// if p.debug {
			// 	defer p.out(p.in("panic handler"))
			// }
			// switch e := e.(type) {
			// case error:
			this.addErr(e);
			// default:
			// 	p.addErr(fmt.Errorf("%v", e))
			// }
			throw this.errs.toException();
		}

		if (!ret.ok) {
			if (this.errs.length == 0) {
				// If parsing fails, but no errors have been recorded, the expected values
				// for the farthest parser position are returned as error.
				var maxFailExpectedMap = new Map<String, Any>();
				for (v in this.maxFailExpected) {
					// maxFailExpectedMap[v] = struct{}{}
					maxFailExpectedMap[v] = {};
				}

				var expected = [];
				var eof = false;
				if (maxFailExpectedMap["!."] != null) {
					// delete(maxFailExpectedMap, "!.")
					maxFailExpectedMap.remove("!.");
					eof = true;
				}
				for (k in maxFailExpectedMap.keys()) {
					expected.push(k);
				}

				// sort.Strings(expected)
				expected.sort((a, b) -> a > b ? 1 : -1);

				if (eof) {
					expected.push("EOF");
				}
				this.addErrAt(this.buildNoMatchError(this.maxFailPos, expected), this.maxFailPos, expected);
			}
			throw this.errs.toException();
		}

		// return val, p.errs.err()
		if (this.errs.length != 0) {
			throw this.errs.toException();
		}
		return ret.val;
	}

	function parseRuleMemoize(rule:Rule):RetValOk {
		var result = this.getMemoizedRule(rule);
		if (result != null) {
			this.restore(savepointClone(result.end));
			return {val: result.v, ok: result.b};
		}

		var startMark = savepointClone(this.pt);
		var ret = this.parseRule(rule);
		this.setMemoizedRule(startMark, rule, ret.val, ret.ok, this.pt);

		return ret;
	}

	function parseRuleWrap(rule:Rule):RetValOk {
		// if p.debug {
		//     defer p.out(p.in("parseRule " + rule.name))
		// }
		var ret:RetValOk = {val: null, ok: false};

		var startMark = savepointClone(this.pt);

		//
		if (this.memoize) {
			ret = this.parseRuleMemoize(rule);
		} else {
			ret = this.parseRule(rule);
		}

		if (ret.ok && this.debug) {
			this.printIndent("MATCH", this.sliceFrom(startMark));
		}
		return ret;
	}

	function parseRule(rule:Rule):RetValOk {
		if (rule.varExists && !this.checkSkipCode()) {
			this.rstack.add(rule);
			this.pushV();
			var ret = this.parseExprWrap(rule.expr);
			this.popV();
			this.rstack.pop();
			return ret;
		} else {
			this.rstack.add(rule);
			var ret = this.parseExprWrap(rule.expr);
			this.rstack.pop();
			return ret;
		}
	}


	function parseExprWrap(expr:AllExpr):RetValOk {
		return this.parseExpr(expr);
	}

	function parseExpr(expr:AllExpr):RetValOk {
		this.stats.ExprCnt++;
		if (this.stats.ExprCnt > this.maxExprCnt) {
			throw Errors.errMaxExprCnt;
		}

		var ret:RetValOk;

		switch (expr) {
			case ActionExpr(expr):
				ret = this.parseActionExpr(expr);
			case AndCodeExpr(expr):
				ret = this.parseAndCodeExpr(expr);
			case AndExpr(expr):
				ret = this.parseAndExpr(expr);
			case AndLogicalExpr(expr):
				ret = this.parseAndLogicalExpr(expr);
			case AnyMatcher(expr):
				ret = this.parseAnyMatcher(expr);
			case CharClassMatcher(expr):
				ret = this.parseCharClassMatcher(expr);
			case ChoiceExpr(expr):
				ret = this.parseChoiceExpr(expr);
			case CodeExpr(expr):
				ret = this.parseCodeExpr(expr);
			case LabeledExpr(expr):
				ret = this.parseLabeledExpr(expr);
			case LitMatcher(expr):
				ret = this.parseLitMatcher(expr);
			case NotCodeExpr(expr):
				ret = this.parseNotCodeExpr(expr);
			case NotExpr(expr):
				ret = this.parseNotExpr(expr);
			case NotLogicalExpr(expr):
				ret = this.parseNotLogicalExpr(expr);
			case OneOrMoreExpr(expr):
				ret = this.parseOneOrMoreExpr(expr);
			case RecoveryExpr(expr):
				ret = this.parseRecoveryExpr(expr);
			case RuleRefExpr(expr):
				ret = this.parseRuleRefExpr(expr);
			case SeqExpr(expr):
				ret = this.parseSeqExpr(expr);
			case ThrowExpr(expr):
				ret = this.parseThrowExpr(expr);
			case ZeroOrMoreExpr(expr):
				ret = this.parseZeroOrMoreExpr(expr);
			case ZeroOrOneExpr(expr):
				ret = this.parseZeroOrOneExpr(expr);
			default:
				throw new Exception("unknown expression type");
				// panic(fmt.Sprintf("unknown expression type %T", expr))
		}
		return ret;
	}

	function parseActionExpr(act:ActionExpr):RetValOk {
		if (this.checkSkipCode()) {
			var ret = this.parseExprWrap(act.expr);
			return { val:null, ok: ret.ok };
		}

		var start = savepointClone(this.pt);
		var ret = this.parseExprWrap(act.expr);
		if (ret.ok) {
			// this.cur.pos = start.position;
			this.cur.pos = positionClone(start);
			this.cur.text = this.sliceFrom(start);
			this._errPos = start;
			var val = act.run(this);
			this._errPos = null;
			ret.val = val;
		}
		if (ret.ok && this.debug) {
			this.printIndent("MATCH", this.sliceFrom(start));
		}
		return ret;
	}

	function parseAndCodeExpr(and:AndCodeExpr):RetValOk {
		var ok = and.run(this);
		return {val: null, ok: ok};
	}

	function parseAndExpr(and:AndExpr):RetValOk {
		return this.parseAndExprBase(and, false);
	}

	function parseAndLogicalExpr(and:AndExpr):RetValOk {
		return this.parseAndExprBase(and, true);
	}

	function parseAndExprBase(and:AndExpr, logical:Bool):RetValOk {
		var pt = savepointClone(this.pt);

		this.scStack.add(true);
		var ret = this.parseExprWrap(and.expr);
		this.scStack.pop();

		var matchedOffset = this.pt.offset;
		this.restore(pt);

		if (logical) {
			return {val: null, ok: ret.ok && this.pt.offset != matchedOffset}
		}
		return {val: null, ok: ret.ok};
	}

	function parseAnyMatcher(act:AnyMatcher):RetValOk {
		// if p.pt.rn == utf8.RuneError && p.pt.w == 0 {
		if (this.pt.rn == 0 && this.pt.w == 0) {
			// EOF - see utf8.DecodeRune
			this.failAt(false, this.pt, ".");
			return {val: null, ok: false};
		}
		this.failAt(true, this.pt, ".");
		this.read();
		return {val: null, ok: true};
	}

	function parseCharClassMatcher(chr:CharClassMatcher):RetValOk {
		var cur = this.pt.rn;

		// can't match EOF
		if (cur == 0 && this.pt.w == 0) { // see utf8.DecodeRune
			this.failAt(false, this.pt, chr.val);
			return {val: null, ok: false};
		}

		if (chr.ignoreCase) {
			cur = String.fromCharCode(cur).toLowerCase().charCodeAt(0);
		}

		if (chr.chars != null) {
			// try to match in the list of available chars
			for (rn in chr.chars) {
				if (rn == cur) {
					if (chr.inverted) {
						this.failAt(false, this.pt, chr.val);
						return {val: null, ok: false};
					}
					this.failAt(true, this.pt, chr.val);
					this.read();
					return {val: null, ok: true};
				}
			}
		}

		if (chr.ranges != null) {
			// try to match in the list of ranges
			for (i in 0...chr.ranges.length) {
				if (i % 2 != 0)
					continue;
				if (cur >= chr.ranges[i] && cur <= chr.ranges[i + 1]) {
					if (chr.inverted) {
						this.failAt(false, this.pt, chr.val);
						return {val: null, ok: false};
					}
					this.failAt(true, this.pt, chr.val);
					this.read();
					return {val: null, ok: true};
				}
			}
		}

		// try to match in the list of Unicode classes
		if (chr.classes != null) {
			for (cl in chr.classes) {
				if (unicodeIs(cl, cur)) {
					if (chr.inverted) {
						this.failAt(false, this.pt, chr.val);
						return {val: null, ok: false};
					}
					this.failAt(true, this.pt, chr.val);
					this.read();
					return {val: null, ok: true};
				}
			}
		}

		if (chr.inverted) {
			this.failAt(true, this.pt, chr.val);
			this.read();
			return {val: null, ok: true};
		}
		this.failAt(true, this.pt, chr.val);
		return {val: null, ok: false};
	}

	function incChoiceAltCnt(ch:ChoiceExpr, altI:Int) {
		// choiceIdent := fmt.Sprintf("%s %d:%d", p.rstack.first().name, ch.pos.line, ch.pos.col)
		var choiceIdent:String = '${this.rstack.first().name}';

		var m = this.stats.ChoiceAltCnt[choiceIdent];
		if (m == null) {
			m = new Map<String, Int>();
			this.stats.ChoiceAltCnt[choiceIdent] = m;
		}
		// We increment altI by 1, so the keys do not start at 0
		var alt:String = Std.string(altI + 1);
		if (altI == _choiceNoMatch) {
			alt = this.choiceNoMatch;
		}
		m[alt]++;
	}

	function parseChoiceExpr(ch:ChoiceExpr):RetValOk {
		// for altI, alt := range ch.alternatives {
		for (altI in 0...ch.alternatives.length) {
			var alt = ch.alternatives[altI];

			var ret = this.parseExprWrap(alt);
			if (ret.ok) {
				this.incChoiceAltCnt(ch, altI);
				return ret;
			}
		}
		this.incChoiceAltCnt(ch, _choiceNoMatch);
		return {val: null, ok: false};
	}

	function parseLabeledExpr(lab:LabeledExpr):RetValOk {
		var startOffset = this.pt.offset;
		var ret = this.parseExprWrap(lab.expr);
		if (ret.ok && lab.label != "" && !this.checkSkipCode()) {
			var m = this.vstack.first();
			m[lab.label] = lab.textCapture ? this.sliceFromOffset(startOffset) : ret.val;
		}
		return ret;
	}

	function parseCodeExpr(code:CodeExpr):RetValOk {
		if (!code.notSkip && this.checkSkipCode()) {
			return {val: null, ok: true};
		}
		return {val: code.run(this), ok: true};
	}

	function parseLitMatcher(lit:LitMatcher):RetValOk {
		var start = savepointClone(this.pt);
		var litVal = lit.val;
		if (lit.ignoreCase) {
			litVal = lit.val.toLowerCase();
		}

		for (want in new StringIteratorUnicode(litVal)) {
			var cur = this.pt.rn;
			if (cur != want) {
				this.failAt(false, start, lit.want);
				this.restore(start);
				return {val: null, ok: false};
			}
			this.read();
		}
		this.failAt(true, start, lit.want);
		return {val: null, ok: true};
	}

	function parseNotCodeExpr(code:NotCodeExpr):RetValOk {
		return {val: null, ok: !code.run(this)};
	}

	function parseNotExpr(not:NotExpr):RetValOk {
		return this.parseNotExprBase(not, false);
	}

	function parseNotLogicalExpr(not:NotLogicalExpr):RetValOk {
		return this.parseNotExprBase(not, true);
	}

	function parseNotExprBase(not:NotExpr, logical:Bool):RetValOk {
		var pt = savepointClone(this.pt);
		this.maxFailInvertExpected = !this.maxFailInvertExpected;

		this.scStack.add(true);
		var ret = this.parseExprWrap(not.expr);
		this.scStack.pop();

		this.maxFailInvertExpected = !this.maxFailInvertExpected;
		var matchedOffset = this.pt.offset;
		this.restore(pt);

		if (logical) {
			return {val: null, ok: !ret.ok && this.pt.offset != matchedOffset}
		}
		return {val: null, ok: !ret.ok};
	}

	function parseOneOrMoreExpr(expr:OneOrMoreExpr):RetValOk {
		var vals:Array<Any> = null;
		var matched = false;

		while (true) {
			var ret = this.parseExprWrap(expr.expr);
			if (!ret.ok) {
				if (vals != null) {
					return {val: vals, ok: matched};
				}
				return {val: null, ok: matched};
			}
			matched = true;
			if (ret.val != null) {
				if (vals == null) {
					vals = [];
				}
				vals.push(ret.val);
			}
		}
	}

	function parseRecoveryExpr(recover:RecoveryExpr):RetValOk {
		this.pushRecovery(recover.failureLabel, recover.recoverExpr);
		var ret = this.parseExprWrap(recover.expr);
		this.popRecovery();
		return ret;
	}

	function parseRuleRefExpr(ref:RuleRefExpr):RetValOk {
		var rule = this.rules[ref.name];
		if (rule == null) {
			// p.addErr(fmt.Errorf("undefined rule: %s", ref.name))
			this.addErr(new Exception('undefined rule: ${ref.name}'));
			return {val: null, ok: false};
		}
		return this.parseRuleWrap(rule);
	}


	function parseSeqExpr(seq:SeqExpr):RetValOk {
		var pt = savepointClone(this.pt);
		for (expr in seq.exprs) {
			var ret = this.parseExprWrap(expr);
			if (!ret.ok) {
				this.restore(pt);
				return {val: null, ok: false};
			}
		}
		return {val: null, ok: true};
	}

	function parseThrowExpr(expr:ThrowExpr):RetValOk {
		// for i := len(this.recoveryStack) - 1; i >= 0; i-- {
		// 	if recoverExpr, ok := this.recoveryStack[i][expr.label]; ok {
		for (i in this.recoveryStack) {
			var recoverExpr = i[expr.label];
			if (recoverExpr != null) {
				return this.parseExprWrap(recoverExpr);
			}
		}
		return {val: null, ok: false};
	}

	function parseZeroOrMoreExpr(expr:ZeroOrMoreExpr):RetValOk {
		var vals:Array<Any> = null;
		while (true) {
			var ret = this.parseExprWrap(expr.expr);
			if (!ret.ok) {
				if (vals != null) {
					return {val: vals, ok: true};
				}
				return {val: null, ok: true};
			}
			if (ret.val != null) {
				if (vals == null) {
					vals = [];
				}
				vals.push(ret.val);
			}
		}
	}

	function parseZeroOrOneExpr(expr:ZeroOrOneExpr):RetValOk {
		var ret = this.parseExprWrap(expr.expr);
		// whether it matched or not, consider it a match
		return {val: ret.val, ok: true};
	}
}

function listJoin(list:Array<String>, sep:String, lastSep:String):String {
    switch (list.length) {
        case 0:
            return "";
        case 1:
            return list[0];
        default:
            // return strings.Join(list[:len(list)-1], sep) + " " + lastSep + " " + list[len(list)-1]
            return list.slice(0, list.length - 1).join(sep) + " " + lastSep + " " + list[list.length - 1];
    }
}
