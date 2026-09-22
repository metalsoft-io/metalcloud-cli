package dns_zone

import (
	"context"
	"fmt"
	"strconv"
	"strings"

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
			Title: "ID",
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

	request := client.DNSZoneAPI.GetDNSZones(ctx).SortBy([]string{"id:ASC"})
	if len(filterIsDefault) > 0 {
		request = request.FilterIsDefault(utils.ProcessFilterStringSlice(filterIsDefault))
	}

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &dnsZonePrintConfig)
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

var dnsRecordSetPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"Name": {
			Order: 2,
		},
		"Type": {
			Order: 3,
		},
		"Records": {
			Order: 4,
			Transformer: func(v interface{}) string {
				if records, ok := v.([]string); ok {
					return strings.Join(records, ", ")
				}
				return fmt.Sprintf("%v", v)
			},
		},
		"Ttl": {
			Title: "TTL",
			Order: 5,
		},
		"Status": {
			Order:       6,
			Transformer: formatter.FormatStatusValue,
		},
		"ZoneName": {
			Title: "Zone",
			Order: 7,
		},
	},
}

// DNSZoneRecords lists DNS record sets. When dnsZoneId is empty, every record
// set of every zone is listed through the global /dns-recordsets endpoint;
// otherwise only the record sets of the given zone are listed.
func DNSZoneRecords(ctx context.Context, dnsZoneId string) error {
	if dnsZoneId == "" {
		return dnsRecordSetsListAll(ctx)
	}

	logger.Get().Info().Msgf("Getting DNS record sets for zone '%s'", dnsZoneId)

	id, err := GetDNSZoneId(dnsZoneId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	result, httpRes, err := client.DNSZoneAPI.ListDNSRecordSetsByZoneId(ctx, id).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(result, &dnsRecordSetPrintConfig)
}

// dnsRecordSetsListAll lists every DNS record set across all zones.
func dnsRecordSetsListAll(ctx context.Context) error {
	logger.Get().Info().Msgf("Listing all DNS record sets")

	client := api.GetApiClient(ctx)

	request := client.DNSRecordSetAPI.ListDNSRecordSets(ctx).SortBy([]string{"id:ASC"})

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &dnsRecordSetPrintConfig)
}

// DNSZoneRecord shows a single DNS record set of a zone.
func DNSZoneRecord(ctx context.Context, dnsZoneId string, recordSetId string) error {
	logger.Get().Info().Msgf("Getting DNS record set '%s' of zone '%s'", recordSetId, dnsZoneId)

	id, err := GetDNSZoneId(dnsZoneId)
	if err != nil {
		return err
	}

	recordSetIdNumeric, err := strconv.ParseInt(recordSetId, 10, 64)
	if err != nil {
		err := fmt.Errorf("invalid DNS record set ID: '%s'", recordSetId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	client := api.GetApiClient(ctx)

	recordSet, httpRes, err := client.DNSZoneAPI.GetDNSRecordSetById(ctx, id, recordSetIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(recordSet, &dnsRecordSetPrintConfig)
}

var dnsZoneNameserverPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Nameserver": {
			Title:    "Nameserver",
			MaxWidth: 80,
			Order:    1,
		},
	},
}

// dnsZoneNameserverRow is the table projection of one nameserver: the endpoint
// returns a bare list of strings, which the tabular formatter cannot render on
// its own. The json and yaml formats keep the plain list.
type dnsZoneNameserverRow struct {
	Nameserver string
}

// DNSZoneNameservers lists the nameservers of a DNS zone.
func DNSZoneNameservers(ctx context.Context, dnsZoneId string) error {
	logger.Get().Info().Msgf("Getting nameservers of DNS zone '%s'", dnsZoneId)

	id, err := GetDNSZoneId(dnsZoneId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	nameservers, httpRes, err := client.DNSZoneAPI.GetDNSZoneNameservers(ctx, id).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	if formatter.IsNativeFormat() {
		return formatter.PrintResult(nameservers, nil)
	}

	rows := make([]dnsZoneNameserverRow, 0, len(nameservers))
	for _, nameserver := range nameservers {
		rows = append(rows, dnsZoneNameserverRow{Nameserver: nameserver})
	}

	return formatter.PrintResult(rows, &dnsZoneNameserverPrintConfig)
}

func GetDNSZoneId(dnsZoneId string) (int64, error) {
	dnsZoneIdNumeric, err := strconv.ParseInt(dnsZoneId, 10, 64)
	if err != nil {
		err := fmt.Errorf("invalid DNS zone ID: '%s'", dnsZoneId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return dnsZoneIdNumeric, nil
}
