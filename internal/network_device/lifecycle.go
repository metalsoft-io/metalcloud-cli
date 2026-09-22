package network_device

import (
	"context"
	"fmt"
	"net/http"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

// NetworkDeviceReplace swaps a network device with a replacement device,
// moving its configuration and connections over.
func NetworkDeviceReplace(ctx context.Context, networkDeviceRef string, newNetworkDeviceRef string) error {
	logger.Get().Info().Msgf("Replacing network device '%s' with '%s'", networkDeviceRef, newNetworkDeviceRef)

	networkDeviceIdNumeric, err := resolveNetworkDeviceId(ctx, networkDeviceRef)
	if err != nil {
		return err
	}

	newNetworkDeviceIdNumeric, err := resolveNetworkDeviceId(ctx, newNetworkDeviceRef)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	result, httpRes, err := client.NetworkDeviceAPI.
		ReplaceNetworkDevice(ctx, networkDeviceIdNumeric).
		SwitchReplace(sdk.SwitchReplace{NewSwitchId: newNetworkDeviceIdNumeric}).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(result, nil)
}

// NetworkDeviceReProvision re-runs provisioning on a network device.
func NetworkDeviceReProvision(ctx context.Context, networkDeviceRef string, reprovisionType string) error {
	logger.Get().Info().Msgf("Re-provisioning network device '%s' (%s)", networkDeviceRef, reprovisionType)

	networkDeviceIdNumeric, err := resolveNetworkDeviceId(ctx, networkDeviceRef)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	jobInfo, httpRes, err := client.NetworkDeviceAPI.
		ReProvisionNetworkDevice(ctx, networkDeviceIdNumeric).
		NetworkEquipmentReprovision(sdk.NetworkEquipmentReprovision{ReprovisionType: reprovisionType}).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(jobInfo, nil)
}

// NetworkDeviceReturnToPlanned moves an archived network device back to the
// planned state so it can be re-installed. The optional configuration supplies
// the identifying MAC/serial and an OS template override.
func NetworkDeviceReturnToPlanned(ctx context.Context, networkDeviceRef string, config []byte) error {
	logger.Get().Info().Msgf("Returning network device '%s' to planned", networkDeviceRef)

	networkDeviceIdNumeric, revision, err := resolveNetworkDeviceIdAndRevision(ctx, networkDeviceRef)
	if err != nil {
		return err
	}

	// The body is always sent: skipping the optional setter would serialize a
	// literal null, which the API rejects with 400.
	returnConfig := sdk.ReturnNetworkDeviceToPlanned{}
	if len(config) > 0 {
		if err := utils.UnmarshalContent(config, &returnConfig); err != nil {
			return err
		}
	}

	client := api.GetApiClient(ctx)

	networkDevice, httpRes, err := client.NetworkDeviceAPI.
		ReturnNetworkDeviceToPlanned(ctx, networkDeviceIdNumeric).
		ReturnNetworkDeviceToPlanned(returnConfig).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(networkDevice, &NetworkDevicePrintConfig)
}

func NetworkDeviceReturnToPlannedConfigExample(ctx context.Context) error {
	returnConfig := sdk.ReturnNetworkDeviceToPlanned{
		ManagementMAC: sdk.PtrString("AA:BB:CC:DD:EE:FF"),
		SerialNumber:  sdk.PtrString("1234567890"),
		OsTemplateId:  sdk.PtrInt64(10),
	}

	return formatter.PrintResult(returnConfig, nil)
}

// NetworkDeviceRevertDefectiveState takes a network device out of the
// defective state and back to its previous status. The API action was renamed
// from revert-failed-state to revert-defective-state and the generated SDK
// still carries the old path, so the request is issued directly.
func NetworkDeviceRevertDefectiveState(ctx context.Context, networkDeviceRef string) error {
	logger.Get().Info().Msgf("Reverting defective state of network device '%s'", networkDeviceRef)

	networkDeviceIdNumeric, revision, err := resolveNetworkDeviceIdAndRevision(ctx, networkDeviceRef)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/api/v2/network-devices/%d/actions/revert-defective-state", networkDeviceIdNumeric)
	body, err := api.RawJSONRequest(ctx, http.MethodPost, path, nil, api.IfMatchHeader(revision))
	if err != nil {
		return err
	}

	return utils.PrintRawObject(body, &NetworkDevicePrintConfig)
}

// NetworkDeviceStartRegistration starts the onboarding registration of a
// network device.
func NetworkDeviceStartRegistration(ctx context.Context, networkDeviceRef string) error {
	logger.Get().Info().Msgf("Starting registration of network device '%s'", networkDeviceRef)

	networkDeviceIdNumeric, revision, err := resolveNetworkDeviceIdAndRevision(ctx, networkDeviceRef)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	networkDevice, httpRes, err := client.NetworkDeviceAPI.
		StartNetworkDeviceRegistration(ctx, networkDeviceIdNumeric).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(networkDevice, &NetworkDevicePrintConfig)
}

// NetworkDeviceMarkInstallationReady marks the physical installation of a
// network device as complete, allowing onboarding to continue.
func NetworkDeviceMarkInstallationReady(ctx context.Context, networkDeviceRef string) error {
	logger.Get().Info().Msgf("Marking network device '%s' installation as ready", networkDeviceRef)

	networkDeviceIdNumeric, revision, err := resolveNetworkDeviceIdAndRevision(ctx, networkDeviceRef)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	networkDevice, httpRes, err := client.NetworkDeviceAPI.
		MarkNetworkDeviceInstallationReady(ctx, networkDeviceIdNumeric).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(networkDevice, &NetworkDevicePrintConfig)
}

// NetworkDeviceRunExtension runs an extension against a network device.
func NetworkDeviceRunExtension(ctx context.Context, networkDeviceRef string, config []byte) error {
	logger.Get().Info().Msgf("Running extension on network device '%s'", networkDeviceRef)

	networkDeviceIdNumeric, revision, err := resolveNetworkDeviceIdAndRevision(ctx, networkDeviceRef)
	if err != nil {
		return err
	}

	var extensionConfig sdk.RunExtensionOnPhysicalDevice
	if err := utils.UnmarshalContent(config, &extensionConfig); err != nil {
		return err
	}
	if extensionConfig.InputArguments == nil {
		extensionConfig.InputArguments = map[string]interface{}{}
	}

	client := api.GetApiClient(ctx)

	jobInfo, httpRes, err := client.NetworkDeviceAPI.
		RunExtensionOnNetworkDevice(ctx, networkDeviceIdNumeric).
		RunExtensionOnPhysicalDevice(extensionConfig).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(jobInfo, nil)
}

func NetworkDeviceRunExtensionConfigExample(ctx context.Context) error {
	extensionConfig := sdk.RunExtensionOnPhysicalDevice{
		ExtensionId:    1,
		InputArguments: map[string]interface{}{"argument1": "value1"},
	}

	return formatter.PrintResult(extensionConfig, nil)
}
