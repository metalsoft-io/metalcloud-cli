package os_template

import (
	"context"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var osTemplatePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"Name": {
			MaxWidth: 30,
			Order:    2,
		},
		"Label": {
			MaxWidth: 30,
			Order:    3,
		},
		"Device": {
			Hidden: true,
			InnerFields: map[string]formatter.RecordFieldConfig{
				"Type": {
					Title: "Device Type",
					Order: 4,
				},
			},
		},
		"Status": {
			Title:       "Status",
			Transformer: formatter.FormatStatusValue,
			Order:       5,
		},
		"Visibility": {
			Order: 6,
		},
		"CreatedAt": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       7,
		},
		"ModifiedAt": {
			Title:       "Updated",
			Transformer: formatter.FormatDateTimeValue,
			Order:       8,
		},
	},
}

func OsTemplateList(ctx context.Context) error {
	logger.Get().Info().Msgf("Listing all OS templates")

	client := api.GetApiClient(ctx)

	osTemplateList, httpRes, err := client.OSTemplateAPI.GetOSTemplates(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(osTemplateList, &osTemplatePrintConfig)
}

func OsTemplateGet(ctx context.Context, osTemplateId string) error {
	logger.Get().Info().Msgf("Get OS template %s details", osTemplateId)

	osTemplate, err := GetOsTemplateByIdOrLabel(ctx, osTemplateId)
	if err != nil {
		return err
	}

	return formatter.PrintResult(osTemplate, &osTemplatePrintConfig)
}

func GetOsTemplateByIdOrLabel(ctx context.Context, osTemplateIdOrLabel string) (*sdk.OSTemplate, error) {
	client := api.GetApiClient(ctx)

	osTemplateId, err := utils.GetFloat32FromString(osTemplateIdOrLabel)
	if err != nil {
		return nil, err
	}

	osTemplateInfo, httpRes, err := client.OSTemplateAPI.GetOSTemplate(ctx, osTemplateId).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, err
	}

	return osTemplateInfo, nil
}
