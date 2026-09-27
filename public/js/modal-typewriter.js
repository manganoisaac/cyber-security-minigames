const MODAL_SPEED_MS = 20;
const MODAL_SELECTOR = ".modal";
const MODAL_TEXT_SELECTOR = ".modal__text";
const MODAL_BUTTON_SELECTOR = ".modal__button";
const MODAL_TYPING_CLASS = "modal__text--typing";

const animations = new WeakMap();
const typed = new WeakSet();

function typewriteModal(modal) {
  const textEl = modal.querySelector(MODAL_TEXT_SELECTOR);

  if (!textEl || animations.has(textEl) || typed.has(textEl)) {
    return;
  }

  const buttonEl = modal.querySelector(MODAL_BUTTON_SELECTOR);
  const fullText = textEl.textContent;

  textEl.textContent = "";
  textEl.classList.add(MODAL_TYPING_CLASS);

  if (buttonEl) {
    buttonEl.disabled = true;
  }

  let i = 0;
  let finished = false;

  function finishTyping() {
    if (finished) {
      return;
    }

    finished = true;

    const interval = animations.get(textEl);

    if (interval) {
      clearInterval(interval);
      animations.delete(textEl);
    }

    textEl.textContent = fullText;
    textEl.classList.remove(MODAL_TYPING_CLASS);
    typed.add(textEl);

    if (buttonEl) {
      buttonEl.disabled = false;
    }
  }

  const interval = setInterval(() => {
    if (!document.body.contains(modal)) {
      clearInterval(interval);
      animations.delete(textEl);
      return;
    }

    if (i >= fullText.length) {
      finishTyping();
      return;
    }

    textEl.textContent = fullText.slice(0, ++i);
  }, MODAL_SPEED_MS);

  animations.set(textEl, interval);

  /*
   * Privacy has longer instructional text, so allow the player to instantly
   * reveal the full message without advancing the modal.
   *
   * They must still press the modal button afterwards to continue.
   */
  if (modal.closest(".privacy-modal-root")) {
    modal.addEventListener(
      "click",
      (event) => {
        if (!finished) {
          event.preventDefault();
          finishTyping();
        }
      },
      { once: true }
    );
  }
}

function typewriteNewModals() {
  document.querySelectorAll(MODAL_SELECTOR).forEach(typewriteModal);
}

document.addEventListener("DOMContentLoaded", typewriteNewModals);
document.addEventListener("htmx:afterSwap", typewriteNewModals);
document.addEventListener("htmx:afterSettle", typewriteNewModals);