package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	metalcloud "github.com/bigstepinc/metal-cloud-sdk-go"
)

//instanceArrayCmds commands affecting instance arrays
var instanceArrayCmds = []Command{

	Command{
		Description:  "Creates an instance array.",
		Subject:      "instance_array",
		AltSubject:   "ia",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("instance_array", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_id_or_label":          c.FlagSet.String("infra", _nilDefaultStr, "(Required) Infrastructure's id or label. Note that the 'label' this be ambiguous in certain situations."),
				"instance_array_instance_count":       c.FlagSet.Int("instance_count", _nilDefaultInt, "(Required) Instance count of this instance array"),
				"instance_array_label":                c.FlagSet.String("label", _nilDefaultStr, "InstanceArray's label"),
				"instance_array_ram_gbytes":           c.FlagSet.Int("ram", _nilDefaultInt, "InstanceArray's minimum RAM (GB)"),
				"instance_array_processor_count":      c.FlagSet.Int("proc", _nilDefaultInt, "InstanceArray's minimum processor count"),
				"instance_array_processor_core_mhz":   c.FlagSet.Int("proc_freq", _nilDefaultInt, "InstanceArray's minimum processor frequency (Mhz)"),
				"instance_array_processor_core_count": c.FlagSet.Int("proc_core_count", _nilDefaultInt, "InstanceArray's minimum processor core count"),
				"instance_array_disk_count":           c.FlagSet.Int("disks", _nilDefaultInt, "InstanceArray's number of local drives"),
				"instance_array_disk_size_mbytes":     c.FlagSet.Int("disk_size", _nilDefaultInt, "InstanceArray's local disk sizes"),
				"instance_array_boot_method":          c.FlagSet.String("boot", _nilDefaultStr, "InstanceArray's boot type:'pxe_iscsi','local_drives'"),
				"instance_array_firewall_not_managed": c.FlagSet.Bool("un_managed_fw", false, "(Flag) If set InstanceArray's firewall management on or off"),
				"volume_template_id":                  c.FlagSet.Int("template", _nilDefaultInt, "InstanceArray's volume template when booting from for local drives"),
				"return_id":                           c.FlagSet.Bool("return_id", false, "(Flag) If set will print the ID of the created Instance Array. Useful for automating tasks."),
			}
		},
		ExecuteFunc: instanceArrayCreateCmd,
	},
	Command{
		Description:  "Lists all instance arrays of an infrastructure.",
		Subject:      "instance_array",
		AltSubject:   "ia",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list instance_array", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_id": c.FlagSet.Int("infra", 0, "(Required) Infrastrucure ID"),
				"format":            c.FlagSet.String("format", "", "The output format. Supported values are 'json','csv'. The default format is human readable."),
			}
		},
		ExecuteFunc: instanceArrayListCmd,
	},
	Command{
		Description:  "Delete instance array.",
		Subject:      "instance_array",
		AltSubject:   "ia",
		Predicate:    "delete",
		AltPredicate: "rm",
		FlagSet:      flag.NewFlagSet("list instance_array", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"instance_array_id": c.FlagSet.Int("id", 0, "(Required) InstanceArray ID"),
				"autoconfirm":       c.FlagSet.Bool("autoconfirm", false, "If true it does not ask for confirmation anymore"),
			}
		},
		ExecuteFunc: instanceArrayDeleteCmd,
	},
	Command{
		Description:  "Edits an instance array.",
		Subject:      "instance_array",
		AltSubject:   "ia",
		Predicate:    "edit",
		AltPredicate: "alter",
		FlagSet:      flag.NewFlagSet("instance_array", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"instance_array_id":                   c.FlagSet.Int("id", _nilDefaultInt, "(Required) InstanceArray's id"),
				"instance_array_instance_count":       c.FlagSet.Int("instance_count", _nilDefaultInt, "Instance count of this instance array"),
				"instance_array_label":                c.FlagSet.String("label", _nilDefaultStr, "(Required) InstanceArray's label"),
				"instance_array_ram_gbytes":           c.FlagSet.Int("ram", _nilDefaultInt, "InstanceArray's minimum RAM (GB)"),
				"instance_array_processor_count":      c.FlagSet.Int("proc", _nilDefaultInt, "InstanceArray's minimum processor count"),
				"instance_array_processor_core_mhz":   c.FlagSet.Int("proc_freq", _nilDefaultInt, "InstanceArray's minimum processor frequency (Mhz)"),
				"instance_array_processor_core_count": c.FlagSet.Int("proc_core_count", _nilDefaultInt, "InstanceArray's minimum processor core count"),
				"instance_array_disk_count":           c.FlagSet.Int("disks", _nilDefaultInt, "InstanceArray's number of local drives"),
				"instance_array_disk_size_mbytes":     c.FlagSet.Int("disk_size", _nilDefaultInt, "InstanceArray's local disk sizes"),
				"instance_array_boot_method":          c.FlagSet.String("boot", _nilDefaultStr, "InstanceArray's boot type:'pxe_iscsi','local_drives'"),
				"instance_array_firewall_not_managed": c.FlagSet.Bool("unmanaged_fw", false, "(Flag) If set InstanceArray's firewall management is off"),
				"volume_template_id":                  c.FlagSet.Int("template", _nilDefaultInt, "InstanceArray's volume template when booting from for local drives"),
				"bSwapExistingInstancesHardware":      c.FlagSet.Bool("swap_existing_hardware", false, "(Flag) If set all the hardware of the Instance objects is swapped to match the new InstanceArray specifications"),
				"no_bKeepDetachingDrives":             c.FlagSet.Bool("do_not_keep_detaching_drives", false, "(Flag) If set and the number of Instance objects is reduced, then the detaching Drive objects will be deleted. If it's set to true, the detaching Drive objects will not be deleted."),
			}
		},
		ExecuteFunc: instanceArrayEditCmd,
	},
}

