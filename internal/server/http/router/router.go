package router

import (
	"github.com/ahmadrezamusthafa/search-engine/internal/server/http/handler"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/index", handler.IndexHandler).Methods("POST")
	r.HandleFunc("/search", handler.SearchHandler).Methods("GET")
	return r
}
