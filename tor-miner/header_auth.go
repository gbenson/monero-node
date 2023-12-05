package miner

import (
	"fmt"
	"net/http"
)

type Authenticator struct {
	Name   string
	Format string
}

var AuthStyleBearer = &Authenticator{"Authorization", "Bearer %s"}

func (h *Authenticator) Authenticate(req *http.Request, token string) {
	if token == "" {
		return
	}
	if h.Format != "" {
		token = fmt.Sprintf(h.Format, token)
	}
	req.Header.Set(h.Name, token)
}
