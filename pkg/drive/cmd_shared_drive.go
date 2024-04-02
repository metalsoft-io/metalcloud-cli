package drive

import (
	"flag"
	"fmt"
	"strings"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
	"github.com/metalsoft-io/metalcloud-cli/internal/colors"
	"github.com/metalsoft-io/metalcloud-cli/internal/command"
	"github.com/metalsoft-io/metalcloud-cli/internal/configuration"
	"github.com/metalsoft-io/tableformatter"
)

var SharedDriveCmds = []command.Command{
	{
		Description:  "Lists all shared drives  of an infrastructure.",
		Subject:      "shared-drive",
		AltSubject:   "shared-drives",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list shared drives", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_id_or_label": c.FlagSet.String("infra", command.NilDefaultStr, colors.Red("(Required)")+" Infrastructure's id or label. Note that the 'label' this be ambiguous in certain situations."),
				"format":                     c.FlagSet.String("format", "", "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
			}
		},
		ExecuteFunc: sharedDriveListCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
}

func sharedDriveListCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	infraIDStr, err := command.GetParam(c, "infrastructure_id_or_label", "infra")
	if err != nil {
		return "", err
	}

	infraID, err := command.GetIDOrDo(*infraIDStr.(*string), func(label string) (int, error) {
		ia, err := client.InfrastructureGetByLabel(label)
		if err != nil {
			return 0, err
		}
		return ia.InfrastructureID, nil
	},
	)
	if err != nil {
		return "", err
	}

	sdList, err := client.SharedDrives(infraID)
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
			FieldName: "STATUS",
			FieldType: tableformatter.TypeString,
			FieldSize: 10,
		},
		{
			FieldName: "SIZE (MB)",
			FieldType: tableformatter.TypeInt,
			FieldSize: 10,
		},
		{
			FieldName: "TYPE",
			FieldType: tableformatter.TypeString,
			FieldSize: 10,
		},
		{
			FieldName: "ATTACHED TO",
			FieldType: tableformatter.TypeString,
			FieldSize: 40,
		},
		{
			FieldName: "IO LIMIT",
			FieldType: tableformatter.TypeString,
			FieldSize: 10,
		},
		{
			FieldName: "WWN",
			FieldType: tableformatter.TypeString,
			FieldSize: 20,
		},
	}

	data := [][]interface{}{}
	for _, sd := range *sdList {
		status := sd.SharedDriveServiceStatus

		if sd.SharedDriveServiceStatus != "ordered" && sd.SharedDriveOperation.SharedDriveServiceStatus == "edit" && sd.SharedDriveOperation.SharedDriveDeployStatus == "not_started" {
			status = "edited"
		}

		if sd.SharedDriveServiceStatus != "ordered" && sd.SharedDriveOperation.SharedDriveServiceStatus == "delete" && sd.SharedDriveOperation.SharedDriveDeployStatus == "not_started" {
			status = "marked for delete"
		}

		attachedInstanceArrays := []string{}

		for _, instanceArrayID := range sd.SharedDriveAttachedInstanceArrays {
			ia, err := client.InstanceArrayGet(instanceArrayID)
			if err != nil {
				return "", err
			}
			attachedInstanceArrays = append(attachedInstanceArrays, fmt.Sprintf("%s (#%d)", ia.InstanceArrayLabel, ia.InstanceArrayID))
		}

		attachedInstanceArraysList := strings.Join(attachedInstanceArrays, ",")

		data = append(data, []interface{}{
			sd.SharedDriveID,
			sd.SharedDriveOperation.SharedDriveLabel,
			status,
			sd.SharedDriveOperation.SharedDriveSizeMbytes,
			sd.SharedDriveOperation.SharedDriveStorageType,
			attachedInstanceArraysList,
			sd.SharedDriveIOLimitPolicy,
			sd.SharedDriveWWN})
	}

	tableformatter.TableSorter(schema).OrderBy(schema[0].FieldName).Sort(data)

	table := tableformatter.Table{
		Data:   data,
		Schema: schema,
	}

	return table.RenderTable("Shared drives", "", command.GetStringParam(c.Arguments["format"]))
}
