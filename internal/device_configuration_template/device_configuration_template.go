package device_configuration_template

import (
	"context"
	"fmt"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var (
	ValidDeviceDrivers   = []string{"cisco_aci51", "cisco_ndfc", "nvidia_ufm", "nexus9000", "cumulus42", "arista_eos", "dell_s4048", "hp5800", "hp5900", "hp5950", "dummy", "junos", "os_10", "sonic_enterprise", "vmware_vds", "cumulus_linux", "brocade", "nvidia_dpu", "dell_s4000", "dell_s6010", "junos18"}
	ValidExecutionTypes  = []string{"cli", "json_patch"}
	ValidLifecycleStages = []string{"provisioning", "configuration", "decommissioning"}
	ValidApplyModes      = []string{"once", "always"}
)

var DeviceConfigurationTemplatePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"Label": {
			Title:    "Label",
			MaxWidth: 30,
			Order:    2,
		},
		"Name": {
			Title:    "Name",
			MaxWidth: 30,
			Order:    3,
		},
		"DeviceDriver": {
			Title: "Driver",
			Order: 4,
		},
		"ExecutionType": {
			Title: "Execution Type",
			Order: 5,
		},
		"Tags": {
			Title:       "Tags",
			Transformer: formatter.FormatStringListValue,
			Order:       6,
		},
		"CreatedTimestamp": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       7,
		},
		"Revision": {
			Title: "Revision",
			Order: 8,
		},
	},
}

var DeviceConfigurationTemplateProfilePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"DeviceConfigurationTemplateId": {
			Title: "Template",
			Order: 2,
		},
		"NetworkDeviceId": {
			Title: "Device",
			Order: 3,
		},
		"NetworkFabricId": {
			Title: "Fabric",
			Order: 4,
		},
		"LifecycleStage": {
			Title: "Lifecycle Stage",
			Order: 5,
		},
		"IsEnabled": {
			Title:       "Enabled",
			Transformer: formatter.FormatBooleanValue,
			Order:       6,
		},
		"Priority": {
			Title: "Priority",
			Order: 7,
		},
		"ApplyMode": {
			Title: "Apply Mode",
			Order: 8,
		},
		"Revision": {
			Title: "Revision",
			Order: 9,
		},
	},
}

var renderedTemplatePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Rendered": {
			Title: "Rendered",
			Order: 1,
		},
	},
}

// ----------------------------------------------------------------------------
// Configuration templates (config sub-resource)
// ----------------------------------------------------------------------------

func DeviceConfigurationTemplateList(ctx context.Context, filterId []string, filterLabel []string, filterName []string) error {
	logger.Get().Info().Msgf("Listing all device configuration templates")

	client := api.GetApiClient(ctx)

	request := client.DeviceConfigurationTemplateAPI.GetDeviceConfigurationTemplates(ctx)

	if len(filterId) > 0 {
		request = request.FilterId(utils.ProcessFilterStringSlice(filterId))
	}
	if len(filterLabel) > 0 {
		request = request.FilterLabel(utils.ProcessFilterStringSlice(filterLabel))
	}
	if len(filterName) > 0 {
		request = request.FilterName(utils.ProcessFilterStringSlice(filterName))
	}

	request = request.SortBy([]string{"id:ASC"})

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &DeviceConfigurationTemplatePrintConfig)
}

