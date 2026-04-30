package extension

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/repo"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

const (
	publicRepositoryURL = "https://github.com/metalsoft-io/metalsoft-extensions.git"
	extensionFileName   = "extension.json"
	readMeFileName      = "README.md"
)

type RepositoryExtensionInfo struct {
	SourcePath    string
	SourceContent string
	Extension     sdk.CreateExtension
}

func cloneExtensionRepository(ctx context.Context, repoUrl string, repoUsername string, repoPassword string) (*object.Tree, error) {
	return repo.CloneRepository(ctx, repoUrl, repoUsername, repoPassword, publicRepositoryURL)
}

func getRepositoryExtensions(tree *object.Tree) map[string]RepositoryExtensionInfo {
	repoMap := make(map[string]RepositoryExtensionInfo)

	files := tree.Files()
	files.ForEach(func(file *object.File) error {
		if file.Mode.IsRegular() {
			depthLevel := strings.Count(file.Name, "/")
			if depthLevel > 0 && depthLevel <= 2 {
				parts := strings.Split(file.Name, "/")
				extensionPrefix := strings.Join(parts[:depthLevel], "/")

				if parts[depthLevel] == readMeFileName {
					return nil // Skip the README file
				}

				var err error
				if parts[depthLevel] == extensionFileName {
					extension := repoMap[extensionPrefix]

					extension.SourcePath = extensionPrefix
					extension.SourceContent, err = file.Contents()
					if err != nil {
						logger.Get().Trace().Msgf("Repo item %s skipped - reading the content failed: %s", file.Name, err.Error())
						return err
					}

					repoMap[extensionPrefix] = extension
				}
			}
		}

		return nil
	})

	return repoMap
}

func processExtensionContent(repoExtension *RepositoryExtensionInfo) error {
	var extensionDefinition sdk.ExtensionDefinition
	err := json.Unmarshal([]byte(repoExtension.SourceContent), &extensionDefinition)
	if err != nil {
		return err
	}

	if extensionDefinition.Name == "" {
		return fmt.Errorf("extension definition is missing or is older format in %s", repoExtension.SourcePath)
	}

	repoExtension.Extension = sdk.CreateExtension{
		Kind:        extensionDefinition.ExtensionType,
		Label:       &extensionDefinition.Label,
		Name:        extensionDefinition.Name,
		Description: *extensionDefinition.Description,
		Definition:  extensionDefinition,
	}

	return nil
}
