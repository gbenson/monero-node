package miner

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"time"
)

const (
	DefaultMinerPath = "xmrig"
	DefaultPool      = "ptkdyo72ibo5edkviouk5w5oct" +
		"xk5d7szizdlgxepfygckeiyt7cdiqd.onion:3333"
	TorProxyAddr   = "127.0.0.1:9050"
	DefaultAPIAddr = "127.0.0.1:3638"
)

var httpAPI = regexp.MustCompile(`HTTP API\s+(\S+:\d+)\n`)

type Runner struct {
	MinerPath string
	MinerArgs []string
	LocalAPI  string
	OnionAPI  string
	isStarted bool
	cmd       *exec.Cmd
	tor       *exec.Cmd
}

func (r *Runner) Run(ctx context.Context) error {
	if r.isStarted {
		return panique("already started")
	}
	r.isStarted = true

	if r.MinerPath == "" {
		r.MinerPath = DefaultMinerPath
	}
	if r.MinerArgs == nil {
		r.MinerArgs = os.Args[1:]
	}
	if r.LocalAPI == "" {
		r.LocalAPI = DefaultAPIAddr
	}

	// Capture the --dry-run output, supplying the default
	// configuration file if necessary.
	out, err := r.dryRun(ctx)
	if strings.Contains(out, "no valid configuration found") {
		opt := []string{"-o", DefaultPool}
		r.MinerArgs = append(opt, r.MinerArgs...)
		out, err = r.dryRun(ctx)
	}
	out = decolor(out)
	if err != nil {
		os.Stdout.WriteString(out)
		return err
	}

	// Specify routing over Tor, if necessary
	if strings.Contains(out, ".onion:") {
		opt := []string{"-x", TorProxyAddr}
		r.MinerArgs = append(opt, r.MinerArgs...)
	}

	// Discover or configure our HTTP API details
	m := httpAPI.FindStringSubmatch(out)
	apiConfigured := len(m) > 1
	if apiConfigured {
		r.LocalAPI = m[1]
	}
	host, port, err := net.SplitHostPort(r.LocalAPI)
	if err != nil {
		os.Stdout.WriteString(out)
		return err
	}
	if !apiConfigured {
		opts := []string{"--http-host", host, "--http-port", port}
		r.MinerArgs = append(r.MinerArgs, opts...)
	}

	// Start Tor
	torStarted := false
	for i := 0; i < 600 && !isTorRunning(); i++ {
		if torStarted {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		torctx, cancel := context.WithCancel(ctx)
		defer cancel()

		r.tor = r.newProxyCommand(torctx)
		defer reap(r.tor, cancel)

		if err = r.startProxy(); err != nil {
			return err
		}

		torStarted = true
	}

	// Enable full remote access to XMRig's API if we can determine
	// the hidden service URL we'll expose it over.  Making this a
	// soft failure enables out-of-container non-root-user testing.
	// XXX

	// Start XMRig
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	r.cmd = r.newMinerCommand(ctx)
	defer reap(r.cmd, cancel)

	r.cmd.Stdout = os.Stdout
	r.cmd.Stderr = os.Stdout
	if err = r.cmd.Start(); err != nil {
		return err
	}

	// XXX http api status reporting
	// XXX kill xmrig if tor dies
	return r.cmd.Wait()
}

func (r *Runner) dryRun(ctx context.Context) (string, error) {
	cmd := r.newMinerCommand(ctx, "--dry-run")
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func (r *Runner) newMinerCommand(ctx context.Context, arg ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, r.MinerPath, append(arg, r.MinerArgs...)...)
	r.MinerPath = cmd.Path
	return cmd
}

func (r *Runner) newProxyCommand(ctx context.Context) *exec.Cmd {
	dir := path.Dir(r.MinerPath)
	if path.Base(dir) == "bin" {
		dir = path.Join(path.Dir(dir), "libexec", "tor-miner")
	}
	return exec.CommandContext(ctx, path.Join(dir, "start-tor"), r.LocalAPI)
}

func (r *Runner) startProxy() error {
	r.tor.Stdout = os.Stderr
	r.tor.Stderr = os.Stderr

	return r.tor.Start()
}

func isTorRunning() bool {
	conn, err := net.Dial("tcp", TorProxyAddr)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func reap(cmd *exec.Cmd, cancel context.CancelFunc) {
	cancel()
	if err := cmd.Wait(); err != nil {
		fmt.Printf("%s: %s\n", cmd.Path, err)
	}
	fmt.Printf("%s: %s\n", cmd.Path, "reaped")
}
