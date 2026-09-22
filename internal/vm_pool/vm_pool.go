package vm_pool

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/internal/container"
	"github.com/metalsoft-io/metalcloud-cli/internal/network_device"
	"github.com/metalsoft-io/metalcloud-cli/internal/vm"
	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var VMPoolPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"SiteId": {
			Title: "Site",
			Order: 2,
		},
		"Name": {
			Title: "Name",
			Order: 3,
		},
		"Type": {
			Title: "Type",
			Order: 4,
		},
		"ManagementHost": {
			Title: "Management Host",
			Order: 5,
		},
		"ManagementPort": {
			Title: "Management Port",
			Order: 6,
		},
		"Status": {
			Title:       "Status",
			Transformer: formatter.FormatStatusValue,
			Order:       7,
		},
	},
}

var vmPoolHostPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"PoolId": {
			Title: "Pool",
			Order: 2,
		},
		"Name": {
			Title:    "Name",
			MaxWidth: 30,
			Order:    3,
		},
		"Address": {
			Title:    "Address",
			MaxWidth: 30,
			Order:    4,
		},
		"FailureDomain": {
			Title: "Failure Domain",
			Order: 5,
		},
		"Architecture": {
			Title: "Arch",
			Order: 6,
		},
		"Status": {
			Title:       "Status",
			Transformer: formatter.FormatStatusValue,
			Order:       7,
		},
		"HealthStatus": {
			Title:       "Health",
			Transformer: formatter.FormatStatusValue,
			Order:       8,
		},
		"AllowVMsToBeCreated": {
			Title: "Allow VMs",
			Order: 9,
		},
		"AllowContainersToBeCreated": {
			Title: "Allow Containers",
			Order: 10,
		},
		"UpdatedTimestamp": {
			Title:       "Updated",
			Transformer: formatter.FormatDateTimeValue,
			Order:       11,
		},
	},
}

var vmPoolHostInterfacePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"HostId": {
			Title: "Host",
			Order: 2,
		},
		"Name": {
			Title:    "Name",
			MaxWidth: 30,
			Order:    3,
		},
		"MacAddress": {
			Title: "MAC Address",
			Order: 4,
		},
		"Status": {
			Title:       "Status",
			Transformer: formatter.FormatStatusValue,
			Order:       5,
		},
	},
}

var vmPoolHostInterfaceNetworkDevicePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"HostInterfaceId": {
			Title: "Host Interface",
			Order: 2,
		},
		"NetworkDeviceId": {
			Title: "Network Device",
			Order: 3,
		},
		"NetworkDeviceInterfaceName": {
			Title:    "Device Interface",
			MaxWidth: 40,
			Order:    4,
		},
	},
}

var vmPoolStatisticsPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"TotalRamGB": {
			Title: "Total RAM GB",
			Order: 1,
		},
		"UsedRamGB": {
			Title: "Used RAM GB",
			Order: 2,
		},
		"FreeRamGB": {
			Title: "Free RAM GB",
			Order: 3,
		},
		"TotalSpaceGB": {
			Title: "Total Disk GB",
			Order: 4,
		},
		"UsedSpaceGB": {
			Title: "Used Disk GB",
			Order: 5,
		},
		"FreeSpaceGB": {
			Title: "Free Disk GB",
			Order: 6,
		},
		"GpuInfo": {
			Title:    "GPUs",
			MaxWidth: 40,
			Order:    7,
		},
	},
}

// vmPoolStatisticsRow is the table projection of sdk.VMPoolStatistics: the
// nested gpuInfo objects are joined into one cell because the tabular
// formatter renders struct-valued fields as empty cells. The json and yaml
// formats keep the full statistics object.
type vmPoolStatisticsRow struct {
	TotalRamGB   float32
	UsedRamGB    float32
	FreeRamGB    float32
	TotalSpaceGB float32
	UsedSpaceGB  float32
	FreeSpaceGB  float32
	GpuInfo      string
}

