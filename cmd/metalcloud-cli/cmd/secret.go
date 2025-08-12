package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/secret"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	secretFlags = struct {
		configSource string
	}{}

	secretCmd = &cobra.Command{
		Use:     "secret [command]",
		Aliases: []string{"sec", "secrets"},
		Short:   "Manage encrypted secrets for secure credential storage",
		Long: `Manage encrypted secrets for secure credential storage.

Secrets provide a secure way to store sensitive information like passwords, API keys,
and other credentials that can be referenced in your infrastructure configurations.
All secret values are encrypted at rest and in transit.

Available Commands:
  list                List all secrets
  get                 Get secret details by ID
  create              Create a new secret
  update              Update an existing secret
  delete              Delete a secret
  config-example      Show example configuration format

Examples:
  # List all secrets
  metalcloud-cli secret list

  # Get details of a specific secret
  metalcloud-cli secret get 123

  # Create a new secret from JSON file
  metalcloud-cli secret create --config-source ./secret.json

  # Create a secret from stdin
  echo '{"name":"my-secret","value":"secret-value","usage":"credential"}' | metalcloud-cli secret create --config-source pipe

  # Update a secret
  metalcloud-cli secret update 123 --config-source ./updated-secret.json

  # Delete a secret
  metalcloud-cli secret delete 123`,
	}

	secretListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all secrets",
		Long: `List all secrets in the current datacenter.

This command displays a table of all secrets with their basic information including:
- Secret ID and name
- Encrypted value (partial display for security)
- Usage type (credential, configuration, etc.)
- Owner information
- Creation and update timestamps

The output is formatted as a table by default and can be filtered or formatted
using global output flags.

Examples:
  # List all secrets
  metalcloud-cli secret list

  # List secrets with JSON output
  metalcloud-cli secret list --output json

  # List secrets with custom formatting
  metalcloud-cli secret list --output table --fields id,name,usage`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VARIABLES_AND_SECRETS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return secret.SecretList(cmd.Context())
		},
	}

	secretGetCmd = &cobra.Command{
		Use:     "get secret_id",
		Aliases: []string{"show"},
		Short:   "Get secret details by ID",
		Long: `Get detailed information about a specific secret by its ID.

This command retrieves and displays comprehensive information about a secret,
including its name, encrypted value (partial display for security), usage type,
owner details, and timestamps. The secret ID must be provided as a numeric value.

Arguments:
  secret_id          Numeric ID of the secret to retrieve (required)

Examples:
  # Get details of secret with ID 123
  metalcloud-cli secret get 123

  # Get secret details with JSON output
  metalcloud-cli secret get 456 --output json

  # Get secret details with custom field selection
  metalcloud-cli secret get 789 --output table --fields id,name,usage,createdTimestamp`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VARIABLES_AND_SECRETS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return secret.SecretGet(cmd.Context(), args[0])
		},
	}

	secretConfigExampleCmd = &cobra.Command{
		Use:   "config-example",
		Short: "Show example secret configuration format",
		Long: `Display an example JSON configuration structure for creating secrets.

This command outputs a sample configuration that can be used as a template
for creating new secrets. The configuration includes all available fields
with example values.

The configuration format includes:
- name: The secret name (required)
- value: The secret value to encrypt (required)  
- usage: The usage type (optional, defaults to "credential")

Available usage types:
- credential: For storing passwords, API keys, tokens
- configuration: For storing configuration values
- certificate: For storing SSL/TLS certificates
- ssh_key: For storing SSH keys

Examples:
  # Show configuration example
  metalcloud-cli secret config-example

  # Save example to file for editing
  metalcloud-cli secret config-example > my-secret.json

  # Use example as template with custom output format
  metalcloud-cli secret config-example --output json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VARIABLES_AND_SECRETS_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			return secret.SecretConfigExample(cmd.Context())
		},
	}

	secretCreateCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Short:   "Create a new secret",
		Long: `Create a new encrypted secret for secure credential storage.

This command creates a new secret that will be encrypted and stored securely.
The secret configuration must be provided through a JSON file or piped input.

Required Flags:
  --config-source    Source of the secret configuration (required)
                     Can be either 'pipe' for stdin input or a path to a JSON file

The configuration file must contain:
  name     (string)   The secret name (required)
  value    (string)   The secret value to encrypt (required)
  usage    (string)   The usage type (optional, defaults to "credential")

Available usage types:
- credential: For storing passwords, API keys, tokens
- configuration: For storing configuration values
- certificate: For storing SSL/TLS certificates
- ssh_key: For storing SSH keys

Examples:
  # Create a secret from JSON file
  metalcloud-cli secret create --config-source ./my-secret.json

  # Create a secret from stdin (pipe)
  echo '{"name":"api-key","value":"sk-1234567890","usage":"credential"}' | metalcloud-cli secret create --config-source pipe

  # Create from file with different usage type
  metalcloud-cli secret create --config-source ./ssh-key.json

Example JSON configuration:
  {
    "name": "database-password",
    "value": "super-secret-password",
    "usage": "credential"
  }`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VARIABLES_AND_SECRETS_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(secretFlags.configSource)
			if err != nil {
				return err
			}

			return secret.SecretCreate(cmd.Context(), config)
		},
	}

	secretUpdateCmd = &cobra.Command{
		Use:     "update secret_id",
		Aliases: []string{"edit"},
		Short:   "Update an existing secret",
		Long: `Update an existing secret with new configuration.

This command updates an existing secret identified by its ID. The updated
configuration must be provided through a JSON file or piped input.

Arguments:
  secret_id          Numeric ID of the secret to update (required)

Required Flags:
  --config-source    Source of the secret configuration updates (required)
                     Can be either 'pipe' for stdin input or a path to a JSON file

The configuration file can contain any combination of updateable fields:
  name     (string)   The secret name (optional)
  value    (string)   The secret value to encrypt (optional)
  usage    (string)   The usage type (optional)

Available usage types:
- credential: For storing passwords, API keys, tokens
- configuration: For storing configuration values
- certificate: For storing SSL/TLS certificates
- ssh_key: For storing SSH keys

Examples:
  # Update a secret from JSON file
  metalcloud-cli secret update 123 --config-source ./updated-secret.json

  # Update secret name and value from stdin
  echo '{"name":"new-api-key","value":"sk-0987654321"}' | metalcloud-cli secret update 456 --config-source pipe

  # Update only the usage type
  metalcloud-cli secret update 789 --config-source ./usage-update.json

Example JSON configuration for update:
  {
    "name": "updated-database-password",
    "value": "new-super-secret-password",
    "usage": "credential"
  }`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VARIABLES_AND_SECRETS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(secretFlags.configSource)
			if err != nil {
				return err
			}

			return secret.SecretUpdate(cmd.Context(), args[0], config)
		},
	}

	secretDeleteCmd = &cobra.Command{
		Use:     "delete secret_id",
		Aliases: []string{"rm"},
		Short:   "Delete a secret",
		Long: `Delete a secret by its ID.

This command permanently removes a secret from the system. The secret ID must
be provided as a numeric value. This action cannot be undone.

Arguments:
  secret_id          Numeric ID of the secret to delete (required)

Examples:
  # Delete a secret by ID
  metalcloud-cli secret delete 123

  # Delete a secret with confirmation
  metalcloud-cli secret delete 456 --auto-approve

Note: Be careful when deleting secrets as this action is irreversible.
Make sure the secret is not being used by any infrastructure configurations
before deletion.`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VARIABLES_AND_SECRETS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return secret.SecretDelete(cmd.Context(), args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(secretCmd)

	secretCmd.AddCommand(secretListCmd)
	secretCmd.AddCommand(secretGetCmd)
	secretCmd.AddCommand(secretConfigExampleCmd)

	secretCmd.AddCommand(secretCreateCmd)
	secretCreateCmd.Flags().StringVar(&secretFlags.configSource, "config-source", "", "Source of the new secret configuration. Can be 'pipe' or path to a JSON file.")
	secretCreateCmd.MarkFlagsOneRequired("config-source")

	secretCmd.AddCommand(secretUpdateCmd)
	secretUpdateCmd.Flags().StringVar(&secretFlags.configSource, "config-source", "", "Source of the secret configuration updates. Can be 'pipe' or path to a JSON file.")
	secretUpdateCmd.MarkFlagsOneRequired("config-source")

	secretCmd.AddCommand(secretDeleteCmd)
}
