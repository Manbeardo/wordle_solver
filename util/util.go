package util

func RuneToIndex(letter rune) int {
	return int(letter) - int('a')
}

func IndexToRune(i int) rune {
	return rune(int('a') + i)
}
