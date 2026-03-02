package server_instance

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

// serverInstanceRaw works around the SDK bug where Links is typed as
// map[string]interface{} but the API may return an array.
type serverInstanceRaw struct {
	Id               float32     `json:"id"`
	Label            string      `json:"label"`
	InfrastructureId float32     `json:"infrastructureId"`
	GroupId          float32     `json:"groupId"`
	ServiceStatus    string      `json:"serviceStatus"`
	CreatedTimestamp string      `json:"createdTimestamp"`
	UpdatedTimestamp string      `json:"updatedTimestamp"`
	Links            interface{} `json:"links,omitempty"`
}

type serverInstanceListRaw struct {
	Data []serverInstanceRaw `json:"data"`
}

var serverInstancePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"Label": {
			MaxWidth: 30,
			Order:    2,
		},
		"InfrastructureId": {
			Title: "Infra ID",
			Order: 3,
		},
		"GroupId": {
			Title: "Group ID",
			Order: 4,
		},
		"ServiceStatus": {
			Title:       "Status",
			Transformer: formatter.FormatStatusValue,
			Order:       5,
		},
		"CreatedTimestamp": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       6,
		},
		"UpdatedTimestamp": {
			Title:       "Updated",
			Transformer: formatter.FormatDateTimeValue,
			Order:       7,
		},
	},
}

func ServerInstanceList(ctx context.Context, infraId string) error {
	logger.Get().Info().Msgf("Listing server instances for infrastructure '%s'", infraId)

	id, err := strconv.ParseInt(infraId, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid infrastructure ID '%s': %w", infraId, err)
	}

	client := api.GetApiClient(ctx)

	_, httpRes, sdkErr := client.ServerInstanceAPI.GetInfrastructureServerInstances(ctx, int32(id)).Execute()

	if httpRes != nil && httpRes.StatusCode >= 400 {
		if err := response_inspector.InspectResponse(httpRes, sdkErr); err != nil {
			return err
		}
	} else if httpRes == nil {
		return sdkErr
	}

	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var raw serverInstanceListRaw
	if err := json.Unmarshal(body, &raw); err != nil {
		return fmt.Errorf("failed to parse server instances: %w", err)
	}

	return formatter.PrintResult(raw.Data, &serverInstancePrintConfig)
}

