package user

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/png"
	"net/http"
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

func GetApiKey(ctx context.Context) error {
	logger.Get().Info().Msgf("Getting API key for current user")

	client := api.GetApiClient(ctx)

	apiKey, httpRes, err := client.UserAPI.GetUserApiKey(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	fmt.Println(apiKey.ApiKey)
	return nil
}

func RegenerateApiKey(ctx context.Context) error {
	logger.Get().Info().Msgf("Regenerating API key for current user")

	revision, err := getCurrentUserRevision(ctx)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	userInfo, httpRes, err := client.UserAPI.RegenerateUserApiKey(ctx).IfMatch(revision).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	fmt.Println("API key has been regenerated. Your previous key is no longer valid.")
	return formatter.PrintResult(userInfo, &userPrintConfig)
}

func Enable2FA(ctx context.Context, token string) error {
	logger.Get().Info().Msgf("Enabling 2FA for current user")

	revision, err := getCurrentUserRevision(ctx)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	recoveryCodes, httpRes, err := client.UserAPI.EnableUser2FA(ctx).
		TwoFactorAuthenticationToken(sdk.TwoFactorAuthenticationToken{
			Token: sdk.PtrString(token),
		}).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	fmt.Printf("Two-factor authentication enabled successfully. Recovery codes: %s\n", strings.Join(recoveryCodes.RecoveryCodes, ", "))
	return nil
}

func Disable2FA(ctx context.Context) error {
	logger.Get().Info().Msgf("Disabling 2FA for current user")

	revision, err := getCurrentUserRevision(ctx)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.UserAPI.DisableUser2FA(ctx).IfMatch(revision).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	fmt.Println("Two-factor authentication disabled.")
	return nil
}

var twoFASecretPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Secret2FA": {
			Title: "Secret",
			Order: 1,
		},
		"QrCode": {
			Title: "QR Code",
			Order: 2,
		},
	},
}

func GenerateUser2FASecret(ctx context.Context) error {
	logger.Get().Info().Msgf("Generating 2FA secret for current user")

	client := api.GetApiClient(ctx)

	secret, httpRes, err := client.UserAPI.GenerateUser2FASecret(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	if !formatter.IsTextFormat() {
		return formatter.PrintResult(secret, &twoFASecretPrintConfig)
	}

	fmt.Printf("Secret: %s\n\n", secret.Secret2FA)

	qrArt, err := renderQRCodeFromDataURI(secret.QrCode)
	if err != nil {
		fmt.Printf("QR Code: %s\n", secret.QrCode)
		return nil
	}

	fmt.Println("QR Code:")
	fmt.Println(qrArt)
	return nil
}

// renderQRCodeFromDataURI decodes a data:image/png;base64,... URI and renders
// the resulting image as Unicode block art suitable for terminal display.
func renderQRCodeFromDataURI(dataURI string) (string, error) {
	// Strip the data URI prefix
	const prefix = "data:image/png;base64,"
	b64Data := dataURI
	if strings.HasPrefix(dataURI, prefix) {
		b64Data = dataURI[len(prefix):]
	}

	imgBytes, err := base64.StdEncoding.DecodeString(b64Data)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	img, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		return "", fmt.Errorf("failed to decode PNG: %w", err)
	}

	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	// Convert image to binary black/white grid
	// Use luminance threshold to determine black vs white
	grid := make([][]bool, height)
	for y := 0; y < height; y++ {
		grid[y] = make([]bool, width)
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(bounds.Min.X+x, bounds.Min.Y+y).RGBA()
			// Luminance (values are 0-65535 from RGBA())
			lum := (299*r + 587*g + 114*b) / 1000
			grid[y][x] = lum < 32768 // true = dark/black
		}
	}

	// Detect the QR module size by finding the first dark-to-light transition
	// in the top-left finder pattern area
	moduleSize := detectModuleSize(grid, width, height)
	if moduleSize < 1 {
		moduleSize = 1
	}

	// Sample the grid at module-level resolution
	modulesW := width / moduleSize
	modulesH := height / moduleSize

	modules := make([][]bool, modulesH)
	for my := 0; my < modulesH; my++ {
		modules[my] = make([]bool, modulesW)
		for mx := 0; mx < modulesW; mx++ {
			// Sample center of each module
			sx := mx*moduleSize + moduleSize/2
			sy := my*moduleSize + moduleSize/2
			if sx < width && sy < height {
				modules[my][mx] = grid[sy][sx]
			}
		}
	}

	// Render using Unicode half-block characters
	// Each character row represents 2 pixel rows
	// ▀ = top black, bottom white
	// ▄ = top white, bottom black
	// █ = both black
	// (space) = both white
	var sb strings.Builder
	for y := 0; y < modulesH; y += 2 {
		for x := 0; x < modulesW; x++ {
			top := modules[y][x]
			bot := false
			if y+1 < modulesH {
				bot = modules[y+1][x]
			}
			switch {
			case top && bot:
				sb.WriteString("█")
			case top && !bot:
				sb.WriteString("▀")
			case !top && bot:
				sb.WriteString("▄")
			default:
				sb.WriteString(" ")
			}
		}
		sb.WriteString("\n")
	}

	return sb.String(), nil
}

