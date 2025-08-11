package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/firmware_policy"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	firmwarePolicyFlags = struct {
		configSource string
	}{}

	firmwarePolicyCmd = &cobra.Command{
		Use:     "firmware-policy [command]",
		Aliases: []string{"fw-policy"},
		Short:   "Manage server firmware upgrade policies and global firmware configurations",
		Long: `Manage server firmware upgrade policies and global firmware configurations.

Firmware policies define rules for automatically upgrading server firmware based on
server properties like OS, server type, or instance groups. Global firmware
configuration controls when and how firmware upgrades are applied system-wide.

Available commands:
  list                    List all firmware policies
  get                     Get firmware policy details
  create                  Create a new firmware policy
  update                  Update an existing firmware policy
  delete                  Delete a firmware policy
  generate-audit          Generate compliance audit for a policy
  apply-with-groups       Apply policies linked to server instance groups
  apply-without-groups    Apply policies not linked to server instance groups
  config-example          Show firmware policy configuration example
  global-config           Manage global firmware configuration

Examples:
  # List all firmware policies
  metalcloud-cli firmware-policy list

  # Create a new policy from JSON file
  metalcloud-cli firmware-policy create --config-source policy.json

  # Apply all policies linked to server groups
  metalcloud-cli firmware-policy apply-with-groups`,
	}

	firmwarePolicyListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all firmware upgrade policies",
		Long: `List all firmware upgrade policies in the system.

This command displays all firmware policies with their basic information including
ID, label, status, action, owner, and timestamps. The output shows which policies
are active and what type of firmware upgrade action they will perform.

No flags are required for this command.

Examples:
  # List all firmware policies
  metalcloud-cli firmware-policy list
  
  # List policies using alias
  metalcloud-cli fw-policy ls`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return firmware_policy.FirmwarePolicyList(cmd.Context())
		},
	}

	firmwarePolicyGetCmd = &cobra.Command{
		Use:     "get policy_id",
		Aliases: []string{"show"},
		Short:   "Get detailed information about a specific firmware policy",
		Long: `Get detailed information about a specific firmware policy including its configuration,
rules, associated server instance groups, and current status.

This command retrieves and displays all available information for a single firmware
policy, including the rules that determine which servers the policy applies to,
the firmware upgrade action to be performed, and any server instance groups
that are linked to this policy.

Required arguments:
  policy_id               The unique identifier (numeric ID) of the firmware policy

Examples:
  # Get details for firmware policy with ID 123
  metalcloud-cli firmware-policy get 123
  
  # Show policy details using alias
  metalcloud-cli fw-policy show 456`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return firmware_policy.FirmwarePolicyGet(cmd.Context(), args[0])
		},
	}

	firmwarePolicyCreateCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Short:   "Create a new firmware upgrade policy",
		Long: `Create a new firmware upgrade policy with the specified configuration.

This command creates a new firmware policy that defines rules for automatically
upgrading server firmware. The policy configuration must be provided via JSON
input that specifies the policy label, action, rules, and optionally associated
server instance groups.

Required flags:
  --config-source         Source of the firmware policy configuration
                          Values: 'pipe' (read from stdin) or path to JSON file

The configuration JSON should include:
  - label: A descriptive name for the policy
  - action: The upgrade action (e.g., "upgrade", "downgrade")
  - rules: Array of rules defining which servers the policy applies to
  - userIdOwner: (optional) User ID of the policy owner
  - serverInstanceGroupIds: (optional) Array of server instance group IDs

Examples:
  # Create policy from JSON file
  metalcloud-cli firmware-policy create --config-source policy.json
  
  # Create policy from stdin
  echo '{"label":"test-policy","action":"upgrade"}' | metalcloud-cli fw-policy create --config-source pipe
  
  # Get configuration example first
  metalcloud-cli firmware-policy config-example > policy.json
  metalcloud-cli firmware-policy create --config-source policy.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(firmwarePolicyFlags.configSource)
			if err != nil {
				return err
			}
			return firmware_policy.FirmwarePolicyCreate(cmd.Context(), config)
		},
	}

	firmwarePolicyUpdateCmd = &cobra.Command{
		Use:     "update policy_id",
		Aliases: []string{"edit"},
		Short:   "Update an existing firmware upgrade policy",
		Long: `Update an existing firmware upgrade policy with new configuration.

This command allows you to modify an existing firmware policy by providing
updated configuration data. You can change the policy's label, action, rules,
and associated server instance groups. The policy ID cannot be changed.

