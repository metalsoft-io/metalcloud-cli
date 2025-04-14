package drive

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/internal/infrastructure"
	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var drivePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"Label": {
			MaxWidth: 30,
			Order:    2,
		},
		"SizeMBytes": {
			Title: "Size (MB)",
			Order: 3,
		},
		"StoragePoolId": {
			Title: "Pool ID",
			Order: 4,
		},
		"InfrastructureId": {
			Title: "Infra ID",
			Order: 5,
		},
		"ServiceStatus": {
			Title:       "Status",
			Transformer: formatter.FormatStatusValue,
			Order:       6,
		},
		"Config": {
			Hidden: true,
			InnerFields: map[string]formatter.RecordFieldConfig{
				"DeployStatus": {
					Title:       "Deploy Status",
					Transformer: formatter.FormatStatusValue,
					Order:       7,
				},
				"DeployType": {
					Title: "Deploy Type",
					Order: 8,
				},
			},
		},
		"WWN": {
			Title: "WWN",
			Order: 9,
		},
	},
}

var driveHostsPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"ServerInstanceGroupId": {
			Title: "Group ID",
			Order: 1,
		},
		"ServerInstanceId": {
			Title: "Server ID",
			Order: 2,
		},
		"ServerId": {
			Title: "HW Server ID",
			Order: 3,
		},
		"Status": {
			Title:       "Status",
			Transformer: formatter.FormatStatusValue,
			Order:       4,
		},
		"TapeEnabled": {
			Title: "Tape Enabled",
			Order: 5,
		},
	},
}

func DriveList(ctx context.Context, infrastructureIdOrLabel string, filterStatus string) error {
	logger.Get().Info().Msgf("Listing drives for infrastructure '%s'", infrastructureIdOrLabel)

	// Get the infrastructure ID from ID or label
	infrastructureInfo, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.DriveAPI.GetInfrastructureDrives(ctx, infrastructureInfo.Id)

	if filterStatus != "" {
		request = request.FilterServiceStatus([]string{"$eq:" + filterStatus})
	}

	driveList, httpRes, err := request.Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(driveList, &drivePrintConfig)
}

