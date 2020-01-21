package client

import (
	"fmt"
	"net/url"
	"time"
)

// Project represents Redmine project model.
type Project struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Identifier  string    `json:"identifier"`
	Description string    `json:"description"`
	Status      int       `json:"status"`
	CreatedOn   time.Time `json:"created_on"`
	Parent      *entity   `json:"parent"`
}

type projectsResponse struct {
	Projects []Project `json:"projects"`
}

type projectResponse struct {
	Project Project `json:"project"`
}

// GetProject fetches project with requested ID.
func (c *Client) GetProject(id int64) (*Project, error) {
	req, err := c.getRequest(fmt.Sprintf("/projects/%v.json", id), "")
	if err != nil {
		return nil, err
	}

	var response projectResponse
	_, err = c.Do(req, &response)
	if err != nil {
		return nil, err
	}

	return &response.Project, nil
}

// GetProjects fetches all projects viewable by currently logged user.
func (c *Client) GetProjects() ([]Project, error) {
	req, err := c.getRequest("/projects.json", "limit=200")
	if err != nil {
		return nil, err
	}

	var response projectsResponse
	_, err = c.Do(req, &response)
	if err != nil {
		return nil, err
	}

	return response.Projects, nil
}

// URL returns project URL.
func (p *Project) URL() string {
	hostname, _ := getCredentials()
	u := url.URL{
		Scheme: "https",
		Host:   hostname,
		Path:   fmt.Sprintf("/projects/%v", p.ID),
	}

	return u.String()
}
