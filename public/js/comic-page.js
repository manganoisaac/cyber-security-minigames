(function () {
  function initComicPage() {
    const page = document.querySelector("[data-comic-page]");

    if (!page || page.dataset.comicInitialised === "true") {
      return;
    }

    page.dataset.comicInitialised = "true";

    const panels = Array.from(
      page.querySelectorAll(".comic-page__panel")
    );

    const continueButton = page.querySelector("[data-comic-continue]");

    if (!panels.length || !continueButton) {
      return;
    }

    const prefersReducedMotion = window.matchMedia(
      "(prefers-reduced-motion: reduce)"
    ).matches;

    let finished = false;
    const timers = [];

    function clearTimers() {
      timers.forEach((timer) => window.clearTimeout(timer));
      timers.length = 0;
    }

    function revealEverything() {
      if (finished) {
        return;
      }

      finished = true;
      clearTimers();

      panels.forEach((panel) => {
        panel.classList.add("is-visible");
      });

      continueButton.classList.add("is-visible");
    }

    if (prefersReducedMotion) {
      revealEverything();
      return;
    }

    const firstPanelDelay = 400;
    const panelInterval = 1400;
    const continueDelay = 1000;

    panels.forEach((panel, index) => {
      const timer = window.setTimeout(() => {
        panel.classList.add("is-visible");
      }, firstPanelDelay + index * panelInterval);

      timers.push(timer);
    });

    const finalPanelTime =
      firstPanelDelay + (panels.length - 1) * panelInterval;

    const continueTimer = window.setTimeout(() => {
      continueButton.classList.add("is-visible");
      finished = true;
    }, finalPanelTime + continueDelay);

    timers.push(continueTimer);

    page.addEventListener("click", (event) => {
      if (finished) {
        return;
      }

      if (event.target.closest("[data-comic-continue]")) {
        return;
      }

      revealEverything();
    });
  }

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", initComicPage);
  } else {
    initComicPage();
  }

  document.addEventListener("htmx:afterSwap", initComicPage);
})();