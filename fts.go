package main

import (
	"strings"
	"unicode"

	"github.com/kljensen/snowball"
)

type Documment struct {
	ID    int
	Title string `xml:"title"`
	URL   string `xml:"url"`
	Text  string `xml:"abstract"`
}

func tokenize(input string) []string {
	return strings.FieldsFunc(input, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}

func lowercaseNormalizer(tokens []string) []string {
	ret := make([]string, len(tokens))
	for idx, token := range tokens {
		ret[idx] = strings.ToLower(token)
	}
	return ret
}

var stopwords = map[string]struct{}{
	"a": {}, "and": {}, "be": {}, "have": {}, "i": {},
	"in": {}, "of": {}, "that": {}, "the": {}, "to": {},
}

func stopwordsFilter(tokens []string, stopwords map[string]struct{}) []string {
	ret := make([]string, 0, len(tokens))
	for _, token := range tokens {
		if _, found := stopwords[token]; found {
			continue
		}
		ret = append(ret, token)
	}
	return ret
}

func stemming(tokens []string) []string {
	r := make([]string, len(tokens))
	for idx, token := range tokens {
		stemmed, _ := snowball.Stem(token, "english", false)
		r[idx] = stemmed
	}
	return r
}

func textToTokens(input string) ([]string, error) {
	return stemming(
		stopwordsFilter(
			lowercaseNormalizer(
				tokenize(input),
			),
			stopwords),
	), nil
}

type index map[string][]int

func (idx index) add(docs []Documment) {
	for id, doc := range docs {
		tokens, err := textToTokens(doc.Text)
		if err != nil {
			continue
		}

		for _, token := range tokens {
			idxs := idx[token]
			// not duplicate id
			if idxs != nil &&
				idxs[len(idxs)-1] == id { // cause here ids increament
				continue
			}
			idx[token] = append(idxs, id)
		}
	}
}

// return list of doc that contain all tokens
func intersectionSelect(tokensDocs [][]int) []int {
	// cause a, b are list of increament int
	intersection := func(a, b []int) []int {
		maxLen := len(a)
		if maxLen < len(b) {
			maxLen = len(b)
		}
		r := make([]int, 0, maxLen)
		i := 0
		j := 0
		for i < len(a) && j < len(b) {
			if a[i] < b[j] {
				i++
			} else if a[i] > b[j] {
				j++
			} else {
				r = append(r, a[i])
				i++
				j++
			}
		}
		return r
	}

	docs := tokensDocs[0]
	for i := 1; i < len(tokensDocs); i++ {
		tokenDocs := tokensDocs[i]
		docs = intersection(docs, tokenDocs)
	}
	return docs
}

func (idx index) query(term string) []int {
	tokens, err := textToTokens(term)
	if err != nil {
		return nil
	}

	tokensDocs := [][]int{}
	for _, token := range tokens {
		if idx[token] == nil {
			// fail-fast here
			// cause term contain token that not yet indexed
			// intersectionSelect will return nothing
			return nil
		}
		tokensDocs = append(tokensDocs, idx[token])
	}

	docs := intersectionSelect(tokensDocs)
	return docs
}