func ServerInstanceGet(ctx context.Context, serverInstanceId string) error {
	logger.Get().Info().Msgf("Get server instance details for %s", serverInstanceId)

	serverInstanceIdNumerical, err := utils.GetFloat32FromString(serverInstanceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	serverInstanceInfo, httpRes, err := client.ServerInstanceAPI.GetServerInstance(ctx, int32(serverInstanceIdNumerical)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(serverInstanceInfo, &serverInstancePrintConfig)
}

func ServerInstancePower(ctx context.Context, serverInstanceId string, action string) error {
	logger.Get().Info().Msgf("Setting power for server instance '%s' to '%s'", serverInstanceId, action)

	validActions := map[string]bool{
		"on":    true,
		"off":   true,
		"reset": true,
		"soft":  true,
	}

	if !validActions[action] {
		return fmt.Errorf("invalid power action: '%s'. Valid actions are: on, off, reset, soft", action)
	}

	instanceId, revision, err := getServerInstanceIdAndRevision(ctx, serverInstanceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	powerSet := sdk.ServerInstancePowerSet{
		PowerCommand: action,
	}

	httpRes, err := client.ServerInstanceAPI.
		SetPowerToServerInstance(ctx, instanceId).
		ServerInstancePowerSet(powerSet).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Power command '%s' sent to server instance '%s'", action, serverInstanceId)
	fmt.Printf("Power command '%s' sent to server instance %s\n", action, serverInstanceId)

	return nil
}

func ServerInstancePowerStatus(ctx context.Context, serverInstanceId string) error {
	logger.Get().Info().Msgf("Getting power status for server instance '%s'", serverInstanceId)

	instanceId, err := getServerInstanceId(serverInstanceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.ServerInstanceAPI.
		GetPowerFromServerInstance(ctx, instanceId).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	if httpRes != nil && httpRes.Body != nil {
		defer httpRes.Body.Close()
		body, err := io.ReadAll(httpRes.Body)
		if err != nil {
			return fmt.Errorf("failed to read power status response: %w", err)
		}

		status := strings.Trim(strings.TrimSpace(string(body)), "\"")
		if status != "" {
			fmt.Printf("Power status for server instance %s: %s\n", serverInstanceId, status)
		} else {
			fmt.Printf("Power status check initiated for server instance %s\n", serverInstanceId)
		}
	}

	return nil
}

func ServerInstanceCredentials(ctx context.Context, serverInstanceId string) error {
	logger.Get().Info().Msgf("Getting credentials for server instance '%s'", serverInstanceId)

	instanceId, err := getServerInstanceId(serverInstanceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	credentials, httpRes, err := client.ServerInstanceAPI.
		GetServerInstanceCredentials(ctx, instanceId).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(credentials, &formatter.PrintConfig{
		FieldsConfig: map[string]formatter.RecordFieldConfig{
			"Username": {
				Title: "Username",
				Order: 1,
			},
			"InitialPassword": {
				Title: "Password",
				Order: 2,
			},
			"PublicSshKey": {
				Title:    "SSH Public Key",
				MaxWidth: 60,
				Order:    3,
			},
		},
	})
}

func ServerInstanceReinstallOS(ctx context.Context, serverInstanceId string) error {
	logger.Get().Info().Msgf("Reinstalling OS for server instance '%s'", serverInstanceId)

	instanceId, revision, err := getServerInstanceIdAndRevision(ctx, serverInstanceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	reinstall := sdk.ServerInstanceReinstallOS{
		PerformAtNextDeploy: true,
		ReinstallOS:         true,
	}

	httpRes, err := client.ServerInstanceAPI.
		ReinstallServerInstanceOS(ctx, instanceId).
		ServerInstanceReinstallOS(reinstall).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("OS reinstall scheduled for server instance '%s'", serverInstanceId)
	fmt.Printf("OS reinstall scheduled for server instance %s (will take effect at next deploy)\n", serverInstanceId)

	return nil
}

func ServerInstanceConfig(ctx context.Context, serverInstanceId string) error {
	logger.Get().Info().Msgf("Getting configuration for server instance '%s'", serverInstanceId)

	instanceId, err := getServerInstanceId(serverInstanceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	config, httpRes, err := client.ServerInstanceAPI.
		GetServerInstanceConfig(ctx, instanceId).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(config, &formatter.PrintConfig{
		FieldsConfig: map[string]formatter.RecordFieldConfig{
			"Label": {
				Title: "Label",
				Order: 1,
			},
			"GroupId": {
				Title: "Group ID",
				Order: 2,
			},
			"ServerTypeId": {
				Title: "Server Type",
				Order: 3,
			},
			"ServerId": {
				Title: "Server ID",
				Order: 4,
			},
			"OsTemplateId": {
				Title: "OS Template",
				Order: 5,
			},
			"Hostname": {
				Title: "Hostname",
				Order: 6,
			},
			"DeployType": {
				Title: "Deploy Type",
				Order: 7,
			},
			"DeployStatus": {
				Title:       "Deploy Status",
				Transformer: formatter.FormatStatusValue,
				Order:       8,
			},
		},
	})
}

func getServerInstanceId(serverInstanceId string) (int32, error) {
	id, err := utils.GetFloat32FromString(serverInstanceId)
	if err != nil {
		return 0, err
	}
	return int32(id), nil
}

func getServerInstanceIdAndRevision(ctx context.Context, serverInstanceId string) (int32, string, error) {
	instanceId, err := getServerInstanceId(serverInstanceId)
	if err != nil {
		return 0, "", err
	}

	client := api.GetApiClient(ctx)

	instance, httpRes, err := client.ServerInstanceAPI.GetServerInstance(ctx, instanceId).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return 0, "", err
	}

	return instanceId, strconv.Itoa(int(instance.Revision)), nil
}
