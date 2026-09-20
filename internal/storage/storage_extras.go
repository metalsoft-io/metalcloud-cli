package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var storageInterfacePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"StorageId": {
			Title: "Storage",
			Order: 2,
		},
		"Name": {
			Title:    "Name",
			MaxWidth: 30,
			Order:    3,
		},
		"Protocols": {
			Title:       "Protocols",
			Transformer: formatter.FormatStringListValue,
			Order:       4,
		},
		"NodeIds": {
			Title:       "Nodes",
			Transformer: formatter.FormatStringListValue,
			Order:       5,
		},
		"IsUplink": {
			Title:       "Uplink",
			Transformer: formatter.FormatBooleanValue,
			Order:       6,
		},
		"UseForDeploys": {
			Title:       "Use For Deploys",
			Transformer: formatter.FormatBooleanValue,
			Order:       7,
		},
		"NetworkEquipmentInterfaceId": {
			Title: "Device Interface",
			Order: 8,
		},
	},
}

var storageScopedAccessUserPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"StorageId": {
			Title: "Storage",
			Order: 2,
		},
		"InfrastructureId": {
			Title: "Infrastructure",
			Order: 3,
		},
		"Username": {
			Title:    "Username",
			MaxWidth: 40,
			Order:    4,
		},
		"CreatedTimestamp": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       5,
		},
		"UpdatedTimestamp": {
			Title:       "Updated",
			Transformer: formatter.FormatDateTimeValue,
			Order:       6,
		},
	},
}

var storageScopedAccessUserCredentialsPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Username": {
			Title: "Username",
			Order: 1,
		},
		"Password": {
			Title: "Password",
			Order: 2,
		},
		"ApiToken": {
			Title:    "API Token",
			MaxWidth: 60,
			Order:    3,
		},
	},
}

var storageStatisticsPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"TotalSpaceGB": {
			Title: "Total (GB)",
			Order: 1,
		},
		"UsedSpaceGB": {
			Title: "Used (GB)",
			Order: 2,
		},
		"FreeSpaceGB": {
			Title: "Free (GB)",
			Order: 3,
		},
	},
}

var storagesStatisticsPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"ActiveCount": {
			Title: "Active",
			Order: 1,
		},
		"ReadyCount": {
			Title: "Ready",
			Order: 2,
		},
		"PendingCount": {
			Title: "Pending",
			Order: 3,
		},
		"MaintenanceCount": {
			Title: "Maintenance",
			Order: 4,
		},
		"ExperimentalCount": {
			Title: "Experimental",
			Order: 5,
		},
		"LowSpaceCount": {
			Title: "Low Space",
			Order: 6,
		},
		"UsedSpace": {
			Title: "Used Space",
			Order: 7,
		},
		"FreeSpace": {
			Title: "Free Space",
			Order: 8,
		},
		"Types": {
			Title:    "Types",
			MaxWidth: 40,
			Order:    9,
		},
	},
}

// storagesStatisticsRow is the table projection of the global storage
// statistics: the per-type counter map is reduced to a single "type: count"
// string because the tabular formatter renders map-valued fields as empty
// cells. The json and yaml formats keep the full response.
type storagesStatisticsRow struct {
	ActiveCount       float32
	ReadyCount        float32
	PendingCount      float32
	MaintenanceCount  float32
	ExperimentalCount float32
	LowSpaceCount     float32
	UsedSpace         float32
	FreeSpace         float32
	Types             string
}

