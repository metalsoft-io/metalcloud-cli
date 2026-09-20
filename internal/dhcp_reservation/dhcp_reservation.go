// Package dhcp_reservation implements the site-scoped DHCP reservation
// commands. A reservation lives under a site and an IP version
// (/api/v2/sites/{siteId}/dhcp/{ipVersion}/reservations), so every command
// takes both: the site by id or label, the IP version as "ipv4" or "ipv6".
package dhcp_reservation

import (
	"context"
	"fmt"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/internal/site"
	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var dhcpReservationPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"SiteId": {
			Title: "Site",
			Order: 2,
		},
		"IpVersion": {
			Title: "IP Version",
			Order: 3,
		},
		"MacAddress": {
			Title: "MAC Address",
			Order: 4,
		},
		"CircuitId": {
			Title: "Circuit ID",
			Order: 5,
		},
		"Giaddr": {
			Title: "Giaddr",
			Order: 6,
		},
		"DeviceType": {
			Title: "Device Type",
			Order: 7,
		},
		"Vendor": {
			Title: "Vendor",
			Order: 8,
		},
		"Priority": {
			Title: "Priority",
			Order: 9,
		},
		"Allocation": {
			Title:       "Allocation",
			MaxWidth:    60,
			Transformer: formatAllocationValue,
			Order:       10,
		},
		"Revision": {
			Title: "Revision",
			Order: 11,
		},
	},
}

// formatAllocationValue renders the polymorphic allocation (auto | manual) as
// a compact single-cell summary for table output.
func formatAllocationValue(value interface{}) string {
	allocation, ok := value.(sdk.DhcpReservationAllocation)
	if !ok {
		if pointer, isPointer := value.(*sdk.DhcpReservationAllocation); isPointer && pointer != nil {
			allocation = *pointer
		} else {
			return fmt.Sprint(value)
		}
	}

	switch {
	case allocation.DhcpManualAllocation != nil:
		manual := allocation.DhcpManualAllocation
		return fmt.Sprintf("manual ip=%s subnetId=%d", manual.Ip, manual.SubnetId)
	case allocation.DhcpAutoAllocation != nil:
		auto := allocation.DhcpAutoAllocation
		pools := make([]string, 0, len(auto.SubnetPoolIds))
		for _, poolId := range auto.SubnetPoolIds {
			pools = append(pools, strconv.FormatInt(int64(poolId), 10))
		}
		return fmt.Sprintf("auto subnetPoolIds=%v", pools)
	default:
		return ""
	}
}

func DhcpReservationList(ctx context.Context, siteIdOrLabel string, ipVersion string) error {
	logger.Get().Info().Msgf("Listing %s DHCP reservations of site '%s'", ipVersion, siteIdOrLabel)

	siteId, ipVersionValue, err := resolveSiteAndIpVersion(ctx, siteIdOrLabel, ipVersion)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.DHCPReservationAPI.
		GetDhcpReservations(ctx, siteId, ipVersionValue).
		SortBy([]string{"id:ASC"})

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &dhcpReservationPrintConfig)
}

func DhcpReservationGet(ctx context.Context, siteIdOrLabel string, ipVersion string, reservationId string) error {
	logger.Get().Info().Msgf("Get %s DHCP reservation '%s' of site '%s'", ipVersion, reservationId, siteIdOrLabel)

	reservation, _, _, err := getDhcpReservation(ctx, siteIdOrLabel, ipVersion, reservationId)
	if err != nil {
		return err
	}

	return formatter.PrintResult(reservation, &dhcpReservationPrintConfig)
}

func DhcpReservationCreate(ctx context.Context, siteIdOrLabel string, ipVersion string, create sdk.DhcpReservationCreate) error {
	logger.Get().Info().Msgf("Creating %s DHCP reservation in site '%s'", ipVersion, siteIdOrLabel)

	siteId, ipVersionValue, err := resolveSiteAndIpVersion(ctx, siteIdOrLabel, ipVersion)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	reservation, httpRes, err := client.DHCPReservationAPI.
		CreateDhcpReservation(ctx, siteId, ipVersionValue).
		DhcpReservationCreate(create).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(reservation, &dhcpReservationPrintConfig)
}

