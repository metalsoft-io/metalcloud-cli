package network_device_configuration_template

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var (
	ValidDeviceTemplateActions      = []string{"add-global-config", "remove-global-config", "add-neighbor", "remove-neighbor"}
	ValidDeviceTemplateNetworkTypes = []string{"underlay", "overlay"}
	ValidNetworkDeviceDrivers       = []string{"cisco_aci51", "nvidia_ufm", "nexus9000", "cumulus42", "arista_eos", "dell_s4048", "hp5800", "hp5900", "hp5950", "dummy", "junos", "os_10", "sonic_enterprise", "vmware_vds", "cumulus_linux", "brocade", "nvidia_dpu", "dell_s4000", "dell_s6010", "junos18"}
	ValidNetworkDevicePositions     = []string{"all", "tor", "north", "spine", "leaf", "other"}
	ValidBgpNumberings              = []string{"numbered", "unnumbered"}
	ValidBgpLinkConfigurations      = []string{"disabled", "active", "passive"}
)

type ConfigFieldGuideEntry struct {
	Field          string `json:"field" yaml:"field"`
	Required       string `json:"required" yaml:"required"`
	AcceptedValues string `json:"acceptedValues" yaml:"acceptedValues"`
	Example        string `json:"example" yaml:"example"`
}

var configFieldGuidePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Field": {
			Title: "Field",
			Order: 1,
		},
		"Required": {
			Title: "Required",
			Order: 2,
		},
		"AcceptedValues": {
			Title:    "Accepted Values",
			MaxWidth: 55,
			Order:    3,
		},
		"Example": {
			Title: "Example",
			Order: 4,
		},
	},
}

var NetworkDeviceConfigurationTemplatePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"NetworkType": {
			Title:    "Network Type",
			MaxWidth: 40,
			Order:    2,
		},
		"NetworkDeviceDriver": {
			Title: "Network Device Driver",
			Order: 3,
		},
		"NetworkDevicePosition": {
			Title: "Network Device Position",
			Order: 4,
		},
		"RemoteNetworkDevicePosition": {
			Title: "Remote Network Device Position",
			Order: 5,
		},
		"BgpNumbering": {
			Title: "BGP Numbering",
			Order: 6,
		},
		"BgpLinkConfiguration": {
			Title: "BGP Link Configuration",
			Order: 7,
		},
		"LibraryLabel": {
			Title: "Library Label",
			Order: 8,
		},
	},
}

// LibrarySummary is one row of the list-libraries output: a distinct library
// label and how many templates carry it.
type LibrarySummary struct {
	LibraryLabel  string `json:"libraryLabel" yaml:"libraryLabel"`
	TemplateCount int    `json:"templateCount" yaml:"templateCount"`
}

var libraryListPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"LibraryLabel": {
			Title: "Library Label",
			Order: 1,
		},
		"TemplateCount": {
			Title: "Templates",
			Order: 2,
		},
	},
}

func NetworkDeviceConfigurationTemplateList(ctx context.Context, filterId []string, filterLibraryLabel []string) error {
	logger.Get().Info().Msgf("Listing all network device configuration templates")

	client := api.GetApiClient(ctx)

	request := client.NetworkDeviceBGPConfigurationTemplateAPI.GetNetworkDeviceBGPConfigurationTemplates(ctx)

	if len(filterId) > 0 {
		request = request.FilterId(utils.ProcessFilterStringSlice(filterId))
	}
	if len(filterLibraryLabel) > 0 {
		request = request.FilterLibraryLabel(filterLibraryLabel)
	}

	request = request.SortBy([]string{"id:ASC"})

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &NetworkDeviceConfigurationTemplatePrintConfig)
}

