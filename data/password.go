package data

// PasswordInstructions is shown in the modal when the password game starts.
const PasswordInstructions = "Click an object, then submit guesses. Green = right spot, orange = right word wrong spot, red = wrong."

// PasswordWinMessage is shown when the player guesses every phrase correctly.
const PasswordWinMessage = "Well done! You cracked every phrase."

// PasswordPhrases is the order the object phrases must be clicked in.
var PasswordPhrases = []string{"Rufus", "object2", "object3"}

// PasswordObjectPhrases maps each object id to its phrase.
var PasswordObjectPhrases = map[int]string{
	1: "object1",
	2: "object2",
	3: "object3",
	4: "Rufus",
	5: "drawing",
}

// InstructionsForPassword returns the password game's instructions.
func InstructionsForPassword() string {
	return PasswordInstructions
}
