package file_share

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

func FileShareList(ctx context.Context, infrastructureIdOrLabel string, filterStatus string) error {
	logger.Get().Info().Msgf("Listing file shares for infrastructure '%s'", infrastructureIdOrLabel)

	// Get the infrastructure ID from ID or label
	infrastructureInfo, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.FileShareAPI.GetInfrastructureFileShares(ctx, infrastructureInfo.Id)

	if filterStatus != "" {
		request = request.FilterServiceStatus([]string{"$eq:" + filterStatus})
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

	fileShare, httpRes, err := client.FileShareAPI.GetInfrastructureFileShare(ctx, infrastructureInfo.Id, fileShareIdNumeric).Execute()
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
	err = json.Unmarshal(config, &fileShareConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	fileShare, httpRes, err := client.FileShareAPI.
		CreateInfrastructureFileShare(ctx, infrastructureInfo.Id).
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

	fileShareIdNumeric, revision, err := getFileShareIdAndRevision(ctx, infrastructureInfo.Id, fileShareId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.FileShareAPI.
		DeleteFileShare(ctx, infrastructureInfo.Id, fileShareIdNumeric).
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

	fileShareIdNumeric, revision, err := getFileShareIdAndRevision(ctx, infrastructureInfo.Id, fileShareId)
	if err != nil {
		return err
	}

	var fileShareConfigUpdate sdk.UpdateFileShare
	err = json.Unmarshal(config, &fileShareConfigUpdate)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	fileShare, httpRes, err := client.FileShareAPI.
		UpdateFileShareConfig(ctx, infrastructureInfo.Id, fileShareIdNumeric).
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
	err = json.Unmarshal(config, &fileShareMetaUpdate)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	fileShare, httpRes, err := client.FileShareAPI.
		PatchFileShareMeta(ctx, infrastructureInfo.Id, fileShareIdNumeric).
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
		GetFileShareHosts(ctx, infrastructureInfo.Id, fileShareIdNumeric).
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
	err = json.Unmarshal(config, &hostsUpdate)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	hosts, httpRes, err := client.FileShareAPI.
		UpdateFileShareInstanceArrayHostsBulk(ctx, infrastructureInfo.Id, fileShareIdNumeric).
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
		GetFileShareConfigInfo(ctx, infrastructureInfo.Id, fileShareIdNumeric).
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

func getFileShareId(fileShareId string) (float32, error) {
	fileShareIdNumeric, err := strconv.ParseFloat(fileShareId, 32)
	if err != nil {
		err := fmt.Errorf("invalid file share ID: '%s'", fileShareId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return float32(fileShareIdNumeric), nil
}

func getFileShareIdAndRevision(ctx context.Context, infrastructureId float32, fileShareId string) (float32, string, error) {
	fileShareIdNumeric, err := getFileShareId(fileShareId)
	if err != nil {
		return 0, "", err
	}

	client := api.GetApiClient(ctx)

	fileShare, httpRes, err := client.FileShareAPI.GetInfrastructureFileShare(ctx, infrastructureId, float32(fileShareIdNumeric)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return 0, "", err
	}

	return float32(fileShareIdNumeric), strconv.Itoa(int(fileShare.Revision)), nil
}
