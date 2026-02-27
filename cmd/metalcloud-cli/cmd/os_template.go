package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/os_template"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	osTemplateFlags = struct {
		configSource string
		deviceType   string
		visibility   string
		repoUrl      string
		repoUsername string
		repoPassword string
		name         string
		label        string
		sourceIso    string
		status       string
		outputPath   string
	}{}

	osTemplateCmd = &cobra.Command{
		Use:     "os-template [command]",
		Aliases: []string{"templates"},
		Short:   "Manage OS templates for server deployments",
		Long: `Manage OS templates used for server deployments.

OS templates define the operating system, installation method, device configuration,
and associated assets needed to deploy operating systems on servers. Templates can
be created from scratch, cloned from repositories, or imported from external sources.

Available commands:
  list                List all available OS templates
  get                 Show detailed information about a specific template
  create              Create a new OS template from JSON configuration
  update              Update an existing OS template
  delete              Delete an OS template
  get-credentials     Show default credentials for a template
  get-assets          List all assets associated with a template
  list-repo           List templates available in a remote repository
  create-from-repo    Create a template by cloning from a repository
  clone               Clone an existing template
  export              Export a template and its assets to a zip archive
  import              Import a template from a zip archive
  example-create      Show example JSON for creating templates`,
	}

	osTemplateListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all available OS templates",
		Long: `List all available OS templates in the system.

This command displays a table of all OS templates with their basic information
including ID, name, label, device type, status, visibility, and timestamps.

The output includes:
  - Template ID (unique identifier)
  - Name (human-readable template name)
  - Label (unique template label)
  - Device Type (server, switch, etc.)
  - Status (ready, active, used, archived)
  - Visibility (public, private)
  - Created/Modified timestamps

Examples:
  # List all OS templates
  metalcloud-cli os-template list
  
  # List templates using alias
  metalcloud-cli templates ls`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_TEMPLATES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return os_template.OsTemplateList(cmd.Context())
		},
	}

	osTemplateGetCmd = &cobra.Command{
		Use:     "get <os_template_id>",
		Aliases: []string{"show"},
		Short:   "Show detailed information about a specific OS template",
		Long: `Display comprehensive details about a specific OS template.

This command retrieves and displays detailed information about an OS template
including its configuration, device settings, installation parameters, OS details,
and metadata.

The template can be identified by its numeric ID.

Examples:
  # Get details for template with ID 123
  metalcloud-cli os-template get 123
  
  # Show template details using alias
  metalcloud-cli templates show 456`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_TEMPLATES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return os_template.OsTemplateGet(cmd.Context(), args[0])
		},
	}

	osTemplateExampleCreateCmd = &cobra.Command{
		Use:     "example-create",
		Aliases: []string{"example"},
		Short:   "Show example JSON for creating OS templates",
		Long: `Display example JSON configuration for creating OS templates.

This command outputs a complete example JSON structure showing all available
fields and their expected values for creating OS templates. The example includes
both template configuration and associated assets.

The output can be used as a starting point for creating custom templates by
modifying the values to match your requirements.

Examples:
  # Show example JSON
  metalcloud-cli os-template example-create
  
  # Save example to file for editing
  metalcloud-cli os-template example-create > template.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_TEMPLATES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			return os_template.OsTemplateExampleCreate(cmd.Context())
		},
	}

	osTemplateCreateCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Short:   "Create a new OS template from JSON configuration",
		Long: `Create a new OS template from JSON configuration.

This command creates a new OS template by reading configuration from a file
or from standard input. The configuration should include both the template
definition and any associated assets.

Required flags:
  --config-source   Source of the template configuration (required)
                   Can be 'pipe' to read from stdin, or path to a JSON file

The JSON configuration should follow the structure shown by the 'example-create'
command, including template definition and optional template assets.

Examples:
  # Create template from file
  metalcloud-cli os-template create --config-source template.json
  
  # Create template from stdin
  cat template.json | metalcloud-cli os-template create --config-source pipe
  
  # Generate example and create template
  metalcloud-cli os-template example-create > template.json
  # Edit template.json with your values
  metalcloud-cli os-template create --config-source template.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_TEMPLATES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(osTemplateFlags.configSource)
			if err != nil {
				return err
			}

			var osTemplateCreateOptions os_template.OsTemplateCreateOptions
			err = utils.UnmarshalContent(config, &osTemplateCreateOptions)
			if err != nil {
				return err
			}

			return os_template.OsTemplateCreate(cmd.Context(), osTemplateCreateOptions)
		},
	}

	osTemplateUpdateCmd = &cobra.Command{
		Use:     "update <os_template_id>",
		Aliases: []string{"edit"},
		Short:   "Update an existing OS template",
		Long: `Update an existing OS template with new configuration.

This command updates an OS template by reading the updated configuration
from a file or from standard input. You can update the template properties
as well as add, update, or delete template assets.

Required arguments:
  os_template_id    The numeric ID of the template to update

Required flags:
  --config-source   Source of the template update configuration (required)
                   Can be 'pipe' to read from stdin, or path to a JSON file

The JSON configuration should include:
  - template: OS template update data (optional)
  - newTemplateAssets: Array of new assets to add (optional)
  - updatedTemplateAssets: Map of asset ID to updated asset data (optional)
  - deletedTemplateAssetIds: Array of asset IDs to delete (optional)

Examples:
  # Update template from file
  metalcloud-cli os-template update 123 --config-source update.json
  
  # Update template from stdin
  cat update.json | metalcloud-cli os-template update 123 --config-source pipe
  
  # Update only template properties (no assets)
  echo '{"template":{"name":"New Name"}}' | metalcloud-cli os-template update 123 --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_TEMPLATES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(osTemplateFlags.configSource)
			if err != nil {
				return err
			}

			var osTemplateUpdateOptions os_template.OsTemplateUpdateOptions
			err = utils.UnmarshalContent(config, &osTemplateUpdateOptions)
			if err != nil {
				return err
			}

			return os_template.OsTemplateUpdate(cmd.Context(), args[0], osTemplateUpdateOptions)
		},
	}

	osTemplateSetStatusCmd = &cobra.Command{
		Use:     "set-status <os_template_id> <status>",
		Aliases: []string{"status"},
		Short:   "Set the status of an OS template",
		Long: `Set the status of an OS template.

This command updates the status of an existing OS template. Valid status values
include: ready, active, used, archived.

Required arguments:
  os_template_id    The numeric ID of the template to update
  status           The new status value (ready, active, used, archived)

Examples:
  # Set template status to active
  metalcloud-cli os-template set-status 123 active
  
  # Archive a template
  metalcloud-cli os-template set-status 456 archived
  
  # Set template to ready status using alias
  metalcloud-cli templates status 789 ready`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_TEMPLATES_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return os_template.OsTemplateSetStatus(cmd.Context(), args[0], args[1])
		},
	}

	osTemplateDeleteCmd = &cobra.Command{
		Use:     "delete <os_template_id>",
		Aliases: []string{"rm"},
		Short:   "Delete an OS template",
		Long: `Delete an OS template from the system.

This command permanently removes an OS template from the system. The template
must not be in use by any active deployments before it can be deleted.

Required arguments:
  os_template_id    The numeric ID of the template to delete

Examples:
  # Delete template with ID 123
  metalcloud-cli os-template delete 123
  
  # Delete template using alias
  metalcloud-cli templates rm 456`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_TEMPLATES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return os_template.OsTemplateDelete(cmd.Context(), args[0])
		},
	}

	osTemplateGetCredentialsCmd = &cobra.Command{
		Use:     "get-credentials <os_template_id>",
		Aliases: []string{"creds"},
		Short:   "Show default credentials for an OS template",
		Long: `Display the default credentials for an OS template.

This command retrieves and displays the default username and password
that are configured for a specific OS template. These credentials are
used for initial access to servers deployed with this template.

Required arguments:
  os_template_id    The numeric ID of the template

Examples:
  # Get credentials for template with ID 123
  metalcloud-cli os-template get-credentials 123
  
  # Get credentials using alias
  metalcloud-cli templates creds 456`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_TEMPLATES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return os_template.OsTemplateGetCredentials(cmd.Context(), args[0])
		},
	}

	osTemplateGetAssetsCmd = &cobra.Command{
		Use:     "get-assets <os_template_id>",
		Aliases: []string{"assets"},
		Short:   "List all assets associated with an OS template",
		Long: `Display all assets associated with a specific OS template.

This command retrieves and displays detailed information about all assets
(files, scripts, configuration files, etc.) that are associated with an
OS template. These assets are used during the deployment process.

Required arguments:
  os_template_id    The numeric ID of the template

The output includes:
  - Asset ID (unique identifier)
  - Usage (build_source_image, build_component, etc.)
  - Filename and MIME type
  - File size
  - Creation timestamp

Examples:
  # List assets for template with ID 123
  metalcloud-cli os-template get-assets 123
  
  # List assets using alias
  metalcloud-cli templates assets 456`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_TEMPLATES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return os_template.OsTemplateGetAssets(cmd.Context(), args[0])
		},
	}

	osTemplateListRepoCmd = &cobra.Command{
		Use:     "list-repo",
		Aliases: []string{"ls-repo"},
		Short:   "List available OS templates from a remote repository",
		Long: `List all available OS templates from a remote repository.

This command retrieves and displays templates available in a remote repository,
showing their basic information and configuration.

Optional flags:
  --repo-url        URL of the repository to list templates from
                   Defaults to the official MetalSoft template repository
  --repo-username   Username for private repository authentication
  --repo-password   Password for private repository authentication

Flag dependencies:
  - If accessing a private repository, both --repo-username and --repo-password
    are required together

Examples:
  # List templates from default public repository
  metalcloud-cli os-template list-repo
  
  # List templates from a custom repository
  metalcloud-cli os-template list-repo --repo-url https://example.com/templates
  
  # List templates from private repository
  metalcloud-cli os-template list-repo --repo-url https://private.com/templates \
    --repo-username user --repo-password pass`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_TEMPLATES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return os_template.OsTemplateListRepo(cmd.Context(), osTemplateFlags.repoUrl, osTemplateFlags.repoUsername, osTemplateFlags.repoPassword)
		},
	}

	osTemplateCloneCmd = &cobra.Command{
		Use:     "clone <os_template_id>",
		Aliases: []string{"copy"},
		Short:   "Clone an existing OS template",
		Long: `Clone an existing OS template to create a new copy.

This command creates a new OS template that is an exact copy of an existing one,
including all its assets. The cloned template is always created with private
visibility. Name and label can be optionally overridden.

Required arguments:
  os_template_id    The numeric ID of the template to clone

Optional flags:
  --name            Name for the cloned template (default: "<original-name> (clone)")
  --label           Label for the cloned template (default: slug of name)

Examples:
  # Clone template with ID 123
  metalcloud-cli os-template clone 123

  # Clone with custom name
  metalcloud-cli os-template clone 123 --name "My Custom Ubuntu"

  # Clone with custom name and label
  metalcloud-cli os-template clone 123 --name "My Custom Ubuntu" --label "my-custom-ubuntu"`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_TEMPLATES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return os_template.OsTemplateClone(cmd.Context(), args[0], osTemplateFlags.name, osTemplateFlags.label)
		},
	}

	osTemplateExportCmd = &cobra.Command{
		Use:     "export <os_template_id>",
		Aliases: []string{"export-to-archive"},
		Short:   "Export an OS template and its assets to a zip archive",
		Long: `Export an OS template and all its assets to a zip archive file.

This command fetches a template by ID and packs its configuration and content
assets into a portable zip archive. URL-based assets (such as ISO links)
are preserved as references without downloading the actual files.

The archive contains:
  - template.yaml: Template configuration in YAML format
  - assets/: Directory containing decoded asset file contents

Required arguments:
  os_template_id    The numeric ID of the template to export

Optional flags:
  --output          Output file path (default: <template-name-slug>.zip)

Examples:
  # Export template with ID 123
  metalcloud-cli os-template export 123

  # Export to a specific file
  metalcloud-cli os-template export 123 --output my-template.zip`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_TEMPLATES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return os_template.OsTemplateExport(cmd.Context(), args[0], osTemplateFlags.outputPath)
		},
	}

	osTemplateImportCmd = &cobra.Command{
		Use:     "import <archive_path>",
		Aliases: []string{"import-from-archive"},
		Short:   "Import an OS template from a zip archive",
		Long: `Import an OS template from a zip archive file.

This command reads a previously exported zip archive and creates a new
OS template with all its assets. The new template is always created with
private visibility.

The archive should contain:
  - template.yaml: Template configuration in YAML format
  - assets/: Directory containing asset file contents

Required arguments:
  archive_path      Path to the zip archive file

Required flags:
  --name            Name for the new template

Optional flags:
  --label           Label for the new template (default: slug of name)

Examples:
  # Import a template
  metalcloud-cli os-template import my-template.zip --name "My Imported Template"

  # Import with custom label
  metalcloud-cli os-template import my-template.zip --name "My Template" --label "my-template-v2"`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_TEMPLATES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return os_template.OsTemplateImport(cmd.Context(), args[0], osTemplateFlags.name, osTemplateFlags.label)
		},
	}

	osTemplateCreateFromRepoCmd = &cobra.Command{
		Use:     "create-from-repo <os_template_path>",
		Aliases: []string{"add-from-repo", "clone-from-repo"},
		Short:   "Create a new OS template by cloning from a repository",
		Long: `Create a new OS template by cloning an existing template from a repository.

This command downloads and creates a local OS template based on a template
available in a remote repository. You can optionally customize the name,
label, and source ISO image during the creation process.

Required arguments:
  os_template_path  Path to the template within the repository
                   Use 'list-repo' command to see available templates

Optional flags:
  --repo-url        URL of the repository to clone from
                   Defaults to the official MetalSoft template repository
  --repo-username   Username for private repository authentication  
  --repo-password   Password for private repository authentication
  --name           Custom name for the new template (overrides original)
  --label          Custom label for the new template (overrides original)
  --source-iso     Custom source ISO image path (overrides original)

Flag dependencies:
  - If accessing a private repository, both --repo-username and --repo-password
    are required together

Examples:
  # Clone template from default public repository
  metalcloud-cli os-template create-from-repo ubuntu/22.04/server
  
  # Clone with custom name and label
  metalcloud-cli os-template create-from-repo ubuntu/22.04/server \
    --name "My Ubuntu 22.04" --label "my-ubuntu-2204"
  
  # Clone from private repository with custom ISO
  metalcloud-cli os-template create-from-repo centos/7/server \
    --repo-url https://private.com/templates \
    --repo-username user --repo-password pass \
    --source-iso /path/to/custom.iso`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_TEMPLATES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return os_template.OsTemplateCreateFromRepo(
				cmd.Context(),
				args[0],
				osTemplateFlags.repoUrl,
				osTemplateFlags.repoUsername,
				osTemplateFlags.repoPassword,
				osTemplateFlags.name,
				osTemplateFlags.label,
				osTemplateFlags.sourceIso,
			)
		},
	}
)

func init() {
	rootCmd.AddCommand(osTemplateCmd)

	osTemplateCmd.AddCommand(osTemplateListCmd)
	osTemplateCmd.AddCommand(osTemplateGetCmd)

	osTemplateCmd.AddCommand(osTemplateExampleCreateCmd)

	osTemplateCmd.AddCommand(osTemplateCreateCmd)
	osTemplateCreateCmd.Flags().StringVar(&osTemplateFlags.configSource, "config-source", "", "Source of the new OS template configuration. Can be 'pipe' or path to a JSON file.")
	osTemplateCreateCmd.MarkFlagsOneRequired("config-source")

	osTemplateCmd.AddCommand(osTemplateUpdateCmd)
	osTemplateUpdateCmd.Flags().StringVar(&osTemplateFlags.configSource, "config-source", "", "Source of the OS template configuration updates. Can be 'pipe' or path to a JSON file.")
	osTemplateUpdateCmd.MarkFlagsOneRequired("config-source")

	osTemplateCmd.AddCommand(osTemplateSetStatusCmd)
	osTemplateCmd.AddCommand(osTemplateDeleteCmd)
	osTemplateCmd.AddCommand(osTemplateGetCredentialsCmd)
	osTemplateCmd.AddCommand(osTemplateGetAssetsCmd)

	osTemplateCmd.AddCommand(osTemplateListRepoCmd)
	osTemplateListRepoCmd.Flags().StringVar(&osTemplateFlags.repoUrl, "repo-url", "", "Private repo to use.")
	osTemplateListRepoCmd.Flags().StringVar(&osTemplateFlags.repoUsername, "repo-username", "", "Private repo username.")
	osTemplateListRepoCmd.Flags().StringVar(&osTemplateFlags.repoPassword, "repo-password", "", "Private repo password.")
	osTemplateListRepoCmd.MarkFlagsRequiredTogether("repo-username", "repo-password")

	osTemplateCmd.AddCommand(osTemplateCloneCmd)
	osTemplateCloneCmd.Flags().StringVar(&osTemplateFlags.name, "name", "", "Name of the cloned OS template.")
	osTemplateCloneCmd.Flags().StringVar(&osTemplateFlags.label, "label", "", "Label of the cloned OS template.")

	osTemplateCmd.AddCommand(osTemplateExportCmd)
	osTemplateExportCmd.Flags().StringVar(&osTemplateFlags.outputPath, "output", "", "Output file path for the exported archive.")

	osTemplateCmd.AddCommand(osTemplateImportCmd)
	osTemplateImportCmd.Flags().StringVar(&osTemplateFlags.name, "name", "", "Name of the new OS template.")
	osTemplateImportCmd.Flags().StringVar(&osTemplateFlags.label, "label", "", "Label of the new OS template.")
	osTemplateImportCmd.MarkFlagsOneRequired("name")

	osTemplateCmd.AddCommand(osTemplateCreateFromRepoCmd)
	osTemplateCreateFromRepoCmd.Flags().StringVar(&osTemplateFlags.repoUrl, "repo-url", "", "Private repo to use.")
	osTemplateCreateFromRepoCmd.Flags().StringVar(&osTemplateFlags.repoUsername, "repo-username", "", "Private repo username.")
	osTemplateCreateFromRepoCmd.Flags().StringVar(&osTemplateFlags.repoPassword, "repo-password", "", "Private repo password.")
	osTemplateCreateFromRepoCmd.Flags().StringVar(&osTemplateFlags.name, "name", "", "Name of the OS template.")
	osTemplateCreateFromRepoCmd.Flags().StringVar(&osTemplateFlags.label, "label", "", "Label of the OS template.")
	osTemplateCreateFromRepoCmd.Flags().StringVar(&osTemplateFlags.sourceIso, "source-iso", "", "The source ISO image path.")
	osTemplateCreateFromRepoCmd.MarkFlagsRequiredTogether("repo-username", "repo-password")
}
