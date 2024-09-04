package handler

import "github.com/ahmadrezamusthafa/search-engine/internal/engine"

type Handler struct {
	Engine *engine.SearchEngine
}

func NewHandler(searchEngine *engine.SearchEngine) *Handler {
	return &Handler{
		Engine: searchEngine,
	}
}
