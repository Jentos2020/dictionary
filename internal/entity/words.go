package entity

type Word struct {
	Data       string `json:"data"`
	Dictionary string `json:"dictionary"`
}

type Words = []Word

type SearchRequest struct {
	Prefix string `json:"prefix"`
}

type SearchResponse struct {
	Words Words  `json:"words"`
	Error string `json:"error,omitempty"`
}
