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
		portEnabled      bool
		portDescription  string
		portIpFamily     string
		portIpAddress    string
		portIpPrefix     int32
		discoverTargets  []string
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

	networkDeviceCreateBulkCmd = &cobra.Command{
		Use:     "create-bulk",
		Aliases: []string{"bulk-create", "new-bulk"},
		Short:   "Create multiple network devices in a single operation",
		Long: `Create multiple network devices at once from a JSON or YAML configuration file.

This command processes a list of network device configurations and creates all devices in sequence.
Each device configuration follows the same format as the single network-device create command.

Required Flags:
  --config-source   Source of the bulk network device configuration (JSON/YAML file path or 'pipe')

Configuration File Format (JSON):
  [
    {
      "siteId": 1,
      "driver": "sonic_enterprise",
      "identifierString": "leaf-01",
      "position": "leaf",
      "managementAddress": "10.0.1.100",
      "managementPort": 22,
      "username": "admin",
      "managementPassword": "password"
    },
    {
      "siteId": 1,
      "driver": "sonic_enterprise",
      "identifierString": "spine-01",
      "position": "spine",
      "managementAddress": "10.0.1.101",
      "managementPort": 22,
      "username": "admin",
      "managementPassword": "password"
    }
  ]

The command will report success/failure for each device and provide a summary at the end.`,
		Example: `  # Create network devices from JSON file
  metalcloud-cli network-device create-bulk --config-source devices.json

  # Create network devices from YAML file
  metalcloud-cli network-device create-bulk --config-source devices.yaml

  # Create network devices from pipe
  cat devices.json | metalcloud-cli network-device create-bulk --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(networkDeviceFlags.configSource)
			if err != nil {
				return err
			}

			return network_device.NetworkDeviceCreateBulk(cmd.Context(), config)
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

By default all discovery types are run (hardware, software, ports). Use --target
to restrict discovery to specific types; repeat the flag for more than one.
Discovered data is persisted to the device inventory.

Arguments:
  network_device_id   The unique identifier of the network device to discover

Flags:
  --target strings   Discovery type(s) to run: hardware, software, ports.
                     Repeatable. Defaults to all three.

Examples:
  # Full discovery (hardware, software and ports)
  metalcloud-cli network-device discover 12345

  # Discover only the ports
  metalcloud-cli network-device discover 12345 --target ports

  # Discover hardware and software only
  metalcloud-cli network-device discover 12345 --target hardware --target software

  # Using alias
  metalcloud-cli switch discover 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceDiscover(cmd.Context(), args[0], networkDeviceFlags.discoverTargets)
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
		Short: "List the interface inventory of a network device",
		Long: `List the interfaces of a network device as MetalSoft has them inventoried:
interface id, name, kind, description, MAC address, LAG membership and tags.

The interface ids reported here are the ones the 'network-device port'
sub-commands take. For the operational state read from the device itself
(link state, negotiated speed, utilization) use
'network-device port live-status' instead.

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

	networkDeviceUpdatePortConfigCmd = &cobra.Command{
		Use:   "update-port-config <network_device_id>",
		Short: "Update the staged config (enable/description) of a port",
		Long: `Update the staged configuration of a single network device port, addressed
by its numeric interface id.

This patches the port's persistent config (applied on the next fabric deploy),
not its live administrative status (use 'set-port-status' for that). You can set
the enabled flag and/or the interface description.

Arguments:
  network_device_id   The unique identifier of the network device

Required Flags:
  --port-id           Numeric interface id of the port to configure

Optional Flags (at least one must be specified):
  --enabled           Whether the port should be enabled (true/false)
  --description       Interface description text

Examples:
  # Enable a port and set its description
  metalcloud-cli network-device update-port-config 12345 --port-id 67890 --enabled --description "to_spine-s00_swp1s0"

  # Only set a description
  metalcloud-cli nd update-port-config 12345 --port-id 67890 --description "to_hgx-su00-h00_enp26s0f0np0"`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var enabled *bool
			if cmd.Flags().Changed("enabled") {
				enabled = &networkDeviceFlags.portEnabled
			}
			var description *string
			if cmd.Flags().Changed("description") {
				description = &networkDeviceFlags.portDescription
			}

			return network_device.NetworkDeviceUpdatePortConfig(cmd.Context(), args[0], networkDeviceFlags.portId, enabled, description)
		},
	}

	networkDeviceAddPortIpCmd = &cobra.Command{
		Use:   "add-port-ip <network_device_id>",
		Short: "Add an IP address to a network device port",
		Long: `Stage a new IP address on a network device port, addressed by its numeric
interface id. Typically used to assign a /32 loopback address to the loopback
interface.

Arguments:
  network_device_id   The unique identifier of the network device

Required Flags:
  --port-id           Numeric interface id of the port
  --address           IP address to add (without prefix, e.g. "10.253.128.1")
  --prefix            Prefix length (e.g. 32 for a loopback /32)

Optional Flags:
  --family            Address family: ipv4 (default) or ipv6

Examples:
  # Add a /32 loopback address
  metalcloud-cli network-device add-port-ip 12345 --port-id 67890 --address 10.253.128.1 --prefix 32

  # Add an IPv6 address
  metalcloud-cli nd add-port-ip 12345 --port-id 67890 --family ipv6 --address 2001:db8::1 --prefix 128`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceAddPortIp(cmd.Context(), args[0], networkDeviceFlags.portId, networkDeviceFlags.portIpFamily, networkDeviceFlags.portIpAddress, networkDeviceFlags.portIpPrefix)
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

	networkDeviceSetAsDefectiveCmd = &cobra.Command{
		Use:          "set-defective <network_device_id>",
		Aliases:      []string{"set-failed"},
		Short:        "Set the network device as defective",
		Long:         `Change the operational status of a network device to defective.`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceSetDefective(cmd.Context(), args[0])
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

	networkDeviceCmd.AddCommand(networkDeviceCreateBulkCmd)
	networkDeviceCreateBulkCmd.Flags().StringVar(&networkDeviceFlags.configSource, "config-source", "", "Source of the bulk network device configuration. Can be 'pipe' or path to a JSON/YAML file containing a list of devices.")
	networkDeviceCreateBulkCmd.MarkFlagRequired("config-source")

	networkDeviceCmd.AddCommand(networkDeviceUpdateCmd)
	networkDeviceUpdateCmd.Flags().StringVar(&networkDeviceFlags.configSource, "config-source", "", "Source of the network device configuration updates. Can be 'pipe' or path to a JSON file.")
	networkDeviceUpdateCmd.MarkFlagRequired("config-source")

	networkDeviceCmd.AddCommand(networkDeviceDeleteCmd)

	networkDeviceCmd.AddCommand(networkDeviceArchiveCmd)

	networkDeviceCmd.AddCommand(networkDeviceDiscoverCmd)
	networkDeviceDiscoverCmd.Flags().StringSliceVar(&networkDeviceFlags.discoverTargets, "target", nil, "Discovery type(s) to run: hardware, software, ports. Repeatable. Defaults to all three.")

	networkDeviceCmd.AddCommand(networkDeviceGetCredentialsCmd)

	networkDeviceCmd.AddCommand(networkDeviceGetPortsCmd)

	networkDeviceCmd.AddCommand(networkDeviceSetPortStatusCmd)
	networkDeviceSetPortStatusCmd.Flags().StringVar(&networkDeviceFlags.portId, "port-id", "", "ID of the port to change status.")
	networkDeviceSetPortStatusCmd.Flags().StringVar(&networkDeviceFlags.portStatusAction, "action", "", "Action to perform on the port (up/down).")
	networkDeviceSetPortStatusCmd.MarkFlagsOneRequired("port-id", "action")

	networkDeviceCmd.AddCommand(networkDeviceUpdatePortConfigCmd)
	networkDeviceUpdatePortConfigCmd.Flags().StringVar(&networkDeviceFlags.portId, "port-id", "", "Numeric interface id of the port to configure.")
	networkDeviceUpdatePortConfigCmd.Flags().BoolVar(&networkDeviceFlags.portEnabled, "enabled", false, "Whether the port should be enabled.")
	networkDeviceUpdatePortConfigCmd.Flags().StringVar(&networkDeviceFlags.portDescription, "description", "", "Interface description text.")
	networkDeviceUpdatePortConfigCmd.MarkFlagRequired("port-id")
	networkDeviceUpdatePortConfigCmd.MarkFlagsOneRequired("enabled", "description")

	networkDeviceCmd.AddCommand(networkDeviceAddPortIpCmd)
	networkDeviceAddPortIpCmd.Flags().StringVar(&networkDeviceFlags.portId, "port-id", "", "Numeric interface id of the port.")
	networkDeviceAddPortIpCmd.Flags().StringVar(&networkDeviceFlags.portIpFamily, "family", "ipv4", "Address family: ipv4 or ipv6.")
	networkDeviceAddPortIpCmd.Flags().StringVar(&networkDeviceFlags.portIpAddress, "address", "", "IP address to add (without prefix).")
	networkDeviceAddPortIpCmd.Flags().Int32Var(&networkDeviceFlags.portIpPrefix, "prefix", 0, "Prefix length (e.g. 32 for a loopback /32).")
	networkDeviceAddPortIpCmd.MarkFlagRequired("port-id")
	networkDeviceAddPortIpCmd.MarkFlagRequired("address")
	networkDeviceAddPortIpCmd.MarkFlagRequired("prefix")

	networkDeviceCmd.AddCommand(networkDeviceResetCmd)

	networkDeviceCmd.AddCommand(networkDeviceSetAsDefectiveCmd)

	networkDeviceCmd.AddCommand(networkDeviceEnableSyslogCmd)
}

// ---------------------------------------------------------------------------
// Network device sub-resources: ports, breakouts, virtual functions, secrets,
// monitoring, drift & snapshots, lifecycle actions and vendor profiles.
// ---------------------------------------------------------------------------

var (
	networkDeviceSubFlags = struct {
		configSource       string
		secretValue        string
		snmpPort           int32
		snmpCommunity      string
		snmpContact        string
		siteIds            []string
		fabricIds          []string
		socketId           string
		filterId           []string
		filterSiteId       []string
		filterStatus       []string
		filterHealthStatus []string
		filterKind         []string
		snapshotKind       string
		reprovisionType    string
		search             string
	}{}

	// --- ports ------------------------------------------------------------

	networkDevicePortCmd = &cobra.Command{
		Use:     "port [command]",
		Aliases: []string{"ports"},
		Short:   "Manage the interfaces of a network device",
		Long: `Manage the interfaces (ports) of a network device.

Ports are addressed by their numeric interface id, as shown by
'metalcloud-cli network-device get-ports <network_device_id>'.`,
	}

	networkDevicePortAddCmd = &cobra.Command{
		Use:     "add <network_device_id>",
		Aliases: []string{"create", "new"},
		Short:   "Create a logical interface on a network device",
		Long: `Create a logical interface (e.g. a loopback or a sub-interface) on a network device.

Required Arguments:
  network_device_id   The numeric id or label of the network device

Required Flags:
  --config-source     Source of the interface configuration
                      Values: 'pipe' for stdin input, or path to a JSON/YAML file

Examples:
  # Create a loopback interface from a file
  metalcloud-cli network-device port add 12345 --config-source loopback.json

  # Create an interface from pipe input
  echo '{"kind":"loopback","name":"Loopback1"}' | metalcloud-cli network-device port add 12345 --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(networkDeviceSubFlags.configSource)
			if err != nil {
				return err
			}

			return network_device.NetworkDevicePortCreate(cmd.Context(), args[0], config)
		},
	}

	networkDevicePortConfigExampleCmd = &cobra.Command{
		Use:          "config-example",
		Short:        "Example configuration for creating a network device interface",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDevicePortConfigExample(cmd.Context())
		},
	}

	networkDevicePortRemoveCmd = &cobra.Command{
		Use:     "remove <network_device_id> <port_id>",
		Aliases: []string{"delete", "rm"},
		Short:   "Delete an interface of a network device",
		Long: `Delete a logical interface of a network device.

Required Arguments:
  network_device_id   The numeric id or label of the network device
  port_id             The numeric interface id of the port

Examples:
  # Delete interface 42 of device 12345
  metalcloud-cli network-device port remove 12345 42`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDevicePortDelete(cmd.Context(), args[0], args[1])
		},
	}

	networkDevicePortGetConfigCmd = &cobra.Command{
		Use:   "get-config <network_device_id> <port_id>",
		Short: "Get the staged configuration of a network device port",
		Long: `Display the staged (desired) configuration of a network device port, including
the admin overrides for description, MTU, enabled state and speed, plus the
optimistic-lock revision of the configuration buffer.

Required Arguments:
  network_device_id   The numeric id or label of the network device
  port_id             The numeric interface id of the port

Examples:
  # Show the staged config of interface 42
  metalcloud-cli network-device port get-config 12345 42`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDevicePortGetConfig(cmd.Context(), args[0], args[1])
		},
	}

	networkDevicePortLiveStatusCmd = &cobra.Command{
		Use:   "live-status <network_device_id>",
		Short: "Query the device for the operational state of its ports",
		Long: `Query the network device itself for the operational state of its physical ports:
admin/link state, negotiated speed and duplex, and traffic utilization.

This differs from 'network-device get-ports', which lists the interface
inventory stored by MetalSoft without contacting the device.

Required Arguments:
  network_device_id   The numeric id or label of the network device

Examples:
  # Show the live port state of device 12345
  metalcloud-cli network-device port live-status 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDevicePortsLiveStatus(cmd.Context(), args[0])
		},
	}

	// --- port IP addresses -------------------------------------------------

	networkDevicePortIpCmd = &cobra.Command{
		Use:     "ip [command]",
		Aliases: []string{"ips"},
		Short:   "Manage the IP addresses of a network device port",
		Long: `Manage the IP addresses staged on a network device port.

Addresses are grouped by address family; the family argument is either
'ipv4' or 'ipv6'.`,
	}

	networkDevicePortIpListCmd = &cobra.Command{
		Use:     "list <network_device_id> <port_id> <family>",
		Aliases: []string{"ls"},
		Short:   "List the IP addresses of a network device port",
		Long: `List the IP addresses of one address family staged on a network device port.

Required Arguments:
  network_device_id   The numeric id or label of the network device
  port_id             The numeric interface id of the port
  family              Address family: ipv4 or ipv6

Examples:
  # List the IPv4 addresses of interface 42
  metalcloud-cli network-device port ip list 12345 42 ipv4`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_READ},
		Args:         cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDevicePortIpList(cmd.Context(), args[0], args[1], args[2])
		},
	}

	networkDevicePortIpGetCmd = &cobra.Command{
		Use:     "get <network_device_id> <port_id> <family> <ip_id>",
		Aliases: []string{"show"},
		Short:   "Get one IP address of a network device port",
		Long: `Display one IP address staged on a network device port.

Required Arguments:
  network_device_id   The numeric id or label of the network device
  port_id             The numeric interface id of the port
  family              Address family: ipv4 or ipv6
  ip_id               The numeric id of the address

Examples:
  # Show IPv4 address 7 of interface 42
  metalcloud-cli network-device port ip get 12345 42 ipv4 7`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_READ},
		Args:         cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDevicePortIpGet(cmd.Context(), args[0], args[1], args[2], args[3])
		},
	}

	networkDevicePortIpRemoveCmd = &cobra.Command{
		Use:     "remove <network_device_id> <port_id> <family> <ip_id>",
		Aliases: []string{"delete", "rm"},
		Short:   "Remove one IP address from a network device port",
		Long: `Remove one IP address from a network device port.

Required Arguments:
  network_device_id   The numeric id or label of the network device
  port_id             The numeric interface id of the port
  family              Address family: ipv4 or ipv6
  ip_id               The numeric id of the address

Examples:
  # Remove IPv4 address 7 from interface 42
  metalcloud-cli network-device port ip remove 12345 42 ipv4 7`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
		Args:         cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDevicePortIpRemove(cmd.Context(), args[0], args[1], args[2], args[3])
		},
	}

	networkDevicePortIpReplaceCmd = &cobra.Command{
		Use:   "replace <network_device_id> <port_id> <family>",
		Short: "Replace the whole IP address set of a network device port",
		Long: `Replace the complete IP address set of one address family on a network device
port with the supplied list. Addresses missing from the list are removed.

Required Arguments:
  network_device_id   The numeric id or label of the network device
  port_id             The numeric interface id of the port
  family              Address family: ipv4 or ipv6

Required Flags:
  --config-source     Source of the address list
                      Values: 'pipe' for stdin input, or path to a JSON/YAML file

Examples:
  # Replace the IPv4 addresses of interface 42
  metalcloud-cli network-device port ip replace 12345 42 ipv4 --config-source ips.json

  # Replace them from pipe input
  echo '{"ips":[{"address":"10.0.0.1","prefixLength":32}]}' | metalcloud-cli network-device port ip replace 12345 42 ipv4 --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
		Args:         cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(networkDeviceSubFlags.configSource)
			if err != nil {
				return err
			}

			return network_device.NetworkDevicePortIpReplace(cmd.Context(), args[0], args[1], args[2], config)
		},
	}

	networkDevicePortIpConfigExampleCmd = &cobra.Command{
		Use:          "config-example",
		Short:        "Example address list for the port IP replace command",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDevicePortIpReplaceConfigExample(cmd.Context())
		},
	}

	// --- port virtual functions -------------------------------------------

	networkDevicePortVirtualFunctionCmd = &cobra.Command{
		Use:     "virtual-function [command]",
		Aliases: []string{"vf"},
		Short:   "Inspect the virtual functions of a network device port",
	}

	networkDevicePortVirtualFunctionListCmd = &cobra.Command{
		Use:     "list <network_device_id> <port_id>",
		Aliases: []string{"ls"},
		Short:   "List the virtual functions of a network device port",
		Long: `List the virtual functions (SR-IOV VFs) exposed by one network device port.

Required Arguments:
  network_device_id   The numeric id or label of the network device
  port_id             The numeric interface id of the port

Examples:
  # List the virtual functions of interface 42
  metalcloud-cli network-device port virtual-function list 12345 42`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDevicePortVirtualFunctionList(cmd.Context(), args[0], args[1])
		},
	}

	networkDevicePortVirtualFunctionGetCmd = &cobra.Command{
		Use:     "get <network_device_id> <port_id> <virtual_function_id>",
		Aliases: []string{"show"},
		Short:   "Get one virtual function of a network device port",
		Long: `Display one virtual function of a network device port.

Required Arguments:
  network_device_id     The numeric id or label of the network device
  port_id               The numeric interface id of the port
  virtual_function_id   The numeric id of the virtual function

Examples:
  # Show virtual function 3 of interface 42
  metalcloud-cli network-device port virtual-function get 12345 42 3`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_READ},
		Args:         cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDevicePortVirtualFunctionGet(cmd.Context(), args[0], args[1], args[2])
		},
	}

	// --- breakouts ---------------------------------------------------------

	networkDeviceBreakoutCmd = &cobra.Command{
		Use:     "breakout [command]",
		Aliases: []string{"breakouts"},
		Short:   "Manage the port breakouts of a network device",
		Long: `Manage the port breakouts of a network device.

A breakout splits one physical port into several child interfaces. The applied
breakout groups are reported by 'get', while 'get-config' and 'update-config'
work on the groups staged for the next deploy.`,
	}

	networkDeviceBreakoutListCmd = &cobra.Command{
		Use:     "list <network_device_id>",
		Aliases: []string{"ls"},
		Short:   "List the breakouts of a network device",
		Long: `List the port breakouts configured on a network device.

Required Arguments:
  network_device_id   The numeric id or label of the network device

Examples:
  # List the breakouts of device 12345
  metalcloud-cli network-device breakout list 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceBreakoutList(cmd.Context(), args[0])
		},
	}

	networkDeviceBreakoutGetCmd = &cobra.Command{
		Use:     "get <network_device_id> <breakout_id>",
		Aliases: []string{"show"},
		Short:   "Get one breakout of a network device",
		Long: `Display one port breakout of a network device, including the applied breakout
groups and the staged configuration.

Required Arguments:
  network_device_id   The numeric id or label of the network device
  breakout_id         The numeric id of the breakout

Examples:
  # Show breakout 7 of device 12345
  metalcloud-cli network-device breakout get 12345 7`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceBreakoutGet(cmd.Context(), args[0], args[1])
		},
	}

	networkDeviceBreakoutGetConfigCmd = &cobra.Command{
		Use:   "get-config <network_device_id> <breakout_id>",
		Short: "Get the staged configuration of a breakout",
		Long: `Display the breakout groups staged for the next deploy, together with the
optimistic-lock revision of the configuration buffer.

Required Arguments:
  network_device_id   The numeric id or label of the network device
  breakout_id         The numeric id of the breakout

Examples:
  # Show the staged breakout groups of breakout 7
  metalcloud-cli network-device breakout get-config 12345 7`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceBreakoutGetConfig(cmd.Context(), args[0], args[1])
		},
	}

	networkDeviceBreakoutCreateCmd = &cobra.Command{
		Use:     "create <network_device_id>",
		Aliases: []string{"new"},
		Short:   "Create a port breakout on a network device",
		Long: `Create a port breakout on a network device, splitting one physical port into
several child interfaces.

Required Arguments:
  network_device_id   The numeric id or label of the network device

Required Flags:
  --config-source     Source of the breakout configuration
                      Values: 'pipe' for stdin input, or path to a JSON/YAML file

Examples:
  # Create a breakout from a file
  metalcloud-cli network-device breakout create 12345 --config-source breakout.json

  # Split Ethernet0 into 4x25G from pipe input
  echo '{"portName":"Ethernet0","groups":[{"numberOfInterfaces":4,"speed":"25G"}]}' | metalcloud-cli network-device breakout create 12345 --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(networkDeviceSubFlags.configSource)
			if err != nil {
				return err
			}

			return network_device.NetworkDeviceBreakoutCreate(cmd.Context(), args[0], config)
		},
	}

	networkDeviceBreakoutUpdateConfigCmd = &cobra.Command{
		Use:   "update-config <network_device_id> <breakout_id>",
		Short: "Stage a new breakout group set on a breakout",
		Long: `Stage a new set of breakout groups on an existing breakout. The staged groups
are applied on the next deploy.

Required Arguments:
  network_device_id   The numeric id or label of the network device
  breakout_id         The numeric id of the breakout

Required Flags:
  --config-source     Source of the breakout group configuration
                      Values: 'pipe' for stdin input, or path to a JSON/YAML file

Examples:
  # Stage a new group set from a file
  metalcloud-cli network-device breakout update-config 12345 7 --config-source groups.json

  # Stage a 2x50G split from pipe input
  echo '{"breakoutGroups":[{"numberOfInterfaces":2,"speed":"50G"}]}' | metalcloud-cli network-device breakout update-config 12345 7 --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(networkDeviceSubFlags.configSource)
			if err != nil {
				return err
			}

			return network_device.NetworkDeviceBreakoutUpdateConfig(cmd.Context(), args[0], args[1], config)
		},
	}

	networkDeviceBreakoutDeleteCmd = &cobra.Command{
		Use:     "delete <network_device_id> <breakout_id>",
		Aliases: []string{"rm"},
		Short:   "Delete a breakout of a network device",
		Long: `Delete a port breakout of a network device, returning the physical port to its
unsplit state on the next deploy.

Required Arguments:
  network_device_id   The numeric id or label of the network device
  breakout_id         The numeric id of the breakout

Examples:
  # Delete breakout 7 of device 12345
  metalcloud-cli network-device breakout delete 12345 7`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceBreakoutDelete(cmd.Context(), args[0], args[1])
		},
	}

	networkDeviceBreakoutConfigExampleCmd = &cobra.Command{
		Use:          "config-example",
		Short:        "Example configuration for creating a breakout",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceBreakoutConfigExample(cmd.Context())
		},
	}

	networkDeviceBreakoutUpdateConfigExampleCmd = &cobra.Command{
		Use:          "update-config-example",
		Short:        "Example configuration for staging breakout groups",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceBreakoutUpdateConfigExample(cmd.Context())
		},
	}

	// --- device virtual functions -----------------------------------------

	networkDeviceVirtualFunctionCmd = &cobra.Command{
		Use:     "virtual-function [command]",
		Aliases: []string{"vf"},
		Short:   "Inspect the virtual functions of a network device",
	}

	networkDeviceVirtualFunctionListCmd = &cobra.Command{
		Use:     "list <network_device_id>",
		Aliases: []string{"ls"},
		Short:   "List all virtual functions of a network device",
		Long: `List every virtual function of a network device, across all of its interfaces.

Required Arguments:
  network_device_id   The numeric id or label of the network device

Examples:
  # List the virtual functions of device 12345
  metalcloud-cli network-device virtual-function list 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceVirtualFunctionList(cmd.Context(), args[0])
		},
	}

	networkDeviceVirtualFunctionGetCmd = &cobra.Command{
		Use:     "get <network_device_id> <virtual_function_id>",
		Aliases: []string{"show"},
		Short:   "Get one virtual function of a network device",
		Long: `Display one virtual function of a network device.

Required Arguments:
  network_device_id     The numeric id or label of the network device
  virtual_function_id   The numeric id of the virtual function

Examples:
  # Show virtual function 3 of device 12345
  metalcloud-cli network-device virtual-function get 12345 3`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceVirtualFunctionGet(cmd.Context(), args[0], args[1])
		},
	}

	// --- secrets -----------------------------------------------------------

	networkDeviceSecretCmd = &cobra.Command{
		Use:     "secret [command]",
		Aliases: []string{"secrets"},
		Short:   "Manage the secrets of a network device",
		Long: `Manage the named secrets stored for a network device.

Secret values are stored encrypted and are only revealed by the
'get-credentials' sub-command.`,
	}

	networkDeviceSecretListCmd = &cobra.Command{
		Use:     "list <network_device_id>",
		Aliases: []string{"ls"},
		Short:   "List the secret names of a network device",
		Long: `List the names of the secrets stored for a network device. Values are not
returned by this command.

Required Arguments:
  network_device_id   The numeric id or label of the network device

Examples:
  # List the secret names of device 12345
  metalcloud-cli network-device secret list 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceSecretList(cmd.Context(), args[0])
		},
	}

	networkDeviceSecretSetCmd = &cobra.Command{
		Use:   "set <network_device_id> <name>",
		Short: "Store a named secret of a network device",
		Long: `Store (or replace) one named secret of a network device. The value is stored
encrypted.

Required Arguments:
  network_device_id   The numeric id or label of the network device
  name                The name of the secret

Required Flags:
  --value             The secret value to store

Examples:
  # Store the enable password of device 12345
  metalcloud-cli network-device secret set 12345 enable_password --value s3cr3t`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceSecretSet(cmd.Context(), args[0], args[1], networkDeviceSubFlags.secretValue)
		},
	}

	networkDeviceSecretGetCredentialsCmd = &cobra.Command{
		Use:   "get-credentials <network_device_id> <name>",
		Short: "Reveal the value of a network device secret",
		Long: `Reveal the unencrypted value of one named secret of a network device.

Required Arguments:
  network_device_id   The numeric id or label of the network device
  name                The name of the secret

Examples:
  # Reveal the enable password of device 12345
  metalcloud-cli network-device secret get-credentials 12345 enable_password`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceSecretGetCredentials(cmd.Context(), args[0], args[1])
		},
	}

	networkDeviceSecretRemoveCmd = &cobra.Command{
		Use:     "remove <network_device_id> <name>",
		Aliases: []string{"delete", "rm"},
		Short:   "Delete one secret of a network device",
		Long: `Delete one named secret of a network device.

Required Arguments:
  network_device_id   The numeric id or label of the network device
  name                The name of the secret

Examples:
  # Delete the enable password of device 12345
  metalcloud-cli network-device secret remove 12345 enable_password`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceSecretRemove(cmd.Context(), args[0], args[1])
		},
	}

	networkDeviceSecretRemoveAllCmd = &cobra.Command{
		Use:   "remove-all <network_device_id>",
		Short: "Delete every secret of a network device",
		Long: `Delete every secret stored for a network device.

Required Arguments:
  network_device_id   The numeric id or label of the network device

Examples:
  # Delete all secrets of device 12345
  metalcloud-cli network-device secret remove-all 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceSecretRemoveAll(cmd.Context(), args[0])
		},
	}

	// --- SNMP monitoring ---------------------------------------------------

	networkDeviceSnmpMonitoringCmd = &cobra.Command{
		Use:   "snmp-monitoring [command]",
		Short: "Manage the SNMP monitoring subscription of network devices",
		Long: `Subscribe network devices to, or unsubscribe them from, SNMP monitoring, and
inspect which monitoring agent polls each device.`,
	}

	networkDeviceSnmpMonitoringEnableCmd = &cobra.Command{
		Use:   "enable <network_device_id>",
		Short: "Subscribe a network device to SNMP monitoring",
		Long: `Subscribe a network device to SNMP monitoring.

Required Arguments:
  network_device_id   The numeric id or label of the network device

Examples:
  # Enable SNMP monitoring for device 12345
  metalcloud-cli network-device snmp-monitoring enable 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceSnmpMonitoringEnable(cmd.Context(), args[0])
		},
	}

	networkDeviceSnmpMonitoringDisableCmd = &cobra.Command{
		Use:   "disable <network_device_id>",
		Short: "Unsubscribe a network device from SNMP monitoring",
		Long: `Unsubscribe a network device from SNMP monitoring.

Required Arguments:
  network_device_id   The numeric id or label of the network device

Examples:
  # Disable SNMP monitoring for device 12345
  metalcloud-cli network-device snmp-monitoring disable 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceSnmpMonitoringDisable(cmd.Context(), args[0])
		},
	}

	networkDeviceSnmpMonitoringEnableBatchCmd = &cobra.Command{
		Use:   "enable-batch [network_device_id...]",
		Short: "Subscribe several network devices to SNMP monitoring",
		Long: `Subscribe several network devices to SNMP monitoring in a single call. Devices
are selected by id (positional arguments), by site or by fabric; at least one
selection must be given.

Optional Arguments:
  network_device_id...   Numeric ids of the network devices to subscribe

Optional Flags:
  --site-id     Numeric site ids whose devices are subscribed (repeatable)
  --fabric-id   Numeric fabric ids whose devices are subscribed (repeatable)

Examples:
  # Subscribe three devices
  metalcloud-cli network-device snmp-monitoring enable-batch 12345 12346 12347

  # Subscribe every device of site 1
  metalcloud-cli network-device snmp-monitoring enable-batch --site-id 1`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceSnmpMonitoringEnableBatch(cmd.Context(), args,
				networkDeviceSubFlags.siteIds, networkDeviceSubFlags.fabricIds)
		},
	}

	networkDeviceSnmpMonitoringDisableBatchCmd = &cobra.Command{
		Use:   "disable-batch [network_device_id...]",
		Short: "Unsubscribe several network devices from SNMP monitoring",
		Long: `Unsubscribe several network devices from SNMP monitoring in a single call.
Devices are selected by id (positional arguments), by site or by fabric; at
least one selection must be given.

Optional Arguments:
  network_device_id...   Numeric ids of the network devices to unsubscribe

Optional Flags:
  --site-id     Numeric site ids whose devices are unsubscribed (repeatable)
  --fabric-id   Numeric fabric ids whose devices are unsubscribed (repeatable)

Examples:
  # Unsubscribe three devices
  metalcloud-cli network-device snmp-monitoring disable-batch 12345 12346 12347

  # Unsubscribe every device of fabric 4
  metalcloud-cli network-device snmp-monitoring disable-batch --fabric-id 4`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceSnmpMonitoringDisableBatch(cmd.Context(), args,
				networkDeviceSubFlags.siteIds, networkDeviceSubFlags.fabricIds)
		},
	}

	networkDeviceSnmpAgentInfoCmd = &cobra.Command{
		Use:   "agent-info [network_device_id...]",
		Short: "Show which monitoring agent polls each network device",
		Long: `Show the monitoring agent allocated to each of the selected network devices.
Devices are selected by id (positional arguments), by site or by fabric; with
no selection the whole allocation map is returned.

Optional Arguments:
  network_device_id...   Numeric ids of the network devices to look up

Optional Flags:
  --site-id     Restrict the lookup to these numeric site ids (repeatable)
  --fabric-id   Restrict the lookup to these numeric fabric ids (repeatable)

Examples:
  # Show the agents of two devices
  metalcloud-cli network-device snmp-monitoring agent-info 12345 12346

  # Show the agents of every device in site 1
  metalcloud-cli network-device snmp-monitoring agent-info --site-id 1`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceSnmpAgentInfo(cmd.Context(), args,
				networkDeviceSubFlags.siteIds, networkDeviceSubFlags.fabricIds)
		},
	}

	// --- SNMP service on the device ---------------------------------------

	networkDeviceSnmpServiceCmd = &cobra.Command{
		Use:   "snmp-service [command]",
		Short: "Manage the SNMP agent running on a network device",
		Long: `Enable or disable the SNMP agent running on the network device itself.

Both operations are asynchronous and return the job that carries them out.`,
	}

	networkDeviceSnmpServiceEnableCmd = &cobra.Command{
		Use:   "enable <network_device_id>",
		Short: "Enable the SNMP agent on a network device",
		Long: `Enable the SNMP agent running on the network device, optionally overriding the
listening port, the community string and the contact information.

Required Arguments:
  network_device_id   The numeric id or label of the network device

Optional Flags:
  --port        SNMP listening port
  --community   SNMP community string
  --contact     SNMP contact information

Examples:
  # Enable the SNMP agent with the stored defaults
  metalcloud-cli network-device snmp-service enable 12345

  # Enable it on a custom port and community
  metalcloud-cli network-device snmp-service enable 12345 --port 1161 --community public`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var port *int32
			if cmd.Flags().Changed("port") {
				port = &networkDeviceSubFlags.snmpPort
			}

			var community *string
			if cmd.Flags().Changed("community") {
				community = &networkDeviceSubFlags.snmpCommunity
			}

			var contact *string
			if cmd.Flags().Changed("contact") {
				contact = &networkDeviceSubFlags.snmpContact
			}

			return network_device.NetworkDeviceSnmpServiceEnable(cmd.Context(), args[0], port, community, contact)
		},
	}

	networkDeviceSnmpServiceDisableCmd = &cobra.Command{
		Use:   "disable <network_device_id>",
		Short: "Disable the SNMP agent on a network device",
		Long: `Disable the SNMP agent running on the network device.

Required Arguments:
  network_device_id   The numeric id or label of the network device

Examples:
  # Disable the SNMP agent of device 12345
  metalcloud-cli network-device snmp-service disable 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceSnmpServiceDisable(cmd.Context(), args[0])
		},
	}

	networkDeviceDisableSyslogCmd = &cobra.Command{
		Use:   "disable-syslog <network_device_id>",
		Short: "Unsubscribe a network device from remote syslog",
		Long: `Unsubscribe a network device from remote syslog, the counterpart of
'enable-syslog'. The operation is asynchronous and returns the job that carries
it out.

Required Arguments:
  network_device_id   The numeric id or label of the network device

Examples:
  # Disable syslog for device 12345
  metalcloud-cli network-device disable-syslog 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceDisableSyslog(cmd.Context(), args[0])
		},
	}

	networkDeviceHealthSummaryCmd = &cobra.Command{
		Use:   "health-summary <network_device_id>",
		Short: "Show the health assessment of a network device",
		Long: `Show the accumulated health assessment of a network device: overall severity,
trend, suspected root causes, key findings and detected issues.

Required Arguments:
  network_device_id   The numeric id or label of the network device

Examples:
  # Show the health summary of device 12345
  metalcloud-cli network-device health-summary 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceHealthSummary(cmd.Context(), args[0])
		},
	}

	networkDeviceSetHealthMonitoringFilterCmd = &cobra.Command{
		Use:   "set-health-monitoring-filter",
		Short: "Set the network device filter of a health monitoring socket",
		Long: `Set the network device filter applied to the health monitoring WebSocket stream
of the given socket. Only devices matching the filter are streamed to that
socket. The endpoint returns no content.

Required Flags:
  --socket-id             The id of the health monitoring WebSocket connection

Optional Flags:
  --filter-id             Restrict the stream to these network device ids
  --filter-site-id        Restrict the stream to these site ids
  --filter-status         Restrict the stream to these device statuses
  --filter-health-status  Restrict the stream to these health statuses

Examples:
  # Stream only the devices of site 1
  metalcloud-cli network-device set-health-monitoring-filter --socket-id abc123 --filter-site-id 1

  # Stream only unhealthy devices
  metalcloud-cli network-device set-health-monitoring-filter --socket-id abc123 --filter-health-status unhealthy`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceSetHealthMonitoringFilter(cmd.Context(),
				networkDeviceSubFlags.socketId,
				networkDeviceSubFlags.filterId,
				networkDeviceSubFlags.filterSiteId,
				networkDeviceSubFlags.filterStatus,
				networkDeviceSubFlags.filterHealthStatus)
		},
	}

	networkDeviceStatisticsCmd = &cobra.Command{
		Use:   "statistics",
		Short: "Show global network device counters",
		Long: `Show the global network device counters: how many network devices and how many
network device ports are registered.

Examples:
  # Show the network device statistics
  metalcloud-cli network-device statistics`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceStatisticsGet(cmd.Context())
		},
	}

	// --- drift & snapshots -------------------------------------------------

	networkDeviceDriftCmd = &cobra.Command{
		Use:   "drift [command]",
		Short: "Inspect the configuration drift of a network device",
		Long: `Inspect and acknowledge the configuration drift detected between the running
configuration of a network device and its target snapshot.`,
	}

	networkDeviceDriftListCmd = &cobra.Command{
		Use:     "list <network_device_id>",
		Aliases: []string{"ls"},
		Short:   "List the configuration drift of a network device",
		Long: `List the configuration drift entries recorded for a network device.

Required Arguments:
  network_device_id   The numeric id or label of the network device

Examples:
  # List the drift history of device 12345
  metalcloud-cli network-device drift list 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceDriftList(cmd.Context(), args[0])
		},
	}

	networkDeviceDriftGetCmd = &cobra.Command{
		Use:     "get <network_device_id> <drift_id>",
		Aliases: []string{"show"},
		Short:   "Get one configuration drift entry",
		Long: `Display one configuration drift entry of a network device, including the
configuration difference that was detected.

Required Arguments:
  network_device_id   The numeric id or label of the network device
  drift_id            The numeric id of the drift entry

Examples:
  # Show drift entry 9 of device 12345
  metalcloud-cli network-device drift get 12345 9`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceDriftGet(cmd.Context(), args[0], args[1])
		},
	}

	networkDeviceDriftAcknowledgeCmd = &cobra.Command{
		Use:     "acknowledge <network_device_id> <drift_id>",
		Aliases: []string{"ack"},
		Short:   "Acknowledge a configuration drift entry",
		Long: `Mark one configuration drift entry of a network device as reviewed, recording
who acknowledged it and when.

Required Arguments:
  network_device_id   The numeric id or label of the network device
  drift_id            The numeric id of the drift entry

Examples:
  # Acknowledge drift entry 9 of device 12345
  metalcloud-cli network-device drift acknowledge 12345 9`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceDriftAcknowledge(cmd.Context(), args[0], args[1])
		},
	}

	networkDeviceSnapshotCmd = &cobra.Command{
		Use:     "snapshot [command]",
		Aliases: []string{"snapshots"},
		Short:   "Inspect the configuration snapshots of a network device",
	}

	networkDeviceSnapshotListCmd = &cobra.Command{
		Use:     "list <network_device_id>",
		Aliases: []string{"ls"},
		Short:   "List the configuration snapshots of a network device",
		Long: `List the configuration snapshots stored for a network device.

Required Arguments:
  network_device_id   The numeric id or label of the network device

Optional Flags:
  --kind   Restrict the listing to one snapshot class

Examples:
  # List all snapshots of device 12345
  metalcloud-cli network-device snapshot list 12345

  # List only the backup snapshots
  metalcloud-cli network-device snapshot list 12345 --kind backup`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceSnapshotList(cmd.Context(), args[0], networkDeviceSubFlags.snapshotKind)
		},
	}

	networkDeviceSyncTargetSnapshotCmd = &cobra.Command{
		Use:   "sync-target-snapshot <network_device_id>",
		Short: "Accept the current configuration as the drift target",
		Long: `Point the drift detection target snapshot of a network device at its latest
snapshot, clearing the drift currently reported for the device.

Required Arguments:
  network_device_id   The numeric id or label of the network device

Examples:
  # Accept the running configuration of device 12345 as the new target
  metalcloud-cli network-device sync-target-snapshot 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceSyncTargetSnapshot(cmd.Context(), args[0])
		},
	}

	// --- lifecycle actions -------------------------------------------------

	networkDeviceReplaceCmd = &cobra.Command{
		Use:   "replace <network_device_id> <new_network_device_id>",
		Short: "Replace a network device with another one",
		Long: `Replace a network device with a replacement device, moving its configuration
and connections over to the new device.

Required Arguments:
  network_device_id       The numeric id or label of the device being replaced
  new_network_device_id   The numeric id or label of the replacement device

Examples:
  # Replace device 12345 with device 12399
  metalcloud-cli network-device replace 12345 12399`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceReplace(cmd.Context(), args[0], args[1])
		},
	}

	networkDeviceReProvisionCmd = &cobra.Command{
		Use:   "re-provision <network_device_id>",
		Short: "Re-run provisioning on a network device",
		Long: `Re-run provisioning on a network device. The operation is asynchronous and
returns the job that carries it out.

Required Arguments:
  network_device_id   The numeric id or label of the network device

Required Flags:
  --type              The type of re-provisioning to perform

Examples:
  # Re-provision device 12345
  metalcloud-cli network-device re-provision 12345 --type full`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceReProvision(cmd.Context(), args[0], networkDeviceSubFlags.reprovisionType)
		},
	}

	networkDeviceReturnToPlannedCmd = &cobra.Command{
		Use:   "return-to-planned <network_device_id>",
		Short: "Move an archived network device back to planned",
		Long: `Move an archived network device back to the planned state so that it can be
installed again.

Required Arguments:
  network_device_id   The numeric id or label of the network device

Optional Flags:
  --config-source     Source of the return configuration (management MAC, serial
                      number, OS template override)
                      Values: 'pipe' for stdin input, or path to a JSON/YAML file

Examples:
  # Return device 12345 to planned
  metalcloud-cli network-device return-to-planned 12345

  # Return it supplying the identifying MAC address
  echo '{"managementMAC":"AA:BB:CC:DD:EE:FF"}' | metalcloud-cli network-device return-to-planned 12345 --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var config []byte
			if networkDeviceSubFlags.configSource != "" {
				var err error
				config, err = utils.ReadConfigFromPipeOrFile(networkDeviceSubFlags.configSource)
				if err != nil {
					return err
				}
			}

			return network_device.NetworkDeviceReturnToPlanned(cmd.Context(), args[0], config)
		},
	}

	networkDeviceReturnToPlannedConfigExampleCmd = &cobra.Command{
		Use:          "return-to-planned-config-example",
		Short:        "Example configuration for the return-to-planned command",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceReturnToPlannedConfigExample(cmd.Context())
		},
	}

	networkDeviceRevertDefectiveStateCmd = &cobra.Command{
		Use:     "revert-defective-state <network_device_id>",
		Aliases: []string{"revert-failed-state"},
		Short:   "Take a network device out of the defective state",
		Long: `Take a network device out of the defective state and back to its previous status,
the counterpart of 'set-defective'.

Required Arguments:
  network_device_id   The numeric id or label of the network device

Examples:
  # Revert the defective state of device 12345
  metalcloud-cli network-device revert-defective-state 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceRevertDefectiveState(cmd.Context(), args[0])
		},
	}

	networkDeviceStartRegistrationCmd = &cobra.Command{
		Use:   "start-registration <network_device_id>",
		Short: "Start the onboarding registration of a network device",
		Long: `Start the onboarding registration of a network device. Supported only by the
drivers that report the 'start-registration' activation action; see
'network-device driver-capabilities'.

Required Arguments:
  network_device_id   The numeric id or label of the network device

Examples:
  # Start the registration of device 12345
  metalcloud-cli network-device start-registration 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceStartRegistration(cmd.Context(), args[0])
		},
	}

	networkDeviceMarkInstallationReadyCmd = &cobra.Command{
		Use:   "mark-installation-ready <network_device_id>",
		Short: "Mark the physical installation of a network device as done",
		Long: `Mark the physical installation of a network device as complete so that
onboarding can continue. Supported only by the drivers that report the
'mark-install-ready' activation action; see 'network-device driver-capabilities'.

Required Arguments:
  network_device_id   The numeric id or label of the network device

Examples:
  # Mark device 12345 as installed
  metalcloud-cli network-device mark-installation-ready 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceMarkInstallationReady(cmd.Context(), args[0])
		},
	}

	networkDeviceRunExtensionCmd = &cobra.Command{
		Use:   "run-extension <network_device_id>",
		Short: "Run an extension against a network device",
		Long: `Run an extension against a network device. The operation is asynchronous and
returns the job that carries it out.

Required Arguments:
  network_device_id   The numeric id or label of the network device

Required Flags:
  --config-source     Source of the extension invocation. Both 'extensionId' and
                      'inputArguments' must be present; pass an empty object (or
                      null) when the extension takes no arguments.
                      Values: 'pipe' for stdin input, or path to a JSON/YAML file

Examples:
  # Run an extension described in a file
  metalcloud-cli network-device run-extension 12345 --config-source extension.json

  # Run extension 7 with no arguments
  echo '{"extensionId":7,"inputArguments":{}}' | metalcloud-cli network-device run-extension 12345 --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(networkDeviceSubFlags.configSource)
			if err != nil {
				return err
			}

			return network_device.NetworkDeviceRunExtension(cmd.Context(), args[0], config)
		},
	}

	networkDeviceRunExtensionConfigExampleCmd = &cobra.Command{
		Use:          "run-extension-config-example",
		Short:        "Example configuration for the run-extension command",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceRunExtensionConfigExample(cmd.Context())
		},
	}

	// --- vendors & drivers -------------------------------------------------

	networkDeviceVendorCmd = &cobra.Command{
		Use:     "vendor [command]",
		Aliases: []string{"vendors"},
		Short:   "Manage the network device vendor profiles",
		Long: `Manage the per-vendor profiles that drive SNMP monitoring, health checks and
configuration backups for each network device driver.`,
	}

	networkDeviceVendorListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List the network device vendor profiles",
		Long: `List the vendor profiles, one per network device driver.

Optional Flags:
  --filter-kind   Restrict the listing to these drivers

Examples:
  # List all vendor profiles
  metalcloud-cli network-device vendor list

  # List the SONiC profile only
  metalcloud-cli network-device vendor list --filter-kind sonic_enterprise`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceVendorList(cmd.Context(), networkDeviceSubFlags.filterKind)
		},
	}

	networkDeviceVendorGetCmd = &cobra.Command{
		Use:     "get <vendor_id>",
		Aliases: []string{"show"},
		Short:   "Get one network device vendor profile",
		Long: `Display one vendor profile: its SNMP OID groups, health check rules and the
optional files included in configuration backups.

Required Arguments:
  vendor_id   The numeric id of the vendor profile

Examples:
  # Show vendor profile 3
  metalcloud-cli network-device vendor get 3`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceVendorGet(cmd.Context(), args[0])
		},
	}

	networkDeviceVendorUpdateCmd = &cobra.Command{
		Use:   "update <vendor_id>",
		Short: "Update one network device vendor profile",
		Long: `Update the SNMP OID groups, health check rules and backup file list of one
vendor profile.

Required Arguments:
  vendor_id         The numeric id of the vendor profile

Required Flags:
  --config-source   Source of the vendor configuration
                    Values: 'pipe' for stdin input, or path to a JSON/YAML file

Examples:
  # Update vendor profile 3 from a file
  metalcloud-cli network-device vendor update 3 --config-source vendor.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(networkDeviceSubFlags.configSource)
			if err != nil {
				return err
			}

			return network_device.NetworkDeviceVendorUpdate(cmd.Context(), args[0], config)
		},
	}

	networkDeviceVendorConfigExampleCmd = &cobra.Command{
		Use:          "config-example",
		Short:        "Example configuration for updating a vendor profile",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceVendorConfigExample(cmd.Context())
		},
	}

	networkDeviceDriverCapabilitiesCmd = &cobra.Command{
		Use:   "driver-capabilities",
		Short: "List the onboarding actions supported by each driver",
		Long: `List, for each network device driver, the onboarding actions it supports
(for example 'start-registration' or 'mark-install-ready').

Optional Flags:
  --search   Restrict the listing to the drivers matching this text

Examples:
  # List the capabilities of every driver
  metalcloud-cli network-device driver-capabilities

  # List the SONiC driver capabilities
  metalcloud-cli network-device driver-capabilities --search sonic`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceDriverCapabilities(cmd.Context(), networkDeviceSubFlags.search)
		},
	}
)

