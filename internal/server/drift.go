package server

import (
	"context"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
)

var serverDriftPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"ServerId": {
			Title: "Server",
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

var serverSnapshotPrintConfig = formatter.PrintConfig{
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

// ServerDriftList lists the configuration drift entries recorded for a server.
func ServerDriftList(ctx context.Context, serverId string) error {
	logger.Get().Info().Msgf("Listing configuration drift of server '%s'", serverId)

	serverIdNumeric, err := GetServerId(serverId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.ServerAPI.
		GetServerDriftHistory(ctx, float32(serverIdNumeric)).
		SortBy([]string{"id:ASC"})

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &serverDriftPrintConfig)
}

// ServerDriftGet shows one configuration drift entry of a server.
func ServerDriftGet(ctx context.Context, serverId string, driftId string) error {
	logger.Get().Info().Msgf("Getting configuration drift %s of server '%s'", driftId, serverId)

	serverIdNumeric, err := GetServerId(serverId)
	if err != nil {
		return err
	}

	driftIdNumeric, err := utils.GetInt64FromString(driftId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	drift, httpRes, err := client.ServerAPI.
		GetServerDriftHistoryById(ctx, float32(serverIdNumeric), float32(driftIdNumeric)).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(drift, &serverDriftPrintConfig)
}

// ServerDriftAcknowledge marks one configuration drift entry as reviewed.
func ServerDriftAcknowledge(ctx context.Context, serverId string, driftId string) error {
	logger.Get().Info().Msgf("Acknowledging configuration drift %s of server '%s'", driftId, serverId)

	serverIdNumeric, err := GetServerId(serverId)
	if err != nil {
		return err
	}

	driftIdNumeric, err := utils.GetInt64FromString(driftId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	drift, httpRes, err := client.ServerAPI.
		AcknowledgeServerDrift(ctx, serverIdNumeric, driftIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(drift, &serverDriftPrintConfig)
}

// ServerSnapshotList lists the configuration snapshots stored for a server.
// kind optionally narrows the listing to one snapshot class.
func ServerSnapshotList(ctx context.Context, serverId string, kind string) error {
	logger.Get().Info().Msgf("Listing snapshots of server '%s'", serverId)

	serverIdNumeric, err := GetServerId(serverId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.ServerAPI.GetServerSnapshots(ctx, serverIdNumeric)
	if kind != "" {
		request = request.Kind(kind)
	}

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &serverSnapshotPrintConfig)
}

// ServerSyncTargetSnapshot points the drift-detection target snapshot of a
// server at its latest snapshot, clearing the reported drift.
func ServerSyncTargetSnapshot(ctx context.Context, serverId string) error {
	logger.Get().Info().Msgf("Syncing target snapshot of server '%s' with the latest snapshot", serverId)

	serverIdNumeric, err := GetServerId(serverId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.ServerAPI.
		ServerSyncTargetSnapshotWithLatestSnapshot(ctx, serverIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Target snapshot of server '%s' synced with the latest snapshot", serverId)
	return nil
}
