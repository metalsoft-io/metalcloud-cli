package user

import (
	"context"
	"fmt"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
)

var userPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"DisplayName": {
			Title: "Name",
			Order: 2,
		},
		"Email": {
			Title: "E-mail",
			Order: 3,
		},
		"CreatedTimestamp": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       4,
		},
		"LastLoginTimestamp": {
			Title:       "Last Login",
			Transformer: formatter.FormatDateTimeValue,
			Order:       5,
		},
		"AccessLevel": {
			Title: "Access",
			Order: 6,
		},
	},
}

func List(ctx context.Context) error {
	logger.Get().Info().Msgf("Listing all users")

	client := api.GetApiClient(ctx)

	userList, httpRes, err := client.UserAPI.GetUsers(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(userList, &userPrintConfig)
}

func Get(ctx context.Context, userId string) error {
	logger.Get().Info().Msgf("Get user '%s'", userId)

	userIdNumber, err := strconv.ParseFloat(userId, 32)
	if err != nil {
		err := fmt.Errorf("invalid user ID: '%s'", userId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	client := api.GetApiClient(ctx)

	userInfo, httpRes, err := client.UserAPI.GetUser(ctx, float32(userIdNumber)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(userInfo, &userPrintConfig)
}
