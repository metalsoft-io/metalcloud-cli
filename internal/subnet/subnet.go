package subnet

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
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

	err := json.Unmarshal(config, &subnetConfig)
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

	err = json.Unmarshal(config, &subnetConfig)
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