// StorageUpdate updates a storage pool from a JSON/YAML configuration.
//
// The request and the response are exchanged as raw JSON: the typed SDK
// Storage model rejects live payloads (see storageRaw — the
// options.fibreChannel* type mismatch), so a typed PATCH would fail while
// decoding a successful response. The configuration is still validated by
// unmarshalling it into sdk.UpdateStorage first.
func StorageUpdate(ctx context.Context, storageId string, config []byte) error {
	logger.Get().Info().Msgf("Updating storage %s", storageId)

	storageIdNumeric, err := getStorageId(storageId)
	if err != nil {
		return err
	}

	var updateStorage sdk.UpdateStorage
	if err := utils.UnmarshalContent(config, &updateStorage); err != nil {
		return err
	}

	body, err := json.Marshal(updateStorage)
	if err != nil {
		return fmt.Errorf("failed to encode storage update: %w", err)
	}

	revision, err := storageRevision(ctx, storageIdNumeric)
	if err != nil {
		return err
	}

	responseBody, err := api.RawJSONRequest(
		ctx,
		http.MethodPatch,
		fmt.Sprintf("/api/v2/storages/%d", storageIdNumeric),
		body,
		api.IfMatchHeader(revision),
	)
	if err != nil {
		return err
	}

	return utils.PrintRawObject(responseBody, &StoragePrintConfig)
}

// storageRevision returns the current revision of a storage pool, read from
// the raw response body.
func storageRevision(ctx context.Context, storageIdNumeric int64) (string, error) {
	body, err := api.RawJSONRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("/api/v2/storages/%d", storageIdNumeric),
		nil,
		nil,
	)
	if err != nil {
		return "", err
	}

	object, err := utils.DecodeRawObject(body)
	if err != nil {
		return "", err
	}

	return utils.RevisionFromRaw(object), nil
}

// StorageGetInterfaces lists the interfaces of a storage pool.
func StorageGetInterfaces(ctx context.Context, storageId string) error {
	logger.Get().Info().Msgf("Getting interfaces for storage %s", storageId)

	storageIdNumeric, err := getStorageId(storageId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.StorageAPI.GetStorageInterfaces(ctx, storageIdNumeric).SortBy([]string{"id:ASC"})

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &storageInterfacePrintConfig)
}

