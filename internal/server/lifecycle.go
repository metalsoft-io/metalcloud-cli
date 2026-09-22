package server

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

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

var serverStatisticsPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"ServerStatus": {
			Title:       "Servers By Status",
			Transformer: formatCountMap,
			MaxWidth:    60,
			Order:       1,
		},
		"Site": {
			Title:       "Servers By Site",
			Transformer: formatCountMap,
			MaxWidth:    60,
			Order:       2,
		},
	},
}

var serverJobInfoPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"ServerId": {
			Title: "ID",
			Order: 1,
		},
		"Revision": {
			Title: "Revision",
			Order: 2,
		},
		"JobInfo": {
			Hidden: true,
			InnerFields: map[string]formatter.RecordFieldConfig{
				"JobId": {
					Title: "Job Id",
					Order: 3,
				},
				"JobGroupId": {
					Title: "Job Group Id",
					Order: 4,
				},
			},
		},
	},
}

// ServerStatistics prints the aggregated server counts by status and by site.
func ServerStatistics(ctx context.Context) error {
	logger.Get().Info().Msgf("Getting server statistics")

	client := api.GetApiClient(ctx)

	statistics, httpRes, err := client.ServerAPI.GetServersStatistics(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(statistics, &serverStatisticsPrintConfig)
}

// ServerHardwareRescan re-reads the hardware inventory of a server. When
// rebootAllowed is set the server may be rebooted to run a full LLDP interface
// discovery.
func ServerHardwareRescan(ctx context.Context, serverId string, rebootAllowed bool) error {
	logger.Get().Info().Msgf("Rescanning hardware of server '%s'", serverId)

	serverIdNumeric, revision, err := getServerIdAndRevision(ctx, serverId)
	if err != nil {
		return err
	}

	// The body is always sent: skipping the optional setter would serialize a
	// literal null, which the API rejects with 400.
	rescanRequest := sdk.HardwareRescanServerRequest{
		RebootAllowed: sdk.PtrBool(rebootAllowed),
	}

	client := api.GetApiClient(ctx)

	response, httpRes, err := client.ServerAPI.
		ServerHardwareRescan(ctx, serverIdNumeric).
		HardwareRescanServerRequest(rescanRequest).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(response, &serverJobInfoPrintConfig)
}

// ServerRegisterProduction registers a server that is already running a
// production workload, keeping the workload in place.
func ServerRegisterProduction(ctx context.Context, serverConfig sdk.RegisterProductionServer) error {
	logger.Get().Info().Msgf("Registering production server")

	client := api.GetApiClient(ctx)

	registrationInfo, httpRes, err := client.ServerAPI.
		RegisterProductionServer(ctx).
		RegisterProductionServer(serverConfig).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(registrationInfo, &formatter.PrintConfig{
		FieldsConfig: map[string]formatter.RecordFieldConfig{
			"ServerId": {
				Title: "ID",
				Order: 1,
			},
			"ServerUUID": {
				Title: "UUID",
				Order: 2,
			},
			"SerialNumber": {
				Title: "S/N",
				Order: 3,
			},
			"JobInfo": {
				Hidden: true,
				InnerFields: map[string]formatter.RecordFieldConfig{
					"JobId": {
						Title: "Job Id",
						Order: 4,
					},
					"JobGroupId": {
						Title: "Job Group Id",
						Order: 5,
					},
				},
			},
		},
	})
}

func ServerRegisterProductionConfigExample(ctx context.Context) error {
	serverConfig := sdk.RegisterProductionServer{
		SiteId:                1,
		ServerUUID:            sdk.PtrString("00000000-0000-0000-0000-000000000000"),
		SerialNumber:          sdk.PtrString("ABC1234"),
		ManagementAddress:     sdk.PtrString("10.0.0.1"),
		Username:              sdk.PtrString("admin"),
		Password:              sdk.PtrString("password"),
		BmcMacAddress:         sdk.PtrString("AA:BB:CC:DD:EE:FF"),
		Vendor:                sdk.PtrString("Dell"),
		Model:                 sdk.PtrString("PowerEdge R640"),
		RegistrationProfileId: sdk.PtrInt64(10),
		Settings: sdk.RegisterProductionServerSettings{
			InfrastructureId: 100,
			OsTemplateId:     sdk.PtrInt64(20),
			InterfaceConnections: []sdk.ServerInterfaceConnection{
				{
					ServerInterfaceMacAddress: "AA:BB:CC:DD:EE:00",
					NetworkDevicePortId:       "Ethernet0",
					NetworkDeviceHostname:     "leaf-01",
				},
			},
		},
	}

	return formatter.PrintResult(serverConfig, nil)
}

// ServerConfigExampleKinds lists the configuration kinds understood by
// ServerConfigExample, in the order they are shown in the help text.
var ServerConfigExampleKinds = []string{
	"register-production",
	"import-unmanaged",
	"connect-interface",
	"set-interfaces-default-fabric",
	"set-interfaces-redundancy-group",
}

// ServerConfigExample prints an example of the configuration accepted by the
// server sub-command named by kind.
func ServerConfigExample(ctx context.Context, kind string) error {
	switch kind {
	case "register-production":
		return ServerRegisterProductionConfigExample(ctx)
	case "import-unmanaged":
		return ServerImportUnmanagedConfigExample(ctx)
	case "connect-interface":
		return ServerConnectInterfaceConfigExample(ctx)
	case "set-interfaces-default-fabric":
		return ServerSetInterfacesDefaultFabricConfigExample(ctx)
	case "set-interfaces-redundancy-group":
		return ServerSetInterfacesRedundancyGroupConfigExample(ctx)
	}

	err := fmt.Errorf("unknown configuration kind '%s', expected one of: %s", kind, strings.Join(ServerConfigExampleKinds, ", "))
	logger.Get().Error().Err(err).Msg("")
	return err
}
