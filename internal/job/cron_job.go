package job

import (
	"context"
	"fmt"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

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

	request := client.JobAPI.GetCronJobs(ctx).SortBy([]string{"id:ASC"})

	cronJobs, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(cronJobs, meta, len(cronJobs), &cronJobPrintConfig)
}

func CronJobGet(ctx context.Context, cronJobId string) error {
	logger.Get().Info().Msgf("Getting cron job '%s'", cronJobId)

	id, err := getCronJobId(cronJobId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	cronJob, httpRes, err := client.JobAPI.GetCronJob(ctx, id).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
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
