package data

// Level instructions shown in the instructions modal at the start of each level.
const (
	InstructionsLevel1 = "Open the emails and compare the sender, links and request. When you think you have found the phishing message, press SELECT on that email."
	InstructionsLevel2 = "Now inspect the phishing email closely. Click every phrase that looks suspicious, then press Submit Selection to check your answer. Selected phrases stay neutral until you submit."
)

// InstructionsIntroLevel1 is the pre-instruction overview shown before the
// level 1 instructions. It frames what the player is about to learn.
const InstructionsIntroLevel1 = "Phishing and social engineering attacks use believable details to trick people into clicking malicious links or handing over information. Somewhere in Angela's inbox is a message that does not belong. Your first task is to identify it."

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
