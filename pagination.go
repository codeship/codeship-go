package codeship

type pagination struct {
	Total   int `json:"total"`
	PerPage int `json:"per_page"`
	Page    int `json:"page"`
}
