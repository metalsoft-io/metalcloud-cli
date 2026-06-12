package server_type

import (
	"context"
	"fmt"
	"net/http"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

type serverTypeRaw struct {
	Id                   interface{} `json:"id"`
	Name                 *string     `json:"name"`
	Label                *string     `json:"label"`
	ProcessorCount       interface{} `json:"processorCount"`
	ProcessorCoreMhz     interface{} `json:"processorCoreMhz"`
	ProcessorNames       interface{} `json:"processorNames"`
	RamGbytes            interface{} `json:"ramGbytes"`
	NetworkInterfaceCount interface{} `json:"networkInterfaceCount"`
	DiskCount            interface{} `json:"diskCount"`
	GpuCount             interface{} `json:"gpuCount"`
}

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

	rawItems, meta, err := utils.FetchAllPagesRaw(func(p float32) (*http.Response, error) {
		_, httpRes, _ := client.ServerTypeAPI.GetServerTypes(ctx).SortBy([]string{"id:ASC"}).Page(p).Limit(100).Execute()
		return httpRes, nil
	})
	if err != nil {
		return err
	}
	records, err := utils.UnmarshalRawItems[serverTypeRaw](rawItems)
	if err != nil {
		return fmt.Errorf("failed to parse server types: %w", err)
	}

	return utils.PrintAllRaw(rawItems, records, meta, len(records), &serverTypePrintConfig)
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
