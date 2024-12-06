package user

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"slices"
	"strconv"
	"strings"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
	metalcloud2 "github.com/metalsoft-io/metal-cloud-sdk2-go"

	"github.com/metalsoft-io/metalcloud-cli/internal/colors"
	"github.com/metalsoft-io/metalcloud-cli/internal/command"
	"github.com/metalsoft-io/metalcloud-cli/internal/configuration"
	"github.com/metalsoft-io/tableformatter"
)

var userFilterProperties = []string{
	"franchise",
	"user_access_level",
	"user_auth_failed_attempts_since_last_login",
	"user_authenticator_created_timestamp",
	"user_authenticator_is_mandatory",
	"user_authenticator_must_change",
	"user_blocked",
	"user_brand",
	"user_created_timestamp",
	"user_custom_prices_json",
	"user_display_name",
	"user_email",
	"user_email_status",
	"user_exclude_from_reports",
	"user_experimental_tags_json",
	"user_external_ids_json",
	"user_gui_settings_json",
	"user_id",
	"user_infrastructure_id_default",
	"user_is_billable",
	"user_is_brand_manager",
	"user_is_datastore_publisher",
	"user_is_suspended",
	"user_is_test_account",
	"user_is_testing_mode",
	"user_kerberos_principal_name",
	"user_language",
	"user_last_login_timestamp",
	"user_last_login_type",
	"user_limits_json",
	"user_password_change_required",
	"user_permissions_json",
	"user_plan_type",
	"user_promotion_tags_json",
}

var UserCmds = []command.Command{
	{
		Description:  "Lists all users.",
		Subject:      "user",
		AltSubject:   "user",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list users", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"format": c.FlagSet.String("format", "", "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
				"filter": c.FlagSet.String("filter", "*", "Properties to use when filtering, for example '+user_is_billable:0 +user_language=en'. Defaults to '*'. Valid filter properties are: "+strings.Join(userFilterProperties, ", ")+"."),
			}
		},
		ExecuteFunc:         userListCmd,
		Endpoint:            configuration.DeveloperEndpoint,
		PermissionsRequired: []string{command.USERS_AND_PERMISSIONS_READ},
	}, {
		Description:  "Shows information about a user.",
		Subject:      "user",
		Predicate:    "show",
		AltSubject:   "user",
		AltPredicate: "get",
		FlagSet:      flag.NewFlagSet("show user", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"user_id": c.FlagSet.String("user-id", "0", colors.Red("(Required)")+" The user ID."),
				"format":  c.FlagSet.String("format", "", "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
				"raw":     c.FlagSet.Bool("raw", false, colors.Green("(Flag)")+" If set returns the user raw object serialized using specified format"),
			}
		},
		ExecuteFunc2: userShowCmd,
	},

	{
		Description:  "Show user limits",
		Subject:      "user",
		Predicate:    "limits-get",
		AltSubject:   "limits",
		AltPredicate: "get",
		FlagSet:      flag.NewFlagSet("show user limits", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"user_id": c.FlagSet.String("user-id", "0", colors.Red("(Required)")+" The user ID."),
				"format":  c.FlagSet.String("format", "json", "The output format. Supported values are 'json','yaml'. The raw object will be returned"),
			}
		},
		ExecuteFunc2: userShowLimitsCmd,
	},
	{
		Description:  "Creates an user in the built-in user database",
		Subject:      "user",
		Predicate:    "create",
		AltSubject:   "user",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("create user", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"user":           c.FlagSet.String("user", command.NilDefaultStr, colors.Red("(Required)")+" The user to use. Usually in email format."),
				"display_name":   c.FlagSet.String("display-name", command.NilDefaultStr, "The user's display name."),
				"role":           c.FlagSet.String("role", command.NilDefaultStr, "The role to set the user to."),
				"account":        c.FlagSet.String("account", command.NilDefaultStr, "The account to assign the user to."),
				"limits_profile": c.FlagSet.String("limits-profile", command.NilDefaultStr, "Pass a JSON file name from which the system will read the limit values."),
				"set_billable":   c.FlagSet.Bool("set-billable", false, colors.Green("(Flag)")+" If set the user will be created with the billable flag set to True."),
				"set_verified":   c.FlagSet.Bool("set-verified", false, colors.Green("(Flag)")+" If set the user will be created with the email verified flag set to True."),
				"no_color":       c.FlagSet.Bool("no-color", false, colors.Green("(Flag)")+" Disable coloring."),
				"return_id":      c.FlagSet.Bool("return-id", false, colors.Green("(Flag)")+" If set will print the ID of the created user. Useful for automating tasks."),
			}
		},
		ExecuteFunc2: userCreateCmd,
	},

	{
		Description:  "Sets user as billable",
		Subject:      "user",
		Predicate:    "billable-flag-set",
		AltSubject:   "billable-flag",
		AltPredicate: "set",
		FlagSet:      flag.NewFlagSet("set user as billable", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"user_id": c.FlagSet.String("user-id", "", colors.Red("(Required)")+" The user ID."),
				"set":     c.FlagSet.Bool("set", false, "If set the user will be set as billable. If not set the user will be set as non-billable."),
			}
		},
		ExecuteFunc2: userSetBillableCmd,
	},
	{
		Description:  "Update user limits",
		Subject:      "user",
		Predicate:    "limits-set",
		AltSubject:   "limits",
		AltPredicate: "set",
		FlagSet:      flag.NewFlagSet("update user limits", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"user_id":        c.FlagSet.String("user-id", "0", colors.Red("(Required)")+" The user ID."),
				"limits_profile": c.FlagSet.String("limits-profile", command.NilDefaultStr, colors.Red("(Required)")+" Pass a JSON file name from which the system will read the limit values."),
			}
		},
		ExecuteFunc2: userUpdateLimitsCmd,
	},
	{
		Description:  "Archive a user",
		Subject:      "user",
		Predicate:    "archive",
		AltSubject:   command.NilDefaultStr,
		AltPredicate: command.NilDefaultStr,
		FlagSet:      flag.NewFlagSet("archive user", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"user_id": c.FlagSet.String("user-id", "", colors.Red("(Required)")+" The user ID."),
			}
		},
		ExecuteFunc2: userArchiveCmd,
	},
	{
		Description:  "Unarchive a user",
		Subject:      "user",
		Predicate:    "unarchive",
		AltSubject:   command.NilDefaultStr,
		AltPredicate: command.NilDefaultStr,
		FlagSet:      flag.NewFlagSet("unarchive user", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"user_id": c.FlagSet.String("user-id", "", colors.Red("(Required)")+" The user ID."),
			}
		},
		ExecuteFunc2: userUnarchiveCmd,
	},
	{
		Description:  "Set a account for a user",
		Subject:      "user",
		Predicate:    "account-set",
		AltSubject:   command.NilDefaultStr,
		AltPredicate: command.NilDefaultStr,
		FlagSet:      flag.NewFlagSet("set account for user", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"user_id":    c.FlagSet.String("user-id", "", colors.Red("(Required)")+" The user ID."),
				"account_id": c.FlagSet.String("account-id", "", colors.Red("(Required)")+" The account ID."),
			}
		},
		ExecuteFunc2: userSetAccountCmd,
	},
}

func userListCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	filter := command.GetStringParam(c.Arguments["filter"])
	iList, err := client.UserSearch(filter)
	if err != nil {
		return "", err
	}

	schema := []tableformatter.SchemaField{
		{
			FieldName: "ID",
			FieldType: tableformatter.TypeInt,
			FieldSize: 6,
		},
		{
			FieldName: "EMAIL",
			FieldType: tableformatter.TypeString,
			FieldSize: 20,
		},
		{
			FieldName: "NAME",
			FieldType: tableformatter.TypeString,
			FieldSize: 20,
		},
		{
			FieldName: "STATUS",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "LAST LOGIN",
			FieldType: tableformatter.TypeString,
			FieldSize: 10,
		},
		{
			FieldName: "CREATED",
			FieldType: tableformatter.TypeString,
			FieldSize: 10,
		},
	}

	data := [][]interface{}{}
	for _, i := range *iList {
		status := ""

	format := command.GetStringParam(c.Arguments["format"])
    if format == "" {
		if i.UserBlocked {
			status = colors.Red("Blocked")
		} else if i.UserIsSuspended {
			status = colors.Red("Suspended")
		} else {
			status = colors.Green("Active")
		}
	}else{
		if i.UserBlocked {
			status = "Blocked"
		} else if i.UserIsSuspended {
			status = "Suspended"
		} else {
			status = "Active"
		}
	}

		data = append(data, []interface{}{
			i.UserID,
			i.UserEmail,
			i.UserDisplayName,
			status,
			i.UserLastLoginTimestamp,
			i.UserCreatedTimestamp,
		})

	}

	tableformatter.TableSorter(schema).OrderBy(
		schema[3].FieldName,
		schema[0].FieldName,
		schema[1].FieldName).Sort(data)

	topLine := "Users"

	table := tableformatter.Table{
		Data:   data,
		Schema: schema,
	}
	return table.RenderTable("Users", topLine, command.GetStringParam(c.Arguments["format"]))
}

