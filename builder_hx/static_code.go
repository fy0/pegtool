//go:generate go run ../bootstrap/cmd/static_code_generator/main.go -- $GOFILE generated_$GOFILE staticCode

//go:build static_code
// +build static_code

package builder

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"
)

// IMPORTANT: All code below this line is added to the parser as static code

class Errors {
    public static var errNoRule:Exception = new Exception("grammar has no rule");
    public static var errInvalidEntrypoint:Exception = new Exception("invalid entrypoint");
    public static var errInvalidEncoding:Exception = new Exception("invalid encoding");
    public static var errMaxExprCnt:Exception = new Exception("max number of expressions parsed");
}

// alternative of Uint64
typedef Uint64A = UInt;
typedef RuneA = UInt;

// Option is a function that can set an option on the parser. It returns
// the previous setting as an Option.
// typedef Option = (p:Parser) -> Option; // 注: 这是个函数
typedef Option = (p:Parser) -> Any; // 注: 这是个函数，因static编译需求返回Any

// MaxExpressions creates an Option to stop parsing after the provided
// number of expressions have been parsed, if the value is 0 then the parser will
// parse for as many steps as needed (possibly an infinite number).
//
// The default for maxExprCnt is 0.
function MaxExpressions(maxExprCnt:Uint64A):Option {
	return function(p:Parser):Option {
		var oldMaxExprCnt = p.maxExprCnt;
		p.maxExprCnt = maxExprCnt;
		return MaxExpressions(oldMaxExprCnt);
	}
}

// Entrypoint creates an Option to set the rule name to use as entrypoint.
// The rule name must have been specified in the -alternate-entrypoints
// if generating the parser with the -optimize-grammar flag, otherwise
// it may have been optimized out. Passing an empty string sets the
// entrypoint to the first rule in the grammar.
//
// The default is to start parsing at the first rule in the grammar.
function Entrypoint(ruleName:String):Option {
	return function(p:Parser):Option {
		var oldEntrypoint = p.entrypoint;
		p.entrypoint = ruleName;
		if (ruleName == "") {
			p.entrypoint = g.rules[0].name;
		}
		return Entrypoint(oldEntrypoint);
	}
}

// ==template== {{ if not .Optimize }}
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
function Statistics(stats:Stats, choiceNoMatch :String): Option {
	return function(p :Parser): Option {
		var oldStats = p.stats;
		p.stats = stats;
		var oldChoiceNoMatch = p.choiceNoMatch;
		p.choiceNoMatch = choiceNoMatch;
		if (p.stats.ChoiceAltCnt == null) {
			p.stats.ChoiceAltCnt = new Map<String, Map<String, Int>>();
		}
		return Statistics(oldStats, oldChoiceNoMatch);
	}
}

// Debug creates an Option to set the debug flag to b. When set to true,
// debugging information is printed to stdout while parsing.
//
// The default is false.
function Debug(b:Bool):Option {
	return function(p:Parser):Option {
		var old = p.debug;
		p.debug = b;
		return Debug(old);
	}
}

// Memoize creates an Option to set the memoize flag to b. When set to true,
// the parser will cache all results so each expression is evaluated only
// once. This guarantees linear parsing time even for pathological cases,
// at the expense of more memory and slower times for typical cases.
//
// The default is false.
function Memoize(b:Bool):Option {
	return function(p:Parser):Option {
		var old = p.memoize;
		p.memoize = b;
		return Memoize(old);
	}
}

// {{ end }} ==template==
// AllowInvalidUTF8 creates an Option to allow invalid UTF-8 bytes.
// Every invalid UTF-8 byte is treated as a utf8.RuneError (U+FFFD)
// by character class matchers and is matched by the any matcher.
// The returned matched value, c.text and c.offset are NOT affected.
//
// The default is false.
function AllowInvalidUTF8(b:Bool):Option {
	return function(p:Parser):Option {
		var old = p.allowInvalidUTF8;
		p.allowInvalidUTF8 = b;
		return AllowInvalidUTF8(old);
	}
}

// Recover creates an Option to set the recover flag to b. When set to
// true, this causes the parser to recover from panics and convert it
// to an error. Setting it to false can be useful while debugging to
// access the full stack trace.
//
// The default is true.
function Recover(b:Bool):Option {
	return function(p:Parser):Option {
		var old = p.recover;
		p.recover = b;
		return Recover(old);
	}
}

// GlobalStore creates an Option to set a key to a certain value in
// the globalStore.
function GlobalStore(key:String, value:Any):Option {
	return function(p:Parser):Option {
		var old = p.cur.globalStore[key];
		p.cur.globalStore[key] = value;
		return GlobalStore(key, old);
	}
}

// ==template== {{ if or .GlobalState (not .Optimize) }}

// InitState creates an Option to set a key to a certain value in
// the global "state" store.
function InitState(key:String, value:Any):Option {
	return function(p:Parser):Option {
		var old = p.cur.state[key];
		p.cur.state[key] = value;
		return InitState(key, old);
	}
}

