package os_template

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
)

// OsTemplateListDirectory lists the OS templates stored in a local directory. It is the offline
// counterpart of OsTemplateListRepo, intended for air-gapped environments where the templates are
// copied to a directory instead of being published through a Git repository.
func OsTemplateListDirectory(ctx context.Context, directory string) error {
	logger.Get().Info().Msgf("Listing all OS templates from directory %s", directory)

	directoryAssets, err := getDirectoryTemplateAssets(directory)
	if err != nil {
		return err
	}

	return printRepositoryTemplates(directoryAssets)
}

// OsTemplateCreateFromDirectory creates a new OS template from a template stored in a local directory.
// It is the offline counterpart of OsTemplateCreateFromRepo.
func OsTemplateCreateFromDirectory(ctx context.Context, sourceTemplate string, directory string, name string, label string, sourceIso string) error {
	logger.Get().Info().Msgf("Creating OS template %s from directory %s", sourceTemplate, directory)

	directoryAssets, err := getDirectoryTemplateAssets(directory)
	if err != nil {
		return err
	}

	return createOsTemplateFromAssets(ctx, directoryAssets, sourceTemplate, fmt.Sprintf("directory %s", directory), name, label, sourceIso)
}

func isLocalDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// getDirectoryTemplateAssets reads the templates stored in a local directory. The directory must use
// the same layout as the template repositories - <vendor>/<os>/<version>/template.yaml along with the
// asset files referenced by it.
func getDirectoryTemplateAssets(dirPath string) (map[string]RepositoryTemplateInfo, error) {
	if !isLocalDirectory(dirPath) {
		return nil, fmt.Errorf("%s is not an accessible local directory", dirPath)
	}

	repoMap := make(map[string]RepositoryTemplateInfo)

	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(dirPath, path)
		if err != nil {
			return err
		}

		// Normalize to forward slashes for consistency
		relPath = filepath.ToSlash(relPath)

		// Only process files at depth 3 (vendor/os/version/filename)
		if strings.Count(relPath, "/") != 3 {
			return nil
		}

		parts := strings.Split(relPath, "/")
		templatePrefix := strings.Join(parts[:3], "/")

		if parts[3] == readMeFileName {
			return nil
		}

		if _, ok := repoMap[templatePrefix]; !ok {
			repoMap[templatePrefix] = RepositoryTemplateInfo{
				Assets: make(map[string]RepositoryTemplateAsset),
			}
		}

		template := repoMap[templatePrefix]

		fileContent, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		if parts[3] == templateFileName {
			template.SourcePath = templatePrefix
			template.SourceContent = string(fileContent)
		} else {
			template.Assets[parts[3]] = RepositoryTemplateAsset{
				ContentBase64: base64.StdEncoding.EncodeToString(fileContent),
			}
		}

		repoMap[templatePrefix] = template
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk local directory %s: %w", dirPath, err)
	}

	return repoMap, nil
}
