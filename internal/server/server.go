package server

import (
	"context"
	"fmt"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
)

var serverPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"ServerId": {
			Title: "#",
		},
		"ServerUUID": {
			Title:    "UUID",
			IsStatus: true,
		},
		"SerialNumber": {
			Title: "S/N",
		},
		"ManagementAddress": {
			Title: "IP",
		},
		"Vendor": {},
		"Model":  {},
		"ServerStatus": {
			Title:    "Status",
			IsStatus: true,
		},
	},
}

func ServerList(ctx context.Context) error {
	logger.Get().Info().Msgf("Listing all servers")

	client, err := system.GetApiClient(ctx)
	if err != nil {
		return err
	}

	serverList, httpRes, err := client.ServerAPI.GetServers(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(serverList, &serverPrintConfig)
}

func ServerGet(ctx context.Context, serverId string) error {
	logger.Get().Info().Msgf("Get server '%s'", serverId)

	client, err := system.GetApiClient(ctx)
	if err != nil {
		return err
	}

	serverIdNumber, err := strconv.ParseFloat(serverId, 32)
	if err != nil {
		err := fmt.Errorf("invalid server ID: '%s'", serverId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	serverInfo, httpRes, err := client.ServerAPI.GetServerInfo(ctx, float32(serverIdNumber)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(serverInfo, &serverPrintConfig)
}
