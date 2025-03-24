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
		Short: "Server type management",
		Long:  `Server type management commands.`,
	}

	serverTypeListCmd = &cobra.Command{
		Use:          "list",
		Aliases:      []string{"ls"},
		Short:        "Lists server types.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SERVERS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_type.ServerTypeList(cmd.Context())
		},
	}

	serverTypeGetCmd = &cobra.Command{
		Use:          "get",
		Aliases:      []string{"show"},
		Short:        "Get server type info.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SERVERS_READ},
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