// {{ end }} ==template==

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

	// ==template== {{ if or .GlobalState (not .Optimize) }}
	// state is a store for arbitrary key,value pairs that the user wants to be
	// tied to the backtracking of the parser.
	// This is always rolled back if a parsing rule fails.
	var state:StoreDict;

	// {{ end }} ==template==
	// globalStore is a general store for the user to store arbitrary key-value
	// pairs that they need to manage and that they do not want tied to the
	// backtracking of the parser. This is only modified by the user and never
	// rolled back by the parser. It is always up to the user to keep this in a
	// consistent state.
	var globalStore:StoreDict;
}

typedef StoreDict = Map<String, Any>;

enum AllExpr {
	ChoiceExpr(e:ChoiceExpr);
	ActionExpr(e:ActionExpr);
	RecoveryExpr(e:RecoveryExpr);
	SeqExpr(e:SeqExpr);
	ThrowExpr(e:ThrowExpr);
	LabeledExpr(e:LabeledExpr);

	AndExpr(e:AndExpr);
	NotExpr(e:NotExpr);
	ZeroOrOneExpr(e:ZeroOrOneExpr);
	ZeroOrMoreExpr(e:ZeroOrMoreExpr);
	OneOrMoreExpr(e:OneOrMoreExpr);
	RuleRefExpr(e:RuleRefExpr);
	// ==template== {{ if or .GlobalState (not .Optimize) }}
	StateCodeExpr(e:StateCodeExpr);
	// {{ end }} ==template==
	AndCodeExpr(e:AndCodeExpr);
	NotCodeExpr(e:NotCodeExpr);

	AnyMatcher(e:AnyMatcher);
	CharClassMatcher(e:CharClassMatcher);
	CustomParserCodeExpr(e:CustomParserCodeExpr);
	LitMatcher(e:LitMatcher);
}

// the AST types...

// {{ if .Nolint }} nolint: structcheck {{else}} ==template== {{ end }}
typedef Grammar = {
	@:optional var pos:Position;
	var rules:Array<Rule>;
}

// {{ if .Nolint }} nolint: structcheck {{else}} ==template== {{ end }}
typedef Rule = {
	var pos:Position;
	var name:String;
	@:optional var displayName:String;
	var expr:AllExpr;

	// ==template== {{ if .LeftRecursion }}
	@:optional var leader:Bool;
	@:optional var leftRecursive:Bool;
	// {{ end }} ==template==
}

// {{ if .Nolint }} nolint: structcheck {{else}} ==template== {{ end }}
typedef ChoiceExpr = {
	var pos:Position;
	var alternatives:Array<Any>;
}

typedef RetValOk = {
	var val:Any;
	var ok:Bool;
}

typedef RetValErr = {
	var val:Any;
	var err:Exception;
}

typedef RetValOkErr = {
	var val:Any;
	var ok:Bool;
	var err:Exception;
}

typedef RetOkErr = {
	var ok:Bool;
	var err:Exception;
}

// {{ if .Nolint }} nolint: structcheck {{else}} ==template== {{ end }}
typedef ActionExpr = {
	var pos:Position;
	var expr:Any;
	var run:(p:Parser) -> RetValErr;
}

// {{ if .Nolint }} nolint: structcheck {{else}} ==template== {{ end }}
typedef RecoveryExpr = {
	var pos:Position;
	var expr:Any;
	var recoverExpr:Any;
	var failureLabel:Array<String>;
}

// {{ if .Nolint }} nolint: structcheck {{else}} ==template== {{ end }}
typedef SeqExpr = {
	var pos:Position;
	var exprs:Array<AllExpr>;
}

// {{ if .Nolint }} nolint: structcheck {{else}} ==template== {{ end }}
typedef ThrowExpr = {
	var pos:Position;
	var label:String;
}

// {{ if .Nolint }} nolint: structcheck {{else}} ==template== {{ end }}
typedef LabeledExpr = {
	var pos:Position;
	var label:String;
	var expr:AllExpr;
}

// {{ if .Nolint }} nolint: structcheck {{else}} ==template== {{ end }}
typedef Expr = {
	var pos:Position;
	var expr:AllExpr;
}

typedef AndExpr = Expr;
typedef NotExpr = Expr;
typedef ZeroOrOneExpr = Expr;
typedef ZeroOrMoreExpr = Expr;
typedef OneOrMoreExpr = Expr;

// {{ if .Nolint }} nolint: structcheck {{else}} ==template== {{ end }}
typedef RuleRefExpr = {
	var pos:Position;
	var name:String;
}

// ==template== {{ if or .GlobalState (not .Optimize) }}

// {{ if .Nolint }} nolint: structcheck {{else}} ==template== {{ end }}
typedef StateCodeExpr = {
	var pos:Position;
	var run:(p:Parser) -> RetValErr;
}

// {{ end }} ==template==

// {{ if .Nolint }} nolint: structcheck {{else}} ==template== {{ end }}
typedef AndCodeExpr = {
	var pos:Position;
	var run:(p:Parser) -> RetOkErr;
}

// {{ if .Nolint }} nolint: structcheck {{else}} ==template== {{ end }}
typedef NotCodeExpr = {
	var pos:Position;
	var run:(p:Parser) -> RetOkErr;
}

