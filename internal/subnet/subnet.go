package subnet

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

type subnetRaw struct {
	Id             interface{} `json:"id"`
	Name           *string     `json:"name"`
	IpVersion      interface{} `json:"ipVersion"`
	NetworkAddress *string     `json:"networkAddress"`
	PrefixLength   interface{} `json:"prefixLength"`
	Netmask        *string     `json:"netmask"`
	IsPool         interface{} `json:"isPool"`
	CreatedAt      interface{} `json:"createdAt"`
}

var SubnetPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
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

	rawItems, meta, err := utils.FetchAllPagesRaw(func(p float32) (*http.Response, error) {
		_, httpRes, _ := client.SubnetAPI.GetSubnets(ctx).SortBy([]string{"id:ASC"}).Page(p).Limit(100).Execute()
		return httpRes, nil
	})
	if err != nil {
		return err
	}
	records, err := utils.UnmarshalRawItems[subnetRaw](rawItems)
	if err != nil {
		return fmt.Errorf("failed to parse subnets: %w", err)
	}

	return utils.PrintAllRaw(rawItems, records, meta, len(records), &SubnetPrintConfig)
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
		UpdateSubnet(ctx, int32(subnetIdNumeric)).
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
		DeleteSubnet(ctx, int32(subnetIdNumeric)).
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
		ParentSubnetId:         sdk.PtrInt32(0),
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
			Title: "#",
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
			Title: "#",
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

	result, httpRes, err := client.SubnetAPI.GetSubnetIps(ctx, id).Execute()
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

	result, httpRes, err := client.SubnetAPI.GetSubnetIpRanges(ctx, id).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(result, &subnetIpRangePrintConfig)
}

func getSubnetId(subnetId string) (float32, error) {
	subnetIdNumeric, err := strconv.ParseFloat(subnetId, 32)
	if err != nil {
		err := fmt.Errorf("invalid subnet ID: '%s'", subnetId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return float32(subnetIdNumeric), nil
}

func getSubnetIdAndRevision(ctx context.Context, subnetId string) (float32, string, error) {
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