// printVMPoolStatistics renders the statistics object: the full payload for
// json/yaml, the flattened projection for the tabular formats.
func printVMPoolStatistics(statistics *sdk.VMPoolStatistics) error {
	if statistics == nil {
		logger.Get().Info().Msg("No statistics returned")
		return nil
	}

	if formatter.IsNativeFormat() {
		return formatter.PrintResult(statistics, nil)
	}

	gpus := make([]string, 0, len(statistics.GpuInfo))
	for _, gpu := range statistics.GpuInfo {
		gpus = append(gpus, fmt.Sprintf("%s x%g", gpu.Name, gpu.Count))
	}

	return formatter.PrintResult(vmPoolStatisticsRow{
		TotalRamGB:   statistics.TotalRamGB,
		UsedRamGB:    statistics.UsedRamGB,
		FreeRamGB:    statistics.FreeRamGB,
		TotalSpaceGB: statistics.TotalSpaceGB,
		UsedSpaceGB:  statistics.UsedSpaceGB,
		FreeSpaceGB:  statistics.FreeSpaceGB,
		GpuInfo:      strings.Join(gpus, ", "),
	}, &vmPoolStatisticsPrintConfig)
}

func VMPoolList(ctx context.Context, filterType []string) error {
	logger.Get().Info().Msgf("Listing all VM pools")

	client := api.GetApiClient(ctx)

	request := client.VMPoolAPI.GetVMPools(ctx).SortBy([]string{"id:ASC"})

	if len(filterType) > 0 {
		request = request.FilterType(utils.ProcessFilterStringSlice(filterType))
	}

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &VMPoolPrintConfig)
}

