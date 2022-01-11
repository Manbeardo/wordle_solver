package wordspace

import (
	"fmt"
	"sort"
	"sync"

	"github.com/manbeardo/wordle_solver/util"
	"github.com/manbeardo/wordle_solver/words"
	"github.com/manbeardo/wordle_solver/wordspace/cache"
	"github.com/manbeardo/wordle_solver/wordspace/params"
)

var WorkerCount = 32
var Verbose = false

type Params = params.Params

type WordSpace struct {
	params           Params
	letter2words     [26]map[string]struct{}
	pos2letter2words [5][26]map[string]struct{}
	allWords         map[string]struct{}
}

type GuessStats struct {
	MaxSpaceSize int
	AvgSpaceSize float64
}

func buildBaseWordSpace() WordSpace {
	idx := WordSpace{
		pos2letter2words: [5][26]map[string]struct{}{},
		letter2words:     [26]map[string]struct{}{},
		allWords:         map[string]struct{}{},
	}
	for pos, letter2words := range idx.pos2letter2words {
		for letterIndex := range letter2words {
			idx.pos2letter2words[pos][letterIndex] = map[string]struct{}{}
		}
	}
	for letterIndex := range idx.letter2words {
		idx.letter2words[letterIndex] = map[string]struct{}{}
	}
	for _, word := range words.List {
		idx.allWords[word] = struct{}{}
		for pos, letter := range word {
			idx.pos2letter2words[pos][util.RuneToIndex(letter)][word] = struct{}{}
			idx.letter2words[util.RuneToIndex(letter)][word] = struct{}{}
		}
	}
	return idx
}

func (idx WordSpace) Size() int {
	return len(idx.allWords)
}

func (idx WordSpace) copy() *WordSpace {
	copy := &WordSpace{
		pos2letter2words: [5][26]map[string]struct{}{},
		letter2words:     [26]map[string]struct{}{},
		allWords:         map[string]struct{}{},
		params:           idx.params,
	}
	for pos, letter2words := range idx.pos2letter2words {
		for letterIndex, words := range letter2words {
			copy.pos2letter2words[pos][letterIndex] = map[string]struct{}{}
			for word := range words {
				copy.pos2letter2words[pos][letterIndex][word] = struct{}{}
			}
		}
	}
	for letterIndex, words := range idx.letter2words {
		copy.letter2words[letterIndex] = map[string]struct{}{}
		for word := range words {
			copy.letter2words[letterIndex][word] = struct{}{}
		}
	}
	for word := range idx.allWords {
		copy.allWords[word] = struct{}{}
	}

	return copy
}

// Removes all words that don't have letter at pos
func (idx WordSpace) withGreenLetter(letter util.Letter, pos int) WordSpace {
	copy := idx.copy()

	for letterIndex, words := range copy.pos2letter2words[pos] {
		if util.IndexToLetter(letterIndex) == letter {
			continue
		}
		for word := range words {
			copy.removeWord(word)
		}
	}
	copy.params = copy.params.WithGreenLetter(letter, pos)

	return *copy
}

// Removes all words that have letter at pos or don't contain letter
func (idx WordSpace) withYellowLetter(letter util.Letter, pos int) WordSpace {
	copy := idx.copy()

	wordsWithLetter := copy.letter2words[letter.AsIndex()]
	for word := range copy.allWords {
		if _, hasLetter := wordsWithLetter[word]; !hasLetter {
			copy.removeWord(word)
		}
	}
	for word := range copy.pos2letter2words[pos][letter.AsIndex()] {
		copy.removeWord(word)
	}
	copy.params = copy.params.WithYellowLetter(letter, pos)

	return *copy
}

// Removes all words that contain letter
func (idx WordSpace) withGrayLetter(letter util.Letter) WordSpace {
	copy := idx.copy()

	for word := range copy.letter2words[letter.AsIndex()] {
		copy.removeWord(word)
	}
	copy.params = copy.params.WithGrayLetter(letter)

	return *copy
}

func (idx WordSpace) GetWords() []string {
	words := make([]string, 0, len(idx.allWords))
	for word := range idx.allWords {
		words = append(words, word)
	}
	sort.Strings(words)
	return words
}

