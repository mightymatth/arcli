package client

import (
	"fmt"
	"net/http"
	"net/url"
)

// SearchItem represents Redmine search item model.
type SearchItem struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	URL         string `json:"url"`
	Description string `json:"description"`
	DateTime    string `json:"datetime"`
}

type searchResponse struct {
	SearchItems []SearchItem `json:"results"`
	TotalCount  int          `json:"total_count"`
}

// GetSearchResults returns search results for given query, offset and limit.
func (c *Client) GetSearchResults(query string, offset, limit int) ([]SearchItem, int, error) {
	req, err := c.getRequest("/search.json",
		fmt.Sprintf("q=%s&offset=%d&limit=%d", url.QueryEscape(query), offset, limit))
	if err != nil {
		return nil, 0, err
	}

	var response searchResponse
	res, err := c.Do(req, &response)
	if err != nil {
		return nil, 0, err
	}

	switch res.StatusCode {
	case http.StatusOK:
		return response.SearchItems, response.TotalCount, nil
	default:
		return nil, 0, fmt.Errorf("could not get search item (status %v)", res.StatusCode)
	}
}
