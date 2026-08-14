package emailtext

import "cyber-security-minigames/types"

// Count returns the number of selectable units in s using SplitUnits,
// the same tokenization the templates use so offsets stay in sync.
func Count(s string) int {
	return len(SplitUnits(s))
}

// PhraseFields returns the email's text fields in render order: Title, To-value,
// From-value, Subject, then Content. This defines the whole-email phrase index
// space used by CorrectPhraseIndexes and PhraseAt.
func PhraseFields(email types.Email) []string {
	return []string{email.Title, email.To, email.From, email.Subject, email.Content}
}

// CorrectSet returns the whole-email indices of suspicious (correct) phrases.
func CorrectSet(email types.Email) map[int]bool {
	set := make(map[int]bool, len(email.CorrectPhraseIndexes))
	for _, idx := range email.CorrectPhraseIndexes {
		set[idx] = true
	}
	return set
}

// PhraseAt returns the phrase at whole-email index idx, walking Title, To-value,
// From-value, Subject, then Content in render order, and whether idx is valid.
func PhraseAt(email types.Email, idx int) (string, bool) {
	for _, field := range PhraseFields(email) {
		if idx < Count(field) {
			if idx < 0 {
				return "", false
			}
			return SplitUnits(field)[idx], true
		}
		idx -= Count(field)
	}
	return "", false
}

// FieldPhrases returns the phrases of the field at field index (0=Title, 1=To,
// 2=From, 3=Subject, 4=Content), in render order.
func FieldPhrases(email types.Email, field int) []string {
	fields := PhraseFields(email)
	if field < 0 || field >= len(fields) {
		return nil
	}
	return SplitUnits(fields[field])
}

// FieldStart returns the whole-email start index of the field at field index.
func FieldStart(email types.Email, field int) int {
	start := 0
	fields := PhraseFields(email)
	for i := 0; i < field && i < len(fields); i++ {
		start += Count(fields[i])
	}
	return start
}
