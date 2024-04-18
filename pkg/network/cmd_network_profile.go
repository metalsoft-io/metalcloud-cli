package network

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
	"github.com/metalsoft-io/metalcloud-cli/internal/colors"
	"github.com/metalsoft-io/metalcloud-cli/internal/command"
	"github.com/metalsoft-io/metalcloud-cli/internal/configuration"
	"github.com/metalsoft-io/metalcloud-cli/internal/objects"
	"github.com/metalsoft-io/tableformatter"
)

var NetworkProfileCmds = []command.Command{
	{
		Description:  "Lists all network profiles.",
		Subject:      "network-profile",
		AltSubject:   "np",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list network_profile", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"datacenter": c.FlagSet.String("datacenter", command.NilDefaultStr, colors.Red("(Required)")+" Network profile datacenter"),
				"format":     c.FlagSet.String("format", command.NilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
			}
		},
		ExecuteFunc: networkProfileListCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
	{
		Description:  "Lists vlans of network profile.",
		Subject:      "network-profile",
		AltSubject:   "np",
		Predicate:    "vlan-list",
		AltPredicate: "vlans",
		FlagSet:      flag.NewFlagSet("vlan-list network_profile", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"network_profile_id": c.FlagSet.Int("id", command.NilDefaultInt, colors.Red("(Required)")+" Network profile's id."),
				"format":             c.FlagSet.String("format", command.NilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
			}
		},
		ExecuteFunc: networkProfileVlansListCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
	{
		Description:  "Get network profile details.",
		Subject:      "network-profile",
		AltSubject:   "np",
		Predicate:    "get",
		AltPredicate: "show",
		FlagSet:      flag.NewFlagSet("Get network profile details.", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"network_profile_id": c.FlagSet.Int("id", command.NilDefaultInt, colors.Red("(Required)")+" Network profile's id."),
				"format":             c.FlagSet.String("format", "yaml", "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
			}
		},
		ExecuteFunc: networkProfileGetCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
	{
		Description:  "Create network profile.",
		Subject:      "network-profile",
		AltSubject:   "np",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("Create network profile", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"read_config_from_file": c.FlagSet.String("f", command.NilDefaultStr, colors.Red("(Required)")+" Read configuration from file in the format specified with --format."),
				"format":                c.FlagSet.String("format", "json", "The input format. Supported values are 'json','yaml'. The default format is json."),
				"return_id":             c.FlagSet.Bool("return-id", false, "Will print the ID of the created object. Useful for automating tasks."),
			}
		},
		ExecuteFunc: networkProfileCreateCmd,
		Endpoint:    configuration.DeveloperEndpoint,
		Example: `
#create file network-profile.yaml:
kind: NetworkProfile
apiVersion: 1.0
label: my-network-profile
dc: my-datacenter
networkType: wan
vlans:
- vlanID: 3510
  portMode: trunk
  provisionSubnetGateways: false
  provisionVXLAN: false
  extConnectionIDs: []
  subnetPools: []
- vlanID: 3511
  portMode: trunk
  provisionSubnetGateways: false
  provisionVXLAN: false
  extConnectionIDs: []
  subnetPools: []
- vlanID: 3512
  portMode: trunk
  provisionSubnetGateways: false
  provisionVXLAN: false
  extConnectionIDs: []
  subnetPools: []
- vlanID: 3513
  portMode: trunk
  provisionSubnetGateways: false
  provisionVXLAN: false
  extConnectionIDs: []
  subnetPools: []
- vlanID: 3642
  portMode: native
  provisionSubnetGateways: false
  provisionVXLAN: false
  extConnectionIDs: []
  subnetPools: []

#create the actual profile from the file: 
metalcloud-cli network-profile create -f ./network-profile.yaml

More details available https://docs.metalsoft.io/en/latest/guides/adding_a_network_profile.html
`,
	},
	{
		Description:  "Delete a network profile.",
		Subject:      "network-profile",
		AltSubject:   "np",
		Predicate:    "delete",
		AltPredicate: "rm",
		FlagSet:      flag.NewFlagSet("delete network profile", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"network_profile_id": c.FlagSet.Int("id", command.NilDefaultInt, colors.Red("(Required)")+" Network profile's id "),
				"autoconfirm":        c.FlagSet.Bool("autoconfirm", false, colors.Green("(Flag)")+" If set it will assume action is confirmed"),
			}
		},
		ExecuteFunc: networkProfileDeleteCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
	{
		Description:  "Add a network profile to an instance array.",
		Subject:      "network-profile",
		AltSubject:   "np",
		Predicate:    "associate",
		AltPredicate: "assign",
		FlagSet:      flag.NewFlagSet("assign network profile to an instance array", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"network_profile_id": c.FlagSet.Int("id", command.NilDefaultInt, colors.Red("(Required)")+" Network profile's id"),
				"network_id":         c.FlagSet.Int("net", command.NilDefaultInt, colors.Red("(Required)")+" Network's id"),
				"instance_array_id":  c.FlagSet.Int("ia", command.NilDefaultInt, colors.Red("(Required)")+" Instance array's id"),
			}
		},
		ExecuteFunc: networkProfileAssociateToInstanceArrayCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
	{
		Description:  "Remove network profile from an instance array.",
		Subject:      "network-profile",
		AltSubject:   "np",
		Predicate:    "remove",
		AltPredicate: "rm",
		FlagSet:      flag.NewFlagSet("remove network profile of an instance array", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"instance_array_id": c.FlagSet.String("ia", command.NilDefaultStr, colors.Red("(Required)")+" Instance array's id"),
				"network_id":        c.FlagSet.String("net", command.NilDefaultStr, colors.Red("(Required)")+" Network's id"),
			}
		},
		ExecuteFunc: networkProfileRemoveFromInstanceArrayCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
}

func networkProfileListCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	datacenter, ok := command.GetStringParamOk(c.Arguments["datacenter"])
	if !ok {
		return "", fmt.Errorf("-datacenter is required")
	}

	npList, err := client.NetworkProfiles(datacenter)
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
			FieldName: "PUBLIC",
			FieldType: tableformatter.TypeBool,
			FieldSize: 10,
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

		for _, vlan := range np.NetworkProfileVLANs {
			if vlan.VlanID != nil {
				if vlans == "" {
					vlans = strconv.Itoa(*vlan.VlanID)

				} else {
					vlans = vlans + "," + strconv.Itoa(*vlan.VlanID)
				}
			}
		}

		data = append(data, []interface{}{
			np.NetworkProfileID,
			colors.Blue(np.NetworkProfileLabel),
			np.NetworkType,
			vlans,
			np.NetworkProfileIsPublic,
			np.NetworkProfileCreatedTimestamp,
			np.NetworkProfileUpdatedTimestamp,
		})
	}

	tableformatter.TableSorter(schema).OrderBy(schema[0].FieldName).Sort(data)

	table := tableformatter.Table{
		Data:   data,
		Schema: schema,
	}

	return table.RenderTable("Network Profiles", "", command.GetStringParam(c.Arguments["format"]))
}

func networkProfileVlansListCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	id, ok := command.GetIntParamOk(c.Arguments["network_profile_id"])
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
			FieldType: tableformatter.TypeString,
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

		vlanid := "auto"
		if vlan.VlanID != nil {
			vlanid = strconv.Itoa(*vlan.VlanID)
		}

		dataConfiguration = append(dataConfiguration, []interface{}{
			vlanid,
			vlan.PortMode,
			ecIds,
			vlan.ProvisionSubnetGateways,
		})
	}

	tableConfiguration := tableformatter.Table{
		Data:   dataConfiguration,
		Schema: schemaConfiguration,
	}

	retConfigTable, err := tableConfiguration.RenderTableFoldable("", "", command.GetStringParam(c.Arguments["format"]), 0)
	if err != nil {
		return "", err
	}

	return retConfigTable, err
}

func networkProfileGetCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	id, ok := command.GetIntParamOk(c.Arguments["network_profile_id"])
	if !ok {
		return "", fmt.Errorf("-id required")
	}

	retNP, err := client.NetworkProfileGet(id)
	if err != nil {
		return "", err
	}

	format := command.GetStringParam(c.Arguments["format"])
	ret, err := objects.RenderRawObject(*retNP, format, "NetworkProfile")
	if err != nil {
		return "", err
	}

	return ret, nil
}

func networkProfileCreateCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	obj, err := objects.ReadSingleObjectFromCommand(c, client)
	if err != nil {
		return "", err
	}
	np := (*obj).(metalcloud.NetworkProfile)

	createdNP, err := client.NetworkProfileCreate(np.DatacenterName, np)
	if err != nil {
		return "", err
	}

	if command.GetBoolParam(c.Arguments["return_id"]) {
		return fmt.Sprintf("%d", createdNP.NetworkProfileID), nil
	}

	return "", err
}

func networkProfileDeleteCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	networkProfileId, ok := command.GetIntParamOk(c.Arguments["network_profile_id"])
	if !ok {
		return "", fmt.Errorf("-id is required")
	}

	confirm := command.GetBoolParam(c.Arguments["autoconfirm"])

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

		confirm, err = command.RequestConfirmation(confirmationMessage)
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

func networkProfileAssociateToInstanceArrayCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	id, ok := command.GetIntParamOk(c.Arguments["network_profile_id"])
	if !ok {
		return "", fmt.Errorf("-id required")
	}

	net, ok := command.GetIntParamOk(c.Arguments["network_id"])
	if !ok {
		return "", fmt.Errorf("-net required")
	}

	ia, ok := command.GetIntParamOk(c.Arguments["instance_array_id"])
	if !ok {
		return "", fmt.Errorf("-net required")
	}

	_, err := client.InstanceArrayNetworkProfileSet(ia, net, id)
	if err != nil {
		return "", err
	}

	return "", nil
}

func networkProfileRemoveFromInstanceArrayCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	instance_array_id, ok := command.GetStringParamOk(c.Arguments["instance_array_id"])
	if !ok {
		return "", fmt.Errorf("-net required")
	}

	ia, err := strconv.Atoi(instance_array_id)
	if err != nil {
		return "", err
	}

	network_id, ok := command.GetStringParamOk(c.Arguments["network_id"])
	if !ok {
		return "", fmt.Errorf("-net required")
	}

	net, err := strconv.Atoi(network_id)
	if err != nil {
		return "", err
	}

	return "", client.InstanceArrayNetworkProfileClear(ia, net)
}
