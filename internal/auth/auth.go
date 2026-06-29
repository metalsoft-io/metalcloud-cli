package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

const (
	defaultUserExternalIdentifier = "objectGUID"
	defaultUsername               = "sAMAccountName"
	defaultEmail                  = "mail"
)

type AuthLdapGroupMapping struct {
	GroupName              string `json:"groupName"`
	RoleName               string `json:"roleName"`
	Priority               int32  `json:"priority"`
	UserExternalIdentifier string `json:"userExternalIdentifier"`
	Username               string `json:"username"`
	Email                  string `json:"email"`
}

type AuthLdapMappingOptions struct {
	RoleName               string
	Priority               int32
	UserExternalIdentifier *string
	Username               *string
	Email                  *string
}

var authLdapMappingPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"GroupName": {
			Title: "Group",
			Order: 1,
		},
		"RoleName": {
			Title:    "Role",
			MaxWidth: 30,
			Order:    2,
		},
		"Priority": {
			Title: "Priority",
			Order: 3,
		},
		"UserExternalIdentifier": {
			Title:    "External Id",
			MaxWidth: 30,
			Order:    4,
		},
		"Username": {
			Title:    "Username",
			MaxWidth: 30,
			Order:    5,
		},
		"Email": {
			Title:    "Email",
			MaxWidth: 30,
			Order:    6,
		},
	},
}

func AuthLdapMappingList(ctx context.Context) error {
	logger.Get().Info().Msgf("Listing all LDAP mappings")

	authConfig, err := getAuthConfig(ctx)
	if err != nil {
		return err
	}

	mappings := getLdapGroupMapping(getLdapConfig(authConfig))

	return formatter.PrintResult(mappings, &authLdapMappingPrintConfig)
}

func AuthLdapMappingAdd(ctx context.Context, groupName string, opts AuthLdapMappingOptions) error {
	logger.Get().Info().Msgf("Adding LDAP mapping for group '%s'", groupName)

	authConfig, err := getAuthConfig(ctx)
	if err != nil {
		return err
	}

	ldapConfig := getLdapConfig(authConfig)
	if ldapConfig == nil {
		return fmt.Errorf("no LDAP configuration found")
	}

	if ldapConfig["groupsMapping"] == nil {
		ldapConfig["groupsMapping"] = []interface{}{}
	}

	for _, mapping := range ldapConfig["groupsMapping"].([]interface{}) {
		if mapping.(map[string]interface{})["groupName"] == groupName {
			return fmt.Errorf("mapping for LDAP group '%s' already exists", groupName)
		}
	}

	newMapping := map[string]interface{}{
		"groupName":              groupName,
		"roleName":               opts.RoleName,
		"priority":               float64(opts.Priority),
		"userExternalIdentifier": defaultUserExternalIdentifier,
		"username":               defaultUsername,
		"email":                  defaultEmail,
	}
	if opts.UserExternalIdentifier != nil && *opts.UserExternalIdentifier != "" {
		newMapping["userExternalIdentifier"] = *opts.UserExternalIdentifier
	}
	if opts.Username != nil && *opts.Username != "" {
		newMapping["username"] = *opts.Username
	}
	if opts.Email != nil && *opts.Email != "" {
		newMapping["email"] = *opts.Email
	}

	ldapConfig["groupsMapping"] = append(ldapConfig["groupsMapping"].([]interface{}), newMapping)

	authConfig, err = patchAuthConfig(ctx, map[string]interface{}{"ldap": ldapConfig})
	if err != nil {
		return err
	}

	mappings := getLdapGroupMapping(getLdapConfig(authConfig))

	return formatter.PrintResult(mappings, &authLdapMappingPrintConfig)
}

func AuthLdapMappingUpdate(ctx context.Context, groupName string, opts AuthLdapMappingOptions) error {
	logger.Get().Info().Msgf("Updating LDAP mapping for group '%s'", groupName)

	authConfig, err := getAuthConfig(ctx)
	if err != nil {
		return err
	}

	ldapConfig := getLdapConfig(authConfig)
	if ldapConfig == nil {
		return fmt.Errorf("no LDAP configuration found")
	}

	if ldapConfig["groupsMapping"] == nil {
		return fmt.Errorf("no LDAP group mapping found")
	}

	matchFound := false
	for _, mapping := range ldapConfig["groupsMapping"].([]interface{}) {
		m := mapping.(map[string]interface{})
		if m["groupName"] == groupName {
			if opts.RoleName != "" {
				m["roleName"] = opts.RoleName
			}
			if opts.Priority != 0 {
				m["priority"] = float64(opts.Priority)
			}
			applyOptionalStringWithDefault(m, "userExternalIdentifier", opts.UserExternalIdentifier, defaultUserExternalIdentifier)
			applyOptionalStringWithDefault(m, "username", opts.Username, defaultUsername)
			applyOptionalStringWithDefault(m, "email", opts.Email, defaultEmail)
			matchFound = true
			break
		}
	}

	if !matchFound {
		return fmt.Errorf("no mapping found for LDAP group '%s'", groupName)
	}

	authConfig, err = patchAuthConfig(ctx, map[string]interface{}{"ldap": ldapConfig})
	if err != nil {
		return err
	}

	mappings := getLdapGroupMapping(getLdapConfig(authConfig))

	return formatter.PrintResult(mappings, &authLdapMappingPrintConfig)
}

