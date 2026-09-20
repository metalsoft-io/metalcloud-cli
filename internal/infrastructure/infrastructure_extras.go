package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var utilizationSummaryPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Kind": {
			Title: "Kind",
			Order: 1,
		},
		"InfrastructureId": {
			Title: "Infrastructure",
			Order: 2,
		},
		"InfrastructureLabel": {
			Title:    "Label",
			MaxWidth: 30,
			Order:    3,
		},
		"ServiceStatus": {
			Title:       "Status",
			Transformer: formatter.FormatStatusValue,
			Order:       4,
		},
		"Resource": {
			Title:    "Resource",
			MaxWidth: 40,
			Order:    5,
		},
		"Quantity": {
			Title: "Quantity",
			Order: 6,
		},
		"MeasurementUnit": {
			Title: "Unit",
			Order: 7,
		},
	},
}

// utilizationSummaryRecord is the table projection of the summarized
// utilization report: the response is a set of nested maps which the tabular
// formatter renders as empty cells, so it is flattened into one row per
// infrastructure and one row per metered resource. The json and yaml formats
// keep the full response.
type utilizationSummaryRecord struct {
	Kind                string
	InfrastructureId    interface{}
	InfrastructureLabel string
	ServiceStatus       string
	Resource            string
	Quantity            interface{}
	MeasurementUnit     string
}

// InfrastructureUpdateMetadata updates the metadata (name, description, tags)
// of an infrastructure from a JSON/YAML configuration.
//
// The SDK request exposes no If-Match setter for this endpoint, so no entity
// tag is sent - unlike the configuration update, metadata is not under
// optimistic concurrency control.
func InfrastructureUpdateMetadata(ctx context.Context, infrastructureIdOrLabel string, config []byte) error {
	logger.Get().Info().Msgf("Update metadata of infrastructure '%s'", infrastructureIdOrLabel)

	infrastructureInfo, err := GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	var updateMeta sdk.UpdateInfrastructureMeta
	if err := utils.UnmarshalContent(config, &updateMeta); err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	updated, httpRes, err := client.InfrastructureAPI.
		UpdateInfrastructureMetadata(ctx, int64(infrastructureInfo.Id)).
		UpdateInfrastructureMeta(updateMeta).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(updated, &infrastructurePrintConfig)
}

// InfrastructureGetUtilizationSummary retrieves the summarized resource
// utilization report for a user over a time range.
//
// The response body is parsed raw because the typed SDK
// InfrastructureResourceUtilizationSummaryResponse model is out of sync with
// the API: it declares infrastructures[].tags as a string, while the API
// returns a list, so decoding a live payload fails with "json: cannot
// unmarshal array into Go struct field ...infrastructures.tags of type
// string". The request body is still built from the typed SDK model.
func InfrastructureGetUtilizationSummary(ctx context.Context, userId int, startTime time.Time, endTime time.Time, infrastructureIds []int) error {
	logger.Get().Info().Msgf("Getting utilization summary for user %d from %s to %s",
		userId, startTime.Format(time.RFC3339), endTime.Format(time.RFC3339))

	request := sdk.GetResourceUtilizationSummarized{
		UserIdOwner:    float32(userId),
		StartTimestamp: startTime.Format(time.RFC3339),
		EndTimestamp:   endTime.Format(time.RFC3339),
	}

	if len(infrastructureIds) > 0 {
		request.InfrastructureIds = make([]int64, 0, len(infrastructureIds))
		for _, infrastructureId := range infrastructureIds {
			request.InfrastructureIds = append(request.InfrastructureIds, int64(infrastructureId))
		}
	}

	// The request body is always sent: an unset body would be serialised as a
	// literal `null`, which the API rejects with 400.
	body, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to encode utilization summary request: %w", err)
	}

	responseBody, err := api.RawJSONRequest(
		ctx,
		http.MethodPost,
		"/api/v2/infrastructures/actions/get/resource-utilization-summarized",
		body,
		nil,
	)
	if err != nil {
		return err
	}

	if formatter.IsNativeFormat() {
		return utils.PrintRawObject(responseBody, nil)
	}

	summary, err := utils.DecodeRawObject(responseBody)
	if err != nil {
		return err
	}

	return formatter.PrintResult(buildUtilizationSummaryRecords(summary), &utilizationSummaryPrintConfig)
}

// buildUtilizationSummaryRecords flattens the raw summary response into table
// rows: one row per infrastructure, one per metered resource and one per
// internet direction.
func buildUtilizationSummaryRecords(summary map[string]any) []utilizationSummaryRecord {
	records := []utilizationSummaryRecord{}

	if infrastructures, ok := summary["infrastructures"].(map[string]any); ok {
		for _, key := range sortedKeys(infrastructures) {
			infra, ok := infrastructures[key].(map[string]any)
			if !ok {
				continue
			}
			records = append(records, utilizationSummaryRecord{
				Kind:                "Infrastructure",
				InfrastructureId:    infra["infrastructureId"],
				InfrastructureLabel: rawString(infra, "infrastructureLabel"),
				ServiceStatus:       rawString(infra, "infrastructureServiceStatus"),
			})
		}
	}

	if resources, ok := summary["resourceUtilization"].(map[string]any); ok {
		for _, key := range sortedKeys(resources) {
			records = append(records, utilizationRecord("Resource", key, resources[key]))
		}
	}

	if internet, ok := summary["internet"].(map[string]any); ok {
		for _, direction := range []string{"upload", "download"} {
			if _, present := internet[direction]; present {
				records = append(records, utilizationRecord("Internet", direction, internet[direction]))
			}
		}
	}

	return records
}

// utilizationRecord builds one metered-resource row out of a raw utilization
// object.
func utilizationRecord(kind string, name string, value any) utilizationSummaryRecord {
	record := utilizationSummaryRecord{Kind: kind, Resource: name}

	if utilization, ok := value.(map[string]any); ok {
		record.Quantity = utilization["quantity"]
		record.MeasurementUnit = rawString(utilization, "measurementUnit")
	}

	return record
}

func sortedKeys(object map[string]any) []string {
	keys := make([]string, 0, len(object))
	for key := range object {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func rawString(object map[string]any, key string) string {
	value, _ := object[key].(string)
	return value
}
