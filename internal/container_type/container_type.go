package container_type

import (
	"context"
	"fmt"
	"net/http"

	"github.com/metalsoft-io/metalcloud-cli/internal/container"
	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var containerTypePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
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
		"DisplayName": {
			Title:    "Display Name",
			MaxWidth: 30,
			Order:    4,
		},
		"CpuCores": {
			Title: "Cores",
			Order: 5,
		},
		"RamGB": {
			Title: "RAM GB",
			Order: 6,
		},
		"IsExperimental": {
			Title: "Experimental",
			Order: 7,
		},
		"ForUnmanagedContainersOnly": {
			Title: "Unmanaged Only",
			Order: 8,
		},
		"Tags": {
			Title:       "Tags",
			Transformer: formatter.FormatStringListValue,
			Order:       9,
		},
	},
}

// ContainerTypeFilters holds the optional list filters accepted by `container-type list`.
type ContainerTypeFilters struct {
	Id          []string
	Label       []string
	Name        []string
	DisplayName []string
}

func ContainerTypeList(ctx context.Context, filters ContainerTypeFilters) error {
	logger.Get().Info().Msgf("Listing container types")

	client := api.GetApiClient(ctx)

	request := client.ContainerTypeAPI.GetContainerTypes(ctx).SortBy([]string{"id:ASC"})

	if len(filters.Id) > 0 {
		request = request.FilterId(utils.ProcessFilterStringSlice(filters.Id))
	}
	if len(filters.Label) > 0 {
		request = request.FilterLabel(utils.ProcessFilterStringSlice(filters.Label))
	}
	if len(filters.Name) > 0 {
		request = request.FilterName(utils.ProcessFilterStringSlice(filters.Name))
	}
	if len(filters.DisplayName) > 0 {
		request = request.FilterDisplayName(utils.ProcessFilterStringSlice(filters.DisplayName))
	}

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &containerTypePrintConfig)
}

func ContainerTypeGet(ctx context.Context, containerTypeId string) error {
	logger.Get().Info().Msgf("Get container type '%s'", containerTypeId)

	containerTypeIdNumerical, err := GetContainerTypeId(containerTypeId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	// The SDK types this path parameter as float32.
	containerType, httpRes, err := client.ContainerTypeAPI.
		GetContainerType(ctx, float32(containerTypeIdNumerical)).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(containerType, &containerTypePrintConfig)
}

func ContainerTypeConfigExample(ctx context.Context) error {
	example := sdk.CreateContainerType{
		Name:                       "container-type-name",
		DisplayName:                sdk.PtrString("Container Type Display Name"),
		Label:                      sdk.PtrString("container-type-label"),
		CpuCores:                   4,
		RamGB:                      16,
		IsExperimental:             sdk.PtrFloat32(0),
		ForUnmanagedContainersOnly: sdk.PtrFloat32(0),
		Tags:                       []string{"tag1", "tag2"},
	}

	if formatter.IsNativeFormat() {
		return formatter.PrintResult(example, nil)
	}

	return formatter.PrintYamlResult(example)
}

func ContainerTypeCreate(ctx context.Context, config []byte) error {
	logger.Get().Info().Msgf("Creating container type")

	var create sdk.CreateContainerType
	if err := utils.UnmarshalContent(config, &create); err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	containerType, httpRes, err := client.ContainerTypeAPI.
		CreateContainerType(ctx).
		CreateContainerType(create).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(containerType, &containerTypePrintConfig)
}

func ContainerTypeUpdate(ctx context.Context, containerTypeId string, config []byte) error {
	logger.Get().Info().Msgf("Updating container type '%s'", containerTypeId)

	containerTypeIdNumerical, err := GetContainerTypeId(containerTypeId)
	if err != nil {
		return err
	}

	var update sdk.UpdateContainerType
	if err := utils.UnmarshalContent(config, &update); err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	containerType, httpRes, err := client.ContainerTypeAPI.
		UpdateContainerType(ctx, containerTypeIdNumerical).
		UpdateContainerType(update).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(containerType, &containerTypePrintConfig)
}

func ContainerTypeDelete(ctx context.Context, containerTypeId string) error {
	logger.Get().Info().Msgf("Deleting container type '%s'", containerTypeId)

	containerTypeIdNumerical, err := GetContainerTypeId(containerTypeId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.ContainerTypeAPI.
		DeleteContainerType(ctx, containerTypeIdNumerical).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Container type '%s' deleted", containerTypeId)
	return nil
}

func ContainerTypeContainers(ctx context.Context, containerTypeId string) error {
	logger.Get().Info().Msgf("Listing containers of container type '%s'", containerTypeId)

	containerTypeIdNumerical, err := GetContainerTypeId(containerTypeId)
	if err != nil {
		return err
	}

	// sdk.Container rejects live payloads (the API returns `hosts` as a string,
	// the SDK declares []string), so this list goes through the raw helpers.
	rawItems, meta, err := utils.FetchAllPagesRaw(func(page float32) (*http.Response, error) {
		return api.DoJSONRequest(ctx, http.MethodGet,
			fmt.Sprintf("/api/v2/container-types/%d/containers?page=%.0f&limit=100&sortBy=id:ASC", containerTypeIdNumerical, page), nil)
	})
	if err != nil {
		return err
	}

	return container.PrintContainersRaw(rawItems, meta)
}

func GetContainerTypeId(containerTypeId string) (int64, error) {
	containerTypeIdNumerical, err := utils.GetInt64FromString(containerTypeId)
	if err != nil {
		err := fmt.Errorf("invalid container type ID: '%s'", containerTypeId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return containerTypeIdNumerical, nil
}
