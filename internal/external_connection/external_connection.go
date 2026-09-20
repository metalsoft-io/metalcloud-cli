package external_connection

import (
	"context"
	"fmt"

	"github.com/metalsoft-io/metalcloud-cli/internal/network_device"
	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var externalConnectionPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"Label": {
			Order:    2,
			MaxWidth: 30,
		},
		"Name": {
			Order:    3,
			MaxWidth: 30,
		},
		"FabricId": {
			Title: "Fabric",
			Order: 4,
		},
		"Revision": {
			Order: 5,
		},
		"CreatedAt": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       6,
		},
	},
}

var externalConnectionInterfacePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"NetworkDeviceId": {
			Title: "Network Device",
			Order: 2,
		},
		"NetworkDeviceInterfaceId": {
			Title: "Interface ID",
			Order: 3,
		},
		"NetworkDeviceInterfaceName": {
			Title:    "Interface",
			Order:    4,
			MaxWidth: 30,
		},
		"Revision": {
			Order: 5,
		},
		"CreatedAt": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       6,
		},
	},
}

var externalConnectionLogicalNetworkPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"ExternalConnectionId": {
			Title: "External Connection",
			Order: 2,
		},
		"LogicalNetworkId": {
			Title: "Logical Network",
			Order: 3,
		},
		"Status": {
			Transformer: formatter.FormatStatusValue,
			Order:       4,
		},
		"CreatedAt": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       5,
		},
	},
}

var networkDeviceInterfacePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"NetworkDeviceId": {
			Title: "Network Device",
			Order: 1,
		},
		"NetworkDeviceInterfaceId": {
			Title: "Interface ID",
			Order: 2,
		},
		"NetworkDeviceInterfaceName": {
			Title:    "Interface",
			Order:    3,
			MaxWidth: 30,
		},
		"ExternalConnectionId": {
			Title: "External Connection",
			Order: 4,
		},
		"ExternalConnectionInterfaceId": {
			Title: "External Connection Interface",
			Order: 5,
		},
	},
}

// ExternalConnectionListFilters carries the optional list filters so the cmd
// layer does not have to pass a long positional argument list.
type ExternalConnectionListFilters struct {
	Id       []string
	FabricId []string
	Label    []string
	Name     []string
}

func ExternalConnectionList(ctx context.Context, filters ExternalConnectionListFilters) error {
	logger.Get().Info().Msgf("Listing external connections")

	client := api.GetApiClient(ctx)

	request := client.ExternalConnectionAPI.GetExternalConnections(ctx).SortBy([]string{"id:ASC"})
	if len(filters.Id) > 0 {
		request = request.FilterId(utils.ProcessFilterStringSlice(filters.Id))
	}
	if len(filters.FabricId) > 0 {
		request = request.FilterFabricId(utils.ProcessFilterStringSlice(filters.FabricId))
	}
	if len(filters.Label) > 0 {
		request = request.FilterLabel(utils.ProcessFilterStringSlice(filters.Label))
	}
	if len(filters.Name) > 0 {
		request = request.FilterName(utils.ProcessFilterStringSlice(filters.Name))
	}

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &externalConnectionPrintConfig)
}

func ExternalConnectionGet(ctx context.Context, externalConnectionIdOrLabel string) error {
	logger.Get().Info().Msgf("Get external connection '%s'", externalConnectionIdOrLabel)

	externalConnection, err := GetExternalConnectionByIdOrLabel(ctx, externalConnectionIdOrLabel)
	if err != nil {
		return err
	}

	return formatter.PrintResult(externalConnection, &externalConnectionPrintConfig)
}

func ExternalConnectionConfigExample(ctx context.Context) error {
	example := sdk.CreateExternalConnection{
		Label:    "dc1-external-connection",
		Name:     "DC1 external connection",
		FabricId: 1,
		ExternalConnectionInterfaces: []sdk.CreateExternalConnectionInterface{
			{NetworkDeviceInterfaceId: 101},
		},
	}

	if formatter.IsNativeFormat() {
		return formatter.PrintResult(example, nil)
	}
	return formatter.PrintYamlResult(example)
}

