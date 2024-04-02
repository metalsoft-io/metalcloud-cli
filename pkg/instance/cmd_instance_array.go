package instance

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
	"github.com/metalsoft-io/metalcloud-cli/internal/colors"
	"github.com/metalsoft-io/metalcloud-cli/internal/command"
	"github.com/metalsoft-io/metalcloud-cli/internal/configuration"
	"github.com/metalsoft-io/tableformatter"
)

var InstanceArrayCmds = []command.Command{
	{
		Description:  "Creates an instance array.",
		Subject:      "instance-array",
		AltSubject:   "ia",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("instance-array", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_id_or_label":    c.FlagSet.String("infra", command.NilDefaultStr, colors.Red("(Required)")+" Infrastructure's id or label. Note that the 'label' this be ambiguous in certain situations."),
				"instance_array_instance_count": c.FlagSet.Int("instance-count", 1, " Instance count of this instance array"),
				"instance_array_label":          c.FlagSet.String("label", command.NilDefaultStr, "InstanceArray's label"),
				"server_type":                   c.FlagSet.String("server-type", command.NilDefaultStr, "InstanceArray's server type."),
				//	"instance_array_ram_gbytes":           c.FlagSet.Int("ram", command.NilDefaultInt, "InstanceArray's minimum RAM (GB)"),
				//	"instance_array_processor_count":      c.FlagSet.Int("proc", command.NilDefaultInt, "InstanceArray's minimum processor count"),
				//	"instance_array_processor_core_mhz":   c.FlagSet.Int("proc-freq", command.NilDefaultInt, "InstanceArray's minimum processor frequency (Mhz)"),
				//	"instance_array_processor_core_count": c.FlagSet.Int("proc-core-count", command.NilDefaultInt, "InstanceArray's minimum processor core count"),
				//	"instance_array_disk_count":           c.FlagSet.Int("disks", command.NilDefaultInt, "InstanceArray's number of local drives"),
				//	"instance_array_disk_size_mbytes":     c.FlagSet.Int("disk-size", command.NilDefaultInt, "InstanceArray's local disks' size in MB"),
				"instance_array_boot_method":          c.FlagSet.String("boot", command.NilDefaultStr, "InstanceArray's boot type:'pxe_iscsi','local_drives'"),
				"instance_array_firewall_not_managed": c.FlagSet.Bool("firewall-management-disabled", false, colors.Green("(Flag)")+" If set InstanceArray's firewall management on or off"),
				"volume_template_id_or_label":         c.FlagSet.String("local-install-template", command.NilDefaultStr, "InstanceArray's volume template when booting from for local drives"),
				"da_volume_template":                  c.FlagSet.String("drive-array-template", command.NilDefaultStr, "The attached DriveArray's  volume template when booting from iscsi drives"),
				"da_volume_disk_size":                 c.FlagSet.Int("drive-array-disk-size", command.NilDefaultInt, "The attached DriveArray's  volume size (in MB) when booting from iscsi drives, If ommited the default size of the volume template will be used."),
				"custom_variables":                    c.FlagSet.String("custom-variables", command.NilDefaultStr, "Comma separated list of custom variables such as 'var1=value,var2=value'. If special characters need to be set use urlencode and pass the encoded string"),
				"return_id":                           c.FlagSet.Bool("return-id", false, colors.Green("(Flag)")+" If set will print the ID of the created Instance Array. Useful for automating tasks."),
			}
		},
		ExecuteFunc: instanceArrayCreateCmd,
		Endpoint:    configuration.UserEndpoint,
	},
	{
		Description:  "Lists all instance arrays of an infrastructure.",
		Subject:      "instance-array",
		AltSubject:   "ia",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list instance_array", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_id_or_label": c.FlagSet.String("infra", command.NilDefaultStr, colors.Red("(Required)")+" Infrastructure's id or label. Note that the 'label' this be ambiguous in certain situations."),
				"format":                     c.FlagSet.String("format", "", "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
			}
		},
		ExecuteFunc: instanceArrayListCmd,
		Endpoint:    configuration.UserEndpoint,
	},
	{
		Description:  "Lists all instances of instance array.",
		Subject:      "instance-array",
		AltSubject:   "ia",
		Predicate:    "instances-list",
		AltPredicate: "instances-ls",
		FlagSet:      flag.NewFlagSet("instances-list instance_array", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"instance_array_id_or_label": c.FlagSet.String("id", command.NilDefaultStr, colors.Red("(Required)")+" InstanceArray's id or label. Note that the label can be ambigous."),
				"format":                     c.FlagSet.String("format", "", "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
			}
		},
		ExecuteFunc: instanceArrayInstancesListCmd,
		Endpoint:    configuration.UserEndpoint,
	},
	{
		Description:  "Delete instance array.",
		Subject:      "instance-array",
		AltSubject:   "ia",
		Predicate:    "delete",
		AltPredicate: "rm",
		FlagSet:      flag.NewFlagSet("list instance_array", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"instance_array_id_or_label": c.FlagSet.String("id", command.NilDefaultStr, colors.Red("(Required)")+" InstanceArray's id or label. Note that the label can be ambigous."),
				"autoconfirm":                c.FlagSet.Bool("autoconfirm", false, colors.Green("(Flag)")+" If set it will assume action is confirmed"),
			}
		},
		ExecuteFunc: instanceArrayDeleteCmd,
		Endpoint:    configuration.UserEndpoint,
	},
	{
		Description:  "Edits an instance array.",
		Subject:      "instance-array",
		AltSubject:   "ia",
		Predicate:    "edit",
		AltPredicate: "update",
		FlagSet:      flag.NewFlagSet("instance_array", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"instance_array_id_or_label":          c.FlagSet.String("id", command.NilDefaultStr, colors.Red("(Required)")+" InstanceArray's id or label. Note that the label can be ambigous."),
				"instance_array_instance_count":       c.FlagSet.Int("instance-count", command.NilDefaultInt, "Instance count of this instance array"),
				"instance_array_label":                c.FlagSet.String("label", command.NilDefaultStr, colors.Red("(Required)")+" InstanceArray's label"),
				"instance_array_ram_gbytes":           c.FlagSet.Int("ram", command.NilDefaultInt, "InstanceArray's minimum RAM (GB)"),
				"instance_array_processor_count":      c.FlagSet.Int("proc", command.NilDefaultInt, "InstanceArray's minimum processor count"),
				"instance_array_processor_core_mhz":   c.FlagSet.Int("proc-freq", command.NilDefaultInt, "InstanceArray's minimum processor frequency (Mhz)"),
				"instance_array_processor_core_count": c.FlagSet.Int("proc-core-count", command.NilDefaultInt, "InstanceArray's minimum processor core count"),
				"instance_array_disk_count":           c.FlagSet.Int("disks", command.NilDefaultInt, "InstanceArray's number of local drives"),
				"instance_array_disk_size_mbytes":     c.FlagSet.Int("disk-size", command.NilDefaultInt, "InstanceArray's local disks' size in MB"),
				"instance_array_boot_method":          c.FlagSet.String("boot", command.NilDefaultStr, "InstanceArray's boot type:'pxe_iscsi','local_drives'"),
				"instance_array_firewall_not_managed": c.FlagSet.Bool("firewall-management-disabled", false, colors.Green("(Flag)")+" If set InstanceArray's firewall management is off"),
				"volume_template_id_or_label":         c.FlagSet.String("local-install-template", command.NilDefaultStr, "InstanceArray's volume template when booting from for local drives"),
				"bSwapExistingInstancesHardware":      c.FlagSet.Bool("swap-existing-hardware", false, colors.Green("(Flag)")+" If set all the hardware of the Instance objects is swapped to match the new InstanceArray specifications"),
				"custom_variables":                    c.FlagSet.String("custom-variables", command.NilDefaultStr, "Comma separated list of custom variables such as 'var1=value,var2=value'. If special characters need to be set use urlencode and pass the encoded string"),
				"no_bKeepDetachingDrives":             c.FlagSet.Bool("do-not-keep-detaching-drives", false, colors.Green("(Flag)")+" If set and the number of Instance objects is reduced, then the detaching Drive objects will be deleted. If it's set to true, the detaching Drive objects will not be deleted."),
			}
		},
		ExecuteFunc: instanceArrayEditCmd,
		Endpoint:    configuration.UserEndpoint,
	},
	{
		Description:  "Get instance array details.",
		Subject:      "instance-array",
		AltSubject:   "ia",
		Predicate:    "get",
		AltPredicate: "show",
		FlagSet:      flag.NewFlagSet("get instance array", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"instance_array_id_or_label": c.FlagSet.String("id", command.NilDefaultStr, colors.Red("(Required)")+" Instance array's id or label. Note that using the 'label' might be ambiguous in certain situations."),
				"show_credentials":           c.FlagSet.Bool("show-credentials", false, colors.Green("(Flag)")+" If set returns the instances' credentials"),
				"show_power_status":          c.FlagSet.Bool("show-power-status", false, colors.Green("(Flag)")+" If set returns the instances' power status"),
				"show_iscsi_credentials":     c.FlagSet.Bool("show-iscsi-credentials", false, colors.Green("(Flag)")+" If set returns the instances' iscsi credentials"),
				"show_custom_variables":      c.FlagSet.Bool("show-custom-variables", false, colors.Green("(Flag)")+" If set returns the instances' custom variables"),
				"format":                     c.FlagSet.String("format", "", "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
			}
		},
		ExecuteFunc: instanceArrayGetCmd,
		Endpoint:    configuration.UserEndpoint,
	},
}

func instanceArrayCreateCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	infra, err := command.GetInfrastructureFromCommand("infra", c, client)
	if err != nil {
		return "", err
	}

	ia, err := argsToInstanceArray(c.Arguments, c, client)
	if err != nil {
		return "", err
	}

	if ia.InstanceArrayLabel == "" {
		return "", fmt.Errorf("-label is required")
	}

	retIA, err := client.InstanceArrayCreate(infra.InfrastructureID, *ia)
	if err != nil {
		return "", err
	}

	if serverTypeLabel, ok := command.GetStringParamOk(c.Arguments["server_type"]); ok {

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

	if driveArrayVolumeTemplateLabel, ok := command.GetStringParamOk(c.Arguments["da_volume_template"]); ok {
		volumeTemplate, err := client.VolumeTemplateGetByLabel(driveArrayVolumeTemplateLabel)
		if err != nil {
			return "", err
		}

		driveSize := command.GetIntParam(c.Arguments["da_volume_disk_size"])
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

func instanceArrayEditCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	retIA, err := command.GetInstanceArrayFromCommand("id", c, client)
	if err != nil {
		return "", err
	}

	err = argsToInstanceArrayOperation(c.Arguments, retIA.InstanceArrayOperation, c, client)
	if err != nil {
		return "", err
	}

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

func instanceArrayListCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	infra, err := command.GetInfrastructureFromCommand("infra", c, client)
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
	return table.RenderTable("Instance Arrays", "", command.GetStringParam(c.Arguments["format"]))
}

func instanceArrayDeleteCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	retIA, err := command.GetInstanceArrayFromCommand("id", c, client)
	if err != nil {
		return "", err
	}

	retInfra, err := client.InfrastructureGet(retIA.InfrastructureID)
	if err != nil {
		return "", err
	}

	confirm := false

	if command.GetBoolParam(c.Arguments["autoconfirm"]) {
		confirm = true
	} else {

		confirmationMessage := fmt.Sprintf("Deleting instance array %s (%d) - from infrastructure %s (%d).  Are you sure? Type \"yes\" to continue:",
			retIA.InstanceArrayLabel, retIA.InstanceArrayID,
			retInfra.InfrastructureLabel, retInfra.InfrastructureID)

		//this is simply so that we don't output a text on the command line under go test
		if strings.HasSuffix(os.Args[0], ".test") {
			confirmationMessage = ""
		}

		confirm, err = command.RequestConfirmation(confirmationMessage)
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

func instanceArrayGetNetworkAttachments(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	retIA, err := command.GetInstanceArrayFromCommand("id", c, client)
	if err != nil {
		return "", err
	}

	if &retIA == nil {
		return "", fmt.Errorf("instance array should not be nil")
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

	return tableNetworkAttachments.RenderTable("", subtitleNetworkAttachmentsRender, command.GetStringParam(c.Arguments["format"]))
}

func instanceArrayGetCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	retIA, err := command.GetInstanceArrayFromCommand("id", c, client)
	if err != nil {
		return "", err
	}

	infra, err := client.InfrastructureGet(retIA.InfrastructureID)
	if err != nil {
		return "", err
	}

	schema := []tableformatter.SchemaField{
		{
			FieldName: "ID",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "STATUS",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "LABEL",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "INFRASTRUCTURE",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "DETAILS",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
	}

	if command.GetBoolParam(c.Arguments["show_custom_variables"]) {
		schema = append(schema, tableformatter.SchemaField{
			FieldName: "CUSTOM_VARS",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		})
	}

	status := colors.Green(retIA.InstanceArrayServiceStatus)
	if retIA.InstanceArrayServiceStatus != "ordered" && retIA.InstanceArrayOperation.InstanceArrayDeployType == "edit" && retIA.InstanceArrayOperation.InstanceArrayDeployStatus == "not_started" {
		status = colors.Blue("edited")
	}

	volumeTemplateName := ""
	if retIA.InstanceArrayOperation.VolumeTemplateID != 0 {
		vt, err := client.VolumeTemplateGet(retIA.InstanceArrayOperation.VolumeTemplateID)
		if err != nil {
			return "", err
		}
		volumeTemplateName = fmt.Sprintf("%s [#%d] ", vt.VolumeTemplateDisplayName, vt.VolumeTemplateID)
	}

	fwMgmtDisabled := ""
	if !retIA.InstanceArrayFirewallManaged {
		fwMgmtDisabled = " fw mgmt disabled"
	}

	details := fmt.Sprintf("%d instances (%d RAM, %d cores, %d disks %s %s%s)",
		retIA.InstanceArrayOperation.InstanceArrayInstanceCount,
		retIA.InstanceArrayOperation.InstanceArrayRAMGbytes,
		retIA.InstanceArrayOperation.InstanceArrayProcessorCount*retIA.InstanceArrayProcessorCoreCount,
		retIA.InstanceArrayOperation.InstanceArrayDiskCount,
		retIA.InstanceArrayOperation.InstanceArrayBootMethod,
		volumeTemplateName,
		fwMgmtDisabled,
	)

	custom_vars := ""
	if command.GetBoolParam(c.Arguments["show_custom_variables"]) {
		custom_vars = command.GetKeyValueStringFromMap(retIA.InstanceArrayCustomVariables)
	}

	data := [][]interface{}{
		{
			"#" + strconv.Itoa(retIA.InstanceArrayID),
			status,
			retIA.InstanceArrayLabel,
			infra.InfrastructureLabel + " #" + strconv.Itoa(retIA.InfrastructureID),
			details,
			custom_vars,
		},
	}

	table := tableformatter.Table{
		Data:   data,
		Schema: schema,
	}

	return table.RenderTable("", "", command.GetStringParam(c.Arguments["format"]))
}

func instanceArrayInstancesListCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	retIA, err := command.GetInstanceArrayFromCommand("id", c, client)
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

	if command.GetBoolParam(c.Arguments["show_credentials"]) {

		schema = append(schema, tableformatter.SchemaField{
			FieldName: "CREDENTIALS",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		})
	}

	if command.GetBoolParam(c.Arguments["show_power_status"]) {

		schema = append(schema, tableformatter.SchemaField{
			FieldName: "POWER",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		})
	}

	if command.GetBoolParam(c.Arguments["show_iscsi_credentials"]) {

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

		if command.GetBoolParam(c.Arguments["show_credentials"]) {
			credentials := ""

			if v := i.InstanceCredentials.SSH; v != nil && v.Username != "" {
				credentials = fmt.Sprintf("SSH (%d) user: %s pass: %s", v.Port, v.Username, v.InitialPassword)
			}

			if v := i.InstanceCredentials.RDP; v != nil && v.Username != "" {
				credentials = fmt.Sprintf("RDP( %d) user: %s pass: %s", v.Port, v.Username, v.InitialPassword)
			}

			dataRow = append(dataRow, credentials)
		}

		if command.GetBoolParam(c.Arguments["show_power_status"]) {
			powerStatus := ""

			pwr, err := client.InstanceServerPowerGet(i.InstanceID)
			if err != nil {
				powerStatus = err.Error()
			} else {
				powerStatus = *pwr
			}

			dataRow = append(dataRow, powerStatus)
		}

		if command.GetBoolParam(c.Arguments["show_iscsi_credentials"]) {
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

	return table.RenderTable("Instances", subtitle, command.GetStringParam(c.Arguments["format"]))
}

func argsToInstanceArray(m map[string]interface{}, c *command.Command, client metalcloud.MetalCloudClient) (*metalcloud.InstanceArray, error) {
	ia := metalcloud.InstanceArray{}

	if v, ok := command.GetIntParamOk(m["instance_array_instance_count"]); ok {
		ia.InstanceArrayInstanceCount = v
	}

	if v, ok := command.GetStringParamOk(m["instance_array_label"]); ok {
		ia.InstanceArrayLabel = v
	}

	if v, ok := command.GetIntParamOk(m["instance_array_ram_gbytes"]); ok {
		ia.InstanceArrayRAMGbytes = v
	}

	if v, ok := command.GetIntParamOk(m["instance_array_processor_count"]); ok {
		ia.InstanceArrayProcessorCount = v
	}

	if v, ok := command.GetIntParamOk(m["instance_array_processor_core_mhz"]); ok {
		ia.InstanceArrayProcessorCoreMHZ = v
	}

	if v, ok := command.GetIntParamOk(m["instance_array_processor_core_count"]); ok {
		ia.InstanceArrayProcessorCoreCount = v
	}

	if v, ok := command.GetIntParamOk(m["instance_array_disk_count"]); ok {
		ia.InstanceArrayDiskCount = v
	}

	if v, ok := command.GetIntParamOk(m["instance_array_disk_size_mbytes"]); ok {
		ia.InstanceArrayDiskSizeMBytes = v
	}

	if v, ok := command.GetStringParamOk(m["instance_array_boot_method"]); ok {
		ia.InstanceArrayBootMethod = v
	}

	if v, ok := command.GetBoolParamOk(m["instance_array_firewall_not_managed"]); ok {
		ia.InstanceArrayFirewallManaged = !v
	}

	if v, ok := command.GetStringParamOk(c.Arguments["volume_template_id_or_label"]); ok {
		vtID, err := command.GetIDOrDo(v, func(label string) (int, error) {
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

	if v, ok := command.GetStringParamOk(c.Arguments["custom_variables"]); ok {

		m, err := command.GetKeyValueMapFromString(v)
		if err != nil {
			return nil, err
		}

		ia.InstanceArrayCustomVariables = m
	}

	return &ia, nil
}

func argsToInstanceArrayOperation(m map[string]interface{}, iao *metalcloud.InstanceArrayOperation, c *command.Command, client metalcloud.MetalCloudClient) error {
	if v, ok := command.GetIntParamOk(m["instance_array_instance_count"]); ok {
		iao.InstanceArrayInstanceCount = v
	}

	if v, ok := command.GetStringParamOk(m["instance_array_label"]); ok {
		iao.InstanceArrayLabel = v
	}

	if v, ok := command.GetIntParamOk(m["instance_array_ram_gbytes"]); ok {
		iao.InstanceArrayRAMGbytes = v
	}

	if v, ok := command.GetIntParamOk(m["instance_array_processor_count"]); ok {
		iao.InstanceArrayProcessorCount = v
	}

	if v, ok := command.GetIntParamOk(m["instance_array_processor_core_mhz"]); ok {
		iao.InstanceArrayProcessorCoreMHZ = v
	}

	if v, ok := command.GetIntParamOk(m["instance_array_processor_core_count"]); ok {
		iao.InstanceArrayProcessorCoreCount = v
	}

	if v, ok := command.GetIntParamOk(m["instance_array_disk_count"]); ok {
		iao.InstanceArrayDiskCount = v
	}

	if v, ok := command.GetIntParamOk(m["instance_array_disk_size_mbytes"]); ok {
		iao.InstanceArrayDiskSizeMBytes = v
	}

	if v, ok := command.GetStringParamOk(m["instance_array_boot_method"]); ok {
		iao.InstanceArrayBootMethod = v
	}

	if v, ok := command.GetBoolParamOk(m["instance_array_firewall_not_managed"]); ok {
		iao.InstanceArrayFirewallManaged = !v
	}

	if v, ok := command.GetStringParamOk(c.Arguments["volume_template_id_or_label"]); ok {
		vtID, err := command.GetIDOrDo(v, func(label string) (int, error) {
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

	if v, ok := command.GetStringParamOk(c.Arguments["custom_variables"]); ok {

		m, err := command.GetKeyValueMapFromString(v)
		if err != nil {
			return err
		}

		iao.InstanceArrayCustomVariables = m
	}
	return nil
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
