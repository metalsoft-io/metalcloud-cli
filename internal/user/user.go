package user

import (
	"context"
	"fmt"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var userPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"DisplayName": {
			Title:    "Name",
			MaxWidth: 30,
			Order:    2,
		},
		"Email": {
			Title:    "E-mail",
			MaxWidth: 50,
			Order:    3,
		},
		"AccessLevel": {
			Title: "Role",
			Order: 4,
		},
		"IsArchived": {
			Title: "Archived",
			Order: 5,
		},
		"CreatedTimestamp": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       6,
		},
		"LastLoginTimestamp": {
			Title:       "Last Login",
			Transformer: formatter.FormatDateTimeValue,
			Order:       7,
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

var userSshKeysPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"SshKey": {
			Title:    "SSH Key",
			MaxWidth: 50,
			Order:    2,
		},
		"Status": {
			Title: "Status",
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

func List(ctx context.Context, archived bool, filterId, filterDisplayName, filterEmail, filterAccountId, filterInfrastructureId, sortBy, search, searchBy string) error {
	logger.Get().Info().Msgf("Listing all users")

	client := api.GetApiClient(ctx)

	request := client.UsersAPI.GetUsers(ctx)

	if !archived {
		request = request.FilterArchived([]string{"false"})
	}
	if filterId != "" {
		request = request.FilterId([]string{filterId})
	}
	if filterDisplayName != "" {
		request = request.FilterDisplayName([]string{filterDisplayName})
	}
	if filterEmail != "" {
		request = request.FilterEmail([]string{filterEmail})
	}
	if filterAccountId != "" {
		request = request.FilterAccountId([]string{filterAccountId})
	}
	if filterInfrastructureId != "" {
		request = request.FilterInfrastructureIdDefault([]string{filterInfrastructureId})
	}
	if sortBy != "" {
		request = request.SortBy([]string{sortBy})
	}
	if search != "" {
		request = request.Search(search)
	}
	if searchBy != "" {
		request = request.SearchBy([]string{searchBy})
	}

	userList, httpRes, err := request.Execute()
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

	userInfo, httpRes, err := client.UsersAPI.GetUser(ctx, userIdNumber).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(userInfo, &userPrintConfig)
}

func Create(ctx context.Context, config []byte) error {
	logger.Get().Info().Msgf("Creating user")

	var userConfig sdk.CreateUser
	err := utils.UnmarshalContent(config, &userConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	userInfo, httpRes, err := client.UsersAPI.CreateUserAuthorized(ctx).CreateUser(userConfig).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(userInfo, &userPrintConfig)
}

// CreateBulk creates multiple users from a JSON or YAML configuration
func CreateBulk(ctx context.Context, config []byte) error {
	logger.Get().Info().Msgf("Creating users in bulk")

	var usersConfig []sdk.CreateUser
	err := utils.UnmarshalContent(config, &usersConfig)
	if err != nil {
		return err
	}

	if len(usersConfig) == 0 {
		return fmt.Errorf("no users found in configuration")
	}

	client := api.GetApiClient(ctx)

	// Track results for reporting
	results := make([]interface{}, 0)
	errors := make([]error, 0)

	logger.Get().Info().Msgf("Creating %d users", len(usersConfig))

	// Process each user
	for i, userConfig := range usersConfig {
		userInfo, httpRes, err := client.UsersAPI.CreateUserAuthorized(ctx).CreateUser(userConfig).Execute()
		if err := response_inspector.InspectResponse(httpRes, err); err != nil {
			logger.Get().Error().Msgf("Failed to create user %d: %s", i+1, err)
			errors = append(errors, fmt.Errorf("user %d (%s): %s", i+1, userConfig.Email, err))
			continue
		}

		results = append(results, userInfo)
		logger.Get().Info().Msgf("Created user %d: %s", i+1, userConfig.Email)
	}

	// Print summary
	logger.Get().Info().Msgf("Bulk user creation complete: %d created, %d failed", len(results), len(errors))

	// Print any errors that occurred
	errorsText := ""
	if len(errors) > 0 {
		logger.Get().Error().Msgf("Errors encountered during bulk creation:")
		for _, err := range errors {
			logger.Get().Error().Msgf("  - %s", err)
			errorsText += fmt.Sprintf("\n  - %s", err)
		}
	}

	// Print the successfully created users
	if len(results) > 0 {
		err = formatter.PrintResult(results, &userPrintConfig)
	}

	if len(errors) > 0 || err != nil {
		if err != nil {
			errorsText += fmt.Sprintf("\n  - %s", err)
		}
		return fmt.Errorf("bulk user creation completed with errors: %s", errorsText)
	}

	return nil
}

func Archive(ctx context.Context, userId string) error {
	logger.Get().Info().Msgf("Archiving user '%s'", userId)

	userIdNumber, revision, err := getUserIdAndRevision(ctx, userId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	userInfo, httpRes, err := client.UsersAPI.ArchiveUser(ctx, userIdNumber).IfMatch(revision).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("User '%s' archived", userId)
	return formatter.PrintResult(userInfo, &userPrintConfig)
}

func Unarchive(ctx context.Context, userId string) error {
	logger.Get().Info().Msgf("Un-archiving user '%s'", userId)

	userIdNumber, revision, err := getUserIdAndRevision(ctx, userId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	userInfo, httpRes, err := client.UsersAPI.UnarchiveUser(ctx, userIdNumber).IfMatch(revision).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("User '%s' un-archived", userId)
	return formatter.PrintResult(userInfo, &userPrintConfig)
}

func GetLimits(ctx context.Context, userId string) error {
	logger.Get().Info().Msgf("Get user '%s' limits", userId)

	userIdNumber, err := getUserId(userId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	userLimits, httpRes, err := client.UsersAPI.GetUserLimits(ctx, userIdNumber).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(userLimits, &userLimitsPrintConfig)
}

func UpdateLimits(ctx context.Context, userId string, config []byte) error {
	logger.Get().Info().Msgf("Updating limits for user '%s'", userId)

	userIdNumber, revision, err := getUserIdAndRevision(ctx, userId)
	if err != nil {
		return err
	}

	var userLimitsConfig sdk.UserLimits
	err = utils.UnmarshalContent(config, &userLimitsConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	userLimits, httpRes, err := client.UsersAPI.UpdateUserLimits(ctx, userIdNumber).UserLimits(userLimitsConfig).IfMatch(revision).Execute()
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
	err = utils.UnmarshalContent(config, &userConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	userConfiguration, httpRes, err := client.UsersAPI.UpdateUserConfig(ctx, userIdNumber).UpdateUser(userConfig).IfMatch(revision).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Configuration updated for user '%s'", userId)
	return formatter.PrintResult(userConfiguration, nil)
}

func ChangeAccount(ctx context.Context, userId string, accountId int) error {
	logger.Get().Info().Msgf("Changing account for user '%s' to account '%d'", userId, accountId)

	userIdNumber, revision, err := getUserIdAndRevision(ctx, userId)
	if err != nil {
		return err
	}

	changeAccount := sdk.ChangeUserAccount{
		NewAccountId: float32(accountId),
	}

	client := api.GetApiClient(ctx)

	userInfo, httpRes, err := client.UsersAPI.ChangeUserAccount(ctx, userIdNumber).ChangeUserAccount(changeAccount).IfMatch(revision).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Account changed for user '%s' to account '%d'", userId, accountId)
	return formatter.PrintResult(userInfo, &userPrintConfig)
}

func GetSSHKeys(ctx context.Context, userId string) error {
	logger.Get().Info().Msgf("Getting SSH keys for user '%s'", userId)

	userIdNumber, err := getUserId(userId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	sshKeys, httpRes, err := client.UsersAPI.GetUserSshKeys(ctx, userIdNumber).Execute()
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

	sshKeyData := sdk.CreateUserSSHKey{
		SshKey: keyContent,
	}

	client := api.GetApiClient(ctx)

	sshKey, httpRes, err := client.UsersAPI.AddUserSshKey(ctx, userIdNumber).CreateUserSSHKey(sshKeyData).Execute()
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

	httpRes, err := client.UsersAPI.DeleteUserSshKey(ctx, userIdNumber, float32(keyIdNumber)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("SSH key '%s' deleted for user '%s'", keyId, userId)
	return nil
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

	suspendInfo, httpRes, err := client.UsersAPI.SuspendUser(ctx, userIdNumber).UserSuspend(suspendReason).IfMatch(revision).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("User '%s' suspended", userId)
	return formatter.PrintResult(suspendInfo, nil)
}

func Unsuspend(ctx context.Context, userId string) error {
	logger.Get().Info().Msgf("Un-suspending user '%s'", userId)

	userIdNumber, revision, err := getUserIdAndRevision(ctx, userId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.UsersAPI.UnsuspendUser(ctx, userIdNumber).IfMatch(revision).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("User '%s' un-suspended", userId)
	return nil
}

func GetPermissions(ctx context.Context, userId string) error {
	logger.Get().Info().Msgf("Getting permissions for user '%s'", userId)

	userIdNumber, err := getUserId(userId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	permissions, httpRes, err := client.UsersAPI.GetUserPermissions(ctx, userIdNumber).Execute()
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

	var permissionsConfig sdk.UpdateUserPermissions
	err = utils.UnmarshalContent(config, &permissionsConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	permissions, httpRes, err := client.UsersAPI.UpdateUserPermissions(ctx, userIdNumber).UpdateUserPermissions(permissionsConfig).IfMatch(revision).Execute()
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

	user, httpRes, err := client.UsersAPI.GetUser(ctx, float32(userIdNumeric)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return 0, "", err
	}

	return float32(userIdNumeric), strconv.Itoa(int(user.Revision)), nil
}
