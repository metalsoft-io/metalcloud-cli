package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/network_fabric_interconnect"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
	"github.com/spf13/cobra"
)

var (
	networkFabricInterconnectFlags = struct {
		configSource        string
		filterStatus        []string
		label               string
		name                string
		description         string
		interconnectType    string
		bgpTemplateId       int64
		transportId         int64
		requireConfirmation bool
		linkIds             []string
	}{}

	networkFabricInterconnectCmd = &cobra.Command{
		Use:     "network-fabric-interconnect [command]",
		Aliases: []string{"fabric-interconnect", "interconnect", "nfi"},
		Short:   "Network fabric interconnect management",
		Long: `Manage network fabric interconnects (e.g. EVPN data center interconnects) that
join two or more network fabrics through BGP.

Command categories:
  Lifecycle:   list, get, create, update, delete, config-example, template
  Links:       get-links, get-link, add-link, remove-link, activate-links, deactivate-links
  Fabrics:     get-fabrics, get-available-fabrics
  Deployment:  deploy, deployment-check, deployment-info, accept-deploy, reject-deploy, detach`,
	}

	networkFabricInterconnectListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List network fabric interconnects",
		Long: `List all network fabric interconnects.

Optional Flags:
  --filter-status strings   Filter by status (e.g. draft, active). Repeatable or comma-separated.

Examples:
  metalcloud network-fabric-interconnect list
  metalcloud nfi list --filter-status active`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRIC_INTERCONNECTS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_fabric_interconnect.InterconnectList(cmd.Context(), networkFabricInterconnectFlags.filterStatus)
		},
	}

	networkFabricInterconnectGetCmd = &cobra.Command{
		Use:     "get interconnect_id_or_label",
		Aliases: []string{"show"},
		Short:   "Get network fabric interconnect details",
		Long: `Get the details of a network fabric interconnect, including its deploy preview.

Required Arguments:
  interconnect_id_or_label   The ID or label of the interconnect

Examples:
  metalcloud network-fabric-interconnect get 12
  metalcloud nfi get dc1-dc2-interconnect`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRIC_INTERCONNECTS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_fabric_interconnect.InterconnectGet(cmd.Context(), args[0])
		},
	}

	networkFabricInterconnectConfigExampleCmd = &cobra.Command{
		Use:   "config-example",
		Short: "Print an example interconnect create configuration",
		Long: `Print an example configuration that can be edited and passed to
'network-fabric-interconnect create --config-source'.

Examples:
  metalcloud network-fabric-interconnect config-example > interconnect.json
  metalcloud nfi config-example -f yaml`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRIC_INTERCONNECTS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_fabric_interconnect.InterconnectConfigExample(cmd.Context())
		},
	}

	networkFabricInterconnectCreateCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Short:   "Create a network fabric interconnect",
		Long: `Create a new network fabric interconnect.

The interconnect can be described either with a configuration file (--config-source)
or with individual flags (--label, --bgp-template-id, ...).

Required Flags (one of):
  --config-source string     Source of the interconnect configuration. Can be 'pipe' or path to a JSON/YAML file.
  --label string             Label of the new interconnect (used together with the flags below)

Optional Flags (when not using --config-source):
  --bgp-template-id int      ID of the BGP interconnect configuration template (required with --label)
  --type string              Interconnect type (default "dci-evpn")
  --name string              Display name
  --description string       Description
  --transport-id int         ID of the transport used by the interconnect

Examples:
  metalcloud network-fabric-interconnect create --config-source interconnect.json
  cat interconnect.yaml | metalcloud nfi create --config-source pipe
  metalcloud nfi create --label dc1-dc2 --bgp-template-id 3 --name "DC1 to DC2"`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRIC_INTERCONNECTS_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			var create sdk.CreateNetworkFabricInterconnect

			if networkFabricInterconnectFlags.configSource != "" {
				config, err := utils.ReadConfigFromPipeOrFile(networkFabricInterconnectFlags.configSource)
				if err != nil {
					return err
				}
				if err := utils.UnmarshalContent(config, &create); err != nil {
					return err
				}
			} else {
				interconnectType, err := sdk.NewNetworkFabricInterconnectTypeFromValue(networkFabricInterconnectFlags.interconnectType)
				if err != nil {
					return err
				}

				create = sdk.CreateNetworkFabricInterconnect{
					InterconnectType:           *interconnectType,
					Label:                      networkFabricInterconnectFlags.label,
					BgpConfigurationTemplateId: networkFabricInterconnectFlags.bgpTemplateId,
				}
				if networkFabricInterconnectFlags.name != "" {
					create.Name = sdk.PtrString(networkFabricInterconnectFlags.name)
				}
				if networkFabricInterconnectFlags.description != "" {
					create.Description = sdk.PtrString(networkFabricInterconnectFlags.description)
				}
				if networkFabricInterconnectFlags.transportId != 0 {
					create.TransportId = sdk.PtrInt64(networkFabricInterconnectFlags.transportId)
				}
			}

			return network_fabric_interconnect.InterconnectCreate(cmd.Context(), create)
		},
	}

	networkFabricInterconnectUpdateCmd = &cobra.Command{
		Use:     "update interconnect_id_or_label",
		Aliases: []string{"edit"},
		Short:   "Update a network fabric interconnect",
		Long: `Update the label, name, description or BGP template of a network fabric interconnect.

Required Arguments:
  interconnect_id_or_label   The ID or label of the interconnect

Required Flags:
  --config-source string     Source of the updated configuration. Can be 'pipe' or path to a JSON/YAML file.

Examples:
  metalcloud network-fabric-interconnect update 12 --config-source update.json
  echo '{"description":"new description"}' | metalcloud nfi update dc1-dc2 --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRIC_INTERCONNECTS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(networkFabricInterconnectFlags.configSource)
			if err != nil {
				return err
			}

			return network_fabric_interconnect.InterconnectUpdate(cmd.Context(), args[0], config)
		},
	}

	networkFabricInterconnectDeleteCmd = &cobra.Command{
		Use:     "delete interconnect_id_or_label",
		Aliases: []string{"rm"},
		Short:   "Delete a network fabric interconnect",
		Long: `Delete a network fabric interconnect.

Required Arguments:
  interconnect_id_or_label   The ID or label of the interconnect

Examples:
  metalcloud network-fabric-interconnect delete 12`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRIC_INTERCONNECTS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_fabric_interconnect.InterconnectDelete(cmd.Context(), args[0])
		},
	}

	networkFabricInterconnectTemplateCmd = &cobra.Command{
		Use:   "template interconnect_type",
		Short: "Show the BGP templates used for an interconnect type",
		Long: `Show the global and neighbor BGP templates that are rendered when an
interconnect of the given type is activated or deactivated.

Required Arguments:
  interconnect_type   The interconnect type (e.g. dci-evpn)

Examples:
  metalcloud network-fabric-interconnect template dci-evpn`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRIC_INTERCONNECTS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_fabric_interconnect.InterconnectTemplateGet(cmd.Context(), args[0])
		},
	}

	networkFabricInterconnectDeployCmd = &cobra.Command{
		Use:   "deploy interconnect_id_or_label",
		Short: "Deploy a network fabric interconnect",
		Long: `Deploy a network fabric interconnect, pushing the BGP configuration to the
network devices referenced by its links. Returns the job that performs the deploy.

Required Arguments:
  interconnect_id_or_label   The ID or label of the interconnect

Optional Flags:
  --require-confirmation     Pause the deploy until it is accepted with 'accept-deploy'
                             (or discarded with 'reject-deploy')

Examples:
  metalcloud network-fabric-interconnect deploy 12
  metalcloud nfi deploy dc1-dc2 --require-confirmation`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRIC_INTERCONNECTS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_fabric_interconnect.InterconnectDeploy(cmd.Context(), args[0], networkFabricInterconnectFlags.requireConfirmation)
		},
	}

	networkFabricInterconnectDeploymentCheckCmd = &cobra.Command{
		Use:     "deployment-check interconnect_id_or_label",
		Aliases: []string{"check"},
		Short:   "Validate whether the interconnect links can be activated",
		Long: `Validate the links of a network fabric interconnect, reporting for each link
whether it can be activated, the resolved template variables and any errors.

Required Arguments:
  interconnect_id_or_label   The ID or label of the interconnect

Optional Flags:
  --link-ids strings         Restrict the check to these link IDs. Repeatable or comma-separated.

Examples:
  metalcloud network-fabric-interconnect deployment-check 12
  metalcloud nfi check dc1-dc2 --link-ids 3,4`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRIC_INTERCONNECTS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_fabric_interconnect.InterconnectDeploymentCheck(cmd.Context(), args[0], networkFabricInterconnectFlags.linkIds)
		},
	}

	networkFabricInterconnectDeploymentInfoCmd = &cobra.Command{
		Use:   "deployment-info interconnect_id_or_label",
		Short: "Show the deployment status and preview of an interconnect",
		Long: `Show the deployment status, the current deploy job and the per-device
configuration preview of a network fabric interconnect.

Required Arguments:
  interconnect_id_or_label   The ID or label of the interconnect

Examples:
  metalcloud network-fabric-interconnect deployment-info 12`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRIC_INTERCONNECTS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_fabric_interconnect.InterconnectDeploymentInfo(cmd.Context(), args[0])
		},
	}

	networkFabricInterconnectAcceptDeployCmd = &cobra.Command{
		Use:   "accept-deploy interconnect_id_or_label",
		Short: "Accept a pending interconnect deploy",
		Long: `Accept a deploy that was started with --require-confirmation so that it proceeds.

Required Arguments:
  interconnect_id_or_label   The ID or label of the interconnect

Examples:
  metalcloud network-fabric-interconnect accept-deploy 12`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRIC_INTERCONNECTS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_fabric_interconnect.InterconnectAcceptDeploy(cmd.Context(), args[0])
		},
	}

	networkFabricInterconnectRejectDeployCmd = &cobra.Command{
		Use:   "reject-deploy interconnect_id_or_label",
		Short: "Reject a pending interconnect deploy",
		Long: `Reject a deploy that was started with --require-confirmation, discarding it.

Required Arguments:
  interconnect_id_or_label   The ID or label of the interconnect

Examples:
  metalcloud network-fabric-interconnect reject-deploy 12`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRIC_INTERCONNECTS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_fabric_interconnect.InterconnectRejectDeploy(cmd.Context(), args[0])
		},
	}

	networkFabricInterconnectDetachCmd = &cobra.Command{
		Use:   "detach interconnect_id_or_label",
		Short: "Detach an interconnect, removing its configuration from the devices",
		Long: `Detach a network fabric interconnect. The BGP configuration is removed from all
linked network devices. Returns the job that performs the detach.

Required Arguments:
  interconnect_id_or_label   The ID or label of the interconnect

Examples:
  metalcloud network-fabric-interconnect detach 12`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRIC_INTERCONNECTS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_fabric_interconnect.InterconnectDetach(cmd.Context(), args[0])
		},
	}

	networkFabricInterconnectFabricsGetCmd = &cobra.Command{
		Use:     "get-fabrics interconnect_id_or_label",
		Aliases: []string{"fabrics"},
		Short:   "List the fabrics attached to an interconnect",
		Long: `List the network fabrics currently attached to a network fabric interconnect.

Required Arguments:
  interconnect_id_or_label   The ID or label of the interconnect

Examples:
  metalcloud network-fabric-interconnect get-fabrics 12`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRIC_INTERCONNECTS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_fabric_interconnect.InterconnectFabricsGet(cmd.Context(), args[0])
		},
	}

	networkFabricInterconnectAvailableFabricsGetCmd = &cobra.Command{
		Use:     "get-available-fabrics interconnect_id_or_label",
		Aliases: []string{"available-fabrics"},
		Short:   "List the fabrics that can be linked to an interconnect",
		Long: `List the network fabrics that are eligible to be linked to a network fabric interconnect.

Required Arguments:
  interconnect_id_or_label   The ID or label of the interconnect

Examples:
  metalcloud network-fabric-interconnect get-available-fabrics 12`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRIC_INTERCONNECTS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_fabric_interconnect.InterconnectAvailableFabricsGet(cmd.Context(), args[0])
		},
	}

	networkFabricInterconnectLinksGetCmd = &cobra.Command{
		Use:     "get-links interconnect_id_or_label",
		Aliases: []string{"links"},
		Short:   "List the links of an interconnect",
		Long: `List the links of a network fabric interconnect. Each link binds one fabric and
one network device (the border device) to the interconnect.

Required Arguments:
  interconnect_id_or_label   The ID or label of the interconnect

Examples:
  metalcloud network-fabric-interconnect get-links 12`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRIC_INTERCONNECTS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_fabric_interconnect.InterconnectLinksGet(cmd.Context(), args[0])
		},
	}

	networkFabricInterconnectLinkGetCmd = &cobra.Command{
		Use:   "get-link interconnect_id_or_label link_id",
		Short: "Get one link of an interconnect",
		Long: `Get the details of a single link of a network fabric interconnect.

Required Arguments:
  interconnect_id_or_label   The ID or label of the interconnect
  link_id                    The ID of the link

Examples:
  metalcloud network-fabric-interconnect get-link 12 3`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRIC_INTERCONNECTS_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_fabric_interconnect.InterconnectLinkGet(cmd.Context(), args[0], args[1])
		},
	}

	networkFabricInterconnectLinkAddCmd = &cobra.Command{
		Use:   "add-link interconnect_id_or_label fabric_id_or_label network_device_id_or_label",
		Short: "Add a link to an interconnect",
		Long: `Add a link to a network fabric interconnect, binding a fabric and one of its
network devices (the border device that will run the interconnect BGP session).

Required Arguments:
  interconnect_id_or_label      The ID or label of the interconnect
  fabric_id_or_label            The ID or name of the fabric to link
  network_device_id_or_label    The ID or identifier of the network device in that fabric

Examples:
  metalcloud network-fabric-interconnect add-link 12 dc1-fabric 45
  metalcloud nfi add-link dc1-dc2 7 border-leaf-01`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRIC_INTERCONNECTS_WRITE},
		Args:         cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_fabric_interconnect.InterconnectLinkAdd(cmd.Context(), args[0], args[1], args[2])
		},
	}

	networkFabricInterconnectLinkRemoveCmd = &cobra.Command{
		Use:   "remove-link interconnect_id_or_label link_id",
		Short: "Remove a link from an interconnect",
		Long: `Remove a link from a network fabric interconnect.

Required Arguments:
  interconnect_id_or_label   The ID or label of the interconnect
  link_id                    The ID of the link to remove

Examples:
  metalcloud network-fabric-interconnect remove-link 12 3`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRIC_INTERCONNECTS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_fabric_interconnect.InterconnectLinkRemove(cmd.Context(), args[0], args[1])
		},
	}

	networkFabricInterconnectLinksActivateCmd = &cobra.Command{
		Use:   "activate-links interconnect_id_or_label link_id...",
		Short: "Activate interconnect links",
		Long: `Activate one or more links of a network fabric interconnect, pushing the BGP
configuration to their network devices. Returns the job that performs the activation.

Required Arguments:
  interconnect_id_or_label   The ID or label of the interconnect
  link_id...                 One or more link IDs to activate

Optional Flags:
  --require-confirmation     Pause until the deploy is accepted with 'accept-deploy'

Examples:
  metalcloud network-fabric-interconnect activate-links 12 3 4`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRIC_INTERCONNECTS_WRITE},
		Args:         cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_fabric_interconnect.InterconnectLinksActivate(cmd.Context(), args[0], args[1:], networkFabricInterconnectFlags.requireConfirmation)
		},
	}

	networkFabricInterconnectLinksDeactivateCmd = &cobra.Command{
		Use:   "deactivate-links interconnect_id_or_label link_id...",
		Short: "Deactivate interconnect links",
		Long: `Deactivate one or more links of a network fabric interconnect, removing the BGP
configuration from their network devices. Returns the job that performs the deactivation.

Required Arguments:
  interconnect_id_or_label   The ID or label of the interconnect
  link_id...                 One or more link IDs to deactivate

Optional Flags:
  --require-confirmation     Pause until the deploy is accepted with 'accept-deploy'

Examples:
  metalcloud network-fabric-interconnect deactivate-links 12 3`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRIC_INTERCONNECTS_WRITE},
		Args:         cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_fabric_interconnect.InterconnectLinksDeactivate(cmd.Context(), args[0], args[1:], networkFabricInterconnectFlags.requireConfirmation)
		},
	}
)

func init() {
	rootCmd.AddCommand(networkFabricInterconnectCmd)

	networkFabricInterconnectCmd.AddCommand(networkFabricInterconnectListCmd)
	networkFabricInterconnectListCmd.Flags().StringSliceVar(&networkFabricInterconnectFlags.filterStatus, "filter-status", nil, "Filter by status (e.g. draft, active).")

	networkFabricInterconnectCmd.AddCommand(networkFabricInterconnectGetCmd)
	networkFabricInterconnectCmd.AddCommand(networkFabricInterconnectConfigExampleCmd)

	networkFabricInterconnectCmd.AddCommand(networkFabricInterconnectCreateCmd)
	networkFabricInterconnectCreateCmd.Flags().StringVar(&networkFabricInterconnectFlags.configSource, "config-source", "", "Source of the new interconnect configuration. Can be 'pipe' or path to a JSON/YAML file.")
	networkFabricInterconnectCreateCmd.Flags().StringVar(&networkFabricInterconnectFlags.label, "label", "", "Label of the new interconnect.")
	networkFabricInterconnectCreateCmd.Flags().StringVar(&networkFabricInterconnectFlags.name, "name", "", "Display name of the new interconnect.")
	networkFabricInterconnectCreateCmd.Flags().StringVar(&networkFabricInterconnectFlags.description, "description", "", "Description of the new interconnect.")
	networkFabricInterconnectCreateCmd.Flags().StringVar(&networkFabricInterconnectFlags.interconnectType, "type", string(sdk.NETWORKFABRICINTERCONNECTTYPE_DCI_EVPN), "Interconnect type.")
	networkFabricInterconnectCreateCmd.Flags().Int64Var(&networkFabricInterconnectFlags.bgpTemplateId, "bgp-template-id", 0, "ID of the BGP interconnect configuration template.")
	networkFabricInterconnectCreateCmd.Flags().Int64Var(&networkFabricInterconnectFlags.transportId, "transport-id", 0, "ID of the transport used by the interconnect.")
	networkFabricInterconnectCreateCmd.MarkFlagsOneRequired("config-source", "label")
	networkFabricInterconnectCreateCmd.MarkFlagsMutuallyExclusive("config-source", "label")
	networkFabricInterconnectCreateCmd.MarkFlagsRequiredTogether("label", "bgp-template-id")

	networkFabricInterconnectCmd.AddCommand(networkFabricInterconnectUpdateCmd)
	networkFabricInterconnectUpdateCmd.Flags().StringVar(&networkFabricInterconnectFlags.configSource, "config-source", "", "Source of the updated interconnect configuration. Can be 'pipe' or path to a JSON/YAML file.")
	networkFabricInterconnectUpdateCmd.MarkFlagsOneRequired("config-source")

	networkFabricInterconnectCmd.AddCommand(networkFabricInterconnectDeleteCmd)
	networkFabricInterconnectCmd.AddCommand(networkFabricInterconnectTemplateCmd)

	networkFabricInterconnectCmd.AddCommand(networkFabricInterconnectDeployCmd)
	networkFabricInterconnectDeployCmd.Flags().BoolVar(&networkFabricInterconnectFlags.requireConfirmation, "require-confirmation", false, "Pause the deploy until it is accepted with 'accept-deploy'.")

	networkFabricInterconnectCmd.AddCommand(networkFabricInterconnectDeploymentCheckCmd)
	networkFabricInterconnectDeploymentCheckCmd.Flags().StringSliceVar(&networkFabricInterconnectFlags.linkIds, "link-ids", nil, "Restrict the check to these link IDs.")

	networkFabricInterconnectCmd.AddCommand(networkFabricInterconnectDeploymentInfoCmd)
	networkFabricInterconnectCmd.AddCommand(networkFabricInterconnectAcceptDeployCmd)
	networkFabricInterconnectCmd.AddCommand(networkFabricInterconnectRejectDeployCmd)
	networkFabricInterconnectCmd.AddCommand(networkFabricInterconnectDetachCmd)

	networkFabricInterconnectCmd.AddCommand(networkFabricInterconnectFabricsGetCmd)
	networkFabricInterconnectCmd.AddCommand(networkFabricInterconnectAvailableFabricsGetCmd)

	networkFabricInterconnectCmd.AddCommand(networkFabricInterconnectLinksGetCmd)
	networkFabricInterconnectCmd.AddCommand(networkFabricInterconnectLinkGetCmd)
	networkFabricInterconnectCmd.AddCommand(networkFabricInterconnectLinkAddCmd)
	networkFabricInterconnectCmd.AddCommand(networkFabricInterconnectLinkRemoveCmd)

	networkFabricInterconnectCmd.AddCommand(networkFabricInterconnectLinksActivateCmd)
	networkFabricInterconnectLinksActivateCmd.Flags().BoolVar(&networkFabricInterconnectFlags.requireConfirmation, "require-confirmation", false, "Pause until the deploy is accepted with 'accept-deploy'.")

	networkFabricInterconnectCmd.AddCommand(networkFabricInterconnectLinksDeactivateCmd)
	networkFabricInterconnectLinksDeactivateCmd.Flags().BoolVar(&networkFabricInterconnectFlags.requireConfirmation, "require-confirmation", false, "Pause until the deploy is accepted with 'accept-deploy'.")
}
