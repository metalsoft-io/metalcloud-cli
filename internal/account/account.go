package account

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var accountPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
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
			Title: "ID",
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

	request := client.AccountAPI.GetAccounts(ctx).SortBy([]string{"id:ASC"})
	if archived {
		// The API excludes archived accounts by default.
		request = request.FilterArchived([]string{"$eq:1"})
	}

	accounts, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(accounts, meta, len(accounts), &accountPrintConfig)
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
	err := utils.UnmarshalContent(config, &accountConfig)
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
	err = utils.UnmarshalContent(config, &accountConfig)
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

func getAccountId(accountId string) (int64, error) {
	accountIdNumber, err := strconv.ParseInt(accountId, 10, 64)
	if err != nil {
		err := fmt.Errorf("invalid account ID: '%s'", accountId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return accountIdNumber, nil
}

func getAccountIdAndRevision(ctx context.Context, accountId string) (int64, string, error) {
	accountIdNumber, err := getAccountId(accountId)
	if err != nil {
		return 0, "", err
	}

	client := api.GetApiClient(ctx)

	account, httpRes, err := client.AccountAPI.
		GetAccount(ctx, accountIdNumber).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return 0, "", err
	}

	return accountIdNumber, strconv.Itoa(int(account.Revision)), nil
}

// accountConfigPrintConfig formats the account configuration object returned by
// the /accounts/{accountId}/config endpoint.
var accountConfigPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Name": {
			Title:    "Name",
			MaxWidth: 30,
			Order:    1,
		},
		"Code": {
			Title:    "Code",
			MaxWidth: 30,
			Order:    2,
		},
		"FiscalNumber": {
			Title:    "Fiscal #",
			MaxWidth: 30,
			Order:    3,
		},
		"ParentAccountId": {
			Title: "Parent Account Id",
			Order: 4,
		},
		"PrimaryContactId": {
			Title: "Primary Contact Id",
			Order: 5,
		},
		"SecondaryContactId": {
			Title: "Secondary Contact Id",
			Order: 6,
		},
		"QuotaProfileId": {
			Title: "Quota Profile Id",
			Order: 7,
		},
		"ApiKeyValidityDuration": {
			Title: "API Key Validity (s)",
			Order: 8,
		},
		"IsArchived": {
			Title: "Archived",
			Order: 9,
		},
		"Revision": {
			Title: "Revision",
			Order: 10,
		},
	},
}

// accountQuotaBreakdownPrintConfig renders the flattened quota breakdown.
var accountQuotaBreakdownPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Limit": {
			Title:    "Limit",
			MaxWidth: 50,
			Order:    1,
		},
		"Effective": {
			Title:    "Effective",
			MaxWidth: 40,
			Order:    2,
		},
		"Account": {
			Title:    "Account",
			MaxWidth: 40,
			Order:    3,
		},
		"ParentAccount": {
			Title:    "Parent Account",
			MaxWidth: 40,
			Order:    4,
		},
	},
}

// accountQuotaBreakdownUsagePrintConfig adds the usage columns, which the API
// only fills in when the breakdown is requested with usage included.
var accountQuotaBreakdownUsagePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Limit": {
			Title:    "Limit",
			MaxWidth: 50,
			Order:    1,
		},
		"Effective": {
			Title:    "Effective",
			MaxWidth: 30,
			Order:    2,
		},
		"Account": {
			Title:    "Account",
			MaxWidth: 30,
			Order:    3,
		},
		"ParentAccount": {
			Title:    "Parent Account",
			MaxWidth: 30,
			Order:    4,
		},
		"AccountUsage": {
			Title:    "Account Usage",
			MaxWidth: 30,
			Order:    5,
		},
		"ParentAccountUsage": {
			Title:    "Parent Account Usage",
			MaxWidth: 30,
			Order:    6,
		},
	},
}

// accountQuotaBreakdownRow is one flattened limit of the quota breakdown. The
// nested scope objects render poorly as table cells, so the text formats get
// one row per limit while JSON and YAML keep the original object.
type accountQuotaBreakdownRow struct {
	Limit              string `json:"limit"`
	Effective          string `json:"effective"`
	Account            string `json:"account"`
	ParentAccount      string `json:"parentAccount"`
	AccountUsage       string `json:"accountUsage"`
	ParentAccountUsage string `json:"parentAccountUsage"`
}

