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
	// bestGuess, bestGuessStats := ws.GetBestGuess(wordspace.ScoreByLowestMaxSize)
	// fmt.Println("best guess (lowest max size):", bestGuess, bestGuessStats)
	// bestGuess, bestGuessStats = ws.GetBestGuess(wordspace.ScoreByLowestAvgSize)
	// fmt.Println("best guess (lowest avg size):", bestGuess, bestGuessStats)

	fmt.Println(wordspace.Get(wordspace.Params{}.
		WithResult(util.Result{
			{Rune: 'r', Color: util.ColorGray},
			{Rune: 'a', Color: util.ColorGray},
			{Rune: 'i', Color: util.ColorGray},
			{Rune: 's', Color: util.ColorGray},
			{Rune: 'e', Color: util.ColorGray},
		}).
		WithResult(util.Result{
			{Rune: 'b', Color: util.ColorGray},
			{Rune: 'l', Color: util.ColorGray},
			{Rune: 'u', Color: util.ColorGreen},
			{Rune: 'd', Color: util.ColorGray},
			{Rune: 'y', Color: util.ColorGray},
		}).
		WithResult(util.Result{
			{Rune: 'c', Color: util.ColorGreen},
			{Rune: 'o', Color: util.ColorGray},
			{Rune: 'm', Color: util.ColorGray},
			{Rune: 'p', Color: util.ColorGray},
			{Rune: 't', Color: util.ColorGray},
		}).
		WithResult(util.Result{
			{Rune: 'c', Color: util.ColorGreen},
			{Rune: 'h', Color: util.ColorGreen},
			{Rune: 'u', Color: util.ColorGreen},
			{Rune: 'c', Color: util.ColorGray},
			{Rune: 'k', Color: util.ColorGreen},
		}),
	).GetBestGuess(wordspace.ScoreByLowestMaxSize))

	// testSolver("abbey")
}

func testSolverForScoringMode(mode wordspace.GuessScoringMode, answer string) {
	params := wordspace.Params{}
	guess := "raise"

	emojiLines := []string{
		fmt.Sprintf("%s ||%s||", getEmojiForGuess(guess, answer), guess),
	}
	for guess != answer {
		var stats wordspace.GuessStats
		params = params.WithGuessForAnswer(guess, answer)
		ws := wordspace.Get(params)
		guess, stats = ws.GetBestGuess(mode)
		emojiLines = append(emojiLines, fmt.Sprintf("%s ||%s||", getEmojiForGuess(guess, answer), guess))
		fmt.Println(guess, stats)
		if ws.Size() < 100 {
			fmt.Println(ws.GetWords())
		}
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
	for _, letterInfo := range util.GetResult(guess, answer) {
		if letterInfo.Color == util.ColorGreen {
			emojis = append(emojis, "🟩")
		} else if letterInfo.Color == util.ColorYellow {
			emojis = append(emojis, "🟨")
		} else if letterInfo.Color == util.ColorGray {
			emojis = append(emojis, "⬜")
		}
	}
	return strings.Join(emojis, " ")
}
