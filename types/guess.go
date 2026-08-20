package types

// GuessStatus is the feedback level for a single password guess.
type GuessStatus int

const (
	GuessIncorrect GuessStatus = iota
	GuessWrongPosition
	GuessCorrect
)

// Guess is a single password guess with its feedback status.
type Guess struct {
	Word   string
	Status GuessStatus
}
