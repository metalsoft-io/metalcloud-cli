package switchcontroller

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

var SwitchControllerCmds = []command.Command{
	{
		Description:  "Create SDN controller.",
		Subject:      "switch-controller",
		AltSubject:   "sw-ctrl",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("Create switch device", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"format":                c.FlagSet.String("format", "json", "The input format. Supported values are 'json','yaml'. The default format is json."),
				"read_config_from_file": c.FlagSet.String("f", command.NilDefaultStr, colors.Red("(Required)")+" Read  configuration from file in the format specified with --format."),
				"read_config_from_pipe": c.FlagSet.Bool("pipe", false, colors.Green("(Flag)")+" If set, read  configuration from pipe instead of from a file. Either this flag or the --raw-config option must be used."),
				"return_id":             c.FlagSet.Bool("return-id", false, "Will print the ID of the created object. Useful for automating tasks."),
			}
		},
		ExecuteFunc: switchControllerCreateCmd,
		Endpoint:    configuration.DeveloperEndpoint,
		Example: `
metalcloud-cli switch-controller create --format yaml --raw-config switch-controller.yaml --return-id

switch-controller.yaml:

identifierString: Cisco ACI 5.1
description: Cisco ACI 5.1 controller
datacenterName: test-aci
provisionerType: sdn
provisionerPosition: leaf
driver: cisco_aci51
managementAddress: 10.255.239.150
managementProtocol: API
managementPort: 22
managementUsername: admin
managementPassword: hello123
`,
	},
	{
		Description:  "Lists registered switch controllers.",
		Subject:      "switch-controller",
		AltSubject:   "sw-ctrl",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list switch controllers", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"format":           c.FlagSet.String("format", command.NilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
				"datacenter_name":  c.FlagSet.String("datacenter", command.NilDefaultStr, "The optional parameter acts as a filter that restricts the returned results to switch devices located in the specified datacenter."),
				"show_credentials": c.FlagSet.Bool("show-credentials", false, colors.Green("(Flag)")+" If set returns the switch management credentials. (Slow for large queries)"),
			}
		},
		ExecuteFunc: switchControllersListCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
	{
		Description:  "Edit switch controller configuration",
		Subject:      "switch-controller",
		AltSubject:   "sw-ctrl",
		Predicate:    "update",
		AltPredicate: "edit",
		FlagSet:      flag.NewFlagSet("Update switch controller configuration", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"network_controller_id_or_identifier_string": c.FlagSet.String("id", command.NilDefaultStr, colors.Red("(Required)")+" Switch controller id or identifier string. "),
				"format":                c.FlagSet.String("format", "json", "The input format. Supported values are 'json','yaml'. The default format is json."),
				"read_config_from_file": c.FlagSet.String("f", command.NilDefaultStr, colors.Red("(Required)")+" Read configuration from file in the format specified with --format."),
				"read_config_from_pipe": c.FlagSet.Bool("pipe", false, colors.Green("(Flag)")+" If set, read  configuration from pipe instead of from a file. Either this flag or the --raw-config option must be used."),
				"return_id":             c.FlagSet.Bool("return-id", false, "Will print the ID of the created object. Useful for automating tasks."),
			}
		},
		ExecuteFunc: switchControllerEditCmd,
		Endpoint:    configuration.DeveloperEndpoint,
		Example: `
metalcloud-cli switch-controller update --id 18 --raw-config update_sw_ctrl.yaml --format yaml

update_sw_ctrl.yaml:

options:
 vrf_shared_name: test1234
fabricConfiguration:
 network_equipment_primary_wan_ipv4_subnet_pool: 192.168.0.0
 network_equipment_primary_wan_ipv4_subnet_prefix_size: 22
 network_equipment_primary_wan_ipv6_subnet_prefix_size: 53
 network_equipment_primary_san_subnet_pool: 192.168.0.0
 network_equipment_primary_san_subnet_prefix_size: 21
 network_equipment_primary_wan_ipv6_subnet_pool: fd1f:8bbb:56b3:800:0:0:0:0  
 network_equipment_description: test
 network_equipment_country: UK
 network_equipment_city: Reading
 `,
	},
	{
		Description:  "Get configuration for a controller.",
		Subject:      "switch-controller",
		AltSubject:   "sw-ctrl",
		Predicate:    "get",
		AltPredicate: "show",
		FlagSet:      flag.NewFlagSet("get a switch controller configuration", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"network_controller_id_or_identifier_string": c.FlagSet.String("id", command.NilDefaultStr, colors.Red("(Required)")+" Switch controller id or identifier string. "),
				"show_credentials":                           c.FlagSet.Bool("show-credentials", false, colors.Green("(Flag)")+" If set returns the switch controller credentials"),
				"format":                                     c.FlagSet.String("format", command.NilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
				"raw":                                        c.FlagSet.Bool("raw", false, colors.Green("(Flag)")+" When set the return will be a full dump of the object. This is useful when copying configurations. Only works with json and yaml formats."),
			}
		},
		ExecuteFunc: switchControllerGetCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
	{
		Description:  "Creates multiple network equipment controller records, based on the fabric configuration of the switch controller.",
		Subject:      "switch-controller",
		AltSubject:   "sw-ctrl",
		Predicate:    "sync",
		AltPredicate: "sync",
		FlagSet:      flag.NewFlagSet("sync switch controller", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"network_controller_id_or_identifier_string": c.FlagSet.String("id", command.NilDefaultStr, colors.Red("(Required)")+" Switch controller id or identifier string. "),
			}
		},
		ExecuteFunc: switchControllerSyncCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
	{
		Description:  "Delete a switch controller.",
		Subject:      "switch-controller",
		AltSubject:   "sw-ctrl",
		Predicate:    "delete",
		AltPredicate: "rm",
		FlagSet:      flag.NewFlagSet("delete switch controller", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"network_controller_id_or_identifier_string": c.FlagSet.String("id", command.NilDefaultStr, colors.Red("(Required)")+" Switch controller id or identifier string. "),
				"autoconfirm": c.FlagSet.Bool("autoconfirm", false, colors.Green("(Flag)")+" If set it will assume action is confirmed"),
			}
		},
		ExecuteFunc: switchControllerDeleteCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
	{
		Description:  "Lists switches managed by a controller.",
		Subject:      "switch-controller",
		AltSubject:   "sw-ctrl",
		Predicate:    "switches",
		AltPredicate: "sw-list",
		FlagSet:      flag.NewFlagSet("list switches managed by a controller", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"format": c.FlagSet.String("format", command.NilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
				"network_controller_id_or_identifier_string": c.FlagSet.String("id", command.NilDefaultStr, colors.Red("(Required)")+" Switch controller id or identifier string. "),
				"raw": c.FlagSet.Bool("raw", false, colors.Green("(Flag)")+" When set the return will be a full dump of the object. This is useful when copying configurations. Only works with json and yaml formats."),
			}
		},
		ExecuteFunc: switchControllerSwitchesListCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
}

func switchControllerCreateCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	var obj metalcloud.SwitchDeviceController

	err := command.GetRawObjectFromCommand(c, &obj)
	if err != nil {
		return "", err
	}

	if obj.DatacenterName == "" {
		return "", fmt.Errorf("datacenter name is required.")
	}

	swCtrl, err := client.SwitchDeviceControllerCreate(obj)
	if err != nil {
		return "", err
	}

	if command.GetBoolParam(c.Arguments["return_id"]) {
		return fmt.Sprintf("%d", swCtrl.NetworkEquipmentControllerID), nil
	}

	return "", err
}

func switchControllersListCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	datacenterName := command.GetStringParam(c.Arguments["datacenter_name"])

	list, err := client.SwitchDeviceControllers(datacenterName)

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
	for _, switchController := range *list {
		credentialsUser := ""
		credentialsPass := ""

		if showCredentials {
			swCtrl, err := client.SwitchDeviceControllerGet(switchController.NetworkEquipmentControllerID, showCredentials)

			if err != nil {
				return "", err
			}

			credentialsUser = fmt.Sprintf("%s", swCtrl.NetworkEquipmentControllerManagementUsername)
			credentialsPass = fmt.Sprintf("%s", swCtrl.NetworkEquipmentControllerManagementPassword)

		}
		data = append(data, []interface{}{
			switchController.NetworkEquipmentControllerID,
			switchController.NetworkEquipmentControllerIdentifierString,
			switchController.DatacenterName,
			switchController.NetworkEquipmentControllerDriver,
			switchController.NetworkEquipmentControllerProvisionerType,
			switchController.NetworkEquipmentControllerManagementAddress,
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
	return table.RenderTable("Switch Controllers", "", command.GetStringParam(c.Arguments["format"]))
}

func switchControllerEditCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	var obj metalcloud.SwitchDeviceController

	err := command.GetRawObjectFromCommand(c, &obj)
	if err != nil {
		return "", err
	}

	retSwCtrl, err := getSwitchControllerFromCommandLine("id", c, client)
	if err != nil {
		return "", err
	}

	networkEquipmentControllerData := map[string]interface{}{
		"datacenter_name":                                   retSwCtrl.DatacenterName,
		"network_equipment_controller_options":              obj.NetworkEquipmentControllerOptions,
		"network_equipment_controller_fabric_configuration": obj.NetworkEquipmentControllerFabricConfiguration,
	}

	updatedSwCtrl, err := client.SwitchDeviceControllerUpdate(retSwCtrl.NetworkEquipmentControllerID, networkEquipmentControllerData)
	if err != nil {
		return "", err
	}

	if command.GetBoolParam(c.Arguments["return_id"]) {
		return fmt.Sprintf("%d", updatedSwCtrl.NetworkEquipmentControllerID), nil
	}

	return "", err
}

func switchControllerGetCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	switchController, err := getSwitchControllerFromCommandLine("id", c, client)
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
			FieldName: "HOSTNAME",
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

	showCredentials := command.GetBoolParam(c.Arguments["show_credentials"])

	if showCredentials {
		credentialsUser = fmt.Sprintf("%s", switchController.NetworkEquipmentControllerManagementUsername)
		credentialsPass = fmt.Sprintf("%s", switchController.NetworkEquipmentControllerManagementPassword)

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
		switchController.NetworkEquipmentControllerID,
		switchController.NetworkEquipmentControllerIdentifierString,
		switchController.DatacenterName,
		switchController.NetworkEquipmentControllerDriver,
		switchController.NetworkEquipmentControllerProvisionerType,
		switchController.NetworkEquipmentControllerManagementAddress,
		credentialsUser,
		credentialsPass,
	}}

	var sb strings.Builder

	format := command.GetStringParam(c.Arguments["format"])

	if command.GetBoolParam(c.Arguments["raw"]) {
		ret, err := objects.RenderRawObject(*switchController, format, "NetworkEquipmentController")
		if err != nil {
			return "", err
		}
		sb.WriteString(ret)
	} else {
		table := tableformatter.Table{
			Data:   data,
			Schema: schema,
		}
		ret, err := table.RenderTransposedTable("properties", "", format)
		if err != nil {
			return "", err
		}
		sb.WriteString(ret)
	}

	return sb.String(), nil
}

