package util

import (
	"fmt"
	"strings"
)

type Result [5]struct {
	Rune  rune
	Color Color
}

func guessHasEveryCorrectLocationOfRune(guess string, answer string, r rune) bool {
	for pos, answerRune := range answer {
		if answerRune == r && []rune(guess)[pos] != r {
			return false
		}
	}
	return true
}

func GetResult(guess string, answer string) Result {
	result := Result{}
	for guessPos, guessRune := range guess {
		result[guessPos].Rune = guessRune
		if []rune(answer)[guessPos] == guessRune {
			result[guessPos].Color = ColorGreen
		} else if strings.ContainsRune(answer, guessRune) && !guessHasEveryCorrectLocationOfRune(guess, answer, guessRune) {
			result[guessPos].Color = ColorYellow
		} else {
			result[guessPos].Color = ColorGray
		}
	}
	return result
}

type Color string

const (
	ColorGray   Color = "gray"
	ColorYellow Color = "yellow"
	ColorGreen  Color = "green"
)

const EmptyLetter = Letter(0)

func RuneToIndex(r rune) int {
	return int(r) - int('a')
}

type Letter rune

func IndexToLetter(i int) Letter {
	return Letter(int('a') + i)
}

func (l Letter) String() string {
	return string([]rune{rune(l)})
}

func (l Letter) AsIndex() int {
	return RuneToIndex(rune(l))
}

type LetterFlagArray [26]bool

func (lfa LetterFlagArray) String() string {
	strs := []string{}
	for letterIndex, isSet := range lfa {
		if isSet {
			strs = append(strs, IndexToLetter(letterIndex).String())
		}
	}
	return fmt.Sprintf("[%s]", strings.Join(strs, " "))
}

func (lfa LetterFlagArray) HasLetter(letter Letter) bool {
	return lfa[letter.AsIndex()]
}

func (lfa *LetterFlagArray) AddLetter(letter Letter) {
	lfa[letter.AsIndex()] = true
}

func (lfa *LetterFlagArray) RemoveLetter(letter Letter) {
	lfa[letter.AsIndex()] = false
}
