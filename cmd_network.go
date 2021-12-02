package main

import (
	"flag"
	"fmt"
	"strconv"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
	"github.com/metalsoft-io/tableformatter"
)

var networkCmds = []Command{
	{
		Description:  "List all networks for an instance array.",
		Subject:      "network",
		AltSubject:   "nw",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list network", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"instance_array_id_or_label": c.FlagSet.String("ia", _nilDefaultStr, "(Required) InstanceArray's id or label. Note that the label can be ambigous."),
				"format":                     c.FlagSet.String("format", "", "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
			}
		},
		ExecuteFunc: networkListCmd,
	},
}

func networkListCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	retIA, err := getInstanceArrayFromCommand("ia", c, client)
	if err != nil {
		return "", err
	}
	if &retIA == nil {
		return "", fmt.Errorf("-ia is required")
	}

	dataNetworkAttachments := [][]interface{}{}

	schemaNetworkAttachments := []tableformatter.SchemaField{

		{
			FieldName: "Port",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "Network",
			FieldType: tableformatter.TypeString,
			FieldSize: 10,
		},
		{
			FieldName: "Profile",
			FieldType: tableformatter.TypeString,
			FieldSize: 10,
		},
	}
	instanceArrayNetworkProfiles, err := client.NetworkProfileListByInstanceArray(retIA.InstanceArrayID)
	if err != nil {
		return "", err
	}
	for _, IAinterface := range retIA.InstanceArrayInterfaces {
		index := strconv.Itoa(IAinterface.InstanceArrayInterfaceIndex + 1)
		net := "unattached"
		profile := ""
		if IAinterface.NetworkID != 0 {
			n, err := client.NetworkGet(IAinterface.NetworkID)
			if err != nil {
				return "", err
			}
			profileId := (*instanceArrayNetworkProfiles)[IAinterface.NetworkID]
			if profileId != 0 {
				networkProfile, err := client.NetworkProfileGet(profileId)
				if err != nil {
					return "", err
				}
				profile = networkProfile.NetworkProfileLabel + " (#" + strconv.Itoa(profileId) + ")"
			}

			net = n.NetworkType + "(#" + strconv.Itoa(IAinterface.NetworkID) + ")"
		}

		IAdataRow := []interface{}{
			"#" + index,
			net,
			profile,
		}

		dataNetworkAttachments = append(dataNetworkAttachments, IAdataRow)
	}

	tableNetworkAttachments := tableformatter.Table{
		Data:   dataNetworkAttachments,
		Schema: schemaNetworkAttachments,
	}
	subtitleNetworkAttachmentsRender := "NETWORK ATTACHEMENTS\n--------------------\nNetworks to which this instance array is attached to:\n"

	return tableNetworkAttachments.RenderTable("", subtitleNetworkAttachmentsRender, getStringParam(c.Arguments["format"]))
}
