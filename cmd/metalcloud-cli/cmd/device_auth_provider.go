package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/device_auth_provider"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	deviceAuthProviderFlags = struct {
		configSource string
		sharedSecret string
		filterSiteId []string
		filterKind   []string
		filterStatus []string
	}{}

	siteDeviceAuthProviderCmd = &cobra.Command{
		Use:     "device-auth-provider [command]",
		Aliases: []string{"device-auth", "auth-provider"},
		Short:   "Manage device authentication providers (e.g. TACACS+)",
		Long: `Manage device authentication providers used by network devices for AAA
(authentication, authorization, accounting). Currently supports TACACS+ providers.

Available Commands:
  list                  List all device auth providers
  get                   Retrieve a specific device auth provider
  create                Create a new device auth provider from JSON
  update                Update an existing device auth provider from JSON
  delete                Delete a device auth provider
  credentials           Retrieve decrypted credentials for a provider
  update-shared-secret  Rotate the shared secret for a provider
  config-example        Print a JSON template suitable for create`,
	}

	siteDeviceAuthProviderListCmd = &cobra.Command{
		Use:          "list [flags...]",
		Aliases:      []string{"ls"},
		Short:        "List all device auth providers",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_DEVICE_AUTH_PROVIDERS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return device_auth_provider.DeviceAuthProviderList(
				cmd.Context(),
				deviceAuthProviderFlags.filterSiteId,
				deviceAuthProviderFlags.filterKind,
				deviceAuthProviderFlags.filterStatus,
			)
		},
	}

	siteDeviceAuthProviderGetCmd = &cobra.Command{
		Use:          "get provider_id_or_label",
		Aliases:      []string{"show"},
		Short:        "Retrieve a device auth provider by ID or label",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_DEVICE_AUTH_PROVIDERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return device_auth_provider.DeviceAuthProviderGet(cmd.Context(), args[0])
		},
	}

	siteDeviceAuthProviderCreateCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Short:   "Create a device auth provider from JSON",
		Long: `Create a new device auth provider from a JSON document.

The configuration source must supply all required fields: label, name, siteId,
kind, ipAddress, port, sharedSecret, username. Run 'config-example' to print a
template.`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_DEVICE_AUTH_PROVIDERS_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(deviceAuthProviderFlags.configSource)
			if err != nil {
				return err
			}
			return device_auth_provider.DeviceAuthProviderCreate(cmd.Context(), config)
		},
	}

	siteDeviceAuthProviderUpdateCmd = &cobra.Command{
		Use:          "update provider_id_or_label",
		Aliases:      []string{"edit"},
		Short:        "Update a device auth provider from JSON",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_DEVICE_AUTH_PROVIDERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(deviceAuthProviderFlags.configSource)
			if err != nil {
				return err
			}
			return device_auth_provider.DeviceAuthProviderUpdate(cmd.Context(), args[0], config)
		},
	}

	siteDeviceAuthProviderDeleteCmd = &cobra.Command{
		Use:          "delete provider_id_or_label",
		Aliases:      []string{"rm"},
		Short:        "Delete a device auth provider",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_DEVICE_AUTH_PROVIDERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return device_auth_provider.DeviceAuthProviderDelete(cmd.Context(), args[0])
		},
	}

	siteDeviceAuthProviderCredentialsCmd = &cobra.Command{
		Use:   "credentials provider_id_or_label",
		Short: "Show the decrypted credentials for a device auth provider",
		Long: `Show the decrypted username, password, and shared secret for a device auth
provider. Secrets are printed in plain text.`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_DEVICE_AUTH_PROVIDER_CREDENTIALS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return device_auth_provider.DeviceAuthProviderGetCredentials(cmd.Context(), args[0])
		},
	}

	siteDeviceAuthProviderUpdateSharedSecretCmd = &cobra.Command{
		Use:   "update-shared-secret provider_id_or_label",
		Short: "Rotate the shared secret for a device auth provider",
		Long: `Rotate the shared secret used to encrypt communication with the auth server.
The provider must be in 'maintenance' or 'disabled' status and must have no
network devices linked to it.`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_DEVICE_AUTH_PROVIDERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return device_auth_provider.DeviceAuthProviderUpdateSharedSecret(
				cmd.Context(),
				args[0],
				deviceAuthProviderFlags.sharedSecret,
			)
		},
	}

	siteDeviceAuthProviderConfigExampleCmd = &cobra.Command{
		Use:          "config-example",
		Short:        "Print a JSON template for creating a device auth provider",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_DEVICE_AUTH_PROVIDERS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return device_auth_provider.DeviceAuthProviderConfigExample(cmd.Context())
		},
	}
)

func init() {
	siteCmd.AddCommand(siteDeviceAuthProviderCmd)

	siteDeviceAuthProviderCmd.AddCommand(siteDeviceAuthProviderListCmd)
	siteDeviceAuthProviderListCmd.Flags().StringSliceVar(&deviceAuthProviderFlags.filterSiteId, "filter-site-id", nil, "Filter providers by site ID.")
	siteDeviceAuthProviderListCmd.Flags().StringSliceVar(&deviceAuthProviderFlags.filterKind, "filter-kind", nil, "Filter providers by kind (e.g. tacacs).")
	siteDeviceAuthProviderListCmd.Flags().StringSliceVar(&deviceAuthProviderFlags.filterStatus, "filter-status", nil, "Filter providers by status (e.g. active, maintenance, disabled).")

	siteDeviceAuthProviderCmd.AddCommand(siteDeviceAuthProviderGetCmd)

	siteDeviceAuthProviderCmd.AddCommand(siteDeviceAuthProviderCreateCmd)
	siteDeviceAuthProviderCreateCmd.Flags().StringVar(&deviceAuthProviderFlags.configSource, "config-source", "", "Source of the new device auth provider configuration. Can be 'pipe' or path to a JSON file.")
	siteDeviceAuthProviderCreateCmd.MarkFlagRequired("config-source")

	siteDeviceAuthProviderCmd.AddCommand(siteDeviceAuthProviderUpdateCmd)
	siteDeviceAuthProviderUpdateCmd.Flags().StringVar(&deviceAuthProviderFlags.configSource, "config-source", "", "Source of the device auth provider updates. Can be 'pipe' or path to a JSON file.")
	siteDeviceAuthProviderUpdateCmd.MarkFlagRequired("config-source")

	siteDeviceAuthProviderCmd.AddCommand(siteDeviceAuthProviderDeleteCmd)

	siteDeviceAuthProviderCmd.AddCommand(siteDeviceAuthProviderCredentialsCmd)

	siteDeviceAuthProviderCmd.AddCommand(siteDeviceAuthProviderUpdateSharedSecretCmd)
	siteDeviceAuthProviderUpdateSharedSecretCmd.Flags().StringVar(&deviceAuthProviderFlags.sharedSecret, "secret", "", "The new shared secret value.")
	siteDeviceAuthProviderUpdateSharedSecretCmd.MarkFlagRequired("secret")

	siteDeviceAuthProviderCmd.AddCommand(siteDeviceAuthProviderConfigExampleCmd)
}
