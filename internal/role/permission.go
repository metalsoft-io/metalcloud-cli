package role

import (
	"context"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
)

var permissionPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Order: 1,
		},
		"Label": {
			Order: 2,
		},
		"Name": {
			Order: 3,
		},
		"Description": {
			MaxWidth: 60,
			Order:    4,
		},
	},
}

func ListPermissions(ctx context.Context) error {
	logger.Get().Info().Msg("Listing all permissions")

	client := api.GetApiClient(ctx)

	permissions, httpRes, err := client.SecurityAPI.GetPermissions(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(permissions.Permissions, &permissionPrintConfig)
}
