package client

import "fmt"

type Issue struct {
	Id          int64  `json:"id"`
	Project     Entity `json:"project"`
	Subject     string `json:"subject"`
	Description string `json:"description"`
}

type IssueResponse struct {
	Issue Issue `json:"issue"`
}

type IssuesResponse struct {
	Issues []Issue `json:"issues"`
}

func (c *Client) GetIssue(id int64) (*Issue, error) {
	req, err := c.getRequest(fmt.Sprintf("/issues/%v.json", id), "")
	if err != nil {
		return nil, err
	}

	var response IssueResponse
	_, err = c.Do(req, &response)
	if err != nil {
		return nil, err
	}

	return &response.Issue, nil
}

func (c *Client) GetMyIssues() ([]Issue, error) {
	return c.GetIssues("assigned_to_id=me")
}

func (c *Client) GetMyWatchedIssues() ([]Issue, error) {
	return c.GetIssues("set_filter=1&sort=updated_on%3Adesc&watcher_id=me")
}

func (c *Client) GetIssues(queryParams string) ([]Issue, error) {
	req, err := c.getRequest("/issues.json", queryParams)
	if err != nil {
		return nil, err
	}

	var response IssuesResponse
	_, err = c.Do(req, &response)
	if err != nil {
		return nil, err
	}

	return response.Issues, nil
}
