package user

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
			Title: "ID",
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

// userLimitsPrintConfig formats the effective quota limits (QuotaProfileLimits)
// returned by the quota-limits-breakdown endpoint.
var userLimitsPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"InfrastructureServerGroupMaxCount": {
			Title: "Max Server Groups / Infra",
			Order: 1,
		},
		"InfrastructureDriveMaxCount": {
			Title: "Max Drives / Infra",
			Order: 2,
		},
		"InfrastructureVmInstanceGroupMaxCount": {
			Title: "Max VM Groups / Infra",
			Order: 3,
		},
		"ServerGroupInstancesMaxCount": {
			Title: "Max Instances / Server Group",
			Order: 4,
		},
		"VmInstanceGroupVmInstancesMaxCount": {
			Title: "Max VMs / VM Group",
			Order: 5,
		},
		"DriveMaxSizeMbytes": {
			Title: "Max Drive Size (MB)",
			Order: 6,
		},
		"UserSshKeysCountMax": {
			Title: "Max SSH Keys",
			Order: 7,
		},
	},
}

var userSshKeysPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
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

func List(ctx context.Context, archived bool, filterId, filterDisplayName, filterEmail, filterAccountId, filterInfrastructureId, sortBy, search, searchBy string) error {
	logger.Get().Info().Msgf("Listing all users")

	client := api.GetApiClient(ctx)

	archivedFilter := "0"
	if archived {
		archivedFilter = "1"
	}

	request := client.UsersAPI.GetUsers(ctx).FilterArchived([]string{archivedFilter})
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
	} else {
		request = request.SortBy([]string{"id:ASC"})
	}
	if search != "" {
		request = request.Search(search)
	}
	if searchBy != "" {
		request = request.SearchBy([]string{searchBy})
	}

	users, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(users, meta, len(users), &userPrintConfig)
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

	// Per-user limits are no longer stored directly on the user. The effective
	// limits are now derived from the user's role/group/account quota profiles and
	// exposed through the quota-limits-breakdown endpoint.
	breakdown, httpRes, err := client.UsersAPI.GetQuotaLimitsBreakdown(ctx, userIdNumber).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(breakdown.Effective, &userLimitsPrintConfig)
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
		NewAccountId: int64(accountId),
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

	keyIdNumber, err := strconv.ParseInt(keyId, 10, 64)
	if err != nil {
		err := fmt.Errorf("invalid SSH key ID: '%s'", keyId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.UsersAPI.DeleteUserSshKey(ctx, userIdNumber, keyIdNumber).Execute()
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

func SetPassword(ctx context.Context, userId string, password string) error {
	logger.Get().Info().Msgf("Setting password for user '%s'", userId)

	userIdNumber, revision, err := getUserIdAndRevision(ctx, userId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)
	cfg := client.GetConfig()

	if len(cfg.Servers) == 0 {
		return fmt.Errorf("no API server configured")
	}
	baseURL := cfg.Servers[0].URL

	body, err := json.Marshal(map[string]string{"password": password})
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	url := fmt.Sprintf("%s/api/v2/users/%v/actions/set-password", baseURL, userIdNumber)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("If-Match", revision)

	if token, ok := ctx.Value(sdk.ContextAccessToken).(string); ok {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	httpRes, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer func() {
		io.Copy(io.Discard, httpRes.Body) //nolint:errcheck
		httpRes.Body.Close()
	}()

	respBody, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}
	httpRes.Body = io.NopCloser(bytes.NewBuffer(respBody))

	if err := response_inspector.InspectResponse(httpRes, nil); err != nil {
		return err
	}

	var userInfo sdk.User
	if err := json.Unmarshal(respBody, &userInfo); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	logger.Get().Info().Msgf("Password set for user '%s'", userId)
	return formatter.PrintResult(userInfo, &userPrintConfig)
}

func getUserId(userId string) (int64, error) {
	userIdNumeric, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		err := fmt.Errorf("invalid user ID: '%s'", userId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return userIdNumeric, nil
}

func getUserIdAndRevision(ctx context.Context, userId string) (int64, string, error) {
	userIdNumeric, err := getUserId(userId)
	if err != nil {
		return 0, "", err
	}

	client := api.GetApiClient(ctx)

	user, httpRes, err := client.UsersAPI.GetUser(ctx, userIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return 0, "", err
	}

	return userIdNumeric, strconv.Itoa(int(user.Revision)), nil
}

// userConfigurationPrintConfig formats the user configuration object returned
// by the /users/{userId}/config endpoint.
var userConfigurationPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"DisplayName": {
			Title:    "Name",
			MaxWidth: 30,
			Order:    1,
		},
		"AccessLevel": {
			Title: "Role",
			Order: 2,
		},
		"EmailStatus": {
			Title: "E-mail Status",
			Order: 3,
		},
		"Language": {
			Title: "Language",
			Order: 4,
		},
		"Brand": {
			Title: "Brand",
			Order: 5,
		},
		"IsBlocked": {
			Title: "Blocked",
			Order: 6,
		},
		"PasswordChangeRequired": {
			Title: "Password Change Required",
			Order: 7,
		},
		"LastLoginTimestamp": {
			Title:       "Last Login",
			Transformer: formatter.FormatDateTimeValue,
			Order:       8,
		},
		"Revision": {
			Title: "Revision",
			Order: 9,
		},
	},
}

// userSuspendReasonsPrintConfig formats the suspend reason history of a user.
var userSuspendReasonsPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"UserId": {
			Title: "User ID",
			Order: 2,
		},
		"Type": {
			Title: "Type",
			Order: 3,
		},
		"PublicComment": {
			Title:    "Public Comment",
			MaxWidth: 40,
			Order:    4,
		},
		"PrivateComment": {
			Title:    "Private Comment",
			MaxWidth: 40,
			Order:    5,
		},
		"CreatedTimestamp": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       6,
		},
		"EndTimestamp": {
			Title:       "Ended",
			Transformer: formatter.FormatDateTimeValue,
			Order:       7,
		},
	},
}

