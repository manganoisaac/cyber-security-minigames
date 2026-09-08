document.addEventListener("DOMContentLoaded", () => {
    const evidenceFiles = {
        E03: {
            title: "Internal Project Announcement",
            content: `METACORTEX INDUSTRIES
INTERNAL STAFF ANNOUNCEMENT

Date: 6 August 2026
Time: 09:10
Subject: Project HELIOS Development Team

Metacortex Industries is pleased to confirm the development team for Project HELIOS, our latest internal research and development initiative.

The project team will include:
- Angela Bennett
- Marcus Lee
- Sophie Carter
- Daniel Ruiz

The team will begin development activities this month, with further project information to be circulated internally as work progresses.

Please join us in congratulating the selected team members.

Research & Development Management
Metacortex Industries`
        },

        E04: {
            title: "Browser History - Eugene's Workstation",
            content: `BROWSER HISTORY EXPORT

Device: MCX-WS17
Assigned User: Eugene
Date: 19 August 2026

19:18:42 - Search: "Angela Bennett Metacortex"
19:21:07 - Public/professional page
19:27:31 - Social page
19:34:18 - Event/attendee page
19:39:52 - Social post
19:46:03 - Metacortex public staff page`
        },

        E08: {
            title: "Email Message - Angela's Inbox",
            content: `EMAIL MESSAGE

To: Angela Bennett
Sent: 14 October 2026 - 15:36
Subject: Account Verification Required

Hi Angela,

We need you to confirm your details using the link below.

[LINK TO ACCOUNT VERIFICATION PAGE]

Please complete the verification as soon as possible.

Metacortex Support`
        },

        E09: {
            title: "Browser Record - Angela's Workstation",
            content: `BROWSER ACTIVITY RECORD

Device: MCX-WS04
Assigned User: Angela Bennett
User Account: angela.b
Date: 14 October 2026

15:41:08 - Navigated to account verification page

Page Title: Metacortex Account Verification
Previous Page: Angela's webmail session
Session Duration: 00:02:17`
        },

        E11: {
            title: "Workstation Session Log - Angela",
            content: `WORKSTATION SESSION LOG

Device: MCX-WS04
Assigned User: Angela Bennett
Date: 14 October 2026

17:12:46 - Project repository accessed
17:16:02 - Email session closed
17:18:31 - User angela.b logged off`
        },

        E12: {
            title: "Recovered Communication - Eugene's Computer",
            content: `RECOVERED COMMUNICATION

Date: 14 October 2026
Time: 20:14
Remote Contact: G. Winston - NURV Technologies

G. Winston:
"Everything is set for 2130. The connection to our server will remain open until the transfer is complete. Send the source-code archive through the internal server as planned, and don't leave anything behind that can be traced back to us."

Local User:
"Understood."

G. Winston:
"Once delivery is confirmed, the agreed payment will be released and the position at NURV is yours."`
        },

        E13: {
            title: "Authentication Log - Internal Network",
            content: `AUTHENTICATION LOG
Date: 14 October 2026

18:07:14 - LOGIN SUCCESS
Account: angela.b
Source Device: MCX-WS17
Target: DEV-SRV02

18:09:31 - ACCESS GRANTED
Account: angela.b
Resource: Development Server`
        },

        E14: {
            title: "Network Security Event",
            content: `NETWORK SECURITY EVENT
Date: 14 October 2026

21:29:58 - OUTBOUND CONNECTION OPENED
Source: DEV-SRV02
Destination: NURV External Server

21:30:02 - FILE TRANSFER STARTED
File: HELIOS_SOURCE_ARCHIVE.zip

21:30:04 - TRANSFER BLOCKED
Action: Network defence controls deployed
Result: Outbound connection terminated
Status: No data transfer completed`
        }
    };

    // ---------------------------------------------------------
    // EVIDENCE FILE VIEWER
    // ---------------------------------------------------------

    const viewer = document.getElementById("evidence-viewer");
    const viewerTitle = document.getElementById("evidence-viewer-title");
    const viewerContent = document.getElementById("evidence-viewer-content");
    const closeViewerButton = document.getElementById("close-evidence-viewer");
    const evidenceButtons = document.querySelectorAll(".timeline-evidence-file");

    const timelineInstructionsModal = document.getElementById(
        "timeline-instructions-modal"
    );
    const beginTimelineButton = document.getElementById(
        "begin-timeline-investigation"
    );
    const timelineSuccessModal = document.getElementById(
        "timeline-success-modal"
    );

    function openEvidenceFile(evidenceId) {
        const evidence = evidenceFiles[evidenceId];

        if (!evidence) {
            return;
        }

        viewerTitle.textContent = evidence.title;
        viewerContent.textContent = evidence.content;
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