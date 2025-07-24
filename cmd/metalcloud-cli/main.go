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

	// Print proxy environment variables
	envVars := []string{
		"HTTP_PROXY",
		"HTTPS_PROXY",
		"NO_PROXY",
		"http_proxy",
		"https_proxy",
		"no_proxy",
	}

	for _, envVar := range envVars {
		value := os.Getenv(envVar)
		if value != "" {
			fmt.Fprintf(os.Stderr, "%s=%s\n", envVar, value)
		}
	}

	err := cmd.Execute()
	if err != nil {
		os.Exit(-1)
	}
}
