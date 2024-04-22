package datacenter

import (
	"flag"
	"fmt"
	"strings"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
	"github.com/metalsoft-io/metalcloud-cli/internal/colors"
	"github.com/metalsoft-io/metalcloud-cli/internal/command"
	"github.com/metalsoft-io/metalcloud-cli/internal/configuration"
	"github.com/metalsoft-io/metalcloud-cli/internal/objects"
	"github.com/metalsoft-io/tableformatter"
)

var DatacenterCmds = []command.Command{
	{
		Description:  "List all datacenters.",
		Subject:      "datacenter",
		AltSubject:   "dc",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list datacenters", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"user_id":       c.FlagSet.String("user", command.NilDefaultStr, "List only specific user's datacenters"),
				"show_inactive": c.FlagSet.Bool("show-inactive", false, colors.Green("(Flag)")+" Set flag if inactive datacenters are to be returned"),
				"show_hidden":   c.FlagSet.Bool("show-hidden", false, colors.Green("(Flag)")+" Set flag if hidden datacenters are to be returned"),
				"format":        c.FlagSet.String("format", command.NilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
			}
		},
		ExecuteFunc: datacenterListCmd,
	},
	{
		Description:  "Create a datacenter.",
		Subject:      "datacenter",
		AltSubject:   "dc",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("Create datacenter", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"read_config_from_file": c.FlagSet.String("f", command.NilDefaultStr, colors.Red("(Required)")+" Read configuration from file in the format specified with --format."),
				"format":                c.FlagSet.String("format", "yaml", "The input format. Supported values are 'json','yaml'. The default format is json."),
				"return_id":             c.FlagSet.Bool("return-id", false, "Will print the ID of the created Datacenter Useful for automating tasks."),
			}
		},
		ExecuteFunc: datacenterCreateCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
	{
		Description:  "Get datacenter",
		Subject:      "datacenter",
		AltSubject:   "dc",
		Predicate:    "get",
		AltPredicate: "show",
		FlagSet:      flag.NewFlagSet("Get datacenter details.", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"datacenter_name":   c.FlagSet.String("id", command.NilDefaultStr, colors.Red("(Required)")+" Label of the datacenter. Also used as an ID."),
				"format":            c.FlagSet.String("format", "yaml", "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
			}
		},
		ExecuteFunc: datacenterGetCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
	{
		Description:  "Get datacenter config URL",
		Subject:      "datacenter",
		AltSubject:   "dc",
		Predicate:    "get-config-url",
		AltPredicate: "show-config-url",
		FlagSet:      flag.NewFlagSet("Get datacenter config URL.", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"datacenter_name":   c.FlagSet.String("id", command.NilDefaultStr, colors.Red("(Required)")+" Label of the datacenter. Also used as an ID."),
			}
		},
		ExecuteFunc: datacenterGetConfigURLCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
	{
		Description:  "Update datacenter config",
		Subject:      "datacenter",
		AltSubject:   "dc",
		Predicate:    "update",
		AltPredicate: "edit",
		FlagSet:      flag.NewFlagSet("Update datacenter config", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"datacenter_name":       c.FlagSet.String("id", command.NilDefaultStr, colors.Red("(Required)")+" Label of the datacenter. Also used as an ID."),
				"read_config_from_file": c.FlagSet.String("f", command.NilDefaultStr, colors.Red("(Required)")+" Read object configuration from file"),
				"format":                c.FlagSet.String("format", "yaml", "The input format. Supported values are 'json','yaml'. The default format is json."),
			}
		},
		ExecuteFunc: datacenterUpdateCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
}

func datacenterListCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	showHidden := command.GetBoolParam(c.Arguments["show_hidden"])
	showInactive := command.GetBoolParam(c.Arguments["show_inactive"])
	userID, userIDProvided := command.GetStringParamOk(c.Arguments["user_id"])

	var dList *map[string]metalcloud.Datacenter
	var err error

	if userIDProvided {
		if id, label, isID := command.IdOrLabelString(userID); isID {
			dList, err = client.DatacentersByUserID(id, !showInactive)
		} else {
			dList, err = client.DatacentersByUserEmail(label, !showInactive)
		}
	} else {
		dList, err = client.Datacenters(!showInactive)
	}

	if err != nil {
		return "", err
	}

	schema := []tableformatter.SchemaField{
		{
			FieldName: "LABEL",
			FieldType: tableformatter.TypeString,
			FieldSize: 15,
		},
		{
			FieldName: "NAME",
			FieldType: tableformatter.TypeString,
			FieldSize: 20,
		},
		{
			FieldName: "OWNER",
			FieldType: tableformatter.TypeString,
			FieldSize: 10,
		},
		{
			FieldName: "PARENT",
			FieldType: tableformatter.TypeString,
			FieldSize: 10,
		},
		{
			FieldName: "FLAGS",
			FieldType: tableformatter.TypeString,
			FieldSize: 20,
		},
	}

	data := [][]interface{}{}
	for _, dc := range *dList {

		if dc.DatacenterHidden && !showHidden {
			continue
		}

		flags := []string{}

		if dc.DatacenterIsMaster {
			flags = append(flags, "MASTER")
		}

		if dc.DatacenterIsMaintenance {
			flags = append(flags, "MAINTENANCE")
		}

		if dc.DatacenterHidden {
			flags = append(flags, "HIDDEN")
		}

		flags = append(flags, dc.DatacenterTags...)

		userStr := ""
		if dc.UserID != 0 {
			user, err := client.UserGet(dc.UserID)
			if err != nil {
				return "", err
			}
			userStr = fmt.Sprintf("%s #%d", user.UserEmail, dc.UserID)
		}

		data = append(data, []interface{}{
			dc.DatacenterName,
			dc.DatacenterDisplayName,
			userStr,
			dc.DatacenterNameParent,
			strings.Join(flags, " "),
		})

	}
	table := tableformatter.Table{
		Data:   data,
		Schema: schema,
	}

	return table.RenderTable("Datacenters", "", command.GetStringParam(c.Arguments["format"]))
}

func datacenterCreateCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	obj, err := objects.ReadSingleObjectFromCommand(c, client)
	if err != nil {
		return "", err
	}
	d := (*obj).(metalcloud.DatacenterWithConfig)

	createdDC, err := client.DatacenterCreateFromDatacenterWithConfig(d)
	if err != nil {
		return "", err
	}

	if command.GetBoolParam(c.Arguments["return_id"]) {
		return fmt.Sprintf("%d", createdDC.Metadata.DatacenterID), nil
	}

	return "", err
}

func datacenterGetCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	datacenterName, ok := command.GetStringParamOk(c.Arguments["datacenter_name"])
	if !ok {
		return "", fmt.Errorf("-id required")
	}

	retDC, err := client.DatacenterWithConfigGet(datacenterName)
	if err != nil {
		return "", err
	}

	format := command.GetStringParam(c.Arguments["format"])
	ret, err := objects.RenderRawObject(*retDC, format, "DatacenterWithConfig")
	if err != nil {
		return "", err
	}

	return ret, nil
}

func datacenterGetConfigURLCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	datacenterName, ok := command.GetStringParamOk(c.Arguments["datacenter_name"])
	if !ok {
		return "", fmt.Errorf("-id required")
	}

	secretConfigURL, err := client.DatacenterAgentsConfigJSONDownloadURL(datacenterName, true)
	if err != nil {
		return "", err
	}

	return secretConfigURL, nil
}

func datacenterUpdateCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	obj, err := objects.ReadSingleObjectFromCommand(c, client)
	if err != nil {
		return "", err
	}
	d := (*obj).(metalcloud.DatacenterWithConfig)

	_, err = client.DatacenterUpdateFromDatacenterWithConfig(d)
	if err != nil {
		return "", err
	}

	return "", err
}
