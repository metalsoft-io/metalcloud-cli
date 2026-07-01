package event

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
)

// eventRaw avoids SDK unmarshal failure: the SDK's typed Event model rejects
// API responses whose enum values it doesn't know (e.g. a numeric EventTypes
// like "109") and requires `severity`, which the API may omit. Parsing the raw
// body into plain string fields keeps `event list`/`event get` tolerant of
// SDK <-> API schema drift.
type eventRaw struct {
	Id                interface{} `json:"id"`
	Type              *string     `json:"type"`
	Severity          *string     `json:"severity"`
	Visibility        *string     `json:"visibility"`
	OccurredTimestamp *string     `json:"occurredTimestamp"`
	Title             *string     `json:"title"`
}

var eventPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
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
		"OccurredTimestamp": {
			Title:       "Time",
			Transformer: formatter.FormatDateTimeValue,
			Order:       5,
		},
		"Title": {
			Order: 6,
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

	if len(flags.FilterId) > 0 {
		request = request.FilterId(flags.FilterId)
	}
	if len(flags.FilterType) > 0 {
		request = request.FilterType(flags.FilterType)
	}
	if len(flags.FilterSeverity) > 0 {
		request = request.FilterLevel(flags.FilterSeverity)
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

	if flags.Page > 0 {
		// Specific page requested — fetch that page window, spanning API pages when limit > 100.
		rawItems, meta, err := utils.FetchPageWindowRaw(func(p, l float32) (*http.Response, error) {
			_, httpRes, _ := request.Page(p).Limit(l).Execute()
			return httpRes, nil
		}, flags.Page, flags.Limit)
		if err != nil {
			return err
		}
		records, err := utils.UnmarshalRawItems[eventRaw](rawItems)
		if err != nil {
			return fmt.Errorf("failed to parse events: %w", err)
		}
		return utils.PrintAllRaw(rawItems, records, meta, len(records), &eventPrintConfig)
	}

	if flags.Limit > 0 {
		// Limit without page — fetch exactly N records, spanning pages as needed.
		rawItems, meta, err := utils.FetchUpToRaw(func(page, limit float32) (*http.Response, error) {
			_, httpRes, _ := request.Page(page).Limit(limit).Execute()
			return httpRes, nil
		}, flags.Limit)
		if err != nil {
			return err
		}
		records, err := utils.UnmarshalRawItems[eventRaw](rawItems)
		if err != nil {
			return fmt.Errorf("failed to parse events: %w", err)
		}
		return utils.PrintAllRaw(rawItems, records, meta, len(records), &eventPrintConfig)
	}

	rawItems, meta, err := utils.FetchAllPagesRaw(func(page float32) (*http.Response, error) {
		_, httpRes, _ := request.Page(page).Limit(100).Execute()
		return httpRes, nil
	})
	if err != nil {
		return err
	}

	records, err := utils.UnmarshalRawItems[eventRaw](rawItems)
	if err != nil {
		return fmt.Errorf("failed to parse events: %w", err)
	}

	return utils.PrintAllRaw(rawItems, records, meta, len(records), &eventPrintConfig)
}

func EventGet(ctx context.Context, eventId string) error {
	logger.Get().Info().Msgf("Get event '%s' details", eventId)

	id, err := getEventId(eventId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	// Raw-body parse: the SDK Event model rejects unknown enum values and requires
	// `severity`, which the API can omit, so SDK unmarshalling fails on valid responses.
	_, httpRes, sdkErr := client.EventAPI.GetEvent(ctx, id).Execute()
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

	var record eventRaw
	if err := json.Unmarshal(body, &record); err != nil {
		return fmt.Errorf("failed to parse event: %w", err)
	}

	return formatter.PrintResult(record, &eventPrintConfig)
}

func getEventId(eventId string) (int64, error) {
	id, err := strconv.ParseInt(eventId, 10, 64)
	if err != nil {
		err := fmt.Errorf("invalid event ID: '%s'", eventId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}
	return id, nil
}
