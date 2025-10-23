package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/network_device"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	networkDeviceFlags = struct {
		filterStatus     []string
		configSource     string
		portId           string
		portStatusAction string
	}{}

	networkDeviceCmd = &cobra.Command{
		Use:     "network-device [command]",
		Aliases: []string{"switch", "nd"},
		Short:   "Manage network devices (switches) in the infrastructure",
		Long: `Network device management commands for switches and other network infrastructure.

Network devices are physical switches that connect servers and provide network connectivity
within the MetalSoft infrastructure. These commands allow you to manage, configure, and
monitor network devices.`,
	}

	networkDeviceListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List network devices with optional status filtering",
		Long: `List all network devices in the infrastructure with optional status filtering.

This command displays all network devices (switches) that are registered in the system.
You can filter the results by device status to focus on specific operational states.

Flags:
  --filter-status   Filter devices by operational status (default: ["active"])
                   Available statuses: active, inactive, maintenance, error, unknown

Examples:
  # List all active network devices (default)
  metalcloud-cli network-device list

  # List devices in maintenance mode
  metalcloud-cli network-device list --filter-status maintenance

  # List devices with multiple statuses
  metalcloud-cli network-device list --filter-status active,maintenance

  # List all devices regardless of status
  metalcloud-cli network-device list --filter-status active,inactive,maintenance,error,unknown`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceList(cmd.Context(), networkDeviceFlags.filterStatus)
		},
	}

	networkDeviceGetCmd = &cobra.Command{
		Use:     "get <network_device_id>",
		Aliases: []string{"show"},
		Short:   "Get detailed information about a specific network device",
		Long: `Display detailed information about a specific network device including its
configuration, status, interfaces, and operational details.

Arguments:
  network_device_id   The unique identifier of the network device

Examples:
  # Get details for network device with ID 12345
  metalcloud-cli network-device get 12345

  # Using alias
  metalcloud-cli switch show 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceGet(cmd.Context(), args[0])
		},
	}

	networkDeviceConfigExampleCmd = &cobra.Command{
		Use:   "config-example",
		Short: "Generate example configuration template for network devices",
		Long: `Generate an example JSON configuration template that can be used to create
or update network devices. This template includes all available configuration
options with example values and documentation.

The generated template can be saved to a file and modified as needed for actual
device configuration.

Examples:
  # Display example configuration
  metalcloud-cli network-device config-example -f json

  # Save example to file
  metalcloud-cli network-device config-example -f json > device-config.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceConfigExample(cmd.Context())
		},
	}

	networkDeviceCreateCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Short:   "Create a new network device with specified configuration",
		Long: `Create a new network device using configuration provided via JSON file or pipe.

The configuration must include device details such as management IP, credentials,
device type, and other operational parameters.

Required Flags:
  --config-source   Source of configuration data (required)
                   Values: 'pipe' for stdin input, or path to JSON file

Use the 'config-example' command to generate an example configuration:

  {
    "siteId": 1,
    "driver": "sonic_enterprise",
    "identifierString": "example",
    "serialNumber": "1234567890",
    "chassisIdentifier": "example",
    "chassisRackId": 1,
    "position": "leaf",
    "isGateway": false,
    "isStorageSwitch": false,
    "isBorderDevice": false,
    "managementMAC": "AA:BB:CC:DD:EE:FF",
    "managementAddress": "1.1.1.1",
    "managementAddressGateway": "1.1.1.1",
    "managementAddressMask": "255.255.255.0",
    "loopbackAddress": "127.0.0.1",
    "vtepAddress": null,
    "asn": 65000,
    "managementPort": 22,
    "username": "admin",
    "managementPassword": "password",
    "syslogEnabled": true
  }

Examples:
  # Create device from JSON file
  metalcloud-cli network-device create --config-source device-config.json

  # Create device from pipe input
  cat device-config.json | metalcloud-cli network-device create --config-source pipe

  # Create device with inline JSON
  echo '{"management_ip":"10.0.1.100","type":"cisco"}' | metalcloud-cli nd create --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(networkDeviceFlags.configSource)
			if err != nil {
				return err
			}

			return network_device.NetworkDeviceCreate(cmd.Context(), config)
		},
	}

	networkDeviceUpdateCmd = &cobra.Command{
		Use:     "update <network_device_id>",
		Aliases: []string{"modify"},
		Short:   "Update configuration of an existing network device",
		Long: `Update the configuration of an existing network device using JSON configuration
