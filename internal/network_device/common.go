package network_device

import (
	"context"
	"fmt"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

// printRecords renders the full API objects for json/yaml output and the
// table-friendly rows for the tabular formats, mirroring utils.PrintAllRaw.
// The tabular formatter renders struct-valued fields as empty cells, so
// records carrying nested objects (breakout groups, OID groups, detected
// issues) are flattened into rows before being tabulated.
func printRecords(records interface{}, rows interface{}, printConfig *formatter.PrintConfig) error {
	if formatter.IsNativeFormat() {
		return formatter.PrintResult(records, nil)
	}

	return formatter.PrintResult(rows, printConfig)
}

// resolveNetworkDeviceId resolves a network device reference (numeric id or
// identifierString label) to its numeric id. A numeric reference is parsed
// locally so the common case costs no extra API call.
func resolveNetworkDeviceId(ctx context.Context, networkDeviceRef string) (int64, error) {
	if networkDeviceIdNumeric, err := strconv.ParseInt(networkDeviceRef, 10, 64); err == nil {
		return networkDeviceIdNumeric, nil
	}

	networkDevice, err := GetNetworkDeviceByIdOrLabel(ctx, networkDeviceRef)
	if err != nil {
		return 0, err
	}

	return networkDeviceNumericId(networkDevice)
}

// resolveNetworkDeviceIdAndRevision resolves a network device reference and
// returns its numeric id together with the current entity revision, for
// requests guarded by optimistic concurrency (If-Match).
func resolveNetworkDeviceIdAndRevision(ctx context.Context, networkDeviceRef string) (int64, string, error) {
	networkDevice, err := GetNetworkDeviceByIdOrLabel(ctx, networkDeviceRef)
	if err != nil {
		return 0, "", err
	}

	networkDeviceIdNumeric, err := networkDeviceNumericId(networkDevice)
	if err != nil {
		return 0, "", err
	}

	return networkDeviceIdNumeric, strconv.FormatInt(networkDevice.Revision, 10), nil
}

// networkDeviceNumericId converts the string id carried by sdk.NetworkDevice
// into the int64 the SDK path parameters expect.
func networkDeviceNumericId(networkDevice *sdk.NetworkDevice) (int64, error) {
	networkDeviceIdNumeric, err := strconv.ParseInt(networkDevice.Id, 10, 64)
	if err != nil {
		err := fmt.Errorf("unexpected network device ID: '%s'", networkDevice.Id)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return networkDeviceIdNumeric, nil
}

// errEmptyBody reports an endpoint that answered without the object the caller
// needs, for example a 204 where a resource was expected.
func errEmptyBody(what string) error {
	return fmt.Errorf("the API returned no %s", what)
}

// validateAddressFamily checks the family path segment accepted by the port IP
// endpoints.
func validateAddressFamily(family string) error {
	if family != "ipv4" && family != "ipv6" {
		return fmt.Errorf("invalid address family: '%s'. Valid families are: ipv4, ipv6", family)
	}

	return nil
}
