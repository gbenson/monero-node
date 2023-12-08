package exporter

import "fmt"

type HostPair Host

const (
	dstBits = 16 << (^uint(0) >> 63) // 16 or 32
	dstMask = (1 << dstBits) - 1
)

func PairHosts(src, dst Host) HostPair {
	return HostPair((src << dstBits) | (dst & dstMask))
}

func (p HostPair) Src() Host {
	return Host(p >> dstBits)
}

func (p HostPair) Dst() Host {
	return Host(p & dstMask)
}

func (p HostPair) String() string {
	return fmt.Sprintf("%s-%s", p.Src(), p.Dst())
}
