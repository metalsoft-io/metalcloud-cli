package server

import (
	"context"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
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