// AccountUnarchive restores a previously archived account.
func AccountUnarchive(ctx context.Context, accountId string) error {
	logger.Get().Info().Msgf("Un-archiving account '%s'", accountId)

	accountIdNumber, revision, err := getAccountIdAndRevision(ctx, accountId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	accountInfo, httpRes, err := client.AccountAPI.
		UnarchiveAccount(ctx, accountIdNumber).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Account '%s' un-archived", accountId)

	if accountInfo == nil {
		return nil
	}

	return formatter.PrintResult(accountInfo, &accountPrintConfig)
}

// AccountGetConfig retrieves the configuration object of an account.
func AccountGetConfig(ctx context.Context, accountId string) error {
	logger.Get().Info().Msgf("Get configuration for account '%s'", accountId)

	accountIdNumber, err := getAccountId(accountId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	accountConfig, httpRes, err := client.AccountAPI.
		AccountControllerGetUserConfiguration(ctx, accountIdNumber).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(accountConfig, &accountConfigPrintConfig)
}

// AccountQuotaBreakdown shows how the effective quota limits of an account are
// derived from the account and parent account quota profiles.
func AccountQuotaBreakdown(ctx context.Context, accountId string, includeUsage bool) error {
	logger.Get().Info().Msgf("Get quota limits breakdown for account '%s'", accountId)

	accountIdNumber, err := getAccountId(accountId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.AccountAPI.GetAccountQuotaLimitsBreakdown(ctx, accountIdNumber)
	if includeUsage {
		request = request.IncludeUsage(true)
	}

	breakdown, httpRes, err := request.Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	if breakdown == nil {
		return fmt.Errorf("no quota limits breakdown returned for account '%s'", accountId)
	}

	// JSON and YAML keep the nested scopes exactly as the API returned them.
	if formatter.IsNativeFormat() {
		return formatter.PrintResult(breakdown, nil)
	}

	printConfig := &accountQuotaBreakdownPrintConfig
	if includeUsage {
		printConfig = &accountQuotaBreakdownUsagePrintConfig
	}

	return formatter.PrintResult(flattenQuotaBreakdown(breakdown), printConfig)
}

// flattenQuotaBreakdown turns the nested breakdown scopes into one row per
// limit so that the tabular formatters can render it.
func flattenQuotaBreakdown(breakdown *sdk.AccountQuotaLimitsBreakdown) []accountQuotaBreakdownRow {
	effective := quotaScopeValues(breakdown.Effective)
	account := quotaScopeValues(breakdown.Account.Get())
	parentAccount := quotaScopeValues(breakdown.ParentAccount.Get())
	accountUsage := quotaScopeValues(breakdown.AccountUsage.Get())
	parentAccountUsage := quotaScopeValues(breakdown.ParentAccountUsage.Get())

	rows := make([]accountQuotaBreakdownRow, 0, len(effective))
	for _, limit := range quotaLimitNames(effective) {
		rows = append(rows, accountQuotaBreakdownRow{
			Limit:              limit,
			Effective:          formatQuotaValue(effective, limit),
			Account:            formatQuotaValue(account, limit),
			ParentAccount:      formatQuotaValue(parentAccount, limit),
			AccountUsage:       formatQuotaValue(accountUsage, limit),
			ParentAccountUsage: formatQuotaValue(parentAccountUsage, limit),
		})
	}

	return rows
}

// quotaScopeValues converts one scope of the breakdown into a generic map.
// Routing through JSON keeps the API field names and tolerates schema drift.
func quotaScopeValues(scope any) map[string]any {
	if scope == nil {
		return nil
	}

	encoded, err := json.Marshal(scope)
	if err != nil {
		return nil
	}

	var values map[string]any
	if err := json.Unmarshal(encoded, &values); err != nil {
		return nil
	}

	return values
}

// quotaLimitNames returns the limit names in the order the SDK model declares
// them, followed by any extra keys the API sent that the SDK does not know.
func quotaLimitNames(effective map[string]any) []string {
	limitsType := reflect.TypeOf(sdk.QuotaProfileLimits{})

	names := make([]string, 0, len(effective))
	known := make(map[string]bool, limitsType.NumField())

	for i := 0; i < limitsType.NumField(); i++ {
		name := strings.Split(limitsType.Field(i).Tag.Get("json"), ",")[0]
		if name == "" || name == "-" {
			continue
		}

		known[name] = true
		names = append(names, name)
	}

	extra := make([]string, 0)
	for name := range effective {
		if !known[name] {
			extra = append(extra, name)
		}
	}
	sort.Strings(extra)

	return append(names, extra...)
}

// formatQuotaValue renders one quota value as a table cell. Values missing from
// a scope (the scope itself may be null) are shown as "-".
func formatQuotaValue(values map[string]any, name string) string {
	if values == nil {
		return "-"
	}

	value, ok := values[name]
	if !ok || value == nil {
		return "-"
	}

	switch typedValue := value.(type) {
	case float64:
		return strconv.FormatFloat(typedValue, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(typedValue)
	case string:
		return typedValue
	case []any:
		items := make([]string, 0, len(typedValue))
		for _, item := range typedValue {
			items = append(items, fmt.Sprintf("%v", item))
		}
		return strings.Join(items, ", ")
	case map[string]any:
		items := make([]string, 0, len(typedValue))
		for key, item := range typedValue {
			items = append(items, fmt.Sprintf("%s=%v", key, item))
		}
		sort.Strings(items)
		return strings.Join(items, ", ")
	default:
		return fmt.Sprintf("%v", typedValue)
	}
}
