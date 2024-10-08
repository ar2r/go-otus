package hw03frequencyanalysis

import (
	"log"
	"regexp"
	"sort"
	"strings"
)

var (
	ClearRegExp = regexp.MustCompile(`^[^a-zA-Zа-яА-Я\-]+|[^a-zA-Zа-яА-Я\-]+$`)
	BlackList   = map[string]struct{}{
		"-": {},
	}
)

type kv struct {
	Key   string
	Value int
}

func Top10(text string) []string {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Что-то пошло не так:", err)
		}
	}()

	Words := map[string]int{}

	for _, word := range strings.Fields(text) {
		word = strings.ToLower(word)
		word = ClearRegExp.ReplaceAllString(word, "")

		if len(word) == 0 || isBlacklisted(word) {
			continue
		}

		Words[word]++
	}

	// Сортировка по частоте слова
	kvSlice := sortByValueDesc(Words)

	// Вывод топ 10
	return getTop(kvSlice, 10)
}

func getTop(kvSlice []kv, topCount int) []string {
	var top []string
	for i := 0; i < len(kvSlice) && i < topCount; i++ {
		top = append(top, kvSlice[i].Key)
	}

	return top
}

func sortByValueDesc(words map[string]int) []kv {
	kvSlice := make([]kv, 0, len(words))
	for k, v := range words {
		kvSlice = append(kvSlice, kv{k, v})
	}
	sort.Slice(kvSlice, func(i, j int) bool {
		if kvSlice[i].Value == kvSlice[j].Value {
			return kvSlice[i].Key < kvSlice[j].Key
		}
		return kvSlice[i].Value > kvSlice[j].Value
	})
	return kvSlice
}

func isBlacklisted(word string) bool {
	_, isExist := BlackList[word]
	return isExist
}
