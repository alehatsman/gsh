package client

// Repo - represents github repository.
type Repo struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"html_url"`
	Stars       int    `json:"stargazers_count"`
}

// SearchResult - represents search response.
type SearchResult struct {
	TotalCount int    `json:"total_count"`
	Items      []Repo `json:"items"`
}