func instanceArrayCreateCmd(c *Command, client MetalCloudClient) (string, error) {

	retInfra, err := getInfrastructureFromCommand(c, client)
	if err != nil {
		return "", err
	}

	ia := argsToInstanceArray(c.Arguments)

	if ia.InstanceArrayLabel == "" {
		return "", fmt.Errorf("-label <instance_array_label> is required")
	}

	retIA, err := client.InstanceArrayCreate(retInfra.InfrastructureID, *ia)
	if err != nil {
		return "", err
	}

	if c.Arguments["return_id"] != nil && *c.Arguments["return_id"].(*bool) {
		return fmt.Sprintf("%d", retIA.InstanceArrayID), nil
	}

	return "", err
}

func instanceArrayEditCmd(c *Command, client MetalCloudClient) (string, error) {

	instanceArrayID := c.Arguments["instance_array_id"]

	if instanceArrayID == nil || *instanceArrayID.(*int) == 0 {
		return "", fmt.Errorf("-id <instance_array_id> is required")
	}

	retIA, err := client.InstanceArrayGet(*instanceArrayID.(*int))
	if err != nil {
		return "", err
	}

	argsToInstanceArrayOperation(c.Arguments, retIA.InstanceArrayOperation)

	var bKeepDetachingDrives *bool
	if v := c.Arguments["not_bKeepDetachingDrives"]; v != nil {
		bVal := !*v.(*bool)
		bKeepDetachingDrives = &bVal
	}

	var bSwapExistingInstancesHardware *bool
	if c.Arguments["bSwapExistingInstancesHardware"] != nil {
		bSwapExistingInstancesHardware = c.Arguments["bSwapExistingInstancesHardware"].(*bool)
	}

	_, err = client.InstanceArrayEdit(*instanceArrayID.(*int),
		*retIA.InstanceArrayOperation,
		bSwapExistingInstancesHardware,
		bKeepDetachingDrives,
		nil,
		nil)

	return "", err
}

