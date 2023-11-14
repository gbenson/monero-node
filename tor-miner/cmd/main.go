package main

import (
	"fmt"
	"os"

	"gbenson.net/tor-miner"
)

func main() {
	app := miner.Runner{}
	if err := app.Run(); err != nil {
		fmt.Println("tor-miner:", err)
		os.Exit(1)
	}
}
