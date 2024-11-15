package main

import (
	"buildenv/command/cli"
	"flag"
	"os"
)

func main() {
	if exit := cli.Listen(); exit {
		os.Exit(0)
	}

	flag.Usage()
}
