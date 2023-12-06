package exporter

import (
	"fmt"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

const (
	DefaultNetworkID = "monero-node_default"
)

type Exporter struct {
	DockerClient    *DockerClient
	DockerNetworkID string
	DockerNetwork   *DockerNetwork
	NetworkDevice   string
	PacketSource    *PacketSource
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
	for {
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

func (e *Exporter) Handle(ctx Context, packet Packet) error {
	return fmt.Errorf("not implemented")
}
