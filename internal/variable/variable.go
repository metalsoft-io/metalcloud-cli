package variable

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

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

	variableList, httpRes, err := client.VariablesAPI.GetVariables(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(variableList, &variablePrintConfig)
}

func VariableGet(ctx context.Context, variableId string) error {
	logger.Get().Info().Msgf("Get variable '%s' details", variableId)

	variableIdNumeric, err := getVariableId(variableId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	variable, httpRes, err := client.VariablesAPI.GetVariable(ctx, variableIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(variable, &variablePrintConfig)
}

func VariableCreate(ctx context.Context, config []byte) error {
	logger.Get().Info().Msgf("Creating variable")

	var variableConfig sdk.CreateVariable
	err := json.Unmarshal(config, &variableConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	variable, httpRes, err := client.VariablesAPI.
		CreateVariable(ctx).
		CreateVariable(variableConfig).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(variable, &variablePrintConfig)
}

func VariableUpdate(ctx context.Context, variableId string, config []byte) error {
	logger.Get().Info().Msgf("Updating variable '%s'", variableId)

	variableIdNumeric, err := getVariableId(variableId)
	if err != nil {
		return err
	}

	var variableConfig sdk.UpdateVariable
	err = json.Unmarshal(config, &variableConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	variable, httpRes, err := client.VariablesAPI.
		UpdateVariable(ctx, variableIdNumeric).
		UpdateVariable(variableConfig).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(variable, &variablePrintConfig)
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
		Name: "example-variable",
		Value: map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
		},
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
