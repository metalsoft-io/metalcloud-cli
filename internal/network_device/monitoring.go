package network_device

import (
	"context"
	"fmt"
	"sort"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var networkDeviceStatisticsPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"DeviceCount": {
			Title: "Devices",
			Order: 1,
		},
		"PortCount": {
			Title: "Ports",
			Order: 2,
		},
	},
}

var networkDeviceAgentInfoPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"NetworkDeviceId": {
			Title: "Device",
			Order: 1,
		},
		"AgentId": {
			Title: "Agent",
			Order: 2,
		},
	},
}

// agentInfoRow is the table projection of one entry of the device -> agent
// allocation map.
type agentInfoRow struct {
	NetworkDeviceId string
	AgentId         string
}

var networkDeviceHealthSummaryPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"OverallSeverity": {
			Title:       "Severity",
			Transformer: formatter.FormatStatusValue,
			Order:       1,
		},
		"TrendDirection": {
			Title: "Trend",
			Order: 2,
		},
		"LastUpdated": {
			Title: "Last Updated",
			Order: 3,
		},
		"SuspectedRootCauses": {
			Title:       "Suspected Root Causes",
			MaxWidth:    60,
			Transformer: formatNetworkDeviceStringList,
			Order:       4,
		},
		"KeyFindings": {
			Title:    "Key Findings",
			MaxWidth: 80,
			Order:    5,
		},
		"DetectedIssues": {
			Title: "Detected Issues",
			Order: 6,
		},
	},
}

// healthSummaryRow is the table projection of a health summary: the nested
// notable-event and detected-issue objects are reduced to a count because the
// tabular formatter renders struct-valued fields as empty cells. The json and
// yaml formats keep the full summary.
type healthSummaryRow struct {
	OverallSeverity     string
	TrendDirection      string
	LastUpdated         string
	SuspectedRootCauses []string
	KeyFindings         string
	DetectedIssues      string
}

