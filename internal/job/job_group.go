package job

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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
			Title: "#",
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

	type jobGroupRaw struct {
		Id                interface{} `json:"id"`
		Type              *string     `json:"type"`
		Description       *string     `json:"description"`
		CreatedTimestamp  *string     `json:"createdTimestamp"`
		FinishedTimestamp *string     `json:"finishedTimestamp"`
	}

	var rawItems []json.RawMessage
	var meta sdk.PaginatedResponseMeta
	var err error

	switch {
	case flags.Page > 0:
		rawItems, meta, err = utils.FetchPageWindowRaw(func(p, l float32) (*http.Response, error) {
			_, httpRes, _ := request.Page(p).Limit(l).Execute()
			return httpRes, nil
		}, flags.Page, flags.Limit)
	case flags.Limit > 0:
		rawItems, meta, err = utils.FetchUpToRaw(func(p, l float32) (*http.Response, error) {
			_, httpRes, _ := request.Page(p).Limit(l).Execute()
			return httpRes, nil
		}, flags.Limit)
	default:
		rawItems, meta, err = utils.FetchAllPagesRaw(func(page float32) (*http.Response, error) {
			_, httpRes, _ := request.Page(page).Limit(100).Execute()
			return httpRes, nil
		})
	}
	if err != nil {
		return err
	}

	records, err := utils.UnmarshalRawItems[jobGroupRaw](rawItems)
	if err != nil {
		return fmt.Errorf("failed to parse job groups: %w", err)
	}

	return utils.PrintAllRaw(rawItems, records, meta, len(records), &jobGroupPrintConfig)
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
