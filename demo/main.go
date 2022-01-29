package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/manbeardo/wordle_solver/util"
	"github.com/manbeardo/wordle_solver/wordspace"
)

func main() {
	for i := 1; i < len(os.Args); i++ {
		testSolver(os.Args[i])
	}
}

func testSolverForOptions(options wordspace.GuessOptions, answer string) {
	params := wordspace.Params{}
	guess := "raise"

	emojiLines := []string{
		fmt.Sprintf(
			"%s ||%s, %d words remain||",
			util.GetEmojiForGuess(guess, answer),
			guess,
			wordspace.Get(params.WithGuessForAnswer(guess, answer)).Size(),
		),
	}
	for guess != answer {
		params = params.WithGuessForAnswer(guess, answer)
		ws := wordspace.Get(params)
		guess, _ = ws.GetBestGuess(options)
		emojiLines = append(
			emojiLines,
			fmt.Sprintf(
				"%s ||%s, %d words remain||",
				util.GetEmojiForGuess(guess, answer),
				guess,
				wordspace.Get(params.WithGuessForAnswer(guess, answer)).Size(),
			),
		)
	}
	fmt.Println(strings.Join(emojiLines, "\n"))
}

func precompurdle(answer string) (bool, string) {
	params := wordspace.Params{}
	guesses := []string{
		"jumby",
		"vozhd",
		"flick",
		"twang",
		"peers",
	}
	logStr := ""
	for _, guess := range guesses {
		params = params.WithGuessForAnswer(guess, answer)
		ws := wordspace.Get(params)
		logStr += fmt.Sprintf(
			"%s ||%s, %d words remain||\n",
			util.GetEmojiForGuess(guess, answer),
			guess,
			ws.Size(),
		)
	}
	ws := wordspace.Get(params)
	logStr += fmt.Sprintf(
		"%s ||%s||\n",
		util.GetEmojiForGuess(ws.GetWords()[0], answer),
		ws.GetWords()[0],
	)
	return ws.Size() == 1, logStr
}

func testSolver(answer string) {
	fmt.Println("Robowordle (dumb mode)")
	_, precompLog := precompurdle(answer)
	fmt.Print(precompLog)

	fmt.Println("Robowordle")
	testSolverForOptions(
		wordspace.GuessOptions{ScoringMode: wordspace.ScoreByLowestMaxSize},
		answer,
	)

	fmt.Println("Robowordle (hard mode)")
	testSolverForOptions(
		wordspace.GuessOptions{ScoringMode: wordspace.ScoreByLowestMaxSize, HardMode: true},
		answer,
	)
}
