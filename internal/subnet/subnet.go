package subnet

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

var SubnetPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"Name": {
			Title: "Name",
			Order: 2,
		},
		"IpVersion": {
			Title: "IP Version",
			Order: 3,
		},
		"NetworkAddress": {
			Title: "Network Address",
			Order: 4,
		},
		"PrefixLength": {
			Title: "Prefix",
			Order: 5,
		},
		"Netmask": {
			Title: "Netmask",
			Order: 6,
		},
		"IsPool": {
			Title: "Pool",
			Order: 7,
		},
		"CreatedAt": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       8,
		},
	},
}

func SubnetList(ctx context.Context) error {
	logger.Get().Info().Msgf("Listing all subnets")

	client := api.GetApiClient(ctx)

	request := client.SubnetAPI.GetSubnets(ctx).SortBy([]string{"id:ASC"})

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &SubnetPrintConfig)
}

func SubnetGet(ctx context.Context, subnetId string) error {
	logger.Get().Info().Msgf("Get subnet %s details", subnetId)

	subnetIdNumeric, err := getSubnetId(subnetId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	subnet, httpRes, err := client.SubnetAPI.GetSubnet(ctx, subnetIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(subnet, &SubnetPrintConfig)
}

func SubnetCreate(ctx context.Context, config []byte) error {
	logger.Get().Info().Msgf("Creating subnet")

	var subnetConfig sdk.CreateSubnet
	err := utils.UnmarshalContent(config, &subnetConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	subnetInfo, httpRes, err := client.SubnetAPI.CreateSubnet(ctx).CreateSubnet(subnetConfig).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(subnetInfo, &SubnetPrintConfig)
}

func SubnetUpdate(ctx context.Context, subnetId string, config []byte) error {
	logger.Get().Info().Msgf("Updating subnet %s", subnetId)

	subnetIdNumeric, revision, err := getSubnetIdAndRevision(ctx, subnetId)
	if err != nil {
		return err
	}

	var subnetConfig sdk.UpdateSubnet
	err = utils.UnmarshalContent(config, &subnetConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	subnetInfo, httpRes, err := client.SubnetAPI.
		UpdateSubnet(ctx, subnetIdNumeric).
		UpdateSubnet(subnetConfig).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(subnetInfo, &SubnetPrintConfig)
}

func SubnetDelete(ctx context.Context, subnetId string) error {
	logger.Get().Info().Msgf("Deleting subnet %s", subnetId)

	subnetIdNumeric, revision, err := getSubnetIdAndRevision(ctx, subnetId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.SubnetAPI.
		DeleteSubnet(ctx, subnetIdNumeric).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Subnet %s deleted", subnetId)
	return nil
}

func SubnetConfigExample(ctx context.Context) error {
	// Example create subnet configuration
	subnetConfiguration := sdk.CreateSubnet{
		Label:                  sdk.PtrString("example-subnet"),
		Name:                   sdk.PtrString("example-subnet"),
		NetworkAddress:         "192.168.1.0",
		PrefixLength:           24,
		IsPool:                 false,
		ParentSubnetId:         sdk.PtrInt64(0),
		DefaultGatewayAddress:  sdk.PtrString("192.168.1.1"),
		AllocationDenylist:     []sdk.AddressRange{},
		ChildOverlapAllowRules: []string{},
		Tags:                   &map[string]string{"tag1": "value1", "tag2": "value2"},
	}

	return formatter.PrintResult(subnetConfiguration, nil)
}

var subnetIpPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"Name": {
			Order: 2,
		},
		"Address": {
			Order: 3,
		},
		"IpVersion": {
			Title: "IP Version",
			Order: 4,
		},
		"SubnetId": {
			Title: "Subnet",
			Order: 5,
		},
	},
}

var subnetIpRangePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"Name": {
			Order: 2,
		},
		"StartAddress": {
			Title: "Start",
			Order: 3,
		},
		"EndAddress": {
			Title: "End",
			Order: 4,
		},
		"IpVersion": {
			Title: "IP Version",
			Order: 5,
		},
		"SubnetId": {
			Title: "Subnet",
			Order: 6,
		},
	},
}

func SubnetIps(ctx context.Context, subnetId string) error {
	logger.Get().Info().Msgf("Getting IPs for subnet '%s'", subnetId)

	id, err := getSubnetId(subnetId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	result, httpRes, err := client.SubnetAPI.GetSubnetIps(ctx, float32(id)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(result, &subnetIpPrintConfig)
}

func SubnetIpRanges(ctx context.Context, subnetId string) error {
	logger.Get().Info().Msgf("Getting IP ranges for subnet '%s'", subnetId)

	id, err := getSubnetId(subnetId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	result, httpRes, err := client.SubnetAPI.GetSubnetIpRanges(ctx, float32(id)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(result, &subnetIpRangePrintConfig)
}

// subnetCapacityPrintConfig renders the capacity report of a subnet. The
// counts are decimal strings because a large IPv6 subnet holds more addresses
// than a JSON number can carry exactly.
var subnetCapacityPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"SubnetId": {
			Title: "Subnet",
			Order: 1,
		},
		"Prefix": {
			Title: "Prefix",
			Order: 2,
		},
		"IsPool": {
			Title: "Pool",
			Order: 3,
		},
		"TotalIpCount": {
			Title: "Total IPs",
			Order: 4,
		},
		"FreeIpCount": {
			Title: "Free IPs",
			Order: 5,
		},
		"UsedIpCount": {
			Title: "Used IPs",
			Order: 6,
		},
		"FreePrefixes": {
			Title:    "Free Prefixes",
			MaxWidth: 60,
			Order:    7,
		},
		"RequestedPrefixLength": {
			Title: "Requested Prefix",
			Order: 8,
		},
		"AllocatablePrefixCount": {
			Title: "Allocatable Blocks",
			Order: 9,
		},
	},
}

// subnetCapacityDisplay is the table view of a subnet capacity report. The SDK
// model stores the counts in NullableString/NullableInt32 wrappers, which the
// table formatter renders as empty cells, so the wrappers are unwrapped here.
// The json/yaml output keeps the untouched SDK object.
type subnetCapacityDisplay struct {
	SubnetId               int64
	Prefix                 string
	IsPool                 bool
	TotalIpCount           string
	FreeIpCount            string
	UsedIpCount            string
	FreePrefixes           []string
	RequestedPrefixLength  string
	AllocatablePrefixCount string
}

func toSubnetCapacityDisplay(capacity *sdk.SubnetCapacity) subnetCapacityDisplay {
	display := subnetCapacityDisplay{
		SubnetId:     capacity.SubnetId,
		Prefix:       capacity.Prefix,
		IsPool:       capacity.IsPool,
		FreePrefixes: capacity.FreePrefixes,
	}

	if value := capacity.TotalIpCount.Get(); value != nil {
		display.TotalIpCount = *value
	}
	if value := capacity.FreeIpCount.Get(); value != nil {
		display.FreeIpCount = *value
	}
	if value := capacity.UsedIpCount.Get(); value != nil {
		display.UsedIpCount = *value
	}
	if value := capacity.AllocatablePrefixCount.Get(); value != nil {
		display.AllocatablePrefixCount = *value
	}
	if value := capacity.RequestedPrefixLength.Get(); value != nil {
		display.RequestedPrefixLength = strconv.FormatInt(int64(*value), 10)
	}

	return display
}

// SubnetCapacity reports how much of a subnet is still available. For a pool,
// prefixLength (when > 0) also asks how many blocks of that size still fit.
func SubnetCapacity(ctx context.Context, subnetId string, prefixLength int) error {
	logger.Get().Info().Msgf("Getting capacity of subnet '%s'", subnetId)

	subnetIdNumeric, err := getSubnetId(subnetId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.SubnetAPI.GetSubnetCapacity(ctx, subnetIdNumeric)
	if prefixLength > 0 {
		request = request.PrefixLength(int32(prefixLength))
	}

	capacity, httpRes, err := request.Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	if formatter.IsNativeFormat() {
		return formatter.PrintResult(capacity, &subnetCapacityPrintConfig)
	}

	return formatter.PrintResult(toSubnetCapacityDisplay(capacity), &subnetCapacityPrintConfig)
}

func SubnetIpRemove(ctx context.Context, subnetId string, ipId string) error {
	logger.Get().Info().Msgf("Removing IP '%s' from subnet '%s'", ipId, subnetId)

	subnetIdNumeric, revision, err := getSubnetIdAndRevision(ctx, subnetId)
	if err != nil {
		return err
	}

	ipIdNumeric, err := utils.GetInt64FromString(ipId)
	if err != nil {
		err = fmt.Errorf("invalid IP ID: '%s'", ipId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.SubnetAPI.
		DeleteSubnetIp(ctx, subnetIdNumeric, ipIdNumeric).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("IP '%s' removed from subnet '%s'", ipId, subnetId)
	return nil
}

func SubnetIpRangeRemove(ctx context.Context, subnetId string, ipRangeId string) error {
	logger.Get().Info().Msgf("Removing IP range '%s' from subnet '%s'", ipRangeId, subnetId)

	subnetIdNumeric, revision, err := getSubnetIdAndRevision(ctx, subnetId)
	if err != nil {
		return err
	}

	ipRangeIdNumeric, err := utils.GetInt64FromString(ipRangeId)
	if err != nil {
		err = fmt.Errorf("invalid IP range ID: '%s'", ipRangeId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.SubnetAPI.
		DeleteSubnetIpRange(ctx, subnetIdNumeric, ipRangeIdNumeric).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("IP range '%s' removed from subnet '%s'", ipRangeId, subnetId)
	return nil
}

func getSubnetId(subnetId string) (int64, error) {
	subnetIdNumeric, err := strconv.ParseInt(subnetId, 10, 64)
	if err != nil {
		err := fmt.Errorf("invalid subnet ID: '%s'", subnetId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return subnetIdNumeric, nil
}

func getSubnetIdAndRevision(ctx context.Context, subnetId string) (int64, string, error) {
	subnetIdNumeric, err := getSubnetId(subnetId)
	if err != nil {
		return 0, "", err
	}

	client := api.GetApiClient(ctx)

	subnet, httpRes, err := client.SubnetAPI.GetSubnet(ctx, subnetIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return 0, "", err
	}

	return subnetIdNumeric, strconv.Itoa(int(subnet.Revision)), nil
}
