package os_template

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
	"gopkg.in/yaml.v3"
)

func OsTemplateImport(ctx context.Context, archivePath string, name string, label string) error {
	logger.Get().Info().Msgf("Importing OS template from %s", archivePath)

	// Extract archive to temp directory
	tmpDir, err := os.MkdirTemp("", "os-template-import-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	if err := extractZip(archivePath, tmpDir); err != nil {
		return fmt.Errorf("failed to extract archive: %w", err)
	}

	// Read and unmarshal template.yaml
	yamlPath := filepath.Join(tmpDir, templateFileName)
	yamlData, err := os.ReadFile(yamlPath)
	if err != nil {
		return fmt.Errorf("failed to read template.yaml: %w", err)
	}

	var createOptions OsTemplateCreateOptions
	if err := yaml.Unmarshal(yamlData, &createOptions); err != nil {
		return fmt.Errorf("failed to parse template.yaml: %w", err)
	}

	// Override name and label
	createOptions.Template.Name = name
	if label != "" {
		createOptions.Template.Label = sdk.PtrString(label)
	} else {
		createOptions.Template.Label = sdk.PtrString(utils.CreateSlug(name))
	}

	// Set visibility to private
	createOptions.Template.Visibility = sdk.PtrString("private")

	// Process content assets: read files and encode
	assetsDir := filepath.Join(tmpDir, "assets")
	for i, asset := range createOptions.TemplateAssets {
		if asset.File.Url != nil && *asset.File.Url != "" {
			// URL-based asset: nothing to process
			continue
		}

		// Content-based asset: read from assets directory
		assetFilePath := filepath.Join(assetsDir, asset.File.Name)
		fileData, err := os.ReadFile(assetFilePath)
		if err != nil {
			return fmt.Errorf("failed to read asset file '%s': %w", asset.File.Name, err)
		}

		contentBase64 := base64.StdEncoding.EncodeToString(fileData)
		checksum := fmt.Sprintf("%x", sha256.Sum256([]byte(contentBase64)))

		createOptions.TemplateAssets[i].File.ContentBase64 = sdk.PtrString(contentBase64)
		createOptions.TemplateAssets[i].File.Checksum = sdk.PtrString(checksum)

		logger.Get().Info().Msgf("Loaded asset '%s' (%d bytes)", asset.File.Name, len(fileData))
	}

	return OsTemplateCreate(ctx, createOptions)
}
