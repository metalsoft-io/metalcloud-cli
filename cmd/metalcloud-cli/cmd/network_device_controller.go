package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/network_device_controller"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
	"github.com/spf13/cobra"
)

var (
	networkDeviceControllerFlags = struct {
		configSource            string
		filterId                []string
		filterSiteId            []string
		filterDatacenterName    []string
		filterManagementAddress []string
		filterIdentifierString  []string
		siteId                  int64
		datacenterName          string
		identifierString        string
		driver                  string
		managementAddress       string
		managementPort          int32
		username                string
		managementPassword      string
		description             string
	}{}

	networkDeviceControllerCmd = &cobra.Command{
		Use:     "network-device-controller [command]",
		Aliases: []string{"nd-controller", "ndc"},
		Short:   "Network device controller management",
		Long: `Manage network device controllers: the fabric controllers (Cisco NDFC/ACI,
NVIDIA UFM, Brocade) that MetalSoft drives instead of talking to each switch directly.

Command categories:
  Lifecycle:   list, get, create, update, delete, config-example
  Operations:  get-credentials, deploy-confirm`,
	}

	networkDeviceControllerListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List network device controllers",
		Long: `List all network device controllers.

Optional Flags:
  --filter-id strings                  Filter by controller ID. Repeatable or comma-separated.
  --filter-site-id strings             Filter by site ID. Repeatable or comma-separated.
  --filter-datacenter-name strings     Filter by datacenter name. Repeatable or comma-separated.
  --filter-management-address strings  Filter by management address. Repeatable or comma-separated.
  --filter-identifier-string strings   Filter by identifier (hostname). Repeatable or comma-separated.

Examples:
  metalcloud network-device-controller list
  metalcloud ndc list --filter-site-id 1`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_DEVICE_CONTROLLERS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device_controller.NetworkDeviceControllerList(cmd.Context(), network_device_controller.NetworkDeviceControllerListFilters{
				Id:                networkDeviceControllerFlags.filterId,
				SiteId:            networkDeviceControllerFlags.filterSiteId,
				DatacenterName:    networkDeviceControllerFlags.filterDatacenterName,
				ManagementAddress: networkDeviceControllerFlags.filterManagementAddress,
				IdentifierString:  networkDeviceControllerFlags.filterIdentifierString,
			})
		},
	}

	networkDeviceControllerGetCmd = &cobra.Command{
		Use:     "get controller_id_or_identifier",
		Aliases: []string{"show"},
		Short:   "Get network device controller details",
		Long: `Get the details of a network device controller.

Required Arguments:
  controller_id_or_identifier   The ID or identifier (hostname) of the controller

Examples:
  metalcloud network-device-controller get 7
  metalcloud ndc get ndfc-controller-01`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_DEVICE_CONTROLLERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device_controller.NetworkDeviceControllerGet(cmd.Context(), args[0])
		},
	}

	networkDeviceControllerConfigExampleCmd = &cobra.Command{
		Use:   "config-example",
		Short: "Print an example controller create configuration",
		Long: `Print an example configuration that can be edited and passed to
'network-device-controller create --config-source'.

Examples:
  metalcloud network-device-controller config-example > controller.json
  metalcloud ndc config-example -f yaml`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_DEVICE_CONTROLLERS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device_controller.NetworkDeviceControllerConfigExample(cmd.Context())
		},
	}

	networkDeviceControllerCreateCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Short:   "Create a network device controller",
		Long: `Register a new network device controller.

The controller can be described either with a configuration file (--config-source)
or with individual flags.

Required Flags (one of):
  --config-source string        Source of the controller configuration. Can be 'pipe' or path to a JSON/YAML file.
  --management-address string   Management address of the controller (used together with the flags below)

Required Flags (when not using --config-source):
  --datacenter-name string      Name of the datacenter the controller manages
  --driver string               Controller driver (cisco_ndfc, cisco_aci51, nvidia_ufm, brocade)
  --username string             Management username
  --management-password string  Management password

Optional Flags (when not using --config-source):
  --management-port int         Management port (default 443)
  --site-id int                 ID of the site the controller belongs to
  --identifier-string string    Identifier (hostname) of the controller
  --description string          Free-text description

Examples:
  metalcloud network-device-controller create --config-source controller.json
  metalcloud ndc create --management-address 10.0.0.50 --datacenter-name dc1 \
      --driver cisco_ndfc --username admin --management-password secret`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_DEVICE_CONTROLLERS_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			var create sdk.CreateNetworkDeviceController

			if networkDeviceControllerFlags.configSource != "" {
				config, err := utils.ReadConfigFromPipeOrFile(networkDeviceControllerFlags.configSource)
				if err != nil {
					return err
				}
				if err := utils.UnmarshalContent(config, &create); err != nil {
					return err
				}
			} else {
				driver, err := sdk.NewSwitchControllerDriverFromValue(networkDeviceControllerFlags.driver)
				if err != nil {
					return err
				}

				create = sdk.CreateNetworkDeviceController{
					DatacenterName:     networkDeviceControllerFlags.datacenterName,
					Driver:             *driver,
					ManagementAddress:  networkDeviceControllerFlags.managementAddress,
					ManagementPort:     networkDeviceControllerFlags.managementPort,
					Username:           networkDeviceControllerFlags.username,
					ManagementPassword: networkDeviceControllerFlags.managementPassword,
				}
				if networkDeviceControllerFlags.siteId != 0 {
					create.SiteId = sdk.PtrInt64(networkDeviceControllerFlags.siteId)
				}
				if networkDeviceControllerFlags.identifierString != "" {
					create.IdentifierString = sdk.PtrString(networkDeviceControllerFlags.identifierString)
				}
				if networkDeviceControllerFlags.description != "" {
					create.Description = sdk.PtrString(networkDeviceControllerFlags.description)
				}
			}

			return network_device_controller.NetworkDeviceControllerCreate(cmd.Context(), create)
		},
	}

	networkDeviceControllerUpdateCmd = &cobra.Command{
		Use:     "update controller_id_or_identifier",
		Aliases: []string{"edit"},
		Short:   "Update a network device controller",
		Long: `Update the settings of a network device controller.

Required Arguments:
  controller_id_or_identifier   The ID or identifier (hostname) of the controller

Required Flags:
  --config-source string        Source of the updated configuration. Can be 'pipe' or path to a JSON/YAML file.

Examples:
  metalcloud network-device-controller update 7 --config-source update.json
  echo '{"description":"new description"}' | metalcloud ndc update ndfc-controller-01 --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_DEVICE_CONTROLLERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(networkDeviceControllerFlags.configSource)
			if err != nil {
				return err
			}

			return network_device_controller.NetworkDeviceControllerUpdate(cmd.Context(), args[0], config)
		},
	}

	networkDeviceControllerDeleteCmd = &cobra.Command{
		Use:     "delete controller_id_or_identifier",
		Aliases: []string{"rm"},
		Short:   "Delete a network device controller",
		Long: `Delete a network device controller.

Required Arguments:
  controller_id_or_identifier   The ID or identifier (hostname) of the controller

Examples:
  metalcloud network-device-controller delete 7`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_DEVICE_CONTROLLERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device_controller.NetworkDeviceControllerDelete(cmd.Context(), args[0])
		},
	}

	networkDeviceControllerGetCredentialsCmd = &cobra.Command{
		Use:     "get-credentials controller_id_or_identifier",
		Aliases: []string{"credentials"},
		Short:   "Get the management credentials of a controller",
		Long: `Get the management credentials used to connect to a network device controller.

Required Arguments:
  controller_id_or_identifier   The ID or identifier (hostname) of the controller

Examples:
  metalcloud network-device-controller get-credentials 7`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_DEVICE_CONTROLLERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device_controller.NetworkDeviceControllerGetCredentials(cmd.Context(), args[0])
		},
	}

	networkDeviceControllerDeployConfirmCmd = &cobra.Command{
		Use:     "deploy-confirm controller_id_or_identifier",
		Aliases: []string{"confirm-deploy"},
		Short:   "Confirm a pending controller deploy",
		Long: `Confirm a pending deploy of the configuration generated for a network device
controller, letting it proceed.

Required Arguments:
  controller_id_or_identifier   The ID or identifier (hostname) of the controller

Examples:
  metalcloud network-device-controller deploy-confirm 7`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_DEVICE_CONTROLLERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device_controller.NetworkDeviceControllerDeployConfirm(cmd.Context(), args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(networkDeviceControllerCmd)

	networkDeviceControllerCmd.AddCommand(networkDeviceControllerListCmd)
	networkDeviceControllerListCmd.Flags().StringSliceVar(&networkDeviceControllerFlags.filterId, "filter-id", nil, "Filter by controller ID.")
	networkDeviceControllerListCmd.Flags().StringSliceVar(&networkDeviceControllerFlags.filterSiteId, "filter-site-id", nil, "Filter by site ID.")
	networkDeviceControllerListCmd.Flags().StringSliceVar(&networkDeviceControllerFlags.filterDatacenterName, "filter-datacenter-name", nil, "Filter by datacenter name.")
	networkDeviceControllerListCmd.Flags().StringSliceVar(&networkDeviceControllerFlags.filterManagementAddress, "filter-management-address", nil, "Filter by management address.")
	networkDeviceControllerListCmd.Flags().StringSliceVar(&networkDeviceControllerFlags.filterIdentifierString, "filter-identifier-string", nil, "Filter by identifier (hostname).")

	networkDeviceControllerCmd.AddCommand(networkDeviceControllerGetCmd)
	networkDeviceControllerCmd.AddCommand(networkDeviceControllerConfigExampleCmd)

	networkDeviceControllerCmd.AddCommand(networkDeviceControllerCreateCmd)
	networkDeviceControllerCreateCmd.Flags().StringVar(&networkDeviceControllerFlags.configSource, "config-source", "", "Source of the new controller configuration. Can be 'pipe' or path to a JSON/YAML file.")
	networkDeviceControllerCreateCmd.Flags().StringVar(&networkDeviceControllerFlags.managementAddress, "management-address", "", "Management address of the controller.")
	networkDeviceControllerCreateCmd.Flags().StringVar(&networkDeviceControllerFlags.datacenterName, "datacenter-name", "", "Name of the datacenter the controller manages.")
	networkDeviceControllerCreateCmd.Flags().StringVar(&networkDeviceControllerFlags.driver, "driver", "", "Controller driver (cisco_ndfc, cisco_aci51, nvidia_ufm, brocade).")
	networkDeviceControllerCreateCmd.Flags().StringVar(&networkDeviceControllerFlags.username, "username", "", "Management username.")
	networkDeviceControllerCreateCmd.Flags().StringVar(&networkDeviceControllerFlags.managementPassword, "management-password", "", "Management password.")
	networkDeviceControllerCreateCmd.Flags().Int32Var(&networkDeviceControllerFlags.managementPort, "management-port", 443, "Management port.")
	networkDeviceControllerCreateCmd.Flags().Int64Var(&networkDeviceControllerFlags.siteId, "site-id", 0, "ID of the site the controller belongs to.")
	networkDeviceControllerCreateCmd.Flags().StringVar(&networkDeviceControllerFlags.identifierString, "identifier-string", "", "Identifier (hostname) of the controller.")
	networkDeviceControllerCreateCmd.Flags().StringVar(&networkDeviceControllerFlags.description, "description", "", "Description of the controller.")
	networkDeviceControllerCreateCmd.MarkFlagsOneRequired("config-source", "management-address")
	networkDeviceControllerCreateCmd.MarkFlagsMutuallyExclusive("config-source", "management-address")
	networkDeviceControllerCreateCmd.MarkFlagsRequiredTogether("management-address", "datacenter-name", "driver", "username", "management-password")

	networkDeviceControllerCmd.AddCommand(networkDeviceControllerUpdateCmd)
	networkDeviceControllerUpdateCmd.Flags().StringVar(&networkDeviceControllerFlags.configSource, "config-source", "", "Source of the updated controller configuration. Can be 'pipe' or path to a JSON/YAML file.")
	networkDeviceControllerUpdateCmd.MarkFlagsOneRequired("config-source")

	networkDeviceControllerCmd.AddCommand(networkDeviceControllerDeleteCmd)
	networkDeviceControllerCmd.AddCommand(networkDeviceControllerGetCredentialsCmd)
	networkDeviceControllerCmd.AddCommand(networkDeviceControllerDeployConfirmCmd)
}
