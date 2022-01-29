package main

import (
	"fmt"

	"github.com/manbeardo/wordle_solver/letterflagarray"
	"github.com/manbeardo/wordle_solver/util"
	"github.com/manbeardo/wordle_solver/wordspace"
	"github.com/manbeardo/wordle_solver/wordspace/params"
)

func main() {
	wsParams := wordspace.Params{WordList: params.WordListAll}
	ws := wordspace.Get(wsParams)
	for ws.Size() > 1 {
		guesses := []string{}
		uniquesInGuess := 0
		for _, word := range ws.GetWords() {
			uniquesInWord := letterflagarray.FromWord(word).Count()
			if uniquesInWord > uniquesInGuess {
				uniquesInGuess = uniquesInWord
				guesses = append(guesses[0:0], word)
			} else if uniquesInWord == uniquesInGuess {
				guesses = append(guesses, word)
			}
		}
		fmt.Println(len(guesses), "viable guesses")

		bestGuess := ""
		bestGuessWsSize := -1
		bestGuessWsParams := wsParams
		for _, guess := range guesses {
			guessWsParams := wsParams
			for _, letter := range guess {
				guessWsParams = guessWsParams.WithAbsentLetter(util.Letter(letter))
			}
			guessWs := wordspace.Get(guessWsParams)
			if guessWs.Size() > bestGuessWsSize {
				bestGuess = guess
				bestGuessWsSize = guessWs.Size()
				bestGuessWsParams = guessWsParams
			}
		}
		fmt.Println(bestGuess, "leaves", bestGuessWsSize, "possible words")
		wsParams = bestGuessWsParams
		ws = wordspace.Get(wsParams)
	}
}
