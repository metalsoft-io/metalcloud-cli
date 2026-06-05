package permission

import (
	"context"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var permissionPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"Name": {
			Title: "Name",
			Order: 2,
		},
		"Label": {
			Title: "Label",
			Order: 3,
		},
		"Type": {
			Title: "Type",
			Order: 4,
		},
		"Description": {
			Title: "Description",
			Order: 5,
		},
	},
}

func PermissionList(ctx context.Context) error {
	logger.Get().Info().Msgf("Listing permissions")

	client := api.GetApiClient(ctx)

	permissionList, httpRes, err := client.SecurityAPI.GetPermissions(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(permissionList, &permissionPrintConfig)
}

func PermissionCreate(ctx context.Context, config []byte) error {
	logger.Get().Info().Msgf("Creating permission")

	var createConfig sdk.CreatePermission
	err := utils.UnmarshalContent(config, &createConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	permission, httpRes, err := client.SecurityAPI.CreatePermission(ctx).CreatePermission(createConfig).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(permission, &permissionPrintConfig)
}

func PermissionDelete(ctx context.Context, permissionName string) error {
	logger.Get().Info().Msgf("Deleting permission '%s'", permissionName)

	client := api.GetApiClient(ctx)

	httpRes, err := client.SecurityAPI.DeletePermission(ctx, permissionName).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Permission '%s' deleted successfully", permissionName)
	return nil
}

func PermissionConfigExample(ctx context.Context) error {
	permissionConfig := sdk.CreatePermission{
		Name:        "custom_read",
		Label:       "Custom Read",
		Description: sdk.PtrString("Example custom permission"),
	}

	return formatter.PrintResult(permissionConfig, nil)
}
