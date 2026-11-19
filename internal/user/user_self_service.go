package user

import (
	"context"
	"fmt"
	"net/http"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

// userPermissionsPrintConfig formats the permission keys of the current user.
var userPermissionsPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Permission": {
			Title: "Permission",
			Order: 1,
		},
	},
}

// userPermissionRecord wraps a permission key so that the tabular formatters
// have a named column to render.
type userPermissionRecord struct {
	Permission string `json:"permission"`
}

// hasUnknownPermission reports whether the SDK replaced a permission it does not
// know about with its generic "unknown" enum value.
func hasUnknownPermission(permissions []sdk.MetalsoftPermissions) bool {
	for _, permission := range permissions {
		if permission == sdk.METALSOFTPERMISSIONS_UNKNOWN_DEFAULT_OPEN_API {
			return true
		}
	}
	return false
}

// GetPermissions lists the permissions of the user owning the API key.
func GetPermissions(ctx context.Context) error {
	logger.Get().Info().Msgf("Getting permissions for current user")

	client := api.GetApiClient(ctx)

	// The response is inspected by hand here: a failed typed decode is expected
	// often enough that it must not be logged as an error before the fallback.
	permissions, httpRes, err := client.UserAPI.GetUserPermissions(ctx).Execute()
	if err == nil && httpRes != nil && httpRes.StatusCode < 400 &&
		permissions != nil && len(permissions.Permissions) > 0 && !hasUnknownPermission(permissions.Permissions) {
		records := make([]userPermissionRecord, 0, len(permissions.Permissions))
		for _, permission := range permissions.Permissions {
			records = append(records, userPermissionRecord{Permission: string(permission)})
		}

		return formatter.PrintResult(records, &userPermissionsPrintConfig)
	}

	// The SDK models the permission keys as an enum generated from the OpenAPI
	// definition. A permission it does not know about is either rejected (the
	// typed response fails to decode) or collapsed into a generic "unknown"
	// value that hides its real name. Read the raw payload instead, which keeps
	// the real names and also surfaces the real HTTP error when there is one.
	logger.Get().Debug().Msg("Reading the user permissions from the raw response")

	body, err := api.RawJSONRequest(ctx, http.MethodGet, "/api/v2/user/permissions", nil, nil)
	if err != nil {
		return err
	}

	rawPermissions, err := utils.DecodeRawObject(body)
	if err != nil {
		return err
	}

	permissionList, ok := rawPermissions["permissions"].([]any)
	if !ok {
		return formatter.PrintResult(rawPermissions, nil)
	}

	records := make([]userPermissionRecord, 0, len(permissionList))
	for _, permission := range permissionList {
		records = append(records, userPermissionRecord{Permission: fmt.Sprintf("%v", permission)})
	}

	return formatter.PrintResult(records, &userPermissionsPrintConfig)
}

// InitiatePasswordReset sends a password reset e-mail to the given address.
func InitiatePasswordReset(ctx context.Context, email string, redirectUrl string) error {
	logger.Get().Info().Msgf("Initiating password reset for '%s'", email)

	requestBody := sdk.PasswordReset{
		Email: email,
	}
	if redirectUrl != "" {
		requestBody.RedirectUrl = sdk.PtrString(redirectUrl)
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.UserAPI.InitiatePasswordReset(ctx).PasswordReset(requestBody).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	fmt.Printf("Password reset e-mail sent to '%s'.\n", email)
	return nil
}

// VerifyEmail consumes the e-mail verification token sent by the platform. It
// is the programmatic equivalent of following the link from the e-mail.
func VerifyEmail(ctx context.Context, token string) error {
	logger.Get().Info().Msgf("Verifying e-mail address with the provided token")

	client := api.GetApiClient(ctx)

	httpRes, err := client.UserAPI.HandleEmailVerify(ctx).Token(token).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	fmt.Println("E-mail address verified.")
	return nil
}

// ResetPassword consumes the password reset token sent by the platform. It is
// the programmatic equivalent of following the link from the e-mail; the new
// password itself is set through the web page the token redirects to.
func ResetPassword(ctx context.Context, token string) error {
	logger.Get().Info().Msgf("Handling password reset with the provided token")

	client := api.GetApiClient(ctx)

	httpRes, err := client.UserAPI.HandleUserResetPassword(ctx).Token(token).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	fmt.Println("Password reset token accepted.")
	return nil
}
