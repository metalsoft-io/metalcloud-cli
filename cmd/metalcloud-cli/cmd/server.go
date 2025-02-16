package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/server"
	"github.com/spf13/cobra"
)

var (
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
			return server.ServerList(cmd.Context())
		},
	}

	serverGetCmd = &cobra.Command{
		Use:          "get",
		Aliases:      []string{"show"},
		Short:        "Get server info.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SERVERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ServerGet(cmd.Context(), args[0])
		},
	}

	// server cleanup policy commands
	serverCleanupPolicyCmd = &cobra.Command{
		Use:     "server cleanup-policy [command]",
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

	serverCmd.AddCommand(serverListCmd)
	serverCmd.AddCommand(serverGetCmd)

	serverCmd.AddCommand(serverCleanupPolicyCmd)
	serverCleanupPolicyCmd.AddCommand(serverCleanupPolicyListCmd)
	serverCleanupPolicyCmd.AddCommand(serverCleanupPolicyGetCmd)
}
