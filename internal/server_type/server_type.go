package server_type

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var serverTypePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"Name": {
			MaxWidth: 30,
			Order:    2,
		},
		"Label": {
			MaxWidth: 30,
			Order:    3,
		},
		"ProcessorCount": {
			Title: "CPU #",
			Order: 4,
		},
		"ProcessorCoreMhz": {
			Title: "CPU MHz",
			Order: 5,
		},
		"ProcessorNames": {
			Title: "CPU Names",
			Order: 6,
		},
		"RamGbytes": {
			Title: "RAM GB",
			Order: 7,
		},
		"NetworkInterfaceCount": {
			Title: "NIC #",
			Order: 8,
		},
		"DiskCount": {
			Title: "Disk #",
			Order: 9,
		},
		"GpuCount": {
			Title: "GPU #",
			Order: 10,
		},
	},
}

func ServerTypeList(ctx context.Context) error {
	logger.Get().Info().Msgf("Listing all server types")

	client := api.GetApiClient(ctx)

	request := client.ServerTypeAPI.GetServerTypes(ctx).SortBy([]string{"id:ASC"})

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &serverTypePrintConfig)
}

func ServerTypeGet(ctx context.Context, serverTypeIdOrLabel string) error {
	logger.Get().Info().Msgf("Get server type %s info", serverTypeIdOrLabel)

	serverType, err := GetServerTypeByIdOrLabel(ctx, serverTypeIdOrLabel)
	if err != nil {
		return err
	}

	return formatter.PrintResult(serverType, &serverTypePrintConfig)
}

func GetServerTypeByIdOrLabel(ctx context.Context, serverTypeIdOrLabel string) (*sdk.ServerType, error) {
	client := api.GetApiClient(ctx)

	serverTypeId, err := utils.GetInt64FromString(serverTypeIdOrLabel)
	if err != nil {
		return nil, err
	}

	serverTypeInfo, httpRes, err := client.ServerTypeAPI.GetServerTypeInfo(ctx, serverTypeId).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, err
	}

	return serverTypeInfo, nil
}

// formatCountMap renders a "key: count" map (as returned by the statistics
// endpoints) as a single sorted, comma separated cell. The formatter flattens
// maps into a []string of "key: value" entries before the transformer runs, so
// both shapes are handled.
func formatCountMap(value interface{}) string {
	var entries []string

	switch counts := value.(type) {
	case []string:
		entries = append(entries, counts...)
	case map[string]interface{}:
		for key, count := range counts {
			entries = append(entries, fmt.Sprintf("%s: %v", key, count))
		}
	default:
		return fmt.Sprintf("%v", value)
	}

	sort.Strings(entries)

	return strings.Join(entries, ", ")
}

// serverTypeStatisticsPrintConfig only renders the server counts in the table
// formats; the nested per-server information and the utilization report are
// too deep for a table and are available through the json/yaml output.
var serverTypeStatisticsPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"ServerTypeIdToServerCount": {
			Title:       "Servers By Type",
			Transformer: formatCountMap,
			MaxWidth:    60,
			Order:       1,
		},
		"ServerTypeIdToServerInformation": {
			Hidden: true,
		},
		"UtilizationReport": {
			Hidden: true,
		},
	},
}

func ServerTypeCreate(ctx context.Context, config []byte) error {
	logger.Get().Info().Msgf("Creating server type")

	var createConfig sdk.CreateServerType
	if err := utils.UnmarshalContent(config, &createConfig); err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	serverType, httpRes, err := client.ServerTypeAPI.
		CreateServerType(ctx).
		CreateServerType(createConfig).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(serverType, &serverTypePrintConfig)
}

