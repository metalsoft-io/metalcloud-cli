package user

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var userPrintConfig = formatter.PrintConfig{
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
		"CreatedTimestamp": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       4,
		},
		"LastLoginTimestamp": {
			Title:       "Last Login",
			Transformer: formatter.FormatDateTimeValue,
			Order:       5,
		},
		"AccessLevel": {
			Title: "Access",
			Order: 6,
		},
	},
}

var userLimitsPrintConfig = formatter.PrintConfig{
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
}

var userApiKeyPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"ApiKey": {
			Title: "API Key",
			Order: 1,
		},
	},
}

var userSshKeysPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"Name": {
			Title: "Name",
			Order: 2,
		},
		"Fingerprint": {
			Title: "Fingerprint",
			Order: 3,
		},
		"CreatedTimestamp": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       4,
		},
	},
}

var userPermissionsPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"UserId": {
			Title: "User ID",
			Order: 1,
		},
		"ResourceType": {
			Title: "Resource Type",
			Order: 2,
		},
		"ResourceId": {
			Title: "Resource ID",
			Order: 3,
		},
		"PermissionLevel": {
			Title: "Permission Level",
			Order: 4,
		},
	},
}

func List(ctx context.Context) error {
	logger.Get().Info().Msgf("Listing all users")

	client := api.GetApiClient(ctx)

	userList, httpRes, err := client.UserAPI.GetUsers(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(userList, &userPrintConfig)
}

func Get(ctx context.Context, userId string) error {
	logger.Get().Info().Msgf("Get user '%s'", userId)

	userIdNumber, err := getUserId(userId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	userInfo, httpRes, err := client.UserAPI.GetUser(ctx, userIdNumber).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(userInfo, &userPrintConfig)
}

func GetLimits(ctx context.Context, userId string) error {
	logger.Get().Info().Msgf("Get user '%s' limits", userId)

	userIdNumber, err := getUserId(userId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	userLimits, httpRes, err := client.UserAPI.GetUserLimits(ctx, userIdNumber).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(userLimits, &userLimitsPrintConfig)
}

func Create(ctx context.Context, config []byte) error {
	logger.Get().Info().Msgf("Creating user")

	var userConfig sdk.CreateUser
	err := json.Unmarshal(config, &userConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	userInfo, httpRes, err := client.UserAPI.CreateUserAuthorized(ctx).CreateUser(userConfig).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(userInfo, &userPrintConfig)
}

func Archive(ctx context.Context, userId string) error {
	logger.Get().Info().Msgf("Archiving user '%s'", userId)

	userIdNumber, revision, err := getUserIdAndRevision(ctx, userId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	userInfo, httpRes, err := client.UserAPI.ArchiveUser(ctx, userIdNumber).IfMatch(revision).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("User '%s' archived", userId)
	return formatter.PrintResult(userInfo, &userPrintConfig)
}

func Unarchive(ctx context.Context, userId string) error {
	logger.Get().Info().Msgf("Unarchiving user '%s'", userId)

	userIdNumber, revision, err := getUserIdAndRevision(ctx, userId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	userInfo, httpRes, err := client.UserAPI.UnarchiveUser(ctx, userIdNumber).IfMatch(revision).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("User '%s' unarchived", userId)
	return formatter.PrintResult(userInfo, &userPrintConfig)
}

func UpdateLimits(ctx context.Context, userId string, config []byte) error {
	logger.Get().Info().Msgf("Updating limits for user '%s'", userId)

	userIdNumber, revision, err := getUserIdAndRevision(ctx, userId)
	if err != nil {
		return err
	}

	var userLimitsConfig sdk.UserLimits
	err = json.Unmarshal(config, &userLimitsConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	userLimits, httpRes, err := client.UserAPI.UpdateUserLimits(ctx, userIdNumber).UserLimits(userLimitsConfig).IfMatch(revision).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Limits updated for user '%s'", userId)
	return formatter.PrintResult(userLimits, &userLimitsPrintConfig)
}

func UpdateConfig(ctx context.Context, userId string, config []byte) error {
	logger.Get().Info().Msgf("Updating configuration for user '%s'", userId)

	userIdNumber, revision, err := getUserIdAndRevision(ctx, userId)
	if err != nil {
		return err
	}

	var userConfig sdk.UpdateUser
	err = json.Unmarshal(config, &userConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	userConfiguration, httpRes, err := client.UserAPI.UpdateUserConfig(ctx, userIdNumber).UpdateUser(userConfig).IfMatch(revision).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Configuration updated for user '%s'", userId)
	return formatter.PrintResult(userConfiguration, nil)
}

func ChangeAccount(ctx context.Context, userId string, accountId float32) error {
	logger.Get().Info().Msgf("Changing account for user '%s' to account '%g'", userId, accountId)

	userIdNumber, revision, err := getUserIdAndRevision(ctx, userId)
	if err != nil {
		return err
	}

	changeAccount := sdk.ChangeUserAccount{
		NewAccountId: accountId,
	}

	client := api.GetApiClient(ctx)

	userInfo, httpRes, err := client.UserAPI.ChangeUserAccount(ctx, userIdNumber).ChangeUserAccount(changeAccount).IfMatch(revision).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Account changed for user '%s' to account '%g'", userId, accountId)
	return formatter.PrintResult(userInfo, &userPrintConfig)
}

func GetSSHKeys(ctx context.Context, userId string) error {
	logger.Get().Info().Msgf("Getting SSH keys for user '%s'", userId)

	userIdNumber, err := getUserId(userId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	sshKeys, httpRes, err := client.UserAPI.GetUserSshKeys(ctx, userIdNumber).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(sshKeys, &userSshKeysPrintConfig)
}

func AddSSHKey(ctx context.Context, userId string, keyContent string) error {
	logger.Get().Info().Msgf("Adding SSH key for user '%s'", userId)

	userIdNumber, err := getUserId(userId)
	if err != nil {
		return err
	}

	sshKeyData := sdk.CreateUserSSHKeyDto{
		SshKey: keyContent,
	}

	client := api.GetApiClient(ctx)

	sshKey, httpRes, err := client.UserAPI.AddUserSshKey(ctx, userIdNumber).CreateUserSSHKeyDto(sshKeyData).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("SSH key added for user '%s'", userId)
	return formatter.PrintResult(sshKey, &userSshKeysPrintConfig)
}

func DeleteSSHKey(ctx context.Context, userId string, keyId string) error {
	logger.Get().Info().Msgf("Deleting SSH key '%s' for user '%s'", keyId, userId)

	userIdNumber, err := getUserId(userId)
	if err != nil {
		return err
	}

	keyIdNumber, err := strconv.ParseFloat(keyId, 32)
	if err != nil {
		err := fmt.Errorf("invalid SSH key ID: '%s'", keyId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.UserAPI.DeleteUserSshKey(ctx, userIdNumber, float32(keyIdNumber)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("SSH key '%s' deleted for user '%s'", keyId, userId)
	return nil
}

func GetAPIKey(ctx context.Context, userId string) error {
	logger.Get().Info().Msgf("Getting API key for user '%s'", userId)

	userIdNumber, err := getUserId(userId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	apiKey, httpRes, err := client.UserAPI.GetUserApiKey(ctx, userIdNumber).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(apiKey, &userApiKeyPrintConfig)
}

func RegenerateAPIKey(ctx context.Context, userId string) error {
	logger.Get().Info().Msgf("Regenerating API key for user '%s'", userId)

	userIdNumber, revision, err := getUserIdAndRevision(ctx, userId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	_, httpRes, err := client.UserAPI.RegenerateUserApiKey(ctx, userIdNumber).IfMatch(revision).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("API key regenerated for user '%s'", userId)

	// Get the new API key to display it
	apiKey, httpRes, err := client.UserAPI.GetUserApiKey(ctx, userIdNumber).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(apiKey, &userApiKeyPrintConfig)
}

func Suspend(ctx context.Context, userId string, reason string) error {
	logger.Get().Info().Msgf("Suspending user '%s'", userId)

	userIdNumber, revision, err := getUserIdAndRevision(ctx, userId)
	if err != nil {
		return err
	}

	suspendReason := sdk.UserSuspend{
		SuspendReason: reason,
	}

	client := api.GetApiClient(ctx)

	suspendInfo, httpRes, err := client.UserAPI.SuspendUser(ctx, userIdNumber).UserSuspend(suspendReason).IfMatch(revision).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("User '%s' suspended", userId)
	return formatter.PrintResult(suspendInfo, nil)
}

func Unsuspend(ctx context.Context, userId string) error {
	logger.Get().Info().Msgf("Unsuspending user '%s'", userId)

	userIdNumber, revision, err := getUserIdAndRevision(ctx, userId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.UserAPI.UnsuspendUser(ctx, userIdNumber).IfMatch(revision).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("User '%s' unsuspended", userId)
	return nil
}

func GetPermissions(ctx context.Context, userId string) error {
	logger.Get().Info().Msgf("Getting permissions for user '%s'", userId)

	userIdNumber, err := getUserId(userId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	permissions, httpRes, err := client.UserAPI.GetUserPermissions(ctx, userIdNumber).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(permissions, &userPermissionsPrintConfig)
}

func UpdatePermissions(ctx context.Context, userId string, config []byte) error {
	logger.Get().Info().Msgf("Updating permissions for user '%s'", userId)

	userIdNumber, revision, err := getUserIdAndRevision(ctx, userId)
	if err != nil {
		return err
	}

	var permissionsConfig sdk.UpdateUserPermissionsDto
	err = json.Unmarshal(config, &permissionsConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	permissions, httpRes, err := client.UserAPI.UpdateUserPermissions(ctx, userIdNumber).UpdateUserPermissionsDto(permissionsConfig).IfMatch(revision).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Permissions updated for user '%s'", userId)
	return formatter.PrintResult(permissions, &userPermissionsPrintConfig)
}

func getUserId(userId string) (float32, error) {
	userIdNumeric, err := strconv.ParseFloat(userId, 32)
	if err != nil {
		err := fmt.Errorf("invalid user ID: '%s'", userId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return float32(userIdNumeric), nil
}

func getUserIdAndRevision(ctx context.Context, userId string) (float32, string, error) {
	userIdNumeric, err := getUserId(userId)
	if err != nil {
		return 0, "", err
	}

	client := api.GetApiClient(ctx)

	user, httpRes, err := client.UserAPI.GetUser(ctx, float32(userIdNumeric)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return 0, "", err
	}

	return float32(userIdNumeric), strconv.Itoa(int(user.Revision)), nil
}
