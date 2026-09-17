(() => {
  "use strict";

  /*
   * ------------------------------------------------------------------------
   * GAME CONFIGURATION
   * ------------------------------------------------------------------------
   * Each scenario has a six-letter cybersecurity key and its own plaintext.
   * The ciphertext in the left terminal is generated automatically from the
   * plaintext using the selected key as a simple columnar-transposition key.
   *
   * continueUrl is temporary until the project's dedicated "what we learned"
   * route exists. Change it later without touching the game engine.
   */
  const GAME_CONFIG = {
    maxAttempts: 6,
    hiddenCharacter: "·",
    validWordsUrl: "/public/data/valid-words-5.txt",
    continueUrl: "/whwl/encryption",

scenarios: [
  {
    id: "email",
    secretWord: "EMAIL",
    plaintextMessage:
      "EUGENE HAS SCHEDULED THE SOURCE CODE TRANSFER FOR 2130. THE ARCHIVE WILL LEAVE THROUGH THE INTERNAL SERVER. BLOCK THE OUTBOUND CONNECTION BEFORE THE TRANSFER BEGINS."
  },
  {
    id: "phish",
    secretWord: "PHISH",
    plaintextMessage:
      "THE PROPRIETARY SOURCE CODE WILL BE EXFILTRATED AT 2130. EUGENE PLANS TO MOVE THE ARCHIVE THROUGH THE INTERNAL SERVER. STOP THE OUTBOUND TRANSFER BEFORE IT STARTS."
  },
  {
    id: "admin",
    secretWord: "ADMIN",
    plaintextMessage:
      "EUGENE WILL USE THE INTERNAL SERVER AT 2130 TO MOVE THE SOURCE CODE ARCHIVE OFF THE NETWORK. BLOCK THE OUTBOUND CONNECTION BEFORE THE FILE TRANSFER CAN COMPLETE."
  },
  {
    id: "login",
    secretWord: "LOGIN",
    plaintextMessage:
      "THE SOURCE CODE ARCHIVE IS READY FOR TRANSFER AT 2130. EUGENE WILL ROUTE IT THROUGH THE INTERNAL SERVER TO AN EXTERNAL DESTINATION. PREVENT THE OUTBOUND CONNECTION."
  },
  {
    id: "token",
    secretWord: "TOKEN",
    plaintextMessage:
      "AT 2130 EUGENE INTENDS TO SEND THE PROPRIETARY SOURCE CODE THROUGH THE INTERNAL SERVER AND OUT OF THE NETWORK. BLOCK THE EXTERNAL TRANSFER BEFORE IT BEGINS."
  },
  {
    id: "shell",
    secretWord: "SHELL",
    plaintextMessage:
      "EUGENE CONFIRMED A 2130 TRANSFER WINDOW FOR THE SOURCE CODE ARCHIVE. THE FILES WILL PASS THROUGH THE INTERNAL SERVER. STOP THE OUTBOUND NETWORK CONNECTION IN TIME."
  },
  {
    id: "virus",
    secretWord: "VIRUS",
    plaintextMessage:
      "THE SOURCE CODE THEFT IS PLANNED FOR 2130. EUGENE WILL USE THE INTERNAL SERVER TO PUSH THE ARCHIVE OUTSIDE THE NETWORK. BLOCK THE OUTBOUND TRANSFER BEFORE THEN."
  },
  {
    id: "proxy",
    secretWord: "PROXY",
    plaintextMessage:
      "EUGENE PLANS TO EXFILTRATE THE SOURCE CODE AT 2130 USING THE INTERNAL SERVER. THE ARCHIVE WILL BE SENT TO AN EXTERNAL DESTINATION. STOP THE OUTBOUND CONNECTION."
  },
  {
    id: "bytes",
    secretWord: "BYTES",
    plaintextMessage:
      "THE SOURCE CODE ARCHIVE WILL BE TRANSFERRED AT 2130 THROUGH THE INTERNAL SERVER. EUGENE INTENDS TO SEND THE FILES OUTSIDE THE NETWORK. BLOCK THE CONNECTION FIRST."
  },
  {
    id: "patch",
    secretWord: "PATCH",
    plaintextMessage:
      "EUGENE HAS PREPARED THE SOURCE CODE FOR EXFILTRATION AT 2130. THE INTERNAL SERVER WILL BE USED TO MOVE THE ARCHIVE OUTSIDE THE NETWORK. STOP THE TRANSFER."
  }
]
  };

  const KEYBOARD_ROWS = [
    ["Q", "W", "E", "R", "T", "Y", "U", "I", "O", "P"],
    ["A", "S", "D", "F", "G", "H", "J", "K", "L"],
    ["ENTER", "Z", "X", "C", "V", "B", "N", "M", "BACKSPACE"]
  ];

  const STATUS_PRIORITY = {
    absent: 1,
    present: 2,
    correct: 3
  };

  function bootstrapEncryptionGame() {
    const root = document.querySelector("#encryption-game");
    if (!root || root.dataset.initialized === "true") {
      return;
    }

    root.dataset.initialized = "true";

    if (window.__caughtInTheNetEncryptionController) {
      window.__caughtInTheNetEncryptionController.abort();
    }

    if (window.__caughtInTheNetEncryptionClockInterval) {
      window.clearInterval(window.__caughtInTheNetEncryptionClockInterval);
    }

    const controller = new AbortController();
    window.__caughtInTheNetEncryptionController = controller;
    const { signal } = controller;

    const encryptedMessage = root.querySelector("#encrypted-message");
    const decryptedMessage = root.querySelector("#decrypted-message");
    const board = root.querySelector("#word-board");
    const keyboard = root.querySelector("#virtual-keyboard");
    const feedback = root.querySelector("#game-feedback");
    const attemptCounter = root.querySelector("#attempt-counter");
    const keyLengthStatus = root.querySelector("#key-length-status");
    const decryptStatus = root.querySelector("#decrypt-status");
    const progressPercent = root.querySelector("#decryption-percent");
    const progressBar = root.querySelector("#decryption-progress-bar");
    const terminalProgress = root.querySelector("#terminal-progress-line");
    const terminalStatus = root.querySelector("#terminal-status-line");
    const instructionsModal = root.querySelector("#instructions-modal");
    const instructionsClose = root.querySelector("#instructions-close");
    const beginButton = root.querySelector("#begin-decryption");
    const restartSessionButton = root.querySelector("#restart-session");
    const resultModal = root.querySelector("#result-modal");
    const resultTitle = root.querySelector("#result-title");
    const resultMessage = root.querySelector("#result-message");
    const resultPrimaryButton = root.querySelector("#result-primary-button");
    const resultSecondaryButton = root.querySelector("#result-secondary-button");
    const recoveredMessageModal = root.querySelector("#recovered-message-modal");
    const recoveredMessageText = root.querySelector("#recovered-message-text");
    const continueStoryButton = root.querySelector("#continue-story");
    const clock = root.querySelector("#desktop-clock");

    let activeScenarioIndex = -1;
    let activeScenario = null;
    let secretWord = "";
    let wordLength = 0;
    let currentGuess = "";
    let attemptIndex = 0;
    let bestGreenCount = 0;
    let gameOver = false;
    let inputLocked = false;
    let instructionsOpen = true;
    let resultMode = null;
    let revealOrder = [];
    let validWords = new Set();
    let wordValidationAvailable = false;
    const keyboardStates = new Map();

    selectScenario(false);
    buildBoard();
    buildKeyboard();
    updateDecryption(0);
    updateStatus();
    updateClock();
    prepareWordList();

    window.__caughtInTheNetEncryptionClockInterval = window.setInterval(updateClock, 1000);

    beginButton.addEventListener("click", closeInstructions, { signal });
    restartSessionButton.addEventListener("click", () => resetGame(true), { signal });

    resultPrimaryButton.addEventListener(
      "click",
      () => {
        if (resultMode === "success") {
          openRecoveredMessage();
        } else if (resultMode === "failure") {
          resetGame(true);
        }
      },
      { signal }
    );

    resultSecondaryButton.addEventListener(
      "click",
      () => {
        resultModal.hidden = true;
      },
      { signal }
    );

    continueStoryButton.addEventListener(
      "click",
      () => {
        root.dispatchEvent(
          new CustomEvent("encryption:continue", {
            bubbles: true,
            detail: {
              scenarioId: activeScenario.id,
              secretWord,
              plaintextMessage: activeScenario.plaintextMessage
            }
          })
        );

        window.location.assign(activeScenario.continueUrl || GAME_CONFIG.continueUrl);
      },
      { signal }
    );

    keyboard.addEventListener(
      "click",
      (event) => {
        const keyButton = event.target.closest("[data-key]");
        if (!keyButton) return;
        handleKey(keyButton.dataset.key);
      },
      { signal }
    );

    document.addEventListener(
      "keydown",
      (event) => {
        if (!document.body.contains(root)) return;
        if (instructionsOpen || gameOver || inputLocked) return;

        if (event.key === "Enter") {
          event.preventDefault();
          handleKey("ENTER");
          return;
        }

        if (event.key === "Backspace") {
          event.preventDefault();
          handleKey("BACKSPACE");
          return;
        }

        if (/^[a-zA-Z]$/.test(event.key)) {
          event.preventDefault();
          handleKey(event.key.toUpperCase());
        }
      },
      { signal }
    );

    async function prepareWordList() {
      beginButton.disabled = true;
      instructionsClose.disabled = true;
      beginButton.textContent = "Loading Word List...";

      try {
        const response = await fetch(GAME_CONFIG.validWordsUrl, { cache: "force-cache" });
        if (!response.ok) {
          throw new Error(`Word list request failed with status ${response.status}`);
        }

        const text = await response.text();
        validWords = new Set(
          text
            .split(/\r?\n/)
            .map((word) => word.trim().toUpperCase())
            .filter((word) => /^[A-Z]+$/.test(word) && word.length === wordLength)
        );

        // Always permit the deliberately chosen scenario answers even if the
        // external list is later swapped for a smaller curated dictionary.
        GAME_CONFIG.scenarios.forEach((scenario) => {
          validWords.add(scenario.secretWord.toUpperCase());
        });

        wordValidationAvailable = validWords.size > 0;
      } catch (error) {
        console.error("Unable to load local word list:", error);
        wordValidationAvailable = false;
      } finally {
        beginButton.disabled = false;
        instructionsClose.disabled = false;
        beginButton.textContent = "Begin Decryption";
      }
    }

    function selectScenario(avoidCurrentScenario = true) {
      const scenarios = GAME_CONFIG.scenarios;
      let nextIndex;

      if (scenarios.length === 1) {
        nextIndex = 0;
      } else {
        do {
          nextIndex = Math.floor(Math.random() * scenarios.length);
        } while (avoidCurrentScenario && nextIndex === activeScenarioIndex);
      }

      activeScenarioIndex = nextIndex;
      activeScenario = scenarios[nextIndex];
      secretWord = activeScenario.secretWord.trim().toUpperCase();
      wordLength = secretWord.length;

      revealOrder = buildRevealOrder(`${activeScenario.id}:${activeScenario.plaintextMessage}`);
      encryptedMessage.textContent = columnarTranspositionEncrypt(
        activeScenario.plaintextMessage,
        secretWord
      );

      board.style.setProperty("--word-length", wordLength);
      keyLengthStatus.textContent = `Key length: ${wordLength}`;
    }

    function buildBoard() {
      board.innerHTML = "";

      for (let rowIndex = 0; rowIndex < GAME_CONFIG.maxAttempts; rowIndex += 1) {
        const row = document.createElement("div");
        row.className = "word-row";
        row.setAttribute("role", "row");

        for (let columnIndex = 0; columnIndex < wordLength; columnIndex += 1) {
          const tile = document.createElement("div");
          tile.className = "word-tile";
          tile.setAttribute("role", "gridcell");
          tile.setAttribute(
            "aria-label",
            `Row ${rowIndex + 1}, character ${columnIndex + 1}`
          );
          row.appendChild(tile);
        }

        board.appendChild(row);
      }
    }

    function buildKeyboard() {
      keyboard.innerHTML = "";

      KEYBOARD_ROWS.forEach((rowKeys) => {
        const row = document.createElement("div");
        row.className = "keyboard-row";

        rowKeys.forEach((key) => {
          const button = document.createElement("button");
          button.type = "button";
          button.className = "keyboard-key";
          button.dataset.key = key;

          if (key === "ENTER" || key === "BACKSPACE") {
            button.classList.add("keyboard-key--wide");
          }

          button.textContent = key === "BACKSPACE" ? "⌫" : key;
          button.setAttribute(
            "aria-label",
            key === "BACKSPACE"
              ? "Backspace"
              : key === "ENTER"
                ? "Submit guess"
                : `Letter ${key}`
          );

          row.appendChild(button);
        });

        keyboard.appendChild(row);
      });
    }

    function closeInstructions() {
      if (beginButton.disabled) return;

      instructionsOpen = false;
      instructionsModal.hidden = true;
      feedback.textContent = wordValidationAvailable
        ? `Enter a ${wordLength}-letter cybersecurity key.`
        : `Enter a ${wordLength}-letter key. Word-list validation is unavailable.`;
      decryptStatus.textContent = "STATUS: ACTIVE";
    }

    function handleKey(key) {
      if (instructionsOpen || gameOver || inputLocked) {
        return;
      }

      if (key === "ENTER") {
        submitGuess();
        return;
      }

      if (key === "BACKSPACE") {
        if (currentGuess.length > 0) {
          currentGuess = currentGuess.slice(0, -1);
          renderCurrentGuess();
        }
        return;
      }

      if (/^[A-Z]$/.test(key) && currentGuess.length < wordLength) {
        currentGuess += key;
        renderCurrentGuess();
      }
    }

    function renderCurrentGuess() {
      const row = board.children[attemptIndex];
      if (!row) return;

      const tiles = row.children;

      for (let i = 0; i < wordLength; i += 1) {
        const tile = tiles[i];
        const character = currentGuess[i] || "";

        tile.textContent = character;
        tile.classList.toggle("word-tile--filled", Boolean(character));
      }
    }

    function submitGuess() {
      if (currentGuess.length !== wordLength) {
        feedback.textContent = `The key must contain exactly ${wordLength} letters.`;
        pulseFeedback();
        return;
      }

      const submittedGuess = currentGuess.toUpperCase();

      if (wordValidationAvailable && !validWords.has(submittedGuess)) {
        feedback.textContent = `“${submittedGuess}” is not in the word list.`;
        pulseFeedback();
        return;
      }

      const statuses = evaluateGuess(submittedGuess, secretWord);
      const greenCount = statuses.filter((status) => status === "correct").length;
      const row = board.children[attemptIndex];
      const tiles = [...row.children];

      inputLocked = true;
      feedback.textContent = "Analysing key fragment...";

      tiles.forEach((tile, index) => {
        window.setTimeout(() => {
          tile.classList.remove("word-tile--filled");
          tile.classList.add(`word-tile--${statuses[index]}`, "word-tile--reveal");
        }, index * 150);
      });

      window.setTimeout(() => {
        updateKeyboard(submittedGuess, statuses);

        if (greenCount > bestGreenCount) {
          bestGreenCount = greenCount;
          updateDecryption(bestGreenCount);
        }

        attemptIndex += 1;
        currentGuess = "";

        if (submittedGuess === secretWord) {
          completeGame();
          return;
        }

        if (attemptIndex >= GAME_CONFIG.maxAttempts) {
          failGame();
          return;
        }

        inputLocked = false;
        feedback.textContent =
          greenCount > 0
            ? `${greenCount} character${greenCount === 1 ? "" : "s"} locked in. Keep going.`
            : "No characters locked in on that attempt.";

        updateStatus();
      }, wordLength * 150 + 180);
    }

    /*
     * Proper Wordle-style duplicate handling:
     * 1. Mark exact matches first.
     * 2. Allocate remaining answer letters to yellow matches only once.
     */
    function evaluateGuess(guess, answer) {
      const result = new Array(wordLength).fill("absent");
      const remaining = answer.split("");

      for (let i = 0; i < wordLength; i += 1) {
        if (guess[i] === answer[i]) {
          result[i] = "correct";
          remaining[i] = null;
        }
      }

      for (let i = 0; i < wordLength; i += 1) {
        if (result[i] === "correct") continue;

        const matchingIndex = remaining.findIndex((letter) => letter === guess[i]);

        if (matchingIndex !== -1) {
          result[i] = "present";
          remaining[matchingIndex] = null;
        }
      }

      return result;
    }

    function updateKeyboard(guess, statuses) {
      [...guess].forEach((letter, index) => {
        const nextStatus = statuses[index];
        const currentStatus = keyboardStates.get(letter);

        if (
          !currentStatus ||
          STATUS_PRIORITY[nextStatus] > STATUS_PRIORITY[currentStatus]
        ) {
          keyboardStates.set(letter, nextStatus);
        }
      });

      keyboardStates.forEach((status, letter) => {
        const key = keyboard.querySelector(`[data-key="${letter}"]`);
        if (!key) return;

        key.classList.remove(
          "keyboard-key--correct",
          "keyboard-key--present",
          "keyboard-key--absent"
        );

        key.classList.add(`keyboard-key--${status}`);
      });
    }

    function buildRevealOrder(seedValue) {
      const indexes = [];
      const plaintext = activeScenario.plaintextMessage;

      for (let i = 0; i < plaintext.length; i += 1) {
        if (/[A-Za-z0-9]/.test(plaintext[i])) {
          indexes.push(i);
        }
      }

      // Seeded shuffle: looks random, remains stable for this scenario, and
      // guarantees previously revealed characters never disappear.
      const random = seededRandom(hashString(seedValue));

      for (let i = indexes.length - 1; i > 0; i -= 1) {
        const j = Math.floor(random() * (i + 1));
        [indexes[i], indexes[j]] = [indexes[j], indexes[i]];
      }

      return indexes;
    }

    function updateDecryption(greenCount) {
      const fraction = Math.max(0, Math.min(1, greenCount / wordLength));
      const percentage = Math.round(fraction * 100);
      const revealCount = Math.round(revealOrder.length * fraction);
      const visibleIndexes = new Set(revealOrder.slice(0, revealCount));
      const plaintext = activeScenario.plaintextMessage;

      let output = "";

      for (let i = 0; i < plaintext.length; i += 1) {
        const character = plaintext[i];

        if (!/[A-Za-z0-9]/.test(character)) {
          output += character;
        } else if (visibleIndexes.has(i)) {
          output += character;
        } else {
          output += GAME_CONFIG.hiddenCharacter;
        }
      }

      decryptedMessage.textContent = output;
      progressPercent.textContent = `${percentage}%`;
      progressBar.style.width = `${percentage}%`;
      terminalProgress.textContent = `[PROGRESS: ${percentage}% | KEY MATCH: ${greenCount}/${wordLength}]`;

      if (percentage === 0) {
        terminalStatus.textContent = "Awaiting valid key fragments...";
      } else if (percentage < 100) {
        terminalStatus.textContent =
          `Recovered ${percentage}% of plaintext. Best key match: ${greenCount}/${wordLength}.`;
      } else {
        terminalStatus.textContent = "DECRYPTION COMPLETE. PLAINTEXT RECOVERED.";
      }
    }

    function completeGame() {
      gameOver = true;
      inputLocked = false;
      bestGreenCount = wordLength;
      updateDecryption(wordLength);

      feedback.textContent = "Key accepted. Decryption complete.";
      decryptStatus.textContent = "STATUS: COMPLETE";
      attemptCounter.textContent = `Solved in ${attemptIndex} attempt${attemptIndex === 1 ? "" : "s"}`;

      resultMode = "success";
      resultTitle.textContent = "Message decrypted";
      resultMessage.textContent =
        "Key accepted. The intercepted transmission has been fully recovered.";
      resultPrimaryButton.textContent = "View Decrypted Message";
      resultSecondaryButton.hidden = true;

      window.setTimeout(() => {
        resultModal.hidden = false;
      }, 350);

      root.dispatchEvent(
        new CustomEvent("encryption:complete", {
          bubbles: true,
          detail: {
            attempts: attemptIndex,
            scenarioId: activeScenario.id,
            secretWord,
            plaintextMessage: activeScenario.plaintextMessage
          }
        })
      );
    }

    function failGame() {
      gameOver = true;
      inputLocked = false;
      feedback.textContent = `Key recovery failed. Please try again.`;
      decryptStatus.textContent = "STATUS: FAILED";
      attemptCounter.textContent = `Attempts used: ${GAME_CONFIG.maxAttempts}/${GAME_CONFIG.maxAttempts}`;
      terminalStatus.textContent =
        `KEY RECOVERY FAILED. Best match: ${bestGreenCount}/${wordLength}.`;

      resultMode = "failure";
      resultTitle.textContent = "Decryption failed";
      resultMessage.textContent =
        `The key was not recovered within ${GAME_CONFIG.maxAttempts} attempts. Retry the decryption to continue the investigation.`;
      resultPrimaryButton.textContent = "Retry";
      resultSecondaryButton.textContent = "Return to Terminal";
      resultSecondaryButton.hidden = false;
      resultModal.hidden = false;

      root.dispatchEvent(
        new CustomEvent("encryption:failed", {
          bubbles: true,
          detail: {
            attempts: GAME_CONFIG.maxAttempts,
            bestGreenCount,
            scenarioId: activeScenario.id
          }
        })
      );
    }

    function openRecoveredMessage() {
      resultModal.hidden = true;
      recoveredMessageText.textContent = activeScenario.plaintextMessage;
      recoveredMessageModal.hidden = false;
    }

    function resetGame(selectNewScenario = true) {
      currentGuess = "";
      attemptIndex = 0;
      bestGreenCount = 0;
      gameOver = false;
      inputLocked = false;
      instructionsOpen = false;
      resultMode = null;
      keyboardStates.clear();

      if (selectNewScenario) {
        selectScenario(true);
      } else {
        revealOrder = buildRevealOrder(`${activeScenario.id}:${activeScenario.plaintextMessage}`);
      }

      buildBoard();
      buildKeyboard();
      updateDecryption(0);
      updateStatus();

      feedback.textContent = `Enter a ${wordLength}-letter cybersecurity key.`;
      decryptStatus.textContent = "STATUS: ACTIVE";
      resultModal.hidden = true;
      recoveredMessageModal.hidden = true;
      resultSecondaryButton.hidden = true;
    }

    function updateStatus() {
      attemptCounter.textContent = `Attempt ${Math.min(attemptIndex + 1, GAME_CONFIG.maxAttempts)} of ${GAME_CONFIG.maxAttempts}`;
      keyLengthStatus.textContent = `Key length: ${wordLength}`;
    }

    /*
     * A real (but deliberately simple) classical transposition cipher.
     * The plaintext is written across rows the width of the keyword, then the
     * columns are read out in alphabetical key order. Duplicate key letters are
     * resolved left-to-right by their original position.
     *
     * This is suitable for the minigame's educational metaphor, but it is NOT
     * representative of secure modern encryption such as AES.
     */
    function columnarTranspositionEncrypt(plaintext, keyword) {
      const cleanText = plaintext.toUpperCase().replace(/[^A-Z0-9]/g, "");
      const cleanKey = keyword.toUpperCase().replace(/[^A-Z]/g, "");

      if (!cleanText || !cleanKey) {
        return "[NO ENCRYPTED DATA]";
      }

      const columnCount = cleanKey.length;
      const rowCount = Math.ceil(cleanText.length / columnCount);
      const paddedText = cleanText.padEnd(rowCount * columnCount, "X");

      const columnOrder = [...cleanKey]
        .map((letter, index) => ({ letter, index }))
        .sort((a, b) => {
          const letterCompare = a.letter.localeCompare(b.letter);
          return letterCompare !== 0 ? letterCompare : a.index - b.index;
        });

      let ciphertext = "";

      columnOrder.forEach(({ index: columnIndex }) => {
        for (let rowIndex = 0; rowIndex < rowCount; rowIndex += 1) {
          ciphertext += paddedText[rowIndex * columnCount + columnIndex];
        }
      });

      return formatCiphertext(ciphertext);
    }

    function formatCiphertext(ciphertext) {
      const groupSize = 5;
      const groupsPerLine = 6;
      const groups = [];

      for (let i = 0; i < ciphertext.length; i += groupSize) {
        groups.push(ciphertext.slice(i, i + groupSize));
      }

      const lines = [];
      for (let i = 0; i < groups.length; i += groupsPerLine) {
        lines.push(groups.slice(i, i + groupsPerLine).join(" "));
      }

      return lines.join("\n");
    }

    function pulseFeedback() {
      feedback.animate(
        [
          { transform: "translateX(0)" },
          { transform: "translateX(-5px)" },
          { transform: "translateX(5px)" },
          { transform: "translateX(0)" }
        ],
        { duration: 180 }
      );
    }

    function updateClock() {
      const now = new Date();
      clock.textContent = now.toLocaleTimeString([], {
        hour: "2-digit",
        minute: "2-digit"
      });
    }

    function hashString(value) {
      let hash = 2166136261;

      for (let i = 0; i < value.length; i += 1) {
        hash ^= value.charCodeAt(i);
        hash = Math.imul(hash, 16777619);
      }

      return hash >>> 0;
    }

    function seededRandom(seed) {
      let state = seed || 1;

      return () => {
        state += 0x6d2b79f5;
        let value = state;
        value = Math.imul(value ^ (value >>> 15), value | 1);
        value ^= value + Math.imul(value ^ (value >>> 7), value | 61);
        return ((value ^ (value >>> 14)) >>> 0) / 4294967296;
      };
    }
  }

  if (!window.__caughtInTheNetEncryptionBootstrapBound) {
    window.__caughtInTheNetEncryptionBootstrapBound = true;

    if (document.readyState === "loading") {
      document.addEventListener("DOMContentLoaded", bootstrapEncryptionGame);
    } else {
      bootstrapEncryptionGame();
    }

    document.addEventListener("htmx:load", bootstrapEncryptionGame);
  } else {
    bootstrapEncryptionGame();
  }
})();
