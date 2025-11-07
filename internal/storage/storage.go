package storage

import (
	"context"
	"fmt"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var StoragePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"SiteId": {
			Title: "Site",
			Order: 2,
		},
		"Driver": {
			Title: "Driver",
			Order: 3,
		},
		"Technology": {
			Title: "Technology",
			Order: 4,
		},
		"Type": {
			Title: "Type",
			Order: 5,
		},
		"Name": {
			Title: "Name",
			Order: 6,
		},
		"Status": {
			Title:       "Status",
			Transformer: formatter.FormatStatusValue,
			Order:       7,
		},
	},
}

func StorageList(ctx context.Context, filterTechnology []string) error {
	logger.Get().Info().Msgf("Listing all storages")

	client := api.GetApiClient(ctx)

	request := client.StorageAPI.GetStorages(ctx)

	if len(filterTechnology) > 0 {
		request = request.FilterTechnologies(utils.ProcessFilterStringSlice(filterTechnology))
	}

	storageList, httpRes, err := request.SortBy([]string{"id:ASC"}).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(storageList, &StoragePrintConfig)
}

func StorageGet(ctx context.Context, storageId string) error {
	logger.Get().Info().Msgf("Get storage %s details", storageId)

	storageIdNumeric, err := getStorageId(storageId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	storage, httpRes, err := client.StorageAPI.GetStorage(ctx, storageIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(storage, &StoragePrintConfig)
}

func StorageCreate(ctx context.Context, config []byte) error {
	logger.Get().Info().Msgf("Creating new storage")

	client := api.GetApiClient(ctx)

	var createStorage sdk.CreateStorage
	err := utils.UnmarshalContent(config, &createStorage)
	if err != nil {
		return err
	}

	response, httpRes, err := client.StorageAPI.CreateStorage(ctx).
		CreateStorage(createStorage).
		Execute()

	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Storage created with ID: %d", int(response.Id))
	return nil
}

func StorageDelete(ctx context.Context, storageId string) error {
	logger.Get().Info().Msgf("Deleting storage %s", storageId)

	storageIdNumeric, err := getStorageId(storageId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.StorageAPI.DeleteStorage(ctx, storageIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Storage %s deleted successfully", storageId)
	return nil
}

func StorageGetCredentials(ctx context.Context, storageId string) error {
	logger.Get().Info().Msgf("Getting credentials for storage %s", storageId)

	storageIdNumeric, err := getStorageId(storageId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	credentials, httpRes, err := client.StorageAPI.GetStorageCredentials(ctx, storageIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(credentials, nil)
}

func StorageGetDrives(ctx context.Context, storageId string, limit float32, page float32) error {
	logger.Get().Info().Msgf("Getting drives for storage %s", storageId)

	storageIdNumeric, err := getStorageId(storageId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.StorageAPI.GetStorageDrives(ctx, storageIdNumeric)

	// Set pagination if provided
	if limit > 0 {
		request = request.Limit(limit)
	}

	if page > 0 {
		request = request.Page(page)
	}

	drives, httpRes, err := request.Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(drives, nil)
}

func StorageGetFileShares(ctx context.Context, storageId string, limit float32, page float32) error {
	logger.Get().Info().Msgf("Getting file shares for storage %s", storageId)

	storageIdNumeric, err := getStorageId(storageId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.StorageAPI.GetStorageFileShares(ctx, storageIdNumeric)

	// Set pagination if provided
	if limit > 0 {
		request = request.Limit(limit)
	}

	if page > 0 {
		request = request.Page(page)
	}

	fileShares, httpRes, err := request.Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(fileShares, nil)
}

func StorageGetBuckets(ctx context.Context, storageId string, limit float32, page float32) error {
	logger.Get().Info().Msgf("Getting buckets for storage %s", storageId)

	storageIdNumeric, err := getStorageId(storageId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.StorageAPI.GetStorageBuckets(ctx, storageIdNumeric)

	// Set pagination if provided
	if limit > 0 {
		request = request.Limit(limit)
	}

	if page > 0 {
		request = request.Page(page)
	}

	buckets, httpRes, err := request.Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(buckets, nil)
}

func StorageGetNetworkDeviceConfigurations(ctx context.Context, storageId string) error {
	logger.Get().Info().Msgf("Getting network device configurations for storage %s", storageId)

	storageIdNumeric, err := getStorageId(storageId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	configs, httpRes, err := client.StorageAPI.GetStorageNetworkDeviceConfigurations(ctx, storageIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(configs, nil)
}

func StorageConfigExample(ctx context.Context) error {
	storageConfiguration := sdk.CreateStorage{
		UserId:              sdk.PtrFloat32(1),
		SiteId:              1,
		Driver:              "netapp",
		Technologies:        []string{"block"},
		Type:                sdk.PtrString("type"),
		Name:                "name",
		IscsiHost:           sdk.PtrString("iscsiHost"),
		IscsiPort:           sdk.PtrFloat32(234),
		ManagementHost:      "storage.host",
		Username:            sdk.PtrString("username"),
		Password:            sdk.PtrString("password"),
		InMaintenance:       sdk.PtrFloat32(0),
		TargetIQN:           sdk.PtrString("targetIQN"),
		SharedDrivePriority: sdk.PtrFloat32(1),
		AlternateSanIPs:     []string{"1.2.3.4"},
		Tags:                []string{"tag1", "tag2"},
		SubnetType:          "subnetType",
		Options: &sdk.UpdateStorageOptions{
			EnableDataReduction:         sdk.PtrBool(true),
			EnableAdvancedDeduplication: sdk.PtrBool(false),
			VolumeName:                  sdk.PtrString("volumeName"),
			ArrayId:                     sdk.PtrString("1"),
			DirectorId:                  sdk.PtrString("1"),
			S3Hostname:                  sdk.PtrString("s3Hostname"),
			S3Port:                      sdk.PtrFloat32(1),
			FibreChannelEnabled:         sdk.PtrBool(false),
		},
	}

	return formatter.PrintResult(storageConfiguration, nil)
}

func getStorageId(storageId string) (float32, error) {
	storageIdNumeric, err := strconv.ParseFloat(storageId, 32)
	if err != nil {
		err := fmt.Errorf("invalid storage ID: '%s'", storageId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return float32(storageIdNumeric), nil
}
