package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
	"github.com/metalsoft-io/tableformatter"
	jq "github.com/savaki/jq"
	"gopkg.in/yaml.v3"
)

//infrastructureCmds commands affecting infrastructures
var datacenterCmds = []Command{

	{
		Description:  "List all datacenters.",
		Subject:      "datacenter",
		AltSubject:   "dc",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list datacenters", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"user_id":       c.FlagSet.String("user", _nilDefaultStr, "List only specific user's datacenters"),
				"show_inactive": c.FlagSet.Bool("show-inactive", false, green("(Flag)")+" Set flag if inactive datacenters are to be returned"),
				"show_hidden":   c.FlagSet.Bool("show-hidden", false, green("(Flag)")+" Set flag if hidden datacenters are to be returned"),
				"format":        c.FlagSet.String("format", "", "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
				"json_path":     c.FlagSet.String("jsonpath", _nilDefaultStr, "Filter the output."),
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
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"datacenter_name":         c.FlagSet.String("id", _nilDefaultStr, red("(Required)")+" Label of the datacenter. Also used as an ID."),
				"datacenter_display_name": c.FlagSet.String("title", _nilDefaultStr, red("(Required)")+" Human readable name of the datacenter. Usually includes the location such as UK,Reading"),
				"read_config_from_file":   c.FlagSet.String("config", _nilDefaultStr, red("(Required)")+" Read datacenter configuration from file"),
				"datacenter_name_parent":  c.FlagSet.String("parent", _nilDefaultStr, "If the datacenter is subordonated to another datacenter such as to a near-edge site."),
				"create_hidden":           c.FlagSet.Bool("hidden", false, green("(Flag)")+" If set, the datacenter will be hidden after creation instead."),
				"is_master":               c.FlagSet.Bool("master", false, green("(Flag)")+" If set, the datacenter will be the master dc."),
				"is_maintenance":          c.FlagSet.Bool("maintenance", false, green("(Flag)")+" If set, the datacenter will be in maintenance."),
				"user_id":                 c.FlagSet.String("user", _nilDefaultStr, "Datacenter's owner. If ommited, the default is a public datacenter."),
				"tags":                    c.FlagSet.String("tags", _nilDefaultStr, "Tags associated with this datacenter, comma separated"),
				"read_config_from_pipe":   c.FlagSet.Bool("pipe", false, green("(Flag)")+" If set, read datacenter configuration from pipe instead of from a file. Either this flag or the -config option must be used."),
				"format":                  c.FlagSet.String("format", "json", "The input format. Supported values are 'json','yaml'. The default format is json."),
				"return_id":               c.FlagSet.Bool("return-id", false, "Will print the ID of the created Datacenter Useful for automating tasks."),
			}
		},
		ExecuteFunc: datacenterCreateCmd,
		Endpoint:    DeveloperEndpoint,
	},
	{
		Description:  "Get datacenter",
		Subject:      "datacenter",
		AltSubject:   "dc",
		Predicate:    "get",
		AltPredicate: "show",
		FlagSet:      flag.NewFlagSet("Get datacenter details.", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"datacenter_name":        c.FlagSet.String("id", _nilDefaultStr, red("(Required)")+" Label of the datacenter. Also used as an ID."),
				"show_secret_config_url": c.FlagSet.Bool("show-config-url", false, green("(Flag)")+" If set returns the secret config url for datacenter agents."),
				"return_config_url":      c.FlagSet.Bool("return-config-url", false, green("(Flag)")+" If set prints the config url of the datacenter. Ignores all other flags. Useful in automation."),
				"format":                 c.FlagSet.String("format", "", "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
				"json_path":              c.FlagSet.String("jsonpath", _nilDefaultStr, "Filter the JSON config."),
			}
		},
		ExecuteFunc: datacenterGetCmd,
		Endpoint:    DeveloperEndpoint,
	},
	{
		Description:  "Update datacenter config",
		Subject:      "datacenter",
		AltSubject:   "dc",
		Predicate:    "update",
		AltPredicate: "edit",
		FlagSet:      flag.NewFlagSet("Update datacenter config", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"datacenter_name":       c.FlagSet.String("id", _nilDefaultStr, red("(Required)")+" Label of the datacenter. Also used as an ID."),
				"read_config_from_file": c.FlagSet.String("config", _nilDefaultStr, red("(Required)")+" Read datacenter configuration from file"),
				"read_config_from_pipe": c.FlagSet.Bool("pipe", false, green("(Flag)")+" If set, read datacenter configuration from pipe instead of from a file. Either this flag or the -config option must be used."),
				"format":                c.FlagSet.String("format", "json", "The input format. Supported values are 'json','yaml'. The default format is json."),
			}
		},
		ExecuteFunc: datacenterUpdateCmd,
		Endpoint:    DeveloperEndpoint,
	},
}

func datacenterListCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	showHidden := getBoolParam(c.Arguments["show_hidden"])
	showInactive := getBoolParam(c.Arguments["show_inactive"])
	userID, userIDProvided := getStringParamOk(c.Arguments["user_id"])

	var dList *map[string]metalcloud.Datacenter
	var err error

	if userIDProvided {
		if id, label, isID := idOrLabelString(userID); isID {
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

	return table.RenderTable("Datacenters", "", getStringParam(c.Arguments["format"]))
}

func datacenterCreateCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	datacenterName, ok := getStringParamOk(c.Arguments["datacenter_name"])

	if !ok {
		return "", fmt.Errorf("id is required")
	}

	datacenterDisplayName, ok := getStringParamOk(c.Arguments["datacenter_display_name"])
	if !ok {
		return "", fmt.Errorf("title is required")
	}

	userID := 0

	userStr, ok := getStringParamOk(c.Arguments["user_id"])
	if ok {
		if id, label, isID := idOrLabelString(userStr); isID {
			userID = id
		} else {
			user, err := client.UserGetByEmail(label)
			if err != nil {
				return "", err
			}
			userID = user.UserID
		}
	}

	datacenterHidden := getBoolParam(c.Arguments["create_hidden"])
	datacenterIsMaster := getBoolParam(c.Arguments["is_master"])
	datacenterInMaintenance := getBoolParam(c.Arguments["is_maintenance"])

	datacenterTags := strings.Split(getStringParam(c.Arguments["tags"]), ",")
	datacenterParent := getStringParam(c.Arguments["datacenter_name_parent"])

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

	readContentfromPipe := getBoolParam((c.Arguments["read_config_from_pipe"]))

	var err error
	content := []byte{}

	if readContentfromPipe {
		content, err = readInputFromPipe()
	} else {

		if configFilePath, ok := getStringParamOk(c.Arguments["read_config_from_file"]); ok {

			content, err = readInputFromFile(configFilePath)
		} else {
			return "", fmt.Errorf("-config <path_to_json_file> or -pipe is required")
		}
	}

	if err != nil {
		return "", err
	}

	if len(content) == 0 {
		return "", fmt.Errorf("Content cannot be empty")
	}

	format := getStringParam(c.Arguments["format"])

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

func datacenterGetCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	datacenterName, ok := getStringParamOk(c.Arguments["datacenter_name"])
	if !ok {
		return "", fmt.Errorf("-id required")
	}

	retDC, err := client.DatacenterGet(datacenterName)
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
			FieldName: "TITLE",
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

	flags := []string{}

	if retDC.DatacenterIsMaster {
		flags = append(flags, "MASTER")
	}

	if retDC.DatacenterIsMaintenance {
		flags = append(flags, "MAINTENANCE")
	}

	if retDC.DatacenterHidden {
		flags = append(flags, "HIDDEN")
	}

	flags = append(flags, retDC.DatacenterTags...)

	userStr := ""
	if retDC.UserID != 0 {
		user, err := client.UserGet(retDC.UserID)
		if err != nil {
			return "", err
		}
		userStr = fmt.Sprintf("%s #%d", user.UserEmail, retDC.UserID)
	}

	showSecretURL := getBoolParam(c.Arguments["show_secret_config_url"])
	secretConfigURL := ""

	if showSecretURL || getBoolParam(c.Arguments["return_config_url"]) {
		schema = append(schema, tableformatter.SchemaField{
			FieldName: "CONFIG_URL",
			FieldType: tableformatter.TypeString,
			FieldSize: 15,
		})
		secretConfigURL, err = client.DatacenterAgentsConfigJSONDownloadURL(datacenterName, true)
		if err != nil {
			return "", err
		}
	}

	format := getStringParam(c.Arguments["format"])

	configStr := ""
	config := metalcloud.DatacenterConfig{}

	schema = append(schema, tableformatter.SchemaField{
		FieldName: "CONFIG",
		FieldType: tableformatter.TypeInterface,
		FieldSize: 15,
	})

	configRet, err := client.DatacenterConfigGet(datacenterName)
	if err != nil {
		return "", err
	}

	config = *configRet

	configBytes, err := yaml.Marshal(*configRet)
	if err != nil {
		return "", err
	}

	configStr = string(configBytes)

	data := [][]interface{}{
		{
			retDC.DatacenterName,
			retDC.DatacenterDisplayName,
			userStr,
			retDC.DatacenterNameParent,
			strings.Join(flags, " "),
			secretConfigURL,
			config,
		},
	}

	var sb strings.Builder

	switch format {
	case "json", "JSON":
		configBytes, err := json.MarshalIndent(config, "", "\t")
		if err != nil {
			return "", err
		}

		if jsonPath, ok := getStringParamOk(c.Arguments["json_path"]); ok {
			op, err := jq.Parse(jsonPath)

			if err != nil {
				return "", err
			}
			result, _ := op.Apply(configBytes)
			data[0][5] = string(result)
		} else {
			data[0][5] = string(configBytes)
		}

		table := tableformatter.Table{
			Data:   data,
			Schema: schema,
		}
		ret, err := table.RenderTableAsJSON()
		if err != nil {
			return "", err
		}
		sb.WriteString(ret)
	case "csv", "CSV":
		table := tableformatter.Table{
			Data:   data,
			Schema: schema,
		}
		ret, err := table.RenderTableAsCSV()
		if err != nil {
			return "", err
		}
		sb.WriteString(ret)

	case "yaml", "YAML":
		configBytes, err := yaml.Marshal(config)
		if err != nil {
			return "", err
		}

		data[0][5] = string(configBytes)

		table := tableformatter.Table{
			Data:   data,
			Schema: schema,
		}
		ret, err := table.RenderTableAsYAML()
		if err != nil {
			return "", err
		}
		sb.WriteString(ret)

	default:

		if getBoolParam(c.Arguments["return_config_url"]) {
			return secretConfigURL, nil
		}

		sb.WriteString("DATACENTER OVERVIEW\n")
		sb.WriteString("-------------------\n")

		sb.WriteString(fmt.Sprintf("Datacenter name (label): %s\n", retDC.DatacenterName))
		sb.WriteString(fmt.Sprintf("Datacenter display name (title): %s\n", retDC.DatacenterDisplayName))
		sb.WriteString(fmt.Sprintf("User: %s\n", userStr))
		sb.WriteString(fmt.Sprintf("Flags: %s\n", strings.Join(flags, " ")))
		sb.WriteString(fmt.Sprintf("Parent: %s\n", retDC.DatacenterNameParent))
		sb.WriteString(fmt.Sprintf("Type: %s\n", retDC.DatacenterType))
		sb.WriteString(fmt.Sprintf("Created: %s\n", retDC.DatacenterCreatedTimestamp))
		sb.WriteString(fmt.Sprintf("Updated: %s\n", retDC.DatacenterUpdatedTimestamp))

		sb.WriteString("\n")
		sb.WriteString("Configuration\n")
		sb.WriteString("---------------\n")
		sb.WriteString(configStr)

		if showSecretURL {
			sb.WriteString("---------------\n")
			sb.WriteString(fmt.Sprintf("Datacenter Agents Secret Config URL: %s\n", secretConfigURL))
		}

	}

	return sb.String(), nil
}

func datacenterUpdateCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	datacenterName, ok := getStringParamOk(c.Arguments["datacenter_name"])

	if !ok {
		return "", fmt.Errorf("id is required")
	}

	readContentfromPipe := getBoolParam((c.Arguments["read_config_from_pipe"]))

	var err error
	content := []byte{}

	if readContentfromPipe {
		content, err = readInputFromPipe()
	} else {

		if configFilePath, ok := getStringParamOk(c.Arguments["read_config_from_file"]); ok {

			content, err = readInputFromFile(configFilePath)
		} else {
			return "", fmt.Errorf("-config <path_to_json_file> or -pipe is required")
		}
	}

	if err != nil {
		return "", err
	}

	if len(content) == 0 {
		return "", fmt.Errorf("Content cannot be empty")
	}

	format := getStringParam(c.Arguments["format"])

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
