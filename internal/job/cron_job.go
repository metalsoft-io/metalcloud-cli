package job

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

// cronJobRaw works around the SDK bug where CronJob.Links is typed as
// map[string]interface{} but the API may return an array.
type cronJobRaw struct {
	Id                float32                  `json:"id"`
	Label             string                   `json:"label"`
	Description       *string                  `json:"description,omitempty"`
	FunctionName      string                   `json:"functionName"`
	Params            interface{}              `json:"params,omitempty"`
	Schedule          string                   `json:"schedule"`
	WaitForCompletion float32                  `json:"waitForCompletion"`
	LifetimeSeconds   float32                  `json:"lifetimeSeconds"`
	Disabled          float32                  `json:"disabled"`
	Links             interface{}              `json:"links,omitempty"`
}

type cronJobListRaw struct {
	Data []cronJobRaw `json:"data"`
}

var cronJobPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"Label": {
			Order: 2,
		},
		"FunctionName": {
			Title: "Function",
			Order: 3,
		},
		"Schedule": {
			Order: 4,
		},
		"Disabled": {
			Order: 5,
			Transformer: func(v interface{}) string {
				switch val := v.(type) {
				case float32:
					if val == 0 {
						return "enabled"
					}
					return "disabled"
				case float64:
					if val == 0 {
						return "enabled"
					}
					return "disabled"
				}
				return fmt.Sprintf("%v", v)
			},
		},
		"LifetimeSeconds": {
			Title: "Lifetime (s)",
			Order: 6,
		},
	},
}

func CronJobList(ctx context.Context) error {
	logger.Get().Info().Msg("Listing cron jobs")

	client := api.GetApiClient(ctx)

	_, httpRes, sdkErr := client.JobAPI.GetCronJobs(ctx).Execute()

	if httpRes != nil && httpRes.StatusCode >= 400 {
		if err := response_inspector.InspectResponse(httpRes, sdkErr); err != nil {
			return err
		}
	} else if httpRes == nil {
		return sdkErr
	}

	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var raw cronJobListRaw
	if err := json.Unmarshal(body, &raw); err != nil {
		return fmt.Errorf("failed to parse cron jobs: %w", err)
	}

	return formatter.PrintResult(raw.Data, &cronJobPrintConfig)
}

func CronJobGet(ctx context.Context, cronJobId string) error {
	logger.Get().Info().Msgf("Getting cron job '%s'", cronJobId)

	id, err := getCronJobId(cronJobId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	_, httpRes, sdkErr := client.JobAPI.GetCronJob(ctx, id).Execute()

	if httpRes != nil && httpRes.StatusCode >= 400 {
		if err := response_inspector.InspectResponse(httpRes, sdkErr); err != nil {
			return err
		}
	} else if httpRes == nil {
		return sdkErr
	}

	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var raw cronJobRaw
	if err := json.Unmarshal(body, &raw); err != nil {
		return fmt.Errorf("failed to parse cron job: %w", err)
	}

	return formatter.PrintResult(raw, &cronJobPrintConfig)
}

func CronJobCreate(ctx context.Context, configBytes []byte) error {
	logger.Get().Info().Msg("Creating cron job")

	var cronJobConfig sdk.CreateCronJob
	if err := utils.UnmarshalContent(configBytes, &cronJobConfig); err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.JobAPI.CreateCronJob(ctx).CreateCronJob(cronJobConfig).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	fmt.Println("Cron job created successfully.")
	return nil
}

func CronJobUpdate(ctx context.Context, cronJobId string, configBytes []byte) error {
	logger.Get().Info().Msgf("Updating cron job '%s'", cronJobId)

	id, err := getCronJobId(cronJobId)
	if err != nil {
		return err
	}

	var cronJobConfig sdk.UpdateCronJob
	if err := utils.UnmarshalContent(configBytes, &cronJobConfig); err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.JobAPI.UpdateCronJob(ctx, id).UpdateCronJob(cronJobConfig).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	fmt.Println("Cron job updated successfully.")
	return nil
}

func CronJobDelete(ctx context.Context, cronJobId string) error {
	logger.Get().Info().Msgf("Deleting cron job '%s'", cronJobId)

	id, err := getCronJobId(cronJobId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.JobAPI.DeleteCronJob(ctx, id).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	fmt.Println("Cron job deleted.")
	return nil
}

func getCronJobId(cronJobId string) (float32, error) {
	id, err := strconv.ParseFloat(cronJobId, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid cron job ID '%s': %w", cronJobId, err)
	}
	return float32(id), nil
}
