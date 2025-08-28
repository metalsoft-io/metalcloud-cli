package extension

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/object"
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
			if strings.Count(file.Name, "/") == 2 {
				parts := strings.Split(file.Name, "/")
				extensionPrefix := strings.Join(parts[:2], "/")

				if parts[2] == readMeFileName {
					return nil // Skip the README file
				}

				var err error
				if parts[2] == extensionFileName {
					extension := repoMap[extensionPrefix]

					extension.SourcePath = extensionPrefix
					extension.SourceContent, err = file.Contents()
					if err != nil {
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
	err := json.Unmarshal([]byte(repoExtension.SourceContent), &repoExtension.Extension)
	if err != nil {
		return err
	}

	if repoExtension.Extension.Name == "" {
		return fmt.Errorf("extension definition is missing or is older format in %s", repoExtension.SourcePath)
	}

	return nil
}
