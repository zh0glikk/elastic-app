package service

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
	"os"
	"sort"
	"strings"
)

type Pair struct {
	Key   string
	Value int64
}

type PairList []Pair

func (p PairList) Len() int {
	return len(p)
}

func (p PairList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p PairList) Less(i, j int) bool {
	return p[i].Value > p[j].Value
}

func Extract(fileTitle string) error {
	dd, err := os.ReadFile(fmt.Sprintf("%s.json", fileTitle))
	if err != nil {
		return errors.Wrap(err, "failed to read file")
	}

	items := make([]gofeed.Item, 0)

	err = json.Unmarshal(dd, &items)
	if err != nil {
		return errors.Wrap(err, "failed to read file")
	}

	words := make([]string, 0)

	for _, i := range items {
		words = append(words, parse(i.Title)...)
		words = append(words, parse(i.Description)...)
	}

	wordsMap := make(map[string]int64)

	for _, w := range words {
		_, ok := wordsMap[w]
		if !ok {
			wordsMap[w] = 0
		}
		wordsMap[w] += 1
	}

	file, err := os.Open("stop-words.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		value := scanner.Text()
		_, ok := wordsMap[value]
		if ok {
			delete(wordsMap, value)
		}
	}

	p := make(PairList, len(wordsMap))

	i := 0

	for k, v := range wordsMap {
		p[i] = Pair{k, v}
		i++
	}

	sort.Sort(p)

	f, err := os.Create(fmt.Sprintf("%s.txt", "extracted"))

	for _, k := range p {
		_, err = fmt.Fprintf(f, "%v\t%v\n", k.Key, k.Value)
		if err != nil {
			return err
		}
	}


	return nil
}

func parse(str string) []string {
	str = strings.Replace(str, ",", "", -1)
	str = strings.Replace(str, ".", "", -1)
	str = strings.Replace(str, "!", "", -1)
	str = strings.Replace(str, "?", "", -1)
	str = strings.Replace(str, "-", "", -1)
	str = strings.Replace(str, "\"", "", -1)
	str = strings.Replace(str, ":", "", -1)
	str = strings.Replace(str, ";", "", -1)

	res := strings.Split(str, " ")

	return res
}