func instanceArrayListCmd(c *Command, client MetalCloudClient) (string, error) {

	infrastructureID := c.Arguments["infrastructure_id"]

	if infrastructureID == nil || *infrastructureID.(*int) == 0 {
		return "", fmt.Errorf("-infra <infrastructure_id> is required")
	}

	iaList, err := client.InstanceArrays(*infrastructureID.(*int))
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
			FieldSize: 20,
		},
		SchemaField{
			FieldName: "STATUS",
			FieldType: TypeString,
			FieldSize: 20,
		},
		SchemaField{
			FieldName: "INST_CNT",
			FieldType: TypeInt,
			FieldSize: 10,
		},
	}

	data := [][]interface{}{}
	for _, ia := range *iaList {
		status := ia.InstanceArrayServiceStatus
		if ia.InstanceArrayServiceStatus != "ordered" && ia.InstanceArrayOperation.InstanceArrayDeployType == "edit" && ia.InstanceArrayOperation.InstanceArrayDeployStatus == "not_started" {
			status = "edited"
		}
		data = append(data, []interface{}{
			ia.InstanceArrayID,
			ia.InstanceArrayOperation.InstanceArrayLabel,
			status,
			ia.InstanceArrayOperation.InstanceArrayInstanceCount})
	}

	var sb strings.Builder

	format := c.Arguments["format"]
	if format == nil {
		var f string
		f = ""
		format = &f
	}

	switch *format.(*string) {
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
		sb.WriteString(fmt.Sprintf("Total: %d Instance Arrays\n\n", len(*iaList)))

	}

	return sb.String(), nil
}

func instanceArrayDeleteCmd(c *Command, client MetalCloudClient) (string, error) {

	instanceArrayID := c.Arguments["instance_array_id"]

	if instanceArrayID == nil || *instanceArrayID.(*int) == 0 {
		return "", fmt.Errorf("-id <instance_array_id> is required")
	}

	retIA, err2 := client.InstanceArrayGet(*instanceArrayID.(*int))
	if err2 != nil {
		return "", err2
	}

	retInfra, err2 := client.InfrastructureGet(retIA.InfrastructureID)
	if err2 != nil {
		return "", err2
	}

	confirm := false

	if c.Arguments["autoconfirm"] != nil && *c.Arguments["autoconfirm"].(*bool) == true {
		confirm = true
	} else {

		confirmationMessage := fmt.Sprintf("Deleting instance array %s (%d) - from infrastructure %s (%d).  Are you sure? Type \"yes\" to continue:",
			retIA.InstanceArrayLabel, retIA.InstanceArrayID,
			retInfra.InfrastructureLabel, retInfra.InfrastructureID)

		//this is simply so that we don't output a text on the command line under go test
		if strings.HasSuffix(os.Args[0], ".test") {
			confirmationMessage = ""
		}

		confirm = requestConfirmation(confirmationMessage)
	}

	if !confirm {
		return "", fmt.Errorf("Operation not confirmed. Aborting")
	}

	err := client.InstanceArrayDelete(*instanceArrayID.(*int))

	return "", err
}

func argsToInstanceArray(m map[string]interface{}) *metalcloud.InstanceArray {
	ia := metalcloud.InstanceArray{}

	if v := m["instance_array_instance_count"]; v != nil && *v.(*int) != _nilDefaultInt {
		ia.InstanceArrayInstanceCount = *v.(*int)
	}

	if v := m["instance_array_label"]; v != nil && *v.(*string) != _nilDefaultStr {
		ia.InstanceArrayLabel = *v.(*string)
	}

	if v := m["instance_array_ram_gbytes"]; v != nil && *v.(*int) != _nilDefaultInt {
		ia.InstanceArrayRAMGbytes = *v.(*int)
	}

	if v := m["instance_array_processor_count"]; v != nil && *v.(*int) != _nilDefaultInt {
		ia.InstanceArrayProcessorCount = *v.(*int)
	}

	if v := m["instance_array_processor_core_mhz"]; v != nil && *v.(*int) != _nilDefaultInt {
		ia.InstanceArrayProcessorCoreMHZ = *v.(*int)
	}

	if v := m["instance_array_processor_core_count"]; v != nil && *v.(*int) != _nilDefaultInt {
		ia.InstanceArrayProcessorCoreCount = *v.(*int)
	}

	if v := m["instance_array_disk_count"]; v != nil && *v.(*int) != _nilDefaultInt {
		ia.InstanceArrayDiskCount = *v.(*int)
	}

	if v := m["instance_array_disk_size_mbytes"]; v != nil && *v.(*int) != _nilDefaultInt {
		ia.InstanceArrayDiskSizeMBytes = *v.(*int)
	}

	if v := m["instance_array_boot_method"]; v != nil && *v.(*string) != _nilDefaultStr {
		ia.InstanceArrayBootMethod = *v.(*string)
	}

	if v := m["instance_array_firewall_not_managed"]; v != nil {
		ia.InstanceArrayFirewallManaged = !(*v.(*bool))
	}

	if v := m["volume_template_id"]; v != nil && *v.(*int) != _nilDefaultInt {
		ia.VolumeTemplateID = *v.(*int)
	}

	return &ia
}

