package client

import (
	"fmt"
	"net/http"
	"net/url"
)

type SearchItem struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	Url         string `json:"url"`
	Description string `json:"description"`
	DateTime    string `json:"datetime"`
}

type SearchResponse struct {
	SearchItems []SearchItem `json:"results"`
	TotalCount  int          `json:"total_count"`
}

func (c *Client) GetSearchResults(query string, offset, limit int) ([]SearchItem, int, error) {
	req, err := c.getRequest("/search.json",
		fmt.Sprintf("q=%s&offset=%d&limit=%d", url.QueryEscape(query), offset, limit))
	if err != nil {
		return nil, 0, err
	}

	var response SearchResponse
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
