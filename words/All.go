package words

import "sort"

func init() {
	All = []string{}
	All = append(All, Answers...)
	All = append(All, NonAnswers...)
	sort.Strings(All)
}

var All []string
