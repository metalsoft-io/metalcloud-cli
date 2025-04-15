package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/server_default_credentials"
	"github.com/spf13/cobra"
)

var (
	serverDefaultCredentialsFlags = struct {
		pageFlag               int
		limitFlag              int
		siteIdFlag             float32
		serverSerialNumberFlag string
		serverMacAddressFlag   string
		usernameFlag           string
		passwordFlag           string
		rackNameFlag           string
		rackPositionLowerFlag  string
		rackPositionUpperFlag  string
		inventoryIdFlag        string
		uuidFlag               string
	}{}

	serverDefaultCredentialsCmd = &cobra.Command{
		Use:     "server-default-credentials [command]",
		Aliases: []string{"srv-dc", "sdc"},
		Short:   "Server default credentials management",
		Long:    `Server default credentials management commands.`,
	}

	serverDefaultCredentialsListCmd = &cobra.Command{
		Use:          "list",
		Aliases:      []string{"ls"},
		Short:        "Lists server default credentials.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_default_credentials.ServerDefaultCredentialsList(cmd.Context(), serverDefaultCredentialsFlags.pageFlag, serverDefaultCredentialsFlags.limitFlag)
		},
	}

	serverDefaultCredentialsGetCmd = &cobra.Command{
		Use:          "get <credentials_id>",
		Short:        "Get server default credentials information.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_default_credentials.ServerDefaultCredentialsGet(cmd.Context(), args[0])
		},
	}

	serverDefaultCredentialsGetCredentialsCmd = &cobra.Command{
		Use:          "get-credentials <credentials_id>",
		Aliases:      []string{"get-password", "password"},
		Short:        "Get server default credentials unencrypted password.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_default_credentials.ServerDefaultCredentialsGetCredentials(cmd.Context(), args[0])
		},
	}

	serverDefaultCredentialsCreateCmd = &cobra.Command{
		Use:          "create",
		Short:        "Create server default credentials.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_default_credentials.ServerDefaultCredentialsCreate(
				cmd.Context(),
				serverDefaultCredentialsFlags.siteIdFlag,
				serverDefaultCredentialsFlags.serverSerialNumberFlag,
				serverDefaultCredentialsFlags.serverMacAddressFlag,
				serverDefaultCredentialsFlags.usernameFlag,
				serverDefaultCredentialsFlags.passwordFlag,
				serverDefaultCredentialsFlags.rackNameFlag,
				serverDefaultCredentialsFlags.rackPositionLowerFlag,
				serverDefaultCredentialsFlags.rackPositionUpperFlag,
				serverDefaultCredentialsFlags.inventoryIdFlag,
				serverDefaultCredentialsFlags.uuidFlag,
			)
		},
	}

	serverDefaultCredentialsDeleteCmd = &cobra.Command{
		Use:          "delete <credentials_id>",
		Aliases:      []string{"rm"},
		Short:        "Delete server default credentials.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_default_credentials.ServerDefaultCredentialsDelete(cmd.Context(), args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(serverDefaultCredentialsCmd)

	// Server Default Credentials commands
	serverDefaultCredentialsCmd.AddCommand(serverDefaultCredentialsListCmd)
	serverDefaultCredentialsCmd.AddCommand(serverDefaultCredentialsGetCmd)
	serverDefaultCredentialsCmd.AddCommand(serverDefaultCredentialsGetCredentialsCmd)
	serverDefaultCredentialsCmd.AddCommand(serverDefaultCredentialsCreateCmd)
	serverDefaultCredentialsCmd.AddCommand(serverDefaultCredentialsDeleteCmd)

	// Add flags for list command
	serverDefaultCredentialsListCmd.Flags().IntVar(&serverDefaultCredentialsFlags.pageFlag, "page", 0, "Page number")
	serverDefaultCredentialsListCmd.Flags().IntVar(&serverDefaultCredentialsFlags.limitFlag, "limit", 0, "Number of records per page (max 100)")

	// Add flags for create command
	serverDefaultCredentialsCreateCmd.Flags().Float32Var(&serverDefaultCredentialsFlags.siteIdFlag, "site-id", 0, "Site ID")
	serverDefaultCredentialsCreateCmd.Flags().StringVar(&serverDefaultCredentialsFlags.serverSerialNumberFlag, "serial", "", "Server serial number")
	serverDefaultCredentialsCreateCmd.Flags().StringVar(&serverDefaultCredentialsFlags.serverMacAddressFlag, "mac", "", "Server MAC address")
	serverDefaultCredentialsCreateCmd.Flags().StringVar(&serverDefaultCredentialsFlags.usernameFlag, "username", "", "Default username")
	serverDefaultCredentialsCreateCmd.Flags().StringVar(&serverDefaultCredentialsFlags.passwordFlag, "password", "", "Default password")
	serverDefaultCredentialsCreateCmd.Flags().StringVar(&serverDefaultCredentialsFlags.rackNameFlag, "rack-name", "", "Default rack name")
	serverDefaultCredentialsCreateCmd.Flags().StringVar(&serverDefaultCredentialsFlags.rackPositionLowerFlag, "rack-position-lower", "", "Default rack position lower unit")
	serverDefaultCredentialsCreateCmd.Flags().StringVar(&serverDefaultCredentialsFlags.rackPositionUpperFlag, "rack-position-upper", "", "Default rack position upper unit")
	serverDefaultCredentialsCreateCmd.Flags().StringVar(&serverDefaultCredentialsFlags.inventoryIdFlag, "inventory-id", "", "Default inventory ID")
	serverDefaultCredentialsCreateCmd.Flags().StringVar(&serverDefaultCredentialsFlags.uuidFlag, "uuid", "", "Default UUID")

	// Required flags
	serverDefaultCredentialsCreateCmd.MarkFlagRequired("site-id")
	serverDefaultCredentialsCreateCmd.MarkFlagRequired("serial")
	serverDefaultCredentialsCreateCmd.MarkFlagRequired("mac")
	serverDefaultCredentialsCreateCmd.MarkFlagRequired("username")
	serverDefaultCredentialsCreateCmd.MarkFlagRequired("password")
}
