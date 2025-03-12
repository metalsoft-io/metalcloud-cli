package server

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var serverPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"ServerId": {
			Title: "#",
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
					Order: 10,
				},
				"Password": {
					Title: "Password",
					Order: 11,
				},
			},
		},
	},
}

type serversWithCredentials struct {
	ServerInfo        sdk.Server
	ServerCredentials sdk.ServerCredentials
}

func ServerList(ctx context.Context, showCredentials bool, filterStatus string, filterType string) error {
	logger.Get().Info().Msgf("Listing all servers")

	client := api.GetApiClient(ctx)

	request := client.ServerAPI.GetServers(ctx)

	if filterStatus != "" {
		request = request.FilterServerStatus(strings.Split(filterStatus, ","))
	}

	if filterType != "" {
		request = request.FilterServerTypeId(strings.Split(filterType, ","))
	}

	serverList, httpRes, err := request.Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	if showCredentials {
		data := make([]serversWithCredentials, 0, len(serverList.Data))

		for _, server := range serverList.Data {
			serverCredentials, httpRes, err := client.ServerAPI.GetServerCredentials(ctx, server.ServerId).Execute()
			if err := response_inspector.InspectResponse(httpRes, err); err != nil {
				return err
			}

			data = append(data, serversWithCredentials{
				ServerInfo:        server,
				ServerCredentials: *serverCredentials,
			})
		}

		return formatter.PrintResult(data, &serverWithCredentialsPrintConfig)
	}

	return formatter.PrintResult(serverList, &serverPrintConfig)
}

func ServerGet(ctx context.Context, serverId string, showCredentials bool) error {
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

	if showCredentials {
		serverCredentials, httpRes, err := client.ServerAPI.GetServerCredentials(ctx, serverInfo.ServerId).Execute()
		if err := response_inspector.InspectResponse(httpRes, err); err != nil {
			return err
		}

		data := serversWithCredentials{
			ServerInfo:        *serverInfo,
			ServerCredentials: *serverCredentials,
		}

		return formatter.PrintResult(data, &serverWithCredentialsPrintConfig)
	}

	return formatter.PrintResult(serverInfo, &serverPrintConfig)
}

func ServerRegister(ctx context.Context, config []byte) error {
	fmt.Printf("Registering server: %s\n", string(config))

	var serverConfig sdk.RegisterServer

	err := json.Unmarshal(config, &serverConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	registrationInfo, httpRes, err := client.ServerAPI.RegisterServer(ctx).RegisterServer(serverConfig).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(registrationInfo, nil)
}
