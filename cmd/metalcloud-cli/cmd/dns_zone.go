package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/dns_zone"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
	"github.com/spf13/cobra"
)

// DNS Zone commands
var (
	dnsZoneFlags = struct {
		configSource  string
		filterDefault []string
		// Create command fields
		description string
		zoneName    string
		zoneType    string
		soaEmail    string
		ttl         int
		isDefault   bool
		nameServers []string
		tags        []string
	}{}

	dnsZoneCmd = &cobra.Command{
		Use:     "dns-zone [command]",
		Aliases: []string{"dns", "zone"},
		Short:   "DNS Zone management",
		Long: `DNS Zone management commands.

This command group provides comprehensive DNS zone management capabilities including
creation, retrieval, updating, and deletion of DNS zones. DNS zones can be
managed individually with their associated record sets.

Available command categories:
  - Basic operations: list, get, create, update, delete
  - Record management: list-records, get-record
  - Information: nameservers

Use "metalcloud-cli dns-zone [command] --help" for detailed information about each command.
`,
	}

	dnsZoneListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List DNS zones",
		Long: `List all DNS zones in the MetalSoft infrastructure.

This command displays information about all DNS zones including their IDs, labels, 
zone names, zone types, status, and other configuration details.

Optional Flags:
  --filter-default    Filter zones by default status (true/false)

Examples:
  # List all DNS zones
  metalcloud-cli dns-zone list

  # Filter zones by status
  metalcloud-cli dns-zone list --filter-default true
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_GLOBAL_CONFIGURATIONS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return dns_zone.DNSZoneList(cmd.Context(), dnsZoneFlags.filterDefault)
		},
	}

	dnsZoneGetCmd = &cobra.Command{
		Use:     "get dns_zone_id",
		Aliases: []string{"show"},
		Short:   "Get detailed DNS zone information",
		Long: `Get detailed information for a specific DNS zone.

This command retrieves comprehensive information about a DNS zone including its
configuration, status, name servers, and other metadata.

Required Arguments:
  dns_zone_id           The ID of the DNS zone to retrieve information for

Examples:
  # Get DNS zone information
  metalcloud-cli dns-zone get 123

  # Get DNS zone information using alias
  metalcloud-cli dns-zone show 123
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_GLOBAL_CONFIGURATIONS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return dns_zone.DNSZoneGet(cmd.Context(), args[0])
		},
	}

	dnsZoneCreateCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Short:   "Create a new DNS zone",
		Long: `Create a new DNS zone in MetalSoft.

You can provide the DNS zone configuration either via command-line flags or by 
specifying a configuration source using the --config-source flag. The configuration 
source can be a path to a JSON file or 'pipe' to read from standard input.

If --config-source is not provided, you must specify at least --zone-name, 
--is-default, and --name-servers, along with any other relevant zone details.

Required Flags (when not using --config-source):
  --zone-name           DNS zone name (without terminating dot)
  --is-default          Whether this is the default DNS zone
  --name-servers        List of name servers (comma-separated)

Optional Flags:
  --config-source       Source of DNS zone configuration (JSON file path or 'pipe')
  --description         DNS zone description
  --zone-type           Zone type (master/slave, default: master)
  --soa-email          Email address of DNS zone administrator
  --ttl                TTL (Time to Live) for the DNS zone
  --tags               Tags for the DNS zone (comma-separated)

Examples:
  # Create using command line flags
  metalcloud-cli dns-zone create --zone-name example.com --is-default true --name-servers ns1.example.com,ns2.example.com

  # Create with additional details
  metalcloud-cli dns-zone create --zone-name test.com --is-default false --name-servers ns1.test.com --description "Test zone" --zone-type master --ttl 300

  # Create using JSON configuration file
  metalcloud-cli dns-zone create --config-source ./zone.json

  # Create using piped JSON configuration
  cat zone.json | metalcloud-cli dns-zone create --config-source pipe
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_GLOBAL_CONFIGURATIONS_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			var zoneConfig sdk.CreateDnsZoneDto

			// If config source is provided, use it
			if dnsZoneFlags.configSource != "" {
				config, err := utils.ReadConfigFromPipeOrFile(dnsZoneFlags.configSource)
				if err != nil {
					return err
				}
				err = utils.UnmarshalContent(config, &zoneConfig)
				if err != nil {
					return err
				}
			} else {
				// Otherwise build config from command line parameters
				zoneConfig = sdk.CreateDnsZoneDto{
					ZoneName:    dnsZoneFlags.zoneName,
					IsDefault:   dnsZoneFlags.isDefault,
					NameServers: dnsZoneFlags.nameServers,
				}

				if dnsZoneFlags.description != "" {
					zoneConfig.Description = &dnsZoneFlags.description
				}
				if dnsZoneFlags.zoneType != "" {
					zoneConfig.ZoneType = &dnsZoneFlags.zoneType
				}
				if dnsZoneFlags.soaEmail != "" {
					zoneConfig.SoaEmail = &dnsZoneFlags.soaEmail
				}
				if dnsZoneFlags.ttl > 0 {
					ttl32 := int32(dnsZoneFlags.ttl)
					zoneConfig.Ttl = &ttl32
				}
				if len(dnsZoneFlags.tags) > 0 {
					zoneConfig.Tags = dnsZoneFlags.tags
				}
			}

			return dns_zone.DNSZoneCreate(cmd.Context(), zoneConfig)
		},
	}

	dnsZoneUpdateCmd = &cobra.Command{
		Use:   "update dns_zone_id",
		Short: "Update DNS zone information",
		Long: `Update DNS zone information.

