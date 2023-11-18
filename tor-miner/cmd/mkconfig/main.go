package main

import (
	"fmt"
	"os"

	"gbenson.net/tor-miner"
)

func main() {
	err := _main()
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

func _main() error {
	if len(os.Args) != 4 {
		fmt.Println("usage: mkconfig PASSPHRASE POOL_URL MONITOR_URL")
		os.Exit(2)
	}

	passphrase := os.Args[1]
	poolURL := os.Args[2]
	monitorURL := os.Args[3]

	config := miner.Config{
		Pool: miner.APIEndpoint{
			URL: poolURL,
		},
		Monitor: miner.APIEndpoint{
			URL: monitorURL,
		},
	}

	sealed, err := config.Seal(passphrase)
	if err != nil {
		return err
	}

	const filename = "sealed.config"
	if err = os.WriteFile(filename, sealed, 0666); err != nil {
		return err
	}

	fmt.Printf("Wrote %s (%d bytes)\n", filename, len(sealed))
	return nil
}
