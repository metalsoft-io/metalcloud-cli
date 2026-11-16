package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

// The safety checks of a firmware job refuse it with a 400 (or a 409 for another firmware job running) whose body
// lists the refusals: { statusCode, message, reasonCode, refusals: [{ componentId?, componentName?, reasonCode,
// message }] }.
type firmwareRefusalBody struct {
	Message    string                   `json:"message"`
	ReasonCode string                   `json:"reasonCode"`
	Refusals   []sdk.FirmwareJobRefusal `json:"refusals"`
}

var firmwareRefusalPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"ComponentId":   {Title: "Component ID", Order: 1},
		"ComponentName": {Title: "Component", Order: 2, MaxWidth: 40},
		"ReasonCode":    {Title: "Reason", Order: 3},
		"Message":       {Title: "Message", Order: 4, MaxWidth: 80},
	},
}

// inspectFirmwareJobResponse is InspectResponse for the firmware upgrade endpoints: a refusal by the safety checks
// becomes one line per refused component, with its reason code, instead of the raw response body.
func inspectFirmwareJobResponse(httpRes *http.Response, err error) error {
	if httpRes != nil && (httpRes.StatusCode == http.StatusBadRequest || httpRes.StatusCode == http.StatusConflict) {
		if refusalErr := firmwareRefusalError(err); refusalErr != nil {
			logger.Get().Error().Err(refusalErr).Msg("")
			return refusalErr
		}
	}
	return response_inspector.InspectResponse(httpRes, err)
}

func firmwareRefusalError(err error) error {
	var apiErr *sdk.GenericOpenAPIError
	if !errors.As(err, &apiErr) {
		return nil
	}
	var body firmwareRefusalBody
	if json.Unmarshal(apiErr.Body(), &body) != nil || len(body.Refusals) == 0 {
		return nil
	}
	lines := []string{"the firmware upgrade was refused by its safety checks:"}
	for _, refusal := range body.Refusals {
		lines = append(lines, "  "+describeFirmwareRefusal(refusal))
	}
	return errors.New(strings.Join(lines, "\n"))
}

// describeFirmwareRefusal is "HGX_FW_GPU_SXM_3 (2320): component_not_healthy: <message>", or the reason alone for a
// refusal of the whole job.
func describeFirmwareRefusal(refusal sdk.FirmwareJobRefusal) string {
	subject := ""
	if refusal.ComponentName != nil {
		subject = *refusal.ComponentName
		if refusal.ComponentId != nil {
			subject += fmt.Sprintf(" (%d)", *refusal.ComponentId)
		}
		subject += ": "
	} else if refusal.ComponentId != nil {
		subject = fmt.Sprintf("component %d: ", *refusal.ComponentId)
	}
	return fmt.Sprintf("%s%s: %s", subject, refusal.ReasonCode, refusal.Message)
}

// summarizeFirmwareChecks is the outcome of a component's safety checks in one cell: the refused checks, else the
// checks that could not be made, else "passed".
func summarizeFirmwareChecks(checks []sdk.FirmwareCheckResult) string {
	if len(checks) == 0 {
		return ""
	}
	refused := []string{}
	warnings := []string{}
	for _, check := range checks {
		switch check.Status {
		case sdk.FIRMWARECHECKSTATUS_REFUSED:
			refused = append(refused, string(check.Check))
		case sdk.FIRMWARECHECKSTATUS_WARNING:
			warnings = append(warnings, string(check.Check))
		}
	}
	switch {
	case len(refused) > 0:
		return "refused: " + strings.Join(refused, ", ")
	case len(warnings) > 0:
		return "not checked: " + strings.Join(warnings, ", ")
	default:
		return "passed"
	}
}

// FirmwareCheckRow is one safety check of an upgrade, as `generate-audit --include-preflight` lists them.
type FirmwareCheckRow struct {
	Component string
	Check     string
	Status    string
	Reason    string
	Message   string
}

var firmwareCheckPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Component": {Title: "Component", Order: 1, MaxWidth: 40},
		"Check":     {Title: "Check", Order: 2},
		"Status":    {Title: "Status", Order: 3},
		"Reason":    {Title: "Reason", Order: 4},
		"Message":   {Title: "Message", Order: 5, MaxWidth: 80},
	},
}

// PlannedFirmwareJobRow is one job an upgrade would run, in order.
type PlannedFirmwareJobRow struct {
	Order         int
	Job           string
	Components    string
	TargetVersion string
	Activations   string
}

var plannedFirmwareJobPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Order":         {Title: "#", Order: 1},
		"Job":           {Title: "Job", Order: 2, MaxWidth: 50},
		"Components":    {Title: "Components", Order: 3, MaxWidth: 60},
		"TargetVersion": {Title: "Target Version", Order: 4},
		"Activations":   {Title: "Activation", Order: 5},
	},
}

// BaseboardRow is a GPU baseboard (HGX tray) of a server.
type BaseboardRow struct {
	Baseboard  string
	Components int
	Classes    string
}

var baseboardPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Baseboard":  {Title: "Baseboard", Order: 1},
		"Components": {Title: "Components", Order: 2},
		"Classes":    {Title: "Classes", Order: 3},
	},
}

// printFirmwareAuditDetails lists, under the audit table, each server's GPU baseboards and, with includePreflight,
// whether an upgrade started now would run, its job-wide checks, its warnings, the jobs it would run in order and
// the checks that did not pass.
func printFirmwareAuditDetails(audits []sdk.ServerFirmwareAudit) error {
	for _, audit := range audits {
		if len(audit.Baseboards) > 0 {
			rows := []BaseboardRow{}
			for _, baseboard := range audit.Baseboards {
				rows = append(rows, BaseboardRow{
					Baseboard:  baseboard.BaseboardId,
					Components: len(baseboard.ComponentIds),
					Classes:    strings.Join(baseboard.ComponentClasses, ", "),
				})
			}
			fmt.Printf("\nServer %d GPU baseboards:\n", audit.ServerId)
			if err := formatter.PrintResult(rows, &baseboardPrintConfig); err != nil {
				return err
			}
		}

		preflight := audit.Preflight
		if preflight == nil {
			continue
		}

		verdict := fmt.Sprintf("would run %d job(s)", len(preflight.PlannedJobs))
		if preflight.JobRefused {
			verdict = "would be refused"
		}
		fmt.Printf("\nServer %d: an upgrade started now %s (checked %s)\n", audit.ServerId, verdict, preflight.CheckedAt.Format("2006-01-02 15:04:05 MST"))

		checks := []FirmwareCheckRow{}
		for _, check := range preflight.Checks {
			checks = append(checks, firmwareCheckRow("(job)", check))
		}
		// Per component, the checks that did not pass; a passed component is summarized in the audit table
		for _, component := range audit.Components {
			for _, check := range component.Checks {
				if check.Status == sdk.FIRMWARECHECKSTATUS_REFUSED || check.Status == sdk.FIRMWARECHECKSTATUS_WARNING {
					checks = append(checks, firmwareCheckRow(component.Name, check))
				}
			}
		}
		if err := formatter.PrintResult(checks, &firmwareCheckPrintConfig); err != nil {
			return err
		}

		if len(preflight.Warnings) > 0 {
			fmt.Println("Not checked (these do not refuse the upgrade):")
			for _, warning := range preflight.Warnings {
				fmt.Printf("  - %s\n", warning)
			}
		}

		if len(preflight.PlannedJobs) > 0 {
			fmt.Println("Jobs, in the order they run:")
			if err := formatter.PrintResult(plannedFirmwareJobRows(preflight.PlannedJobs), &plannedFirmwareJobPrintConfig); err != nil {
				return err
			}
		}
	}
	return nil
}

