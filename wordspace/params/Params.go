package params

import (
	"strings"

	"github.com/manbeardo/wordle_solver/util"
)

type Params struct {
	GrayLetters   [26]bool
	YellowLetters [5][26]bool
	GreenLetters  [5]rune
}

func (params Params) WithGuessForAnswer(guess string, answer string) Params {
	for answerPos, answerLetter := range answer {
		if []rune(guess)[answerPos] == answerLetter {
			params = params.WithGreenLetter(answerLetter, answerPos)
		} else if strings.ContainsRune(guess, answerLetter) {
			params = params.WithYellowLetter(answerLetter, answerPos)
		} else {
			params = params.WithGrayLetter(answerLetter)
		}
	}
	return params
}

func (params Params) WithGreenLetter(letter rune, pos int) Params {
	params.GreenLetters[pos] = letter

	return params
}

func (params Params) WithoutGreenLetter(pos int) Params {
	params.GreenLetters[pos] = rune(0)

	return params
}

func (params Params) HasGreenLetter(letter rune, pos int) bool {
	return params.GreenLetters[pos] == letter
}

func (params Params) HasGreenLetterAnywhere(letter rune) bool {
	for _, greenLetter := range params.GreenLetters {
		if greenLetter == letter {
			return true
		}
	}
	return false
}

func (params Params) WithYellowLetter(letter rune, pos int) Params {
	params.YellowLetters[pos][util.RuneToIndex(letter)] = true

	return params
}

func (params Params) WithoutYellowLetter(letter rune, pos int) Params {
	params.YellowLetters[pos][util.RuneToIndex(letter)] = false

	return params
}

func (params Params) HasYellowLetter(letter rune, pos int) bool {
	return params.YellowLetters[pos][util.RuneToIndex(letter)]
}

func (params Params) WithGrayLetter(letter rune) Params {
	params.GrayLetters[util.RuneToIndex(letter)] = true

	return params
}

func (params Params) WithoutGrayLetter(letter rune) Params {
	params.GrayLetters[util.RuneToIndex(letter)] = false

	return params
}
