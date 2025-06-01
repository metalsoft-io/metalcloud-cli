package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/cmd"
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
)

var version string
var allowDevelop string

func main() {
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChannel
		os.Exit(-1)
	}()

	cmd.Version = version
	system.AllowDevelop = allowDevelop == "true" || allowDevelop == "yes"

	err := cmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(-1)
	}
}
