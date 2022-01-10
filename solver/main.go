package main

import (
	"fmt"

	"github.com/manbeardo/wordle_solver/wordspace"
)

func main() {
	ws := wordspace.Get(wordspace.Params{})
	fmt.Println("soare", ws.GetGuessStats("soare"))
	ws = ws.
		WithGrayLetter('s').
		WithGrayLetter('o').
		WithGrayLetter('a').
		WithGreenLetter('r', 3).
		WithYellowLetter('e', 4)
	fmt.Println(ws.GetBestGuess())
	ws = ws.
		WithYellowLetter('e', 0).
		WithYellowLetter('y', 1).
		WithYellowLetter('e', 1).
		WithGreenLetter('e', 2).
		WithGreenLetter('r', 3).
		WithGrayLetter('s')
	fmt.Println(ws.GetBestGuess())

	ws = wordspace.Get(wordspace.Params{}.WithGuessForAnswer("soare", "query"))
	fmt.Println(ws.GetBestGuess())
	ws = wordspace.Get(ws.GetParams().WithGuessForAnswer("eyers", "query"))
	fmt.Println(ws.GetBestGuess())
}