func firmwareCheckRow(component string, check sdk.FirmwareCheckResult) FirmwareCheckRow {
	row := FirmwareCheckRow{
		Component: component,
		Check:     string(check.Check),
		Status:    string(check.Status),
		Message:   derefString(check.Message),
	}
	if check.ReasonCode != nil {
		row.Reason = string(*check.ReasonCode)
	}
	return row
}

func plannedFirmwareJobRows(jobs []sdk.ServerFirmwareAuditPlannedJob) []PlannedFirmwareJobRow {
	rows := []PlannedFirmwareJobRow{}
	for i, job := range jobs {
		name := strings.Join(job.ComponentNames, ", ")
		if job.BundleId != nil {
			bundle := "bundle " + *job.BundleId
			if job.BundleName != nil {
				bundle = *job.BundleName
			}
			name = bundle
			if job.StepIndex != nil && job.StepCount != nil {
				name += fmt.Sprintf(" step %d/%d", *job.StepIndex, *job.StepCount)
			}
			if len(job.StepIds) > 0 {
				name += " (" + strings.Join(job.StepIds, ", ") + ")"
			}
		}
		activations := strings.Join(job.Activations, ", then ")
		if activations == "" {
			activations = "host restart"
		}
		rows = append(rows, PlannedFirmwareJobRow{
			Order:         i + 1,
			Job:           name,
			Components:    strings.Join(collapseComponentNames(job.ComponentNames), ", "),
			TargetVersion: job.TargetVersion,
			Activations:   activations,
		})
	}
	return rows
}

// collapseComponentNames collapses consecutive members of a numbered series, keeping their order:
// HGX_FW_GPU_SXM_1 … HGX_FW_GPU_SXM_8 → HGX_FW_GPU_SXM_1..8. A gap starts a new range.
func collapseComponentNames(names []string) []string {
	result := []string{}
	prefix, first, last := "", -1, -1
	flush := func() {
		switch {
		case first < 0:
		case first == last:
			result = append(result, fmt.Sprintf("%s%d", prefix, first))
		default:
			result = append(result, fmt.Sprintf("%s%d..%d", prefix, first, last))
		}
		first, last = -1, -1
	}
	for _, name := range names {
		digits := len(name)
		for digits > 0 && name[digits-1] >= '0' && name[digits-1] <= '9' {
			digits--
		}
		if digits == len(name) {
			flush()
			result = append(result, name)
			continue
		}
		number := 0
		fmt.Sscanf(name[digits:], "%d", &number)
		if first >= 0 && name[:digits] == prefix && number == last+1 {
			last = number
			continue
		}
		flush()
		prefix, first, last = name[:digits], number, number
	}
	flush()
	return result
}

// ServerFirmwareSearchDryRunOptions selects where the dry run searches: a baseline, explicit catalogs, or, when
// neither is set, the server's target baseline.
type ServerFirmwareSearchDryRunOptions struct {
	BaselineId         int64
	CatalogIds         []int64
	IncludeAllVersions bool
}

// FirmwareMatchRow is one binary the search matched or rejected for a component.
type FirmwareMatchRow struct {
	ComponentId string
	Component   string
	Vendor      string
	Result      string
	BinaryId    string
	Binary      string
	Version     string
	Bundle      string
	Rules       string
	Reason      string
}

var firmwareMatchPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"ComponentId": {Title: "Component ID", Order: 1},
		"Component":   {Title: "Component", Order: 2, MaxWidth: 32},
		"Vendor":      {Title: "Vendor", Order: 3},
		"Result":      {Title: "Result", Order: 4},
		"BinaryId":    {Title: "Binary ID", Order: 5},
		"Binary":      {Title: "Binary", Order: 6, MaxWidth: 40},
		"Version":     {Title: "Version", Order: 7},
		"Bundle":      {Title: "Bundle Release", Order: 8},
		"Rules":       {Title: "Rules", Order: 9, MaxWidth: 60},
		"Reason":      {Title: "Reason", Order: 10, MaxWidth: 80},
	},
}

