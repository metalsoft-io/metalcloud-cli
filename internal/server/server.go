package server

import (
	"context"
	"fmt"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
)

var serverPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"ServerId": {
			Title: "#",
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
		"ManagementAddress": {
			Title: "IP",
			Order: 4,
		},
		"Vendor": {
			Order: 5,
		},
		"Model": {
			Order: 6,
		},
		"ServerStatus": {
			Title:       "Status",
			Transformer: formatter.FormatStatusValue,
			Order:       7,
		},
	},
}

func ServerList(ctx context.Context) error {
	logger.Get().Info().Msgf("Listing all servers")

	client := api.GetApiClient(ctx)

	serverList, httpRes, err := client.ServerAPI.GetServers(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(serverList, &serverPrintConfig)
}

func ServerGet(ctx context.Context, serverId string) error {
	logger.Get().Info().Msgf("Get server '%s'", serverId)

	serverIdNumber, err := strconv.ParseFloat(serverId, 32)
	if err != nil {
		err := fmt.Errorf("invalid server ID: '%s'", serverId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	client := api.GetApiClient(ctx)

	serverInfo, httpRes, err := client.ServerAPI.GetServerInfo(ctx, float32(serverIdNumber)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(serverInfo, &serverPrintConfig)
}