// userDelegatesPrintConfig formats the UserInfo records returned by the
// parent/child delegate endpoints.
var userDelegatesPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
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
		"AccountId": {
			Title: "Account ID",
			Order: 5,
		},
		"IsArchived": {
			Title: "Archived",
			Order: 6,
		},
		"CreatedTimestamp": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       7,
		},
	},
}

// GetConfiguration retrieves the configuration object of a user.
func GetConfiguration(ctx context.Context, userId string) error {
	logger.Get().Info().Msgf("Get configuration for user '%s'", userId)

	userIdNumber, err := getUserId(userId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	userConfiguration, httpRes, err := client.UsersAPI.GetUserConfiguration(ctx, userIdNumber).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(userConfiguration, &userConfigurationPrintConfig)
}

// UpdateMeta replaces the metadata (GUI settings) of a user.
func UpdateMeta(ctx context.Context, userId string, config []byte) error {
	logger.Get().Info().Msgf("Updating metadata for user '%s'", userId)

	userIdNumber, err := getUserId(userId)
	if err != nil {
		return err
	}

	var userMeta sdk.UserMeta
	if err := utils.UnmarshalContent(config, &userMeta); err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	updatedMeta, httpRes, err := client.UsersAPI.UpdateUserMeta(ctx, userIdNumber).UserMeta(userMeta).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Metadata updated for user '%s'", userId)
	return formatter.PrintResult(updatedMeta, nil)
}

// GetSSHKey retrieves a single SSH key of a user.
func GetSSHKey(ctx context.Context, userId string, keyId string) error {
	logger.Get().Info().Msgf("Get SSH key '%s' of user '%s'", keyId, userId)

	userIdNumber, err := getUserId(userId)
	if err != nil {
		return err
	}

	keyIdNumber, err := getSshKeyId(keyId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	sshKey, httpRes, err := client.UsersAPI.GetUserSshKey(ctx, userIdNumber, keyIdNumber).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(sshKey, &userSshKeysPrintConfig)
}

// GetSuspendReasons lists the suspend reasons recorded for a user.
func GetSuspendReasons(ctx context.Context, userId string) error {
	logger.Get().Info().Msgf("Get suspend reasons for user '%s'", userId)

	userIdNumber, err := getUserId(userId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	suspendReasons, httpRes, err := client.UsersAPI.GetUserSuspendReasons(ctx, userIdNumber).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	if suspendReasons == nil {
		return formatter.PrintResult([]sdk.UserSuspendReason{}, &userSuspendReasonsPrintConfig)
	}

	return formatter.PrintResult(suspendReasons.Data, &userSuspendReasonsPrintConfig)
}

// AddDelegate grants a delegate user access to the resources of a user.
func AddDelegate(ctx context.Context, userId string, delegateId string) error {
	logger.Get().Info().Msgf("Adding delegate '%s' to user '%s'", delegateId, userId)

	userIdNumber, err := getUserId(userId)
	if err != nil {
		return err
	}

	delegateIdNumber, err := getDelegateId(delegateId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	userInfo, httpRes, err := client.UsersAPI.AddUserDelegate(ctx, userIdNumber, delegateIdNumber).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Delegate '%s' added to user '%s'", delegateId, userId)
	return formatter.PrintResult(userInfo, &userPrintConfig)
}

// RemoveDelegate revokes the delegate access previously granted to a user.
func RemoveDelegate(ctx context.Context, userId string, delegateId string) error {
	logger.Get().Info().Msgf("Removing delegate '%s' from user '%s'", delegateId, userId)

	userIdNumber, err := getUserId(userId)
	if err != nil {
		return err
	}

	delegateIdNumber, err := getDelegateId(delegateId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	userInfo, httpRes, err := client.UsersAPI.RemoveUserDelegate(ctx, userIdNumber, delegateIdNumber).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Delegate '%s' removed from user '%s'", delegateId, userId)
	return formatter.PrintResult(userInfo, &userPrintConfig)
}

// GetParentDelegates lists the users that delegated access to the given user.
func GetParentDelegates(ctx context.Context, userId string) error {
	logger.Get().Info().Msgf("Get parent delegates of user '%s'", userId)

	userIdNumber, err := getUserId(userId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	delegates, httpRes, err := client.UsersAPI.GetUserParentDelegates(ctx, userIdNumber).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return printDelegates(delegates)
}

// GetChildDelegates lists the users the given user delegated access to.
func GetChildDelegates(ctx context.Context, userId string) error {
	logger.Get().Info().Msgf("Get child delegates of user '%s'", userId)

	userIdNumber, err := getUserId(userId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	delegates, httpRes, err := client.UsersAPI.GetUserChildDelegates(ctx, userIdNumber).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return printDelegates(delegates)
}

// ResendEmailVerification sends the e-mail verification message again.
func ResendEmailVerification(ctx context.Context, userId string, redirectUrl string) error {
	logger.Get().Info().Msgf("Resending e-mail verification for user '%s'", userId)

	userIdNumber, err := getUserId(userId)
	if err != nil {
		return err
	}

	// The body is optional in the API definition, but the SDK sends a literal
	// `null` payload when it is not set, which the API rejects with 400.
	requestBody := sdk.ResendUserVerificationEmail{}
	if redirectUrl != "" {
		requestBody.RedirectUrl = sdk.PtrString(redirectUrl)
	}

	client := api.GetApiClient(ctx)

	userInfo, httpRes, err := client.UsersAPI.ResendEmailVerification(ctx, userIdNumber).
		ResendUserVerificationEmail(requestBody).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("E-mail verification resent for user '%s'", userId)
	return formatter.PrintResult(userInfo, &userPrintConfig)
}

// ResendInvitation sends the account invitation message again.
func ResendInvitation(ctx context.Context, userId string, redirectUrl string) error {
	logger.Get().Info().Msgf("Resending invitation for user '%s'", userId)

	userIdNumber, err := getUserId(userId)
	if err != nil {
		return err
	}

	requestBody := sdk.ResendUserInvitation{}
	if redirectUrl != "" {
		requestBody.RedirectUrl = sdk.PtrString(redirectUrl)
	}

	client := api.GetApiClient(ctx)

	userInfo, httpRes, err := client.UsersAPI.ResendUserInvitation(ctx, userIdNumber).
		ResendUserInvitation(requestBody).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Invitation resent for user '%s'", userId)
	return formatter.PrintResult(userInfo, &userPrintConfig)
}

// SendPasswordReset sends a password reset message to a user as administrator.
func SendPasswordReset(ctx context.Context, userId string, redirectUrl string) error {
	logger.Get().Info().Msgf("Sending password reset for user '%s'", userId)

	userIdNumber, err := getUserId(userId)
	if err != nil {
		return err
	}

	requestBody := sdk.PasswordResetByAdmin{}
	if redirectUrl != "" {
		requestBody.RedirectUrl = sdk.PtrString(redirectUrl)
	}

	client := api.GetApiClient(ctx)

	userInfo, httpRes, err := client.UsersAPI.SendPasswordResetByAdmin(ctx, userIdNumber).
		PasswordResetByAdmin(requestBody).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Password reset sent for user '%s'", userId)
	return formatter.PrintResult(userInfo, &userPrintConfig)
}

// Delete archives a user and irreversibly removes their personal information.
// Unlike Archive, this operation cannot be undone.
func Delete(ctx context.Context, userId string) error {
	logger.Get().Info().Msgf("Deleting user '%s'", userId)

	userIdNumber, revision, err := getUserIdAndRevision(ctx, userId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	userInfo, httpRes, err := client.UsersAPI.DeleteUser(ctx, userIdNumber).IfMatch(revision).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("User '%s' deleted", userId)

	// The endpoint may answer with an empty body once the personal information
	// has been removed.
	if userInfo == nil {
		return nil
	}

	return formatter.PrintResult(userInfo, &userPrintConfig)
}

func printDelegates(delegates *sdk.UserList) error {
	if delegates == nil {
		return formatter.PrintResult([]sdk.UserInfo{}, &userDelegatesPrintConfig)
	}

	return formatter.PrintResult(delegates.Data, &userDelegatesPrintConfig)
}

func getSshKeyId(keyId string) (int64, error) {
	keyIdNumeric, err := utils.GetInt64FromString(keyId)
	if err != nil {
		err := fmt.Errorf("invalid SSH key ID: '%s'", keyId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return keyIdNumeric, nil
}

func getDelegateId(delegateId string) (int64, error) {
	delegateIdNumeric, err := utils.GetInt64FromString(delegateId)
	if err != nil {
		err := fmt.Errorf("invalid delegate user ID: '%s'", delegateId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return delegateIdNumeric, nil
}
