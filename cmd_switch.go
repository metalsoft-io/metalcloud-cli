package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	metalcloud "github.com/bigstepinc/metal-cloud-sdk-go"
	"github.com/metalsoft-io/tableformatter"
)

var switchCmds = []Command{

	{
		Description:  "Lists registered switches.",
		Subject:      "switch",
		AltSubject:   "sw",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list switches", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"format":           c.FlagSet.String("format", _nilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
				"datacenter_name":  c.FlagSet.String("datacenter", "", "The optional parameter acts as a filter that restricts the returned results to switch devices located in the specified datacenter."),
				"switch_type":      c.FlagSet.String("switch-type", "", "The optional parameter acts as a filter that restricts the returned results to switch devices of the specified type."),
				"show_credentials": c.FlagSet.Bool("show-credentials", false, "(Flag) If set returns the switch management credentials. (Slow for large queries)"),
			}
		},
		ExecuteFunc: switchListCmd,
		Endpoint:    DeveloperEndpoint,
	},
	{
		Description:  "Create switch device.",
		Subject:      "switch",
		AltSubject:   "sw",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("Create switch device", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"overwrite_hostname_from_switch": c.FlagSet.Bool("retrieve-hostname-from-switch", false, "(Flag) Retrieve the hostname from the equipment instead of configuration file."),
				"format":                         c.FlagSet.String("format", "json", "The input format. Supported values are 'json','yaml'. The default format is json."),
				"read_config_from_file":          c.FlagSet.String("raw-config", _nilDefaultStr, "(Required) Read  configuration from file in the format specified with --format."),
				"read_config_from_pipe":          c.FlagSet.Bool("pipe", false, "(Flag) If set, read  configuration from pipe instead of from a file. Either this flag or the --raw-config option must be used."),
				"return_id":                      c.FlagSet.Bool("return-id", false, "Will print the ID of the created object. Useful for automating tasks."),
			}
		},
		ExecuteFunc: switchCreateCmd,
		Endpoint:    DeveloperEndpoint,
	},
	{
		Description:  "Get a switch device.",
		Subject:      "switch",
		AltSubject:   "sw",
		Predicate:    "get",
		AltPredicate: "show",
		FlagSet:      flag.NewFlagSet("get a switch device", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"network_device_id_or_identifier_string": c.FlagSet.String("id", _nilDefaultStr, "(Required) Switch's id or identifier string. "),
				"show_credentials":                       c.FlagSet.Bool("show-credentials", false, "(Flag) If set returns the switch credentials"),
				"format":                                 c.FlagSet.String("format", "", "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
				"raw":                                    c.FlagSet.Bool("raw", false, "(Flag) When set the return will be a full dump of the object. This is useful when copying configurations. Only works with json and yaml formats."),
			}
		},
		ExecuteFunc: switchGetCmd,
		Endpoint:    DeveloperEndpoint,
	},
	{
		Description:  "Delete a switch.",
		Subject:      "switch",
		AltSubject:   "sw",
		Predicate:    "delete",
		AltPredicate: "rm",
		FlagSet:      flag.NewFlagSet("delete switch", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"network_device_id_or_identifier_string": c.FlagSet.String("id", _nilDefaultStr, "(Required) Switch's id or identifier string. "),
				"autoconfirm":                            c.FlagSet.Bool("autoconfirm", false, "If true it does not ask for confirmation anymore"),
			}
		},
		ExecuteFunc: switchDeleteCmd,
		Endpoint:    DeveloperEndpoint,
	},
}

func switchListCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	datacenterName := getStringParam(c.Arguments["datacenter_name"])
	switchType := getStringParam(c.Arguments["switch_type"])

	list, err := client.SwitchDevices(datacenterName, switchType)

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
			FieldName: "IDENTIFIER",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "DATACENTER",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "DRIVER",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "PROVISIONER",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "MGMT IP",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
	}

	showCredentials := false
	if c.Arguments["show_credentials"] != nil && *c.Arguments["show_credentials"].(*bool) {
		showCredentials = true

		schema = append(schema, tableformatter.SchemaField{
			FieldName: "MGMT_USER",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		})

		schema = append(schema, tableformatter.SchemaField{
			FieldName: "MGMT_PASS",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		})

	}

	data := [][]interface{}{}
	for _, s := range *list {

		credentialsUser := ""
		credentialsPass := ""

		if showCredentials {

			sw, err := client.SwitchDeviceGet(s.NetworkEquipmentID, showCredentials)

			if err != nil {
				return "", err
			}

			credentialsUser = fmt.Sprintf("%s", sw.NetworkEquipmentManagementUsername)
			credentialsPass = fmt.Sprintf("%s", sw.NetworkEquipmentManagementPassword)

		}
		data = append(data, []interface{}{
			s.NetworkEquipmentID,
			s.NetworkEquipmentIdentifierString,
			s.DatacenterName,
			s.NetworkEquipmentDriver,
			s.NetworkEquipmentProvisionerType,
			s.NetworkEquipmentManagementAddress,
			credentialsUser,
			credentialsPass,
		})

	}

	tableformatter.TableSorter(schema).OrderBy(
		schema[0].FieldName,
		schema[1].FieldName,
		schema[2].FieldName).Sort(data)

	table := tableformatter.Table{
		Data:   data,
		Schema: schema,
	}
	return table.RenderTable("Switches", "", getStringParam(c.Arguments["format"]))
}

func switchCreateCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	var obj metalcloud.SwitchDevice

	err := getRawObjectFromCommand(c, &obj)
	if err != nil {
		return "", err
	}

	ret, err := client.SwitchDeviceCreate(obj, getBoolParam(c.Arguments["overwrite_hostname_from_switch"]))
	if err != nil {
		return "", err
	}

	if getBoolParam(c.Arguments["return_id"]) {
		return fmt.Sprintf("%d", ret.NetworkEquipmentID), nil
	}

	return "", err
}

func switchGetCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	retSW, err := getSwitchFromCommandLine("id", c, client)
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
			FieldName: "IDENTIFIER",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "DATACENTER",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "DRIVER",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "PROVISIONER",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "MGMT IP",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
	}

	credentialsUser := ""
	credentialsPass := ""

	showCredentials := getBoolParam(c.Arguments["show_credentials"])

	if showCredentials {
		credentialsUser = fmt.Sprintf("%s", retSW.NetworkEquipmentManagementUsername)
		credentialsPass = fmt.Sprintf("%s", retSW.NetworkEquipmentManagementPassword)

		schema = append(schema, tableformatter.SchemaField{
			FieldName: "MGMT_USER",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		})

		schema = append(schema, tableformatter.SchemaField{
			FieldName: "MGMT_PASS",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		})

	}

	data := [][]interface{}{{
		retSW.NetworkEquipmentID,
		retSW.NetworkEquipmentIdentifierString,
		retSW.DatacenterName,
		retSW.NetworkEquipmentDriver,
		retSW.NetworkEquipmentProvisionerType,
		retSW.NetworkEquipmentManagementAddress,
		credentialsUser,
		credentialsPass,
	}}

	var sb strings.Builder

	format := getStringParam(c.Arguments["format"])

	if getBoolParam(c.Arguments["raw"]) {
		ret, err := tableformatter.RenderRawObject(*retSW, format, "NetworkEquipment")
		if err != nil {
			return "", err
		}
		sb.WriteString(ret)
	} else {
		table := tableformatter.Table{
			Data:   data,
			Schema: schema,
		}
		ret, err := table.RenderTransposedTable("switch device", "", format)
		if err != nil {
			return "", err
		}
		sb.WriteString(ret)
	}

	return sb.String(), nil
}

func switchDeleteCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	retSW, err := getSwitchFromCommandLine("id", c, client)
	if err != nil {
		return "", err
	}
	confirm := false

	if getBoolParam(c.Arguments["autoconfirm"]) {
		confirm = true
	} else {

		confirmationMessage := fmt.Sprintf("Deleting switch %s (%d).  Are you sure? Type \"yes\" to continue:",
			retSW.NetworkEquipmentIdentifierString,
			retSW.NetworkEquipmentID)

		//this is simply so that we don't output a text on the command line under go test
		if strings.HasSuffix(os.Args[0], ".test") {
			confirmationMessage = ""
		}

		confirm, err = requestConfirmation(confirmationMessage)
		if err != nil {
			return "", err
		}

	}

	if !confirm {
		return "", fmt.Errorf("Operation not confirmed. Aborting")
	}

	err = client.SwitchDeviceDelete(retSW.NetworkEquipmentID)

	return "", err
}

func getSwitchFromCommandLine(paramName string, c *Command, client metalcloud.MetalCloudClient) (*metalcloud.SwitchDevice, error) {
	m, err := getParam(c, "network_device_id_or_identifier_string", paramName)
	if err != nil {
		return nil, err
	}

	showCredentials := getBoolParam(c.Arguments["show_credentials"])

	var retSW *metalcloud.SwitchDevice

	id, label, isID := idOrLabel(m)

	if isID {
		retSW, err = client.SwitchDeviceGet(id, showCredentials)

	} else {
		retSW, err = client.SwitchDeviceGetByIdentifierString(label, showCredentials)
	}

	if err != nil {
		return nil, err
	}

	return retSW, nil
}
