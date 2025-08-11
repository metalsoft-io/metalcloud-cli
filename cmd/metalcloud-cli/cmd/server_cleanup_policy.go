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
		Short:   "Manage server cleanup policies for automated server maintenance",
		Long: `Manage server cleanup policies that define automated maintenance procedures for servers.

Server cleanup policies control how and when servers are automatically cleaned up,
including configuration of cleanup schedules, retention policies, and maintenance actions.

Available Commands:
  list    List all server cleanup policies
  get     Get details of a specific server cleanup policy

Examples:
  # List all server cleanup policies
  metalcloud-cli server-cleanup-policy list

  # Get details of a specific policy
  metalcloud-cli server-cleanup-policy get policy-123

  # Using short aliases
  metalcloud-cli scp list
  metalcloud-cli srv-cp get policy-123`,
	}

	serverCleanupPolicyListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all server cleanup policies",
		Long: `List all server cleanup policies configured in the system.

This command displays a table of all available server cleanup policies with their
key attributes including ID, name, status, and configuration summary.

Output Format:
  By default, output is formatted as a table. Use global flags to change output format:
  --format=json    JSON output
  --format=csv     CSV output
  --format=yaml    YAML output

Required Permissions:
  - server_cleanup_policies:read

Examples:
  # List all server cleanup policies in table format
  metalcloud-cli server-cleanup-policy list

  # List policies in JSON format
  metalcloud-cli server-cleanup-policy list --format=json

  # List policies with custom output format
  metalcloud-cli scp ls --format=csv`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_CLEANUP_POLICIES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_cleanup_policy.CleanupPolicyList(cmd.Context())
		},
	}

	serverCleanupPolicyGetCmd = &cobra.Command{
		Use:     "get <policy-id>",
		Aliases: []string{"show"},
		Short:   "Get detailed information about a specific server cleanup policy",
		Long: `Get detailed information about a specific server cleanup policy by its ID.

This command retrieves and displays comprehensive information about a server cleanup
policy including its configuration, schedule, retention settings, and associated
maintenance actions.

Arguments:
  policy-id    The unique identifier of the server cleanup policy to retrieve.
               This can be either the numeric ID or the policy name.

Output Format:
  By default, output is formatted as a detailed table. Use global flags to change output format:
  --format=json    JSON output (recommended for programmatic use)
  --format=yaml    YAML output
  --format=csv     CSV output (limited detail)

Required Permissions:
  - server_cleanup_policies:read

Examples:
  # Get policy details by numeric ID
  metalcloud-cli server-cleanup-policy get 12345

  # Get policy details by name
  metalcloud-cli server-cleanup-policy get "weekly-maintenance"

  # Get policy details in JSON format
  metalcloud-cli scp get 12345 --format=json

  # Using alias commands
  metalcloud-cli srv-cp show policy-name`,
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
