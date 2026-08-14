package emailtext

import "strings"

// Delimiter separates the game's selectable units (phrases) inside email
// fields. Change it here to change how every field is tokenized.
const Delimiter = "\n"

// SplitUnits splits s into the game's selectable units on Delimiter.
func SplitUnits(s string) []string {
	return strings.Split(s, Delimiter)
}
