// Code is generated; DO NOT EDIT.
import haxe.Exception;
import haxe.ds.GenericStack;
import haxe.io.BytesBuffer;
import haxe.ds.HashMap;
import haxe.iterators.StringIteratorUnicode;

class ParserCustomData {
	public function new() {}
}

var g:Grammar = {
	rules: [
		{
			name: "start",
			varExists: true,
			expr: ActionExpr({
				run: _call_onstart_1,
				expr: SeqExpr({
					exprs: [
						LabeledExpr({
							label: "first",
							expr: RuleRefExpr({name: "word",}),
							textCapture: false,
						}),
						LabeledExpr({
							label: "rest",
							expr: ZeroOrMoreExpr({
								expr: ActionExpr({
									run: _call_onstart_7,
									expr: SeqExpr({
										exprs: [
											LitMatcher({val: ",", want: "\",\"",}),
											LabeledExpr({
												label: "next",
												expr: RuleRefExpr({name: "word",}),
												textCapture: false,
											}),
										],
									}),
								}),
							}),
							textCapture: false,
						}),
						NotExpr({
							expr: AnyMatcher({}),
						}),
					],
				}),
			}),
		},
		{
			name: "word",
			expr: ActionExpr({
				run: _call_onword_1,
				expr: OneOrMoreExpr({
					expr: CharClassMatcher({
						val: "[\\p{L}]",
						classes: [
							{R16: [
								[65,90,1], [97,122,1], [170,181,11], [186,192,6], [193,214,1], [216,246,1],
								[248,705,1], [710,721,1], [736,740,1], [748,750,2], [880,884,1], [886,887,1],
								[890,893,1], [895,902,7], [904,906,1], [908,910,2], [911,929,1], [931,1013,1],
								[1015,1153,1], [1162,1327,1], [1329,1366,1], [1369,1376,7], [1377,1416,1],
								[1488,1514,1], [1519,1522,1], [1568,1610,1], [1646,1647,1], [1649,1747,1],
								[1749,1765,16], [1766,1774,8], [1775,1786,11], [1787,1788,1], [1791,1808,17],
								[1810,1839,1], [1869,1957,1], [1969,1994,25], [1995,2026,1], [2036,2037,1],
								[2042,2048,6], [2049,2069,1], [2074,2084,10], [2088,2112,24], [2113,2136,1],
								[2144,2154,1], [2160,2183,1], [2185,2190,1], [2208,2249,1], [2308,2361,1],
								[2365,2384,19], [2392,2401,1], [2417,2432,1], [2437,2444,1], [2447,2448,1],
								[2451,2472,1], [2474,2480,1], [2482,2486,4], [2487,2489,1], [2493,2510,17],
								[2524,2525,1], [2527,2529,1], [2544,2545,1], [2556,2565,9], [2566,2570,1],
								[2575,2576,1], [2579,2600,1], [2602,2608,1], [2610,2611,1], [2613,2614,1],
								[2616,2617,1], [2649,2652,1], [2654,2674,20], [2675,2676,1], [2693,2701,1],
								[2703,2705,1], [2707,2728,1], [2730,2736,1], [2738,2739,1], [2741,2745,1],
								[2749,2768,19], [2784,2785,1], [2809,2821,12], [2822,2828,1], [2831,2832,1],
								[2835,2856,1], [2858,2864,1], [2866,2867,1], [2869,2873,1], [2877,2908,31],
								[2909,2911,2], [2912,2913,1], [2929,2947,18], [2949,2954,1], [2958,2960,1],
								[2962,2965,1], [2969,2970,1], [2972,2974,2], [2975,2979,4], [2980,2984,4],
								[2985,2986,1], [2990,3001,1], [3024,3077,53], [3078,3084,1], [3086,3088,1],
								[3090,3112,1], [3114,3129,1], [3133,3160,27], [3161,3162,1], [3165,3168,3],
								[3169,3200,31], [3205,3212,1], [3214,3216,1], [3218,3240,1], [3242,3251,1],
								[3253,3257,1], [3261,3293,32], [3294,3296,2], [3297,3313,16], [3314,3332,18],
								[3333,3340,1], [3342,3344,1], [3346,3386,1], [3389,3406,17], [3412,3414,1],
								[3423,3425,1], [3450,3455,1], [3461,3478,1], [3482,3505,1], [3507,3515,1],
								[3517,3520,3], [3521,3526,1], [3585,3632,1], [3634,3635,1], [3648,3654,1],
								[3713,3714,1], [3716,3718,2], [3719,3722,1], [3724,3747,1], [3749,3751,2],
								[3752,3760,1], [3762,3763,1], [3773,3776,3], [3777,3780,1], [3782,3804,22],
								[3805,3807,1], [3840,3904,64], [3905,3911,1], [3913,3948,1], [3976,3980,1],
								[4096,4138,1], [4159,4176,17], [4177,4181,1], [4186,4189,1], [4193,4197,4],
								[4198,4206,8], [4207,4208,1], [4213,4225,1], [4238,4256,18], [4257,4293,1],
								[4295,4301,6], [4304,4346,1], [4348,4680,1], [4682,4685,1], [4688,4694,1],
								[4696,4698,2], [4699,4701,1], [4704,4744,1], [4746,4749,1], [4752,4784,1],
								[4786,4789,1], [4792,4798,1], [4800,4802,2], [4803,4805,1], [4808,4822,1],
								[4824,4880,1], [4882,4885,1], [4888,4954,1], [4992,5007,1], [5024,5109,1],
								[5112,5117,1], [5121,5740,1], [5743,5759,1], [5761,5786,1], [5792,5866,1],
								[5873,5880,1], [5888,5905,1], [5919,5937,1], [5952,5969,1], [5984,5996,1],
								[5998,6000,1], [6016,6067,1], [6103,6108,5], [6176,6264,1], [6272,6276,1],
								[6279,6312,1], [6314,6320,6], [6321,6389,1], [6400,6430,1], [6480,6509,1],
								[6512,6516,1], [6528,6571,1], [6576,6601,1], [6656,6678,1], [6688,6740,1],
								[6823,6917,94], [6918,6963,1], [6981,6988,1], [7043,7072,1], [7086,7087,1],
								[7098,7141,1], [7168,7203,1], [7245,7247,1], [7258,7293,1], [7296,7304,1],
								[7312,7354,1], [7357,7359,1], [7401,7404,1], [7406,7411,1], [7413,7414,1],
								[7418,7424,6], [7425,7615,1], [7680,7957,1], [7960,7965,1], [7968,8005,1],
								[8008,8013,1], [8016,8023,1], [8025,8031,2], [8032,8061,1], [8064,8116,1],
								[8118,8124,1], [8126,8130,4], [8131,8132,1], [8134,8140,1], [8144,8147,1],
								[8150,8155,1], [8160,8172,1], [8178,8180,1], [8182,8188,1], [8305,8319,14],
								[8336,8348,1], [8450,8455,5], [8458,8467,1], [8469,8473,4], [8474,8477,1],
								[8484,8490,2], [8491,8493,1], [8495,8505,1], [8508,8511,1], [8517,8521,1],
								[8526,8579,53], [8580,11264,2684], [11265,11492,1], [11499,11502,1],
								[11506,11507,1], [11520,11557,1], [11559,11565,6], [11568,11623,1],
								[11631,11648,17], [11649,11670,1], [11680,11686,1], [11688,11694,1],
								[11696,11702,1], [11704,11710,1], [11712,11718,1], [11720,11726,1],
								[11728,11734,1], [11736,11742,1], [11823,12293,470], [12294,12337,43],
								[12338,12341,1], [12347,12348,1], [12353,12438,1], [12445,12447,1],
								[12449,12538,1], [12540,12543,1], [12549,12591,1], [12593,12686,1],
								[12704,12735,1], [12784,12799,1], [13312,19903,1], [19968,42124,1],
								[42192,42237,1], [42240,42508,1], [42512,42527,1], [42538,42539,1],
								[42560,42606,1], [42623,42653,1], [42656,42725,1], [42775,42783,1],
								[42786,42888,1], [42891,42954,1], [42960,42961,1], [42963,42965,2],
								[42966,42969,1], [42994,43009,1], [43011,43013,1], [43015,43018,1],
								[43020,43042,1], [43072,43123,1], [43138,43187,1], [43250,43255,1],
								[43259,43261,2], [43262,43274,12], [43275,43301,1], [43312,43334,1],
								[43360,43388,1], [43396,43442,1], [43471,43488,17], [43489,43492,1],
								[43494,43503,1], [43514,43518,1], [43520,43560,1], [43584,43586,1],
								[43588,43595,1], [43616,43638,1], [43642,43646,4], [43647,43695,1],
								[43697,43701,4], [43702,43705,3], [43706,43709,1], [43712,43714,2],
								[43739,43741,1], [43744,43754,1], [43762,43764,1], [43777,43782,1],
								[43785,43790,1], [43793,43798,1], [43808,43814,1], [43816,43822,1],
								[43824,43866,1], [43868,43881,1], [43888,44002,1], [44032,55203,1],
								[55216,55238,1], [55243,55291,1], [63744,64109,1], [64112,64217,1],
								[64256,64262,1], [64275,64279,1], [64285,64287,2], [64288,64296,1],
								[64298,64310,1], [64312,64316,1], [64318,64320,2], [64321,64323,2],
								[64324,64326,2], [64327,64433,1], [64467,64829,1], [64848,64911,1],
								[64914,64967,1], [65008,65019,1], [65136,65140,1], [65142,65276,1],
								[65313,65338,1], [65345,65370,1], [65382,65470,1], [65474,65479,1],
								[65482,65487,1], [65490,65495,1], [65498,65500,1],
							], R32: [
								[65536,65547,1], [65549,65574,1], [65576,65594,1], [65596,65597,1],
								[65599,65613,1], [65616,65629,1], [65664,65786,1], [66176,66204,1],
								[66208,66256,1], [66304,66335,1], [66349,66368,1], [66370,66377,1],
								[66384,66421,1], [66432,66461,1], [66464,66499,1], [66504,66511,1],
								[66560,66717,1], [66736,66771,1], [66776,66811,1], [66816,66855,1],
								[66864,66915,1], [66928,66938,1], [66940,66954,1], [66956,66962,1],
								[66964,66965,1], [66967,66977,1], [66979,66993,1], [66995,67001,1],
								[67003,67004,1], [67072,67382,1], [67392,67413,1], [67424,67431,1],
								[67456,67461,1], [67463,67504,1], [67506,67514,1], [67584,67589,1],
								[67592,67594,2], [67595,67637,1], [67639,67640,1], [67644,67647,3],
								[67648,67669,1], [67680,67702,1], [67712,67742,1], [67808,67826,1],
								[67828,67829,1], [67840,67861,1], [67872,67897,1], [67968,68023,1],
								[68030,68031,1], [68096,68112,16], [68113,68115,1], [68117,68119,1],
								[68121,68149,1], [68192,68220,1], [68224,68252,1], [68288,68295,1],
								[68297,68324,1], [68352,68405,1], [68416,68437,1], [68448,68466,1],
								[68480,68497,1], [68608,68680,1], [68736,68786,1], [68800,68850,1],
								[68864,68899,1], [69248,69289,1], [69296,69297,1], [69376,69404,1],
								[69415,69424,9], [69425,69445,1], [69488,69505,1], [69552,69572,1],
								[69600,69622,1], [69635,69687,1], [69745,69746,1], [69749,69763,14],
								[69764,69807,1], [69840,69864,1], [69891,69926,1], [69956,69959,3],
								[69968,70002,1], [70006,70019,13], [70020,70066,1], [70081,70084,1],
								[70106,70108,2], [70144,70161,1], [70163,70187,1], [70207,70208,1],
								[70272,70278,1], [70280,70282,2], [70283,70285,1], [70287,70301,1],
								[70303,70312,1], [70320,70366,1], [70405,70412,1], [70415,70416,1],
								[70419,70440,1], [70442,70448,1], [70450,70451,1], [70453,70457,1],
								[70461,70480,19], [70493,70497,1], [70656,70708,1], [70727,70730,1],
								[70751,70753,1], [70784,70831,1], [70852,70853,1], [70855,71040,185],
								[71041,71086,1], [71128,71131,1], [71168,71215,1], [71236,71296,60],
								[71297,71338,1], [71352,71424,72], [71425,71450,1], [71488,71494,1],
								[71680,71723,1], [71840,71903,1], [71935,71942,1], [71945,71948,3],
								[71949,71955,1], [71957,71958,1], [71960,71983,1], [71999,72001,2],
								[72096,72103,1], [72106,72144,1], [72161,72163,2], [72192,72203,11],
								[72204,72242,1], [72250,72272,22], [72284,72329,1], [72349,72368,19],
								[72369,72440,1], [72704,72712,1], [72714,72750,1], [72768,72818,50],
								[72819,72847,1], [72960,72966,1], [72968,72969,1], [72971,73008,1],
								[73030,73056,26], [73057,73061,1], [73063,73064,1], [73066,73097,1],
								[73112,73440,328], [73441,73458,1], [73474,73476,2], [73477,73488,1],
								[73490,73523,1], [73648,73728,80], [73729,74649,1], [74880,75075,1],
								[77712,77808,1], [77824,78895,1], [78913,78918,1], [82944,83526,1],
								[92160,92728,1], [92736,92766,1], [92784,92862,1], [92880,92909,1],
								[92928,92975,1], [92992,92995,1], [93027,93047,1], [93053,93071,1],
								[93760,93823,1], [93952,94026,1], [94032,94099,67], [94100,94111,1],
								[94176,94177,1], [94179,94208,29], [94209,100343,1], [100352,101589,1],
								[101632,101640,1], [110576,110579,1], [110581,110587,1], [110589,110590,1],
								[110592,110882,1], [110898,110928,30], [110929,110930,1], [110933,110948,15],
								[110949,110951,1], [110960,111355,1], [113664,113770,1], [113776,113788,1],
								[113792,113800,1], [113808,113817,1], [119808,119892,1], [119894,119964,1],
								[119966,119967,1], [119970,119973,3], [119974,119977,3], [119978,119980,1],
								[119982,119993,1], [119995,119997,2], [119998,120003,1], [120005,120069,1],
								[120071,120074,1], [120077,120084,1], [120086,120092,1], [120094,120121,1],
								[120123,120126,1], [120128,120132,1], [120134,120138,4], [120139,120144,1],
								[120146,120485,1], [120488,120512,1], [120514,120538,1], [120540,120570,1],
								[120572,120596,1], [120598,120628,1], [120630,120654,1], [120656,120686,1],
								[120688,120712,1], [120714,120744,1], [120746,120770,1], [120772,120779,1],
								[122624,122654,1], [122661,122666,1], [122928,122989,1], [123136,123180,1],
								[123191,123197,1], [123214,123536,322], [123537,123565,1], [123584,123627,1],
								[124112,124139,1], [124896,124902,1], [124904,124907,1], [124909,124910,1],
								[124912,124926,1], [124928,125124,1], [125184,125251,1], [125259,126464,1205],
								[126465,126467,1], [126469,126495,1], [126497,126498,1], [126500,126503,3],
								[126505,126514,1], [126516,126519,1], [126521,126523,2], [126530,126535,5],
								[126537,126541,2], [126542,126543,1], [126545,126546,1], [126548,126551,3],
								[126553,126561,2], [126562,126564,2], [126567,126570,1], [126572,126578,1],
								[126580,126583,1], [126585,126588,1], [126590,126592,2], [126593,126601,1],
								[126603,126619,1], [126625,126627,1], [126629,126633,1], [126635,126651,1],
								[131072,173791,1], [173824,177977,1], [177984,178205,1], [178208,183969,1],
								[183984,191456,1], [194560,195101,1], [196608,201546,1], [201552,205743,1],
							], LatinOffset: 6},
						],
					}),
				}),
			}),
		},
	],
}

private function _call_onstart_7(p:Parser):Any {
	var stack = p.vstack.first();
	var c = p.cur;
	return (function(c:Current, next:Any):Any {
		return next;
		return null;
	})(p.cur, stack["next"]);
}

private function _call_onstart_1(p:Parser):Any {
	var stack = p.vstack.first();
	var c = p.cur;
	return (function(c:Current, first, rest:Any):Any {
		var values:Array<String> = [Std.string(first)];
	var tail:Array<Dynamic> = cast rest;
	for (item in tail) {
		values.push(Std.string(item));
	}
	return values;
		return null;
	})(p.cur, stack["first"], stack["rest"]);
}

private function _call_onword_1(p:Parser):Any {
	var c = p.cur;
	return (function(c:Current):Any {
		return c.text;
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
		this.entrypoint = "start";

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