func ExternalConnectionCreate(ctx context.Context, create sdk.CreateExternalConnection) error {
	logger.Get().Info().Msgf("Creating external connection '%s'", create.Label)

	client := api.GetApiClient(ctx)

	externalConnection, httpRes, err := client.ExternalConnectionAPI.
		CreateExternalConnection(ctx).
		CreateExternalConnection(create).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(externalConnection, &externalConnectionPrintConfig)
}

func ExternalConnectionUpdate(ctx context.Context, externalConnectionIdOrLabel string, config []byte) error {
	logger.Get().Info().Msgf("Updating external connection '%s'", externalConnectionIdOrLabel)

	var update sdk.UpdateExternalConnection
	if err := utils.UnmarshalContent(config, &update); err != nil {
		return err
	}

	externalConnectionId, revision, err := resolveExternalConnectionIdAndRevision(ctx, externalConnectionIdOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	externalConnection, httpRes, err := client.ExternalConnectionAPI.
		UpdateExternalConnection(ctx, externalConnectionId).
		UpdateExternalConnection(update).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(externalConnection, &externalConnectionPrintConfig)
}

func ExternalConnectionDelete(ctx context.Context, externalConnectionIdOrLabel string) error {
	logger.Get().Info().Msgf("Deleting external connection '%s'", externalConnectionIdOrLabel)

	externalConnectionId, err := ResolveExternalConnectionNumericId(ctx, externalConnectionIdOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.ExternalConnectionAPI.DeleteExternalConnection(ctx, externalConnectionId).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("External connection '%s' deleted", externalConnectionIdOrLabel)
	return nil
}

func ExternalConnectionInterfaceList(ctx context.Context, externalConnectionIdOrLabel string) error {
	logger.Get().Info().Msgf("Listing interfaces of external connection '%s'", externalConnectionIdOrLabel)

	externalConnectionId, err := ResolveExternalConnectionNumericId(ctx, externalConnectionIdOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.ExternalConnectionAPI.
		GetExternalConnectionInterfaces(ctx, externalConnectionId).
		SortBy([]string{"id:ASC"})

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &externalConnectionInterfacePrintConfig)
}

func ExternalConnectionInterfaceGet(ctx context.Context, externalConnectionIdOrLabel string, interfaceId string) error {
	logger.Get().Info().Msgf("Get interface '%s' of external connection '%s'", interfaceId, externalConnectionIdOrLabel)

	connectionInterface, err := getExternalConnectionInterface(ctx, externalConnectionIdOrLabel, interfaceId)
	if err != nil {
		return err
	}

	return formatter.PrintResult(connectionInterface, &externalConnectionInterfacePrintConfig)
}

func ExternalConnectionInterfaceAdd(ctx context.Context, externalConnectionIdOrLabel string, networkDeviceInterfaceId string) error {
	logger.Get().Info().Msgf("Adding network device interface '%s' to external connection '%s'",
		networkDeviceInterfaceId, externalConnectionIdOrLabel)

	externalConnectionId, err := ResolveExternalConnectionNumericId(ctx, externalConnectionIdOrLabel)
	if err != nil {
		return err
	}

	interfaceIdNumeric, err := parseId(networkDeviceInterfaceId, "network device interface")
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	connectionInterface, httpRes, err := client.ExternalConnectionAPI.
		CreateExternalConnectionInterface(ctx, externalConnectionId).
		CreateExternalConnectionInterface(sdk.CreateExternalConnectionInterface{
			NetworkDeviceInterfaceId: interfaceIdNumeric,
		}).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(connectionInterface, &externalConnectionInterfacePrintConfig)
}

func ExternalConnectionInterfaceUpdate(ctx context.Context, externalConnectionIdOrLabel string, interfaceId string, networkDeviceInterfaceId string) error {
	logger.Get().Info().Msgf("Updating interface '%s' of external connection '%s'", interfaceId, externalConnectionIdOrLabel)

	externalConnectionId, err := ResolveExternalConnectionNumericId(ctx, externalConnectionIdOrLabel)
	if err != nil {
		return err
	}

	connectionInterface, err := getExternalConnectionInterface(ctx, externalConnectionIdOrLabel, interfaceId)
	if err != nil {
		return err
	}

	interfaceIdNumeric, err := parseId(networkDeviceInterfaceId, "network device interface")
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	// The interface carries its own revision; the parent external connection
	// revision is not accepted by this endpoint.
	updated, httpRes, err := client.ExternalConnectionAPI.
		UpdateExternalConnectionInterface(ctx, externalConnectionId, connectionInterface.Id).
		UpdateExternalConnectionInterface(sdk.UpdateExternalConnectionInterface{
			NetworkDeviceInterfaceId: interfaceIdNumeric,
		}).
		IfMatch(connectionInterface.Revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(updated, &externalConnectionInterfacePrintConfig)
}

func ExternalConnectionInterfaceRemove(ctx context.Context, externalConnectionIdOrLabel string, interfaceId string) error {
	logger.Get().Info().Msgf("Removing interface '%s' from external connection '%s'", interfaceId, externalConnectionIdOrLabel)

	externalConnectionId, err := ResolveExternalConnectionNumericId(ctx, externalConnectionIdOrLabel)
	if err != nil {
		return err
	}

	connectionInterface, err := getExternalConnectionInterface(ctx, externalConnectionIdOrLabel, interfaceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.ExternalConnectionAPI.
		DeleteExternalConnectionInterface(ctx, externalConnectionId, connectionInterface.Id).
		IfMatch(connectionInterface.Revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Interface '%s' removed from external connection '%s'", interfaceId, externalConnectionIdOrLabel)
	return nil
}

func ExternalConnectionLogicalNetworkList(ctx context.Context, externalConnectionIdOrLabel string) error {
	logger.Get().Info().Msgf("Listing logical networks of external connection '%s'", externalConnectionIdOrLabel)

	externalConnectionId, err := ResolveExternalConnectionNumericId(ctx, externalConnectionIdOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.ExternalConnectionAPI.
		GetExternalConnectionLogicalNetworks(ctx, externalConnectionId).
		SortBy([]string{"id:ASC"})

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &externalConnectionLogicalNetworkPrintConfig)
}

func ExternalConnectionLogicalNetworkGet(ctx context.Context, externalConnectionIdOrLabel string, id string) error {
	logger.Get().Info().Msgf("Get logical network '%s' of external connection '%s'", id, externalConnectionIdOrLabel)

	externalConnectionId, err := ResolveExternalConnectionNumericId(ctx, externalConnectionIdOrLabel)
	if err != nil {
		return err
	}

	idNumeric, err := parseId(id, "external connection logical network")
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	logicalNetwork, httpRes, err := client.ExternalConnectionAPI.
		GetExternalConnectionLogicalNetworkById(ctx, externalConnectionId, idNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(logicalNetwork, &externalConnectionLogicalNetworkPrintConfig)
}

func ExternalConnectionLogicalNetworkAdd(ctx context.Context, externalConnectionIdOrLabel string, logicalNetworkId string) error {
	logger.Get().Info().Msgf("Attaching logical network '%s' to external connection '%s'", logicalNetworkId, externalConnectionIdOrLabel)

	externalConnectionId, err := ResolveExternalConnectionNumericId(ctx, externalConnectionIdOrLabel)
	if err != nil {
		return err
	}

	logicalNetworkIdNumeric, err := parseId(logicalNetworkId, "logical network")
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	logicalNetwork, httpRes, err := client.ExternalConnectionAPI.
		CreateExternalConnectionLogicalNetwork(ctx, externalConnectionId).
		CreateExternalConnectionLogicalNetwork(sdk.CreateExternalConnectionLogicalNetwork{
			LogicalNetworkId: logicalNetworkIdNumeric,
		}).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(logicalNetwork, &externalConnectionLogicalNetworkPrintConfig)
}

func ExternalConnectionLogicalNetworkRemove(ctx context.Context, externalConnectionIdOrLabel string, id string) error {
	logger.Get().Info().Msgf("Detaching logical network '%s' from external connection '%s'", id, externalConnectionIdOrLabel)

	externalConnectionId, err := ResolveExternalConnectionNumericId(ctx, externalConnectionIdOrLabel)
	if err != nil {
		return err
	}

	idNumeric, err := parseId(id, "external connection logical network")
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.ExternalConnectionAPI.
		DeleteExternalConnectionLogicalNetwork(ctx, externalConnectionId, idNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Logical network '%s' detached from external connection '%s'", id, externalConnectionIdOrLabel)
	return nil
}

// NetworkDeviceInterfacesGet lists the interfaces of a network device together
// with the external connection (if any) each interface belongs to.
func NetworkDeviceInterfacesGet(ctx context.Context, networkDeviceIdOrLabel string) error {
	logger.Get().Info().Msgf("Listing external connection interfaces of network device '%s'", networkDeviceIdOrLabel)

	networkDevice, err := network_device.GetNetworkDeviceByIdOrLabel(ctx, networkDeviceIdOrLabel)
	if err != nil {
		return err
	}

	networkDeviceId, err := utils.GetFloat32FromString(networkDevice.Id)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	interfaces, httpRes, err := client.ExternalConnectionAPI.
		GetNetworkDeviceInterfacesAndExternalConnections(ctx, networkDeviceId).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(interfaces, &networkDeviceInterfacePrintConfig)
}

// GetExternalConnectionByIdOrLabel resolves an external connection by numeric
// ID first and falls back to an exact label match.
func GetExternalConnectionByIdOrLabel(ctx context.Context, externalConnectionIdOrLabel string) (*sdk.ExternalConnection, error) {
	client := api.GetApiClient(ctx)

	if idNumeric, err := utils.GetInt64FromString(externalConnectionIdOrLabel); err == nil {
		externalConnection, httpRes, err := client.ExternalConnectionAPI.GetExternalConnectionById(ctx, idNumeric).Execute()
		if err = response_inspector.InspectResponse(httpRes, err); err == nil {
			return externalConnection, nil
		}
	}

	list, httpRes, err := client.ExternalConnectionAPI.
		GetExternalConnections(ctx).
		FilterLabel([]string{externalConnectionIdOrLabel}).
		Execute()
	if err = response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, err
	}

	for i := range list.Data {
		if list.Data[i].Label == externalConnectionIdOrLabel {
			return &list.Data[i], nil
		}
	}

	err = fmt.Errorf("external connection '%s' not found", externalConnectionIdOrLabel)
	logger.Get().Error().Err(err).Msg("")
	return nil, err
}

// ResolveExternalConnectionNumericId resolves an external connection ID or
// label to its numeric ID.
func ResolveExternalConnectionNumericId(ctx context.Context, externalConnectionIdOrLabel string) (int64, error) {
	externalConnectionId, _, err := resolveExternalConnectionIdAndRevision(ctx, externalConnectionIdOrLabel)
	return externalConnectionId, err
}

func resolveExternalConnectionIdAndRevision(ctx context.Context, externalConnectionIdOrLabel string) (int64, string, error) {
	externalConnection, err := GetExternalConnectionByIdOrLabel(ctx, externalConnectionIdOrLabel)
	if err != nil {
		return 0, "", err
	}

	externalConnectionId, err := utils.GetInt64FromString(externalConnection.Id)
	if err != nil {
		return 0, "", fmt.Errorf("invalid external connection ID %q: %w", externalConnection.Id, err)
	}

	return externalConnectionId, externalConnection.Revision, nil
}

func getExternalConnectionInterface(ctx context.Context, externalConnectionIdOrLabel string, interfaceId string) (*sdk.ExternalConnectionInterface, error) {
	externalConnectionId, err := ResolveExternalConnectionNumericId(ctx, externalConnectionIdOrLabel)
	if err != nil {
		return nil, err
	}

	interfaceIdNumeric, err := parseId(interfaceId, "external connection interface")
	if err != nil {
		return nil, err
	}

	client := api.GetApiClient(ctx)

	connectionInterface, httpRes, err := client.ExternalConnectionAPI.
		GetExternalConnectionInterfaceById(ctx, externalConnectionId, interfaceIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, err
	}

	return connectionInterface, nil
}

func parseId(value string, what string) (int64, error) {
	id, err := utils.GetInt64FromString(value)
	if err != nil {
		err = fmt.Errorf("invalid %s ID: '%s'", what, value)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}
	return id, nil
}
