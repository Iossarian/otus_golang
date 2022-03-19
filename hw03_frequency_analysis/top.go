package hw03

import (
	"sort"
	"strings"
)

func Top10(text string) []string {
	explodedSlice := strings.Fields(text)
	dictionaryMap := make(map[string]int)
	for _, word := range explodedSlice {
		if dictionaryMap[word] == 0 {
			dictionaryMap[word] = 1
		} else {
			dictionaryMap[word]++
		}
	}
	dictionaries := sortByRepetitions(dictionaryMap)
	return sortLexicographically(dictionaries)
}

type dictionary struct {
	Word        string
	RepeatCount int
}

func sortLexicographically(dictionaries []dictionary) []string {
	groupedByRepeatsCountDict := make(map[int][]string)
	for _, v := range dictionaries {
		groupedByRepeatsCountDict[v.RepeatCount] = append(groupedByRepeatsCountDict[v.RepeatCount], v.Word)
	}
	for _, words := range groupedByRepeatsCountDict {
		sort.SliceStable(words, func(i, j int) bool {
			return words[i] < words[j]
		})
	}
	result := make([]string, 0, 10)
	for i := len(groupedByRepeatsCountDict) + 1; i > 0; i-- {
		if len(result) < 10 {
			result = append(result, groupedByRepeatsCountDict[i]...)
		}
	}

	return result
}

func sortByRepetitions(dictionaryMap map[string]int) []dictionary {
	dictionaries := make([]dictionary, 0, len(dictionaryMap))
	for k, v := range dictionaryMap {
		dictionaries = append(dictionaries, dictionary{k, v})
	}
	sort.Slice(dictionaries, func(i, j int) bool {
		return dictionaries[i].RepeatCount > dictionaries[j].RepeatCount
	})

	return dictionaries
}