// {{ if .Nolint }} nolint: structcheck {{else}} ==template== {{ end }}
typedef LitMatcher = {
	var pos:Position;
	var val:String;
	var ignoreCase:Bool;
	var want:String;
}

// {{ if .Nolint }} nolint: structcheck {{else}} ==template== {{ end }}
typedef CustomParserCodeExpr = {
	var pos:Position;
	var run:(p:Parser) -> RetValOkErr;
}

// {{ if .Nolint }} nolint: structcheck {{else}} ==template== {{ end }}
typedef CharClassMatcher = {
	var pos:Position;
	var val:String;
	// basicLatinChars [128]bool // TODO
	var chars:Array<UInt>;
	var ranges:Array<UInt>;
	@:optional var classes:Array<Unicode.RangeTable>;
	var ignoreCase:Bool;
	var inverted:Bool;
}

typedef AnyMatcher = Position; // {{ if .Nolint }} nolint: structcheck {{else}} ==template== {{ end }}
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
				return this.lst[0].Inner.toString();
			default:
				var arr = [];
				for (i in this.lst) {
					arr.push(i.Inner.toString());
				}
				return arr.join('\n');
		}
	}

	public function toException():Exception {
		if (this.lst.length == 0) {
			return null;
		}
		return new Exception(this.err().error());
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
function parserErrorError(p:ParserError):String {
	return p.prefix + ": " + p.Inner.toString();
}

// {{ if .Nolint }} nolint: structcheck,deadcode {{else}} ==template== {{ end }}
typedef ResultTuple = {
	var v: Any;
	var b: Bool;
	var end: Savepoint;
}

// {{ if .Nolint }} nolint: varcheck {{else}} ==template== {{ end }}
var choiceNoMatch = -1; // const

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


// ==template== {{ if .LeftRecursion }}
typedef RuleWithExpsStack = {
	var rule: Rule;
	var estack: GenericStack<Any>;
}
// {{ end }} ==template==

// {{ if .Nolint }} nolint: structcheck,maligned {{else}} ==template== {{ end }}
class Parser {
	public var filename:String;
	public var pt:Savepoint;
	public var cur:Current;

	public var data:String;
	public var errs = new ErrList(); // x *errList

	public var depth:Int;
	public var recover:Bool;
	// ==template== {{ if not .Optimize }}
	public var debug:Bool;

	public var memoize:Bool;
	// {{ end }} ==template==
	// ==template== {{ if or .LeftRecursion (not .Optimize) }}
	// memoization table for the packrat algorithm:
	// map[offset in source] map[expression or rule] {value, match}
	public var memo:Map<Int, Map<Any, ResultTuple>>; // map[int]map[any]resultTuple

	// {{ end }} ==template==
	// rules table, maps the rule identifier to the rule node
	public var rules:Map<String, Rule>;
	// variables stack, map of label to value
	public var vstack = new GenericStack<Map<String, Any>>();
	// rule stack, allows identification of the current rule in errors
	public var rstack = new GenericStack<Rule>();

	// parse fail
	public var maxFailPos:Position;
	public var maxFailExpected:Array<String>;
	public var maxFailInvertExpected:Bool = false;

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
			// ==template== {{ if or .GlobalState (not .Optimize) }}
			state: new StoreDict(),
			// {{ end }} ==template==
			globalStore: new StoreDict(),
		};

		this.maxFailPos = {col: 1, line: 1, offset: 0};
		this.maxFailExpected = []; // make([]string, 0, 20),
		this.stats = stats;
		// start rule is rule [0] unless an alternate entrypoint is specified
		this.entrypoint = g.rules[0].name;

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

	// ==template== {{ if not .Optimize }}
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

	// {{ end }} ==template==

	function addErr(err:Exception) {
		// p.addErrAt(err, p.pt.position,[])
		this.addErrAt(err, this.pt, []);
	}

	function addErrAt(err:Exception, pos:Position, expected:Array<String>) {
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
		// rn, n := utf8.DecodeRune(p.data[p.pt.offset:])
		// p.pt.rn = rn
		// p.pt.w = n
		// p.pt.col++
		var rn = this.data.charAt(this.pt.offset);
		var rnB = Bytes.ofString(rn);
        var _rnInt = rn.charCodeAt(0);
		this.pt.rn = _rnInt != null ? _rnInt : 0; // static compile requires !null
		this.pt.w = rnB.length;
		this.pt.col ++;

		// if rn == '\n' {
        if (rn == '\n') {
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
			// ==template== {{ if not .Optimize }}
			if (this.debug) {
				this.out(this.in_("restore"));
			}
			// {{ end }} ==template==
			return;
		}
		this.pt = pt;
		// ==template== {{ if not .Optimize }}
		if (this.debug) {
			this.out(this.in_("restore"));
		}
		// {{ end }} ==template==
	}

	// ==template== {{ if or .GlobalState (not .Optimize) }}
	// clone and return parser current state.
	function cloneState():StoreDict {
		// ==template== {{ if not .Optimize }}
		if (this.debug) {
			// defer this.out(p.in("cloneState"));
		}
		// {{ end }} ==template==

		var state = statePool.get();

		for (i in this.cur.state.keyValueIterator()) {
			state.set(i.key, i.value);
		}

		// TODO: 不理解这个Cloner，没有看到有对应实现
		// for k, v := range p.cur.state {
		// 	if c, ok := v.(Cloner); ok {
		// 		state[k] = c.Clone()
		// 	} else {
		// 		state[k] = v
		// 	}
		// }
		// return state
		return state;
	}

	// restore parser current state to the state storeDict.
	// every restoreState should applied only one time for every cloned state
	function restoreState(state:StoreDict) {
		storeDictDiscard(this.cur.state);
		this.cur.state = state;
		// ==template== {{ if not .Optimize }}
		if (this.debug) {
			this.out(this.in_("restoreState"));
		}
		// {{ end }} ==template==
	}

	// {{ end }} ==template==

	function sliceFrom(start:Savepoint):String {
		// return p.data[start.position.offset:p.pt.position.offset]
		return this.data.substring(start.offset, this.pt.offset);
	}
	

	// ==template== {{ if or .LeftRecursion (not .Optimize) }}
	function getMemoized(node:Any):{result:ResultTuple, ok:Bool} {
		var m = this.memo[this.pt.offset];

		if (m != null) {
			var m2 = m[node];
			if (m2 != null) {
				return {
					result: m2,
					ok: true,
				}
			}
		}

		return {
			// ResultTuple{v:null, b: false, end: null}
			result: null,
			ok: false,
		}
	}

	function setMemoized(pt:Savepoint, node:Any, tuple:ResultTuple) {
		if (this.memo == null) {
			this.memo = new Map<Int, Map<Any, ResultTuple>>();
		}
		var m = this.memo[pt.offset];
		if (m == null) {
			m = new Map<Any, ResultTuple>();
			this.memo[pt.offset] = m;
		}
		m[node] = tuple;
	}
	// {{ end }} ==template==

	public function buildRulesTable(g:Grammar):Void {
		this.rules = new Map<String, Rule>();
		for (r in g.rules) {
			this.rules[r.name] = r;
		}
	}

	// {{ if .Nolint }} nolint: gocyclo {{else}} ==template== {{ end }}
	public function parse(g:Grammar):Any {
		if (g.rules.length == 0) {
			this.addErr(Errors.errNoRule);
			// return nil, p.errs.err()
			throw this.errs.toException();
		}

		// TODO : not super critical but this could be generated
		this.buildRulesTable(g);
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
			// ==template== {{ if not .Optimize }}
			// if p.debug {
			// 	defer p.out(p.in("panic handler"))
			// }
			// {{ end }} ==template==
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
				// p.addErrAt(errors.New("no match found, expected: " + listJoin(expected, ", ", "or")), p.maxFailPos, expected)
				this.addErrAt(new Exception("no match found, expected: " + listJoin(expected, ", ", "or")), this.maxFailPos, expected);
			}
			throw this.errs.toException();
		}

		// return val, p.errs.err()
		if (this.errs.length != 0) {
			throw this.errs.toException();
		}
		return ret.val;
	}

	// ==template== {{ if .LeftRecursion }}
	function parseRuleRecursiveLeader(rule:Rule):RetValOk {
		return {val: null, ok: false};
	}

	function parseRuleRecursiveNoLeader(rule:Rule):RetValOk {
		return this.parseRule(rule);
	}

	// {{ end }} ==template==
	// ==template== {{ if not .Optimize }}
	function parseRuleMemoize(rule:Rule):RetValOk {
		var ret = this.getMemoized(rule); // ResultTuple
		if (ret.ok) {
			this.restore(ret.result.end);
			return {
				val: ret.result.v,
				ok: ret.result.b,
			}
		}

		var startMark = this.pt; // TODO: 这里需要值拷贝吗
		var ret2 = this.parseRule(rule);
		this.setMemoized(startMark, rule, {v: ret2.val, b: ret2.ok, end: this.pt});

		return ret2;
	}

	// {{ end }} ==template==

	function parseRuleWrap(rule:Rule):RetValOk {
		// ==template== {{ if not .Optimize }}
		// if p.debug {
		//     defer p.out(p.in("parseRule " + rule.name))
		// }
		// {{ end }} ==template==
		var ret:RetValOk = {val: null, ok: false};

		// ==template== {{ if not .Optimize }}
		var startMark = this.pt;
		// {{ end }} ==template==

		// ==template== {{ if and .LeftRecursion (not .Optimize) }}
		if (this.memoize || rule.leftRecursive) {
			if (rule.leader) {
				ret = this.parseRuleRecursiveLeader(rule);
			} else if (this.memoize && !rule.leftRecursive) {
				ret = this.parseRuleMemoize(rule);
			} else {
				ret = this.parseRuleRecursiveNoLeader(rule);
			}
		} else {
			ret = this.parseRule(rule);
		}
		// {{ else if not .Optimize }}
		if (this.memoize) {
			ret = this.parseRuleMemoize(rule);
		} else {
			ret = this.parseRule(rule);
		}
		// {{ else if .LeftRecursion }}
		if (rule.leftRecursive) {
			if (rule.leader) {
				ret = this.parseRuleRecursiveLeader(rule);
			} else {
				ret = this.parseRuleRecursiveNoLeader(rule);
			}
		} else {
			ret = this.parseRule(rule);
		}
		// {{ else }}
		ret = this.parseRule(rule);
		// {{ end }} ==template==

		// ==template== {{ if not .Optimize }}
		if (ret.ok && this.debug) {
			this.printIndent("MATCH", this.sliceFrom(startMark));
		}
		// {{ end }} ==template==
		return ret;
	}

	function parseRule(rule:Rule):RetValOk {
		this.rstack.add(rule);
		this.pushV();

		var ret = this.parseExprWrap(rule.expr);
		this.popV();
		this.rstack.pop();
		return ret;
	}

	function parseExprWrap(expr:AllExpr):RetValOk {
		// ==template== {{ if not .Optimize }}
		var pt:Savepoint = {
			line: 0,
			col: 0,
			offset: 0,
			rn: 0,
			w: 0
		};

		// ==template== {{ if .LeftRecursion }}
		var isLeftRecursion = this.rstack.first().leftRecursive;
		if (this.memoize && !isLeftRecursion) {
			// {{ else }}
			// if (this.memoize) { // NOTE: 进行调整
			// {{ end }} ==template==
			var ret = this.getMemoized(expr);
			if (ret.ok) {
				this.restore(ret.result.end);
				return {val: ret.result.v, ok: ret.result.b};
			}
			pt = this.pt;
		}

		// {{ end }} ==template==
		var ret = this.parseExpr(expr);

		// ==template== {{ if not .Optimize }}
		// ==template== {{ if .LeftRecursion }}
		if (this.memoize && !isLeftRecursion) {
			// {{ else }}
			// if (this.memoize) { // NOTE: 进行调整
			// {{ end }} ==template==
			// this.setMemoized(pt, expr, resultTuple{val, ok, p.pt})
			this.setMemoized(pt, expr, {v: ret.val, b: ret.ok, end: this.pt});
		}
		// {{ end }} ==template==
		return ret;
	}

	// {{ if .Nolint }} nolint: gocyclo {{else}} ==template== {{ end }}
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
			case AnyMatcher(expr):
				ret = this.parseAnyMatcher(expr);
			case CharClassMatcher(expr):
				ret = this.parseCharClassMatcher(expr);
			case ChoiceExpr(expr):
				ret = this.parseChoiceExpr(expr);
			case CustomParserCodeExpr(expr):
				ret = this.parseCustomParserCodeExpr(expr);
			case LabeledExpr(expr):
				ret = this.parseLabeledExpr(expr);
			case LitMatcher(expr):
				ret = this.parseLitMatcher(expr);
			case NotCodeExpr(expr):
				ret = this.parseNotCodeExpr(expr);
			case NotExpr(expr):
				ret = this.parseNotExpr(expr);
			case OneOrMoreExpr(expr):
				ret = this.parseOneOrMoreExpr(expr);
			case RecoveryExpr(expr):
				ret = this.parseRecoveryExpr(expr);
			case RuleRefExpr(expr):
				ret = this.parseRuleRefExpr(expr);
			case SeqExpr(expr):
				ret = this.parseSeqExpr(expr);
			// ==template== {{ if or .GlobalState (not .Optimize) }}
			case StateCodeExpr(expr):
				ret = this.parseStateCodeExpr(expr);
			// {{ end }} ==template==
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
		// ==template== {{ if not .Optimize }}
		if (this.debug) {
			// defer this.out(this.in_("parseActionExpr"))
		}
		// {{ end }} ==template==
		var start = savepointClone(this.pt);
		var ret = this.parseExprWrap(act.expr);
		if (ret.ok) {
			// this.cur.pos = start.position;
			this.cur.pos = {line: start.line, col: start.col, offset: start.offset};
			this.cur.text = this.sliceFrom(start);
			// ==template== {{ if or .GlobalState (not .Optimize) }}
			var state = this.cloneState();
			// {{ end }} ==template==
			var RetValErr = act.run(this);
			if (RetValErr.err != null) {
				this.addErrAt(RetValErr.err, {line: start.line, col: start.col, offset: start.offset}, []);
			}
			// ==template== {{ if or .GlobalState (not .Optimize) }}
			this.restoreState(state);
			// {{ end }} ==template==

			ret.val = RetValErr.val;
		}
		// ==template== {{ if not .Optimize }}
		if (ret.ok && this.debug) {
			this.printIndent("MATCH", this.sliceFrom(start));
		}
		// {{ end }} ==template==
		return ret;
	}

	function parseAndCodeExpr(and:AndCodeExpr):RetValOk {
		// ==template== {{ if or .GlobalState (not .Optimize) }}
		var state = this.cloneState();
		// {{ end }} ==template==

		// ok, err := and.run(p);
		var ret = and.run(this);
		if (ret.err != null) {
			this.addErr(ret.err);
		}
		// ==template== {{ if or .GlobalState (not .Optimize) }}
		this.restoreState(state);
		// {{ end }} ==template==

		// ==template== {{ if not .Optimize }}
		if (this.debug) {
			// defer p.out(p.in("parseAndCodeExpr"))
			this.out(this.in_("parseAndCodeExpr"));
		}
		// {{ end }} ==template==

		return {val: null, ok: true};
	}

	function parseAndExpr(and:AndExpr):RetValOk {
		var pt = this.pt;
		// ==template== {{ if or .GlobalState (not .Optimize) }}
		var state = this.cloneState();
		// {{ end }} ==template==
		this.pushV();
		var ret = this.parseExprWrap(and.expr);
		this.popV(); 
		// ==template== {{ if or .GlobalState (not .Optimize) }}
		this.restoreState(state);
		// {{ end }} ==template==
		this.restore(pt);

		// ==template== {{ if not .Optimize }}
		if (this.debug) {
			this.out(this.in_("parseAndExpr"));
		}
		// {{ end }} ==template==
		return {val: null, ok: ret.ok};
	}

	function parseAnyMatcher(act:AnyMatcher):RetValOk {
		// ==template== {{ if not .Optimize }}
		// if this.debug {
		// 	defer this.out(this.in("parseAnyMatcher"))
		// }
		// {{ end }} ==template==
		// if p.pt.rn == utf8.RuneError && p.pt.w == 0 {
		if (this.pt.rn == 0 && this.pt.w == 0) {
			// EOF - see utf8.DecodeRune
			this.failAt(false, this.pt, ".");
			return {val: null, ok: false};
		}
		var start = savepointClone(this.pt);
		this.read();
		this.failAt(true, start, ".");
		return {val: this.sliceFrom(start), ok: true};
	}

	// {{ if .Nolint }} nolint: gocyclo {{else}} ==template== {{ end }}
	function parseCharClassMatcher(chr:CharClassMatcher):RetValOk {
		// ==template== {{ if not .Optimize }}
		// if (this.debug) {
		// 	defer this.out(this.in("parseCharClassMatcher"))
		// }
		// {{ end }} ==template==
		var cur = this.pt.rn;
		var start = savepointClone(this.pt);

		// ==template== {{ if .BasicLatinLookupTable }}
		if (cur < 128) {
			// if (chr.basicLatinChars[cur] != chr.inverted) {
			// 	this.read()
			// 	this.failAt(true, start.position, chr.val)
			// 	return this.sliceFrom(start), true
			// }
			this.failAt(false, start, chr.val);
			return {val: null, ok: false};
		}
		// {{ end }} ==template==

		// can't match EOF
		if (cur == 0 && this.pt.w == 0) { // see utf8.DecodeRune
			this.failAt(false, start, chr.val);
			return {val: null, ok: false};
		}

		if (chr.ignoreCase) {
			cur = String.fromCharCode(cur).toLowerCase().charCodeAt(0);
		}

		// try to match in the list of available chars
		for (rn in chr.chars) {
			if (rn == cur) {
				if (chr.inverted) {
					this.failAt(false, start, chr.val);
					return {val: null, ok: false};
				}
				this.read();
				this.failAt(true, start, chr.val);
				return {val: this.sliceFrom(start), ok: true};
			}
		}

		// try to match in the list of ranges
		for (i in 0...chr.ranges.length) {
			if (i % 2 != 0)
				continue;
			if (cur >= chr.ranges[i] && cur <= chr.ranges[i + 1]) {
				if (chr.inverted) {
					this.failAt(false, start, chr.val);
					return {val: null, ok: false};
				}
				this.read();
				this.failAt(true, start, chr.val);
				return {val: this.sliceFrom(start), ok: true};
			}
		}

		// try to match in the list of Unicode classes
		if (chr.classes != null) {
			for (cl in chr.classes) {
				if (Unicode.Is(cl, cur)) {
					if (chr.inverted) {
						this.failAt(false, start, chr.val);
						return {val: null, ok: false};
					}
					this.read();
					this.failAt(true, start, chr.val);
					return {val: this.sliceFrom(start), ok: true};
				}
			}
		}

		if (chr.inverted) {
			this.read();
			this.failAt(true, start, chr.val);
			return {val: this.sliceFrom(start), ok: true};
		}
		this.failAt(true, start, chr.val);
		return {val: null, ok: false};
	}

	// ==template== {{ if not .Optimize }}
	function incChoiceAltCnt(ch:ChoiceExpr, altI:Int) {
		// choiceIdent := fmt.Sprintf("%s %d:%d", p.rstack.first().name, ch.pos.line, ch.pos.col)
		var choiceIdent:String = '${this.rstack.first().name} ${ch.pos.line}:${ch.pos.col}';

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

	// {{ end }} ==template==

	function parseChoiceExpr(ch:ChoiceExpr):RetValOk {
		// ==template== {{ if not .Optimize }}
		// if p.debug {
		// 	defer p.out(p.in("parseChoiceExpr"))
		// }
		// {{ end }} ==template==

		// for altI, alt := range ch.alternatives {
		for (altI in 0...ch.alternatives.length) {
			var alt = ch.alternatives[altI];

			// dummy assignment to prevent compile error if optimized
			var _ = altI;

			// ==template== {{ if or .GlobalState (not .Optimize) }}
			var state = this.cloneState();
			// {{ end }} ==template==

			this.pushV();
			var ret = this.parseExprWrap(alt);
			this.popV();
			if (ret.ok) {
				// ==template== {{ if not .Optimize }}
				this.incChoiceAltCnt(ch, altI);
				// {{ end }} ==template==
				return ret;
			}
			// ==template== {{ if or .GlobalState (not .Optimize) }}
			this.restoreState(state);
			// {{ end }} ==template==
		}
		// ==template== {{ if not .Optimize }}
		this.incChoiceAltCnt(ch, _choiceNoMatch);
		// {{ end }} ==template==

		return {val: null, ok: false};
	}

	function parseLabeledExpr(lab:LabeledExpr):RetValOk {
		// ==template== {{ if not .Optimize }}
		// if this.debug {
		// 	defer this.out(this.in("parseLabeledExpr"))
		// }
		// {{ end }} ==template==
		this.pushV();
		var ret = this.parseExprWrap(lab.expr);
		this.popV();
		if (ret.ok && lab.label != "") {
			var m = this.vstack.first();
			m[lab.label] = ret.val;
		}
		return ret;
	}

	function parseCustomParserCodeExpr(code:CustomParserCodeExpr):RetValOk {
		// ==template== {{ if not .Optimize }}
		// if p.debug {
		// 	defer p.out(p.in("parseCustomParserCodeExpr"))
		// }
		// {{ end }} ==template==
		// val, ok, err := code.run(p)
		var ret = code.run(this);
		if (ret.err != null) {
			this.addErr(ret.err);
			return {val: null, ok: false};
		}

		return {val: ret.val, ok: ret.ok};
	}

	function parseLitMatcher(lit:LitMatcher):RetValOk {
		// ==template== {{ if not .Optimize }}
		// if p.debug {
		// 	defer p.out(p.in("parseLitMatcher"))
		// }
		// {{ end }} ==template==
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
		return {val: this.sliceFrom(start), ok: true};
	}

	function parseNotCodeExpr(code:NotCodeExpr):RetValOk {
		// ==template== {{ if not .Optimize }}
		// if p.debug {
		// 	defer p.out(p.in("parseNotCodeExpr"))
		// }
		// {{ end }} ==template==

		// ==template== {{ if or .GlobalState (not .Optimize) }}
		var state = this.cloneState();
		// {{ end }} ==template==

		var ret = code.run(this);
		if (ret.err != null) {
			this.addErr(ret.err);
		}
		// ==template== {{ if or .GlobalState (not .Optimize) }}
		this.restoreState(state);
		// {{ end }} ==template==

		return {val: null, ok: !ret.ok};
	}

	function parseNotExpr(not:NotExpr):RetValOk {
		// ==template== {{ if not .Optimize }}
		// if this.debug {
		// 	defer this.out(this.in("parseNotExpr"))
		// }
		// {{ end }} ==template==
		var pt = savepointClone(this.pt);
		// ==template== {{ if or .GlobalState (not .Optimize) }}
		var state = this.cloneState();
		// {{ end }} ==template==
		this.pushV();
		this.maxFailInvertExpected = !this.maxFailInvertExpected;
		var ret = this.parseExprWrap(not.expr);
		this.maxFailInvertExpected = !this.maxFailInvertExpected;
		this.popV();
		// ==template== {{ if or .GlobalState (not .Optimize) }}
		this.restoreState(state);
		// {{ end }} ==template==
		this.restore(pt);

		return {val: null, ok: !ret.ok};
	}

	function parseOneOrMoreExpr(expr:OneOrMoreExpr):RetValOk {
		// ==template== {{ if not .Optimize }}
		// if this.debug {
		// 	defer this.out(this.in("parseOneOrMoreExpr"))
		// }
		// {{ end }} ==template==
		var vals:Array<Any> = [];

		while (true) {
			this.pushV();
			var ret = this.parseExprWrap(expr.expr);
			this.popV();
			if (!ret.ok) {
				if (vals.length == 0) {
					// did not match once, no match
					return {val: null, ok: false};
				}
				return {val: vals, ok: true};
			}
			vals.push(ret.val);
		}
	}

	function parseRecoveryExpr(recover:RecoveryExpr):RetValOk {
		// ==template== {{ if not .Optimize }}
		// if this.debug {
		// 	defer this.out(this.in("parseRecoveryExpr (" + strings.Join(recover.failureLabel, ",") + ")"))
		// }
		// {{ end }} ==template==

		this.pushRecovery(recover.failureLabel, recover.recoverExpr);
		var ret = this.parseExprWrap(recover.expr);
		this.popRecovery();

		return ret;
	}

	function parseRuleRefExpr(ref:RuleRefExpr):RetValOk {
		// ==template== {{ if not .Optimize }}
		// if this.debug {
		// 	defer this.out(this.in("parseRuleRefExpr " + ref.name))
		// }
		// {{ end }} ==template==
		if (ref.name == "") {
			// panic(fmt.Sprintf("%s: invalid rule: missing name", ref.pos))
			throw '${ref.pos}: invalid rule: missing name';
		}

		var rule = this.rules[ref.name];
		if (rule == null) {
			// p.addErr(fmt.Errorf("undefined rule: %s", ref.name))
			this.addErr(new Exception('undefined rule: ${ref.name}'));
			return {val: null, ok: false};
		}
		return this.parseRuleWrap(rule);
	}

	function parseSeqExpr(seq:SeqExpr):RetValOk {
		// ==template== {{ if not .Optimize }}
		// if this.debug {
		//     defer this.out(this.in("parseSeqExpr"))
		// }
		// {{ end }} ==template==
		var vals = [];

		var pt = savepointClone(this.pt);
		// ==template== {{ if or .GlobalState (not .Optimize) }}
		var state = this.cloneState();
		// {{ end }} ==template==
		for (expr in seq.exprs) {
			var ret = this.parseExprWrap(expr);
			if (!ret.ok) {
				// ==template== {{ if or .GlobalState (not .Optimize) }}
				this.restoreState(state);
				// {{ end }} ==template==
				this.restore(pt);
				return {val: null, ok: false};
			}
			vals.push(ret.val);
		}

		return {val: vals, ok: true};
	}

	// ==template== {{ if or .GlobalState (not .Optimize) }}
	function parseStateCodeExpr(code:StateCodeExpr):RetValOk {
		// ==template== {{ if not .Optimize }}
		// if p.debug {
		// 	defer p.out(p.in("parseStateCodeExpr"))
		// }
		// {{ end }} ==template==
		var ret = code.run(this);
		if (ret.err != null) {
			this.addErr(ret.err);
		}
		return {val: null, ok: true};
	}

	// {{ end }} ==template==

	function parseThrowExpr(expr:ThrowExpr):RetValOk {
		// ==template== {{ if not .Optimize }}
		// if this.debug {
		// 	defer this.out(this.in("parseThrowExpr"))
		// }
		// {{ end }} ==template==

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
		// ==template== {{ if not .Optimize }}
		// if this.debug {
		// 	defer this.out(this.in("parseZeroOrMoreExpr"))
		// }
		// {{ end }} ==template==
		var vals:Array<Any> = [];

		while (true) {
			this.pushV();
			var ret = this.parseExprWrap(expr.expr);
			this.popV();
			if (!ret.ok) {
				return {val: vals, ok: true};
			}
			vals.push(ret.val);
		}
	}

	function parseZeroOrOneExpr(expr:ZeroOrOneExpr):RetValOk {
		// ==template== {{ if not .Optimize }}
		// if this.debug {
		// 	defer this.out(this.in("parseZeroOrOneExpr"))
		// }
		// {{ end }} ==template==
		this.pushV();
		var ret = this.parseExprWrap(expr.expr);
		this.popV();
		// whether it matched or not, consider it a match
		return {val: ret.val, ok: true};
	}
}

