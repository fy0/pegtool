#include "parser.c"

int main(void) {
    const char *expected[] = {
        "alpha",
        "\344\270\226\347\225\214",
        "\360\220\220\200"
    };
    PegtoolResult result = pegtool_parse("alpha,\344\270\226\347\225\214,\360\220\220\200");
    size_t index;

    if (!result.ok || result.value.type != PEGTOOL_VALUE_ARRAY || result.value.as.array.length != 3) {
        fprintf(stderr, "unexpected parse result: %s\n", result.error == NULL ? "value" : result.error);
        pegtool_result_free(&result);
        return 1;
    }
    for (index = 0; index < result.value.as.array.length; index++) {
        const PegtoolValue *value = &result.value.as.array.items[index];
        if (value->type != PEGTOOL_VALUE_STRING || strcmp(value->as.string.data, expected[index]) != 0) {
            pegtool_result_free(&result);
            return 2;
        }
    }
    printf("%s|%s|%s\n", expected[0], expected[1], expected[2]);
    pegtool_result_free(&result);
    return 0;
}
