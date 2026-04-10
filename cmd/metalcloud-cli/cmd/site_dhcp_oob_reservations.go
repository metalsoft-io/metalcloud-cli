package cmd

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/site"
	"github.com/spf13/cobra"
)

var (
	dhcpOobReservationsFlags = struct {
		sortBy     string
		mac        string
		ip         string
		jsonSource string
		fromFile   string
	}{}

	siteDhcpOobReservationsCmd = &cobra.Command{
		Use:     "dhcp-oob-reservations [command]",
		Aliases: []string{"dhcp-oob"},
		Short:   "Manage DHCP Option82 OOB IP reservations for a site",
		Long: `Manage DHCP Option82 to IP address mappings in the site's server policy configuration.

These reservations map MAC addresses (DHCP Option82) to static IP addresses for
out-of-band (OOB) management interfaces. The mappings are stored in the site
configuration under serverPolicy.dhcpOption82ToIPMapping.

Available Commands:
  list      List all DHCP OOB reservations for a site
  add       Add a MAC-to-IP reservation entry
  remove    Remove a reservation entry by MAC address
  replace   Replace all reservation entries from JSON input`,
	}

	siteDhcpOobReservationsListCmd = &cobra.Command{
		Use:     "list site_id_or_name",
		Aliases: []string{"ls"},
		Short:   "List all DHCP OOB reservations for a site",
		Long: `List all DHCP Option82 to IP address mappings configured for a site.

Displays a table of MAC address to IP address reservations from the site's
serverPolicy.dhcpOption82ToIPMapping configuration.

Examples:
  # List reservations sorted by MAC address (default)
  metalcloud-cli site dhcp-oob-reservations list site-01

  # List reservations sorted by IP address
  metalcloud-cli site dhcp-oob-reservations list site-01 --sort-by ip

  # List in JSON format
  metalcloud-cli site dhcp-oob-reservations list site-01 --format json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SITES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return site.DhcpOobReservationsList(cmd.Context(), args[0], dhcpOobReservationsFlags.sortBy)
		},
	}

	siteDhcpOobReservationsAddCmd = &cobra.Command{
		Use:   "add site_id_or_name",
		Short: "Add a DHCP OOB reservation entry",
		Long: `Add a new MAC address to IP address mapping to the site's DHCP OOB reservations.

The IP address must belong to an OOB subnet assigned to the site. The MAC address
must be in standard format (e.g., AA:BB:CC:DD:EE:FF or aa:bb:cc:dd:ee:ff).

Entries can be added one at a time using --mac and --ip flags, or in bulk
using --from-json with a JSON string or --from-file with a path to a JSON file.
The JSON format is an object with MAC addresses as keys and IP addresses as values.

Examples:
  # Add a single reservation
  metalcloud-cli site dhcp-oob-reservations add site-01 --mac AA:BB:CC:DD:EE:FF --ip 10.0.0.100

  # Add multiple reservations from JSON string
  metalcloud-cli site dhcp-oob-reservations add site-01 --from-json '{"AA:BB:CC:DD:EE:FF":"10.0.0.100","11:22:33:44:55:66":"10.0.0.101"}'

  # Add multiple reservations from a JSON file
  metalcloud-cli site dhcp-oob-reservations add site-01 --from-file reservations.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SITES_WRITE},
		Args:         cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			hasMac := cmd.Flags().Changed("mac")
			hasIP := cmd.Flags().Changed("ip")
			hasJSON := cmd.Flags().Changed("from-json")
			hasFile := cmd.Flags().Changed("from-file")

			bulkSources := countTrue(hasJSON, hasFile)
			if bulkSources > 1 {
				return fmt.Errorf("--from-json and --from-file are mutually exclusive")
			}
			if bulkSources > 0 && (hasMac || hasIP) {
				return fmt.Errorf("--from-json/--from-file cannot be used with --mac or --ip")
			}
			if bulkSources == 0 && (!hasMac || !hasIP) {
				return fmt.Errorf("either --mac and --ip, --from-json, or --from-file must be provided")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := parseDhcpEntries(cmd)
			if err != nil {
				return err
			}
			return site.DhcpOobReservationsAdd(cmd.Context(), args[0], entries)
		},
	}

	siteDhcpOobReservationsRemoveCmd = &cobra.Command{
		Use:   "remove site_id_or_name",
		Short: "Remove a DHCP OOB reservation entry by MAC address",
		Long: `Remove one or more MAC address entries from the site's DHCP OOB reservations.

MAC addresses can be specified with --mac flags, or loaded from a JSON file via
--from-file. The file should contain either a JSON array of MAC address strings
or a JSON object whose keys are MAC addresses.

Examples:
  # Remove a single reservation
  metalcloud-cli site dhcp-oob-reservations remove site-01 --mac AA:BB:CC:DD:EE:FF

  # Remove multiple reservations
  metalcloud-cli site dhcp-oob-reservations remove site-01 --mac AA:BB:CC:DD:EE:FF --mac 11:22:33:44:55:66

  # Remove reservations listed in a file (array format)
  metalcloud-cli site dhcp-oob-reservations remove site-01 --from-file remove-macs.json

  # Remove reservations listed in a file (object format - keys are used)
  metalcloud-cli site dhcp-oob-reservations remove site-01 --from-file reservations.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SITES_WRITE},
		Args:         cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			hasMac := cmd.Flags().Changed("mac")
			hasFile := cmd.Flags().Changed("from-file")

			if hasMac && hasFile {
				return fmt.Errorf("--mac and --from-file are mutually exclusive")
			}
			if !hasMac && !hasFile {
				return fmt.Errorf("either --mac or --from-file must be provided")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			macs, err := parseRemoveMACs(cmd)
			if err != nil {
				return err
			}
			return site.DhcpOobReservationsRemove(cmd.Context(), args[0], macs)
		},
	}

	siteDhcpOobReservationsReplaceCmd = &cobra.Command{
		Use:   "replace site_id_or_name",
		Short: "Replace all DHCP OOB reservation entries",
		Long: `Replace the entire DHCP Option82 to IP address mapping with the provided entries.

This removes all existing reservations and sets the mapping to the provided JSON
object. All IP addresses must belong to OOB subnets assigned to the site.
Input can be a JSON string via --from-json or a JSON file via --from-file.

Examples:
  # Replace all reservations from JSON string
  metalcloud-cli site dhcp-oob-reservations replace site-01 --from-json '{"AA:BB:CC:DD:EE:FF":"10.0.0.100","11:22:33:44:55:66":"10.0.0.101"}'

  # Replace all reservations from a file
  metalcloud-cli site dhcp-oob-reservations replace site-01 --from-file reservations.json

  # Clear all reservations
  metalcloud-cli site dhcp-oob-reservations replace site-01 --from-json '{}'`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SITES_WRITE},
		Args:         cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			hasJSON := cmd.Flags().Changed("from-json")
			hasFile := cmd.Flags().Changed("from-file")

			if hasJSON && hasFile {
				return fmt.Errorf("--from-json and --from-file are mutually exclusive")
			}
			if !hasJSON && !hasFile {
				return fmt.Errorf("either --from-json or --from-file must be provided")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := parseDhcpEntries(cmd)
			if err != nil {
				return err
			}
			return site.DhcpOobReservationsReplace(cmd.Context(), args[0], entries)
		},
	}
)

var macRegex = regexp.MustCompile(`^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`)

func isValidMAC(mac string) bool {
	return macRegex.MatchString(mac)
}

func countTrue(vals ...bool) int {
	n := 0
	for _, v := range vals {
		if v {
			n++
		}
	}
	return n
}

func readJSONFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file '%s': %w", path, err)
	}
	return data, nil
}

func parseMACToIPJSON(data []byte) (map[string]string, error) {
	var entries map[string]string
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("invalid JSON: expected an object with MAC-to-IP mappings: %w", err)
	}

	for mac, ip := range entries {
		if !isValidMAC(mac) {
			return nil, fmt.Errorf("invalid MAC address format: '%s'", mac)
		}
		if net.ParseIP(ip) == nil {
			return nil, fmt.Errorf("invalid IP address '%s' for MAC '%s'", ip, mac)
		}
	}

	return entries, nil
}

// parseDhcpEntries resolves MAC-to-IP entries from --mac/--ip, --from-json, or --from-file.
func parseDhcpEntries(cmd *cobra.Command) (map[string]string, error) {
	if cmd.Flags().Changed("from-file") {
		data, err := readJSONFile(dhcpOobReservationsFlags.fromFile)
		if err != nil {
			return nil, err
		}
		return parseMACToIPJSON(data)
	}

	if cmd.Flags().Changed("from-json") {
		return parseMACToIPJSON([]byte(dhcpOobReservationsFlags.jsonSource))
	}

	mac := strings.TrimSpace(dhcpOobReservationsFlags.mac)
	ip := strings.TrimSpace(dhcpOobReservationsFlags.ip)

	if !isValidMAC(mac) {
		return nil, fmt.Errorf("invalid MAC address format: '%s'", mac)
	}
	if net.ParseIP(ip) == nil {
		return nil, fmt.Errorf("invalid IP address: '%s'", ip)
	}

	return map[string]string{mac: ip}, nil
}

// parseRemoveMACs resolves MAC addresses from --mac flags or --from-file.
// The file can contain a JSON array of MAC strings or a JSON object whose keys are MACs.
func parseRemoveMACs(cmd *cobra.Command) ([]string, error) {
	if cmd.Flags().Changed("from-file") {
		filePath, _ := cmd.Flags().GetString("from-file")
		data, err := readJSONFile(filePath)
		if err != nil {
			return nil, err
		}

		// Try array first
		var macList []string
		if err := json.Unmarshal(data, &macList); err == nil {
			for _, mac := range macList {
				if !isValidMAC(mac) {
					return nil, fmt.Errorf("invalid MAC address format: '%s'", mac)
				}
			}
			return macList, nil
		}

		// Fall back to object (extract keys)
		var macMap map[string]interface{}
		if err := json.Unmarshal(data, &macMap); err != nil {
			return nil, fmt.Errorf("invalid JSON: expected an array of MAC addresses or an object with MAC keys: %w", err)
		}
		macs := make([]string, 0, len(macMap))
		for mac := range macMap {
			if !isValidMAC(mac) {
				return nil, fmt.Errorf("invalid MAC address format: '%s'", mac)
			}
			macs = append(macs, mac)
		}
		return macs, nil
	}

	macs, _ := cmd.Flags().GetStringSlice("mac")
	for _, mac := range macs {
		if !isValidMAC(mac) {
			return nil, fmt.Errorf("invalid MAC address format: '%s'", mac)
		}
	}
	return macs, nil
}

func init() {
	siteCmd.AddCommand(siteDhcpOobReservationsCmd)

	// List
	siteDhcpOobReservationsCmd.AddCommand(siteDhcpOobReservationsListCmd)
	siteDhcpOobReservationsListCmd.Flags().StringVar(&dhcpOobReservationsFlags.sortBy, "sort-by", "mac", "Sort results by field: mac or ip")

	// Add
	siteDhcpOobReservationsCmd.AddCommand(siteDhcpOobReservationsAddCmd)
	siteDhcpOobReservationsAddCmd.Flags().StringVar(&dhcpOobReservationsFlags.mac, "mac", "", "MAC address (e.g., AA:BB:CC:DD:EE:FF)")
	siteDhcpOobReservationsAddCmd.Flags().StringVar(&dhcpOobReservationsFlags.ip, "ip", "", "IP address to map to the MAC address")
	siteDhcpOobReservationsAddCmd.Flags().StringVar(&dhcpOobReservationsFlags.jsonSource, "from-json", "", "JSON object with MAC-to-IP mappings")
	siteDhcpOobReservationsAddCmd.Flags().StringVar(&dhcpOobReservationsFlags.fromFile, "from-file", "", "Path to a JSON file with MAC-to-IP mappings")

	// Remove
	siteDhcpOobReservationsCmd.AddCommand(siteDhcpOobReservationsRemoveCmd)
	siteDhcpOobReservationsRemoveCmd.Flags().StringSlice("mac", nil, "MAC address(es) to remove (can be specified multiple times)")
	siteDhcpOobReservationsRemoveCmd.Flags().String("from-file", "", "Path to a JSON file with MAC addresses (array or object keys)")

	// Replace
	siteDhcpOobReservationsCmd.AddCommand(siteDhcpOobReservationsReplaceCmd)
	siteDhcpOobReservationsReplaceCmd.Flags().StringVar(&dhcpOobReservationsFlags.jsonSource, "from-json", "", "JSON object with MAC-to-IP mappings")
	siteDhcpOobReservationsReplaceCmd.Flags().StringVar(&dhcpOobReservationsFlags.fromFile, "from-file", "", "Path to a JSON file with MAC-to-IP mappings")
}
