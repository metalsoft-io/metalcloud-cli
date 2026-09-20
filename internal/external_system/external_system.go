// Package external_system manages the external systems registered with the
// platform (/api/v2/external-systems).
package external_system

import (
	"context"
	"fmt"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var externalSystemPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"Label": {
			Title:    "Label",
			Order:    2,
			MaxWidth: 30,
		},
		"Name": {
			Title:    "Name",
			Order:    3,
			MaxWidth: 40,
		},
		"Revision": {
			Title: "Revision",
			Order: 4,
		},
		"CreatedAt": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       5,
		},
		"UpdatedAt": {
			Title:       "Updated",
			Transformer: formatter.FormatDateTimeValue,
			Order:       6,
		},
	},
}

func ExternalSystemList(ctx context.Context, filterLabel []string) error {
	logger.Get().Info().Msgf("Listing external systems")

	client := api.GetApiClient(ctx)

	request := client.ExternalSystemAPI.GetExternalSystems(ctx).SortBy([]string{"id:ASC"})
	if len(filterLabel) > 0 {
		request = request.FilterLabel(utils.ProcessFilterStringSlice(filterLabel))
	}

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &externalSystemPrintConfig)
}

func ExternalSystemGet(ctx context.Context, externalSystemId string) error {
	logger.Get().Info().Msgf("Get external system '%s'", externalSystemId)

	externalSystem, err := GetExternalSystemById(ctx, externalSystemId)
	if err != nil {
		return err
	}

	return formatter.PrintResult(externalSystem, &externalSystemPrintConfig)
}

func ExternalSystemConfigExample(ctx context.Context) error {
	example := sdk.CreateExternalSystem{
		Label: "my-external-system",
		Name:  "My External System",
		Annotations: map[string]interface{}{
			"owner": "platform-team",
		},
	}

	if formatter.IsNativeFormat() {
		return formatter.PrintResult(example, nil)
	}
	return formatter.PrintYamlResult(example)
}

func ExternalSystemCreate(ctx context.Context, create sdk.CreateExternalSystem) error {
	logger.Get().Info().Msgf("Creating external system '%s'", create.Label)

	client := api.GetApiClient(ctx)

	externalSystem, httpRes, err := client.ExternalSystemAPI.
		CreateExternalSystem(ctx).
		CreateExternalSystem(create).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(externalSystem, &externalSystemPrintConfig)
}

func ExternalSystemUpdate(ctx context.Context, externalSystemId string, config []byte) error {
	logger.Get().Info().Msgf("Updating external system '%s'", externalSystemId)

	var update sdk.UpdateExternalSystem
	if err := utils.UnmarshalContent(config, &update); err != nil {
		return err
	}

	// PATCH requires an If-Match entity tag, so the current revision is read first.
	current, err := GetExternalSystemById(ctx, externalSystemId)
	if err != nil {
		return err
	}

	numericId, err := externalSystemNumericId(current.Id)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	externalSystem, httpRes, err := client.ExternalSystemAPI.
		UpdateExternalSystem(ctx, numericId).
		UpdateExternalSystem(update).
		IfMatch(current.Revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(externalSystem, &externalSystemPrintConfig)
}

func ExternalSystemDelete(ctx context.Context, externalSystemId string) error {
	logger.Get().Info().Msgf("Deleting external system '%s'", externalSystemId)

	numericId, err := externalSystemNumericId(externalSystemId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	// DeleteExternalSystem exposes no IfMatch setter, so no revision is sent.
	httpRes, err := client.ExternalSystemAPI.DeleteExternalSystem(ctx, numericId).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("External system '%s' deleted", externalSystemId)
	return nil
}

// GetExternalSystemById fetches an external system by its numeric ID. The
// model carries the ID as a string while the endpoints take an int64, so the
// value is converted on the way in.
func GetExternalSystemById(ctx context.Context, externalSystemId string) (*sdk.ExternalSystem, error) {
	numericId, err := externalSystemNumericId(externalSystemId)
	if err != nil {
		return nil, err
	}

	client := api.GetApiClient(ctx)

	externalSystem, httpRes, err := client.ExternalSystemAPI.GetExternalSystemById(ctx, numericId).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, err
	}

	if externalSystem == nil {
		err := fmt.Errorf("external system '%s' not found", externalSystemId)
		logger.Get().Error().Err(err).Msg("")
		return nil, err
	}

	return externalSystem, nil
}

func externalSystemNumericId(externalSystemId string) (int64, error) {
	numericId, err := utils.GetInt64FromString(externalSystemId)
	if err != nil {
		err = fmt.Errorf("invalid external system ID: '%s'", externalSystemId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}
	return numericId, nil
}
