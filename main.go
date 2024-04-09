package main

import (
	"blockchain/cmd"
	"os"
)

func main() {
	defer os.Exit(0)
	cli := cmd.CommandLine{}
	cli.Run()
}
