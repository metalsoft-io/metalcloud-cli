package extension

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
	"github.com/spf13/viper"
)

var extensionPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"Label": {
			MaxWidth: 20,
			Order:    2,
		},
		"Name": {
			MaxWidth: 30,
			Order:    3,
		},
		"Status": {
			Title:       "Status",
			Transformer: formatter.FormatStatusValue,
			Order:       4,
		},
		"Kind": {
			Order: 5,
		},
		"Description": {
			MaxWidth: 50,
			Order:    6,
		},
	},
}

func ExtensionList(ctx context.Context, filterLabel []string, filterName []string, filterStatus []string, filterKind []string, filterPublic string) error {
	logger.Get().Info().Msgf("Listing extensions")

	client := api.GetApiClient(ctx)
	request := client.ExtensionAPI.GetExtensions(ctx).SortBy([]string{"id:ASC"})

	// Apply filters if provided
	if len(filterLabel) > 0 {
		request = request.FilterLabel(utils.ProcessFilterStringSlice(filterLabel))
	}

	if len(filterName) > 0 {
		request = request.FilterName(utils.ProcessFilterStringSlice(filterName))
	}

	if len(filterStatus) > 0 {
		request = request.FilterStatus(utils.ProcessFilterStringSlice(filterStatus))
	}

	// if len(filterKind) > 0 {
	// 	request = request.FilterKind(utils.ProcessFilterStringSlice(filterKind))
	// }

	if filterPublic != "" {
		request = request.FilterIsPublic([]string{filterPublic})
	}

	extensions := make([]sdk.ExtensionInfo, 0)

	page := float32(1)

	// Loop through all pages and collect extensions list
	for {
		request = request.Page(page)

		extensionList, httpRes, err := request.Execute()
		if err := response_inspector.InspectResponse(httpRes, err); err != nil {
			return err
		}

		extensions = append(extensions, extensionList.Data...)

		if *extensionList.Meta.TotalPages <= *extensionList.Meta.CurrentPage {
			break // No more pages to process
		}

		page++
	}

	// Workaround until the API supports this filter - filter out the extensions by kind, if needed
	if len(filterKind) > 0 {
		filteredExtensions := make([]sdk.ExtensionInfo, 0)
		for _, ext := range extensions {
			if slices.Contains(filterKind, *ext.Kind) {
				filteredExtensions = append(filteredExtensions, ext)
			}
		}
		extensions = filteredExtensions
	}

	return formatter.PrintResult(extensions, &extensionPrintConfig)
}

func ExtensionGet(ctx context.Context, extensionId string) error {
	logger.Get().Info().Msgf("Get extension '%s'", extensionId)

	extension, err := GetExtensionByIdOrLabel(ctx, extensionId)
	if err != nil {
		return err
	}

	if strings.ToLower(viper.GetString(formatter.ConfigFormat)) == "text" {
		// If the output format is text, print the basic information followed by the inputs
		err := formatter.PrintResult(*extension, &extensionPrintConfig)
		if err != nil {
			return err
		}

		if len(extension.Definition.Inputs) > 0 {
			err := formatter.PrintResult(toExtensionInputs(extension.Definition.Inputs), &formatter.PrintConfig{
				FieldsConfig: map[string]formatter.RecordFieldConfig{
					"Label": {
						Title: "Input Label",
						Order: 1,
					},
					"Name": {
						Title: "Input Name",
						Order: 2,
					},
					"InputType": {
						Title: "Input Type",
						Order: 3,
					},
					"DefaultValue": {
						Title: "Default Value",
						Order: 4,
					},
				},
			})
			if err != nil {
				return fmt.Errorf("failed to print inputs: %w", err)
			}
		}

		return nil
	} else {
		return formatter.PrintResult(*extension, &extensionPrintConfig)
	}
}

func toExtensionInputs(dataItems []sdk.ExtensionDefinitionInputsDataItem) []extensionInput {
	result := []extensionInput{}

	for _, dataItem := range dataItems {
		result = append(result, toExtensionInput(dataItem))
	}

	return result
}

func valueOf(defaultValue *sdk.ExtensionInputStringDefaultValue) any {
	if defaultValue != nil {
		if defaultValue.Bool != nil {
			return *defaultValue.Bool
		}
		if defaultValue.Int32 != nil {
			return *defaultValue.Int32
		}
		if defaultValue.String != nil {
			return *defaultValue.String
		}
	}
	return nil
}

type extensionInput struct {
	Label        string
	Name         string
	InputType    sdk.ExtensionInputType
	DefaultValue any
}

