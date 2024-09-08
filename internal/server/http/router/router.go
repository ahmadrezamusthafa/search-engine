package router

import (
	"embed"
	"github.com/ahmadrezamusthafa/search-engine/internal/server/http/handler"
	"github.com/gorilla/mux"
	"net/http"
)

//go:embed web/*
var staticFiles embed.FS

func NewRouter(h *handler.Handler) *mux.Router {
	r := mux.NewRouter()

	staticFS := http.FS(staticFiles)
	staticHandler := http.FileServer(staticFS)
	r.PathPrefix("/web").Handler(staticHandler)

	r.HandleFunc("/index", h.IndexHandler).Methods("POST")
	r.HandleFunc("/search", h.SearchHandler).Methods("GET")
	return r
}
