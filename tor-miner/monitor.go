package miner

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"time"
)

type Monitor struct {
	cmds     []*exec.Cmd
	localAPI *APIEndpoint
	onionAPI *APIEndpoint
	receiver *APIEndpoint
	hostInfo *InstanceMetadata
}

func monitor(cmds []*exec.Cmd,
	localAPI, onionAPI, receiver *APIEndpoint) error {

	m := Monitor{
		cmds:     cmds,
		localAPI: localAPI,
		onionAPI: onionAPI,
		receiver: receiver,
	}

	time.Sleep(1 * time.Second)
	for {
		err := m.mainLoop()
		if err != nil {
			return err
		}
		time.Sleep(1 * time.Minute)
	}
}

func (m *Monitor) mainLoop() error {
	r := Report{
		receiver: m.receiver,
		HostInfo: m.hostInfo,
		MinerAPI: m.onionAPI,
	}

	// Ensure our subprocesses are still running.  Failure of either
	// is treated as unrecoverable: we try to report the error and
	// then terminate.  Recovery is our invoker's problem.
	for _, cmd := range m.cmds {
		if isRunning(cmd.Process) {
			continue
		}

		err := fmt.Errorf("%v exited", cmd)
		r.ReportGoError(err)
		return err
	}

	// Get the miner's status
	res, err := m.localAPI.Get("/2/summary")
	if err != nil {
		return r.ReportGoError(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return r.ReportHTTPError(res)
	}

	ctype := res.Header.Get("Content-Type")
	switch ctype {
	case "application/json":
		err = nil
	case "":
		err = errors.New("unspecified Content-Type")
	default:
		err = fmt.Errorf("%q: unexpected Content-Type", ctype)
	}
	if err != nil {
		return r.ReportGoError(err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return r.ReportGoError(err)
	}

	if err = json.Unmarshal(body, &r.MinerStatus); err != nil {
		r.MinerStatus = body
		return r.ReportGoError(err)
	}

	if err = r.Send(); err != nil {
		return r.ReportGoError(err)
	}

	return nil
}