provided via file or pipe. Only the specified fields will be updated; other
configuration will remain unchanged.

Arguments:
  network_device_id   The unique identifier of the network device to update

Required Flags:
  --config-source   Source of configuration updates (required)
                   Values: 'pipe' for stdin input, or path to JSON file

Examples:
  # Update device from JSON file
  metalcloud-cli network-device update 12345 --config-source updates.json

  # Update device from pipe input
  cat updates.json | metalcloud-cli network-device update 12345 --config-source pipe

  # Update specific field
  echo '{"management_ip":"10.0.1.101"}' | metalcloud-cli nd update 12345 --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
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
		Use:     "delete <network_device_id>",
		Aliases: []string{"rm"},
		Short:   "Delete a network device from the infrastructure",
		Long: `Delete a network device from the infrastructure. This operation will remove
the device from management and monitoring. The physical device will no longer
be controlled by the system.

WARNING: This operation is irreversible. Ensure the device is not in use
before deletion.

Arguments:
  network_device_id   The unique identifier of the network device to delete

Examples:
  # Delete network device
  metalcloud-cli network-device delete 12345

  # Using alias
  metalcloud-cli switch rm 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceDelete(cmd.Context(), args[0])
		},
	}

	networkDeviceArchiveCmd = &cobra.Command{
		Use:   "archive <network_device_id>",
		Short: "Archive a network device (soft delete with history preservation)",
		Long: `Archive a network device, which performs a soft delete operation while
preserving the device's operational history and configuration for audit purposes.

Archived devices are no longer active in the infrastructure but their data
is retained for compliance and historical analysis.

Arguments:
  network_device_id   The unique identifier of the network device to archive

