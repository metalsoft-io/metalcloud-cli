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
		Short:   "Firmware policy management",
		Long:    `Firmware policy management commands.`,
	}

	firmwarePolicyListCmd = &cobra.Command{
		Use:          "list",
		Aliases:      []string{"ls"},
		Short:        "List all firmware policies",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return firmware_policy.FirmwarePolicyList(cmd.Context())
		},
	}

	firmwarePolicyGetCmd = &cobra.Command{
		Use:          "get policy_id",
		Aliases:      []string{"show"},
		Short:        "Get firmware policy details",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return firmware_policy.FirmwarePolicyGet(cmd.Context(), args[0])
		},
	}

	firmwarePolicyCreateCmd = &cobra.Command{
		Use:          "create",
		Aliases:      []string{"new"},
		Short:        "Create a new firmware policy",
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
		Use:          "update policy_id",
		Aliases:      []string{"edit"},
		Short:        "Update a firmware policy",
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
		Use:          "delete policy_id",
		Aliases:      []string{"rm"},
		Short:        "Delete a firmware policy",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return firmware_policy.FirmwarePolicyDelete(cmd.Context(), args[0])
		},
	}

	firmwarePolicyAuditCmd = &cobra.Command{
		Use:          "generate-audit policy_id",
		Aliases:      []string{"audit"},
		Short:        "Generate audit for a firmware policy",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return firmware_policy.FirmwarePolicyGenerateAudit(cmd.Context(), args[0])
		},
	}

	firmwarePolicyApplyWithGroupsCmd = &cobra.Command{
		Use:          "apply-with-groups",
		Short:        "Apply all firmware policies linked to server instance groups",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			return firmware_policy.FirmwarePolicyApplyWithGroups(cmd.Context())
		},
	}

	firmwarePolicyApplyWithoutGroupsCmd = &cobra.Command{
		Use:          "apply-without-groups",
		Short:        "Apply all firmware policies not linked to server instance groups",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			return firmware_policy.FirmwarePolicyApplyWithoutGroups(cmd.Context())
		},
	}

	firmwarePolicyConfigExampleCmd = &cobra.Command{
		Use:          "config-example",
		Short:        "Show firmware policy configuration example",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return firmware_policy.FirmwarePolicyConfigExample(cmd.Context())
		},
	}

	// Global firmware configuration commands
	globalFirmwareConfigCmd = &cobra.Command{
		Use:          "global-config",
		Aliases:      []string{"global"},
		Short:        "Manage global firmware configuration",
		SilenceUsage: true,
	}

	globalFirmwareConfigGetCmd = &cobra.Command{
		Use:          "get",
		Short:        "Get global firmware configuration",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return firmware_policy.GetGlobalFirmwareConfiguration(cmd.Context())
		},
	}

	globalFirmwareConfigUpdateCmd = &cobra.Command{
		Use:          "update",
		Aliases:      []string{"edit"},
		Short:        "Update global firmware configuration",
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
		Use:          "config-example",
		Short:        "Show global firmware configuration example",
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
