package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/network_device"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	networkDeviceFlags = struct {
		filterStatus     string
		configSource     string
		portId           string
		portStatusAction string
	}{}

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
			return network_device.NetworkDeviceList(cmd.Context(), networkDeviceFlags.filterStatus)
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

	networkDeviceConfigExampleCmd = &cobra.Command{
		Use:          "config-example",
		Short:        "Get network device configuration example.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SWITCHES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceConfigExample(cmd.Context())
		},
	}

	networkDeviceCreateCmd = &cobra.Command{
		Use:          "create",
		Aliases:      []string{"new"},
		Short:        "Create a network device.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SWITCHES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(networkDeviceFlags.configSource)
			if err != nil {
				return err
			}

			return network_device.NetworkDeviceCreate(cmd.Context(), config)
		},
	}

	networkDeviceUpdateCmd = &cobra.Command{
		Use:          "update network_device_id",
		Aliases:      []string{"modify"},
		Short:        "Update a network device.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SWITCHES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(networkDeviceFlags.configSource)
			if err != nil {
				return err
			}

			return network_device.NetworkDeviceUpdate(cmd.Context(), args[0], config)
		},
	}

	networkDeviceDeleteCmd = &cobra.Command{
		Use:          "delete network_device_id",
		Aliases:      []string{"rm"},
		Short:        "Delete a network device.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SWITCHES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceDelete(cmd.Context(), args[0])
		},
	}

	networkDeviceArchiveCmd = &cobra.Command{
		Use:          "archive network_device_id",
		Short:        "Archive a network device.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SWITCHES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceArchive(cmd.Context(), args[0])
		},
	}

	networkDeviceDiscoverCmd = &cobra.Command{
		Use:          "discover network_device_id",
		Short:        "Discover network device interfaces, hardware and software configuration.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SWITCHES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceDiscover(cmd.Context(), args[0])
		},
	}

	networkDeviceGetCredentialsCmd = &cobra.Command{
		Use:          "get-credentials network_device_id",
		Short:        "Get network device credentials.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SWITCHES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceGetCredentials(cmd.Context(), args[0])
		},
	}

	networkDeviceGetPortsCmd = &cobra.Command{
		Use:          "get-ports network_device_id",
		Short:        "Get port statistics for network device directly from the device.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SWITCHES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceGetPorts(cmd.Context(), args[0])
		},
	}

	networkDeviceGetInventoryPortsCmd = &cobra.Command{
		Use:          "get-inventory-ports network_device_id",
		Short:        "Get all ports for network device from the inventory (cached).",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SWITCHES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceGetInventoryPorts(cmd.Context(), args[0])
		},
	}

	networkDeviceSetPortStatusCmd = &cobra.Command{
		Use:          "set-port-status network_device_id",
		Short:        "Set port status for a network device.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SWITCHES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceSetPortStatus(cmd.Context(), args[0], networkDeviceFlags.portId, networkDeviceFlags.portStatusAction)
		},
	}

	networkDeviceResetCmd = &cobra.Command{
		Use:          "reset network_device_id",
		Short:        "Reset a network device to default state and destroy all configurations.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SWITCHES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceReset(cmd.Context(), args[0])
		},
	}

	networkDeviceChangeStatusCmd = &cobra.Command{
		Use:          "change-status network_device_id status",
		Short:        "Change the status of a network device.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SWITCHES_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceChangeStatus(cmd.Context(), args[0], args[1])
		},
	}

	networkDeviceEnableSyslogCmd = &cobra.Command{
		Use:          "enable-syslog network_device_id",
		Short:        "Enable remote syslog for a network device.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SWITCHES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceEnableSyslog(cmd.Context(), args[0])
		},
	}

	networkDeviceGetDefaultsCmd = &cobra.Command{
		Use:          "get-defaults site_label",
		Short:        "Get network device defaults for a site.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SWITCHES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceGetDefaults(cmd.Context(), args[0])
		},
	}

	networkDeviceGetIscsiBootServersCmd = &cobra.Command{
		Use:          "get-iscsi-boot-servers network_device_id",
		Short:        "Get iSCSI boot servers connected through this network device.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SWITCHES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceGetIscsiBootServers(cmd.Context(), args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(networkDeviceCmd)

	networkDeviceCmd.AddCommand(networkDeviceListCmd)
	networkDeviceListCmd.Flags().StringVar(&networkDeviceFlags.filterStatus, "filter-status", "", "Filter the result by network device status.")

	networkDeviceCmd.AddCommand(networkDeviceGetCmd)

	networkDeviceCmd.AddCommand(networkDeviceConfigExampleCmd)

	networkDeviceCmd.AddCommand(networkDeviceCreateCmd)
	networkDeviceCreateCmd.Flags().StringVar(&networkDeviceFlags.configSource, "config-source", "", "Source of the new network device configuration. Can be 'pipe' or path to a JSON file.")
	networkDeviceCreateCmd.MarkFlagsOneRequired("config-source")

	networkDeviceCmd.AddCommand(networkDeviceUpdateCmd)
	networkDeviceUpdateCmd.Flags().StringVar(&networkDeviceFlags.configSource, "config-source", "", "Source of the network device configuration updates. Can be 'pipe' or path to a JSON file.")
	networkDeviceUpdateCmd.MarkFlagsOneRequired("config-source")

	networkDeviceCmd.AddCommand(networkDeviceDeleteCmd)

	networkDeviceCmd.AddCommand(networkDeviceArchiveCmd)

	networkDeviceCmd.AddCommand(networkDeviceDiscoverCmd)

	networkDeviceCmd.AddCommand(networkDeviceGetCredentialsCmd)

	networkDeviceCmd.AddCommand(networkDeviceGetPortsCmd)

	networkDeviceCmd.AddCommand(networkDeviceGetInventoryPortsCmd)

	networkDeviceCmd.AddCommand(networkDeviceSetPortStatusCmd)
	networkDeviceSetPortStatusCmd.Flags().StringVar(&networkDeviceFlags.portId, "port-id", "", "ID of the port to change status.")
	networkDeviceSetPortStatusCmd.Flags().StringVar(&networkDeviceFlags.portStatusAction, "action", "", "Action to perform on the port (up/down).")
	networkDeviceSetPortStatusCmd.MarkFlagsOneRequired("port-id", "action")

	networkDeviceCmd.AddCommand(networkDeviceResetCmd)

	networkDeviceCmd.AddCommand(networkDeviceChangeStatusCmd)

	networkDeviceCmd.AddCommand(networkDeviceEnableSyslogCmd)

	networkDeviceCmd.AddCommand(networkDeviceGetDefaultsCmd)

	networkDeviceCmd.AddCommand(networkDeviceGetIscsiBootServersCmd)
}
