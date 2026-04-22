package auth

import (
	"context"
	"fmt"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
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

func getAuthConfig(ctx context.Context) (map[string]interface{}, error) {
	client := api.GetApiClient(ctx)

	configuration, httpRes, err := client.ConfigurationAPI.GetConfiguration(ctx).Filter("auth").Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, err
	}

	if configuration == nil {
		logger.Get().Warn().Msg("No configuration found")
		return nil, nil
	}

	authConfig, ok := configuration["auth"]
	if !ok {
		authConfig = configuration
	}

	if authConfig == nil {
		logger.Get().Warn().Msg("No auth configuration found")
		return nil, nil
	}

	return authConfig.(map[string]interface{}), nil
}

func patchAuthConfig(ctx context.Context, authConfigChange map[string]interface{}) (map[string]interface{}, error) {
	client := api.GetApiClient(ctx)

	authConfig, httpRes, err := client.ConfigurationAPI.PatchConfiguration(ctx, "auth").Body(authConfigChange).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, err
	}

	return authConfig, nil
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
