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
		showCredentials bool
		filterStatus    []string
		filterType      []string
		configSource    string

		// list
		listFilterId                         []string
		listFilterSiteId                     []string
		listFilterVendor                     []string
		listFilterModel                      []string
		listFilterSerialNumber               []string
		listFilterManagementAddress          []string
		listFilterAdministrationState        []string
		listFilterPowerStatus                []string
		listFilterDhcpStatus                 []string
		listFilterClass                      []string
		listFilterDiskWipe                   []string
		listFilterRequiresManualCleaning     []string
		listFilterDriftDetectionSyncStatus   []string
		listFilterFirmwareBaselineCompliance []string
		listFilterHealthState                []string
		listFilterFirmwareBackupEnabled      []string
		listFilterFirmwareDriftDetection     []string
		listFilterTags                       []string
		listFilterInfrastructureId           []string
		listFilterOwnerId                    []string
		listPage                             int
		listLimit                            int

		siteId                int
		managementAddress     string
		username              string
		password              string
		serialNumber          string
		model                 string
		vendor                string
		serverStatus          string
		serverUuid            string
		bmcMacAddress         string
		registrationProfileId int
		rackName              string
		rackPositionLower     string
		rackPositionUpper     string
		inventoryId           string

		// replace
		replacementServerId string
		acceptDowngrade     bool

		// snapshots
		snapshotKind string

		// hardware-rescan
		rebootAllowed bool

		// register-production
		registerProductionConfigSource   string
		registerProductionSiteId         int
		registerProductionInfrastructure int
		registerProductionOsTemplateId   int
		registerProductionAddress        string
		registerProductionUsername       string
		registerProductionPassword       string
		registerProductionSerialNumber   string
		registerProductionUuid           string
		registerProductionBmcMacAddress  string
		registerProductionVendor         string
		registerProductionModel          string
		registerProductionProfileId      int

		// connect-interface
		connectInterfaceConfigSource string
		connectInterfaceId           int
		connectInterfacePortId       string
		connectInterfaceHostname     string

		// set-interfaces-default-fabric
		defaultFabricConfigSource string
		defaultFabricInterfaceIds []string
		defaultFabricId           int
		defaultFabricClear        bool

		// set-interfaces-redundancy-group
		redundancyGroupConfigSource string
		redundancyGroupInterfaceIds []string
		redundancyGroupIndex        int
		redundancyGroupClear        bool

		// import-unmanaged
		importUnmanagedConfigSource string

		// firmware schedule-upgrade-batch
		firmwareBatchScheduleTimestamp    string
		firmwareBatchConfirmationRequired bool

		// firmware components
		firmwareComponentClass     []string
		firmwareComponentBaseboard []string

		// firmware generate-audit
		firmwareIncludePreflight bool

		// firmware search-dry-run
		firmwareDryRunBaselineId  int64
		firmwareDryRunCatalogIds  []int64
		firmwareDryRunAllVersions bool
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
  - Planned servers: register --status planned, register-planned, set-registration-ready
  - Power management: power (on, off, reset, cycle, soft, status)
  - Maintenance: re-register, factory-reset, archive
  - Replacement: replacement-candidates, replace, cleanup
  - Security: update-ipmi-credentials, enable-snmp, enable-syslog
  - Remote access: vnc-info, console-info
  - Firmware: firmware subcommands for component management and upgrades
  - Information: capabilities, statistics
  - Drift detection: drift subcommands, snapshots, sync-target-snapshot
  - Hardware: hardware-rescan, connect-interface, set-interfaces-default-fabric,
    set-interfaces-redundancy-group
  - Onboarding: register-production, import-unmanaged, config-example

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

The table output also shows the operational fields: health state, details and last check, firmware
baseline compliance and target/current baselines, drift detection status, backup, drift detection and
health monitoring settings, and tags. These are independent of the lifecycle status.

Filters are applied server-side. Every --filter-* flag accepts a comma-separated list or can be
repeated; multiple values are combined with OR. A single value is sent as-is, so it may also use the
API operator syntax (e.g. '$not:$eq:value' or '$ilike:value').
Without --page/--limit all matching servers are returned.

Optional Flags:
  --show-credentials                      Display server IPMI credentials (username and password)
  --filter-status                         Filter by server status (e.g., active, registered, provisioning)
  --filter-type                           Filter by server type ID
  --filter-id                             Filter by server ID
  --filter-site-id                        Filter by site ID
  --filter-tags                           Filter by tag (e.g., a rack or room label)
  --filter-vendor                         Filter by vendor
  --filter-model                          Filter by model
  --filter-serial-number                  Filter by serial number
  --filter-management-address             Filter by BMC management address
  --filter-administration-state           Filter by administration state
  --filter-power-status                   Filter by power status
  --filter-dhcp-status                    Filter by DHCP status
  --filter-class                          Filter by server class
  --filter-disk-wipe                      Filter by disk wipe setting
  --filter-requires-manual-cleaning       Filter by manual cleaning requirement
  --filter-drift-detection-sync-status    Filter by drift detection sync status
  --filter-firmware-baseline-compliance   Filter by firmware baseline compliance status (unknown, compliant, non_compliant)
  --filter-health-state                   Filter by health state (healthy, warnings_present, critical_warnings_present, failed)
  --filter-firmware-backup-enabled        Filter by firmware backup setting
  --filter-firmware-drift-detection       Filter by firmware drift detection setting
  --filter-infrastructure-id              Filter by the infrastructure the server is allocated to
  --filter-owner-id                       Filter by the owner (user ID) of the infrastructure
  --page                                  Page number to retrieve (default: return all records)
  --limit                                 Number of records per page, or maximum records when used without --page

