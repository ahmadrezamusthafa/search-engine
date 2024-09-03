package tokenizer

import (
	"github.com/ahmadrezamusthafa/search-engine/common/util"
	"github.com/ahmadrezamusthafa/search-engine/internal/structs"
	"strings"
)

func Tokenize(content structs.Content) []string {
	strContent := strings.ToLower(content.String)
	tokens := strings.Fields(strContent)

	if len(content.Object) > 0 {
		if len(content.ObjectIndexes) > 0 {
			for _, index := range content.ObjectIndexes {
				if v, ok := content.Object[index]; ok {
					tokens = append(tokens, util.InterfaceToString(v))
				}
			}
		} else {
			for _, v := range content.Object {
				tokens = append(tokens, util.InterfaceToString(v))
			}
		}
	}

	return tokens
}
