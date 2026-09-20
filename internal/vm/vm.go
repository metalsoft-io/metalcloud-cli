package vm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

// VMListPrintConfig drives the table columns for VM list output.
var VMListPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"Name": {
			Title:    "Name",
			MaxWidth: 30,
			Order:    2,
		},
		"InfrastructureId": {
			Title: "Infra",
			Order: 3,
		},
		"Host": {
			Title:    "Host",
			MaxWidth: 40,
			Order:    4,
		},
		"CpuCores": {
			Title: "CPU Cores",
			Order: 5,
		},
		"RamGB": {
			Title: "RAM (GB)",
			Order: 6,
		},
		"DiskSizeGB": {
			Title: "Disk (GB)",
			Order: 7,
		},
		"TypeId": {
			Title: "Type",
			Order: 8,
		},
		"PoolId": {
			Title: "Pool",
			Order: 9,
		},
		"AdministrationState": {
			Title: "Admin State",
			Order: 10,
		},
		"PowerState": {
			Title:       "Power",
			Transformer: formatter.FormatStatusValue,
			Order:       11,
		},
	},
}

// vmSummary is the table projection of a VM.
//
// The generated sdk.VM model cannot decode live responses: the API returns
// `hosts` as a plain string while the SDK declares it as []string, so its
// strict UnmarshalJSON rejects every payload with "json: cannot unmarshal
// string into Go struct field _VM.hosts of type []string". Every endpoint that
// answers with VMs therefore goes through the raw-body helpers; json and yaml
// output stays lossless, the table uses these fields.
type vmSummary struct {
	Id                  int64   `json:"id"`
	Name                string  `json:"name"`
	InfrastructureId    int64   `json:"infrastructureId"`
	Host                string  `json:"host"`
	CpuCores            float32 `json:"cpuCores"`
	RamGB               float32 `json:"ramGB"`
	DiskSizeGB          float32 `json:"diskSizeGB"`
	TypeId              int64   `json:"typeId"`
	PoolId              int64   `json:"poolId"`
	AdministrationState string  `json:"administrationState"`
	PowerState          string  `json:"powerState"`
}

// PrintVMsRaw renders a page of raw VM objects. It is shared with the vm-pool
// command group, which lists the VMs of a pool and of a cluster host.
func PrintVMsRaw(rawItems []json.RawMessage, meta sdk.PaginatedResponseMeta) error {
	records, err := utils.UnmarshalRawItems[vmSummary](rawItems)
	if err != nil {
		return fmt.Errorf("failed to parse VMs: %w", err)
	}

	return utils.PrintAllRaw(rawItems, records, meta, len(records), &VMListPrintConfig)
}

// printVMRaw renders one raw VM object: the full payload for json/yaml, the
// table projection otherwise. formatter.PrintResult dumps every key of a map
// when it renders a table, so the table path needs the struct.
func printVMRaw(body []byte) error {
	if len(body) == 0 {
		logger.Get().Info().Msg("No VM returned")
		return nil
	}

	if formatter.IsNativeFormat() {
		return utils.PrintRawObject(body, &VMListPrintConfig)
	}

	var record vmSummary
	if err := json.Unmarshal(body, &record); err != nil {
		return fmt.Errorf("failed to parse VM: %w", err)
	}

	return formatter.PrintResult(record, &VMListPrintConfig)
}

// VMFilters holds the optional list filters accepted by `vm list`. Every entry
// is passed through utils.ProcessFilterStringSlice, so bare values are turned
// into the SDK filter DSL ($eq:value).
type VMFilters struct {
	Id                  []string
	SiteId              []string
	Name                []string
	Address             []string
	Host                []string
	Hosts               []string
	TypeId              []string
	PoolId              []string
	AdministrationState []string
	NumaNodes           []string
	InfrastructureId    []string
}

func (f VMFilters) queryValues() url.Values {
	values := url.Values{}

	add := func(key string, filter []string) {
		for _, value := range utils.ProcessFilterStringSlice(filter) {
			values.Add(key, value)
		}
	}

	add("filter.id", f.Id)
	add("filter.siteId", f.SiteId)
	add("filter.name", f.Name)
	add("filter.address", f.Address)
	add("filter.host", f.Host)
	add("filter.hosts", f.Hosts)
	add("filter.typeId", f.TypeId)
	add("filter.poolId", f.PoolId)
	add("filter.administrationState", f.AdministrationState)
	add("filter.numaNodes", f.NumaNodes)
	add("filter.infrastructureId", f.InfrastructureId)

	return values
}

// VMList lists every VM visible to the caller, across all pools and
// infrastructures, applying the optional filters.
func VMList(ctx context.Context, filters VMFilters) error {
	logger.Get().Info().Msgf("Listing VMs")

	values := filters.queryValues()

	rawItems, meta, err := utils.FetchAllPagesRaw(func(page float32) (*http.Response, error) {
		return api.DoJSONRequest(ctx, http.MethodGet, "/api/v2/vms?"+pageQuery(values, page), nil)
	})
	if err != nil {
		return err
	}

	return PrintVMsRaw(rawItems, meta)
}

