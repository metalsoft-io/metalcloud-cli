package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	metalcloud "github.com/bigstepinc/metal-cloud-sdk-go"
)

//driveArrayCmds commands affecting instance arrays
var driveArrayCmds = []Command{

	Command{
		Description:  "Creates a drive array.",
		Subject:      "drive_array",
		AltSubject:   "da",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("drive_array", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_id":                      c.FlagSet.Int("infra", _nilDefaultInt, "(Required) Infrastrucure ID"),
				"instance_array_id":                      c.FlagSet.Int("ia", _nilDefaultInt, "(Required) The id of the instance array it is attached to. It can be zero for unattached Drive Arrays"),
				"drive_array_label":                      c.FlagSet.String("label", _nilDefaultStr, "(Required) The label of the drive array"),
				"drive_array_storage_type":               c.FlagSet.String("type", _nilDefaultStr, "Possible values: iscsi_ssd, iscsi_hdd"),
				"drive_size_mbytes_default":              c.FlagSet.Int("size", _nilDefaultInt, "(Optional, default = 40960) Drive arrays's size in MBytes"),
				"drive_array_count":                      c.FlagSet.Int("count", _nilDefaultInt, "DriveArrays's drive count. Use this only for unconnected DriveArrays."),
				"drive_array_expand_with_instance_array": c.FlagSet.Bool("expand_with_ia", true, "Auto-expand when the connected instance array expands"),
				"volume_template_id":                     c.FlagSet.Int("template", _nilDefaultInt, "DriveArrays's volume template to clone when creating Drives"),
				"return_id":                              c.FlagSet.Bool("return_id", false, "(Optional) Will print the ID of the created Drive Array. Useful for automating tasks."),
			}
		},
		ExecuteFunc: driveArrayCreateCmd,
	},
	Command{
		Description:  "Edit a drive array.",
		Subject:      "drive_array",
		AltSubject:   "da",
		Predicate:    "edit",
		AltPredicate: "alter",
		FlagSet:      flag.NewFlagSet("edit_drive_array", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"drive_array_id":                         c.FlagSet.Int("id", _nilDefaultInt, "(Required) Drive Array's ID"),
				"instance_array_id":                      c.FlagSet.Int("ia", _nilDefaultInt, "(Required) The id of the instance array it is attached to. It can be zero for unattached Drive Arrays"),
				"drive_array_label":                      c.FlagSet.String("label", _nilDefaultStr, "(Required) The label of the drive array"),
				"drive_array_storage_type":               c.FlagSet.String("type", _nilDefaultStr, "Possible values: iscsi_ssd, iscsi_hdd"),
				"drive_size_mbytes_default":              c.FlagSet.Int("size", _nilDefaultInt, "(Optional, default = 40960) Drive arrays's size in MBytes"),
				"drive_array_count":                      c.FlagSet.Int("count", _nilDefaultInt, "DriveArrays's drive count. Use this only for unconnected DriveArrays."),
				"drive_array_expand_with_instance_array": c.FlagSet.Bool("expand_with_ia", true, "Auto-expand when the connected instance array expands"),
				"volume_template_id":                     c.FlagSet.Int("template", _nilDefaultInt, "DriveArrays's volume template to clone when creating Drives"),
			}
		},
		ExecuteFunc: driveArrayEditCmd,
	},
	Command{
		Description:  "Lists all drive arrays of an infrastructure.",
		Subject:      "drive_array",
		AltSubject:   "da",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list drive_array", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_id": c.FlagSet.Int("infra", 0, "(Required) Infrastrucure ID"),
				"format":            c.FlagSet.String("format", "", "The output format. Supported values are 'json','csv'. The default format is human readable."),
			}
		},
		ExecuteFunc: driveArrayListCmd,
	},
	Command{
		Description:  "Delete a drive array.",
		Subject:      "drive_array",
		AltSubject:   "da",
		Predicate:    "delete",
		AltPredicate: "rm",
		FlagSet:      flag.NewFlagSet("delete drive_array", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"drive_array_id": c.FlagSet.Int("id", 0, "(Required) DriveArray's ID"),
				"autoconfirm":    c.FlagSet.Bool("autoconfirm", false, "If true it does not ask for confirmation anymore"),
			}
		},
		ExecuteFunc: driveArrayDeleteCmd,
	},
}

