package user

import (
	"flag"
	"strings"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
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

		if i.UserBlocked {
			status = colors.Red("Blocked")
		} else if i.UserIsSuspended {
			status = colors.Red("Suspended")
		} else {
			status = colors.Green("Active")
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
