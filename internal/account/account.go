package account

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

type accountRaw struct {
	Id                 interface{} `json:"id"`
	Revision           json.Number `json:"revision"`
	Name               *string     `json:"name"`
	Code               *string     `json:"code"`
	FiscalNumber       *string     `json:"fiscalNumber"`
	IsArchived         interface{} `json:"isArchived"`
	PrimaryContactId   interface{} `json:"primaryContactId"`
	SecondaryContactId interface{} `json:"secondaryContactId"`
}

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

func AccountList(ctx context.Context, archived bool) error {
	logger.Get().Info().Msgf("Listing all accounts")

	client := api.GetApiClient(ctx)

	rawItems, meta, err := utils.FetchAllPagesRaw(func(page float32) (*http.Response, error) {
		request := client.AccountAPI.GetAccounts(ctx).SortBy([]string{"id:ASC"})
		if archived {
			// The API excludes archived accounts by default.
			request = request.FilterArchived([]string{"$eq:1"})
		}
		_, httpRes, _ := request.Page(page).Limit(100).Execute()
		return httpRes, nil
	})
	if err != nil {
		return err
	}

	records, err := utils.UnmarshalRawItems[accountRaw](rawItems)
	if err != nil {
		return fmt.Errorf("failed to parse accounts: %w", err)
	}

	return utils.PrintAllRaw(rawItems, records, meta, len(records), &accountPrintConfig)
}

func AccountGet(ctx context.Context, accountId string) error {
	logger.Get().Info().Msgf("Get account '%s'", accountId)

	accountIdNumber, err := getAccountId(accountId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	_, httpRes, sdkErr := client.AccountAPI.GetAccount(ctx, accountIdNumber).Execute()

	record, err := parseAccountResponse(httpRes, sdkErr)
	if err != nil {
		return err
	}

	return formatter.PrintResult(record, &accountPrintConfig)
}

func AccountCreate(ctx context.Context, config []byte) error {
	logger.Get().Info().Msgf("Creating account")

	var accountConfig sdk.CreateAccount
	err := utils.UnmarshalContent(config, &accountConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	_, httpRes, sdkErr := client.AccountAPI.CreateAccount(ctx).CreateAccount(accountConfig).Execute()

	record, err := parseAccountResponse(httpRes, sdkErr)
	if err != nil {
		return err
	}

	return formatter.PrintResult(record, &accountPrintConfig)
}

func AccountUpdate(ctx context.Context, accountId string, config []byte) error {
	logger.Get().Info().Msgf("Updating account '%s'", accountId)

	accountIdNumber, revision, err := getAccountIdAndRevision(ctx, accountId)
	if err != nil {
		return err
	}

	var accountConfig sdk.UpdateAccount
	err = utils.UnmarshalContent(config, &accountConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	_, httpRes, sdkErr := client.AccountAPI.
		UpdateAccountConfig(ctx, accountIdNumber).
		UpdateAccount(accountConfig).
		IfMatch(revision).
		Execute()

	record, err := parseAccountResponse(httpRes, sdkErr)
	if err != nil {
		return err
	}

	return formatter.PrintResult(record, &accountPrintConfig)
}

func AccountArchive(ctx context.Context, accountId string) error {
	logger.Get().Info().Msgf("Archiving account '%s'", accountId)

	accountIdNumber, revision, err := getAccountIdAndRevision(ctx, accountId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	_, httpRes, sdkErr := client.AccountAPI.
		ArchiveAccount(ctx, accountIdNumber).
		IfMatch(revision).
		Execute()
	if _, err := parseAccountResponse(httpRes, sdkErr); err != nil {
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

	_, httpRes, sdkErr := client.AccountAPI.
		GetAccount(ctx, float32(accountIdNumber)).
		Execute()

	record, err := parseAccountResponse(httpRes, sdkErr)
	if err != nil {
		return 0, "", err
	}

	return float32(accountIdNumber), record.Revision.String(), nil
}

// parseAccountResponse reads an account response body without SDK type unmarshalling.
// The SDK Account model marks `limits` as required, but the API omits it, so SDK
// unmarshalling fails on otherwise valid responses.
func parseAccountResponse(httpRes *http.Response, sdkErr error) (*accountRaw, error) {
	if httpRes != nil && httpRes.StatusCode >= 400 {
		return nil, response_inspector.InspectResponse(httpRes, sdkErr)
	}
	if httpRes == nil {
		return nil, sdkErr
	}

	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var record accountRaw
	if err := json.Unmarshal(body, &record); err != nil {
		return nil, fmt.Errorf("failed to parse account: %w", err)
	}
	return &record, nil
}
