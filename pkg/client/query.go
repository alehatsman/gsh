package client

import (
	"net/url"
	"strconv"
)

type SearchQuery struct {
	Query   string
	Sort    string
	Order   string
	PerPage int
}

func (sq *SearchQuery) Encode() string {
	values := url.Values{}
	values.Add("q", sq.Query)
	values.Add("sort", sq.Sort)
	values.Add("order", sq.Order)
	values.Add("per_page", strconv.Itoa(sq.PerPage))
	return values.Encode()
}