func getSwitchControllerFromCommandLine(paramName string, c *command.Command, client metalcloud.MetalCloudClient) (*metalcloud.SwitchDeviceController, error) {
	return getSwitchControllerFromCommandLineWithPrivateParam("network_controller_id_or_identifier_string", paramName, c, client)
}

func getSwitchControllerFromCommandLineWithPrivateParam(private_paramName string, public_paramName string, c *command.Command, client metalcloud.MetalCloudClient) (*metalcloud.SwitchDeviceController, error) {
	m, err := command.GetParam(c, private_paramName, public_paramName)
	if err != nil {
		return nil, err
	}

	showCredentials := command.GetBoolParam(c.Arguments["show_credentials"])

	var switchController *metalcloud.SwitchDeviceController

	id, label, isID := command.IdOrLabel(m)

	if isID {
		switchController, err = client.SwitchDeviceControllerGet(id, showCredentials)
	} else {
		switchController, err = client.SwitchDeviceControllerGetByIdentifierString(label, showCredentials)
	}

	if err != nil {
		return nil, err
	}

	return switchController, nil
}

func switchControllerSyncCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	retSWCtrl, err := getSwitchControllerFromCommandLine("id", c, client)
	if err != nil {
		return "", err
	}

	_, err = client.SwitchDeviceControllerSync(retSWCtrl.NetworkEquipmentControllerID)
	return "", err
}

func switchControllerDeleteCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	retSWCtrl, err := getSwitchControllerFromCommandLine("id", c, client)
	if err != nil {
		return "", err
	}
	confirm := false

	if command.GetBoolParam(c.Arguments["autoconfirm"]) {
		confirm = true
	} else {
		confirmationMessage := fmt.Sprintf("Deleting switch controller %s (%d).  Are you sure? Type \"yes\" to continue:",
			retSWCtrl.NetworkEquipmentControllerIdentifierString,
			retSWCtrl.NetworkEquipmentControllerID)

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

	err = client.SwitchDeviceControllerDelete(retSWCtrl.NetworkEquipmentControllerID)

	return "", err
}

func switchControllerSwitchesListCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	switchControllerID := command.GetStringParam(c.Arguments["network_controller_id_or_identifier_string"])

	if switchControllerID == "" {
		return "", fmt.Errorf("id must be specified")
	}

	controllerSwitchDevices, err := client.SwitchDeviceControllerSwitches(switchControllerID)

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
	for _, s := range controllerSwitchDevices {

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
	return table.RenderTable("Switches", "", command.GetStringParam(c.Arguments["format"]))
}
