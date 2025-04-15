package vm

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

var vmPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"VmId": {
			Title: "#",
			Order: 1,
		},
		"HostId": {
			Title: "Host",
			Order: 2,
		},
		"VmStatus": {
			Title:       "Status",
			Transformer: formatter.FormatStatusValue,
			Order:       3,
		},
		"VmName": {
			Title: "Name",
			Order: 4,
		},
		"VmDescription": {
			Title: "Description",
			Order: 5,
		},
		"VmCPUs": {
			Title: "CPUs",
			Order: 6,
		},
		"VmMemoryMB": {
			Title: "Memory (MB)",
			Order: 7,
		},
	},
}

func GetVMId(vmId string) (float32, error) {
	vmIdNumeric, err := strconv.ParseFloat(vmId, 32)
	if err != nil {
		err := fmt.Errorf("invalid VM ID: '%s'", vmId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return float32(vmIdNumeric), nil
}

func VMGet(ctx context.Context, vmId string) error {
	logger.Get().Info().Msgf("Getting VM '%s'", vmId)

	vmIdNumeric, err := GetVMId(vmId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	vmInfo, httpRes, err := client.VMAPI.GetVM(ctx, vmIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(vmInfo, &vmPrintConfig)
}

func VMPowerStatus(ctx context.Context, vmId string) error {
	logger.Get().Info().Msgf("Getting power status for VM '%s'", vmId)

	vmIdNumeric, err := GetVMId(vmId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	powerStatus, httpRes, err := client.VMAPI.GetVMPowerStatus(ctx, vmIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Power status for VM '%s' is '%s'", vmId, powerStatus)

	return formatter.PrintResult(powerStatus, nil)
}

func VMStart(ctx context.Context, vmId string) error {
	logger.Get().Info().Msgf("Starting VM '%s'", vmId)

	vmIdNumeric, err := GetVMId(vmId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.VMAPI.StartVM(ctx, vmIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("VM '%s' started", vmId)

	return nil
}

func VMShutdown(ctx context.Context, vmId string) error {
	logger.Get().Info().Msgf("Shutting down VM '%s'", vmId)

	vmIdNumeric, err := GetVMId(vmId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.VMAPI.ShutdownVM(ctx, vmIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("VM '%s' shutdown initiated", vmId)

	return nil
}

func VMReboot(ctx context.Context, vmId string) error {
	logger.Get().Info().Msgf("Rebooting VM '%s'", vmId)

	vmIdNumeric, err := GetVMId(vmId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.VMAPI.RebootVM(ctx, vmIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("VM '%s' reboot initiated", vmId)

	return nil
}

func VMUpdate(ctx context.Context, vmId string, config []byte) error {
	logger.Get().Info().Msgf("Updating VM '%s'", vmId)

	var updateConfig sdk.UpdateVM
	err := json.Unmarshal(config, &updateConfig)
	if err != nil {
		return err
	}

	vmIdNumeric, err := GetVMId(vmId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	vmInfo, httpRes, err := client.VMAPI.UpdateVM(ctx, vmIdNumeric).UpdateVM(updateConfig).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("VM '%s' updated", vmId)

	return formatter.PrintResult(vmInfo, &vmPrintConfig)
}

func VMRemoteConsoleInfo(ctx context.Context, vmId string) error {
	logger.Get().Info().Msgf("Getting remote console info for VM '%s'", vmId)

	vmIdNumeric, err := GetVMId(vmId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	consoleInfo, httpRes, err := client.VMAPI.GetVMRemoteConsoleInfo(ctx, vmIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(consoleInfo, &formatter.PrintConfig{
		FieldsConfig: map[string]formatter.RecordFieldConfig{
			"ActiveConnections": {
				Title: "Active Connections",
				Order: 1,
			},
			"ConsoleUrl": {
				Title: "Console URL",
				Order: 2,
			},
		},
	})
}
