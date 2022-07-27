package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
	"github.com/metalsoft-io/tableformatter"
)

var switchPairCmds = []Command{

	{
		Description:  "Lists switch pairs.",
		Subject:      "switch-pair",
		AltSubject:   "sw-pair",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list switch pairs", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"format": c.FlagSet.String("format", _nilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
			}
		},
		ExecuteFunc: switchPairListCmd,
		Endpoint:    DeveloperEndpoint,
	},
	{
		Description:  "Create switch pair.",
		Subject:      "switch-pair",
		AltSubject:   "sw-pair",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("Create a switch pair", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"network_device_id_or_identifier_string1": c.FlagSet.String("switch1", _nilDefaultStr, red("(Required)")+" First Switch's id or identifier string. "),
				"network_device_id_or_identifier_string2": c.FlagSet.String("switch2", _nilDefaultStr, red("(Required)")+" Second Switch's id or identifier string. "),
				"type":      c.FlagSet.String("type", "mlag", "The type of link. The default and only link type supported is `mlag`"),
				"return_id": c.FlagSet.Bool("return-id", false, "Will print the ID of the created object. Useful for automating tasks."),
			}
		},
		ExecuteFunc: switchPairCreateCmd,
		Endpoint:    DeveloperEndpoint,
	},
	{
		Description:  "Delete a switch pair.",
		Subject:      "switch-pair",
		AltSubject:   "sw-pair",
		Predicate:    "delete",
		AltPredicate: "rm",
		FlagSet:      flag.NewFlagSet("Delete switch pair", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"network_device_id_or_identifier_string1": c.FlagSet.String("switch1", _nilDefaultStr, red("(Required)")+" First Switch's id or identifier string. "),
				"network_device_id_or_identifier_string2": c.FlagSet.String("switch2", _nilDefaultStr, red("(Required)")+" Second Switch's id or identifier string. "),
				"type":        c.FlagSet.String("type", "mlag", "The type of link. The default and only link type supported is `mlag`"),
				"autoconfirm": c.FlagSet.Bool("autoconfirm", false, green("(Flag)")+" If set it will assume action is confirmed"),
			}
		},
		ExecuteFunc: switchPairDeleteCmd,
		Endpoint:    DeveloperEndpoint,
	},
}

func switchPairListCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	list, err := client.SwitchDeviceLinks()

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
			FieldName: "Switch1",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "Switch2",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "TYPE",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
	}

	data := [][]interface{}{}

	for _, s := range *list {

		sw1, err := client.SwitchDeviceGet(s.NetworkEquipmentID1, false)
		if err != nil {
			return "", err
		}

		sw2, err := client.SwitchDeviceGet(s.NetworkEquipmentID2, false)
		if err != nil {
			return "", err
		}

		data = append(data, []interface{}{
			s.NetworkEquipmentLinkID,
			fmt.Sprintf("%s (#%d)", sw1.NetworkEquipmentIdentifierString, sw1.NetworkEquipmentID),
			fmt.Sprintf("%s (#%d)", sw2.NetworkEquipmentIdentifierString, sw2.NetworkEquipmentID),
			s.NetworkEquipmentLinkType,
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
	return table.RenderTable("Switch links", "", getStringParam(c.Arguments["format"]))

}

func switchPairCreateCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	sw1, err := getSwitchFromCommandLineWithPrivateParam("network_device_id_or_identifier_string1", "switch1", c, client)
	if err != nil {
		return "", err
	}

	sw2, err := getSwitchFromCommandLineWithPrivateParam("network_device_id_or_identifier_string2", "switch2", c, client)
	if err != nil {
		return "", err
	}

	t := getStringParam(c.Arguments["type"])

	ret, err := client.SwitchDeviceLinkCreate(sw1.NetworkEquipmentID, sw2.NetworkEquipmentID, t)
	if err != nil {
		return "", err
	}

	if getBoolParam(c.Arguments["return_id"]) {
		return fmt.Sprintf("%d", ret.NetworkEquipmentLinkID), nil
	}

	return "", err
}

func switchPairDeleteCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	sw1, err := getSwitchFromCommandLineWithPrivateParam("network_device_id_or_identifier_string1", "switch1", c, client)
	if err != nil {
		return "", err
	}

	sw2, err := getSwitchFromCommandLineWithPrivateParam("network_device_id_or_identifier_string2", "switch2", c, client)
	if err != nil {
		return "", err
	}

	t := getStringParam(c.Arguments["type"])

	_, err = client.SwitchDeviceLinkGet(sw1.NetworkEquipmentID, sw2.NetworkEquipmentID, t)
	if err != nil {
		return "", err
	}

	confirm := false

	if getBoolParam(c.Arguments["autoconfirm"]) {
		confirm = true
	} else {

		confirmationMessage := fmt.Sprintf("Deleting switch pair %s - %s.  Are you sure? Type \"yes\" to continue:",
			sw1.NetworkEquipmentIdentifierString,
			sw2.NetworkEquipmentIdentifierString)

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

	err = client.SwitchDeviceLinkDelete(sw1.NetworkEquipmentID, sw2.NetworkEquipmentID, t)

	return "", err
}
