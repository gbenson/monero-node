package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"gbenson.net/monero-node/node-exporter"
)

func main() {
	log.SetFlags(log.Lshortfile)

	if len(os.Args) > 2 {
		fmt.Println("usage: node-exporter [DOCKER_NETWORK]")
		os.Exit(2)
	}

	e := &exporter.Exporter{}
	if len(os.Args) == 2 {
		e.DockerNetworkID = os.Args[1]
	}

	err := e.Run(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}
