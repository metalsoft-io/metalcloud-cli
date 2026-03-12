package cmd

import (
	"github.com/spf13/cobra"
)

var networkConfigurationCmd = &cobra.Command{
	Use:     "network-configuration [command]",
	Aliases: []string{"net-config", "nc"},
	Short:   "Manage network configuration templates",
	Long: `Network configuration commands.

Manage network device configuration templates used to deploy configurations to network devices.
Available subcommands:
  device-template              Manage network device configuration templates
  link-aggregation-template    Manage network device link aggregation configuration templates`,
}

func init() {
	rootCmd.AddCommand(networkConfigurationCmd)
}