func NetworkDeviceConfigurationTemplateConfigExample(ctx context.Context) error {
	if formatter.IsNativeFormat() {
		type templateExample struct {
			Action                      string `json:"action" yaml:"action"`
			NetworkType                 string `json:"networkType" yaml:"networkType"`
			NetworkDeviceDriver         string `json:"networkDeviceDriver" yaml:"networkDeviceDriver"`
			NetworkDevicePosition       string `json:"networkDevicePosition" yaml:"networkDevicePosition"`
			RemoteNetworkDevicePosition string `json:"remoteNetworkDevicePosition" yaml:"remoteNetworkDevicePosition"`
			BgpNumbering                string `json:"bgpNumbering" yaml:"bgpNumbering"`
			BgpLinkConfiguration        string `json:"bgpLinkConfiguration" yaml:"bgpLinkConfiguration"`
			ExecutionType               string `json:"executionType" yaml:"executionType"`
			LibraryLabel                string `json:"libraryLabel" yaml:"libraryLabel"`
			Preparation                 string `json:"preparation,omitempty" yaml:"preparation,omitempty"`
			Configuration               string `json:"configuration" yaml:"configuration"`
		}
		return formatter.PrintResult(templateExample{
			Action:                      strings.Join(ValidDeviceTemplateActions, "|"),
			NetworkType:                 strings.Join(ValidDeviceTemplateNetworkTypes, "|"),
			NetworkDeviceDriver:         strings.Join(ValidNetworkDeviceDrivers, "|"),
			NetworkDevicePosition:       strings.Join(ValidNetworkDevicePositions, "|"),
			RemoteNetworkDevicePosition: strings.Join(ValidNetworkDevicePositions, "|"),
			BgpNumbering:                strings.Join(ValidBgpNumberings, "|"),
			BgpLinkConfiguration:        strings.Join(ValidBgpLinkConfigurations, "|"),
			ExecutionType:               "cli",
			LibraryLabel:                "<string>",
			Configuration:               "<base64 encoded commands>",
		}, nil)
	}

	entries := []ConfigFieldGuideEntry{
		{Field: "action", Required: "yes", AcceptedValues: strings.Join(ValidDeviceTemplateActions, ", "), Example: ValidDeviceTemplateActions[0]},
		{Field: "networkType", Required: "yes", AcceptedValues: strings.Join(ValidDeviceTemplateNetworkTypes, ", "), Example: ValidDeviceTemplateNetworkTypes[0]},
		{Field: "networkDeviceDriver", Required: "yes", AcceptedValues: strings.Join(ValidNetworkDeviceDrivers, ", "), Example: "junos"},
		{Field: "networkDevicePosition", Required: "yes", AcceptedValues: strings.Join(ValidNetworkDevicePositions, ", "), Example: ValidNetworkDevicePositions[0]},
		{Field: "remoteNetworkDevicePosition", Required: "yes", AcceptedValues: strings.Join(ValidNetworkDevicePositions, ", "), Example: ValidNetworkDevicePositions[0]},
		{Field: "bgpNumbering", Required: "yes", AcceptedValues: strings.Join(ValidBgpNumberings, ", "), Example: ValidBgpNumberings[0]},
		{Field: "bgpLinkConfiguration", Required: "yes", AcceptedValues: strings.Join(ValidBgpLinkConfigurations, ", "), Example: ValidBgpLinkConfigurations[1]},
		{Field: "executionType", Required: "yes", AcceptedValues: "cli", Example: "cli"},
		{Field: "libraryLabel", Required: "yes", AcceptedValues: "any string", Example: "my-template"},
		{Field: "preparation", Required: "no", AcceptedValues: "base64 encoded commands", Example: ""},
		{Field: "configuration", Required: "yes", AcceptedValues: "base64 encoded commands", Example: ""},
	}
	return formatter.PrintResult(entries, &configFieldGuidePrintConfig)
}

