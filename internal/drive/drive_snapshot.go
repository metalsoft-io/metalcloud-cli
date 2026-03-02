package drive

import (
	"context"
	"fmt"

	"github.com/metalsoft-io/metalcloud-cli/internal/infrastructure"
	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var snapshotPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Name": {
			Title: "Name",
			Order: 1,
		},
		"CreatedTimestamp": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       2,
		},
	},
}

func DriveSnapshotList(ctx context.Context, infrastructureIdOrLabel string, driveId string) error {
	logger.Get().Info().Msgf("Listing snapshots for drive '%s' in infrastructure '%s'", driveId, infrastructureIdOrLabel)

	infrastructureInfo, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	driveIdNumeric, err := getDriveId(driveId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	snapshots, httpRes, err := client.DriveAPI.GetDriveSnapshots(ctx, infrastructureInfo.Id, driveIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(snapshots, &snapshotPrintConfig)
}

func DriveSnapshotCreate(ctx context.Context, infrastructureIdOrLabel string, driveId string) error {
	logger.Get().Info().Msgf("Creating snapshot for drive '%s' in infrastructure '%s'", driveId, infrastructureIdOrLabel)

	infrastructureInfo, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	driveIdNumeric, err := getDriveId(driveId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	snapshot, httpRes, err := client.DriveAPI.CreateDriveSnapshot(ctx, infrastructureInfo.Id, driveIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	fmt.Printf("Snapshot '%s' created for drive %s\n", snapshot.Name, driveId)
	return formatter.PrintResult(snapshot, &snapshotPrintConfig)
}

func DriveSnapshotDelete(ctx context.Context, infrastructureIdOrLabel string, driveId string, snapshotName string) error {
	logger.Get().Info().Msgf("Deleting snapshot '%s' for drive '%s' in infrastructure '%s'", snapshotName, driveId, infrastructureIdOrLabel)

	infrastructureInfo, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	driveIdNumeric, err := getDriveId(driveId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.DriveAPI.
		DeleteDriveSnapshot(ctx, infrastructureInfo.Id, driveIdNumeric).
		DeleteSharedDriveSnapshot(sdk.DeleteSharedDriveSnapshot{
			Name: snapshotName,
		}).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	fmt.Printf("Snapshot '%s' deleted for drive %s\n", snapshotName, driveId)
	return nil
}

func DriveSnapshotRestore(ctx context.Context, infrastructureIdOrLabel string, driveId string, snapshotName string) error {
	logger.Get().Info().Msgf("Restoring snapshot '%s' for drive '%s' in infrastructure '%s'", snapshotName, driveId, infrastructureIdOrLabel)

	infrastructureInfo, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	driveIdNumeric, err := getDriveId(driveId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.DriveAPI.
		RestoreDriveToSnapshot(ctx, infrastructureInfo.Id, driveIdNumeric).
		RestoreSharedDriveSnapshot(sdk.RestoreSharedDriveSnapshot{
			Name: snapshotName,
		}).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	fmt.Printf("Drive %s restored to snapshot '%s'\n", driveId, snapshotName)
	return nil
}