func DriveGet(ctx context.Context, infrastructureIdOrLabel string, driveId string) error {
	logger.Get().Info().Msgf("Get drive '%s' details", driveId)

	// Get the infrastructure ID from ID or label
	infrastructureInfo, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	driveIdNumeric, err := getDriveId(driveId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	drive, httpRes, err := client.DriveAPI.GetInfrastructureDrive(ctx, infrastructureInfo.Id, driveIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(drive, &drivePrintConfig)
}

func DriveCreate(ctx context.Context, infrastructureIdOrLabel string, config []byte) error {
	logger.Get().Info().Msgf("Creating drive for infrastructure '%s'", infrastructureIdOrLabel)

	// Get the infrastructure ID from ID or label
	infrastructureInfo, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	var driveConfig sdk.CreateSharedDrive
	err = json.Unmarshal(config, &driveConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	drive, httpRes, err := client.DriveAPI.
		CreateDrive(ctx, infrastructureInfo.Id).
		CreateSharedDrive(driveConfig).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(drive, &drivePrintConfig)
}

func DriveDelete(ctx context.Context, infrastructureIdOrLabel string, driveId string) error {
	logger.Get().Info().Msgf("Deleting drive '%s'", driveId)

	// Get the infrastructure ID from ID or label
	infrastructureInfo, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	driveIdNumeric, revision, err := getDriveIdAndRevision(ctx, infrastructureInfo.Id, driveId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.DriveAPI.
		DeleteDrive(ctx, infrastructureInfo.Id, driveIdNumeric).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Drive '%s' deleted", driveId)
	return nil
}

func DriveUpdateConfig(ctx context.Context, infrastructureIdOrLabel string, driveId string, config []byte) error {
	logger.Get().Info().Msgf("Updating drive '%s' configuration", driveId)

	// Get the infrastructure ID from ID or label
	infrastructureInfo, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	driveIdNumeric, revision, err := getDriveIdAndRevision(ctx, infrastructureInfo.Id, driveId)
	if err != nil {
		return err
	}

	var driveConfigUpdate sdk.UpdateSharedDrive
	err = json.Unmarshal(config, &driveConfigUpdate)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	drive, httpRes, err := client.DriveAPI.
		PatchDriveConfig(ctx, infrastructureInfo.Id, driveIdNumeric).
		UpdateSharedDrive(driveConfigUpdate).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(drive, &drivePrintConfig)
}

func DriveUpdateMeta(ctx context.Context, infrastructureIdOrLabel string, driveId string, config []byte) error {
	logger.Get().Info().Msgf("Updating drive '%s' metadata", driveId)

	// Get the infrastructure ID from ID or label
	infrastructureInfo, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	driveIdNumeric, err := getDriveId(driveId)
	if err != nil {
		return err
	}

	var driveMetaUpdate sdk.UpdateSharedDriveMeta
	err = json.Unmarshal(config, &driveMetaUpdate)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	drive, httpRes, err := client.DriveAPI.
		PatchDriveMeta(ctx, infrastructureInfo.Id, driveIdNumeric).
		UpdateSharedDriveMeta(driveMetaUpdate).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(drive, &drivePrintConfig)
}

func DriveGetHosts(ctx context.Context, infrastructureIdOrLabel string, driveId string) error {
	logger.Get().Info().Msgf("Getting hosts for drive '%s'", driveId)

	// Get the infrastructure ID from ID or label
	infrastructureInfo, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	driveIdNumeric, err := getDriveId(driveId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	hosts, httpRes, err := client.DriveAPI.
		GetDriveHosts(ctx, infrastructureInfo.Id, driveIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(hosts, &driveHostsPrintConfig)
}

func DriveUpdateHosts(ctx context.Context, infrastructureIdOrLabel string, driveId string, config []byte) error {
	logger.Get().Info().Msgf("Updating hosts for drive '%s'", driveId)

	// Get the infrastructure ID from ID or label
	infrastructureInfo, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	driveIdNumeric, err := getDriveId(driveId)
	if err != nil {
		return err
	}

	var hostsUpdate sdk.SharedDriveHostsModifyBulk
	err = json.Unmarshal(config, &hostsUpdate)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	hosts, httpRes, err := client.DriveAPI.
		UpdateDriveServerInstanceGroupHostsBulk(ctx, infrastructureInfo.Id, driveIdNumeric).
		SharedDriveHostsModifyBulk(hostsUpdate).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(hosts, &driveHostsPrintConfig)
}

func DriveGetConfigInfo(ctx context.Context, infrastructureIdOrLabel string, driveId string) error {
	logger.Get().Info().Msgf("Getting configuration info for drive '%s'", driveId)

	// Get the infrastructure ID from ID or label
	infrastructureInfo, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	driveIdNumeric, err := getDriveId(driveId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	configInfo, httpRes, err := client.DriveAPI.
		GetDriveConfigInfo(ctx, infrastructureInfo.Id, driveIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(configInfo, &formatter.PrintConfig{
		FieldsConfig: map[string]formatter.RecordFieldConfig{
			"Label": {
				Title: "Label",
				Order: 1,
			},
			"SizeMBytes": {
				Title: "Size (MB)",
				Order: 2,
			},
			"DriveDeployId": {
				Title: "Deploy ID",
				Order: 3,
			},
			"DeployType": {
				Title: "Deploy Type",
				Order: 4,
			},
			"DeployStatus": {
				Title:       "Deploy Status",
				Transformer: formatter.FormatStatusValue,
				Order:       5,
			},
			"LogicalNetworkId": {
				Title: "Logical Network ID",
				Order: 6,
			},
			"UpdatedTimestamp": {
				Title:       "Updated",
				Transformer: formatter.FormatDateTimeValue,
				Order:       7,
			},
		},
	})
}

func getDriveId(driveId string) (float32, error) {
	driveIdNumeric, err := strconv.ParseFloat(driveId, 32)
	if err != nil {
		err := fmt.Errorf("invalid drive ID: '%s'", driveId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return float32(driveIdNumeric), nil
}

func getDriveIdAndRevision(ctx context.Context, infrastructureId float32, driveId string) (float32, string, error) {
	driveIdNumeric, err := getDriveId(driveId)
	if err != nil {
		return 0, "", err
	}

	client := api.GetApiClient(ctx)

	drive, httpRes, err := client.DriveAPI.GetInfrastructureDrive(ctx, infrastructureId, float32(driveIdNumeric)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return 0, "", err
	}

	return float32(driveIdNumeric), strconv.Itoa(int(drive.Revision)), nil
}