func DeviceConfigurationTemplateGet(ctx context.Context, deviceConfigurationTemplateId string) error {
	logger.Get().Info().Msgf("Get device configuration template %s details", deviceConfigurationTemplateId)

	id, err := parseId(deviceConfigurationTemplateId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	template, httpRes, err := client.DeviceConfigurationTemplateAPI.
		GetDeviceConfigurationTemplate(ctx, id).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(template, &DeviceConfigurationTemplatePrintConfig)
}

func DeviceConfigurationTemplateCreate(ctx context.Context, config []byte) error {
	logger.Get().Info().Msgf("Creating device configuration template")

	var payload sdk.CreateDeviceConfigurationTemplate
	if err := utils.UnmarshalContent(config, &payload); err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	template, httpRes, err := client.DeviceConfigurationTemplateAPI.
		CreateDeviceConfigurationTemplate(ctx).
		CreateDeviceConfigurationTemplate(payload).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(template, &DeviceConfigurationTemplatePrintConfig)
}

func DeviceConfigurationTemplateUpdate(ctx context.Context, deviceConfigurationTemplateId string, config []byte) error {
	logger.Get().Info().Msgf("Updating device configuration template %s", deviceConfigurationTemplateId)

	id, err := parseId(deviceConfigurationTemplateId)
	if err != nil {
		return err
	}

	var payload sdk.UpdateDeviceConfigurationTemplate
	if err := utils.UnmarshalContent(config, &payload); err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	existing, httpRes, err := client.DeviceConfigurationTemplateAPI.
		GetDeviceConfigurationTemplate(ctx, id).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	template, httpRes, err := client.DeviceConfigurationTemplateAPI.
		UpdateDeviceConfigurationTemplate(ctx, id).
		UpdateDeviceConfigurationTemplate(payload).
		IfMatch(existing.Revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(template, &DeviceConfigurationTemplatePrintConfig)
}

func DeviceConfigurationTemplateDelete(ctx context.Context, deviceConfigurationTemplateId string) error {
	logger.Get().Info().Msgf("Deleting device configuration template %s", deviceConfigurationTemplateId)

	id, err := parseId(deviceConfigurationTemplateId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	existing, httpRes, err := client.DeviceConfigurationTemplateAPI.
		GetDeviceConfigurationTemplate(ctx, id).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	httpRes, err = client.DeviceConfigurationTemplateAPI.
		DeleteDeviceConfigurationTemplate(ctx, id).
		IfMatch(existing.Revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Device configuration template %s deleted", deviceConfigurationTemplateId)
	return nil
}

func DeviceConfigurationTemplateRender(ctx context.Context, config []byte) error {
	logger.Get().Info().Msgf("Rendering device configuration template content")

	var payload sdk.RenderDeviceConfigurationTemplate
	if err := utils.UnmarshalContent(config, &payload); err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	rendered, httpRes, err := client.DeviceConfigurationTemplateAPI.
		RenderDeviceConfigurationTemplate(ctx).
		RenderDeviceConfigurationTemplate(payload).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(rendered, &renderedTemplatePrintConfig)
}

func DeviceConfigurationTemplateRenderSaved(ctx context.Context, deviceConfigurationTemplateId string, config []byte) error {
	logger.Get().Info().Msgf("Rendering saved device configuration template %s", deviceConfigurationTemplateId)

	id, err := parseId(deviceConfigurationTemplateId)
	if err != nil {
		return err
	}

	var payload sdk.RenderSavedDeviceConfigurationTemplate
	if err := utils.UnmarshalContent(config, &payload); err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	rendered, httpRes, err := client.DeviceConfigurationTemplateAPI.
		RenderSavedDeviceConfigurationTemplate(ctx, id).
		RenderSavedDeviceConfigurationTemplate(payload).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(rendered, &renderedTemplatePrintConfig)
}

func DeviceConfigurationTemplateConfigExample(ctx context.Context) error {
	example := sdk.CreateDeviceConfigurationTemplate{
		Label:         "my-device-template",
		Name:          sdk.PtrString("My device template"),
		Description:   sdk.PtrString("Example device configuration template"),
		DeviceDriver:  sdk.SwitchDriver(ValidDeviceDrivers[0]),
		ExecutionType: sdk.NetworkTemplateExecutionType(ValidExecutionTypes[0]),
		TemplateContent: sdk.PtrString("hostname {{ hostname }}"),
		CustomVariablesJson: map[string]interface{}{
			"hostname": "switch-01",
		},
		Tags: []string{"example"},
	}

	return formatter.PrintResult(example, nil)
}

// ----------------------------------------------------------------------------
// Configuration template profiles (profile sub-resource)
// ----------------------------------------------------------------------------

func DeviceConfigurationTemplateProfileList(ctx context.Context, filterId []string, filterTemplateId []string, filterNetworkDeviceId []string, filterNetworkFabricId []string) error {
	logger.Get().Info().Msgf("Listing all device configuration template profiles")

	client := api.GetApiClient(ctx)

	request := client.DeviceConfigurationTemplateAPI.GetDeviceConfigurationTemplateProfiles(ctx)

	if len(filterId) > 0 {
		request = request.FilterId(utils.ProcessFilterStringSlice(filterId))
	}
	if len(filterTemplateId) > 0 {
		request = request.FilterDeviceConfigurationTemplateId(utils.ProcessFilterStringSlice(filterTemplateId))
	}
	if len(filterNetworkDeviceId) > 0 {
		request = request.FilterNetworkDeviceId(utils.ProcessFilterStringSlice(filterNetworkDeviceId))
	}
	if len(filterNetworkFabricId) > 0 {
		request = request.FilterNetworkFabricId(utils.ProcessFilterStringSlice(filterNetworkFabricId))
	}

	request = request.SortBy([]string{"id:ASC"})

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &DeviceConfigurationTemplateProfilePrintConfig)
}

func DeviceConfigurationTemplateProfileGet(ctx context.Context, profileId string) error {
	logger.Get().Info().Msgf("Get device configuration template profile %s details", profileId)

	id, err := parseId(profileId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	profile, httpRes, err := client.DeviceConfigurationTemplateAPI.
		GetDeviceConfigurationTemplateProfile(ctx, id).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(profile, &DeviceConfigurationTemplateProfilePrintConfig)
}

func DeviceConfigurationTemplateProfileCreate(ctx context.Context, config []byte) error {
	logger.Get().Info().Msgf("Creating device configuration template profile")

	var payload sdk.CreateDeviceConfigurationTemplateProfile
	if err := utils.UnmarshalContent(config, &payload); err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	profile, httpRes, err := client.DeviceConfigurationTemplateAPI.
		CreateDeviceConfigurationTemplateProfile(ctx).
		CreateDeviceConfigurationTemplateProfile(payload).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(profile, &DeviceConfigurationTemplateProfilePrintConfig)
}

func DeviceConfigurationTemplateProfileUpdate(ctx context.Context, profileId string, config []byte) error {
	logger.Get().Info().Msgf("Updating device configuration template profile %s", profileId)

	id, err := parseId(profileId)
	if err != nil {
		return err
	}

	var payload sdk.UpdateDeviceConfigurationTemplateProfile
	if err := utils.UnmarshalContent(config, &payload); err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	existing, httpRes, err := client.DeviceConfigurationTemplateAPI.
		GetDeviceConfigurationTemplateProfile(ctx, id).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	profile, httpRes, err := client.DeviceConfigurationTemplateAPI.
		UpdateDeviceConfigurationTemplateProfile(ctx, id).
		UpdateDeviceConfigurationTemplateProfile(payload).
		IfMatch(existing.Revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(profile, &DeviceConfigurationTemplateProfilePrintConfig)
}

func DeviceConfigurationTemplateProfileDelete(ctx context.Context, profileId string) error {
	logger.Get().Info().Msgf("Deleting device configuration template profile %s", profileId)

	id, err := parseId(profileId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	existing, httpRes, err := client.DeviceConfigurationTemplateAPI.
		GetDeviceConfigurationTemplateProfile(ctx, id).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	httpRes, err = client.DeviceConfigurationTemplateAPI.
		DeleteDeviceConfigurationTemplateProfile(ctx, id).
		IfMatch(existing.Revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Device configuration template profile %s deleted", profileId)
	return nil
}

func DeviceConfigurationTemplateProfileRender(ctx context.Context, profileId string, config []byte) error {
	logger.Get().Info().Msgf("Rendering device configuration template profile %s", profileId)

	id, err := parseId(profileId)
	if err != nil {
		return err
	}

	var payload sdk.RenderDeviceConfigurationTemplateProfile
	if err := utils.UnmarshalContent(config, &payload); err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	rendered, httpRes, err := client.DeviceConfigurationTemplateAPI.
		RenderDeviceConfigurationTemplateProfile(ctx, id).
		RenderDeviceConfigurationTemplateProfile(payload).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(rendered, &renderedTemplatePrintConfig)
}

func DeviceConfigurationTemplateProfileFindApplicable(ctx context.Context, config []byte) error {
	logger.Get().Info().Msgf("Finding applicable device configuration template profiles")

	var payload sdk.FindApplicableDeviceConfigurationTemplateProfiles
	if err := utils.UnmarshalContent(config, &payload); err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	result, httpRes, err := client.DeviceConfigurationTemplateAPI.
		FindApplicableDeviceConfigurationTemplateProfiles(ctx).
		FindApplicableDeviceConfigurationTemplateProfiles(payload).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(result, nil)
}

func DeviceConfigurationTemplateProfileRenderApplicable(ctx context.Context, config []byte) error {
	logger.Get().Info().Msgf("Rendering applicable device configuration template profiles")

	var payload sdk.RenderApplicableDeviceConfigurationTemplateProfiles
	if err := utils.UnmarshalContent(config, &payload); err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	result, httpRes, err := client.DeviceConfigurationTemplateAPI.
		RenderApplicableDeviceConfigurationTemplateProfiles(ctx).
		RenderApplicableDeviceConfigurationTemplateProfiles(payload).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(result, nil)
}

func DeviceConfigurationTemplateProfileBulkAssign(ctx context.Context, config []byte) error {
	logger.Get().Info().Msgf("Bulk assigning device configuration template profiles")

	var payload sdk.BulkAssignDeviceConfigurationTemplateProfile
	if err := utils.UnmarshalContent(config, &payload); err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	result, httpRes, err := client.DeviceConfigurationTemplateAPI.
		BulkAssignDeviceConfigurationTemplateProfiles(ctx).
		BulkAssignDeviceConfigurationTemplateProfile(payload).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(result, nil)
}

func DeviceConfigurationTemplateProfileConfigExample(ctx context.Context) error {
	example := sdk.CreateDeviceConfigurationTemplateProfile{
		DeviceConfigurationTemplateId: 1,
		NetworkDeviceId:               *sdk.NewNullableInt64(sdk.PtrInt64(100)),
		LifecycleStage:                ptrLifecycleStage(ValidLifecycleStages[1]),
		Variables: map[string]interface{}{
			"hostname": "switch-01",
		},
		IsEnabled: sdk.PtrBool(true),
		Priority:  sdk.PtrFloat32(100),
		ApplyMode: ptrApplyMode(ValidApplyModes[0]),
		Tags:      []string{"example"},
	}

	return formatter.PrintResult(example, nil)
}

// ----------------------------------------------------------------------------
// Helpers
// ----------------------------------------------------------------------------

func parseId(id string) (int64, error) {
	idNumeric, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		err := fmt.Errorf("invalid device configuration template ID: '%s'", id)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return idNumeric, nil
}

func ptrLifecycleStage(value string) *sdk.DeviceConfigurationProfileLifecycleStage {
	stage := sdk.DeviceConfigurationProfileLifecycleStage(value)
	return &stage
}

func ptrApplyMode(value string) *sdk.DeviceConfigurationProfileApplyMode {
	mode := sdk.DeviceConfigurationProfileApplyMode(value)
	return &mode
}
