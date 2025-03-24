package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/network_device"
	"github.com/spf13/cobra"
)

var (
	networkDeviceCmd = &cobra.Command{
		Use:     "network-device [command]",
		Aliases: []string{"switch", "nd"},
		Short:   "Network device management",
		Long:    `Network device commands.`,
	}

	networkDeviceListCmd = &cobra.Command{
		Use:          "list",
		Aliases:      []string{"ls"},
		Short:        "List all network devices.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SWITCHES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceList(cmd.Context())
		},
	}

	networkDeviceGetCmd = &cobra.Command{
		Use:          "get network_device_id",
		Aliases:      []string{"show"},
		Short:        "Get network device details.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SWITCHES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceGet(cmd.Context(), args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(networkDeviceCmd)

	networkDeviceCmd.AddCommand(networkDeviceListCmd)

	networkDeviceCmd.AddCommand(networkDeviceGetCmd)
}
