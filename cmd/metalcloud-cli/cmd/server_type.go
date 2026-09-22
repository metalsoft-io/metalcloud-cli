package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/server_type"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

// Server Type commands
var (
	serverTypeFlags = struct {
		createConfigSource string
		updateConfigSource string

		statisticsSiteId          int64
		statisticsUserId          int
		statisticsMaxResults      int
		statisticsInstanceArrayId int
	}{}

	serverTypeCmd = &cobra.Command{
		Use:   "server-type [command]",
		Short: "Manage server types and hardware configurations",
		Long: `Manage server types and view detailed hardware specifications.

Server types define the hardware configurations available for provisioning,
including CPU, memory, storage, and network interface specifications.

Available Commands:
  list            List all available server types
  get             Get detailed information about a specific server type
  create          Create a new server type
  update          Update an existing server type
  delete          Delete a server type
  clean-unused    Remove the server types that are no longer used
  statistics      Get the server availability statistics of a site
  config-example  Show a server type creation configuration example

Use "metalcloud server-type [command] --help" for more information about a command.`,
	}

	serverTypeListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all available server types",
		Long: `List all available server types with their hardware specifications.

This command displays server types in a tabular format showing key hardware 
characteristics including CPU count, RAM, storage, network interfaces, and GPU information.

The output includes:
- Server type ID and name
- Processor specifications (count, speed, names)
- Memory configuration (RAM in GB)
- Storage information (disk count)
- Network interface details
- GPU count (if applicable)

Examples:
  # List all server types
  metalcloud server-type list

  # List server types (using alias)
  metalcloud server-type ls

Required Permissions:
  - Server Types Read

Output Format:
  The command outputs data in table format by default. Use global output flags
  to change the format (--output json, --output yaml, etc.).`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_TYPES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_type.ServerTypeList(cmd.Context())
		},
	}

	serverTypeGetCmd = &cobra.Command{
		Use:     "get <server-type-id>",
		Aliases: []string{"show"},
		Short:   "Get detailed information about a specific server type",
		Long: `Get detailed information about a specific server type by its ID.

This command retrieves comprehensive hardware specifications for a specific server type,
including detailed processor information, memory configuration, storage details,
network interface specifications, GPU information (if applicable), and other
hardware characteristics.

The detailed output includes:
- Server type ID, name, and label
- Complete processor specifications (count, speed, core count, names)
- Memory configuration (RAM in GB)
- Storage information (disk count and disk groups)
- Network interface details (count, speeds, total capacity)
- GPU information (count and detailed GPU info)
- Server class and boot type
- Various flags (experimental, unmanaged servers only, etc.)
- Allowed vendor SKU IDs
- Tags associated with the server type

Arguments:
  server-type-id    The numeric ID of the server type to retrieve

Examples:
  # Get information about server type with ID 123
  metalcloud server-type get 123

  # Get server type information (using alias)
  metalcloud server-type show 456

  # Get server type info with JSON output
  metalcloud server-type get 789 --output json

Required Permissions:
  - Server Types Read

Output Format:
  The command outputs data in table format by default. Use global output flags
  to change the format (--output json, --output yaml, etc.).`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_TYPES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_type.ServerTypeGet(cmd.Context(), args[0])
		},
	}

	serverTypeCreateCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Short:   "Create a new server type",
		Long: `Create a new server type from a JSON or YAML configuration.

The configuration describes the hardware the server type stands for: the
processor, memory, disk and network interface specifications together with the
server class and the optional GPU and disk group details.

Required Flags:
  --config-source    Source of the new server type configuration. Can be 'pipe' or path to a JSON file.

Examples:
  # Create a server type from a JSON file
  metalcloud-cli server-type create --config-source ./server-type.json

  # Create a server type from piped configuration
  metalcloud-cli server-type config-example | metalcloud-cli server-type create --config-source pipe
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_TYPES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(serverTypeFlags.createConfigSource)
			if err != nil {
				return err
			}
			return server_type.ServerTypeCreate(cmd.Context(), config)
		},
	}

	serverTypeUpdateCmd = &cobra.Command{
		Use:   "update server_type_id",
		Short: "Update an existing server type",
		Long: `Update an existing server type from a JSON or YAML configuration.

Only the descriptive attributes of a server type can be updated: its label,
name, description, experimental flag, tags and allowed vendor SKU ids. The
hardware specification itself is derived from the registered servers.

Required Arguments:
  server_type_id     The numeric ID of the server type to update

Required Flags:
  --config-source    Source of the server type update configuration. Can be 'pipe' or path to a JSON file.