func userShowCmd(ctx context.Context, c *command.Command, client *metalcloud2.APIClient) (string, error) {
	userID := command.GetStringParam(c.Arguments["user_id"])
	tempID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return "", fmt.Errorf("invalid user id '%s'", userID)
	}
	numericID := int(tempID)
	result, _, err := client.UsersApi.GetUser(ctx, userID, nil)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return "", err
	}
	account := ""
	if result.AccountId != "" {
		accountInfo, _, err := client.AccountsApi.GetAccount(ctx, result.AccountId, nil)
		if err != nil {
			return "", fmt.Errorf("can't get account info: %w", err)
		}
		account = fmt.Sprintf("%s(%s)", accountInfo.Name, result.AccountId)
	}

	schema := []tableformatter.SchemaField{
		{
			FieldName: "ID",
			FieldType: tableformatter.TypeInt,
			FieldSize: 6,
		},
		{
			FieldName: "EMAIL",
			FieldType: tableformatter.TypeString,
			FieldSize: 20,
		},
		{
			FieldName: "NAME",
			FieldType: tableformatter.TypeString,
			FieldSize: 20,
		},
		{
			FieldName: "ACCOUNT",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "ACCESS LEVEL",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "CREATED",
			FieldType: tableformatter.TypeDateTime,
			FieldSize: 10,
		},
		{
			FieldName: "LAST LOGIN",
			FieldType: tableformatter.TypeDateTime,
			FieldSize: 10,
		},
		{
			FieldName: "IS BILLABLE",
			FieldType: tableformatter.TypeBool,
			FieldSize: 5,
		},
	}
	data := [][]interface{}{}
	data = append(data, []interface{}{
		numericID,
		result.Email,
		result.DisplayName,
		account,
		result.AccessLevel,
		result.CreatedTimestamp,
		result.LastLoginTimestamp,
		result.IsBillable,
	})

	format := command.GetStringParam(c.Arguments["format"])
	//check if format lowercase is one json, yaml or ""
	if !slices.Contains([]string{"json", "yaml", "JSON", "YAML", "csv", "CSV", ""}, format) {
		return "", fmt.Errorf("invalid format '%s' given. Valid values are json, JSON, yaml or YAML", format)
	}
	var sb strings.Builder

	if command.GetBoolParam(c.Arguments["raw"]) {
		ret, err := tableformatter.RenderRawObject(result, format, "User")
		if err != nil {
			return "", err
		}
		sb.WriteString(ret)
	} else {

		table := tableformatter.Table{
			Data:   data,
			Schema: schema,
		}

		ret, err := table.RenderTransposedTable("user details", "", format)
		if err != nil {
			return "", err
		}

		sb.WriteString(ret)
	}
	return sb.String(), nil

}

func userShowLimitsCmd(ctx context.Context, c *command.Command, client *metalcloud2.APIClient) (string, error) {
	userID := command.GetStringParam(c.Arguments["user_id"])
	result, _, err := client.UsersApi.GetUserLimits(ctx, userID)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return "", err
	}
	format := command.GetStringParam(c.Arguments["format"])
	if format == "" {
		return "", fmt.Errorf("format is required")
	}
	if !slices.Contains([]string{"json", "yaml", "JSON", "YAML"}, format) {
		return "", fmt.Errorf("invalid format '%s' given. Valid values are json, JSON, yaml or YAML", format)
	}
	var sb strings.Builder

	ret, err := tableformatter.RenderRawObject(result, format, "User")
	if err != nil {
		return "", err
	}
	sb.WriteString(ret)
	return sb.String(), nil
}

func userCreateCmd(ctx context.Context, c *command.Command, client *metalcloud2.APIClient) (string, error) {
	userName, ok := command.GetStringParamOk(c.Arguments["user"])
	if !ok {
		return "", fmt.Errorf("user is required")
	}
	displayName, ok := command.GetStringParamOk(c.Arguments["display_name"])
	if !ok {
		displayName = userName // TODO: check if this is correct
	}
	role := command.GetStringParam(c.Arguments["role"])
	account := command.GetStringParam(c.Arguments["account"])
	accountId, _ := strconv.ParseInt(account, 10, 64)
	limitsProfile := command.GetStringParam(c.Arguments["limits_profile"])
	emailVerified := command.GetBoolParam(c.Arguments["set_verified"])

	ret, _, err := client.UsersApi.CreateUser(ctx, metalcloud2.CreateUserDto{
		DisplayName:   displayName,
		Email:         userName,
		AccessLevel:   role,
		EmailVerified: emailVerified,
		AccountId:     float64(accountId),
	})

	if err != nil {
		return "", err
	}
	id := ret.Id

	if limitsProfile != "" {
		limits, err := readLimitsProfile(c)
		if err != nil {
			return "", err
		}
		_, _, err = client.UsersApi.UpdateUserLimits(ctx, limits, id)
		if err != nil {
			return "", fmt.Errorf("can't set user limits: %w", err)
		}
	}
	if c.Arguments["set_billable"] != nil && *c.Arguments["set_billable"].(*bool) {
		billable := true
		_, _, err = client.UsersApi.UpdateUser(ctx, metalcloud2.UpdateUserDto{
			IsBillable: &billable,
		}, id)
		if err != nil {
			return "", fmt.Errorf("can't set user billable: %w", err)
		}
	}
	if c.Arguments["return_id"] != nil && *c.Arguments["return_id"].(*bool) {
		return ret.Id, nil
	}

	return "", nil
}

