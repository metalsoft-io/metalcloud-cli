package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
	"github.com/spf13/viper"
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

/* cspell:disable */

func InfrastructureGetUtilization(ctx context.Context, userId int, startTime time.Time, endTime time.Time, siteIds []int, infrastructureIds []int,
	showInstances, showDrives, showSubnets bool) error {
	logger.Get().Info().Msgf("Getting utilization report for user %d from %s to %s", userId, startTime.Format(time.RFC3339), endTime.Format(time.RFC3339))

	client := api.GetApiClient(ctx)

	request := sdk.GetResourceUtilizationDetailed{
		UserIdOwner:    float32(userId),
		StartTimestamp: startTime.Format(time.RFC3339),
		EndTimestamp:   endTime.Format(time.RFC3339),
	}

	if len(siteIds) > 0 {
		request.SiteIds = make([]float32, 0, len(siteIds))
		for _, siteId := range siteIds {
			siteIdFloat := float32(siteId)
			request.SiteIds = append(request.SiteIds, siteIdFloat)
		}
	}

	if len(infrastructureIds) > 0 {
		request.InfrastructureIds = make([]float32, 0, len(infrastructureIds))
		for _, infrastructureId := range infrastructureIds {
			infrastructureIdFloat := float32(infrastructureId)
			request.InfrastructureIds = append(request.InfrastructureIds, infrastructureIdFloat)
		}
	}

	utilization, httpRes, err := client.InfrastructureAPI.
		GetInfrastructureResourceUtilizationDetailed(ctx).
		GetResourceUtilizationDetailed(request).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	if strings.ToLower(viper.GetString(formatter.ConfigFormat)) == "json" {
		// Do not change the output format if JSON is requested
		return formatter.PrintResult(utilization, nil)
	}

	reportData := []utilReportRecord{}

	logger.Get().Debug().Msgf("Processing utilization for %d infrastructures", len(utilization.Infrastructures))

	for _, infrastructure := range utilization.Infrastructures {
		logger.Get().Debug().Msgf("Processing utilization for infrastructure %v", infrastructure)

		infrastructureId := infrastructure.InfrastructureId
		infrastructureLabel := infrastructure.InfrastructureLabel
		infrastructureServiceStatus := infrastructure.InfrastructureServiceStatus

		// Add infrastructure-level data
		reportData = append(reportData, utilReportRecord{
			InfrastructureId:            infrastructureId,
			InfrastructureLabel:         infrastructureLabel,
			InfrastructureServiceStatus: infrastructureServiceStatus,
		})

		// Get detailed report for this infrastructure
		infrastructureIdStr := fmt.Sprintf("%v", infrastructureId)
		if infraDetails, infraExists := (*utilization.DetailedReport)[infrastructureIdStr]; infraExists {
			// Process dataList for this infrastructure
			if showInstances {
				for _, item := range infraDetails.Instance {
					reportData = append(reportData, utilReportRecord{
						InfrastructureId:            infrastructureId,
						InfrastructureLabel:         infrastructureLabel,
						InfrastructureServiceStatus: infrastructureServiceStatus,
						//
						Kind:              "Instance",
						Id:                item.Id,
						Label:             item.Label,
						StartTime:         item.StartTimestamp,
						EndTime:           item.EndTimestamp,
						MeasurementPeriod: item.MeasurementPeriod,
						MeasurementUnit:   item.MeasurementUnit,
						Quantity:          item.Quantity,
						//
						ServerTypeId:               item.ServerTypeId,
						ServerId:                   item.ServerId,
						ServerTypeName:             item.ServerTypeName,
						OperatingSystemType:        item.OperatingSystemType,
						OperatingSystemVersion:     item.OperatingSystemVersion,
						OperatingSystemDisplayName: item.OperatingSystemDisplayName,
						OperatingSystemTemplateId:  item.OperatingSystemTemplateId,
						OriginalStartTimestamp:     item.OriginalStartTimestamp,
					})
				}
			}

			if showDrives {
				for _, item := range infraDetails.Drive {
					reportData = append(reportData, utilReportRecord{
						InfrastructureId:            infrastructureId,
						InfrastructureLabel:         infrastructureLabel,
						InfrastructureServiceStatus: infrastructureServiceStatus,
						//
						Kind:              "Drive",
						Id:                item.Id,
						Label:             item.Label,
						StartTime:         item.StartTimestamp,
						EndTime:           item.EndTimestamp,
						MeasurementPeriod: item.MeasurementPeriod,
						MeasurementUnit:   item.MeasurementUnit,
						Quantity:          item.Quantity,
						//
						DriveSizeMbytes:  item.DriveSizeMbytes,
						DriveStorageType: item.DriveStorageType,
					})
				}

				for _, item := range infraDetails.SharedDrive {
					reportData = append(reportData, utilReportRecord{
						InfrastructureId:            infrastructureId,
						InfrastructureLabel:         infrastructureLabel,
						InfrastructureServiceStatus: infrastructureServiceStatus,
						//
						Kind:              "Shared Drive",
						Id:                item.Id,
						Label:             item.Label,
						StartTime:         item.StartTimestamp,
						EndTime:           item.EndTimestamp,
						MeasurementPeriod: item.MeasurementPeriod,
						MeasurementUnit:   item.MeasurementUnit,
						Quantity:          item.Quantity,
						//
						DriveSizeMbytes:  item.SharedDriveSizeMbytes,
						DriveStorageType: item.SharedDriveStorageType,
					})
				}
			}

			if showSubnets {
				for _, item := range infraDetails.Subnet {
					reportData = append(reportData, utilReportRecord{
						InfrastructureId:            infrastructureId,
						InfrastructureLabel:         infrastructureLabel,
						InfrastructureServiceStatus: infrastructureServiceStatus,
						//
						Kind:              "Subnet",
						Id:                item.Id,
						Label:             item.Label,
						StartTime:         item.StartTimestamp,
						EndTime:           item.EndTimestamp,
						MeasurementPeriod: item.MeasurementPeriod,
						MeasurementUnit:   item.MeasurementUnit,
						Quantity:          item.Quantity,
						//
						SubnetIpCount:    item.SubnetIpCount,
						SubnetPrefixSize: item.SubnetPrefixSize,
						SubnetType:       item.SubnetType,
					})
				}
			}
		}
	}

	logger.Get().Debug().Msgf("Utilization report data: %v", reportData)

	return formatter.PrintResult(reportData, utilReportRecordPrintConfig(showInstances, showDrives, showSubnets))
}

