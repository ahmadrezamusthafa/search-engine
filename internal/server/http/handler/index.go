package handler

import (
	"encoding/json"
	apiresponse "github.com/ahmadrezamusthafa/search-engine/common/api-response"
	"github.com/ahmadrezamusthafa/search-engine/common/util"
	"github.com/ahmadrezamusthafa/search-engine/internal/structs"
	"github.com/ahmadrezamusthafa/search-engine/pkg/tokenizer"
	"io"
	"net/http"
)

func (h *Handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			response := apiresponse.APIResponse{
				Status:  "error",
				Message: util.CapitalizeFirstWord(err.Error()),
			}
			apiresponse.RespondJSON(w, http.StatusInternalServerError, response)
		}
	}()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	defer r.Body.Close()

	var doc structs.Document
	if err = json.Unmarshal(body, &doc); err != nil {
		return
	}

	tokens := tokenizer.Tokenize(doc.Content, doc.StopWords...)
	h.Engine.StoreDocument(doc.ID, tokens)

	response := apiresponse.APIResponse{
		Status:  "success",
		Message: "Indexed successfully",
	}
	apiresponse.RespondJSON(w, http.StatusOK, response)
}
