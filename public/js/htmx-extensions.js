// ============================================================
// htmx extensions
// ============================================================

// ------------------------------------------------------------
// Boosted-link transitions
// Give boosted <a> tags that lack an hx-swap a smooth
// transition so full-page loads don't snap.
// ------------------------------------------------------------
htmx.on("htmx:beforeRequest", (event) => {
  const elt = event.detail.elt;
  if (
    event.detail.boosted === true &&
    elt.tageName === "A" &&
    !elt.hasAttribute("hx-swap")
  ) {
    elt.setAttribute("hx-swap", "transition:true");
  }
});

// ------------------------------------------------------------
// Email selected-word count
// Word toggles pass the current number of selected words to the
// server so it can keep the submit-button badge in sync. The count
// is read from the DOM before the request is configured, so it's
// always current regardless of which word was clicked.
// ------------------------------------------------------------
htmx.on("htmx:configRequest", (event) => {
  const elt = event.detail.elt;
  if (elt.matches(".email-word, .email-word--selected")) {
    event.detail.parameters["selectedCount"] =
      document.querySelectorAll(".email-word--selected").length;
  } else if (elt.matches(".email-submit")) {
    event.detail.parameters["selectedIndexes"] = [
      ...document.querySelectorAll(".email-word--selected"),
    ]
      .map((el) => el.id.replace("email-word-", ""))
      .join(",");
  }
});

// ------------------------------------------------------------
// Email sidebar scroll preservation
// Selecting an email re-renders the whole list (outerHTML swap),
// which resets .email-list__sidebar scrollTop to 0. Capture the
// position before the swap and restore it after.
// ------------------------------------------------------------
let emailSidebarScroll = 0;

htmx.on("htmx:beforeSwap", () => {
  const sidebar = document.querySelector(".email-list__sidebar");
  if (sidebar) emailSidebarScroll = sidebar.scrollTop;
});

htmx.on("htmx:afterSwap", () => {
  const sidebar = document.querySelector(".email-list__sidebar");
  if (sidebar) sidebar.scrollTop = emailSidebarScroll;
});
