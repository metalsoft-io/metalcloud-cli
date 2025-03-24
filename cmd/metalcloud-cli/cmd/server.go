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
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SERVERS_READ},
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
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SERVERS_READ},
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
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SERVERS_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(serverFlags.configSource)
			if err != nil {
				return err
			}

			return server.ServerRegister(cmd.Context(), config)
		},
	}
)

// Server Type commands
var (
	serverTypeCmd = &cobra.Command{
		Use:   "type [command]",
		Short: "Server type management",
		Long:  `Server type management commands.`,
	}

	serverTypeListCmd = &cobra.Command{
		Use:          "list",
		Aliases:      []string{"ls"},
		Short:        "Lists server types.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SERVERS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerTypeList(cmd.Context())
		},
	}

	serverTypeGetCmd = &cobra.Command{
		Use:          "get",
		Aliases:      []string{"show"},
		Short:        "Get server type info.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SERVERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerTypeGet(cmd.Context(), args[0])
		},
	}
)

// Server Cleanup Policy commands
var (
	serverCleanupPolicyCmd = &cobra.Command{
		Use:     "cleanup-policy [command]",
		Aliases: []string{"cp"},
		Short:   "Server cleanup policy management",
		Long:    `Server cleanup policy management commands.`,
	}

	serverCleanupPolicyListCmd = &cobra.Command{
		Use:          "list",
		Aliases:      []string{"ls"},
		Short:        "Lists server cleanup policies.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SERVERS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.CleanupPolicyList(cmd.Context())
		},
	}

	serverCleanupPolicyGetCmd = &cobra.Command{
		Use:          "get",
		Aliases:      []string{"show"},
		Short:        "Get server cleanup policy.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SERVERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.CleanupPolicyGet(cmd.Context(), args[0])
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

	// Server Type commands
	serverCmd.AddCommand(serverTypeCmd)
	serverTypeCmd.AddCommand(serverTypeListCmd)
	serverTypeCmd.AddCommand(serverTypeGetCmd)

	// Server Cleanup Policy commands
	serverCmd.AddCommand(serverCleanupPolicyCmd)
	serverCleanupPolicyCmd.AddCommand(serverCleanupPolicyListCmd)
	serverCleanupPolicyCmd.AddCommand(serverCleanupPolicyGetCmd)
}
