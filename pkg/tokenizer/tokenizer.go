package tokenizer

import (
	"bytes"
	"github.com/ahmadrezamusthafa/search-engine/common/util"
	"github.com/ahmadrezamusthafa/search-engine/internal/structs"
	"strings"
)

func Tokenize(content structs.Content) []string {
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
	tokens := removeStopWords(filteredBuf.String())
	filteredBuf.Reset()

	return tokens
}

func removeStopWords(text string) []string {
	words := strings.Fields(text)
	filteredWords := make([]string, 0, len(words))

	for _, word := range words {
		if _, found := stopWords[word]; !found {
			filteredWords = append(filteredWords, word)
		}
	}

	return filteredWords
}
