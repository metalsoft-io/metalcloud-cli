package role

import (
	"context"
	"encoding/json"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

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

	roles, httpRes, err := client.SecurityAPI.GetRoles(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(roles, &rolePrintConfig)
}

func Get(ctx context.Context, roleName string) error {
	logger.Get().Info().Msgf("Getting role '%s'", roleName)

	client := api.GetApiClient(ctx)

	role, httpRes, err := client.SecurityAPI.GetRole(ctx, roleName).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(role, &rolePrintConfig)
}

func Create(ctx context.Context, config []byte) error {
	logger.Get().Info().Msg("Creating role")

	var createRole sdk.CreateRole
	if err := json.Unmarshal(config, &createRole); err != nil {
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

	role, httpRes, err := client.SecurityAPI.DeleteRole(ctx, roleName).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Role '%s' deleted", roleName)
	return formatter.PrintResult(role, &rolePrintConfig)
}

func Update(ctx context.Context, roleName string, config []byte) error {
	logger.Get().Info().Msgf("Updating role '%s'", roleName)

	var editRole sdk.EditRole
	if err := json.Unmarshal(config, &editRole); err != nil {
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
