package router

import (
	"github.com/ahmadrezamusthafa/search-engine/internal/indexer"
	"github.com/ahmadrezamusthafa/search-engine/internal/searcher"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/index", indexer.IndexHandler).Methods("POST")
	r.HandleFunc("/search", searcher.SearchHandler).Methods("GET")
	return r
}