// StorageGetInterface shows one interface of a storage pool.
func StorageGetInterface(ctx context.Context, storageId string, interfaceId string) error {
	logger.Get().Info().Msgf("Getting interface %s of storage %s", interfaceId, storageId)

	storageIdNumeric, interfaceIdNumeric, err := getStorageInterfaceIds(storageId, interfaceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	storageInterface, httpRes, err := client.StorageAPI.
		GetStorageInterface(ctx, storageIdNumeric, interfaceIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(storageInterface, &storageInterfacePrintConfig)
}

// StorageUpdateInterface updates one interface of a storage pool from a
// JSON/YAML configuration. The current revision is fetched first and sent as
// the If-Match entity tag.
func StorageUpdateInterface(ctx context.Context, storageId string, interfaceId string, config []byte) error {
	logger.Get().Info().Msgf("Updating interface %s of storage %s", interfaceId, storageId)

	storageIdNumeric, interfaceIdNumeric, err := getStorageInterfaceIds(storageId, interfaceId)
	if err != nil {
		return err
	}

	var updateInterface sdk.UpdateStorageInterface
	if err := utils.UnmarshalContent(config, &updateInterface); err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	current, httpRes, err := client.StorageAPI.
		GetStorageInterface(ctx, storageIdNumeric, interfaceIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}
	if current == nil {
		return fmt.Errorf("storage interface '%s' not found", interfaceId)
	}

	updated, httpRes, err := client.StorageAPI.
		UpdateStorageInterface(ctx, storageIdNumeric, interfaceIdNumeric).
		UpdateStorageInterface(updateInterface).
		IfMatch(strconv.FormatInt(current.Revision, 10)).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(updated, &storageInterfacePrintConfig)
}

// StorageGetScopedAccessUsers lists the scoped access users of a storage pool.
func StorageGetScopedAccessUsers(ctx context.Context, storageId string) error {
	logger.Get().Info().Msgf("Getting scoped access users for storage %s", storageId)

	storageIdNumeric, err := getStorageId(storageId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.StorageAPI.GetStorageScopedAccessUsers(ctx, storageIdNumeric).SortBy([]string{"id:ASC"})

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &storageScopedAccessUserPrintConfig)
}

// StorageGetScopedAccessUser shows one scoped access user of a storage pool.
func StorageGetScopedAccessUser(ctx context.Context, storageId string, userId string) error {
	logger.Get().Info().Msgf("Getting scoped access user %s of storage %s", userId, storageId)

	storageIdNumeric, userIdNumeric, err := getStorageScopedAccessUserIds(storageId, userId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	user, httpRes, err := client.StorageAPI.
		GetStorageScopedAccessUser(ctx, storageIdNumeric, userIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(user, &storageScopedAccessUserPrintConfig)
}

// StorageGetScopedAccessUserCredentials shows the credentials of one scoped
// access user of a storage pool.
func StorageGetScopedAccessUserCredentials(ctx context.Context, storageId string, userId string) error {
	logger.Get().Info().Msgf("Getting credentials of scoped access user %s of storage %s", userId, storageId)

	storageIdNumeric, userIdNumeric, err := getStorageScopedAccessUserIds(storageId, userId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	credentials, httpRes, err := client.StorageAPI.
		GetStorageScopedAccessUserCredentials(ctx, storageIdNumeric, userIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(credentials, &storageScopedAccessUserCredentialsPrintConfig)
}

// StorageGetStatistics shows the capacity statistics of one storage pool.
func StorageGetStatistics(ctx context.Context, storageId string) error {
	logger.Get().Info().Msgf("Getting statistics for storage %s", storageId)

	storageIdNumeric, err := getStorageId(storageId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	statistics, httpRes, err := client.StorageAPI.GetStorageStatistics(ctx, storageIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(statistics, &storageStatisticsPrintConfig)
}

// StoragesGetStatistics shows the aggregated statistics of all storage pools.
func StoragesGetStatistics(ctx context.Context) error {
	logger.Get().Info().Msgf("Getting statistics for all storages")

	client := api.GetApiClient(ctx)

	statistics, httpRes, err := client.StorageAPI.GetStoragesStatistics(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}
	if statistics == nil {
		return fmt.Errorf("no storage statistics returned")
	}

	if formatter.IsNativeFormat() {
		return formatter.PrintResult(statistics, nil)
	}

	row := storagesStatisticsRow{
		ActiveCount:       statistics.ActiveCount,
		ReadyCount:        statistics.ReadyCount,
		PendingCount:      statistics.PendingCount,
		MaintenanceCount:  statistics.MaintenanceCount,
		ExperimentalCount: statistics.ExperimentalCount,
		LowSpaceCount:     statistics.LowSpaceCount,
		UsedSpace:         statistics.UsedSpace,
		FreeSpace:         statistics.FreeSpace,
		Types:             formatStorageTypeCounts(statistics.Types),
	}

	return formatter.PrintResult(row, &storagesStatisticsPrintConfig)
}

// formatStorageTypeCounts renders the per-type counter map as a stable,
// comma-separated "type: count" list.
func formatStorageTypeCounts(types map[string]interface{}) string {
	if len(types) == 0 {
		return ""
	}

	keys := make([]string, 0, len(types))
	for key := range types {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	entries := make([]string, 0, len(keys))
	for _, key := range keys {
		entries = append(entries, fmt.Sprintf("%s: %v", key, types[key]))
	}

	return strings.Join(entries, ", ")
}

func getStorageInterfaceIds(storageId string, interfaceId string) (int64, int64, error) {
	storageIdNumeric, err := getStorageId(storageId)
	if err != nil {
		return 0, 0, err
	}

	interfaceIdNumeric, err := strconv.ParseInt(interfaceId, 10, 64)
	if err != nil {
		err := fmt.Errorf("invalid storage interface ID: '%s'", interfaceId)
		logger.Get().Error().Err(err).Msg("")
		return 0, 0, err
	}

	return storageIdNumeric, interfaceIdNumeric, nil
}

func getStorageScopedAccessUserIds(storageId string, userId string) (int64, int64, error) {
	storageIdNumeric, err := getStorageId(storageId)
	if err != nil {
		return 0, 0, err
	}

	userIdNumeric, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		err := fmt.Errorf("invalid storage scoped access user ID: '%s'", userId)
		logger.Get().Error().Err(err).Msg("")
		return 0, 0, err
	}

	return storageIdNumeric, userIdNumeric, nil
}
