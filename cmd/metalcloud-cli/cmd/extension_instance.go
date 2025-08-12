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
		Short:   "Manage extension instances within infrastructure deployments",
		Long: `Manage extension instances within infrastructure deployments.

Extension instances are concrete deployments of application extensions within a specific infrastructure.
They represent the running or configured state of an extension with specific input variables and configurations.

Each extension instance is tied to an infrastructure and can be configured with custom input
variables that define its behavior. Extension instances go through various lifecycle states
including deployment, running, and deletion.

Available Commands:
  list     List all extension instances in an infrastructure
  get      Retrieve detailed extension instance information
  create   Deploy new extension instance in infrastructure
  update   Modify existing extension instance configuration
  delete   Remove extension instance from infrastructure

Examples:
  metalcloud extension-instance list my-infrastructure
  metalcloud extension-instance create my-infra --extension-id 123 --label "web-server"
  metalcloud extension-instance update inst456 --config-source updated-config.json
  metalcloud extension-instance delete inst456`,
	}

	extensionInstanceListCmd = &cobra.Command{
		Use:     "list infrastructure_id_or_label",
		Aliases: []string{"ls"},
		Short:   "List all extension instances in an infrastructure",
		Long: `List all extension instances deployed within a specific infrastructure.

This command displays all extension instances that are currently deployed or configured
within the specified infrastructure. Extension instances represent active deployments
of extensions (workflows, applications, or actions) with their current status,
configuration, and input variables.

The output includes instance details such as:
- Instance ID and label
- Associated extension information
- Current status and state
- Input variables and configuration
- Deployment timestamps

Arguments:
  infrastructure_id_or_label    The unique ID or label of the infrastructure

Examples:
  # List extension instances by infrastructure ID
  metalcloud extension-instance list 12345
  
  # List extension instances by infrastructure label
  metalcloud extension-instance list production-infrastructure
  
  # List instances in staging environment
  metalcloud extension-instance ls staging-env`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EXTENSION_INSTANCES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return extension_instance.ExtensionInstanceList(cmd.Context(), args[0])
		},
	}

	extensionInstanceGetCmd = &cobra.Command{
		Use:     "get extension_instance_id",
		Aliases: []string{"show"},
		Short:   "Retrieve detailed information about a specific extension instance",
		Long: `Retrieve detailed information about a specific extension instance by ID.

This command displays comprehensive information about an extension instance including
its configuration, current status, associated extension details, input variables,
and deployment history. The output provides insights into the instance's current
state and operational status within its infrastructure.

Information displayed includes:
- Instance metadata and identifiers
- Associated extension details
- Current operational status
- Input variables and their values
- Configuration parameters
- Deployment and update timestamps
- Infrastructure context

Arguments:
  extension_instance_id    The unique ID of the extension instance to retrieve

Examples:
  # Get extension instance details by ID
  metalcloud extension-instance get 12345
  
  # Show instance information with alias
  metalcloud extension-instance show 67890
  
  # Get instance details for troubleshooting
  metalcloud ext-inst get instance-123`,
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
		Short:   "Deploy new extension instance in specified infrastructure",
		Long: `Deploy a new extension instance in the specified infrastructure.

This command creates and deploys a new extension instance within an infrastructure.
Extension instances are concrete deployments of extensions (workflows, applications,
or actions) with specific configurations and input variables.

You can provide the configuration using two methods:

Method 1: Configuration file or pipe (--config-source)
  Use a JSON file or pipe the configuration via stdin. This method allows for
  complete configuration including complex input variables.

Method 2: Individual flags (--extension-id with optional flags)
  Specify the extension ID directly with optional label and input variables.
  This method is suitable for simple configurations.

Arguments:
  infrastructure_id_or_label    The unique ID or label of the target infrastructure

Required Flags (mutually exclusive):
  --config-source string        Source of configuration (pipe or JSON file path)
  --extension-id int           ID of the extension to instantiate

Optional Flags (only with --extension-id):
  --label string               Custom label for the extension instance
  --input-variable strings     Input variables in 'label=value' format (repeatable)

Flag Dependencies:
- --config-source and --extension-id are mutually exclusive
- One of --config-source or --extension-id is required
- --label and --input-variable only work with --extension-id

JSON Configuration Format:
  {
    "extensionId": 123,
    "label": "optional-instance-label",
    "inputVariables": [
      {"label": "variable1", "value": "value1"},
      {"label": "variable2", "value": "value2"}
    ]
  }

Examples:
  # Create from JSON file
  metalcloud extension-instance create my-infra --config-source ./config.json

  # Create from pipe
  echo '{"extensionId": 123, "label": "web-app"}' | metalcloud extension-instance create my-infra --config-source pipe

  # Create using individual flags
  metalcloud extension-instance create my-infra --extension-id 123 --label "database-server"

  # Create with input variables
  metalcloud extension-instance create my-infra --extension-id 123 --input-variable "env=production" --input-variable "replicas=3"

  # Create minimal instance (auto-generated label)
  metalcloud ext-inst create prod-infra --extension-id 456`,
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
		Use:     "update extension_instance_id",
		Aliases: []string{"edit"},
		Short:   "Modify existing extension instance configuration",
		Long: `Modify existing extension instance configuration with updated parameters.

This command allows you to update the configuration of an existing extension instance.
The updated configuration must be provided through the --config-source flag, which
accepts either 'pipe' for stdin input or a path to a JSON file containing the
updated configuration.

Use this command to modify input variables, change labels, or update other
configurable parameters of a deployed extension instance. The instance will
be reconfigured with the new settings while maintaining its deployment state.

Arguments:
  extension_instance_id    The unique ID of the extension instance to update

Required Flags:
  --config-source string   Source of the updated configuration (required)
                          Can be 'pipe' for stdin or path to a JSON file

JSON Configuration Format:
  {
    "label": "updated-instance-label",
    "inputVariables": [
      {"label": "variable1", "value": "new-value1"},
      {"label": "variable2", "value": "new-value2"}
    ]
  }

Examples:
  # Update from JSON file
  metalcloud extension-instance update 12345 --config-source ./updated-config.json
  
  # Update from pipe
  echo '{"label": "new-label"}' | metalcloud extension-instance update 12345 --config-source pipe
  
  # Update input variables
  metalcloud extension-instance update 12345 --config-source ./new-variables.json
  
  # Edit with alias
  metalcloud ext-inst edit 67890 --config-source updated-config.json`,
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
		Use:     "delete extension_instance_id",
		Aliases: []string{"rm"},
		Short:   "Remove extension instance from infrastructure",
		Long: `Remove an extension instance from its infrastructure deployment.

This command permanently deletes an extension instance from the specified infrastructure.
The extension instance will be stopped, undeployed, and removed from the infrastructure's
configuration. This action cannot be undone.

Before deletion, ensure that:
- The extension instance is not critical to ongoing operations
- Any dependent services or workflows are properly handled
- You have appropriate backups or documentation if needed

The deletion process will:
1. Stop the running extension instance
2. Clean up associated resources
3. Remove the instance from infrastructure configuration
4. Update the infrastructure state

Arguments:
  extension_instance_id    The unique ID of the extension instance to delete

Examples:
  # Delete extension instance by ID
  metalcloud extension-instance delete 12345
  
  # Remove instance using alias
  metalcloud extension-instance rm 67890
  
  # Delete with short command alias
  metalcloud ext-inst delete instance-123
  
  # Remove specific instance
  metalcloud ext-inst rm failed-deployment-456`,
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
