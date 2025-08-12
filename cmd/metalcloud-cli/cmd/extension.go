package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/extension"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	extensionFlags = struct {
		definitionSource string
		filterLabel      []string
		filterName       []string
		filterStatus     []string
		filterKind       []string
		filterPublic     string
	}{}

	extensionCmd = &cobra.Command{
		Use:     "extension [command]",
		Aliases: []string{"ext", "extensions"},
		Short:   "Manage platform extensions for workflows, applications, and actions",
		Long: `Manage platform extensions including workflows, applications, and actions.

Extensions are modular components that extend the platform's functionality. They can be:
- workflows: Define automated sequences of operations
- applications: Provide custom application deployment logic
- actions: Implement specific operational tasks

Extension lifecycle includes draft, active, and archived states. Only published extensions
become active and available for use across the platform.

Available Commands:
  list     List and filter extensions
  get      Retrieve detailed extension information
  create   Create new extension from definition
  update   Modify existing extension properties
  publish  Activate draft extension for platform use
  archive  Deactivate published extension

Examples:
  metalcloud extension list --filter-kind workflow --filter-status active
  metalcloud extension create my-workflow workflow "Custom deployment workflow" --definition-source definition.json
  metalcloud extension update ext123 "Updated Name" "New description"
  metalcloud extension publish ext123`,
	}

	extensionListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List and filter platform extensions",
		Long: `List and filter platform extensions with various criteria.

This command displays all extensions accessible to your account, including workflows,
applications, and actions. Use filters to narrow down the results based on labels,
names, status, kind, and public visibility.

Extension kinds:
- workflow: Automated sequences of operations
- application: Custom application deployment logic  
- action: Specific operational tasks

Extension statuses:
- draft: Extension is being developed, not yet active
- active: Extension is published and available for use
- archived: Extension is deactivated and no longer available

Flags:
  --filter-label strings    Filter by extension labels (can specify multiple)
  --filter-name strings     Filter by extension names (can specify multiple)
  --filter-status strings   Filter by status: draft, active, archived (can specify multiple)
  --filter-kind strings     Filter by kind: workflow, application, action (can specify multiple)
  --filter-public string    Filter by public visibility: true or false

Examples:
  # List all extensions
  metalcloud extension list
  
  # List only workflow extensions
  metalcloud extension list --filter-kind workflow
  
  # List active and draft extensions
  metalcloud extension list --filter-status active --filter-status draft
  
  # List public workflows
  metalcloud extension list --filter-kind workflow --filter-public true
  
  # List extensions with specific labels
  metalcloud extension list --filter-label production --filter-label critical`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EXTENSIONS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return extension.ExtensionList(
				cmd.Context(),
				extensionFlags.filterLabel,
				extensionFlags.filterName,
				extensionFlags.filterStatus,
				extensionFlags.filterKind,
				extensionFlags.filterPublic,
			)
		},
	}

	extensionGetCmd = &cobra.Command{
		Use:     "get extension_id_or_label",
		Aliases: []string{"show"},
		Short:   "Retrieve detailed information about a specific extension",
		Long: `Retrieve detailed information about a specific extension by ID or label.

This command displays comprehensive information about an extension including its
metadata, definition, current status, and configuration. The extension can be
identified by either its unique ID or label.

Arguments:
  extension_id_or_label    The unique ID or label of the extension to retrieve

Examples:
  # Get extension by ID
  metalcloud extension get 12345
  
  # Get extension by label
  metalcloud extension get my-workflow-v1
  
  # Show extension details
  metalcloud extension show production-deployment`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EXTENSIONS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return extension.ExtensionGet(cmd.Context(), args[0])
		},
	}

	extensionCreateCmd = &cobra.Command{
		Use:     "create name kind description",
		Aliases: []string{"new"},
		Short:   "Create new extension from definition",
		Long: `Create a new extension with the specified name, kind, and description.

This command creates a new extension in draft status. The extension definition must be
provided through the --definition-source flag, which accepts either 'pipe' for stdin
input or a path to a JSON file containing the extension definition.

Extension kinds:
- workflow: Automated sequences of operations
- application: Custom application deployment logic
- action: Specific operational tasks

The newly created extension will be in draft status and must be published before
it becomes available for use on the platform.

Arguments:
  name          The name of the extension to create
  kind          The extension type (workflow, application, action)
  description   Description of the extension's purpose and functionality

Required Flags:
  --definition-source string   Source of the extension definition (required)
                              Can be 'pipe' for stdin or path to a JSON file

Examples:
  # Create extension from JSON file
  metalcloud extension create my-workflow workflow "Custom deployment workflow" --definition-source workflow.json
  
  # Create extension from stdin
  cat definition.json | metalcloud extension create my-app application "Custom app logic" --definition-source pipe
  
  # Create action extension
  metalcloud extension create cleanup-action action "Cleanup resources" --definition-source ./actions/cleanup.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EXTENSIONS_WRITE},
		Args:         cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			definition, err := utils.ReadConfigFromPipeOrFile(extensionFlags.definitionSource)
			if err != nil {
				return err
			}

			return extension.ExtensionCreate(cmd.Context(), args[0], args[1], args[2], definition)
		},
	}

	extensionUpdateCmd = &cobra.Command{
		Use:     "update extension_id_or_label [name [description]]",
		Aliases: []string{"edit"},
		Short:   "Modify existing extension properties and definition",
		Long: `Modify existing extension properties including name, description, and definition.

This command allows you to update various properties of an existing extension.
You can update the name, description, and/or the extension definition. All
parameters are optional, allowing you to update only specific properties.

Arguments:
  extension_id_or_label    The unique ID or label of the extension to update
  name                     New name for the extension (optional)
  description              New description for the extension (optional)

Optional Flags:
  --definition-source string   Source of the updated extension definition
                              Can be 'pipe' for stdin or path to a JSON file

Flag Dependencies:
- --definition-source is independent of other parameters
- name and description are positional arguments

Examples:
  # Update only the name
  metalcloud extension update ext123 "New Extension Name"
  
  # Update name and description
  metalcloud extension update ext123 "New Name" "Updated description"
  
  # Update only the definition
  metalcloud extension update ext123 --definition-source updated-definition.json
  
  # Update name, description, and definition
  metalcloud extension update ext123 "New Name" "New description" --definition-source definition.json
  
  # Update definition from stdin
  cat new-definition.json | metalcloud extension update ext123 --definition-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EXTENSIONS_WRITE},
		Args:         cobra.RangeArgs(1, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := ""
			if len(args) > 1 {
				name = args[1]
			}

			description := ""
			if len(args) > 2 {
				description = args[2]
			}

			var definition []byte
			var err error
			if extensionFlags.definitionSource != "" {
				definition, err = utils.ReadConfigFromPipeOrFile(extensionFlags.definitionSource)
				if err != nil {
					return err
				}
			}

			return extension.ExtensionUpdate(cmd.Context(), args[0], name, description, definition)
		},
	}

	extensionPublishCmd = &cobra.Command{
		Use:   "publish extension_id_or_label",
		Short: "Activate draft extension for platform use",
		Long: `Activate a draft extension making it available for use across the platform.

This command publishes a draft extension, changing its status from draft to active.
Only published extensions are available for use in workflows, applications, and
actions. Once published, an extension cannot be modified directly - you must
create a new version or archive and recreate it.

Publishing validates the extension definition and ensures it meets all platform
requirements before making it available to users.

Arguments:
  extension_id_or_label    The unique ID or label of the draft extension to publish

Requirements:
- Extension must be in draft status
- Extension definition must be valid
- User must have write permissions for extensions

Examples:
  # Publish extension by ID
  metalcloud extension publish 12345
  
  # Publish extension by label
  metalcloud extension publish my-workflow-v1`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EXTENSIONS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return extension.ExtensionPublish(cmd.Context(), args[0])
		},
	}

	extensionArchiveCmd = &cobra.Command{
		Use:   "archive extension_id_or_label",
		Short: "Deactivate published extension and make it unavailable",
		Long: `Deactivate a published extension making it unavailable for use across the platform.

This command archives an active extension, changing its status from active to archived.
Archived extensions are no longer available for use in workflows, applications, and
actions, but their definitions and history are preserved.

Archiving is useful when you want to retire an extension without permanently deleting
it. Archived extensions can be viewed but cannot be used or modified. To reactivate
an archived extension, you must create a new version.

Arguments:
  extension_id_or_label    The unique ID or label of the active extension to archive

Requirements:
- Extension must be in active status
- User must have write permissions for extensions
- Extension should not be actively used in critical workflows

Examples:
  # Archive extension by ID
  metalcloud extension archive 12345
  
  # Archive extension by label
  metalcloud extension archive deprecated-workflow-v1`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EXTENSIONS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return extension.ExtensionArchive(cmd.Context(), args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(extensionCmd)

	extensionCmd.AddCommand(extensionListCmd)
	extensionListCmd.Flags().StringSliceVar(&extensionFlags.filterLabel, "filter-label", nil, "Filter extensions by label")
	extensionListCmd.Flags().StringSliceVar(&extensionFlags.filterName, "filter-name", nil, "Filter extensions by name")
	extensionListCmd.Flags().StringSliceVar(&extensionFlags.filterStatus, "filter-status", nil, "Filter extensions by status (draft, active, archived)")
	extensionListCmd.Flags().StringSliceVar(&extensionFlags.filterKind, "filter-kind", nil, "Filter extensions by kind (application, workflow, action)")
	extensionListCmd.Flags().StringVar(&extensionFlags.filterPublic, "filter-public", "", "Filter extensions by public status (true/false)")

	extensionCmd.AddCommand(extensionGetCmd)

	extensionCmd.AddCommand(extensionCreateCmd)
	extensionCreateCmd.Flags().StringVar(&extensionFlags.definitionSource, "definition-source", "", "Source of the extension definition. Can be 'pipe' or path to a JSON file.")
	extensionCreateCmd.MarkFlagRequired("definition-source")

	extensionCmd.AddCommand(extensionUpdateCmd)
	extensionUpdateCmd.Flags().StringVar(&extensionFlags.definitionSource, "definition-source", "", "Source of the updated extension definition. Can be 'pipe' or path to a JSON file.")

	extensionCmd.AddCommand(extensionPublishCmd)
	extensionCmd.AddCommand(extensionArchiveCmd)
}
