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

type scheduledJobRaw struct {
	Id              interface{} `json:"id"`
	Label           *string     `json:"label"`
	FunctionName    *string     `json:"functionName"`
	Schedule        *string     `json:"schedule"`
	Disabled        interface{} `json:"disabled"`
	LifetimeSeconds interface{} `json:"lifetimeSeconds"`
}

var scheduledJobPrintConfig = formatter.PrintConfig{
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

func ScheduledJobList(ctx context.Context) error {
	logger.Get().Info().Msg("Listing scheduled jobs")

	client := api.GetApiClient(ctx)

	request := client.JobAPI.GetScheduledJobs(ctx).SortBy([]string{"id:ASC"})

	rawItems, meta, err := utils.FetchAllPagesRaw(func(page float32) (*http.Response, error) {
		_, httpRes, _ := request.Page(page).Limit(100).Execute()
		return httpRes, nil
	})
	if err != nil {
		return err
	}

	scheduledJobs, err := utils.UnmarshalRawItems[scheduledJobRaw](rawItems)
	if err != nil {
		return fmt.Errorf("failed to parse scheduled jobs: %w", err)
	}

	return utils.PrintAllRaw(rawItems, scheduledJobs, meta, len(scheduledJobs), &scheduledJobPrintConfig)
}

func ScheduledJobGet(ctx context.Context, scheduledJobId string) error {
	logger.Get().Info().Msgf("Getting scheduled job '%s'", scheduledJobId)

	id, err := getScheduledJobId(scheduledJobId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	// Raw-body parse: see scheduledJobRaw — the `params` type mismatch breaks typed decoding.
	_, httpRes, sdkErr := client.JobAPI.GetScheduledJob(ctx, id).Execute()
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

	var scheduledJob scheduledJobRaw
	if err := json.Unmarshal(body, &scheduledJob); err != nil {
		return fmt.Errorf("failed to parse scheduled job: %w", err)
	}

	return formatter.PrintResult(scheduledJob, &scheduledJobPrintConfig)
}

func ScheduledJobCreate(ctx context.Context, configBytes []byte) error {
	logger.Get().Info().Msg("Creating scheduled job")

	var scheduledJobConfig sdk.ScheduledCronJob
	if err := utils.UnmarshalContent(configBytes, &scheduledJobConfig); err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	_, httpRes, err := client.JobAPI.CreateScheduledJob(ctx).ScheduledCronJob(scheduledJobConfig).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	fmt.Println("Scheduled job created successfully.")
	return nil
}

func ScheduledJobUpdate(ctx context.Context, scheduledJobId string, configBytes []byte) error {
	logger.Get().Info().Msgf("Updating scheduled job '%s'", scheduledJobId)

	id, err := getScheduledJobId(scheduledJobId)
	if err != nil {
		return err
	}

	var scheduledJobConfig sdk.UpdateScheduledJob
	if err := utils.UnmarshalContent(configBytes, &scheduledJobConfig); err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	_, httpRes, err := client.JobAPI.UpdateScheduledJob(ctx, float32(id)).UpdateScheduledJob(scheduledJobConfig).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	fmt.Println("Scheduled job updated successfully.")
	return nil
}

func ScheduledJobDelete(ctx context.Context, scheduledJobId string) error {
	logger.Get().Info().Msgf("Deleting scheduled job '%s'", scheduledJobId)

	id, err := getScheduledJobId(scheduledJobId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.JobAPI.DeleteScheduledJob(ctx, id).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	fmt.Println("Scheduled job deleted.")
	return nil
}

func getScheduledJobId(scheduledJobId string) (int64, error) {
	id, err := strconv.ParseInt(scheduledJobId, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid scheduled job ID '%s': %w", scheduledJobId, err)
	}
	return id, nil
}
