package variable

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

// variableRaw avoids SDK unmarshal failure: the SDK's typed Variable model types
// `value` as map[string]interface{}, but the API may return a plain string (or
// other JSON scalar), failing deserialization of the whole response. Decoding
// `value` into interface{} keeps `variable list`/`variable get` tolerant of the
// SDK <-> API type mismatch.
type variableRaw struct {
	Id               interface{} `json:"id"`
	UserIdOwner      interface{} `json:"userIdOwner"`
	Name             *string     `json:"name"`
	Value            interface{} `json:"value"`
	Usage            *string     `json:"usage"`
	CreatedTimestamp *string     `json:"createdTimestamp"`
	UpdatedTimestamp *string     `json:"updatedTimestamp"`
}

var variablePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"Name": {
			MaxWidth: 30,
			Order:    2,
		},
		"Value": {
			MaxWidth: 50,
			Order:    3,
		},
		"Usage": {
			Title: "Usage Type",
			Order: 4,
		},
		"UserIdOwner": {
			Title: "Owner ID",
			Order: 5,
		},
		"CreatedTimestamp": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       6,
		},
		"UpdatedTimestamp": {
			Title:       "Updated",
			Transformer: formatter.FormatDateTimeValue,
			Order:       7,
		},
	},
}

func VariableList(ctx context.Context) error {
	logger.Get().Info().Msgf("Listing all variables")

	client := api.GetApiClient(ctx)

	request := client.VariablesAPI.GetVariables(ctx).SortBy([]string{"id:ASC"})

	// Raw-body parse: the typed Variable model types `value` as a map, but the
	// API may return a string, so SDK unmarshalling fails on valid responses.
	rawItems, meta, err := utils.FetchAllPagesRaw(func(page float32) (*http.Response, error) {
		_, httpRes, _ := request.Page(page).Limit(100).Execute()
		return httpRes, nil
	})
	if err != nil {
		return err
	}

	records, err := utils.UnmarshalRawItems[variableRaw](rawItems)
	if err != nil {
		return fmt.Errorf("failed to parse variables: %w", err)
	}

	return utils.PrintAllRaw(rawItems, records, meta, len(records), &variablePrintConfig)
}

func VariableGet(ctx context.Context, variableId string) error {
	logger.Get().Info().Msgf("Get variable '%s' details", variableId)

	variableIdNumeric, err := getVariableId(variableId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	// Raw-body parse: see VariableList — the `value` type mismatch breaks typed decoding.
	_, httpRes, sdkErr := client.VariablesAPI.GetVariable(ctx, variableIdNumeric).Execute()
	if httpRes != nil && httpRes.StatusCode >= 400 {
		return response_inspector.InspectResponse(httpRes, sdkErr)
	}
	if httpRes == nil {
		return sdkErr
	}

	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var record variableRaw
	if err := json.Unmarshal(body, &record); err != nil {
		return fmt.Errorf("failed to parse variable: %w", err)
	}

	return formatter.PrintResult(record, &variablePrintConfig)
}

// variableConfigToJSON normalizes a create/update config (JSON or YAML) into a
// JSON body, bypassing the SDK's CreateVariable/UpdateVariable types whose `value`
// field is a `oneOf` union (VariableValue) that fails to round-trip some valid
// JSON values (e.g. an empty object). Passing a plain map preserves any value shape.
func variableConfigToJSON(config []byte) ([]byte, error) {
	var payload map[string]interface{}
	if err := utils.UnmarshalContent(config, &payload); err != nil {
		return nil, err
	}
	return json.Marshal(payload)
}

// parseAndPrintVariable parses a single-variable raw API response into the
// tolerant variableRaw struct and prints it, sidestepping the VariableValue union.
func parseAndPrintVariable(httpRes *http.Response, reqErr error) error {
	if httpRes == nil {
		return reqErr
	}
	defer httpRes.Body.Close()

	if httpRes.StatusCode >= 400 {
		return response_inspector.InspectResponse(httpRes, reqErr)
	}

	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var record variableRaw
	if err := json.Unmarshal(body, &record); err != nil {
		return fmt.Errorf("failed to parse variable: %w", err)
	}

	return formatter.PrintResult(record, &variablePrintConfig)
}

func VariableCreate(ctx context.Context, config []byte) error {
	logger.Get().Info().Msgf("Creating variable")

	body, err := variableConfigToJSON(config)
	if err != nil {
		return err
	}

	httpRes, reqErr := api.DoJSONRequest(ctx, http.MethodPost, "/api/v2/variables", body)
	return parseAndPrintVariable(httpRes, reqErr)
}

func VariableUpdate(ctx context.Context, variableId string, config []byte) error {
	logger.Get().Info().Msgf("Updating variable '%s'", variableId)

	// Validate the ID is numeric, then use the original string in the path to
	// avoid float precision loss for large IDs.
	if _, err := getVariableId(variableId); err != nil {
		return err
	}

	body, err := variableConfigToJSON(config)
	if err != nil {
		return err
	}

	httpRes, reqErr := api.DoJSONRequest(ctx, http.MethodPut, "/api/v2/variables/"+variableId, body)
	return parseAndPrintVariable(httpRes, reqErr)
}

func VariableDelete(ctx context.Context, variableId string) error {
	logger.Get().Info().Msgf("Deleting variable '%s'", variableId)

	variableIdNumeric, err := getVariableId(variableId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.VariablesAPI.
		DeleteVariable(ctx, variableIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Variable '%s' deleted", variableId)
	return nil
}

func VariableConfigExample(ctx context.Context) error {
	// Example create variable configuration
	variableConfiguration := sdk.CreateVariable{
		Name:  "example-variable",
		Value: sdk.VariableValue{String: sdk.PtrString("value1")},
		Usage: (*sdk.VariableUsageType)(sdk.PtrString("general")),
	}

	return formatter.PrintResult(variableConfiguration, nil)
}

func getVariableId(variableId string) (float32, error) {
	variableIdNumeric, err := strconv.ParseFloat(variableId, 32)
	if err != nil {
		err := fmt.Errorf("invalid variable ID: '%s'", variableId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return float32(variableIdNumeric), nil
}
