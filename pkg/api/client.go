package api

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
	"github.com/spf13/cobra"
)

type ContextKey string

const (
	ApiClientContextKey       ContextKey = "apiClient"
	UserIdContextKey          ContextKey = "userId"
	UserAccessLevelContextKey ContextKey = "userAccessLevel"
)

func GetApiClient(ctx context.Context) *sdk.APIClient {
	client := ctx.Value(ApiClientContextKey)
	if client == nil {
		err := fmt.Errorf("SDK client not found in context")
		logger.Get().Error().Err(err).Msg("")
		cobra.CheckErr(err)
	}

	apiClient, ok := client.(*sdk.APIClient)
	if !ok {
		err := fmt.Errorf("invalid SDK client in context")
		logger.Get().Error().Err(err).Msg("")
		cobra.CheckErr(err)
	}

	return apiClient
}

func GetApiClientE(ctx context.Context) (*sdk.APIClient, error) {
	client := ctx.Value(ApiClientContextKey)
	if client == nil {
		err := fmt.Errorf("SDK client not found in context")
		logger.Get().Error().Err(err).Msg("")
		return nil, err
	}

	apiClient, ok := client.(*sdk.APIClient)
	if !ok {
		err := fmt.Errorf("invalid SDK client in context")
		logger.Get().Error().Err(err).Msg("")
		return nil, err
	}

	return apiClient, nil
}

func SetApiClient(ctx context.Context, apiEndpoint string, apiKey string, debug bool, insecure bool) context.Context {
	// Initialize API client using the arguments from the command line or environment variables
	cfg := sdk.NewConfiguration()
	cfg.UserAgent = "metalcloud-cli"
	cfg.Servers = []sdk.ServerConfiguration{
		{
			URL:         apiEndpoint,
			Description: "MetalSoft",
		},
	}

	if insecure {
		cfg.HTTPClient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	}

	// Set debug mode
	cfg.Debug = debug

	// Create API client
	apiClient := sdk.NewAPIClient(cfg)

	ctx = context.WithValue(ctx, ApiClientContextKey, apiClient)
	ctx = context.WithValue(ctx, sdk.ContextAccessToken, apiKey)

	return ctx
}

// DoJSONRequest issues an HTTP request against the configured API endpoint using
// the same base URL, auth token, and HTTP client as the SDK, but bypassing the
// typed SDK request/response models. Use it for endpoints whose generated types
// reject valid payloads due to schema drift — e.g. polymorphic `oneOf` unions
// that fail to match an empty object. `body` may be nil for bodyless requests.
// The caller owns closing the returned response body.
func DoJSONRequest(ctx context.Context, method, path string, body []byte) (*http.Response, error) {
	return DoJSONRequestWithHeaders(ctx, method, path, body, nil)
}

// DoJSONRequestWithHeaders is DoJSONRequest with additional request headers
// (e.g. If-Match for optimistic concurrency). Headers set here override the
// defaults. The caller owns closing the returned response body.
func DoJSONRequestWithHeaders(ctx context.Context, method, path string, body []byte, headers map[string]string) (*http.Response, error) {
	apiClient, err := GetApiClientE(ctx)
	if err != nil {
		return nil, err
	}

	cfg := apiClient.GetConfig()
	if cfg == nil || len(cfg.Servers) == 0 {
		return nil, fmt.Errorf("no API server configured")
	}
	url := strings.TrimRight(cfg.Servers[0].URL, "/") + path

	var reqBody io.Reader
	if body != nil {
		reqBody = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, err
	}

	if token, ok := ctx.Value(sdk.ContextAccessToken).(string); ok && token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	if cfg.UserAgent != "" {
		req.Header.Set("User-Agent", cfg.UserAgent)
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return httpClient.Do(req)
}

func GetUserId(ctx context.Context) string {
	userId := ctx.Value(UserIdContextKey)
	if userId == nil {
		err := fmt.Errorf("user ID not found in context")
		logger.Get().Error().Err(err).Msg("")
		cobra.CheckErr(err)
	}

	userIdStr, ok := userId.(string)
	if !ok {
		err := fmt.Errorf("invalid user ID in context")
		logger.Get().Error().Err(err).Msg("")
		cobra.CheckErr(err)
	}

	return userIdStr
}

func SetUserId(ctx context.Context, userId string) context.Context {
	return context.WithValue(ctx, UserIdContextKey, userId)
}

func GetUserAccessLevel(ctx context.Context) string {
	level := ctx.Value(UserAccessLevelContextKey)
	if level == nil {
		return ""
	}
	levelStr, ok := level.(string)
	if !ok {
		return ""
	}
	return levelStr
}

func SetUserAccessLevel(ctx context.Context, level string) context.Context {
	return context.WithValue(ctx, UserAccessLevelContextKey, level)
}

// IsAdminAccessLevel reports whether the given user access level grants
// administrative scope. Values are sourced from the MetalSoft user model.
func IsAdminAccessLevel(level string) bool {
	switch level {
	case "root", "admin", "full_admin":
		return true
	}
	return false
}
