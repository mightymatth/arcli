package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/mightymatth/arcli/config"

	"github.com/spf13/viper"
)

type Client struct {
	BaseURL    *url.URL
	UserAgent  string
	httpClient *http.Client
}

var RClient *Client

func init() {
	RClient = &Client{
		BaseURL:    &url.URL{Scheme: "https"},
		UserAgent:  "arcli",
		httpClient: &http.Client{},
	}
}

func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	hostname, apiKey := getCredentials()
	c.BaseURL.Host = hostname

	u := c.BaseURL.ResolveReference(&url.URL{Path: path})
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("X-Redmine-API-Key", apiKey)
	return req, nil
}

func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func getCredentials() (hostname, apiKey string) {
	hostname = viper.GetString(config.Hostname)
	apiKey = viper.GetString(config.ApiKey)

	if hostname == "" || apiKey == "" {
		fmt.Println("You are not logged in.")
		os.Exit(1)
	}

	return
}
