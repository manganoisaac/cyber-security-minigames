package data

// PasswordInstructions is shown in the modal when the password game starts.
const PasswordInstructions = `You've been called into Eugene Belford's office. His computer is locked,
and your objective is simple: guess his password and gain access. 
Search the office for evidence, photos, drawings, calendars, paperwork.
Anything you find gets logged automatically in your notepad. 
When you're ready, hit the login screen and start testing your clues. Correct
characters lock in green, wrong ones turn red.
Good luck, Detective. It won't take much digging to prove that personal
information makes for lousy passwords`

// PasswordWinMessage is shown when the player guesses every phrase correctly.
const PasswordWinMessage = "Well done! You cracked every phrase."

// PasswordPhrases is the order the object phrases must be clicked in.
var PasswordPhrases = []string{"Rufus", "0904"}

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
