package dns_zone

import (
	"context"
	"fmt"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var dnsZonePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"Label": {
			Title: "Label",
			Order: 2,
		},
		"ZoneName": {
			Title: "Zone Name",
			Order: 3,
		},
		"ZoneType": {
			Title: "Type",
			Order: 4,
		},
		"Status": {
			Title:       "Status",
			Transformer: formatter.FormatStatusValue,
			Order:       5,
		},
		"IsDefault": {
			Title: "Default",
			Order: 6,
		},
		"SoaEmail": {
			Title: "SOA Email",
			Order: 7,
		},
		"Ttl": {
			Title: "TTL",
			Order: 8,
		},
		"Description": {
			Title: "Description",
			Order: 9,
		},
	},
}

func DNSZoneList(ctx context.Context, filterIsDefault []string) error {
	logger.Get().Info().Msgf("Listing DNS zones")

	client := api.GetApiClient(ctx)

	request := client.DNSZoneAPI.GetDNSZones(ctx)

	if len(filterIsDefault) > 0 {
		request = request.FilterIsDefault(utils.ProcessFilterStringSlice(filterIsDefault))
	}

	dnsZoneList, httpRes, err := request.Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(dnsZoneList, &dnsZonePrintConfig)
}

func DNSZoneGet(ctx context.Context, dnsZoneId string) error {
	logger.Get().Info().Msgf("Get DNS zone '%s'", dnsZoneId)

	dnsZoneIdNumeric, err := GetDNSZoneId(dnsZoneId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	dnsZone, httpRes, err := client.DNSZoneAPI.GetDNSZoneById(ctx, dnsZoneIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(dnsZone, &dnsZonePrintConfig)
}

func DNSZoneCreate(ctx context.Context, zoneConfig sdk.CreateDnsZone) error {
	logger.Get().Info().Msgf("Creating DNS zone")

	client := api.GetApiClient(ctx)

	dnsZone, httpRes, err := client.DNSZoneAPI.CreateDNSZone(ctx).CreateDnsZone(zoneConfig).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(dnsZone, &dnsZonePrintConfig)
}

func DNSZoneUpdate(ctx context.Context, dnsZoneId string, config []byte) error {
	logger.Get().Info().Msgf("Updating DNS zone '%s'", dnsZoneId)

	var updateConfig sdk.UpdateDnsZone
	err := utils.UnmarshalContent(config, &updateConfig)
	if err != nil {
		return err
	}

	dnsZoneIdNumeric, err := GetDNSZoneId(dnsZoneId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	dnsZone, httpRes, err := client.DNSZoneAPI.UpdateDNSZone(ctx, dnsZoneIdNumeric).UpdateDnsZone(updateConfig).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(dnsZone, &dnsZonePrintConfig)
}

func DNSZoneDelete(ctx context.Context, dnsZoneId string) error {
	logger.Get().Info().Msgf("Deleting DNS zone '%s'", dnsZoneId)

	dnsZoneIdNumeric, err := GetDNSZoneId(dnsZoneId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.DNSZoneAPI.DeleteDNSZone(ctx, dnsZoneIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("DNS zone '%s' deleted successfully", dnsZoneId)
	return nil
}

func GetDNSZoneId(dnsZoneId string) (float32, error) {
	dnsZoneIdNumeric, err := strconv.ParseFloat(dnsZoneId, 32)
	if err != nil {
		err := fmt.Errorf("invalid DNS zone ID: '%s'", dnsZoneId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return float32(dnsZoneIdNumeric), nil
}
