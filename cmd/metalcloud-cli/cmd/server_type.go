package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/server_type"
	"github.com/spf13/cobra"
)

// Server Type commands
var (
	serverTypeCmd = &cobra.Command{
		Use:   "server-type [command]",
		Short: "Manage server types and hardware configurations",
		Long: `Manage server types and view detailed hardware specifications.

Server types define the hardware configurations available for provisioning,
including CPU, memory, storage, and network interface specifications.

Available Commands:
  list    List all available server types
  get     Get detailed information about a specific server type

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
)

func init() {
	rootCmd.AddCommand(serverTypeCmd)

	// Server Type commands
	serverTypeCmd.AddCommand(serverTypeListCmd)
	serverTypeCmd.AddCommand(serverTypeGetCmd)
}