func driveArrayCreateCmd(c *Command, client MetalCloudClient) (string, error) {

	infrastructureID := c.Arguments["infrastructure_id"]

	if infrastructureID == nil || *infrastructureID.(*int) == 0 {
		return "", fmt.Errorf("-infra <infrastructure_id> is required")
	}

	if c.Arguments["instance_array_id"] == nil {
		return "", fmt.Errorf("-ia <instance_array_id> is required. Use 0 for unattached")
	}

	if c.Arguments["volume_template_id"] == nil {
		return "", fmt.Errorf("-template <volume_template_id> is required. Use 0 for unformatted drive")
	}

	da := argsToDriveArray(c.Arguments)

	if da.DriveArrayLabel == "" {
		return "", fmt.Errorf("-label <drive_array_label> is required")
	}

	retDA, err := client.DriveArrayCreate(*infrastructureID.(*int), *da)
	if err != nil {
		return "", err
	}

	if c.Arguments["return_id"] != nil && *c.Arguments["return_id"].(*bool) == true {
		return fmt.Sprintf("%d", retDA.DriveArrayID), nil
	}

	return "", err
}

func driveArrayEditCmd(c *Command, client MetalCloudClient) (string, error) {

	id := c.Arguments["drive_array_id"]

	if id == nil || *id.(*int) == 0 {
		return "", fmt.Errorf("-id <drive_array_id> is required")
	}

	retDA, err := client.DriveArrayGet(*id.(*int))
	if err != nil {
		return "", err
	}

	dao := retDA.DriveArrayOperation

	argsToDriveArrayOperation(c.Arguments, dao)

	_, err = client.DriveArrayEdit(*id.(*int), *dao)

	return "", err
}