This command updates DNS zone configuration using a JSON configuration file or 
piped JSON data. The configuration must be provided via the --config-source flag.

Required Arguments:
  dns_zone_id           The ID of the DNS zone to update

Required Flags:
  --config-source       Source of the DNS zone update configuration. Can be 'pipe' or path to a JSON file

Examples:
  # Update DNS zone using JSON configuration file
  metalcloud-cli dns-zone update 123 --config-source ./zone-update.json

  # Update DNS zone using piped JSON configuration
  echo '{"description": "Updated description"}' | metalcloud-cli dns-zone update 123 --config-source pipe
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_GLOBAL_CONFIGURATIONS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(dnsZoneFlags.configSource)
			if err != nil {
				return err
			}
			return dns_zone.DNSZoneUpdate(cmd.Context(), args[0], config)
		},
	}

	dnsZoneDeleteCmd = &cobra.Command{
		Use:     "delete dns_zone_id",
		Aliases: []string{"rm", "remove"},
		Short:   "Delete a DNS zone",
		Long: `Delete a DNS zone from MetalSoft infrastructure.

This command permanently deletes a DNS zone and all its associated DNS records.
This action cannot be undone, so use with caution.

Required Arguments:
  dns_zone_id           The ID of the DNS zone to delete

Examples:
  # Delete a DNS zone
  metalcloud-cli dns-zone delete 123

  # Delete a DNS zone using alias
  metalcloud-cli dns-zone rm 123
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_GLOBAL_CONFIGURATIONS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return dns_zone.DNSZoneDelete(cmd.Context(), args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(dnsZoneCmd)

	// DNS Zone commands
	dnsZoneCmd.AddCommand(dnsZoneListCmd)
	dnsZoneListCmd.Flags().StringSliceVar(&dnsZoneFlags.filterDefault, "filter-default", nil, "Filter the result by default status.")

	dnsZoneCmd.AddCommand(dnsZoneGetCmd)

	dnsZoneCmd.AddCommand(dnsZoneCreateCmd)
	dnsZoneCreateCmd.Flags().StringVar(&dnsZoneFlags.configSource, "config-source", "", "Source of the new DNS zone configuration. Can be 'pipe' or path to a JSON file.")
	dnsZoneCreateCmd.Flags().StringVar(&dnsZoneFlags.description, "description", "", "DNS zone description")
	dnsZoneCreateCmd.Flags().StringVar(&dnsZoneFlags.zoneName, "zone-name", "", "DNS zone name (without terminating dot)")
	dnsZoneCreateCmd.Flags().StringVar(&dnsZoneFlags.zoneType, "zone-type", "master", "Zone type (master/slave)")
	dnsZoneCreateCmd.Flags().StringVar(&dnsZoneFlags.soaEmail, "soa-email", "", "Email address of DNS zone administrator")
	dnsZoneCreateCmd.Flags().IntVar(&dnsZoneFlags.ttl, "ttl", 0, "TTL (Time to Live) for the DNS zone")
	dnsZoneCreateCmd.Flags().BoolVar(&dnsZoneFlags.isDefault, "is-default", false, "Whether this is the default DNS zone")
	dnsZoneCreateCmd.Flags().StringSliceVar(&dnsZoneFlags.nameServers, "name-servers", nil, "Name servers for the DNS zone")
	dnsZoneCreateCmd.Flags().StringSliceVar(&dnsZoneFlags.tags, "tags", nil, "Tags for the DNS zone")
	dnsZoneCreateCmd.MarkFlagsOneRequired("config-source", "zone-name")
	dnsZoneCreateCmd.MarkFlagsMutuallyExclusive("config-source", "zone-name")
	dnsZoneCreateCmd.MarkFlagsRequiredTogether("zone-name", "name-servers")

	dnsZoneCmd.AddCommand(dnsZoneUpdateCmd)
	dnsZoneUpdateCmd.Flags().StringVar(&dnsZoneFlags.configSource, "config-source", "", "Source of the DNS zone update configuration. Can be 'pipe' or path to a JSON file.")
	dnsZoneUpdateCmd.MarkFlagsOneRequired("config-source")

	dnsZoneCmd.AddCommand(dnsZoneDeleteCmd)
}
