package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/server_cleanup_policy"
	"github.com/spf13/cobra"
)

// Server Cleanup Policy commands
var (
	serverCleanupPolicyCmd = &cobra.Command{
		Use:     "server-cleanup-policy [command]",
		Aliases: []string{"srv-cp", "scp"},
		Short:   "Server cleanup policy management",
		Long:    `Server cleanup policy management commands.`,
	}

	serverCleanupPolicyListCmd = &cobra.Command{
		Use:          "list",
		Aliases:      []string{"ls"},
		Short:        "Lists server cleanup policies.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_CLEANUP_POLICIES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_cleanup_policy.CleanupPolicyList(cmd.Context())
		},
	}

	serverCleanupPolicyGetCmd = &cobra.Command{
		Use:          "get",
		Aliases:      []string{"show"},
		Short:        "Get server cleanup policy.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_CLEANUP_POLICIES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_cleanup_policy.CleanupPolicyGet(cmd.Context(), args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(serverCleanupPolicyCmd)

	// Server Cleanup Policy commands
	serverCleanupPolicyCmd.AddCommand(serverCleanupPolicyListCmd)
	serverCleanupPolicyCmd.AddCommand(serverCleanupPolicyGetCmd)
}
