package cmd

import (
	"fmt"
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/extension_instance"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
	"github.com/spf13/cobra"
)

var (
	extensionInstanceFlags = struct {
		configSource   string
		extensionId    int
		label          string
		inputVariables []string
	}{}

	extensionInstanceCmd = &cobra.Command{
		Use:     "extension-instance [command]",
		Aliases: []string{"ext-inst"},
		Short:   "Extension instance management",
		Long:    `Extension instance management commands.`,
	}

	extensionInstanceListCmd = &cobra.Command{
		Use:          "list infrastructure_id_or_label",
		Aliases:      []string{"ls"},
		Short:        "List all extension instances in an infrastructure.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EXTENSION_INSTANCES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return extension_instance.ExtensionInstanceList(cmd.Context(), args[0])
		},
	}

	extensionInstanceGetCmd = &cobra.Command{
		Use:          "get extension_instance_id",
		Aliases:      []string{"show"},
		Short:        "Get extension instance details.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EXTENSION_INSTANCES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return extension_instance.ExtensionInstanceGet(cmd.Context(), args[0])
		},
	}

	extensionInstanceCreateCmd = &cobra.Command{
		Use:     "create infrastructure_id_or_label",
		Aliases: []string{"new"},
		Short:   "Create a new extension instance in the specified infrastructure",
		Long: `Create a new extension instance in the specified infrastructure.

Extension instances are deployments of extensions within an infrastructure. They can be created
using either a configuration file/pipe or by specifying individual parameters via flags.

Configuration methods:
  1. Using --config-source: Provide a JSON configuration file or pipe the configuration
  2. Using individual flags: Specify --extension-id and optionally --label and --input-variable

Examples:
  # Create from JSON file
  metalcloud-cli extension-instance create my-infra --config-source ./config.json

  # Create from pipe
  echo '{"extensionId": 123, "label": "my-instance"}' | metalcloud-cli extension-instance create my-infra --config-source pipe

  # Create using individual flags
  metalcloud-cli extension-instance create my-infra --extension-id 123 --label "my-instance"

  # Create with input variables
  metalcloud-cli extension-instance create my-infra --extension-id 123 --input-variable "env=production" --input-variable "replicas=3"

JSON Configuration Format:
  {
    "extensionId": 123,
    "label": "optional-instance-label",
    "inputVariables": [
      {"label": "variable1", "value": "value1"},
      {"label": "variable2", "value": "value2"}
    ]
  }`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EXTENSION_INSTANCES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var payload sdk.CreateExtensionInstance

			// Check if config-source is provided
			if extensionInstanceFlags.configSource != "" {
				config, err := utils.ReadConfigFromPipeOrFile(extensionInstanceFlags.configSource)
				if err != nil {
					return err
				}

				if err := utils.UnmarshalContent(config, &payload); err != nil {
					return fmt.Errorf("invalid config: %w", err)
				}
			} else {
				// If config-source is not provided, use individual flags
				if extensionInstanceFlags.extensionId == 0 {
					return fmt.Errorf("extension-id is required when config-source is not provided")
				}

				payload.ExtensionId = sdk.PtrFloat32(float32(extensionInstanceFlags.extensionId))
				if extensionInstanceFlags.label != "" {
					payload.Label = &extensionInstanceFlags.label
				}
				payload.InputVariables = make([]sdk.ExtensionVariable, len(extensionInstanceFlags.inputVariables))
				for i, inputVar := range extensionInstanceFlags.inputVariables {
					parts := strings.Split(inputVar, "=")
					if len(parts) != 2 {
						return fmt.Errorf("invalid input variable format: %s, expected 'label=value'", inputVar)
					}

					payload.InputVariables[i] = sdk.ExtensionVariable{
						Label: parts[0],
						Value: parts[1],
					}
				}
			}

			return extension_instance.ExtensionInstanceCreate(cmd.Context(), args[0], payload)

		},
	}

	extensionInstanceUpdateCmd = &cobra.Command{
		Use:          "update extension_instance_id",
		Aliases:      []string{"edit"},
		Short:        "Update extension instance configuration.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EXTENSION_INSTANCES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(extensionInstanceFlags.configSource)
			if err != nil {
				return err
			}
			return extension_instance.ExtensionInstanceUpdate(cmd.Context(), args[0], config)
		},
	}

	extensionInstanceDeleteCmd = &cobra.Command{
		Use:          "delete extension_instance_id",
		Aliases:      []string{"rm"},
		Short:        "Delete extension instance.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EXTENSION_INSTANCES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return extension_instance.ExtensionInstanceDelete(cmd.Context(), args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(extensionInstanceCmd)

	extensionInstanceCmd.AddCommand(extensionInstanceListCmd)
	extensionInstanceCmd.AddCommand(extensionInstanceGetCmd)

	extensionInstanceCmd.AddCommand(extensionInstanceCreateCmd)
	extensionInstanceCreateCmd.Flags().StringVar(&extensionInstanceFlags.configSource, "config-source", "", "Source of the new extension instance configuration. Can be 'pipe' or path to a JSON file.")
	extensionInstanceCreateCmd.Flags().IntVar(&extensionInstanceFlags.extensionId, "extension-id", 0, "The extension ID to create an instance of.")
	extensionInstanceCreateCmd.Flags().StringVar(&extensionInstanceFlags.label, "label", "", "The extension instance label (optional, will be auto-generated if not provided).")
	extensionInstanceCreateCmd.Flags().StringArrayVar(&extensionInstanceFlags.inputVariables, "input-variable", []string{}, "Input variables in format 'label=value'. Can be specified multiple times.")
	extensionInstanceCreateCmd.MarkFlagsMutuallyExclusive("config-source", "extension-id")
	extensionInstanceCreateCmd.MarkFlagsOneRequired("config-source", "extension-id")

	extensionInstanceCmd.AddCommand(extensionInstanceUpdateCmd)
	extensionInstanceUpdateCmd.Flags().StringVar(&extensionInstanceFlags.configSource, "config-source", "", "Source of the extension instance configuration updates. Can be 'pipe' or path to a JSON file.")
	extensionInstanceUpdateCmd.MarkFlagsOneRequired("config-source")

	extensionInstanceCmd.AddCommand(extensionInstanceDeleteCmd)
}