Required arguments:
  policy_id               The unique identifier (numeric ID) of the firmware policy

Required flags:
  --config-source         Source of the firmware policy configuration updates
                          Values: 'pipe' (read from stdin) or path to JSON file

The configuration JSON can include any of these fields:
  - label: (optional) Updated descriptive name for the policy
  - action: (optional) Updated upgrade action (e.g., "upgrade", "downgrade") 
  - rules: (optional) Updated array of rules defining server selection criteria
  - userIdOwner: (optional) Updated user ID of the policy owner
  - serverInstanceGroupIds: (optional) Updated array of server instance group IDs

Note: Only provide the fields you want to update. Missing fields will retain
their current values.

Examples:
  # Update policy from JSON file
  metalcloud-cli firmware-policy update 123 --config-source policy-updates.json
  
  # Update policy label only via stdin
  echo '{"label":"updated-policy-name"}' | metalcloud-cli fw-policy update 456 --config-source pipe
  
  # Update policy action and rules
  metalcloud-cli firmware-policy update 789 --config-source new-config.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(firmwarePolicyFlags.configSource)
			if err != nil {
				return err
			}
			return firmware_policy.FirmwarePolicyUpdate(cmd.Context(), args[0], config)
		},
	}

	firmwarePolicyDeleteCmd = &cobra.Command{
		Use:     "delete policy_id",
		Aliases: []string{"rm"},
		Short:   "Delete a firmware upgrade policy permanently",
		Long: `Delete a firmware upgrade policy permanently from the system.

This command removes a firmware policy and all its associated configuration
including rules and server instance group associations. This action cannot
be undone and will stop any automated firmware upgrades controlled by this policy.

Required arguments:
  policy_id               The unique identifier (numeric ID) of the firmware policy

Warning: This action is irreversible. Make sure you no longer need the policy
before deleting it. Consider getting a backup of the policy configuration first
using the 'get' command.

Examples:
  # Delete firmware policy with ID 123
  metalcloud-cli firmware-policy delete 123
  
  # Delete policy using alias
  metalcloud-cli fw-policy rm 456
  
  # Get policy details before deletion (recommended)
  metalcloud-cli firmware-policy get 123
  metalcloud-cli firmware-policy delete 123`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return firmware_policy.FirmwarePolicyDelete(cmd.Context(), args[0])
		},
	}

	firmwarePolicyAuditCmd = &cobra.Command{
		Use:     "generate-audit policy_id",
		Aliases: []string{"audit"},
		Short:   "Generate compliance audit report for a firmware policy",
		Long: `Generate a compliance audit report for a firmware policy to analyze server firmware status.

This command analyzes the current firmware status of all servers that match the
specified policy rules and generates a detailed compliance report. The audit shows
which servers are compliant with the policy requirements, which need updates,
and provides detailed firmware version information.

Required arguments:
  policy_id               The unique identifier (numeric ID) of the firmware policy

The audit report includes:
  - List of servers matching the policy rules
  - Current firmware versions for each server component
  - Compliance status for each server
  - Recommended firmware updates
  - Servers that would be affected by policy execution

Examples:
  # Generate audit for firmware policy with ID 123
  metalcloud-cli firmware-policy generate-audit 123
  
  # Generate audit using alias
  metalcloud-cli fw-policy audit 456
  
  # Save audit results to file
  metalcloud-cli firmware-policy generate-audit 789 > audit-report.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return firmware_policy.FirmwarePolicyGenerateAudit(cmd.Context(), args[0])
		},
	}

	firmwarePolicyApplyWithGroupsCmd = &cobra.Command{
		Use:   "apply-with-groups",
		Short: "Apply all firmware policies linked to server instance groups",
		Long: `Apply all firmware policies that are linked to server instance groups.

This command executes all active firmware policies that have server instance groups
associated with them. It will only affect servers that belong to the specified
server instance groups in each policy's configuration.

The command respects the global firmware configuration settings for timing and
scheduling constraints. If global firmware upgrades are disabled or outside
the configured time window, the command may be blocked.

No flags or arguments are required for this command.

Prerequisites:
  - At least one firmware policy must exist with server instance groups assigned
  - Global firmware configuration must allow policy execution
  - Servers in the target groups must be accessible and eligible for firmware updates

Examples:
  # Apply all policies linked to server instance groups
  metalcloud-cli firmware-policy apply-with-groups
  
  # Apply policies using alias
  metalcloud-cli fw-policy apply-with-groups
  
  # Check global config before applying
  metalcloud-cli firmware-policy global-config get
  metalcloud-cli firmware-policy apply-with-groups`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			return firmware_policy.FirmwarePolicyApplyWithGroups(cmd.Context())
		},
	}

	firmwarePolicyApplyWithoutGroupsCmd = &cobra.Command{
		Use:   "apply-without-groups",
		Short: "Apply all firmware policies not linked to server instance groups",
		Long: `Apply all firmware policies that are not linked to server instance groups.

This command executes all active firmware policies that do not have server instance
groups associated with them. These policies will apply their rules globally across
all servers in the system that match the policy criteria.

The command respects the global firmware configuration settings for timing and
scheduling constraints. If global firmware upgrades are disabled or outside
the configured time window, the command may be blocked.

No flags or arguments are required for this command.

Prerequisites:
  - At least one firmware policy must exist without server instance groups assigned
  - Global firmware configuration must allow policy execution
  - Target servers must be accessible and eligible for firmware updates

Examples:
  # Apply all policies not linked to server instance groups
  metalcloud-cli firmware-policy apply-without-groups
  
  # Apply policies using alias
  metalcloud-cli fw-policy apply-without-groups
  
  # Check which policies will be affected first
  metalcloud-cli firmware-policy list
  metalcloud-cli firmware-policy apply-without-groups`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			return firmware_policy.FirmwarePolicyApplyWithoutGroups(cmd.Context())
		},
	}

	firmwarePolicyConfigExampleCmd = &cobra.Command{
		Use:   "config-example",
		Short: "Show example firmware policy configuration in JSON format",
		Long: `Show an example firmware policy configuration in JSON format.

This command displays a sample JSON configuration that can be used as a template
for creating new firmware policies. The example includes all available fields
with sample values and explains the structure of policy rules and server
instance group associations.

The example output can be saved to a file and modified to create your own
firmware policy configurations.

No flags or arguments are required for this command.

Examples:
  # Show example configuration
  metalcloud-cli firmware-policy config-example
  
  # Save example to file for editing
  metalcloud-cli fw-policy config-example > my-policy.json
  
  # Use example as template for creating policy
  metalcloud-cli firmware-policy config-example > policy.json
  # Edit policy.json with your values
  metalcloud-cli firmware-policy create --config-source policy.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return firmware_policy.FirmwarePolicyConfigExample(cmd.Context())
		},
	}

	// Global firmware configuration commands
	globalFirmwareConfigCmd = &cobra.Command{
		Use:     "global-config",
		Aliases: []string{"global"},
		Short:   "Manage global firmware configuration settings",
		Long: `Manage global firmware configuration settings that control system-wide firmware upgrade behavior.

The global firmware configuration defines when firmware upgrades can be executed,
whether they are enabled system-wide, and other global constraints that affect
all firmware policies. This configuration acts as a master control for the
entire firmware upgrade system.

Available subcommands:
  get                     Get current global firmware configuration
  update                  Update global firmware configuration
  config-example          Show example global configuration

Examples:
  # Get current global configuration
  metalcloud-cli firmware-policy global-config get
  
  # Update global configuration from file
  metalcloud-cli firmware-policy global-config update --config-source config.json`,
		SilenceUsage: true,
	}

	globalFirmwareConfigGetCmd = &cobra.Command{
		Use:   "get",
		Short: "Get current global firmware configuration settings",
		Long: `Get current global firmware configuration settings that control system-wide firmware upgrade behavior.

This command retrieves and displays the global firmware configuration which includes
settings such as whether firmware upgrades are enabled globally, upgrade time windows,
scheduling constraints, and other system-wide policies that affect all firmware
upgrade operations.

The global configuration acts as a master switch and constraint system for all
firmware policies. Even if individual policies are active, they must comply with
the global configuration settings.

No flags or arguments are required for this command.

Examples:
  # Get current global firmware configuration
  metalcloud-cli firmware-policy global-config get
  
  # Get global config using alias
  metalcloud-cli fw-policy global get`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return firmware_policy.GetGlobalFirmwareConfiguration(cmd.Context())
		},
	}

	globalFirmwareConfigUpdateCmd = &cobra.Command{
		Use:     "update",
		Aliases: []string{"edit"},
		Short:   "Update global firmware configuration settings",
		Long: `Update global firmware configuration settings that control system-wide firmware upgrade behavior.

This command allows you to modify the global firmware configuration which acts as
a master control for all firmware upgrade operations. You can enable/disable
firmware upgrades globally, set time windows for when upgrades can occur,
and configure other system-wide constraints.

Required flags:
  --config-source         Source of the global firmware configuration updates
                          Values: 'pipe' (read from stdin) or path to JSON file

The configuration JSON can include any of these fields:
  - activated: (optional) Boolean to enable/disable firmware upgrades globally
  - upgradeStartTime: (optional) ISO 8601 timestamp for upgrade window start
  - upgradeEndTime: (optional) ISO 8601 timestamp for upgrade window end
  - other global firmware settings: (varies based on API specification)

Note: Only provide the fields you want to update. Missing fields will retain
their current values.

Examples:
  # Update global config from JSON file
  metalcloud-cli firmware-policy global-config update --config-source global-config.json
  
  # Enable firmware upgrades globally via stdin
  echo '{"activated":true}' | metalcloud-cli fw-policy global update --config-source pipe
  
  # Set upgrade time window
  metalcloud-cli firmware-policy global-config update --config-source time-window.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(firmwarePolicyFlags.configSource)
			if err != nil {
				return err
			}
			return firmware_policy.UpdateGlobalFirmwareConfiguration(cmd.Context(), config)
		},
	}

	globalFirmwareConfigExampleCmd = &cobra.Command{
		Use:   "config-example",
		Short: "Show example global firmware configuration in JSON format",
		Long: `Show an example global firmware configuration in JSON format.

This command displays a sample JSON configuration that can be used as a template
for updating the global firmware configuration. The example includes all available
fields with sample values that control system-wide firmware upgrade behavior.

The example output can be saved to a file and modified to update your global
firmware configuration settings.

No flags or arguments are required for this command.

Examples:
  # Show example global configuration
  metalcloud-cli firmware-policy global-config config-example
  
  # Save example to file for editing
  metalcloud-cli fw-policy global config-example > global-config.json
  
  # Use example as template for updating global config
  metalcloud-cli firmware-policy global-config config-example > config.json
  # Edit config.json with your values
  metalcloud-cli firmware-policy global-config update --config-source config.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return firmware_policy.GlobalFirmwareConfigExample(cmd.Context())
		},
	}
)