func NetworkDeviceConfigurationTemplateGet(ctx context.Context, networkDeviceConfigurationTemplateId string) error {
	logger.Get().Info().Msgf("Get network device configuration template %s details", networkDeviceConfigurationTemplateId)

	networkDeviceConfigurationTemplateIdNumeric, err := getNetworkDeviceConfigurationTemplateId(networkDeviceConfigurationTemplateId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	networkDeviceConfigurationTemplate, httpRes, err := client.NetworkDeviceBGPConfigurationTemplateAPI.
		GetNetworkDeviceBGPConfigurationTemplate(ctx, networkDeviceConfigurationTemplateIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}
	return formatter.PrintResult(networkDeviceConfigurationTemplate, &NetworkDeviceConfigurationTemplatePrintConfig)
}

func NetworkDeviceConfigurationTemplateCreate(ctx context.Context, config []byte) error {
	logger.Get().Info().Msgf("Creating network device configuration template")

	var networkDeviceConfigurationTemplateConfig sdk.CreateNetworkDeviceBGPConfigurationTemplate
	err := utils.UnmarshalContent(config, &networkDeviceConfigurationTemplateConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	networkDeviceConfigurationTemplateInfo, httpRes, err := client.NetworkDeviceBGPConfigurationTemplateAPI.
		CreateNetworkDeviceBGPConfigurationTemplate(ctx).
		CreateNetworkDeviceBGPConfigurationTemplate(networkDeviceConfigurationTemplateConfig).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(networkDeviceConfigurationTemplateInfo, &NetworkDeviceConfigurationTemplatePrintConfig)
}

// NetworkDeviceConfigurationTemplateImportLibrary bulk-imports every template
// descriptor file found in dir as a network device configuration template,
// forcing them all under a single libraryLabel so they form one library.
//
// Each file is a JSON/YAML descriptor with the same fields as `config-example`
// (the configuration/preparation fields are base64-encoded commands); the file's
// own libraryLabel, if any, is overridden with libraryLabel. Files are processed
// in name order and an unreadable or invalid file is counted as a failure and
// skipped so one bad file does not abort the rest.
func NetworkDeviceConfigurationTemplateImportLibrary(ctx context.Context, libraryLabel string, dir string, dryRun bool) error {
	logger.Get().Info().Msgf("Importing network device configuration template library %q from %q", libraryLabel, dir)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read template directory %q: %w", dir, err)
	}

	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		switch strings.ToLower(filepath.Ext(entry.Name())) {
		case ".json", ".yaml", ".yml":
			files = append(files, entry.Name())
		}
	}
	sort.Strings(files)

	if len(files) == 0 {
		return fmt.Errorf("no template files (*.json, *.yaml, *.yml) found in %q", dir)
	}

	client := api.GetApiClient(ctx)

	created, failed := 0, 0
	for _, name := range files {
		content, readErr := utils.ReadConfigFromFile(filepath.Join(dir, name))
		if readErr != nil {
			failed++
			logger.Get().Error().Msgf("[%s] read failed: %s", name, readErr.Error())
			continue
		}

		var tmpl sdk.CreateNetworkDeviceBGPConfigurationTemplate
		if err := utils.UnmarshalContent(content, &tmpl); err != nil {
			failed++
			logger.Get().Error().Msgf("[%s] invalid: %s", name, err.Error())
			continue
		}
		tmpl.LibraryLabel = libraryLabel // every template joins the same library

		if dryRun {
			created++
			logger.Get().Info().Msgf("[%s] would import into library %q", name, libraryLabel)
			continue
		}

		info, httpRes, createErr := client.NetworkDeviceBGPConfigurationTemplateAPI.
			CreateNetworkDeviceBGPConfigurationTemplate(ctx).
			CreateNetworkDeviceBGPConfigurationTemplate(tmpl).
			Execute()
		if err := response_inspector.InspectResponse(httpRes, createErr); err != nil {
			failed++
			logger.Get().Error().Msgf("[%s] import failed: %s", name, err.Error())
			continue
		}
		created++
		logger.Get().Info().Msgf("[%s] imported template id=%d", name, info.Id)
	}

	verb, suffix := "imported", ""
	if dryRun {
		verb, suffix = "would import", " (dry-run, no changes made)"
	}
	logger.Get().Info().Msgf("Summary: library=%q, %s=%d, failed=%d%s", libraryLabel, verb, created, failed, suffix)

	if failed > 0 {
		return fmt.Errorf("library import completed with %d failure(s)", failed)
	}
	return nil
}

// NetworkDeviceConfigurationTemplateExportLibrary writes every template that
// carries libraryLabel to dir, one JSON descriptor per template. Each file is
// the inverse of import-library's input (the create fields only, no id or
// timestamps), so an exported directory can be re-imported as-is. dir is created
// if it does not exist.
func NetworkDeviceConfigurationTemplateExportLibrary(ctx context.Context, libraryLabel string, dir string) error {
	logger.Get().Info().Msgf("Exporting network device configuration template library %q to %q", libraryLabel, dir)

	templates, err := fetchAllTemplates(ctx)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory %q: %w", dir, err)
	}

	exported := 0
	for _, t := range templates {
		if t.LibraryLabel != libraryLabel {
			continue
		}

		descriptor := sdk.CreateNetworkDeviceBGPConfigurationTemplate{
			Action:                      t.Action,
			NetworkType:                 t.NetworkType,
			NetworkDeviceDriver:         t.NetworkDeviceDriver,
			NetworkDevicePosition:       t.NetworkDevicePosition,
			RemoteNetworkDevicePosition: t.RemoteNetworkDevicePosition,
			BgpNumbering:                t.BgpNumbering,
			BgpLinkConfiguration:        t.BgpLinkConfiguration,
			ExecutionType:               t.ExecutionType,
			LibraryLabel:                t.LibraryLabel,
			Preparation:                 t.Preparation,
			Configuration:               t.Configuration,
		}

		data, err := json.MarshalIndent(descriptor, "", "  ")
		if err != nil {
			return fmt.Errorf("template id=%d: %w", t.Id, err)
		}
		data = append(data, '\n')

		name := fmt.Sprintf("template-%d.json", t.Id)
		if err := os.WriteFile(filepath.Join(dir, name), data, 0o644); err != nil {
			return fmt.Errorf("failed to write %q: %w", name, err)
		}
		exported++
		logger.Get().Info().Msgf("[%s] exported template id=%d", name, t.Id)
	}

	if exported == 0 {
		return fmt.Errorf("no templates found in library %q", libraryLabel)
	}

	logger.Get().Info().Msgf("Summary: library=%q, exported=%d, directory=%q", libraryLabel, exported, dir)
	return nil
}

