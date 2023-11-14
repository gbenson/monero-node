package miner

import "fmt"

// Panique is like panic but softer and more French.
func panique(msg any) error {
	return fmt.Errorf("panique: %v", msg)
}