func VMPoolGet(ctx context.Context, vmPoolId string) error {
	logger.Get().Info().Msgf("Get VM pool %s details", vmPoolId)

	vmPoolIdNumeric, err := getVMPoolId(vmPoolId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	vmPool, httpRes, err := client.VMPoolAPI.GetVMPool(ctx, float32(vmPoolIdNumeric)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(vmPool, &VMPoolPrintConfig)
}

func VMPoolCreate(ctx context.Context, config []byte) error {
	logger.Get().Info().Msgf("Creating new VM pool")

	client := api.GetApiClient(ctx)

	var createVMPool sdk.CreateVMPool
	err := utils.UnmarshalContent(config, &createVMPool)
	if err != nil {
		return err
	}

	response, httpRes, err := client.VMPoolAPI.CreateVMPool(ctx).
		CreateVMPool(createVMPool).
		Execute()

	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("VM pool created with ID: %d", int(response.Id))
	return nil
}

func VMPoolDelete(ctx context.Context, vmPoolId string) error {
	logger.Get().Info().Msgf("Deleting VM pool %s", vmPoolId)

	vmPoolIdNumeric, err := getVMPoolId(vmPoolId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.VMPoolAPI.DeleteVMPool(ctx, vmPoolIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("VM pool %s deleted successfully", vmPoolId)
	return nil
}

func VMPoolGetCredentials(ctx context.Context, vmPoolId string) error {
	logger.Get().Info().Msgf("Getting credentials for VM pool %s", vmPoolId)

	vmPoolIdNumeric, err := getVMPoolId(vmPoolId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	credentials, httpRes, err := client.VMPoolAPI.GetVMPoolCredentials(ctx, vmPoolIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(credentials, nil)
}

func VMPoolGetVMs(ctx context.Context, vmPoolId string, limit float32, page float32) error {
	logger.Get().Info().Msgf("Getting VMs for VM pool %s", vmPoolId)

	vmPoolIdNumeric, err := getVMPoolId(vmPoolId)
	if err != nil {
		return err
	}

	// The strict sdk.VM model rejects live payloads ("hosts" comes back as a
	// plain string), so VM listings go through the raw-body helpers.
	path := fmt.Sprintf("/api/v2/vm-pools/%d/vms", vmPoolIdNumeric)

	if limit > 0 || page > 0 {
		rawItems, meta, err := vm.FetchVMPageRaw(ctx, path, page, limit)
		if err != nil {
			return err
		}

		return vm.PrintVMsRaw(rawItems, meta)
	}

	rawItems, meta, err := vm.FetchAllVMsRaw(ctx, path)
	if err != nil {
		return err
	}

	return vm.PrintVMsRaw(rawItems, meta)
}

func VMPoolGetClusterHosts(ctx context.Context, vmPoolId string, limit float32, page float32) error {
	logger.Get().Info().Msgf("Getting cluster hosts for VM pool %s", vmPoolId)

	vmPoolIdNumeric, err := getVMPoolId(vmPoolId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.VMPoolAPI.GetVMPoolClusterHosts(ctx, vmPoolIdNumeric)

	if limit > 0 || page > 0 {
		if limit > 0 {
			request = request.Limit(limit)
		}
		if page > 0 {
			request = request.Page(page)
		}

		result, httpRes, err := request.Execute()
		if err := response_inspector.InspectResponse(httpRes, err); err != nil {
			return err
		}

		return utils.PrintAll(result.Data, result.Meta, len(result.Data), &vmPoolHostPrintConfig)
	}

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &vmPoolHostPrintConfig)
}

func VMPoolGetClusterHostVMs(ctx context.Context, vmPoolId string, hostId string, limit float32, page float32) error {
	logger.Get().Info().Msgf("Getting VMs for cluster host %s in VM pool %s", hostId, vmPoolId)

	vmPoolIdNumeric, err := getVMPoolId(vmPoolId)
	if err != nil {
		return err
	}

	hostIdNumeric, err := parseId(hostId, "cluster host ID")
	if err != nil {
		return err
	}

	// See VMPoolGetVMs: the strict sdk.VM model cannot decode live payloads.
	path := fmt.Sprintf("/api/v2/vm-pools/%d/cluster-hosts/%d/vms", vmPoolIdNumeric, hostIdNumeric)

	if limit > 0 || page > 0 {
		rawItems, meta, err := vm.FetchVMPageRaw(ctx, path, page, limit)
		if err != nil {
			return err
		}

		return vm.PrintVMsRaw(rawItems, meta)
	}

	rawItems, meta, err := vm.FetchAllVMsRaw(ctx, path)
	if err != nil {
		return err
	}

	return vm.PrintVMsRaw(rawItems, meta)
}

func VMPoolGetClusterHostInterfaces(ctx context.Context, vmPoolId string, hostId string) error {
	logger.Get().Info().Msgf("Getting interfaces for cluster host %s in VM pool %s", hostId, vmPoolId)

	vmPoolIdNumeric, hostIdNumeric, err := getVMPoolAndHostId(vmPoolId, hostId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	// GetVMPoolClusterHostInterfaces returns a flat []VMPoolHostInterfaces — no Page/Limit methods.
	interfaces, httpRes, err := client.VMPoolAPI.GetVMPoolClusterHostInterfaces(ctx, vmPoolIdNumeric, hostIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(interfaces, &vmPoolHostInterfacePrintConfig)
}

func VMPoolConfigExample(ctx context.Context) error {
	vmPoolConfig := sdk.CreateVMPool{
		SiteId:         1,
		ManagementHost: "vcenter.example.com",
		ManagementPort: 443,
		Name:           "VM-Pool-Example",
		Description:    sdk.PtrString("Example VM pool for testing"),
		Type:           "vmware",
		Certificate:    sdk.PtrString("-----BEGIN CERTIFICATE-----\nMIID...certificate content...==\n-----END CERTIFICATE-----"),
		PrivateKey:     sdk.PtrString("-----BEGIN PRIVATE KEY-----\nMIIE...key content...==\n-----END PRIVATE KEY-----"),
		InMaintenance:  sdk.PtrFloat32(0),
		IsExperimental: sdk.PtrFloat32(0),
		Tags:           []string{"test", "example", "vmware"},
	}

	return formatter.PrintResult(vmPoolConfig, nil)
}

func VMPoolImportVMs(ctx context.Context, vmPoolId string, importVMs sdk.VMPoolImportVMs) error {
	logger.Get().Info().Msgf("Importing VMs into VM pool %s", vmPoolId)

	vmPoolIdNumeric, err := getVMPoolId(vmPoolId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.VMPoolAPI.ImportVMPoolVMs(ctx, vmPoolIdNumeric).
		VMPoolImportVMs(importVMs).
		Execute()

	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("VMs imported successfully into VM pool %s", vmPoolId)
	return nil
}

// VMPoolUpdate patches a VM pool with the given configuration. The endpoint
// takes no If-Match header, so no revision is fetched first.
func VMPoolUpdate(ctx context.Context, vmPoolId string, config []byte) error {
	logger.Get().Info().Msgf("Updating VM pool %s", vmPoolId)

	vmPoolIdNumeric, err := getVMPoolId(vmPoolId)
	if err != nil {
		return err
	}

	var updateVMPool sdk.UpdateVMPool
	if err := utils.UnmarshalContent(config, &updateVMPool); err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	vmPool, httpRes, err := client.VMPoolAPI.UpdateVMPool(ctx, vmPoolIdNumeric).
		UpdateVMPool(updateVMPool).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(vmPool, &VMPoolPrintConfig)
}

// VMPoolUpdateConfigExample prints a populated sdk.UpdateVMPool payload that can
// be edited and fed back to `vm-pool update --config-source`.
func VMPoolUpdateConfigExample(ctx context.Context) error {
	updateConfig := sdk.UpdateVMPool{
		Description:    sdk.PtrString("Updated VM pool description"),
		ManagementHost: sdk.PtrString("vcenter.example.com"),
		ManagementPort: sdk.PtrFloat32(443),
		Username:       sdk.PtrString("administrator@vsphere.local"),
		Password:       sdk.PtrString("password"),
		InMaintenance:  sdk.PtrFloat32(0),
		IsExperimental: sdk.PtrFloat32(0),
		Tags:           []string{"production", "vmware"},
	}

	return formatter.PrintResult(updateConfig, nil)
}

// VMPoolSync triggers a discovery pass on the VM pool (on VMware VCF this finds
// new virtual distributed switches) and prints the resulting job info.
func VMPoolSync(ctx context.Context, vmPoolId string) error {
	logger.Get().Info().Msgf("Syncing VM pool %s", vmPoolId)

	vmPoolIdNumeric, err := getVMPoolId(vmPoolId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	jobInfo, httpRes, err := client.VMPoolAPI.SyncVMPool(ctx, vmPoolIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	if jobInfo == nil {
		logger.Get().Info().Msgf("VM pool %s sync started", vmPoolId)
		return nil
	}

	return formatter.PrintResult(jobInfo, nil)
}

// VMPoolRefresh re-reads the VM pool information from the hypervisor (on VMware
// VCF this reports any new datastores) and prints the refreshed pool.
func VMPoolRefresh(ctx context.Context, vmPoolId string) error {
	logger.Get().Info().Msgf("Refreshing information of VM pool %s", vmPoolId)

	vmPoolIdNumeric, err := getVMPoolId(vmPoolId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	vmPool, httpRes, err := client.VMPoolAPI.RefreshVMPoolInformation(ctx, vmPoolIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	if vmPool == nil {
		logger.Get().Info().Msgf("VM pool %s information refreshed", vmPoolId)
		return nil
	}

	return formatter.PrintResult(vmPool, &VMPoolPrintConfig)
}

// VMPoolGetStatistics prints the aggregated RAM, disk and GPU usage of a pool.
func VMPoolGetStatistics(ctx context.Context, vmPoolId string) error {
	logger.Get().Info().Msgf("Getting statistics for VM pool %s", vmPoolId)

	vmPoolIdNumeric, err := getVMPoolId(vmPoolId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	statistics, httpRes, err := client.VMPoolAPI.GetVmPoolStatistics(ctx, vmPoolIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return printVMPoolStatistics(statistics)
}

// VMPoolGetContainers lists the containers running in a VM pool.
//
// The generated sdk.Container model cannot decode live responses (the API
// returns `hosts` as a plain string while the SDK declares []string), so the
// listing goes through the raw-body helpers, exactly like `container list`.
func VMPoolGetContainers(ctx context.Context, vmPoolId string) error {
	logger.Get().Info().Msgf("Getting containers for VM pool %s", vmPoolId)

	vmPoolIdNumeric, err := getVMPoolId(vmPoolId)
	if err != nil {
		return err
	}

	rawItems, meta, err := utils.FetchAllPagesRaw(func(page float32) (*http.Response, error) {
		return api.DoJSONRequest(ctx, http.MethodGet,
			fmt.Sprintf("/api/v2/vm-pools/%d/containers?page=%.0f&limit=100&sortBy=id:ASC", vmPoolIdNumeric, page), nil)
	})
	if err != nil {
		return err
	}

	return container.PrintContainersRaw(rawItems, meta)
}

// VMPoolGetClusterHost prints one cluster host of a VM pool.
func VMPoolGetClusterHost(ctx context.Context, vmPoolId string, hostId string) error {
	logger.Get().Info().Msgf("Getting cluster host %s in VM pool %s", hostId, vmPoolId)

	vmPoolIdNumeric, hostIdNumeric, err := getVMPoolAndHostId(vmPoolId, hostId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	host, httpRes, err := client.VMPoolAPI.GetVMPoolClusterHost(ctx, vmPoolIdNumeric, hostIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(host, &vmPoolHostPrintConfig)
}

// VMPoolUpdateClusterHost patches a cluster host of a VM pool.
func VMPoolUpdateClusterHost(ctx context.Context, vmPoolId string, hostId string, config []byte) error {
	logger.Get().Info().Msgf("Updating cluster host %s in VM pool %s", hostId, vmPoolId)

	vmPoolIdNumeric, hostIdNumeric, err := getVMPoolAndHostId(vmPoolId, hostId)
	if err != nil {
		return err
	}

	var updateHost sdk.UpdateVMPoolHost
	if err := utils.UnmarshalContent(config, &updateHost); err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	host, httpRes, err := client.VMPoolAPI.UpdateVMPoolClusterHost(ctx, vmPoolIdNumeric, hostIdNumeric).
		UpdateVMPoolHost(updateHost).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(host, &vmPoolHostPrintConfig)
}

// VMPoolUpdateClusterHostConfigExample prints a populated sdk.UpdateVMPoolHost
// payload for `vm-pool update-cluster-host --config-source`.
func VMPoolUpdateClusterHostConfigExample(ctx context.Context) error {
	updateConfig := sdk.UpdateVMPoolHost{
		AllowVMsToBeCreated:        sdk.PtrBool(true),
		AllowContainersToBeCreated: sdk.PtrBool(false),
	}

	return formatter.PrintResult(updateConfig, nil)
}

// VMPoolGetClusterHostStatistics prints the RAM, disk and GPU usage of a single
// cluster host.
func VMPoolGetClusterHostStatistics(ctx context.Context, vmPoolId string, hostId string) error {
	logger.Get().Info().Msgf("Getting statistics for cluster host %s in VM pool %s", hostId, vmPoolId)

	vmPoolIdNumeric, hostIdNumeric, err := getVMPoolAndHostId(vmPoolId, hostId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	statistics, httpRes, err := client.VMPoolAPI.GetVMPoolClusterHostStatistics(ctx, vmPoolIdNumeric, hostIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return printVMPoolStatistics(statistics)
}

// VMPoolGetClusterHostContainers lists the containers running on one cluster
// host. Like VMPoolGetContainers it decodes raw bodies, because the strict SDK
// container model rejects live payloads.
func VMPoolGetClusterHostContainers(ctx context.Context, vmPoolId string, hostId string) error {
	logger.Get().Info().Msgf("Getting containers for cluster host %s in VM pool %s", hostId, vmPoolId)

	vmPoolIdNumeric, hostIdNumeric, err := getVMPoolAndHostId(vmPoolId, hostId)
	if err != nil {
		return err
	}

	rawItems, meta, err := utils.FetchAllPagesRaw(func(page float32) (*http.Response, error) {
		return api.DoJSONRequest(ctx, http.MethodGet,
			fmt.Sprintf("/api/v2/vm-pools/%d/cluster-hosts/%d/containers?page=%.0f&limit=100&sortBy=id:ASC",
				vmPoolIdNumeric, hostIdNumeric, page), nil)
	})
	if err != nil {
		return err
	}

	return container.PrintContainersRaw(rawItems, meta)
}

// VMPoolGetClusterHostInterface prints one network interface of a cluster host.
func VMPoolGetClusterHostInterface(ctx context.Context, vmPoolId string, hostId string, interfaceId string) error {
	logger.Get().Info().Msgf("Getting interface %s of cluster host %s in VM pool %s", interfaceId, hostId, vmPoolId)

	vmPoolIdNumeric, hostIdNumeric, interfaceIdNumeric, err := getVMPoolHostAndInterfaceId(vmPoolId, hostId, interfaceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	hostInterface, httpRes, err := client.VMPoolAPI.
		GetVMPoolClusterHostInterface(ctx, vmPoolIdNumeric, hostIdNumeric, interfaceIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(hostInterface, &vmPoolHostInterfacePrintConfig)
}

// VMPoolUpdateClusterHostInterface patches one network interface of a cluster
// host. The only writable field is the interface status.
func VMPoolUpdateClusterHostInterface(ctx context.Context, vmPoolId string, hostId string, interfaceId string, config []byte) error {
	logger.Get().Info().Msgf("Updating interface %s of cluster host %s in VM pool %s", interfaceId, hostId, vmPoolId)

	vmPoolIdNumeric, hostIdNumeric, interfaceIdNumeric, err := getVMPoolHostAndInterfaceId(vmPoolId, hostId, interfaceId)
	if err != nil {
		return err
	}

	var updateInterface sdk.UpdateVMPoolHostInterface
	if err := utils.UnmarshalContent(config, &updateInterface); err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	hostInterface, httpRes, err := client.VMPoolAPI.
		UpdateVMPoolClusterHostInterface(ctx, vmPoolIdNumeric, hostIdNumeric, interfaceIdNumeric).
		UpdateVMPoolHostInterface(updateInterface).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(hostInterface, &vmPoolHostInterfacePrintConfig)
}

// VMPoolUpdateClusterHostInterfaceConfigExample prints a populated
// sdk.UpdateVMPoolHostInterface payload.
func VMPoolUpdateClusterHostInterfaceConfigExample(ctx context.Context) error {
	updateConfig := sdk.UpdateVMPoolHostInterface{
		Status: sdk.VMPOOLHOSTINTERFACESTATUS_MANAGED,
	}

	return formatter.PrintResult(updateConfig, nil)
}

// VMPoolGetClusterHostInterfaceNetworkDevices lists the network device
// assignments of one cluster host interface.
func VMPoolGetClusterHostInterfaceNetworkDevices(ctx context.Context, vmPoolId string, hostId string, interfaceId string) error {
	logger.Get().Info().Msgf("Getting network devices of interface %s of cluster host %s in VM pool %s", interfaceId, hostId, vmPoolId)

	vmPoolIdNumeric, hostIdNumeric, interfaceIdNumeric, err := getVMPoolHostAndInterfaceId(vmPoolId, hostId, interfaceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.VMPoolAPI.
		GetVMPoolClusterHostInterfaceNetworkDevices(ctx, vmPoolIdNumeric, hostIdNumeric, interfaceIdNumeric).
		SortBy([]string{"id:ASC"})

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &vmPoolHostInterfaceNetworkDevicePrintConfig)
}

// VMPoolGetClusterHostInterfaceNetworkDevice prints one network device
// assignment of a cluster host interface.
func VMPoolGetClusterHostInterfaceNetworkDevice(ctx context.Context, vmPoolId string, hostId string, interfaceId string, assignmentId string) error {
	logger.Get().Info().Msgf("Getting network device assignment %s of interface %s of cluster host %s in VM pool %s",
		assignmentId, interfaceId, hostId, vmPoolId)

	vmPoolIdNumeric, hostIdNumeric, interfaceIdNumeric, err := getVMPoolHostAndInterfaceId(vmPoolId, hostId, interfaceId)
	if err != nil {
		return err
	}

	assignmentIdNumeric, err := parseId(assignmentId, "network device assignment ID")
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	assignment, httpRes, err := client.VMPoolAPI.
		GetVMPoolClusterHostInterfaceNetworkDevice(ctx, vmPoolIdNumeric, hostIdNumeric, interfaceIdNumeric, assignmentIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(assignment, &vmPoolHostInterfaceNetworkDevicePrintConfig)
}

// VMPoolAddClusterHostInterfaceNetworkDevice links a cluster host interface to
// an interface of a network device (switch). networkDeviceRef accepts either a
// network device ID or its label.
func VMPoolAddClusterHostInterfaceNetworkDevice(ctx context.Context, vmPoolId string, hostId string, interfaceId string,
	networkDeviceRef string, networkDeviceInterfaceName string, config []byte) error {
	logger.Get().Info().Msgf("Adding network device to interface %s of cluster host %s in VM pool %s", interfaceId, hostId, vmPoolId)

	vmPoolIdNumeric, hostIdNumeric, interfaceIdNumeric, err := getVMPoolHostAndInterfaceId(vmPoolId, hostId, interfaceId)
	if err != nil {
		return err
	}

	var createAssignment sdk.CreateVMPoolHostInterfaceNetworkDevice
	if len(config) > 0 {
		if err := utils.UnmarshalContent(config, &createAssignment); err != nil {
			return err
		}
	} else {
		networkDevice, err := network_device.GetNetworkDeviceByIdOrLabel(ctx, networkDeviceRef)
		if err != nil {
			return err
		}

		networkDeviceId, err := utils.GetInt64FromString(networkDevice.Id)
		if err != nil {
			return fmt.Errorf("invalid network device ID %q: %w", networkDevice.Id, err)
		}

		createAssignment = sdk.CreateVMPoolHostInterfaceNetworkDevice{
			NetworkDeviceId:            networkDeviceId,
			NetworkDeviceInterfaceName: networkDeviceInterfaceName,
		}
	}

	client := api.GetApiClient(ctx)

	assignment, httpRes, err := client.VMPoolAPI.
		CreateVMPoolClusterHostInterfaceNetworkDevice(ctx, vmPoolIdNumeric, hostIdNumeric, interfaceIdNumeric).
		CreateVMPoolHostInterfaceNetworkDevice(createAssignment).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(assignment, &vmPoolHostInterfaceNetworkDevicePrintConfig)
}

// VMPoolRemoveClusterHostInterfaceNetworkDevice deletes one network device
// assignment of a cluster host interface.
func VMPoolRemoveClusterHostInterfaceNetworkDevice(ctx context.Context, vmPoolId string, hostId string, interfaceId string, assignmentId string) error {
	logger.Get().Info().Msgf("Removing network device assignment %s from interface %s of cluster host %s in VM pool %s",
		assignmentId, interfaceId, hostId, vmPoolId)

	vmPoolIdNumeric, hostIdNumeric, interfaceIdNumeric, err := getVMPoolHostAndInterfaceId(vmPoolId, hostId, interfaceId)
	if err != nil {
		return err
	}

	assignmentIdNumeric, err := parseId(assignmentId, "network device assignment ID")
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.VMPoolAPI.
		DeleteVMPoolClusterHostInterfaceNetworkDevice(ctx, vmPoolIdNumeric, hostIdNumeric, interfaceIdNumeric, assignmentIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Network device assignment %s removed", assignmentId)
	return nil
}

// VMPoolAddClusterHostInterfaceNetworkDeviceConfigExample prints a populated
// sdk.CreateVMPoolHostInterfaceNetworkDevice payload.
func VMPoolAddClusterHostInterfaceNetworkDeviceConfigExample(ctx context.Context) error {
	createConfig := sdk.CreateVMPoolHostInterfaceNetworkDevice{
		NetworkDeviceId:            1,
		NetworkDeviceInterfaceName: "Ethernet1/1",
	}

	return formatter.PrintResult(createConfig, nil)
}

func getVMPoolAndHostId(vmPoolId string, hostId string) (int64, int64, error) {
	vmPoolIdNumeric, err := getVMPoolId(vmPoolId)
	if err != nil {
		return 0, 0, err
	}

	hostIdNumeric, err := parseId(hostId, "cluster host ID")
	if err != nil {
		return 0, 0, err
	}

	return vmPoolIdNumeric, hostIdNumeric, nil
}

func getVMPoolHostAndInterfaceId(vmPoolId string, hostId string, interfaceId string) (int64, int64, int64, error) {
	vmPoolIdNumeric, hostIdNumeric, err := getVMPoolAndHostId(vmPoolId, hostId)
	if err != nil {
		return 0, 0, 0, err
	}

	interfaceIdNumeric, err := parseId(interfaceId, "cluster host interface ID")
	if err != nil {
		return 0, 0, 0, err
	}

	return vmPoolIdNumeric, hostIdNumeric, interfaceIdNumeric, nil
}

func parseId(value string, label string) (int64, error) {
	idNumeric, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		err := fmt.Errorf("invalid %s: '%s'", label, value)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return idNumeric, nil
}

func getVMPoolId(vmPoolId string) (int64, error) {
	vmPoolIdNumeric, err := strconv.ParseInt(vmPoolId, 10, 64)
	if err != nil {
		err := fmt.Errorf("invalid VM pool ID: '%s'", vmPoolId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return vmPoolIdNumeric, nil
}
