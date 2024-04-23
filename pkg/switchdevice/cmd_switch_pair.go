package switchdevice

import (
	"flag"
	"fmt"
	"os"
	"strings"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
	"github.com/metalsoft-io/metalcloud-cli/internal/colors"
	"github.com/metalsoft-io/metalcloud-cli/internal/command"
	"github.com/metalsoft-io/metalcloud-cli/internal/configuration"
	"github.com/metalsoft-io/metalcloud-cli/internal/objects"
	"github.com/metalsoft-io/tableformatter"
)

var SwitchPairCmds = []command.Command{
	{
		Description:  "Create switch pair.",
		Subject:      "switch-pair",
		AltSubject:   "sw-pair",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("Create a switch pair", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"read_config_from_file": c.FlagSet.String("f", command.NilDefaultStr, colors.Red("(Required)")+" Read configuration from file in the format specified with --format."),
				"format":                c.FlagSet.String("format", "yaml", "The input format. Supported values are 'json','yaml'. The default format is json."),
				"return_id":             c.FlagSet.Bool("return-id", false, "Will print the ID of the created object. Useful for automating tasks."),
			}
		},
		ExecuteFunc: switchPairCreateCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
	{
		Description:  "Lists switch pairs.",
		Subject:      "switch-pair",
		AltSubject:   "sw-pair",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list switch pairs", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"format": c.FlagSet.String("format", command.NilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
			}
		},
		ExecuteFunc: switchPairListCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
	{
		Description:  "List a switch pair.",
		Subject:      "switch-pair",
		AltSubject:   "sw-pair",
		Predicate:    "get",
		AltPredicate: "show",
		FlagSet:      flag.NewFlagSet("list switch pair", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"network_device_id1": c.FlagSet.Int("switch1", command.NilDefaultInt, colors.Red("(Required)")+" First Switch's id or identifier string. "),
				"network_device_id2": c.FlagSet.Int("switch2", command.NilDefaultInt, colors.Red("(Required)")+" Second Switch's id or identifier string. "),
				"type":               c.FlagSet.String("type", "mlag", "The type of link. The default and only link type supported is `mlag`"),
				"format":             c.FlagSet.String("format", "yaml", "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
			}
		},
		ExecuteFunc: switchPairGetCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
	{
		Description:  "Update switch pair.",
		Subject:      "switch-pair",
		AltSubject:   "sw-pair",
		Predicate:    "update",
		AltPredicate: "edit",
		FlagSet:      flag.NewFlagSet("Update a switch pair", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"read_config_from_file": c.FlagSet.String("f", command.NilDefaultStr, colors.Red("(Required)")+" Read configuration from file in the format specified with --format."),
				"format":                c.FlagSet.String("format", "yaml", "The input format. Supported values are 'json','yaml'. The default format is json."),
			}
		},
		ExecuteFunc: switchPairUpdateCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
	{
		Description:  "Delete a switch pair.",
		Subject:      "switch-pair",
		AltSubject:   "sw-pair",
		Predicate:    "delete",
		AltPredicate: "rm",
		FlagSet:      flag.NewFlagSet("Delete switch pair", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"network_device_id_or_identifier_string1": c.FlagSet.String("switch1", command.NilDefaultStr, colors.Red("(Required)")+" First Switch's id or identifier string. "),
				"network_device_id_or_identifier_string2": c.FlagSet.String("switch2", command.NilDefaultStr, colors.Red("(Required)")+" Second Switch's id or identifier string. "),
				"type":        c.FlagSet.String("type", "mlag", "The type of link. The default and only link type supported is `mlag`"),
				"autoconfirm": c.FlagSet.Bool("autoconfirm", false, colors.Green("(Flag)")+" If set it will assume action is confirmed"),
			}
		},
		ExecuteFunc: switchPairDeleteCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
}

func switchPairCreateCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	obj, err := objects.ReadSingleObjectFromCommand(c, client)
	if err != nil {
		return "", err
	}
	da := (*obj).(metalcloud.SwitchDeviceLink)

	createdSDL, err := client.SwitchDeviceLinkCreate(da.NetworkEquipmentID1, da.NetworkEquipmentID2, da.NetworkEquipmentLinkType)
	if err != nil {
		return "", err
	}

	if command.GetBoolParam(c.Arguments["return_id"]) {
		return fmt.Sprintf("%d", createdSDL.NetworkEquipmentLinkID), nil
	}

	return "", err
}

func switchPairListCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
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
	return table.RenderTable("Switch links", "", command.GetStringParam(c.Arguments["format"]))
}

func switchPairGetCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	swID1, ok := command.GetIntParamOk(c.Arguments["network_device_id1"])
	if !ok {
		return "", fmt.Errorf("-switch1 required")
	}

	swID2, ok := command.GetIntParamOk(c.Arguments["network_device_id2"])
	if !ok {
		return "", fmt.Errorf("-switch2 required")
	}

	t, ok := command.GetStringParamOk(c.Arguments["type"])
	if !ok {
		return "", fmt.Errorf("-type required")
	}

	sdl, err := client.SwitchDeviceLinkGet(swID1, swID2, t)
	if err != nil {
		return "", err
	}

	format := command.GetStringParam(c.Arguments["format"])
	ret, err := objects.RenderRawObject(*sdl, format, "SwitchDeviceLink")
	if err != nil {
		return "", err
	}

	return ret, nil
}

func switchPairUpdateCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	obj, err := objects.ReadSingleObjectFromCommand(c, client)
	if err != nil {
		return "", err
	}
	da := (*obj).(metalcloud.SwitchDeviceLink)

	_, err = client.SwitchDeviceLinkUpdate(da.NetworkEquipmentID1, da.NetworkEquipmentID2, da.NetworkEquipmentLinkType)
	if err != nil {
		return "", err
	}

	return "", err
}

func switchPairDeleteCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	sw1, err := getSwitchFromCommandLineWithPrivateParam("network_device_id_or_identifier_string1", "switch1", c, client)
	if err != nil {
		return "", err
	}

	sw2, err := getSwitchFromCommandLineWithPrivateParam("network_device_id_or_identifier_string2", "switch2", c, client)
	if err != nil {
		return "", err
	}

	t := command.GetStringParam(c.Arguments["type"])

	_, err = client.SwitchDeviceLinkGet(sw1.NetworkEquipmentID, sw2.NetworkEquipmentID, t)
	if err != nil {
		return "", err
	}

	confirm := false

	if command.GetBoolParam(c.Arguments["autoconfirm"]) {
		confirm = true
	} else {

		confirmationMessage := fmt.Sprintf("Deleting switch pair %s - %s.  Are you sure? Type \"yes\" to continue:",
			sw1.NetworkEquipmentIdentifierString,
			sw2.NetworkEquipmentIdentifierString)

		//this is simply so that we don't output a text on the command line under go test
		if strings.HasSuffix(os.Args[0], ".test") {
			confirmationMessage = ""
		}

		confirm, err = command.RequestConfirmation(confirmationMessage)
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
