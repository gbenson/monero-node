package miner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Report struct {
	receiver    *APIEndpoint
	HostInfo    any          `json:"miner_host,omitempty"`
	MinerStatus any          `json:"miner_status,omitempty"`
	MinerAPI    *APIEndpoint `json:"miner_api,omitempty"`
	Error       any          `json:"error,omitempty"`
}

type GoError struct {
	Message string `json:"message,omitempty"`
	Error   error  `json:"go_error,omitempty"`
	Context any    `json:"context,omitempty"`
}

type HTTPError struct {
	Code        int    `json:"code,omitempty"`
	Message     string `json:"message,omitempty"`
	ContentType string `json:"content_type,omitempty"`
	Body        any    `json:"response_body,omitempty"`
	Context     any    `json:"context,omitempty"`
}

func (r *Report) ReportHTTPError(res *http.Response) error {
	fmt.Println("tor-miner: HTTP error:", res.Status)

	httpError := &HTTPError{
		Code:        res.StatusCode,
		Message:     res.Status,
		ContentType: res.Header.Get("Content-Type"),
		Context:     r.Error,
	}
	r.Error = httpError

	// Trim status code from start of message
	code, msg, ok := strings.Cut(res.Status, " ")
	if ok && code == fmt.Sprint(res.StatusCode) {
		httpError.Message = msg
	}

	// Include the response body if possible
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return r.ReportGoError(err)
	}
	httpError.Body = body

	// Send it
	if err = r.Send(); err == nil {
		return nil
	}

	err = fmt.Errorf("%w [handling: HTTP %v]", err, res.Status)
	fmt.Println("tor-miner:", err)
	return nil
}

func (r *Report) ReportGoError(err error) error {
	fmt.Println("tor-miner: ", err)

	goError := &GoError{
		Message: err.Error(),
		Error:   err,
		Context: r.Error,
	}
	r.Error = goError
	savedError := err

	// Avoid failure if err isn't marshalable.
	if d, err := json.Marshal(err); err != nil || string(d) == "{}" {
		goError.Error = nil
	}

	// Send it
	if err = r.Send(); err == nil {
		return nil
	}

	err = fmt.Errorf("%w [reporting: %w]", err, savedError)
	fmt.Println("tor-miner:", err)
	return nil
}

func (r *Report) Send() error {
	body, err := json.Marshal(r)
	if err != nil {
		return err
	}

	res, err := r.receiver.Post("/recv", "application/json",
		bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		fmt.Printf("tor-miner: %s: %s\n", res.Request.URL, res.Status)
	}
	return nil
}
