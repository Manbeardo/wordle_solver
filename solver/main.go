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

	testSolver("abbey")
}

func testSolverForScoringMode(mode wordspace.GuessScoringMode, answer string) {
	params := wordspace.Params{}
	guess := "serai"

	emojiLines := []string{
		fmt.Sprintf("%s ||%s||", getEmojiForGuess(guess, answer), guess),
	}
	for guess != answer {
		var stats wordspace.GuessStats
		params = params.WithGuessForAnswer(guess, answer)
		guess, stats = wordspace.Get(params).GetBestGuess(mode)
		emojiLines = append(emojiLines, fmt.Sprintf("%s ||%s||", getEmojiForGuess(guess, answer), guess))
		fmt.Println(guess, stats)
	}
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
