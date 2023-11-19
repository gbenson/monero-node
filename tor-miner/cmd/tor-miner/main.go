package main

import (
	"context"
	"fmt"
	"os"

	"gbenson.net/tor-miner"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: tor-miner CONFIG_PASSPHRASE [XMRIG_ARGS...]")
		os.Exit(2)
	}

	app := miner.Runner{}
	if err := app.Run(context.Background()); err != nil {
		fmt.Println("tor-miner:", err)
		os.Exit(1)
	}
}