func argsToInstanceArrayOperation(m map[string]interface{}, iao *metalcloud.InstanceArrayOperation) {

	if v := m["instance_array_instance_count"]; v != nil && *v.(*int) != _nilDefaultInt {
		iao.InstanceArrayInstanceCount = *v.(*int)
	}

	if v := m["instance_array_label"]; v != nil && *v.(*string) != _nilDefaultStr {
		iao.InstanceArrayLabel = *v.(*string)
	}

	if v := m["instance_array_ram_gbytes"]; v != nil && *v.(*int) != _nilDefaultInt {
		iao.InstanceArrayRAMGbytes = *v.(*int)
	}

	if v := m["instance_array_processor_count"]; v != nil && *v.(*int) != _nilDefaultInt {
		iao.InstanceArrayProcessorCount = *v.(*int)
	}

	if v := m["instance_array_processor_core_mhz"]; v != nil && *v.(*int) != _nilDefaultInt {
		iao.InstanceArrayProcessorCoreMHZ = *v.(*int)
	}

	if v := m["instance_array_processor_core_count"]; v != nil && *v.(*int) != _nilDefaultInt {
		iao.InstanceArrayProcessorCoreCount = *v.(*int)
	}

	if v := m["instance_array_disk_count"]; v != nil && *v.(*int) != _nilDefaultInt {
		iao.InstanceArrayDiskCount = *v.(*int)
	}

	if v := m["instance_array_disk_size_mbytes"]; v != nil && *v.(*int) != _nilDefaultInt {
		iao.InstanceArrayDiskSizeMBytes = *v.(*int)
	}

	if v := m["instance_array_boot_method"]; v != nil && *v.(*string) != _nilDefaultStr {
		iao.InstanceArrayBootMethod = *v.(*string)
	}

	if v := m["instance_array_firewall_not_managed"]; v != nil {
		iao.InstanceArrayFirewallManaged = !*v.(*bool)
	}

	if v := m["volume_template_id"]; v != nil && *v.(*int) != _nilDefaultInt {
		iao.VolumeTemplateID = *v.(*int)
	}
}

func copyInstanceArrayToOperation(ia metalcloud.InstanceArray, iao *metalcloud.InstanceArrayOperation) {

	iao.InstanceArrayID = ia.InstanceArrayID
	iao.InstanceArrayLabel = ia.InstanceArrayLabel
	iao.InstanceArrayBootMethod = ia.InstanceArrayBootMethod
	iao.InstanceArrayInstanceCount = ia.InstanceArrayInstanceCount
	iao.InstanceArrayRAMGbytes = ia.InstanceArrayRAMGbytes
	iao.InstanceArrayProcessorCount = ia.InstanceArrayProcessorCount
	iao.InstanceArrayProcessorCoreMHZ = ia.InstanceArrayProcessorCoreMHZ
	iao.InstanceArrayDiskCount = ia.InstanceArrayDiskCount
	iao.InstanceArrayDiskSizeMBytes = ia.InstanceArrayDiskSizeMBytes
	iao.InstanceArrayDiskTypes = ia.InstanceArrayDiskTypes
	iao.ClusterID = ia.ClusterID
	iao.InstanceArrayFirewallManaged = ia.InstanceArrayFirewallManaged
	iao.InstanceArrayFirewallRules = ia.InstanceArrayFirewallRules
	iao.VolumeTemplateID = ia.VolumeTemplateID
}

func copyInstanceArrayInterfaceToOperation(i metalcloud.InstanceArrayInterface, io *metalcloud.InstanceArrayInterfaceOperation) {
	io.InstanceArrayInterfaceLAGGIndexes = i.InstanceArrayInterfaceLAGGIndexes
	io.InstanceArrayInterfaceIndex = i.InstanceArrayInterfaceIndex
	io.NetworkID = i.NetworkID
}
