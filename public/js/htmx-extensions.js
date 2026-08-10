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
// Email badge counts
// The submit-button badge counts selected words, the checked-email
// badge counts checked checkboxes. Both counts are read from the
// DOM before the request is configured.
//
// Word toggles need a server-side increment/decrement because the
// word's --selected class is only applied after the response swaps
// in, so the DOM count lags the click. A native checkbox flips its
// checked state before the click handler runs instead, so by the
// time the request is configured the DOM count already reflects the
// toggle — the checked-email badge needs no adjustment.
// ------------------------------------------------------------
htmx.on("htmx:configRequest", (event) => {
  const elt = event.detail.elt;
  if (elt.matches(".email-word, .email-word--selected")) {
    event.detail.parameters["selectedCount"] = document.querySelectorAll(
      ".email-word--selected",
    ).length;
  } else if (elt.matches(".email-submit")) {
    event.detail.parameters["selectedIndexes"] = [
      ...document.querySelectorAll(".email-word--selected"),
    ]
      .map((el) => el.id.replace("email-word-", ""))
      .join(",");
  } else if (elt.matches(".email-checked")) {
    event.detail.parameters["checkedCount"] = document.querySelectorAll(
      ".email-item__checkbox:checked",
    ).length;
  } else if (elt.closest(".email-item")) {
    event.detail.parameters["checkedCount"] = document.querySelectorAll(
      ".email-item__checkbox:checked",
    ).length;
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