func readLimitsProfile(cmd *command.Command) (metalcloud2.UserLimitsDto, error) {
	limitsProfile, ok := command.GetStringParamOk(cmd.Arguments["limits_profile"])
	if !ok {
		return metalcloud2.UserLimitsDto{}, fmt.Errorf("limits-profile is required %w", nil)
	}
	var err error
	var data []byte

	if limitsProfile == "-" {
		data, err = configuration.ReadInputFromPipe()
	} else {
		data, err = configuration.ReadInputFromFile(limitsProfile)
	}
	if err != nil {
		return metalcloud2.UserLimitsDto{}, fmt.Errorf("can't read limits profile: %w", err)
	}

	var limits metalcloud2.UserLimitsDto
	err = json.Unmarshal(data, &limits)
	if err != nil {
		return metalcloud2.UserLimitsDto{}, fmt.Errorf("can't read limits profile: %w", err)
	}
	return limits, nil
}

func userSetBillableCmd(ctx context.Context, c *command.Command, client *metalcloud2.APIClient) (string, error) {
	userID, ok := command.GetStringParamOk(c.Arguments["user_id"])
	if !ok || userID == "" {
		return "", fmt.Errorf("user-id is required")
	}

	billable := command.GetBoolParam(c.Arguments["set"])
	_, _, err := client.UsersApi.UpdateUser(ctx, metalcloud2.UpdateUserDto{
		IsBillable: &billable,
	}, userID)
	if err != nil {
		return "", fmt.Errorf("can't set user billable: %w", err)
	}
	return "", nil
}

func userUpdateLimitsCmd(ctx context.Context, c *command.Command, client *metalcloud2.APIClient) (string, error) {
	userID, ok := command.GetStringParamOk(c.Arguments["user_id"])
	if !ok {
		return "", fmt.Errorf("user-id is required")
	}

	limits, err := readLimitsProfile(c)
	if err != nil {
		return "", fmt.Errorf("can't set user limits: %w", err)
	}
	_, _, err = client.UsersApi.UpdateUserLimits(ctx, limits, userID)
	if err != nil {
		if swaggerErr, ok := err.(metalcloud2.GenericSwaggerError); ok {
			return "", fmt.Errorf("can't set user limits: %w, %s", err, string(swaggerErr.Body()))
		}
		return "", fmt.Errorf("can't set user limits: %w", err)
	}
	return "", nil
}

func userArchiveCmd(ctx context.Context, c *command.Command, client *metalcloud2.APIClient) (string, error) {
	userID, ok := command.GetStringParamOk(c.Arguments["user_id"])
	if !ok {
		return "", fmt.Errorf("user-id is required")
	}
	_, _, err := client.UsersApi.ArchiveUser(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("can't archive user: %w", err)
	}
	return "", nil
}

func userUnarchiveCmd(ctx context.Context, c *command.Command, client *metalcloud2.APIClient) (string, error) {
	userID, ok := command.GetStringParamOk(c.Arguments["user_id"])
	if !ok {
		return "", fmt.Errorf("user-id is required")
	}
	_, _, err := client.UsersApi.UnarchiveUser(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("can't unarchive user: %w", err)
	}
	return "", nil
}

func userSetAccountCmd(ctx context.Context, c *command.Command, client *metalcloud2.APIClient) (string, error) {
	userID, ok := command.GetStringParamOk(c.Arguments["user_id"])
	if !ok {
		return "", fmt.Errorf("user-id is required")
	}
	accountID, ok := command.GetStringParamOk(c.Arguments["account_id"])
	if !ok {
		return "", fmt.Errorf("account-id is required")
	}
	aID, err := strconv.ParseFloat(accountID, 64)
	if err != nil {
		return "", fmt.Errorf("invalid account id '%s'", accountID)
	}
	_, _, err = client.UsersApi.ChangeUserAccount(ctx, metalcloud2.ChangeUserAccountDto{
		NewAccountId: aID,
	}, userID)
	if err != nil {
		return "", fmt.Errorf("can't set user account: %w", err)
	}
	return "", nil
}
