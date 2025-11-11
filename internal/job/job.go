package job

import (
	"context"
	"fmt"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
)

var jobPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"JobId": {
			Title: "#",
			Order: 1,
		},
		"Status": {
			Order: 2,
		},
		"FunctionName": {
			Order: 3,
		},
		"CreatedTimestamp": {
			Title: "Created",
			Order: 4,
		},
		"JobGroupId": {
			Title: "Group",
			Order: 5,
		},
	},
}

type ListFlags struct {
	FilterJobId      []string
	FilterStatus     []string
	FilterJobGroupId []string
	SortBy           []string
}

func JobList(ctx context.Context, flags ListFlags) error {
	logger.Get().Info().Msg("Listing jobs")

	client := api.GetApiClient(ctx)
	request := client.JobAPI.GetJobs(ctx)

	if len(flags.FilterJobId) > 0 {
		request = request.FilterJobId(flags.FilterJobId)
	}
	if len(flags.FilterStatus) > 0 {
		request = request.FilterStatus(flags.FilterStatus)
	}
	if len(flags.FilterJobGroupId) > 0 {
		request = request.FilterJobGroupId(flags.FilterJobGroupId)
	}
	if len(flags.SortBy) > 0 {
		request = request.SortBy(flags.SortBy)
	}

	jobs, httpRes, err := request.Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(jobs, &jobPrintConfig)
}

func JobGet(ctx context.Context, jobId string) error {
	logger.Get().Info().Msgf("Get job '%s' details", jobId)

	id, err := getJobId(jobId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)
	job, httpRes, err := client.JobAPI.GetJob(ctx, id).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(job, &jobPrintConfig)
}

func getJobId(jobId string) (float32, error) {
	id, err := strconv.ParseFloat(jobId, 32)
	if err != nil {
		err := fmt.Errorf("invalid job ID: '%s'", jobId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}
	return float32(id), nil
}
