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
}
