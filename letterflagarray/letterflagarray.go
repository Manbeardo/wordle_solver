package letterflagarray

import (
	"fmt"
	"strings"

	"github.com/manbeardo/wordle_solver/util"
)

type LetterFlagArray [26]bool

func FromWord(word string) LetterFlagArray {
	lfa := LetterFlagArray{}
	for _, letter := range word {
		lfa.AddLetter(util.Letter(letter))
	}
	return lfa
}

func Merge(lfas ...LetterFlagArray) LetterFlagArray {
	result := LetterFlagArray{}
	for i := 0; i < 26; i++ {
		letterIsUsed := false
		for _, lfa := range lfas {
			letterIsUsed = letterIsUsed || lfa[i]
		}
		result[i] = letterIsUsed
	}
	return result
}

func (lfa LetterFlagArray) String() string {
	strs := []string{}
	for letterIndex, isSet := range lfa {
		if isSet {
			strs = append(strs, util.IndexToLetter(letterIndex).String())
		}
	}
	return fmt.Sprintf("[%s]", strings.Join(strs, " "))
}

func (lfa LetterFlagArray) HasLetter(letter util.Letter) bool {
	return lfa[letter.AsIndex()]
}

func (lfa *LetterFlagArray) AddLetter(letter util.Letter) {
	lfa[letter.AsIndex()] = true
}

func (lfa *LetterFlagArray) RemoveLetter(letter util.Letter) {
	lfa[letter.AsIndex()] = false
}

func (lfa LetterFlagArray) Count() int {
	count := 0
	for i := 0; i < 26; i++ {
		if lfa[i] {
			count++
		}
	}
	return count
}
