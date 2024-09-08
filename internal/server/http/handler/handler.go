package handler

import (
	"github.com/ahmadrezamusthafa/search-engine/internal/engine"
)

type Handler struct {
	SearchEngine engine.ISearchEngine
}

func NewHandler(searchEngine engine.ISearchEngine) *Handler {
	return &Handler{
		SearchEngine: searchEngine,
	}
}
