(() => {
  const clueDefinitions = [
    {
      id: "mia",
      notebook: "Mia",
      password: "Mia",
      image: "/public/images/password/drawing.png",
      imageAlt: "Close-up of a child's drawing found in Eugene's office",
      targetClass: "password-game__zoom-target--mia",
      feedback:
        'Correct! The drawing is signed "Mia". Eugene\'s daughter\'s name is personal information that could easily become part of a password.'
    },
    {
      id: "april9",
      notebook: "April 9",
      password: "April9",
      image: "/public/images/password/calendar.png",
      imageAlt: "Close-up of Eugene's April calendar",
      targetClass: "password-game__zoom-target--april9",
      feedback:
        "Correct! April 9 is marked as Mia's birthday. Birthdays and important dates are commonly reused in passwords."
    },
    {
      id: "birthday",
      notebook: "04/08/1990",
      password: "04081990",
      image: "/public/images/password/id.png",
      imageAlt: "Close-up of Eugene's employee identification card",
      targetClass: "password-game__zoom-target--birthday",
      feedback:
        "Correct! Eugene's ID reveals his date of birth: 04/08/1990. Dates of birth are predictable personal details and risky password material."
    },
    {
      id: "jersey",
      notebook: "20",
      password: "20",
      image: "/public/images/password/jersey.png",
      imageAlt: "Close-up of Eugene's framed sports jersey",
      targetClass: "password-game__zoom-target--jersey",
      feedback:
        "Correct! Eugene's framed jersey shows the number 20. Favourite or meaningful numbers are often included in passwords."
    },
    {
      id: "rufus",
      notebook: "Rufus",
      password: "Rufus",
      image: "/public/images/password/rufus.png",
      imageAlt: "Close-up of Eugene's dog photograph",
      targetClass: "password-game__zoom-target--rufus",
      feedback:
        'Correct! Eugene\'s dog is named "Rufus". Pet names are a common and easy-to-guess password choice.'
    }
  ];

  const clueById = new Map(clueDefinitions.map((clue) => [clue.id, clue]));

  function shuffle(values) {
    const result = [...values];
    for (let index = result.length - 1; index > 0; index -= 1) {
      const swapIndex = Math.floor(Math.random() * (index + 1));
      [result[index], result[swapIndex]] = [result[swapIndex], result[index]];
    }
    return result;
  }

  function initPasswordGame() {
    const root = document.getElementById("password-game");
    if (!root || root.dataset.initialized === "true") {
      return;
    }
    root.dataset.initialized = "true";

    const officePhase = root.querySelector("#password-phase-office");
    const loginPhase = root.querySelector("#password-phase-login");
    const officeHotspots = [...root.querySelectorAll(".password-game__office-hotspot")];

    const zoom = root.querySelector("#password-zoom");
    const zoomBackdrop = root.querySelector("#password-zoom-backdrop");
    const zoomFrame = root.querySelector("#password-zoom-frame");
    const zoomImage = root.querySelector("#password-zoom-image");
    const zoomTarget = root.querySelector("#password-zoom-target");

    const notebookIcon = root.querySelector("#password-notebook-icon");
    const notebookCount = root.querySelector("#password-notebook-count");
    const notebookOverlay = root.querySelector("#password-notebook-overlay");
    const notebookBackdrop = root.querySelector("#password-notebook-backdrop");
    const notebookLines = root.querySelector("#password-notebook-lines");

    const phase2Clues = root.querySelector("#password-phase2-clues");
    const passwordSlots = [...root.querySelectorAll(".password-game__password-slot")];
    const loginButton = root.querySelector("#password-login-button");
    const loginStatus = root.querySelector("#password-login-status");

    const modal = root.querySelector("#password-game-modal");
    const modalEyebrow = root.querySelector("#password-modal-eyebrow");
    const modalTitle = root.querySelector("#password-modal-title");
    const modalContent = root.querySelector("#password-modal-content");
    const modalPrimary = root.querySelector("#password-modal-primary");

    let activeZoomClueId = null;
    let discoveredOrder = [];
    let typingTimer = null;
    let notebookWasAutoOpened = false;
    let pendingPhaseOneCompletion = false;
    let slots = [null, null, null];

    // Each run uses three different clues in a new order.
    const passwordSolution = shuffle(clueDefinitions.map((clue) => clue.id)).slice(0, 3);

    function showModal({ eyebrow = "PASSWORD INVESTIGATION", title, html, button = "Continue", onPrimary }) {
      modalEyebrow.textContent = eyebrow;
      modalTitle.textContent = title;
      modalContent.innerHTML = html;
      modalPrimary.textContent = button;
      modalPrimary.onclick = () => {
        if (typeof onPrimary === "function") {
          onPrimary();
        }
      };
      modal.hidden = false;
    }

    function hideModal() {
      modal.hidden = true;
    }

    function showIntro() {
      showModal({
        eyebrow: "PASSWORDS & SOCIAL ENGINEERING",
        title: "Why personal information matters",
        html: `
      <p>Passwords are one of the most common ways of protecting accounts, but a password is only useful if it is difficult for someone else to predict. People often build passwords from things that are easy to remember, such as names, birthdays, pets, interests and meaningful numbers.</p>
      <p>Attackers may collect personal information from social media, workplaces, conversations or other sources and use it to guess passwords or make other attacks more convincing. This is one way social engineering can help an attacker get around security without directly attacking the technology.</p>
      <p>In this activity, you will investigate Eugene's office and see how ordinary personal details can be combined into a predictable password.</p>
    `,
        button: "Continue",
        onPrimary: showPhaseOneInstructions
      });
    }

    function showPhaseOneInstructions() {
      showModal({
        eyebrow: "PHASE 1 - OFFICE INVESTIGATION",
        title: "Find Eugene's password clues",
        html: `
          <p>Search Eugene's office for objects that may reveal personal information.</p>
          <ol>
            <li>Click an object in the room to inspect it more closely.</li>
            <li>In the close-up, click the specific detail that could be used in a password.</li>
            <li>Each successful clue is recorded in your Case Notes.</li>
          </ol>
          <p>Find all <strong>five</strong> visible clues before attempting the login.</p>
        `,
        button: "Begin Investigation",
        onPrimary: hideModal
      });
    }

    function openZoom(clueId) {
      const clue = clueById.get(clueId);
      if (!clue) {
        return;
      }

      activeZoomClueId = clueId;
      zoomImage.src = clue.image;
      zoomImage.alt = clue.imageAlt;
      zoomTarget.className = `password-game__zoom-target ${clue.targetClass}`;
      zoomTarget.setAttribute("aria-label", `Select the ${clue.notebook} clue`);
      zoom.hidden = false;
    }

    function closeZoom() {
      zoom.hidden = true;
      activeZoomClueId = null;
      zoomImage.src = "";
      zoomTarget.className = "password-game__zoom-target";
    }

    function renderNotebook({ animateClueId = null } = {}) {
      if (typingTimer) {
        window.clearInterval(typingTimer);
        typingTimer = null;
      }

      notebookLines.innerHTML = "";

      discoveredOrder.forEach((clueId, index) => {
        const clue = clueById.get(clueId);
        const line = document.createElement("p");
        line.className = "password-game__notebook-line";
        line.style.setProperty("--line-index", index);
        line.dataset.clueId = clueId;
        line.textContent = clueId === animateClueId ? "" : clue.notebook;
        notebookLines.appendChild(line);
      });

      if (!animateClueId) {
        return;
      }

      const clue = clueById.get(animateClueId);
      const targetLine = notebookLines.querySelector(`[data-clue-id="${animateClueId}"]`);
      if (!clue || !targetLine) {
        return;
      }

      let characterIndex = 0;
      typingTimer = window.setInterval(() => {
        characterIndex += 1;
        targetLine.textContent = clue.notebook.slice(0, characterIndex);

        if (characterIndex >= clue.notebook.length) {
          window.clearInterval(typingTimer);
          typingTimer = null;
        }
      }, 55);
    }

    function openNotebook({ animateClueId = null, automatic = false } = {}) {
      notebookWasAutoOpened = automatic;
      renderNotebook({ animateClueId });
      notebookOverlay.hidden = false;
    }

    function closeNotebook() {
      if (typingTimer) {
        window.clearInterval(typingTimer);
        typingTimer = null;
        renderNotebook();
      }

      notebookOverlay.hidden = true;
      notebookWasAutoOpened = false;

      if (pendingPhaseOneCompletion) {
        pendingPhaseOneCompletion = false;
        showPhaseOneComplete();
      }
    }

    function showClueFeedback(clueId, newlyDiscovered) {
      const clue = clueById.get(clueId);
      if (!clue) {
        return;
      }

      if (!newlyDiscovered) {
        showModal({
          eyebrow: "CASE NOTES",
          title: "Clue already recorded",
          html: `<p><strong>${clue.notebook}</strong> is already in your Case Notes.</p>`,
          button: "Return to Office",
          onPrimary: hideModal
        });
        return;
      }

      showModal({
        eyebrow: "CLUE FOUND",
        title: "Correct",
        html: `<p>${clue.feedback}</p>`,
        button: "Add to Case Notes",
        onPrimary: () => {
          hideModal();
          pendingPhaseOneCompletion = discoveredOrder.length === clueDefinitions.length;
          openNotebook({ animateClueId: clueId, automatic: true });
        }
      });
    }

    function collectActiveZoomClue() {
      if (!activeZoomClueId) {
        return;
      }

      const clueId = activeZoomClueId;
      const newlyDiscovered = !discoveredOrder.includes(clueId);

      if (newlyDiscovered) {
        discoveredOrder.push(clueId);
        notebookCount.textContent = `${discoveredOrder.length} / ${clueDefinitions.length}`;

        const officeHotspot = root.querySelector(`.password-game__office-hotspot[data-clue-id="${clueId}"]`);
        if (officeHotspot) {
          officeHotspot.classList.add("password-game__office-hotspot--found");
          officeHotspot.setAttribute("aria-label", `${officeHotspot.getAttribute("aria-label")} - clue found`);
        }
      }

      closeZoom();
      showClueFeedback(clueId, newlyDiscovered);
    }

    function showPhaseOneComplete() {
      showModal({
        eyebrow: "PHASE 1 COMPLETE",
        title: "All five clues found",
        html: `
          <p>You found all the visible personal information in Eugene's office that he might have used to create his password.</p>
          <p>Now use your Case Notes to reconstruct a three-part password.</p>
          <ol>
            <li>Click a clue to place it in the next empty password box.</li>
            <li><strong>Green</strong> means the clue is correct in that position. <strong>Red</strong> means it is not.</li>
            <li>Click a selected clue again to remove it and try another combination.</li>
          </ol>
        `,
        button: "Go to Login",
        onPrimary: beginPhaseTwo
      });
    }

    function beginPhaseTwo() {
      hideModal();
      officePhase.hidden = true;
      loginPhase.hidden = false;
      renderPhaseTwoClues();
      renderSlots();
    }

    function renderPhaseTwoClues() {
      phase2Clues.innerHTML = "";

      discoveredOrder.forEach((clueId, index) => {
        const clue = clueById.get(clueId);
        const button = document.createElement("button");
        button.type = "button";
        button.className = "password-game__phase2-clue";
        button.dataset.clueId = clueId;
        button.style.setProperty("--line-index", index);
        button.textContent = clue.notebook;
        button.setAttribute("aria-label", `Use ${clue.notebook} in the password`);
        button.addEventListener("click", () => toggleClueInPassword(clueId));
        phase2Clues.appendChild(button);
      });
    }

    function toggleClueInPassword(clueId) {
      const existingIndex = slots.indexOf(clueId);
      if (existingIndex !== -1) {
        slots[existingIndex] = null;
        renderSlots();
        return;
      }

      const emptyIndex = slots.indexOf(null);
      if (emptyIndex === -1) {
        loginStatus.textContent = "All three boxes are full. Remove a clue before adding another.";
        return;
      }

      slots[emptyIndex] = clueId;
      renderSlots();
    }

    function removeSlot(index) {
      if (!slots[index]) {
        return;
      }
      slots[index] = null;
      renderSlots();
    }

    function renderSlots() {
      passwordSlots.forEach((slotElement, index) => {
        const clueId = slots[index];
        slotElement.classList.remove(
          "password-game__password-slot--correct",
          "password-game__password-slot--incorrect"
        );

        if (!clueId) {
          slotElement.textContent = "";
          slotElement.disabled = true;
          return;
        }

        const clue = clueById.get(clueId);
        slotElement.textContent = clue.password;
        slotElement.disabled = false;
        slotElement.classList.add(
          passwordSolution[index] === clueId
            ? "password-game__password-slot--correct"
            : "password-game__password-slot--incorrect"
        );
      });

      const selectedIds = new Set(slots.filter(Boolean));
      root.querySelectorAll(".password-game__phase2-clue").forEach((button) => {
        button.classList.toggle(
          "password-game__phase2-clue--selected",
          selectedIds.has(button.dataset.clueId)
        );
      });

      const allCorrect = slots.every(
        (clueId, index) => clueId !== null && clueId === passwordSolution[index]
      );

      loginButton.disabled = !allCorrect;

      if (allCorrect) {
        loginStatus.textContent = "Password reconstructed. Login is ready.";
        loginStatus.classList.add("password-game__login-status--ready");
      } else {
        const filled = slots.filter(Boolean).length;
        loginStatus.textContent =
          filled === 0
            ? "Match all three password parts to enable login."
            : `${filled} of 3 password boxes filled. Green is correct; red needs another clue.`;
        loginStatus.classList.remove("password-game__login-status--ready");
      }
    }

    function completeGame() {
      const passwordText = passwordSolution
        .map((clueId) => clueById.get(clueId).password)
        .join("");

      showModal({
        eyebrow: "ACCESS GRANTED",
        title: "Password cracked",
        html: `
          <p>You reconstructed Eugene's password:</p>
          <p class="password-game__recovered-password">${passwordText}</p>
          <p>Personal details made the password memorable for Eugene, but they also made it much easier for someone else to guess.</p>
        `,
        button: "Continue Investigation",
        onPrimary: () => {
          window.location.href = "/whwl/password";
        }
      });
    }

    officeHotspots.forEach((hotspot) => {
      hotspot.addEventListener("click", () => openZoom(hotspot.dataset.clueId));
    });

    zoomBackdrop.addEventListener("click", closeZoom);
    zoomTarget.addEventListener("click", collectActiveZoomClue);

    notebookIcon.addEventListener("click", () => openNotebook());
    notebookBackdrop.addEventListener("click", closeNotebook);

    passwordSlots.forEach((slotElement, index) => {
      slotElement.addEventListener("click", () => removeSlot(index));
    });

    loginButton.addEventListener("click", completeGame);

    document.addEventListener("keydown", (event) => {
      if (event.key !== "Escape") {
        return;
      }

      if (!zoom.hidden) {
        closeZoom();
        return;
      }

      if (!notebookOverlay.hidden) {
        closeNotebook();
      }
    });

    // Keep the zoom target aligned with the image rather than accidentally
    // inheriting browser button dimensions.
    zoomFrame.addEventListener("click", (event) => {
      event.stopPropagation();
    });

    showIntro();
  }

  document.addEventListener("DOMContentLoaded", initPasswordGame);
  document.addEventListener("htmx:afterSwap", initPasswordGame);
})();
