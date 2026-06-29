package role

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

// roleRaw avoids SDK unmarshal failure: the SDK's typed Role model decodes
// `permissions` into the strict MetalsoftPermissions enum, which rejects
// permission values it doesn't know (e.g. network_profile_allowed_for_user_read).
// Parsing the raw body into plain string fields keeps `role list`/`role get`
// tolerant of SDK <-> API permission-enum drift.
type roleRaw struct {
	Id             interface{} `json:"id"`
	Name           *string     `json:"name"`
	Label          *string     `json:"label"`
	Description    *string     `json:"description"`
	Type           *string     `json:"type"`
	Permissions    []string    `json:"permissions"`
	QuotaProfileId interface{} `json:"quotaProfileId"`
	UsersWithRole  *float32    `json:"usersWithRole"`
}

var rolePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Order: 1,
		},
		"Label": {
			Order: 2,
		},
		"Name": {
			MaxWidth: 30,
			Order:    3,
		},
		"Description": {
			MaxWidth: 50,
			Order:    4,
		},
		"UsersWithRole": {
			Title: "Users",
			Order: 5,
		},
	},
}

func List(ctx context.Context) error {
	logger.Get().Info().Msg("Listing all roles")

	client := api.GetApiClient(ctx)

	// Raw-body parse: the typed Role model rejects unknown permission enum
	// values, so SDK unmarshalling fails on otherwise-valid responses.
	_, httpRes, sdkErr := client.SecurityAPI.GetRoles(ctx).Execute()
	if httpRes == nil {
		return sdkErr
	}

	rawItems, meta, err := utils.ParseRawPage(httpRes)
	if err != nil {
		return err
	}

	records, err := utils.UnmarshalRawItems[roleRaw](rawItems)
	if err != nil {
		return fmt.Errorf("failed to parse roles: %w", err)
	}

	return utils.PrintAllRaw(rawItems, records, meta, len(records), &rolePrintConfig)
}

func Get(ctx context.Context, roleName string) error {
	logger.Get().Info().Msgf("Getting role '%s'", roleName)

	client := api.GetApiClient(ctx)

	// Raw-body parse: see List — the MetalsoftPermissions enum drift breaks typed decoding.
	_, httpRes, sdkErr := client.SecurityAPI.GetRole(ctx, roleName).Execute()
	if httpRes != nil && httpRes.StatusCode >= 400 {
		return response_inspector.InspectResponse(httpRes, sdkErr)
	}
	if httpRes == nil {
		return sdkErr
	}

	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var record roleRaw
	if err := json.Unmarshal(body, &record); err != nil {
		return fmt.Errorf("failed to parse role: %w", err)
	}

	return formatter.PrintResult(record, &rolePrintConfig)
}

func Create(ctx context.Context, config []byte) error {
	logger.Get().Info().Msg("Creating role")

	var createRole sdk.CreateRole
	if err := utils.UnmarshalContent(config, &createRole); err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	role, httpRes, err := client.SecurityAPI.CreateRole(ctx).CreateRole(createRole).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(role, &rolePrintConfig)
}

func Delete(ctx context.Context, roleName string) error {
	logger.Get().Info().Msgf("Deleting role '%s'", roleName)

	client := api.GetApiClient(ctx)

	httpRes, err := client.SecurityAPI.DeleteRole(ctx, roleName).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Role '%s' deleted", roleName)
	return nil
}

func Update(ctx context.Context, roleName string, config []byte) error {
	logger.Get().Info().Msgf("Updating role '%s'", roleName)

	var editRole sdk.EditRole
	if err := utils.UnmarshalContent(config, &editRole); err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	role, httpRes, err := client.SecurityAPI.UpdateRole(ctx, roleName).EditRole(editRole).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Role '%s' updated", roleName)
	return formatter.PrintResult(role, &rolePrintConfig)
}
