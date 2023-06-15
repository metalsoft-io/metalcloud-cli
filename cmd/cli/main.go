package main

// to regenerate the interface and mocking object for the underlying sdk client run 'go generate'. Make sure you have pulled or used go get on the sdk

//go:generate mockgen -source=../metal-cloud-sdk-go/metal_cloud_client.go -destination=helpers/mock_client.go

import (
	"fmt"
	"os"

	"github.com/metalsoft-io/metalcloud-cli/internal/configuration"
	"github.com/metalsoft-io/metalcloud-cli/internal/command"
	"github.com/metalsoft-io/tableformatter"

)

func main() {
	configuration.SetConsoleIOChannel(os.Stdin, os.Stdout)

	clients, err := initClients()
	if err != nil {
		fmt.Fprintf(configuration.GetStdout(), "Could not initialize metal cloud client %s\n", err)
		os.Exit(-1)
	}

	if len(os.Args) < 2 {
		//fmt.Fprintf(GetStdout(), "Invalid command! Use 'help' for a list of commands.\n")
		fmt.Fprintf(configuration.GetStdout(), "%s\n", getHelp(clients, false))
		os.Exit(-1)
	}

	if os.Args[1] == "help" {
		fmt.Fprintf(configuration.GetStdout(), "%s\n", getHelp(clients, false))
		os.Exit(0)
	}

	if len(os.Args) == 1 {
		fmt.Fprint(configuration.GetStdout(), "Invalid command! Use 'help' for a list of commands\n")
		os.Exit(-1)
	}

	tableformatter.DefaultFoldAtLength = 1000

	commands := getCommands(clients)

	err = command.ExecuteCommand(os.Args, commands, clients)

	if err != nil {
		fmt.Fprintf(configuration.GetStdout(), "%s\n", err)
		os.Exit(-2)
	}
}
