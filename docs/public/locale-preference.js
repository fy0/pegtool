(() => {
  const storageKey = "pegtool-language";
  const english = "en";
  const chinese = "zh";

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
    return pathname === "/en" || pathname.startsWith("/en/")
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
  const isRoot = pathname === "/" || pathname === "/index.html";
  const storedPreference = readPreference();

  if (isRoot) {
    const preferredLocale = storedPreference || localeFromBrowser();
    writePreference(preferredLocale);

    if (preferredLocale === english) {
      window.location.replace(
        `/en/${window.location.search}${window.location.hash}`,
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
