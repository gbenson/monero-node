package miner

import (
	"errors"
	"io"
	"net/http"
	"regexp"
	"time"
)

type Resolver struct {
	URL            string         `json:"url"`
	ResponseParser *regexp.Regexp `json:"response_parser"`
}

func (r *Resolver) GetPoolURL() (string, error) {
	httpClient := http.Client{Timeout: 1 * time.Second}
	res, err := httpClient.Get(r.URL)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", nil
	}

	matches := r.ResponseParser.FindSubmatch(body)
	if len(matches) < 2 {
		return "", errors.New("unexpected response")
	}

	return string(matches[1]), nil
}