type utilReportRecord = struct {
	InfrastructureId            interface{}
	InfrastructureLabel         interface{}
	InfrastructureServiceStatus interface{}
	Kind                        interface{}
	StartTime                   interface{}
	EndTime                     interface{}
	Id                          interface{}
	Label                       interface{}
	MeasurementPeriod           interface{}
	MeasurementUnit             interface{}
	Quantity                    interface{}

	//INSTANCE
	ServerTypeId               interface{}
	ServerId                   interface{}
	ServerTypeName             interface{}
	OperatingSystemType        interface{}
	OperatingSystemVersion     interface{}
	OperatingSystemDisplayName interface{}
	OperatingSystemTemplateId  interface{}
	OriginalStartTimestamp     interface{}
	//(SHARED) DRIVE
	DriveSizeMbytes  interface{}
	DriveStorageType interface{}
	//SharedDriveSizeMbytes  interface{}
	//SharedDriveStorageType interface{}
	//SUBNET
	SubnetIpCount    interface{}
	SubnetPrefixSize interface{}
	SubnetType       interface{}
}

func utilReportRecordPrintConfig(showInstances, showDrives, showSubnets bool) *formatter.PrintConfig {
	printConfig := formatter.PrintConfig{
		FieldsConfig: map[string]formatter.RecordFieldConfig{
			"InfrastructureId": {
				Title: "Infra ID",
				Order: 1,
			},
			"InfrastructureLabel": {
				Title: "Infra Label",
				Order: 2,
			},
			"InfrastructureServiceStatus": {
				Title: "Infra Status",
				Order: 3,
			},
		},
	}

	var orderIndex = len(printConfig.FieldsConfig)
	if showInstances || showDrives || showSubnets {
		orderIndex++
		printConfig.FieldsConfig["Kind"] = formatter.RecordFieldConfig{
			Title: "Kind",
			Order: orderIndex,
		}

		orderIndex++
		printConfig.FieldsConfig["Id"] = formatter.RecordFieldConfig{
			Title: "ID",
			Order: orderIndex,
		}

		orderIndex++
		printConfig.FieldsConfig["Label"] = formatter.RecordFieldConfig{
			Title: "Label",
			Order: orderIndex,
		}

		orderIndex++
		printConfig.FieldsConfig["StartTime"] = formatter.RecordFieldConfig{
			Title: "Start Time",
			Order: orderIndex,
			//Transformer: formatter.FormatDateTimeValue,
		}

		orderIndex++
		printConfig.FieldsConfig["EndTime"] = formatter.RecordFieldConfig{
			Title: "End Time",
			Order: orderIndex,
			//Transformer: formatter.FormatDateTimeValue,
		}

		orderIndex++
		printConfig.FieldsConfig["MeasurementPeriod"] = formatter.RecordFieldConfig{
			Title:       "Measurement",
			Order:       orderIndex,
			Transformer: formatter.FormatIntegerValue,
		}

		orderIndex++
		printConfig.FieldsConfig["MeasurementUnit"] = formatter.RecordFieldConfig{
			Title: "Unit",
			Order: orderIndex,
		}

		orderIndex++
		printConfig.FieldsConfig["Quantity"] = formatter.RecordFieldConfig{
			Title:       "Quantity",
			Order:       orderIndex,
			Transformer: formatter.FormatIntegerValue,
		}
	}

	if showInstances {
		//INSTANCE:
		orderIndex++
		printConfig.FieldsConfig["ServerTypeId"] = formatter.RecordFieldConfig{
			Title: "Server Type ID",
			Order: orderIndex,
		}

		orderIndex++
		printConfig.FieldsConfig["ServerId"] = formatter.RecordFieldConfig{
			Title: "Server ID",
			Order: orderIndex,
		}

		orderIndex++
		printConfig.FieldsConfig["ServerTypeName"] = formatter.RecordFieldConfig{
			Title: "Server Type Name",
			Order: orderIndex,
		}

		orderIndex++
		printConfig.FieldsConfig["OperatingSystemType"] = formatter.RecordFieldConfig{
			Title: "OS Type",
			Order: orderIndex,
		}

		orderIndex++
		printConfig.FieldsConfig["OperatingSystemVersion"] = formatter.RecordFieldConfig{
			Title: "OS Version",
			Order: orderIndex,
		}

		orderIndex++
		printConfig.FieldsConfig["OperatingSystemDisplayName"] = formatter.RecordFieldConfig{
			Title: "OS Display Name",
			Order: orderIndex,
		}

		orderIndex++
		printConfig.FieldsConfig["OperatingSystemTemplateId"] = formatter.RecordFieldConfig{
			Title: "OS Template ID",
			Order: orderIndex,
		}

		orderIndex++
		printConfig.FieldsConfig["OriginalStartTimestamp"] = formatter.RecordFieldConfig{
			Title: "Original Start TS",
			Order: orderIndex,
		}
	}

	if showDrives {
		//DRIVE:
		orderIndex++
		printConfig.FieldsConfig["DriveSizeMbytes"] = formatter.RecordFieldConfig{
			Title: "Drive Size Mbytes",
			Order: orderIndex,
		}

		orderIndex++
		printConfig.FieldsConfig["DriveStorageType"] = formatter.RecordFieldConfig{
			Title: "Drive Storage Type",
			Order: orderIndex,
		}

		//SHARED DRIVE:
		/*
			orderIndex++
			printConfig.FieldsConfig["SharedDriveSizeMbytes"] = formatter.RecordFieldConfig{
				Title: "Shared Drive Size Mbytes",
				Order: orderIndex,
			}

			orderIndex++
			printConfig.FieldsConfig["SharedDriveStorageType"] = formatter.RecordFieldConfig{
				Title: "Shared Drive Storage Type",
				Order: orderIndex,
			}
		*/
	}

	//SUBNETS:
	if showSubnets {
		orderIndex++
		printConfig.FieldsConfig["SubnetIpCount"] = formatter.RecordFieldConfig{
			Title: "Subnet Ip Count",
			Order: orderIndex,
		}
		orderIndex++
		printConfig.FieldsConfig["SubnetPrefixSize"] = formatter.RecordFieldConfig{
			Title: "Subnet Prefix Size",
			Order: orderIndex,
		}
		orderIndex++
		printConfig.FieldsConfig["SubnetType"] = formatter.RecordFieldConfig{
			Title: "Subnet Type",
			Order: orderIndex,
		}
	}

	return &printConfig
}
