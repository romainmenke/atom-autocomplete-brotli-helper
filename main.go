package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"
)

type thesaurus map[string][]string
type brotli []string

func main() {

	thesF, err := os.Open("data/thes.json")
	if err != nil {
		panic(err)
	}
	defer thesF.Close()

	thesB, err := ioutil.ReadAll(thesF)
	if err != nil {
		panic(err)
	}

	thes := thesaurus{}
	err = json.Unmarshal(thesB, &thes)
	if err != nil {
		panic(err)
	}

	brotliF, err := os.Open("data/brotli.json")
	if err != nil {
		panic(err)
	}
	defer brotliF.Close()

	brotliB, err := ioutil.ReadAll(brotliF)
	if err != nil {
		panic(err)
	}

	brot := brotli{}
	err = json.Unmarshal(brotliB, &brot)
	if err != nil {
		panic(err)
	}

	brotliThesaurus := thesaurus{}

	re := regexp.MustCompile("[0-9]+")

	for word, synonyms := range thes {

		if len(word) < 2 {
			continue
		}

		if len(word) > 12 {
			continue
		}

		if strings.HasPrefix(word, "-") {
			continue
		}

		if strings.Contains(word, " ") {
			continue
		}

		if re.MatchString(word) {
			continue
		}

		for _, synonym := range synonyms {
			for _, brotliWord := range brot {
				if synonym == brotliWord {

					brotliSynonyms, ok := brotliThesaurus[word]
					if ok {
						brotliSynonyms = append(brotliSynonyms, brotliWord)
						brotliThesaurus[word] = brotliSynonyms
					} else {
						brotliSynonyms = []string{brotliWord}
						brotliThesaurus[word] = brotliSynonyms
					}
				}
			}
		}
	}

WORD_LOOP:
	for word, synonyms := range brotliThesaurus {
		for _, brotliWord := range brot {
			if word == brotliWord {

				for _, synonym := range synonyms {
					if synonym == brotliWord {
						continue WORD_LOOP
					}
				}

				brotliSynonyms, ok := brotliThesaurus[word]
				if ok {
					brotliSynonyms = append(brotliSynonyms, brotliWord)
					brotliThesaurus[word] = brotliSynonyms
				}
			}
		}
	}

	fmt.Println(len(brotliThesaurus))

	outB, err := json.Marshal(brotliThesaurus)
	if err != nil {
		panic(err)
	}

	out, err := os.Create("data/brotli_thes.json")
	if err != nil {
		panic(err)
	}

	defer out.Close()

	_, err = io.Copy(out, bufio.NewReader(bytes.NewBuffer(outB)))
	if err != nil {
		panic(err)
	}

	list := []string{}
	for word, _ := range brotliThesaurus {
		list = append(list, word)
	}

	sort.Strings(list)

	outB, err = json.Marshal(list)
	if err != nil {
		panic(err)
	}

	out, err = os.Create("data/brotli_list.json")
	if err != nil {
		panic(err)
	}

	defer out.Close()

	_, err = io.Copy(out, bufio.NewReader(bytes.NewBuffer(outB)))
	if err != nil {
		panic(err)
	}

}
