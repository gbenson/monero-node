package exporter

import "fmt"

type Host int

const (
	NilHost       Host = iota * 0
	UnknownHost        = iota + 'u'
	ExternalHost       = iota + 'i'
	LocalMiner         = iota + 'l'
	MoneroNode         = iota + 'm'
	P2PoolNode         = iota + 'p'
	P2PoolTorNode      = iota + 't'
)

var hostStrings = map[Host]string{
	UnknownHost:   "unknown",
	ExternalHost:  "internet",
	LocalMiner:    "local-miner",
	MoneroNode:    "monero",
	P2PoolNode:    "p2pool",
	P2PoolTorNode: "p2pool-tor",
}

var NamedHosts = map[string]Host{
	"monerod":    MoneroNode,
	"p2pool":     P2PoolNode,
	"p2pool-tor": P2PoolTorNode,
}

func HostFromName(hostname string) (Host, error) {
	result, ok := NamedHosts[hostname]
	if ok {
		return result, nil
	}
	return NilHost, UnknownHostError(hostname)
}

func (host Host) String() string {
	result, ok := hostStrings[host]
	if !ok {
		result = fmt.Sprintf("Host(%d)", host)
	}
	return result
}
