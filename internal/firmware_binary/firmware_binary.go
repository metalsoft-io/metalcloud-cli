package firmware_binary

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var firmwareBinaryPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"Name": {
			MaxWidth: 30,
			Order:    2,
		},
		"CatalogId": {
			Title: "Catalog ID",
			Order: 3,
		},
		"PackageVersion": {
			Title: "Version",
			Order: 4,
		},
		"UpdateSeverity": {
			Title: "Severity",
			Order: 5,
		},
		"RebootRequired": {
			Title: "Reboot Required",
			Order: 6,
		},
		"VendorReleaseTimestamp": {
			Title:       "Vendor Release",
			Transformer: formatter.FormatDateTimeValue,
			Order:       7,
		},
		"PackageId": {
			Title: "Package ID",
			Order: 8,
		},
		"ExternalId": {
			Title: "External ID",
			Order: 9,
		},
		"VendorDownloadUrl": {
			Title: "Download URL",
			Order: 10,
		},
		"VendorInfoUrl": {
			Title: "Info URL",
			Order: 11,
		},
		"ComponentClass": {
			Title: "Class",
			Order: 12,
		},
		"IsBundle": {
			Title: "Bundle",
			Order: 13,
		},
	},
}

// The steps of a bundle, in order, as `firmware-binary get` lists them.
var firmwareBundleStepPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id":             {Title: "Step", Order: 1},
		"ComponentClass": {Title: "Class", Order: 2},
		"TargetPattern":  {Title: "Target Pattern", Order: 3, MaxWidth: 32},
		"PackageFile":    {Title: "Package File", Order: 4, MaxWidth: 40},
		"PackageVersion": {Title: "Package Version", Order: 5},
		"VerifyVersion":  {Title: "Verify Version", Order: 6},
		"Activation":     {Title: "Activation", Order: 7},
		"DependsOnStep":  {Title: "After", Order: 8},
	},
}

var firmwarePrerequisitePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Platform":       {Title: "Platform", Order: 1},
		"ComponentClass": {Title: "Class", Order: 2},
		"ComponentModel": {Title: "Model", Order: 3},
		"ComponentName":  {Title: "Name", Order: 4},
		"MinimumVersion": {Title: "Minimum Version", Order: 5},
	},
}

// ListFlags holds the server-side filters and pagination options of `firmware-binary list`.
// Every filter maps 1:1 to a `filter.*` query parameter of GET /api/v2/firmware/binary.
type ListFlags struct {
	FilterId                     []string
	FilterCatalogId              []string
	FilterExternalId             []string
	FilterPackageId              []string
	FilterPackageVersion         []string
	FilterUpdateSeverity         []string
	FilterRebootRequired         []string
	FilterCreatedTimestamp       []string
	FilterVendorReleaseTimestamp []string
	FilterIsBundle               []string
	FilterComponentClass         []string

	Page  int
	Limit int
}

func FirmwareBinaryList(ctx context.Context, flags ListFlags) error {
	logger.Get().Info().Msgf("Listing firmware binaries")

	client := api.GetApiClient(ctx)

	request := client.FirmwareBinaryAPI.GetFirmwareBinaries(ctx).SortBy([]string{"id:ASC"})

	apply := func(values []string, set func([]string) sdk.FirmwareBinaryAPIGetFirmwareBinariesRequest) {
		if len(values) > 0 {
			request = set(utils.ProcessFilterStringSlice(values))
		}
	}

	apply(flags.FilterId, request.FilterId)
	apply(flags.FilterCatalogId, request.FilterCatalogId)
	apply(flags.FilterExternalId, request.FilterExternalId)
	apply(flags.FilterPackageId, request.FilterPackageId)
	apply(flags.FilterPackageVersion, request.FilterPackageVersion)
	apply(flags.FilterUpdateSeverity, request.FilterUpdateSeverity)
	apply(flags.FilterRebootRequired, request.FilterRebootRequired)
	apply(flags.FilterCreatedTimestamp, request.FilterCreatedTimestamp)
	apply(flags.FilterVendorReleaseTimestamp, request.FilterVendorReleaseTimestamp)
	apply(flags.FilterIsBundle, request.FilterIsBundle)
	apply(flags.FilterComponentClass, request.FilterComponentClass)

	switch {
	case flags.Page > 0:
		records, meta, err := utils.FetchPageWindow(request, flags.Page, flags.Limit)
		if err != nil {
			return err
		}
		return utils.PrintAll(records, meta, len(records), &firmwareBinaryPrintConfig)
	case flags.Limit > 0:
		records, meta, err := utils.FetchUpTo(request, flags.Limit)
		if err != nil {
			return err
		}
		return utils.PrintAll(records, meta, len(records), &firmwareBinaryPrintConfig)
	default:
		records, meta, err := utils.FetchAllPages(request)
		if err != nil {
			return err
		}
		return utils.PrintAll(records, meta, len(records), &firmwareBinaryPrintConfig)
	}
}

