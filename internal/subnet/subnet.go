package subnet

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

	subnetList, httpRes, err := client.SubnetAPI.GetSubnets(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(subnetList, &SubnetPrintConfig)
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

// subnetIpRaw works around the SDK bug where Links is typed as
// map[string]interface{} but the API may return an array.
type subnetIpRaw struct {
	Id        float32     `json:"id"`
	Name      string      `json:"name"`
	Address   string      `json:"address"`
	IpVersion string      `json:"ipVersion"`
	SubnetId  float32     `json:"subnetId"`
	Links     interface{} `json:"links,omitempty"`
}

type subnetIpListRaw struct {
	Data []subnetIpRaw `json:"data"`
}

// subnetIpRangeRaw works around the SDK bug where Links is typed as
// map[string]interface{} but the API may return an array.
type subnetIpRangeRaw struct {
	Id           float32     `json:"id"`
	Name         string      `json:"name"`
	StartAddress string      `json:"startAddress"`
	EndAddress   string      `json:"endAddress"`
	IpVersion    string      `json:"ipVersion"`
	SubnetId     float32     `json:"subnetId"`
	Links        interface{} `json:"links,omitempty"`
}

type subnetIpRangeListRaw struct {
	Data []subnetIpRangeRaw `json:"data"`
}

func SubnetIps(ctx context.Context, subnetId string) error {
	logger.Get().Info().Msgf("Getting IPs for subnet '%s'", subnetId)

	id, err := getSubnetId(subnetId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	_, httpRes, sdkErr := client.SubnetAPI.GetSubnetIps(ctx, id).Execute()

	if httpRes != nil && httpRes.StatusCode >= 400 {
		if err := response_inspector.InspectResponse(httpRes, sdkErr); err != nil {
			return err
		}
	} else if httpRes == nil {
		return sdkErr
	}

	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var raw subnetIpListRaw
	if err := json.Unmarshal(body, &raw); err != nil {
		return fmt.Errorf("failed to parse subnet IPs: %w", err)
	}

	return formatter.PrintResult(raw.Data, &subnetIpPrintConfig)
}

func SubnetIpRanges(ctx context.Context, subnetId string) error {
	logger.Get().Info().Msgf("Getting IP ranges for subnet '%s'", subnetId)

	id, err := getSubnetId(subnetId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	_, httpRes, sdkErr := client.SubnetAPI.GetSubnetIpRanges(ctx, id).Execute()

	if httpRes != nil && httpRes.StatusCode >= 400 {
		if err := response_inspector.InspectResponse(httpRes, sdkErr); err != nil {
			return err
		}
	} else if httpRes == nil {
		return sdkErr
	}

	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var raw subnetIpRangeListRaw
	if err := json.Unmarshal(body, &raw); err != nil {
		return fmt.Errorf("failed to parse subnet IP ranges: %w", err)
	}

	return formatter.PrintResult(raw.Data, &subnetIpRangePrintConfig)
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
