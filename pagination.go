package codeship

import (
	"net/url"
	"strconv"
)

type pagination struct {
	Total   int `json:"total,omitempty"`
	PerPage int `json:"per_page,omitempty"`
	Page    int `json:"page,omitempty"`
}

// PaginationOptions structure for providing pagination options for list requests
type PaginationOptions struct {
	PerPage int
	Page    int
}

func paginate(path string, opts *PaginationOptions) (string, error) {
	if opts == nil {
		return path, nil
	}

	u, err := url.Parse(path)
	if err != nil {
		return path, err
	}

	q := u.Query()
	if opts.Page > 0 {
		q.Add("page", strconv.Itoa(opts.Page))
	}
	if opts.PerPage > 0 {
		q.Add("per_page", strconv.Itoa(opts.PerPage))
	}

	u.RawQuery = q.Encode()
	return u.String(), nil
}
