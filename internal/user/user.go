package user

import (
	"context"
	"fmt"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
)

func List(ctx context.Context) error {
	logger.Get().Info().Msgf("Listing all users")

	client, err := system.GetApiClient(ctx)
	if err != nil {
		return err
	}

	userList, httpRes, err := client.UserAPI.GetUsers(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(userList, nil)
}

func Get(ctx context.Context, userId string) error {
	logger.Get().Info().Msgf("Get user '%s'", userId)

	client, err := system.GetApiClient(ctx)
	if err != nil {
		return err
	}

	userIdNumber, err := strconv.ParseFloat(userId, 32)
	if err != nil {
		err := fmt.Errorf("invalid user ID: '%s'", userId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	userInfo, httpRes, err := client.UserAPI.GetUser(ctx, float32(userIdNumber)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(userInfo, nil)
}
