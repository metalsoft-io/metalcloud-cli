package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
	"github.com/metalsoft-io/tableformatter"
)

//instanceArrayCmds commands affecting instance arrays
var instanceArrayCmds = []Command{

	{
		Description:  "Creates an instance array.",
		Subject:      "instance-array",
		AltSubject:   "ia",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("instance-array", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_id_or_label":          c.FlagSet.String("infra", _nilDefaultStr, "(Required) Infrastructure's id or label. Note that the 'label' this be ambiguous in certain situations."),
				"instance_array_instance_count":       c.FlagSet.Int("instance-count", _nilDefaultInt, "(Required) Instance count of this instance array"),
				"instance_array_label":                c.FlagSet.String("label", _nilDefaultStr, "InstanceArray's label"),
				"server_type":                         c.FlagSet.String("server-type", _nilDefaultStr, "InstanceArray's server type."),
				"instance_array_ram_gbytes":           c.FlagSet.Int("ram", _nilDefaultInt, "InstanceArray's minimum RAM (GB)"),
				"instance_array_processor_count":      c.FlagSet.Int("proc", _nilDefaultInt, "InstanceArray's minimum processor count"),
				"instance_array_processor_core_mhz":   c.FlagSet.Int("proc-freq", _nilDefaultInt, "InstanceArray's minimum processor frequency (Mhz)"),
				"instance_array_processor_core_count": c.FlagSet.Int("proc-core-count", _nilDefaultInt, "InstanceArray's minimum processor core count"),
				"instance_array_disk_count":           c.FlagSet.Int("disks", _nilDefaultInt, "InstanceArray's number of local drives"),
				"instance_array_disk_size_mbytes":     c.FlagSet.Int("disk-size", _nilDefaultInt, "InstanceArray's local disks' size in MB"),
				"instance_array_boot_method":          c.FlagSet.String("boot", _nilDefaultStr, "InstanceArray's boot type:'pxe_iscsi','local_drives'"),
				"instance_array_firewall_not_managed": c.FlagSet.Bool("firewall-management-disabled", false, "(Flag) If set InstanceArray's firewall management on or off"),
				"volume_template_id_or_label":         c.FlagSet.String("local-install-template", _nilDefaultStr, "InstanceArray's volume template when booting from for local drives"),
				"da_volume_template":                  c.FlagSet.String("drive-array-template", _nilDefaultStr, "The attached DriveArray's  volume template when booting from iscsi drives"),
				"da_volume_disk_size":                 c.FlagSet.Int("drive-array-disk-size", _nilDefaultInt, "The attached DriveArray's  volume size (in MB) when booting from iscsi drives, If ommited the default size of the volume template will be used."),
				"return_id":                           c.FlagSet.Bool("return-id", false, "(Flag) If set will print the ID of the created Instance Array. Useful for automating tasks."),
			}
		},
		ExecuteFunc: instanceArrayCreateCmd,
	},
	{
		Description:  "Lists all instance arrays of an infrastructure.",
		Subject:      "instance-array",
		AltSubject:   "ia",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list instance_array", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_id_or_label": c.FlagSet.String("infra", _nilDefaultStr, "(Required) Infrastructure's id or label. Note that the 'label' this be ambiguous in certain situations."),
				"format":                     c.FlagSet.String("format", "", "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
			}
		},
		ExecuteFunc: instanceArrayListCmd,
	},
	{
		Description:  "Delete instance array.",
		Subject:      "instance-array",
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
	{
		Description:  "Edits an instance array.",
		Subject:      "instance-array",
		AltSubject:   "ia",
		Predicate:    "edit",
		AltPredicate: "update",
		FlagSet:      flag.NewFlagSet("instance_array", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"instance_array_id_or_label":          c.FlagSet.String("id", _nilDefaultStr, "(Required) InstanceArray's id or label. Note that the label can be ambigous."),
				"instance_array_instance_count":       c.FlagSet.Int("instance-count", _nilDefaultInt, "Instance count of this instance array"),
				"instance_array_label":                c.FlagSet.String("label", _nilDefaultStr, "(Required) InstanceArray's label"),
				"instance_array_ram_gbytes":           c.FlagSet.Int("ram", _nilDefaultInt, "InstanceArray's minimum RAM (GB)"),
				"instance_array_processor_count":      c.FlagSet.Int("proc", _nilDefaultInt, "InstanceArray's minimum processor count"),
				"instance_array_processor_core_mhz":   c.FlagSet.Int("proc-freq", _nilDefaultInt, "InstanceArray's minimum processor frequency (Mhz)"),
				"instance_array_processor_core_count": c.FlagSet.Int("proc-core-count", _nilDefaultInt, "InstanceArray's minimum processor core count"),
				"instance_array_disk_count":           c.FlagSet.Int("disks", _nilDefaultInt, "InstanceArray's number of local drives"),
				"instance_array_disk_size_mbytes":     c.FlagSet.Int("disk-size", _nilDefaultInt, "InstanceArray's local disks' size in MB"),
				"instance_array_boot_method":          c.FlagSet.String("boot", _nilDefaultStr, "InstanceArray's boot type:'pxe_iscsi','local_drives'"),
				"instance_array_firewall_not_managed": c.FlagSet.Bool("firewall-management-disabled", false, "(Flag) If set InstanceArray's firewall management is off"),
				"volume_template_id_or_label":         c.FlagSet.String("local-install-template", _nilDefaultStr, "InstanceArray's volume template when booting from for local drives"),
				"bSwapExistingInstancesHardware":      c.FlagSet.Bool("swap-existing-hardware", false, "(Flag) If set all the hardware of the Instance objects is swapped to match the new InstanceArray specifications"),
				"no_bKeepDetachingDrives":             c.FlagSet.Bool("do-not-keep-detaching-drives", false, "(Flag) If set and the number of Instance objects is reduced, then the detaching Drive objects will be deleted. If it's set to true, the detaching Drive objects will not be deleted."),
			}
		},
		ExecuteFunc: instanceArrayEditCmd,
	},
	{
		Description:  "Get instance array details.",
		Subject:      "instance-array",
		AltSubject:   "ia",
		Predicate:    "get",
		AltPredicate: "show",
		FlagSet:      flag.NewFlagSet("get instance array", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"instance_array_id_or_label": c.FlagSet.String("id", _nilDefaultStr, "(Required) Instance array's id or label. Note that using the 'label' might be ambiguous in certain situations."),
				"show_credentials":           c.FlagSet.Bool("show-credentials", false, "(Flag) If set returns the instances' credentials"),
				"show_power_status":          c.FlagSet.Bool("show-power-status", false, "(Flag) If set returns the instances' power status"),
				"show_iscsi_credentials":     c.FlagSet.Bool("show-iscsi-credentials", false, "(Flag) If set returns the instances' iscsi credentials"),
				"format":                     c.FlagSet.String("format", "", "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
			}
		},
		ExecuteFunc: instanceArrayGetCmd,
	},
}

func instanceArrayCreateCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	infra, err := getInfrastructureFromCommand("infra", c, client)
	if err != nil {
		return "", err
	}

	ia := argsToInstanceArray(c.Arguments, c, client)

	if ia.InstanceArrayLabel == "" {
		return "", fmt.Errorf("-label is required")
	}

	retIA, err := client.InstanceArrayCreate(infra.InfrastructureID, *ia)
	if err != nil {
		return "", err
	}

	if serverTypeLabel, ok := getStringParamOk(c.Arguments["server_type"]); ok {

		serverType, err := client.ServerTypeGetByLabel(serverTypeLabel)
		if err != nil {
			return "", err
		}

		stMatches := metalcloud.ServerTypeMatches{
			ServerTypes: map[int]metalcloud.ServerTypeMatch{
				serverType.ServerTypeID: {
					ServerCount: retIA.InstanceArrayInstanceCount,
				},
			},
		}
		retIA.InstanceArrayProcessorCoreCount = serverType.ServerProcessorCoreCount
		retIA.InstanceArrayProcessorCount = serverType.ServerProcessorCount
		retIA.InstanceArrayRAMGbytes = serverType.ServerRAMGbytes

		bFalse := false
		_, err = client.InstanceArrayEdit(retIA.InstanceArrayID, *retIA.InstanceArrayOperation, &bFalse, &bFalse, &stMatches, nil)
		if err != nil {
			return "", err
		}
	}

	if driveArrayVolumeTemplateLabel, ok := getStringParamOk(c.Arguments["da_volume_template"]); ok {
		volumeTemplate, err := client.VolumeTemplateGetByLabel(driveArrayVolumeTemplateLabel)
		if err != nil {
			return "", err
		}

		driveSize := getIntParam(c.Arguments["da_volume_disk_size"])
		if driveSize == 0 {
			driveSize = volumeTemplate.VolumeTemplateSizeMBytes
		}

		da := metalcloud.DriveArray{
			VolumeTemplateID:                  volumeTemplate.VolumeTemplateID,
			DriveSizeMBytesDefault:            driveSize,
			InstanceArrayID:                   retIA.InstanceArrayID,
			DriveArrayExpandWithInstanceArray: true,
		}
		_, err = client.DriveArrayCreate(retIA.InfrastructureID, da)
		if err != nil {
			return "", err
		}
	}

	if c.Arguments["return_id"] != nil && *c.Arguments["return_id"].(*bool) {
		return fmt.Sprintf("%d", retIA.InstanceArrayID), nil
	}

	return "", err
}

func instanceArrayEditCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	retIA, err := getInstanceArrayFromCommand("id", c, client)
	if err != nil {
		return "", err
	}

	argsToInstanceArrayOperation(c.Arguments, retIA.InstanceArrayOperation, c, client)

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

func instanceArrayListCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	infra, err := getInfrastructureFromCommand("infra", c, client)
	if err != nil {
		return "", err
	}

	iaList, err := client.InstanceArrays(infra.InfrastructureID)
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
			FieldSize: 15,
		},
		{
			FieldName: "STATUS",
			FieldType: tableformatter.TypeString,
			FieldSize: 10,
		},
		{
			FieldName: "INST_CNT",
			FieldType: tableformatter.TypeInt,
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

	tableformatter.TableSorter(schema).OrderBy(schema[0].FieldName).Sort(data)

	table := tableformatter.Table{
		Data:   data,
		Schema: schema,
	}
	return table.RenderTable("Instance Arrays", "", getStringParam(c.Arguments["format"]))
}

func instanceArrayDeleteCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	retIA, err := getInstanceArrayFromCommand("id", c, client)
	if err != nil {
		return "", err
	}

	retInfra, err := client.InfrastructureGet(retIA.InfrastructureID)
	if err != nil {
		return "", err
	}

	confirm := false

	if getBoolParam(c.Arguments["autoconfirm"]) {
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

func instanceArrayGetCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	retIA, err := getInstanceArrayFromCommand("id", c, client)
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
			FieldName: "SUBDOMAIN",
			FieldType: tableformatter.TypeString,
			FieldSize: 10,
		},
		{
			FieldName: "WAN_IP",
			FieldType: tableformatter.TypeString,
			FieldSize: 10,
		},
		{
			FieldName: "DETAILS",
			FieldType: tableformatter.TypeString,
			FieldSize: 10,
		},
		{
			FieldName: "STATUS",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
	}

	if getBoolParam(c.Arguments["show_credentials"]) {

		schema = append(schema, tableformatter.SchemaField{
			FieldName: "CREDENTIALS",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		})
	}

	if getBoolParam(c.Arguments["show_power_status"]) {

		schema = append(schema, tableformatter.SchemaField{
			FieldName: "POWER",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		})
	}

	if getBoolParam(c.Arguments["show_iscsi_credentials"]) {

		schema = append(schema, tableformatter.SchemaField{
			FieldName: "ISCSI",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		})
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

		if getBoolParam(c.Arguments["show_credentials"]) {
			credentials := ""

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

			pwr, err := client.InstanceServerPowerGet(i.InstanceID)
			if err != nil {
				powerStatus = err.Error()
			} else {
				powerStatus = *pwr
			}

			dataRow = append(dataRow, powerStatus)
		}

		if getBoolParam(c.Arguments["show_iscsi_credentials"]) {
			iscsiCreds := ""
			if v := i.InstanceCredentials.ISCSI; v != nil {
				iscsiCreds = fmt.Sprintf("Initiator IQN: %s Username: %s Password: %s ", v.InitiatorIQN, v.Username, v.Password)
			}
			dataRow = append(dataRow, iscsiCreds)
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

	tableformatter.TableSorter(schema).OrderBy(schema[0].FieldName).Sort(data)

	table := tableformatter.Table{
		Data:   data,
		Schema: schema,
	}
	return table.RenderTable("Instances", subtitle, getStringParam(c.Arguments["format"]))
}

func argsToInstanceArray(m map[string]interface{}, c *Command, client metalcloud.MetalCloudClient) *metalcloud.InstanceArray {
	ia := metalcloud.InstanceArray{}

	if v, ok := getIntParamOk(m["instance_array_instance_count"]); ok {
		ia.InstanceArrayInstanceCount = v
	}

	if v, ok := getStringParamOk(m["instance_array_label"]); ok {
		ia.InstanceArrayLabel = v
	}

	if v, ok := getIntParamOk(m["instance_array_ram_gbytes"]); ok {
		ia.InstanceArrayRAMGbytes = v
	}

	if v, ok := getIntParamOk(m["instance_array_processor_count"]); ok {
		ia.InstanceArrayProcessorCount = v
	}

	if v, ok := getIntParamOk(m["instance_array_processor_core_mhz"]); ok {
		ia.InstanceArrayProcessorCoreMHZ = v
	}

	if v, ok := getIntParamOk(m["instance_array_processor_core_count"]); ok {
		ia.InstanceArrayProcessorCoreCount = v
	}

	if v, ok := getIntParamOk(m["instance_array_disk_count"]); ok {
		ia.InstanceArrayDiskCount = v
	}

	if v, ok := getIntParamOk(m["instance_array_disk_size_mbytes"]); ok {
		ia.InstanceArrayDiskSizeMBytes = v
	}

	if v, ok := getStringParamOk(m["instance_array_boot_method"]); ok {
		ia.InstanceArrayBootMethod = v
	}

	if v, ok := getBoolParamOk(m["instance_array_firewall_not_managed"]); ok {
		ia.InstanceArrayFirewallManaged = !v
	}

	if v, ok := getStringParamOk(c.Arguments["volume_template_id_or_label"]); ok {
		vtID, err := getIDOrDo(v, func(label string) (int, error) {
			vt, err := client.VolumeTemplateGetByLabel(label)
			if err != nil {
				return 0, err
			}
			return vt.VolumeTemplateID, nil
		},
		)
		if err != nil {
			ia.VolumeTemplateID = 0
		}
		ia.VolumeTemplateID = vtID
	}

	return &ia
}

func argsToInstanceArrayOperation(m map[string]interface{}, iao *metalcloud.InstanceArrayOperation, c *Command, client metalcloud.MetalCloudClient) {
	if v, ok := getIntParamOk(m["instance_array_instance_count"]); ok {
		iao.InstanceArrayInstanceCount = v
	}

	if v, ok := getStringParamOk(m["instance_array_label"]); ok {
		iao.InstanceArrayLabel = v
	}

	if v, ok := getIntParamOk(m["instance_array_ram_gbytes"]); ok {
		iao.InstanceArrayRAMGbytes = v
	}

	if v, ok := getIntParamOk(m["instance_array_processor_count"]); ok {
		iao.InstanceArrayProcessorCount = v
	}

	if v, ok := getIntParamOk(m["instance_array_processor_core_mhz"]); ok {
		iao.InstanceArrayProcessorCoreMHZ = v
	}

	if v, ok := getIntParamOk(m["instance_array_processor_core_count"]); ok {
		iao.InstanceArrayProcessorCoreCount = v
	}

	if v, ok := getIntParamOk(m["instance_array_disk_count"]); ok {
		iao.InstanceArrayDiskCount = v
	}

	if v, ok := getIntParamOk(m["instance_array_disk_size_mbytes"]); ok {
		iao.InstanceArrayDiskSizeMBytes = v
	}

	if v, ok := getStringParamOk(m["instance_array_boot_method"]); ok {
		iao.InstanceArrayBootMethod = v
	}

	if v, ok := getBoolParamOk(m["instance_array_firewall_not_managed"]); ok {
		iao.InstanceArrayFirewallManaged = !v
	}

	if v, ok := getStringParamOk(c.Arguments["volume_template_id_or_label"]); ok {
		vtID, err := getIDOrDo(v, func(label string) (int, error) {
			vt, err := client.VolumeTemplateGetByLabel(label)
			if err != nil {
				return 0, err
			}
			return vt.VolumeTemplateID, nil
		},
		)
		if err != nil {
			iao.VolumeTemplateID = 0
		}
		iao.VolumeTemplateID = vtID
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

func getInstanceArrayFromCommand(paramName string, c *Command, client metalcloud.MetalCloudClient) (*metalcloud.InstanceArray, error) {

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