// ==template== {{ if or .GlobalState (not .Optimize) }}

// Cloner is implemented by any value that has a Clone method, which returns a
// copy of the value. This is mainly used for types which are not passed by
// value (e.g map, slice, chan) or structs that contain such types.
//
// This is used in conjunction with the global state feature to create proper
// copies of the state to allow the parser to properly restore the state in
// the case of backtracking.
class ObjectPool<T> {
    var pool:Array<T> = [];
    var create:Void->T = null;

    public function new(create:Void->T) {
        this.create = create;
    }

    public function get():T {
        if (pool.length > 0) {
            return pool.pop();
        } else {
            // 注: 没找到 T extends Type 的写法，就这样吧
            return this.create();
        }
    }

	public function put(obj:T):Void {
        pool.push(obj);
    }
}

var statePool = new ObjectPool<StoreDict>(() -> {
    return new StoreDict();
});

// {{ end }} ==template==

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

// ==template== {{ if .LeftRecursion }}

func (p *parser) parseRuleRecursiveLeader(rule *rule) (any, bool) {
	result, ok := p.getMemoized(rule)
	if ok {
		p.restore(result.end)
		return result.v, result.b
	}

	// ==template== {{ if not .Optimize }}
	if p.debug {
		defer p.out(p.in("recursive " + rule.name))
	}
	// {{ end }} ==template==

	var (
		depth      = 0
		startMark  = p.pt
		lastResult = resultTuple{nil, false, startMark}
		lastErrors = *p.errs
	)

	for {
		// ==template== {{ if or .GlobalState (not .Optimize) }}
		lastState := p.cloneState()
		// {{ end }} ==template==
		p.setMemoized(startMark, rule, lastResult)
		val, ok := p.parseRule(rule)
		endMark := p.pt
		// ==template== {{ if not .Optimize }}
		if p.debug {
			p.printIndent("RECURSIVE", fmt.Sprintf(
				"Rule %s depth %d: %t -> %s",
				rule.name, depth, ok, string(p.sliceFrom(startMark))))
		}
		// {{ end }} ==template==
		if (!ok) || (endMark.offset <= lastResult.end.offset && depth != 0) {
			// ==template== {{ if or .GlobalState (not .Optimize) }}
			p.restoreState(lastState)
			// {{ end }} ==template==
			*p.errs = lastErrors
			break
		}
		lastResult = resultTuple{val, ok, endMark}
		lastErrors = *p.errs
		p.restore(startMark)
		depth++
	}

	p.restore(lastResult.end)
	p.setMemoized(startMark, rule, lastResult)
	return lastResult.v, lastResult.b
}

// {{ end }} ==template==
