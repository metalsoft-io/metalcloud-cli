package file_share

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

func FileShareSnapshotList(ctx context.Context, infrastructureIdOrLabel string, fileShareId string) error {
	logger.Get().Info().Msgf("Listing snapshots for file share '%s' in infrastructure '%s'", fileShareId, infrastructureIdOrLabel)

	infrastructureInfo, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	fileShareIdNumeric, err := getFileShareId(fileShareId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	snapshots, httpRes, err := client.FileShareAPI.GetFileShareSnapshots(ctx, infrastructureInfo.Id, fileShareIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(snapshots, &snapshotPrintConfig)
}

func FileShareSnapshotCreate(ctx context.Context, infrastructureIdOrLabel string, fileShareId string) error {
	logger.Get().Info().Msgf("Creating snapshot for file share '%s' in infrastructure '%s'", fileShareId, infrastructureIdOrLabel)

	infrastructureInfo, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	fileShareIdNumeric, err := getFileShareId(fileShareId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	snapshot, httpRes, err := client.FileShareAPI.CreateFileShareSnapshot(ctx, infrastructureInfo.Id, fileShareIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	fmt.Printf("Snapshot '%s' created for file share %s\n", snapshot.Name, fileShareId)
	return formatter.PrintResult(snapshot, &snapshotPrintConfig)
}

func FileShareSnapshotDelete(ctx context.Context, infrastructureIdOrLabel string, fileShareId string, snapshotName string) error {
	logger.Get().Info().Msgf("Deleting snapshot '%s' for file share '%s' in infrastructure '%s'", snapshotName, fileShareId, infrastructureIdOrLabel)

	infrastructureInfo, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	fileShareIdNumeric, err := getFileShareId(fileShareId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.FileShareAPI.
		DeleteFileShareSnapshot(ctx, infrastructureInfo.Id, fileShareIdNumeric).
		DeleteFileShareSnapshot(sdk.DeleteFileShareSnapshot{
			Name: snapshotName,
		}).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	fmt.Printf("Snapshot '%s' deleted for file share %s\n", snapshotName, fileShareId)
	return nil
}

func FileShareSnapshotRestore(ctx context.Context, infrastructureIdOrLabel string, fileShareId string, snapshotName string) error {
	logger.Get().Info().Msgf("Restoring snapshot '%s' for file share '%s' in infrastructure '%s'", snapshotName, fileShareId, infrastructureIdOrLabel)

	infrastructureInfo, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	fileShareIdNumeric, err := getFileShareId(fileShareId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.FileShareAPI.
		RestoreFileShareToSnapshot(ctx, infrastructureInfo.Id, fileShareIdNumeric).
		RestoreFileShareSnapshot(sdk.RestoreFileShareSnapshot{
			Name: snapshotName,
		}).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	fmt.Printf("File share %s restored to snapshot '%s'\n", fileShareId, snapshotName)
	return nil
}
