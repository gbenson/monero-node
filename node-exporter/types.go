package exporter

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"

	"github.com/google/gopacket"
)

type (
	Context = context.Context

	DockerClient  = client.Client
	DockerNetwork = types.NetworkResource

	Packet       = gopacket.Packet
	PacketSource = gopacket.PacketSource
)
