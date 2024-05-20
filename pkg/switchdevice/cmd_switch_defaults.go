package switchdevice

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
	"github.com/metalsoft-io/metalcloud-cli/internal/colors"
	"github.com/metalsoft-io/metalcloud-cli/internal/command"
	"github.com/metalsoft-io/metalcloud-cli/internal/configuration"
	"github.com/metalsoft-io/tableformatter"
)

var SwitchDefaultsCmds = []command.Command{
	{
		Description:  "Lists switch defaults.",
		Subject:      "switch-defaults",
		AltSubject:   "sw-defaults",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list switch defaults", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"datacenter_name": c.FlagSet.String("datacenter", command.NilDefaultStr, colors.Red("(Required)")+"The datacenter name for which to retrieve the switch defaults."),
				"format":          c.FlagSet.String("format", command.NilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
				"raw":             c.FlagSet.Bool("raw", false, colors.Green("(Flag)")+" When set the return will be a full dump of the object. This is useful when copying configurations. Only works with json and yaml formats."),
			}
		},
		ExecuteFunc: switchDefaultsListCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
	{
		Description:  "Create switch defaults.",
		Subject:      "switch-defaults",
		AltSubject:   "sw-defaults",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("Create switch defaults", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"read_config_from_file": c.FlagSet.String("raw-config", command.NilDefaultStr, colors.Red("(Required)")+" Read configuration from a yaml file."),
			}
		},
		ExecuteFunc: switchDefaultsCreateCmd,
		Endpoint:    configuration.DeveloperEndpoint,
		Example: `
This command supports one or more switch defaults records in a yaml file format. The file format uses "---" separator between records. Here are 3 examples. The first one contains only the required properties.

datacenterName: test
serialNumber: ASDVDSF43GFD3221
managementMacAddress: 00:a0:c9:14:c8:33
---
datacenterName: test
serialNumber: ASDVDSF43GFD
managementMacAddress: 00:a0:c9:14:c8:31
asn: 65201
identifierString: test_switch1
isPartOfMlagPair: false
loopbackAddressIpv4: 99.107.105.233
loopbackAddressIpv6: 01ba:dff2:f830:c06a:bd1c:cfc3:8cbc:9635
partOfMlagPair: true
mlagDomainId: 5
mlagPartnerHostname: example-hostname
mlagPartnerVlanId: 40
mlagPeerLinkPortChannelId: 3
mlagSystemMac: 0B:B5:56:6E:A9:9D
position: leaf
skipInitialConfiguration: true
volumeTemplateID: 1
vtepAddressIpv4: 216.161.109.102
vtepAddressIpv6: 39c3:4824:d65b:e028:e86f:e4dc:97de:9456
customVariables:
  key1: value1
  key2: value2
---
datacenterName: test
serialNumber: ASDVDSF43GFD322
managementMacAddress: 00:a0:c9:14:c8:32
asn: 65202
identifierString: test_switch2
isPartOfMlagPair: false
loopbackAddressIpv4: 99.107.105.233
loopbackAddressIpv6: 01ba:dff2:f830:c06a:bd1c:cfc3:8cbc:9635
partOfMlagPair: true
mlagDomainId: 5
mlagPartnerHostname: example-hostname2
mlagPartnerVlanId: 40
mlagPeerLinkPortChannelId: 3
mlagSystemMac: 0B:B5:56:6E:A9:92
position: leaf
skipInitialConfiguration: true
volumeTemplateID: 1
vtepAddressIpv4: 216.161.109.102
vtepAddressIpv6: 39c3:4824:d65b:e028:e86f:e4dc:97de:9456
customVariables:
  key1: value1
  key2: value2
  key3: value3
`,
	},
	{
		Description:  "Delete switch defaults.",
		Subject:      "switch-defaults",
		AltSubject:   "sw-defaults",
		Predicate:    "delete",
		AltPredicate: "rm",
		FlagSet:      flag.NewFlagSet("Delete switch defaults", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"switch_defaults_ids": c.FlagSet.String("ids", command.NilDefaultStr, colors.Red("(Required)")+" Comma separated list of switch defaults ids to be removed."),
			}
		},
		ExecuteFunc: switchDefaultsDeleteCmd,
		Endpoint:    configuration.DeveloperEndpoint,
		Example: `
metalcloud-cli switch-defaults delete -ids "834, 835, 836"
`,
	},
}

func switchDefaultsListCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	var datacenterName string
	if datacenterNameValue, ok := command.GetStringParamOk(c.Arguments["datacenter_name"]); !ok {
		return "", fmt.Errorf("-datacenter is required")
	} else {
		datacenterName = datacenterNameValue
	}

	list, err := client.SwitchDeviceDefaults(datacenterName)
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
			FieldName: "Serial Number",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "Management MAC",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "Position",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "Identifier",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "ASN",
			FieldType: tableformatter.TypeInt,
			FieldSize: 6,
		},
		{
			FieldName: "Loopback IPv4",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "Loopback IPv6",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "VTEP IPv4",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "VTEP IPv6",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
	}

	data := [][]interface{}{}

	for _, obj := range *list {
		serialNumber := obj.NetworkEquipmentSerialNumber
		if serialNumber == nil {
			serialNumber = new(string)
		}

		macAddress := obj.NetworkEquipmentManagementMacAddress
		if macAddress == nil {
			macAddress = new(string)
		}

		position := obj.NetworkEquipmentPosition
		if position == nil {
			position = new(string)
		}

		identifierString := obj.NetworkEquipmentIdentifierString
		if identifierString == nil {
			identifierString = new(string)
		}

		asn := obj.NetworkEquipmentAsn
		if asn == nil {
			asn = new(int)
		}

		loopbackAddressIpv4 := obj.NetworkEquipmentLoopbackAddressIpv4
		if loopbackAddressIpv4 == nil {
			loopbackAddressIpv4 = new(string)
		}

		loopbackAddressIpv6 := obj.NetworkEquipmentLoopbackAddressIpv6
		if loopbackAddressIpv6 == nil {
			loopbackAddressIpv6 = new(string)
		}

		vtepAddressIpv4 := obj.NetworkEquipmentVtepAddressIpv4
		if vtepAddressIpv4 == nil {
			vtepAddressIpv4 = new(string)
		}

		vtepAddressIpv6 := obj.NetworkEquipmentVtepAddressIpv6
		if vtepAddressIpv6 == nil {
			vtepAddressIpv6 = new(string)
		}

		data = append(data, []interface{}{
			obj.NetworkEquipmentDefaultsID,
			*serialNumber,
			*macAddress,
			*position,
			*identifierString,
			*asn,
			*loopbackAddressIpv4,
			*loopbackAddressIpv6,
			*vtepAddressIpv4,
			*vtepAddressIpv6,
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
	return table.RenderTable("Switch defaults", "", command.GetStringParam(c.Arguments["format"]))

}

func switchDefaultsCreateCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	var defaults []metalcloud.SwitchDeviceDefaults

	filePath, ok := command.GetStringParamOk(c.Arguments["read_config_from_file"])
	if !ok {
		return "", fmt.Errorf("-raw-config is required")
	}

	_, err := os.Open(filePath)

	if err != nil {
		return "", fmt.Errorf("error reading file %s: %s", filePath, err.Error())
	}

	defaults, err = getMultipleSwitchDefaultsFromYamlFile(filePath)
	if err != nil {
		return "", err
	}

	for idx, obj := range defaults {
		if obj.DatacenterName == "" {
			return "", fmt.Errorf("datacenter name is required for switch defaults #%d.", idx+1)
		}

		if obj.NetworkEquipmentSerialNumber == new(string) && obj.NetworkEquipmentManagementMacAddress == new(string) {
			return "", fmt.Errorf("at least one of serial number or management MAC address must be provided for switch defaults #%d.", idx+1)
		}
	}

	err = client.SwitchDeviceDefaultsCreate(defaults)
	if err != nil {
		return "", err
	}

	fmt.Printf("Created %d switch defaults.\n", len(defaults))
	return "", err
}

func switchDefaultsDeleteCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	ids := ""
	if idsValue, ok := command.GetStringParamOk(c.Arguments["switch_defaults_ids"]); !ok {
		return "", fmt.Errorf("-ids is required")
	} else {
		ids = idsValue
	}

	defaultsIds := []int{}
	idsList := strings.Split(ids, ",")
	for _, id := range idsList {
		value, err := strconv.Atoi(strings.TrimSpace(id))
		if err != nil {
			return "", fmt.Errorf("invalid switch defaults id: %s", id)
		}

		defaultsIds = append(defaultsIds, value)
	}

	err := client.SwitchDeviceDefaultsDelete(defaultsIds)
	return "", err
}

func getMultipleSwitchDefaultsFromYamlFile(filePath string) ([]metalcloud.SwitchDeviceDefaults, error) {
	file, err := os.Open(filePath)
	defer file.Close()

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []metalcloud.SwitchDeviceDefaults{}, nil
		} else {
			return []metalcloud.SwitchDeviceDefaults{}, err
		}
	}

	decoder := yaml.NewDecoder(file)

	defaults := []metalcloud.SwitchDeviceDefaults{}

	for true {
		var defaultObj metalcloud.SwitchDeviceDefaults

		err = decoder.Decode(&defaultObj)
		if err == nil {
			defaults = append(defaults, defaultObj)
		}

		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			} else {
				return nil, fmt.Errorf("Error while reading %s: %v. Make sure the file is of the yaml format.", filePath, err)
			}
		}
	}

	return defaults, nil
}
