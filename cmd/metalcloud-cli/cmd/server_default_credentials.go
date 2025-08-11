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
		Short:   "Manage server default credentials and authentication settings",
		Long: `Manage server default credentials and authentication settings for bare metal servers.

Server default credentials store authentication information (username, password) and optional
server metadata (rack location, inventory ID, UUID) that can be used during server provisioning
and management operations. These credentials are encrypted and stored securely.

Available commands:
  list           List all server default credentials
  get            Get detailed information about specific credentials
  get-credentials Retrieve unencrypted password for credentials
  create         Create new server default credentials
  delete         Delete existing server default credentials

Examples:
  # List all server default credentials
  metalcloud-cli server-default-credentials list

  # Get specific credentials information
  metalcloud-cli server-default-credentials get 123

  # Create credentials with required fields only
  metalcloud-cli sdc create --site-id 1 --serial "ABC123" --mac "00:11:22:33:44:55" --username "admin" --password "secret"`,
	}

	serverDefaultCredentialsListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all server default credentials",
		Long: `List all server default credentials with pagination support.

This command displays server default credentials in a tabular format showing ID, site ID,
server serial number, MAC address, username, and optional metadata like rack information.

Flags:
  --page    Page number for paginated results (default: 0, which returns all results)
  --limit   Number of records per page, maximum 100 (default: 0, which returns all results)

Examples:
  # List all server default credentials
  metalcloud-cli server-default-credentials list

  # List credentials with pagination (page 2, 10 records per page)
  metalcloud-cli sdc list --page 2 --limit 10

  # List first 25 credentials
  metalcloud-cli srv-dc ls --limit 25`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_DEFAULT_CREDENTIALS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_default_credentials.ServerDefaultCredentialsList(cmd.Context(), serverDefaultCredentialsFlags.pageFlag, serverDefaultCredentialsFlags.limitFlag)
		},
	}

	serverDefaultCredentialsGetCmd = &cobra.Command{
		Use:   "get <credentials_id>",
		Short: "Get detailed information about specific server default credentials",
		Long: `Get detailed information about specific server default credentials.

This command retrieves comprehensive information about a specific set of server default
credentials, including all metadata fields but not the actual password (use get-credentials
for that). The output includes ID, site ID, server serial number, MAC address, username,
and any optional metadata like rack information, inventory ID, and UUID.

Arguments:
  credentials_id    The ID of the server default credentials to retrieve (required)

Examples:
  # Get information about credentials with ID 123
  metalcloud-cli server-default-credentials get 123

  # Get credentials info using short alias
  metalcloud-cli sdc get 456

  # Get credentials info using alternate alias
  metalcloud-cli srv-dc get 789`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_DEFAULT_CREDENTIALS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_default_credentials.ServerDefaultCredentialsGet(cmd.Context(), args[0])
		},
	}

	serverDefaultCredentialsGetCredentialsCmd = &cobra.Command{
		Use:     "get-credentials <credentials_id>",
		Aliases: []string{"get-password", "password"},
		Short:   "Retrieve unencrypted password for server default credentials",
		Long: `Retrieve the unencrypted password for specific server default credentials.

This command returns the decrypted username and password for a specific set of server
default credentials. Use this when you need the actual password values for authentication
or configuration purposes. The password is decrypted server-side and transmitted securely.

Arguments:
  credentials_id    The ID of the server default credentials to retrieve password for (required)

Examples:
  # Get password for credentials with ID 123
  metalcloud-cli server-default-credentials get-credentials 123

  # Get password using alias
  metalcloud-cli sdc get-password 456

  # Get password using short alias
  metalcloud-cli srv-dc password 789`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_DEFAULT_CREDENTIALS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_default_credentials.ServerDefaultCredentialsGetCredentials(cmd.Context(), args[0])
		},
	}

	serverDefaultCredentialsCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create new server default credentials",
		Long: `Create new server default credentials for a server.

This command creates a new set of default credentials that can be used during server
provisioning and management. The password is encrypted before storage and can be retrieved
using the get-credentials command.

Required Flags:
  --site-id     Site ID where the server is located (required)
  --serial      Server serial number for identification (required)
  --mac         Server MAC address for network identification (required)
  --username    Default username for server authentication (required)
  --password    Default password for server authentication (required)

Optional Flags:
  --rack-name              Default rack name where server is located
  --rack-position-lower    Default rack position lower unit (e.g., "1")
  --rack-position-upper    Default rack position upper unit (e.g., "2")
  --inventory-id           Default inventory ID for asset tracking
  --uuid                   Default UUID for server identification

Examples:
  # Create credentials with required fields only
  metalcloud-cli server-default-credentials create \
    --site-id 1 \
    --serial "ABC123456" \
    --mac "00:11:22:33:44:55" \
    --username "admin" \
    --password "securepassword"

  # Create credentials with rack information
  metalcloud-cli sdc create \
    --site-id 2 \
    --serial "DEF789012" \
    --mac "aa:bb:cc:dd:ee:ff" \
    --username "root" \
    --password "complexpass123" \
    --rack-name "R1-A" \
    --rack-position-lower "10" \
    --rack-position-upper "12"

  # Create credentials with all optional metadata
  metalcloud-cli srv-dc create \
    --site-id 3 \
    --serial "GHI345678" \
    --mac "11:22:33:44:55:66" \
    --username "operator" \
    --password "mypassword" \
    --rack-name "R2-B" \
    --rack-position-lower "5" \
    --rack-position-upper "6" \
    --inventory-id "INV-2024-001" \
    --uuid "550e8400-e29b-41d4-a716-446655440000"`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_DEFAULT_CREDENTIALS_WRITE},
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
		Use:     "delete <credentials_id>",
		Aliases: []string{"rm"},
		Short:   "Delete server default credentials",
		Long: `Delete server default credentials by ID.

This command permanently removes a set of server default credentials from the system.
Once deleted, the credentials cannot be recovered and will no longer be available for
server provisioning or management operations.

Arguments:
  credentials_id    The ID of the server default credentials to delete (required)

Examples:
  # Delete credentials with ID 123
  metalcloud-cli server-default-credentials delete 123

  # Delete using short alias
  metalcloud-cli sdc rm 456

  # Delete using alternate alias
  metalcloud-cli srv-dc delete 789`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_DEFAULT_CREDENTIALS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_default_credentials.ServerDefaultCredentialsDelete(cmd.Context(), args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(serverDefaultCredentialsCmd)

	serverDefaultCredentialsCmd.AddCommand(serverDefaultCredentialsListCmd)
	serverDefaultCredentialsListCmd.Flags().IntVar(&serverDefaultCredentialsFlags.pageFlag, "page", 0, "Page number")
	serverDefaultCredentialsListCmd.Flags().IntVar(&serverDefaultCredentialsFlags.limitFlag, "limit", 0, "Number of records per page (max 100)")

	serverDefaultCredentialsCmd.AddCommand(serverDefaultCredentialsGetCmd)
	serverDefaultCredentialsCmd.AddCommand(serverDefaultCredentialsGetCredentialsCmd)

	serverDefaultCredentialsCmd.AddCommand(serverDefaultCredentialsCreateCmd)
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
	serverDefaultCredentialsCreateCmd.MarkFlagRequired("site-id")
	serverDefaultCredentialsCreateCmd.MarkFlagRequired("serial")
	serverDefaultCredentialsCreateCmd.MarkFlagRequired("mac")
	serverDefaultCredentialsCreateCmd.MarkFlagRequired("username")
	serverDefaultCredentialsCreateCmd.MarkFlagRequired("password")

	serverDefaultCredentialsCmd.AddCommand(serverDefaultCredentialsDeleteCmd)
}
