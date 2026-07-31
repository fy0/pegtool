import { defineConfig, type DefaultTheme } from "vitepress";

const pegLanguage = {
  name: "peg",
  scopeName: "source.peg",
  patterns: [
    { include: "#comments" },
    { include: "#strings" },
    { include: "#character-classes" },
    {
      match: "\\b[A-Za-z_][A-Za-z0-9_]*(?=\\s*(?:=|<-|←|⟵))",
      name: "entity.name.function.rule.peg",
    },
    {
      match: "\\b[A-Za-z_][A-Za-z0-9_]*(?=\\s*:)",
      name: "variable.other.label.peg",
    },
    {
      match: "<-|←|⟵|=|/|&&|!!|[&!?*+]",
      name: "keyword.operator.peg",
    },
    {
      match: "%\\{[A-Za-z_][A-Za-z0-9_]*\\}|//\\{[^}]+\\}",
      name: "keyword.control.recovery.peg",
    },
  ],
  repository: {
    comments: {
      patterns: [
        { begin: "/\\*", end: "\\*/", name: "comment.block.peg" },
        { match: "//(?!\\{).*$", name: "comment.line.double-slash.peg" },
      ],
    },
    strings: {
      patterns: [
        { begin: "\"", end: "\"", name: "string.quoted.double.peg" },
        { begin: "'", end: "'", name: "string.quoted.single.peg" },
        { begin: "`", end: "`", name: "string.quoted.other.raw.peg" },
      ],
    },
    "character-classes": {
      patterns: [
        {
          begin: "\\[",
          end: "\\]i?",
          name: "string.regexp.character-class.peg",
        },
      ],
    },
  },
};

const zhThemeConfig = {
  nav: [
    { text: "指南", link: "/guide/getting-started" },
    { text: "目标语言", link: "/targets/" },
    { text: "参数选择", link: "/guide/performance-options" },
    { text: "参考", link: "/reference/cli" },
  ],
  sidebar: [
    {
      text: "开始",
      items: [
        { text: "项目概览", link: "/guide/overview" },
        { text: "快速开始", link: "/guide/getting-started" },
        { text: "从 Pigeon 迁移", link: "/guide/fork-semantics" },
      ],
    },
    {
      text: "编写文法",
      items: [
        { text: "PEG 语法", link: "/guide/grammar" },
        { text: "动作与值", link: "/guide/actions-and-values" },
        { text: "错误处理", link: "/guide/error-handling" },
        { text: "多文件与入口规则", link: "/guide/multiple-grammars" },
      ],
    },
    {
      text: "生成目标",
      items: [
        { text: "目标语言概览", link: "/targets/" },
        { text: "Go", link: "/targets/go" },
        { text: "Haxe", link: "/targets/haxe" },
        { text: "TypeScript", link: "/targets/typescript" },
        { text: "C#", link: "/targets/csharp" },
        { text: "C99", link: "/targets/c" },
        { text: "Rust", link: "/targets/rust" },
      ],
    },
    {
      text: "性能与参考",
      items: [
        { text: "性能参数选择", link: "/guide/performance-options" },
        { text: "命令行参考", link: "/reference/cli" },
        { text: "生成 API", link: "/reference/generated-api" },
        { text: "项目开发", link: "/reference/development" },
      ],
    },
  ],
  outline: { level: [2, 3], label: "本页内容" },
  docFooter: { prev: "上一篇", next: "下一篇" },
  darkModeSwitchLabel: "切换深色模式",
  sidebarMenuLabel: "菜单",
  returnToTopLabel: "返回顶部",
  langMenuLabel: "切换语言",
  skipToContentLabel: "跳到正文",
  lastUpdated: {
    text: "最后更新",
    formatOptions: { dateStyle: "medium", timeStyle: "short" },
  },
  footer: {
    message: "基于 BSD 3-Clause License 发布",
    copyright: "pegtool 文档",
  },
} satisfies DefaultTheme.Config;