func DhcpReservationUpdate(ctx context.Context, siteIdOrLabel string, ipVersion string, reservationId string, config []byte) error {
	logger.Get().Info().Msgf("Updating %s DHCP reservation '%s' of site '%s'", ipVersion, reservationId, siteIdOrLabel)

	var update sdk.DhcpReservationUpdate
	if err := utils.UnmarshalContent(config, &update); err != nil {
		return err
	}

	current, siteId, reservationIdNumeric, err := getDhcpReservation(ctx, siteIdOrLabel, ipVersion, reservationId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	reservation, httpRes, err := client.DHCPReservationAPI.
		UpdateDhcpReservation(ctx, siteId, current.IpVersion, reservationIdNumeric).
		DhcpReservationUpdate(update).
		IfMatch(strconv.FormatInt(current.Revision, 10)).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(reservation, &dhcpReservationPrintConfig)
}

func DhcpReservationDelete(ctx context.Context, siteIdOrLabel string, ipVersion string, reservationId string) error {
	logger.Get().Info().Msgf("Deleting %s DHCP reservation '%s' of site '%s'", ipVersion, reservationId, siteIdOrLabel)

	current, siteId, reservationIdNumeric, err := getDhcpReservation(ctx, siteIdOrLabel, ipVersion, reservationId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.DHCPReservationAPI.
		DeleteDhcpReservation(ctx, siteId, current.IpVersion, reservationIdNumeric).
		IfMatch(strconv.FormatInt(current.Revision, 10)).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("DHCP reservation '%s' deleted from site '%s'", reservationId, siteIdOrLabel)
	return nil
}

// DhcpReservationConfigExample prints one example create body per allocation
// kind (auto and manual).
func DhcpReservationConfigExample(ctx context.Context) error {
	examples := []sdk.DhcpReservationCreate{
		{
			MacAddress: sdk.PtrString("AA:BB:CC:DD:EE:FF"),
			Priority:   sdk.PtrInt32(100),
			Allocation: sdk.DhcpManualAllocationAsDhcpReservationAllocation(&sdk.DhcpManualAllocation{
				Kind:     "manual",
				Ip:       "192.168.1.10",
				SubnetId: 1,
			}),
		},
		{
			CircuitId:  sdk.PtrString("leaf-01:swp1"),
			DeviceType: sdk.PtrString("server"),
			Allocation: sdk.DhcpAutoAllocationAsDhcpReservationAllocation(&sdk.DhcpAutoAllocation{
				Kind:          "auto",
				SubnetPoolIds: []float32{1},
			}),
		},
	}

	if formatter.IsNativeFormat() {
		return formatter.PrintResult(examples, nil)
	}
	return formatter.PrintYamlResult(examples)
}

// getDhcpReservation fetches one reservation and returns it together with the
// resolved numeric site and reservation ids.
func getDhcpReservation(ctx context.Context, siteIdOrLabel string, ipVersion string, reservationId string) (*sdk.DhcpReservation, int64, int64, error) {
	siteId, ipVersionValue, err := resolveSiteAndIpVersion(ctx, siteIdOrLabel, ipVersion)
	if err != nil {
		return nil, 0, 0, err
	}

	reservationIdNumeric, err := utils.GetInt64FromString(reservationId)
	if err != nil {
		err = fmt.Errorf("invalid DHCP reservation ID: '%s'", reservationId)
		logger.Get().Error().Err(err).Msg("")
		return nil, 0, 0, err
	}

	client := api.GetApiClient(ctx)

	reservation, httpRes, err := client.DHCPReservationAPI.
		GetDhcpReservation(ctx, siteId, ipVersionValue, reservationIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, 0, 0, err
	}

	return reservation, siteId, reservationIdNumeric, nil
}

// resolveSiteAndIpVersion turns the user-supplied site reference (id or label)
// and IP version into the values the API path segments need. The IP version is
// validated against the SDK enum so a typo fails before the request is sent.
func resolveSiteAndIpVersion(ctx context.Context, siteIdOrLabel string, ipVersion string) (int64, string, error) {
	siteInfo, err := site.GetSiteByIdOrLabel(ctx, siteIdOrLabel)
	if err != nil {
		return 0, "", err
	}

	ipVersionValue, err := sdk.NewIpVersionFromValue(ipVersion)
	if err != nil {
		logger.Get().Error().Err(err).Msg("")
		return 0, "", err
	}

	return int64(siteInfo.Id), string(*ipVersionValue), nil
}
