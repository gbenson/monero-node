package miner

import (
	"fmt"
	"io"
	"net/http"
	net_url "net/url"
	"strings"
	"time"
)

var httpClient = http.Client{Timeout: 30 * time.Second}

type APIEndpoint struct {
	URL         string `json:"url"`
	AccessToken string `json:"access_token,omitempty"`
}

// Analogue of http.Client.Get
func (api *APIEndpoint) Get(route string) (*http.Response, error) {

	req, err := api.NewRequest("GET", route, nil)
	if err != nil {
		return nil, err
	}
	return api.Do(req)
}

// Analogue of http.Client.Post
func (api *APIEndpoint) Post(route, contentType string,
	body io.Reader) (*http.Response, error) {

	req, err := api.NewRequest("POST", route, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return api.Do(req)
}

// Analogue of http.NewRequest
func (api *APIEndpoint) NewRequest(
	method, route string,
	body io.Reader) (*http.Request, error) {

	url, err := net_url.Parse(api.URL)
	if err != nil {
		return nil, err
	}
	url = url.JoinPath(strings.TrimLeft(route, "/"))

	req, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		return nil, err
	}
	if api.AccessToken != "" {
		req.Header.Set("Authorization", "Bearer "+api.AccessToken)
	}
	return req, err
}

// Analogue of http.Client.Do
func (api *APIEndpoint) Do(req *http.Request) (*http.Response, error) {
	fmt.Printf("tor-miner: %s %s\n", req.Method, req.URL)
	return httpClient.Do(req)
}