Examples:
  # List all servers
  metalcloud-cli server list

  # List servers with IPMI credentials
  metalcloud-cli server list --show-credentials

  # Filter servers by status
  metalcloud-cli server list --filter-status active,registered

  # Every server in one site carrying a given tag
  metalcloud-cli server list --filter-site-id 2 --filter-tags rack-12

  # Every server that is not compliant with its firmware baseline
  metalcloud-cli server list --filter-firmware-baseline-compliance non_compliant

  # Every server whose health is degraded or failed
  metalcloud-cli server list --filter-health-state warnings_present,critical_warnings_present,failed

  # Second page of 50 servers
  metalcloud-cli server list --page 2 --limit 50

  # Combine filters
  metalcloud-cli server list --filter-status active --filter-type 1
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerList(cmd.Context(), server.ListFlags{
				ShowCredentials:                  serverFlags.showCredentials,
				FilterStatus:                     serverFlags.filterStatus,
				FilterType:                       serverFlags.filterType,
				FilterId:                         serverFlags.listFilterId,
				FilterSiteId:                     serverFlags.listFilterSiteId,
				FilterVendor:                     serverFlags.listFilterVendor,
				FilterModel:                      serverFlags.listFilterModel,
				FilterSerialNumber:               serverFlags.listFilterSerialNumber,
				FilterManagementAddress:          serverFlags.listFilterManagementAddress,
				FilterAdministrationState:        serverFlags.listFilterAdministrationState,
				FilterPowerStatus:                serverFlags.listFilterPowerStatus,
				FilterDhcpStatus:                 serverFlags.listFilterDhcpStatus,
				FilterClass:                      serverFlags.listFilterClass,
				FilterDiskWipe:                   serverFlags.listFilterDiskWipe,
				FilterRequiresManualCleaning:     serverFlags.listFilterRequiresManualCleaning,
				FilterDriftDetectionSyncStatus:   serverFlags.listFilterDriftDetectionSyncStatus,
				FilterFirmwareBaselineCompliance: serverFlags.listFilterFirmwareBaselineCompliance,
				FilterHealthState:                serverFlags.listFilterHealthState,
				FilterFirmwareBackupEnabled:      serverFlags.listFilterFirmwareBackupEnabled,
				FilterFirmwareDriftDetection:     serverFlags.listFilterFirmwareDriftDetection,
				FilterTags:                       serverFlags.listFilterTags,
				FilterInfrastructureId:           serverFlags.listFilterInfrastructureId,
				FilterOwnerId:                    serverFlags.listFilterOwnerId,
				Page:                             serverFlags.listPage,
				Limit:                            serverFlags.listLimit,
			})
		},
	}

	serverGetCmd = &cobra.Command{
		Use:     "get server_id",
		Aliases: []string{"show"},
		Short:   "Get detailed server information",
		Long: `Get detailed information for a specific server.

This command retrieves comprehensive information about a server including its configuration,
status, hardware details, and optionally IPMI credentials.

The table output also shows the operational fields: health state, details and last check, firmware
baseline compliance and target/current baselines, drift detection status, backup, drift detection and
health monitoring settings, and tags. These are independent of the lifecycle status.

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
		Aliases: []string{"new", "create"},
		Short:   "Create a server record and optionally register it",
		Long: `Create a server record in MetalSoft and optionally start its registration.

The --status flag controls what happens after the record is created:
  registering            (default) Start the registration right away. The BMC must be
                         reachable, so --management-address and the BMC credentials
                         are normally needed.
  planned                Only record the server, for example ahead of its delivery.
                         Only --site-id and --serial-number are required; all other
                         details (BMC address, credentials, rack, inventory) are optional
                         and can be filled in later with "server update". Start the
                         registration later with "server register-planned".
  registration_ready     Record the server as ready to be registered, for example
                         when it is to be picked up by ZTP/DHCP discovery.

You can provide the server configuration either via command-line flags or with
--config-source (a JSON/YAML file path or 'pipe'). Run
"metalcloud-cli server config-example register" for an example.

Required Flags (when not using --config-source):
  --site-id                Site ID where the server is located
  --serial-number          Server serial number

Optional Flags:
  --config-source          Source of server configuration (file path or 'pipe')
  --status                 Initial status: registering (default), planned, registration_ready
  --management-address     BMC management IP address
  --username               BMC username
  --password               BMC password
  --model                  Server model
  --vendor                 Server vendor
  --uuid                   Server UUID
  --bmc-mac                BMC MAC address
  --registration-profile-id  Registration profile to use
  --rack-name              Rack name
  --rack-position-lower    Rack position lower unit
  --rack-position-upper    Rack position upper unit
  --inventory-id           Inventory ID

The rack and inventory details are not accepted when the server is created, so they
are applied with an update of the new server right after it is created.

Examples:
  # Create and register a server right away
  metalcloud-cli server register --site-id 1 --serial-number ABC123 --management-address 10.0.0.1 --username admin --password secret

  # Record a planned server with only its serial number
  metalcloud-cli server register --site-id 1 --serial-number ABC123 --status planned

  # Record a planned server with its BMC defaults and rack position
  metalcloud-cli server register --site-id 1 --serial-number ABC123 --status planned \
    --username root --password calvin --rack-name R01 --rack-position-lower 10 --rack-position-upper 11

  # Register using a configuration file
  metalcloud-cli server register --config-source ./server.json

  # Register using piped configuration
  cat server.json | metalcloud-cli server register --config-source pipe
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			var serverConfig sdk.CreateServer

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
				serverConfig = buildCreateServerFromFlags(cmd)
			}

			return server.ServerRegister(cmd.Context(), serverConfig, buildRackUpdateFromFlags())
		},
	}

	serverRegisterPlannedCmd = &cobra.Command{
		Use:   "register-planned server_id",
		Short: "Start the registration of a planned server",
		Long: `Start the registration of a server in the "planned" state.

The server must have a reachable BMC address and valid BMC credentials. Set them
first with "server update" if they were not provided when the server was created.

Required Arguments:
  server_id              The ID of the planned server

Examples:
  # Start the registration of planned server 123
  metalcloud-cli server register-planned 123
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerRegisterPlanned(cmd.Context(), args[0])
		},
	}

	serverSetRegistrationReadyCmd = &cobra.Command{
		Use:   "set-registration-ready server_id",
		Short: "Set a server to the registration_ready status",
		Long: `Set the status of a server to "registration_ready".

Required Arguments:
  server_id              The ID of the server

Examples:
  # Mark planned server 123 as ready to be registered
  metalcloud-cli server set-registration-ready 123
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerSetRegistrationReady(cmd.Context(), args[0])
		},
	}

	serverReplacementCandidatesCmd = &cobra.Command{
		Use:   "replacement-candidates server_id",
		Short: "List the servers that can replace a server in its instance",
		Long: `List the available servers that can replace the given server in its instance.

Eligible candidates are listed first, best first. Ineligible candidates are listed
with the reasons they were skipped. Use json or yaml output to also see the
requirements (RAM, cores, disks, GPUs, network links) the candidates were checked against.

Required Arguments:
  server_id              The ID of the server to replace

Examples:
  # List replacement candidates for server 123
  metalcloud-cli server replacement-candidates 123

  # Include the replacement requirements
  metalcloud-cli server replacement-candidates 123 -f json
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerReplacementCandidates(cmd.Context(), args[0])
		},
	}

	serverReplaceCmd = &cobra.Command{
		Use:   "replace server_id",
		Short: "Replace a server in its instance with another server",
		Long: `Replace a server in its instance with another server.

The replacement server is put in the instance and the infrastructure is deployed,
which moves the instance onto the new server with the same networking and OS. The
old server becomes defective and should be wiped with "server cleanup".

Use "server replacement-candidates" to find an eligible replacement server.

Required Arguments:
  server_id                  The ID of the server to replace

Required Flags:
  --replacement-server-id    The ID of the server to put in its place

Optional Flags:
  --accept-downgrade         Accept a replacement with less RAM, cores or disk

Examples:
  # Replace server 123 with server 456
  metalcloud-cli server replace 123 --replacement-server-id 456

  # Replace with a smaller server
  metalcloud-cli server replace 123 --replacement-server-id 456 --accept-downgrade
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerReplace(cmd.Context(), args[0], serverFlags.replacementServerId, serverFlags.acceptDowngrade)
		},
	}

	serverCleanupCmd = &cobra.Command{
		Use:   "cleanup server_id",
		Short: "Wipe a server that was replaced in its instance",
		Long: `Wipe a server that was replaced in its instance.

The server stays defective after the cleanup; set it available by hand once it is fixed.

Required Arguments:
  server_id              The ID of the replaced server

Examples:
  # Clean up replaced server 123
  metalcloud-cli server cleanup 123
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerCleanupReplaced(cmd.Context(), args[0])
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
		Use:   "power <server_id> <on|off|reset|cycle|soft|status>",
		Short: "Control server power state",
		Long: `Control server power state.

This command allows you to control the power state of a server by sending
power management commands to the server's BMC/IPMI interface.

Arguments:
  server_id              The ID of the server to control
  action                 Power action to perform (on, off, reset, cycle, soft, status)

Valid Actions:
  on                    Power on the server
  off                   Hard power off the server
  reset                 Hard reset the server
  cycle                 Power cycle the server (off then on)
  soft                  Soft power off the server (graceful shutdown)
  status                Get the current power status of the server

Subcommands:
  status                Get the current power status of a server

Examples:
  # Power on server
  metalcloud-cli server power 123 on

  # Hard power off server
  metalcloud-cli server power 123 off

  # Reset server
  metalcloud-cli server power 123 reset

  # Power cycle server
  metalcloud-cli server power 123 cycle

  # Soft power off (graceful shutdown)
  metalcloud-cli server power 123 soft

  # Get power status
  metalcloud-cli server power 123 status
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if args[1] == "status" {
				return server.ServerPowerStatus(cmd.Context(), args[0])
			}
			return server.ServerPower(cmd.Context(), args[0], args[1])
		},
	}

	serverPowerStatusCmd = &cobra.Command{
		Use:   "status server_id",
		Short: "Get server power status",
		Long: `Get the current power status of a server.

This command retrieves the current power state of the specified server
from its BMC/IPMI interface.

Required Arguments:
  server_id              The ID of the server to query

Examples:
  # Get power status for server with ID 123
  metalcloud-cli server power status 123
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

	serverIdentifyCmd = &cobra.Command{
		Use:          "identify server_id",
		Short:        "Identify a server by blinking its chassis LED",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerIdentify(cmd.Context(), args[0])
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
  - Auditing: generate-audit, search-dry-run

Use "metalcloud-cli server firmware [command] --help" for detailed information about each command.
`,
	}

	serverFirmwareComponentsListCmd = &cobra.Command{
		Use:     "components server_id",
		Aliases: []string{"list-components", "ls"},
		Short:   "List firmware components for a server",
		Long: `List firmware components for a server.

This command displays all firmware components available on the specified server,
including their IDs, names, class, vendor, model, versions and, for the components of a
GPU baseboard (an HGX tray), the baseboard they belong to.

Required Arguments:
  server_id              The ID of the server to query

Optional Flags:
  --class                Only the components of these classes (bios, bmc, gpu, nic, hgx_hmc, ...)
  --baseboard            Only the components of these GPU baseboards

Examples:
  # List firmware components for server with ID 123
  metalcloud-cli server firmware components 123

  # Only the GPUs
  metalcloud-cli server firmware components 123 --class gpu

  # Using short alias
  metalcloud-cli server fw ls 123
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerFirmwareComponentsListFiltered(cmd.Context(), args[0], server.ComponentsListFlags{
				FilterComponentClass: serverFlags.firmwareComponentClass,
				FilterBaseboardId:    serverFlags.firmwareComponentBaseboard,
			})
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

When the target version comes from a firmware package that other components share (for example
an NVIDIA bundle step that updates every GPU of an HGX tray in one job), the same target is set
on every component the package updates, and setting it to null resets all of them. The text
output lists the other components whose target changed.

Examples:
  # Set the target version of a component
  echo '{"firmwareTargetVersion": "96.00.A5.00.01"}' | metalcloud-cli server firmware update-component 123 456 --config-source pipe

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

Before the jobs are created the safety checks run (server state, other firmware jobs,
BMC reachability, prerequisites, component health). Components the checks refuse are
left out and listed with their reason; when nothing is left the upgrade is refused and
each reason is printed. Use 'generate-audit --include-preflight' to see the outcome first.

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
  echo '{"serverIds": [123, 456]}' | metalcloud-cli server firmware generate-audit --config-source pipe

  # Also run the safety checks an upgrade would run now and list the jobs it would create
  echo '{"serverIds": [123]}' | metalcloud-cli server firmware generate-audit --config-source pipe --include-preflight

Optional Flags:
  --include-preflight    Run the upgrade safety checks (BMC checks included) and list the planned
                         jobs, in order, for each server. Nothing is changed.
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(serverFlags.configSource)
			if err != nil {
				return err
			}
			return server.ServerFirmwareGenerateAuditWithPreflight(cmd.Context(), config, serverFlags.firmwareIncludePreflight)
		},
	}

	serverFirmwareSearchDryRunCmd = &cobra.Command{
		Use:     "search-dry-run server_id",
		Aliases: []string{"explain-matching"},
		Short:   "Explain which firmware binaries match each component of a server",
		Long: `Explain which firmware binaries match each component of a server.

Runs the binary search a firmware upgrade uses and shows, for each component, the binaries
that matched (with the rules that matched them) and the ones rejected (with the reason).
Nothing is stored or scheduled.

Required Arguments:
  server_id              The ID of the server

Optional Flags:
  --baseline-id          Search this firmware baseline
  --catalog-ids          Search these firmware catalogs (instead of a baseline)
  --all-versions         Include every matching version, not only the latest

Without --baseline-id or --catalog-ids the server's target baseline is searched.

Examples:
  # Search the server's target baseline
  metalcloud-cli server firmware search-dry-run 123

  # Search two catalogs, every version
  metalcloud-cli server firmware search-dry-run 123 --catalog-ids 10,11 --all-versions
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerFirmwareSearchDryRun(cmd.Context(), args[0], server.ServerFirmwareSearchDryRunOptions{
				BaselineId:         serverFlags.firmwareDryRunBaselineId,
				CatalogIds:         serverFlags.firmwareDryRunCatalogIds,
				IncludeAllVersions: serverFlags.firmwareDryRunAllVersions,
			})
		},
	}

	// --- drift & snapshots --------------------------------------------------

	serverDriftCmd = &cobra.Command{
		Use:   "drift [command]",
		Short: "Inspect the configuration drift of a server",
		Long: `Inspect and acknowledge the configuration drift detected between the running
configuration of a server and its target snapshot.

Use "metalcloud-cli server drift [command] --help" for detailed information about each command.
`,
	}

	serverDriftListCmd = &cobra.Command{
		Use:     "list server_id",
		Aliases: []string{"ls"},
		Short:   "List the configuration drift of a server",
		Long: `List the configuration drift entries recorded for a server.

Required Arguments:
  server_id              The ID of the server

Examples:
  # List the drift history of server 123
  metalcloud-cli server drift list 123
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerDriftList(cmd.Context(), args[0])
		},
	}

	serverDriftGetCmd = &cobra.Command{
		Use:     "get server_id drift_id",
		Aliases: []string{"show"},
		Short:   "Get one configuration drift entry of a server",
		Long: `Display one configuration drift entry of a server, including the configuration
difference that was detected.

Required Arguments:
  server_id              The ID of the server
  drift_id               The ID of the drift entry

Examples:
  # Show drift entry 9 of server 123
  metalcloud-cli server drift get 123 9
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerDriftGet(cmd.Context(), args[0], args[1])
		},
	}

	serverDriftAcknowledgeCmd = &cobra.Command{
		Use:     "acknowledge server_id drift_id",
		Aliases: []string{"ack"},
		Short:   "Acknowledge a configuration drift entry of a server",
		Long: `Mark one configuration drift entry of a server as reviewed, recording who
acknowledged it and when.

Required Arguments:
  server_id              The ID of the server
  drift_id               The ID of the drift entry

Examples:
  # Acknowledge drift entry 9 of server 123
  metalcloud-cli server drift acknowledge 123 9
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerDriftAcknowledge(cmd.Context(), args[0], args[1])
		},
	}

	serverSnapshotsCmd = &cobra.Command{
		Use:     "snapshots server_id",
		Aliases: []string{"snapshot"},
		Short:   "List the configuration snapshots of a server",
		Long: `List the configuration snapshots stored for a server.

Required Arguments:
  server_id              The ID of the server

Optional Flags:
  --kind                 Restrict the listing to one snapshot class

Examples:
  # List all snapshots of server 123
  metalcloud-cli server snapshots 123

  # List only the backup snapshots of server 123
  metalcloud-cli server snapshots 123 --kind backup
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerSnapshotList(cmd.Context(), args[0], serverFlags.snapshotKind)
		},
	}

	serverSyncTargetSnapshotCmd = &cobra.Command{
		Use:   "sync-target-snapshot server_id",
		Short: "Accept the current configuration of a server as the drift target",
		Long: `Point the drift detection target snapshot of a server at its latest snapshot,
clearing the drift currently reported for the server.

Required Arguments:
  server_id              The ID of the server

Examples:
  # Accept the current configuration of server 123 as the new target
  metalcloud-cli server sync-target-snapshot 123
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerSyncTargetSnapshot(cmd.Context(), args[0])
		},
	}

	// --- hardware -----------------------------------------------------------

	serverHardwareRescanCmd = &cobra.Command{
		Use:   "hardware-rescan server_id",
		Short: "Re-read the hardware inventory of a server",
		Long: `Re-read the hardware inventory of a server.

The rescan updates the recorded hardware configuration of the server. By default
the server is not rebooted, which preserves the existing switch connections; pass
--reboot-allowed to let MetalSoft reboot the server and run a full LLDP interface
discovery.

Required Arguments:
  server_id              The ID of the server to rescan

Optional Flags:
  --reboot-allowed       Allow the server to be rebooted for full LLDP interface discovery

Examples:
  # Rescan the hardware of server 123 without rebooting it
  metalcloud-cli server hardware-rescan 123

  # Rescan the hardware of server 123, allowing a reboot
  metalcloud-cli server hardware-rescan 123 --reboot-allowed
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerHardwareRescan(cmd.Context(), args[0], serverFlags.rebootAllowed)
		},
	}

	serverStatisticsCmd = &cobra.Command{
		Use:     "statistics",
		Aliases: []string{"stats"},
		Short:   "Get aggregated server statistics",
		Long: `Get the aggregated server counts, grouped by server status and by site.

Examples:
  # Show the server statistics
  metalcloud-cli server statistics

  # Show the server statistics as JSON
  metalcloud-cli server statistics -f json
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerStatistics(cmd.Context())
		},
	}

	serverConnectInterfaceCmd = &cobra.Command{
		Use:   "connect-interface server_id",
		Short: "Record the network device port a server interface is cabled to",
		Long: `Record the network device port a server interface is cabled to.

You can provide the connection either via command-line flags or by specifying a
configuration source using the --config-source flag. The configuration source can
be a path to a JSON file or 'pipe' to read from standard input.

Required Arguments:
  server_id              The ID of the server

Required Flags (when not using --config-source):
  --interface-id         The ID of the server interface to connect
  --port-id              The network device port name (e.g. Ethernet0)
  --hostname             The network device hostname

Optional Flags:
  --config-source        Source of the connection configuration. Can be 'pipe' or path to a JSON file.

Flag Dependencies:
  --config-source and --interface-id are mutually exclusive
  --interface-id, --port-id and --hostname must be used together

Examples:
  # Connect interface 1 of server 123 to port Ethernet0 of leaf-01
  metalcloud-cli server connect-interface 123 --interface-id 1 --port-id Ethernet0 --hostname leaf-01

  # Connect using a JSON configuration file
  metalcloud-cli server connect-interface 123 --config-source ./connect.json

  # Show the expected configuration
  metalcloud-cli server config-example connect-interface
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var connectConfig sdk.ServerConnectInterface

			if serverFlags.connectInterfaceConfigSource != "" {
				config, err := utils.ReadConfigFromPipeOrFile(serverFlags.connectInterfaceConfigSource)
				if err != nil {
					return err
				}
				if err := utils.UnmarshalContent(config, &connectConfig); err != nil {
					return err
				}
			} else {
				connectConfig = sdk.ServerConnectInterface{
					ServerInterfaceId:     int64(serverFlags.connectInterfaceId),
					NetworkDevicePortId:   serverFlags.connectInterfacePortId,
					NetworkDeviceHostname: serverFlags.connectInterfaceHostname,
				}
			}

			return server.ServerConnectInterface(cmd.Context(), args[0], connectConfig)
		},
	}

	serverSetInterfacesDefaultFabricCmd = &cobra.Command{
		Use:   "set-interfaces-default-fabric server_id",
		Short: "Set the default fabric of some server interfaces",
		Long: `Assign the default fabric of the given server interfaces, or clear it.

You can provide the configuration either via command-line flags or by specifying a
configuration source using the --config-source flag. The configuration source can
be a path to a JSON file or 'pipe' to read from standard input.

Required Arguments:
  server_id              The ID of the server

Required Flags (when not using --config-source):
  --interface-id         ID of a server interface; repeat or comma-separate for several
  --fabric-id            The ID of the fabric to set as default (or use --clear-fabric)

Optional Flags:
  --config-source        Source of the configuration. Can be 'pipe' or path to a JSON file.
  --clear-fabric         Clear the default fabric of the given interfaces

Flag Dependencies:
  --config-source and --interface-id are mutually exclusive
  --fabric-id and --clear-fabric are mutually exclusive

Examples:
  # Set fabric 10 as the default fabric of interfaces 1 and 2 of server 123
  metalcloud-cli server set-interfaces-default-fabric 123 --interface-id 1,2 --fabric-id 10

  # Clear the default fabric of interface 1 of server 123
  metalcloud-cli server set-interfaces-default-fabric 123 --interface-id 1 --clear-fabric

  # Use a JSON configuration file
  metalcloud-cli server set-interfaces-default-fabric 123 --config-source ./fabric.json
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var fabricConfig sdk.ServerInterfacesDefaultFabric

			if serverFlags.defaultFabricConfigSource != "" {
				config, err := utils.ReadConfigFromPipeOrFile(serverFlags.defaultFabricConfigSource)
				if err != nil {
					return err
				}
				if err := utils.UnmarshalContent(config, &fabricConfig); err != nil {
					return err
				}
			} else {
				interfaceIds, err := utils.GetInt64SliceFromStrings(serverFlags.defaultFabricInterfaceIds)
				if err != nil {
					return err
				}

				fabricConfig = sdk.ServerInterfacesDefaultFabric{
					ServerInterfaceIds: interfaceIds,
				}

				if serverFlags.defaultFabricClear {
					fabricConfig.DefaultFabricId.Set(nil)
				} else {
					fabricConfig.DefaultFabricId.Set(sdk.PtrInt64(int64(serverFlags.defaultFabricId)))
				}
			}

			return server.ServerSetInterfacesDefaultFabric(cmd.Context(), args[0], fabricConfig)
		},
	}

	serverSetInterfacesRedundancyGroupCmd = &cobra.Command{
		Use:   "set-interfaces-redundancy-group server_id",
		Short: "Set the redundancy group of some server interfaces",
		Long: `Group the given server interfaces into a redundancy group, or remove them from one.

You can provide the configuration either via command-line flags or by specifying a
configuration source using the --config-source flag. The configuration source can
be a path to a JSON file or 'pipe' to read from standard input.

Required Arguments:
  server_id              The ID of the server

Required Flags (when not using --config-source):
  --interface-id         ID of a server interface; repeat or comma-separate for several
  --group-index          The redundancy group index (or use --clear-group)

Optional Flags:
  --config-source        Source of the configuration. Can be 'pipe' or path to a JSON file.
  --clear-group          Remove the given interfaces from their redundancy group

Flag Dependencies:
  --config-source and --interface-id are mutually exclusive
  --group-index and --clear-group are mutually exclusive

Examples:
  # Put interfaces 1 and 2 of server 123 in redundancy group 1
  metalcloud-cli server set-interfaces-redundancy-group 123 --interface-id 1,2 --group-index 1

  # Remove interface 1 of server 123 from its redundancy group
  metalcloud-cli server set-interfaces-redundancy-group 123 --interface-id 1 --clear-group

  # Use a JSON configuration file
  metalcloud-cli server set-interfaces-redundancy-group 123 --config-source ./redundancy.json
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var redundancyConfig sdk.ServerInterfacesRedundancyGroup

			if serverFlags.redundancyGroupConfigSource != "" {
				config, err := utils.ReadConfigFromPipeOrFile(serverFlags.redundancyGroupConfigSource)
				if err != nil {
					return err
				}
				if err := utils.UnmarshalContent(config, &redundancyConfig); err != nil {
					return err
				}
			} else {
				interfaceIds, err := utils.GetInt64SliceFromStrings(serverFlags.redundancyGroupInterfaceIds)
				if err != nil {
					return err
				}

				redundancyConfig = sdk.ServerInterfacesRedundancyGroup{
					ServerInterfaceIds: interfaceIds,
				}

				if serverFlags.redundancyGroupClear {
					redundancyConfig.RedundancyGroupIndex.Set(nil)
				} else {
					redundancyConfig.RedundancyGroupIndex.Set(sdk.PtrFloat32(float32(serverFlags.redundancyGroupIndex)))
				}
			}

			return server.ServerSetInterfacesRedundancyGroup(cmd.Context(), args[0], redundancyConfig)
		},
	}

	// --- onboarding ---------------------------------------------------------

	serverRegisterProductionCmd = &cobra.Command{
		Use:   "register-production",
		Short: "Register a server that is already running a production workload",
		Long: `Register a server that is already running a production workload, keeping the
workload in place.

The server is attached to an existing infrastructure instead of being wiped and
re-provisioned. You can provide the configuration either via command-line flags
or by specifying a configuration source using the --config-source flag. The
configuration source can be a path to a JSON file or 'pipe' to read from
standard input.

Required Flags (when not using --config-source):
  --site-id              Site ID where the server is located
  --infrastructure-id    ID of the infrastructure the server belongs to

Optional Flags:
  --config-source        Source of the configuration. Can be 'pipe' or path to a JSON file.
  --os-template-id       ID of the OS template already installed on the server
  --management-address   IPMI/BMC management IP address
  --username             IPMI/BMC username
  --password             IPMI/BMC password
  --serial-number        Server serial number
  --uuid                 Server UUID
  --bmc-mac-address      BMC MAC address
  --vendor               Server vendor
  --model                Server model
  --registration-profile-id  ID of the server registration profile to use

Flag Dependencies:
  --config-source and --site-id are mutually exclusive
  --site-id and --infrastructure-id must be used together

Examples:
  # Register a production server with its interface connections from a JSON file
  metalcloud-cli server register-production --config-source ./production-server.json

  # Register a production server using flags
  metalcloud-cli server register-production --site-id 1 --infrastructure-id 100 \
    --management-address 10.0.0.1 --username admin --password secret

  # Show the expected configuration
  metalcloud-cli server config-example register-production
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			var serverConfig sdk.RegisterProductionServer

			if serverFlags.registerProductionConfigSource != "" {
				config, err := utils.ReadConfigFromPipeOrFile(serverFlags.registerProductionConfigSource)
				if err != nil {
					return err
				}
				if err := utils.UnmarshalContent(config, &serverConfig); err != nil {
					return err
				}
			} else {
				serverConfig = sdk.RegisterProductionServer{
					SiteId: int64(serverFlags.registerProductionSiteId),
					Settings: sdk.RegisterProductionServerSettings{
						InfrastructureId: int64(serverFlags.registerProductionInfrastructure),
					},
				}

				if serverFlags.registerProductionOsTemplateId > 0 {
					serverConfig.Settings.OsTemplateId = sdk.PtrInt64(int64(serverFlags.registerProductionOsTemplateId))
				}
				if serverFlags.registerProductionAddress != "" {
					serverConfig.ManagementAddress = sdk.PtrString(serverFlags.registerProductionAddress)
				}
				if serverFlags.registerProductionUsername != "" {
					serverConfig.Username = sdk.PtrString(serverFlags.registerProductionUsername)
				}
				if serverFlags.registerProductionPassword != "" {
					serverConfig.Password = sdk.PtrString(serverFlags.registerProductionPassword)
				}
				if serverFlags.registerProductionSerialNumber != "" {
					serverConfig.SerialNumber = sdk.PtrString(serverFlags.registerProductionSerialNumber)
				}
				if serverFlags.registerProductionUuid != "" {
					serverConfig.ServerUUID = sdk.PtrString(serverFlags.registerProductionUuid)
				}
				if serverFlags.registerProductionBmcMacAddress != "" {
					serverConfig.BmcMacAddress = sdk.PtrString(serverFlags.registerProductionBmcMacAddress)
				}
				if serverFlags.registerProductionVendor != "" {
					serverConfig.Vendor = sdk.PtrString(serverFlags.registerProductionVendor)
				}
				if serverFlags.registerProductionModel != "" {
					serverConfig.Model = sdk.PtrString(serverFlags.registerProductionModel)
				}
				if serverFlags.registerProductionProfileId > 0 {
					serverConfig.RegistrationProfileId = sdk.PtrInt64(int64(serverFlags.registerProductionProfileId))
				}
			}

			return server.ServerRegisterProduction(cmd.Context(), serverConfig)
		},
	}

	serverImportUnmanagedCmd = &cobra.Command{
		Use:   "import-unmanaged",
		Short: "Import a server whose lifecycle MetalSoft does not manage",
		Long: `Import a server that MetalSoft does not manage the lifecycle of, using an
externally supplied hardware description.

The configuration must describe the site, the server type and the server
interfaces together with the network device ports they are cabled to.

Required Flags:
  --config-source        Source of the import configuration. Can be 'pipe' or path to a JSON file.

Examples:
  # Import an unmanaged server from a JSON file
  metalcloud-cli server import-unmanaged --config-source ./unmanaged-server.json

  # Import an unmanaged server from piped configuration
  metalcloud-cli server config-example import-unmanaged | metalcloud-cli server import-unmanaged --config-source pipe
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(serverFlags.importUnmanagedConfigSource)
			if err != nil {
				return err
			}
			return server.ServerImportUnmanaged(cmd.Context(), config)
		},
	}

	serverConfigExampleCmd = &cobra.Command{
		Use:     "config-example kind",
		Aliases: []string{"example"},
		Short:   "Show a server configuration example",
		Long: `Show an example of the configuration accepted by the server sub-commands that
take a --config-source flag.

Required Arguments:
  kind                   The configuration kind. One of:
                           register
                           register-production
                           import-unmanaged
                           connect-interface
                           set-interfaces-default-fabric
                           set-interfaces-redundancy-group

Examples:
  # Show the production server registration configuration
  metalcloud-cli server config-example register-production

  # Save the unmanaged server import configuration for editing
  metalcloud-cli server config-example import-unmanaged > unmanaged-server.json
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerConfigExample(cmd.Context(), args[0])
		},
	}

	// --- firmware batch operations -----------------------------------------

	serverFirmwareUpgradeBatchCmd = &cobra.Command{
		Use:   "upgrade-batch server_id...",
		Short: "Upgrade the firmware of several servers at once",
		Long: `Upgrade the firmware of several servers at once.

A firmware upgrade job is started for each of the given servers. The result lists
the job started for every server that could be upgraded and the error reported for
every server that could not.

Required Arguments:
  server_id...           The IDs of the servers to upgrade

Examples:
  # Upgrade the firmware of servers 123, 124 and 125
  metalcloud-cli server firmware upgrade-batch 123 124 125
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerFirmwareUpgradeBatch(cmd.Context(), args)
		},
	}

	serverFirmwareScheduleUpgradeBatchCmd = &cobra.Command{
		Use:   "schedule-upgrade-batch server_id...",
		Short: "Schedule a firmware upgrade for several servers at once",
		Long: `Schedule a firmware upgrade for several servers at once.

The API models this as an action on a single server, so the first server ID given
is used to address the endpoint while the full list of servers is sent in the
request body. The result lists the error reported for every server that could not
be scheduled.

Required Arguments:
  server_id...           The IDs of the servers to schedule the upgrade for

Optional Flags:
  --schedule-timestamp       When the firmware upgrade should run (RFC 3339 timestamp)
  --confirmation-required    Require a confirmation before the scheduled upgrade runs

Examples:
  # Schedule an upgrade for servers 123 and 124 as soon as possible
  metalcloud-cli server firmware schedule-upgrade-batch 123 124

  # Schedule an upgrade for a maintenance window, with confirmation
  metalcloud-cli server firmware schedule-upgrade-batch 123 124 \
    --schedule-timestamp 2024-01-01T10:00:00Z --confirmation-required
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerFirmwareScheduleUpgradeBatch(cmd.Context(), args,
				serverFlags.firmwareBatchScheduleTimestamp,
				serverFlags.firmwareBatchConfirmationRequired)
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
	serverListCmd.Flags().StringSliceVar(&serverFlags.listFilterId, "filter-id", nil, "Filter the result by server ID.")
	serverListCmd.Flags().StringSliceVar(&serverFlags.listFilterSiteId, "filter-site-id", nil, "Filter the result by site ID.")
	serverListCmd.Flags().StringSliceVar(&serverFlags.listFilterTags, "filter-tags", nil, "Filter the result by server tag.")
	serverListCmd.Flags().StringSliceVar(&serverFlags.listFilterVendor, "filter-vendor", nil, "Filter the result by vendor.")
	serverListCmd.Flags().StringSliceVar(&serverFlags.listFilterModel, "filter-model", nil, "Filter the result by model.")
	serverListCmd.Flags().StringSliceVar(&serverFlags.listFilterSerialNumber, "filter-serial-number", nil, "Filter the result by serial number.")
	serverListCmd.Flags().StringSliceVar(&serverFlags.listFilterManagementAddress, "filter-management-address", nil, "Filter the result by management address.")
	serverListCmd.Flags().StringSliceVar(&serverFlags.listFilterAdministrationState, "filter-administration-state", nil, "Filter the result by administration state.")
	serverListCmd.Flags().StringSliceVar(&serverFlags.listFilterPowerStatus, "filter-power-status", nil, "Filter the result by power status.")
	serverListCmd.Flags().StringSliceVar(&serverFlags.listFilterDhcpStatus, "filter-dhcp-status", nil, "Filter the result by DHCP status.")
	serverListCmd.Flags().StringSliceVar(&serverFlags.listFilterClass, "filter-class", nil, "Filter the result by server class.")
	serverListCmd.Flags().StringSliceVar(&serverFlags.listFilterDiskWipe, "filter-disk-wipe", nil, "Filter the result by disk wipe setting.")
	serverListCmd.Flags().StringSliceVar(&serverFlags.listFilterRequiresManualCleaning, "filter-requires-manual-cleaning", nil, "Filter the result by manual cleaning requirement.")
	serverListCmd.Flags().StringSliceVar(&serverFlags.listFilterDriftDetectionSyncStatus, "filter-drift-detection-sync-status", nil, "Filter the result by drift detection sync status.")
	serverListCmd.Flags().StringSliceVar(&serverFlags.listFilterFirmwareBaselineCompliance, "filter-firmware-baseline-compliance", nil, "Filter the result by firmware baseline compliance status (unknown, compliant, non_compliant).")
	serverListCmd.Flags().StringSliceVar(&serverFlags.listFilterHealthState, "filter-health-state", nil, "Filter the result by health state (healthy, warnings_present, critical_warnings_present, failed).")
	serverListCmd.Flags().StringSliceVar(&serverFlags.listFilterFirmwareBackupEnabled, "filter-firmware-backup-enabled", nil, "Filter the result by firmware backup setting.")
	serverListCmd.Flags().StringSliceVar(&serverFlags.listFilterFirmwareDriftDetection, "filter-firmware-drift-detection", nil, "Filter the result by firmware drift detection setting.")
	serverListCmd.Flags().StringSliceVar(&serverFlags.listFilterInfrastructureId, "filter-infrastructure-id", nil, "Filter the result by the infrastructure the server is allocated to.")
	serverListCmd.Flags().StringSliceVar(&serverFlags.listFilterOwnerId, "filter-owner-id", nil, "Filter the result by the owner (user ID) of the infrastructure.")
	serverListCmd.Flags().IntVar(&serverFlags.listPage, "page", 0, "Page number to retrieve (default: return all records).")
	serverListCmd.Flags().IntVar(&serverFlags.listLimit, "limit", 0, "Number of records per page, or maximum records without --page (default: return all records).")

	serverCmd.AddCommand(serverGetCmd)
	serverGetCmd.Flags().BoolVar(&serverFlags.showCredentials, "show-credentials", false, "If set returns the server IPMI credentials.")

	serverCmd.AddCommand(serverRegisterCmd)
	serverRegisterCmd.Flags().StringVar(&serverFlags.configSource, "config-source", "", "Source of the new server configuration. Can be 'pipe' or path to a JSON/YAML file.")
	serverRegisterCmd.Flags().IntVar(&serverFlags.siteId, "site-id", 0, "Site ID")
	serverRegisterCmd.Flags().StringVar(&serverFlags.serialNumber, "serial-number", "", "Serial number")
	serverRegisterCmd.Flags().StringVar(&serverFlags.serverStatus, "status", "", "Initial server status: registering (default), planned or registration_ready")
	serverRegisterCmd.Flags().StringVar(&serverFlags.managementAddress, "management-address", "", "BMC management address")
	serverRegisterCmd.Flags().StringVar(&serverFlags.username, "username", "", "BMC username")
	serverRegisterCmd.Flags().StringVar(&serverFlags.password, "password", "", "BMC password")
	serverRegisterCmd.Flags().StringVar(&serverFlags.model, "model", "", "Model")
	serverRegisterCmd.Flags().StringVar(&serverFlags.vendor, "vendor", "", "Vendor")
	serverRegisterCmd.Flags().StringVar(&serverFlags.serverUuid, "uuid", "", "Server UUID")
	serverRegisterCmd.Flags().StringVar(&serverFlags.bmcMacAddress, "bmc-mac", "", "BMC MAC address")
	serverRegisterCmd.Flags().IntVar(&serverFlags.registrationProfileId, "registration-profile-id", 0, "Registration profile ID")
	serverRegisterCmd.Flags().StringVar(&serverFlags.rackName, "rack-name", "", "Rack name")
	serverRegisterCmd.Flags().StringVar(&serverFlags.rackPositionLower, "rack-position-lower", "", "Rack position lower unit")
	serverRegisterCmd.Flags().StringVar(&serverFlags.rackPositionUpper, "rack-position-upper", "", "Rack position upper unit")
	serverRegisterCmd.Flags().StringVar(&serverFlags.inventoryId, "inventory-id", "", "Inventory ID")
	serverRegisterCmd.MarkFlagsOneRequired("config-source", "site-id")
	serverRegisterCmd.MarkFlagsMutuallyExclusive("config-source", "site-id")
	serverRegisterCmd.MarkFlagsRequiredTogether("site-id", "serial-number")
	for _, flag := range []string{"status", "management-address", "username", "password", "model", "vendor", "uuid", "bmc-mac", "registration-profile-id"} {
		serverRegisterCmd.MarkFlagsMutuallyExclusive("config-source", flag)
	}

	serverCmd.AddCommand(serverRegisterPlannedCmd)

	serverCmd.AddCommand(serverSetRegistrationReadyCmd)

	serverCmd.AddCommand(serverReplacementCandidatesCmd)

	serverCmd.AddCommand(serverReplaceCmd)
	serverReplaceCmd.Flags().StringVar(&serverFlags.replacementServerId, "replacement-server-id", "", "The ID of the server to put in its place")
	serverReplaceCmd.Flags().BoolVar(&serverFlags.acceptDowngrade, "accept-downgrade", false, "Accept a replacement with less RAM, cores or disk than the instance had")
	serverReplaceCmd.MarkFlagRequired("replacement-server-id")

	serverCmd.AddCommand(serverCleanupCmd)

	serverCmd.AddCommand(serverReRegisterCmd)

	serverCmd.AddCommand(serverFactoryResetCmd)

	serverCmd.AddCommand(serverArchiveCmd)

	serverCmd.AddCommand(serverDeleteCmd)

	serverCmd.AddCommand(serverPowerCmd)
	serverPowerCmd.AddCommand(serverPowerStatusCmd)

	serverCmd.AddCommand(serverUpdateCmd)
	serverUpdateCmd.Flags().StringVar(&serverFlags.configSource, "config-source", "", "Source of the server update configuration. Can be 'pipe' or path to a JSON file.")
	serverUpdateCmd.MarkFlagsOneRequired("config-source")

	serverCmd.AddCommand(serverUpdateIpmiCredentialsCmd)

	serverCmd.AddCommand(serverEnableSnmpCmd)

	serverCmd.AddCommand(serverEnableSyslogCmd)

	serverCmd.AddCommand(serverVncInfoCmd)

	serverCmd.AddCommand(serverRemoteConsoleInfoCmd)

	serverCmd.AddCommand(serverCapabilitiesCmd)

	serverCmd.AddCommand(serverIdentifyCmd)

	// Firmware commands
	serverCmd.AddCommand(serverFirmwareCmd)

	// Firmware component listing
	serverFirmwareCmd.AddCommand(serverFirmwareComponentsListCmd)
	serverFirmwareComponentsListCmd.Flags().StringSliceVar(&serverFlags.firmwareComponentClass, "class", nil, "Only the components of these classes.")
	serverFirmwareComponentsListCmd.Flags().StringSliceVar(&serverFlags.firmwareComponentBaseboard, "baseboard", nil, "Only the components of these GPU baseboards.")

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
	serverFirmwareGenerateAuditCmd.Flags().BoolVar(&serverFlags.firmwareIncludePreflight, "include-preflight", false, "Run the upgrade safety checks and list the planned jobs for each server.")

	// Explain the firmware binary search
	serverFirmwareCmd.AddCommand(serverFirmwareSearchDryRunCmd)
	serverFirmwareSearchDryRunCmd.Flags().Int64Var(&serverFlags.firmwareDryRunBaselineId, "baseline-id", 0, "Search this firmware baseline.")
	serverFirmwareSearchDryRunCmd.Flags().Int64SliceVar(&serverFlags.firmwareDryRunCatalogIds, "catalog-ids", nil, "Search these firmware catalogs.")
	serverFirmwareSearchDryRunCmd.Flags().BoolVar(&serverFlags.firmwareDryRunAllVersions, "all-versions", false, "Include every matching version, not only the latest.")
	serverFirmwareSearchDryRunCmd.MarkFlagsMutuallyExclusive("baseline-id", "catalog-ids")

	// Batch firmware operations
	serverFirmwareCmd.AddCommand(serverFirmwareUpgradeBatchCmd)

	serverFirmwareCmd.AddCommand(serverFirmwareScheduleUpgradeBatchCmd)
	serverFirmwareScheduleUpgradeBatchCmd.Flags().StringVar(&serverFlags.firmwareBatchScheduleTimestamp, "schedule-timestamp", "", "When the firmware upgrade should run (RFC 3339 timestamp).")
	serverFirmwareScheduleUpgradeBatchCmd.Flags().BoolVar(&serverFlags.firmwareBatchConfirmationRequired, "confirmation-required", false, "Require a confirmation before the scheduled upgrade runs.")

	// Drift detection
	serverCmd.AddCommand(serverDriftCmd)
	serverDriftCmd.AddCommand(serverDriftListCmd)
	serverDriftCmd.AddCommand(serverDriftGetCmd)
	serverDriftCmd.AddCommand(serverDriftAcknowledgeCmd)

	serverCmd.AddCommand(serverSnapshotsCmd)
	serverSnapshotsCmd.Flags().StringVar(&serverFlags.snapshotKind, "kind", "", "Restrict the listing to one snapshot class.")

	serverCmd.AddCommand(serverSyncTargetSnapshotCmd)

	// Hardware
	serverCmd.AddCommand(serverHardwareRescanCmd)
	serverHardwareRescanCmd.Flags().BoolVar(&serverFlags.rebootAllowed, "reboot-allowed", false, "Allow the server to be rebooted for full LLDP interface discovery.")

	serverCmd.AddCommand(serverStatisticsCmd)

	serverCmd.AddCommand(serverConnectInterfaceCmd)
	serverConnectInterfaceCmd.Flags().StringVar(&serverFlags.connectInterfaceConfigSource, "config-source", "", "Source of the connection configuration. Can be 'pipe' or path to a JSON file.")
	serverConnectInterfaceCmd.Flags().IntVar(&serverFlags.connectInterfaceId, "interface-id", 0, "The ID of the server interface to connect.")
	serverConnectInterfaceCmd.Flags().StringVar(&serverFlags.connectInterfacePortId, "port-id", "", "The network device port name.")
	serverConnectInterfaceCmd.Flags().StringVar(&serverFlags.connectInterfaceHostname, "hostname", "", "The network device hostname.")
	serverConnectInterfaceCmd.MarkFlagsOneRequired("config-source", "interface-id")
	serverConnectInterfaceCmd.MarkFlagsMutuallyExclusive("config-source", "interface-id")
	serverConnectInterfaceCmd.MarkFlagsMutuallyExclusive("config-source", "port-id")
	serverConnectInterfaceCmd.MarkFlagsMutuallyExclusive("config-source", "hostname")
	serverConnectInterfaceCmd.MarkFlagsRequiredTogether("interface-id", "port-id", "hostname")

	serverCmd.AddCommand(serverSetInterfacesDefaultFabricCmd)
	serverSetInterfacesDefaultFabricCmd.Flags().StringVar(&serverFlags.defaultFabricConfigSource, "config-source", "", "Source of the configuration. Can be 'pipe' or path to a JSON file.")
	serverSetInterfacesDefaultFabricCmd.Flags().StringSliceVar(&serverFlags.defaultFabricInterfaceIds, "interface-id", nil, "IDs of the server interfaces to change.")
	serverSetInterfacesDefaultFabricCmd.Flags().IntVar(&serverFlags.defaultFabricId, "fabric-id", 0, "The ID of the fabric to set as default.")
	serverSetInterfacesDefaultFabricCmd.Flags().BoolVar(&serverFlags.defaultFabricClear, "clear-fabric", false, "Clear the default fabric of the given interfaces.")
	serverSetInterfacesDefaultFabricCmd.MarkFlagsOneRequired("config-source", "interface-id")
	serverSetInterfacesDefaultFabricCmd.MarkFlagsMutuallyExclusive("config-source", "interface-id")
	serverSetInterfacesDefaultFabricCmd.MarkFlagsMutuallyExclusive("config-source", "fabric-id")
	serverSetInterfacesDefaultFabricCmd.MarkFlagsMutuallyExclusive("config-source", "clear-fabric")
	serverSetInterfacesDefaultFabricCmd.MarkFlagsMutuallyExclusive("fabric-id", "clear-fabric")

	serverCmd.AddCommand(serverSetInterfacesRedundancyGroupCmd)
	serverSetInterfacesRedundancyGroupCmd.Flags().StringVar(&serverFlags.redundancyGroupConfigSource, "config-source", "", "Source of the configuration. Can be 'pipe' or path to a JSON file.")
	serverSetInterfacesRedundancyGroupCmd.Flags().StringSliceVar(&serverFlags.redundancyGroupInterfaceIds, "interface-id", nil, "IDs of the server interfaces to change.")
	serverSetInterfacesRedundancyGroupCmd.Flags().IntVar(&serverFlags.redundancyGroupIndex, "group-index", 0, "The redundancy group index.")
	serverSetInterfacesRedundancyGroupCmd.Flags().BoolVar(&serverFlags.redundancyGroupClear, "clear-group", false, "Remove the given interfaces from their redundancy group.")
	serverSetInterfacesRedundancyGroupCmd.MarkFlagsOneRequired("config-source", "interface-id")
	serverSetInterfacesRedundancyGroupCmd.MarkFlagsMutuallyExclusive("config-source", "interface-id")
	serverSetInterfacesRedundancyGroupCmd.MarkFlagsMutuallyExclusive("config-source", "group-index")
	serverSetInterfacesRedundancyGroupCmd.MarkFlagsMutuallyExclusive("config-source", "clear-group")
	serverSetInterfacesRedundancyGroupCmd.MarkFlagsMutuallyExclusive("group-index", "clear-group")

	// Onboarding
	serverCmd.AddCommand(serverRegisterProductionCmd)
	serverRegisterProductionCmd.Flags().StringVar(&serverFlags.registerProductionConfigSource, "config-source", "", "Source of the production server configuration. Can be 'pipe' or path to a JSON file.")
	serverRegisterProductionCmd.Flags().IntVar(&serverFlags.registerProductionSiteId, "site-id", 0, "Site ID where the server is located.")
	serverRegisterProductionCmd.Flags().IntVar(&serverFlags.registerProductionInfrastructure, "infrastructure-id", 0, "ID of the infrastructure the server belongs to.")
	serverRegisterProductionCmd.Flags().IntVar(&serverFlags.registerProductionOsTemplateId, "os-template-id", 0, "ID of the OS template already installed on the server.")
	serverRegisterProductionCmd.Flags().StringVar(&serverFlags.registerProductionAddress, "management-address", "", "IPMI/BMC management IP address.")
	serverRegisterProductionCmd.Flags().StringVar(&serverFlags.registerProductionUsername, "username", "", "IPMI/BMC username.")
	serverRegisterProductionCmd.Flags().StringVar(&serverFlags.registerProductionPassword, "password", "", "IPMI/BMC password.")
	serverRegisterProductionCmd.Flags().StringVar(&serverFlags.registerProductionSerialNumber, "serial-number", "", "Server serial number.")
	serverRegisterProductionCmd.Flags().StringVar(&serverFlags.registerProductionUuid, "uuid", "", "Server UUID.")
	serverRegisterProductionCmd.Flags().StringVar(&serverFlags.registerProductionBmcMacAddress, "bmc-mac-address", "", "BMC MAC address.")
	serverRegisterProductionCmd.Flags().StringVar(&serverFlags.registerProductionVendor, "vendor", "", "Server vendor.")
	serverRegisterProductionCmd.Flags().StringVar(&serverFlags.registerProductionModel, "model", "", "Server model.")
	serverRegisterProductionCmd.Flags().IntVar(&serverFlags.registerProductionProfileId, "registration-profile-id", 0, "ID of the server registration profile to use.")
	serverRegisterProductionCmd.MarkFlagsOneRequired("config-source", "site-id")
	serverRegisterProductionCmd.MarkFlagsMutuallyExclusive("config-source", "site-id")
	serverRegisterProductionCmd.MarkFlagsMutuallyExclusive("config-source", "infrastructure-id")
	serverRegisterProductionCmd.MarkFlagsRequiredTogether("site-id", "infrastructure-id")

	serverCmd.AddCommand(serverImportUnmanagedCmd)
	serverImportUnmanagedCmd.Flags().StringVar(&serverFlags.importUnmanagedConfigSource, "config-source", "", "Source of the import configuration. Can be 'pipe' or path to a JSON file.")
	serverImportUnmanagedCmd.MarkFlagsOneRequired("config-source")

	serverCmd.AddCommand(serverConfigExampleCmd)
}

// buildCreateServerFromFlags builds the server create request from the
// register flags, setting only the optional fields whose flags were given.
func buildCreateServerFromFlags(cmd *cobra.Command) sdk.CreateServer {
	config := sdk.CreateServer{
		SiteId:       int64(serverFlags.siteId),
		SerialNumber: sdk.PtrString(serverFlags.serialNumber),
	}

	optionalString := func(flag string, value string) *string {
		if cmd.Flags().Changed(flag) {
			return sdk.PtrString(value)
		}
		return nil
	}

	config.ServerStatus = optionalString("status", serverFlags.serverStatus)
	config.ManagementAddress = optionalString("management-address", serverFlags.managementAddress)
	config.Username = optionalString("username", serverFlags.username)
	config.Password = optionalString("password", serverFlags.password)
	config.Model = optionalString("model", serverFlags.model)
	config.Vendor = optionalString("vendor", serverFlags.vendor)
	config.ServerUUID = optionalString("uuid", serverFlags.serverUuid)
	config.BmcMacAddress = optionalString("bmc-mac", serverFlags.bmcMacAddress)
	if cmd.Flags().Changed("registration-profile-id") {
		config.RegistrationProfileId = sdk.PtrInt64(int64(serverFlags.registrationProfileId))
	}

	return config
}

// buildRackUpdateFromFlags returns the follow-up update carrying the rack and
// inventory details, or nil when none of them were given.
func buildRackUpdateFromFlags() *sdk.UpdateServer {
	update := sdk.UpdateServer{}
	set := false

	if serverFlags.rackName != "" {
		update.RackName = sdk.PtrString(serverFlags.rackName)
		set = true
	}
	if serverFlags.rackPositionLower != "" {
		update.RackPositionLowerUnit = sdk.PtrString(serverFlags.rackPositionLower)
		set = true
	}
	if serverFlags.rackPositionUpper != "" {
		update.RackPositionUpperUnit = sdk.PtrString(serverFlags.rackPositionUpper)
		set = true
	}
	if serverFlags.inventoryId != "" {
		update.InventoryId = sdk.PtrString(serverFlags.inventoryId)
		set = true
	}

	if !set {
		return nil
	}
	return &update
}
