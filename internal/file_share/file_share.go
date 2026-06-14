package file_share

import (
	"context"
	"fmt"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/internal/infrastructure"
	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var fileSharePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"Label": {
			MaxWidth: 30,
			Order:    2,
		},
		"SizeGB": {
			Title: "Size (GB)",
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
		"Endpoint": {
			Title: "Endpoint",
			Order: 9,
		},
	},
}

var fileShareHostsPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"InstanceArrayId": {
			Title: "Array ID",
			Order: 1,
		},
		"InstanceId": {
			Title: "Instance ID",
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
		"AccessRight": {
			Title: "Access Right",
			Order: 5,
		},
	},
}

func FileShareList(ctx context.Context, infrastructureIdOrLabel string, filterStatus []string) error {
	logger.Get().Info().Msgf("Listing file shares for infrastructure '%s'", infrastructureIdOrLabel)

	// Get the infrastructure ID from ID or label
	infrastructureInfo, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.FileShareAPI.GetInfrastructureFileShares(ctx, int64(infrastructureInfo.Id))

	if len(filterStatus) > 0 {
		request = request.FilterServiceStatus(utils.ProcessFilterStringSlice(filterStatus))
	}

	fileShareList, httpRes, err := request.Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(fileShareList, &fileSharePrintConfig)
}

func FileShareGet(ctx context.Context, infrastructureIdOrLabel string, fileShareId string) error {
	logger.Get().Info().Msgf("Get file share '%s' details", fileShareId)

	// Get the infrastructure ID from ID or label
	infrastructureInfo, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	fileShareIdNumeric, err := getFileShareId(fileShareId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	fileShare, httpRes, err := client.FileShareAPI.GetInfrastructureFileShare(ctx, int64(infrastructureInfo.Id), fileShareIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(fileShare, &fileSharePrintConfig)
}

func FileShareCreate(ctx context.Context, infrastructureIdOrLabel string, config []byte) error {
	logger.Get().Info().Msgf("Creating file share for infrastructure '%s'", infrastructureIdOrLabel)

	// Get the infrastructure ID from ID or label
	infrastructureInfo, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	var fileShareConfig sdk.CreateFileShare
	err = utils.UnmarshalContent(config, &fileShareConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	fileShare, httpRes, err := client.FileShareAPI.
		CreateInfrastructureFileShare(ctx, int64(infrastructureInfo.Id)).
		CreateFileShare(fileShareConfig).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(fileShare, &fileSharePrintConfig)
}

func FileShareDelete(ctx context.Context, infrastructureIdOrLabel string, fileShareId string) error {
	logger.Get().Info().Msgf("Deleting file share '%s'", fileShareId)

	// Get the infrastructure ID from ID or label
	infrastructureInfo, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	fileShareIdNumeric, revision, err := getFileShareIdAndRevision(ctx, int64(infrastructureInfo.Id), fileShareId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.FileShareAPI.
		DeleteFileShare(ctx, int64(infrastructureInfo.Id), fileShareIdNumeric).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("File share '%s' deleted", fileShareId)
	return nil
}

func FileShareUpdateConfig(ctx context.Context, infrastructureIdOrLabel string, fileShareId string, config []byte) error {
	logger.Get().Info().Msgf("Updating file share '%s' configuration", fileShareId)

	// Get the infrastructure ID from ID or label
	infrastructureInfo, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	fileShareIdNumeric, revision, err := getFileShareIdAndRevision(ctx, int64(infrastructureInfo.Id), fileShareId)
	if err != nil {
		return err
	}

	var fileShareConfigUpdate sdk.UpdateFileShare
	err = utils.UnmarshalContent(config, &fileShareConfigUpdate)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	fileShare, httpRes, err := client.FileShareAPI.
		UpdateFileShareConfig(ctx, int64(infrastructureInfo.Id), fileShareIdNumeric).
		UpdateFileShare(fileShareConfigUpdate).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(fileShare, &fileSharePrintConfig)
}

func FileShareUpdateMeta(ctx context.Context, infrastructureIdOrLabel string, fileShareId string, config []byte) error {
	logger.Get().Info().Msgf("Updating file share '%s' metadata", fileShareId)

	// Get the infrastructure ID from ID or label
	infrastructureInfo, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	fileShareIdNumeric, err := getFileShareId(fileShareId)
	if err != nil {
		return err
	}

	var fileShareMetaUpdate sdk.UpdateFileShareMeta
	err = utils.UnmarshalContent(config, &fileShareMetaUpdate)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	fileShare, httpRes, err := client.FileShareAPI.
		PatchFileShareMeta(ctx, int64(infrastructureInfo.Id), fileShareIdNumeric).
		UpdateFileShareMeta(fileShareMetaUpdate).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(fileShare, &fileSharePrintConfig)
}

func FileShareGetHosts(ctx context.Context, infrastructureIdOrLabel string, fileShareId string) error {
	logger.Get().Info().Msgf("Getting hosts for file share '%s'", fileShareId)

	// Get the infrastructure ID from ID or label
	infrastructureInfo, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	fileShareIdNumeric, err := getFileShareId(fileShareId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	hosts, httpRes, err := client.FileShareAPI.
		GetFileShareHosts(ctx, int64(infrastructureInfo.Id), fileShareIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(hosts, &fileShareHostsPrintConfig)
}

func FileShareUpdateHosts(ctx context.Context, infrastructureIdOrLabel string, fileShareId string, config []byte) error {
	logger.Get().Info().Msgf("Updating hosts for file share '%s'", fileShareId)

	// Get the infrastructure ID from ID or label
	infrastructureInfo, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	fileShareIdNumeric, err := getFileShareId(fileShareId)
	if err != nil {
		return err
	}

	var hostsUpdate sdk.FileShareHostsModifyBulk
	err = utils.UnmarshalContent(config, &hostsUpdate)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	hosts, httpRes, err := client.FileShareAPI.
		UpdateFileShareInstanceArrayHostsBulk(ctx, int64(infrastructureInfo.Id), fileShareIdNumeric).
		FileShareHostsModifyBulk(hostsUpdate).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(hosts, &fileShareHostsPrintConfig)
}

func FileShareGetConfigInfo(ctx context.Context, infrastructureIdOrLabel string, fileShareId string) error {
	logger.Get().Info().Msgf("Getting configuration info for file share '%s'", fileShareId)

	// Get the infrastructure ID from ID or label
	infrastructureInfo, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	fileShareIdNumeric, err := getFileShareId(fileShareId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	configInfo, httpRes, err := client.FileShareAPI.
		GetFileShareConfigInfo(ctx, int64(infrastructureInfo.Id), fileShareIdNumeric).
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
			"SizeGB": {
				Title: "Size (GB)",
				Order: 2,
			},
			"FileShareDeployId": {
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

func getFileShareId(fileShareId string) (int64, error) {
	fileShareIdNumeric, err := strconv.ParseInt(fileShareId, 10, 64)
	if err != nil {
		err := fmt.Errorf("invalid file share ID: '%s'", fileShareId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return fileShareIdNumeric, nil
}

func getFileShareIdAndRevision(ctx context.Context, infrastructureId int64, fileShareId string) (int64, string, error) {
	fileShareIdNumeric, err := getFileShareId(fileShareId)
	if err != nil {
		return 0, "", err
	}

	client := api.GetApiClient(ctx)

	fileShare, httpRes, err := client.FileShareAPI.GetInfrastructureFileShare(ctx, infrastructureId, fileShareIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return 0, "", err
	}

	return fileShareIdNumeric, strconv.Itoa(int(fileShare.Revision)), nil
}