Examples:
  # Update a server type from a JSON file
  metalcloud-cli server-type update 123 --config-source ./server-type-update.json

  # Update a server type from piped configuration
  echo '{"label":"m-32-128-2","description":"Updated"}' | metalcloud-cli server-type update 123 --config-source pipe
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_TYPES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(serverTypeFlags.updateConfigSource)
			if err != nil {
				return err
			}
			return server_type.ServerTypeUpdate(cmd.Context(), args[0], config)
		},
	}

	serverTypeDeleteCmd = &cobra.Command{
		Use:     "delete server_type_id",
		Aliases: []string{"rm"},
		Short:   "Delete a server type",
		Long: `Delete a server type.

A server type can only be deleted while no server references it.

Required Arguments:
  server_type_id     The numeric ID of the server type to delete

Examples:
  # Delete the server type with ID 123
  metalcloud-cli server-type delete 123

  # Delete using the alias
  metalcloud-cli server-type rm 123
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_TYPES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_type.ServerTypeDelete(cmd.Context(), args[0])
		},
	}

	serverTypeCleanUnusedCmd = &cobra.Command{
		Use:   "clean-unused",
		Short: "Remove the server types that are no longer used",
		Long: `Remove every server type that no server references anymore.

Server types are created automatically when servers with a new hardware
configuration are registered; this command cleans up the ones left behind when
those servers are removed.

Examples:
  # Remove all unused server types
  metalcloud-cli server-type clean-unused
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_TYPES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_type.ServerTypeCleanUnused(cmd.Context())
		},
	}

	serverTypeStatisticsCmd = &cobra.Command{
		Use:     "statistics [server_type_id...]",
		Aliases: []string{"stats"},
		Short:   "Get the server availability statistics of a site",
		Long: `Get the server availability statistics of a site.

The report contains the number of available servers of each server type in the
site, the details of those servers and the utilization report grouped by RAM,
server type name, product name and owner. Use the json or yaml output format to
see the nested per-server information and the utilization report.

Optional Arguments:
  server_type_id...              Restrict the report to these server type IDs

Required Flags:
  --site-id                      The ID of the site to report on

Optional Flags:
  --user-id                      Only count the resources owned by this user ID
  --max-results-per-server-type  Maximum number of servers returned per server type
  --instance-array-id            Treat only the active instances of this instance array as available

Examples:
  # Statistics for every server type in site 1
  metalcloud-cli server-type statistics --site-id 1

  # Statistics for two server types only
  metalcloud-cli server-type statistics 12 13 --site-id 1

  # Statistics for the servers owned by user 5, at most 10 servers per type
  metalcloud-cli server-type statistics --site-id 1 --user-id 5 --max-results-per-server-type 10
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_TYPES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_type.ServerTypeStatistics(cmd.Context(),
				serverTypeFlags.statisticsSiteId,
				args,
				serverTypeFlags.statisticsUserId,
				serverTypeFlags.statisticsMaxResults,
				serverTypeFlags.statisticsInstanceArrayId)
		},
	}

	serverTypeConfigExampleCmd = &cobra.Command{
		Use:   "config-example",
		Short: "Show a server type creation configuration example",
		Long: `Show an example of the configuration accepted by 'server-type create'.

Examples:
  # Write the example to a file and edit it
  metalcloud-cli server-type config-example > server-type.json
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_TYPES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_type.ServerTypeConfigExample(cmd.Context())
		},
	}

	serverTypeUpdateConfigExampleCmd = &cobra.Command{
		Use:   "update-config-example",
		Short: "Show a server type update configuration example",
		Long: `Show an example of the configuration accepted by 'server-type update'.

Examples:
  # Write the example to a file and edit it
  metalcloud-cli server-type update-config-example > server-type-update.json
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_TYPES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_type.ServerTypeUpdateConfigExample(cmd.Context())
		},
	}
)

func init() {
	rootCmd.AddCommand(serverTypeCmd)

	// Server Type commands
	serverTypeCmd.AddCommand(serverTypeListCmd)
	serverTypeCmd.AddCommand(serverTypeGetCmd)

	serverTypeCmd.AddCommand(serverTypeCreateCmd)
	serverTypeCreateCmd.Flags().StringVar(&serverTypeFlags.createConfigSource, "config-source", "", "Source of the new server type configuration. Can be 'pipe' or path to a JSON file.")
	serverTypeCreateCmd.MarkFlagsOneRequired("config-source")

	serverTypeCmd.AddCommand(serverTypeUpdateCmd)
	serverTypeUpdateCmd.Flags().StringVar(&serverTypeFlags.updateConfigSource, "config-source", "", "Source of the server type update configuration. Can be 'pipe' or path to a JSON file.")
	serverTypeUpdateCmd.MarkFlagsOneRequired("config-source")

	serverTypeCmd.AddCommand(serverTypeDeleteCmd)

	serverTypeCmd.AddCommand(serverTypeCleanUnusedCmd)

	serverTypeCmd.AddCommand(serverTypeStatisticsCmd)
	serverTypeStatisticsCmd.Flags().Int64Var(&serverTypeFlags.statisticsSiteId, "site-id", 0, "The ID of the site to report on.")
	serverTypeStatisticsCmd.Flags().IntVar(&serverTypeFlags.statisticsUserId, "user-id", 0, "Only count the resources owned by this user ID.")
	serverTypeStatisticsCmd.Flags().IntVar(&serverTypeFlags.statisticsMaxResults, "max-results-per-server-type", 0, "Maximum number of servers returned per server type.")
	serverTypeStatisticsCmd.Flags().IntVar(&serverTypeFlags.statisticsInstanceArrayId, "instance-array-id", 0, "Treat only the active instances of this instance array as available.")
	serverTypeStatisticsCmd.MarkFlagRequired("site-id")

	serverTypeCmd.AddCommand(serverTypeConfigExampleCmd)
	serverTypeCmd.AddCommand(serverTypeUpdateConfigExampleCmd)
}
