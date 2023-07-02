package switchcontroller

import (
	"flag"
	"fmt"
	"strings"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
	"github.com/metalsoft-io/metalcloud-cli/internal/colors"
	"github.com/metalsoft-io/metalcloud-cli/internal/command"
	"github.com/metalsoft-io/metalcloud-cli/internal/configuration"
	"github.com/metalsoft-io/tableformatter"
)

var SwitchControllerCmds = []command.Command{
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
				"datacenter_name":  c.FlagSet.String("datacenter", "", "The optional parameter acts as a filter that restricts the returned results to switch devices located in the specified datacenter."),
				"show_credentials": c.FlagSet.Bool("show_credentials", false, colors.Green("(Flag)")+" If set returns the switch management credentials. (Slow for large queries)"),
				"no_color":         c.FlagSet.Bool("no_color", false, " Disable coloring."),
			}
		},
		ExecuteFunc: switchControllersListCmd,
		Endpoint:    configuration.DeveloperEndpoint,
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
				"network_controller_id_or_identifier_string": c.FlagSet.String("id", command.NilDefaultStr, colors.Red("(Required)")+" Switch id or identifier string. "),
				"show_credentials":                           c.FlagSet.Bool("show_credentials", false, colors.Green("(Flag)")+" If set returns the switch credentials"),
				"format":                                     c.FlagSet.String("format", "", "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
				"raw":                                        c.FlagSet.Bool("raw", false, colors.Green("(Flag)")+" When set the return will be a full dump of the object. This is useful when copying configurations. Only works with json and yaml formats."),
				"no_color":                                   c.FlagSet.Bool("no_color", false, " Disable coloring."),
			}
		},
		ExecuteFunc: switchControllerGetCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
	// 	{
	// 		Description:  "Create switch device.",
	// 		Subject:      "switch",
	// 		AltSubject:   "sw",
	// 		Predicate:    "create",
	// 		AltPredicate: "new",
	// 		FlagSet:      flag.NewFlagSet("Create switch device", flag.ExitOnError),
	// 		InitFunc: func(c *command.Command) {
	// 			c.Arguments = map[string]interface{}{
	// 				"overwrite_hostname_from_switch": c.FlagSet.Bool("retrieve-hostname-from-switch", false, colors.Green("(Flag)")+" Retrieve the hostname from the equipment instead of configuration file."),
	// 				"format":                         c.FlagSet.String("format", "json", "The input format. Supported values are 'json','yaml'. The default format is json."),
	// 				"read_config_from_file":          c.FlagSet.String("raw-config", command.NilDefaultStr, colors.Red("(Required)")+" Read  configuration from file in the format specified with --format."),
	// 				"read_config_from_pipe":          c.FlagSet.Bool("pipe", false, colors.Green("(Flag)")+" If set, read  configuration from pipe instead of from a file. Either this flag or the --raw-config option must be used."),
	// 				"return_id":                      c.FlagSet.Bool("return-id", false, "Will print the ID of the created object. Useful for automating tasks."),
	// 			}
	// 		},
	// 		ExecuteFunc: switchCreateCmd,
	// 		Endpoint:    configuration.DeveloperEndpoint,
	// 		Example: `
	// metalcloud-cli switch create --format yaml --raw-config switch.yml --return-id

	// #Example configurations:

	// Dell OS10 - Data leaf
	// ===========================================================================
	// identifierString: swd003
	// datacenterName: dc-prod
	// provisionerType: evpnvxlanl2
	// provisionerPosition: leaf
	// driver: os_10
	// managementUsername: metalsoftadmin
	// managementPassword: XXXXX
	// managementAddress: 172.16.2.23
	// managementPort: 22
	// managementProtocol: ssh
	// managementMACAddress: "00:00:00:00:00:00"
	// primaryWANIPv4SubnetPrefixSize: 22
	// primaryWANIPv6SubnetPoolID: 4
	// primarySANSubnetPrefixSize: 21
	// quarantineSubnetStart: 172.16.240.64
	// quarantineSubnetEnd: 172.16.240.127
	// quarantineSubnetPrefixSize: 26
	// quarantineSubnetGateway: 172.16.240.65
	// requiresOSInstall: false
	// isBorderDevice: false
	// isStorageSwitch: false
	// networkTypesAllowed:
	// - wan
	// - quarantine
	// ===========================================================================

	// Dell OS10 - Storage leaf
	// ===========================================================================
	// identifierString: sws001
	// datacenterName: dc-prod
	// provisionerType: evpnvxlanl2
	// provisionerPosition: leaf
	// driver: os_10
	// managementUsername: metalsoftadmin
	// managementPassword: XXXXX
	// managementAddress: 172.16.2.33
	// managementPort: 22
	// managementProtocol: ssh
	// managementMACAddress: "00:00:00:00:00:00"
	// primaryWANIPv4SubnetPrefixSize: 22
	// primaryWANIPv6SubnetPoolID: 4
	// primarySANSubnetPrefixSize: 21
	// quarantineSubnetStart: 172.16.240.0
	// quarantineSubnetEnd: 172.16.240.63
	// quarantineSubnetPrefixSize: 26
	// quarantineSubnetGateway: 172.16.240.1
	// requiresOSInstall: false
	// isBorderDevice: false
	// isStorageSwitch: true
	// networkTypesAllowed:
	// - san
	// ===========================================================================

	// Dell OS10 - Border leaf
	// ===========================================================================
	// identifierString: swb001
	// datacenterName: dc-prod
	// provisionerType: evpnvxlanl2
	// provisionerPosition: other
	// driver: os_10
	// managementUsername: metalsoftadmin
	// managementPassword: XXXXXXXX
	// managementAddress: 172.16.2.19
	// managementPort: 22
	// managementProtocol: ssh
	// managementMACAddress: "00:00:00:00:00:00"
	// primaryWANIPv4SubnetPrefixSize: 22
	// primaryWANIPv6SubnetPoolID: 4
	// primarySANSubnetPrefixSize: 21
	// quarantineSubnetStart: 172.16.240.0
	// quarantineSubnetEnd: 172.16.240.63
	// quarantineSubnetPrefixSize: 26
	// quarantineSubnetGateway: 172.16.240.1
	// requiresOSInstall: false
	// isBorderDevice: true
	// isStorageSwitch: false
	// networkTypesAllowed:
	// - wan
	// ===========================================================================

	// JunOS leaf
	// ===========================================================================
	// identifierString: juniper-virtual-chassis-QFX5100
	// datacenterName: dc-internal-test
	// provisionerType: vlan
	// provisionerPosition: leaf
	// driver: junos18
	// managementUsername: root
	// managementPassword: "xxxxx"
	// managementAddress: 10.0.5.10
	// managementPort: 22
	// managementProtocol: ssh
	// managementMACAddress: "00:00:00:00:00:00"
	// primaryWANIPv4SubnetPool: 192.168.253.0
	// primaryWANIPv4SubnetPrefixSize: 24
	// primaryWANIPv6SubnetPrefixSize: 48
	// primaryWANIPv6SubnetPool: fddf:d958:fb10:0000:0000:0000:0000:0000
	// primarySANSubnetPool: 100.64.0.1
	// primarySANSubnetPrefixSize: 21
	// quarantineSubnetStart: 192.168.254.0
	// quarantineSubnetEnd: 192.168.254.255
	// quarantineSubnetPrefixSize: 24
	// quarantineSubnetGateway: 192.168.254.1
	// requiresOSInstall: false
	// isBorderDevice: false
	// isStorageSwitch: false
	// networkTypesAllowed:
	// - wan
	// - lan
	// - san
	// - quarantine
	// ===========================================================================

	// HP VPLS TOR
	// ===========================================================================
	// identifierString: US_CHG_QTS01_01_MJ40_ML43_01
	// description: ToR switch
	// #the datacenter label
	// datacenter: dc-prod
	// provisionerType: vpls
	// provisionerPosition: tor
	// driver: hp5900
	// managementAddress: 10.0.2.1
	// managementProtocol: ssh
	// managementPort: 22
	// managementUsername: msprov1
	// managementPassword: XXXX
	// primaryWANIPv6SubnetPool: 2a02:cb80:1000:0000:0000:0000:0000:0000
	// primaryWANIPv6SubnetPrefixSize: 53
	// primarySANSubnetPool: 100.64.0.0
	// primarySANSubnetPrefixSize: 21
	// primaryWANIPv4SubnetPool: 172.24.0.0
	// primaryWANIPv4SubnetPrefixSize: 22
	// quarantineSubnetStart: 172.16.0.2
	// quarantineSubnetEnd: 172.16.0.254
	// quarantineSubnetPrefixSize: 24
	// quarantineSubnetGateway: 172.16.0.1
	// requiresOSInstall: false
	// volumeTemplateID: 0
	// ===========================================================================
	// 		`,
	// 	},
	// 	{
	// 		Description:  "Edit switch device.",
	// 		Subject:      "switch",
	// 		AltSubject:   "sw",
	// 		Predicate:    "edit",
	// 		AltPredicate: "update",
	// 		FlagSet:      flag.NewFlagSet("Edit switch device", flag.ExitOnError),
	// 		InitFunc: func(c *command.Command) {
	// 			c.Arguments = map[string]interface{}{
	// 				"network_device_id_or_identifier_string": c.FlagSet.String("id", command.NilDefaultStr, colors.Red("(Required)")+" Switch id or identifier string. "),
	// 				"overwrite_hostname_from_switch":         c.FlagSet.Bool("retrieve-hostname-from-switch", false, colors.Green("(Flag)")+" Retrieve the hostname from the equipment instead of configuration file."),
	// 				"format":                                 c.FlagSet.String("format", "json", "The input format. Supported values are 'json','yaml'. The default format is json."),
	// 				"read_config_from_file":                  c.FlagSet.String("raw-config", command.NilDefaultStr, colors.Red("(Required)")+" Read  configuration from file in the format specified with --format."),
	// 				"read_config_from_pipe":                  c.FlagSet.Bool("pipe", false, colors.Green("(Flag)")+" If set, read  configuration from pipe instead of from a file. Either this flag or the --raw-config option must be used."),
	// 				"return_id":                              c.FlagSet.Bool("return-id", false, "Will print the ID of the created object. Useful for automating tasks."),
	// 			}
	// 		},
	// 		ExecuteFunc: switchEditCmd,
	// 		Endpoint:    configuration.DeveloperEndpoint,
	// 	},
	// 	{
	// 		Description:  "Delete a switch.",
	// 		Subject:      "switch",
	// 		AltSubject:   "sw",
	// 		Predicate:    "delete",
	// 		AltPredicate: "rm",
	// 		FlagSet:      flag.NewFlagSet("delete switch", flag.ExitOnError),
	// 		InitFunc: func(c *command.Command) {
	// 			c.Arguments = map[string]interface{}{
	// 				"network_device_id_or_identifier_string": c.FlagSet.String("id", command.NilDefaultStr, colors.Red("(Required)")+" Switch id or identifier string. "),
	// 				"autoconfirm":                            c.FlagSet.Bool("autoconfirm", false, colors.Green("(Flag)")+" If set it will assume action is confirmed"),
	// 			}
	// 		},
	// 		ExecuteFunc: switchDeleteCmd,
	// 		Endpoint:    configuration.DeveloperEndpoint,
	// 	},
	// 	{
	// 		Description:  "Lists switch interfaces.",
	// 		Subject:      "switch",
	// 		AltSubject:   "sw",
	// 		Predicate:    "interfaces",
	// 		AltPredicate: "intf",
	// 		FlagSet:      flag.NewFlagSet("list switch interfaces", flag.ExitOnError),
	// 		InitFunc: func(c *command.Command) {
	// 			c.Arguments = map[string]interface{}{
	// 				"network_device_id_or_identifier_string": c.FlagSet.String("id", command.NilDefaultStr, colors.Red("(Required)")+" Switch id or identifier string. "),
	// 				"format":                                 c.FlagSet.String("format", command.NilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
	// 				"raw":                                    c.FlagSet.Bool("raw", false, colors.Green("(Flag)")+" When set the return will be a full dump of the object. This is useful when copying configurations. Only works with json and yaml formats."),
	// 			}
	// 		},
	// 		ExecuteFunc: switchInterfacesListCmd,
	// 		Endpoint:    configuration.DeveloperEndpoint,
	// 	},
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
		ret, err := tableformatter.RenderRawObject(*switchController, format, "NetworkEquipmentController")
		if err != nil {
			return "", err
		}
		sb.WriteString(ret)
	} else {
		table := tableformatter.Table{
			Data:   data,
			Schema: schema,
		}
		ret, err := table.RenderTransposedTable("switch device controller", "", format)
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
		fmt.Println("AICI")
		return nil, err
	}

	return switchController, nil
}
