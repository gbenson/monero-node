package miner

import (
	"fmt"
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
	// Ensure our subprocesses are still running.  Failure of either
	// is treated as unrecoverable: we try to report the error and
	// then terminate.  Recovery is our invoker's problem.
	for _, cmd := range m.cmds {
		if isRunning(cmd.Process) {
			continue
		}

		err := fmt.Errorf("%v exited", cmd)
		return err
	}

	return nil
}