func toExtensionInput(dataItem sdk.ExtensionDefinitionInputsDataItem) extensionInput {
	var extensionInput extensionInput

	if dataItem.ExtensionInputBoolean != nil {
		extensionInput.Label = dataItem.ExtensionInputBoolean.Label
		extensionInput.Name = dataItem.ExtensionInputBoolean.Name
		extensionInput.InputType = dataItem.ExtensionInputBoolean.InputType
		extensionInput.DefaultValue = valueOf(dataItem.ExtensionInputBoolean.DefaultValue)
	} else if dataItem.ExtensionInputInteger != nil {
		extensionInput.Label = dataItem.ExtensionInputInteger.Label
		extensionInput.Name = dataItem.ExtensionInputInteger.Name
		extensionInput.InputType = dataItem.ExtensionInputInteger.InputType
		extensionInput.DefaultValue = valueOf(dataItem.ExtensionInputInteger.DefaultValue)
	} else if dataItem.ExtensionInputString != nil {
		extensionInput.Label = dataItem.ExtensionInputString.Label
		extensionInput.Name = dataItem.ExtensionInputString.Name
		extensionInput.InputType = dataItem.ExtensionInputString.InputType
		extensionInput.DefaultValue = valueOf(dataItem.ExtensionInputString.DefaultValue)
	} else if dataItem.ExtensionInputServerType != nil {
		extensionInput.Label = dataItem.ExtensionInputServerType.Label
		extensionInput.Name = dataItem.ExtensionInputServerType.Name
		extensionInput.InputType = dataItem.ExtensionInputServerType.InputType
		extensionInput.DefaultValue = valueOf(dataItem.ExtensionInputServerType.DefaultValue)
	} else if dataItem.ExtensionInputOsTemplate != nil {
		extensionInput.Label = dataItem.ExtensionInputOsTemplate.Label
		extensionInput.Name = dataItem.ExtensionInputOsTemplate.Name
		extensionInput.InputType = dataItem.ExtensionInputOsTemplate.InputType
		extensionInput.DefaultValue = valueOf(dataItem.ExtensionInputOsTemplate.DefaultValue)
	}

	return extensionInput
}

func ExtensionCreate(ctx context.Context, name string, kind string, description string, config []byte) error {
	logger.Get().Info().Msgf("Create extension '%s'", name)

	var definition sdk.ExtensionDefinition
	err := utils.UnmarshalContent(config, &definition)
	if err != nil {
		return fmt.Errorf("invalid definition JSON: %v", err)
	}

	label := "" // Optional
	createExtension := sdk.CreateExtension{
		Name:        name,
		Kind:        kind,
		Description: description,
		Definition:  definition,
	}

	if len(label) > 0 {
		createExtension.Label = &label
	}

	client := api.GetApiClient(ctx)

	extensionInfo, httpRes, err := client.ExtensionAPI.CreateExtension(ctx).CreateExtension(createExtension).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(extensionInfo, &extensionPrintConfig)
}

func ExtensionUpdate(ctx context.Context, extensionId string, name string, description string, config []byte) error {
	logger.Get().Info().Msgf("Update extension '%s'", extensionId)

	extension, err := GetExtensionByIdOrLabel(ctx, extensionId)
	if err != nil {
		return err
	}

	var definition sdk.ExtensionDefinition
	if len(config) > 0 {
		err = utils.UnmarshalContent(config, &definition)
		if err != nil {
			return err
		}
	} else {
		definition = extension.Definition
	}

	if name == "" {
		name = extension.Name
	}

	if description == "" {
		description = extension.Description
	}

	updateExtension := sdk.UpdateExtension{
		Name:        name,
		Description: description,
		Definition:  definition,
	}

	if extension.Label != nil {
		updateExtension.Label = extension.Label
	}

	client := api.GetApiClient(ctx)

	updatedExtension, httpRes, err := client.ExtensionAPI.UpdateExtension(ctx, extension.Id).
		UpdateExtension(updateExtension).
		IfMatch(fmt.Sprintf("%.0f", extension.Revision)).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(updatedExtension, &extensionPrintConfig)
}

