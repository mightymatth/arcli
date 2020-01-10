package client

import (
	"net/http"
	"net/url"

	"github.com/mightymatth/arcli/config"

	"github.com/spf13/viper"
)

type User struct {
	Id        int64  `json:"id"`
	Username  string `json:"login"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"mail"`
	ApiKey    string `json:"api_key"`
}

type UserApiResponse struct {
	User User `json:"user"`
}

func (c *Client) NewAuthRequest(username, password string) (*http.Request, error) {
	c.BaseURL.Host = viper.GetString(config.Hostname)
	u := c.BaseURL.ResolveReference(&url.URL{Path: "/users/current.json"})
	u.User = url.UserPassword(username, password)

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	return req, nil
}

func (c *Client) GetUser() (*User, error) {
	req, err := c.newRequest("GET", "/users/current.json", nil)
	if err != nil {
		return nil, err
	}
	var userResponse UserApiResponse
	_, err = c.Do(req, &userResponse)
	if err != nil {
		return nil, err
	}
	return &(userResponse.User), nil
}
