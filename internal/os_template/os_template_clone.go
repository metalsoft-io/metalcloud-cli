package os_template

import (
	"context"
	"fmt"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

func OsTemplateClone(ctx context.Context, osTemplateId string, name string, label string) error {
	logger.Get().Info().Msgf("Cloning OS template %s", osTemplateId)

	osTemplate, err := GetOsTemplateByIdOrLabel(ctx, osTemplateId)
	if err != nil {
		return err
	}

	templateCreate := convertOSTemplateToCreate(osTemplate)

	// Override name and label if provided, otherwise append " (clone)"
	if name != "" {
		templateCreate.Name = name
	} else {
		templateCreate.Name = osTemplate.Name + " (clone)"
	}

	if label != "" {
		templateCreate.Label = sdk.PtrString(label)
	} else {
		templateCreate.Label = sdk.PtrString(utils.CreateSlug(templateCreate.Name))
	}

	templateCreate.Visibility = sdk.PtrString("private")

	client := api.GetApiClient(ctx)

	// List assets for this template
	templateAssetList, httpRes, err := client.TemplateAssetAPI.
		GetTemplateAssets(ctx).
		FilterTemplateId([]string{"$eq:" + fmt.Sprintf("%d", osTemplate.Id)}).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return fmt.Errorf("failed to list template assets: %w", err)
	}

	var assetCreates []sdk.TemplateAssetCreate

	for _, asset := range templateAssetList.Data {
		assetCreate := convertTemplateAssetToCreate(&asset)

		if asset.File.Url != nil && *asset.File.Url != "" {
			// URL-based asset: preserve as-is
			logger.Get().Info().Msgf("Asset '%s' is URL-based, preserving URL reference", asset.File.Name)
		} else {
			// Content-based asset: fetch full content
			fullAsset, httpRes, err := client.TemplateAssetAPI.
				GetTemplateAsset(ctx, float32(asset.Id)).
				Execute()
			if err := response_inspector.InspectResponse(httpRes, err); err != nil {
				return fmt.Errorf("failed to get asset %d content: %w", asset.Id, err)
			}

			if fullAsset.File.ContentBase64 == nil || *fullAsset.File.ContentBase64 == "" {
				return fmt.Errorf("content asset '%s' (ID %d) has no content - cannot clone", asset.File.Name, asset.Id)
			}

			assetCreate.File.ContentBase64 = fullAsset.File.ContentBase64
			assetCreate.File.Checksum = fullAsset.File.Checksum
		}

		assetCreates = append(assetCreates, assetCreate)
	}

	createOptions := OsTemplateCreateOptions{
		Template:       templateCreate,
		TemplateAssets: assetCreates,
	}

	return OsTemplateCreate(ctx, createOptions)
}