// ServerFirmwareSearchDryRun runs the binary search a firmware upgrade would use on a server and shows, per
// component, the binaries that matched (with the rules) and the ones that were rejected (with the reason).
// Nothing is stored or scheduled.
func ServerFirmwareSearchDryRun(ctx context.Context, serverId string, options ServerFirmwareSearchDryRunOptions) error {
	logger.Get().Info().Msgf("Explaining the firmware search for server '%s'", serverId)

	serverIdNumeric, err := GetServerId(serverId)
	if err != nil {
		return err
	}

	request := sdk.ServerFirmwareSearchDryRun{}
	if options.BaselineId > 0 && len(options.CatalogIds) > 0 {
		return fmt.Errorf("search either a baseline or catalogs, not both")
	}
	if options.BaselineId > 0 {
		request.BaselineId = sdk.PtrInt64(options.BaselineId)
	}
	for _, catalogId := range options.CatalogIds {
		request.CatalogIds = append(request.CatalogIds, float32(catalogId))
	}
	if options.IncludeAllVersions {
		request.IncludeAllVersions = sdk.PtrBool(true)
	}

	client := api.GetApiClient(ctx)
	result, httpRes, err := client.ServerFirmwareAPI.ServerFirmwareSearchDryRun(ctx, serverIdNumeric).
		ServerFirmwareSearchDryRun(request).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	if formatter.IsNativeFormat() {
		return formatter.PrintResult(result, nil)
	}
	if result.Baseline != nil {
		fmt.Printf("Baseline: %s\n", result.Baseline.Name)
	}
	return formatter.PrintResult(flattenFirmwareMatchReport(result), &firmwareMatchPrintConfig)
}

// flattenFirmwareMatchReport lists the matched binaries of each component, then the rejected ones; a component
// with neither gets one row with the search's note.
func flattenFirmwareMatchReport(result *sdk.IServerFirmwareBinaryResponse) []FirmwareMatchRow {
	if result.MatchReport == nil {
		return []FirmwareMatchRow{}
	}
	reports := []sdk.FirmwareComponentMatchReport{}
	for _, report := range *result.MatchReport {
		reports = append(reports, report)
	}
	// Components with a matching binary first, then by name
	sort.SliceStable(reports, func(i, j int) bool {
		matchedI, matchedJ := len(reports[i].Matched) > 0, len(reports[j].Matched) > 0
		if matchedI != matchedJ {
			return matchedI
		}
		return derefString(reports[i].ComponentName) < derefString(reports[j].ComponentName)
	})

	rows := []FirmwareMatchRow{}
	for _, report := range reports {
		component := derefString(report.ComponentName)
		if component == "" {
			component = fmt.Sprintf("Component %d", report.ComponentId)
		}
		row := func(result string, match sdk.FirmwareBinaryMatch) FirmwareMatchRow {
			rules := []string{}
			for _, rule := range match.Rules {
				rules = append(rules, string(rule))
			}
			bundle := derefString(match.BundleVersion)
			if bundle == derefString(match.PackageVersion) {
				bundle = ""
			}
			return FirmwareMatchRow{
				ComponentId: fmt.Sprintf("%d", report.ComponentId),
				Component:   component,
				Vendor:      report.Vendor,
				Result:      result,
				BinaryId:    fmt.Sprintf("%d", match.BinaryId),
				Binary:      match.Name,
				Version:     derefString(match.PackageVersion),
				Bundle:      bundle,
				Rules:       strings.Join(rules, ", "),
				Reason:      derefString(match.Reason),
			}
		}
		for _, match := range report.Matched {
			rows = append(rows, row("matched", match))
		}
		for _, match := range report.Rejected {
			rows = append(rows, row("rejected", match))
		}
		if len(report.Matched) == 0 && len(report.Rejected) == 0 {
			rows = append(rows, FirmwareMatchRow{
				ComponentId: fmt.Sprintf("%d", report.ComponentId),
				Component:   component,
				Vendor:      report.Vendor,
				Result:      "none",
				Reason:      derefString(report.Note),
			})
		}
	}
	return rows
}
