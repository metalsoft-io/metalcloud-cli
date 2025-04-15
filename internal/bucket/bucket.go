package bucket

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

var bucketPrintConfig = formatter.PrintConfig{
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
		"Subdomain": {
			Title: "Subdomain",
			Order: 9,
		},
		"SubdomainPermanent": {
			Title: "Permanent Subdomain",
			Order: 10,
		},
		"Endpoint": {
			Title: "Endpoint",
			Order: 11,
		},
	},
}

func BucketList(ctx context.Context, infrastructureIdOrLabel string, filterStatus string) error {
	logger.Get().Info().Msgf("Listing buckets for infrastructure '%s'", infrastructureIdOrLabel)

	// Get the infrastructure ID from ID or label
	infrastructureInfo, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.BucketAPI.GetInfrastructureBuckets(ctx, infrastructureInfo.Id)

	if filterStatus != "" {
		request = request.FilterServiceStatus([]string{"$eq:" + filterStatus})
	}

	bucketList, httpRes, err := request.Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(bucketList, &bucketPrintConfig)
}

func BucketGet(ctx context.Context, infrastructureIdOrLabel string, bucketId string) error {
	logger.Get().Info().Msgf("Get bucket '%s' details", bucketId)

	// Get the infrastructure ID from ID or label
	infrastructureInfo, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	bucketIdNumeric, err := getBucketId(bucketId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	bucket, httpRes, err := client.BucketAPI.GetInfrastructureBucket(ctx, infrastructureInfo.Id, bucketIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(bucket, &bucketPrintConfig)
}

func BucketCreate(ctx context.Context, infrastructureIdOrLabel string, config []byte) error {
	logger.Get().Info().Msgf("Creating bucket for infrastructure '%s'", infrastructureIdOrLabel)

	// Get the infrastructure ID from ID or label
	infrastructureInfo, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	var bucketConfig sdk.CreateBucket
	err = json.Unmarshal(config, &bucketConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	bucket, httpRes, err := client.BucketAPI.
		CreateInfrastructureBucket(ctx, infrastructureInfo.Id).
		CreateBucket(bucketConfig).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(bucket, &bucketPrintConfig)
}

func BucketDelete(ctx context.Context, infrastructureIdOrLabel string, bucketId string) error {
	logger.Get().Info().Msgf("Deleting bucket '%s'", bucketId)

	// Get the infrastructure ID from ID or label
	infrastructureInfo, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	bucketIdNumeric, revision, err := getBucketIdAndRevision(ctx, infrastructureInfo.Id, bucketId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.BucketAPI.
		DeleteBucket(ctx, infrastructureInfo.Id, bucketIdNumeric).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Bucket '%s' deleted", bucketId)
	return nil
}

func BucketUpdateConfig(ctx context.Context, infrastructureIdOrLabel string, bucketId string, config []byte) error {
	logger.Get().Info().Msgf("Updating bucket '%s' configuration", bucketId)

	// Get the infrastructure ID from ID or label
	infrastructureInfo, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	bucketIdNumeric, revision, err := getBucketIdAndRevision(ctx, infrastructureInfo.Id, bucketId)
	if err != nil {
		return err
	}

	var bucketConfigUpdate sdk.UpdateBucket
	err = json.Unmarshal(config, &bucketConfigUpdate)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	bucket, httpRes, err := client.BucketAPI.
		UpdateBucket(ctx, infrastructureInfo.Id, bucketIdNumeric).
		UpdateBucket(bucketConfigUpdate).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(bucket, &bucketPrintConfig)
}

func BucketUpdateMeta(ctx context.Context, infrastructureIdOrLabel string, bucketId string, config []byte) error {
	logger.Get().Info().Msgf("Updating bucket '%s' metadata", bucketId)

	// Get the infrastructure ID from ID or label
	infrastructureInfo, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	bucketIdNumeric, err := getBucketId(bucketId)
	if err != nil {
		return err
	}

	var bucketMetaUpdate sdk.UpdateBucketMeta
	err = json.Unmarshal(config, &bucketMetaUpdate)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	bucket, httpRes, err := client.BucketAPI.
		UpdateBucketMeta(ctx, infrastructureInfo.Id, bucketIdNumeric).
		UpdateBucketMeta(bucketMetaUpdate).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(bucket, &bucketPrintConfig)
}

func BucketGetConfigInfo(ctx context.Context, infrastructureIdOrLabel string, bucketId string) error {
	logger.Get().Info().Msgf("Getting configuration info for bucket '%s'", bucketId)

	// Get the infrastructure ID from ID or label
	infrastructureInfo, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	bucketIdNumeric, err := getBucketId(bucketId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	configInfo, httpRes, err := client.BucketAPI.
		GetBucketConfigInfo(ctx, infrastructureInfo.Id, bucketIdNumeric).
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
			"DeployType": {
				Title: "Deploy Type",
				Order: 3,
			},
			"DeployStatus": {
				Title:       "Deploy Status",
				Transformer: formatter.FormatStatusValue,
				Order:       4,
			},
			"LogicalNetworkId": {
				Title: "Logical Network ID",
				Order: 5,
			},
			"UpdatedTimestamp": {
				Title:       "Updated",
				Transformer: formatter.FormatDateTimeValue,
				Order:       6,
			},
		},
	})
}

func BucketGetCredentials(ctx context.Context, infrastructureIdOrLabel string, bucketId string) error {
	logger.Get().Info().Msgf("Getting credentials for bucket '%s'", bucketId)

	// Get the infrastructure ID from ID or label
	infrastructureInfo, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	bucketIdNumeric, err := getBucketId(bucketId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	credentials, httpRes, err := client.BucketAPI.
		GetBucketCredentials(ctx, infrastructureInfo.Id, bucketIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(credentials, &formatter.PrintConfig{
		FieldsConfig: map[string]formatter.RecordFieldConfig{
			"accessKeyId": {
				Title: "Access Key ID",
				Order: 1,
			},
			"secretKey": {
				Title: "Secret Key",
				Order: 2,
			},
			"endpoint": {
				Title: "Endpoint",
				Order: 3,
			},
		},
	})
}

func getBucketId(bucketId string) (float32, error) {
	bucketIdNumeric, err := strconv.ParseFloat(bucketId, 32)
	if err != nil {
		err := fmt.Errorf("invalid bucket ID: '%s'", bucketId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return float32(bucketIdNumeric), nil
}

func getBucketIdAndRevision(ctx context.Context, infrastructureId float32, bucketId string) (float32, string, error) {
	bucketIdNumeric, err := getBucketId(bucketId)
	if err != nil {
		return 0, "", err
	}

	client := api.GetApiClient(ctx)

	bucket, httpRes, err := client.BucketAPI.GetInfrastructureBucket(ctx, infrastructureId, bucketIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return 0, "", err
	}

	return bucketIdNumeric, strconv.Itoa(int(bucket.Revision)), nil
}
