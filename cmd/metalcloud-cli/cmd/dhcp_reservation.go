package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/dhcp_reservation"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
	"github.com/spf13/cobra"
)

var (
	dhcpReservationFlags = struct {
		configSource string
		macAddress   string
		circuitId    string
		deviceType   string
		vendor       string
		giaddr       string
		priority     int
		ip           string
		subnetId     int
		subnetPoolId []int
	}{}

	dhcpReservationCmd = &cobra.Command{
		Use:     "dhcp-reservation [command]",
		Aliases: []string{"dhcp-reservations"},
		Short:   "Manage site DHCP reservations",
		Long: `Manage the DHCP reservations of a site.

A DHCP reservation pins how the site's DHCP service answers a matching request:
either with a fixed IP (a manual allocation) or with the next free IP of one or
more IPAM OOB subnet pools (an auto allocation). Requests are matched on MAC
address, Option 82 circuit id, relay gateway address, device type and vendor.

Reservations are scoped to a site and an IP version, so every command takes
both: the site by ID or label, and 'ipv4' or 'ipv6'.

Available Commands:
  list            List the reservations of a site and IP version
  get             Get one reservation
  create          Create a reservation
  update          Update a reservation
  delete          Delete a reservation
  config-example  Print example create configurations

Examples:
  metalcloud-cli dhcp-reservation list dc-1 ipv4
  metalcloud-cli dhcp-reservation get dc-1 ipv4 12`,
	}

	dhcpReservationListCmd = &cobra.Command{
		Use:     "list site_id_or_label ip_version",
		Aliases: []string{"ls"},
		Short:   "List the DHCP reservations of a site",
		Long: `List all DHCP reservations of one site and IP version.

Required Arguments:
  site_id_or_label  The ID or label of the site
  ip_version        'ipv4' or 'ipv6'

Examples:
  metalcloud-cli dhcp-reservation list 1 ipv4
  metalcloud-cli dhcp-reservation ls dc-1 ipv6`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SITES_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return dhcp_reservation.DhcpReservationList(cmd.Context(), args[0], args[1])
		},
	}

	dhcpReservationGetCmd = &cobra.Command{
		Use:     "get site_id_or_label ip_version reservation_id",
		Aliases: []string{"show"},
		Short:   "Get one DHCP reservation",
		Long: `Display the details of a single DHCP reservation.

Required Arguments:
  site_id_or_label  The ID or label of the site
  ip_version        'ipv4' or 'ipv6'
  reservation_id    The ID of the reservation

Examples:
  metalcloud-cli dhcp-reservation get 1 ipv4 12`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SITES_READ},
		Args:         cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return dhcp_reservation.DhcpReservationGet(cmd.Context(), args[0], args[1], args[2])
		},
	}

	dhcpReservationConfigExampleCmd = &cobra.Command{
		Use:     "config-example",
		Aliases: []string{"example"},
		Short:   "Print example DHCP reservation configurations",
		Long: `Print one example create body per allocation kind (manual and auto).

Edit one of them and pass it to 'create' via --config-source.

Examples:
  metalcloud-cli dhcp-reservation config-example
  metalcloud-cli dhcp-reservation config-example -f yaml > reservation.yaml`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SITES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return dhcp_reservation.DhcpReservationConfigExample(cmd.Context())
		},
	}

	dhcpReservationCreateCmd = &cobra.Command{
		Use:     "create site_id_or_label ip_version",
		Aliases: []string{"new"},
		Short:   "Create a DHCP reservation",
		Long: `Create a DHCP reservation in a site for one IP version.

The reservation can be described either by individual flags or by a JSON/YAML
configuration passed with --config-source.

A manual allocation pins one IP: pass --ip together with --subnet-id. An auto
allocation takes the next free IP of one or more OOB subnet pools: pass
--subnet-pool-id (repeatable). The two forms are mutually exclusive.

Required Arguments:
  site_id_or_label  The ID or label of the site
  ip_version        'ipv4' or 'ipv6'

Required Flags (one of):
  --ip + --subnet-id   Pin this IP from this subnet (manual allocation)
  --subnet-pool-id     Allocate from these pools, in order (auto allocation)
  --config-source      'pipe' to read from stdin, or a path to a JSON/YAML file

Optional Flags:
  --mac-address   MAC address to match, as six colon-separated hex octets
  --circuit-id    Option 82 circuit id to match
  --device-type   Device type guard/filter
  --vendor        Vendor guard/filter
  --giaddr        Relay gateway address (giaddr) to match
  --priority      Tiebreaker among equally specific matches (lower wins)

Examples:
  metalcloud-cli dhcp-reservation create 1 ipv4 --mac-address AA:BB:CC:DD:EE:FF --ip 192.168.1.10 --subnet-id 3
  metalcloud-cli dhcp-reservation create 1 ipv4 --circuit-id leaf-01:swp1 --subnet-pool-id 3 --subnet-pool-id 4
  metalcloud-cli dhcp-reservation create 1 ipv4 --config-source reservation.yaml`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SITES_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var createConfig sdk.DhcpReservationCreate

			if dhcpReservationFlags.configSource != "" {
				config, err := utils.ReadConfigFromPipeOrFile(dhcpReservationFlags.configSource)
				if err != nil {
					return err
				}
				if err := utils.UnmarshalContent(config, &createConfig); err != nil {
					return err
				}
			} else {
				if len(dhcpReservationFlags.subnetPoolId) > 0 {
					subnetPoolIds := make([]float32, 0, len(dhcpReservationFlags.subnetPoolId))
					for _, poolId := range dhcpReservationFlags.subnetPoolId {
						subnetPoolIds = append(subnetPoolIds, float32(poolId))
					}
					createConfig.Allocation = sdk.DhcpAutoAllocationAsDhcpReservationAllocation(&sdk.DhcpAutoAllocation{
						Kind:          "auto",
						SubnetPoolIds: subnetPoolIds,
					})
				} else {
					createConfig.Allocation = sdk.DhcpManualAllocationAsDhcpReservationAllocation(&sdk.DhcpManualAllocation{
						Kind:     "manual",
						Ip:       dhcpReservationFlags.ip,
						SubnetId: int64(dhcpReservationFlags.subnetId),
					})
				}

				if dhcpReservationFlags.macAddress != "" {
					createConfig.MacAddress = &dhcpReservationFlags.macAddress
				}
				if dhcpReservationFlags.circuitId != "" {
					createConfig.CircuitId = &dhcpReservationFlags.circuitId
				}
				if dhcpReservationFlags.deviceType != "" {
					createConfig.DeviceType = &dhcpReservationFlags.deviceType
				}
				if dhcpReservationFlags.vendor != "" {
					createConfig.Vendor = &dhcpReservationFlags.vendor
				}
				if dhcpReservationFlags.giaddr != "" {
					createConfig.Giaddr = &dhcpReservationFlags.giaddr
				}
				if dhcpReservationFlags.priority != 0 {
					createConfig.Priority = sdk.PtrInt32(int32(dhcpReservationFlags.priority))
				}
			}

			return dhcp_reservation.DhcpReservationCreate(cmd.Context(), args[0], args[1], createConfig)
		},
	}

	dhcpReservationUpdateCmd = &cobra.Command{
		Use:     "update site_id_or_label ip_version reservation_id",
		Aliases: []string{"edit"},
		Short:   "Update a DHCP reservation",
		Long: `Update an existing DHCP reservation.

The update replaces the reservation (the endpoint is a PUT), so the
configuration must describe the reservation in full, allocation included. The
reservation's current revision is sent as the If-Match entity tag, so a
concurrent change is rejected instead of being overwritten.

Required Arguments:
  site_id_or_label  The ID or label of the site
  ip_version        'ipv4' or 'ipv6'
  reservation_id    The ID of the reservation

Required Flags:
  --config-source  'pipe' to read from stdin, or a path to a JSON/YAML file

Examples:
  metalcloud-cli dhcp-reservation update 1 ipv4 12 --config-source reservation.yaml
  metalcloud-cli dhcp-reservation get 1 ipv4 12 -f json > r.json
  metalcloud-cli dhcp-reservation update 1 ipv4 12 --config-source r.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SITES_WRITE},
		Args:         cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(dhcpReservationFlags.configSource)
			if err != nil {
				return err
			}

			return dhcp_reservation.DhcpReservationUpdate(cmd.Context(), args[0], args[1], args[2], config)
		},
	}

	dhcpReservationDeleteCmd = &cobra.Command{
		Use:     "delete site_id_or_label ip_version reservation_id",
		Aliases: []string{"rm", "del"},
		Short:   "Delete a DHCP reservation",
		Long: `Delete a DHCP reservation from a site.

The reservation's current revision is sent as the If-Match entity tag, so a
concurrent change is rejected instead of being overwritten.

Required Arguments:
  site_id_or_label  The ID or label of the site
  ip_version        'ipv4' or 'ipv6'
  reservation_id    The ID of the reservation

Examples:
  metalcloud-cli dhcp-reservation delete 1 ipv4 12
  metalcloud-cli dhcp-reservation rm dc-1 ipv6 13`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SITES_WRITE},
		Args:         cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return dhcp_reservation.DhcpReservationDelete(cmd.Context(), args[0], args[1], args[2])
		},
	}
)

func init() {
	rootCmd.AddCommand(dhcpReservationCmd)

	dhcpReservationCmd.AddCommand(dhcpReservationListCmd)
	dhcpReservationCmd.AddCommand(dhcpReservationGetCmd)
	dhcpReservationCmd.AddCommand(dhcpReservationConfigExampleCmd)

	dhcpReservationCmd.AddCommand(dhcpReservationCreateCmd)
	dhcpReservationCreateCmd.Flags().StringVar(&dhcpReservationFlags.configSource, "config-source", "", "Source of the new DHCP reservation configuration. Can be 'pipe' or path to a JSON/YAML file.")
	dhcpReservationCreateCmd.Flags().StringVar(&dhcpReservationFlags.macAddress, "mac-address", "", "MAC address to match, as six colon-separated hex octets.")
	dhcpReservationCreateCmd.Flags().StringVar(&dhcpReservationFlags.circuitId, "circuit-id", "", "Option 82 circuit id to match.")
	dhcpReservationCreateCmd.Flags().StringVar(&dhcpReservationFlags.deviceType, "device-type", "", "Device type guard/filter.")
	dhcpReservationCreateCmd.Flags().StringVar(&dhcpReservationFlags.vendor, "vendor", "", "Vendor guard/filter.")
	dhcpReservationCreateCmd.Flags().StringVar(&dhcpReservationFlags.giaddr, "giaddr", "", "Relay gateway address (giaddr) to match.")
	dhcpReservationCreateCmd.Flags().IntVar(&dhcpReservationFlags.priority, "priority", 0, "Tiebreaker among equally specific matches (lower = higher precedence).")
	dhcpReservationCreateCmd.Flags().StringVar(&dhcpReservationFlags.ip, "ip", "", "Manual allocation: the IP address to pin for this reservation.")
	dhcpReservationCreateCmd.Flags().IntVar(&dhcpReservationFlags.subnetId, "subnet-id", 0, "Manual allocation: the OOB subnet id that contains the pinned IP.")
	dhcpReservationCreateCmd.Flags().IntSliceVar(&dhcpReservationFlags.subnetPoolId, "subnet-pool-id", nil, "Auto allocation: OOB subnet pool ids to allocate the next free IP from, in order.")
	dhcpReservationCreateCmd.MarkFlagsOneRequired("config-source", "ip", "subnet-pool-id")
	dhcpReservationCreateCmd.MarkFlagsRequiredTogether("ip", "subnet-id")
	dhcpReservationCreateCmd.MarkFlagsMutuallyExclusive("config-source", "ip")
	dhcpReservationCreateCmd.MarkFlagsMutuallyExclusive("config-source", "subnet-id")
	dhcpReservationCreateCmd.MarkFlagsMutuallyExclusive("config-source", "subnet-pool-id")
	dhcpReservationCreateCmd.MarkFlagsMutuallyExclusive("ip", "subnet-pool-id")

	dhcpReservationCmd.AddCommand(dhcpReservationUpdateCmd)
	dhcpReservationUpdateCmd.Flags().StringVar(&dhcpReservationFlags.configSource, "config-source", "", "Source of the DHCP reservation configuration. Can be 'pipe' or path to a JSON/YAML file.")
	dhcpReservationUpdateCmd.MarkFlagsOneRequired("config-source")

	dhcpReservationCmd.AddCommand(dhcpReservationDeleteCmd)
}
