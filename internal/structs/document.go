package structs

type Document struct {
	ID        string   `json:"id"`
	Content   Content  `json:"content"`
	StopWords []string `json:"stop_words"`
}
