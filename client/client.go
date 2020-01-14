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
	HttpClient *http.Client
	UserAgent  string
}

func (c *Client) getRequest(path string, queryParams string) (*http.Request, error) {
	hostname, apiKey := getCredentials()
	u := url.URL{
		Scheme:   "https",
		Host:     hostname,
		Path:     path,
		RawQuery: queryParams,
	}

	var buf io.ReadWriter
	req, err := http.NewRequest("GET", u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("X-Redmine-API-Key", apiKey)

	return req, nil
}

func (c *Client) postRequest(path string, body interface{}) (*http.Request, error) {
	hostname, apiKey := getCredentials()
	u := url.URL{
		Scheme: "https",
		Host:   hostname,
		Path:   path,
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest("POST", u.String(), buf)
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

func (c *Client) deleteRequest(path string) (*http.Request, error) {
	hostname, apiKey := getCredentials()
	u := url.URL{
		Scheme: "https",
		Host:   hostname,
		Path:   path,
	}

	var buf io.ReadWriter
	req, err := http.NewRequest("DELETE", u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("X-Redmine-API-Key", apiKey)

	return req, nil
}

func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.HttpClient.Do(req)
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

type Entity struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type EntityId struct {
	Id int64 `json:"id"`
}

func (e EntityId) String() string {
	switch e.Id {
	case 0:
		return "-"
	default:
		return fmt.Sprintf("%v", e.Id)
	}
}
