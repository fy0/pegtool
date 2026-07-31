(() => {
  const storageKey = "pegtool-language";
  const english = "en";
  const chinese = "zh";
  const scriptSuffix = "/locale-preference.js";
  const scriptPath = new URL(
    document.currentScript?.src || window.location.href,
    window.location.href,
  ).pathname;
  const basePath = scriptPath.endsWith(scriptSuffix)
    ? scriptPath.slice(0, -scriptSuffix.length)
    : "";

  function stripBase(pathname) {
    if (basePath && pathname === basePath) {
      return "/";
    }
    if (basePath && pathname.startsWith(`${basePath}/`)) {
      return pathname.slice(basePath.length);
    }
    return pathname;
  }

  function readPreference() {
    try {
      const value = window.localStorage.getItem(storageKey);
      return value === english || value === chinese ? value : null;
    } catch {
      return null;
    }
  }

  function writePreference(value) {
    try {
      window.localStorage.setItem(storageKey, value);
    } catch {
      // Language detection still works when storage is unavailable.
    }
  }

  function localeFromPath(pathname) {
    const relativePath = stripBase(pathname);
    return relativePath === "/en" || relativePath.startsWith("/en/")
      ? english
      : chinese;
  }

  function localeFromBrowser() {
    const languages =
      Array.isArray(window.navigator.languages) &&
      window.navigator.languages.length > 0
        ? window.navigator.languages
        : [window.navigator.language];

    for (const language of languages) {
      const normalized = String(language).toLowerCase();
      if (normalized === "zh" || normalized.startsWith("zh-")) {
        return chinese;
      }
      if (normalized === "en" || normalized.startsWith("en-")) {
        return english;
      }
    }

    return english;
  }

  const pathname = window.location.pathname;
  const relativePath = stripBase(pathname);
  const isRoot = relativePath === "/" || relativePath === "/index.html";
  const storedPreference = readPreference();

  if (isRoot) {
    const preferredLocale = storedPreference || localeFromBrowser();
    writePreference(preferredLocale);

    if (preferredLocale === english) {
      window.location.replace(
        `${basePath}/en/${window.location.search}${window.location.hash}`,
      );
    }
  } else if (!storedPreference) {
    writePreference(localeFromPath(pathname));
  }

  document.addEventListener("click", (event) => {
    const element = event.target instanceof Element ? event.target : null;
    const link = element?.closest("a[href]");

    if (!link || !link.closest(".translations")) {
      return;
    }

    const target = new URL(link.href, window.location.href);
    if (target.origin === window.location.origin) {
      writePreference(localeFromPath(target.pathname));
    }
  });
})();