Examples:
  # Archive network device
  metalcloud-cli network-device archive 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceArchive(cmd.Context(), args[0])
		},
	}

	networkDeviceDiscoverCmd = &cobra.Command{
		Use:   "discover <network_device_id>",
		Short: "Discover and inventory network device interfaces and configuration",
		Long: `Initiate discovery process for a network device to automatically detect and
inventory its interfaces, hardware components, and software configuration.

This process connects to the device using its management interface and gathers
detailed information about:
- Physical interfaces and their status
- Hardware components and capabilities
- Software version and configuration
- VLAN and networking setup

Arguments:
  network_device_id   The unique identifier of the network device to discover

Examples:
  # Discover device interfaces and configuration
  metalcloud-cli network-device discover 12345

  # Using alias
  metalcloud-cli switch discover 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceDiscover(cmd.Context(), args[0])
		},
	}

	networkDeviceGetCredentialsCmd = &cobra.Command{
		Use:   "get-credentials <network_device_id>",
		Short: "Retrieve management credentials for a network device",
		Long: `Retrieve the management credentials (username/password) configured for
accessing a network device. This information is used by the system to
connect to the device for configuration and monitoring.

Note: This command may require elevated permissions and credentials will
be displayed in plain text.

Arguments:
  network_device_id   The unique identifier of the network device

Examples:
  # Get device credentials
  metalcloud-cli network-device get-credentials 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceGetCredentials(cmd.Context(), args[0])
		},
	}

	networkDeviceGetPortsCmd = &cobra.Command{
		Use:   "get-ports <network_device_id>",
		Short: "Get real-time port statistics directly from the network device",
		Long: `Retrieve real-time port statistics and status information directly from
the network device. This provides current operational data including:
- Port status (up/down)
- Traffic statistics (bytes, packets)
- Error counters
- Link speed and duplex settings

This data is fetched directly from the device rather than cached information.

Arguments:
  network_device_id   The unique identifier of the network device

Examples:
  # Get current port statistics
  metalcloud-cli network-device get-ports 12345

  # Using alias
  metalcloud-cli switch get-ports 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceGetPorts(cmd.Context(), args[0])
		},
	}

	networkDeviceSetPortStatusCmd = &cobra.Command{
		Use:   "set-port-status <network_device_id>",
		Short: "Enable or disable a specific port on the network device",
		Long: `Set the administrative status of a specific port on the network device.
This allows you to enable (bring up) or disable (bring down) individual
ports for maintenance or troubleshooting purposes.

Arguments:
  network_device_id   The unique identifier of the network device

Required Flags (both must be specified):
  --port-id    ID or name of the port to modify
  --action     Action to perform on the port
               Values: 'up' (enable port), 'down' (disable port)

Examples:
  # Bring port down for maintenance
  metalcloud-cli network-device set-port-status 12345 --port-id eth0/1 --action down

  # Bring port back up
  metalcloud-cli network-device set-port-status 12345 --port-id eth0/1 --action up

  # Using port number
  metalcloud-cli nd set-port-status 12345 --port-id 24 --action up`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceSetPortStatus(cmd.Context(), args[0], networkDeviceFlags.portId, networkDeviceFlags.portStatusAction)
		},
	}

	networkDeviceResetCmd = &cobra.Command{
		Use:   "reset <network_device_id>",
		Short: "Reset network device to factory defaults (destructive operation)",
		Long: `Reset a network device to its factory default state, destroying all
custom configurations, VLANs, and settings. This is a destructive operation
that will:
- Remove all VLANs and network configurations
- Reset interface configurations
- Clear all custom settings
- Restore factory default credentials

WARNING: This operation is irreversible and will cause network disruption.
Ensure all connected services are properly migrated before performing this reset.

Arguments:
  network_device_id   The unique identifier of the network device to reset

Examples:
  # Reset device to factory defaults
  metalcloud-cli network-device reset 12345

  # Confirm the operation is intentional
  metalcloud-cli switch reset 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceReset(cmd.Context(), args[0])
		},
	}

	networkDeviceSetAsFailedCmd = &cobra.Command{
		Use:          "set-failed <network_device_id>",
		Short:        "Set the network device as failed",
		Long:         `Change the operational status of a network device to failed.`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceSetFailed(cmd.Context(), args[0])
		},
	}

	networkDeviceEnableSyslogCmd = &cobra.Command{
		Use:   "enable-syslog <network_device_id>",
		Short: "Enable remote syslog forwarding on the network device",
		Long: `Enable remote syslog forwarding on the network device to send system logs
and events to a centralized syslog server. This helps with centralized
monitoring and troubleshooting.

The device will be configured to forward its system logs, including:
- Interface status changes
- Configuration changes
- System events and errors
- Security events

Arguments:
  network_device_id   The unique identifier of the network device

Examples:
  # Enable syslog forwarding
  metalcloud-cli network-device enable-syslog 12345

  # Using alias
  metalcloud-cli switch enable-syslog 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceEnableSyslog(cmd.Context(), args[0])
		},
	}

	networkDeviceGetDefaultsCmd = &cobra.Command{
		Use:   "get-defaults <site_id>",
		Short: "Get default network device configuration settings for a site",
		Long: `Retrieve the default configuration settings and templates that are applied
to new network devices when they are added to a specific site.

These defaults include standard configurations for:
- Management network settings
- VLAN configurations
- Security policies
- Monitoring settings
- Device-specific parameters

Arguments:
  site_id   The unique identifier of the site

Examples:
  # Get defaults for site
  metalcloud-cli network-device get-defaults site-123

  # View defaults for current site
  metalcloud-cli nd get-defaults my-datacenter`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceGetDefaults(cmd.Context(), args[0])
		},
	}

	networkDeviceAddDefaultsCmd = &cobra.Command{
		Use:   "add-defaults",
		Short: "Add network device default configuration",
		Long: `Add network device default configuration that will be applied to new
devices when they are added to sites. These defaults provide consistent
baseline configurations across your infrastructure.

Default configurations can include:
- Management network settings and credentials
- Standard VLAN configurations
- Security policies and access controls
- Monitoring and logging settings
- Device-specific operational parameters
- Network topology preferences

The configuration is provided via JSON file or pipe input and will be merged
with existing defaults, allowing for incremental updates.

Required Flags:
  --config-source   Source of default configuration data (required)
                   Values: 'pipe' for stdin input, or path to JSON file

Use the 'example-defaults' command to see the configuration format:

Examples:
  # Add defaults from JSON file
  metalcloud-cli network-device add-defaults --config-source defaults.json

  # Add defaults from pipe input
  cat site-defaults.json | metalcloud-cli network-device add-defaults --config-source pipe

  # Update specific default settings
  echo '{"syslogEnabled": true, "managementPort": 22}' | metalcloud-cli nd add-defaults --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(networkDeviceFlags.configSource)
			if err != nil {
				return err
			}

			return network_device.NetworkDeviceAddDefaults(cmd.Context(), config)
		},
	}

	networkDeviceExampleDefaultsCmd = &cobra.Command{
		Use:          "example-defaults",
		Short:        "Network device default configuration example",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceExampleDefaults(cmd.Context())
		},
	}
)

