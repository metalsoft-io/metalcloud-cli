package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	metalcloud "github.com/bigstepinc/metal-cloud-sdk-go"
)

//InstanceArrayCmds commands affecting instance arrays
var InstanceArrayCmds = []Command{

	Command{
		Description:  "Creates an instance array.",
		Subject:      "instance_array",
		AltSubject:   "ia",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("instance_array", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_id":                   c.FlagSet.Int("infra", 0, "(Required) Infrastrucure ID"),
				"instance_array_instance_count":       c.FlagSet.Int("instance_count", 1, "(Required) Instance count of this instance array"),
				"instance_array_label":                c.FlagSet.String("label", "", "InstanceArray's label"),
				"instance_array_ram_gbytes":           c.FlagSet.Int("ram", 1, "InstanceArray's minimum RAM (GB)"),
				"instance_array_processor_count":      c.FlagSet.Int("proc", 1, "InstanceArray's minimum processor count"),
				"instance_array_processor_core_mhz":   c.FlagSet.Int("proc_freq", 1000, "InstanceArray's minimum processor frequency (Mhz)"),
				"instance_array_processor_core_count": c.FlagSet.Int("proc_core_count", 1, "InstanceArray's minimum processor core count"),
				"instance_array_disk_count":           c.FlagSet.Int("disks", 1, "InstanceArray's number of local drives"),
				"instance_array_disk_size_mbytes":     c.FlagSet.Int("disk_size", 1, "InstanceArray's local disk sizes"),
				"instance_array_boot_method":          c.FlagSet.String("boot", "", "InstanceArray's boot type:'pxe_iscsi','local_drives'"),
				"instance_array_firewall_managed":     c.FlagSet.Bool("managed_fw", true, "InstanceArray's firewall management on or off"),
				"volume_template_id":                  c.FlagSet.Int("volume_template_id", 0, "InstanceArray's volume template when booting from for local drives"),
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
				"format":            c.FlagSet.String("format", "", "The output format. Supproted values are 'json','csv'. The default format is human readable."),
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
				"instance_array_id":                   c.FlagSet.Int("id", 0, "(Required) InstanceArray's id"),
				"instance_array_instance_count":       c.FlagSet.Int("instance_count", 0, "Instance count of this instance array"),
				"instance_array_label":                c.FlagSet.String("label", "", "(Required) InstanceArray's label"),
				"instance_array_ram_gbytes":           c.FlagSet.Int("ram", 1, "InstanceArray's minimum RAM (GB)"),
				"instance_array_processor_count":      c.FlagSet.Int("proc", 1, "InstanceArray's minimum processor count"),
				"instance_array_processor_core_mhz":   c.FlagSet.Int("proc_freq", 1000, "InstanceArray's minimum processor frequency (Mhz)"),
				"instance_array_processor_core_count": c.FlagSet.Int("proc_core_count", 1, "InstanceArray's minimum processor core count"),
				"instance_array_disk_count":           c.FlagSet.Int("disks", 1, "InstanceArray's number of local drives"),
				"instance_array_disk_size_mbytes":     c.FlagSet.Int("disk_size", 1, "InstanceArray's local disk sizes"),
				"instance_array_boot_method":          c.FlagSet.String("boot", "", "InstanceArray's boot type:'pxe_iscsi','local_drives'"),
				"instance_array_firewall_managed":     c.FlagSet.Bool("managed_fw", true, "InstanceArray's firewall management on or off"),
				"volume_template_id":                  c.FlagSet.Int("volume_template_id", 0, "InstanceArray's volume template when booting from for local drives"),
				"bSwapExistingInstancesHardware":      c.FlagSet.Bool("swap_existing_hardware", false, "If true, all the hardware of the Instance objects is swapped to match the new InstanceArray specifications"),
				"bKeepDetachingDrives":                c.FlagSet.Bool("keep_detaching_drives", true, "If false and the number of Instance objects is reduced, then the detaching Drive objects will be deleted. If it's set to true, the detaching Drive objects will not be deleted."),
				//		"objServerTypeMatches":                c.FlagSet.Int("server_type_id", 0, "If not null then the instances of this InstanceArray will be matched with the server configuration provided in the parameter (through the server_type_id property of a ServerType object). The InstanceArray properties detailing the minimum hardware configuration will be ignored."),
			}
		},
		ExecuteFunc: instanceArrayEditCmd,
	},
}

