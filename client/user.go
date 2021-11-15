package client

import (
	"context"
	"net/http"
	"net/url"

	"github.com/mightymatth/arcli/config"

	"github.com/spf13/viper"
)

// User represents user model in Redmine.
type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"login"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"mail"`
	APIKey    string `json:"api_key"`
}

// UserAPIResponse response when user is being fetched.
type UserAPIResponse struct {
	User User `json:"user"`
}

// NewAuthRequest fetches user credentials for given username and password. Method uses
// simple basic authentication.
func (c *Client) NewAuthRequest(ctx context.Context, username, password string) (*http.Request, error) {
	u, err := url.Parse(viper.GetString(config.Host))
	if err != nil {
		return nil, err
	}

	u.Path = "/users/current.json"
	u.User = url.UserPassword(username, password)

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	return req, nil
}

// GetUser fetches data of currently logged user.
func (c *Client) GetUser() (*User, error) {
	req, err := c.getRequest("/users/current.json", "")
	if err != nil {
		return nil, err
	}
	var userResponse UserAPIResponse
	_, err = c.Do(req, &userResponse)
	if err != nil {
		return nil, err
	}
	return &(userResponse.User), nil
}
