package os_template

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
	"gopkg.in/yaml.v3"
)

func OsTemplateExport(ctx context.Context, osTemplateId string, outputPath string) error {
	logger.Get().Info().Msgf("Exporting OS template %s", osTemplateId)

	osTemplate, err := GetOsTemplateByIdOrLabel(ctx, osTemplateId)
	if err != nil {
		return err
	}

	templateCreate := convertOSTemplateToCreate(osTemplate)

	client := api.GetApiClient(ctx)

	// List assets for this template
	templateAssetList, httpRes, err := client.TemplateAssetAPI.
		GetTemplateAssets(ctx).
		FilterTemplateId([]string{"$eq:" + fmt.Sprintf("%d", osTemplate.Id)}).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return fmt.Errorf("failed to list template assets: %w", err)
	}

	// Create temp directory for building the archive
	tmpDir, err := os.MkdirTemp("", "os-template-export-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	assetsDir := filepath.Join(tmpDir, "assets")
	if err := os.MkdirAll(assetsDir, 0755); err != nil {
		return fmt.Errorf("failed to create assets directory: %w", err)
	}

	var assetCreates []sdk.TemplateAssetCreate

	for _, asset := range templateAssetList.Data {
		assetCreate := convertTemplateAssetToCreate(&asset)

		if asset.File.Url != nil && *asset.File.Url != "" {
			// URL-based asset (e.g., ISO link): preserve URL, no file to save
			logger.Get().Info().Msgf("Asset '%s' is URL-based, preserving URL reference", asset.File.Name)
		} else {
			// Content-based asset: fetch full content via individual GET
			fullAsset, httpRes, err := client.TemplateAssetAPI.
				GetTemplateAsset(ctx, float32(asset.Id)).
				Execute()
			if err := response_inspector.InspectResponse(httpRes, err); err != nil {
				return fmt.Errorf("failed to get asset %d content: %w", asset.Id, err)
			}

			if fullAsset.File.ContentBase64 == nil || *fullAsset.File.ContentBase64 == "" {
				return fmt.Errorf("content asset '%s' (ID %d) has no content - cannot export", asset.File.Name, asset.Id)
			}

			// Decode and write content to file
			decoded, err := base64.StdEncoding.DecodeString(*fullAsset.File.ContentBase64)
			if err != nil {
				return fmt.Errorf("failed to decode content for asset '%s': %w", asset.File.Name, err)
			}

			assetFilePath := filepath.Join(assetsDir, asset.File.Name)
			if err := os.WriteFile(assetFilePath, decoded, 0644); err != nil {
				return fmt.Errorf("failed to write asset file '%s': %w", asset.File.Name, err)
			}

			logger.Get().Info().Msgf("Exported asset '%s' (%d bytes)", asset.File.Name, len(decoded))

			// Clear content and checksum from metadata (will be recomputed on import)
			assetCreate.File.ContentBase64 = nil
			assetCreate.File.Checksum = nil
		}

		assetCreates = append(assetCreates, assetCreate)
	}

	// Build the OsTemplateCreateOptions
	createOptions := OsTemplateCreateOptions{
		Template:       templateCreate,
		TemplateAssets: assetCreates,
	}

	// Serialize to YAML
	yamlData, err := yaml.Marshal(createOptions)
	if err != nil {
		return fmt.Errorf("failed to marshal template to YAML: %w", err)
	}

	yamlPath := filepath.Join(tmpDir, templateFileName)
	if err := os.WriteFile(yamlPath, yamlData, 0644); err != nil {
		return fmt.Errorf("failed to write template.yaml: %w", err)
	}

	// Determine output path
	if outputPath == "" {
		timestamp := time.Now().Format("20060102150405")
		outputPath = fmt.Sprintf("%d_%s_%s.zip", osTemplate.Id, utils.CreateSlug(osTemplate.Name), timestamp)
	}

	// Create the zip archive
	if err := createZip(tmpDir, outputPath); err != nil {
		return fmt.Errorf("failed to create archive: %w", err)
	}

	logger.Get().Info().Msgf("Template exported to %s", outputPath)
	fmt.Printf("Template '%s' exported to %s\n", osTemplate.Name, outputPath)

	return nil
}

func convertOSTemplateToCreate(t *sdk.OSTemplate) sdk.OSTemplateCreate {
	create := sdk.OSTemplateCreate{
		Name:                  t.Name,
		Description:           t.Description,
		Label:                 t.Label,
		Device:                t.Device,
		Install:               t.Install,
		ImageCertSerialNumber: t.ImageCertSerialNumber,
		Os:                    t.Os,
		Visibility:            &t.Visibility,
		Tags:                  t.Tags,
	}

	if t.ImageBuild != nil {
		create.ImageBuild = *t.ImageBuild
	}

	return create
}

func convertTemplateAssetToCreate(a *sdk.TemplateAsset) sdk.TemplateAssetCreate {
	return sdk.TemplateAssetCreate{
		TemplateId: 0, // Will be set after template creation on import
		Usage:      a.Usage,
		File:       a.File,
		Tags:       a.Tags,
	}
}
