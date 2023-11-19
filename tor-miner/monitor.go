package miner

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

type Monitor struct {
	cmds     []*exec.Cmd
	localAPI *APIEndpoint
	onionAPI *APIEndpoint
}

func monitor(ctx context.Context, cmds []*exec.Cmd,
	localAPI, onionAPI *APIEndpoint) error {

	m := Monitor{
		cmds:     cmds,
		localAPI: localAPI,
		onionAPI: onionAPI,
	}

	time.Sleep(1 * time.Second)
	for {
		err := m.mainLoop(ctx)
		if err != nil {
			return err
		}
		time.Sleep(1 * time.Minute)
	}
}

func (m *Monitor) mainLoop(ctx context.Context) error {
	// Ensure our subprocesses are still running.  Failure of either
	// is treated as unrecoverable: we try to report the error and
	// then terminate.  Recovery is our invoker's problem.
	for _, cmd := range m.cmds {
		if isRunning(cmd.Process) {
			continue
		}

		err := fmt.Errorf("%v exited", cmd)
		m.reportGoError(ctx, err)
		return err
	}

	// Get the miner's status
	res, err := m.localAPI.Get("/2/summary")
	if err != nil {
		return m.reportGoError(ctx, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return m.reportHTTPError(ctx, res)
	}

	if ct := res.Header.Get("Content-Type"); ct != "application/json" {
		err = fmt.Errorf("unhandled Content-Type %q", ct)
		return m.reportGoError(ctx, err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return m.reportGoError(ctx, err)
	}

	var report any
	err = json.Unmarshal(body, &report)
	if err != nil {
		return m.reportGoError(ctx, err)
	}

	return m.reportStatus(ctx, report)
}

func (m *Monitor) reportStatus(ctx context.Context, rr any) error {
	// Insert the connection details
	if m.onionAPI != nil {
		switch r := rr.(type) {
		case map[string]interface{}:
			if _, ok := r["status"]; ok {
				r["status"] = "OK"
			}
			r["http_api"] = m.onionAPI

		default:
			err := fmt.Errorf("unhandled type %T", rr)
			return m.reportGoError(ctx, err)
		}
	}

	// Marshal the report body
	body, err := json.Marshal(rr)
	if err != nil {
		return m.reportGoError(ctx, err)
	}

	fmt.Println(string(body))
	return panique("not implemented")
}

func (m *Monitor) reportHTTPError(ctx context.Context,
	res *http.Response) error {
	fmt.Println("tor-miner: HTTP error:", res.Status, res.Body)

	report := map[string]any{
		"status":  "error",
		"type":    "http",
		"code":    res.StatusCode,
		"message": res.Status,
	}

	// Trim res.StatusCode from start of report["message"]
	code, msg, ok := strings.Cut(res.Status, " ")
	if ok && code == fmt.Sprint(res.StatusCode) {
		report["message"] = msg
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		err = fmt.Errorf("error: %w [handling: HTTP %v]", err, res.Status)
		return m.reportGoError(ctx, err)
	}
	report["body"] = string(body)

	return m.reportStatus(ctx, report)
}

func (m *Monitor) reportGoError(ctx context.Context, err error) error {
	fmt.Println("tor-miner:", err)

	report := map[string]any{
		"status":  "error",
		"type":    "golang",
		"message": err.Error(),
	}

	savedError := err
	if err = m.reportStatus(ctx, report); err == nil {
		return nil
	}
	err = fmt.Errorf("error: %w [handling: %w]", err, savedError)

	fmt.Println("tor-miner:", err)
	return nil
}
