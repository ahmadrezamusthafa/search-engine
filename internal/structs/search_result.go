package structs

type SearchResult struct {
	ID    string  `json:"id"`
	Score float64 `json:"score"`
}