func instanceArrayCreateCmd(c *Command, client MetalCloudClient) (string, error) {

	infrastructureID := c.Arguments["infrastructure_id"]

	if infrastructureID == nil || *infrastructureID.(*int) == 0 {
		return "", fmt.Errorf("-infra <infrastructure_id> is required")
	}

	ia := argsToInstanceArray(c.Arguments)

	if ia.InstanceArrayLabel == "" {
		return "", fmt.Errorf("-label <instance_array_label> is required")
	}

	_, err := client.InstanceArrayCreate(*infrastructureID.(*int), *ia)

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
	if c.Arguments["bKeepDetachingDrives"] != nil {
		bKeepDetachingDrives = c.Arguments["bKeepDetachingDrives"].(*bool)
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

	autoConfirm := c.Arguments["autoconfirm"]

	confirm := false

	if autoConfirm == nil || *autoConfirm.(*bool) == false {
		fmt.Printf("Deleting instance array %s (%d) - from infrastructure %s (%d).  Are you sure? Type \"yes\" to continue:",
			retIA.InstanceArrayLabel, retIA.InstanceArrayID,
			retInfra.InfrastructureLabel, retInfra.InfrastructureID)
		reader := bufio.NewReader(os.Stdin)
		yes, _ := reader.ReadString('\n')

		if yes == "yes\n" {
			confirm = true
		}

	} else {
		confirm = true
	}

	if confirm {
		err := client.InstanceArrayDelete(*instanceArrayID.(*int))
		if err != nil {
			return "", err
		}
	}

	return "", nil
}

func argsToInstanceArray(m map[string]interface{}) *metalcloud.InstanceArray {
	ia := metalcloud.InstanceArray{}

	if v := m["instance_array_instance_count"]; v != nil {
		ia.InstanceArrayInstanceCount = *v.(*int)
	}

	if v := m["instance_array_label"]; v != nil {
		ia.InstanceArrayLabel = *v.(*string)
	}

	if v := m["instance_array_ram_gbytes"]; v != nil {
		ia.InstanceArrayRAMGbytes = *v.(*int)
	}

	if v := m["instance_array_processor_count"]; v != nil {
		ia.InstanceArrayProcessorCount = *v.(*int)
	}

	if v := m["instance_array_processor_core_mhz"]; v != nil {
		ia.InstanceArrayProcessorCoreMHZ = *v.(*int)
	}

	if v := m["instance_array_processor_core_count"]; v != nil {
		ia.InstanceArrayProcessorCoreCount = *v.(*int)
	}

	if v := m["instance_array_disk_count"]; v != nil {
		ia.InstanceArrayDiskCount = *v.(*int)
	}

	if v := m["instance_array_disk_size_mbytes"]; v != nil {
		ia.InstanceArrayDiskSizeMBytes = *v.(*int)
	}

	if v := m["instance_array_boot_method"]; v != nil {
		ia.InstanceArrayBootMethod = *v.(*string)
	}

	if v := m["instance_array_firewall_managed"]; v != nil {
		ia.InstanceArrayFirewallManaged = *v.(*bool)
	}

	if v := m["volume_template_id"]; v != nil {
		ia.VolumeTemplateID = *v.(*int)
	}

	return &ia
}

func argsToInstanceArrayOperation(m map[string]interface{}, iao *metalcloud.InstanceArrayOperation) {

	if v := m["instance_array_instance_count"]; v != nil {
		iao.InstanceArrayInstanceCount = *v.(*int)
	}

	if v := m["instance_array_label"]; v != nil {
		iao.InstanceArrayLabel = *v.(*string)
	}

	if v := m["instance_array_ram_gbytes"]; v != nil {
		iao.InstanceArrayRAMGbytes = *v.(*int)
	}

	if v := m["instance_array_processor_count"]; v != nil {
		iao.InstanceArrayProcessorCount = *v.(*int)
	}

	if v := m["instance_array_processor_core_mhz"]; v != nil {
		iao.InstanceArrayProcessorCoreMHZ = *v.(*int)
	}

	if v := m["instance_array_processor_core_count"]; v != nil {
		iao.InstanceArrayProcessorCoreCount = *v.(*int)
	}

	if v := m["instance_array_disk_count"]; v != nil {
		iao.InstanceArrayDiskCount = *v.(*int)
	}

	if v := m["instance_array_disk_size_mbytes"]; v != nil {
		iao.InstanceArrayDiskSizeMBytes = *v.(*int)
	}

	if v := m["instance_array_boot_method"]; v != nil {
		iao.InstanceArrayBootMethod = *v.(*string)
	}

	if v := m["instance_array_firewall_managed"]; v != nil {
		iao.InstanceArrayFirewallManaged = *v.(*bool)
	}

	if v := m["volume_template_id"]; v != nil {
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