func FirmwareBinaryGet(ctx context.Context, firmwareBinaryId string) error {
	logger.Get().Info().Msgf("Get firmware binary '%s' details", firmwareBinaryId)

	firmwareBinaryIdNumeric, err := getFirmwareBinaryId(firmwareBinaryId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	firmwareBinary, httpRes, err := client.FirmwareBinaryAPI.GetFirmwareBinary(ctx, firmwareBinaryIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	if err := formatter.PrintResult(firmwareBinary, &firmwareBinaryPrintConfig); err != nil {
		return err
	}
	// JSON and YAML already carry the steps and prerequisites; the text table lists them below
	if !formatter.IsTextFormat() {
		return nil
	}
	if len(firmwareBinary.BundleSteps) > 0 {
		fmt.Printf("\nBundle steps (%d):\n", len(firmwareBinary.BundleSteps))
		if err := formatter.PrintResult(firmwareBinary.BundleSteps, &firmwareBundleStepPrintConfig); err != nil {
			return err
		}
	}
	if len(firmwareBinary.Prerequisites) > 0 {
		fmt.Printf("\nPrerequisites (%d):\n", len(firmwareBinary.Prerequisites))
		if err := formatter.PrintResult(firmwareBinary.Prerequisites, &firmwarePrerequisitePrintConfig); err != nil {
			return err
		}
	}
	return nil
}

// FirmwareBinaryUpdate changes the fields of a firmware binary given in the configuration and keeps the others. The
// API replaces the whole binary (PUT), so the current binary is read and the configuration is applied over it.
func FirmwareBinaryUpdate(ctx context.Context, firmwareBinaryId string, config []byte) error {
	logger.Get().Info().Msgf("Updating firmware binary '%s'", firmwareBinaryId)

	firmwareBinaryIdNumeric, err := getFirmwareBinaryId(firmwareBinaryId)
	if err != nil {
		return err
	}

	var changes map[string]interface{}
	if err := utils.UnmarshalContent(config, &changes); err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	current, httpRes, err := client.FirmwareBinaryAPI.GetFirmwareBinary(ctx, firmwareBinaryIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	update, err := mergeFirmwareBinaryUpdate(current, changes)
	if err != nil {
		return err
	}

	firmwareBinary, httpRes, err := client.FirmwareBinaryAPI.
		SetFirmwareBinary(ctx, firmwareBinaryIdNumeric).
		CreateFirmwareBinary(update).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(firmwareBinary, &firmwareBinaryPrintConfig)
}

// mergeFirmwareBinaryUpdate applies the changed fields over the current binary. A field set to null clears it.
func mergeFirmwareBinaryUpdate(current *sdk.FirmwareBinary, changes map[string]interface{}) (sdk.CreateFirmwareBinary, error) {
	var update sdk.CreateFirmwareBinary

	currentJson, err := json.Marshal(current)
	if err != nil {
		return update, err
	}
	var merged map[string]interface{}
	if err := json.Unmarshal(currentJson, &merged); err != nil {
		return update, err
	}
	// Read-only fields of the binary, not part of the update
	for _, readOnly := range []string{"id", "links", "createdTimestamp", "updatedTimestamp"} {
		delete(merged, readOnly)
	}
	for key, value := range changes {
		if value == nil {
			delete(merged, key)
		} else {
			merged[key] = value
		}
	}

	mergedJson, err := json.Marshal(merged)
	if err != nil {
		return update, err
	}
	if err := json.Unmarshal(mergedJson, &update); err != nil {
		return update, fmt.Errorf("invalid firmware binary configuration: %w", err)
	}
	// Fields the update does not have: a changed one is a mistake (the API would reject it); one the binary was
	// read with (e.g. its vendor) is not part of the update
	for key := range update.AdditionalProperties {
		if _, changed := changes[key]; changed {
			return update, fmt.Errorf("unknown firmware binary field %q", key)
		}
	}
	update.AdditionalProperties = nil
	return update, nil
}

func FirmwareBinaryCreate(ctx context.Context, config []byte) error {
	logger.Get().Info().Msgf("Creating firmware binary")

	var firmwareBinaryConfig sdk.CreateFirmwareBinary
	err := utils.UnmarshalContent(config, &firmwareBinaryConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	firmwareBinary, httpRes, err := client.FirmwareBinaryAPI.
		CreateFirmwareBinary(ctx).
		CreateFirmwareBinary(firmwareBinaryConfig).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(firmwareBinary, &firmwareBinaryPrintConfig)
}

func FirmwareBinaryDelete(ctx context.Context, firmwareBinaryId string) error {
	logger.Get().Info().Msgf("Deleting firmware binary '%s'", firmwareBinaryId)

	firmwareBinaryIdNumeric, err := getFirmwareBinaryId(firmwareBinaryId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.FirmwareBinaryAPI.
		DeleteFirmwareBinary(ctx, firmwareBinaryIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Firmware binary '%s' deleted", firmwareBinaryId)
	return nil
}

// FirmwareBundleConfigExample prints the configuration of a firmware bundle: its ordered steps and prerequisites.
func FirmwareBundleConfigExample(ctx context.Context) error {
	return formatter.PrintResult(firmwareBundleConfigExample(), nil)
}

func FirmwareBinaryConfigExample(ctx context.Context) error {
	// Example create firmware binary configuration
	exampleReleaseTimestamp := time.Date(2024, time.April, 1, 12, 0, 0, 0, time.UTC)
	firmwareBinaryConfiguration := sdk.CreateFirmwareBinary{
		Name:                   "BIOS-R740-2.15.0",
		CatalogId:              1,
		VendorDownloadUrl:      "https://dell.com/downloads/firmware/R740/BIOS-2.15.0.bin",
		VendorInfoUrl:          sdk.PtrString("https://dell.com/support/firmware/R740/BIOS/2.15.0"),
		ExternalId:             sdk.PtrString("DELL-R740-BIOS-2.15.0"),
		PackageId:              sdk.PtrString("BIOS"),
		PackageVersion:         sdk.PtrString("2.15.0"),
		RebootRequired:         true,
		UpdateSeverity:         "recommended",
		VendorReleaseTimestamp: &exampleReleaseTimestamp,
		// BIOS packages have no PCI ids; for device firmware (NIC, RAID, GPU, ...)
		// list the PCI ids of the supported devices, e.g.
		// {Pci: &sdk.PciIds{VendorId: sdk.PtrString("14e4"), DeviceId: sdk.PtrString("16d7")}}.
		VendorSupportedDevices: []sdk.VendorSupportedDevice{},
		VendorSupportedSystems: []sdk.ServerFirmwareBinaryVendorSupportedSystems{
			{Model: sdk.PtrString("R740")},
		},
	}

	return formatter.PrintResult(firmwareBinaryConfiguration, nil)
}

func getFirmwareBinaryId(firmwareBinaryId string) (int64, error) {
	firmwareBinaryIdNumeric, err := strconv.ParseInt(firmwareBinaryId, 10, 64)
	if err != nil {
		err := fmt.Errorf("invalid firmware binary ID: '%s'", firmwareBinaryId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return firmwareBinaryIdNumeric, nil
}

// firmwareBundleConfigExample is an NVIDIA HGX bundle: its steps, in order, and the host firmware it needs first.
// Bundles are usually imported from a manifest with `firmware-catalog create --vendor nvidia`.
func firmwareBundleConfigExample() sdk.CreateFirmwareBinary {
	releaseTimestamp := time.Date(2025, time.June, 12, 0, 0, 0, 0, time.UTC)
	step := func(id string, class sdk.ServerComponentClass, pattern string, activation sdk.FirmwareBundleStepActivation, verifyVersion string, dependsOn *string) sdk.FirmwareBundleStep {
		result := sdk.FirmwareBundleStep{
			Id:             id,
			ComponentClass: class,
			TargetPattern:  pattern,
			PackageFile:    "nvfw_HGX_H100x8_250612.fwpkg",
			PackageVersion: "25.06.1",
			Sha256:         "f013d7e18fbde4f8f4cd8fce030ec358d5b50c02d2f711ce3dca0179dfe979c2",
			Activation:     activation,
			DependsOnStep:  *sdk.NewNullableString(dependsOn),
		}
		if verifyVersion != "" {
			// The version the component reports once the step is active: what the bundle is offered as
			result.VerifyVersion = sdk.PtrString(verifyVersion)
		}
		return result
	}

	return sdk.CreateFirmwareBinary{
		Name:                   "NVIDIA HGX H100 8-GPU 25.06.1",
		CatalogId:              1,
		ExternalId:             sdk.PtrString("nvfw_HGX_H100x8_250612.fwpkg"),
		PackageId:              sdk.PtrString("HGX H100 8-GPU"),
		PackageVersion:         sdk.PtrString("25.06.1"),
		VendorDownloadUrl:      "https://repo.example.com/nvidia/hgx-h100/25.06.1/nvfw_HGX_H100x8_250612.fwpkg",
		VendorInfoUrl:          sdk.PtrString("https://docs.nvidia.com/dgx/dgxh100-fw-rn/index.html"),
		RebootRequired:         true,
		UpdateSeverity:         "recommended",
		VendorReleaseTimestamp: &releaseTimestamp,
		VendorSupportedDevices: []sdk.VendorSupportedDevice{},
		VendorSupportedSystems: []sdk.ServerFirmwareBinaryVendorSupportedSystems{
			{AdditionalProperties: map[string]interface{}{"id": "NVIDIA DGXH100"}},
		},
		IsBundle: sdk.PtrBool(true),
		Sha256:   sdk.PtrString("f013d7e18fbde4f8f4cd8fce030ec358d5b50c02d2f711ce3dca0179dfe979c2"),
		BundleSteps: []sdk.FirmwareBundleStep{
			step("hmc", sdk.SERVERCOMPONENTCLASS_HGX_HMC, "^HGX_FW_(HMC|BMC)_0$", sdk.FIRMWAREBUNDLESTEPACTIVATION_HGX_RESET, "HGX-22.10-1-rc80", nil),
			step("gpu", sdk.SERVERCOMPONENTCLASS_GPU, "^HGX_FW_GPU_SXM_\\d+$", sdk.FIRMWAREBUNDLESTEPACTIVATION_CHASSIS_POWER_CYCLE, "96.00.A5.00.01", sdk.PtrString("hmc")),
		},
		Prerequisites: []sdk.FirmwarePrerequisite{
			{Platform: sdk.PtrString("NVIDIA DGXH100"), ComponentClass: sdk.SERVERCOMPONENTCLASS_BIOS, MinimumVersion: "01.02.05"},
			{ComponentClass: sdk.SERVERCOMPONENTCLASS_BMC, MinimumVersion: "23.06.00"},
		},
	}
}
