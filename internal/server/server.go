package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

// serverRaw works around the SDK bug where Server.ServerMetricsMetadata is typed
// as a map but the API returns an array, which makes typed unmarshalling fail.
// It keeps the original response body so json/yaml output stays complete while
// the table formatter reads the typed fields below.
type serverRaw struct {
	ServerId          float32     `json:"serverId"`
	SiteId            float32     `json:"siteId"`
	ServerTypeId      float32     `json:"serverTypeId"`
	ServerUUID        string      `json:"serverUUID"`
	SerialNumber      string      `json:"serialNumber"`
	ManagementAddress string      `json:"managementAddress"`
	Vendor            string      `json:"vendor"`
	Model             string      `json:"model"`
	ServerStatus      string      `json:"serverStatus"`
	Revision          float32     `json:"revision"`
	Links             interface{} `json:"links,omitempty"`

	// Operational fields shown in the table formats. Pointers keep a field the
	// API did not return distinct from a zero value, and render as empty.
	HealthState                      *string     `json:"healthState,omitempty"`
	HealthDetails                    []string    `json:"healthDetails,omitempty"`
	HealthLastCheckedTimestamp       *string     `json:"healthLastCheckedTimestamp,omitempty"`
	HealthMonitoringEnabled          interface{} `json:"healthMonitoringEnabled,omitempty"`
	FirmwareBaselineComplianceStatus *string     `json:"firmwareBaselineComplianceStatus,omitempty"`
	TargetFirmwareBaselineId         *float64    `json:"targetFirmwareBaselineId,omitempty"`
	CurrentFirmwareBaselineId        *float64    `json:"currentFirmwareBaselineId,omitempty"`
	DriftDetectionSyncStatus         *string     `json:"driftDetectionSyncStatus,omitempty"`
	FirmwareBackupEnabled            *bool       `json:"firmwareBackupEnabled,omitempty"`
	FirmwareDriftDetectionEnabled    *bool       `json:"firmwareDriftDetectionEnabled,omitempty"`
	Tags                             []string    `json:"tags,omitempty"`

	body json.RawMessage
}

// serverRawFields has the same layout as serverRaw but none of its methods, so it
// can be used for the default (un)marshalling of the typed fields.
type serverRawFields serverRaw

func (s *serverRaw) UnmarshalJSON(data []byte) error {
	var fields serverRawFields
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	*s = serverRaw(fields)
	s.body = append(json.RawMessage(nil), data...)

	return nil
}

func (s serverRaw) MarshalJSON() ([]byte, error) {
	if len(s.body) > 0 {
		return s.body, nil
	}

	return json.Marshal(serverRawFields(s))
}

var serverPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"ServerId": {
			Title: "ID",
			Order: 1,
		},
		"SiteId": {
			Title: "Site",
			Order: 2,
		},
		"ServerTypeId": {
			Title: "Type",
			Order: 3,
		},
		"ServerUUID": {
			Title: "UUID",
			Order: 4,
		},
		"SerialNumber": {
			Title: "S/N",
			Order: 5,
		},
		"ManagementAddress": {
			Title: "IP",
			Order: 6,
		},
		"Vendor": {
			Order: 7,
		},
		"Model": {
			Order: 8,
		},
		"ServerStatus": {
			Title:       "Status",
			Transformer: formatter.FormatStatusValue,
			Order:       9,
		},
		"HealthState": {
			Title:       "Health",
			Transformer: formatter.FormatHealthValue,
			Order:       10,
		},
		"HealthDetails": {
			Title:       "Health details",
			Transformer: formatter.FormatStringListValue,
			MaxWidth:    60,
			Order:       11,
		},
		"HealthLastCheckedTimestamp": {
			Title:       "Health checked",
			Transformer: formatter.FormatDateTimeValue,
			Order:       12,
		},
		"FirmwareBaselineComplianceStatus": {
			Title: "FW compliance",
			Order: 13,
		},
		"TargetFirmwareBaselineId": {
			Title:       "Target FW baseline",
			Transformer: formatter.FormatIdValue,
			Order:       14,
		},
		"CurrentFirmwareBaselineId": {
			Title:       "Current FW baseline",
			Transformer: formatter.FormatIdValue,
			Order:       15,
		},
		"DriftDetectionSyncStatus": {
			Title: "Drift sync",
			Order: 16,
		},
		"FirmwareBackupEnabled": {
			Title:       "FW backup",
			Transformer: formatter.FormatBooleanValue,
			Order:       17,
		},
		"FirmwareDriftDetectionEnabled": {
			Title:       "FW drift detection",
			Transformer: formatter.FormatBooleanValue,
			Order:       18,
		},
		"HealthMonitoringEnabled": {
			Title:       "Health monitoring",
			Transformer: formatter.FormatBooleanValue,
			Order:       19,
		},
		"Tags": {
			Title:       "Tags",
			Transformer: formatter.FormatStringListValue,
			Order:       20,
		},
	},
}

var serverWithCredentialsPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"ServerInfo": {
			Hidden:      true,
			InnerFields: serverPrintConfig.FieldsConfig,
		},
		"ServerCredentials": {
			Hidden: true,
			InnerFields: map[string]formatter.RecordFieldConfig{
				"Username": {
					Title: "User",
					Order: 21,
				},
				"Password": {
					Title: "Password",
					Order: 22,
				},
			},
		},
	},
}

// serverWithCredentials keeps the untagged field names the `server get
// --show-credentials` json/yaml output has always used.
type serverWithCredentials struct {
	ServerInfo        serverRaw
	ServerCredentials sdk.ServerCredentials
}

// ListFlags holds the server-side filters and pagination options of `server list`.
// Every filter maps 1:1 to a `filter.*` query parameter of GET /api/v2/servers.
type ListFlags struct {
	ShowCredentials bool

	FilterStatus                     []string
	FilterType                       []string
	FilterId                         []string
	FilterSiteId                     []string
	FilterVendor                     []string
	FilterModel                      []string
	FilterSerialNumber               []string
	FilterManagementAddress          []string
	FilterAdministrationState        []string
	FilterPowerStatus                []string
	FilterDhcpStatus                 []string
	FilterClass                      []string
	FilterDiskWipe                   []string
	FilterRequiresManualCleaning     []string
	FilterDriftDetectionSyncStatus   []string
	FilterFirmwareBaselineCompliance []string
	FilterHealthState                []string
	FilterFirmwareBackupEnabled      []string
	FilterFirmwareDriftDetection     []string
	FilterTags                       []string
	FilterInfrastructureId           []string
	FilterOwnerId                    []string

	Page  int
	Limit int
}

func ServerList(ctx context.Context, flags ListFlags) error {
	logger.Get().Info().Msgf("Listing servers")

	client := api.GetApiClient(ctx)

	request := client.ServerAPI.GetServers(ctx)

	apply := func(values []string, set func([]string) sdk.ServerAPIGetServersRequest) {
		if len(values) > 0 {
			request = set(utils.ProcessFilterStringSlice(values))
		}
	}

	apply(flags.FilterStatus, request.FilterServerStatus)
	apply(flags.FilterType, request.FilterServerTypeId)
	apply(flags.FilterId, request.FilterServerId)
	apply(flags.FilterSiteId, request.FilterSiteId)
	apply(flags.FilterVendor, request.FilterVendor)
	apply(flags.FilterModel, request.FilterModel)
	apply(flags.FilterSerialNumber, request.FilterSerialNumber)
	apply(flags.FilterManagementAddress, request.FilterManagementAddress)
	apply(flags.FilterAdministrationState, request.FilterAdministrationState)
	apply(flags.FilterPowerStatus, request.FilterPowerStatus)
	apply(flags.FilterDhcpStatus, request.FilterServerDhcpStatus)
	apply(flags.FilterClass, request.FilterServerClass)
	apply(flags.FilterDiskWipe, request.FilterServerDiskWipe)
	apply(flags.FilterRequiresManualCleaning, request.FilterRequiresManualCleaning)
	apply(flags.FilterDriftDetectionSyncStatus, request.FilterDriftDetectionSyncStatus)
	apply(flags.FilterFirmwareBaselineCompliance, request.FilterFirmwareBaselineComplianceStatus)
	apply(flags.FilterHealthState, request.FilterHealthState)
	apply(flags.FilterFirmwareBackupEnabled, request.FilterFirmwareBackupEnabled)
	apply(flags.FilterFirmwareDriftDetection, request.FilterFirmwareDriftDetectionEnabled)
	apply(flags.FilterTags, request.FilterTags)
	apply(flags.FilterInfrastructureId, request.FilterInstanceInfrastructureInfrastructureId)
	apply(flags.FilterOwnerId, request.FilterInstanceInfrastructureUserIdOwner)

	// Raw-body parse: the SDK Server model types ServerMetricsMetadata as a map but
	// the API returns an array, so typed unmarshalling fails on valid responses.
	var (
		rawItems []json.RawMessage
		meta     sdk.PaginatedResponseMeta
		err      error
	)

	switch {
	case flags.Page > 0:
		// Specific page requested; spans API pages when limit exceeds the API cap.
		rawItems, meta, err = utils.FetchPageWindowRaw(func(p, l float32) (*http.Response, error) {
			_, httpRes, _ := request.Page(p).Limit(l).Execute()
			return httpRes, nil
		}, flags.Page, flags.Limit)
	case flags.Limit > 0:
		rawItems, meta, err = utils.FetchUpToRaw(func(p, l float32) (*http.Response, error) {
			_, httpRes, _ := request.Page(p).Limit(l).Execute()
			return httpRes, nil
		}, flags.Limit)
	default:
		rawItems, meta, err = utils.FetchAllPagesRaw(func(p float32) (*http.Response, error) {
			_, httpRes, _ := request.Page(p).Limit(100).Execute()
			return httpRes, nil
		})
	}
	if err != nil {
		return err
	}

	servers, err := utils.UnmarshalRawItems[serverRaw](rawItems)
	if err != nil {
		return fmt.Errorf("failed to parse servers: %w", err)
	}

	if flags.ShowCredentials {
		type serverRawWithCredentials struct {
			ServerInfo        serverRaw             `json:"serverInfo"`
			ServerCredentials sdk.ServerCredentials `json:"serverCredentials"`
		}

		data := make([]serverRawWithCredentials, 0, len(servers))

		for _, server := range servers {
			serverCredentials, httpRes, err := client.ServerAPI.GetServerCredentials(ctx, int64(server.ServerId)).Execute()
			if err := response_inspector.InspectResponse(httpRes, err); err != nil {
				return err
			}

			data = append(data, serverRawWithCredentials{
				ServerInfo:        server,
				ServerCredentials: *serverCredentials,
			})
		}

		return utils.PrintAll(data, meta, len(data), &serverWithCredentialsPrintConfig)
	}

	return utils.PrintAllRaw(rawItems, servers, meta, len(servers), &serverPrintConfig)
}

