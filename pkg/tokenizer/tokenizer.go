package tokenizer

import (
	"bytes"
	"github.com/ahmadrezamusthafa/search-engine/common/util"
	"github.com/ahmadrezamusthafa/search-engine/internal/structs"
	"strings"
)

func Tokenize(content structs.Content, stopWords ...string) []string {
	buf := bytes.Buffer{}
	buf.WriteString(strings.ToLower(content.String))

	if len(content.Object) > 0 {
		if len(content.ObjectIndexes) > 0 {
			for _, index := range content.ObjectIndexes {
				if v, ok := content.Object[index]; ok {
					if buf.Len() > 0 {
						buf.WriteString(" ")
					}
					buf.WriteString(util.InterfaceToString(v))
				}
			}
		} else {
			for _, v := range content.Object {
				if buf.Len() > 0 {
					buf.WriteString(" ")
				}
				buf.WriteString(util.InterfaceToString(v))
			}
		}
	}

	filteredBuf := bytes.Buffer{}
	for _, r := range buf.String() {
		if r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == ' ' {
			filteredBuf.WriteRune(r)
		}
	}

	buf.Reset()
	tokens := removeStopWords(filteredBuf.String(), stopWords...)
	filteredBuf.Reset()

	return tokens
}

func removeStopWords(text string, stopWords ...string) []string {
	words := strings.Fields(text)
	filteredWords := make([]string, 0, len(words))
	customStopWords := make(map[string]interface{})

	for _, stopWord := range stopWords {
		customStopWords[stopWord] = nil
	}

	for _, word := range words {
		_, isCustomStopWordFound := customStopWords[word]
		if _, isDefaultStopWordFound := defaultStopWords[word]; !isDefaultStopWordFound && !isCustomStopWordFound {
			filteredWords = append(filteredWords, word)
		}
	}

	return filteredWords
}