// NetworkDeviceConfigurationTemplateListLibraries prints every distinct library
// label across all network device configuration templates, with the number of
// templates in each.
func NetworkDeviceConfigurationTemplateListLibraries(ctx context.Context) error {
	logger.Get().Info().Msgf("Listing network device configuration template libraries")

	templates, err := fetchAllTemplates(ctx)
	if err != nil {
		return err
	}

	counts := map[string]int{}
	for _, t := range templates {
		counts[t.LibraryLabel]++
	}

	libraries := make([]LibrarySummary, 0, len(counts))
	for label, n := range counts {
		libraries = append(libraries, LibrarySummary{LibraryLabel: label, TemplateCount: n})
	}
	sort.Slice(libraries, func(i, j int) bool { return libraries[i].LibraryLabel < libraries[j].LibraryLabel })

	return formatter.PrintResult(libraries, &libraryListPrintConfig)
}

// fetchAllTemplates returns every network device configuration template. Both
// export-library and list-libraries group/filter by libraryLabel client-side so
// they don't depend on the server-side filter's match semantics.
func fetchAllTemplates(ctx context.Context) ([]sdk.NetworkDeviceBGPConfigurationTemplate, error) {
	client := api.GetApiClient(ctx)
	request := client.NetworkDeviceBGPConfigurationTemplateAPI.
		GetNetworkDeviceBGPConfigurationTemplates(ctx).
		SortBy([]string{"id:ASC"})

	templates, _, err := utils.FetchAllPages(request)
	if err != nil {
		return nil, err
	}
	return templates, nil
}

func NetworkDeviceConfigurationTemplateUpdate(ctx context.Context, networkDeviceConfigurationTemplateId string, config []byte) error {
	logger.Get().Info().Msgf("Updating network device configuration template %s", networkDeviceConfigurationTemplateId)

	networkDeviceConfigurationTemplateIdNumeric, err := getNetworkDeviceConfigurationTemplateId(networkDeviceConfigurationTemplateId)
	if err != nil {
		return err
	}

	var networkDeviceConfigurationTemplateConfig sdk.UpdateNetworkDeviceBGPConfigurationTemplate
	err = utils.UnmarshalContent(config, &networkDeviceConfigurationTemplateConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	networkDeviceConfigurationTemplateInfo, httpRes, err := client.NetworkDeviceBGPConfigurationTemplateAPI.
		UpdateNetworkDeviceBGPConfigurationTemplate(ctx, networkDeviceConfigurationTemplateIdNumeric).
		UpdateNetworkDeviceBGPConfigurationTemplate(networkDeviceConfigurationTemplateConfig).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(networkDeviceConfigurationTemplateInfo, &NetworkDeviceConfigurationTemplatePrintConfig)
}

func NetworkDeviceConfigurationTemplateDelete(ctx context.Context, networkDeviceConfigurationTemplateId string) error {
	logger.Get().Info().Msgf("Deleting network device configuration template %s", networkDeviceConfigurationTemplateId)

	networkDeviceConfigurationTemplateIdNumeric, err := getNetworkDeviceConfigurationTemplateId(networkDeviceConfigurationTemplateId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.NetworkDeviceBGPConfigurationTemplateAPI.
		DeleteNetworkDeviceBGPConfigurationTemplate(ctx, networkDeviceConfigurationTemplateIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Network device configuration template %s deleted", networkDeviceConfigurationTemplateId)
	return nil
}

func getNetworkDeviceConfigurationTemplateId(networkDeviceConfigurationTemplateId string) (int64, error) {
	networkDeviceConfigurationTemplateIdNumeric, err := strconv.ParseInt(networkDeviceConfigurationTemplateId, 10, 64)
	if err != nil {
		err := fmt.Errorf("invalid network device configuration template ID: '%s'", networkDeviceConfigurationTemplateId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return networkDeviceConfigurationTemplateIdNumeric, nil
}
