package server

import (
	"context"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

// ServerConnectInterface records the network device port a server interface is
// cabled to.
func ServerConnectInterface(ctx context.Context, serverId string, connectConfig sdk.ServerConnectInterface) error {
	logger.Get().Info().Msgf("Connecting interface of server '%s'", serverId)

	serverIdNumeric, err := GetServerId(serverId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	// The body is always sent: skipping the optional setter would serialize a
	// literal null, which the API rejects with 400.
	httpRes, err := client.ServerAPI.
		ConnectServerInterface(ctx, serverIdNumeric).
		ServerConnectInterface(connectConfig).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Interface %d of server '%s' connected to port '%s' of '%s'",
		connectConfig.ServerInterfaceId, serverId, connectConfig.NetworkDevicePortId, connectConfig.NetworkDeviceHostname)

	return nil
}

func ServerConnectInterfaceConfigExample(ctx context.Context) error {
	connectConfig := sdk.ServerConnectInterface{
		ServerInterfaceId:     1,
		NetworkDevicePortId:   "Ethernet0",
		NetworkDeviceHostname: "leaf-01",
	}

	return formatter.PrintResult(connectConfig, nil)
}

// ServerSetInterfacesDefaultFabric assigns (or, when the fabric id is null,
// clears) the default fabric of the given server interfaces.
func ServerSetInterfacesDefaultFabric(ctx context.Context, serverId string, fabricConfig sdk.ServerInterfacesDefaultFabric) error {
	logger.Get().Info().Msgf("Setting the default fabric of the interfaces of server '%s'", serverId)

	serverIdNumeric, err := GetServerId(serverId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	// The body is always sent: skipping the optional setter would serialize a
	// literal null, which the API rejects with 400.
	httpRes, err := client.ServerAPI.
		SetServerInterfacesDefaultFabric(ctx, serverIdNumeric).
		ServerInterfacesDefaultFabric(fabricConfig).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Default fabric set for the interfaces of server '%s'", serverId)

	return nil
}

func ServerSetInterfacesDefaultFabricConfigExample(ctx context.Context) error {
	fabricConfig := sdk.ServerInterfacesDefaultFabric{
		ServerInterfaceIds: []int64{1, 2},
	}
	fabricConfig.DefaultFabricId.Set(sdk.PtrInt64(10))

	return formatter.PrintResult(fabricConfig, nil)
}

// ServerSetInterfacesRedundancyGroup groups the given server interfaces into a
// redundancy group, or removes them from one when the index is null.
func ServerSetInterfacesRedundancyGroup(ctx context.Context, serverId string, redundancyConfig sdk.ServerInterfacesRedundancyGroup) error {
	logger.Get().Info().Msgf("Setting the redundancy group of the interfaces of server '%s'", serverId)

	serverIdNumeric, err := GetServerId(serverId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	// The body is always sent: skipping the optional setter would serialize a
	// literal null, which the API rejects with 400.
	httpRes, err := client.ServerAPI.
		SetServerInterfacesRedundancyGroup(ctx, serverIdNumeric).
		ServerInterfacesRedundancyGroup(redundancyConfig).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Redundancy group set for the interfaces of server '%s'", serverId)

	return nil
}

func ServerSetInterfacesRedundancyGroupConfigExample(ctx context.Context) error {
	redundancyConfig := sdk.ServerInterfacesRedundancyGroup{
		ServerInterfaceIds: []int64{1, 2},
	}
	redundancyConfig.RedundancyGroupIndex.Set(sdk.PtrFloat32(1))

	return formatter.PrintResult(redundancyConfig, nil)
}