// FetchAllVMsRaw walks every page of a VM listing endpoint and returns the raw
// items. path must be an API path without query string; the pagination and
// sorting parameters are appended by this helper. It is shared with the vm-pool
// command group.
func FetchAllVMsRaw(ctx context.Context, path string) ([]json.RawMessage, sdk.PaginatedResponseMeta, error) {
	return utils.FetchAllPagesRaw(func(page float32) (*http.Response, error) {
		return api.DoJSONRequest(ctx, http.MethodGet, path+"?"+pageQuery(url.Values{}, page), nil)
	})
}

// FetchVMPageRaw fetches one explicit page of a VM listing endpoint.
func FetchVMPageRaw(ctx context.Context, path string, page, limit float32) ([]json.RawMessage, sdk.PaginatedResponseMeta, error) {
	return utils.FetchPageWindowRaw(func(p, l float32) (*http.Response, error) {
		values := url.Values{}
		values.Set("page", fmt.Sprintf("%.0f", p))
		values.Set("limit", fmt.Sprintf("%.0f", l))
		values.Set("sortBy", "id:ASC")
		return api.DoJSONRequest(ctx, http.MethodGet, path+"?"+values.Encode(), nil)
	}, int(page), int(limit))
}

// pageQuery clones values and adds the standard pagination parameters.
func pageQuery(values url.Values, page float32) string {
	pageValues := url.Values{}
	for key, entries := range values {
		pageValues[key] = entries
	}
	pageValues.Set("page", fmt.Sprintf("%.0f", page))
	pageValues.Set("limit", "100")
	pageValues.Set("sortBy", "id:ASC")

	return pageValues.Encode()
}

func GetVMId(vmId string) (int64, error) {
	vmIdNumeric, err := strconv.ParseInt(vmId, 10, 64)
	if err != nil {
		err := fmt.Errorf("invalid VM ID: '%s'", vmId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return vmIdNumeric, nil
}

func VMGet(ctx context.Context, vmId string) error {
	logger.Get().Info().Msgf("Getting VM '%s'", vmId)

	vmIdNumeric, err := GetVMId(vmId)
	if err != nil {
		return err
	}

	body, err := api.RawJSONRequest(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/vms/%d", vmIdNumeric), nil, nil)
	if err != nil {
		return err
	}

	return printVMRaw(body)
}

func VMPowerStatus(ctx context.Context, vmId string) error {
	logger.Get().Info().Msgf("Getting power status for VM '%s'", vmId)

	vmIdNumeric, err := GetVMId(vmId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	powerStatus, httpRes, err := client.VMAPI.GetVMPowerStatus(ctx, vmIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Power status for VM '%s' is '%s'", vmId, powerStatus)

	return formatter.PrintResult(powerStatus, nil)
}

func VMStart(ctx context.Context, vmId string) error {
	logger.Get().Info().Msgf("Starting VM '%s'", vmId)

	vmIdNumeric, err := GetVMId(vmId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.VMAPI.StartVM(ctx, vmIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("VM '%s' started", vmId)

	return nil
}

func VMShutdown(ctx context.Context, vmId string) error {
	logger.Get().Info().Msgf("Shutting down VM '%s'", vmId)

	vmIdNumeric, err := GetVMId(vmId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.VMAPI.ShutdownVM(ctx, vmIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("VM '%s' shutdown initiated", vmId)

	return nil
}

func VMReboot(ctx context.Context, vmId string) error {
	logger.Get().Info().Msgf("Rebooting VM '%s'", vmId)

	vmIdNumeric, err := GetVMId(vmId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.VMAPI.RebootVM(ctx, vmIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("VM '%s' reboot initiated", vmId)

	return nil
}

func VMUpdate(ctx context.Context, vmId string, config []byte) error {
	logger.Get().Info().Msgf("Updating VM '%s'", vmId)

	var updateConfig sdk.UpdateVM
	err := utils.UnmarshalContent(config, &updateConfig)
	if err != nil {
		return err
	}

	vmIdNumeric, err := GetVMId(vmId)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(updateConfig)
	if err != nil {
		return fmt.Errorf("failed to encode VM update: %w", err)
	}

	body, err := api.RawJSONRequest(ctx, http.MethodPatch,
		fmt.Sprintf("/api/v2/vms/%d", vmIdNumeric), payload, nil)
	if err != nil {
		return err
	}

	logger.Get().Info().Msgf("VM '%s' updated", vmId)

	return printVMRaw(body)
}

func VMRemoteConsoleInfo(ctx context.Context, vmId string) error {
	logger.Get().Info().Msgf("Getting remote console info for VM '%s'", vmId)

	vmIdNumeric, err := GetVMId(vmId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	consoleInfo, httpRes, err := client.VMAPI.GetVMRemoteConsoleInfo(ctx, vmIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(consoleInfo, &formatter.PrintConfig{
		FieldsConfig: map[string]formatter.RecordFieldConfig{
			"ActiveConnections": {
				Title: "Active Connections",
				Order: 1,
			},
			"ConsoleUrl": {
				Title: "Console URL",
				Order: 2,
			},
		},
	})
}
