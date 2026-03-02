package job

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
)

var jobExceptionPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"ExceptionId": {
			Title: "#",
			Order: 1,
		},
		"JobId": {
			Title: "Job ID",
			Order: 2,
		},
		"ArchiveId": {
			Title: "Archive ID",
			Order: 3,
		},
		"Exception": {
			Title:    "Exception",
			MaxWidth: 80,
			Order:    4,
		},
		"CreatedTimestamp": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       5,
		},
	},
}

var jobArchivePrintConfig = formatter.PrintConfig{
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
		"Type": {
			Order: 4,
		},
		"InfrastructureId": {
			Title: "Infra ID",
			Order: 5,
		},
		"JobGroupId": {
			Title: "Group",
			Order: 6,
		},
	},
}

var jobStatisticsPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"StatusToCount": {
			Title: "Status Counts",
			Order: 1,
		},
		"ArchivedCount": {
			Title: "Archived",
			Order: 2,
		},
	},
}

// TODO: jobExceptionRaw is a local type that accepts links as interface{} to work
// around an SDK bug where links is typed as map[string]interface{} but the
// API returns an array.
type jobExceptionRaw struct {
	ExceptionId      int32       `json:"exceptionId"`
	ArchiveId        *int32      `json:"archiveId,omitempty"`
	JobId            *int32      `json:"jobId,omitempty"`
	Exception        interface{} `json:"exception,omitempty"`
	CreatedTimestamp string      `json:"createdTimestamp"`
	Links            interface{} `json:"links,omitempty"`
}

type jobExceptionListRaw struct {
	Data []jobExceptionRaw `json:"data"`
}

func JobExceptions(ctx context.Context, jobId string) error {
	logger.Get().Info().Msgf("Getting exceptions for job '%s'", jobId)

	id, err := getJobId(jobId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	// The SDK's JobException.Links is typed as map[string]interface{} but the API
	// returns an array, causing unmarshaling to fail. We make the call, check for
	// HTTP-level errors, then decode the body ourselves with flexible types.
	_, httpRes, sdkErr := client.JobAPI.GetJobExceptions(ctx, id).Execute()

	// Only use InspectResponse for real HTTP errors (4xx/5xx), not SDK unmarshal errors
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

	var raw jobExceptionListRaw
	if err := json.Unmarshal(body, &raw); err != nil {
		return fmt.Errorf("failed to parse job exceptions: %w", err)
	}

	return formatter.PrintResult(raw, &jobExceptionPrintConfig)
}

func JobStatistics(ctx context.Context) error {
	logger.Get().Info().Msg("Getting job statistics")

	client := api.GetApiClient(ctx)

	stats, httpRes, err := client.JobAPI.GetJobsStatistics(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	// Print status counts in a readable format
	fmt.Println("Job Statistics:")
	for status, count := range stats.StatusToCount {
		fmt.Printf("  %-25s %v\n", status+":", count)
	}
	fmt.Printf("  %-25s %d\n", "archived:", stats.ArchivedCount)

	return nil
}

type ArchiveListFlags struct {
	FilterJobId      []string
	FilterStatus     []string
	FilterJobGroupId []string
	SortBy           []string
}

func JobListArchived(ctx context.Context, flags ArchiveListFlags) error {
	logger.Get().Info().Msg("Listing archived jobs")

	client := api.GetApiClient(ctx)
	request := client.JobAPI.GetJobsFromArchive(ctx)

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

	return formatter.PrintResult(jobs, &jobArchivePrintConfig)
}
