package job

import (
	"context"
	"fmt"
	"strconv"
	"time"

	sdk "github.com/metalsoft-io/metalcloud-sdk-go"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
)

var jobGroupPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"Type": {
			Title: "Type",
			Order: 2,
		},
		"Description": {
			MaxWidth: 60,
			Order:    3,
		},
		"CreatedTimestamp": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       4,
		},
		"FinishedTimestamp": {
			Title:       "Finished",
			Transformer: formatter.FormatDateTimeValue,
			Order:       5,
		},
	},
}

type GroupListFlags struct {
	FilterJobGroupId []string
	FilterType       []string
	SortBy           []string
	Page             int
	Limit            int
}

func JobGroupList(ctx context.Context, flags GroupListFlags) error {
	logger.Get().Info().Msg("Listing job groups")

	client := api.GetApiClient(ctx)
	request := client.JobAPI.GetJobGroups(ctx)

	if len(flags.FilterJobGroupId) > 0 {
		request = request.FilterId(flags.FilterJobGroupId)
	}
	if len(flags.FilterType) > 0 {
		request = request.FilterType(flags.FilterType)
	}
	if len(flags.SortBy) > 0 {
		request = request.SortBy(flags.SortBy)
	} else {
		request = request.SortBy([]string{"id:ASC"})
	}

	switch {
	case flags.Page > 0:
		records, meta, err := utils.FetchPageWindow(request, flags.Page, flags.Limit)
		if err != nil {
			return err
		}
		return utils.PrintAll(records, meta, len(records), &jobGroupPrintConfig)
	case flags.Limit > 0:
		records, meta, err := utils.FetchUpTo(request, flags.Limit)
		if err != nil {
			return err
		}
		return utils.PrintAll(records, meta, len(records), &jobGroupPrintConfig)
	default:
		records, meta, err := utils.FetchAllPages(request)
		if err != nil {
			return err
		}
		return utils.PrintAll(records, meta, len(records), &jobGroupPrintConfig)
	}
}

func JobGroupGet(ctx context.Context, groupId string) error {
	logger.Get().Info().Msgf("Get job group '%s' details", groupId)

	id, err := getJobGroupId(groupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)
	group, httpRes, err := client.JobAPI.GetJobGroup(ctx, id).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(group, &jobGroupPrintConfig)
}

func JobGroupWait(ctx context.Context, groupId string) error {
	logger.Get().Info().Msgf("Waiting for job group '%s' to finish", groupId)

	id, err := getJobGroupId(groupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	group, httpRes, err := client.JobAPI.GetJobGroup(ctx, id).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	if err := formatter.PrintResult(group, &jobGroupPrintConfig); err != nil {
		return err
	}

	if jobGroupFinished(group) {
		return nil
	}

	fmt.Printf("Waiting for job group %d to finish...\n", id)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}

		group, httpRes, err = client.JobAPI.GetJobGroup(ctx, id).Execute()
		if err := response_inspector.InspectResponse(httpRes, err); err != nil {
			fmt.Printf("Error while waiting for job group to finish: %s", err.Error())
		} else if jobGroupFinished(group) {
			return formatter.PrintResult(group, &jobGroupPrintConfig)
		}
	}
}

func jobGroupFinished(group *sdk.JobGroup) bool {
	return group != nil && group.FinishedTimestamp != nil && *group.FinishedTimestamp != ""
}

func getJobGroupId(groupId string) (int64, error) {
	id, err := strconv.ParseInt(groupId, 10, 64)
	if err != nil {
		err := fmt.Errorf("invalid job group ID: '%s'", groupId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}
	return id, nil
}
