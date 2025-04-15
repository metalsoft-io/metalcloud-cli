package vm_type

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

var VMTypePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"Label": {
			Title:    "Label",
			MaxWidth: 30,
			Order:    2,
		},
		"Name": {
			Title:    "Name",
			MaxWidth: 30,
			Order:    3,
		},
		"DisplayName": {
			Title:    "Display Name",
			MaxWidth: 30,
			Order:    4,
		},
		"CpuCores": {
			Title: "CPU Cores",
			Order: 5,
		},
		"RamGB": {
			Title: "RAM (GB)",
			Order: 6,
		},
		"IsExperimental": {
			Title: "Experimental",
			Order: 7,
		},
		"ForUnmanagedVMsOnly": {
			Title: "For Unmanaged VMs Only",
			Order: 8,
		},
	},
}

func VMTypeList(ctx context.Context, limit float32, page float32) error {
	logger.Get().Info().Msg("Listing VM types")

	client := api.GetApiClient(ctx)

	request := client.VMTypeAPI.GetVMTypes(ctx)

	// Set pagination if provided
	if limit > 0 {
		request = request.Limit(limit)
	}

	if page > 0 {
		request = request.Page(page)
	}

	// Sort by id ascending
	request = request.SortBy([]string{"id:ASC"})

	vmTypes, httpRes, err := request.Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(vmTypes, &VMTypePrintConfig)
}

func VMTypeGet(ctx context.Context, vmTypeId string) error {
	logger.Get().Info().Msgf("Getting VM type %s details", vmTypeId)

	vmTypeIdNumeric, err := getVMTypeId(vmTypeId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	vmType, httpRes, err := client.VMTypeAPI.GetVMType(ctx, vmTypeIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(vmType, &VMTypePrintConfig)
}

func VMTypeCreate(ctx context.Context, config []byte) error {
	logger.Get().Info().Msg("Creating new VM type")

	client := api.GetApiClient(ctx)

	var createVMType sdk.CreateVMType

	err := json.Unmarshal(config, &createVMType)
	if err != nil {
		return err
	}

	response, httpRes, err := client.VMTypeAPI.CreateVMType(ctx).
		CreateVMType(createVMType).
		Execute()

	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("VM type created with ID: %d", int(response.Id))
	return nil
}

func VMTypeUpdate(ctx context.Context, vmTypeId string, config []byte) error {
	logger.Get().Info().Msgf("Updating VM type %s", vmTypeId)

	vmTypeIdNumeric, err := getVMTypeId(vmTypeId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	var updateVMType sdk.UpdateVMType

	err = json.Unmarshal(config, &updateVMType)
	if err != nil {
		return err
	}

	response, httpRes, err := client.VMTypeAPI.UpdateVMType(ctx, vmTypeIdNumeric).
		UpdateVMType(updateVMType).
		Execute()

	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("VM type %s updated successfully", vmTypeId)
	return formatter.PrintResult(response, &VMTypePrintConfig)
}

func VMTypeDelete(ctx context.Context, vmTypeId string) error {
	logger.Get().Info().Msgf("Deleting VM type %s", vmTypeId)

	vmTypeIdNumeric, err := getVMTypeId(vmTypeId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.VMTypeAPI.DeleteVMType(ctx, vmTypeIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("VM type %s deleted successfully", vmTypeId)
	return nil
}

func VMTypeGetVMs(ctx context.Context, vmTypeId string, limit float32, page float32) error {
	logger.Get().Info().Msgf("Getting VMs for VM type %s", vmTypeId)

	vmTypeIdNumeric, err := getVMTypeId(vmTypeId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.VMTypeAPI.GetVMsByVMType(ctx, vmTypeIdNumeric)

	// Set pagination if provided
	if limit > 0 {
		request = request.Limit(limit)
	}

	if page > 0 {
		request = request.Page(page)
	}

	vms, httpRes, err := request.Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(vms, nil)
}

func VMTypeConfigExample(ctx context.Context) error {
	vmTypeConfig := sdk.CreateVMType{
		Name:                "example-vm-type",
		DisplayName:         sdk.PtrString("Example VM Type"),
		Label:               sdk.PtrString("Example"),
		CpuCores:            4,
		RamGB:               8,
		IsExperimental:      sdk.PtrFloat32(0),
		ForUnmanagedVMsOnly: sdk.PtrFloat32(0),
		Tags:                []string{"example", "test"},
	}

	return formatter.PrintResult(vmTypeConfig, nil)
}

func getVMTypeId(vmTypeId string) (float32, error) {
	vmTypeIdNumeric, err := strconv.ParseFloat(vmTypeId, 32)
	if err != nil {
		err := fmt.Errorf("invalid VM type ID: '%s'", vmTypeId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return float32(vmTypeIdNumeric), nil
}
