package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/server_cleanup_policy"
	"github.com/spf13/cobra"
)

// Server Cleanup Policy commands
var (
	serverCleanupPolicyFlags = struct {
		label               string
		cleanupDrives       bool
		recreateRaid        bool
		disableEmbeddedNics bool
		raidOneDrive        string
		raidTwoDrives       string
		raidEvenDrives      string
		raidOddDrives       string
		skipRaidActions     string
	}{}

	serverCleanupPolicyCmd = &cobra.Command{
		Use:     "server-cleanup-policy [command]",
		Aliases: []string{"srv-cp", "scp"},
		Short:   "Manage server cleanup policies for automated server maintenance",
		Long: `Manage server cleanup policies that define automated maintenance procedures for servers.

Server cleanup policies control how and when servers are automatically cleaned up,
including configuration of cleanup schedules, retention policies, and maintenance actions.

Available Commands:
  list      List all server cleanup policies
  get       Get details of a specific server cleanup policy
  create    Create a new server cleanup policy
  update    Update an existing server cleanup policy
  delete    Delete a server cleanup policy

Examples:
  # List all server cleanup policies
  metalcloud-cli server-cleanup-policy list

  # Get details of a specific policy
  metalcloud-cli server-cleanup-policy get policy-123

  # Create a new server cleanup policy
  metalcloud-cli server-cleanup-policy create --label "my-policy" --cleanup-drives 1 --recreate-raid 1

  # Update an existing policy
  metalcloud-cli server-cleanup-policy update 123 --label "updated-policy"

  # Delete a policy
  metalcloud-cli server-cleanup-policy delete 123

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

	serverCleanupPolicyCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new server cleanup policy",
		Long: `Create a new server cleanup policy with specified configuration.

This command creates a new server cleanup policy that defines automated maintenance
procedures for servers. The policy configuration includes cleanup behaviors for
drives, RAID settings, and embedded NICs.

Required Flags:
  --label                                Label for the server cleanup policy
  --cleanup-drives                       Enable cleanup drives for OOB enabled servers
  --recreate-raid                        Enable RAID recreation
  --disable-embedded-nics                Enable disabling embedded NICs
  --raid-one-drive                       RAID configuration for single drive (e.g. "raid0")
  --raid-two-drives                      RAID configuration for two drives (e.g. "raid1")
  --raid-even-drives                     RAID configuration for even number of drives (>2)
  --raid-odd-drives                      RAID configuration for odd number of drives (>1)
  --skip-raid-actions                    Comma-separated list of RAID actions to skip

Required Permissions:
  - server_cleanup_policies:write

Examples:
  # Create a basic cleanup policy
  metalcloud-cli server-cleanup-policy create --label "basic-cleanup" \
    --cleanup-drives --recreate-raid --disable-embedded-nics \
    --raid-one-drive "raid0" --raid-two-drives "raid1" \
    --raid-even-drives "raid10" --raid-odd-drives "raid5" \
    --skip-raid-actions "cleanup"`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_CLEANUP_POLICIES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_cleanup_policy.CleanupPolicyCreate(cmd.Context(),
				serverCleanupPolicyFlags.label,
				serverCleanupPolicyFlags.cleanupDrives,
				serverCleanupPolicyFlags.recreateRaid,
				serverCleanupPolicyFlags.disableEmbeddedNics,
				serverCleanupPolicyFlags.raidOneDrive,
				serverCleanupPolicyFlags.raidTwoDrives,
				serverCleanupPolicyFlags.raidEvenDrives,
				serverCleanupPolicyFlags.raidOddDrives,
				serverCleanupPolicyFlags.skipRaidActions)
		},
	}

	serverCleanupPolicyUpdateCmd = &cobra.Command{
		Use:   "update <policy-id>",
		Short: "Update an existing server cleanup policy",
		Long: `Update an existing server cleanup policy by its ID.

This command updates the configuration of an existing server cleanup policy.
Only the flags that are provided will be updated, other settings remain unchanged.

Arguments:
  policy-id    The unique identifier of the server cleanup policy to update.
               This must be the numeric ID of the policy.

Optional Flags:
  --label                                New label for the server cleanup policy
  --cleanup-drives                       Enable cleanup drives for OOB enabled servers
  --recreate-raid                        Enable RAID recreation
  --disable-embedded-nics                Enable disabling embedded NICs
  --raid-one-drive                       RAID configuration for single drive (e.g. "raid0")
  --raid-two-drives                      RAID configuration for two drives (e.g. "raid1")
  --raid-even-drives                     RAID configuration for even number of drives (>2)
  --raid-odd-drives                      RAID configuration for odd number of drives (>1)
  --skip-raid-actions                    Comma-separated list of RAID actions to skip

Required Permissions:
  - server_cleanup_policies:write

Examples:
  # Update policy label only
  metalcloud-cli server-cleanup-policy update 123 --label "updated-policy"

  # Update multiple settings
  metalcloud-cli scp update 123 --cleanup-drives --recreate-raid

  # Update RAID configurations
  metalcloud-cli srv-cp update 123 --raid-one-drive "RAID1" --raid-two-drives "RAID10"`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_CLEANUP_POLICIES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_cleanup_policy.CleanupPolicyUpdate(cmd.Context(), args[0],
				serverCleanupPolicyFlags.label,
				serverCleanupPolicyFlags.cleanupDrives,
				serverCleanupPolicyFlags.recreateRaid,
				serverCleanupPolicyFlags.disableEmbeddedNics,
				serverCleanupPolicyFlags.raidOneDrive,
				serverCleanupPolicyFlags.raidTwoDrives,
				serverCleanupPolicyFlags.raidEvenDrives,
				serverCleanupPolicyFlags.raidOddDrives,
				serverCleanupPolicyFlags.skipRaidActions,
				cmd)
		},
	}

	serverCleanupPolicyDeleteCmd = &cobra.Command{
		Use:     "delete <policy-id>",
		Aliases: []string{"rm", "remove"},
		Short:   "Delete a server cleanup policy",
		Long: `Delete a server cleanup policy by its ID.

This command permanently removes a server cleanup policy from the system.
This action cannot be undone, so use with caution.

Arguments:
  policy-id    The unique identifier of the server cleanup policy to delete.
               This must be the numeric ID of the policy.

Required Permissions:
  - server_cleanup_policies:write

Examples:
  # Delete a policy by ID
  metalcloud-cli server-cleanup-policy delete 123

  # Using aliases
  metalcloud-cli scp rm 456
  metalcloud-cli srv-cp remove 789`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_CLEANUP_POLICIES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_cleanup_policy.CleanupPolicyDelete(cmd.Context(), args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(serverCleanupPolicyCmd)

	serverCleanupPolicyCmd.AddCommand(serverCleanupPolicyListCmd)

	serverCleanupPolicyCmd.AddCommand(serverCleanupPolicyGetCmd)

	serverCleanupPolicyCmd.AddCommand(serverCleanupPolicyCreateCmd)
	serverCleanupPolicyCreateCmd.Flags().StringVar(&serverCleanupPolicyFlags.label, "label", "", "Label for the server cleanup policy")
	serverCleanupPolicyCreateCmd.Flags().BoolVar(&serverCleanupPolicyFlags.cleanupDrives, "cleanup-drives", false, "Enable cleanup drives for OOB enabled servers")
	serverCleanupPolicyCreateCmd.Flags().BoolVar(&serverCleanupPolicyFlags.recreateRaid, "recreate-raid", false, "Enable RAID recreation")
	serverCleanupPolicyCreateCmd.Flags().BoolVar(&serverCleanupPolicyFlags.disableEmbeddedNics, "disable-embedded-nics", false, "Enable disabling embedded NICs")
	serverCleanupPolicyCreateCmd.Flags().StringVar(&serverCleanupPolicyFlags.raidOneDrive, "raid-one-drive", "", "RAID configuration for single drive")
	serverCleanupPolicyCreateCmd.Flags().StringVar(&serverCleanupPolicyFlags.raidTwoDrives, "raid-two-drives", "", "RAID configuration for two drives")
	serverCleanupPolicyCreateCmd.Flags().StringVar(&serverCleanupPolicyFlags.raidEvenDrives, "raid-even-drives", "", "RAID configuration for even number of drives (>2)")
	serverCleanupPolicyCreateCmd.Flags().StringVar(&serverCleanupPolicyFlags.raidOddDrives, "raid-odd-drives", "", "RAID configuration for odd number of drives (>1)")
	serverCleanupPolicyCreateCmd.Flags().StringVar(&serverCleanupPolicyFlags.skipRaidActions, "skip-raid-actions", "", "Comma-separated list of RAID actions to skip")
	serverCleanupPolicyCreateCmd.MarkFlagRequired("label")
	serverCleanupPolicyCreateCmd.MarkFlagRequired("raid-one-drive")
	serverCleanupPolicyCreateCmd.MarkFlagRequired("raid-two-drives")
	serverCleanupPolicyCreateCmd.MarkFlagRequired("raid-even-drives")
	serverCleanupPolicyCreateCmd.MarkFlagRequired("raid-odd-drives")

	serverCleanupPolicyCmd.AddCommand(serverCleanupPolicyUpdateCmd)
	serverCleanupPolicyUpdateCmd.Flags().StringVar(&serverCleanupPolicyFlags.label, "label", "", "New label for the server cleanup policy")
	serverCleanupPolicyUpdateCmd.Flags().BoolVar(&serverCleanupPolicyFlags.cleanupDrives, "cleanup-drives", false, "Enable cleanup drives for OOB enabled servers")
	serverCleanupPolicyUpdateCmd.Flags().BoolVar(&serverCleanupPolicyFlags.recreateRaid, "recreate-raid", false, "Enable RAID recreation")
	serverCleanupPolicyUpdateCmd.Flags().BoolVar(&serverCleanupPolicyFlags.disableEmbeddedNics, "disable-embedded-nics", false, "Enable disabling embedded NICs")
	serverCleanupPolicyUpdateCmd.Flags().StringVar(&serverCleanupPolicyFlags.raidOneDrive, "raid-one-drive", "", "RAID configuration for single drive")
	serverCleanupPolicyUpdateCmd.Flags().StringVar(&serverCleanupPolicyFlags.raidTwoDrives, "raid-two-drives", "", "RAID configuration for two drives")
	serverCleanupPolicyUpdateCmd.Flags().StringVar(&serverCleanupPolicyFlags.raidEvenDrives, "raid-even-drives", "", "RAID configuration for even number of drives (>2)")
	serverCleanupPolicyUpdateCmd.Flags().StringVar(&serverCleanupPolicyFlags.raidOddDrives, "raid-odd-drives", "", "RAID configuration for odd number of drives (>1)")
	serverCleanupPolicyUpdateCmd.Flags().StringVar(&serverCleanupPolicyFlags.skipRaidActions, "skip-raid-actions", "", "Comma-separated list of RAID actions to skip")

	serverCleanupPolicyCmd.AddCommand(serverCleanupPolicyDeleteCmd)
}
