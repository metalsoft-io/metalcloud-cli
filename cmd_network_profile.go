package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
	"github.com/metalsoft-io/tableformatter"
)

var networkProfilesCmds = []Command{
	{
		Description:  "Lists all network profiles.",
		Subject:      "network-profile",
		AltSubject:   "np",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list network_profile", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"datacenter": c.FlagSet.String("datacenter", GetDatacenter(), "(Required) Network profile datacenter"),
				"format":     c.FlagSet.String("format", "", "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
			}
		},
		ExecuteFunc: networkProfileListCmd,
	},
	{
		Description:  "Lists vlans of network profile.",
		Subject:      "network-profile",
		AltSubject:   "np",
		Predicate:    "vlans-list",
		AltPredicate: "vlans-ls",
		FlagSet:      flag.NewFlagSet("vlans-list network_profile", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"network_profile_id": c.FlagSet.Int("id", _nilDefaultInt, "(Required) Network profile's id."),
				"format":             c.FlagSet.String("format", "", "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
			}
		},
		ExecuteFunc: networkProfileVlansListCmd,
	},
	{
		Description:  "Get network profile details.",
		Subject:      "network-profile",
		AltSubject:   "np",
		Predicate:    "get",
		AltPredicate: "show",
		FlagSet:      flag.NewFlagSet("Get network profile details.", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"network_profile_id": c.FlagSet.Int("id", _nilDefaultInt, "(Required) Network profile's id."),
				"format":             c.FlagSet.String("format", _nilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
			}
		},
		ExecuteFunc: networkProfileGetCmd,
		Endpoint:    ExtendedEndpoint,
	},
	{
		Description:  "Create network profile.",
		Subject:      "network-profile",
		AltSubject:   "np",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("Create network profile", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"datacenter":            c.FlagSet.String("datacenter", GetDatacenter(), "(Required) Label of the datacenter. Also used as an ID."),
				"format":                c.FlagSet.String("format", "json", "The input format. Supported values are 'json','yaml'. The default format is json."),
				"read_config_from_file": c.FlagSet.String("raw-config", _nilDefaultStr, "(Required) Read  configuration from file in the format specified with --format."),
				"read_config_from_pipe": c.FlagSet.Bool("pipe", false, "(Flag) If set, read  configuration from pipe instead of from a file. Either this flag or the --raw-config option must be used."),
				"return_id":             c.FlagSet.Bool("return-id", false, "Will print the ID of the created object. Useful for automating tasks."),
			}
		},
		ExecuteFunc: networkProfileCreateCmd,
		Endpoint:    DeveloperEndpoint,
	},
	{
		Description:  "Delete a network profile.",
		Subject:      "network-profile",
		AltSubject:   "np",
		Predicate:    "delete",
		AltPredicate: "rm",
		FlagSet:      flag.NewFlagSet("delete network profile", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"network_profile_id": c.FlagSet.Int("id", _nilDefaultInt, "(Required) Network profile's id "),
				"autoconfirm":        c.FlagSet.Bool("autoconfirm", false, "If true it does not ask for confirmation anymore"),
			}
		},
		ExecuteFunc: networkProfileDeleteCmd,
		Endpoint:    DeveloperEndpoint,
	},
	{
		Description:  "Add a network profile to an instance array.",
		Subject:      "network-profile",
		AltSubject:   "np",
		Predicate:    "associate",
		AltPredicate: "assign",
		FlagSet:      flag.NewFlagSet("assign network profile to an instance array", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"network_profile_id": c.FlagSet.Int("id", _nilDefaultInt, "(Required) Network profile's id"),
				"network_id":         c.FlagSet.Int("net", _nilDefaultInt, "(Required) Network's id"),
				"instance_array_id":  c.FlagSet.Int("ia", _nilDefaultInt, "(Required) Instance array's id"),
			}
		},
		ExecuteFunc: networkProfileAssociateToInstanceArrayCmd,
		Endpoint:    DeveloperEndpoint,
	},
	{
		Description:  "Remove (unassign) profile to an instance array.",
		Subject:      "network-profile",
		AltSubject:   "np",
		Predicate:    "disassociate",
		AltPredicate: "unassign",
		FlagSet:      flag.NewFlagSet("disassociate network profile to an instance array", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"instance_array_id": c.FlagSet.String("ia", _nilDefaultStr, "(Required) Instance array's id"),
				"network_id":        c.FlagSet.String("net", _nilDefaultStr, "(Required) Network's id"),
			}
		},
		ExecuteFunc: networkProfileUnassociateToInstanceArrayCmd,
		Endpoint:    DeveloperEndpoint,
	},
}

func networkProfileListCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	datacenter := c.Arguments["datacenter"]
	npList, err := client.NetworkProfiles(*datacenter.(*string))
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
			FieldName: "LABEL",
			FieldType: tableformatter.TypeString,
			FieldSize: 30,
		},
		{
			FieldName: "NETWORK TYPE",
			FieldType: tableformatter.TypeString,
			FieldSize: 30,
		},
		{
			FieldName: "VLANs",
			FieldType: tableformatter.TypeInterface,
			FieldSize: 30,
		},
		{
			FieldName: "CREATED",
			FieldType: tableformatter.TypeString,
			FieldSize: 10,
		},
		{
			FieldName: "UPDATED",
			FieldType: tableformatter.TypeString,
			FieldSize: 10,
		},
	}

	data := [][]interface{}{}
	for _, np := range *npList {
		vlans := ""

		for index, vlan := range np.NetworkProfileVLANs {
			if index == 0 {
				vlans = strconv.Itoa(vlan.VlanID)
			} else {
				vlans = vlans + "," + strconv.Itoa(vlan.VlanID)
			}
		}

		data = append(data, []interface{}{
			np.NetworkProfileID,
			np.NetworkProfileLabel,
			np.NetworkType,
			vlans,
			np.NetworkProfileCreatedTimestamp,
			np.NetworkProfileUpdatedTimestamp,
		})
	}

	tableformatter.TableSorter(schema).OrderBy(schema[0].FieldName).Sort(data)

	table := tableformatter.Table{
		Data:   data,
		Schema: schema,
	}

	return table.RenderTable("Network Profiles", "", getStringParam(c.Arguments["format"]))
}

func networkProfileVlansListCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {
	id, ok := getIntParamOk(c.Arguments["network_profile_id"])
	if !ok {
		return "", fmt.Errorf("-id required")
	}

	retNP, err := client.NetworkProfileGet(id)
	if err != nil {
		return "", err
	}

	schemaConfiguration := []tableformatter.SchemaField{
		{
			FieldName: "VLAN",
			FieldType: tableformatter.TypeInt,
			FieldSize: 6,
		},
		{
			FieldName: "Port mode",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "External connections",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "Provision subnet gateways",
			FieldType: tableformatter.TypeBool,
			FieldSize: 6,
		},
	}

	dataConfiguration := [][]interface{}{}
	networkProfileVlans := retNP.NetworkProfileVLANs

	for _, vlan := range networkProfileVlans {

		externalConnectionIDs := vlan.ExternalConnectionIDs
		ecIds := ""
		for index, ecId := range externalConnectionIDs {

			retEC, err := client.ExternalConnectionGet(ecId)
			if err != nil {
				return "", err
			}

			if index == 0 {
				ecIds = retEC.ExternalConnectionLabel + " (#" + strconv.Itoa(ecId) + ")"
			} else {
				ecIds = ecIds + ", " + retEC.ExternalConnectionLabel + " (#" + strconv.Itoa(ecId) + ")"
			}
		}

		dataConfiguration = append(dataConfiguration, []interface{}{
			vlan.VlanID,
			vlan.PortMode,
			ecIds,
			vlan.ProvisionSubnetGateways,
		})
	}

	tableConfiguration := tableformatter.Table{
		Data:   dataConfiguration,
		Schema: schemaConfiguration,
	}

	retConfigTable, err := tableConfiguration.RenderTableFoldable("", "", getStringParam(c.Arguments["format"]), 0)
	if err != nil {
		return "", err
	}

	return retConfigTable, err
}

func networkProfileGetCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	id, ok := getIntParamOk(c.Arguments["network_profile_id"])
	if !ok {
		return "", fmt.Errorf("-id required")
	}

	retNP, err := client.NetworkProfileGet(id)
	if err != nil {
		return "", err
	}

	schema := []tableformatter.SchemaField{
		{
			FieldName: "ID",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "LABEL",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "DATACENTER",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
	}

	data := [][]interface{}{
		{
			"#" + strconv.Itoa(retNP.NetworkProfileID),
			retNP.NetworkProfileLabel,
			retNP.DatacenterName,
		},
	}

	table := tableformatter.Table{
		Data:   data,
		Schema: schema,
	}

	retOverviewTable, err := table.RenderTableFoldable("", "", getStringParam(c.Arguments["format"]), 0)
	if err != nil {
		return "", err
	}

	return retOverviewTable, err
}

func networkProfileCreateCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {
	datacenter := c.Arguments["datacenter"]

	readContentfromPipe := getBoolParam((c.Arguments["read_config_from_pipe"]))

	var err error
	content := []byte{}

	if readContentfromPipe {
		content, err = readInputFromPipe()
	} else {

		if configFilePath, ok := getStringParamOk(c.Arguments["read_config_from_file"]); ok {

			content, err = readInputFromFile(configFilePath)
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

	format := getStringParam(c.Arguments["format"])

	var npConf metalcloud.NetworkProfile
	switch format {
	case "json":
		err := json.Unmarshal(content, &npConf)
		if err != nil {
			return "", err
		}
	case "yaml":
		err := yaml.Unmarshal(content, &npConf)
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("input format \"%s\" not supported", format)
	}

	ret, err := client.NetworkProfileCreate(*datacenter.(*string), npConf)
	if err != nil {
		return "", err
	}

	if c.Arguments["return_id"] != nil && *c.Arguments["return_id"].(*bool) {
		return fmt.Sprintf("%s", ret.DatacenterName), nil
	}

	return "", err

}

func networkProfileDeleteCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {
	networkProfileId, ok := getIntParamOk(c.Arguments["network_profile_id"])
	if !ok {
		return "", fmt.Errorf("-id is required (network_profile_id (NPI) number returned by get network profile")
	}

	confirm := getBoolParam(c.Arguments["autoconfirm"])

	networkProfile, err := client.NetworkProfileGet(networkProfileId)
	if err != nil {
		return "", err
	}

	if !confirm {

		confirmationMessage := fmt.Sprintf("Deleting network profile %s (%d).  Are you sure? Type \"yes\" to continue:",
			networkProfile.NetworkProfileLabel, networkProfile.NetworkProfileID)

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

	err = client.NetworkProfileDelete(networkProfileId)

	return "", err
}

func networkProfileAssociateToInstanceArrayCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {
	id, ok := getIntParamOk(c.Arguments["network_profile_id"])
	if !ok {
		return "", fmt.Errorf("-id required")
	}

	net, ok := getIntParamOk(c.Arguments["network_id"])
	if !ok {
		return "", fmt.Errorf("-net required")
	}

	ia, ok := getIntParamOk(c.Arguments["instance_array_id"])
	if !ok {
		return "", fmt.Errorf("-net required")
	}

	_, err := client.InstanceArrayNetworkProfileSet(ia, net, id)
	if err != nil {
		return "", err
	}

	return "", nil
}

func networkProfileUnassociateToInstanceArrayCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	instance_array_id, ok := getStringParamOk(c.Arguments["instance_array_id"])
	if !ok {
		return "", fmt.Errorf("-net required")
	}

	ia, err := strconv.Atoi(instance_array_id)
	if err != nil {
		return "", err
	}

	network_id, ok := getStringParamOk(c.Arguments["network_id"])
	if !ok {
		return "", fmt.Errorf("-net required")
	}

	net, err := strconv.Atoi(network_id)
	if err != nil {
		return "", err
	}

	return "", client.InstanceArrayNetworkProfileClear(ia, net)
}
