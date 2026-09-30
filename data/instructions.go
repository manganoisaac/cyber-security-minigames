package data

// Level instructions shown in the instructions modal at the start of each level.
const (
	InstructionsLevel1 = `Open the emails one at a time and compare the sender address, links, wording and request. A familiar name or professional-looking message does not automatically make an email legitimate. When you believe you have found the phishing message, press SELECT on that email. You will receive feedback and can keep investigating if your choice is incorrect.`

	InstructionsLevel2 = `Now analyse the phishing email in detail. Click every phrase or detail that you believe is a phishing warning sign, then press Submit Selection to check your choices.

Selected phrases remain neutral until you submit, so use the content and context of the message rather than relying on colour feedback. If your selection is incomplete or includes incorrect indicators, use the feedback to reconsider the email and try again.`
)

// InstructionsIntroLevel1 is the pre-instruction overview shown before the
// level 1 instructions. It frames what the player is about to learn.
const InstructionsIntroLevel1 = `Phishing is a form of social engineering that uses deceptive messages to persuade people to click malicious links, open unsafe attachments, reveal sensitive information or take another action that benefits an attacker.

A phishing message may imitate a trusted person or organisation, and targeted attacks can use personal information about the victim to make the message appear more relevant and believable. Warning signs can appear in the sender address, link destination, urgency, wording or the request being made.

Somewhere in Angela's inbox is a phishing message. Compare the emails carefully and identify the message that does not belong.`

// InstructionsForLevel returns the instruction text for a level, or "" if unknown.
func InstructionsForLevel(level int) string {
	switch level {
	case 1:
		return InstructionsLevel1
	case 2:
		return InstructionsLevel2
	}
	return ""
}
