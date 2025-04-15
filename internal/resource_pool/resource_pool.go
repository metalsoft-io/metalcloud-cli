package resource_pool

import (
	"context"
	"fmt"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"

	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var resourcePoolPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"ResourcePoolId": {
			Title: "#",
			Order: 1,
		},
		"ResourcePoolLabel": {
			Title:    "Label",
			MaxWidth: 30,
			Order:    2,
		},
		"ResourcePoolDescription": {
			Title:    "Description",
			MaxWidth: 50,
			Order:    3,
		},
	},
}

// ResourcePoolList lists all resource pools
func ResourcePoolList(ctx context.Context, page int, limit int, search string) error {
	logger.Get().Info().Msgf("Listing resource pools")

	client := api.GetApiClient(ctx)

	req := client.ResourcePoolAPI.GetResourcePools(ctx)

	if page > 0 {
		req = req.Page(float32(page))
	}

	if limit > 0 {
		req = req.Limit(float32(limit))
	}

	if search != "" {
		req = req.Search(search)
	}

	poolList, httpRes, err := req.Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(poolList.Data, &resourcePoolPrintConfig)
}

// ResourcePoolGet retrieves a specific resource pool's information
func ResourcePoolGet(ctx context.Context, poolId string) error {
	logger.Get().Info().Msgf("Getting resource pool '%s'", poolId)

	poolIdNumber, err := strconv.ParseFloat(poolId, 32)
	if err != nil {
		err := fmt.Errorf("invalid resource pool ID: '%s'", poolId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	client := api.GetApiClient(ctx)

	pool, httpRes, err := client.ResourcePoolAPI.GetResourcePool(ctx, float32(poolIdNumber)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(pool, &resourcePoolPrintConfig)
}

// ResourcePoolCreate creates a new resource pool
func ResourcePoolCreate(ctx context.Context, label string, description string) error {
	logger.Get().Info().Msgf("Creating resource pool '%s'", label)

	client := api.GetApiClient(ctx)

	createResourcePool := sdk.NewCreateResourcePool(label, description)

	pool, httpRes, err := client.ResourcePoolAPI.CreateResourcePool(ctx).
		CreateResourcePool(*createResourcePool).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(pool, &resourcePoolPrintConfig)
}

// ResourcePoolDelete deletes a resource pool
func ResourcePoolDelete(ctx context.Context, poolId string) error {
	logger.Get().Info().Msgf("Deleting resource pool '%s'", poolId)

	poolIdNumber, err := strconv.ParseFloat(poolId, 32)
	if err != nil {
		err := fmt.Errorf("invalid resource pool ID: '%s'", poolId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.ResourcePoolAPI.DeleteResourcePool(ctx, float32(poolIdNumber)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Resource pool with ID %s deleted successfully", poolId)
	return nil
}

// ResourcePoolGetUsers retrieves users that have access to a resource pool
func ResourcePoolGetUsers(ctx context.Context, poolId string) error {
	logger.Get().Info().Msgf("Getting users for resource pool '%s'", poolId)

	poolIdNumber, err := strconv.ParseFloat(poolId, 32)
	if err != nil {
		err := fmt.Errorf("invalid resource pool ID: '%s'", poolId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	client := api.GetApiClient(ctx)

	users, httpRes, err := client.ResourcePoolAPI.GetResourcePoolUsers(ctx, float32(poolIdNumber)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(users, nil)
}

// ResourcePoolAddUser adds a user to a resource pool
func ResourcePoolAddUser(ctx context.Context, poolId string, userId string) error {
	logger.Get().Info().Msgf("Adding user '%s' to resource pool '%s'", userId, poolId)

	poolIdNumber, err := strconv.ParseFloat(poolId, 32)
	if err != nil {
		err := fmt.Errorf("invalid resource pool ID: '%s'", poolId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	userIdNumber, err := strconv.ParseFloat(userId, 32)
	if err != nil {
		err := fmt.Errorf("invalid user ID: '%s'", userId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.ResourcePoolAPI.AddResourcePoolUser(ctx, float32(poolIdNumber), float32(userIdNumber)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("User '%s' successfully added to resource pool '%s'", userId, poolId)
	return nil
}

// ResourcePoolRemoveUser removes a user from a resource pool
func ResourcePoolRemoveUser(ctx context.Context, poolId string, userId string) error {
	logger.Get().Info().Msgf("Removing user '%s' from resource pool '%s'", userId, poolId)

	poolIdNumber, err := strconv.ParseFloat(poolId, 32)
	if err != nil {
		err := fmt.Errorf("invalid resource pool ID: '%s'", poolId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	userIdNumber, err := strconv.ParseFloat(userId, 32)
	if err != nil {
		err := fmt.Errorf("invalid user ID: '%s'", userId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.ResourcePoolAPI.RemoveResourcePoolUser(ctx, float32(poolIdNumber), float32(userIdNumber)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("User '%s' successfully removed from resource pool '%s'", userId, poolId)
	return nil
}

// ResourcePoolGetServers retrieves servers that are part of a resource pool
func ResourcePoolGetServers(ctx context.Context, poolId string) error {
	logger.Get().Info().Msgf("Getting servers for resource pool '%s'", poolId)

	poolIdNumber, err := strconv.ParseFloat(poolId, 32)
	if err != nil {
		err := fmt.Errorf("invalid resource pool ID: '%s'", poolId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	client := api.GetApiClient(ctx)

	servers, httpRes, err := client.ResourcePoolAPI.GetResourcePoolServers(ctx, float32(poolIdNumber)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(servers, nil)
}

// ResourcePoolAddServer adds a server to a resource pool
func ResourcePoolAddServer(ctx context.Context, poolId string, serverId string) error {
	logger.Get().Info().Msgf("Adding server '%s' to resource pool '%s'", serverId, poolId)

	poolIdNumber, err := strconv.ParseFloat(poolId, 32)
	if err != nil {
		err := fmt.Errorf("invalid resource pool ID: '%s'", poolId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	serverIdNumber, err := strconv.ParseFloat(serverId, 32)
	if err != nil {
		err := fmt.Errorf("invalid server ID: '%s'", serverId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.ResourcePoolAPI.AddServerToResourcePool(ctx, float32(poolIdNumber), float32(serverIdNumber)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Server '%s' successfully added to resource pool '%s'", serverId, poolId)
	return nil
}

// ResourcePoolRemoveServer removes a server from a resource pool
func ResourcePoolRemoveServer(ctx context.Context, poolId string, serverId string) error {
	logger.Get().Info().Msgf("Removing server '%s' from resource pool '%s'", serverId, poolId)

	poolIdNumber, err := strconv.ParseFloat(poolId, 32)
	if err != nil {
		err := fmt.Errorf("invalid resource pool ID: '%s'", poolId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	serverIdNumber, err := strconv.ParseFloat(serverId, 32)
	if err != nil {
		err := fmt.Errorf("invalid server ID: '%s'", serverId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.ResourcePoolAPI.RemoveServerFromResourcePool(ctx, float32(poolIdNumber), float32(serverIdNumber)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Server '%s' successfully removed from resource pool '%s'", serverId, poolId)
	return nil
}

// ResourcePoolGetSubnetPools retrieves subnet pools that are part of a resource pool
func ResourcePoolGetSubnetPools(ctx context.Context, poolId string) error {
	logger.Get().Info().Msgf("Getting subnet pools for resource pool '%s'", poolId)

	poolIdNumber, err := strconv.ParseFloat(poolId, 32)
	if err != nil {
		err := fmt.Errorf("invalid resource pool ID: '%s'", poolId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	client := api.GetApiClient(ctx)

	subnetPools, httpRes, err := client.ResourcePoolAPI.GetResourcePoolSubnetPools(ctx, float32(poolIdNumber)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(subnetPools, nil)
}

// ResourcePoolAddSubnetPool adds a subnet pool to a resource pool
func ResourcePoolAddSubnetPool(ctx context.Context, poolId string, subnetPoolId string) error {
	logger.Get().Info().Msgf("Adding subnet pool '%s' to resource pool '%s'", subnetPoolId, poolId)

	poolIdNumber, err := strconv.ParseFloat(poolId, 32)
	if err != nil {
		err := fmt.Errorf("invalid resource pool ID: '%s'", poolId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	subnetPoolIdNumber, err := strconv.ParseFloat(subnetPoolId, 32)
	if err != nil {
		err := fmt.Errorf("invalid subnet pool ID: '%s'", subnetPoolId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.ResourcePoolAPI.AddSubnetPoolToResourcePool(ctx, float32(poolIdNumber), float32(subnetPoolIdNumber)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Subnet pool '%s' successfully added to resource pool '%s'", subnetPoolId, poolId)
	return nil
}

// ResourcePoolRemoveSubnetPool removes a subnet pool from a resource pool
func ResourcePoolRemoveSubnetPool(ctx context.Context, poolId string, subnetPoolId string) error {
	logger.Get().Info().Msgf("Removing subnet pool '%s' from resource pool '%s'", subnetPoolId, poolId)

	poolIdNumber, err := strconv.ParseFloat(poolId, 32)
	if err != nil {
		err := fmt.Errorf("invalid resource pool ID: '%s'", poolId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	subnetPoolIdNumber, err := strconv.ParseFloat(subnetPoolId, 32)
	if err != nil {
		err := fmt.Errorf("invalid subnet pool ID: '%s'", subnetPoolId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.ResourcePoolAPI.RemoveSubnetPoolFromResourcePool(ctx, float32(poolIdNumber), float32(subnetPoolIdNumber)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Subnet pool '%s' successfully removed from resource pool '%s'", subnetPoolId, poolId)
	return nil
}
