document.addEventListener("DOMContentLoaded", () => {

    const analysisClock = document.getElementById("analysis-clock");

    function updateAnalysisClock() {
        if (!analysisClock) {
            return;
        }

        const now = new Date();
        analysisClock.textContent = now.toLocaleTimeString([], {
            hour: "2-digit",
            minute: "2-digit"
        });
    }

    updateAnalysisClock();

    if (window.__caughtInTheNetAnalysisClockInterval) {
        window.clearInterval(window.__caughtInTheNetAnalysisClockInterval);
    }

    window.__caughtInTheNetAnalysisClockInterval = window.setInterval(
        updateAnalysisClock,
        1000
    );

    // ---------------------------------------------------------
    // EVIDENCE DATA
    // ---------------------------------------------------------

    const evidenceFiles = {
        E01: {
            title: "Senior Developer Promotion Rejected - Eugene",
            content: `PROMOTION REVIEW OUTCOME

Date: 12 June 2026
Employee: Eugene
Position Applied For: Senior Developer
Outcome: Application unsuccessful
Reason: Strong technical performance, but further team leadership and collaborative project experience required`
        },

        E02: {
            title: "Project HELIOS Application Rejected - Eugene",
            content: `PROJECT HELIOS APPLICATION OUTCOME

Date: 31 July 2026
Employee: Eugene
Project: Project HELIOS
Outcome: Not selected for the development team`
        },

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

        E05: {
            title: "Angela's Public Information",
            content: `PUBLIC INFORMATION - ANGELA BENNETT

Name: Angela Bennett
Employer: MetaCortex Industries
Location: South Brisbane
Current Project: Project HELIOS
Work Contact: David Walker
Travel Plans: Away from work 14-28 October for a trip to Japan

Source: Public social media and professional information`
        },

        E06: {
            title: "Encrypted File transfer.enc - Eugene's Workstation",
            content: `RECOVERED FILE RECORD

Device: MCX-WS17 - Eugene's Workstation
User: Eugene
File: transfer.enc
Location: C:\\Users\\eugene\\Documents\\transfer.enc
File Type: Encrypted data
Created: 14 October 2026 - 20:15:06
Modified: 14 October 2026 - 20:15:06`
        },

        E07: {
            title: "NURV Senior Developer Research - Eugene's Workstation",
            content: `BROWSER HISTORY - EUGENE'S WORKSTATION

Device: MCX-WS17 - Eugene's Workstation
User: Eugene

3 September 2026: Viewed NURV Technologies - About
3 September 2026: Viewed NURV Research & Development
5 September 2026: Viewed NURV Careers
5 September 2026: Viewed Software Development
8 September 2026: Viewed Senior Developer position`
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

        E10: {
            title: "Project HELIOS Source-Code Testing - Angela's Inbox",
            content: `PROJECT HELIOS EMAIL

Date: 9 October 2026
Time: 10:24
From: Dr Marcus Lee
To: Project HELIOS Development Team
Subject: Project HELIOS - Final Source-Code Testing

Message:
Project HELIOS is entering its final development stage.

Angela Bennett will coordinate source-code testing and repository access for the development team before the final release.`
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
    // PAGE ELEMENTS
    // ---------------------------------------------------------

    const evidenceList = document.getElementById("analysis-evidence-list");
    const evidenceButtons = document.querySelectorAll(
        ".analysis-evidence-file"
    );

    const evidenceViewer = document.getElementById(
        "analysis-evidence-viewer"
    );
    const evidenceViewerTitle = document.getElementById(
        "analysis-evidence-viewer-title"
    );
    const evidenceViewerContent = document.getElementById(
        "analysis-evidence-viewer-content"
    );
    const closeEvidenceViewerButton = document.getElementById(
        "close-analysis-evidence-viewer"
    );

    // ---------------------------------------------------------
    // RANDOMISE EVIDENCE FILE POSITIONS
    // ---------------------------------------------------------

    function shuffleEvidenceFiles() {
        const files = Array.from(
            evidenceList.querySelectorAll(".analysis-evidence-file")
        );

        for (let i = files.length - 1; i > 0; i--) {
            const randomIndex = Math.floor(Math.random() * (i + 1));

            [files[i], files[randomIndex]] = [
                files[randomIndex],
                files[i]
            ];
        }

        files.forEach((file) => {
            evidenceList.appendChild(file);
        });
    }

    shuffleEvidenceFiles();

    // ---------------------------------------------------------
    // EVIDENCE FILE VIEWER
    // ---------------------------------------------------------

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

        evidenceViewerTitle.textContent = evidence.title;
        renderEvidenceContent(evidenceViewerContent, evidence.content);
        evidenceViewer.hidden = false;
    }

    function closeEvidenceFile() {
        evidenceViewer.hidden = true;
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

    closeEvidenceViewerButton.addEventListener(
        "click",
        closeEvidenceFile
    );

    document.addEventListener("keydown", (event) => {
        if (event.key === "Escape" && !evidenceViewer.hidden) {
            closeEvidenceFile();
        }
    });

    // ---------------------------------------------------------
    // RAW EVIDENCE DRAG / DROP
    // ---------------------------------------------------------

    const correlationSlots = document.querySelectorAll(
        ".correlation-slot"
    );

    const correlateButton = document.getElementById(
        "correlate-button"
    );

    const correlationFeedback = document.getElementById(
        "correlation-feedback"
    );

    const workingTheoryList = document.getElementById(
        "working-theory-list"
    );

    const workingTheoryEmpty = document.getElementById(
        "working-theory-empty"
    );

    const theoryCount = document.getElementById(
        "theory-count"
    );

    const caseFileCulprit = document.getElementById(
        "case-file-culprit"
    );

    const caseFileMethod = document.getElementById(
        "case-file-method"
    );

    const caseFileGoal = document.getElementById(
        "case-file-goal"
    );

    const caseFileMotive = document.getElementById(
        "case-file-motive"
    );

    const solveCaseButton = document.getElementById(
        "solve-case-button"
    );

    const caseSolvedModal = document.getElementById(
        "case-solved-modal"
    );

    const closeCaseButton = document.getElementById(
        "close-case-button"
    );

    const analysisIntroModal1 = document.getElementById(
        "analysis-intro-modal-1"
    );

    const analysisIntroModal2 = document.getElementById(
        "analysis-intro-modal-2"
    );

    const analysisIntroModal3 = document.getElementById(
        "analysis-intro-modal-3"
    );

    const analysisIntroNext1 = document.getElementById(
        "analysis-intro-next-1"
    );

    const analysisIntroNext2 = document.getElementById(
        "analysis-intro-next-2"
    );

    const analysisIntroBegin = document.getElementById(
        "analysis-intro-begin"
    );

    const investigationHelpButton = document.getElementById(
        "investigation-help-button"
    );

    const investigationHintsModal = document.getElementById(
        "investigation-hints-modal"
    );

    const investigationHintsClose = document.getElementById(
        "investigation-hints-close"
    );

    const caseFileAnswers = {
        culprit: caseFileCulprit,
        method: caseFileMethod,
        goal: caseFileGoal,
        motive: caseFileMotive
    };

    const deductionModal = document.getElementById(
        "deduction-established-modal"
    );

    const deductionModalTitle = document.getElementById(
        "deduction-established-title"
    );

    const deductionModalDescription = document.getElementById(
        "deduction-established-description"
    );

    const deductionModalEvidence = document.getElementById(
        "deduction-established-evidence"
    );

    const addDeductionToCaseFileButton = document.getElementById(
        "add-deduction-to-case-file"
    );

    let pendingDeductionId = null;

    const completedDeductions = new Set();

    function updateSolveCaseButton() {
        const requiredDeductions = [
            "culprit",
            "method",
            "goal",
            "motive"
        ];

        const caseComplete = requiredDeductions.every(
            (deductionId) =>
                completedDeductions.has(deductionId)
        );

        solveCaseButton.disabled = !caseComplete;
    }

    let draggedEvidenceId = null;
    let draggedTheoryId = null;
    let draggedFromSlot = null;
    let dragPreview = null;

    function updateCorrelateButton() {
        const bothSlotsFilled = Array.from(correlationSlots).every(
            (slot) =>
                slot.dataset.evidenceId ||
                slot.dataset.theoryId
        );

        correlateButton.disabled = !bothSlotsFilled;
    }

    const theoryDefinitions = {
        C1: {
            title: "Suspicious Account Activity",
            text: "Angela's account was used from MCX-WS17 after her legitimate session had ended.",
            evidence: ["E11", "E13"],
            category: "culprit",
        },

        C2: {
            title: "Eugene's Workstation Used",
            text: "Angela's account was used from MCX-WS17, Eugene's assigned workstation.",
            evidence: ["E06", "E13"],
            category: "culprit",
        },

        C3: {
            title: "Planned Theft Evidence",
            text: "Evidence connected to the planned theft was recovered from Eugene's workstation, MCX-WS17.",
            evidence: ["E06", "E12"],
            category: "culprit",
        },

        C4: {
            title: "Eugene's Workstation Linked to Account Misuse",
            text: "The suspicious use of Angela's account occurred from Eugene's workstation after Angela had logged off.",
            evidence: ["E06", "E11", "E13"],
            category: "culprit"
        },

        C5: {
            title: "Eugene's Workstation Linked to the Theft Plan",
            text: "Angela's account activity and evidence of the planned theft both point to Eugene's workstation.",
            evidence: ["E06", "E12", "E13"],
            category: "culprit"
        },

        M1: {
            title: "Angela Researched Online",
            text: "Eugene researched publicly available information about Angela before the phishing attack.",
            evidence: ["E04", "E05"],
            category: "method"
        },

        M2: {
            title: "Public Information Reused",
            text: "Public information about Angela was reused to personalise the phishing message.",
            evidence: ["E05", "E08"],
            category: "method"
        },

        M3: {
            title: "Phishing Link Opened",
            text: "Angela opened the phishing link shortly after receiving the message.",
            evidence: ["E08", "E09"],
            category: "method"
        },

        M4: {
            title: "Personalised Phishing Linked to Eugene",
            text: "Eugene's research into Angela was connected to a phishing message personalised with her public information.",
            evidence: ["E04", "E05", "E08"],
            category: "method"
        },

        M5: {
            title: "Personalised Phishing Link Opened",
            text: "A phishing message using Angela's public information was opened on her workstation.",
            evidence: ["E05", "E08", "E09"],
            category: "method"
        },

        G1: {
            title: "HELIOS Source Code Connected to NURV",
            text: "The planned source-code transfer to NURV appears connected to Project HELIOS.",
            evidence: ["E10", "E12"],
            category: "goal"
        },

        G2: {
            title: "Planned Transfer Matches Network Event",
            text: "The planned 21:30 source-code transfer to NURV matches the later blocked network transfer.",
            evidence: ["E12", "E14"],
            category: "goal"
        },

        G3: {
            title: "HELIOS Source Code Was the Target",
            text: "Project HELIOS source code was the target of the blocked outbound transfer.",
            evidence: ["E10", "E14"],
            category: "goal"
        },

        MOT1: {
            title: "Repeated Career Setbacks",
            text: "Eugene experienced repeated career setbacks, including a rejected promotion and rejection from Project HELIOS.",
            evidence: ["E01", "E02"],
            category: "motive"
        },

        MOT2: {
            title: "NURV Job Research After HELIOS Rejection",
            text: "After being rejected from Project HELIOS, Eugene researched senior developer opportunities at NURV Technologies.",
            evidence: ["E02", "E07"],
            category: "motive"
        },

        MOT3: {
            title: "NURV Interest Connected to a Deal",
            text: "Eugene's interest in NURV was connected to an agreement involving payment and a position at the company.",
            evidence: ["E07", "E12"],
            category: "motive"
        },

        MOT4: {
            title: "Career Setbacks Followed by NURV Job Research",
            text: "Eugene's repeated career setbacks were followed by research into employment opportunities at NURV Technologies.",
            evidence: ["E01", "E02", "E07"],
            category: "motive"
        },

        MOT5: {
            title: "HELIOS Rejection Linked to the NURV Deal",
            text: "Eugene's rejection from Project HELIOS was followed by NURV job research and an agreement involving payment and a position.",
            evidence: ["E02", "E07", "E12"],
            category: "motive"
        }
    };

    const deductionDefinitions = {
        culprit: {
            title: "Culprit — Eugene",
            description:
                "Angela's account was used from Eugene's workstation after her legitimate session had ended, while evidence of the planned theft was also recovered from Eugene's computer.",
            evidence: ["E06", "E11", "E12", "E13"],
            answer: "Eugene",
            category: "culprit"
        },

        method: {
            title: "Method — Targeted Phishing",
            description:
                "Eugene researched publicly available information about Angela, that information was reused to personalise a phishing message, and Angela opened the link on her workstation.",
            evidence: ["E04", "E05", "E08", "E09"],
            answer: "Targeted phishing using Angela's public information",
            category: "method"
        },

        goal: {
            title: "Goal — Steal Project HELIOS Source Code",
            description:
                "The evidence shows that the planned transfer to NURV targeted Project HELIOS source code, and the outbound transfer was blocked before any data was successfully sent.",
            evidence: ["E10", "E12", "E14"],
            answer: "Steal Project HELIOS source code",
            category: "goal"
        },

        motive: {
            title: "Motive — Career Frustration and a NURV Deal",
            description:
                "Eugene suffered repeated career setbacks, including rejection from Project HELIOS, then researched employment at NURV and entered an agreement involving both payment and a position at the company.",
            evidence: ["E01", "E02", "E07", "E12"],
            answer: "Career frustration and a NURV payment/job offer",
            category: "motive"
        }
    };

    const correlationRules = {
        "E11|E13": {
            result: "C1",
            status: "THEORY FORMED"
        },

        "E06|E13": {
            result: "C2",
            status: "THEORY FORMED"
        },

        "E06|E12": {
            result: "C3",
            status: "THEORY FORMED"
        },

        "C1|E06": {
            result: "C4",
            status: "THEORY STRENGTHENED"
        },

        "C2|E11": {
            result: "C4",
            status: "THEORY STRENGTHENED"
        },

        "C2|E12": {
            result: "C5",
            status: "THEORY STRENGTHENED"
        },

        "C3|E13": {
            result: "C5",
            status: "THEORY STRENGTHENED"
        },

        "C4|E12": {
            deduction: "culprit",
            status: "DEDUCTION ESTABLISHED"
        },

        "C5|E11": {
            deduction: "culprit",
            status: "DEDUCTION ESTABLISHED"
        },

        "C1|C3": {
            deduction: "culprit",
            status: "DEDUCTION ESTABLISHED"
        },

        "E04|E05": {
            result: "M1",
            status: "THEORY FORMED"
        },

        "E05|E08": {
            result: "M2",
            status: "THEORY FORMED"
        },

        "E08|E09": {
            result: "M3",
            status: "THEORY FORMED"
        },

        "E08|M1": {
            result: "M4",
            status: "THEORY STRENGTHENED"
        },

        "E04|M2": {
            result: "M4",
            status: "THEORY STRENGTHENED"
        },

        "E09|M2": {
            result: "M5",
            status: "THEORY STRENGTHENED"
        },

        "E05|M3": {
            result: "M5",
            status: "THEORY STRENGTHENED"
        },

        "E09|M4": {
            deduction: "method",
            status: "DEDUCTION ESTABLISHED"
        },

        "E04|M5": {
            deduction: "method",
            status: "DEDUCTION ESTABLISHED"
        },

        "M1|M3": {
            deduction: "method",
            status: "DEDUCTION ESTABLISHED"
        },

        "E10|E12": {
            result: "G1",
            status: "THEORY FORMED"
        },

        "E12|E14": {
            result: "G2",
            status: "THEORY FORMED"
        },

        "E10|E14": {
            result: "G3",
            status: "THEORY FORMED"
        },

        "E14|G1": {
            deduction: "goal",
            status: "DEDUCTION ESTABLISHED"
        },

        "E10|G2": {
            deduction: "goal",
            status: "DEDUCTION ESTABLISHED"
        },

        "E12|G3": {
            deduction: "goal",
            status: "DEDUCTION ESTABLISHED"
        },

        "E01|E02": {
            result: "MOT1",
            status: "THEORY FORMED"
        },

        "E02|E07": {
            result: "MOT2",
            status: "THEORY FORMED"
        },

        "E07|E12": {
            result: "MOT3",
            status: "THEORY FORMED"
        },

        "E01|E07": {
            status: "POSSIBLE CONNECTION",
            message:
                "Eugene researched employment at NURV after his promotion was rejected. Look for more evidence connecting his career setbacks to NURV."
        },

        "E01|E12": {
            status: "POSSIBLE CONNECTION",
            message:
                "Eugene experienced a career setback and was later offered a position at NURV. Look for evidence showing what happened between these events."
        },

        "E02|E12": {
            status: "POSSIBLE CONNECTION",
            message:
                "Eugene was rejected from Project HELIOS and later entered an agreement with NURV. Look for evidence connecting these events."
        },

        "E07|MOT1": {
            result: "MOT4",
            status: "THEORY STRENGTHENED"
        },

        "E01|MOT2": {
            result: "MOT4",
            status: "THEORY STRENGTHENED"
        },

        "E02|MOT3": {
            result: "MOT5",
            status: "THEORY STRENGTHENED"
        },

        "E12|MOT4": {
            deduction: "motive",
            status: "DEDUCTION ESTABLISHED"
        },

        "E01|MOT5": {
            deduction: "motive",
            status: "DEDUCTION ESTABLISHED"
        },

        "MOT1|MOT3": {
            deduction: "motive",
            status: "DEDUCTION ESTABLISHED"
        }
    };

    const formedTheories = new Set();
    const retiredTheories = new Set();

    function getSlotItemId(slot) {
        return (
            slot.dataset.evidenceId ||
            slot.dataset.theoryId ||
            null
        );
    }

    function createPairKey(firstId, secondId) {
        return [firstId, secondId]
            .sort()
            .join("|");
    }

    function clearCorrelationFeedback() {
        correlationFeedback.textContent =
            "Select two pieces of information to compare.";
    }

    function updateTheoryCount() {
        theoryCount.textContent =
            `${formedTheories.size} ${formedTheories.size === 1 ? "theory" : "theories"}`;
    }

    function isEvidenceSubset(smallerEvidence, largerEvidence) {
        return (
            smallerEvidence.length < largerEvidence.length &&
            smallerEvidence.every((evidenceId) =>
                largerEvidence.includes(evidenceId)
            )
        );
    }

    function retireSupersededTheories(newTheoryId) {
        const newTheory = theoryDefinitions[newTheoryId];

        if (!newTheory) {
            return;
        }

        Object.entries(theoryDefinitions).forEach(
            ([existingTheoryId, existingTheory]) => {
                if (existingTheoryId === newTheoryId) {
                    return;
                }

                if (existingTheory.category !== newTheory.category) {
                    return;
                }

                if (
                    !isEvidenceSubset(
                        existingTheory.evidence,
                        newTheory.evidence
                    )
                ) {
                    return;
                }

                const existingCard = document.querySelector(
                    `.working-theory-card[data-theory-id="${existingTheoryId}"]`
                );

                if (existingCard) {
                    existingCard.remove();
                }

                formedTheories.delete(existingTheoryId);
                retiredTheories.add(existingTheoryId);
            }
        );

        updateTheoryCount();
    }

    function retireTheoryCategory(category) {
        Object.entries(theoryDefinitions).forEach(
            ([theoryId, theory]) => {
                if (theory.category !== category) {
                    return;
                }

                const card = document.querySelector(
                    `.working-theory-card[data-theory-id="${theoryId}"]`
                );

                if (card) {
                    card.remove();
                }

                formedTheories.delete(theoryId);
                retiredTheories.add(theoryId);
            }
        );

        updateTheoryCount();

        if (formedTheories.size === 0) {
            workingTheoryEmpty.style.display = "";
        }
    }

    function populateCaseFileEntry(answerElement, deduction) {
        answerElement.innerHTML = "";

        const answer = document.createElement("strong");
        answer.className = "case-file-answer__result";
        answer.textContent = deduction.answer;

        const support = document.createElement("div");
        support.className = "case-file-answer__support";

        const supportLabel = document.createElement("span");
        supportLabel.className = "case-file-answer__support-label";
        supportLabel.textContent = "Supported by:";

        const supportList = document.createElement("ul");

        deduction.evidence.forEach((evidenceId) => {
            const item = document.createElement("li");
            item.textContent = evidenceFiles[evidenceId].title;
            supportList.appendChild(item);
        });

        support.append(supportLabel, supportList);
        answerElement.append(answer, support);
    }

    function establishDeduction(deductionId) {
        const deduction = deductionDefinitions[deductionId];

        if (!deduction) {
            return false;
        }

        pendingDeductionId = deductionId;

        deductionModalTitle.textContent = deduction.title;
        deductionModalDescription.textContent =
            deduction.description;

        deductionModalEvidence.innerHTML = "";

        deduction.evidence.forEach((evidenceId) => {
            const item = document.createElement("li");
            item.textContent = evidenceFiles[evidenceId].title;
            deductionModalEvidence.appendChild(item);
        });

        deductionModal.hidden = false;

        return true;
    }

    addDeductionToCaseFileButton.addEventListener(
        "click",
        () => {
            if (!pendingDeductionId) {
                return;
            }

            const deduction =
                deductionDefinitions[pendingDeductionId];

            if (!deduction) {
                return;
            }

            const caseFileAnswer =
                caseFileAnswers[pendingDeductionId];

            if (!caseFileAnswer) {
                return;
            }

            populateCaseFileEntry(
                caseFileAnswer,
                deduction
            );

            const caseFileEntry = caseFileAnswer.closest(
                ".case-file-entry"
            );

            if (caseFileEntry) {
                caseFileEntry.classList.add(
                    "case-file-entry--complete"
                );
            }

            completedDeductions.add(
                pendingDeductionId
            );

            retireTheoryCategory(
                deduction.category
            );

            updateSolveCaseButton();

            deductionModal.hidden = true;
            pendingDeductionId = null;
        }
    );

    function addWorkingTheory(theoryId) {
        if (
            formedTheories.has(theoryId) ||
            retiredTheories.has(theoryId)
        ) {
            return false;
        }

        const theory = theoryDefinitions[theoryId];

        if (!theory) {
            return false;
        }

        formedTheories.add(theoryId);

        const card = document.createElement("div");
        card.className = "working-theory-card";
        card.dataset.theoryId = theoryId;
        card.draggable = true;

        card.addEventListener("dragstart", (event) => {
            clearCorrelationFeedback();

            draggedTheoryId = theoryId;
            draggedFromSlot = null;

            event.dataTransfer.effectAllowed = "copy";
            event.dataTransfer.setData(
                "text/plain",
                theoryId
            );
        });

        card.addEventListener("dragend", () => {
            draggedTheoryId = null;
            draggedFromSlot = null;
        });

        const title = document.createElement("strong");
        title.textContent = theory.title;

        const description = document.createElement("span");
        description.textContent = theory.text;

        const support = document.createElement("div");
        support.className = "working-theory-card__support";

        const supportLabel = document.createElement("span");
        supportLabel.className = "working-theory-card__support-label";
        supportLabel.textContent = "Supported by:";

        const supportList = document.createElement("ul");

        theory.evidence.forEach((evidenceId) => {
            const evidenceItem = document.createElement("li");
            evidenceItem.textContent = evidenceFiles[evidenceId].title;
            supportList.appendChild(evidenceItem);
        });

        support.append(supportLabel, supportList);

        card.append(title, description, support);
        workingTheoryList.appendChild(card);

        workingTheoryEmpty.style.display = "none";

        retireSupersededTheories(theoryId);

        updateTheoryCount();

        return true;
    }

    function removeDragPreview() {
        if (dragPreview) {
            dragPreview.remove();
            dragPreview = null;
        }
    }

    function createDragPreview(sourceButton) {
        removeDragPreview();

        dragPreview = sourceButton.cloneNode(true);
        dragPreview.className = "analysis-drag-preview";

        document.body.appendChild(dragPreview);

        return dragPreview;
    }

    function createCorrelationEvidenceItem(evidenceId) {
        const evidence = evidenceFiles[evidenceId];

        const item = document.createElement("button");
        item.type = "button";
        item.className = "correlation-item";
        item.draggable = true;
        item.dataset.evidenceId = evidenceId;

        const icon = document.createElement("span");
        icon.className = "analysis-file-icon";
        icon.setAttribute("aria-hidden", "true");
        icon.textContent = "TXT";

        const name = document.createElement("span");
        name.className = "correlation-item__name";
        name.textContent = evidence.title;

        item.append(icon, name);

        item.addEventListener("dragstart", (event) => {
            clearCorrelationFeedback();

            draggedEvidenceId = evidenceId;
            draggedFromSlot = item.closest(".correlation-slot");

            event.dataTransfer.effectAllowed = "move";
            event.dataTransfer.setData(
                "text/plain",
                evidenceId
            );

            const preview = createDragPreview(item);

            event.dataTransfer.setDragImage(
                preview,
                20,
                20
            );
        });

        item.addEventListener("dragend", () => {
            draggedEvidenceId = null;
            draggedFromSlot = null;

            removeDragPreview();
        });

        return item;
    }

    function createCorrelationTheoryItem(theoryId) {
        const theory = theoryDefinitions[theoryId];

        const item = document.createElement("button");
        item.type = "button";
        item.className = "correlation-item correlation-item--theory";
        item.draggable = true;
        item.dataset.theoryId = theoryId;

        const icon = document.createElement("span");
        icon.className = "correlation-theory-icon";
        icon.setAttribute("aria-hidden", "true");
        icon.textContent = "?";

        const name = document.createElement("span");
        name.className = "correlation-item__name";
        name.textContent = theory.title;

        item.append(icon, name);

        item.addEventListener("dragstart", (event) => {
            clearCorrelationFeedback();

            draggedTheoryId = theoryId;
            draggedFromSlot = item.closest(".correlation-slot");

            event.dataTransfer.effectAllowed = "move";
            event.dataTransfer.setData(
                "text/plain",
                theoryId
            );
        });

        item.addEventListener("dragend", () => {
            draggedTheoryId = null;
            draggedFromSlot = null;

            removeDragPreview();
        });

        return item;
    }

    function setEvidenceInUse(evidenceId, inUse) {
        const evidenceButton = document.querySelector(
            `.analysis-evidence-file[data-evidence-id="${evidenceId}"]`
        );

        if (!evidenceButton) {
            return;
        }

        if (inUse) {
            evidenceButton.style.display = "none";
        } else {
            evidenceButton.style.display = "";
        }
    }

    function setTheoryInUse(theoryId, inUse) {
        const theoryCard = document.querySelector(
            `.working-theory-card[data-theory-id="${theoryId}"]`
        );

        if (theoryCard) {
            theoryCard.style.display = inUse ? "none" : "";
        }

        const theory = theoryDefinitions[theoryId];

        if (!theory) {
            return;
        }

        theory.evidence.forEach((evidenceId) => {
            setEvidenceInUse(evidenceId, inUse);
        });
    }

    function releaseSlotEvidence(slot) {
        const existingEvidenceId = slot.dataset.evidenceId;

        if (!existingEvidenceId) {
            return;
        }

        setEvidenceInUse(existingEvidenceId, false);

        delete slot.dataset.evidenceId;
    }

    function releaseSlotTheory(slot) {
        const existingTheoryId = slot.dataset.theoryId;

        if (!existingTheoryId) {
            return;
        }

        setTheoryInUse(existingTheoryId, false);

        delete slot.dataset.theoryId;
    }

    function emptyCorrelationSlot(slot) {
        releaseSlotEvidence(slot);
        releaseSlotTheory(slot);

        slot.innerHTML = "<span>DROP EVIDENCE OR THEORY</span>";
        slot.classList.remove("correlation-slot--filled");

        updateCorrelateButton();
    }

    function resetCorrelationWorkspace() {
        correlationSlots.forEach((slot) => {
            emptyCorrelationSlot(slot);
        });
    }

    function placeEvidenceInSlot(slot, evidenceId) {
        // If this slot already contains evidence, return that evidence
        // to the Evidence Files list before replacing it.
        releaseSlotEvidence(slot);

        slot.innerHTML = "";

        const item = createCorrelationEvidenceItem(evidenceId);

        slot.appendChild(item);
        slot.classList.add("correlation-slot--filled");
        slot.dataset.evidenceId = evidenceId;

        // Temporarily remove the selected evidence from the left-hand list.
        setEvidenceInUse(evidenceId, true);

        updateCorrelateButton();
    }

    function placeTheoryInSlot(slot, theoryId) {
        releaseSlotEvidence(slot);
        releaseSlotTheory(slot);

        slot.innerHTML = "";

        const item = createCorrelationTheoryItem(theoryId);

        slot.appendChild(item);
        slot.classList.add("correlation-slot--filled");
        slot.dataset.theoryId = theoryId;

        setTheoryInUse(theoryId, true);

        updateCorrelateButton();
    }

    evidenceButtons.forEach((button) => {
        button.draggable = true;

        button.addEventListener("dragstart", (event) => {
            clearCorrelationFeedback();

            draggedEvidenceId = button.dataset.evidenceId;

            button.classList.add(
                "analysis-evidence-file--dragging"
            );

            event.dataTransfer.effectAllowed = "copy";
            event.dataTransfer.setData(
                "text/plain",
                draggedEvidenceId
            );

            const preview = createDragPreview(button);

            event.dataTransfer.setDragImage(
                preview,
                20,
                20
            );
        });

        button.addEventListener("dragend", () => {
            draggedEvidenceId = null;

            button.classList.remove(
                "analysis-evidence-file--dragging"
            );

            removeDragPreview();

            correlationSlots.forEach((slot) => {
                slot.classList.remove(
                    "correlation-slot--dragover"
                );
            });
        });
    });

    evidenceList.addEventListener("dragover", (event) => {
        if (!draggedFromSlot || !draggedEvidenceId) {
            return;
        }

        event.preventDefault();
        event.dataTransfer.dropEffect = "move";
    });

    evidenceList.addEventListener("drop", (event) => {
        if (!draggedFromSlot || !draggedEvidenceId) {
            return;
        }

        event.preventDefault();

        emptyCorrelationSlot(draggedFromSlot);
    });

    workingTheoryList.addEventListener("dragover", (event) => {
        if (!draggedFromSlot || !draggedTheoryId) {
            return;
        }

        event.preventDefault();
        event.dataTransfer.dropEffect = "move";
    });

    workingTheoryList.addEventListener("drop", (event) => {
        if (!draggedFromSlot || !draggedTheoryId) {
            return;
        }

        event.preventDefault();

        emptyCorrelationSlot(draggedFromSlot);
    });

    correlationSlots.forEach((slot) => {
        slot.addEventListener("dragover", (event) => {
            if (
                (!draggedEvidenceId && !draggedTheoryId) ||
                draggedFromSlot
            ) {
                return;
            }

            event.preventDefault();
            event.dataTransfer.dropEffect = "copy";

            slot.classList.add(
                "correlation-slot--dragover"
            );
        });

        slot.addEventListener("dragleave", () => {
            slot.classList.remove(
                "correlation-slot--dragover"
            );
        });

        slot.addEventListener("drop", (event) => {
            if (
                (!draggedEvidenceId && !draggedTheoryId) ||
                draggedFromSlot
            ) {
                return;
            }

            event.preventDefault();

            slot.classList.remove(
                "correlation-slot--dragover"
            );

            if (draggedTheoryId) {
                placeTheoryInSlot(
                    slot,
                    draggedTheoryId
                );
            } else {
                placeEvidenceInSlot(
                    slot,
                    draggedEvidenceId
                );
            }
        });
    });

    // ---------------------------------------------------------
    // CORRELATION PROCESSING
    // ---------------------------------------------------------

    correlateButton.addEventListener("click", () => {
        correlateButton.addEventListener("click", () => {
            const firstItemId = getSlotItemId(
                correlationSlots[0]
            );

            const secondItemId = getSlotItemId(
                correlationSlots[1]
            );

            if (!firstItemId || !secondItemId) {
                return;
            }

            const pairKey = createPairKey(
                firstItemId,
                secondItemId
            );

            console.log("Correlation pair:", pairKey);

            const rule = correlationRules[pairKey];

            if (!rule) {
                correlationFeedback.innerHTML = `
			<strong>NO CLEAR CONNECTION</strong><br>
			These two pieces of evidence do not establish a useful connection yet.
		`;

                resetCorrelationWorkspace();
                return;
            }

            if (rule.deduction) {
                const wasEstablished = establishDeduction(
                    rule.deduction
                );

                if (wasEstablished) {
                    correlationFeedback.innerHTML = `
			<strong>${rule.status}</strong><br>
			The evidence identifies Eugene as the person responsible.
		`;
                }

                resetCorrelationWorkspace();
                return;
            }

            if (rule.message && !rule.result) {
                correlationFeedback.innerHTML = `
		<strong>${rule.status}</strong><br>
		${rule.message}
	`;

                resetCorrelationWorkspace();
                return;
            }

            const wasAdded = addWorkingTheory(rule.result);

            if (wasAdded) {
                correlationFeedback.innerHTML = `
			<strong>${rule.status}</strong><br>
			${theoryDefinitions[rule.result].text}
		`;
            } else {
                correlationFeedback.innerHTML = `
			<strong>CONNECTION ALREADY ESTABLISHED</strong><br>
			This theory is already in your Working Theories.
		`;
            }

            resetCorrelationWorkspace();
        });

    }, true);

    // ---------------------------------------------------------
    // PHASE 2 INTRODUCTION
    // ---------------------------------------------------------

    analysisIntroModal1.hidden = false;

    analysisIntroNext1.addEventListener("click", () => {
        analysisIntroModal1.hidden = true;
        analysisIntroModal2.hidden = false;
    });

    analysisIntroNext2.addEventListener("click", () => {
        analysisIntroModal2.hidden = true;
        analysisIntroModal3.hidden = false;
    });

    analysisIntroBegin.addEventListener("click", () => {
        analysisIntroModal3.hidden = true;
    });

    // ---------------------------------------------------------
    // INVESTIGATION HINTS
    // ---------------------------------------------------------

    investigationHelpButton.addEventListener("click", () => {
        investigationHintsModal.hidden = false;
    });

    investigationHintsClose.addEventListener("click", () => {
        investigationHintsModal.hidden = true;
    });

    // ---------------------------------------------------------
    // FINAL CASE SOLVED
    // ---------------------------------------------------------

    solveCaseButton.addEventListener("click", () => {
        if (solveCaseButton.disabled) {
            return;
        }

        caseSolvedModal.hidden = false;
    });

    closeCaseButton.addEventListener("click", () => {
        window.location.href = "/whwl/digital-investigation";
    });
});