const enThemeConfig = {
  nav: [
    { text: "Guide", link: "/en/guide/getting-started" },
    { text: "Targets", link: "/en/targets/" },
    { text: "Performance", link: "/en/guide/performance-options" },
    { text: "Reference", link: "/en/reference/cli" },
  ],
  sidebar: [
    {
      text: "Getting Started",
      items: [
        { text: "Overview", link: "/en/guide/overview" },
        { text: "Quick Start", link: "/en/guide/getting-started" },
        { text: "Migrating from Pigeon", link: "/en/guide/fork-semantics" },
      ],
    },
    {
      text: "Writing Grammars",
      items: [
        { text: "PEG Syntax", link: "/en/guide/grammar" },
        { text: "Actions and Values", link: "/en/guide/actions-and-values" },
        { text: "Error Handling", link: "/en/guide/error-handling" },
        {
          text: "Multiple Grammars and Entry Points",
          link: "/en/guide/multiple-grammars",
        },
      ],
    },
    {
      text: "Generation Targets",
      items: [
        { text: "Target Overview", link: "/en/targets/" },
        { text: "Go", link: "/en/targets/go" },
        { text: "Haxe", link: "/en/targets/haxe" },
        { text: "TypeScript", link: "/en/targets/typescript" },
        { text: "C#", link: "/en/targets/csharp" },
        { text: "C99", link: "/en/targets/c" },
        { text: "Rust", link: "/en/targets/rust" },
      ],
    },
    {
      text: "Performance and Reference",
      items: [
        {
          text: "Choosing Performance Options",
          link: "/en/guide/performance-options",
        },
        { text: "CLI Reference", link: "/en/reference/cli" },
        { text: "Generated API", link: "/en/reference/generated-api" },
        { text: "Development", link: "/en/reference/development" },
      ],
    },
  ],
  outline: { level: [2, 3], label: "On this page" },
  docFooter: { prev: "Previous page", next: "Next page" },
  darkModeSwitchLabel: "Toggle dark mode",
  sidebarMenuLabel: "Menu",
  returnToTopLabel: "Return to top",
  langMenuLabel: "Change language",
  skipToContentLabel: "Skip to content",
  lastUpdated: {
    text: "Last updated",
    formatOptions: { dateStyle: "medium", timeStyle: "short" },
  },
  footer: {
    message: "Released under the BSD 3-Clause License",
    copyright: "pegtool Documentation",
  },
} satisfies DefaultTheme.Config;

export default defineConfig({
  title: "pegtool",
  locales: {
    root: {
      label: "简体中文",
      lang: "zh-CN",
      description: "面向 Go、Haxe、TypeScript、C#、C99 与 Rust 的 PEG 解析器生成器",
      themeConfig: zhThemeConfig,
    },
    en: {
      label: "English",
      lang: "en-US",
      link: "/en/",
      description:
        "A PEG parser generator for Go, Haxe, TypeScript, C#, C99, and Rust",
      themeConfig: enThemeConfig,
    },
  },
  cleanUrls: true,
  lastUpdated: true,
  head: [
    ["link", { rel: "icon", type: "image/svg+xml", href: "/pegtool-mark.svg" }],
    ["meta", { name: "theme-color", content: "#0f766e" }],
    ["script", { src: "/locale-preference.js" }],
  ],
  markdown: {
    lineNumbers: true,
    languages: [pegLanguage],
  },
  themeConfig: {
    logo: "/pegtool-mark.svg",
    siteTitle: "pegtool",
    search: {
      provider: "local",
      options: {
        locales: {
          root: {
            translations: {
              button: {
                buttonText: "搜索文档",
                buttonAriaLabel: "搜索文档",
              },
              modal: {
                noResultsText: "没有找到相关内容",
                resetButtonTitle: "清除查询",
                footer: {
                  selectText: "选择",
                  navigateText: "切换",
                  closeText: "关闭",
                },
              },
            },
          },
          en: {
            translations: {
              button: {
                buttonText: "Search",
                buttonAriaLabel: "Search documentation",
              },
              modal: {
                noResultsText: "No results found",
                resetButtonTitle: "Reset search",
                footer: {
                  selectText: "Select",
                  navigateText: "Navigate",
                  closeText: "Close",
                },
              },
            },
          },
        },
      },
    },
    socialLinks: [
      { icon: "github", link: "https://github.com/fy0/pegtool" },
    ],
  },
});
