package network_device_controller

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

var networkDeviceControllerPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"IdentifierString": {
			Title:    "Identifier",
			Order:    2,
			MaxWidth: 30,
		},
		"SiteId": {
			Title: "Site",
			Order: 3,
		},
		"DatacenterName": {
			Title:    "Datacenter",
			Order:    4,
			MaxWidth: 30,
		},
		"Driver": {
			Order: 5,
		},
		"ManagementAddress": {
			Title: "Management Address",
			Order: 6,
		},
		"ManagementPort": {
			Title: "Port",
			Order: 7,
		},
		"Username": {
			Order: 8,
		},
		"Status": {
			Transformer: formatter.FormatStatusValue,
			Order:       9,
		},
		"Revision": {
			Order: 10,
		},
	},
}

var networkDeviceControllerCredentialsPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Hostname": {
			Order:    1,
			MaxWidth: 30,
		},
		"Host": {
			Order: 2,
		},
		"Port": {
			Order: 3,
		},
		"Datacenter": {
			Order:    4,
			MaxWidth: 30,
		},
		"Driver": {
			Order: 5,
		},
		"Username": {
			Order: 6,
		},
		"Password": {
			Order: 7,
		},
	},
}

// NetworkDeviceControllerListFilters carries the optional list filters.
type NetworkDeviceControllerListFilters struct {
	Id                []string
	SiteId            []string
	DatacenterName    []string
	ManagementAddress []string
	IdentifierString  []string
}

func NetworkDeviceControllerList(ctx context.Context, filters NetworkDeviceControllerListFilters) error {
	logger.Get().Info().Msgf("Listing network device controllers")

	client := api.GetApiClient(ctx)

	request := client.NetworkDeviceControllerAPI.GetNetworkDeviceControllers(ctx).SortBy([]string{"id:ASC"})
	if len(filters.Id) > 0 {
		request = request.FilterId(utils.ProcessFilterStringSlice(filters.Id))
	}
	if len(filters.SiteId) > 0 {
		request = request.FilterSiteId(utils.ProcessFilterStringSlice(filters.SiteId))
	}
	if len(filters.DatacenterName) > 0 {
		request = request.FilterDatacenterName(utils.ProcessFilterStringSlice(filters.DatacenterName))
	}
	if len(filters.ManagementAddress) > 0 {
		request = request.FilterManagementAddress(utils.ProcessFilterStringSlice(filters.ManagementAddress))
	}
	if len(filters.IdentifierString) > 0 {
		request = request.FilterIdentifierString(utils.ProcessFilterStringSlice(filters.IdentifierString))
	}

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &networkDeviceControllerPrintConfig)
}

func NetworkDeviceControllerGet(ctx context.Context, controllerIdOrIdentifier string) error {
	logger.Get().Info().Msgf("Get network device controller '%s'", controllerIdOrIdentifier)

	controller, err := GetNetworkDeviceControllerByIdOrLabel(ctx, controllerIdOrIdentifier)
	if err != nil {
		return err
	}

	return formatter.PrintResult(controller, &networkDeviceControllerPrintConfig)
}

func NetworkDeviceControllerConfigExample(ctx context.Context) error {
	example := sdk.CreateNetworkDeviceController{
		SiteId:             sdk.PtrInt64(1),
		DatacenterName:     "dc1",
		IdentifierString:   sdk.PtrString("ndfc-controller-01"),
		Driver:             sdk.SWITCHCONTROLLERDRIVER_CISCO_NDFC,
		ManagementAddress:  "10.0.0.50",
		ManagementPort:     443,
		Username:           "admin",
		ManagementPassword: "password",
		Description:        sdk.PtrString("Cisco NDFC controller for DC1"),
	}

	if formatter.IsNativeFormat() {
		return formatter.PrintResult(example, nil)
	}
	return formatter.PrintYamlResult(example)
}

func NetworkDeviceControllerCreate(ctx context.Context, create sdk.CreateNetworkDeviceController) error {
	logger.Get().Info().Msgf("Creating network device controller '%s'", create.ManagementAddress)

	client := api.GetApiClient(ctx)

	controller, httpRes, err := client.NetworkDeviceControllerAPI.
		CreateNetworkDeviceController(ctx).
		CreateNetworkDeviceController(create).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(controller, &networkDeviceControllerPrintConfig)
}