func AuthLdapMappingRemove(ctx context.Context, groupName string) error {
	logger.Get().Info().Msgf("Removing LDAP mapping for group '%s'", groupName)

	authConfig, err := getAuthConfig(ctx)
	if err != nil {
		return err
	}

	ldapConfig := getLdapConfig(authConfig)
	if ldapConfig == nil {
		return fmt.Errorf("no LDAP configuration found")
	}

	if ldapConfig["groupsMapping"] == nil {
		return fmt.Errorf("no LDAP group mapping found")
	}

	matchFound := false
	newGroupsMapping := []interface{}{}
	for _, mapping := range ldapConfig["groupsMapping"].([]interface{}) {
		if mapping.(map[string]interface{})["groupName"] == groupName {
			matchFound = true
		} else {
			newGroupsMapping = append(newGroupsMapping, mapping)
		}
	}

	if !matchFound {
		return fmt.Errorf("no mapping found for LDAP group '%s'", groupName)
	}

	ldapConfig["groupsMapping"] = newGroupsMapping

	authConfig, err = patchAuthConfig(ctx, map[string]interface{}{"ldap": ldapConfig})
	if err != nil {
		return err
	}

	mappings := getLdapGroupMapping(getLdapConfig(authConfig))

	return formatter.PrintResult(mappings, &authLdapMappingPrintConfig)
}

// The v7.4 SDK models the LDAP mapping fields (groupName/roleName/priority and the
// profileMapping attributes) as map[string]interface{}, which cannot represent the
// API's scalar values — so typed (un)marshalling of the auth configuration fails on
// valid responses. We read and write the auth config as raw JSON to work around it.

func getAuthConfig(ctx context.Context) (map[string]interface{}, error) {
	httpRes, reqErr := authConfigRequest(ctx, http.MethodGet, "/api/v2/config?filter=auth", nil)
	config, err := parseAuthConfigResponse(httpRes, reqErr)
	if err != nil {
		return nil, err
	}

	if config == nil {
		logger.Get().Warn().Msg("No auth configuration found")
		return nil, nil
	}

	authConfig, ok := config["auth"].(map[string]interface{})
	if !ok || authConfig == nil {
		logger.Get().Warn().Msg("No auth configuration found")
		return nil, nil
	}

	return authConfig, nil
}

func patchAuthConfig(ctx context.Context, authConfigChange map[string]interface{}) (map[string]interface{}, error) {
	httpRes, reqErr := authConfigRequest(ctx, http.MethodPut, "/api/v2/config/auth", authConfigChange)
	return parseAuthConfigResponse(httpRes, reqErr)
}

// authConfigRequest performs a direct HTTP request against the configuration API,
// bypassing the SDK's typed configuration models (see note above).
func authConfigRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	client := api.GetApiClient(ctx)
	cfg := client.GetConfig()

	baseURL, err := cfg.ServerURL(0, nil)
	if err != nil {
		return nil, err
	}

	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, baseURL+path, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if auth, ok := ctx.Value(sdk.ContextAccessToken).(string); ok {
		req.Header.Set("Authorization", "Bearer "+auth)
	}

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return resp, err
	}

	// Re-wrap the body so response_inspector can read it for error formatting.
	respBody, readErr := io.ReadAll(resp.Body)
	resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewBuffer(respBody))
	if readErr != nil {
		return resp, readErr
	}

	return resp, nil
}

// parseAuthConfigResponse inspects the HTTP response for errors and parses its body
// into a generic map without SDK type unmarshalling.
func parseAuthConfigResponse(httpRes *http.Response, reqErr error) (map[string]interface{}, error) {
	if err := response_inspector.InspectResponse(httpRes, reqErr); err != nil {
		return nil, err
	}
	if httpRes == nil {
		return nil, fmt.Errorf("no response from server")
	}

	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	if len(body) == 0 {
		return nil, nil
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse auth configuration: %w", err)
	}

	return result, nil
}

func getLdapConfig(authConfig map[string]interface{}) map[string]interface{} {
	if authConfig == nil {
		return nil
	}

	ldapConfig, ok := authConfig["ldap"]
	if !ok || ldapConfig == nil {
		logger.Get().Warn().Msg("No LDAP configuration found")
		return nil
	}

	return ldapConfig.(map[string]interface{})
}

func getLdapGroupMapping(ldapConfig map[string]interface{}) []AuthLdapGroupMapping {
	groupsMapping, ok := ldapConfig["groupsMapping"]
	if !ok || groupsMapping == nil {
		logger.Get().Warn().Msg("No LDAP group mappings found")
		return nil
	}

	mappings := []AuthLdapGroupMapping{}
	for _, mapping := range groupsMapping.([]interface{}) {
		m := mapping.(map[string]interface{})
		groupMapping := AuthLdapGroupMapping{
			GroupName:              m["groupName"].(string),
			RoleName:               m["roleName"].(string),
			Priority:               int32(m["priority"].(float64)),
			UserExternalIdentifier: stringValue(m["userExternalIdentifier"]),
			Username:               stringValue(m["username"]),
			Email:                  stringValue(m["email"]),
		}
		mappings = append(mappings, groupMapping)
	}

	return mappings
}

func applyOptionalStringWithDefault(m map[string]interface{}, key string, value *string, defaultValue string) {
	if value == nil {
		return
	}
	if *value == "" {
		m[key] = defaultValue
		return
	}
	m[key] = *value
}

func stringValue(v interface{}) string {
	if v == nil {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		return ""
	}
	return s
}
