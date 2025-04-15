package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/server"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

// Server commands
var (
	serverFlags = struct {
		showCredentials bool
		filterStatus    string
		filterType      string
		configSource    string
		powerAction     string
	}{}

	serverCmd = &cobra.Command{
		Use:     "server [command]",
		Aliases: []string{"srv", "servers"},
		Short:   "Server management",
		Long:    `Server management commands.`,
	}

	serverListCmd = &cobra.Command{
		Use:          "list",
		Aliases:      []string{"ls"},
		Short:        "List all servers.",
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
		Use:          "get server_id",
		Aliases:      []string{"show"},
		Short:        "Get server info.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerGet(cmd.Context(), args[0],
				serverFlags.showCredentials)
		},
	}

	serverRegisterCmd = &cobra.Command{
		Use:          "register",
		Aliases:      []string{"new"},
		Short:        "Register a server.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(serverFlags.configSource)
			if err != nil {
				return err
			}

			return server.ServerRegister(cmd.Context(), config)
		},
	}

	serverReRegisterCmd = &cobra.Command{
		Use:          "re-register server_id",
		Short:        "Re-register an existing server.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerReRegister(cmd.Context(), args[0])
		},
	}

	serverFactoryResetCmd = &cobra.Command{
		Use:          "factory-reset server_id",
		Short:        "Reset a server to factory defaults.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerFactoryReset(cmd.Context(), args[0])
		},
	}

	serverArchiveCmd = &cobra.Command{
		Use:          "archive server_id",
		Short:        "Archive a server.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerArchive(cmd.Context(), args[0])
		},
	}

	serverDeleteCmd = &cobra.Command{
		Use:          "delete server_id",
		Aliases:      []string{"rm"},
		Short:        "Delete a server.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerDelete(cmd.Context(), args[0])
		},
	}

	serverPowerCmd = &cobra.Command{
		Use:          "power server_id",
		Short:        "Control server power state.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerPower(cmd.Context(), args[0], serverFlags.powerAction)
		},
	}

	serverPowerStatusCmd = &cobra.Command{
		Use:          "power-status server_id",
		Short:        "Get server power status.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerPowerStatus(cmd.Context(), args[0])
		},
	}

	serverUpdateCmd = &cobra.Command{
		Use:          "update server_id",
		Short:        "Update server information.",
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
		Use:          "update-ipmi-credentials server_id username password",
		Short:        "Update server IPMI credentials.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerUpdateIpmiCredentials(cmd.Context(), args[0], args[1], args[2])
		},
	}

	serverEnableSnmpCmd = &cobra.Command{
		Use:          "enable-snmp server_id",
		Short:        "Enable SNMP on server.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerEnableSnmp(cmd.Context(), args[0])
		},
	}

	serverEnableSyslogCmd = &cobra.Command{
		Use:          "enable-syslog server_id",
		Short:        "Enable remote syslog for a server.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerEnableSyslog(cmd.Context(), args[0])
		},
	}

	serverVncInfoCmd = &cobra.Command{
		Use:          "vnc-info server_id",
		Short:        "Get server VNC information.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerVncInfo(cmd.Context(), args[0])
		},
	}

	serverRemoteConsoleInfoCmd = &cobra.Command{
		Use:          "console-info server_id",
		Short:        "Get server remote console information.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerRemoteConsoleInfo(cmd.Context(), args[0])
		},
	}

	serverCapabilitiesCmd = &cobra.Command{
		Use:          "capabilities server_id",
		Short:        "Get server capabilities.",
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
		Long:    `Server firmware management commands.`,
	}

	serverFirmwareComponentsListCmd = &cobra.Command{
		Use:          "components server_id",
		Aliases:      []string{"list-components", "ls"},
		Short:        "List firmware components for a server.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerFirmwareComponentsList(cmd.Context(), args[0])
		},
	}

	serverFirmwareComponentGetCmd = &cobra.Command{
		Use:          "component-info server_id component_id",
		Aliases:      []string{"get-component"},
		Short:        "Get firmware component information.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerFirmwareComponentGet(cmd.Context(), args[0], args[1])
		},
	}

	serverFirmwareComponentUpdateCmd = &cobra.Command{
		Use:          "update-component server_id component_id",
		Short:        "Update firmware component settings.",
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
		Use:          "update-info server_id",
		Short:        "Update firmware information for a server.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerFirmwareUpdateInfo(cmd.Context(), args[0])
		},
	}

	serverFirmwareInventoryCmd = &cobra.Command{
		Use:          "inventory server_id",
		Short:        "Get firmware inventory from redfish.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerFirmwareInventory(cmd.Context(), args[0])
		},
	}

	serverFirmwareUpgradeCmd = &cobra.Command{
		Use:          "upgrade server_id",
		Short:        "Upgrade firmware for all components on a server.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerFirmwareUpgrade(cmd.Context(), args[0])
		},
	}

	serverFirmwareComponentUpgradeCmd = &cobra.Command{
		Use:          "upgrade-component server_id component_id",
		Short:        "Upgrade firmware for a specific component.",
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
		Use:          "schedule-upgrade server_id",
		Short:        "Schedule a firmware upgrade for a server.",
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
		Use:          "fetch-versions server_id",
		Short:        "Fetch available firmware versions for a server.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerFirmwareFetchVersions(cmd.Context(), args[0])
		},
	}

	serverFirmwareGenerateAuditCmd = &cobra.Command{
		Use:          "generate-audit",
		Short:        "Generate firmware upgrade audit for servers.",
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
	serverListCmd.Flags().StringVar(&serverFlags.filterStatus, "filter-status", "", "Filter the result by server status.")
	serverListCmd.Flags().StringVar(&serverFlags.filterType, "filter-type", "", "Filter the result by server type.")

	serverCmd.AddCommand(serverGetCmd)
	serverGetCmd.Flags().BoolVar(&serverFlags.showCredentials, "show-credentials", false, "If set returns the server IPMI credentials.")

	serverCmd.AddCommand(serverRegisterCmd)
	serverRegisterCmd.Flags().StringVar(&serverFlags.configSource, "config-source", "", "Source of the new server configuration. Can be 'pipe' or path to a JSON file.")
	serverRegisterCmd.MarkFlagsOneRequired("config-source")

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
