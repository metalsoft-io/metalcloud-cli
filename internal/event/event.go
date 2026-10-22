package event

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/viper"
)

// eventRaw avoids SDK unmarshal failure: the SDK's typed Event model rejects
// API responses whose enum values it doesn't know (e.g. a numeric EventTypes
// like "109") and requires `severity`, which the API no longer returns (it is `level`). Parsing the raw
// body into plain string fields keeps `event list`/`event get` tolerant of
// SDK <-> API schema drift.
type eventRaw struct {
	Id                interface{} `json:"id"`
	Type              *string     `json:"type"`
	Level             *string     `json:"level"`
	Visibility        *string     `json:"visibility"`
	OccurredTimestamp *string     `json:"occurredTimestamp"`
	Title             *string     `json:"title"`

	// Fields shown in the table formats so that an exported audit trail says what happened,
	// to which entity and who did it. Identifiers are kept loosely typed because the API
	// has been seen to send them both as strings and as numbers.
	Message                   *string     `json:"message"`
	UserId                    interface{} `json:"userId"`
	UserIdAuthenticated       interface{} `json:"userIdAuthenticated"`
	InfrastructureId          interface{} `json:"infrastructureId"`
	ServerId                  interface{} `json:"serverId"`
	NetworkDeviceId           interface{} `json:"networkDeviceId"`
	NetworkDeviceControllerId interface{} `json:"networkDeviceControllerId"`
	StorageId                 interface{} `json:"storageId"`
	VmPoolId                  interface{} `json:"vmPoolId"`
	SiteId                    interface{} `json:"siteId"`
	JobId                     interface{} `json:"jobId"`
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
		"Level": {
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
		"Message": {
			MaxWidth: 60,
			Order:    7,
		},
		"UserId": {
			Title: "User ID",
			Order: 8,
		},
		"UserIdAuthenticated": {
			Title: "Authenticated User ID",
			Order: 9,
		},
		"InfrastructureId": {
			Title: "Infrastructure ID",
			Order: 10,
		},
		"ServerId": {
			Title: "Server ID",
			Order: 11,
		},
		"NetworkDeviceId": {
			Title: "Network Device ID",
			Order: 12,
		},
		"NetworkDeviceControllerId": {
			Title: "Network Device Controller ID",
			Order: 13,
		},
		"StorageId": {
			Title: "Storage ID",
			Order: 14,
		},
		"VmPoolId": {
			Title: "VM Pool ID",
			Order: 15,
		},
		"SiteId": {
			Title: "Site ID",
			Order: 16,
		},
		"JobId": {
			Title: "Job ID",
			Order: 17,
		},
	},
}

// optionalEventColumns are the entity and actor columns that only apply to some events.
// The text and markdown tables leave out the ones that are empty for every event shown, so a
// listing of server events does not carry a wall of empty storage, VM pool and job columns.
var optionalEventColumns = []struct {
	field string
	value func(e eventRaw) interface{}
}{
	{"UserId", func(e eventRaw) interface{} { return e.UserId }},
	{"UserIdAuthenticated", func(e eventRaw) interface{} { return e.UserIdAuthenticated }},
	{"InfrastructureId", func(e eventRaw) interface{} { return e.InfrastructureId }},
	{"ServerId", func(e eventRaw) interface{} { return e.ServerId }},
	{"NetworkDeviceId", func(e eventRaw) interface{} { return e.NetworkDeviceId }},
	{"NetworkDeviceControllerId", func(e eventRaw) interface{} { return e.NetworkDeviceControllerId }},
	{"StorageId", func(e eventRaw) interface{} { return e.StorageId }},
	{"VmPoolId", func(e eventRaw) interface{} { return e.VmPoolId }},
	{"SiteId", func(e eventRaw) interface{} { return e.SiteId }},
	{"JobId", func(e eventRaw) interface{} { return e.JobId }},
}

func isEmptyValue(value interface{}) bool {
	if value == nil {
		return true
	}

	text, ok := value.(string)
	return ok && strings.TrimSpace(text) == ""
}

// eventPrintConfigFor returns the print config to render events with. For the text and
// markdown tables, optional columns that are empty for all of the events are dropped and the
// remaining columns renumbered. CSV, JSON and YAML keep every column so that exports have a
// stable layout whatever the data.
func eventPrintConfigFor(events []eventRaw) *formatter.PrintConfig {
	switch strings.ToLower(viper.GetString(formatter.ConfigFormat)) {
	case "text", "md":
	default:
		return &eventPrintConfig
	}

	fields := make(map[string]formatter.RecordFieldConfig, len(eventPrintConfig.FieldsConfig))
	for name, config := range eventPrintConfig.FieldsConfig {
		fields[name] = config
	}

	for _, column := range optionalEventColumns {
		used := false
		for _, event := range events {
			if !isEmptyValue(column.value(event)) {
				used = true
				break
			}
		}
		if !used {
			delete(fields, column.field)
		}
	}

	names := make([]string, 0, len(fields))
	for name := range fields {
		names = append(names, name)
	}
	sort.Slice(names, func(i, j int) bool { return fields[names[i]].Order < fields[names[j]].Order })

	for i, name := range names {
		config := fields[name]
		config.Order = i + 1
		fields[name] = config
	}

	return &formatter.PrintConfig{FieldsConfig: fields}
}

type ListFlags struct {
	FilterId                  []string
	FilterType                []string
	FilterSeverity            []string
	FilterVisibility          []string
	FilterInfrastructureId    []string
	FilterUserId              []string
	FilterUserIdAuthenticated []string
	FilterServerId            []string
	FilterNetworkDeviceId     []string
	FilterJobId               []string
	FilterSiteId              []string
	SortBy                    []string
	Page                      int
	Limit                     int
	Search                    string
	SearchBy                  []string
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
	if len(flags.FilterUserIdAuthenticated) > 0 {
		request = request.FilterUserIdAuthenticated(flags.FilterUserIdAuthenticated)
	}
	if len(flags.FilterNetworkDeviceId) > 0 {
		request = request.FilterNetworkDeviceId(flags.FilterNetworkDeviceId)
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
		return utils.PrintAllRaw(rawItems, records, meta, len(records), eventPrintConfigFor(records))
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
		return utils.PrintAllRaw(rawItems, records, meta, len(records), eventPrintConfigFor(records))
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

	return utils.PrintAllRaw(rawItems, records, meta, len(records), eventPrintConfigFor(records))
}

func EventGet(ctx context.Context, eventId string) error {
	logger.Get().Info().Msgf("Get event '%s' details", eventId)

	id, err := getEventId(eventId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	// Raw-body parse: the SDK Event model rejects unknown enum values and may require
	// fields the API no longer returns, so SDK unmarshalling fails on valid responses.
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

	return formatter.PrintResult(record, eventPrintConfigFor([]eventRaw{record}))
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