// detectModuleSize examines the top row of the QR code to determine the pixel
// size of each module by measuring the first dark run in the finder pattern.
func detectModuleSize(grid [][]bool, width, height int) int {
	if height == 0 || width == 0 {
		return 1
	}

	// Find the first row that has dark pixels (skip any white border)
	startRow := 0
	for startRow < height {
		hasDark := false
		for x := 0; x < width; x++ {
			if grid[startRow][x] {
				hasDark = true
				break
			}
		}
		if hasDark {
			break
		}
		startRow++
	}
	if startRow >= height {
		return 1
	}

	// In this row, find the first dark pixel
	startCol := 0
	for startCol < width && !grid[startRow][startCol] {
		startCol++
	}
	if startCol >= width {
		return 1
	}

	// Count consecutive dark pixels — this is one module width in the finder pattern
	count := 0
	for startCol+count < width && grid[startRow][startCol+count] {
		count++
	}

	// The finder pattern starts with 7 modules of dark-light-dark-light-dark-light-dark
	// The first run is 1 module wide, but we might be hitting the full 7-module pattern
	// on a zoomed-in image. For a standard QR finder pattern, the first dark run at the
	// top-left corner spans exactly 7 modules (the full top bar of the finder pattern).
	// So divide by 7 to get the module size.
	if count >= 7 {
		return count / 7
	}

	return count
}

func getCurrentUserRevision(ctx context.Context) (string, error) {
	userIdStr := api.GetUserId(ctx)
	userIdNumeric, err := getUserId(userIdStr)
	if err != nil {
		return "", err
	}

	client := api.GetApiClient(ctx)

	userInfo, httpRes, err := client.UsersAPI.GetUser(ctx, userIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", int(userInfo.Revision)), nil
}

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

// GetPermissions lists the permissions of the user owning the API key.
func GetPermissions(ctx context.Context) error {
	logger.Get().Info().Msgf("Getting permissions for current user")

	client := api.GetApiClient(ctx)

	// The response is inspected by hand here: a failed typed decode is expected
	// often enough that it must not be logged as an error before the fallback.
	permissions, httpRes, err := client.UserAPI.GetUserPermissions(ctx).Execute()
	if err == nil && httpRes != nil && httpRes.StatusCode < 400 &&
		permissions != nil && len(permissions.Permissions) > 0 {
		records := make([]userPermissionRecord, 0, len(permissions.Permissions))
		for _, permission := range permissions.Permissions {
			records = append(records, userPermissionRecord{Permission: string(permission)})
		}

		return formatter.PrintResult(records, &userPermissionsPrintConfig)
	}

	// The SDK models the permission keys as a strict enum generated from the
	// OpenAPI definition, so a single permission the SDK does not know about
	// makes the whole typed response fail to decode. Read the raw payload
	// instead, which also surfaces the real HTTP error when there is one.
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

// ChangePassword updates the password of the user owning the API key.
func ChangePassword(ctx context.Context, config []byte) error {
	logger.Get().Info().Msgf("Changing password for current user")

	var passwordUpdate sdk.UserUpdatePassword
	if err := utils.UnmarshalContent(config, &passwordUpdate); err != nil {
		return err
	}

	revision, err := getCurrentUserRevision(ctx)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	userInfo, httpRes, err := client.UserAPI.UpdateUserPassword(ctx).
		UserUpdatePassword(passwordUpdate).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	fmt.Println("Password changed.")
	return formatter.PrintResult(userInfo, &userPrintConfig)
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

// InitiateEmailChange starts the e-mail change flow for the current user. The
// new address becomes active only after the verification link is followed.
func InitiateEmailChange(ctx context.Context, email string, redirectUrl string) error {
	logger.Get().Info().Msgf("Initiating e-mail change to '%s' for current user", email)

	requestBody := sdk.ChangeUserEmail{
		Email: email,
	}
	if redirectUrl != "" {
		requestBody.RedirectUrl = sdk.PtrString(redirectUrl)
	}

	client := api.GetApiClient(ctx)

	userInfo, httpRes, err := client.UserAPI.InitiateEmailChange(ctx).ChangeUserEmail(requestBody).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	fmt.Printf("Verification e-mail sent to '%s'. The address changes once the link is followed.\n", email)
	return formatter.PrintResult(userInfo, &userPrintConfig)
}

// RegenerateJwtSalt rotates the JWT salt of the current user, which invalidates
// every session and token issued so far.
func RegenerateJwtSalt(ctx context.Context) error {
	logger.Get().Info().Msgf("Regenerating JWT salt for current user")

	revision, err := getCurrentUserRevision(ctx)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	userInfo, httpRes, err := client.UserAPI.RegenerateUserJwtSalt(ctx).IfMatch(revision).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	fmt.Println("JWT salt regenerated. All existing sessions and tokens are now invalid.")
	return formatter.PrintResult(userInfo, &userPrintConfig)
}

// VerifyEmail consumes the e-mail verification token sent by the platform. It
// is the programmatic equivalent of following the link from the e-mail.
func VerifyEmail(ctx context.Context, token string) error {
	logger.Get().Info().Msgf("Verifying e-mail address with the provided token")

	client := api.GetApiClient(ctx)

	httpRes, err := client.UserAPI.UserControllerHandleEmailVerify(ctx).Token(token).Execute()
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

	httpRes, err := client.UserAPI.UserControllerHandleUserResetPassword(ctx).Token(token).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	fmt.Println("Password reset token accepted.")
	return nil
}
