package os_template

import (
	"context"
	"fmt"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
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

var osTemplateCredentialsPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Username": {
			Title: "Username",
			Order: 1,
		},
		"Password": {
			Title: "Password",
			Order: 2,
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

type OsTemplateCreateOptions struct {
	Template       sdk.OSTemplateCreate      `json:"template"`
	TemplateAssets []sdk.TemplateAssetCreate `json:"templateAssets"`
}

func OsTemplateExampleCreate(ctx context.Context) error {
	osTemplateCreateOptions := OsTemplateCreateOptions{
		Template: sdk.OSTemplateCreate{
			Name:        "OS Template Name",
			Label:       sdk.PtrString("os-template-label - optional"),
			Description: sdk.PtrString("OS template description. - optional"),
			Device: sdk.OSTemplateDevice{
				Type:         "server",
				BootMode:     "uefi",
				Architecture: "x86_64",
			},
			Install: sdk.OSTemplateInstall{
				Method:      "oob",
				DriveType:   "local_drive",
				ReadyMethod: "wait_for_power_off",
				OnieStrings: []string{
					"tempor officia elit proident",
					"magna v",
				},
			},
			ImageBuild: sdk.OSTemplateImageBuild{
				Required: true,
				Provider: sdk.PtrString("xorriso - optional"),
			},
			Os: sdk.OSTemplateOs{
				Name:    "Ubuntu",
				Version: "22.04",
				Credential: sdk.OSTemplateOsCredential{
					Username:     "root",
					PasswordType: "plain",
					Password:     sdk.PtrString("rqi|password - optional"),
				},
				SshPort: sdk.PtrInt32(22),
			},
			Visibility: sdk.PtrString("public"),
			Tags: []string{
				"tag1",
				"tag2",
			},
		},
		TemplateAssets: []sdk.TemplateAssetCreate{
			{
				TemplateId: 0, // This will be set after the template is created
				Usage:      "build_source_image",
				File: sdk.TemplateAssetFile{
					Name:             "name.iso",
					MimeType:         "application/octet-stream",
					TemplatingEngine: false,
					Url:              sdk.PtrString("http://my.repo.com/test.iso"),
					Path:             "/name.iso",
				},
				Tags: []string{
					"tag1",
					"tag2",
				},
			},
			{
				TemplateId: 0, // This will be set after the template is created
				Usage:      "build_component",
				File: sdk.TemplateAssetFile{
					Name:             "name.xml",
					MimeType:         "text/plain",
					Checksum:         sdk.PtrString("checksum - optional"),
					ContentBase64:    sdk.PtrString("contentBase64 - optional"),
					TemplatingEngine: true,
					Path:             "/name.xml",
				},
				Tags: []string{
					"tag1",
					"tag2",
				},
			},
		},
	}

	return formatter.PrintResult(osTemplateCreateOptions, nil)
}

func OsTemplateCreate(ctx context.Context, osTemplateCreateOptions OsTemplateCreateOptions) error {
	logger.Get().Info().Msgf("Creating OS template")

	client := api.GetApiClient(ctx)

	osTemplate, httpRes, err := client.OSTemplateAPI.CreateOSTemplate(ctx).OSTemplateCreate(osTemplateCreateOptions.Template).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}
	logger.Get().Info().Msgf("Template %d created", osTemplate.Id)

	if osTemplateCreateOptions.TemplateAssets != nil {
		for _, asset := range osTemplateCreateOptions.TemplateAssets {
			asset.TemplateId = osTemplate.Id

			newAsset, httpRes, err := client.TemplateAssetAPI.CreateTemplateAsset(ctx).TemplateAssetCreate(asset).Execute()
			if err := response_inspector.InspectResponse(httpRes, err); err != nil {
				return err
			}
			logger.Get().Info().Msgf("Template asset %d created", newAsset.Id)
		}
	}

	return formatter.PrintResult(osTemplate, &osTemplatePrintConfig)
}

type OsTemplateUpdateOptions struct {
	Template                *sdk.OSTemplateUpdate             `json:"template"`
	NewTemplateAssets       []sdk.TemplateAssetCreate         `json:"newTemplateAssets"`
	UpdatedTemplateAssets   map[int32]sdk.TemplateAssetCreate `json:"updatedTemplateAssets"`
	DeletedTemplateAssetIds []int32                           `json:"deletedTemplateAssetIds"`
}

