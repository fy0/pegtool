package builder_c

const runtimePreamble = `
#include <limits.h>
#include <stddef.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

typedef enum PegtoolValueType {
    PEGTOOL_VALUE_NULL = 0,
    PEGTOOL_VALUE_STRING,
    PEGTOOL_VALUE_INTEGER,
    PEGTOOL_VALUE_NUMBER,
    PEGTOOL_VALUE_BOOLEAN,
    PEGTOOL_VALUE_ARRAY,
    PEGTOOL_VALUE_POINTER
} PegtoolValueType;

typedef struct PegtoolValue PegtoolValue;

typedef struct PegtoolString {
    char *data;
    size_t length;
} PegtoolString;

typedef struct PegtoolArray {
    PegtoolValue *items;
    size_t length;
    size_t capacity;
} PegtoolArray;

struct PegtoolValue {
    PegtoolValueType type;
    union {
        PegtoolString string;
        long long integer;
        double number;
        int boolean;
        PegtoolArray array;
        void *pointer;
    } as;
};

typedef struct PegtoolPosition {
    size_t line;
    size_t col;
    size_t offset;
} PegtoolPosition;

typedef struct PegtoolCurrent {
    PegtoolPosition pos;
    char *text;
    size_t text_length;
    void *data;
} PegtoolCurrent;

typedef struct PegtoolParser PegtoolParser;
typedef PegtoolValue (*PegtoolActionFn)(PegtoolParser *parser);
typedef int (*PegtoolPredicateFn)(PegtoolParser *parser);

typedef struct PegtoolUnicodeRange {
    uint32_t lo;
    uint32_t hi;
    uint32_t stride;
} PegtoolUnicodeRange;

typedef enum PegtoolExprKind {
    PEGTOOL_EXPR_ACTION = 0,
    PEGTOOL_EXPR_AND_CODE,
    PEGTOOL_EXPR_AND,
    PEGTOOL_EXPR_AND_LOGICAL,
    PEGTOOL_EXPR_ANY,
    PEGTOOL_EXPR_CHAR_CLASS,
    PEGTOOL_EXPR_CHOICE,
    PEGTOOL_EXPR_CODE,
    PEGTOOL_EXPR_LABELED,
    PEGTOOL_EXPR_LITERAL,
    PEGTOOL_EXPR_NOT_CODE,
    PEGTOOL_EXPR_NOT,
    PEGTOOL_EXPR_NOT_LOGICAL,
    PEGTOOL_EXPR_ONE_OR_MORE,
    PEGTOOL_EXPR_RECOVERY,
    PEGTOOL_EXPR_RULE_REF,
    PEGTOOL_EXPR_RULE_INDEX,
    PEGTOOL_EXPR_SEQUENCE,
    PEGTOOL_EXPR_THROW,
    PEGTOOL_EXPR_ZERO_OR_MORE,
    PEGTOOL_EXPR_ZERO_OR_ONE
} PegtoolExprKind;

typedef struct PegtoolExpr PegtoolExpr;

struct PegtoolExpr {
    PegtoolExprKind kind;
    const PegtoolExpr *inner;
    const PegtoolExpr *const *items;
    size_t item_count;
    PegtoolActionFn action;
    PegtoolPredicateFn predicate;
    const PegtoolExpr *recover_expr;
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
    const PegtoolUnicodeRange *classes;
    size_t class_count;
    int ignore_case;
    int inverted;
    int not_skip;
    PegtoolPosition pos;
    int reserved;
};

typedef struct PegtoolRule {
    const char *name;
    const char *display_name;
    int var_exists;
    const PegtoolExpr *expr;
    PegtoolPosition pos;
    int reserved;
} PegtoolRule;

typedef struct PegtoolGrammar {
    const PegtoolRule *rules;
    size_t rule_count;
} PegtoolGrammar;

typedef struct PegtoolParserOptions {
    const char *filename;
    const char *entrypoint;
    void *custom_data;
    size_t max_expressions;
} PegtoolParserOptions;

typedef struct PegtoolResult {
    int ok;
    PegtoolValue value;
    char *error;
    PegtoolPosition position;
} PegtoolResult;

typedef struct PegtoolSavepoint {
    size_t line;
    size_t col;
    size_t offset;
    size_t index;
} PegtoolSavepoint;

typedef struct PegtoolMatch {
    PegtoolValue value;
    int ok;
} PegtoolMatch;

typedef struct PegtoolBinding {
    const char *name;
    PegtoolValue value;
} PegtoolBinding;

typedef struct PegtoolScope {
    PegtoolBinding *bindings;
    size_t length;
    size_t capacity;
} PegtoolScope;

typedef struct PegtoolRecovery {
    const PegtoolExpr *expr;
} PegtoolRecovery;

struct PegtoolParser {
    const char *filename;
    const unsigned char *input;
    size_t input_length;
    const char *entrypoint;
    size_t max_expressions;
    size_t expression_count;
    PegtoolSavepoint pt;
    PegtoolCurrent cur;
    const PegtoolGrammar *grammar;
    PegtoolScope *vstack;
    size_t vstack_length;
    size_t vstack_capacity;
    PegtoolRecovery *recovery_stack;
    size_t recovery_length;
    size_t recovery_capacity;
    size_t skip_depth;
    int invert_failures;
    PegtoolSavepoint max_fail_pos;
    char **max_fail_expected;
    size_t expected_length;
    size_t expected_capacity;
    char *error;
    PegtoolPosition error_position;
};

typedef struct PegtoolCaseMap {
    uint32_t from;
    uint32_t to;
} PegtoolCaseMap;

PegtoolValue pegtool_value_null(void);
PegtoolValue pegtool_value_string(const char *value);
PegtoolValue pegtool_value_string_n(const char *value, size_t length);
PegtoolValue pegtool_value_integer(long long value);
PegtoolValue pegtool_value_number(double value);
PegtoolValue pegtool_value_boolean(int value);
PegtoolValue pegtool_value_array(void);
PegtoolValue pegtool_value_pointer(void *value);
PegtoolValue pegtool_value_clone(const PegtoolValue *value);
int pegtool_value_array_push(PegtoolValue *array, PegtoolValue value);
void pegtool_value_free(PegtoolValue *value);
const PegtoolValue *pegtool_parser_label(PegtoolParser *parser, const char *name);
void pegtool_parser_add_error(PegtoolParser *parser, const char *message);
PegtoolResult pegtool_parse(const char *input);
PegtoolResult pegtool_parse_n(const char *input, size_t length);
PegtoolResult pegtool_parse_with_options(const char *input, const PegtoolParserOptions *options);
PegtoolResult pegtool_parse_with_options_n(const char *input, size_t length, const PegtoolParserOptions *options);
void pegtool_result_free(PegtoolResult *result);
`

