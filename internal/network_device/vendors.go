package network_device

import (
	"context"
	"fmt"
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

var networkDeviceVendorPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"Kind": {
			Title: "Driver",
			Order: 2,
		},
		"OidGroups": {
			Title:    "OID Groups",
			MaxWidth: 40,
			Order:    3,
		},
		"HealthCheckRules": {
			Title:    "Health Rules",
			MaxWidth: 20,
			Order:    4,
		},
		"OptionalFilesToBackup": {
			Title:    "Backup File Sets",
			MaxWidth: 40,
			Order:    5,
		},
		"Revision": {
			Order: 6,
		},
	},
}

var networkDeviceDriverCapabilitiesPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Driver": {
			Order: 1,
		},
		"ActivationActions": {
			Title:       "Activation Actions",
			MaxWidth:    60,
			Transformer: formatNetworkDeviceStringList,
			Order:       2,
		},
	},
}

// vendorRow is the table projection of a vendor profile: the nested OID group,
// health rule and backup file objects are summarized because the tabular
// formatter renders struct-valued fields as empty cells.
type vendorRow struct {
	Id                    int64
	Kind                  string
	OidGroups             string
	HealthCheckRules      string
	OptionalFilesToBackup string
	Revision              int64
}

func toVendorRow(vendor sdk.NetworkDeviceVendors) vendorRow {
	oidGroups := make([]string, 0, len(vendor.OidGroups))
	for _, oidGroup := range vendor.OidGroups {
		oidGroups = append(oidGroups, oidGroup.Name)
	}

	fileSets := make([]string, 0, len(vendor.OptionalFilesToBackup.Files))
	for label := range vendor.OptionalFilesToBackup.Files {
		fileSets = append(fileSets, label)
	}
	sort.Strings(fileSets)

	return vendorRow{
		Id:                    vendor.Id,
		Kind:                  vendor.Kind,
		OidGroups:             strings.Join(oidGroups, ", "),
		HealthCheckRules:      fmt.Sprintf("%d rule(s)", len(vendor.HealthCheckRules.Rules)),
		OptionalFilesToBackup: strings.Join(fileSets, ", "),
		Revision:              vendor.Revision,
	}
}

// NetworkDeviceVendorList lists the per-vendor monitoring and backup profiles.
func NetworkDeviceVendorList(ctx context.Context, filterKind []string) error {
	logger.Get().Info().Msgf("Listing network device vendors")

	client := api.GetApiClient(ctx)

	request := client.NetworkDeviceAPI.GetNetworkDeviceVendors(ctx).SortBy([]string{"id:ASC"})
	if len(filterKind) > 0 {
		request = request.FilterKind(utils.ProcessFilterStringSlice(filterKind))
	}

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	rows := make([]vendorRow, 0, len(records))
	for _, vendor := range records {
		rows = append(rows, toVendorRow(vendor))
	}

	if formatter.IsNativeFormat() {
		return utils.PrintAll(records, meta, len(records), nil)
	}

	return utils.PrintAll(rows, meta, len(rows), &networkDeviceVendorPrintConfig)
}

// NetworkDeviceVendorGet shows one vendor profile.
func NetworkDeviceVendorGet(ctx context.Context, vendorId string) error {
	logger.Get().Info().Msgf("Getting network device vendor '%s'", vendorId)

	vendor, err := getNetworkDeviceVendor(ctx, vendorId)
	if err != nil {
		return err
	}

	return printRecords(vendor, toVendorRow(*vendor), &networkDeviceVendorPrintConfig)
}

// NetworkDeviceVendorUpdate updates the SNMP OID groups, health check rules
// and backup file list of one vendor profile.
func NetworkDeviceVendorUpdate(ctx context.Context, vendorId string, config []byte) error {
	logger.Get().Info().Msgf("Updating network device vendor '%s'", vendorId)

	vendor, err := getNetworkDeviceVendor(ctx, vendorId)
	if err != nil {
		return err
	}

	var vendorConfig sdk.UpdateNetworkDeviceVendorsDto
	if err := utils.UnmarshalContent(config, &vendorConfig); err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	updatedVendor, httpRes, err := client.NetworkDeviceAPI.
		UpdateNetworkDeviceVendor(ctx, vendor.Id).
		UpdateNetworkDeviceVendorsDto(vendorConfig).
		IfMatch(strconv.FormatInt(vendor.Revision, 10)).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	if updatedVendor == nil {
		return errEmptyBody("vendor profile")
	}

	return printRecords(updatedVendor, toVendorRow(*updatedVendor), &networkDeviceVendorPrintConfig)
}

func NetworkDeviceVendorConfigExample(ctx context.Context) error {
	vendorConfig := sdk.UpdateNetworkDeviceVendorsDto{
		OidGroups: []sdk.SNMPOIDGroupConfig{
			{
				Name: "interfaces",
				Oids: map[string]interface{}{"1.3.6.1.2.1.2.2.1.10": "ifInOctets"},
				Mapping: sdk.SNMPOIDToMetricMapping{
					Metric: "interface_in_octets",
					Type:   sdk.NETWORKDEVICEMETRICVALUETYPE_INTEGER,
				},
			},
		},
		OptionalFilesToBackup: &sdk.NetworkDeviceOptionalFilesToBackup{
			Files: map[string]interface{}{"default": []string{"/etc/sonic/config_db.json"}},
		},
	}

	return formatter.PrintResult(vendorConfig, nil)
}

// NetworkDeviceDriverCapabilities lists, per network device driver, the
// onboarding actions the driver supports.
func NetworkDeviceDriverCapabilities(ctx context.Context, search string) error {
	logger.Get().Info().Msgf("Listing network device driver capabilities")

	client := api.GetApiClient(ctx)

	request := client.NetworkDeviceAPI.GetNetworkDeviceDriverCapabilities(ctx)
	if search != "" {
		request = request.Search(search)
	}

	capabilities, httpRes, err := request.Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(capabilities, &networkDeviceDriverCapabilitiesPrintConfig)
}

func getNetworkDeviceVendor(ctx context.Context, vendorId string) (*sdk.NetworkDeviceVendors, error) {
	vendorIdNumeric, err := utils.GetInt64FromString(vendorId)
	if err != nil {
		return nil, err
	}

	client := api.GetApiClient(ctx)

	vendor, httpRes, err := client.NetworkDeviceAPI.GetNetworkDeviceVendor(ctx, vendorIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, err
	}
	if vendor == nil {
		return nil, errEmptyBody("vendor profile")
	}

	return vendor, nil
}
