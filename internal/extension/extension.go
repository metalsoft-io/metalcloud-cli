package extension

import (
	"context"
	"fmt"
	"slices"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
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

	return formatter.PrintResult(*extension, &extensionPrintConfig)
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