const runtimeCore = `
static char *_pegtool_copy_string(const char *value, size_t length) {
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

PegtoolValue pegtool_value_null(void) {
    PegtoolValue value;
    memset(&value, 0, sizeof(value));
    value.type = PEGTOOL_VALUE_NULL;
    return value;
}

PegtoolValue pegtool_value_string_n(const char *text, size_t length) {
    PegtoolValue value = pegtool_value_null();
    if (text == NULL) {
        length = 0;
    }
    char *copy = _pegtool_copy_string(text, length);
    if (copy == NULL) {
        return value;
    }
    value.type = PEGTOOL_VALUE_STRING;
    value.as.string.data = copy;
    value.as.string.length = length;
    return value;
}

PegtoolValue pegtool_value_string(const char *text) {
    return pegtool_value_string_n(text, text == NULL ? 0 : strlen(text));
}

PegtoolValue pegtool_value_integer(long long integer) {
    PegtoolValue value = pegtool_value_null();
    value.type = PEGTOOL_VALUE_INTEGER;
    value.as.integer = integer;
    return value;
}

PegtoolValue pegtool_value_number(double number) {
    PegtoolValue value = pegtool_value_null();
    value.type = PEGTOOL_VALUE_NUMBER;
    value.as.number = number;
    return value;
}

PegtoolValue pegtool_value_boolean(int boolean) {
    PegtoolValue value = pegtool_value_null();
    value.type = PEGTOOL_VALUE_BOOLEAN;
    value.as.boolean = boolean != 0;
    return value;
}

PegtoolValue pegtool_value_array(void) {
    PegtoolValue value = pegtool_value_null();
    value.type = PEGTOOL_VALUE_ARRAY;
    return value;
}

PegtoolValue pegtool_value_pointer(void *pointer) {
    PegtoolValue value = pegtool_value_null();
    value.type = PEGTOOL_VALUE_POINTER;
    value.as.pointer = pointer;
    return value;
}

int pegtool_value_array_push(PegtoolValue *array, PegtoolValue value) {
    PegtoolValue *items;
    size_t capacity;
    if (array == NULL || array->type != PEGTOOL_VALUE_ARRAY) {
        return 0;
    }
    if (array->as.array.length == array->as.array.capacity) {
        capacity = array->as.array.capacity == 0 ? 4 : array->as.array.capacity * 2;
        if (capacity < array->as.array.capacity || capacity > SIZE_MAX / sizeof(PegtoolValue)) {
            return 0;
        }
        items = (PegtoolValue *)realloc(array->as.array.items, capacity * sizeof(PegtoolValue));
        if (items == NULL) {
            return 0;
        }
        array->as.array.items = items;
        array->as.array.capacity = capacity;
    }
    array->as.array.items[array->as.array.length++] = value;
    return 1;
}

PegtoolValue pegtool_value_clone(const PegtoolValue *source) {
    PegtoolValue value;
    size_t index;
    if (source == NULL) {
        return pegtool_value_null();
    }
    switch (source->type) {
    case PEGTOOL_VALUE_STRING:
        return pegtool_value_string_n(source->as.string.data, source->as.string.length);
    case PEGTOOL_VALUE_INTEGER:
        return pegtool_value_integer(source->as.integer);
    case PEGTOOL_VALUE_NUMBER:
        return pegtool_value_number(source->as.number);
    case PEGTOOL_VALUE_BOOLEAN:
        return pegtool_value_boolean(source->as.boolean);
    case PEGTOOL_VALUE_POINTER:
        return pegtool_value_pointer(source->as.pointer);
    case PEGTOOL_VALUE_ARRAY:
        value = pegtool_value_array();
        for (index = 0; index < source->as.array.length; index++) {
            PegtoolValue item = pegtool_value_clone(&source->as.array.items[index]);
            if (!pegtool_value_array_push(&value, item)) {
                pegtool_value_free(&item);
                pegtool_value_free(&value);
                return pegtool_value_null();
            }
        }
        return value;
    case PEGTOOL_VALUE_NULL:
    default:
        return pegtool_value_null();
    }
}

void pegtool_value_free(PegtoolValue *value) {
    size_t index;
    if (value == NULL) {
        return;
    }
    if (value->type == PEGTOOL_VALUE_STRING) {
        free(value->as.string.data);
    } else if (value->type == PEGTOOL_VALUE_ARRAY) {
        for (index = 0; index < value->as.array.length; index++) {
            pegtool_value_free(&value->as.array.items[index]);
        }
        free(value->as.array.items);
    }
    *value = pegtool_value_null();
}

static PegtoolPosition _pegtool_position(PegtoolSavepoint point) {
    PegtoolPosition value;
    value.line = point.line;
    value.col = point.col;
    value.offset = point.offset;
    return value;
}

void pegtool_parser_add_error(PegtoolParser *parser, const char *message) {
    if (parser == NULL || parser->error != NULL) {
        return;
    }
    parser->error = _pegtool_copy_string(message == NULL ? "parser error" : message,
                                        message == NULL ? 12 : strlen(message));
    parser->error_position = _pegtool_position(parser->pt);
}

static void _pegtool_out_of_memory(PegtoolParser *parser) {
    pegtool_parser_add_error(parser, "out of memory");
}

static int _pegtool_scope_push(PegtoolParser *parser) {
    PegtoolScope *scopes;
    size_t capacity;
    if (parser->vstack_length == parser->vstack_capacity) {
        capacity = parser->vstack_capacity == 0 ? 4 : parser->vstack_capacity * 2;
        scopes = (PegtoolScope *)realloc(parser->vstack, capacity * sizeof(PegtoolScope));
        if (scopes == NULL) {
            _pegtool_out_of_memory(parser);
            return 0;
        }
        parser->vstack = scopes;
        parser->vstack_capacity = capacity;
    }
    memset(&parser->vstack[parser->vstack_length], 0, sizeof(PegtoolScope));
    parser->vstack_length++;
    return 1;
}

static void _pegtool_scope_pop(PegtoolParser *parser) {
    PegtoolScope *scope;
    size_t index;
    if (parser->vstack_length == 0) {
        return;
    }
    scope = &parser->vstack[parser->vstack_length - 1];
    for (index = 0; index < scope->length; index++) {
        pegtool_value_free(&scope->bindings[index].value);
    }
    free(scope->bindings);
    memset(scope, 0, sizeof(*scope));
    parser->vstack_length--;
}

static int _pegtool_scope_set(PegtoolParser *parser, const char *name, PegtoolValue value) {
    PegtoolScope *scope;
    PegtoolBinding *bindings;
    size_t index;
    size_t capacity;
    if (parser->vstack_length == 0) {
        pegtool_value_free(&value);
        pegtool_parser_add_error(parser, "label used outside a value scope");
        return 0;
    }
    scope = &parser->vstack[parser->vstack_length - 1];
    for (index = 0; index < scope->length; index++) {
        if (strcmp(scope->bindings[index].name, name) == 0) {
            pegtool_value_free(&scope->bindings[index].value);
            scope->bindings[index].value = value;
            return 1;
        }
    }
    if (scope->length == scope->capacity) {
        capacity = scope->capacity == 0 ? 4 : scope->capacity * 2;
        bindings = (PegtoolBinding *)realloc(scope->bindings, capacity * sizeof(PegtoolBinding));
        if (bindings == NULL) {
            pegtool_value_free(&value);
            _pegtool_out_of_memory(parser);
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

const PegtoolValue *pegtool_parser_label(PegtoolParser *parser, const char *name) {
    PegtoolScope *scope;
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

static void _pegtool_clear_expected(PegtoolParser *parser) {
    size_t index;
    for (index = 0; index < parser->expected_length; index++) {
        free(parser->max_fail_expected[index]);
    }
    parser->expected_length = 0;
}

static void _pegtool_fail_at(PegtoolParser *parser, int matched, PegtoolSavepoint point, const char *want) {
    char **items;
    char *expected;
    size_t length;
    size_t capacity;
    if ((matched != 0) != (parser->invert_failures != 0) || point.offset < parser->max_fail_pos.offset) {
        return;
    }
    if (point.offset > parser->max_fail_pos.offset) {
        parser->max_fail_pos = point;
        _pegtool_clear_expected(parser);
    }
    want = want == NULL ? "" : want;
    length = strlen(want) + (parser->invert_failures ? 1 : 0);
    expected = (char *)malloc(length + 1);
    if (expected == NULL) {
        _pegtool_out_of_memory(parser);
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
            _pegtool_out_of_memory(parser);
            return;
        }
        parser->max_fail_expected = items;
        parser->expected_capacity = capacity;
    }
    parser->max_fail_expected[parser->expected_length++] = expected;
}

static uint32_t _pegtool_decode(const unsigned char *input, size_t length, size_t index, size_t *width) {
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

static uint32_t _pegtool_current(PegtoolParser *parser, size_t *width) {
    return _pegtool_decode(parser->input, parser->input_length, parser->pt.index, width);
}

static void _pegtool_advance(PegtoolParser *parser) {
    size_t width;
    uint32_t current = _pegtool_current(parser, &width);
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

static uint32_t _pegtool_lower(uint32_t value) {
    size_t low = 0;
    size_t high = _pegtool_case_map_count;
    while (low < high) {
        size_t middle = low + (high - low) / 2;
        if (_pegtool_case_map[middle].from < value) {
            low = middle + 1;
        } else {
            high = middle;
        }
    }
    if (low < _pegtool_case_map_count && _pegtool_case_map[low].from == value) {
        return _pegtool_case_map[low].to;
    }
    return value;
}

static char *_pegtool_slice(PegtoolParser *parser, PegtoolSavepoint start) {
    size_t length = parser->pt.index - start.index;
    return _pegtool_copy_string((const char *)parser->input + start.index, length);
}

static int _pegtool_set_current(PegtoolParser *parser, PegtoolSavepoint start) {
    char *text = _pegtool_slice(parser, start);
    if (text == NULL) {
        _pegtool_out_of_memory(parser);
        return 0;
    }
    free(parser->cur.text);
    parser->cur.text = text;
    parser->cur.text_length = parser->pt.index - start.index;
    parser->cur.pos = _pegtool_position(start);
    return 1;
}

static PegtoolMatch _pegtool_match(PegtoolValue value, int ok) {
    PegtoolMatch match;
    match.value = value;
    match.ok = ok;
    return match;
}

static PegtoolMatch _pegtool_null_match(int ok) {
    return _pegtool_match(pegtool_value_null(), ok);
}

static int _pegtool_skipping(PegtoolParser *parser) {
    return parser->skip_depth > 0;
}

static PegtoolMatch _pegtool_parse_expr(PegtoolParser *parser, const PegtoolExpr *expr);

static const PegtoolRule *_pegtool_find_rule(PegtoolParser *parser, const char *name) {
    size_t index;
    for (index = 0; index < parser->grammar->rule_count; index++) {
        if (strcmp(parser->grammar->rules[index].name, name) == 0) {
            return &parser->grammar->rules[index];
        }
    }
    return NULL;
}

static PegtoolMatch _pegtool_parse_rule(PegtoolParser *parser, const PegtoolRule *rule) {
    PegtoolMatch result;
    int scoped = rule->var_exists && !_pegtool_skipping(parser);
    if (scoped && !_pegtool_scope_push(parser)) {
        return _pegtool_null_match(0);
    }
    result = _pegtool_parse_expr(parser, rule->expr);
    if (scoped) {
        _pegtool_scope_pop(parser);
    }
    return result;
}

static PegtoolMatch _pegtool_parse_action(PegtoolParser *parser, const PegtoolExpr *expr) {
    PegtoolSavepoint start;
    PegtoolMatch result;
    if (_pegtool_skipping(parser)) {
        result = _pegtool_parse_expr(parser, expr->inner);
        pegtool_value_free(&result.value);
        return _pegtool_null_match(result.ok);
    }
    start = parser->pt;
    result = _pegtool_parse_expr(parser, expr->inner);
    if (result.ok) {
        pegtool_value_free(&result.value);
        if (!_pegtool_set_current(parser, start)) {
            return _pegtool_null_match(0);
        }
        result.value = expr->action(parser);
    }
    return result;
}

static PegtoolMatch _pegtool_parse_and(PegtoolParser *parser, const PegtoolExpr *expr, int logical) {
    PegtoolSavepoint start = parser->pt;
    PegtoolMatch result;
    int consumed;
    parser->skip_depth++;
    result = _pegtool_parse_expr(parser, expr->inner);
    parser->skip_depth--;
    consumed = parser->pt.offset != start.offset;
    parser->pt = start;
    pegtool_value_free(&result.value);
    return _pegtool_null_match(logical ? result.ok && consumed : result.ok);
}

static PegtoolMatch _pegtool_parse_not(PegtoolParser *parser, const PegtoolExpr *expr, int logical) {
    PegtoolSavepoint start = parser->pt;
    PegtoolMatch result;
    int consumed;
    parser->invert_failures = !parser->invert_failures;
    parser->skip_depth++;
    result = _pegtool_parse_expr(parser, expr->inner);
    parser->skip_depth--;
    parser->invert_failures = !parser->invert_failures;
    consumed = parser->pt.offset != start.offset;
    parser->pt = start;
    pegtool_value_free(&result.value);
    return _pegtool_null_match(logical ? !result.ok && consumed : !result.ok);
}

static PegtoolMatch _pegtool_parse_any(PegtoolParser *parser) {
    size_t width;
    (void)_pegtool_current(parser, &width);
    if (width == 0) {
        _pegtool_fail_at(parser, 0, parser->pt, ".");
        return _pegtool_null_match(0);
    }
    _pegtool_fail_at(parser, 1, parser->pt, ".");
    _pegtool_advance(parser);
    return _pegtool_null_match(1);
}

static PegtoolMatch _pegtool_parse_char_class(PegtoolParser *parser, const PegtoolExpr *expr) {
    size_t width;
    size_t index;
    uint32_t current = _pegtool_current(parser, &width);
    int found = 0;
    int matched;
    if (width == 0) {
        _pegtool_fail_at(parser, 0, parser->pt, expr->value);
        return _pegtool_null_match(0);
    }
    if (expr->ignore_case) {
        current = _pegtool_lower(current);
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
        const PegtoolUnicodeRange *range = &expr->classes[index];
        if (current >= range->lo && current <= range->hi &&
            (current - range->lo) % range->stride == 0) {
            found = 1;
        }
    }
    matched = expr->inverted ? !found : found;
    _pegtool_fail_at(parser, matched, parser->pt, expr->value);
    if (!matched) {
        return _pegtool_null_match(0);
    }
    _pegtool_advance(parser);
    return _pegtool_null_match(1);
}

static PegtoolMatch _pegtool_parse_choice(PegtoolParser *parser, const PegtoolExpr *expr) {
    size_t index;
    for (index = 0; index < expr->item_count; index++) {
        PegtoolMatch result = _pegtool_parse_expr(parser, expr->items[index]);
        if (result.ok) {
            return result;
        }
        pegtool_value_free(&result.value);
    }
    return _pegtool_null_match(0);
}

static PegtoolMatch _pegtool_parse_code(PegtoolParser *parser, const PegtoolExpr *expr) {
    if (!expr->not_skip && _pegtool_skipping(parser)) {
        return _pegtool_null_match(1);
    }
    return _pegtool_match(expr->action(parser), 1);
}

static PegtoolMatch _pegtool_parse_labeled(PegtoolParser *parser, const PegtoolExpr *expr) {
    PegtoolSavepoint start = parser->pt;
    PegtoolMatch result = _pegtool_parse_expr(parser, expr->inner);
    if (result.ok && expr->label != NULL && !_pegtool_skipping(parser)) {
        PegtoolValue value;
        if (expr->text_capture) {
            value = pegtool_value_string_n((const char *)parser->input + start.index,
                                          parser->pt.index - start.index);
        } else {
            value = pegtool_value_clone(&result.value);
        }
        if (!_pegtool_scope_set(parser, expr->label, value)) {
            pegtool_value_free(&result.value);
            return _pegtool_null_match(0);
        }
    }
    return result;
}

static PegtoolMatch _pegtool_parse_literal(PegtoolParser *parser, const PegtoolExpr *expr) {
    PegtoolSavepoint start = parser->pt;
    size_t index;
    for (index = 0; index < expr->rune_count; index++) {
        size_t width;
        uint32_t current = _pegtool_current(parser, &width);
        if (expr->ignore_case) {
            current = _pegtool_lower(current);
        }
        if (width == 0 || current != expr->runes[index]) {
            _pegtool_fail_at(parser, 0, start, expr->want);
            parser->pt = start;
            return _pegtool_null_match(0);
        }
        _pegtool_advance(parser);
    }
    _pegtool_fail_at(parser, 1, start, expr->want);
    return _pegtool_null_match(1);
}

static PegtoolMatch _pegtool_finish_values(PegtoolValue *values, int ok) {
    if (values->as.array.length == 0) {
        pegtool_value_free(values);
        return _pegtool_null_match(ok);
    }
    {
        PegtoolValue result = *values;
        *values = pegtool_value_null();
        return _pegtool_match(result, ok);
    }
}

static int _pegtool_collect_value(PegtoolParser *parser, PegtoolValue *values, PegtoolValue *value) {
    if (value->type == PEGTOOL_VALUE_NULL) {
        return 1;
    }
    if (!pegtool_value_array_push(values, *value)) {
        _pegtool_out_of_memory(parser);
        return 0;
    }
    *value = pegtool_value_null();
    return 1;
}

static PegtoolMatch _pegtool_parse_one_or_more(PegtoolParser *parser, const PegtoolExpr *expr) {
    PegtoolValue values = pegtool_value_array();
    int matched = 0;
    for (;;) {
        size_t before = parser->pt.offset;
        PegtoolMatch result = _pegtool_parse_expr(parser, expr->inner);
        if (!result.ok) {
            pegtool_value_free(&result.value);
            return _pegtool_finish_values(&values, matched);
        }
        if (parser->pt.offset == before) {
            pegtool_value_free(&result.value);
            pegtool_value_free(&values);
            pegtool_parser_add_error(parser, "repetition matched empty input");
            return _pegtool_null_match(0);
        }
        matched = 1;
        if (!_pegtool_collect_value(parser, &values, &result.value)) {
            pegtool_value_free(&result.value);
            pegtool_value_free(&values);
            return _pegtool_null_match(0);
        }
    }
}

static PegtoolMatch _pegtool_parse_zero_or_more(PegtoolParser *parser, const PegtoolExpr *expr) {
    PegtoolValue values = pegtool_value_array();
    for (;;) {
        size_t before = parser->pt.offset;
        PegtoolMatch result = _pegtool_parse_expr(parser, expr->inner);
        if (!result.ok) {
            pegtool_value_free(&result.value);
            return _pegtool_finish_values(&values, 1);
        }
        if (parser->pt.offset == before) {
            pegtool_value_free(&result.value);
            pegtool_value_free(&values);
            pegtool_parser_add_error(parser, "repetition matched empty input");
            return _pegtool_null_match(0);
        }
        if (!_pegtool_collect_value(parser, &values, &result.value)) {
            pegtool_value_free(&result.value);
            pegtool_value_free(&values);
            return _pegtool_null_match(0);
        }
    }
}

static PegtoolMatch _pegtool_parse_zero_or_one(PegtoolParser *parser, const PegtoolExpr *expr) {
    PegtoolMatch result = _pegtool_parse_expr(parser, expr->inner);
    if (!result.ok) {
        pegtool_value_free(&result.value);
        return _pegtool_null_match(1);
    }
    return _pegtool_match(result.value, 1);
}

static PegtoolMatch _pegtool_parse_sequence(PegtoolParser *parser, const PegtoolExpr *expr) {
    PegtoolSavepoint start = parser->pt;
    PegtoolValue values = pegtool_value_array();
    size_t index;
    for (index = 0; index < expr->item_count; index++) {
        PegtoolMatch result = _pegtool_parse_expr(parser, expr->items[index]);
        if (!result.ok) {
            pegtool_value_free(&result.value);
            pegtool_value_free(&values);
            parser->pt = start;
            return _pegtool_null_match(0);
        }
        if (!_pegtool_skipping(parser) && !_pegtool_collect_value(parser, &values, &result.value)) {
            pegtool_value_free(&result.value);
            pegtool_value_free(&values);
            parser->pt = start;
            return _pegtool_null_match(0);
        }
        pegtool_value_free(&result.value);
    }
    return _pegtool_finish_values(&values, 1);
}

static PegtoolMatch _pegtool_parse_rule_ref(PegtoolParser *parser, const PegtoolExpr *expr) {
    const PegtoolRule *rule = _pegtool_find_rule(parser, expr->name);
    if (rule == NULL) {
        pegtool_parser_add_error(parser, "undefined rule");
        return _pegtool_null_match(0);
    }
    return _pegtool_parse_rule(parser, rule);
}

static PegtoolMatch _pegtool_parse_rule_index(PegtoolParser *parser, const PegtoolExpr *expr) {
    if (expr->index >= parser->grammar->rule_count) {
        pegtool_parser_add_error(parser, "undefined rule index");
        return _pegtool_null_match(0);
    }
    return _pegtool_parse_rule(parser, &parser->grammar->rules[expr->index]);
}

static int _pegtool_recovery_push(PegtoolParser *parser, const PegtoolExpr *expr) {
    PegtoolRecovery *items;
    size_t capacity;
    if (parser->recovery_length == parser->recovery_capacity) {
        capacity = parser->recovery_capacity == 0 ? 4 : parser->recovery_capacity * 2;
        items = (PegtoolRecovery *)realloc(parser->recovery_stack, capacity * sizeof(PegtoolRecovery));
        if (items == NULL) {
            _pegtool_out_of_memory(parser);
            return 0;
        }
        parser->recovery_stack = items;
        parser->recovery_capacity = capacity;
    }
    parser->recovery_stack[parser->recovery_length++].expr = expr;
    return 1;
}

static PegtoolMatch _pegtool_parse_recovery(PegtoolParser *parser, const PegtoolExpr *expr) {
    PegtoolMatch result;
    if (!_pegtool_recovery_push(parser, expr)) {
        return _pegtool_null_match(0);
    }
    result = _pegtool_parse_expr(parser, expr->inner);
    parser->recovery_length--;
    return result;
}

static PegtoolMatch _pegtool_parse_throw(PegtoolParser *parser, const PegtoolExpr *expr) {
    size_t stack_index;
    for (stack_index = parser->recovery_length; stack_index > 0; stack_index--) {
        const PegtoolExpr *recovery = parser->recovery_stack[stack_index - 1].expr;
        size_t label_index;
        for (label_index = 0; label_index < recovery->failure_label_count; label_index++) {
            if (strcmp(recovery->failure_labels[label_index], expr->label) == 0) {
                PegtoolMatch result = _pegtool_parse_expr(parser, recovery->recover_expr);
                if (result.ok) {
                    return result;
                }
                pegtool_value_free(&result.value);
            }
        }
    }
    return _pegtool_null_match(0);
}

static PegtoolMatch _pegtool_parse_expr(PegtoolParser *parser, const PegtoolExpr *expr) {
    parser->expression_count++;
    if (parser->expression_count > parser->max_expressions) {
        pegtool_parser_add_error(parser, "max number of expressions parsed");
        return _pegtool_null_match(0);
    }
    if (expr == NULL) {
        pegtool_parser_add_error(parser, "invalid expression");
        return _pegtool_null_match(0);
    }
    switch (expr->kind) {
    case PEGTOOL_EXPR_ACTION: return _pegtool_parse_action(parser, expr);
    case PEGTOOL_EXPR_AND_CODE: return _pegtool_null_match(expr->predicate(parser) != 0);
    case PEGTOOL_EXPR_AND: return _pegtool_parse_and(parser, expr, 0);
    case PEGTOOL_EXPR_AND_LOGICAL: return _pegtool_parse_and(parser, expr, 1);
    case PEGTOOL_EXPR_ANY: return _pegtool_parse_any(parser);
    case PEGTOOL_EXPR_CHAR_CLASS: return _pegtool_parse_char_class(parser, expr);
    case PEGTOOL_EXPR_CHOICE: return _pegtool_parse_choice(parser, expr);
    case PEGTOOL_EXPR_CODE: return _pegtool_parse_code(parser, expr);
    case PEGTOOL_EXPR_LABELED: return _pegtool_parse_labeled(parser, expr);
    case PEGTOOL_EXPR_LITERAL: return _pegtool_parse_literal(parser, expr);
    case PEGTOOL_EXPR_NOT_CODE: return _pegtool_null_match(expr->predicate(parser) == 0);
    case PEGTOOL_EXPR_NOT: return _pegtool_parse_not(parser, expr, 0);
    case PEGTOOL_EXPR_NOT_LOGICAL: return _pegtool_parse_not(parser, expr, 1);
    case PEGTOOL_EXPR_ONE_OR_MORE: return _pegtool_parse_one_or_more(parser, expr);
    case PEGTOOL_EXPR_RECOVERY: return _pegtool_parse_recovery(parser, expr);
    case PEGTOOL_EXPR_RULE_REF: return _pegtool_parse_rule_ref(parser, expr);
    case PEGTOOL_EXPR_RULE_INDEX: return _pegtool_parse_rule_index(parser, expr);
    case PEGTOOL_EXPR_SEQUENCE: return _pegtool_parse_sequence(parser, expr);
    case PEGTOOL_EXPR_THROW: return _pegtool_parse_throw(parser, expr);
    case PEGTOOL_EXPR_ZERO_OR_MORE: return _pegtool_parse_zero_or_more(parser, expr);
    case PEGTOOL_EXPR_ZERO_OR_ONE: return _pegtool_parse_zero_or_one(parser, expr);
    default:
        pegtool_parser_add_error(parser, "unknown expression type");
        return _pegtool_null_match(0);
    }
}

static int _pegtool_compare_strings(const void *left, const void *right) {
    const char *const *a = (const char *const *)left;
    const char *const *b = (const char *const *)right;
    return strcmp(*a, *b);
}

static void _pegtool_no_match_error(PegtoolParser *parser) {
    size_t index;
    size_t unique_count = 0;
    size_t message_length;
    size_t used;
    char *message;
    const char *filename = parser->filename == NULL ? "" : parser->filename;
    if (parser->expected_length > 1) {
        qsort(parser->max_fail_expected, parser->expected_length, sizeof(char *), _pegtool_compare_strings);
    }
    for (index = 0; index < parser->expected_length; index++) {
        char *value = parser->max_fail_expected[index];
        if (strcmp(value, "!.") == 0) {
            char *replacement = _pegtool_copy_string("EOF", 3);
            if (replacement == NULL) {
                _pegtool_out_of_memory(parser);
                return;
            }
            free(value);
            parser->max_fail_expected[index] = replacement;
        }
    }
    if (parser->expected_length > 1) {
        qsort(parser->max_fail_expected, parser->expected_length, sizeof(char *), _pegtool_compare_strings);
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
        _pegtool_out_of_memory(parser);
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
    parser->error_position = _pegtool_position(parser->max_fail_pos);
}

static void _pegtool_parser_cleanup(PegtoolParser *parser) {
    while (parser->vstack_length > 0) {
        _pegtool_scope_pop(parser);
    }
    free(parser->vstack);
    free(parser->recovery_stack);
    _pegtool_clear_expected(parser);
    free(parser->max_fail_expected);
    free(parser->cur.text);
}

static PegtoolResult _pegtool_parse_grammar(const PegtoolGrammar *grammar,
                                          const char *default_entrypoint,
                                          const char *input,
                                          size_t length,
                                          const PegtoolParserOptions *options) {
    PegtoolParser parser;
    PegtoolResult output;
    const PegtoolRule *start;
    PegtoolMatch result;
    memset(&parser, 0, sizeof(parser));
    memset(&output, 0, sizeof(output));
    output.value = pegtool_value_null();
    parser.filename = options != NULL && options->filename != NULL ? options->filename : "";
    parser.input = (const unsigned char *)(input == NULL ? "" : input);
    parser.input_length = input == NULL ? 0 : length;
    parser.entrypoint = options != NULL && options->entrypoint != NULL
                        ? options->entrypoint : default_entrypoint;
    parser.max_expressions = options != NULL && options->max_expressions != 0
                             ? options->max_expressions : (size_t)INT_MAX;
    parser.pt.line = 1;
    parser.pt.col = 1;
    parser.cur.pos = _pegtool_position(parser.pt);
    parser.cur.data = options == NULL ? NULL : options->custom_data;
    parser.max_fail_pos = parser.pt;
    parser.grammar = grammar;

    if (grammar == NULL || grammar->rule_count == 0) {
        pegtool_parser_add_error(&parser, "grammar has no rule");
        result = _pegtool_null_match(0);
    } else {
        start = _pegtool_find_rule(&parser, parser.entrypoint);
        if (start == NULL) {
            pegtool_parser_add_error(&parser, "invalid entrypoint");
            result = _pegtool_null_match(0);
        } else {
            result = _pegtool_parse_rule(&parser, start);
        }
    }
    if (!result.ok && parser.error == NULL) {
        _pegtool_no_match_error(&parser);
    }
    if (parser.error == NULL) {
        output.ok = 1;
        output.value = result.value;
        result.value = pegtool_value_null();
    } else {
        output.error = parser.error;
        output.position = parser.error_position;
        parser.error = NULL;
    }
    pegtool_value_free(&result.value);
    free(parser.error);
    _pegtool_parser_cleanup(&parser);
    return output;
}

void pegtool_result_free(PegtoolResult *result) {
    if (result == NULL) {
        return;
    }
    pegtool_value_free(&result->value);
    free(result->error);
    memset(result, 0, sizeof(*result));
    result->value = pegtool_value_null();
}
`

const runtimeFooter = `
PegtoolResult pegtool_parse_with_options_n(const char *input, size_t length,
                                         const PegtoolParserOptions *options) {
    return _pegtool_parse_grammar(&{{ .GrammarVarName }}, "{{ .Entrypoint }}", input, length, options);
}

PegtoolResult pegtool_parse_with_options(const char *input, const PegtoolParserOptions *options) {
    return pegtool_parse_with_options_n(input, input == NULL ? 0 : strlen(input), options);
}

PegtoolResult pegtool_parse_n(const char *input, size_t length) {
    return pegtool_parse_with_options_n(input, length, NULL);
}

PegtoolResult pegtool_parse(const char *input) {
    return pegtool_parse_with_options(input, NULL);
}
`
