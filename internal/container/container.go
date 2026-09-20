package container

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var containerPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"Name": {
			MaxWidth: 30,
			Order:    2,
		},
		"SiteId": {
			Title: "Site",
			Order: 3,
		},
		"InfrastructureId": {
			Title: "Infra",
			Order: 4,
		},
		"ContainerInstanceId": {
			Title: "Instance",
			Order: 5,
		},
		"TypeId": {
			Title: "Type",
			Order: 6,
		},
		"PoolId": {
			Title: "Pool",
			Order: 7,
		},
		"Host": {
			MaxWidth: 30,
			Order:    8,
		},
		"CpuCores": {
			Title: "Cores",
			Order: 9,
		},
		"RamGB": {
			Title: "RAM GB",
			Order: 10,
		},
		"DiskSizeGB": {
			Title: "Disk GB",
			Order: 11,
		},
		"AdministrationState": {
			Title:       "Admin State",
			Transformer: formatter.FormatStatusValue,
			Order:       12,
		},
		"PowerState": {
			Title:       "Power",
			Transformer: formatter.FormatStatusValue,
			Order:       13,
		},
		"CreatedTimestamp": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       14,
		},
	},
}

var containerRemoteConsoleInfoPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"ActiveConnections": {
			Title: "Active Connections",
			Order: 1,
		},
	},
}

// containerSummary is the table projection of a container.
//
// The generated sdk.Container model cannot be used to decode live responses:
// the API returns `hosts` as a plain string while the SDK declares it as
// []string, so its strict UnmarshalJSON rejects every payload with
// "json: cannot unmarshal string into Go struct field _Container.hosts of type
// []string". Every endpoint that answers with containers therefore goes through
// the raw-body helpers; json/yaml output stays lossless, the table uses these
// fields.
type containerSummary struct {
	Id                  int64   `json:"id"`
	Name                string  `json:"name"`
	SiteId              int64   `json:"siteId"`
	InfrastructureId    int64   `json:"infrastructureId"`
	ContainerInstanceId int64   `json:"containerInstanceId"`
	TypeId              int64   `json:"typeId"`
	PoolId              int64   `json:"poolId"`
	Host                string  `json:"host"`
	CpuCores            float32 `json:"cpuCores"`
	RamGB               float32 `json:"ramGB"`
	DiskSizeGB          float32 `json:"diskSizeGB"`
	AdministrationState string  `json:"administrationState"`
	PowerState          string  `json:"powerState"`
	CreatedTimestamp    string  `json:"createdTimestamp"`
}

// PrintConfig exposes the container table layout so other packages that list
// containers (e.g. the container-type command group) render them identically.
func PrintConfig() *formatter.PrintConfig {
	return &containerPrintConfig
}

// ContainerFilters holds the optional list filters accepted by `container list`.
type ContainerFilters struct {
	Id                  []string
	SiteId              []string
	Name                []string
	Host                []string
	TypeId              []string
	PoolId              []string
	AdministrationState []string
	InfrastructureId    []string
}

func (f ContainerFilters) queryValues() url.Values {
	values := url.Values{}

	add := func(key string, filter []string) {
		for _, value := range utils.ProcessFilterStringSlice(filter) {
			values.Add(key, value)
		}
	}

	add("filter.id", f.Id)
	add("filter.siteId", f.SiteId)
	add("filter.name", f.Name)
	add("filter.host", f.Host)
	add("filter.typeId", f.TypeId)
	add("filter.poolId", f.PoolId)
	add("filter.administrationState", f.AdministrationState)
	add("filter.infrastructureId", f.InfrastructureId)

	return values
}

func ContainerList(ctx context.Context, filters ContainerFilters) error {
	logger.Get().Info().Msgf("Listing containers")

	values := filters.queryValues()

	rawItems, meta, err := utils.FetchAllPagesRaw(func(page float32) (*http.Response, error) {
		pageValues := url.Values{}
		for key, entries := range values {
			pageValues[key] = entries
		}
		pageValues.Set("page", fmt.Sprintf("%.0f", page))
		pageValues.Set("limit", "100")
		pageValues.Set("sortBy", "id:ASC")

		return api.DoJSONRequest(ctx, http.MethodGet, "/api/v2/containers?"+pageValues.Encode(), nil)
	})
	if err != nil {
		return err
	}

	return PrintContainersRaw(rawItems, meta)
}

