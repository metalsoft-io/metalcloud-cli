package account

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var accountPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"Name": {
			MaxWidth: 30,
			Order:    2,
		},
		"Code": {
			MaxWidth: 30,
			Order:    3,
		},
		"FiscalNumber": {
			Title:    "Fiscal #",
			MaxWidth: 30,
			Order:    4,
		},
		"IsArchived": {
			Title: "Archived",
			Order: 5,
		},
		"PrimaryContactId": {
			Title: "Primary Contact Id",
			Order: 6,
		},
		"SecondaryContactId": {
			Title: "Secondary Contact Id",
			Order: 7,
		},
	},
}

var accountUsersPrintConfig = formatter.PrintConfig{
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
		"AccessLevel": {
			Title: "Access",
			Order: 4,
		},
	},
}

func AccountList(ctx context.Context) error {
	logger.Get().Info().Msgf("Listing all accounts")

	client := api.GetApiClient(ctx)

	accountList, httpRes, err := client.AccountAPI.GetAccounts(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(accountList, &accountPrintConfig)
}

func AccountGet(ctx context.Context, accountId string) error {
	logger.Get().Info().Msgf("Get account '%s'", accountId)

	accountIdNumber, err := getAccountId(accountId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	accountInfo, httpRes, err := client.AccountAPI.GetAccount(ctx, accountIdNumber).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(accountInfo, &accountPrintConfig)
}

func AccountCreate(ctx context.Context, config []byte) error {
	logger.Get().Info().Msgf("Creating account")

	var accountConfig sdk.CreateAccount
	err := json.Unmarshal(config, &accountConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	accountInfo, httpRes, err := client.AccountAPI.CreateAccount(ctx).CreateAccount(accountConfig).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(accountInfo, &accountPrintConfig)
}

func AccountUpdate(ctx context.Context, accountId string, config []byte) error {
	logger.Get().Info().Msgf("Updating account '%s'", accountId)

	accountIdNumber, revision, err := getAccountIdAndRevision(ctx, accountId)
	if err != nil {
		return err
	}

	var accountConfig sdk.UpdateAccount
	err = json.Unmarshal(config, &accountConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	accountInfo, httpRes, err := client.AccountAPI.
		UpdateAccountConfig(ctx, accountIdNumber).
		UpdateAccount(accountConfig).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(accountInfo, &accountPrintConfig)
}

func AccountArchive(ctx context.Context, accountId string) error {
	logger.Get().Info().Msgf("Archiving account '%s'", accountId)

	accountIdNumber, revision, err := getAccountIdAndRevision(ctx, accountId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	_, httpRes, err := client.AccountAPI.
		ArchiveAccount(ctx, accountIdNumber).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Account '%s' archived", accountId)
	return nil
}

func AccountGetUsers(ctx context.Context, accountId string) error {
	logger.Get().Info().Msgf("Getting users for account '%s'", accountId)

	accountIdNumber, err := getAccountId(accountId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	users, httpRes, err := client.AccountAPI.
		GetAccountUsers(ctx, accountIdNumber).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(users, &accountUsersPrintConfig)
}

func getAccountId(accountId string) (float32, error) {
	accountIdNumber, err := strconv.ParseFloat(accountId, 32)
	if err != nil {
		err := fmt.Errorf("invalid account ID: '%s'", accountId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return float32(accountIdNumber), nil
}

func getAccountIdAndRevision(ctx context.Context, accountId string) (float32, string, error) {
	accountIdNumber, err := getAccountId(accountId)
	if err != nil {
		return 0, "", err
	}

	client := api.GetApiClient(ctx)

	account, httpRes, err := client.AccountAPI.
		GetAccount(ctx, float32(accountIdNumber)).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return 0, "", err
	}

	return float32(accountIdNumber), strconv.Itoa(int(account.Revision)), nil
}
