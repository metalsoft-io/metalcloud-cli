package site

import (
	"context"
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

type DHCPOOBReservation struct {
	MACAddress string
	IPAddress  string
}

var dhcpOOBReservationPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"MACAddress": {
			Title: "MAC Address",
			Order: 1,
		},
		"IPAddress": {
			Title: "IP Address",
			Order: 2,
		},
	},
}

// DhcpOobReservationsList retrieves and displays the DHCP Option82 to IP mapping from the site config.
func DhcpOobReservationsList(ctx context.Context, siteIdOrName string, sortBy string) error {
	logger.Get().Info().Msgf("Listing DHCP OOB reservations for site '%s'", siteIdOrName)

	mapping, err := getDhcpOption82Mapping(ctx, siteIdOrName)
	if err != nil {
		return err
	}

	reservations := mapToReservations(mapping)

	switch strings.ToLower(sortBy) {
	case "ip":
		sort.Slice(reservations, func(i, j int) bool {
			return reservations[i].IPAddress < reservations[j].IPAddress
		})
	default:
		sort.Slice(reservations, func(i, j int) bool {
			return reservations[i].MACAddress < reservations[j].MACAddress
		})
	}

	return formatter.PrintResult(reservations, &dhcpOOBReservationPrintConfig)
}

// DhcpOobReservationsAdd adds one or more MAC-to-IP entries to the DHCP Option82 mapping.
func DhcpOobReservationsAdd(ctx context.Context, siteIdOrName string, entries map[string]string) error {
	logger.Get().Info().Msgf("Adding DHCP OOB reservations for site '%s'", siteIdOrName)

	siteInfo, err := GetSiteByIdOrLabel(ctx, siteIdOrName)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	siteConfig, httpRes, err := client.SiteAPI.GetSiteConfig(ctx, float32(siteInfo.Id)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	mapping := getOrInitMapping(siteConfig)

	// Check for duplicate MACs
	for mac := range entries {
		normalizedMAC := strings.ToLower(mac)
		for existingMAC := range mapping {
			if strings.ToLower(existingMAC) == normalizedMAC {
				return fmt.Errorf("MAC address '%s' already exists in DHCP OOB reservations", mac)
			}
		}
	}

	// Check for duplicate IPs within the new entries
	if err := checkDuplicateIPs(entries); err != nil {
		return err
	}

	// Check new IPs don't collide with existing reservations
	for mac, ip := range entries {
		for existingMAC, existingIP := range mapping {
			if fmt.Sprintf("%v", existingIP) == ip {
				return fmt.Errorf("IP address '%s' (for MAC '%s') is already reserved by MAC '%s'", ip, mac, existingMAC)
			}
		}
	}

	// Validate IPs belong to OOB subnets
	if err := validateIPsInOOBSubnets(ctx, client, siteInfo, entries); err != nil {
		return err
	}

	for mac, ip := range entries {
		mapping[mac] = ip
	}

	return updateDhcpOption82Mapping(ctx, client, siteInfo, mapping)
}

// DhcpOobReservationsRemove removes one or more MAC entries from the DHCP Option82 mapping.
func DhcpOobReservationsRemove(ctx context.Context, siteIdOrName string, macAddresses []string) error {
	logger.Get().Info().Msgf("Removing DHCP OOB reservations for site '%s'", siteIdOrName)

	siteInfo, err := GetSiteByIdOrLabel(ctx, siteIdOrName)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	siteConfig, httpRes, err := client.SiteAPI.GetSiteConfig(ctx, float32(siteInfo.Id)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	mapping := getOrInitMapping(siteConfig)

	for _, mac := range macAddresses {
		found := false
		normalizedMAC := strings.ToLower(mac)
		for existingMAC := range mapping {
			if strings.ToLower(existingMAC) == normalizedMAC {
				delete(mapping, existingMAC)
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("MAC address '%s' not found in DHCP OOB reservations", mac)
		}
	}

	return updateDhcpOption82Mapping(ctx, client, siteInfo, mapping)
}

// DhcpOobReservationsReplace replaces the entire DHCP Option82 mapping with the provided entries.
func DhcpOobReservationsReplace(ctx context.Context, siteIdOrName string, entries map[string]string) error {
	logger.Get().Info().Msgf("Replacing DHCP OOB reservations for site '%s'", siteIdOrName)

	siteInfo, err := GetSiteByIdOrLabel(ctx, siteIdOrName)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	// Check for duplicate IPs within the new entries
	if err := checkDuplicateIPs(entries); err != nil {
		return err
	}

	// Validate IPs belong to OOB subnets
	if err := validateIPsInOOBSubnets(ctx, client, siteInfo, entries); err != nil {
		return err
	}

	mapping := make(map[string]interface{}, len(entries))
	for mac, ip := range entries {
		mapping[mac] = ip
	}

	return updateDhcpOption82Mapping(ctx, client, siteInfo, mapping)
}

func getDhcpOption82Mapping(ctx context.Context, siteIdOrName string) (map[string]interface{}, error) {
	siteInfo, err := GetSiteByIdOrLabel(ctx, siteIdOrName)
	if err != nil {
		return nil, err
	}

	client := api.GetApiClient(ctx)

	siteConfig, httpRes, err := client.SiteAPI.GetSiteConfig(ctx, float32(siteInfo.Id)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, err
	}

	serverPolicy, ok := siteConfig.GetServerPolicyOk()
	if !ok || serverPolicy == nil {
		return map[string]interface{}{}, nil
	}

	return serverPolicy.GetDhcpOption82ToIPMapping(), nil
}

func getOrInitMapping(siteConfig *sdk.SiteConfig) map[string]interface{} {
	serverPolicy, ok := siteConfig.GetServerPolicyOk()
	if !ok || serverPolicy == nil {
		return make(map[string]interface{})
	}
	mapping := serverPolicy.GetDhcpOption82ToIPMapping()
	if mapping == nil {
		return make(map[string]interface{})
	}
	return mapping
}

func mapToReservations(mapping map[string]interface{}) []DHCPOOBReservation {
	reservations := make([]DHCPOOBReservation, 0, len(mapping))
	for mac, ip := range mapping {
		ipStr := fmt.Sprintf("%v", ip)
		reservations = append(reservations, DHCPOOBReservation{
			MACAddress: mac,
			IPAddress:  ipStr,
		})
	}
	return reservations
}

func updateDhcpOption82Mapping(ctx context.Context, client *sdk.APIClient, siteInfo *sdk.Site, mapping map[string]interface{}) error {
	currentSite, httpRes, err := client.SiteAPI.GetSite(ctx, float32(siteInfo.Id)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	configUpdate := sdk.SiteConfigUpdate{
		ServerPolicy: &sdk.ServerPolicyUpdate{
			DhcpOption82ToIPMapping: mapping,
		},
	}

	updatedConfig, httpRes, err := client.SiteAPI.UpdateSiteConfig(ctx, float32(siteInfo.Id)).
		SiteConfigUpdate(configUpdate).
		IfMatch(strconv.Itoa(int(currentSite.Revision))).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	resultMapping := getOrInitMapping(updatedConfig)
	reservations := mapToReservations(resultMapping)

	sort.Slice(reservations, func(i, j int) bool {
		return reservations[i].MACAddress < reservations[j].MACAddress
	})

	return formatter.PrintResult(reservations, &dhcpOOBReservationPrintConfig)
}

func validateIPsInOOBSubnets(ctx context.Context, client *sdk.APIClient, siteInfo *sdk.Site, entries map[string]string) error {
	oobSubnets, err := getOOBSubnetsForSite(ctx, client, siteInfo)
	if err != nil {
		return fmt.Errorf("failed to retrieve OOB subnets for site: %w", err)
	}

	if len(oobSubnets) == 0 {
		return fmt.Errorf("no OOB subnets found for site '%s'; cannot validate IP addresses", siteInfo.Name)
	}

	for mac, ip := range entries {
		parsedIP := net.ParseIP(ip)
		if parsedIP == nil {
			return fmt.Errorf("invalid IP address '%s' for MAC '%s'", ip, mac)
		}

		if !ipBelongsToAnySubnet(parsedIP, oobSubnets) {
			return fmt.Errorf("IP address '%s' for MAC '%s' does not belong to any OOB subnet assigned to site '%s'", ip, mac, siteInfo.Name)
		}
	}

	return nil
}

func getOOBSubnetsForSite(ctx context.Context, client *sdk.APIClient, siteInfo *sdk.Site) ([]*net.IPNet, error) {
	siteIdStr := strconv.Itoa(int(siteInfo.Id))

	subnets, httpRes, err := client.SubnetAPI.GetSubnets(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, err
	}

	var oobSubnets []*net.IPNet
	for _, subnet := range subnets.Data {
		tags := subnet.GetTags()
		if tags[tagOOB] == "" {
			continue
		}
		if tags[tagSiteID] != siteIdStr {
			continue
		}

		_, ipNet, err := net.ParseCIDR(fmt.Sprintf("%s/%d", subnet.NetworkAddress, subnet.PrefixLength))
		if err != nil {
			logger.Get().Warn().Msgf("Failed to parse subnet CIDR %s/%d: %v", subnet.NetworkAddress, subnet.PrefixLength, err)
			continue
		}

		oobSubnets = append(oobSubnets, ipNet)
	}

	return oobSubnets, nil
}

const (
	tagOOB    = "metalcloud/oob"
	tagSiteID = "metalcloud/site-id"
)

func checkDuplicateIPs(entries map[string]string) error {
	seen := make(map[string]string, len(entries))
	for mac, ip := range entries {
		if existingMAC, exists := seen[ip]; exists {
			return fmt.Errorf("duplicate IP address '%s' assigned to both MAC '%s' and MAC '%s'", ip, existingMAC, mac)
		}
		seen[ip] = mac
	}
	return nil
}

func ipBelongsToAnySubnet(ip net.IP, subnets []*net.IPNet) bool {
	for _, subnet := range subnets {
		if subnet.Contains(ip) {
			return true
		}
	}
	return false
}