func ExtensionListRepo(ctx context.Context, repoUrl string, repoUsername string, repoPassword string) error {
	logger.Get().Info().Msgf("Listing extensions from repository")

	tree, err := cloneExtensionRepository(ctx, repoUrl, repoUsername, repoPassword)
	if err != nil {
		return fmt.Errorf("failed to clone OS template repository: %w", err)
	}

	// This map stores all files for an extension and will be used to check if their information is correct
	repoMap := make(map[string]RepositoryExtensionInfo)
	for templatePrefix, repoTemplate := range getRepositoryExtensions(tree) {
		err = processExtensionContent(&repoTemplate)
		if err != nil {
			// Ignore extension with errors - they may be using old format
			logger.Get().Warn().Msgf("Ignoring extension %s - error processing its content: %v", templatePrefix, err)
			continue
		}

		repoMap[templatePrefix] = repoTemplate
	}

	// Convert the map to slice for printing
	repoExtensionsSlice := make([]RepositoryExtensionInfo, 0, len(repoMap))
	for _, repoExtension := range repoMap {
		repoExtensionsSlice = append(repoExtensionsSlice, repoExtension)
	}

	// Order the extensions by SourcePath
	slices.SortStableFunc(repoExtensionsSlice, func(a, b RepositoryExtensionInfo) int {
		if a.SourcePath < b.SourcePath {
			return -1
		} else if a.SourcePath > b.SourcePath {
			return 1
		}
		return 0
	})

	return formatter.PrintResult(repoExtensionsSlice, &formatter.PrintConfig{
		FieldsConfig: map[string]formatter.RecordFieldConfig{
			"SourcePath": {
				Title: "Path",
				Order: 1,
			},
			"Extension": {
				Hidden: true,
				InnerFields: map[string]formatter.RecordFieldConfig{
					"Name": {
						Order: 2,
					},
					"Label": {
						Order: 3,
					},
					"Kind": {
						Order: 4,
					},
				},
			},
		},
	})
}

func ExtensionCreateFromRepo(ctx context.Context, extensionPath string, repoUrl string, repoUsername string, repoPassword string, name string, label string) error {
	logger.Get().Info().Msgf("Creating extension from repository path '%s'", extensionPath)

	tree, err := cloneExtensionRepository(ctx, repoUrl, repoUsername, repoPassword)
	if err != nil {
		return fmt.Errorf("failed to clone extension repository: %w", err)
	}

	repoMap := getRepositoryExtensions(tree)

	extension, ok := repoMap[extensionPath]
	if !ok {
		return fmt.Errorf("extension %s not found in repository", extensionPath)
	}

	err = processExtensionContent(&extension)
	if err != nil {
		return fmt.Errorf("error processing extension content: %w", err)
	}

	if name != "" {
		extension.Extension.Name = name
	}
	if label != "" {
		extension.Extension.Label = sdk.PtrString(label)
	}

	client := api.GetApiClient(ctx)

	extensionInfo, httpRes, err := client.ExtensionAPI.CreateExtension(ctx).CreateExtension(extension.Extension).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(extensionInfo, &extensionPrintConfig)
}

func ExtensionPublish(ctx context.Context, extensionId string) error {
	logger.Get().Info().Msgf("Publishing extension '%s'", extensionId)

	extension, err := GetExtensionByIdOrLabel(ctx, extensionId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.ExtensionAPI.PublishExtension(ctx, extension.Id).
		IfMatch(fmt.Sprintf("%.0f", extension.Revision)).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Extension '%s' published successfully", extensionId)
	return nil
}

func ExtensionArchive(ctx context.Context, extensionId string) error {
	logger.Get().Info().Msgf("Archiving extension '%s'", extensionId)

	extension, err := GetExtensionByIdOrLabel(ctx, extensionId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.ExtensionAPI.ArchiveExtension(ctx, extension.Id).
		IfMatch(fmt.Sprintf("%.0f", extension.Revision)).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Extension '%s' archived successfully", extensionId)
	return nil
}

func GetExtensionByIdOrLabel(ctx context.Context, extensionIdOrLabel string) (*sdk.Extension, error) {
	client := api.GetApiClient(ctx)

	// First try to get by ID
	extensionIdFloat, err := strconv.ParseFloat(extensionIdOrLabel, 32)
	if err == nil {
		extensionInfo, httpRes, err := client.ExtensionAPI.GetExtension(ctx, float32(extensionIdFloat)).Execute()
		logger.Get().Info().Msgf("Extension '%s' get by ID:\n err: %v\n httpRes: %v\n StatusCode: %v", extensionIdOrLabel, err, httpRes, httpRes.StatusCode)

		if err == nil && httpRes != nil && httpRes.StatusCode == 200 {
			return extensionInfo, nil
		}
		/* RM: IMO
		if err != nil && httpRes != nil && httpRes.StatusCode != 400 {
			return nil, err
		}
		*/
		// If not found by ID, continue to search by label
	}

	// Try to get by label
	extensions, httpRes, err := client.ExtensionAPI.GetExtensions(ctx).
		FilterLabel([]string{extensionIdOrLabel}).
		Execute()
	if err = response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, err
	}

	if len(extensions.Data) == 0 {
		err := fmt.Errorf("extension '%s' not found", extensionIdOrLabel)
		logger.Get().Error().Err(err).Msg("")
		return nil, err
	}

	extensionInfo, httpRes, err := client.ExtensionAPI.GetExtension(ctx, *extensions.Data[0].Id).Execute()
	if err = response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, err
	}

	return extensionInfo, nil
}
