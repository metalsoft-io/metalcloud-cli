package infrastructure

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

var infrastructurePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"Label": {
			MaxWidth: 30,
			Order:    2,
		},
		"Config": {
			Hidden: true,
			InnerFields: map[string]formatter.RecordFieldConfig{
				"Label": {
					Title:    "Config Label",
					MaxWidth: 30,
					Order:    3,
				},
				"DeployStatus": {
					Title:       "Deploy Status",
					Transformer: formatter.FormatStatusValue,
					Order:       9,
				},
				"InfrastructureDeployId": {
					Title: "Deploy ID",
					Order: 10,
				},
			},
		},
		"ServiceStatus": {
			Title:       "Status",
			Transformer: formatter.FormatStatusValue,
			Order:       4,
		},
		"UserIdOwner": {
			Title: "Owner",
			Order: 5,
		},
		"SiteId": {
			Title: "Site",
			Order: 6,
		},
		"CreatedTimestamp": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       7,
		},
		"UpdatedTimestamp": {
			Title:       "Updated",
			Transformer: formatter.FormatDateTimeValue,
			Order:       8,
		},
	},
}

func InfrastructureList(ctx context.Context, showAll bool, showOrdered bool, showDeleted bool) error {
	logger.Get().Info().Msgf("Listing all infrastructures")

	client := api.GetApiClient(ctx)

	request := client.InfrastructureAPI.GetInfrastructures(ctx)

	if !showAll {
		userId := api.GetUserId(ctx)
		request = request.FilterUserIdOwner([]string{"$eq:" + userId})
	}

	statusFilters := []string{}
	if !showOrdered {
		statusFilters = append(statusFilters, "$not:$eq:ordered")
	}
	if !showDeleted {
		statusFilters = append(statusFilters, "$not:$eq:deleted")
	}

	if len(statusFilters) > 0 {
		request = request.FilterServiceStatus(statusFilters)
	}

	request = request.SortBy([]string{"id:ASC"})

	infrastructureList, httpRes, err := request.Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(infrastructureList, &infrastructurePrintConfig)
}

func InfrastructureGet(ctx context.Context, infrastructureIdOrLabel string) error {
	logger.Get().Info().Msgf("Get infrastructure '%s'", infrastructureIdOrLabel)

	infrastructureInfo, err := GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	return formatter.PrintResult(infrastructureInfo, &infrastructurePrintConfig)
}

