package miner

import (
	"context"
	"crypto/rand"
	"encoding/base32"
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
	TorProxyAddr      = "127.0.0.1:9050"
	TorStartupTimeout = time.Minute
	DefaultAPIAddr    = "127.0.0.1:3638"
)

var httpAPI = regexp.MustCompile(`HTTP API\s+(\S+:\d+)\n`)

type Runner struct {
	MinerPath string
	MinerArgs []string

	LocalAddr   string
	OnionAddr   string
	AccessToken string

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
	if r.LocalAddr == "" {
		r.LocalAddr = DefaultAPIAddr
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
		r.LocalAddr = m[1]
	}
	host, port, err := net.SplitHostPort(r.LocalAddr)
	if err != nil {
		os.Stdout.WriteString(out)
		return err
	}
	if !apiConfigured {
		opts := []string{"--http-host", host, "--http-port", port}
		r.MinerArgs = append(r.MinerArgs, opts...)

		if r.AccessToken == "" {
			var buf [32]byte
			_, err := rand.Read(buf[:])
			if err != nil {
				return err
			}

			r.AccessToken = strings.ToLower(strings.TrimRight(
				base32.StdEncoding.EncodeToString(buf[:]), "="))
		}

		opts = []string{"--http-access-token", r.AccessToken}
		r.MinerArgs = append(r.MinerArgs, opts...)
	}

	// Start Tor
	for start := time.Now(); !isTorRunning(); {
		if r.tor == nil {
			torctx, cancel := context.WithCancel(ctx)
			defer cancel()

			r.tor = r.newProxyCommand(torctx)
			defer reap(r.tor, cancel)

			fmt.Println("tor-miner: starting Tor")
			if err = r.startProxy(); err != nil {
				return err
			}
		}
		if time.Now().Sub(start) > TorStartupTimeout {
			fmt.Println("tor-miner: timed out waiting for Tor")
			break // maybe it'll work? let the monitor handle it
		}
		time.Sleep(100 * time.Millisecond)

		if !isRunning(r.tor.Process) {
			fmt.Println("tor-miner: Tor startup failed")
			break // as above, let the monitor figure this out
		}

	}

	// Enable full remote access to XMRig's API if we have a hidden
	// service URL to expose it over.  Having this a soft requirement
	// enables out-of-container testing by non-root users.
	bytes, err := os.ReadFile("/var/lib/tor/hidden_service/hostname")
	if err != nil {
		fmt.Println("tor-miner:", err)
	} else {
		s := strings.TrimSpace(string(bytes))
		if strings.HasSuffix(s, ".onion") {
			r.OnionAddr = s + ":" + port
		}
	}
	if !apiConfigured && r.OnionAddr != "" && r.AccessToken != "" {
		r.MinerArgs = append(r.MinerArgs, "--http-no-restricted")
	}

	// Start XMRig
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	r.cmd = r.newMinerCommand(ctx)
	defer reap(r.cmd, cancel)

	r.cmd.Stdout = os.Stdout
	r.cmd.Stderr = os.Stdout
	fmt.Println("tor-miner: starting XMrig")
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
	return exec.CommandContext(ctx, path.Join(dir, "start-tor"), r.LocalAddr)
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
