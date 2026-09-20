package job

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	sdk "github.com/metalsoft-io/metalcloud-sdk-go"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
)

type jobRaw struct {
	JobId            interface{} `json:"jobId"`
	Status           *string     `json:"status"`
	FunctionName     *string     `json:"functionName"`
	CreatedTimestamp *string     `json:"createdTimestamp"`
	JobGroupId       interface{} `json:"jobGroupId"`
}

var jobPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"JobId": {
			Title: "ID",
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
	Page             int
	Limit            int
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

	switch {
	case flags.Page > 0:
		records, meta, err := utils.FetchPageWindow(request, flags.Page, flags.Limit)
		if err != nil {
			return err
		}
		return utils.PrintAll(records, meta, len(records), &jobPrintConfig)
	case flags.Limit > 0:
		records, meta, err := utils.FetchUpTo(request, flags.Limit)
		if err != nil {
			return err
		}
		return utils.PrintAll(records, meta, len(records), &jobPrintConfig)
	default:
		records, meta, err := utils.FetchAllPages(request)
		if err != nil {
			return err
		}
		return utils.PrintAll(records, meta, len(records), &jobPrintConfig)
	}
}

func JobGet(ctx context.Context, jobId string) error {
	logger.Get().Info().Msgf("Get job '%s' details", jobId)

	id, err := getJobId(jobId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	// Raw-body parse: see jobRaw — the missing required `links` property breaks
	// typed decoding.
	_, httpRes, sdkErr := client.JobAPI.GetJob(ctx, id).Execute()
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

	var job jobRaw
	if err := json.Unmarshal(body, &job); err != nil {
		return fmt.Errorf("failed to parse job: %w", err)
	}

	return formatter.PrintResult(job, &jobPrintConfig)
}

func JobSkip(ctx context.Context, jobId string) error {
	logger.Get().Info().Msgf("Skipping job '%s'", jobId)

	id, err := getJobId(jobId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.JobAPI.SkipJob(ctx, id).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	fmt.Printf("Job %s has been skipped.\n", jobId)
	return nil
}

func JobRetry(ctx context.Context, jobId string, retryEvenIfSuccessful bool) error {
	logger.Get().Info().Msgf("Retrying job '%s'", jobId)

	id, err := getJobId(jobId)
	if err != nil {
		return err
	}

	retryInfo := sdk.NewJobRetryInfo()
	if retryEvenIfSuccessful {
		retryInfo.SetRetryEvenIfSuccessful(retryEvenIfSuccessful)
	}

	client := api.GetApiClient(ctx)

	// The retry body is mandatory for the SDK request: leaving it unset sends a
	// literal `null` body which the API rejects with 400.
	httpRes, err := client.JobAPI.RetryJob(ctx, id).JobRetryInfo(*retryInfo).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	fmt.Printf("Job %s has been retried.\n", jobId)
	return nil
}

func JobKill(ctx context.Context, jobId string) error {
	logger.Get().Info().Msgf("Killing job '%s'", jobId)

	commandInfo := sdk.NewJobCommandInfo()
	commandInfo.SetCommand("kill")
	commandInfo.SetExecuteImmediately(true)

	if err := JobIssueCommand(ctx, jobId, *commandInfo); err != nil {
		return err
	}

	fmt.Printf("Job %s has been killed.\n", jobId)
	return nil
}

// JobIssueCommand issues an operational command (e.g. "kill") for a job.
func JobIssueCommand(ctx context.Context, jobId string, commandInfo sdk.JobCommandInfo) error {
	logger.Get().Info().Msgf("Issuing command for job '%s'", jobId)

	id, err := getJobId(jobId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	// The command body is mandatory for the SDK request: leaving it unset sends
	// a literal `null` body which the API rejects with 400.
	httpRes, err := client.JobAPI.IssueCommandForJob(ctx, id).JobCommandInfo(commandInfo).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return nil
}

func getJobId(jobId string) (int64, error) {
	id, err := strconv.ParseInt(jobId, 10, 64)
	if err != nil {
		err := fmt.Errorf("invalid job ID: '%s'", jobId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}
	return id, nil
}
