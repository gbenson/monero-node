package miner

import "regexp"

var ansi = regexp.MustCompile(`\x1B\[\d+(?:;\d+)*m`)

func decolor(s string) string {
	return ansi.ReplaceAllLiteralString(s, "")
}
