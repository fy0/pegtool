<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useData } from "vitepress";

const { lang } = useData();
const pageUrl = ref("");
const copied = ref(false);
const isEnglish = computed(() => lang.value.startsWith("en"));
const migrationUrl = computed(
  () =>
    pageUrl.value ||
    (isEnglish.value
      ? "/en/guide/fork-semantics"
      : "/guide/fork-semantics"),
);
const title = computed(() =>
  isEnglish.value ? "Code agent migration prompt" : "Code agent 迁移 Prompt",
);
const copyLabel = computed(() => {
  if (isEnglish.value) {
    return copied.value ? "Copied" : "Copy prompt";
  }
  return copied.value ? "已复制" : "复制 Prompt";
});

onMounted(() => {
  pageUrl.value = window.location.href.split("#", 1)[0];
});

const prompt = computed(() =>
  isEnglish.value
    ? `You are the code agent responsible for migrating a parser. Migrate the PEG grammar in the current repository from upstream Pigeon to pegtool.

Migration guide: ${migrationUrl.value}

Follow these steps:
1. Read the migration guide, then inspect the repository's .peg grammars, generation commands, public parser wrappers, generated files, and related tests.
2. Identify the Pigeon version, target language, and current public API before making changes. Do not perform an unverified global text replacement.
3. Replace the generation command with pegtool and select the correct -t option for the target language.
4. Follow the migration guide for action return values, predicates, terminal values, sequence/repetition values, text capture, custom data, entry points, and left recursion.
5. Preserve the parser's public API and application semantics where practical. Add explicit wrappers in the grammar initializer or regular source files when needed.
6. Do not edit the generated parser by hand to hide generator or grammar problems. Regenerate it after changing the grammar.
7. Compile the generated parser, run the existing tests, and add regression tests for behavior affected by the migration.
8. In the final report, list changed files, behavior differences, the generation command, test results, and remaining risks.

Constraints:
- Do not reduce Unicode character sets or modify unrelated parser algorithms.
- Do not enable -optimize-ref-expr-by-index by default. Use it only when benchmarks confirm the benefit and a large generated diff is acceptable.
- Choose memoization according to the target and workload; do not enable it mechanically.`
    : `你是负责 parser 迁移的 code agent。请把当前仓库中基于原版 Pigeon 的 PEG 文法迁移到 pegtool。

迁移规范：${migrationUrl.value}

请按以下步骤工作：
1. 阅读迁移规范，并检查仓库中的 .peg 文法、生成命令、公开 parser wrapper、生成文件和相关测试。
2. 先识别当前使用的 Pigeon 版本、目标语言和现有公开 API，不要做未经验证的全局文本替换。
3. 将生成命令改为 pegtool，并根据目标语言选择正确的 -t 参数。
4. 对照迁移规范处理 action 返回值、predicate、终结符值、sequence/repetition 值、文本捕获、自定义数据、入口规则和左递归。
5. 尽量保持现有 parser 的公开 API 和业务语义；需要 wrapper 时在文法 initializer 或普通源码中显式提供。
6. 不要手工修改生成 parser 来掩盖生成器或文法问题；修改文法后重新生成。
7. 编译生成 parser，运行现有测试，并为迁移涉及的行为补充回归测试。
8. 最后报告修改文件、行为差异、生成命令、测试结果和尚存风险。

约束：
- 不删减 Unicode 字符集，不改动与迁移无关的 parser 算法。
- 不默认启用 -optimize-ref-expr-by-index；只有基准确认且能接受大面积生成 diff 时才使用。
- memoization 必须按目标和负载选择，不能机械开启。`,
);

async function copyPrompt() {
  await navigator.clipboard.writeText(prompt.value);
  copied.value = true;
  window.setTimeout(() => {
    copied.value = false;
  }, 1800);
}
</script>

<template>
  <section class="migration-prompt" :aria-label="title">
    <div class="migration-prompt__header">
      <strong>{{ title }}</strong>
      <button type="button" @click="copyPrompt">
        {{ copyLabel }}
      </button>
    </div>
    <pre><code>{{ prompt }}</code></pre>
  </section>
</template>

<style scoped>
.migration-prompt {
  margin: 20px 0 32px;
  border: 1px solid var(--vp-c-divider);
  border-radius: 8px;
  background: var(--vp-c-bg-soft);
  overflow: hidden;
}

.migration-prompt__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 12px 16px;
  border-bottom: 1px solid var(--vp-c-divider);
}

.migration-prompt button {
  flex: none;
  min-height: 32px;
  padding: 0 12px;
  border: 1px solid var(--vp-c-brand-1);
  border-radius: 6px;
  color: var(--vp-c-brand-1);
  background: transparent;
  font: inherit;
  cursor: pointer;
}

.migration-prompt button:hover {
  color: var(--vp-c-bg);
  background: var(--vp-c-brand-1);
}

.migration-prompt pre {
  max-height: 520px;
  margin: 0;
  padding: 16px;
  overflow: auto;
  background: transparent;
  white-space: pre-wrap;
  overflow-wrap: anywhere;
}

.migration-prompt code {
  color: var(--vp-c-text-1);
  font-size: 13px;
  line-height: 1.65;
}
</style>
