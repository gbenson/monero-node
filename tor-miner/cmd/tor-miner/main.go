package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"gbenson.net/tor-miner"
)

func main() {
	app := miner.Runner{}
	if err := app.Run(context.Background()); err != nil {
		if errors.Is(err, miner.UsageError) {
			fmt.Println(err)
			os.Exit(2)
		}

		fmt.Println("tor-miner:", err)
		os.Exit(1)
	}
}