func init() {
	rootCmd.AddCommand(firmwarePolicyCmd)

	// Basic policy commands
	firmwarePolicyCmd.AddCommand(firmwarePolicyListCmd)
	firmwarePolicyCmd.AddCommand(firmwarePolicyGetCmd)
	firmwarePolicyCmd.AddCommand(firmwarePolicyConfigExampleCmd)
	firmwarePolicyCmd.AddCommand(firmwarePolicyAuditCmd)

	// Policy modification commands
	firmwarePolicyCmd.AddCommand(firmwarePolicyCreateCmd)
	firmwarePolicyCreateCmd.Flags().StringVar(&firmwarePolicyFlags.configSource, "config-source", "", "Source of the new firmware policy configuration. Can be 'pipe' or path to a JSON file.")
	firmwarePolicyCreateCmd.MarkFlagsOneRequired("config-source")

	firmwarePolicyCmd.AddCommand(firmwarePolicyUpdateCmd)
	firmwarePolicyUpdateCmd.Flags().StringVar(&firmwarePolicyFlags.configSource, "config-source", "", "Source of the firmware policy configuration updates. Can be 'pipe' or path to a JSON file.")
	firmwarePolicyUpdateCmd.MarkFlagsOneRequired("config-source")

	firmwarePolicyCmd.AddCommand(firmwarePolicyDeleteCmd)

	// Apply commands
	firmwarePolicyCmd.AddCommand(firmwarePolicyApplyWithGroupsCmd)
	firmwarePolicyCmd.AddCommand(firmwarePolicyApplyWithoutGroupsCmd)

	// Global firmware configuration commands
	firmwarePolicyCmd.AddCommand(globalFirmwareConfigCmd)
	globalFirmwareConfigCmd.AddCommand(globalFirmwareConfigGetCmd)
	globalFirmwareConfigCmd.AddCommand(globalFirmwareConfigExampleCmd)

	globalFirmwareConfigCmd.AddCommand(globalFirmwareConfigUpdateCmd)
	globalFirmwareConfigUpdateCmd.Flags().StringVar(&firmwarePolicyFlags.configSource, "config-source", "", "Source of the global firmware configuration updates. Can be 'pipe' or path to a JSON file.")
	globalFirmwareConfigUpdateCmd.MarkFlagsOneRequired("config-source")
}