func NetworkDeviceControllerUpdate(ctx context.Context, controllerIdOrIdentifier string, config []byte) error {
	logger.Get().Info().Msgf("Updating network device controller '%s'", controllerIdOrIdentifier)

	var update sdk.UpdateNetworkDeviceController
	if err := utils.UnmarshalContent(config, &update); err != nil {
		return err
	}

	controllerId, revision, err := resolveControllerIdAndRevision(ctx, controllerIdOrIdentifier)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	controller, httpRes, err := client.NetworkDeviceControllerAPI.
		UpdateNetworkDeviceController(ctx, controllerId).
		UpdateNetworkDeviceController(update).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(controller, &networkDeviceControllerPrintConfig)
}

func NetworkDeviceControllerDelete(ctx context.Context, controllerIdOrIdentifier string) error {
	logger.Get().Info().Msgf("Deleting network device controller '%s'", controllerIdOrIdentifier)

	controllerId, revision, err := resolveControllerIdAndRevision(ctx, controllerIdOrIdentifier)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.NetworkDeviceControllerAPI.
		DeleteNetworkDeviceController(ctx, controllerId).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Network device controller '%s' deleted", controllerIdOrIdentifier)
	return nil
}

func NetworkDeviceControllerGetCredentials(ctx context.Context, controllerIdOrIdentifier string) error {
	logger.Get().Info().Msgf("Get credentials of network device controller '%s'", controllerIdOrIdentifier)

	controllerId, err := ResolveNetworkDeviceControllerNumericId(ctx, controllerIdOrIdentifier)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	credentials, httpRes, err := client.NetworkDeviceControllerAPI.
		GetNetworkDeviceControllerCredentials(ctx, controllerId).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(credentials, &networkDeviceControllerCredentialsPrintConfig)
}

// NetworkDeviceControllerDeployConfirm confirms a pending deploy of the
// controller's configuration. The endpoint returns no body.
func NetworkDeviceControllerDeployConfirm(ctx context.Context, controllerIdOrIdentifier string) error {
	logger.Get().Info().Msgf("Confirming deploy of network device controller '%s'", controllerIdOrIdentifier)

	controllerId, err := ResolveNetworkDeviceControllerNumericId(ctx, controllerIdOrIdentifier)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.NetworkDeviceControllerAPI.
		NetworkDeviceControllerDeployConfirm(ctx, controllerId).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Deploy of network device controller '%s' confirmed", controllerIdOrIdentifier)
	return nil
}

// GetNetworkDeviceControllerByIdOrLabel resolves a controller by numeric ID
// first and falls back to an exact identifier (hostname) match.
func GetNetworkDeviceControllerByIdOrLabel(ctx context.Context, controllerIdOrIdentifier string) (*sdk.NetworkDeviceController, error) {
	client := api.GetApiClient(ctx)

	if idNumeric, err := utils.GetInt64FromString(controllerIdOrIdentifier); err == nil {
		controller, httpRes, err := client.NetworkDeviceControllerAPI.GetNetworkDeviceController(ctx, idNumeric).Execute()
		if err = response_inspector.InspectResponse(httpRes, err); err == nil {
			return controller, nil
		}
	}

	list, httpRes, err := client.NetworkDeviceControllerAPI.
		GetNetworkDeviceControllers(ctx).
		FilterIdentifierString([]string{controllerIdOrIdentifier}).
		Execute()
	if err = response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, err
	}

	for i := range list.Data {
		if list.Data[i].IdentifierString == controllerIdOrIdentifier {
			return &list.Data[i], nil
		}
	}

	err = fmt.Errorf("network device controller '%s' not found", controllerIdOrIdentifier)
	logger.Get().Error().Err(err).Msg("")
	return nil, err
}

// ResolveNetworkDeviceControllerNumericId resolves a controller ID or
// identifier to its numeric ID.
func ResolveNetworkDeviceControllerNumericId(ctx context.Context, controllerIdOrIdentifier string) (int64, error) {
	controllerId, _, err := resolveControllerIdAndRevision(ctx, controllerIdOrIdentifier)
	return controllerId, err
}

// resolveControllerIdAndRevision returns the numeric controller ID and its
// revision as an If-Match entity tag. The SDK models the revision as a number,
// so it is rendered back to a string here.
func resolveControllerIdAndRevision(ctx context.Context, controllerIdOrIdentifier string) (int64, string, error) {
	controller, err := GetNetworkDeviceControllerByIdOrLabel(ctx, controllerIdOrIdentifier)
	if err != nil {
		return 0, "", err
	}

	controllerId, err := utils.GetInt64FromString(controller.Id)
	if err != nil {
		return 0, "", fmt.Errorf("invalid network device controller ID %q: %w", controller.Id, err)
	}

	return controllerId, strconv.FormatInt(controller.Revision, 10), nil
}
