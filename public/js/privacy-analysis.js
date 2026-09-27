(() => {
  if (window.__caughtInTheNetPrivacyAnalysisInitialised) {
    return;
  }

  window.__caughtInTheNetPrivacyAnalysisInitialised = true;

  let activePostId = null;
  let activeCorrectOption = -1;
  let activeFeedback = "";
  let solved = false;

  function getModal() {
    return document.getElementById("privacy-analysis-modal");
  }

  function closeAnalysis() {
    const modal = getModal();
    if (!modal) {
      return;
    }

    modal.hidden = true;
    document.body.classList.remove("privacy-analysis-open");
    activePostId = null;
    activeCorrectOption = -1;
    activeFeedback = "";
    solved = false;
  }

  function openAnalysis(button) {
    const modal = getModal();
    if (!modal || button.disabled) {
      return;
    }

    const article = button.closest(".post-item");
    const postText = article?.querySelector(".post-item__content")?.textContent?.trim();
    const analysisData = article?.querySelector(".post-item__analysis-data");
    const optionNodes = analysisData?.querySelectorAll("[data-analysis-option]") ?? [];

    activePostId = Number.parseInt(button.dataset.postId ?? "", 10);
    activeCorrectOption = Number.parseInt(button.dataset.correctOption ?? "", 10);
    activeFeedback = button.dataset.feedback ?? "";
    solved = false;

    const postPreview = modal.querySelector("#privacy-analysis-post");
    const optionsContainer = modal.querySelector("#privacy-analysis-options");
    const feedback = modal.querySelector("#privacy-analysis-feedback");
    const continueButton = modal.querySelector("#privacy-analysis-continue");

    postPreview.textContent = postText || "Post unavailable";
    optionsContainer.replaceChildren();

    optionNodes.forEach((node, index) => {
      const optionButton = document.createElement("button");
      optionButton.type = "button";
      optionButton.className = "privacy-analysis-modal__option";
      optionButton.dataset.optionIndex = String(index);
      optionButton.textContent = node.textContent.trim();
      optionsContainer.appendChild(optionButton);
    });

    feedback.hidden = true;
    feedback.textContent = "";
    feedback.className = "privacy-analysis-modal__feedback";
    continueButton.hidden = true;

    modal.hidden = false;
    document.body.classList.add("privacy-analysis-open");

    const firstOption = optionsContainer.querySelector("button");
    firstOption?.focus();
  }

  function handleOptionSelection(optionButton) {
    if (solved) {
      return;
    }

    const modal = getModal();
    if (!modal) {
      return;
    }

    const selectedIndex = Number.parseInt(optionButton.dataset.optionIndex ?? "", 10);
    const feedback = modal.querySelector("#privacy-analysis-feedback");
    const continueButton = modal.querySelector("#privacy-analysis-continue");

    if (selectedIndex !== activeCorrectOption) {
      optionButton.classList.add("privacy-analysis-modal__option--incorrect");
      optionButton.disabled = true;

      feedback.textContent = "Not quite. Re-read the post and think about what a stranger could actually learn from it.";
      feedback.className = "privacy-analysis-modal__feedback privacy-analysis-modal__feedback--incorrect";
      feedback.hidden = false;
      return;
    }

    solved = true;
    optionButton.classList.add("privacy-analysis-modal__option--correct");

    modal.querySelectorAll(".privacy-analysis-modal__option").forEach((button) => {
      button.disabled = true;
    });

    feedback.textContent = `Correct. ${activeFeedback}`;
    feedback.className = "privacy-analysis-modal__feedback privacy-analysis-modal__feedback--correct";
    feedback.hidden = false;
    continueButton.hidden = false;
    continueButton.focus();
  }

  async function commitAnalysis() {
    if (!solved || !Number.isInteger(activePostId)) {
      return;
    }

    const postId = activePostId;
    closeAnalysis();

    await htmx.ajax("POST", `/privacy/clickpost/${postId}`, {
      target: `#privacy-post-${postId}`,
      swap: "outerHTML"
    });
  }

  document.addEventListener("click", (event) => {
    const analyseButton = event.target.closest("[data-privacy-analyse]");
    if (analyseButton) {
      openAnalysis(analyseButton);
      return;
    }

    const optionButton = event.target.closest(".privacy-analysis-modal__option");
    if (optionButton) {
      handleOptionSelection(optionButton);
      return;
    }

    if (event.target.closest("#privacy-analysis-continue")) {
      commitAnalysis();
      return;
    }

    if (event.target.closest("[data-privacy-analysis-close]")) {
      closeAnalysis();
    }
  });

  document.addEventListener("keydown", (event) => {
    if (event.key === "Escape" && !getModal()?.hidden) {
      closeAnalysis();
    }
  });
})();
