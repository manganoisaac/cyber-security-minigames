const MODAL_SPEED_MS = 20;
const MODAL_SELECTOR = ".modal";
const MODAL_TEXT_SELECTOR = ".modal__text";
const MODAL_BUTTON_SELECTOR = ".modal__button";
const MODAL_TYPING_CLASS = "modal__text--typing";

const animations = new WeakMap();
const typed = new WeakSet();

function typewriteModal(modal) {
  const textEl = modal.querySelector(MODAL_TEXT_SELECTOR);
  if (!textEl || animations.has(textEl) || typed.has(textEl)) return;

  const buttonEl = modal.querySelector(MODAL_BUTTON_SELECTOR);
  const fullText = textEl.textContent;

  textEl.textContent = "";
  textEl.classList.add(MODAL_TYPING_CLASS);
  if (buttonEl) buttonEl.disabled = true;

  let i = 0;
  const interval = setInterval(() => {
    if (!document.body.contains(modal)) {
      clearInterval(interval);
      animations.delete(textEl);
      return;
    }
    if (i >= fullText.length) {
      clearInterval(interval);
      animations.delete(textEl);
      typed.add(textEl);
      textEl.classList.remove(MODAL_TYPING_CLASS);
      if (buttonEl) buttonEl.disabled = false;
      return;
    }
    textEl.textContent = fullText.slice(0, ++i);
  }, MODAL_SPEED_MS);

  animations.set(textEl, interval);
}

function typewriteNewModals() {
  document.querySelectorAll(MODAL_SELECTOR).forEach(typewriteModal);
}

document.addEventListener("DOMContentLoaded", typewriteNewModals);
document.addEventListener("htmx:afterSwap", typewriteNewModals);
document.addEventListener("htmx:afterSettle", typewriteNewModals);
