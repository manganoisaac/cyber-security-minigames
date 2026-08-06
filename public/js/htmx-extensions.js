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
// Email sidebar scroll preservation
// Selecting an email re-renders the whole list (outerHTML swap),
// which resets .email-sidebar scrollTop to 0. Capture the
// position before the swap and restore it after.
// ------------------------------------------------------------
let emailSidebarScroll = 0;

htmx.on("htmx:beforeSwap", () => {
  const sidebar = document.querySelector(".email-sidebar");
  if (sidebar) emailSidebarScroll = sidebar.scrollTop;
});

htmx.on("htmx:afterSwap", () => {
  const sidebar = document.querySelector(".email-sidebar");
  if (sidebar) sidebar.scrollTop = emailSidebarScroll;
});
