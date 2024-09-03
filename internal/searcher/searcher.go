package searcher

import (
	"errors"
	apiresponse "github.com/ahmadrezamusthafa/search-engine/common/api-response"
	"github.com/ahmadrezamusthafa/search-engine/common/util"
	"github.com/ahmadrezamusthafa/search-engine/internal/storage"
	"net/http"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
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

	query := r.URL.Query().Get("query")
	if query == "" {
		err = errors.New("query parameter 'query' is required")
		return
	}

	results := storage.Search(query)
	response := apiresponse.APIResponse{
		Status: "success",
		Data:   results,
	}
	apiresponse.RespondJSON(w, http.StatusOK, response)
}