// PrintContainersRaw renders a page of raw container objects. It is shared with
// the container-type command group, which lists the containers of a type.
func PrintContainersRaw(rawItems []json.RawMessage, meta sdk.PaginatedResponseMeta) error {
	records, err := utils.UnmarshalRawItems[containerSummary](rawItems)
	if err != nil {
		return fmt.Errorf("failed to parse containers: %w", err)
	}

	return utils.PrintAllRaw(rawItems, records, meta, len(records), &containerPrintConfig)
}

func ContainerGet(ctx context.Context, containerId string) error {
	logger.Get().Info().Msgf("Get container '%s'", containerId)

	containerIdNumerical, err := GetContainerId(containerId)
	if err != nil {
		return err
	}

	body, err := api.RawJSONRequest(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/containers/%d", containerIdNumerical), nil, nil)
	if err != nil {
		return err
	}

	return printContainerRaw(body)
}

func ContainerUpdate(ctx context.Context, containerId string, config []byte) error {
	logger.Get().Info().Msgf("Update container '%s'", containerId)

	containerIdNumerical, err := GetContainerId(containerId)
	if err != nil {
		return err
	}

	var update sdk.UpdateContainer
	if err := utils.UnmarshalContent(config, &update); err != nil {
		return err
	}

	payload, err := json.Marshal(update)
	if err != nil {
		return fmt.Errorf("failed to encode container update: %w", err)
	}

	body, err := api.RawJSONRequest(ctx, http.MethodPatch,
		fmt.Sprintf("/api/v2/containers/%d", containerIdNumerical), payload, nil)
	if err != nil {
		return err
	}

	return printContainerRaw(body)
}

// printContainerRaw renders one raw container object: the full payload for
// json/yaml, the table projection otherwise. formatter.PrintResult dumps every
// key of a map when it renders a table, so the table path needs the struct.
func printContainerRaw(body []byte) error {
	if formatter.IsNativeFormat() {
		return utils.PrintRawObject(body, &containerPrintConfig)
	}

	var record containerSummary
	if err := json.Unmarshal(body, &record); err != nil {
		return fmt.Errorf("failed to parse container: %w", err)
	}

	return formatter.PrintResult(record, &containerPrintConfig)
}

func ContainerPowerStatus(ctx context.Context, containerId string) error {
	logger.Get().Info().Msgf("Get power status of container '%s'", containerId)

	containerIdNumerical, err := GetContainerId(containerId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	powerStatus, httpRes, err := client.ContainerAPI.GetContainerPowerStatus(ctx, containerIdNumerical).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(powerStatus, nil)
}

func ContainerRemoteConsoleInfo(ctx context.Context, containerId string) error {
	logger.Get().Info().Msgf("Get remote console info of container '%s'", containerId)

	containerIdNumerical, err := GetContainerId(containerId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	consoleInfo, httpRes, err := client.ContainerAPI.GetContainerRemoteConsoleInfo(ctx, containerIdNumerical).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(consoleInfo, &containerRemoteConsoleInfoPrintConfig)
}

// ContainerPowerControl runs one of the 'start', 'shutdown' or 'reboot' actions
// against a container. These endpoints answer 204 with no body.
func ContainerPowerControl(ctx context.Context, containerId string, action string) error {
	logger.Get().Info().Msgf("Performing '%s' action on container '%s'", action, containerId)

	containerIdNumerical, err := GetContainerId(containerId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	var httpRes *http.Response

	switch action {
	case "start":
		httpRes, err = client.ContainerAPI.StartContainer(ctx, containerIdNumerical).Execute()
	case "shutdown":
		httpRes, err = client.ContainerAPI.ShutdownContainer(ctx, containerIdNumerical).Execute()
	case "reboot":
		httpRes, err = client.ContainerAPI.RebootContainer(ctx, containerIdNumerical).Execute()
	default:
		err := fmt.Errorf("unsupported power action: '%s'. Use 'start', 'shutdown' or 'reboot'", action)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Container '%s' power action '%s' successful", containerId, action)
	return nil
}

func GetContainerId(containerId string) (int64, error) {
	containerIdNumerical, err := utils.GetInt64FromString(containerId)
	if err != nil {
		err := fmt.Errorf("invalid container ID: '%s'", containerId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return containerIdNumerical, nil
}
