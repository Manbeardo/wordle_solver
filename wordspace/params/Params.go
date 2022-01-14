package params

import (
	"fmt"
	"strings"

	"github.com/manbeardo/wordle_solver/util"
)

type Params struct {
	AbsentLetters    util.LetterFlagArray
	ElsewhereLetters [5]util.LetterFlagArray
	CorrectLetters   [5]util.Letter
}

func (params Params) WithGuessForAnswer(guess string, answer string) Params {
	return params.WithResult(util.GetResult(guess, answer))
}

func (params Params) WithResult(result util.Result) Params {
	greenLetters := map[rune]struct{}{}
	for _, letterInfo := range result {
		if letterInfo.Color == util.ColorGreen {
			greenLetters[letterInfo.Rune] = struct{}{}
		}
	}

	for pos, letterInfo := range result {
		letter := util.Letter(letterInfo.Rune)
		if letterInfo.Color == util.ColorGreen {
			params = params.WithCorrectLetter(letter, pos)
		} else if letterInfo.Color == util.ColorYellow {
			params = params.WithElsewhereLetter(letter, pos)
		} else if letterInfo.Color == util.ColorGray {
			_, hasGreenInResult := greenLetters[letterInfo.Rune]
			if hasGreenInResult {
				params = params.WithElsewhereLetter(letter, pos)
				for i := 0; i < 5; i++ {
					if result[i].Rune != letterInfo.Rune {
						params = params.WithElsewhereLetter(letter, i)
					}
				}
			} else {
				params = params.WithAbsentLetter(letter)
			}
		} else {
			panic("unrecognized color")
		}
	}

	return params
}

func (params Params) WithCorrectLetter(letter util.Letter, pos int) Params {
	params.CorrectLetters[pos] = letter

	return params
}

func (params Params) WithoutCorrectLetter(pos int) Params {
	params.CorrectLetters[pos] = util.Letter(0)

	return params
}

func (params Params) HasCorrectLetter(letter util.Letter, pos int) bool {
	return params.CorrectLetters[pos] == letter
}

func (params Params) WithElsewhereLetter(letter util.Letter, pos int) Params {
	params.ElsewhereLetters[pos].AddLetter(letter)

	return params
}

func (params Params) WithoutElsewhereLetter(letter util.Letter, pos int) Params {
	params.ElsewhereLetters[pos].RemoveLetter(letter)

	return params
}

func (params Params) HasElsewhereLetter(letter util.Letter, pos int) bool {
	return params.ElsewhereLetters[pos][letter.AsIndex()]
}

func (params Params) WithAbsentLetter(letter util.Letter) Params {
	params.AbsentLetters.AddLetter(letter)

	return params
}

func (params Params) WithoutAbsentLetter(letter util.Letter) Params {
	params.AbsentLetters.RemoveLetter(letter)

	return params
}

func (params Params) String() string {
	yellowStrs := []string{}
	for pos, yellow := range params.ElsewhereLetters {
		yellowStrs = append(yellowStrs, fmt.Sprintf("%d: %s", pos, yellow.String()))
	}
	greenStrs := []string{}
	for pos, green := range params.CorrectLetters {
		if green != util.EmptyLetter {
			greenStrs = append(greenStrs, fmt.Sprintf("%d: %s", pos, green.String()))
		}
	}
	return fmt.Sprintf(
		"gray: %s yellow: {%s} green: {%s}",
		params.AbsentLetters.String(),
		strings.Join(yellowStrs, " "),
		strings.Join(greenStrs, " "),
	)
}
