package exporter

import (
	"context"
	"log"
	"net"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

const (
	DefaultNetworkID = "monero-node_default"
	StratumPort      = 3333
)

type Exporter struct {
	DockerClient    *DockerClient
	DockerNetworkID string
	DockerNetwork   *DockerNetwork
	NetworkDevice   string
	PacketSource    *PacketSource

	knownHosts map[string]Host
	byteCounts map[HostPair]int
}

func (e *Exporter) Run(ctx Context) error {
	if ctx == nil {
		panic("nil context")
	}

	ps, err := e.packetSource(ctx)
	if err != nil {
		return err
	}

	return e.handlePackets(ctx, ps)
}

func (e *Exporter) packetSource(ctx Context) (*PacketSource, error) {
	if e.PacketSource != nil {
		return e.PacketSource, nil
	}

	device, err := e.networkDevice(ctx)
	if err != nil {
		return nil, err
	}

	handle, err := pcap.OpenLive(device, 65535, true, pcap.BlockForever)
	if err != nil {
		return nil, err
	}
	log.Println("listening on:", device)

	if err := handle.SetBPFFilter("tcp"); err != nil {
		log.Println("warning:", err)
	}

	e.PacketSource = gopacket.NewPacketSource(handle, handle.LinkType())
	return e.PacketSource, nil
}

func (e *Exporter) networkDevice(ctx Context) (string, error) {
	if e.NetworkDevice != "" {
		return e.NetworkDevice, nil
	}

	network, err := e.dockerNetwork(ctx)
	if err != nil {
		return "", err
	}

	e.NetworkDevice = "br-" + network.ID[:12]
	return e.NetworkDevice, nil
}

func (e *Exporter) dockerNetwork(ctx Context) (*DockerNetwork, error) {
	if e.DockerNetwork != nil {
		return e.DockerNetwork, nil
	}

	client, err := e.dockerClient(ctx)
	if err != nil {
		return nil, err
	}

	netID := e.DockerNetworkID
	if netID == "" {
		netID = DefaultNetworkID
	}

	opts := types.NetworkInspectOptions{}
	network, err := client.NetworkInspect(ctx, netID, opts)
	if err != nil {
		return nil, err
	}

	e.DockerNetwork = &network
	return e.DockerNetwork, nil
}

func (e *Exporter) dockerClient(ctx Context) (*DockerClient, error) {
	if e.DockerClient != nil {
		return e.DockerClient, nil
	}

	dc, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, err
	}

	e.DockerClient = dc
	return e.DockerClient, nil
}

func (e *Exporter) handlePackets(ctx Context, ps *PacketSource) error {
	e.Reset()

	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)
	go func() {
		cancel(e.exportMetrics(ctx))
	}()

	for {
		err := context.Cause(ctx)
		if err != nil {
			return err
		}

		packet, err := ps.NextPacket()
		if err != nil {
			return err
		}

		err = e.Handle(ctx, packet)
		if err != nil {
			return err
		}
	}
}

func (e *Exporter) Reset() {
	e.byteCounts = make(map[HostPair]int)
}

func (e *Exporter) Handle(ctx Context, packet Packet) error {
	layer := packet.Layer(layers.LayerTypeIPv4)
	if layer == nil {
		return nil
	}
	ip, ok := layer.(*layers.IPv4)
	if !ok {
		log.Printf("gopacket error")
		return UnhandledPacketError(packet)
	}

	layer = packet.Layer(layers.LayerTypeTCP)
	if layer == nil {
		return nil
	}
	tcp, ok := layer.(*layers.TCP)
	if !ok {
		log.Printf("gopacket error")
		return UnhandledPacketError(packet)
	}

	src, err := e.Categorize(ctx, ip.SrcIP, tcp.SrcPort, tcp.DstPort)
	if err != nil {
		return err
	}
	dst, err := e.Categorize(ctx, ip.DstIP, tcp.DstPort, tcp.SrcPort)
	if err != nil {
		return err
	}

	e.byteCounts[PairHosts(src, dst)] += len(packet.Data())
	return nil
}

func (e *Exporter) exportMetrics(ctx Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		}
	}
}

func (e *Exporter) Categorize(ctx Context, ip net.IP,
	port, peerPort layers.TCPPort) (Host, error) {

	result, err := e.categorize(ctx, ip, port, peerPort)
	if err == nil && result == LocalMiner &&
		port != StratumPort && peerPort != StratumPort {
		err = UnhandledAddressError(ip)
	}
	if err != nil && result != NilHost {
		log.Println("warning:", err)
		err = nil
	}
	return result, err
}

func (e *Exporter) categorize(ctx Context, ip net.IP,
	port, peerPort layers.TCPPort) (Host, error) {

	if !ip.IsPrivate() {
		if ip.IsGlobalUnicast() {
			return ExternalHost, nil
		}
		return NilHost, UnhandledAddressError(ip)
	}

	wantCacheKey, err := knownHostsKey(ip)
	if err != nil {
		return NilHost, err
	}

	result, found := e.knownHosts[wantCacheKey]
	if found {
		return result, nil
	}

	log.Printf("%s: unknown host", ip)
	log.Println("scanning Docker network")

	network, err := e.dockerNetwork(ctx)
	if err != nil {
		return NilHost, err
	}

	for _, endpoint := range network.Containers {
		if endpoint.IPv4Address == "" {
			continue
		}

		gotIP, _, err := net.ParseCIDR(endpoint.IPv4Address)
		if err != nil {
			log.Println("warning:", err)
			continue
		}

		gotCacheKey, err := knownHostsKey(gotIP)
		if err != nil {
			log.Println("warning:", err)
			continue
		}

		host, err := HostFromName(endpoint.Name)
		if err != nil {
			log.Println("warning:", err)
			host = UnknownHost
		}

		e.knowHost(gotCacheKey, host)

		if gotCacheKey == wantCacheKey {
			result = host
		}
	}

	if result != NilHost {
		return result, nil
	}
	log.Printf("%s: not found in %s", ip, network.Name)

	if peerPort == StratumPort {
		log.Printf("%s: assuming to be local miner", ip)
		result = LocalMiner
		e.knowHost(wantCacheKey, result)

		return result, nil
	}

	return UnknownHost, UnhandledAddressError(ip)
}

func knownHostsKey(ip net.IP) (string, error) {
	bytes := ip.To4()
	if bytes == nil {
		return "", UnhandledAddressError(ip)
	}
	return string(bytes), nil
}

func (e *Exporter) knowHost(key string, host Host) {
	cachedHost, found := e.knownHosts[key]
	if found && cachedHost == host {
		return
	}

	if e.knownHosts == nil {
		e.knownHosts = make(map[string]Host)
	}

	e.knownHosts[key] = host
	log.Println(net.IP(key), "=>", host)
}
