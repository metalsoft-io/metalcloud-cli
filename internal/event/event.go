package event

import (
	"context"
	"fmt"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
)

var eventPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Event": {
			Hidden: true,
			InnerFields: map[string]formatter.RecordFieldConfig{
				"Id": {
					Title: "#",
					Order: 1,
				},
				"Type": {
					Order: 2,
				},
				"Severity": {
					Order: 3,
				},
				"Visibility": {
					Order: 4,
				},
				"CreatedAt": {
					Title: "Created",
					Order: 5,
				},
				"Title": {
					Order: 6,
				},
			},
		},
	},
}

type ListFlags struct {
	FilterId               []string
	FilterType             []string
	FilterSeverity         []string
	FilterVisibility       []string
	FilterInfrastructureId []string
	FilterUserId           []string
	FilterServerId         []string
	FilterJobId            []string
	FilterSiteId           []string
	SortBy                 []string
	Page                   int
	Limit                  int
	Search                 string
	SearchBy               []string
}

func EventList(ctx context.Context, flags ListFlags) error {
	logger.Get().Info().Msg("Listing events")

	client := api.GetApiClient(ctx)
	request := client.EventAPI.GetEvents(ctx)

	if flags.Page > 0 {
		request = request.Page(float32(flags.Page))
	}
	if flags.Limit > 0 {
		request = request.Limit(float32(flags.Limit))
	}
	if len(flags.FilterId) > 0 {
		request = request.FilterId(flags.FilterId)
	}
	if len(flags.FilterType) > 0 {
		request = request.FilterType(flags.FilterType)
	}
	if len(flags.FilterSeverity) > 0 {
		request = request.FilterSeverity(flags.FilterSeverity)
	}
	if len(flags.FilterVisibility) > 0 {
		request = request.FilterVisibility(flags.FilterVisibility)
	}
	if len(flags.FilterInfrastructureId) > 0 {
		request = request.FilterInfrastructureId(flags.FilterInfrastructureId)
	}
	if len(flags.FilterUserId) > 0 {
		request = request.FilterUserId(flags.FilterUserId)
	}
	if len(flags.FilterServerId) > 0 {
		request = request.FilterServerId(flags.FilterServerId)
	}
	if len(flags.FilterJobId) > 0 {
		request = request.FilterJobId(flags.FilterJobId)
	}
	if len(flags.FilterSiteId) > 0 {
		request = request.FilterSiteId(flags.FilterSiteId)
	}
	if len(flags.SortBy) > 0 {
		request = request.SortBy(flags.SortBy)
	}
	if flags.Search != "" {
		request = request.Search(flags.Search)
	}
	if len(flags.SearchBy) > 0 {
		request = request.SearchBy(flags.SearchBy)
	}

	events, httpRes, err := request.Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(events, &eventPrintConfig)
}

func EventGet(ctx context.Context, eventId string) error {
	logger.Get().Info().Msgf("Get event '%s' details", eventId)

	id, err := getEventId(eventId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)
	event, httpRes, err := client.EventAPI.GetEvent(ctx, id).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(event, &eventPrintConfig)
}

func getEventId(eventId string) (float32, error) {
	id, err := strconv.ParseFloat(eventId, 32)
	if err != nil {
		err := fmt.Errorf("invalid event ID: '%s'", eventId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}
	return float32(id), nil
}
