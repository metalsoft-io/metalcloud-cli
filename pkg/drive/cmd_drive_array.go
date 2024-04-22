package drive

import (
	"flag"
	"fmt"
	"os"
	"strings"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
	"github.com/metalsoft-io/metalcloud-cli/internal/colors"
	"github.com/metalsoft-io/metalcloud-cli/internal/command"
	"github.com/metalsoft-io/metalcloud-cli/internal/objects"
	"github.com/metalsoft-io/tableformatter"
)

var DriveArrayCmds = []command.Command{
	{
		Description:  "Creates a drive array.",
		Subject:      "drive-array",
		AltSubject:   "da",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("drive-array", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"read_config_from_file": c.FlagSet.String("f", command.NilDefaultStr, colors.Red("(Required)")+" Read configuration from file in the format specified with --format."),
				"format":                c.FlagSet.String("format", "yaml", "The input format. Supported values are 'json','yaml'. The default format is json."),
				"return_id":             c.FlagSet.Bool("return-id", false, "(Optional) Will print the ID of the created Drive Array. Useful for automating tasks."),
			}
		},
		ExecuteFunc: driveArrayCreateCmd,
	},
	{
		Description:  "Edit a drive array.",
		Subject:      "drive-array",
		AltSubject:   "da",
		Predicate:    "update",
		AltPredicate: "edit",
		FlagSet:      flag.NewFlagSet("update_drive_array", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"read_config_from_file": c.FlagSet.String("f", command.NilDefaultStr, colors.Red("(Required)")+" Read configuration from file in the format specified with --format."),
				"format":                c.FlagSet.String("format", "yaml", "The input format. Supported values are 'json','yaml'. The default format is json."),
			}
		},
		ExecuteFunc: driveArrayUpdateCmd,
	},
	{
		Description:  "Lists all drive arrays of an infrastructure.",
		Subject:      "drive-array",
		AltSubject:   "da",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list drive_array", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_id_or_label": c.FlagSet.String("infra", command.NilDefaultStr, colors.Red("(Required)")+" Infrastructure's id or label. Note that the 'label' this be ambiguous in certain situations."),
				"format":                     c.FlagSet.String("format", command.NilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
			}
		},
		ExecuteFunc: driveArrayListCmd,
	},
	{
		Description:  "Delete a drive array.",
		Subject:      "drive-array",
		AltSubject:   "da",
		Predicate:    "delete",
		AltPredicate: "rm",
		FlagSet:      flag.NewFlagSet("delete drive_array", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"drive_array_id_or_label": c.FlagSet.String("id", command.NilDefaultStr, colors.Red("(Required)")+" Drive Array's ID or label. Note that using the label can be ambiguous and is slower."),
				"autoconfirm":             c.FlagSet.Bool("autoconfirm", false, colors.Green("(Flag)")+" If set it will assume action is confirmed"),
			}
		},
		ExecuteFunc: driveArrayDeleteCmd,
	},
	{
		Description:  "Gets a drive array.",
		Subject:      "drive-array",
		AltSubject:   "da",
		Predicate:    "get",
		AltPredicate: "show",
		FlagSet:      flag.NewFlagSet("show drive_array", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"drive_array_id_or_label": c.FlagSet.String("id", command.NilDefaultStr, colors.Red("(Required)")+" Drive Array's ID or label. Note that using the label can be ambiguous and is slower."),
				"format":                  c.FlagSet.String("format", "yaml", "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
			}
		},
		ExecuteFunc: driveArrayGetCmd,
	},
	{
		Description:  "Lists a drive array's drives.",
		Subject:      "drive-array",
		AltSubject:   "da",
		Predicate:    "list-drives",
		AltPredicate: "show-drives",
		FlagSet:      flag.NewFlagSet("show drive_array drives", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"drive_array_id_or_label": c.FlagSet.String("id", command.NilDefaultStr, colors.Red("(Required)")+" Drive Array's ID or label. Note that using the label can be ambiguous and is slower."),
				"format":                  c.FlagSet.String("format", "yaml", "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
			}
		},
		ExecuteFunc: driveArrayDrivesGetCmd,
	},
}

func driveArrayCreateCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	obj, err := objects.ReadSingleObjectFromCommand(c, client)
	if err != nil {
		return "", err
	}
	da := (*obj).(metalcloud.DriveArray)

	createdDA, err := client.DriveArrayCreate(da.InfrastructureID, da)
	if err != nil {
		return "", err
	}

	if command.GetBoolParam(c.Arguments["return_id"]) {
		return fmt.Sprintf("%d", createdDA.DriveArrayID), nil
	}

	return "", err
}

func driveArrayUpdateCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	obj, err := objects.ReadSingleObjectFromCommand(c, client)
	if err != nil {
		return "", err
	}
	da := (*obj).(metalcloud.DriveArray)

	_, err = client.DriveArrayEdit(da.DriveArrayID, *da.DriveArrayOperation)
	if err != nil {
		return "", err
	}

	return "", err
}

func driveArrayListCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
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

	daList, err := client.DriveArrays(infraID)
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
			FieldSize: 30,
		},
		{
			FieldName: "DRV_CNT",
			FieldType: tableformatter.TypeInt,
			FieldSize: 10,
		},
		{
			FieldName: "TEMPLATE",
			FieldType: tableformatter.TypeString,
			FieldSize: 25,
		},
	}

	data := [][]interface{}{}
	for _, da := range *daList {
		status := da.DriveArrayServiceStatus

		if da.DriveArrayServiceStatus != "ordered" && da.DriveArrayOperation.DriveArrayDeployType == "edit" && da.DriveArrayOperation.DriveArrayDeployStatus == "not_started" {
			status = "edited"
		}

		if da.DriveArrayServiceStatus != "ordered" && da.DriveArrayOperation.DriveArrayDeployType == "delete" && da.DriveArrayOperation.DriveArrayDeployStatus == "not_started" {
			status = "marked for delete"
		}

		volumeTemplateName := ""
		if da.VolumeTemplateID != 0 {
			vt, err := client.VolumeTemplateGet(da.DriveArrayOperation.VolumeTemplateID)
			if err != nil {
				return "", err
			}

			volumeTemplateName = fmt.Sprintf("%s (#%d)", vt.VolumeTemplateDisplayName, vt.VolumeTemplateID)
		}

		instanceArrayLabel := ""
		if da.DriveArrayOperation.InstanceArrayID != nil && da.DriveArrayOperation.InstanceArrayID != 0 {
			var instanceArrayID int

			switch da.DriveArrayOperation.InstanceArrayID.(type) {
			case int:
				instanceArrayID = da.DriveArrayOperation.InstanceArrayID.(int)
			case float64:
				instanceArrayID = int(da.DriveArrayOperation.InstanceArrayID.(float64))
			default:
				return "", fmt.Errorf("Instance array ID type invalid.")
			}

			ia, err := client.InstanceArrayGet(instanceArrayID)
			if err != nil {
				return "", err
			}
			instanceArrayLabel = fmt.Sprintf("%s (#%d)", ia.InstanceArrayLabel, ia.InstanceArrayID)
		}

		data = append(data, []interface{}{
			da.DriveArrayID,
			da.DriveArrayOperation.DriveArrayLabel,
			status,
			da.DriveArrayOperation.DriveSizeMBytesDefault,
			da.DriveArrayOperation.DriveArrayStorageType,
			instanceArrayLabel,
			da.DriveArrayOperation.DriveArrayCount,
			volumeTemplateName})
	}

	tableformatter.TableSorter(schema).OrderBy(schema[0].FieldName).Sort(data)

	table := tableformatter.Table{
		Data:   data,
		Schema: schema,
	}

	return table.RenderTable("Drive Arrays", "", command.GetStringParam(c.Arguments["format"]))
}

func driveArrayDeleteCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	retDA, err := command.GetDriveArrayFromCommand("id", c, client)
	if err != nil {
		return "", err
	}

	var retIA *metalcloud.InstanceArray

	if retDA.InstanceArrayID != 0 {
		retIA, err = client.InstanceArrayGet(retDA.InstanceArrayID)
		if err != nil {
			return "", err
		}
	}

	retInfra, err2 := client.InfrastructureGet(retDA.InfrastructureID)
	if err2 != nil {
		return "", err2
	}

	confirm, err := command.ConfirmCommand(c, func() string {

		var confirmationMessage string

		if retIA != nil {
			confirmationMessage = fmt.Sprintf("Deleting drive array %s (%d), attached to instance array (%s, %d) - from infrastructure %s (%d).  Are you sure? Type \"yes\" to continue:",
				retDA.DriveArrayLabel, retDA.DriveArrayID,
				retIA.InstanceArrayLabel, retIA.InstanceArrayID,
				retInfra.InfrastructureLabel, retInfra.InfrastructureID)
		} else {
			confirmationMessage = fmt.Sprintf("Deleting drive array %s (%d), unattached - from infrastructure %s (%d).  Are you sure? Type \"yes\" to continue:",
				retDA.DriveArrayLabel, retDA.DriveArrayID,
				retInfra.InfrastructureLabel, retInfra.InfrastructureID)
		}

		//this is simply so that we don't output a text on the command line
		if strings.HasSuffix(os.Args[0], ".test") {
			confirmationMessage = ""
		}

		return confirmationMessage
	})
	if err != nil {
		return "", err
	}

	if confirm {
		return "", client.DriveArrayDelete(retDA.DriveArrayID)
	}

	return "", fmt.Errorf("operation not confirmed. Aborting")
}

func driveArrayGetCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	driveArray, err := command.GetDriveArrayFromCommand("id", c, client)
	if err != nil {
		return "", err
	}

	format := command.GetStringParam(c.Arguments["format"])
	ret, err := objects.RenderRawObject(*driveArray, format, "DriveArray")
	if err != nil {
		return "", err
	}

	return ret, nil
}

func driveArrayDrivesGetCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	driveArray, err := command.GetDriveArrayFromCommand("id", c, client)
	if err != nil {
		return "", err
	}

	drives, err := client.DriveArrayDrives(driveArray.DriveArrayID)
	if err != nil {
		return "", err
	}

	format := command.GetStringParam(c.Arguments["format"])
	ret, err := objects.RenderRawObject(*drives, format, "DriveArrayDrives")
	if err != nil {
		return "", err
	}

	return ret, nil
}
