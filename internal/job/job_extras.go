package job

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"

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

var jobGroupStatisticsPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"GroupId": {
			Title: "ID",
			Order: 1,
		},
		"GroupType": {
			Title: "Type",
			Order: 2,
		},
		"JobsTotal": {
			Title: "Total",
			Order: 3,
		},
		"JobsCompleted": {
			Title: "Completed",
			Order: 4,
		},
		"JobsThrownError": {
			Title: "Errored",
			Order: 5,
		},
		"GroupCreatedTimestamp": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       6,
		},
		"GroupCompletedTimestamp": {
			Title:       "Completed At",
			Transformer: formatter.FormatDateTimeValue,
			Order:       7,
		},
	},
}

var scheduledJobFunctionPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Name": {
			Title:    "Name",
			MaxWidth: 40,
			Order:    1,
		},
		"Description": {
			Title:    "Description",
			MaxWidth: 80,
			Order:    2,
		},
		"ParamsSchema": {
			Title:    "Parameters",
			MaxWidth: 60,
			Order:    3,
		},
	},
}

// scheduledJobFunctionRow is the table projection of a supported scheduled job
// function: the JSON-schema parameter object is reduced to the list of its
// property names because the tabular formatter renders map-valued fields
// poorly. The json and yaml formats keep the full schema.
type scheduledJobFunctionRow struct {
	Name         string
	Description  string
	ParamsSchema string
}

// JobGetArchived shows a single job read from the job archive.
//
// The response body is parsed raw: the typed SDK JobArchive model requires a
// `links` property that the API does not always return (the same reason JobGet
// parses raw — see jobRaw).
func JobGetArchived(ctx context.Context, jobId string) error {
	logger.Get().Info().Msgf("Get archived job '%s' details", jobId)

	id, err := getJobId(jobId)
	if err != nil {
		return err
	}

	body, err := api.RawJSONRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("/api/v2/jobs/archive/%d", id),
		nil,
		nil,
	)
	if err != nil {
		return err
	}

	return utils.PrintRawObject(body, &jobArchivePrintConfig)
}

// JobGroupStatistics shows the job counters of one job group.
func JobGroupStatistics(ctx context.Context, groupId string) error {
	logger.Get().Info().Msgf("Getting statistics for job group '%s'", groupId)

	id, err := getJobGroupId(groupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	statistics, httpRes, err := client.JobAPI.GetJobGroupStatistics(ctx, id).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(statistics, &jobGroupStatisticsPrintConfig)
}

// ScheduledJobSupportedFunctions lists the functions that can be used when
// creating a scheduled job.
func ScheduledJobSupportedFunctions(ctx context.Context) error {
	logger.Get().Info().Msg("Getting supported scheduled job functions")

	client := api.GetApiClient(ctx)

	functions, httpRes, err := client.JobAPI.GetScheduledJobsSupportedFunctions(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	if formatter.IsNativeFormat() {
		return formatter.PrintResult(functions, nil)
	}

	rows := make([]scheduledJobFunctionRow, 0, len(functions))
	for _, function := range functions {
		row := scheduledJobFunctionRow{
			Name:         function.Name,
			ParamsSchema: formatScheduledJobParams(function.ParamsSchema),
		}
		if function.Description != nil {
			row.Description = *function.Description
		}
		rows = append(rows, row)
	}

	return formatter.PrintResult(rows, &scheduledJobFunctionPrintConfig)
}

// formatScheduledJobParams renders the JSON-schema parameter object of a
// scheduled job function as a stable, comma-separated list of property names.
func formatScheduledJobParams(schema map[string]interface{}) string {
	if len(schema) == 0 {
		return ""
	}

	properties, ok := schema["properties"].(map[string]interface{})
	if !ok || len(properties) == 0 {
		return ""
	}

	names := make([]string, 0, len(properties))
	for name := range properties {
		names = append(names, name)
	}
	sort.Strings(names)

	return strings.Join(names, ", ")
}
