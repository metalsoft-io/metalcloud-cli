package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/server"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
	"github.com/spf13/cobra"
)

// Server commands
var (
	serverFlags = struct {
		showCredentials   bool
		filterStatus      []string
		filterType        []string
		configSource      string
		powerAction       string
		siteId            int
		managementAddress string
		username          string
		password          string
		serialNumber      string
		model             string
		vendor            string
	}{}

	serverCmd = &cobra.Command{
		Use:     "server [command]",
		Aliases: []string{"srv", "servers"},
		Short:   "Server management",
		Long: `Server management commands.

This command group provides comprehensive server management capabilities including
registration, power control, firmware management, and monitoring. Servers can be
managed individually or in bulk operations.

Available command categories:
  - Basic operations: list, get, register, update, delete
  - Power management: power, power-status  
  - Maintenance: re-register, factory-reset, archive
  - Security: update-ipmi-credentials, enable-snmp, enable-syslog
  - Remote access: vnc-info, console-info
  - Firmware: firmware subcommands for component management and upgrades
  - Information: capabilities

Use "metalcloud-cli server [command] --help" for detailed information about each command.
`,
	}

	serverListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List servers",
		Long: `List all servers in the MetalSoft infrastructure.

This command displays information about all servers including their IDs, site locations, 
types, UUIDs, serial numbers, management addresses, vendors, models, and current status.

Optional Flags:
  --show-credentials     Display server IPMI credentials (username and password)
  --filter-status        Filter servers by status (e.g., active, registered, provisioning)
  --filter-type          Filter servers by type ID

Examples:
  # List all servers
  metalcloud-cli server list

  # List servers with IPMI credentials
  metalcloud-cli server list --show-credentials

  # Filter servers by status
  metalcloud-cli server list --filter-status active,registered

  # Filter servers by type
  metalcloud-cli server list --filter-type 1,2,3

  # Combine filters
  metalcloud-cli server list --filter-status active --filter-type 1
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerList(cmd.Context(),
				serverFlags.showCredentials,
				serverFlags.filterStatus,
				serverFlags.filterType)
		},
	}

	serverGetCmd = &cobra.Command{
		Use:     "get server_id",
		Aliases: []string{"show"},
		Short:   "Get detailed server information",
		Long: `Get detailed information for a specific server.

This command retrieves comprehensive information about a server including its configuration,
status, hardware details, and optionally IPMI credentials.

Required Arguments:
  server_id              The ID of the server to retrieve information for

Optional Flags:
  --show-credentials     Include IPMI credentials (username and password) in the output

Examples:
  # Get basic server information
  metalcloud-cli server get 123

  # Get server information including IPMI credentials
  metalcloud-cli server get 123 --show-credentials
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerGet(cmd.Context(), args[0],
				serverFlags.showCredentials)
		},
	}

	serverRegisterCmd = &cobra.Command{
		Use:     "register",
		Aliases: []string{"new"},
		Short:   "Register a new server in MetalSoft",
		Long: `Register a new server in MetalSoft.

You can provide the server configuration either via command-line flags or by specifying a configuration source using the --config-source flag. 
The configuration source can be a path to a JSON file or 'pipe' to read from standard input.

If --config-source is not provided, you must specify at least --site-id and --management-address, along with any other relevant server details such as --username, --password, --serial-number, --model, and --vendor.

Flag Dependencies:
  --config-source and --site-id are mutually exclusive
  --site-id and --management-address must be used together when not using --config-source

Required Flags (when not using --config-source):
  --site-id              Site ID where the server is located
  --management-address   IPMI/BMC management IP address

Optional Flags:
  --config-source        Source of server configuration (JSON file path or 'pipe')
  --username            IPMI/BMC username
  --password            IPMI/BMC password
  --serial-number       Server serial number
  --model               Server model
  --vendor              Server vendor

Examples:
  # Register using command line flags
  metalcloud-cli server register --site-id 1 --management-address 10.0.0.1 --username admin --password secret

  # Register with additional server details
  metalcloud-cli server register --site-id 1 --management-address 10.0.0.1 --username admin --password secret --serial-number ABC123 --model PowerEdge --vendor Dell

  # Register using JSON configuration file
  metalcloud-cli server register --config-source ./server.json

  # Register using piped JSON configuration
  cat server.json | metalcloud-cli server register --config-source pipe
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			var serverConfig sdk.RegisterServer

			// If config source is provided, use it
			if serverFlags.configSource != "" {
				config, err := utils.ReadConfigFromPipeOrFile(serverFlags.configSource)
				if err != nil {
					return err
				}
				err = utils.UnmarshalContent(config, &serverConfig)
				if err != nil {
					return err
				}
			} else {
				// Otherwise build config from command line parameters
				serverConfig = sdk.RegisterServer{
					SiteId:            float32(serverFlags.siteId),
					ManagementAddress: &serverFlags.managementAddress,
					Username:          &serverFlags.username,
					Password:          &serverFlags.password,
					SerialNumber:      &serverFlags.serialNumber,
					Model:             &serverFlags.model,
					Vendor:            &serverFlags.vendor,
				}
			}

			return server.ServerRegister(cmd.Context(), serverConfig)
		},
	}

	serverReRegisterCmd = &cobra.Command{
		Use:   "re-register server_id",
		Short: "Re-register an existing server",
		Long: `Re-register an existing server in MetalSoft.

This command initiates the re-registration process for a server that is already
registered in the system. This is typically used when a server needs to be
re-discovered or when its configuration has changed.

Required Arguments:
  server_id              The ID of the server to re-register

Examples:
  # Re-register server with ID 123
  metalcloud-cli server re-register 123
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerReRegister(cmd.Context(), args[0])
		},
	}

	serverFactoryResetCmd = &cobra.Command{
		Use:   "factory-reset server_id",
		Short: "Reset a server to factory defaults",
		Long: `Reset a server to factory defaults.

This command initiates a factory reset operation on the specified server,
restoring it to its original configuration. This operation is irreversible
and will remove all custom configurations.

Required Arguments:
  server_id              The ID of the server to factory reset

Examples:
  # Factory reset server with ID 123
  metalcloud-cli server factory-reset 123
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerFactoryReset(cmd.Context(), args[0])
		},
	}

	serverArchiveCmd = &cobra.Command{
		Use:   "archive server_id",
		Short: "Archive a server",
		Long: `Archive a server.

This command moves a server to an archived state, effectively removing it from
active use while preserving its information for historical purposes. Archived
servers are no longer available for deployment but can still be referenced.

Required Arguments:
  server_id              The ID of the server to archive

Examples:
  # Archive server with ID 123
  metalcloud-cli server archive 123
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerArchive(cmd.Context(), args[0])
		},
	}

	serverDeleteCmd = &cobra.Command{
		Use:     "delete server_id",
		Aliases: []string{"rm"},
		Short:   "Delete a server",
		Long: `Delete a server from MetalSoft.

This command permanently removes a server from the MetalSoft system. 
This operation is irreversible and will delete all server information.

Required Arguments:
  server_id              The ID of the server to delete

Examples:
  # Delete server with ID 123
  metalcloud-cli server delete 123

  # Delete server using alias
  metalcloud-cli server rm 123
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerDelete(cmd.Context(), args[0])
		},
	}

	serverPowerCmd = &cobra.Command{
		Use:   "power server_id",
		Short: "Control server power state",
		Long: `Control server power state.

This command allows you to control the power state of a server by sending
power management commands to the server's BMC/IPMI interface.

Required Arguments:
  server_id              The ID of the server to control

Required Flags:
  --action              Power action to perform

Valid Actions:
  on                    Power on the server
  off                   Hard power off the server
  reset                 Hard reset the server
  cycle                 Power cycle the server (off then on)
  soft                  Soft power off the server (graceful shutdown)

Examples:
  # Power on server
  metalcloud-cli server power 123 --action on

  # Hard power off server
  metalcloud-cli server power 123 --action off

  # Reset server
  metalcloud-cli server power 123 --action reset

  # Power cycle server
  metalcloud-cli server power 123 --action cycle

  # Soft power off (graceful shutdown)
  metalcloud-cli server power 123 --action soft
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerPower(cmd.Context(), args[0], serverFlags.powerAction)
		},
	}

	serverPowerStatusCmd = &cobra.Command{
		Use:   "power-status server_id",
		Short: "Get server power status",
		Long: `Get the current power status of a server.

This command retrieves the current power state of the specified server
from its BMC/IPMI interface.

Required Arguments:
  server_id              The ID of the server to query

Examples:
  # Get power status for server with ID 123
  metalcloud-cli server power-status 123
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerPowerStatus(cmd.Context(), args[0])
		},
	}

	serverUpdateCmd = &cobra.Command{
		Use:   "update server_id",
		Short: "Update server information",
		Long: `Update server information.

This command updates server configuration using a JSON configuration file or 
piped JSON data. The configuration must be provided via the --config-source flag.

Required Arguments:
  server_id              The ID of the server to update

Required Flags:
  --config-source        Source of the server update configuration. Can be 'pipe' or path to a JSON file

Examples:
  # Update server using JSON configuration file
  metalcloud-cli server update 123 --config-source ./server-update.json

  # Update server using piped JSON configuration
  echo '{"vendor": "Dell", "model": "PowerEdge R740"}' | metalcloud-cli server update 123 --config-source pipe
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(serverFlags.configSource)
			if err != nil {
				return err
			}
			return server.ServerUpdate(cmd.Context(), args[0], config)
		},
	}

	serverUpdateIpmiCredentialsCmd = &cobra.Command{
		Use:   "update-ipmi-credentials server_id username password",
		Short: "Update server IPMI credentials",
		Long: `Update server IPMI credentials.

This command updates the IPMI/BMC username and password for the specified server.
The credentials are used for server management operations like power control
and hardware monitoring.

Required Arguments:
  server_id              The ID of the server to update
  username               New IPMI/BMC username
  password               New IPMI/BMC password

Examples:
  # Update IPMI credentials for server with ID 123
  metalcloud-cli server update-ipmi-credentials 123 admin newpassword
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerUpdateIpmiCredentials(cmd.Context(), args[0], args[1], args[2])
		},
	}

	serverEnableSnmpCmd = &cobra.Command{
		Use:   "enable-snmp server_id",
		Short: "Enable SNMP on server",
		Long: `Enable SNMP on server.

This command enables SNMP (Simple Network Management Protocol) monitoring
on the specified server, allowing network management systems to collect
server metrics and status information.

Required Arguments:
  server_id              The ID of the server to enable SNMP on

Examples:
  # Enable SNMP for server with ID 123
  metalcloud-cli server enable-snmp 123
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerEnableSnmp(cmd.Context(), args[0])
		},
	}

	serverEnableSyslogCmd = &cobra.Command{
		Use:   "enable-syslog server_id",
		Short: "Enable remote syslog for a server",
		Long: `Enable remote syslog for a server.

This command enables remote syslog forwarding on the specified server,
allowing the server to send system log messages to a remote syslog server
for centralized logging and monitoring.

Required Arguments:
  server_id              The ID of the server to enable syslog on

Examples:
  # Enable syslog for server with ID 123
  metalcloud-cli server enable-syslog 123
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerEnableSyslog(cmd.Context(), args[0])
		},
	}

	serverVncInfoCmd = &cobra.Command{
		Use:   "vnc-info server_id",
		Short: "Get server VNC information",
		Long: `Get server VNC connection information.

This command retrieves VNC (Virtual Network Computing) connection details
for the specified server, including active sessions, maximum sessions,
port information, timeout settings, and status.

Required Arguments:
  server_id              The ID of the server to query

Examples:
  # Get VNC information for server with ID 123
  metalcloud-cli server vnc-info 123
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerVncInfo(cmd.Context(), args[0])
		},
	}

	serverRemoteConsoleInfoCmd = &cobra.Command{
		Use:   "console-info server_id",
		Short: "Get server remote console information",
		Long: `Get server remote console information.

This command retrieves remote console connection details for the specified server,
including active connections and console access information.

Required Arguments:
  server_id              The ID of the server to query

Examples:
  # Get remote console information for server with ID 123
  metalcloud-cli server console-info 123
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerRemoteConsoleInfo(cmd.Context(), args[0])
		},
	}

	serverCapabilitiesCmd = &cobra.Command{
		Use:   "capabilities server_id",
		Short: "Get server capabilities",
		Long: `Get server capabilities.

This command retrieves information about the capabilities supported by the
specified server, including firmware upgrade support, VNC capabilities,
and other available features.

Required Arguments:
  server_id              The ID of the server to query

Examples:
  # Get capabilities for server with ID 123
  metalcloud-cli server capabilities 123
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerCapabilities(cmd.Context(), args[0])
		},
	}

	// New Firmware Management Commands
	serverFirmwareCmd = &cobra.Command{
		Use:     "firmware [command]",
		Aliases: []string{"fw"},
		Short:   "Server firmware management",
		Long: `Server firmware management commands.

This command group provides comprehensive firmware management capabilities for servers
including component listing, upgrades, scheduling, and auditing. Firmware operations
can be performed on individual components or entire servers.

Available commands:
  - Information: components, component-info, inventory, fetch-versions
  - Updates: update-info, update-component
  - Upgrades: upgrade, upgrade-component, schedule-upgrade
  - Auditing: generate-audit

Use "metalcloud-cli server firmware [command] --help" for detailed information about each command.
`,
	}

	serverFirmwareComponentsListCmd = &cobra.Command{
		Use:     "components server_id",
		Aliases: []string{"list-components", "ls"},
		Short:   "List firmware components for a server",
		Long: `List firmware components for a server.

This command displays all firmware components available on the specified server,
including their IDs, names, current versions, and status information.

Required Arguments:
  server_id              The ID of the server to query

Examples:
  # List firmware components for server with ID 123
  metalcloud-cli server firmware components 123

  # Using short alias
  metalcloud-cli server fw ls 123
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerFirmwareComponentsList(cmd.Context(), args[0])
		},
	}

	serverFirmwareComponentGetCmd = &cobra.Command{
		Use:     "component-info server_id component_id",
		Aliases: []string{"get-component"},
		Short:   "Get firmware component information",
		Long: `Get detailed information for a specific firmware component.

This command retrieves comprehensive information about a firmware component
including its current version, available updates, and configuration options.

Required Arguments:
  server_id              The ID of the server
  component_id           The ID of the firmware component

Examples:
  # Get firmware component information
  metalcloud-cli server firmware component-info 123 456

  # Using alias
  metalcloud-cli server firmware get-component 123 456
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerFirmwareComponentGet(cmd.Context(), args[0], args[1])
		},
	}

	serverFirmwareComponentUpdateCmd = &cobra.Command{
		Use:   "update-component server_id component_id",
		Short: "Update firmware component settings",
		Long: `Update firmware component settings.

This command updates configuration settings for a specific firmware component
using a JSON configuration file or piped JSON data.

Required Arguments:
  server_id              The ID of the server
  component_id           The ID of the firmware component to update

Required Flags:
  --config-source        Source of the component update configuration. Can be 'pipe' or path to a JSON file

Examples:
  # Update component using JSON configuration file
  metalcloud-cli server firmware update-component 123 456 --config-source ./component-config.json

  # Update component using piped JSON configuration
  echo '{"setting": "value"}' | metalcloud-cli server firmware update-component 123 456 --config-source pipe
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(serverFlags.configSource)
			if err != nil {
				return err
			}
			return server.ServerFirmwareComponentUpdate(cmd.Context(), args[0], args[1], config)
		},
	}

	serverFirmwareUpdateInfoCmd = &cobra.Command{
		Use:   "update-info server_id",
		Short: "Update firmware information for a server",
		Long: `Update firmware information for a server.

This command refreshes and updates the firmware information for the specified server,
retrieving the latest firmware details and status from the server's BMC.

Required Arguments:
  server_id              The ID of the server to update firmware information for

Examples:
  # Update firmware information for server with ID 123
  metalcloud-cli server firmware update-info 123
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerFirmwareUpdateInfo(cmd.Context(), args[0])
		},
	}

	serverFirmwareInventoryCmd = &cobra.Command{
		Use:   "inventory server_id",
		Short: "Get firmware inventory from redfish",
		Long: `Get firmware inventory from redfish.

This command retrieves the firmware inventory for the specified server using
the Redfish API, providing detailed information about all firmware components
installed on the server.

Required Arguments:
  server_id              The ID of the server to query

Examples:
  # Get firmware inventory for server with ID 123
  metalcloud-cli server firmware inventory 123
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerFirmwareInventory(cmd.Context(), args[0])
		},
	}

	serverFirmwareUpgradeCmd = &cobra.Command{
		Use:   "upgrade server_id",
		Short: "Upgrade firmware for all components on a server",
		Long: `Upgrade firmware for all components on a server.

This command initiates a firmware upgrade process for all upgradeable components
on the specified server. The system will automatically determine which components
need updates and apply the latest available firmware versions.

Required Arguments:
  server_id              The ID of the server to upgrade

Examples:
  # Upgrade all firmware components for server with ID 123
  metalcloud-cli server firmware upgrade 123
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerFirmwareUpgrade(cmd.Context(), args[0])
		},
	}

	serverFirmwareComponentUpgradeCmd = &cobra.Command{
		Use:   "upgrade-component server_id component_id",
		Short: "Upgrade firmware for a specific component",
		Long: `Upgrade firmware for a specific component.

This command initiates a firmware upgrade for a specific firmware component
on the specified server. The upgrade configuration must be provided via
a JSON configuration file or piped JSON data.

Required Arguments:
  server_id              The ID of the server
  component_id           The ID of the firmware component to upgrade

Required Flags:
  --config-source        Source of the firmware upgrade configuration. Can be 'pipe' or path to a JSON file

Examples:
  # Upgrade component using JSON configuration file
  metalcloud-cli server firmware upgrade-component 123 456 --config-source ./upgrade-config.json

  # Upgrade component using piped JSON configuration
  echo '{"version": "2.1.0"}' | metalcloud-cli server firmware upgrade-component 123 456 --config-source pipe
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(serverFlags.configSource)
			if err != nil {
				return err
			}
			return server.ServerFirmwareComponentUpgrade(cmd.Context(), args[0], args[1], config)
		},
	}

	serverFirmwareScheduleUpgradeCmd = &cobra.Command{
		Use:   "schedule-upgrade server_id",
		Short: "Schedule a firmware upgrade for a server",
		Long: `Schedule a firmware upgrade for a server.

This command schedules a firmware upgrade for the specified server using a JSON 
configuration file or piped JSON data. The configuration should specify the 
schedule details and upgrade parameters.

Required Arguments:
  server_id              The ID of the server to schedule upgrade for

Required Flags:
  --config-source        Source of the schedule configuration. Can be 'pipe' or path to a JSON file

Examples:
  # Schedule upgrade using JSON configuration file
  metalcloud-cli server firmware schedule-upgrade 123 --config-source ./schedule-config.json

  # Schedule upgrade using piped JSON configuration
  echo '{"schedule": "2024-01-01T10:00:00Z"}' | metalcloud-cli server firmware schedule-upgrade 123 --config-source pipe
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(serverFlags.configSource)
			if err != nil {
				return err
			}
			return server.ServerFirmwareScheduleUpgrade(cmd.Context(), args[0], config)
		},
	}

	serverFirmwareFetchVersionsCmd = &cobra.Command{
		Use:   "fetch-versions server_id",
		Short: "Fetch available firmware versions for a server",
		Long: `Fetch available firmware versions for a server.

This command retrieves and displays all available firmware versions for the 
specified server's components. This helps identify which firmware updates 
are available before performing upgrades.

Required Arguments:
  server_id              The ID of the server to query

Examples:
  # Fetch available firmware versions for server with ID 123
  metalcloud-cli server firmware fetch-versions 123
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerFirmwareFetchVersions(cmd.Context(), args[0])
		},
	}

	serverFirmwareGenerateAuditCmd = &cobra.Command{
		Use:   "generate-audit",
		Short: "Generate firmware upgrade audit for servers",
		Long: `Generate firmware upgrade audit for servers.

This command generates a comprehensive firmware upgrade audit report for servers
using a JSON configuration file or piped JSON data. The audit helps identify
which servers need firmware updates and provides detailed upgrade recommendations.

Required Flags:
  --config-source        Source of the audit configuration. Can be 'pipe' or path to a JSON file

Examples:
  # Generate audit using JSON configuration file
  metalcloud-cli server firmware generate-audit --config-source ./audit-config.json

  # Generate audit using piped JSON configuration
  echo '{"servers": [123, 456]}' | metalcloud-cli server firmware generate-audit --config-source pipe
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(serverFlags.configSource)
			if err != nil {
				return err
			}
			return server.ServerFirmwareGenerateAudit(cmd.Context(), config)
		},
	}
)

