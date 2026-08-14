package data

// Level instructions shown in the instructions modal at the start of each level.
const (
	InstructionsLevel1 = "Level 1 Instructions: Click on the orange icon to find the phishing email"
	InstructionsLevel2 = "Level 2 Instructions: Click on the suspicious phrases and click the submit button to check"
)

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
