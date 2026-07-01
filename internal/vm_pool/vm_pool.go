package vm_pool

import (
	"context"
	"fmt"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/internal/vm"
	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var VMPoolPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"SiteId": {
			Title: "Site",
			Order: 2,
		},
		"Name": {
			Title: "Name",
			Order: 3,
		},
		"Type": {
			Title: "Type",
			Order: 4,
		},
		"ManagementHost": {
			Title: "Management Host",
			Order: 5,
		},
		"ManagementPort": {
			Title: "Management Port",
			Order: 6,
		},
		"Status": {
			Title:       "Status",
			Transformer: formatter.FormatStatusValue,
			Order:       7,
		},
	},
}

func VMPoolList(ctx context.Context, filterType []string) error {
	logger.Get().Info().Msgf("Listing all VM pools")

	client := api.GetApiClient(ctx)

	request := client.VMPoolAPI.GetVMPools(ctx).SortBy([]string{"id:ASC"})

	if len(filterType) > 0 {
		request = request.FilterType(utils.ProcessFilterStringSlice(filterType))
	}

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &VMPoolPrintConfig)
}

func VMPoolGet(ctx context.Context, vmPoolId string) error {
	logger.Get().Info().Msgf("Get VM pool %s details", vmPoolId)

	vmPoolIdNumeric, err := getVMPoolId(vmPoolId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	vmPool, httpRes, err := client.VMPoolAPI.GetVMPool(ctx, float32(vmPoolIdNumeric)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(vmPool, &VMPoolPrintConfig)
}

func VMPoolCreate(ctx context.Context, config []byte) error {
	logger.Get().Info().Msgf("Creating new VM pool")

	client := api.GetApiClient(ctx)

	var createVMPool sdk.CreateVMPool
	err := utils.UnmarshalContent(config, &createVMPool)
	if err != nil {
		return err
	}

	response, httpRes, err := client.VMPoolAPI.CreateVMPool(ctx).
		CreateVMPool(createVMPool).
		Execute()

	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("VM pool created with ID: %d", int(response.Id))
	return nil
}

func VMPoolDelete(ctx context.Context, vmPoolId string) error {
	logger.Get().Info().Msgf("Deleting VM pool %s", vmPoolId)

	vmPoolIdNumeric, err := getVMPoolId(vmPoolId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.VMPoolAPI.DeleteVMPool(ctx, vmPoolIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("VM pool %s deleted successfully", vmPoolId)
	return nil
}

func VMPoolGetCredentials(ctx context.Context, vmPoolId string) error {
	logger.Get().Info().Msgf("Getting credentials for VM pool %s", vmPoolId)

	vmPoolIdNumeric, err := getVMPoolId(vmPoolId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	credentials, httpRes, err := client.VMPoolAPI.GetVMPoolCredentials(ctx, vmPoolIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(credentials, nil)
}

func VMPoolGetVMs(ctx context.Context, vmPoolId string, limit float32, page float32) error {
	logger.Get().Info().Msgf("Getting VMs for VM pool %s", vmPoolId)

	vmPoolIdNumeric, err := getVMPoolId(vmPoolId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.VMPoolAPI.GetVMPoolVMs(ctx, vmPoolIdNumeric)

	if limit > 0 || page > 0 {
		if limit > 0 {
			request = request.Limit(limit)
		}
		if page > 0 {
			request = request.Page(page)
		}

		result, httpRes, err := request.Execute()
		if err := response_inspector.InspectResponse(httpRes, err); err != nil {
			return err
		}

		return utils.PrintAll(result.Data, result.Meta, len(result.Data), &vm.VMListPrintConfig)
	}

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &vm.VMListPrintConfig)
}

func VMPoolGetClusterHosts(ctx context.Context, vmPoolId string, limit float32, page float32) error {
	logger.Get().Info().Msgf("Getting cluster hosts for VM pool %s", vmPoolId)

	vmPoolIdNumeric, err := getVMPoolId(vmPoolId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.VMPoolAPI.GetVMPoolClusterHosts(ctx, vmPoolIdNumeric)

	if limit > 0 || page > 0 {
		if limit > 0 {
			request = request.Limit(limit)
		}
		if page > 0 {
			request = request.Page(page)
		}

		result, httpRes, err := request.Execute()
		if err := response_inspector.InspectResponse(httpRes, err); err != nil {
			return err
		}

		return utils.PrintAll(result.Data, result.Meta, len(result.Data), nil)
	}

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), nil)
}

func VMPoolGetClusterHostVMs(ctx context.Context, vmPoolId string, hostId string, limit float32, page float32) error {
	logger.Get().Info().Msgf("Getting VMs for cluster host %s in VM pool %s", hostId, vmPoolId)

	vmPoolIdNumeric, err := getVMPoolId(vmPoolId)
	if err != nil {
		return err
	}

	hostIdNumeric, err := getVMPoolId(hostId) // reusing the same conversion function
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.VMPoolAPI.GetVMPoolClusterHostVMs(ctx, vmPoolIdNumeric, hostIdNumeric)

	if limit > 0 || page > 0 {
		if limit > 0 {
			request = request.Limit(limit)
		}
		if page > 0 {
			request = request.Page(page)
		}

		result, httpRes, err := request.Execute()
		if err := response_inspector.InspectResponse(httpRes, err); err != nil {
			return err
		}

		return utils.PrintAll(result.Data, result.Meta, len(result.Data), &vm.VMListPrintConfig)
	}

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &vm.VMListPrintConfig)
}

func VMPoolGetClusterHostInterfaces(ctx context.Context, vmPoolId string, hostId string) error {
	logger.Get().Info().Msgf("Getting interfaces for cluster host %s in VM pool %s", hostId, vmPoolId)

	vmPoolIdNumeric, err := getVMPoolId(vmPoolId)
	if err != nil {
		return err
	}

	hostIdNumeric, err := getVMPoolId(hostId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	// GetVMPoolClusterHostInterfaces returns a flat []VMPoolHostInterfaces — no Page/Limit methods.
	interfaces, httpRes, err := client.VMPoolAPI.GetVMPoolClusterHostInterfaces(ctx, vmPoolIdNumeric, hostIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(interfaces, nil)
}

func VMPoolConfigExample(ctx context.Context) error {
	vmPoolConfig := sdk.CreateVMPool{
		SiteId:         1,
		ManagementHost: "vcenter.example.com",
		ManagementPort: 443,
		Name:           "VM-Pool-Example",
		Description:    sdk.PtrString("Example VM pool for testing"),
		Type:           "vmware",
		Certificate:    sdk.PtrString("-----BEGIN CERTIFICATE-----\nMIID...certificate content...==\n-----END CERTIFICATE-----"),
		PrivateKey:     sdk.PtrString("-----BEGIN PRIVATE KEY-----\nMIIE...key content...==\n-----END PRIVATE KEY-----"),
		InMaintenance:  sdk.PtrFloat32(0),
		IsExperimental: sdk.PtrFloat32(0),
		Tags:           []string{"test", "example", "vmware"},
	}

	return formatter.PrintResult(vmPoolConfig, nil)
}

func VMPoolImportVMs(ctx context.Context, vmPoolId string, importVMs sdk.VMPoolImportVMs) error {
	logger.Get().Info().Msgf("Importing VMs into VM pool %s", vmPoolId)

	vmPoolIdNumeric, err := getVMPoolId(vmPoolId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.VMPoolAPI.ImportVMPoolVMs(ctx, vmPoolIdNumeric).
		VMPoolImportVMs(importVMs).
		Execute()

	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("VMs imported successfully into VM pool %s", vmPoolId)
	return nil
}

func getVMPoolId(vmPoolId string) (int64, error) {
	vmPoolIdNumeric, err := strconv.ParseInt(vmPoolId, 10, 64)
	if err != nil {
		err := fmt.Errorf("invalid VM pool ID: '%s'", vmPoolId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return vmPoolIdNumeric, nil
}
