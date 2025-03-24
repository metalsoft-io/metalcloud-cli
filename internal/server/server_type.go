package server

import (
	"context"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var serverTypePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"Name": {
			MaxWidth: 30,
			Order:    2,
		},
		"Label": {
			MaxWidth: 30,
			Order:    3,
		},
		"ProcessorCount": {
			Title: "CPU #",
			Order: 4,
		},
		"ProcessorCoreMhz": {
			Title: "CPU MHz",
			Order: 5,
		},
		"ProcessorNames": {
			Title: "CPU Names",
			Order: 6,
		},
		"RamGbytes": {
			Title: "RAM GB",
			Order: 7,
		},
		"NetworkInterfaceCount": {
			Title: "NIC #",
			Order: 8,
		},
		"DiskCount": {
			Title: "Disk #",
			Order: 9,
		},
		"GpuCount": {
			Title: "GPU #",
			Order: 10,
		},
	},
}

func ServerTypeList(ctx context.Context) error {
	logger.Get().Info().Msgf("Listing all server types")

	client := api.GetApiClient(ctx)

	typesList, httpRes, err := client.ServerTypeAPI.GetServerTypes(ctx).SortBy([]string{"id:ASC"}).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(typesList, &serverTypePrintConfig)
}

func ServerTypeGet(ctx context.Context, serverTypeIdOrLabel string) error {
	logger.Get().Info().Msgf("Get server type %s info", serverTypeIdOrLabel)

	serverType, err := GetServerTypeByIdOrLabel(ctx, serverTypeIdOrLabel)
	if err != nil {
		return err
	}

	return formatter.PrintResult(serverType, &serverTypePrintConfig)
}

func GetServerTypeByIdOrLabel(ctx context.Context, serverTypeIdOrLabel string) (*sdk.ServerType, error) {
	client := api.GetApiClient(ctx)

	serverTypeId, err := utils.GetFloat32FromString(serverTypeIdOrLabel)
	if err != nil {
		return nil, err
	}

	serverTypeInfo, httpRes, err := client.ServerTypeAPI.GetServerTypeInfo(ctx, serverTypeId).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, err
	}

	return serverTypeInfo, nil
}
