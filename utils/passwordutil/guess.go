package passwordutil

import "cyber-security-minigames/types"

// GuessResults classifies each guess word against the correct phrases: a guess
// matching the phrase at its own slot is correct, matching a phrase in another
// slot is wrongPosition, and matching nothing is incorrect.
func GuessResults(guesses, correctPhrases []string) []types.Guess {
	results := make([]types.Guess, len(guesses))
	for i, guess := range guesses {
		status := types.GuessIncorrect
		for j, phrase := range correctPhrases {
			if phrase == guess {
				if i == j {
					status = types.GuessCorrect
				} else {
					status = types.GuessWrongPosition
				}
				break
			}
		}
		results[i] = types.Guess{Word: guess, Status: status}
	}
	return results
}

// AllCorrect reports whether every guess landed in its right slot.
func AllCorrect(results []types.Guess) bool {
	for _, r := range results {
		if r.Status != types.GuessCorrect {
			return false
		}
	}
	return true
}
