package main

import (
	"fmt"
	"strings"

	"github.com/manbeardo/wordle_solver/util"
	"github.com/manbeardo/wordle_solver/wordspace"
)

func main() {
	// wordspace.Verbose = true
	// ws := wordspace.Get(wordspace.Params{})
	// bestGuess, bestGuessStats := ws.GetBestGuess()
	// fmt.Println("best guess:", bestGuess, bestGuessStats)

	testSolver("favor")
}

func testSolverForScoringMode(mode wordspace.GuessScoringMode, answer string) {
	params := wordspace.Params{}
	emojiLines := []string{}
	for guess := "serai"; guess != answer; {
		var stats wordspace.GuessStats
		emojiLines = append(emojiLines, getEmojiForGuess(guess, answer)+" "+guess)
		params = params.WithGuessForAnswer(guess, answer)
		guess, stats = wordspace.Get(params).GetBestGuess(mode)
		fmt.Println(guess, stats)
	}
	emojiLines = append(emojiLines, getEmojiForGuess(answer, answer)+" "+answer)
	fmt.Println(answer, "correct!")
	fmt.Println(strings.Join(emojiLines, "\n"))
}

func testSolver(answer string) {
	fmt.Println("scored by lowest max size")
	testSolverForScoringMode(wordspace.ScoreByLowestMaxSize, answer)

	fmt.Println("scored by lowest avg size")
	testSolverForScoringMode(wordspace.ScoreByLowestAvgSize, answer)
}

func getEmojiForGuess(guess string, answer string) string {
	emojis := []string{}
	for guessPos, guessRune := range guess {
		guessLetter := util.Letter(guessRune)
		if []util.Letter(answer)[guessPos] == guessLetter {
			emojis = append(emojis, "🟩")
		} else if strings.ContainsRune(answer, guessRune) {
			emojis = append(emojis, "🟨")
		} else {
			emojis = append(emojis, "⬜")
		}
	}
	return strings.Join(emojis, " ")
}