func driveArrayListCmd(c *Command, client MetalCloudClient) (string, error) {

	infrastructureID := c.Arguments["infrastructure_id"]

	if infrastructureID == nil || *infrastructureID.(*int) == 0 {
		return "", fmt.Errorf("-infra <infrastructure_id> is required")
	}

	daList, err := client.DriveArrays(*infrastructureID.(*int))
	if err != nil {
		return "", err
	}

	schema := []SchemaField{
		SchemaField{
			FieldName: "ID",
			FieldType: TypeInt,
			FieldSize: 6,
		},
		SchemaField{
			FieldName: "LABEL",
			FieldType: TypeString,
			FieldSize: 30,
		},
		SchemaField{
			FieldName: "STATUS",
			FieldType: TypeString,
			FieldSize: 10,
		},
		SchemaField{
			FieldName: "SIZE (MB)",
			FieldType: TypeInt,
			FieldSize: 10,
		},
		SchemaField{
			FieldName: "TYPE",
			FieldType: TypeString,
			FieldSize: 10,
		},
		SchemaField{
			FieldName: "ATTACHED TO",
			FieldType: TypeString,
			FieldSize: 30,
		},
		SchemaField{
			FieldName: "DRV_CNT",
			FieldType: TypeInt,
			FieldSize: 10,
		},
		SchemaField{
			FieldName: "TEMPLATE",
			FieldType: TypeString,
			FieldSize: 25,
		},
	}

	data := [][]interface{}{}
	for _, da := range *daList {
		status := da.DriveArrayServiceStatus

		if da.DriveArrayServiceStatus != "ordered" && da.DriveArrayOperation.DriveArrayDeployType == "edit" && da.DriveArrayOperation.DriveArrayDeployStatus == "not_started" {
			status = "edited"
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
		if da.InstanceArrayID != 0 {
			ia, err := client.InstanceArrayGet(da.DriveArrayOperation.InstanceArrayID)
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

	var sb strings.Builder

	format := "text"
	if v := c.Arguments["format"]; v != _nilDefaultStr {
		format = *v.(*string)
	}

	switch format {
	case "json", "JSON":
		ret, err := GetTableAsJSONString(data, schema)
		if err != nil {
			return "", err
		}
		sb.WriteString(ret)
	case "csv", "CSV":
		ret, err := GetTableAsCSVString(data, schema)
		if err != nil {
			return "", err
		}
		sb.WriteString(ret)

	default:

		sb.WriteString(GetTableAsString(data, schema))
		sb.WriteString(fmt.Sprintf("Total: %d Drive Arrays\n\n", len(*daList)))

	}

	return sb.String(), nil
}

func driveArrayDeleteCmd(c *Command, client MetalCloudClient) (string, error) {

	driveArrayID := c.Arguments["drive_array_id"]

	if driveArrayID == nil || *driveArrayID.(*int) == 0 {
		return "", fmt.Errorf("-id <drive_array_id> is required")
	}

	retDA, err2 := client.DriveArrayGet(*driveArrayID.(*int))
	if err2 != nil {
		return "", err2
	}

	var retIA *metalcloud.InstanceArray
	if retDA.InstanceArrayID != 0 {
		retIA, err2 = client.InstanceArrayGet(retDA.InstanceArrayID)
	}

	retInfra, err2 := client.InfrastructureGet(retIA.InfrastructureID)
	if err2 != nil {
		return "", err2
	}

	confirm := false

	if c.Arguments["autoconfirm"] != nil && *c.Arguments["autoconfirm"].(*bool) == true {
		confirm = true
	} else {

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

		confirm = requestConfirmation(confirmationMessage)

	}

	if !confirm {
		return "", fmt.Errorf("Operation not confirmed. Aborting")
	}

	err := client.DriveArrayDelete(*driveArrayID.(*int))

	return "", err
}

func argsToDriveArray(m map[string]interface{}) *metalcloud.DriveArray {
	obj := metalcloud.DriveArray{}

	if v := m["drive_array_id"]; v != nil && *v.(*int) != _nilDefaultInt {
		obj.DriveArrayID = *v.(*int)
	}

	if v := m["drive_array_label"]; v != nil && *v.(*string) != _nilDefaultStr {
		obj.DriveArrayLabel = *v.(*string)
	}

	if v := m["volume_template_id"]; v != nil && *v.(*int) != _nilDefaultInt {
		obj.VolumeTemplateID = *v.(*int)
	}

	if v := m["drive_array_storage_type"]; v != nil && *v.(*string) != _nilDefaultStr {
		obj.DriveArrayStorageType = *v.(*string)
	}

	if v := m["drive_array_count"]; v != nil && *v.(*int) != _nilDefaultInt {
		obj.DriveArrayCount = *v.(*int)
	}

	if v := m["drive_size_mbytes_default"]; v != nil && *v.(*int) != _nilDefaultInt {
		obj.DriveSizeMBytesDefault = *v.(*int)
	}

	if v := m["drive_array_expand_with_instance_array"]; v != nil {
		obj.DriveArrayExpandWithInstanceArray = *v.(*bool)
	}

	if v := m["instance_array_id"]; v != nil && *v.(*int) != _nilDefaultInt {
		obj.InstanceArrayID = *v.(*int)
	}

	return &obj
}

func argsToDriveArrayOperation(m map[string]interface{}, dao *metalcloud.DriveArrayOperation) {

	if v := m["drive_array_id"]; v != nil && *v.(*int) != _nilDefaultInt {
		dao.DriveArrayID = *v.(*int)
	}

	if v := m["drive_array_label"]; v != nil && *v.(*string) != _nilDefaultStr {
		dao.DriveArrayLabel = *v.(*string)
	}

	if v := m["volume_template_id"]; v != nil && *v.(*int) != _nilDefaultInt {
		dao.VolumeTemplateID = *v.(*int)
	}

	if v := m["drive_array_storage_type"]; v != nil && *v.(*string) != _nilDefaultStr {
		dao.DriveArrayStorageType = *v.(*string)
	}

	if v := m["drive_array_count"]; v != nil && *v.(*int) != _nilDefaultInt {
		dao.DriveArrayCount = *v.(*int)
	}

	if v := m["drive_size_mbytes_default"]; v != nil && *v.(*int) != _nilDefaultInt {
		dao.DriveSizeMBytesDefault = *v.(*int)
	}

	if v := m["drive_array_expand_with_instance_array"]; v != nil {
		dao.DriveArrayExpandWithInstanceArray = *v.(*bool)
	}

	if v := m["instance_array_id"]; v != nil && *v.(*int) != _nilDefaultInt {
		dao.InstanceArrayID = *v.(*int)
	}

}