func init() {
	rootCmd.AddCommand(serverCmd)

	// Server commands
	serverCmd.AddCommand(serverListCmd)
	serverListCmd.Flags().BoolVar(&serverFlags.showCredentials, "show-credentials", false, "If set returns the server IPMI credentials.")
	serverListCmd.Flags().StringSliceVar(&serverFlags.filterStatus, "filter-status", nil, "Filter the result by server status.")
	serverListCmd.Flags().StringSliceVar(&serverFlags.filterType, "filter-type", nil, "Filter the result by server type.")

	serverCmd.AddCommand(serverGetCmd)
	serverGetCmd.Flags().BoolVar(&serverFlags.showCredentials, "show-credentials", false, "If set returns the server IPMI credentials.")

	serverCmd.AddCommand(serverRegisterCmd)
	serverRegisterCmd.Flags().StringVar(&serverFlags.configSource, "config-source", "", "Source of the new server configuration. Can be 'pipe' or path to a JSON file.")
	serverRegisterCmd.Flags().IntVar(&serverFlags.siteId, "site-id", 0, "Site ID")
	serverRegisterCmd.Flags().StringVar(&serverFlags.managementAddress, "management-address", "", "Management address")
	serverRegisterCmd.Flags().StringVar(&serverFlags.username, "username", "", "Username")
	serverRegisterCmd.Flags().StringVar(&serverFlags.password, "password", "", "Password")
	serverRegisterCmd.Flags().StringVar(&serverFlags.serialNumber, "serial-number", "", "Serial number")
	serverRegisterCmd.Flags().StringVar(&serverFlags.model, "model", "", "Model")
	serverRegisterCmd.Flags().StringVar(&serverFlags.vendor, "vendor", "", "Vendor")
	serverRegisterCmd.MarkFlagsOneRequired("config-source", "site-id")
	serverRegisterCmd.MarkFlagsMutuallyExclusive("config-source", "site-id")
	serverRegisterCmd.MarkFlagsRequiredTogether("site-id", "management-address")

	serverCmd.AddCommand(serverReRegisterCmd)

	serverCmd.AddCommand(serverFactoryResetCmd)

	serverCmd.AddCommand(serverArchiveCmd)

	serverCmd.AddCommand(serverDeleteCmd)

	serverCmd.AddCommand(serverPowerCmd)
	serverPowerCmd.Flags().StringVar(&serverFlags.powerAction, "action", "", "Power action: on, off, reset, cycle, soft")
	serverPowerCmd.MarkFlagsOneRequired("action")

	serverCmd.AddCommand(serverPowerStatusCmd)

	serverCmd.AddCommand(serverUpdateCmd)
	serverUpdateCmd.Flags().StringVar(&serverFlags.configSource, "config-source", "", "Source of the server update configuration. Can be 'pipe' or path to a JSON file.")
	serverUpdateCmd.MarkFlagsOneRequired("config-source")

	serverCmd.AddCommand(serverUpdateIpmiCredentialsCmd)

	serverCmd.AddCommand(serverEnableSnmpCmd)

	serverCmd.AddCommand(serverEnableSyslogCmd)

	serverCmd.AddCommand(serverVncInfoCmd)

	serverCmd.AddCommand(serverRemoteConsoleInfoCmd)

	serverCmd.AddCommand(serverCapabilitiesCmd)

	// Firmware commands
	serverCmd.AddCommand(serverFirmwareCmd)

	// Firmware component listing
	serverFirmwareCmd.AddCommand(serverFirmwareComponentsListCmd)

	// Get firmware component info
	serverFirmwareCmd.AddCommand(serverFirmwareComponentGetCmd)

	// Update firmware component
	serverFirmwareCmd.AddCommand(serverFirmwareComponentUpdateCmd)
	serverFirmwareComponentUpdateCmd.Flags().StringVar(&serverFlags.configSource, "config-source", "", "Source of the component update configuration. Can be 'pipe' or path to a JSON file.")
	serverFirmwareComponentUpdateCmd.MarkFlagsOneRequired("config-source")

	// Update firmware info
	serverFirmwareCmd.AddCommand(serverFirmwareUpdateInfoCmd)

	// Get firmware inventory
	serverFirmwareCmd.AddCommand(serverFirmwareInventoryCmd)

	// Upgrade firmware
	serverFirmwareCmd.AddCommand(serverFirmwareUpgradeCmd)

	// Upgrade component firmware
	serverFirmwareCmd.AddCommand(serverFirmwareComponentUpgradeCmd)
	serverFirmwareComponentUpgradeCmd.Flags().StringVar(&serverFlags.configSource, "config-source", "", "Source of the firmware upgrade configuration. Can be 'pipe' or path to a JSON file.")
	serverFirmwareComponentUpgradeCmd.MarkFlagsOneRequired("config-source")

	// Schedule firmware upgrade
	serverFirmwareCmd.AddCommand(serverFirmwareScheduleUpgradeCmd)
	serverFirmwareScheduleUpgradeCmd.Flags().StringVar(&serverFlags.configSource, "config-source", "", "Source of the schedule configuration. Can be 'pipe' or path to a JSON file.")
	serverFirmwareScheduleUpgradeCmd.MarkFlagsOneRequired("config-source")

	// Fetch firmware versions
	serverFirmwareCmd.AddCommand(serverFirmwareFetchVersionsCmd)

	// Generate firmware audit
	serverFirmwareCmd.AddCommand(serverFirmwareGenerateAuditCmd)
	serverFirmwareGenerateAuditCmd.Flags().StringVar(&serverFlags.configSource, "config-source", "", "Source of the audit configuration. Can be 'pipe' or path to a JSON file.")
	serverFirmwareGenerateAuditCmd.MarkFlagsOneRequired("config-source")
}
