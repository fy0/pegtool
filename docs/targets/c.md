# C99

C target 输出可移植的 C99 源码和显式所有权 API。生成 parser 使用标准 C 库，不要求额外 runtime。

## 最小可运行示例

在空目录中手写下面两个文件。`parser.c` 由 Pegtool 生成。

### `identifier.peg`

```peg
Start = Identifier !. {
    return pigeon_value_string_n(c->text, c->text_length);
}

Identifier = [\pL_] [\pL\p{Nd}_]*
```

### `main.c`

```c
#include "parser.c"

int main(void) {
    PigeonResult result = pigeon_parse("name_世界");
    if (!result.ok) {
        fprintf(stderr, "%s\n", result.error);
        pigeon_result_free(&result);
        return 1;
    }
    printf("%s\n", result.value.as.string.data);
    pigeon_result_free(&result);
    return 0;
}
```

生成并编译：

```powershell
pegtool -t c -o parser.c identifier.peg
gcc -std=c99 -Wall -Wextra -Werror -pedantic-errors main.c -o identifier
.\identifier.exe
```

输出：

```text
name_世界
```

## 公开 API

| 函数 | 用途 |
| --- | --- |
| `pigeon_parse` | 解析 NUL 结尾字符串 |
| `pigeon_parse_n` | 解析指定字节长度，可包含 NUL |
| `pigeon_parse_with_options` | 使用 `PigeonParserOptions` |
| `pigeon_parse_with_options_n` | options 与显式长度组合 |
| `pigeon_result_free` | 释放成功值或错误结果拥有的内存 |

每次 parse 后都应调用 `pigeon_result_free`，无论成功还是失败。

## Action 与所有权

```peg
word = [\p{L}]+ {
    return pigeon_value_string_n(c->text, c->text_length);
}
```

标签是 `const PigeonValue *`，action 返回拥有所有权的 `PigeonValue`。需要保存标签值时调用 `pigeon_value_clone`；数组通过 `pigeon_value_array` 和 `pigeon_value_array_push` 构建。失败路径必须释放已经取得所有权的临时值。

与内部参数或 C99 关键字冲突的标签会被生成为 `_pigeon_label_<name>`。`-receiver-name` 可以修改当前匹配对象默认使用的 `c` 名称。

## Options

`PigeonParserOptions` 支持 `filename`、`entrypoint`、`custom_data` 和 `max_expressions`。C target 当前没有 memo 开关；`-optimize-parser` 也没有 C 专用效果。
