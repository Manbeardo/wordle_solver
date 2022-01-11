package params

import (
	"fmt"
	"strings"

	"github.com/manbeardo/wordle_solver/util"
)

type Params struct {
	GrayLetters   util.LetterFlagArray
	YellowLetters [5]util.LetterFlagArray
	GreenLetters  [5]util.Letter
}

func (params Params) WithGuessForAnswer(guess string, answer string) Params {
	for guessPos, guessRune := range guess {
		guessLetter := util.Letter(guessRune)
		if []util.Letter(answer)[guessPos] == guessLetter {
			params = params.WithGreenLetter(guessLetter, guessPos)
		} else if strings.ContainsRune(answer, guessRune) {
			params = params.WithYellowLetter(guessLetter, guessPos)
		} else {
			params = params.WithGrayLetter(guessLetter)
		}
	}
	return params
}

func (params Params) WithGreenLetter(letter util.Letter, pos int) Params {
	params.GreenLetters[pos] = letter

	return params
}

func (params Params) WithoutGreenLetter(pos int) Params {
	params.GreenLetters[pos] = util.Letter(0)

	return params
}

func (params Params) HasGreenLetter(letter util.Letter, pos int) bool {
	return params.GreenLetters[pos] == letter
}

func (params Params) HasGreenLetterAnywhere(letter util.Letter) bool {
	for _, greenLetter := range params.GreenLetters {
		if greenLetter == letter {
			return true
		}
	}
	return false
}

func (params Params) WithYellowLetter(letter util.Letter, pos int) Params {
	params.YellowLetters[pos].AddLetter(letter)

	return params
}

func (params Params) WithoutYellowLetter(letter util.Letter, pos int) Params {
	params.YellowLetters[pos].RemoveLetter(letter)

	return params
}

func (params Params) HasYellowLetter(letter util.Letter, pos int) bool {
	return params.YellowLetters[pos][letter.AsIndex()]
}

func (params Params) WithGrayLetter(letter util.Letter) Params {
	params.GrayLetters.AddLetter(letter)

	return params
}

func (params Params) WithoutGrayLetter(letter util.Letter) Params {
	params.GrayLetters.RemoveLetter(letter)

	return params
}

func (params Params) String() string {
	yellowStrs := []string{}
	for pos, yellow := range params.YellowLetters {
		yellowStrs = append(yellowStrs, fmt.Sprintf("%d: %s", pos, yellow.String()))
	}
	greenStrs := []string{}
	for pos, green := range params.GreenLetters {
		if green != util.EmptyLetter {
			greenStrs = append(greenStrs, fmt.Sprintf("%d: %s", pos, green.String()))
		}
	}
	return fmt.Sprintf(
		"gray: %s yellow: {%s} green: {%s}",
		params.GrayLetters.String(),
		strings.Join(yellowStrs, " "),
		strings.Join(greenStrs, " "),
	)
}