func init() {
	rootCmd.AddCommand(networkDeviceCmd)

	networkDeviceCmd.AddCommand(networkDeviceListCmd)
	networkDeviceListCmd.Flags().StringSliceVar(&networkDeviceFlags.filterStatus, "filter-status", []string{"active"}, "Filter the result by network device status.")

	networkDeviceCmd.AddCommand(networkDeviceGetCmd)

	networkDeviceCmd.AddCommand(networkDeviceConfigExampleCmd)

	networkDeviceCmd.AddCommand(networkDeviceCreateCmd)
	networkDeviceCreateCmd.Flags().StringVar(&networkDeviceFlags.configSource, "config-source", "", "Source of the new network device configuration. Can be 'pipe' or path to a JSON file.")
	networkDeviceCreateCmd.MarkFlagRequired("config-source")

	networkDeviceCmd.AddCommand(networkDeviceUpdateCmd)
	networkDeviceUpdateCmd.Flags().StringVar(&networkDeviceFlags.configSource, "config-source", "", "Source of the network device configuration updates. Can be 'pipe' or path to a JSON file.")
	networkDeviceUpdateCmd.MarkFlagRequired("config-source")

	networkDeviceCmd.AddCommand(networkDeviceDeleteCmd)

	networkDeviceCmd.AddCommand(networkDeviceArchiveCmd)

	networkDeviceCmd.AddCommand(networkDeviceDiscoverCmd)

	networkDeviceCmd.AddCommand(networkDeviceGetCredentialsCmd)

	networkDeviceCmd.AddCommand(networkDeviceGetPortsCmd)

	networkDeviceCmd.AddCommand(networkDeviceSetPortStatusCmd)
	networkDeviceSetPortStatusCmd.Flags().StringVar(&networkDeviceFlags.portId, "port-id", "", "ID of the port to change status.")
	networkDeviceSetPortStatusCmd.Flags().StringVar(&networkDeviceFlags.portStatusAction, "action", "", "Action to perform on the port (up/down).")
	networkDeviceSetPortStatusCmd.MarkFlagsOneRequired("port-id", "action")

	networkDeviceCmd.AddCommand(networkDeviceResetCmd)

	networkDeviceCmd.AddCommand(networkDeviceSetAsFailedCmd)

	networkDeviceCmd.AddCommand(networkDeviceEnableSyslogCmd)

	networkDeviceCmd.AddCommand(networkDeviceGetDefaultsCmd)

	networkDeviceCmd.AddCommand(networkDeviceAddDefaultsCmd)
	networkDeviceAddDefaultsCmd.Flags().StringVar(&networkDeviceFlags.configSource, "config-source", "", "Source of the network device default configuration. Can be 'pipe' or path to a JSON file.")
	networkDeviceAddDefaultsCmd.MarkFlagRequired("config-source")

	networkDeviceCmd.AddCommand(networkDeviceExampleDefaultsCmd)
}