// NetworkDeviceSnmpMonitoringEnable subscribes a network device to SNMP
// monitoring.
func NetworkDeviceSnmpMonitoringEnable(ctx context.Context, networkDeviceRef string) error {
	logger.Get().Info().Msgf("Enabling SNMP monitoring for network device '%s'", networkDeviceRef)

	networkDeviceIdNumeric, err := resolveNetworkDeviceId(ctx, networkDeviceRef)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.NetworkDeviceAPI.
		EnableNetworkDeviceSnmpMonitoring(ctx, networkDeviceIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("SNMP monitoring enabled for network device '%s'", networkDeviceRef)
	return nil
}

// NetworkDeviceSnmpMonitoringDisable unsubscribes a network device from SNMP
// monitoring.
func NetworkDeviceSnmpMonitoringDisable(ctx context.Context, networkDeviceRef string) error {
	logger.Get().Info().Msgf("Disabling SNMP monitoring for network device '%s'", networkDeviceRef)

	networkDeviceIdNumeric, err := resolveNetworkDeviceId(ctx, networkDeviceRef)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.NetworkDeviceAPI.
		DisableNetworkDeviceSnmpMonitoring(ctx, networkDeviceIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("SNMP monitoring disabled for network device '%s'", networkDeviceRef)
	return nil
}

// NetworkDeviceSnmpMonitoringEnableBatch subscribes several network devices to
// SNMP monitoring at once, selected by device id, site id and/or fabric id.
func NetworkDeviceSnmpMonitoringEnableBatch(ctx context.Context, networkDeviceIds []string, siteIds []string, fabricIds []string) error {
	changeStatus, err := buildSnmpMonitoringSelection(networkDeviceIds, siteIds, fabricIds)
	if err != nil {
		return err
	}

	logger.Get().Info().Msgf("Enabling SNMP monitoring in batch")

	client := api.GetApiClient(ctx)

	httpRes, err := client.NetworkDeviceAPI.
		EnableNetworkDeviceSnmpMonitoringBatch(ctx).
		NetworkDeviceSNMPMonitoringChangeStatus(*changeStatus).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("SNMP monitoring enabled in batch")
	return nil
}

// NetworkDeviceSnmpMonitoringDisableBatch unsubscribes several network devices
// from SNMP monitoring at once.
func NetworkDeviceSnmpMonitoringDisableBatch(ctx context.Context, networkDeviceIds []string, siteIds []string, fabricIds []string) error {
	changeStatus, err := buildSnmpMonitoringSelection(networkDeviceIds, siteIds, fabricIds)
	if err != nil {
		return err
	}

	logger.Get().Info().Msgf("Disabling SNMP monitoring in batch")

	client := api.GetApiClient(ctx)

	httpRes, err := client.NetworkDeviceAPI.
		DisableNetworkDeviceSnmpMonitoringBatch(ctx).
		NetworkDeviceSNMPMonitoringChangeStatus(*changeStatus).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("SNMP monitoring disabled in batch")
	return nil
}

// buildSnmpMonitoringSelection converts the device/site/fabric id selections
// into the request body. The body is always sent: an SDK action request whose
// optional body setter is skipped serializes a literal null, which the API
// rejects with 400.
func buildSnmpMonitoringSelection(networkDeviceIds []string, siteIds []string, fabricIds []string) (*sdk.NetworkDeviceSNMPMonitoringChangeStatus, error) {
	if len(networkDeviceIds) == 0 && len(siteIds) == 0 && len(fabricIds) == 0 {
		return nil, fmt.Errorf("no selection given: specify at least one network device id, site id or fabric id")
	}

	changeStatus := sdk.NetworkDeviceSNMPMonitoringChangeStatus{}

	if len(networkDeviceIds) > 0 {
		ids, err := utils.GetInt64SliceFromStrings(networkDeviceIds)
		if err != nil {
			return nil, fmt.Errorf("invalid network device id: %w", err)
		}
		changeStatus.NetworkDeviceIds = ids
	}

	if len(siteIds) > 0 {
		ids, err := utils.GetInt64SliceFromStrings(siteIds)
		if err != nil {
			return nil, fmt.Errorf("invalid site id: %w", err)
		}
		changeStatus.SiteIds = ids
	}

	if len(fabricIds) > 0 {
		ids, err := utils.GetInt64SliceFromStrings(fabricIds)
		if err != nil {
			return nil, fmt.Errorf("invalid fabric id: %w", err)
		}
		changeStatus.FabricIds = ids
	}

	return &changeStatus, nil
}

// NetworkDeviceSnmpServiceEnable turns on the SNMP agent running on the device
// itself, optionally overriding port, community and contact.
func NetworkDeviceSnmpServiceEnable(ctx context.Context, networkDeviceRef string, port *int32, community *string, contact *string) error {
	logger.Get().Info().Msgf("Enabling SNMP service on network device '%s'", networkDeviceRef)

	networkDeviceIdNumeric, err := resolveNetworkDeviceId(ctx, networkDeviceRef)
	if err != nil {
		return err
	}

	// The body is always sent: skipping the optional setter would serialize a
	// literal null, which the API rejects with 400.
	snmpConfig := sdk.NetworkDeviceSNMPConfig{
		Port:      port,
		Community: community,
		Contact:   contact,
	}

	client := api.GetApiClient(ctx)

	jobInfo, httpRes, err := client.NetworkDeviceAPI.
		EnableNetworkDeviceSnmpService(ctx, networkDeviceIdNumeric).
		NetworkDeviceSNMPConfig(snmpConfig).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(jobInfo, nil)
}

// NetworkDeviceSnmpServiceDisable turns off the SNMP agent running on the
// device itself.
func NetworkDeviceSnmpServiceDisable(ctx context.Context, networkDeviceRef string) error {
	logger.Get().Info().Msgf("Disabling SNMP service on network device '%s'", networkDeviceRef)

	networkDeviceIdNumeric, err := resolveNetworkDeviceId(ctx, networkDeviceRef)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	jobInfo, httpRes, err := client.NetworkDeviceAPI.
		DisableNetworkDeviceSnmpService(ctx, networkDeviceIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(jobInfo, nil)
}

// NetworkDeviceDisableSyslog unsubscribes a network device from remote syslog.
func NetworkDeviceDisableSyslog(ctx context.Context, networkDeviceRef string) error {
	logger.Get().Info().Msgf("Disabling syslog for network device '%s'", networkDeviceRef)

	networkDeviceIdNumeric, err := resolveNetworkDeviceId(ctx, networkDeviceRef)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	jobInfo, httpRes, err := client.NetworkDeviceAPI.
		DisableNetworkDeviceSyslog(ctx, networkDeviceIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(jobInfo, nil)
}

// NetworkDeviceSnmpAgentInfo shows which monitoring agent polls each of the
// selected network devices.
func NetworkDeviceSnmpAgentInfo(ctx context.Context, networkDeviceIds []string, siteIds []string, fabricIds []string) error {
	logger.Get().Info().Msgf("Getting SNMP monitoring agent allocation")

	client := api.GetApiClient(ctx)

	request := client.NetworkDeviceAPI.GetNetworkDeviceSNMPMonitoringAgentInfoBatch(ctx)

	if len(networkDeviceIds) > 0 {
		ids, err := utils.GetInt64SliceFromStrings(networkDeviceIds)
		if err != nil {
			return fmt.Errorf("invalid network device id: %w", err)
		}
		request = request.NetworkDeviceIds(ids)
	}

	if len(siteIds) > 0 {
		ids, err := utils.GetInt64SliceFromStrings(siteIds)
		if err != nil {
			return fmt.Errorf("invalid site id: %w", err)
		}
		request = request.SiteIds(ids)
	}

	if len(fabricIds) > 0 {
		ids, err := utils.GetInt64SliceFromStrings(fabricIds)
		if err != nil {
			return fmt.Errorf("invalid fabric id: %w", err)
		}
		request = request.FabricIds(ids)
	}

	agentInfo, httpRes, err := request.Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	rows := make([]agentInfoRow, 0, len(agentInfo.AllocationInfo))
	for networkDeviceId, agentId := range agentInfo.AllocationInfo {
		rows = append(rows, agentInfoRow{NetworkDeviceId: networkDeviceId, AgentId: agentId})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].NetworkDeviceId < rows[j].NetworkDeviceId })

	return printRecords(agentInfo, rows, &networkDeviceAgentInfoPrintConfig)
}

// NetworkDeviceHealthSummary shows the accumulated health assessment of a
// network device.
func NetworkDeviceHealthSummary(ctx context.Context, networkDeviceRef string) error {
	logger.Get().Info().Msgf("Getting health summary of network device '%s'", networkDeviceRef)

	networkDeviceIdNumeric, err := resolveNetworkDeviceId(ctx, networkDeviceRef)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	healthSummary, httpRes, err := client.NetworkDeviceAPI.
		GetNetworkDeviceHealthSummary(ctx, networkDeviceIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	// The endpoint answers 204 with no body when health monitoring has not yet
	// written a summary for the device; the SDK then returns a nil summary.
	if healthSummary == nil {
		logger.Get().Info().Msgf("No health summary recorded for network device '%s'", networkDeviceRef)
		return nil
	}

	row := healthSummaryRow{
		OverallSeverity:     healthSummary.OverallSeverity,
		TrendDirection:      healthSummary.TrendDirection,
		LastUpdated:         healthSummary.LastUpdated,
		SuspectedRootCauses: healthSummary.SuspectedRootCauses,
		KeyFindings:         healthSummary.KeyFindings,
		DetectedIssues:      fmt.Sprintf("%d issue(s)", len(healthSummary.DetectedIssues)),
	}

	return printRecords(healthSummary, row, &networkDeviceHealthSummaryPrintConfig)
}

// NetworkDeviceSetHealthMonitoringFilter sets the network device filter applied
// to the health monitoring WebSocket stream of the given socket. The endpoint
// is served as a GET and returns no body.
func NetworkDeviceSetHealthMonitoringFilter(ctx context.Context, socketId string, filterId []string, filterSiteId []string, filterStatus []string, filterHealthStatus []string) error {
	logger.Get().Info().Msgf("Setting health monitoring filter for socket '%s'", socketId)

	client := api.GetApiClient(ctx)

	request := client.NetworkDeviceAPI.SetNetworkDeviceHealthMonitoringFilter(ctx).SocketId(socketId)

	if len(filterId) > 0 {
		request = request.FilterId(utils.ProcessFilterStringSlice(filterId))
	}
	if len(filterSiteId) > 0 {
		request = request.FilterSiteId(utils.ProcessFilterStringSlice(filterSiteId))
	}
	if len(filterStatus) > 0 {
		request = request.FilterStatus(utils.ProcessFilterStringSlice(filterStatus))
	}
	if len(filterHealthStatus) > 0 {
		request = request.FilterHealthStatus(utils.ProcessFilterStringSlice(filterHealthStatus))
	}

	httpRes, err := request.Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Health monitoring filter set for socket '%s'", socketId)
	return nil
}

// NetworkDeviceStatisticsGet shows global network device counters.
func NetworkDeviceStatisticsGet(ctx context.Context) error {
	logger.Get().Info().Msgf("Getting network device statistics")

	client := api.GetApiClient(ctx)

	statistics, httpRes, err := client.NetworkDeviceAPI.GetNetworkDeviceStatistics(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(statistics, &networkDeviceStatisticsPrintConfig)
}
