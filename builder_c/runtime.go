package builder_c

const runtimePreamble = `
#include <limits.h>
#include <stddef.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

typedef enum PigeonValueType {
    PIGEON_VALUE_NULL = 0,
    PIGEON_VALUE_STRING,
    PIGEON_VALUE_INTEGER,
    PIGEON_VALUE_NUMBER,
    PIGEON_VALUE_BOOLEAN,
    PIGEON_VALUE_ARRAY,
    PIGEON_VALUE_POINTER
} PigeonValueType;

typedef struct PigeonValue PigeonValue;

typedef struct PigeonString {
    char *data;
    size_t length;
} PigeonString;

typedef struct PigeonArray {
    PigeonValue *items;
    size_t length;
    size_t capacity;
} PigeonArray;

struct PigeonValue {
    PigeonValueType type;
    union {
        PigeonString string;
        long long integer;
        double number;
        int boolean;
        PigeonArray array;
        void *pointer;
    } as;
};

typedef struct PigeonPosition {
    size_t line;
    size_t col;
    size_t offset;
} PigeonPosition;

typedef struct PigeonCurrent {
    PigeonPosition pos;
    char *text;
    size_t text_length;
    void *data;
} PigeonCurrent;

typedef struct PigeonParser PigeonParser;
typedef PigeonValue (*PigeonActionFn)(PigeonParser *parser);
typedef int (*PigeonPredicateFn)(PigeonParser *parser);

typedef struct PigeonUnicodeRange {
    uint32_t lo;
    uint32_t hi;
    uint32_t stride;
} PigeonUnicodeRange;

typedef enum PigeonExprKind {
    PIGEON_EXPR_ACTION = 0,
    PIGEON_EXPR_AND_CODE,
    PIGEON_EXPR_AND,
    PIGEON_EXPR_AND_LOGICAL,
    PIGEON_EXPR_ANY,
    PIGEON_EXPR_CHAR_CLASS,
    PIGEON_EXPR_CHOICE,
    PIGEON_EXPR_CODE,
    PIGEON_EXPR_LABELED,
    PIGEON_EXPR_LITERAL,
    PIGEON_EXPR_NOT_CODE,
    PIGEON_EXPR_NOT,
    PIGEON_EXPR_NOT_LOGICAL,
    PIGEON_EXPR_ONE_OR_MORE,
    PIGEON_EXPR_RECOVERY,
    PIGEON_EXPR_RULE_REF,
    PIGEON_EXPR_RULE_INDEX,
    PIGEON_EXPR_SEQUENCE,
    PIGEON_EXPR_THROW,
    PIGEON_EXPR_ZERO_OR_MORE,
    PIGEON_EXPR_ZERO_OR_ONE
} PigeonExprKind;

typedef struct PigeonExpr PigeonExpr;

struct PigeonExpr {
    PigeonExprKind kind;
    const PigeonExpr *inner;
    const PigeonExpr *const *items;
    size_t item_count;
    PigeonActionFn action;
    PigeonPredicateFn predicate;
    const PigeonExpr *recover_expr;
    const char *const *failure_labels;
    size_t failure_label_count;
    const char *label;
    int text_capture;
    const char *name;
    size_t index;
    const uint32_t *runes;
    size_t rune_count;
    const char *value;
    const char *want;
    const uint32_t *chars;
    size_t chars_count;
    const uint32_t *ranges;
    size_t ranges_count;
    const PigeonUnicodeRange *classes;
    size_t class_count;
    int ignore_case;
    int inverted;
    int not_skip;
    PigeonPosition pos;
    int reserved;
};

typedef struct PigeonRule {
    const char *name;
    const char *display_name;
    int var_exists;
    const PigeonExpr *expr;
    PigeonPosition pos;
    int reserved;
} PigeonRule;

typedef struct PigeonGrammar {
    const PigeonRule *rules;
    size_t rule_count;
} PigeonGrammar;

typedef struct PigeonParserOptions {
    const char *filename;
    const char *entrypoint;
    void *custom_data;
    size_t max_expressions;
} PigeonParserOptions;

typedef struct PigeonResult {
    int ok;
    PigeonValue value;
    char *error;
    PigeonPosition position;
} PigeonResult;

typedef struct PigeonSavepoint {
    size_t line;
    size_t col;
    size_t offset;
    size_t index;
} PigeonSavepoint;

typedef struct PigeonMatch {
    PigeonValue value;
    int ok;
} PigeonMatch;

typedef struct PigeonBinding {
    const char *name;
    PigeonValue value;
} PigeonBinding;

typedef struct PigeonScope {
    PigeonBinding *bindings;
    size_t length;
    size_t capacity;
} PigeonScope;

typedef struct PigeonRecovery {
    const PigeonExpr *expr;
} PigeonRecovery;

struct PigeonParser {
    const char *filename;
    const unsigned char *input;
    size_t input_length;
    const char *entrypoint;
    size_t max_expressions;
    size_t expression_count;
    PigeonSavepoint pt;
    PigeonCurrent cur;
    const PigeonGrammar *grammar;
    PigeonScope *vstack;
    size_t vstack_length;
    size_t vstack_capacity;
    PigeonRecovery *recovery_stack;
    size_t recovery_length;
    size_t recovery_capacity;
    size_t skip_depth;
    int invert_failures;
    PigeonSavepoint max_fail_pos;
    char **max_fail_expected;
    size_t expected_length;
    size_t expected_capacity;
    char *error;
    PigeonPosition error_position;
};

typedef struct PigeonCaseMap {
    uint32_t from;
    uint32_t to;
} PigeonCaseMap;

PigeonValue pigeon_value_null(void);
PigeonValue pigeon_value_string(const char *value);
PigeonValue pigeon_value_string_n(const char *value, size_t length);
PigeonValue pigeon_value_integer(long long value);
PigeonValue pigeon_value_number(double value);
PigeonValue pigeon_value_boolean(int value);
PigeonValue pigeon_value_array(void);
PigeonValue pigeon_value_pointer(void *value);
PigeonValue pigeon_value_clone(const PigeonValue *value);
int pigeon_value_array_push(PigeonValue *array, PigeonValue value);
void pigeon_value_free(PigeonValue *value);
const PigeonValue *pigeon_parser_label(PigeonParser *parser, const char *name);
void pigeon_parser_add_error(PigeonParser *parser, const char *message);
PigeonResult pigeon_parse(const char *input);
PigeonResult pigeon_parse_n(const char *input, size_t length);
PigeonResult pigeon_parse_with_options(const char *input, const PigeonParserOptions *options);
PigeonResult pigeon_parse_with_options_n(const char *input, size_t length, const PigeonParserOptions *options);
void pigeon_result_free(PigeonResult *result);
`

