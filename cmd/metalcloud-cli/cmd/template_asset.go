package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/template_asset"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	templateAssetFlags = struct {
		configSource string
		templateId   []string
		usage        []string
		mimeType     []string
	}{}

	templateAssetCmd = &cobra.Command{
		Use:     "template-asset [command]",
		Aliases: []string{"template-assets", "assets"},
		Short:   "Manage template assets for OS templates",
		Long: `Manage template assets associated with OS templates. Template assets are files that are 
used during the deployment process of operating systems on different device types.

Template assets have different usages depending on the device type:

SERVER device type supports:
  - build_source_image: Source image for building (requires file.url)
  - build_component: Component to be copied during build (uses file.path for ISO location)
  - secondary_image: Secondary image mounted in 2nd virtual media (requires file.url)

NETWORK_DEVICE type supports:
  - source_image: NOS image for ONIE install process (requires file.url, uses file.path for HTTP URL)
  - switch_ztp_config: Configuration for ZTP process (uses file.path for HTTP URL)
  - generic: General purpose asset for ZTP process (uses file.path for HTTP URL)

VM device type supports:
  - metadata_source_image: VM image metadata (uses file.path to identify metadata)
  - generic: General purpose asset (uses file.path to identify role)

Available commands: list, get, config-example, create, update, delete`,
	}

	templateAssetListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List template assets with optional filtering",
		Long: `List all template assets or filter by specific criteria.

Filtering Options:
  --template-id: Filter assets by one or more template IDs
  --usage: Filter assets by usage type (e.g., build_source_image, secondary_image, source_image, etc.)
  --mime-type: Filter assets by file MIME type (e.g., image/png, image/jpeg, application/octet-stream)

Examples:
  # List all template assets
  metalcloud-cli template-asset list

  # List assets for specific template
  metalcloud-cli template-asset list --template-id 123

  # List assets with specific usage
  metalcloud-cli template-asset list --usage build_source_image

  # List image assets only
  metalcloud-cli template-asset list --mime-type image/png,image/jpeg

  # Combine multiple filters
  metalcloud-cli template-asset list --template-id 123,456 --usage secondary_image`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_TEMPLATES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return template_asset.TemplateAssetList(
				cmd.Context(),
				templateAssetFlags.templateId,
				templateAssetFlags.usage,
				templateAssetFlags.mimeType)
		},
	}

	templateAssetGetCmd = &cobra.Command{
		Use:     "get template_asset_id",
		Aliases: []string{"show"},
		Short:   "Get detailed information about a specific template asset",
		Long: `Retrieve and display detailed information about a specific template asset by its ID.

Arguments:
  template_asset_id: The unique identifier of the template asset (required)

Examples:
  # Get details for template asset with ID 123
  metalcloud-cli template-asset get 123

  # Using alias
  metalcloud-cli template-asset show 456`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_TEMPLATES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return template_asset.TemplateAssetGet(cmd.Context(), args[0])
		},
	}

	templateAssetConfigExampleCmd = &cobra.Command{
		Use:   "config-example",
		Short: "Display template asset configuration example",
		Long: `Display a template asset configuration example that can be used as a template 
for creating or updating template assets.

The example shows the JSON structure with all available fields:
  - templateId: ID of the OS template that this asset belongs to (required)
  - usage: Usage type of the asset (required, e.g., logo, icon, build_source_image)
  - file: File configuration object (required)
    - name: Name of the file (required)
    - mimeType: MIME type of the file (required, e.g., image/png, application/octet-stream)
    - checksum: File checksum for integrity verification (optional)
    - contentBase64: Base64-encoded file content (optional, for direct upload)
    - templatingEngine: Whether templating is enabled for the file (optional, default: false)
    - url: External URL where the file can be downloaded (optional)
  - tags: Array of tags for categorization (optional)

Examples:
  # Display configuration example
  metalcloud-cli template-asset config-example

  # Save example to file for editing
  metalcloud-cli template-asset config-example > asset-config.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_TEMPLATES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return template_asset.TemplateAssetConfigExample(cmd.Context())
		},
	}

	templateAssetCreateCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Short:   "Create a new template asset from configuration",
		Long: `Create a new template asset using a configuration file or piped input.

Required Flags:
  --config-source: Source of the template asset configuration (required)
                   Can be 'pipe' for piped input or path to a JSON file

Configuration Format:
The configuration must be in JSON format with the following structure:
  {
    "templateId": 123,                              // Required: ID of the OS template
    "usage": "logo",                                // Required: Usage type (logo, icon, build_source_image, etc.)
    "file": {                                       // Required: File configuration
      "name": "example.png",                        // Required: File name
      "mimeType": "image/png",                      // Required: MIME type
      "checksum": "sha256-hash",                    // Optional: File checksum
      "contentBase64": "base64-content",            // Optional: Base64 encoded file content
      "templatingEngine": false,                    // Optional: Enable templating (default: false)
      "url": "https://example.com/file.png"        // Optional: External URL
    },
    "tags": ["branding", "image"]                   // Optional: Tags for categorization
  }

Examples:
  # Create from JSON file
  metalcloud-cli template-asset create --config-source asset-config.json

  # Create from pipe
  echo '{"templateId": 123, "usage": "logo", "file": {"name": "logo.png", "mimeType": "image/png"}}' | metalcloud-cli template-asset create --config-source pipe

  # Generate example and create
  metalcloud-cli template-asset config-example > config.json
  # Edit config.json with your values
  metalcloud-cli template-asset create --config-source config.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_TEMPLATES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(templateAssetFlags.configSource)
			if err != nil {
				return err
			}

			return template_asset.TemplateAssetCreate(cmd.Context(), config)
		},
	}

	templateAssetUpdateCmd = &cobra.Command{
		Use:     "update template_asset_id",
		Aliases: []string{"edit"},
		Short:   "Update an existing template asset with new configuration",
		Long: `Update an existing template asset using a configuration file or piped input.

Arguments:
  template_asset_id: The unique identifier of the template asset to update (required)

Required Flags:
  --config-source: Source of the template asset configuration updates (required)
                   Can be 'pipe' for piped input or path to a JSON file

Configuration Format:
The configuration must be in JSON format with the same structure as create:
  {
    "templateId": 123,                              // Required: ID of the OS template
    "usage": "logo",                                // Required: Usage type (logo, icon, build_source_image, etc.)
    "file": {                                       // Required: File configuration
      "name": "updated-example.png",                // Required: File name
      "mimeType": "image/png",                      // Required: MIME type
      "checksum": "new-sha256-hash",                // Optional: File checksum
      "contentBase64": "updated-base64-content",    // Optional: Base64 encoded file content
      "templatingEngine": false,                    // Optional: Enable templating (default: false)
      "url": "https://example.com/updated-file.png" // Optional: External URL
    },
    "tags": ["branding", "updated"]                 // Optional: Tags for categorization
  }

Examples:
  # Update from JSON file
  metalcloud-cli template-asset update 123 --config-source updated-config.json

  # Update from pipe
  echo '{"templateId": 123, "usage": "icon", "file": {"name": "new-icon.png", "mimeType": "image/png"}}' | metalcloud-cli template-asset update 456 --config-source pipe

  # Generate example, edit, and update
  metalcloud-cli template-asset config-example > update-config.json
  # Edit update-config.json with your new values
  metalcloud-cli template-asset update 789 --config-source update-config.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_TEMPLATES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(templateAssetFlags.configSource)
			if err != nil {
				return err
			}

			return template_asset.TemplateAssetUpdate(cmd.Context(), args[0], config)
		},
	}

	templateAssetDeleteCmd = &cobra.Command{
		Use:     "delete template_asset_id",
		Aliases: []string{"rm"},
		Short:   "Delete a template asset",
		Long: `Delete an existing template asset by its ID.

Arguments:
  template_asset_id: The unique identifier of the template asset to delete (required)

Examples:
  # Delete template asset with ID 123
  metalcloud-cli template-asset delete 123

  # Using alias
  metalcloud-cli template-asset rm 456

Warning:
This operation is irreversible. The template asset and all its associated data will be permanently removed.`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_TEMPLATES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return template_asset.TemplateAssetDelete(cmd.Context(), args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(templateAssetCmd)

	// List command with filter options
	templateAssetCmd.AddCommand(templateAssetListCmd)
	templateAssetListCmd.Flags().StringSliceVar(&templateAssetFlags.templateId, "template-id", nil, "Filter assets by template ID.")
	templateAssetListCmd.Flags().StringSliceVar(&templateAssetFlags.usage, "usage", nil, "Filter assets by usage type (e.g., logo, icon, etc.).")
	templateAssetListCmd.Flags().StringSliceVar(&templateAssetFlags.mimeType, "mime-type", nil, "Filter assets by file MIME type (e.g., image/png, image/jpeg, etc.).")

	// Get command
	templateAssetCmd.AddCommand(templateAssetGetCmd)

	// Config example command
	templateAssetCmd.AddCommand(templateAssetConfigExampleCmd)

	// Create command
	templateAssetCmd.AddCommand(templateAssetCreateCmd)
	templateAssetCreateCmd.Flags().StringVar(&templateAssetFlags.configSource, "config-source", "", "Source of the new template asset configuration. Can be 'pipe' or path to a JSON file.")
	templateAssetCreateCmd.MarkFlagsOneRequired("config-source")

	// Update command
	templateAssetCmd.AddCommand(templateAssetUpdateCmd)
	templateAssetUpdateCmd.Flags().StringVar(&templateAssetFlags.configSource, "config-source", "", "Source of the template asset configuration updates. Can be 'pipe' or path to a JSON file.")
	templateAssetUpdateCmd.MarkFlagsOneRequired("config-source")

	// Delete command
	templateAssetCmd.AddCommand(templateAssetDeleteCmd)
}
