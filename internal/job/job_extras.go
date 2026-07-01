package job

import (
	"context"
	"fmt"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
)

var jobExceptionPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"ExceptionId": {
			Title: "ID",
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
			Title: "ID",
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

func JobExceptions(ctx context.Context, jobId string) error {
	logger.Get().Info().Msgf("Getting exceptions for job '%s'", jobId)

	id, err := getJobId(jobId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.JobAPI.GetJobExceptions(ctx, id)

	exceptions, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(exceptions, meta, len(exceptions), &jobExceptionPrintConfig)
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
	Page             int
	Limit            int
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

	switch {
	case flags.Page > 0:
		records, meta, err := utils.FetchPageWindow(request, flags.Page, flags.Limit)
		if err != nil {
			return err
		}
		return utils.PrintAll(records, meta, len(records), &jobArchivePrintConfig)
	case flags.Limit > 0:
		records, meta, err := utils.FetchUpTo(request, flags.Limit)
		if err != nil {
			return err
		}
		return utils.PrintAll(records, meta, len(records), &jobArchivePrintConfig)
	default:
		records, meta, err := utils.FetchAllPages(request)
		if err != nil {
			return err
		}
		return utils.PrintAll(records, meta, len(records), &jobArchivePrintConfig)
	}
}