func (idx WordSpace) GetGuessStats(guess string) GuessStats {
	maxSize := 0
	sizeSum := 0
	answerChan := make(chan string)
	statsLock := &sync.Mutex{}
	workerWg := &sync.WaitGroup{}
	go func() {
		for _, answer := range idx.GetWords() {
			answerChan <- answer
		}
		close(answerChan)
	}()
	for i := 0; i < WorkerCount; i++ {
		workerWg.Add(1)
		go func() {
			defer workerWg.Done()
			for answer := range answerChan {
				size := GetSize(idx.params.WithGuessForAnswer(guess, answer))
				statsLock.Lock()
				if size > maxSize {
					maxSize = size
				}
				sizeSum += size
				statsLock.Unlock()
			}
		}()
	}
	workerWg.Wait()
	return GuessStats{
		MaxSpaceSize: maxSize,
		AvgSpaceSize: float64(sizeSum) / float64(len(idx.allWords)),
	}
}

func (idx WordSpace) GetBestGuess() (string, GuessStats) {
	bestGuess, bestStats := "", GuessStats{MaxSpaceSize: len(words.List) + 1}
	for _, guess := range words.List {
		guessStats := idx.GetGuessStats(guess)
		_, guessIsPossibleAnswer := idx.allWords[guess]
		isNewBest := guessStats.MaxSpaceSize < bestStats.MaxSpaceSize ||
			(guessStats.MaxSpaceSize == bestStats.MaxSpaceSize && guessStats.AvgSpaceSize < bestStats.AvgSpaceSize) ||
			(guessStats.MaxSpaceSize == bestStats.MaxSpaceSize && guessStats.AvgSpaceSize == bestStats.AvgSpaceSize && guessIsPossibleAnswer)
		if isNewBest {
			bestGuess, bestStats = guess, guessStats
		}
		if Verbose {
			toPrint := []interface{}{}
			if isNewBest {
				toPrint = append(toPrint, "[NEW BEST!]")
			}
			toPrint = append(toPrint, guess, guessStats)
			fmt.Println(toPrint...)
		}
	}
	return bestGuess, bestStats
}

func (idx WordSpace) GetParams() Params {
	return idx.params
}

func (idx *WordSpace) removeWord(word string) {
	delete(idx.allWords, word)
	for pos, letter := range word {
		delete(idx.pos2letter2words[pos][util.RuneToIndex(letter)], word)
		delete(idx.letter2words[util.RuneToIndex(letter)], word)
	}
}

func Get(p Params) WordSpace {
	cache.Lock(p)
	defer cache.Unlock(p)

	wordSpaceFromCache, isCached := cache.Get(p)
	if isCached {
		return wordSpaceFromCache.(WordSpace)
	}

	wordSpace := getWordSpaceImpl(p)
	cache.Set(p, wordSpace)
	return wordSpace
}

func GetSize(p Params) int {
	cache.LockSize(p)
	defer cache.UnlockSize(p)

	sizeFromCache, isCached := cache.GetSize(p)
	if isCached {
		return sizeFromCache
	}

	size := Get(p).Size()
	cache.SetSize(p, size)
	return size
}

func getWordSpaceImpl(p Params) WordSpace {
	// recursively unroll the changes
	for pos, letter := range p.GreenLetters {
		if letter != util.EmptyLetter {
			return Get(p.WithoutGreenLetter(pos)).withGreenLetter(letter, pos)
		}
	}
	for pos, letterIndexList := range p.YellowLetters {
		for letterIndex, isSet := range letterIndexList {
			if isSet {
				letter := util.IndexToLetter(letterIndex)
				return Get(p.WithoutYellowLetter(letter, pos)).withYellowLetter(letter, pos)
			}
		}
	}
	for letterIndex, isSet := range p.GrayLetters {
		if isSet {
			letter := util.IndexToLetter(letterIndex)
			return Get(p.WithoutGrayLetter(letter)).withGrayLetter(letter)
		}
	}

	// empty params, build the base space
	return buildBaseWordSpace()
}