func init() {
	// ports
	networkDeviceCmd.AddCommand(networkDevicePortCmd)

	networkDevicePortCmd.AddCommand(networkDevicePortAddCmd)
	networkDevicePortAddCmd.Flags().StringVar(&networkDeviceSubFlags.configSource, "config-source", "", "Source of the new interface configuration. Can be 'pipe' or path to a JSON file.")
	networkDevicePortAddCmd.MarkFlagRequired("config-source")

	networkDevicePortCmd.AddCommand(networkDevicePortConfigExampleCmd)
	networkDevicePortCmd.AddCommand(networkDevicePortRemoveCmd)
	networkDevicePortCmd.AddCommand(networkDevicePortGetConfigCmd)
	networkDevicePortCmd.AddCommand(networkDevicePortLiveStatusCmd)

	networkDevicePortCmd.AddCommand(networkDevicePortIpCmd)
	networkDevicePortIpCmd.AddCommand(networkDevicePortIpListCmd)
	networkDevicePortIpCmd.AddCommand(networkDevicePortIpGetCmd)
	networkDevicePortIpCmd.AddCommand(networkDevicePortIpRemoveCmd)

	networkDevicePortIpCmd.AddCommand(networkDevicePortIpReplaceCmd)
	networkDevicePortIpReplaceCmd.Flags().StringVar(&networkDeviceSubFlags.configSource, "config-source", "", "Source of the address list. Can be 'pipe' or path to a JSON file.")
	networkDevicePortIpReplaceCmd.MarkFlagRequired("config-source")

	networkDevicePortIpCmd.AddCommand(networkDevicePortIpConfigExampleCmd)

	networkDevicePortCmd.AddCommand(networkDevicePortVirtualFunctionCmd)
	networkDevicePortVirtualFunctionCmd.AddCommand(networkDevicePortVirtualFunctionListCmd)
	networkDevicePortVirtualFunctionCmd.AddCommand(networkDevicePortVirtualFunctionGetCmd)

	// breakouts
	networkDeviceCmd.AddCommand(networkDeviceBreakoutCmd)
	networkDeviceBreakoutCmd.AddCommand(networkDeviceBreakoutListCmd)
	networkDeviceBreakoutCmd.AddCommand(networkDeviceBreakoutGetCmd)
	networkDeviceBreakoutCmd.AddCommand(networkDeviceBreakoutGetConfigCmd)

	networkDeviceBreakoutCmd.AddCommand(networkDeviceBreakoutCreateCmd)
	networkDeviceBreakoutCreateCmd.Flags().StringVar(&networkDeviceSubFlags.configSource, "config-source", "", "Source of the new breakout configuration. Can be 'pipe' or path to a JSON file.")
	networkDeviceBreakoutCreateCmd.MarkFlagRequired("config-source")

	networkDeviceBreakoutCmd.AddCommand(networkDeviceBreakoutUpdateConfigCmd)
	networkDeviceBreakoutUpdateConfigCmd.Flags().StringVar(&networkDeviceSubFlags.configSource, "config-source", "", "Source of the breakout group configuration. Can be 'pipe' or path to a JSON file.")
	networkDeviceBreakoutUpdateConfigCmd.MarkFlagRequired("config-source")

	networkDeviceBreakoutCmd.AddCommand(networkDeviceBreakoutDeleteCmd)
	networkDeviceBreakoutCmd.AddCommand(networkDeviceBreakoutConfigExampleCmd)
	networkDeviceBreakoutCmd.AddCommand(networkDeviceBreakoutUpdateConfigExampleCmd)

	// virtual functions
	networkDeviceCmd.AddCommand(networkDeviceVirtualFunctionCmd)
	networkDeviceVirtualFunctionCmd.AddCommand(networkDeviceVirtualFunctionListCmd)
	networkDeviceVirtualFunctionCmd.AddCommand(networkDeviceVirtualFunctionGetCmd)

	// secrets
	networkDeviceCmd.AddCommand(networkDeviceSecretCmd)
	networkDeviceSecretCmd.AddCommand(networkDeviceSecretListCmd)

	networkDeviceSecretCmd.AddCommand(networkDeviceSecretSetCmd)
	networkDeviceSecretSetCmd.Flags().StringVar(&networkDeviceSubFlags.secretValue, "value", "", "The secret value to store.")
	networkDeviceSecretSetCmd.MarkFlagRequired("value")

	networkDeviceSecretCmd.AddCommand(networkDeviceSecretGetCredentialsCmd)
	networkDeviceSecretCmd.AddCommand(networkDeviceSecretRemoveCmd)
	networkDeviceSecretCmd.AddCommand(networkDeviceSecretRemoveAllCmd)

	// SNMP monitoring
	networkDeviceCmd.AddCommand(networkDeviceSnmpMonitoringCmd)
	networkDeviceSnmpMonitoringCmd.AddCommand(networkDeviceSnmpMonitoringEnableCmd)
	networkDeviceSnmpMonitoringCmd.AddCommand(networkDeviceSnmpMonitoringDisableCmd)

	networkDeviceSnmpMonitoringCmd.AddCommand(networkDeviceSnmpMonitoringEnableBatchCmd)
	networkDeviceSnmpMonitoringEnableBatchCmd.Flags().StringSliceVar(&networkDeviceSubFlags.siteIds, "site-id", nil, "Numeric site ids whose network devices are subscribed. Repeatable.")
	networkDeviceSnmpMonitoringEnableBatchCmd.Flags().StringSliceVar(&networkDeviceSubFlags.fabricIds, "fabric-id", nil, "Numeric fabric ids whose network devices are subscribed. Repeatable.")

	networkDeviceSnmpMonitoringCmd.AddCommand(networkDeviceSnmpMonitoringDisableBatchCmd)
	networkDeviceSnmpMonitoringDisableBatchCmd.Flags().StringSliceVar(&networkDeviceSubFlags.siteIds, "site-id", nil, "Numeric site ids whose network devices are unsubscribed. Repeatable.")
	networkDeviceSnmpMonitoringDisableBatchCmd.Flags().StringSliceVar(&networkDeviceSubFlags.fabricIds, "fabric-id", nil, "Numeric fabric ids whose network devices are unsubscribed. Repeatable.")

	networkDeviceSnmpMonitoringCmd.AddCommand(networkDeviceSnmpAgentInfoCmd)
	networkDeviceSnmpAgentInfoCmd.Flags().StringSliceVar(&networkDeviceSubFlags.siteIds, "site-id", nil, "Restrict the lookup to these numeric site ids. Repeatable.")
	networkDeviceSnmpAgentInfoCmd.Flags().StringSliceVar(&networkDeviceSubFlags.fabricIds, "fabric-id", nil, "Restrict the lookup to these numeric fabric ids. Repeatable.")

	// SNMP service
	networkDeviceCmd.AddCommand(networkDeviceSnmpServiceCmd)
	networkDeviceSnmpServiceCmd.AddCommand(networkDeviceSnmpServiceEnableCmd)
	networkDeviceSnmpServiceEnableCmd.Flags().Int32Var(&networkDeviceSubFlags.snmpPort, "port", 161, "SNMP listening port.")
	networkDeviceSnmpServiceEnableCmd.Flags().StringVar(&networkDeviceSubFlags.snmpCommunity, "community", "", "SNMP community string.")
	networkDeviceSnmpServiceEnableCmd.Flags().StringVar(&networkDeviceSubFlags.snmpContact, "contact", "", "SNMP contact information.")

	networkDeviceSnmpServiceCmd.AddCommand(networkDeviceSnmpServiceDisableCmd)

	// syslog, health and statistics
	networkDeviceCmd.AddCommand(networkDeviceDisableSyslogCmd)
	networkDeviceCmd.AddCommand(networkDeviceHealthSummaryCmd)

	networkDeviceCmd.AddCommand(networkDeviceSetHealthMonitoringFilterCmd)
	networkDeviceSetHealthMonitoringFilterCmd.Flags().StringVar(&networkDeviceSubFlags.socketId, "socket-id", "", "The id of the health monitoring WebSocket connection.")
	networkDeviceSetHealthMonitoringFilterCmd.Flags().StringSliceVar(&networkDeviceSubFlags.filterId, "filter-id", nil, "Restrict the stream to these network device ids.")
	networkDeviceSetHealthMonitoringFilterCmd.Flags().StringSliceVar(&networkDeviceSubFlags.filterSiteId, "filter-site-id", nil, "Restrict the stream to these site ids.")
	networkDeviceSetHealthMonitoringFilterCmd.Flags().StringSliceVar(&networkDeviceSubFlags.filterStatus, "filter-status", nil, "Restrict the stream to these network device statuses.")
	networkDeviceSetHealthMonitoringFilterCmd.Flags().StringSliceVar(&networkDeviceSubFlags.filterHealthStatus, "filter-health-status", nil, "Restrict the stream to these health statuses.")
	networkDeviceSetHealthMonitoringFilterCmd.MarkFlagRequired("socket-id")

	networkDeviceCmd.AddCommand(networkDeviceStatisticsCmd)

	// drift & snapshots
	networkDeviceCmd.AddCommand(networkDeviceDriftCmd)
	networkDeviceDriftCmd.AddCommand(networkDeviceDriftListCmd)
	networkDeviceDriftCmd.AddCommand(networkDeviceDriftGetCmd)
	networkDeviceDriftCmd.AddCommand(networkDeviceDriftAcknowledgeCmd)

	networkDeviceCmd.AddCommand(networkDeviceSnapshotCmd)
	networkDeviceSnapshotCmd.AddCommand(networkDeviceSnapshotListCmd)
	networkDeviceSnapshotListCmd.Flags().StringVar(&networkDeviceSubFlags.snapshotKind, "kind", "", "Restrict the listing to one snapshot class.")

	networkDeviceCmd.AddCommand(networkDeviceSyncTargetSnapshotCmd)

	// lifecycle actions
	networkDeviceCmd.AddCommand(networkDeviceReplaceCmd)

	networkDeviceCmd.AddCommand(networkDeviceReProvisionCmd)
	networkDeviceReProvisionCmd.Flags().StringVar(&networkDeviceSubFlags.reprovisionType, "type", "", "The type of re-provisioning to perform.")
	networkDeviceReProvisionCmd.MarkFlagRequired("type")

	networkDeviceCmd.AddCommand(networkDeviceReturnToPlannedCmd)
	networkDeviceReturnToPlannedCmd.Flags().StringVar(&networkDeviceSubFlags.configSource, "config-source", "", "Source of the return configuration. Can be 'pipe' or path to a JSON file.")

	networkDeviceCmd.AddCommand(networkDeviceReturnToPlannedConfigExampleCmd)
	networkDeviceCmd.AddCommand(networkDeviceRevertDefectiveStateCmd)
	networkDeviceCmd.AddCommand(networkDeviceStartRegistrationCmd)
	networkDeviceCmd.AddCommand(networkDeviceMarkInstallationReadyCmd)

	networkDeviceCmd.AddCommand(networkDeviceRunExtensionCmd)
	networkDeviceRunExtensionCmd.Flags().StringVar(&networkDeviceSubFlags.configSource, "config-source", "", "Source of the extension invocation. Can be 'pipe' or path to a JSON file.")
	networkDeviceRunExtensionCmd.MarkFlagRequired("config-source")

	networkDeviceCmd.AddCommand(networkDeviceRunExtensionConfigExampleCmd)

	// vendors & drivers
	networkDeviceCmd.AddCommand(networkDeviceVendorCmd)
	networkDeviceVendorCmd.AddCommand(networkDeviceVendorListCmd)
	networkDeviceVendorListCmd.Flags().StringSliceVar(&networkDeviceSubFlags.filterKind, "filter-kind", nil, "Filter the vendor profiles by network device driver.")

	networkDeviceVendorCmd.AddCommand(networkDeviceVendorGetCmd)

	networkDeviceVendorCmd.AddCommand(networkDeviceVendorUpdateCmd)
	networkDeviceVendorUpdateCmd.Flags().StringVar(&networkDeviceSubFlags.configSource, "config-source", "", "Source of the vendor configuration. Can be 'pipe' or path to a JSON file.")
	networkDeviceVendorUpdateCmd.MarkFlagRequired("config-source")

	networkDeviceVendorCmd.AddCommand(networkDeviceVendorConfigExampleCmd)

	networkDeviceCmd.AddCommand(networkDeviceDriverCapabilitiesCmd)
	networkDeviceDriverCapabilitiesCmd.Flags().StringVar(&networkDeviceSubFlags.search, "search", "", "Restrict the listing to the drivers matching this text.")
}