func InfrastructureCreate(ctx context.Context, siteId string, infrastructureLabel string) error {
	logger.Get().Info().Msgf("Create infrastructure '%s'", infrastructureLabel)

	siteIdNumber, err := strconv.ParseFloat(siteId, 32)
	if err != nil {
		err := fmt.Errorf("invalid site ID: '%s'", siteId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	createInfrastructure := sdk.InfrastructureCreate{
		Label:  sdk.PtrString(infrastructureLabel),
		SiteId: float32(siteIdNumber),
		Meta:   nil,
	}

	client := api.GetApiClient(ctx)

	infrastructureInfo, httpRes, err := client.InfrastructureAPI.CreateInfrastructure(ctx).InfrastructureCreate(createInfrastructure).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(infrastructureInfo, &infrastructurePrintConfig)
}

func InfrastructureUpdate(ctx context.Context, infrastructureIdOrLabel string, label string, customVariables string) error {
	logger.Get().Info().Msgf("Update infrastructure '%s'", infrastructureIdOrLabel)

	infrastructureInfo, err := GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	updateInfrastructure := sdk.UpdateInfrastructure{}

	if label != "" {
		updateInfrastructure.Label = &label
	} else {
		updateInfrastructure.Label = &infrastructureInfo.Label
	}

	if customVariables != "" {
		err = json.Unmarshal([]byte(customVariables), &updateInfrastructure.CustomVariables)
		if err != nil {
			logger.Get().Error().Err(err).Msg("")
			return err
		}
	}

	client := api.GetApiClient(ctx)

	infrastructureInfo, httpRes, err := client.InfrastructureAPI.UpdateInfrastructureConfiguration(ctx, infrastructureInfo.Id).
		UpdateInfrastructure(updateInfrastructure).
		IfMatch(strconv.Itoa(int(*infrastructureInfo.Config.Revision))).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(infrastructureInfo, &infrastructurePrintConfig)
}

func InfrastructureDelete(ctx context.Context, infrastructureIdOrLabel string) error {
	logger.Get().Info().Msgf("Delete infrastructure '%s'", infrastructureIdOrLabel)

	infrastructureInfo, err := GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.InfrastructureAPI.
		DeleteInfrastructure(ctx, infrastructureInfo.Id).
		IfMatch(strconv.Itoa(int(infrastructureInfo.Revision))).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return nil
}

func InfrastructureDeploy(ctx context.Context, infrastructureIdOrLabel string, allowDataLoss bool, attemptSoftShutdown bool, attemptHardShutdown bool, softShutdownTimeout int, forceShutdown bool) error {
	logger.Get().Info().Msgf("Deploy infrastructure '%s'", infrastructureIdOrLabel)

	infrastructureInfo, err := GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	infrastructureDeployOptions := sdk.InfrastructureDeployOptions{
		AllowDataLoss: allowDataLoss,
		ShutdownOptions: sdk.InfrastructureDeployShutdownOptions{
			AttemptSoftShutdown: attemptSoftShutdown,
			AttemptHardShutdown: attemptHardShutdown,
			SoftShutdownTimeout: float32(softShutdownTimeout),
			ForceShutdown:       forceShutdown,
		},
		ServerTypeIdToPreferredServerIds: map[string]interface{}{},
	}

	client := api.GetApiClient(ctx)

	infrastructureInfo, httpRes, err := client.InfrastructureAPI.
		DeployInfrastructure(ctx, infrastructureInfo.Id).
		InfrastructureDeployOptions(infrastructureDeployOptions).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(infrastructureInfo, &infrastructurePrintConfig)
}

func InfrastructureRevert(ctx context.Context, infrastructureIdOrLabel string) error {
	logger.Get().Info().Msgf("Revert infrastructure '%s' changes", infrastructureIdOrLabel)

	infrastructureInfo, err := GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.InfrastructureAPI.
		RevertInfrastructure(ctx, infrastructureInfo.Id).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return nil
}

func GetInfrastructureByIdOrLabel(ctx context.Context, infrastructureIdOrLabel string) (*sdk.Infrastructure, error) {
	client := api.GetApiClient(ctx)

	infrastructureList, httpRes, err := client.InfrastructureAPI.GetInfrastructures(ctx).Search(infrastructureIdOrLabel).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, err
	}

	if len(infrastructureList.Data) == 0 {
		err := fmt.Errorf("infrastructure '%s' not found", infrastructureIdOrLabel)
		logger.Get().Error().Err(err).Msg("")
		return nil, err
	}

	var infrastructureInfo sdk.Infrastructure
	for _, infrastructure := range infrastructureList.Data {
		if infrastructure.Label == infrastructureIdOrLabel {
			infrastructureInfo = infrastructure
			break
		}

		if strconv.Itoa(int(infrastructure.Id)) == infrastructureIdOrLabel {
			infrastructureInfo = infrastructure
			break
		}
	}

	if infrastructureInfo.Id == 0 {
		err := fmt.Errorf("infrastructure '%s' not found", infrastructureIdOrLabel)
		logger.Get().Error().Err(err).Msg("")
		return nil, err
	}

	return &infrastructureInfo, nil
}

func InfrastructureGetUsers(ctx context.Context, infrastructureIdOrLabel string) error {
	logger.Get().Info().Msgf("Getting users for infrastructure '%s'", infrastructureIdOrLabel)

	infrastructureInfo, err := GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	usersPaginated, httpRes, err := client.InfrastructureAPI.
		GetInfrastructureUsers(ctx, infrastructureInfo.Id).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(usersPaginated, &formatter.PrintConfig{
		FieldsConfig: map[string]formatter.RecordFieldConfig{
			"Id": {
				Title: "#",
				Order: 1,
			},
			"DisplayName": {
				Title: "Name",
				Order: 2,
			},
			"Email": {
				Title: "E-mail",
				Order: 3,
			},
			"AccessLevel": {
				Title: "Access",
				Order: 4,
			},
		},
	})
}

func InfrastructureAddUser(ctx context.Context, infrastructureIdOrLabel string, userEmail string, createMissing string) error {
	logger.Get().Info().Msgf("Adding user to infrastructure '%s'", infrastructureIdOrLabel)

	infrastructureInfo, err := GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	addUserConfig := sdk.AddUserToInfrastructure{
		UserEmail:         userEmail,
		CreateIfNotExists: strings.ToLower(createMissing) == "true",
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.InfrastructureAPI.
		AddInfrastructureUser(ctx, infrastructureInfo.Id).
		AddUserToInfrastructure(addUserConfig).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("User added to infrastructure '%s'", infrastructureIdOrLabel)
	return nil
}

func InfrastructureRemoveUser(ctx context.Context, infrastructureIdOrLabel string, userId string) error {
	logger.Get().Info().Msgf("Removing user '%s' from infrastructure '%s'", userId, infrastructureIdOrLabel)

	infrastructureInfo, err := GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	userIdNumber, err := strconv.ParseFloat(userId, 32)
	if err != nil {
		err := fmt.Errorf("invalid user ID: '%s'", userId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.InfrastructureAPI.
		RemoveInfrastructureUser(ctx, infrastructureInfo.Id, float32(userIdNumber)).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("User '%s' removed from infrastructure '%s'", userId, infrastructureIdOrLabel)
	return nil
}

func InfrastructureGetUserLimits(ctx context.Context, infrastructureIdOrLabel string) error {
	logger.Get().Info().Msgf("Getting user limits for infrastructure '%s'", infrastructureIdOrLabel)

	infrastructureInfo, err := GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	userLimits, httpRes, err := client.InfrastructureAPI.
		GetInfrastructureUserLimits(ctx, infrastructureInfo.Id).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(userLimits, &formatter.PrintConfig{
		FieldsConfig: map[string]formatter.RecordFieldConfig{
			"ComputeNodesInstancesToProvisionLimit": {
				Title: "Compute Nodes Limit",
				Order: 1,
			},
			"DrivesAttachedToInstancesLimit": {
				Title: "Drives Limit",
				Order: 2,
			},
			"InfrastructuresLimit": {
				Title: "Infrastructures Limit",
				Order: 3,
			},
		},
	})
}

func InfrastructureGetStatistics(ctx context.Context, infrastructureIdOrLabel string) error {
	logger.Get().Info().Msgf("Getting statistics for infrastructure '%s'", infrastructureIdOrLabel)

	infrastructureInfo, err := GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	statistics, httpRes, err := client.InfrastructureAPI.
		GetInfrastructureStatistics(ctx, infrastructureInfo.Id).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(statistics, &formatter.PrintConfig{
		FieldsConfig: map[string]formatter.RecordFieldConfig{
			"GroupId": {
				Title: "Group Id",
				Order: 1,
			},
			"GroupCreatedTimestamp": {
				Title:       "Created",
				Transformer: formatter.FormatDateTimeValue,
				Order:       2,
			},
			"GroupCompletedTimestamp": {
				Title:       "Completed",
				Transformer: formatter.FormatDateTimeValue,
				Order:       3,
			},
			"JobsThrownError": {
				Title: "Errors",
				Order: 4,
			},
			"JobsCompleted": {
				Title: "Completed",
				Order: 5,
			},
		},
	})
}

func InfrastructureGetConfigInfo(ctx context.Context, infrastructureIdOrLabel string) error {
	logger.Get().Info().Msgf("Getting configuration info for infrastructure '%s'", infrastructureIdOrLabel)

	infrastructureInfo, err := GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	configInfo, httpRes, err := client.InfrastructureAPI.
		GetInfrastructureConfigInfo(ctx, infrastructureInfo.Id).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(configInfo, &formatter.PrintConfig{
		FieldsConfig: map[string]formatter.RecordFieldConfig{
			"Label": {
				Title: "Label",
				Order: 1,
			},
			"InfrastructureDeployId": {
				Title: "Deploy ID",
				Order: 2,
			},
			"DeployType": {
				Title: "Deploy Type",
				Order: 3,
			},
			"DeployStatus": {
				Title:       "Deploy Status",
				Transformer: formatter.FormatStatusValue,
				Order:       4,
			},
			"UpdatedTimestamp": {
				Title:       "Updated",
				Transformer: formatter.FormatDateTimeValue,
				Order:       5,
			},
		},
	})
}

func InfrastructureGetAllStatistics(ctx context.Context) error {
	logger.Get().Info().Msgf("Getting statistics for all infrastructures")

	client := api.GetApiClient(ctx)

	statistics, httpRes, err := client.InfrastructureAPI.
		GetAllInfrastructureStatistics(ctx).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(statistics, &formatter.PrintConfig{
		FieldsConfig: map[string]formatter.RecordFieldConfig{
			"InfrastructureCount": {
				Title: "Total Infrastructures",
				Order: 1,
			},
			"InfrastructureServiceStatus": {
				Hidden: true,
				InnerFields: map[string]formatter.RecordFieldConfig{
					"Active": {
						Title: "Active",
						Order: 2,
					},
					"Ordered": {
						Title: "Ordered",
						Order: 3,
					},
				},
			},
			"OngoingInfrastructureCount": {
				Title: "Ongoing Deployment",
				Order: 4,
			},
			"InfrastructureDeployOngoingStatusCount": {
				Hidden: true,
				InnerFields: map[string]formatter.RecordFieldConfig{
					"ThrownError": {
						Title: "Errored",
						Order: 5,
					},
					"ThrownErrorRetrying": {
						Title: "Retrying",
						Order: 6,
					},
				},
			},
		},
	})
}
