document.addEventListener("DOMContentLoaded", () => {
    const evidenceFiles = {
        E03: {
            title: "Project HELIOS Development Team Announcement",
            content: `PROJECT HELIOS DEVELOPMENT TEAM ANNOUNCEMENT

Date: 6 August 2026
Time: 09:10

Selected Team Members:
- Angela Bennett
- Marcus Lee
- Sophie Carter
- Daniel Ruiz`
        },

        E04: {
            title: "Angela Research - Eugene's Workstation",
            content: `ANGELA RESEARCH - EUGENE'S WORKSTATION

Device: MCX-WS17 - Eugene's Workstation
User: Eugene
Date: 19 August 2026

19:18:42: Search "Angela Bennett Metacortex"
19:21:07: Viewed public/professional profile
19:27:31: Viewed social media page
19:34:18: Viewed event/attendee page
19:39:52: Viewed social media post
19:46:03: Viewed Metacortex public staff page`
        },

        E08: {
            title: "Phishing Email Received - Angela's Inbox",
            content: `PHISHING EMAIL RECEIVED - ANGELA'S INBOX

To: Angela Bennett
Sent: 14 October 2026 - 15:36
Subject: Project Helios Access Review
From: David Walker <david.walker@metac0rtex.com.au>

Message:
Before you head off next week, please confirm your Project Helios repository access by signing in before 4:00pm today.

Displayed Link: https://portal.metacortex.com.au/helios
Actual Link: https://metacortex-access-review.com/helios`
        },

        E09: {
            title: "Phishing Link Opened - Angela's Workstation",
            content: `PHISHING LINK OPENED - ANGELA'S WORKSTATION

Device: MCX-WS04 - Angela's Workstation
User Account: angela.b
Date: 14 October 2026

15:41:08: Opened link from suspicious email
Destination: https://metacortex-access-review.com/helios
Previous Page: Angela's Webmail
Time on Page: 2 minutes 17 seconds`
        },

        E11: {
            title: "Angela Logged Off - Angela's Workstation",
            content: `ANGELA BENNETT - WORKSTATION SESSION

Date: 14 October 2026
User Account: angela.b
Source Device: MCX-WS04 - Angela's Workstation

17:12:46: Project HELIOS repository accessed
17:16:02: Email closed
17:18:31: angela.b logged off`
        },

        E12: {
            title: "NURV Source-Code Transfer Plan - Eugene's Computer",
            content: `RECOVERED COMMUNICATION

Date: 14 October 2026
Time: 20:14
Contact: G. Winston - NURV Technologies

G. Winston:
Everything is set for 21:30. The connection will remain open until the transfer is complete. Send the source-code archive through the internal server and don't leave any traces. Payment and your position at NURV will be arranged after delivery.

Local User:
Understood.`
        },

        E13: {
            title: "angela.b Login - Eugene's Workstation",
            content: `AUTHENTICATION LOG - INTERNAL NETWORK

Date: 14 October 2026
User Account: angela.b
Source Device: MCX-WS17 - Eugene's Workstation
Destination: Development Server

18:07:14: Login successful
18:09:31: Development Server access granted`
        },

        E14: {
            title: "HELIOS Transfer to NURV Blocked",
            content: `NETWORK EVENT - DEVELOPMENT SERVER

Date: 14 October 2026
Source: Development Server
Destination: NURV Technologies

21:29:58: Connection opened to NURV Technologies
21:30:02: HELIOS_SOURCE_ARCHIVE.zip transfer started
21:30:04: Transfer blocked by network defence controls

Result: No data successfully transferred`
        },
    };

    // ---------------------------------------------------------
    // EVIDENCE FILE VIEWER
    // ---------------------------------------------------------

    const investigationClock = document.getElementById("investigation-clock");
    const viewer = document.getElementById("evidence-viewer");
    const viewerTitle = document.getElementById("evidence-viewer-title");
    const viewerContent = document.getElementById("evidence-viewer-content");
    const closeViewerButton = document.getElementById("close-evidence-viewer");
    const evidenceButtons = document.querySelectorAll(".timeline-evidence-file");

    const timelineIntroModal = document.getElementById(
        "timeline-intro-modal"
    );
    const continueToTimelineInstructionsButton = document.getElementById(
        "continue-to-timeline-instructions"
    );
    const timelineInstructionsModal = document.getElementById(
        "timeline-instructions-modal"
    );
    const beginTimelineButton = document.getElementById(
        "begin-timeline-investigation"
    );
    const timelineSuccessModal = document.getElementById(
        "timeline-success-modal"
    );

    function renderEvidenceContent(element, content) {
        element.textContent = "";

        const lines = content.split("\n");

        lines.forEach((line, index) => {
            const labelMatch = line.match(/^([A-Za-z][A-Za-z /&'-]{1,40}):\s*(.*)$/);

            if (labelMatch) {
                const label = document.createElement("strong");
                label.textContent = `${labelMatch[1]}:`;

                element.appendChild(label);

                if (labelMatch[2]) {
                    element.appendChild(
                        document.createTextNode(` ${labelMatch[2]}`)
                    );
                }
            } else {
                element.appendChild(document.createTextNode(line));
            }

            if (index < lines.length - 1) {
                element.appendChild(document.createTextNode("\n"));
            }
        });
    }

    function openEvidenceFile(evidenceId) {
        const evidence = evidenceFiles[evidenceId];

        if (!evidence) {
            return;
        }

        viewerTitle.textContent = evidence.title;
        renderEvidenceContent(viewerContent, evidence.content);
        viewer.hidden = false;
    }

    function closeEvidenceFile() {
        viewer.hidden = true;
    }

    evidenceButtons.forEach((button) => {
        button.addEventListener("dblclick", () => {
            openEvidenceFile(button.dataset.evidenceId);
        });

        button.addEventListener("keydown", (event) => {
            if (event.key === "Enter") {
                openEvidenceFile(button.dataset.evidenceId);
            }
        });
    });

    closeViewerButton.addEventListener("click", closeEvidenceFile);

    continueToTimelineInstructionsButton.addEventListener("click", () => {
        timelineIntroModal.hidden = true;
        timelineInstructionsModal.hidden = false;
    });

    beginTimelineButton.addEventListener("click", () => {
        timelineInstructionsModal.hidden = true;
    });

    const beginEvidenceAnalysisButton = document.getElementById(
        "begin-evidence-analysis"
    );

    beginEvidenceAnalysisButton.addEventListener("click", () => {
        window.location.href = "/digital-investigation/analysis";
    });

    document.addEventListener("keydown", (event) => {
        if (event.key === "Escape" && !viewer.hidden) {
            closeEvidenceFile();
        }
    });

    // ---------------------------------------------------------
    // TIMELINE DRAG AND DROP
    // ---------------------------------------------------------

    const evidencePool = document.getElementById("timeline-evidence-pool");
    const timelineSlots = document.querySelectorAll(".timeline-slot");
    const checkTimelineButton = document.getElementById("check-timeline");
    const timelineFeedback = document.getElementById("timeline-feedback");

    const correctTimelineOrder = [
        "E03",
        "E04",
        "E08",
        "E09",
        "E11",
        "E13",
        "E12",
        "E14"
    ];

    let draggedFile = null;

    function shuffleEvidenceFiles() {
        const files = Array.from(
            evidencePool.querySelectorAll(".timeline-evidence-file")
        );

        for (let i = files.length - 1; i > 0; i--) {
            const randomIndex = Math.floor(Math.random() * (i + 1));
            [files[i], files[randomIndex]] = [files[randomIndex], files[i]];
        }

        files.forEach((file) => {
            evidencePool.appendChild(file);
        });
    }

    shuffleEvidenceFiles();

    function clearSlotFeedback(slot) {
        if (!slot) {
            return;
        }

        slot.classList.remove(
            "timeline-slot--correct",
            "timeline-slot--incorrect",
            "timeline-slot--dragover"
        );
    }

    evidenceButtons.forEach((button) => {
        button.addEventListener("dragstart", (event) => {
            draggedFile = button;
            button.classList.add("timeline-evidence-file--dragging");

            const sourceSlot = button.closest(".timeline-slot");
            clearSlotFeedback(sourceSlot);

            const dragPreview = button.cloneNode(true);
            dragPreview.classList.add("timeline-drag-preview");

            document.body.appendChild(dragPreview);

            event.dataTransfer.setDragImage(dragPreview, 28, 28);

            setTimeout(() => {
                dragPreview.remove();
            }, 0);
        });

        button.addEventListener("dragend", () => {
            button.classList.remove("timeline-evidence-file--dragging");
            draggedFile = null;

            timelineSlots.forEach((slot) => {
                slot.classList.remove("timeline-slot--dragover");
            });
        });
    });

    timelineSlots.forEach((slot) => {
        slot.addEventListener("dragover", (event) => {
            event.preventDefault();

            if (draggedFile) {
                slot.classList.add("timeline-slot--dragover");
            }
        });

        slot.addEventListener("dragleave", () => {
            slot.classList.remove("timeline-slot--dragover");
        });

        slot.addEventListener("drop", (event) => {
            event.preventDefault();
            slot.classList.remove("timeline-slot--dragover");

            if (!draggedFile) {
                return;
            }

            const sourceSlot = draggedFile.closest(".timeline-slot");
            const existingFile = slot.querySelector(".timeline-evidence-file");

            clearSlotFeedback(sourceSlot);
            clearSlotFeedback(slot);

            if (existingFile && existingFile !== draggedFile) {
                if (sourceSlot) {
                    sourceSlot.appendChild(existingFile);
                } else {
                    evidencePool.appendChild(existingFile);
                }
            }

            slot.appendChild(draggedFile);
        });
    });

    evidencePool.addEventListener("dragover", (event) => {
        event.preventDefault();
    });

    evidencePool.addEventListener("drop", (event) => {
        event.preventDefault();

        if (!draggedFile) {
            return;
        }

        const sourceSlot = draggedFile.closest(".timeline-slot");
        clearSlotFeedback(sourceSlot);

        evidencePool.appendChild(draggedFile);
    });

    function updateInvestigationClock() {
        if (!investigationClock) {
            return;
        }

        const now = new Date();
        investigationClock.textContent = now.toLocaleTimeString([], {
            hour: "2-digit",
            minute: "2-digit"
        });
    }

    updateInvestigationClock();

    if (window.__caughtInTheNetInvestigationClockInterval) {
        window.clearInterval(window.__caughtInTheNetInvestigationClockInterval);
    }

    window.__caughtInTheNetInvestigationClockInterval = window.setInterval(
        updateInvestigationClock,
        1000
    );

    // ---------------------------------------------------------
    // CHECK TIMELINE
    // ---------------------------------------------------------

    checkTimelineButton.addEventListener("click", () => {
        const filledSlots = Array.from(timelineSlots).filter((slot) =>
            slot.querySelector(".timeline-evidence-file")
        );

        if (filledSlots.length < correctTimelineOrder.length) {
            timelineFeedback.textContent =
                "Place all eight evidence files on the timeline before checking.";
            return;
        }

        let correctCount = 0;

        timelineSlots.forEach((slot, index) => {
            slot.classList.remove(
                "timeline-slot--correct",
                "timeline-slot--incorrect"
            );

            const file = slot.querySelector(".timeline-evidence-file");

            if (file && file.dataset.evidenceId === correctTimelineOrder[index]) {
                slot.classList.add("timeline-slot--correct");
                correctCount++;
            } else {
                slot.classList.add("timeline-slot--incorrect");

                setTimeout(() => {
                    slot.classList.remove("timeline-slot--incorrect");
                }, 900);
            }
        });

        if (correctCount === correctTimelineOrder.length) {
            timelineFeedback.textContent =
                "Incident timeline successfully reconstructed.";

            setTimeout(() => {
                timelineSuccessModal.hidden = false;
            }, 350);
        } else {
            timelineFeedback.textContent =
                `${correctCount} of 8 events are in the correct position. Keep the green events in place and try again.`;
        }
    });

});