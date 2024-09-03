package structs

type Content struct {
	String        string                 `json:"string"`
	Object        map[string]interface{} `json:"object"`
	ObjectIndexes []string               `json:"object_indexes"`
}
