package template_asset

import (
	"context"
	"fmt"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var templateAssetPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"TemplateId": {
			Title: "Template ID",
			Order: 2,
		},
		"Usage": {
			Title: "Usage",
			Order: 3,
		},
		"File": {
			Hidden: true,
			InnerFields: map[string]formatter.RecordFieldConfig{
				"Name": {
					Title: "Filename",
					Order: 4,
				},
				"MimeType": {
					Title: "MIME Type",
					Order: 5,
				},
			},
		},
		"CreatedAt": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       6,
		},
		"ModifiedAt": {
			Title:       "Updated",
			Transformer: formatter.FormatDateTimeValue,
			Order:       7,
		},
	},
}

func TemplateAssetList(ctx context.Context, templateId []string, usage []string, mimeType []string) error {
	logger.Get().Info().Msgf("Listing template assets")

	client := api.GetApiClient(ctx)

	request := client.TemplateAssetAPI.GetTemplateAssets(ctx)

	// Apply filters if provided
	if len(templateId) > 0 {
		request = request.FilterTemplateId(utils.ProcessFilterStringSlice(templateId))
	}

	if len(usage) > 0 {
		request = request.FilterUsage(utils.ProcessFilterStringSlice(usage))
	}

	if len(mimeType) > 0 {
		request = request.FilterFileMimeType(utils.ProcessFilterStringSlice(mimeType))
	}

	request = request.SortBy([]string{"id:ASC"})

	templateAssetList, httpRes, err := request.Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(templateAssetList, &templateAssetPrintConfig)
}

func TemplateAssetGet(ctx context.Context, templateAssetId string) error {
	logger.Get().Info().Msgf("Get template asset %s details", templateAssetId)

	templateAssetIdNumeric, err := getTemplateAssetId(templateAssetId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	templateAsset, httpRes, err := client.TemplateAssetAPI.GetTemplateAsset(ctx, templateAssetIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(templateAsset, &templateAssetPrintConfig)
}

func TemplateAssetConfigExample(ctx context.Context) error {
	// Example template asset configuration
	templateAssetConfiguration := sdk.TemplateAssetCreate{
		TemplateId: 1,
		Usage:      "logo",
		File: sdk.TemplateAssetFile{
			Name:             "example-logo.png",
			MimeType:         "image/png",
			Checksum:         sdk.PtrString("example-checksum"),
			ContentBase64:    sdk.PtrString("base64-encoded-content-here"),
			TemplatingEngine: false,
			Url:              sdk.PtrString("https://example.com/logo.png"),
		},
		Tags: []string{"branding", "image"},
	}

	return formatter.PrintResult(templateAssetConfiguration, nil)
}

func TemplateAssetCreate(ctx context.Context, config []byte) error {
	logger.Get().Info().Msgf("Creating template asset")

	var templateAssetConfig sdk.TemplateAssetCreate
	err := utils.UnmarshalContent(config, &templateAssetConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	templateAssetInfo, httpRes, err := client.TemplateAssetAPI.CreateTemplateAsset(ctx).TemplateAssetCreate(templateAssetConfig).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(templateAssetInfo, &templateAssetPrintConfig)
}

func TemplateAssetUpdate(ctx context.Context, templateAssetId string, config []byte) error {
	logger.Get().Info().Msgf("Updating template asset %s", templateAssetId)

	templateAssetIdNumeric, err := getTemplateAssetId(templateAssetId)
	if err != nil {
		return err
	}

	var templateAssetConfig sdk.TemplateAssetCreate
	err = utils.UnmarshalContent(config, &templateAssetConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	templateAssetInfo, httpRes, err := client.TemplateAssetAPI.
		UpdateTemplateAsset(ctx, templateAssetIdNumeric).
		TemplateAssetCreate(templateAssetConfig).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(templateAssetInfo, &templateAssetPrintConfig)
}

func TemplateAssetDelete(ctx context.Context, templateAssetId string) error {
	logger.Get().Info().Msgf("Deleting template asset %s", templateAssetId)

	templateAssetIdNumeric, err := getTemplateAssetId(templateAssetId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.TemplateAssetAPI.
		DeleteTemplateAsset(ctx, templateAssetIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Template asset %s deleted", templateAssetId)
	return nil
}

func getTemplateAssetId(templateAssetId string) (float32, error) {
	templateAssetIdNumeric, err := strconv.ParseFloat(templateAssetId, 32)
	if err != nil {
		err := fmt.Errorf("invalid template asset ID: '%s'", templateAssetId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return float32(templateAssetIdNumeric), nil
}