func ServerTypeUpdate(ctx context.Context, serverTypeId string, config []byte) error {
	logger.Get().Info().Msgf("Updating server type '%s'", serverTypeId)

	var updateConfig sdk.UpdateServerType
	if err := utils.UnmarshalContent(config, &updateConfig); err != nil {
		return err
	}

	serverTypeIdNumeric, err := utils.GetInt64FromString(serverTypeId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	serverType, httpRes, err := client.ServerTypeAPI.
		UpdateServerType(ctx, serverTypeIdNumeric).
		UpdateServerType(updateConfig).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(serverType, &serverTypePrintConfig)
}

func ServerTypeDelete(ctx context.Context, serverTypeId string) error {
	logger.Get().Info().Msgf("Deleting server type '%s'", serverTypeId)

	serverTypeIdNumeric, err := utils.GetInt64FromString(serverTypeId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.ServerTypeAPI.DeleteServerType(ctx, serverTypeIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Server type '%s' deleted", serverTypeId)

	return nil
}

// ServerTypeCleanUnused removes the server types no server is using anymore.
func ServerTypeCleanUnused(ctx context.Context) error {
	logger.Get().Info().Msgf("Removing unused server types")

	client := api.GetApiClient(ctx)

	httpRes, err := client.ServerTypeAPI.RemoveUnusedServerTypes(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Unused server types removed")

	return nil
}

// ServerTypeStatistics returns the server availability statistics of a site,
// optionally restricted to the given server types.
func ServerTypeStatistics(ctx context.Context, siteId int64, serverTypeIds []string, userId int, maximumResultsPerServerType int, instanceArrayId int) error {
	logger.Get().Info().Msgf("Getting server type statistics for site %d", siteId)

	serverTypeIdsNumeric, err := utils.GetInt64SliceFromStrings(serverTypeIds)
	if err != nil {
		return err
	}

	// The body is always sent: skipping the optional setter would serialize a
	// literal null, which the API rejects with 400.
	options := sdk.ServerTypeStatisticsBatchOptions{
		SiteId: siteId,
	}

	if len(serverTypeIdsNumeric) > 0 {
		options.ServerTypeIds = serverTypeIdsNumeric
	}

	if userId > 0 {
		options.UserIdOwner = sdk.PtrFloat32(float32(userId))
	}

	if maximumResultsPerServerType > 0 {
		options.MaximumResultsPerServerType = sdk.PtrFloat32(float32(maximumResultsPerServerType))
	}

	if instanceArrayId > 0 {
		options.InstanceArrayId = sdk.PtrInt64(int64(instanceArrayId))
	}

	client := api.GetApiClient(ctx)

	statistics, httpRes, err := client.ServerTypeAPI.
		GetServerTypesStatisticsBatch(ctx).
		ServerTypeStatisticsBatchOptions(options).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(statistics, &serverTypeStatisticsPrintConfig)
}

func ServerTypeConfigExample(ctx context.Context) error {
	createConfig := sdk.CreateServerType{
		Name:                   "M.32.128.2",
		Label:                  "m-32-128-2",
		Description:            sdk.PtrString("Example server type"),
		RamGbytes:              128,
		ProcessorCount:         2,
		ProcessorCoreMhz:       2400,
		ProcessorCoreCount:     16,
		ProcessorNames:         []string{"Intel Xeon Gold 6226R"},
		NetworkInterfaceSpeeds: []float32{10000, 10000},
		DiskCount:              2,
		ServerClass:            "bigdata",
		BootType:               sdk.PtrString("uefi"),
		IsExperimental:         sdk.PtrFloat32(0),
		GpuCount:               sdk.PtrFloat32(0),
		Tags:                   []string{"example"},
		AllowedVendorSkuIds:    []string{},
	}

	return formatter.PrintResult(createConfig, nil)
}

func ServerTypeUpdateConfigExample(ctx context.Context) error {
	updateConfig := sdk.UpdateServerType{
		Label:               "m-32-128-2",
		Name:                sdk.PtrString("M.32.128.2"),
		Description:         sdk.PtrString("Example server type"),
		IsExperimental:      sdk.PtrFloat32(0),
		Tags:                []string{"example"},
		AllowedVendorSkuIds: []string{},
	}

	return formatter.PrintResult(updateConfig, nil)
}
