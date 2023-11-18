package miner

import (
	"os"
	"syscall"
)

// https://pkg.go.dev/os#FindProcess
func isRunning(p *os.Process) bool {
	return p.Signal(syscall.Signal(0)) == nil
}
