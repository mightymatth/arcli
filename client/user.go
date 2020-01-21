package client

import (
	"context"
	"net/http"
	"net/url"

	"github.com/mightymatth/arcli/config"

	"github.com/spf13/viper"
)

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"login"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"mail"`
	APIKey    string `json:"api_key"`
}

type UserAPIResponse struct {
	User User `json:"user"`
}

func (c *Client) NewAuthRequest(ctx context.Context, username, password string) (*http.Request, error) {
	u := url.URL{
		Scheme: "https",
		Host:   viper.GetString(config.Hostname),
		Path:   "/users/current.json",
		User:   url.UserPassword(username, password),
	}

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	return req, nil
}

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
