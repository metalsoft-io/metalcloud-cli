package job

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

type cronJobRaw struct {
	Id              interface{} `json:"id"`
	Label           *string     `json:"label"`
	FunctionName    *string     `json:"functionName"`
	Schedule        *string     `json:"schedule"`
	Disabled        interface{} `json:"disabled"`
	LifetimeSeconds interface{} `json:"lifetimeSeconds"`
}

var cronJobPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
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
				// The formatter normalizes whole numeric values to int64, so
				// handle both float and integer representations.
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
				case int64:
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

	request := client.JobAPI.GetCronJobs(ctx).SortBy([]string{"id:ASC"})

	rawItems, meta, err := utils.FetchAllPagesRaw(func(page float32) (*http.Response, error) {
		_, httpRes, _ := request.Page(page).Limit(100).Execute()
		return httpRes, nil
	})
	if err != nil {
		return err
	}

	cronJobs, err := utils.UnmarshalRawItems[cronJobRaw](rawItems)
	if err != nil {
		return fmt.Errorf("failed to parse cron jobs: %w", err)
	}

	return utils.PrintAllRaw(rawItems, cronJobs, meta, len(cronJobs), &cronJobPrintConfig)
}

func CronJobGet(ctx context.Context, cronJobId string) error {
	logger.Get().Info().Msgf("Getting cron job '%s'", cronJobId)

	id, err := getCronJobId(cronJobId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	// Raw-body parse: see cronJobRaw — the `params` type mismatch breaks typed decoding.
	_, httpRes, sdkErr := client.JobAPI.GetCronJob(ctx, id).Execute()
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

	var cronJob cronJobRaw
	if err := json.Unmarshal(body, &cronJob); err != nil {
		return fmt.Errorf("failed to parse cron job: %w", err)
	}

	return formatter.PrintResult(cronJob, &cronJobPrintConfig)
}

func CronJobCreate(ctx context.Context, configBytes []byte) error {
	logger.Get().Info().Msg("Creating cron job")

	var cronJobConfig sdk.CreateCronJob
	if err := utils.UnmarshalContent(configBytes, &cronJobConfig); err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	_, httpRes, err := client.JobAPI.CreateCronJob(ctx).CreateCronJob(cronJobConfig).Execute()
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

	_, httpRes, err := client.JobAPI.UpdateCronJob(ctx, float32(id)).UpdateCronJob(cronJobConfig).Execute()
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

func getCronJobId(cronJobId string) (int64, error) {
	id, err := strconv.ParseInt(cronJobId, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid cron job ID '%s': %w", cronJobId, err)
	}
	return id, nil
}
