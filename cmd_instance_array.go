package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	metalcloud "github.com/bigstepinc/metal-cloud-sdk-go"
	interfaces "github.com/bigstepinc/metalcloud-cli/interfaces"
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
				"instance_array_disk_size_mbytes":     c.FlagSet.Int("disk_size", _nilDefaultInt, "InstanceArray's local disks' size in MB"),
				"instance_array_boot_method":          c.FlagSet.String("boot", _nilDefaultStr, "InstanceArray's boot type:'pxe_iscsi','local_drives'"),
				"instance_array_firewall_not_managed": c.FlagSet.Bool("firewall_management_disabled", false, "(Flag) If set InstanceArray's firewall management on or off"),
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
				"infrastructure_id_or_label": c.FlagSet.String("infra", _nilDefaultStr, "(Required) Infrastructure's id or label. Note that the 'label' this be ambiguous in certain situations."),
				"format":                     c.FlagSet.String("format", "", "The output format. Supported values are 'json','csv'. The default format is human readable."),
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
				"instance_array_id_or_label": c.FlagSet.String("id", _nilDefaultStr, "(Required) InstanceArray's id or label. Note that the label can be ambigous."),
				"autoconfirm":                c.FlagSet.Bool("autoconfirm", false, "If true it does not ask for confirmation anymore"),
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
				"instance_array_id_or_label":          c.FlagSet.String("id", _nilDefaultStr, "(Required) InstanceArray's id or label. Note that the label can be ambigous."),
				"instance_array_instance_count":       c.FlagSet.Int("instance_count", _nilDefaultInt, "Instance count of this instance array"),
				"instance_array_label":                c.FlagSet.String("label", _nilDefaultStr, "(Required) InstanceArray's label"),
				"instance_array_ram_gbytes":           c.FlagSet.Int("ram", _nilDefaultInt, "InstanceArray's minimum RAM (GB)"),
				"instance_array_processor_count":      c.FlagSet.Int("proc", _nilDefaultInt, "InstanceArray's minimum processor count"),
				"instance_array_processor_core_mhz":   c.FlagSet.Int("proc_freq", _nilDefaultInt, "InstanceArray's minimum processor frequency (Mhz)"),
				"instance_array_processor_core_count": c.FlagSet.Int("proc_core_count", _nilDefaultInt, "InstanceArray's minimum processor core count"),
				"instance_array_disk_count":           c.FlagSet.Int("disks", _nilDefaultInt, "InstanceArray's number of local drives"),
				"instance_array_disk_size_mbytes":     c.FlagSet.Int("disk_size", _nilDefaultInt, "InstanceArray's local disks' size in MB"),
				"instance_array_boot_method":          c.FlagSet.String("boot", _nilDefaultStr, "InstanceArray's boot type:'pxe_iscsi','local_drives'"),
				"instance_array_firewall_not_managed": c.FlagSet.Bool("firewall_management_disabled", false, "(Flag) If set InstanceArray's firewall management is off"),
				"volume_template_id":                  c.FlagSet.Int("template", _nilDefaultInt, "InstanceArray's volume template when booting from for local drives"),
				"bSwapExistingInstancesHardware":      c.FlagSet.Bool("swap_existing_hardware", false, "(Flag) If set all the hardware of the Instance objects is swapped to match the new InstanceArray specifications"),
				"no_bKeepDetachingDrives":             c.FlagSet.Bool("do_not_keep_detaching_drives", false, "(Flag) If set and the number of Instance objects is reduced, then the detaching Drive objects will be deleted. If it's set to true, the detaching Drive objects will not be deleted."),
			}
		},
		ExecuteFunc: instanceArrayEditCmd,
	},
	Command{
		Description:  "Get an instance array.",
		Subject:      "instance_array",
		AltSubject:   "ia",
		Predicate:    "get",
		AltPredicate: "show",
		FlagSet:      flag.NewFlagSet("get instance array", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"instance_array_id_or_label": c.FlagSet.String("id", _nilDefaultStr, "(Required) Instance array's id or label. Note that using the 'label' might be ambiguous in certain situations."),
				"show_credentials":           c.FlagSet.Bool("show_credentials", false, "(Flag) If set returns the instances' credentials"),
				"show_power_status":          c.FlagSet.Bool("show_power_status", false, "(Flag) If set returns the instances' power status"),
				"format":                     c.FlagSet.String("format", "", "The output format. Supported values are 'json','csv'. The default format is human readable."),
			}
		},
		ExecuteFunc: instanceArrayGetCmd,
	},
}

func instanceArrayCreateCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	infra, err := getInfrastructureFromCommand("infra", c, client)
	if err != nil {
		return "", err
	}

	ia := argsToInstanceArray(c.Arguments)

	if ia.InstanceArrayLabel == "" {
		return "", fmt.Errorf("-label <instance_array_label> is required")
	}

	retIA, err := client.InstanceArrayCreate(infra.InfrastructureID, *ia)
	if err != nil {
		return "", err
	}

	if c.Arguments["return_id"] != nil && *c.Arguments["return_id"].(*bool) {
		return fmt.Sprintf("%d", retIA.InstanceArrayID), nil
	}

	return "", err
}

func instanceArrayEditCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	retIA, err := getInstanceArrayFromCommand("id", c, client)
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

	_, err = client.InstanceArrayEdit(
		retIA.InstanceArrayID,
		*retIA.InstanceArrayOperation,
		bSwapExistingInstancesHardware,
		bKeepDetachingDrives,
		nil,
		nil)

	return "", err
}

func instanceArrayListCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	infra, err := getInfrastructureFromCommand("infra", c, client)
	if err != nil {
		return "", err
	}

	iaList, err := client.InstanceArrays(infra.InfrastructureID)
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
			FieldSize: 15,
		},
		SchemaField{
			FieldName: "STATUS",
			FieldType: TypeString,
			FieldSize: 10,
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
		if ia.InstanceArrayServiceStatus != "ordered" && ia.InstanceArrayOperation.InstanceArrayDeployType == "delete" && ia.InstanceArrayOperation.InstanceArrayDeployStatus == "not_started" {
			status = "marked for delete"
		}
		data = append(data, []interface{}{
			ia.InstanceArrayID,
			ia.InstanceArrayOperation.InstanceArrayLabel,
			status,
			ia.InstanceArrayOperation.InstanceArrayInstanceCount})
	}

	return renderTable("Instance Arrays", "", getStringParam(c.Arguments["format"]), data, schema)
}

func instanceArrayDeleteCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	retIA, err := getInstanceArrayFromCommand("id", c, client)
	if err != nil {
		return "", err
	}

	retInfra, err := client.InfrastructureGet(retIA.InfrastructureID)
	if err != nil {
		return "", err
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

		confirm, err = requestConfirmation(confirmationMessage)
		if err != nil {
			return "", err
		}
	}

	if !confirm {
		return "", fmt.Errorf("Operation not confirmed. Aborting")
	}

	err = client.InstanceArrayDelete(retIA.InstanceArrayID)

	return "", err
}

func instanceArrayGetCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	retIA, err := getInstanceArrayFromCommand("id", c, client)
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
			FieldName: "SUBDOMAIN",
			FieldType: TypeString,
			FieldSize: 10,
		},
		SchemaField{
			FieldName: "WAN_IP",
			FieldType: TypeString,
			FieldSize: 10,
		},
		SchemaField{
			FieldName: "DETAILS",
			FieldType: TypeString,
			FieldSize: 10,
		},
		SchemaField{
			FieldName: "STATUS",
			FieldType: TypeString,
			FieldSize: 5,
		},
	}

	data := [][]interface{}{}

	iList, err := client.InstanceArrayInstances(retIA.InstanceArrayID)
	if err != nil {
		return "", err
	}

	for _, i := range *iList {
		status := i.InstanceServiceStatus
		if i.InstanceServiceStatus != "ordered" && i.InstanceOperation.InstanceDeployType == "edit" && i.InstanceOperation.InstanceDeployStatus == "not_started" {
			status = "edited"
		}

		volumeTemplateName := ""
		if i.InstanceOperation.TemplateIDOrigin != 0 {
			vt, err := client.VolumeTemplateGet(i.InstanceOperation.TemplateIDOrigin)
			if err != nil {
				return "", err
			}
			volumeTemplateName = fmt.Sprintf("%s [#%d] ", vt.VolumeTemplateDisplayName, vt.VolumeTemplateID)
		}

		serverType := ""
		if i.ServerTypeID != 0 {
			st, err := client.ServerTypeGet(i.ServerTypeID)
			if err != nil {
				return "", err
			}
			serverType = st.ServerTypeDisplayName
		}

		details := fmt.Sprintf("%s (#%d) %s",
			serverType,
			i.ServerID,
			volumeTemplateName,
		)

		wanIP := ""
		for _, p := range i.InstanceInterfaces {
			if p.NetworkID != 0 {

				n, err := client.NetworkGet(p.NetworkID)
				if err != nil {
					return "", err
				}

				if n.NetworkType == "wan" {
					for _, iip := range p.InstanceInterfaceIPs {
						if iip.IPType == "ipv4" {
							wanIP = iip.IPHumanReadable
							break
						}
					}
				}
			}
			if wanIP != "" {
				break
			}
		}

		dataRow := []interface{}{
			i.InstanceID,
			i.InstanceSubdomainPermanent,
			wanIP,
			details,
			status,
		}

		if c.Arguments["show_credentials"] != nil && *c.Arguments["show_credentials"].(*bool) {
			credentials := ""
			schema = append(schema, SchemaField{
				FieldName: "CREDENTIALS",
				FieldType: TypeString,
				FieldSize: 5,
			})

			if v := i.InstanceCredentials.SSH; v != nil && v.Username != "" {
				credentials = fmt.Sprintf("SSH (%d) user: %s pass: %s", v.Port, v.Username, v.InitialPassword)
			}

			if v := i.InstanceCredentials.RDP; v != nil && v.Username != "" {
				credentials = fmt.Sprintf("RDP( %d) user: %s pass: %s", v.Port, v.Username, v.InitialPassword)
			}

			dataRow = append(dataRow, credentials)
		}

		if getBoolParam(c.Arguments["show_power_status"]) {
			powerStatus := ""
			schema = append(schema, SchemaField{
				FieldName: "POWER",
				FieldType: TypeString,
				FieldSize: 5,
			})

			pwr, err := client.InstanceServerPowerGet(i.InstanceID)
			if err != nil {
				powerStatus = err.Error()
			} else {
				powerStatus = *pwr
			}

			dataRow = append(dataRow, powerStatus)
		}

		data = append(data, dataRow)

	}

	infra, err := client.InfrastructureGet(retIA.InfrastructureID)
	if err != nil {
		return "", err
	}
	subtitle := fmt.Sprintf("Instances of instance array %s (#%d) of infrastructure %s (#%d):",
		retIA.InstanceArrayLabel,
		retIA.InstanceArrayID,
		infra.InfrastructureLabel,
		infra.InfrastructureID)

	return renderTable("Instances", subtitle, getStringParam(c.Arguments["format"]), data, schema)
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

func getInstanceArrayFromCommand(paramName string, c *Command, client interfaces.MetalCloudClient) (*metalcloud.InstanceArray, error) {

	m, err := getParam(c, "instance_array_id_or_label", paramName)
	if err != nil {
		return nil, err
	}

	id, label, isID := idOrLabel(m)

	if isID {
		return client.InstanceArrayGet(id)
	}

	return client.InstanceArrayGetByLabel(label)
}
