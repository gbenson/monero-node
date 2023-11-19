package miner

import (
	"context"
	"encoding/base32"
	"fmt"
	"net"
	net_url "net/url"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"time"
)

const (
	DefaultMinerPath  = "xmrig"
	TorProxyAddr      = "127.0.0.1:9050"
	TorStartupTimeout = time.Minute
	DefaultAPIAddr    = "127.0.0.1:3638"
)

var httpAPI = regexp.MustCompile(`HTTP API\s+(\S+:\d+)\n`)

type Runner struct {
	MinerPath string
	MinerArgs []string

	LocalAPI APIEndpoint
	onionAPI *APIEndpoint

	isStarted bool
	cmd       *exec.Cmd
	tor       *exec.Cmd
}

func (r *Runner) Run(ctx context.Context) error {
	if len(os.Args) < 2 {
		return panique("invocation error")
	}
	if r.isStarted {
		return panique("already started")
	}
	r.isStarted = true

	config, err := DefaultConfig(os.Args[1])
	if err != nil {
		return err
	}

	if r.MinerPath == "" {
		r.MinerPath = DefaultMinerPath
	}
	if r.MinerArgs == nil {
		r.MinerArgs = os.Args[2:]
	}
	if r.LocalAPI.URL == "" {
		r.LocalAPI.URL = "http://" + DefaultAPIAddr
	}

	// Capture the --dry-run output, supplying the default
	// configuration file if necessary.
	out, err := r.dryRun(ctx)
	if strings.Contains(out, "no valid configuration found") {
		opt := []string{"-o", config.Pool.URL}
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
		r.LocalAPI.URL = "http://" + m[1]
	}
	url, err := net_url.Parse(r.LocalAPI.URL)
	if err != nil {
		os.Stdout.WriteString(out)
		return err
	}
	port := url.Port()
	if !apiConfigured {
		host := url.Hostname()
		opts := []string{"--http-host", host, "--http-port", port}
		r.MinerArgs = append(r.MinerArgs, opts...)

		if r.LocalAPI.AccessToken == "" {
			buf, err := randomBytes(32)
			if err != nil {
				return err
			}

			r.LocalAPI.AccessToken = strings.ToLower(
				strings.TrimRight(
					base32.StdEncoding.EncodeToString(buf), "="))
		}

		opts = []string{"--http-access-token", r.LocalAPI.AccessToken}
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
			r.onionAPI = &APIEndpoint{
				URL:         fmt.Sprintf("http://%s:%s", s, port),
				AccessToken: r.LocalAPI.AccessToken,
			}
		}
	}
	if !apiConfigured && r.onionAPI != nil && r.onionAPI.AccessToken != "" {
		r.MinerArgs = append(r.MinerArgs, "--http-no-restricted")
	}

	// Start XMRig
	rigctx, cancel := context.WithCancel(ctx)
	defer cancel()

	r.cmd = r.newMinerCommand(rigctx)
	defer reap(r.cmd, cancel)

	r.cmd.Stdout = os.Stdout
	r.cmd.Stderr = os.Stdout
	fmt.Println("tor-miner: starting XMrig")
	if err = r.cmd.Start(); err != nil {
		return err
	}

	// Begin monitoring
	cmds := []*exec.Cmd{r.cmd}
	if r.tor != nil {
		cmds = append(cmds, r.tor)
	}
	return monitor(cmds, &r.LocalAPI, r.onionAPI, &config.Monitor)
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
	url, err := net_url.Parse(r.LocalAPI.URL)
	if err != nil {
		panic(err)
	}
	return exec.CommandContext(ctx, path.Join(dir, "start-tor"), url.Host)
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
