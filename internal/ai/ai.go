// Package ai wraps the platform AI assistant (Eli) endpoint,
// POST /api/v2/ai/generate.
package ai

import (
	"context"
	"fmt"
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var aiResponsePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Result": {
			Title: "Result",
			Order: 1,
		},
		"Steps": {
			Title: "Steps",
			Order: 2,
		},
	},
}

// generateRequestInput mirrors sdk.AIGenerateRequest but without its
// required-property validation, so a --config-source file may carry only the
// prompt and let --datacenter supply the rest.
type generateRequestInput struct {
	Datacenter *string `json:"datacenter,omitempty"`
	Prompt     *string `json:"prompt,omitempty"`
}

// AIGenerate asks the platform AI assistant for a response to prompt.
// datacenter scopes the question to a datacenter and is required by the API.
func AIGenerate(ctx context.Context, datacenter string, prompt string) error {
	logger.Get().Info().Msgf("Generating AI response for datacenter '%s'", datacenter)

	request, err := buildGenerateRequest(datacenter, prompt)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	// The body is always set: the SDK sends a literal "null" when the optional
	// body setter is skipped, which the API rejects with 400.
	response, httpRes, err := client.AIAPI.
		GenerateEliResponse(ctx).
		AIGenerateRequest(*request).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	if response == nil {
		err := fmt.Errorf("the AI endpoint returned no response")
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	return printGenerateResponse(response)
}

// AIGenerateFromConfig reads the request from a JSON/YAML document. Values not
// present in the document fall back to the datacenter and prompt arguments.
func AIGenerateFromConfig(ctx context.Context, datacenter string, config []byte) error {
	var input generateRequestInput
	if err := utils.UnmarshalContent(config, &input); err != nil {
		return err
	}

	if input.Datacenter != nil && *input.Datacenter != "" {
		datacenter = *input.Datacenter
	}

	prompt := ""
	if input.Prompt != nil {
		prompt = *input.Prompt
	}

	return AIGenerate(ctx, datacenter, prompt)
}

func buildGenerateRequest(datacenter string, prompt string) (*sdk.AIGenerateRequest, error) {
	prompt = strings.TrimSpace(prompt)
	if prompt == "" {
		err := fmt.Errorf("the AI prompt cannot be empty")
		logger.Get().Error().Err(err).Msg("")
		return nil, err
	}

	datacenter = strings.TrimSpace(datacenter)
	if datacenter == "" {
		err := fmt.Errorf("a datacenter is required, pass it with --datacenter")
		logger.Get().Error().Err(err).Msg("")
		return nil, err
	}

	return sdk.NewAIGenerateRequest(datacenter, prompt), nil
}

// printGenerateResponse keeps the full object for json/yaml output and prints
// the answer as plain text otherwise, because the response is prose that a
// table cell would truncate.
func printGenerateResponse(response *sdk.AIGenerateResponse) error {
	if formatter.IsNativeFormat() {
		return formatter.PrintResult(response, &aiResponsePrintConfig)
	}

	if formatter.IsTextFormat() {
		fmt.Println(strings.TrimRight(response.Result, "\n"))
		if steps := strings.TrimSpace(response.Steps); steps != "" {
			fmt.Println()
			fmt.Println(steps)
		}
		return nil
	}

	return formatter.PrintResult(response, &aiResponsePrintConfig)
}
