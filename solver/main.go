package main

import (
	"fmt"

	"github.com/manbeardo/wordle_solver/wordspace"
)

func main() {
	// wordspace.Verbose = true
	// ws := wordspace.Get(wordspace.Params{})
	// bestGuess, bestGuessStats := ws.GetBestGuess()
	// fmt.Println("best guess:", bestGuess, bestGuessStats)

	ws := wordspace.Get(
		wordspace.Params{}.
			WithGrayLetter('s').
			WithGrayLetter('e').
			WithGrayLetter('r').
			WithGrayLetter('a').
			WithGrayLetter('i').
			WithGrayLetter('p').
			WithGrayLetter('h').
			WithYellowLetter('o', 2).
			WithGrayLetter('n').
			WithGreenLetter('y', 4).
			WithGrayLetter('m').
			WithGrayLetter('u').
			WithYellowLetter('l', 2).
			WithGrayLetter('c').
			WithGrayLetter('t').
			WithGrayLetter('d').
			WithYellowLetter('w', 1).
			WithGrayLetter('a').
			WithGrayLetter('n').
			WithGrayLetter('g'),
	)
	fmt.Println(ws.GetBestGuess())
}
