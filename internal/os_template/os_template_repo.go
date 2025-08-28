package os_template

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/metalsoft-io/metalcloud-cli/pkg/repo"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
	"gopkg.in/yaml.v3"
)

const (
	publicRepositoryURL = "https://github.com/metalsoft-io/os-templates.git"
	templateFileName    = "template.yaml"
	readMeFileName      = "README.md"
)

type RepositoryTemplateInfo struct {
	SourcePath    string
	SourceContent string
	OsTemplate    OsTemplateCreateOptions
	Assets        map[string]RepositoryTemplateAsset
}

type RepositoryTemplateAsset struct {
	ContentBase64 string
}

func cloneOsTemplateRepository(ctx context.Context, repoUrl string, repoUsername string, repoPassword string) (*object.Tree, error) {
	return repo.CloneRepository(ctx, repoUrl, repoUsername, repoPassword, publicRepositoryURL)
}

func getRepositoryTemplateAssets(tree *object.Tree) map[string]RepositoryTemplateInfo {
	repoMap := make(map[string]RepositoryTemplateInfo)

	files := tree.Files()
	files.ForEach(func(file *object.File) error {
		if file.Mode.IsRegular() {
			if strings.Count(file.Name, "/") == 3 {
				parts := strings.Split(file.Name, "/")
				templatePrefix := strings.Join(parts[:3], "/")

				if parts[3] == readMeFileName {
					return nil // Skip the README file
				}

				if _, ok := repoMap[templatePrefix]; !ok {
					repoMap[templatePrefix] = RepositoryTemplateInfo{
						Assets: make(map[string]RepositoryTemplateAsset),
					}
				}

				template := repoMap[templatePrefix]

				var err error
				if parts[3] == templateFileName {
					template.SourcePath = templatePrefix
					template.SourceContent, err = file.Contents()
					if err != nil {
						return err
					}
				} else {
					isBinary, err := file.IsBinary()
					if err != nil {
						return err
					}

					if isBinary {
						reader, err := file.Reader()
						if err != nil {
							return err
						}

						buf := new(bytes.Buffer)
						_, err = buf.ReadFrom(reader)
						if err != nil {
							return err
						}

						template.Assets[parts[3]] = RepositoryTemplateAsset{
							ContentBase64: base64.StdEncoding.EncodeToString(buf.Bytes()),
						}
					} else {
						fileContent, err := file.Contents()
						if err != nil {
							return err
						}

						template.Assets[parts[3]] = RepositoryTemplateAsset{
							ContentBase64: base64.StdEncoding.EncodeToString([]byte(fileContent)),
						}
					}
				}

				repoMap[templatePrefix] = template
			}
		}

		return nil
	})

	return repoMap
}

func processTemplateContent(repoTemplate *RepositoryTemplateInfo) error {
	err := yaml.Unmarshal([]byte(repoTemplate.SourceContent), &repoTemplate.OsTemplate)
	if err != nil {
		return err
	}

	if repoTemplate.OsTemplate.Template.Name == "" {
		return fmt.Errorf("template definition is missing or is older format in %s", repoTemplate.SourcePath)
	}

	for i, asset := range repoTemplate.OsTemplate.TemplateAssets {
		if asset.File.Url != nil && *asset.File.Url != "" {
			// If the asset has a URL, we don't need to process it further
			continue
		}

		if _, ok := repoTemplate.Assets[asset.File.Name]; ok {
			// If the asset is already in the repository, we add its content and checksum
			repoTemplate.OsTemplate.TemplateAssets[i].File.ContentBase64 = sdk.PtrString(repoTemplate.Assets[asset.File.Name].ContentBase64)
			checksum := fmt.Sprintf("%x", sha256.Sum256([]byte(repoTemplate.Assets[asset.File.Name].ContentBase64)))
			repoTemplate.OsTemplate.TemplateAssets[i].File.Checksum = sdk.PtrString(checksum)
		}
	}

	return nil
}
