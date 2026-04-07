package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/network_device"
	"github.com/spf13/cobra"
)

var (
	networkDeviceDefaultSecretsFlags = struct {
		pageFlag                     int
		limitFlag                    int
		siteIdFlag                   float32
		macAddressOrSerialNumberFlag string
		secretNameFlag               string
		secretValueFlag              string
		csvFileFlag                  string
	}{}

	networkDeviceDefaultSecretsCmd = &cobra.Command{
		Use:     "default-secrets [command]",
		Aliases: []string{"ds"},
		Short:   "Manage network device default secrets",
		Long: `Manage network device default secrets for network devices (switches).

Network device default secrets store secret values (such as passwords or keys) associated
with a specific network device identified by MAC address or serial number. These secrets
are encrypted and stored securely.

Examples:
  # List all network device default secrets
  metalcloud-cli network-device default-secrets list

  # Get specific secret information
  metalcloud-cli network-device default-secrets get 123

  # Create a new secret
  metalcloud-cli nd ds create --site-id 1 --mac-or-serial "AA:BB:CC:DD:EE:FF" --secret-name "admin_password" --secret-value "s3cur3"`,
	}

	networkDeviceDefaultSecretsListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all network device default secrets",
		Long: `List all network device default secrets with pagination support.

Flags:
  --page    Page number for paginated results (default: 0, which returns all results)
  --limit   Number of records per page, maximum 100 (default: 0, which returns all results)

Examples:
  # List all network device default secrets
  metalcloud-cli network-device default-secrets list

  # List with pagination
  metalcloud-cli nd ds list --page 2 --limit 10`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceDefaultSecretsList(cmd.Context(), networkDeviceDefaultSecretsFlags.pageFlag, networkDeviceDefaultSecretsFlags.limitFlag)
		},
	}

	networkDeviceDefaultSecretsGetCmd = &cobra.Command{
		Use:     "get <secrets_id>",
		Aliases: []string{"show"},
		Short:   "Get detailed information about a specific network device default secret",
		Long: `Get detailed information about a specific network device default secret.

This returns metadata about the secret (ID, site, MAC/serial, name, timestamps)
but not the actual secret value. Use get-credentials to retrieve the secret value.

Arguments:
  secrets_id    The ID of the network device default secrets to retrieve (required)

Examples:
  # Get secret information
  metalcloud-cli network-device default-secrets get 123

  # Using alias
  metalcloud-cli nd ds show 456`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceDefaultSecretsGet(cmd.Context(), args[0])
		},
	}

	networkDeviceDefaultSecretsGetCredentialsCmd = &cobra.Command{
		Use:     "get-credentials <secrets_id>",
		Aliases: []string{"get-secret"},
		Short:   "Retrieve the unencrypted secret value",
		Long: `Retrieve the unencrypted secret value for a specific network device default secret.

The secret value is decrypted server-side and returned in plain text.

Arguments:
  secrets_id    The ID of the network device default secrets (required)

Examples:
  # Get the secret value
  metalcloud-cli network-device default-secrets get-credentials 123

  # Using alias
  metalcloud-cli nd ds get-secret 456`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceDefaultSecretsGetCredentials(cmd.Context(), args[0])
		},
	}

	networkDeviceDefaultSecretsCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new network device default secret",
		Long: `Create a new network device default secret.

Required Flags:
  --site-id          Site ID where the network device is located
  --mac-or-serial    MAC address or serial number of the network device
  --secret-name      Name of the secret
  --secret-value     Value of the secret

Examples:
  # Create a new secret
  metalcloud-cli network-device default-secrets create \
    --site-id 1 \
    --mac-or-serial "AA:BB:CC:DD:EE:FF" \
    --secret-name "admin_password" \
    --secret-value "s3cur3"

  # Using alias
  metalcloud-cli nd ds create --site-id 2 --mac-or-serial "SN123456" --secret-name "enable_secret" --secret-value "mypass"`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceDefaultSecretsCreate(
				cmd.Context(),
				networkDeviceDefaultSecretsFlags.siteIdFlag,
				networkDeviceDefaultSecretsFlags.macAddressOrSerialNumberFlag,
				networkDeviceDefaultSecretsFlags.secretNameFlag,
				networkDeviceDefaultSecretsFlags.secretValueFlag,
			)
		},
	}

	networkDeviceDefaultSecretsUpdateCmd = &cobra.Command{
		Use:   "update <secrets_id>",
		Short: "Update an existing network device default secret",
		Long: `Update the secret value of an existing network device default secret.

Arguments:
  secrets_id    The ID of the network device default secrets to update (required)

Required Flags:
  --secret-value    New value of the secret

Examples:
  # Update the secret value
  metalcloud-cli network-device default-secrets update 123 --secret-value "new_s3cur3"

  # Using alias
  metalcloud-cli nd ds update 456 --secret-value "updated_password"`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceDefaultSecretsUpdate(
				cmd.Context(),
				args[0],
				networkDeviceDefaultSecretsFlags.secretValueFlag,
			)
		},
	}

	networkDeviceDefaultSecretsBatchCreateCmd = &cobra.Command{
		Use:   "batch-create",
		Short: "Create multiple network device default secrets from a CSV file",
		Long: `Create multiple network device default secrets from a CSV file.

The CSV file must have a header row with the following columns:
  siteId, macAddressOrSerialNumber, secretName, secretValue

Each subsequent row creates one secret.

Required Flags:
  --csv-file    Path to the CSV file

Example CSV file:
  siteId,macAddressOrSerialNumber,secretName,secretValue
  1,AA:BB:CC:DD:EE:FF,admin_password,s3cur3
  1,SN123456,enable_secret,mypass
  2,11:22:33:44:55:66,snmp_community,public

Examples:
  # Batch create secrets from a CSV file
  metalcloud-cli network-device default-secrets batch-create --csv-file secrets.csv

  # Using alias
  metalcloud-cli nd ds batch-create --csv-file /path/to/secrets.csv`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceDefaultSecretsBatchCreate(
				cmd.Context(),
				networkDeviceDefaultSecretsFlags.csvFileFlag,
			)
		},
	}

	networkDeviceDefaultSecretsDeleteCmd = &cobra.Command{
		Use:     "delete <secrets_id>",
		Aliases: []string{"rm"},
		Short:   "Delete a network device default secret",
		Long: `Delete a network device default secret by ID.

This operation permanently removes the secret and cannot be undone.

Arguments:
  secrets_id    The ID of the network device default secrets to delete (required)

Examples:
  # Delete a secret
  metalcloud-cli network-device default-secrets delete 123

  # Using alias
  metalcloud-cli nd ds rm 456`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SWITCHES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device.NetworkDeviceDefaultSecretsDelete(cmd.Context(), args[0])
		},
	}
)

