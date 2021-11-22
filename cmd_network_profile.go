package main

import (
	"flag"

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
			FieldType: tableformatter.TypeInt,
			FieldSize: 30,
		},
		{
			FieldName: "CREATED",
			FieldType: tableformatter.TypeString,
			FieldSize: 10,
		},
		{
			FieldName: "DELETED",
			FieldType: tableformatter.TypeString,
			FieldSize: 10,
		},
	}

	data := [][]interface{}{}
	for _, np := range *npList {
		// status := np.DriveArrayServiceStatus

		// if da.DriveArrayServiceStatus != "ordered" && da.DriveArrayOperation.DriveArrayDeployType == "edit" && da.DriveArrayOperation.DriveArrayDeployStatus == "not_started" {
		// 	status = "edited"
		// }

		// if da.DriveArrayServiceStatus != "ordered" && da.DriveArrayOperation.DriveArrayDeployType == "delete" && da.DriveArrayOperation.DriveArrayDeployStatus == "not_started" {
		// 	status = "marked for delete"
		// }

		// volumeTemplateName := ""
		// if da.VolumeTemplateID != 0 {
		// 	vt, err := client.VolumeTemplateGet(da.DriveArrayOperation.VolumeTemplateID)
		// 	if err != nil {
		// 		return "", err
		// 	}

		// 	volumeTemplateName = fmt.Sprintf("%s (#%d)", vt.VolumeTemplateDisplayName, vt.VolumeTemplateID)
		// }

		// instanceArrayLabel := ""
		// if da.DriveArrayOperation.InstanceArrayID != 0 {
		// 	ia, err := client.InstanceArrayGet(da.DriveArrayOperation.InstanceArrayID)
		// 	if err != nil {
		// 		return "", err
		// 	}
		// 	instanceArrayLabel = fmt.Sprintf("%s (#%d)", ia.InstanceArrayLabel, ia.InstanceArrayID)
		// }

		data = append(data, []interface{}{
			np.NetworkProfileID,
			// da.DriveArrayOperation.DriveArrayLabel,
			// status,
			// da.DriveArrayOperation.DriveSizeMBytesDefault,
			// da.DriveArrayOperation.DriveArrayStorageType,
			// instanceArrayLabel,
			// da.DriveArrayOperation.DriveArrayCount,
			// volumeTemplateName
		})
	}

	tableformatter.TableSorter(schema).OrderBy(schema[0].FieldName).Sort(data)

	table := tableformatter.Table{
		Data:   data,
		Schema: schema,
	}

	return table.RenderTable("Network Profiles", "", getStringParam(c.Arguments["format"]))
}
