package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/variable"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	variableFlags = struct {
		configSource string
	}{}

	variableCmd = &cobra.Command{
		Use:     "variable [command]",
		Aliases: []string{"var", "vars"},
		Short:   "Manage variables for infrastructure configuration",
		Long: `Manage variables that can be used across your infrastructure configurations.
Variables store key-value pairs that can be referenced in templates, scripts, and other configurations.

Available Commands:
  list           List all variables
  get            Get details of a specific variable
  create         Create a new variable
  update         Update an existing variable
  delete         Delete a variable
  config-example Show example configuration format

Examples:
  # List all variables
  metalcloud-cli variable list

  # Get details of a specific variable
  metalcloud-cli variable get 123

  # Create a new variable from a JSON file
  metalcloud-cli variable create --config-source /path/to/config.json

  # Create a new variable from stdin
  echo '{"name":"my-var","value":{"key1":"value1"}}' | metalcloud-cli variable create --config-source pipe`,
	}

	variableListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all variables",
		Long: `List all variables available in the current account.

This command displays all variables that have been created, showing their ID, name, 
value (truncated if long), usage type, owner, and timestamps.

Required Permissions:
  VARIABLES_AND_SECRETS_READ

Examples:
  # List all variables
  metalcloud-cli variable list
  
  # List variables using alias
  metalcloud-cli var ls`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VARIABLES_AND_SECRETS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return variable.VariableList(cmd.Context())
		},
	}

	variableGetCmd = &cobra.Command{
		Use:     "get variable_id",
		Aliases: []string{"show"},
		Short:   "Get details of a specific variable",
		Long: `Get detailed information about a specific variable by its ID.

This command retrieves and displays all details of a variable including its ID, name,
complete value, usage type, owner information, and timestamps.

Required Arguments:
  variable_id    Numeric ID of the variable to retrieve

Required Permissions:
  VARIABLES_AND_SECRETS_READ

Examples:
  # Get variable details by ID
  metalcloud-cli variable get 123
  
  # Get variable details using alias
  metalcloud-cli var show 456`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VARIABLES_AND_SECRETS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return variable.VariableGet(cmd.Context(), args[0])
		},
	}

	variableConfigExampleCmd = &cobra.Command{
		Use:   "config-example",
		Short: "Show variable configuration example",
		Long: `Display an example configuration format for creating or updating variables.

This command outputs a JSON template showing the structure and available fields
for variable configuration. Use this template as a reference when creating
configuration files or JSON input for variable operations.

Configuration Fields:
  name     (required) - The variable name (string)
  value    (required) - The variable value (key-value object)
  usage    (optional) - Variable usage type (string)

Required Permissions:
  VARIABLES_AND_SECRETS_WRITE

Examples:
  # Show configuration example
  metalcloud-cli variable config-example
  
  # Save example to file for editing
  metalcloud-cli variable config-example > my-variable.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VARIABLES_AND_SECRETS_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			return variable.VariableConfigExample(cmd.Context())
		},
	}

	variableCreateCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Short:   "Create a new variable",
		Long: `Create a new variable with specified configuration.

This command creates a new variable using configuration provided through a JSON file
or piped input. The configuration must include at minimum the variable name and value.

Required Flags:
  --config-source    Source of the variable configuration data

Configuration Source Options:
  pipe              Read configuration from stdin (use with echo or cat)
  /path/to/file     Read configuration from specified JSON file

Configuration Format:
  The configuration must be valid JSON with the following structure:
  {
    "name": "variable-name",           // Required: Variable name
    "value": {                         // Required: Key-value pairs
      "key1": "value1",
      "key2": "value2"
    },
    "usage": "general"                 // Optional: Usage type
  }

Required Permissions:
  VARIABLES_AND_SECRETS_WRITE

Examples:
  # Create variable from JSON file
  metalcloud-cli variable create --config-source /path/to/config.json
  
  # Create variable from stdin using pipe
  echo '{"name":"my-var","value":{"env":"prod","region":"us-east"}}' | metalcloud-cli variable create --config-source pipe
  
  # Create variable from file using cat
  cat my-variable.json | metalcloud-cli variable create --config-source pipe
  
  # Using alias
  metalcloud-cli var new --config-source /tmp/variable.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VARIABLES_AND_SECRETS_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(variableFlags.configSource)
			if err != nil {
				return err
			}

			return variable.VariableCreate(cmd.Context(), config)
		},
	}

	variableUpdateCmd = &cobra.Command{
		Use:     "update variable_id",
		Aliases: []string{"edit"},
		Short:   "Update an existing variable",
		Long: `Update an existing variable with new configuration values.

This command updates a variable's configuration using data provided through a JSON file
or piped input. You can modify the variable's name, value, and usage type.

Required Arguments:
  variable_id    Numeric ID of the variable to update

Required Flags:
  --config-source    Source of the variable configuration updates

Configuration Source Options:
  pipe              Read configuration from stdin (use with echo or cat)
  /path/to/file     Read configuration from specified JSON file

Configuration Format:
  The configuration must be valid JSON. You can include any of these fields:
  {
    "name": "new-variable-name",       // Optional: New variable name
    "value": {                         // Optional: New key-value pairs
      "key1": "new-value1",
      "key2": "new-value2"
    },
    "usage": "specific"                // Optional: New usage type
  }

  Note: Only the fields you specify will be updated. Omitted fields remain unchanged.

Required Permissions:
  VARIABLES_AND_SECRETS_WRITE

Examples:
  # Update variable from JSON file
  metalcloud-cli variable update 123 --config-source /path/to/update.json
  
  # Update variable from stdin using pipe
  echo '{"name":"updated-name","value":{"new-key":"new-value"}}' | metalcloud-cli variable update 123 --config-source pipe
  
  # Update only the value using cat
  cat value-update.json | metalcloud-cli variable update 456 --config-source pipe
  
  # Using alias
  metalcloud-cli var edit 789 --config-source /tmp/variable-update.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VARIABLES_AND_SECRETS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(variableFlags.configSource)
			if err != nil {
				return err
			}

			return variable.VariableUpdate(cmd.Context(), args[0], config)
		},
	}

	variableDeleteCmd = &cobra.Command{
		Use:     "delete variable_id",
		Aliases: []string{"rm"},
		Short:   "Delete a variable",
		Long: `Delete an existing variable by its ID.

This command permanently removes a variable from your account. The variable
will no longer be available for use in configurations, templates, or scripts.

Required Arguments:
  variable_id    Numeric ID of the variable to delete

Required Permissions:
  VARIABLES_AND_SECRETS_WRITE

Examples:
  # Delete variable by ID
  metalcloud-cli variable delete 123
  
  # Delete variable using alias
  metalcloud-cli var rm 456`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VARIABLES_AND_SECRETS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return variable.VariableDelete(cmd.Context(), args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(variableCmd)

	variableCmd.AddCommand(variableListCmd)
	variableCmd.AddCommand(variableGetCmd)
	variableCmd.AddCommand(variableConfigExampleCmd)

	variableCmd.AddCommand(variableCreateCmd)
	variableCreateCmd.Flags().StringVar(&variableFlags.configSource, "config-source", "", "Source of the new variable configuration. Can be 'pipe' or path to a JSON file.")
	variableCreateCmd.MarkFlagsOneRequired("config-source")

	variableCmd.AddCommand(variableUpdateCmd)
	variableUpdateCmd.Flags().StringVar(&variableFlags.configSource, "config-source", "", "Source of the variable configuration updates. Can be 'pipe' or path to a JSON file.")
	variableUpdateCmd.MarkFlagsOneRequired("config-source")

	variableCmd.AddCommand(variableDeleteCmd)
}
