package hw03

import (
	"sort"
	"strings"
)

func Top10(text string) []string {
	if len(text) == 0 {
		return make([]string, 0)
	}
	explodedSlice := strings.Fields(text)
	dictionaryMap := make(map[string]int)
	for _, word := range explodedSlice {
		dictionaryMap[word]++
	}

	return sortByRepetitions(dictionaryMap)
}

type dictionary struct {
	Word        string
	RepeatCount int
}

func sortByRepetitions(dictionaryMap map[string]int) []string {
	dictionaries := make([]dictionary, 0, len(dictionaryMap))
	for k, v := range dictionaryMap {
		dictionaries = append(dictionaries, dictionary{k, v})
	}
	sort.SliceStable(dictionaries, func(i, j int) bool {
		if dictionaries[i].RepeatCount == dictionaries[j].RepeatCount {
			return dictionaries[i].Word < dictionaries[j].Word
		}
		return dictionaries[i].RepeatCount > dictionaries[j].RepeatCount
	})

	return fillResult(dictionaries)
}

func fillResult(dictionaries []dictionary) []string {
	result := make([]string, 0, 10)
	dictionaryLen := len(dictionaries)
	for i := 0; i < 10; i++ {
		if i >= dictionaryLen {
			break
		}
		result = append(result, dictionaries[i].Word)
	}

	return result
}
