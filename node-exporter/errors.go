package exporter

import (
	"errors"
	"fmt"
	"log"
	"net"
)

var (
	ErrUnhandledPacket = errors.New("unhandled packet")
)

func UnhandledAddressError(ip net.IP) error {
	return fmt.Errorf("%s: unhandled address", ip)
}

func UnhandledPacketError(packet Packet) error {
	log.Output(2, packet.Dump())
	return ErrUnhandledPacket
}
