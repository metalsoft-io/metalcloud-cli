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

var jobGroupPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"JobGroup": {
			Hidden: true,
			InnerFields: map[string]formatter.RecordFieldConfig{
				"JobGroupId": {
					Title: "#",
					Order: 1,
				},
				"Status": {
					Order: 2,
				},
				"CreatedTimestamp": {
					Title: "Created",
					Order: 3,
				},
			},
		},
	},
}

type GroupListFlags struct {
	FilterJobGroupId []string
	FilterType       []string
	SortBy           []string
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
	}

	groups, httpRes, err := request.Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(groups, &jobGroupPrintConfig)
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

func getJobGroupId(groupId string) (float32, error) {
	id, err := strconv.ParseFloat(groupId, 32)
	if err != nil {
		err := fmt.Errorf("invalid job group ID: '%s'", groupId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}
	return float32(id), nil
}
