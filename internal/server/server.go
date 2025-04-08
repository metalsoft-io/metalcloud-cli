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

	serverIdNumeric, err := getServerId(serverId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	serverInfo, httpRes, err := client.ServerAPI.GetServerInfo(ctx, serverIdNumeric).Execute()
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
	logger.Get().Info().Msgf("Registering server")

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

	return formatter.PrintResult(response, nil)
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

	fmt.Printf("Server '%s' deleted.\n", serverId)
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

	fmt.Printf("Server '%s' archived.\n", serverId)
	return nil
}

func ServerPower(ctx context.Context, serverId string, action string) error {
	logger.Get().Info().Msgf("Setting power state for server '%s' to '%s'", serverId, action)

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

	fmt.Printf("Power state for server '%s' set to '%s'.\n", serverId, action)
	return nil
}

func ServerPowerStatus(ctx context.Context, serverId string) error {
	logger.Get().Info().Msgf("Getting power state for server '%s'", serverId)

	serverIdNumeric, err := getServerId(serverId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	powerStatus, httpRes, err := client.ServerAPI.GetServerPowerStatus(ctx, serverIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	fmt.Printf("Power state for server '%s' is '%s'.\n", serverId, powerStatus)
	return formatter.PrintResult(powerStatus, nil)
}

func ServerUpdate(ctx context.Context, serverId string, config []byte) error {
	logger.Get().Info().Msgf("Updating server '%s'", serverId)

	var updateConfig sdk.UpdateServer
	err := json.Unmarshal(config, &updateConfig)
	if err != nil {
		return err
	}

	serverIdNumeric, revision, err := getServerIdAndRevision(ctx, serverId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	serverInfo, httpRes, err := client.ServerAPI.UpdateServer(ctx, serverIdNumeric).UpdateServer(updateConfig).IfMatch(revision).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
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

	fmt.Printf("IPMI credentials for server '%s' updated.\n", serverId)
	return formatter.PrintResult(serverCredentials, nil)
}

func ServerEnableSnmp(ctx context.Context, serverId string) error {
	logger.Get().Info().Msgf("Enabling SNMP for server '%s'", serverId)

	serverIdNumeric, revision, err := getServerIdAndRevision(ctx, serverId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	result, httpRes, err := client.ServerAPI.UpdateServerEnableSnmp(ctx, serverIdNumeric).IfMatch(revision).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	fmt.Printf("SNMP enabled for server '%s'. Result: %v\n", serverId, result)
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

	fmt.Printf("Syslog enabled for server '%s'.\n", serverId)
	return nil
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

	fmt.Printf("Factory reset initiated for server '%s'.\n", serverId)
	return nil
}

func ServerVncInfo(ctx context.Context, serverId string) error {
	logger.Get().Info().Msgf("Getting VNC info for server '%s'", serverId)

	serverIdNumeric, err := getServerId(serverId)
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

	serverIdNumeric, err := getServerId(serverId)
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

	serverIdNumeric, err := getServerId(serverId)
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

func getServerId(serverId string) (float32, error) {
	serverIdNumeric, err := strconv.ParseFloat(serverId, 32)
	if err != nil {
		err := fmt.Errorf("invalid server ID: '%s'", serverId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return float32(serverIdNumeric), nil
}

func getServerIdAndRevision(ctx context.Context, serverId string) (float32, string, error) {
	serverIdNumeric, err := getServerId(serverId)
	if err != nil {
		return 0, "", err
	}

	client := api.GetApiClient(ctx)

	server, httpRes, err := client.ServerAPI.GetServerInfo(ctx, float32(serverIdNumeric)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return 0, "", err
	}

	return float32(serverIdNumeric), strconv.Itoa(int(server.Revision)), nil
}