func init() {
	networkDeviceCmd.AddCommand(networkDeviceDefaultSecretsCmd)

	networkDeviceDefaultSecretsCmd.AddCommand(networkDeviceDefaultSecretsListCmd)
	networkDeviceDefaultSecretsListCmd.Flags().IntVar(&networkDeviceDefaultSecretsFlags.pageFlag, "page", 0, "Page number")
	networkDeviceDefaultSecretsListCmd.Flags().IntVar(&networkDeviceDefaultSecretsFlags.limitFlag, "limit", 0, "Number of records per page (max 100)")

	networkDeviceDefaultSecretsCmd.AddCommand(networkDeviceDefaultSecretsGetCmd)

	networkDeviceDefaultSecretsCmd.AddCommand(networkDeviceDefaultSecretsGetCredentialsCmd)

	networkDeviceDefaultSecretsCmd.AddCommand(networkDeviceDefaultSecretsCreateCmd)
	networkDeviceDefaultSecretsCreateCmd.Flags().Float32Var(&networkDeviceDefaultSecretsFlags.siteIdFlag, "site-id", 0, "Site ID")
	networkDeviceDefaultSecretsCreateCmd.Flags().StringVar(&networkDeviceDefaultSecretsFlags.macAddressOrSerialNumberFlag, "mac-or-serial", "", "MAC address or serial number of the network device")
	networkDeviceDefaultSecretsCreateCmd.Flags().StringVar(&networkDeviceDefaultSecretsFlags.secretNameFlag, "secret-name", "", "Name of the secret")
	networkDeviceDefaultSecretsCreateCmd.Flags().StringVar(&networkDeviceDefaultSecretsFlags.secretValueFlag, "secret-value", "", "Value of the secret")
	networkDeviceDefaultSecretsCreateCmd.MarkFlagRequired("site-id")
	networkDeviceDefaultSecretsCreateCmd.MarkFlagRequired("mac-or-serial")
	networkDeviceDefaultSecretsCreateCmd.MarkFlagRequired("secret-name")
	networkDeviceDefaultSecretsCreateCmd.MarkFlagRequired("secret-value")

	networkDeviceDefaultSecretsCmd.AddCommand(networkDeviceDefaultSecretsUpdateCmd)
	networkDeviceDefaultSecretsUpdateCmd.Flags().StringVar(&networkDeviceDefaultSecretsFlags.secretValueFlag, "secret-value", "", "New value of the secret")
	networkDeviceDefaultSecretsUpdateCmd.MarkFlagRequired("secret-value")

	networkDeviceDefaultSecretsCmd.AddCommand(networkDeviceDefaultSecretsBatchCreateCmd)
	networkDeviceDefaultSecretsBatchCreateCmd.Flags().StringVar(&networkDeviceDefaultSecretsFlags.csvFileFlag, "csv-file", "", "Path to the CSV file")
	networkDeviceDefaultSecretsBatchCreateCmd.MarkFlagRequired("csv-file")

	networkDeviceDefaultSecretsCmd.AddCommand(networkDeviceDefaultSecretsDeleteCmd)
}
