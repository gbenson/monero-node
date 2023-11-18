package miner

type APIEndpoint struct {
	URL         string `json:"url"`
	AccessToken string `json:"access_token,omitempty"`
}