func ServerGet(ctx context.Context, serverId string, showCredentials bool) error {
	logger.Get().Info().Msgf("Get server '%s'", serverId)

	serverIdNumeric, err := GetServerId(serverId)
	if err != nil {
		return err
	}

	serverInfo, err := getServerRaw(ctx, serverIdNumeric)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	if showCredentials {
		serverCredentials, httpRes, err := client.ServerAPI.GetServerCredentials(ctx, int64(serverInfo.ServerId)).Execute()
		if err := response_inspector.InspectResponse(httpRes, err); err != nil {
			return err
		}

		data := serverWithCredentials{
			ServerInfo:        serverInfo,
			ServerCredentials: *serverCredentials,
		}

		return formatter.PrintResult(data, &serverWithCredentialsPrintConfig)
	}

	return formatter.PrintResult(serverInfo, &serverPrintConfig)
}

func ServerRegister(ctx context.Context, serverConfig sdk.RegisterServer) error {
	logger.Get().Info().Msgf("Registering server")

	client := api.GetApiClient(ctx)

	registrationInfo, httpRes, err := client.ServerAPI.RegisterServer(ctx).RegisterServer(serverConfig).Execute()
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

func ServerReRegister(ctx context.Context, serverId string) error {
	logger.Get().Info().Msgf("Re-registering server '%s'", serverId)

	serverIdNumeric, revision, err := getServerIdAndRevision(ctx, serverId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	response, httpRes, err := client.ServerAPI.ReRegisterServer(ctx, serverIdNumeric).IfMatch(revision).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(response, &formatter.PrintConfig{
		FieldsConfig: map[string]formatter.RecordFieldConfig{
			"ServerId": {
				Title: "ID",
				Order: 1,
			},
			"JobInfo": {
				Hidden: true,
				InnerFields: map[string]formatter.RecordFieldConfig{
					"JobId": {
						Title: "Job Id",
						Order: 2,
					},
					"JobGroupId": {
						Title: "Job Group Id",
						Order: 3,
					},
				},
			},
		},
	})
}

func ServerFactoryReset(ctx context.Context, serverId string) error {
	logger.Get().Info().Msgf("Factory resetting server '%s'", serverId)

	serverIdNumeric, revision, err := getServerIdAndRevision(ctx, serverId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.ServerAPI.ResetServerToFactoryDefaults(ctx, serverIdNumeric).IfMatch(revision).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Factory reset initiated for server '%s'", serverId)

	return nil
}

func ServerArchive(ctx context.Context, serverId string) error {
	logger.Get().Info().Msgf("Archiving server '%s'", serverId)

	serverIdNumeric, revision, err := getServerIdAndRevision(ctx, serverId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.ServerAPI.ArchiveServer(ctx, serverIdNumeric).IfMatch(revision).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Server '%s' archived", serverId)

	return nil
}

func ServerDelete(ctx context.Context, serverId string) error {
	logger.Get().Info().Msgf("Deleting server '%s'", serverId)

	serverIdNumeric, revision, err := getServerIdAndRevision(ctx, serverId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.ServerAPI.DeleteServer(ctx, serverIdNumeric).IfMatch(revision).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Server '%s' deleted", serverId)

	return nil
}

func ServerPower(ctx context.Context, serverId string, action string) error {
	logger.Get().Info().Msgf("Setting power status for server '%s' to '%s'", serverId, action)

	validActions := map[string]bool{
		"on":    true,
		"off":   true,
		"reset": true,
		"cycle": true,
		"soft":  true,
	}

	if !validActions[action] {
		return fmt.Errorf("invalid power action: '%s'. Valid actions are: on, off, reset, cycle, soft", action)
	}

	powerSet := sdk.ServerPowerSet{
		PowerCommand: action,
	}

	serverIdNumeric, revision, err := getServerIdAndRevision(ctx, serverId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.ServerAPI.SetServerPowerState(ctx, serverIdNumeric).ServerPowerSet(powerSet).IfMatch(revision).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Power status for server '%s' set to '%s'", serverId, action)

	return nil
}

func ServerPowerStatus(ctx context.Context, serverId string) error {
	logger.Get().Info().Msgf("Getting power status for server '%s'", serverId)

	serverIdNumeric, err := GetServerId(serverId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	powerStatus, httpRes, err := client.ServerAPI.GetServerPowerStatus(ctx, serverIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Power status for server '%s' is '%s'", serverId, powerStatus)

	return formatter.PrintResult(powerStatus, nil)
}

func ServerUpdate(ctx context.Context, serverId string, config []byte) error {
	logger.Get().Info().Msgf("Updating server '%s'", serverId)

	var updateConfig sdk.UpdateServer
	err := utils.UnmarshalContent(config, &updateConfig)
	if err != nil {
		return err
	}

	serverIdNumeric, revision, err := getServerIdAndRevision(ctx, serverId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	_, httpRes, sdkErr := client.ServerAPI.UpdateServer(ctx, serverIdNumeric).UpdateServer(updateConfig).IfMatch(revision).Execute()

	serverInfo, err := parseServerRaw(httpRes, sdkErr)
	if err != nil {
		return err
	}

	return formatter.PrintResult(serverInfo, &serverPrintConfig)
}

func ServerUpdateIpmiCredentials(ctx context.Context, serverId string, username string, password string) error {
	logger.Get().Info().Msgf("Updating IPMI credentials for server '%s'", serverId)

	serverIdNumeric, revision, err := getServerIdAndRevision(ctx, serverId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	credentials := sdk.UpdateServerIpmiCredentials{
		Username: sdk.PtrString(username),
		Password: sdk.PtrString(password),
	}

	serverCredentials, httpRes, err := client.ServerAPI.UpdateServerIpmiCredentials(ctx, serverIdNumeric).UpdateServerIpmiCredentials(credentials).IfMatch(revision).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("IPMI credentials for server '%s' updated", serverId)

	return formatter.PrintResult(serverCredentials, nil)
}

func ServerEnableSnmp(ctx context.Context, serverId string) error {
	logger.Get().Info().Msgf("Enabling SNMP for server '%s'", serverId)

	serverIdNumeric, revision, err := getServerIdAndRevision(ctx, serverId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	_, httpRes, err := client.ServerAPI.UpdateServerEnableSnmp(ctx, serverIdNumeric).IfMatch(revision).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("SNMP enabled for server '%s'", serverId)

	return nil
}

func ServerEnableSyslog(ctx context.Context, serverId string) error {
	logger.Get().Info().Msgf("Enabling syslog for server '%s'", serverId)

	serverIdNumeric, revision, err := getServerIdAndRevision(ctx, serverId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.ServerAPI.EnableServerSyslog(ctx, serverIdNumeric).IfMatch(revision).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Syslog enabled for server '%s'", serverId)

	return nil
}

func ServerVncInfo(ctx context.Context, serverId string) error {
	logger.Get().Info().Msgf("Getting VNC info for server '%s'", serverId)

	serverIdNumeric, err := GetServerId(serverId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	vncInfo, httpRes, err := client.ServerAPI.GetServerVNCInfo(ctx, serverIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(vncInfo, &formatter.PrintConfig{
		FieldsConfig: map[string]formatter.RecordFieldConfig{
			"ActiveSessions": {
				Title: "Active Sessions",
				Order: 1,
			},
			"MaxSessions": {
				Title: "Max Sessions",
				Order: 2,
			},
			"Port": {
				Title: "Port",
				Order: 3,
			},
			"Timeout": {
				Title: "Timeout",
				Order: 4,
			},
			"Enable": {
				Title: "Status",
				Order: 5,
			},
		},
	})
}

func ServerRemoteConsoleInfo(ctx context.Context, serverId string) error {
	logger.Get().Info().Msgf("Getting remote console info for server '%s'", serverId)

	serverIdNumeric, err := GetServerId(serverId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	consoleInfo, httpRes, err := client.ServerAPI.GetServerRemoteConsoleInfo(ctx, serverIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(consoleInfo, &formatter.PrintConfig{
		FieldsConfig: map[string]formatter.RecordFieldConfig{
			"ActiveConnections": {
				Title: "Active Connections",
				Order: 1,
			},
		},
	})
}

func ServerCapabilities(ctx context.Context, serverId string) error {
	logger.Get().Info().Msgf("Getting capabilities for server '%s'", serverId)

	serverIdNumeric, err := GetServerId(serverId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	capabilities, httpRes, err := client.ServerAPI.GetServerCapabilities(ctx, serverIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(capabilities, &formatter.PrintConfig{
		FieldsConfig: map[string]formatter.RecordFieldConfig{
			"FirmwareUpgradeSupported": {
				Title: "Firmware Upgrade",
				Order: 1,
			},
			"FirmwareUpgradeApplyOnRebootSupported": {
				Title: "Apply Firmware Upgrade On Reboot",
				Order: 2,
			},
			"VncEnabled": {
				Title: "VNC Enabled",
				Order: 3,
			},
		},
	})
}

func ServerIdentify(ctx context.Context, serverId string) error {
	logger.Get().Info().Msgf("Identifying server '%s'", serverId)

	id, err := GetServerId(serverId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.ServerAPI.IdentifyServer(ctx, id).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	fmt.Println("Server identification initiated (chassis LED blinking).")
	return nil
}

func GetServerId(serverId string) (int64, error) {
	serverIdNumeric, err := strconv.ParseInt(serverId, 10, 64)
	if err != nil {
		err := fmt.Errorf("invalid server ID: '%s'", serverId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return int64(serverIdNumeric), nil
}

func getServerIdAndRevision(ctx context.Context, serverId string) (int64, string, error) {
	serverIdNumeric, err := GetServerId(serverId)
	if err != nil {
		return 0, "", err
	}

	server, err := getServerRaw(ctx, serverIdNumeric)
	if err != nil {
		return 0, "", err
	}

	return serverIdNumeric, strconv.Itoa(int(server.Revision)), nil
}

// getServerRaw fetches a single server with a raw-body parse: the SDK Server model
// fails typed unmarshalling on valid responses (ServerMetricsMetadata map vs array,
// dpuInfo entries without a password).
func getServerRaw(ctx context.Context, serverIdNumeric int64) (serverRaw, error) {
	client := api.GetApiClient(ctx)

	_, httpRes, sdkErr := client.ServerAPI.GetServerInfo(ctx, serverIdNumeric).Execute()

	return parseServerRaw(httpRes, sdkErr)
}

// parseServerRaw turns a server response into serverRaw, ignoring the SDK's own
// typed-decoding error while still surfacing API errors.
func parseServerRaw(httpRes *http.Response, sdkErr error) (serverRaw, error) {
	var server serverRaw

	if httpRes != nil && httpRes.StatusCode >= 400 {
		if err := response_inspector.InspectResponse(httpRes, sdkErr); err != nil {
			return server, err
		}
	} else if httpRes == nil {
		return server, sdkErr
	}

	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return server, fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(body, &server); err != nil {
		return server, fmt.Errorf("failed to parse server: %w", err)
	}

	return server, nil
}
