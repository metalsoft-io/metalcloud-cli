package template

import (
	"context"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var templatePrintConfig = formatter.PrintConfig{
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

func TemplateList(ctx context.Context) error {
	logger.Get().Info().Msgf("Listing all templates")

	client := api.GetApiClient(ctx)

	templateList, httpRes, err := client.OSTemplateAPI.GetOSTemplates(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(templateList, &templatePrintConfig)
}

func TemplateGet(ctx context.Context, templateId string) error {
	logger.Get().Info().Msgf("Get template %s details", templateId)

	template, err := GetTemplateByIdOrLabel(ctx, templateId)
	if err != nil {
		return err
	}

	return formatter.PrintResult(template, &templatePrintConfig)
}

func GetTemplateByIdOrLabel(ctx context.Context, templateIdOrLabel string) (*sdk.OSTemplate, error) {
	client := api.GetApiClient(ctx)

	templateId, err := utils.GetFloat32FromString(templateIdOrLabel)
	if err != nil {
		return nil, err
	}

	templateInfo, httpRes, err := client.OSTemplateAPI.GetOSTemplate(ctx, templateId).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, err
	}

	return templateInfo, nil
}
