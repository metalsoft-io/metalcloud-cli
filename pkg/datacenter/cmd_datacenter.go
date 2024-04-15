package datacenter

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
	"github.com/metalsoft-io/metalcloud-cli/internal/colors"
	"github.com/metalsoft-io/metalcloud-cli/internal/command"
	"github.com/metalsoft-io/metalcloud-cli/internal/configuration"
	"github.com/metalsoft-io/metalcloud-cli/internal/objects"
	"github.com/metalsoft-io/tableformatter"
	"gopkg.in/yaml.v3"
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
				"format":        c.FlagSet.String("format", "", "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
				"json_path":     c.FlagSet.String("jsonpath", command.NilDefaultStr, "Filter the output."),
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
				"datacenter_name":         c.FlagSet.String("id", command.NilDefaultStr, colors.Red("(Required)")+" Label of the datacenter. Also used as an ID."),
				"datacenter_display_name": c.FlagSet.String("title", command.NilDefaultStr, colors.Red("(Required)")+" Human readable name of the datacenter. Usually includes the location such as UK,Reading"),
				"read_config_from_file":   c.FlagSet.String("raw-config", command.NilDefaultStr, colors.Red("(Required)")+" Read datacenter configuration from file"),
				"datacenter_name_parent":  c.FlagSet.String("parent", command.NilDefaultStr, "If the datacenter is subordonated to another datacenter such as to a near-edge site."),
				"create_hidden":           c.FlagSet.Bool("hidden", false, colors.Green("(Flag)")+" If set, the datacenter will be hidden after creation instead."),
				"is_master":               c.FlagSet.Bool("master", false, colors.Green("(Flag)")+" If set, the datacenter will be the master dc."),
				"is_maintenance":          c.FlagSet.Bool("maintenance", false, colors.Green("(Flag)")+" If set, the datacenter will be in maintenance."),
				"user_id":                 c.FlagSet.String("user", command.NilDefaultStr, "Datacenter's owner. If ommited, the default is a public datacenter."),
				"tags":                    c.FlagSet.String("tags", command.NilDefaultStr, "Tags associated with this datacenter, comma separated"),
				"read_config_from_pipe":   c.FlagSet.Bool("pipe", false, colors.Green("(Flag)")+" If set, read datacenter configuration from pipe instead of from a file. Either this flag or the -config option must be used."),
				"format":                  c.FlagSet.String("format", "json", "The input format. Supported values are 'json','yaml'. The default format is json."),
				"return_id":               c.FlagSet.Bool("return-id", false, "Will print the ID of the created Datacenter Useful for automating tasks."),
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
				"datacenter_name": c.FlagSet.String("id", command.NilDefaultStr, colors.Red("(Required)")+" Label of the datacenter. Also used as an ID."),
				//	"show_secret_config_url": c.FlagSet.Bool("show-config-url", false, colors.Green("(Flag)")+" If set returns the secret config url for datacenter agents."),
				"return_config_url": c.FlagSet.Bool("return-config-url", false, colors.Green("(Flag)")+" If set prints the config url of the datacenter. Ignores all other flags. Useful in automation."),
				"format":            c.FlagSet.String("format", "yaml", "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
				"json_path":         c.FlagSet.String("jsonpath", command.NilDefaultStr, "Filter the JSON config."),
			}
		},
		ExecuteFunc: datacenterGetCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
	{
		Description:  "Get site controller config url",
		Subject:      "datacenter",
		AltSubject:   "dc",
		Predicate:    "get",
		AltPredicate: "show",
		FlagSet:      flag.NewFlagSet("Get site controller config url.", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"datacenter_name": c.FlagSet.String("id", command.NilDefaultStr, colors.Red("(Required)")+" Label of the datacenter. Also used as an ID."),
			}
		},
		ExecuteFunc: datacenterGetConfigUrlCmd,
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
				"read_config_from_file": c.FlagSet.String("raw-config", command.NilDefaultStr, colors.Red("(Required)")+" Read datacenter configuration from file"),
				"read_config_from_pipe": c.FlagSet.Bool("pipe", false, colors.Green("(Flag)")+" If set, read datacenter configuration from pipe instead of from a file. Either this flag or the -config option must be used."),
				"format":                c.FlagSet.String("format", "json", "The input format. Supported values are 'json','yaml'. The default format is json."),
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

	datacenterName, ok := command.GetStringParamOk(c.Arguments["datacenter_name"])

	if !ok {
		return "", fmt.Errorf("id is required")
	}

	datacenterDisplayName, ok := command.GetStringParamOk(c.Arguments["datacenter_display_name"])
	if !ok {
		return "", fmt.Errorf("title is required")
	}

	userID := 0

	userStr, ok := command.GetStringParamOk(c.Arguments["user_id"])
	if ok {
		if id, label, isID := command.IdOrLabelString(userStr); isID {
			userID = id
		} else {
			user, err := client.UserGetByEmail(label)
			if err != nil {
				return "", err
			}
			userID = user.UserID
		}
	}

	datacenterHidden := command.GetBoolParam(c.Arguments["create_hidden"])
	datacenterIsMaster := command.GetBoolParam(c.Arguments["is_master"])
	datacenterInMaintenance := command.GetBoolParam(c.Arguments["is_maintenance"])

	datacenterTags := strings.Split(command.GetStringParam(c.Arguments["tags"]), ",")
	datacenterParent := command.GetStringParam(c.Arguments["datacenter_name_parent"])

	dc := metalcloud.Datacenter{
		DatacenterName:          datacenterName,
		DatacenterDisplayName:   datacenterDisplayName,
		UserID:                  userID,
		DatacenterIsMaster:      datacenterIsMaster,
		DatacenterIsMaintenance: datacenterInMaintenance,
		DatacenterHidden:        datacenterHidden,
		DatacenterTags:          datacenterTags,
		DatacenterNameParent:    datacenterParent,
	}

	readContentfromPipe := command.GetBoolParam((c.Arguments["read_config_from_pipe"]))

	var err error
	content := []byte{}

	if readContentfromPipe {
		content, err = configuration.ReadInputFromPipe()
	} else {

		if configFilePath, ok := command.GetStringParamOk(c.Arguments["read_config_from_file"]); ok {

			content, err = configuration.ReadInputFromFile(configFilePath)
		} else {
			return "", fmt.Errorf("-raw-config <path_to_json_file> or -pipe is required")
		}
	}

	if err != nil {
		return "", err
	}

	if len(content) == 0 {
		return "", fmt.Errorf("Content cannot be empty")
	}

	format := command.GetStringParam(c.Arguments["format"])

	var dcConf metalcloud.DatacenterConfig
	switch format {
	case "json":
		err := json.Unmarshal(content, &dcConf)
		if err != nil {
			return "", err
		}
	case "yaml":
		err := yaml.Unmarshal(content, &dcConf)
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("input format \"%s\" not supported", format)
	}

	ret, err := client.DatacenterCreate(dc, dcConf)
	if err != nil {
		return "", err
	}

	if c.Arguments["return_id"] != nil && *c.Arguments["return_id"].(*bool) {
		return fmt.Sprintf("%s", ret.DatacenterName), nil
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

	var sb strings.Builder

	format := command.GetStringParam(c.Arguments["format"])

	ret, err := objects.RenderRawObject(*retDC, format, "DatacenterWithConfig")
	if err != nil {
		return "", err
	}
	sb.WriteString(ret)

	return sb.String(), nil
}

func datacenterGetConfigUrlCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	datacenterName, ok := command.GetStringParamOk(c.Arguments["datacenter_name"])
	if !ok {
		return "", fmt.Errorf("-id required")
	}

	return client.DatacenterAgentsConfigJSONDownloadURL(datacenterName, true)
}

func datacenterUpdateCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	datacenterName, ok := command.GetStringParamOk(c.Arguments["datacenter_name"])

	if !ok {
		return "", fmt.Errorf("id is required")
	}

	readContentfromPipe := command.GetBoolParam((c.Arguments["read_config_from_pipe"]))

	var err error
	content := []byte{}

	if readContentfromPipe {
		content, err = configuration.ReadInputFromPipe()
	} else {

		if configFilePath, ok := command.GetStringParamOk(c.Arguments["read_config_from_file"]); ok {

			content, err = configuration.ReadInputFromFile(configFilePath)
		} else {
			return "", fmt.Errorf("-raw-config <path_to_json_file> or -pipe is required")
		}
	}

	if err != nil {
		return "", err
	}

	if len(content) == 0 {
		return "", fmt.Errorf("Content cannot be empty")
	}

	format := command.GetStringParam(c.Arguments["format"])

	var dcConf metalcloud.DatacenterConfig
	switch format {
	case "json":
		err := json.Unmarshal(content, &dcConf)
		if err != nil {
			return "", err
		}
	case "yaml":
		err := yaml.Unmarshal(content, &dcConf)
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("input format \"%s\" not supported", format)
	}

	err = client.DatacenterConfigUpdate(datacenterName, dcConf)
	if err != nil {
		return "", err
	}

	return "", err
}
