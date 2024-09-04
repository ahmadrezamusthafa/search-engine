package router

import (
	"github.com/ahmadrezamusthafa/search-engine/internal/server/http/handler"
	"github.com/gorilla/mux"
)

func NewRouter(h *handler.Handler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/index", h.IndexHandler).Methods("POST")
	r.HandleFunc("/search", h.SearchHandler).Methods("GET")
	return r
}
