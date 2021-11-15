package client

import (
	"fmt"
	"net/url"

	"github.com/mightymatth/arcli/config"

	"github.com/spf13/viper"
)

// Issue correspond with issue in Redmine.
type Issue struct {
	ID          int64  `json:"id"`
	Project     entity `json:"project"`
	Subject     string `json:"subject"`
	Description string `json:"description"`
}

type issueResponse struct {
	Issue Issue `json:"issue"`
}

type issuesResponse struct {
	Issues []Issue `json:"issues"`
}

// GetIssue fetches issue with requested ID.
func (c *Client) GetIssue(id int64) (*Issue, error) {
	req, err := c.getRequest(fmt.Sprintf("/issues/%v.json", id), "")
	if err != nil {
		return nil, err
	}

	var response issueResponse
	_, err = c.Do(req, &response)
	if err != nil {
		return nil, err
	}

	return &response.Issue, nil
}

// GetMyIssues fetches issues assigned only to currently logged user.
func (c *Client) GetMyIssues() ([]Issue, error) {
	params := fmt.Sprintf("assigned_to_id=%v", viper.GetString(config.UserID))
	return c.GetIssues(params)
}

// GetMyRelatedIssues fetches issues assigned to currently logged user.
func (c *Client) GetMyRelatedIssues() ([]Issue, error) {
	return c.GetIssues("assigned_to_id=me")
}

// GetMyWatchedIssues fetches issues that currently logged user watches.
func (c *Client) GetMyWatchedIssues() ([]Issue, error) {
	return c.GetIssues("set_filter=1&sort=updated_on%3Adesc&watcher_id=me")
}

// GetIssues fetches issues with rules defined in queryParams.
func (c *Client) GetIssues(queryParams string) ([]Issue, error) {
	req, err := c.getRequest("/issues.json", queryParams)
	if err != nil {
		return nil, err
	}

	var response issuesResponse
	_, err = c.Do(req, &response)
	if err != nil {
		return nil, err
	}

	return response.Issues, nil
}

// URL returns issue URL.
func (i *Issue) URL() string {
	host, _ := getCredentials()

	u, err := url.Parse(host)
	if err != nil {
		return ""
	}

	u.Path = fmt.Sprintf("/issues/%v", i.ID)

	return u.String()
}
