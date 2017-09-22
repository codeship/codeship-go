package codeship

type pagination struct {
	Total   int `json:"total,omitempty"`
	PerPage int `json:"per_page,omitempty"`
	Page    int `json:"page,omitempty"`
}