const runtimeCore = `
static char *_pigeon_copy_string(const char *value, size_t length) {
    char *copy;
    if (length == SIZE_MAX) {
        return NULL;
    }
    copy = (char *)malloc(length + 1);
    if (copy == NULL) {
        return NULL;
    }
    if (length > 0 && value != NULL) {
        memcpy(copy, value, length);
    }
    copy[length] = '\0';
    return copy;
}

PigeonValue pigeon_value_null(void) {
    PigeonValue value;
    memset(&value, 0, sizeof(value));
    value.type = PIGEON_VALUE_NULL;
    return value;
}

PigeonValue pigeon_value_string_n(const char *text, size_t length) {
    PigeonValue value = pigeon_value_null();
    if (text == NULL) {
        length = 0;
    }
    char *copy = _pigeon_copy_string(text, length);
    if (copy == NULL) {
        return value;
    }
    value.type = PIGEON_VALUE_STRING;
    value.as.string.data = copy;
    value.as.string.length = length;
    return value;
}

PigeonValue pigeon_value_string(const char *text) {
    return pigeon_value_string_n(text, text == NULL ? 0 : strlen(text));
}

PigeonValue pigeon_value_integer(long long integer) {
    PigeonValue value = pigeon_value_null();
    value.type = PIGEON_VALUE_INTEGER;
    value.as.integer = integer;
    return value;
}

PigeonValue pigeon_value_number(double number) {
    PigeonValue value = pigeon_value_null();
    value.type = PIGEON_VALUE_NUMBER;
    value.as.number = number;
    return value;
}

PigeonValue pigeon_value_boolean(int boolean) {
    PigeonValue value = pigeon_value_null();
    value.type = PIGEON_VALUE_BOOLEAN;
    value.as.boolean = boolean != 0;
    return value;
}

PigeonValue pigeon_value_array(void) {
    PigeonValue value = pigeon_value_null();
    value.type = PIGEON_VALUE_ARRAY;
    return value;
}

PigeonValue pigeon_value_pointer(void *pointer) {
    PigeonValue value = pigeon_value_null();
    value.type = PIGEON_VALUE_POINTER;
    value.as.pointer = pointer;
    return value;
}

int pigeon_value_array_push(PigeonValue *array, PigeonValue value) {
    PigeonValue *items;
    size_t capacity;
    if (array == NULL || array->type != PIGEON_VALUE_ARRAY) {
        return 0;
    }
    if (array->as.array.length == array->as.array.capacity) {
        capacity = array->as.array.capacity == 0 ? 4 : array->as.array.capacity * 2;
        if (capacity < array->as.array.capacity || capacity > SIZE_MAX / sizeof(PigeonValue)) {
            return 0;
        }
        items = (PigeonValue *)realloc(array->as.array.items, capacity * sizeof(PigeonValue));
        if (items == NULL) {
            return 0;
        }
        array->as.array.items = items;
        array->as.array.capacity = capacity;
    }
    array->as.array.items[array->as.array.length++] = value;
    return 1;
}

PigeonValue pigeon_value_clone(const PigeonValue *source) {
    PigeonValue value;
    size_t index;
    if (source == NULL) {
        return pigeon_value_null();
    }
    switch (source->type) {
    case PIGEON_VALUE_STRING:
        return pigeon_value_string_n(source->as.string.data, source->as.string.length);
    case PIGEON_VALUE_INTEGER:
        return pigeon_value_integer(source->as.integer);
    case PIGEON_VALUE_NUMBER:
        return pigeon_value_number(source->as.number);
    case PIGEON_VALUE_BOOLEAN:
        return pigeon_value_boolean(source->as.boolean);
    case PIGEON_VALUE_POINTER:
        return pigeon_value_pointer(source->as.pointer);
    case PIGEON_VALUE_ARRAY:
        value = pigeon_value_array();
        for (index = 0; index < source->as.array.length; index++) {
            PigeonValue item = pigeon_value_clone(&source->as.array.items[index]);
            if (!pigeon_value_array_push(&value, item)) {
                pigeon_value_free(&item);
                pigeon_value_free(&value);
                return pigeon_value_null();
            }
        }
        return value;
    case PIGEON_VALUE_NULL:
    default:
        return pigeon_value_null();
    }
}

void pigeon_value_free(PigeonValue *value) {
    size_t index;
    if (value == NULL) {
        return;
    }
    if (value->type == PIGEON_VALUE_STRING) {
        free(value->as.string.data);
    } else if (value->type == PIGEON_VALUE_ARRAY) {
        for (index = 0; index < value->as.array.length; index++) {
            pigeon_value_free(&value->as.array.items[index]);
        }
        free(value->as.array.items);
    }
    *value = pigeon_value_null();
}

static PigeonPosition _pigeon_position(PigeonSavepoint point) {
    PigeonPosition value;
    value.line = point.line;
    value.col = point.col;
    value.offset = point.offset;
    return value;
}

void pigeon_parser_add_error(PigeonParser *parser, const char *message) {
    if (parser == NULL || parser->error != NULL) {
        return;
    }
    parser->error = _pigeon_copy_string(message == NULL ? "parser error" : message,
                                        message == NULL ? 12 : strlen(message));
    parser->error_position = _pigeon_position(parser->pt);
}

static void _pigeon_out_of_memory(PigeonParser *parser) {
    pigeon_parser_add_error(parser, "out of memory");
}

static int _pigeon_scope_push(PigeonParser *parser) {
    PigeonScope *scopes;
    size_t capacity;
    if (parser->vstack_length == parser->vstack_capacity) {
        capacity = parser->vstack_capacity == 0 ? 4 : parser->vstack_capacity * 2;
        scopes = (PigeonScope *)realloc(parser->vstack, capacity * sizeof(PigeonScope));
        if (scopes == NULL) {
            _pigeon_out_of_memory(parser);
            return 0;
        }
        parser->vstack = scopes;
        parser->vstack_capacity = capacity;
    }
    memset(&parser->vstack[parser->vstack_length], 0, sizeof(PigeonScope));
    parser->vstack_length++;
    return 1;
}

static void _pigeon_scope_pop(PigeonParser *parser) {
    PigeonScope *scope;
    size_t index;
    if (parser->vstack_length == 0) {
        return;
    }
    scope = &parser->vstack[parser->vstack_length - 1];
    for (index = 0; index < scope->length; index++) {
        pigeon_value_free(&scope->bindings[index].value);
    }
    free(scope->bindings);
    memset(scope, 0, sizeof(*scope));
    parser->vstack_length--;
}

static int _pigeon_scope_set(PigeonParser *parser, const char *name, PigeonValue value) {
    PigeonScope *scope;
    PigeonBinding *bindings;
    size_t index;
    size_t capacity;
    if (parser->vstack_length == 0) {
        pigeon_value_free(&value);
        pigeon_parser_add_error(parser, "label used outside a value scope");
        return 0;
    }
    scope = &parser->vstack[parser->vstack_length - 1];
    for (index = 0; index < scope->length; index++) {
        if (strcmp(scope->bindings[index].name, name) == 0) {
            pigeon_value_free(&scope->bindings[index].value);
            scope->bindings[index].value = value;
            return 1;
        }
    }
    if (scope->length == scope->capacity) {
        capacity = scope->capacity == 0 ? 4 : scope->capacity * 2;
        bindings = (PigeonBinding *)realloc(scope->bindings, capacity * sizeof(PigeonBinding));
        if (bindings == NULL) {
            pigeon_value_free(&value);
            _pigeon_out_of_memory(parser);
            return 0;
        }
        scope->bindings = bindings;
        scope->capacity = capacity;
    }
    scope->bindings[scope->length].name = name;
    scope->bindings[scope->length].value = value;
    scope->length++;
    return 1;
}

const PigeonValue *pigeon_parser_label(PigeonParser *parser, const char *name) {
    PigeonScope *scope;
    size_t index;
    if (parser == NULL || name == NULL || parser->vstack_length == 0) {
        return NULL;
    }
    scope = &parser->vstack[parser->vstack_length - 1];
    for (index = scope->length; index > 0; index--) {
        if (strcmp(scope->bindings[index - 1].name, name) == 0) {
            return &scope->bindings[index - 1].value;
        }
    }
    return NULL;
}

static void _pigeon_clear_expected(PigeonParser *parser) {
    size_t index;
    for (index = 0; index < parser->expected_length; index++) {
        free(parser->max_fail_expected[index]);
    }
    parser->expected_length = 0;
}

static void _pigeon_fail_at(PigeonParser *parser, int matched, PigeonSavepoint point, const char *want) {
    char **items;
    char *expected;
    size_t length;
    size_t capacity;
    if ((matched != 0) != (parser->invert_failures != 0) || point.offset < parser->max_fail_pos.offset) {
        return;
    }
    if (point.offset > parser->max_fail_pos.offset) {
        parser->max_fail_pos = point;
        _pigeon_clear_expected(parser);
    }
    want = want == NULL ? "" : want;
    length = strlen(want) + (parser->invert_failures ? 1 : 0);
    expected = (char *)malloc(length + 1);
    if (expected == NULL) {
        _pigeon_out_of_memory(parser);
        return;
    }
    if (parser->invert_failures) {
        expected[0] = '!';
        memcpy(expected + 1, want, length);
    } else {
        memcpy(expected, want, length + 1);
    }
    if (parser->expected_length == parser->expected_capacity) {
        capacity = parser->expected_capacity == 0 ? 8 : parser->expected_capacity * 2;
        items = (char **)realloc(parser->max_fail_expected, capacity * sizeof(char *));
        if (items == NULL) {
            free(expected);
            _pigeon_out_of_memory(parser);
            return;
        }
        parser->max_fail_expected = items;
        parser->expected_capacity = capacity;
    }
    parser->max_fail_expected[parser->expected_length++] = expected;
}

static uint32_t _pigeon_decode(const unsigned char *input, size_t length, size_t index, size_t *width) {
    unsigned char a;
    unsigned char b;
    unsigned char c;
    unsigned char d;
    uint32_t value;
    if (index >= length) {
        *width = 0;
        return 0;
    }
    a = input[index];
    if (a < 0x80) {
        *width = 1;
        return a;
    }
    if (a >= 0xC2 && a <= 0xDF && index + 1 < length) {
        b = input[index + 1];
        if ((b & 0xC0) == 0x80) {
            *width = 2;
            return ((uint32_t)(a & 0x1F) << 6) | (uint32_t)(b & 0x3F);
        }
    } else if (a >= 0xE0 && a <= 0xEF && index + 2 < length) {
        b = input[index + 1];
        c = input[index + 2];
        if ((b & 0xC0) == 0x80 && (c & 0xC0) == 0x80 &&
            !(a == 0xE0 && b < 0xA0) && !(a == 0xED && b >= 0xA0)) {
            value = ((uint32_t)(a & 0x0F) << 12) | ((uint32_t)(b & 0x3F) << 6) | (uint32_t)(c & 0x3F);
            *width = 3;
            return value;
        }
    } else if (a >= 0xF0 && a <= 0xF4 && index + 3 < length) {
        b = input[index + 1];
        c = input[index + 2];
        d = input[index + 3];
        if ((b & 0xC0) == 0x80 && (c & 0xC0) == 0x80 && (d & 0xC0) == 0x80 &&
            !(a == 0xF0 && b < 0x90) && !(a == 0xF4 && b >= 0x90)) {
            value = ((uint32_t)(a & 0x07) << 18) | ((uint32_t)(b & 0x3F) << 12) |
                    ((uint32_t)(c & 0x3F) << 6) | (uint32_t)(d & 0x3F);
            *width = 4;
            return value;
        }
    }
    *width = 1;
    return 0xFFFD;
}

static uint32_t _pigeon_current(PigeonParser *parser, size_t *width) {
    return _pigeon_decode(parser->input, parser->input_length, parser->pt.index, width);
}

static void _pigeon_advance(PigeonParser *parser) {
    size_t width;
    uint32_t current = _pigeon_current(parser, &width);
    if (width == 0) {
        return;
    }
    parser->pt.index += width;
    parser->pt.offset += width;
    if (current == 10) {
        parser->pt.line++;
        parser->pt.col = 1;
    } else {
        parser->pt.col++;
    }
}

static uint32_t _pigeon_lower(uint32_t value) {
    size_t low = 0;
    size_t high = _pigeon_case_map_count;
    while (low < high) {
        size_t middle = low + (high - low) / 2;
        if (_pigeon_case_map[middle].from < value) {
            low = middle + 1;
        } else {
            high = middle;
        }
    }
    if (low < _pigeon_case_map_count && _pigeon_case_map[low].from == value) {
        return _pigeon_case_map[low].to;
    }
    return value;
}

static char *_pigeon_slice(PigeonParser *parser, PigeonSavepoint start) {
    size_t length = parser->pt.index - start.index;
    return _pigeon_copy_string((const char *)parser->input + start.index, length);
}

static int _pigeon_set_current(PigeonParser *parser, PigeonSavepoint start) {
    char *text = _pigeon_slice(parser, start);
    if (text == NULL) {
        _pigeon_out_of_memory(parser);
        return 0;
    }
    free(parser->cur.text);
    parser->cur.text = text;
    parser->cur.text_length = parser->pt.index - start.index;
    parser->cur.pos = _pigeon_position(start);
    return 1;
}

static PigeonMatch _pigeon_match(PigeonValue value, int ok) {
    PigeonMatch match;
    match.value = value;
    match.ok = ok;
    return match;
}

static PigeonMatch _pigeon_null_match(int ok) {
    return _pigeon_match(pigeon_value_null(), ok);
}

static int _pigeon_skipping(PigeonParser *parser) {
    return parser->skip_depth > 0;
}

static PigeonMatch _pigeon_parse_expr(PigeonParser *parser, const PigeonExpr *expr);

static const PigeonRule *_pigeon_find_rule(PigeonParser *parser, const char *name) {
    size_t index;
    for (index = 0; index < parser->grammar->rule_count; index++) {
        if (strcmp(parser->grammar->rules[index].name, name) == 0) {
            return &parser->grammar->rules[index];
        }
    }
    return NULL;
}

static PigeonMatch _pigeon_parse_rule(PigeonParser *parser, const PigeonRule *rule) {
    PigeonMatch result;
    int scoped = rule->var_exists && !_pigeon_skipping(parser);
    if (scoped && !_pigeon_scope_push(parser)) {
        return _pigeon_null_match(0);
    }
    result = _pigeon_parse_expr(parser, rule->expr);
    if (scoped) {
        _pigeon_scope_pop(parser);
    }
    return result;
}

static PigeonMatch _pigeon_parse_action(PigeonParser *parser, const PigeonExpr *expr) {
    PigeonSavepoint start;
    PigeonMatch result;
    if (_pigeon_skipping(parser)) {
        result = _pigeon_parse_expr(parser, expr->inner);
        pigeon_value_free(&result.value);
        return _pigeon_null_match(result.ok);
    }
    start = parser->pt;
    result = _pigeon_parse_expr(parser, expr->inner);
    if (result.ok) {
        pigeon_value_free(&result.value);
        if (!_pigeon_set_current(parser, start)) {
            return _pigeon_null_match(0);
        }
        result.value = expr->action(parser);
    }
    return result;
}

static PigeonMatch _pigeon_parse_and(PigeonParser *parser, const PigeonExpr *expr, int logical) {
    PigeonSavepoint start = parser->pt;
    PigeonMatch result;
    int consumed;
    parser->skip_depth++;
    result = _pigeon_parse_expr(parser, expr->inner);
    parser->skip_depth--;
    consumed = parser->pt.offset != start.offset;
    parser->pt = start;
    pigeon_value_free(&result.value);
    return _pigeon_null_match(logical ? result.ok && consumed : result.ok);
}

static PigeonMatch _pigeon_parse_not(PigeonParser *parser, const PigeonExpr *expr, int logical) {
    PigeonSavepoint start = parser->pt;
    PigeonMatch result;
    int consumed;
    parser->invert_failures = !parser->invert_failures;
    parser->skip_depth++;
    result = _pigeon_parse_expr(parser, expr->inner);
    parser->skip_depth--;
    parser->invert_failures = !parser->invert_failures;
    consumed = parser->pt.offset != start.offset;
    parser->pt = start;
    pigeon_value_free(&result.value);
    return _pigeon_null_match(logical ? !result.ok && consumed : !result.ok);
}

static PigeonMatch _pigeon_parse_any(PigeonParser *parser) {
    size_t width;
    (void)_pigeon_current(parser, &width);
    if (width == 0) {
        _pigeon_fail_at(parser, 0, parser->pt, ".");
        return _pigeon_null_match(0);
    }
    _pigeon_fail_at(parser, 1, parser->pt, ".");
    _pigeon_advance(parser);
    return _pigeon_null_match(1);
}

static PigeonMatch _pigeon_parse_char_class(PigeonParser *parser, const PigeonExpr *expr) {
    size_t width;
    size_t index;
    uint32_t current = _pigeon_current(parser, &width);
    int found = 0;
    int matched;
    if (width == 0) {
        _pigeon_fail_at(parser, 0, parser->pt, expr->value);
        return _pigeon_null_match(0);
    }
    if (expr->ignore_case) {
        current = _pigeon_lower(current);
    }
    for (index = 0; index < expr->chars_count; index++) {
        if (expr->chars[index] == current) {
            found = 1;
            break;
        }
    }
    for (index = 0; !found && index + 1 < expr->ranges_count; index += 2) {
        if (current >= expr->ranges[index] && current <= expr->ranges[index + 1]) {
            found = 1;
        }
    }
    for (index = 0; !found && index < expr->class_count; index++) {
        const PigeonUnicodeRange *range = &expr->classes[index];
        if (current >= range->lo && current <= range->hi &&
            (current - range->lo) % range->stride == 0) {
            found = 1;
        }
    }
    matched = expr->inverted ? !found : found;
    _pigeon_fail_at(parser, matched, parser->pt, expr->value);
    if (!matched) {
        return _pigeon_null_match(0);
    }
    _pigeon_advance(parser);
    return _pigeon_null_match(1);
}

static PigeonMatch _pigeon_parse_choice(PigeonParser *parser, const PigeonExpr *expr) {
    size_t index;
    for (index = 0; index < expr->item_count; index++) {
        PigeonMatch result = _pigeon_parse_expr(parser, expr->items[index]);
        if (result.ok) {
            return result;
        }
        pigeon_value_free(&result.value);
    }
    return _pigeon_null_match(0);
}

static PigeonMatch _pigeon_parse_code(PigeonParser *parser, const PigeonExpr *expr) {
    if (!expr->not_skip && _pigeon_skipping(parser)) {
        return _pigeon_null_match(1);
    }
    return _pigeon_match(expr->action(parser), 1);
}

static PigeonMatch _pigeon_parse_labeled(PigeonParser *parser, const PigeonExpr *expr) {
    PigeonSavepoint start = parser->pt;
    PigeonMatch result = _pigeon_parse_expr(parser, expr->inner);
    if (result.ok && expr->label != NULL && !_pigeon_skipping(parser)) {
        PigeonValue value;
        if (expr->text_capture) {
            value = pigeon_value_string_n((const char *)parser->input + start.index,
                                          parser->pt.index - start.index);
        } else {
            value = pigeon_value_clone(&result.value);
        }
        if (!_pigeon_scope_set(parser, expr->label, value)) {
            pigeon_value_free(&result.value);
            return _pigeon_null_match(0);
        }
    }
    return result;
}

static PigeonMatch _pigeon_parse_literal(PigeonParser *parser, const PigeonExpr *expr) {
    PigeonSavepoint start = parser->pt;
    size_t index;
    for (index = 0; index < expr->rune_count; index++) {
        size_t width;
        uint32_t current = _pigeon_current(parser, &width);
        if (expr->ignore_case) {
            current = _pigeon_lower(current);
        }
        if (width == 0 || current != expr->runes[index]) {
            _pigeon_fail_at(parser, 0, start, expr->want);
            parser->pt = start;
            return _pigeon_null_match(0);
        }
        _pigeon_advance(parser);
    }
    _pigeon_fail_at(parser, 1, start, expr->want);
    return _pigeon_null_match(1);
}

static PigeonMatch _pigeon_finish_values(PigeonValue *values, int ok) {
    if (values->as.array.length == 0) {
        pigeon_value_free(values);
        return _pigeon_null_match(ok);
    }
    {
        PigeonValue result = *values;
        *values = pigeon_value_null();
        return _pigeon_match(result, ok);
    }
}

static int _pigeon_collect_value(PigeonParser *parser, PigeonValue *values, PigeonValue *value) {
    if (value->type == PIGEON_VALUE_NULL) {
        return 1;
    }
    if (!pigeon_value_array_push(values, *value)) {
        _pigeon_out_of_memory(parser);
        return 0;
    }
    *value = pigeon_value_null();
    return 1;
}

static PigeonMatch _pigeon_parse_one_or_more(PigeonParser *parser, const PigeonExpr *expr) {
    PigeonValue values = pigeon_value_array();
    int matched = 0;
    for (;;) {
        size_t before = parser->pt.offset;
        PigeonMatch result = _pigeon_parse_expr(parser, expr->inner);
        if (!result.ok) {
            pigeon_value_free(&result.value);
            return _pigeon_finish_values(&values, matched);
        }
        if (parser->pt.offset == before) {
            pigeon_value_free(&result.value);
            pigeon_value_free(&values);
            pigeon_parser_add_error(parser, "repetition matched empty input");
            return _pigeon_null_match(0);
        }
        matched = 1;
        if (!_pigeon_collect_value(parser, &values, &result.value)) {
            pigeon_value_free(&result.value);
            pigeon_value_free(&values);
            return _pigeon_null_match(0);
        }
    }
}

static PigeonMatch _pigeon_parse_zero_or_more(PigeonParser *parser, const PigeonExpr *expr) {
    PigeonValue values = pigeon_value_array();
    for (;;) {
        size_t before = parser->pt.offset;
        PigeonMatch result = _pigeon_parse_expr(parser, expr->inner);
        if (!result.ok) {
            pigeon_value_free(&result.value);
            return _pigeon_finish_values(&values, 1);
        }
        if (parser->pt.offset == before) {
            pigeon_value_free(&result.value);
            pigeon_value_free(&values);
            pigeon_parser_add_error(parser, "repetition matched empty input");
            return _pigeon_null_match(0);
        }
        if (!_pigeon_collect_value(parser, &values, &result.value)) {
            pigeon_value_free(&result.value);
            pigeon_value_free(&values);
            return _pigeon_null_match(0);
        }
    }
}

static PigeonMatch _pigeon_parse_zero_or_one(PigeonParser *parser, const PigeonExpr *expr) {
    PigeonMatch result = _pigeon_parse_expr(parser, expr->inner);
    if (!result.ok) {
        pigeon_value_free(&result.value);
        return _pigeon_null_match(1);
    }
    return _pigeon_match(result.value, 1);
}

static PigeonMatch _pigeon_parse_sequence(PigeonParser *parser, const PigeonExpr *expr) {
    PigeonSavepoint start = parser->pt;
    PigeonValue values = pigeon_value_array();
    size_t index;
    for (index = 0; index < expr->item_count; index++) {
        PigeonMatch result = _pigeon_parse_expr(parser, expr->items[index]);
        if (!result.ok) {
            pigeon_value_free(&result.value);
            pigeon_value_free(&values);
            parser->pt = start;
            return _pigeon_null_match(0);
        }
        if (!_pigeon_skipping(parser) && !_pigeon_collect_value(parser, &values, &result.value)) {
            pigeon_value_free(&result.value);
            pigeon_value_free(&values);
            parser->pt = start;
            return _pigeon_null_match(0);
        }
        pigeon_value_free(&result.value);
    }
    return _pigeon_finish_values(&values, 1);
}

static PigeonMatch _pigeon_parse_rule_ref(PigeonParser *parser, const PigeonExpr *expr) {
    const PigeonRule *rule = _pigeon_find_rule(parser, expr->name);
    if (rule == NULL) {
        pigeon_parser_add_error(parser, "undefined rule");
        return _pigeon_null_match(0);
    }
    return _pigeon_parse_rule(parser, rule);
}

static PigeonMatch _pigeon_parse_rule_index(PigeonParser *parser, const PigeonExpr *expr) {
    if (expr->index >= parser->grammar->rule_count) {
        pigeon_parser_add_error(parser, "undefined rule index");
        return _pigeon_null_match(0);
    }
    return _pigeon_parse_rule(parser, &parser->grammar->rules[expr->index]);
}

static int _pigeon_recovery_push(PigeonParser *parser, const PigeonExpr *expr) {
    PigeonRecovery *items;
    size_t capacity;
    if (parser->recovery_length == parser->recovery_capacity) {
        capacity = parser->recovery_capacity == 0 ? 4 : parser->recovery_capacity * 2;
        items = (PigeonRecovery *)realloc(parser->recovery_stack, capacity * sizeof(PigeonRecovery));
        if (items == NULL) {
            _pigeon_out_of_memory(parser);
            return 0;
        }
        parser->recovery_stack = items;
        parser->recovery_capacity = capacity;
    }
    parser->recovery_stack[parser->recovery_length++].expr = expr;
    return 1;
}

static PigeonMatch _pigeon_parse_recovery(PigeonParser *parser, const PigeonExpr *expr) {
    PigeonMatch result;
    if (!_pigeon_recovery_push(parser, expr)) {
        return _pigeon_null_match(0);
    }
    result = _pigeon_parse_expr(parser, expr->inner);
    parser->recovery_length--;
    return result;
}

static PigeonMatch _pigeon_parse_throw(PigeonParser *parser, const PigeonExpr *expr) {
    size_t stack_index;
    for (stack_index = parser->recovery_length; stack_index > 0; stack_index--) {
        const PigeonExpr *recovery = parser->recovery_stack[stack_index - 1].expr;
        size_t label_index;
        for (label_index = 0; label_index < recovery->failure_label_count; label_index++) {
            if (strcmp(recovery->failure_labels[label_index], expr->label) == 0) {
                PigeonMatch result = _pigeon_parse_expr(parser, recovery->recover_expr);
                if (result.ok) {
                    return result;
                }
                pigeon_value_free(&result.value);
            }
        }
    }
    return _pigeon_null_match(0);
}

static PigeonMatch _pigeon_parse_expr(PigeonParser *parser, const PigeonExpr *expr) {
    parser->expression_count++;
    if (parser->expression_count > parser->max_expressions) {
        pigeon_parser_add_error(parser, "max number of expressions parsed");
        return _pigeon_null_match(0);
    }
    if (expr == NULL) {
        pigeon_parser_add_error(parser, "invalid expression");
        return _pigeon_null_match(0);
    }
    switch (expr->kind) {
    case PIGEON_EXPR_ACTION: return _pigeon_parse_action(parser, expr);
    case PIGEON_EXPR_AND_CODE: return _pigeon_null_match(expr->predicate(parser) != 0);
    case PIGEON_EXPR_AND: return _pigeon_parse_and(parser, expr, 0);
    case PIGEON_EXPR_AND_LOGICAL: return _pigeon_parse_and(parser, expr, 1);
    case PIGEON_EXPR_ANY: return _pigeon_parse_any(parser);
    case PIGEON_EXPR_CHAR_CLASS: return _pigeon_parse_char_class(parser, expr);
    case PIGEON_EXPR_CHOICE: return _pigeon_parse_choice(parser, expr);
    case PIGEON_EXPR_CODE: return _pigeon_parse_code(parser, expr);
    case PIGEON_EXPR_LABELED: return _pigeon_parse_labeled(parser, expr);
    case PIGEON_EXPR_LITERAL: return _pigeon_parse_literal(parser, expr);
    case PIGEON_EXPR_NOT_CODE: return _pigeon_null_match(expr->predicate(parser) == 0);
    case PIGEON_EXPR_NOT: return _pigeon_parse_not(parser, expr, 0);
    case PIGEON_EXPR_NOT_LOGICAL: return _pigeon_parse_not(parser, expr, 1);
    case PIGEON_EXPR_ONE_OR_MORE: return _pigeon_parse_one_or_more(parser, expr);
    case PIGEON_EXPR_RECOVERY: return _pigeon_parse_recovery(parser, expr);
    case PIGEON_EXPR_RULE_REF: return _pigeon_parse_rule_ref(parser, expr);
    case PIGEON_EXPR_RULE_INDEX: return _pigeon_parse_rule_index(parser, expr);
    case PIGEON_EXPR_SEQUENCE: return _pigeon_parse_sequence(parser, expr);
    case PIGEON_EXPR_THROW: return _pigeon_parse_throw(parser, expr);
    case PIGEON_EXPR_ZERO_OR_MORE: return _pigeon_parse_zero_or_more(parser, expr);
    case PIGEON_EXPR_ZERO_OR_ONE: return _pigeon_parse_zero_or_one(parser, expr);
    default:
        pigeon_parser_add_error(parser, "unknown expression type");
        return _pigeon_null_match(0);
    }
}

static int _pigeon_compare_strings(const void *left, const void *right) {
    const char *const *a = (const char *const *)left;
    const char *const *b = (const char *const *)right;
    return strcmp(*a, *b);
}

static void _pigeon_no_match_error(PigeonParser *parser) {
    size_t index;
    size_t unique_count = 0;
    size_t message_length;
    size_t used;
    char *message;
    const char *filename = parser->filename == NULL ? "" : parser->filename;
    if (parser->expected_length > 1) {
        qsort(parser->max_fail_expected, parser->expected_length, sizeof(char *), _pigeon_compare_strings);
    }
    for (index = 0; index < parser->expected_length; index++) {
        char *value = parser->max_fail_expected[index];
        if (strcmp(value, "!.") == 0) {
            char *replacement = _pigeon_copy_string("EOF", 3);
            if (replacement == NULL) {
                _pigeon_out_of_memory(parser);
                return;
            }
            free(value);
            parser->max_fail_expected[index] = replacement;
        }
    }
    if (parser->expected_length > 1) {
        qsort(parser->max_fail_expected, parser->expected_length, sizeof(char *), _pigeon_compare_strings);
    }
    for (index = 0; index < parser->expected_length; index++) {
        if (unique_count == 0 ||
            strcmp(parser->max_fail_expected[index], parser->max_fail_expected[unique_count - 1]) != 0) {
            parser->max_fail_expected[unique_count++] = parser->max_fail_expected[index];
        } else {
            free(parser->max_fail_expected[index]);
        }
    }
    parser->expected_length = unique_count;
    message_length = strlen(filename) + 96;
    for (index = 0; index < parser->expected_length; index++) {
        message_length += strlen(parser->max_fail_expected[index]) + 2;
    }
    message = (char *)malloc(message_length);
    if (message == NULL) {
        _pigeon_out_of_memory(parser);
        return;
    }
    used = (size_t)snprintf(message, message_length, "%s%s%lu:%lu: no match found, expected: ",
                            filename, filename[0] == '\0' ? "" : ":",
                            (unsigned long)parser->max_fail_pos.line,
                            (unsigned long)parser->max_fail_pos.col);
    for (index = 0; index < parser->expected_length && used < message_length; index++) {
        used += (size_t)snprintf(message + used, message_length - used, "%s%s",
                                 index == 0 ? "" : ", ", parser->max_fail_expected[index]);
    }
    parser->error = message;
    parser->error_position = _pigeon_position(parser->max_fail_pos);
}

static void _pigeon_parser_cleanup(PigeonParser *parser) {
    while (parser->vstack_length > 0) {
        _pigeon_scope_pop(parser);
    }
    free(parser->vstack);
    free(parser->recovery_stack);
    _pigeon_clear_expected(parser);
    free(parser->max_fail_expected);
    free(parser->cur.text);
}

static PigeonResult _pigeon_parse_grammar(const PigeonGrammar *grammar,
                                          const char *default_entrypoint,
                                          const char *input,
                                          size_t length,
                                          const PigeonParserOptions *options) {
    PigeonParser parser;
    PigeonResult output;
    const PigeonRule *start;
    PigeonMatch result;
    memset(&parser, 0, sizeof(parser));
    memset(&output, 0, sizeof(output));
    output.value = pigeon_value_null();
    parser.filename = options != NULL && options->filename != NULL ? options->filename : "";
    parser.input = (const unsigned char *)(input == NULL ? "" : input);
    parser.input_length = input == NULL ? 0 : length;
    parser.entrypoint = options != NULL && options->entrypoint != NULL
                        ? options->entrypoint : default_entrypoint;
    parser.max_expressions = options != NULL && options->max_expressions != 0
                             ? options->max_expressions : (size_t)INT_MAX;
    parser.pt.line = 1;
    parser.pt.col = 1;
    parser.cur.pos = _pigeon_position(parser.pt);
    parser.cur.data = options == NULL ? NULL : options->custom_data;
    parser.max_fail_pos = parser.pt;
    parser.grammar = grammar;

    if (grammar == NULL || grammar->rule_count == 0) {
        pigeon_parser_add_error(&parser, "grammar has no rule");
        result = _pigeon_null_match(0);
    } else {
        start = _pigeon_find_rule(&parser, parser.entrypoint);
        if (start == NULL) {
            pigeon_parser_add_error(&parser, "invalid entrypoint");
            result = _pigeon_null_match(0);
        } else {
            result = _pigeon_parse_rule(&parser, start);
        }
    }
    if (!result.ok && parser.error == NULL) {
        _pigeon_no_match_error(&parser);
    }
    if (parser.error == NULL) {
        output.ok = 1;
        output.value = result.value;
        result.value = pigeon_value_null();
    } else {
        output.error = parser.error;
        output.position = parser.error_position;
        parser.error = NULL;
    }
    pigeon_value_free(&result.value);
    free(parser.error);
    _pigeon_parser_cleanup(&parser);
    return output;
}

void pigeon_result_free(PigeonResult *result) {
    if (result == NULL) {
        return;
    }
    pigeon_value_free(&result->value);
    free(result->error);
    memset(result, 0, sizeof(*result));
    result->value = pigeon_value_null();
}
`

const runtimeFooter = `
PigeonResult pigeon_parse_with_options_n(const char *input, size_t length,
                                         const PigeonParserOptions *options) {
    return _pigeon_parse_grammar(&{{ .GrammarVarName }}, "{{ .Entrypoint }}", input, length, options);
}

PigeonResult pigeon_parse_with_options(const char *input, const PigeonParserOptions *options) {
    return pigeon_parse_with_options_n(input, input == NULL ? 0 : strlen(input), options);
}

PigeonResult pigeon_parse_n(const char *input, size_t length) {
    return pigeon_parse_with_options_n(input, length, NULL);
}

PigeonResult pigeon_parse(const char *input) {
    return pigeon_parse_with_options(input, NULL);
}
`
