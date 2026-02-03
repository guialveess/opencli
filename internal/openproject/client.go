package openproject

import (
	"encoding/json"
	"net/http"
	"time"
)

type Client struct {
	BaseURL string
	Token   string
	Project string
	HTTP    *http.Client
}

func NewClient(baseURL, token, project string) *Client {
	return &Client{
		BaseURL: baseURL,
		Token:   token,
		Project: project,
		HTTP: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c *Client) newRequest(method, path string) (*http.Request, error) {
	req, err := http.NewRequest(method, c.BaseURL+path, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth("apikey", c.Token)
	req.Header.Set("Accept", "application/hal+json")

	return req, nil
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (c *Client) GetCurrentUser() (*User, error) {
	req, err := c.newRequest("GET", "/api/v3/users/me")
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
