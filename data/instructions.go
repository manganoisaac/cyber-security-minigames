package data

// Level instructions shown in the instructions modal at the start of each level.
const (
	InstructionsLevel1 = "Click on the orange icon to find the phishing email. For the best experience, play in full screen."
	InstructionsLevel2 = "Click on the suspicious phrases and click the submit button to check"
)

// InstructionsIntroLevel1 is the pre-instruction overview shown before the
// level 1 instructions. It frames what the player is about to learn.
const InstructionsIntroLevel1 = "Phishing & Social Engineering - Deceptive emails use believable details to get you to click a malicious link or hand over information. Somewhere in Angela's inbox is a message that doesn't belong. Open the emails and compare the sender, the links and the request to work out which one is the phishing attempt. Spotting the warning signs is how you stop an attacker before they get in."

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