func OsTemplateUpdate(ctx context.Context, osTemplateId string, osTemplateUpdateOptions OsTemplateUpdateOptions) error {
	logger.Get().Info().Msgf("Updating OS template %s", osTemplateId)

	osTemplateIdNumeric, revision, err := getOsTemplateIdAndRevision(ctx, osTemplateId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	if osTemplateUpdateOptions.Template != nil {
		_, httpRes, err := client.OSTemplateAPI.
			UpdateOSTemplate(ctx, osTemplateIdNumeric).
			OSTemplateUpdate(*osTemplateUpdateOptions.Template).
			IfMatch(revision).
			Execute()
		if err := response_inspector.InspectResponse(httpRes, err); err != nil {
			return err
		}
	}

	if osTemplateUpdateOptions.NewTemplateAssets != nil {
		for _, asset := range osTemplateUpdateOptions.NewTemplateAssets {
			asset.TemplateId = int32(osTemplateIdNumeric)

			newAsset, httpRes, err := client.TemplateAssetAPI.CreateTemplateAsset(ctx).TemplateAssetCreate(asset).Execute()
			if err := response_inspector.InspectResponse(httpRes, err); err != nil {
				return err
			}
			logger.Get().Info().Msgf("Template asset %d created", newAsset.Id)
		}
	}

	if osTemplateUpdateOptions.UpdatedTemplateAssets != nil {
		for assetId, asset := range osTemplateUpdateOptions.UpdatedTemplateAssets {
			asset.TemplateId = int32(osTemplateIdNumeric)

			_, httpRes, err := client.TemplateAssetAPI.
				UpdateTemplateAsset(ctx, float32(assetId)).
				TemplateAssetCreate(asset).
				Execute()
			if err := response_inspector.InspectResponse(httpRes, err); err != nil {
				return err
			}
			logger.Get().Info().Msgf("Template asset %d updated", assetId)
		}
	}

	if osTemplateUpdateOptions.DeletedTemplateAssetIds != nil {
		for _, assetId := range osTemplateUpdateOptions.DeletedTemplateAssetIds {
			httpRes, err := client.TemplateAssetAPI.
				DeleteTemplateAsset(ctx, float32(assetId)).
				Execute()
			if err := response_inspector.InspectResponse(httpRes, err); err != nil {
				return err
			}
			logger.Get().Info().Msgf("Template asset %d deleted", assetId)
		}
	}

	return nil
}

func OsTemplateDelete(ctx context.Context, osTemplateId string) error {
	logger.Get().Info().Msgf("Deleting OS template %s", osTemplateId)

	osTemplateIdNumeric, err := getOsTemplateId(osTemplateId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.OSTemplateAPI.
		DeleteOSTemplate(ctx, osTemplateIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("OS template %s deleted", osTemplateId)
	return nil
}

func OsTemplateGetCredentials(ctx context.Context, osTemplateId string) error {
	logger.Get().Info().Msgf("Getting credentials for OS template %s", osTemplateId)

	osTemplateIdNumeric, err := getOsTemplateId(osTemplateId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	credentials, httpRes, err := client.OSTemplateAPI.
		GetOSTemplateCredentials(ctx, osTemplateIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(credentials, &osTemplateCredentialsPrintConfig)
}

func OsTemplateGetAssets(ctx context.Context, osTemplateId string) error {
	logger.Get().Info().Msgf("Getting assets for OS template %s", osTemplateId)

	osTemplateIdNumeric, err := getOsTemplateId(osTemplateId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	// Use the template asset API with filter for this template ID
	templateAssetList, httpRes, err := client.TemplateAssetAPI.
		GetTemplateAssets(ctx).
		FilterTemplateId([]string{"$eq:" + fmt.Sprintf("%d", int32(osTemplateIdNumeric))}).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(templateAssetList, &formatter.PrintConfig{
		FieldsConfig: map[string]formatter.RecordFieldConfig{
			"Id": {
				Title: "#",
				Order: 1,
			},
			"Usage": {
				Title: "Usage",
				Order: 2,
			},
			"File": {
				Hidden: true,
				InnerFields: map[string]formatter.RecordFieldConfig{
					"Name": {
						Title: "Filename",
						Order: 3,
					},
					"MimeType": {
						Title: "MIME Type",
						Order: 4,
					},
					"Size": {
						Title: "Size",
						Order: 5,
					},
				},
			},
			"CreatedAt": {
				Title:       "Created",
				Transformer: formatter.FormatDateTimeValue,
				Order:       6,
			},
		},
	})
}

func GetOsTemplateByIdOrLabel(ctx context.Context, osTemplateIdOrLabel string) (*sdk.OSTemplate, error) {
	client := api.GetApiClient(ctx)

	osTemplateId, err := getOsTemplateId(osTemplateIdOrLabel)
	if err != nil {
		return nil, err
	}

	osTemplateInfo, httpRes, err := client.OSTemplateAPI.GetOSTemplate(ctx, osTemplateId).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, err
	}

	return osTemplateInfo, nil
}

func getOsTemplateId(osTemplateId string) (float32, error) {
	osTemplateIdNumeric, err := strconv.ParseFloat(osTemplateId, 32)
	if err != nil {
		err := fmt.Errorf("invalid OS template ID: '%s'", osTemplateId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return float32(osTemplateIdNumeric), nil
}

func getOsTemplateIdAndRevision(ctx context.Context, osTemplateId string) (float32, string, error) {
	osTemplateIdNumeric, err := getOsTemplateId(osTemplateId)
	if err != nil {
		return 0, "", err
	}

	client := api.GetApiClient(ctx)

	osTemplate, httpRes, err := client.OSTemplateAPI.GetOSTemplate(ctx, float32(osTemplateIdNumeric)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return 0, "", err
	}

	return float32(osTemplateIdNumeric), strconv.Itoa(int(osTemplate.Revision)), nil
}
