package network_device

import (
	"context"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
)

var networkDeviceDriftPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"NetworkDeviceId": {
			Title: "Device",
			Order: 2,
		},
		"SnapshotId": {
			Title:    "Snapshot",
			MaxWidth: 20,
			Order:    3,
		},
		"ConfigurationDrift": {
			Title:    "Drift",
			MaxWidth: 60,
			Order:    4,
		},
		"AcknowledgedBy": {
			Title: "Acknowledged By",
			Order: 5,
		},
		"AcknowledgeTimestamp": {
			Title:       "Acknowledged",
			Transformer: formatter.FormatDateTimeValue,
			Order:       6,
		},
		"CreatedTimestamp": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       7,
		},
	},
}

var networkDeviceSnapshotPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Oid": {
			Title:    "OID",
			MaxWidth: 20,
			Order:    1,
		},
		"Kind": {
			Order: 2,
		},
		"Message": {
			Title:    "Message",
			MaxWidth: 60,
			Order:    3,
		},
		"Timestamp": {
			Title: "Timestamp",
			Order: 4,
		},
	},
}

// NetworkDeviceDriftList lists the configuration drift entries recorded for a
// network device.
func NetworkDeviceDriftList(ctx context.Context, networkDeviceRef string) error {
	logger.Get().Info().Msgf("Listing configuration drift of network device '%s'", networkDeviceRef)

	networkDeviceIdNumeric, err := resolveNetworkDeviceId(ctx, networkDeviceRef)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.NetworkDeviceAPI.
		GetNetworkDeviceDriftHistory(ctx, float32(networkDeviceIdNumeric)).
		SortBy([]string{"id:ASC"})

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &networkDeviceDriftPrintConfig)
}

// NetworkDeviceDriftGet shows one configuration drift entry of a network device.
func NetworkDeviceDriftGet(ctx context.Context, networkDeviceRef string, driftId string) error {
	logger.Get().Info().Msgf("Getting configuration drift %s of network device '%s'", driftId, networkDeviceRef)

	networkDeviceIdNumeric, err := resolveNetworkDeviceId(ctx, networkDeviceRef)
	if err != nil {
		return err
	}

	driftIdNumeric, err := utils.GetInt64FromString(driftId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	drift, httpRes, err := client.NetworkDeviceAPI.
		GetNetworkDeviceDriftHistoryById(ctx, float32(networkDeviceIdNumeric), float32(driftIdNumeric)).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(drift, &networkDeviceDriftPrintConfig)
}

// NetworkDeviceDriftAcknowledge marks one configuration drift entry as reviewed.
func NetworkDeviceDriftAcknowledge(ctx context.Context, networkDeviceRef string, driftId string) error {
	logger.Get().Info().Msgf("Acknowledging configuration drift %s of network device '%s'", driftId, networkDeviceRef)

	networkDeviceIdNumeric, err := resolveNetworkDeviceId(ctx, networkDeviceRef)
	if err != nil {
		return err
	}

	driftIdNumeric, err := utils.GetInt64FromString(driftId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	drift, httpRes, err := client.NetworkDeviceAPI.
		AcknowledgeNetworkDeviceDrift(ctx, networkDeviceIdNumeric, driftIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(drift, &networkDeviceDriftPrintConfig)
}

// NetworkDeviceSnapshotList lists the configuration snapshots stored for a
// network device. kind optionally narrows the listing to one snapshot class.
func NetworkDeviceSnapshotList(ctx context.Context, networkDeviceRef string, kind string) error {
	logger.Get().Info().Msgf("Listing snapshots of network device '%s'", networkDeviceRef)

	networkDeviceIdNumeric, err := resolveNetworkDeviceId(ctx, networkDeviceRef)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.NetworkDeviceAPI.GetNetworkDeviceSnapshots(ctx, networkDeviceIdNumeric)
	if kind != "" {
		request = request.Kind(kind)
	}

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &networkDeviceSnapshotPrintConfig)
}

// NetworkDeviceSyncTargetSnapshot points the drift-detection target snapshot of
// a network device at its latest snapshot, clearing the reported drift.
func NetworkDeviceSyncTargetSnapshot(ctx context.Context, networkDeviceRef string) error {
	logger.Get().Info().Msgf("Syncing target snapshot of network device '%s' with the latest snapshot", networkDeviceRef)

	networkDeviceIdNumeric, err := resolveNetworkDeviceId(ctx, networkDeviceRef)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.NetworkDeviceAPI.
		NetworkDeviceSyncTargetSnapshotWithLatestSnapshot(ctx, float32(networkDeviceIdNumeric)).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Target snapshot of network device '%s' synced with the latest snapshot", networkDeviceRef)
	return nil
}
