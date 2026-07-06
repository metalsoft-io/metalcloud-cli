package license

import (
	"context"
	"fmt"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var licenseStatusPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Valid":      {Title: "Valid", Transformer: formatter.FormatBooleanValue, Order: 1},
		"Expiration": {Title: "Expiration", Transformer: formatter.FormatDateTimeValue, Order: 2},
		"Hostname":   {Title: "Hostname", Order: 3},
		"Signature":  {Title: "Signature", MaxWidth: 40, Order: 4},
	},
}

var licenseAllowancePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Servers":   {Title: "Servers", Order: 1},
		"Switches":  {Title: "Switches", Order: 2},
		"Devices":   {Title: "Devices", Order: 3},
		"StorageGB": {Title: "Storage (GB)", Order: 4},
	},
}

var licensedProductsPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"ServersAreLicensed":  {Title: "Servers", Transformer: formatter.FormatBooleanValue, Order: 1},
		"SwitchesAreLicensed": {Title: "Switches", Transformer: formatter.FormatBooleanValue, Order: 2},
		"VmsAreLicensed":      {Title: "VMs", Transformer: formatter.FormatBooleanValue, Order: 3},
		"StoragesAreLicensed": {Title: "Storage", Transformer: formatter.FormatBooleanValue, Order: 4},
	},
}

var addLicenseResponsePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Ok": {Title: "Installed", Transformer: formatter.FormatBooleanValue, Order: 1},
	},
}

func LicenseStatus(ctx context.Context) error {
	logger.Get().Info().Msg("Getting license status")

	client := api.GetApiClient(ctx)

	status, httpRes, err := client.SystemAPI.GetLicenseStatus(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(status, &licenseStatusPrintConfig)
}

func LicenseAllowance(ctx context.Context) error {
	logger.Get().Info().Msg("Getting license allowance")

	client := api.GetApiClient(ctx)

	allowance, httpRes, err := client.SystemAPI.GetLicenseAllowance(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(allowance, &licenseAllowancePrintConfig)
}

func LicensedProducts(ctx context.Context) error {
	logger.Get().Info().Msg("Getting licensed products")

	client := api.GetApiClient(ctx)

	products, httpRes, err := client.SystemAPI.GetLicensedProducts(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(products, &licensedProductsPrintConfig)
}

// LicenseGet returns the license document currently installed on the system.
// The document is a Base64-encoded signed blob; in text output it is printed
// verbatim so it can be piped or saved, while structured formats wrap it.
func LicenseGet(ctx context.Context) error {
	logger.Get().Info().Msg("Getting current license")

	client := api.GetApiClient(ctx)

	license, httpRes, err := client.SystemAPI.GetLicense(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	if formatter.IsTextFormat() {
		fmt.Println(license.License)
		return nil
	}

	return formatter.PrintResult(license, &formatter.PrintConfig{
		FieldsConfig: map[string]formatter.RecordFieldConfig{
			"License": {Title: "License", Order: 1},
		},
	})
}

// LicenseRequest returns the Base64-encoded license request document that must
// be sent to MetalSoft in order to obtain a license. In text output it is
// printed verbatim so it can be piped or saved.
func LicenseRequest(ctx context.Context) error {
	logger.Get().Info().Msg("Getting license request")

	client := api.GetApiClient(ctx)

	request, httpRes, err := client.SystemAPI.GetLicenseRequest(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	if formatter.IsTextFormat() {
		fmt.Println(request.LicenseRequest)
		return nil
	}

	return formatter.PrintResult(request, &formatter.PrintConfig{
		FieldsConfig: map[string]formatter.RecordFieldConfig{
			"LicenseRequest": {Title: "License Request", Order: 1},
		},
	})
}

// LicenseAdd installs a license on the system. The license must be provided as
// a Base64-encoded signed document; it is forwarded verbatim to preserve the
// original signed bytes.
func LicenseAdd(ctx context.Context, license string) error {
	logger.Get().Info().Msg("Adding license")

	client := api.GetApiClient(ctx)

	payload := sdk.NewAddLicense(license)

	response, httpRes, err := client.SystemAPI.AddLicense(ctx).AddLicense(*payload).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(response, &addLicenseResponsePrintConfig)
}
