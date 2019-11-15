package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	metalcloud "github.com/bigstepinc/metal-cloud-sdk-go"
)

//infrastructureCmds commands affecting infrastructures
var infrastructureCmds = []Command{

	Command{
		Description:  "Creates an infrastructure.",
		Subject:      "infrastructure",
		AltSubject:   "infra",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("create infrastructure", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_label": c.FlagSet.String("label", "", "(Required) Infrastructure's label"),
				"datacenter":           c.FlagSet.String("dc", GetDatacenter(), "(Required) Infrastructure datacenter"),
				"return_id":            c.FlagSet.Bool("return_id", false, "(Flag) If set will print the ID of the created infrastructure. Useful for automating tasks."),
			}
		},
		ExecuteFunc: infrastructureCreateCmd,
	},
	Command{
		Description:  "Lists all infrastructures.",
		Subject:      "infrastructure",
		AltSubject:   "infra",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list infrastructure", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"format": c.FlagSet.String("format", "", "The output format. Supported values are 'json','csv'. The default format is human readable."),
			}
		},
		ExecuteFunc: infrastructureListCmd,
	},
	Command{
		Description:  "Delete an infrastructure.",
		Subject:      "infrastructure",
		AltSubject:   "infra",
		Predicate:    "delete",
		AltPredicate: "rm",
		FlagSet:      flag.NewFlagSet("delete infrastructure", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_id_or_label": c.FlagSet.String("id", _nilDefaultStr, "(Required) Infrastructure's id or label. Note that using the 'label' might be ambiguous in certain situations."),
				"autoconfirm":                c.FlagSet.Bool("autoconfirm", false, "(Flag) If set it does not ask for confirmation anymore"),
			}
		},
		ExecuteFunc: infrastructureDeleteCmd,
	},
	Command{
		Description:  "Deploy an infrastructure.",
		Subject:      "infrastructure",
		AltSubject:   "infra",
		Predicate:    "deploy",
		AltPredicate: "apply",
		FlagSet:      flag.NewFlagSet("deploy infrastructure", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_id_or_label":     c.FlagSet.String("id", _nilDefaultStr, "(Required) Infrastructure's id or label. Note that using the 'label' might be ambiguous in certain situations."),
				"no_hard_shutdown_after_timeout": c.FlagSet.Bool("no_hard_shutdown_after_timeout", false, "(Flag) If set do not force a hard power off after timeout expired and the server is not powered off."),
				"no_attempt_soft_shutdown":       c.FlagSet.Bool("no_attempt_soft_shutdown", false, "(Flag) If set,do not atempt a soft (ACPI) power off of all the servers in the infrastructure before the deploy"),
				"soft_shutdown_timeout_seconds":  c.FlagSet.Int("soft_shutdown_timeout_seconds", 180, "(Optional, default 180) Timeout to wait if hard_shutdown_after_timeout is set."),
				"allow_data_loss":                c.FlagSet.Bool("allow_data_loss", false, "(Flag) If set, deploy will throw error if data loss is expected."),
				"skip_ansible":                   c.FlagSet.Bool("skip_ansible", false, "(Flag) If set, some automatic provisioning steps will be skipped. This parameter should generally be ignored."),
				"autoconfirm":                    c.FlagSet.Bool("autoconfirm", false, "(Flag) If set operation procedes without asking for confirmation"),
			}
		},
		ExecuteFunc: infrastructureDeployCmd,
	},
	Command{
		Description:  "Get an infrastructure.",
		Subject:      "infrastructure",
		AltSubject:   "infra",
		Predicate:    "get",
		AltPredicate: "show",
		FlagSet:      flag.NewFlagSet("get infrastructure", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_id_or_label": c.FlagSet.String("id", _nilDefaultStr, "(Required) Infrastructure's id or label. Note that using the 'label' might be ambiguous in certain situations."),
				"format":                     c.FlagSet.String("format", "", "The output format. Supported values are 'json','csv'. The default format is human readable."),
			}
		},
		ExecuteFunc: infrastructureGetCmd,
	},
	Command{
		Description:  "Revert all changes of an infrastructure.",
		Subject:      "infrastructure",
		AltSubject:   "infra",
		Predicate:    "revert",
		AltPredicate: "undo",
		FlagSet:      flag.NewFlagSet("deploy infrastructure", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_id_or_label": c.FlagSet.String("id", _nilDefaultStr, "(Required) Infrastructure's id or label. Note that using the 'label' might be ambiguous in certain situations."),
				"autoconfirm":                c.FlagSet.Bool("autoconfirm", false, "(Flag) If set it does not ask for confirmation anymore"),
			}
		},
		ExecuteFunc: infrastructureRevertCmd,
	},
}

func infrastructureCreateCmd(c *Command, client MetalCloudClient) (string, error) {

	infrastructureLabel := c.Arguments["infrastructure_label"]

	if infrastructureLabel == nil || *infrastructureLabel.(*string) == "" {
		return "", fmt.Errorf("-label <infrastructure_label> is required")
	}

	datacenter := c.Arguments["datacenter"]

	ia := metalcloud.Infrastructure{
		InfrastructureLabel: *infrastructureLabel.(*string),
		DatacenterName:      *datacenter.(*string),
	}

	retInfra, err := client.InfrastructureCreate(ia)
	if err != nil {
		return "", err
	}

	if c.Arguments["return_id"] != nil && *c.Arguments["return_id"].(*bool) {
		return fmt.Sprintf("%d", retInfra.InfrastructureID), nil
	}

	return "", nil
}

func infrastructureListCmd(c *Command, client MetalCloudClient) (string, error) {

	iList, err := client.Infrastructures()
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
			FieldSize: 40,
		},
		SchemaField{
			FieldName: "OWNER",
			FieldType: TypeString,
			FieldSize: 35,
		},
		SchemaField{
			FieldName: "REL.",
			FieldType: TypeString,
			FieldSize: 10,
		},
		SchemaField{
			FieldName: "STATUS",
			FieldType: TypeString,
			FieldSize: 10,
		},
		SchemaField{
			FieldName: "DATACENTER",
			FieldType: TypeString,
			FieldSize: 20,
		},
	}

	user := GetUserEmail()

	data := [][]interface{}{}
	for _, i := range *iList {
		relation := "OWNER"
		if i.UserEmailOwner != user {
			relation = "_DELEGATE"
		}
		data = append(data, []interface{}{
			i.InfrastructureID,
			i.InfrastructureOperation.InfrastructureLabel,
			i.UserEmailOwner,
			relation,
			i.InfrastructureServiceStatus,
			i.DatacenterName,
		})

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
		sb.WriteString(fmt.Sprintf("Infrastructures I have access to (as %s) in datacenter %s\n", user, GetDatacenter()))

		TableSorter(schema).OrderBy(
			schema[3].FieldName,
			schema[0].FieldName,
			schema[1].FieldName).Sort(data)

		sb.WriteString(GetTableAsString(data, schema))

		sb.WriteString(fmt.Sprintf("Total: %d Infrastructures\n\n", len(*iList)))
	}

	return sb.String(), nil
}

func infrastructureDeleteCmd(c *Command, client MetalCloudClient) (string, error) {

	retInfra, err := getInfrastructureFromCommand(c, client)
	if err != nil {
		return "", err
	}

	confirm := false

	if c.Arguments["autoconfirm"] != nil && *c.Arguments["autoconfirm"].(*bool) == true {
		confirm = true
	} else {

		confirmationMessage := fmt.Sprintf("Deleting infrastructure %s (%d). Are you sure? Type \"yes\" to continue:", retInfra.InfrastructureLabel, retInfra.InfrastructureID)

		//this is simply so that we don't output a text on the command line under go test
		if strings.HasSuffix(os.Args[0], ".test") {
			confirmationMessage = ""
		}

		confirm = requestConfirmation(confirmationMessage)
	}

	if !confirm {
		return "", fmt.Errorf("Operation not confirmed. Aborting")
	}

	err = client.InfrastructureDelete(retInfra.InfrastructureID)

	return "", err
}

func infrastructureDeployCmd(c *Command, client MetalCloudClient) (string, error) {

	retInfra, err := getInfrastructureFromCommand(c, client)
	if err != nil {
		return "", err
	}

	confirm := false

	if c.Arguments["autoconfirm"] != nil && *c.Arguments["autoconfirm"].(*bool) {
		confirm = true
	} else {

		confirmationMessage := fmt.Sprintf("Deploying infrastructure %s (%d). Are you sure? Type \"yes\" to continue:", retInfra.InfrastructureLabel, retInfra.InfrastructureID)
		//this is simply so that we don't output a text on the command line under go test
		if strings.HasSuffix(os.Args[0], ".test") {
			confirmationMessage = ""
		}

		confirm = requestConfirmation(confirmationMessage)
	}

	if !confirm {
		return "", fmt.Errorf("Operation not confirmed. Aborting")
	}

	timeout := 180
	if c.Arguments["soft_shutdown_timeout_seconds"] != nil {
		timeout = *c.Arguments["soft_shutdown_timeout_seconds"].(*int)
	}

	NoHardShutdownAfterTimeout := c.Arguments["no_hard_shutdown_after_timeout"] != nil && *c.Arguments["no_hard_shutdown_after_timeout"].(*bool)
	NoAttemptSoftShutdown := c.Arguments["no_attempt_soft_shutdown"] != nil && *c.Arguments["no_attempt_soft_shutdown"].(*bool)

	shutDownOptions := metalcloud.ShutdownOptions{
		HardShutdownAfterTimeout:   !NoHardShutdownAfterTimeout,
		AttemptSoftShutdown:        !NoAttemptSoftShutdown,
		SoftShutdownTimeoutSeconds: timeout,
	}

	err = client.InfrastructureDeploy(
		retInfra.InfrastructureID,
		shutDownOptions,
		c.Arguments["allow_data_loss"] != nil && *c.Arguments["allow_data_loss"].(*bool),
		c.Arguments["skip_ansible"] != nil && *c.Arguments["skip_ansible"].(*bool),
	)

	return "", err
}

func infrastructureRevertCmd(c *Command, client MetalCloudClient) (string, error) {

	retInfra, err := getInfrastructureFromCommand(c, client)
	if err != nil {
		return "", err
	}

	confirm := false

	if c.Arguments["autoconfirm"] != nil && *c.Arguments["autoconfirm"].(*bool) == true {
		confirm = true
	} else {

		fmt.Printf("Reverting infrastructure %s (%d) to the deployed state. Are you sure? Type \"yes\" to continue:", retInfra.InfrastructureLabel, retInfra.InfrastructureID)
		reader := bufio.NewReader(os.Stdin)
		yes, _ := reader.ReadString('\n')

		if yes == "yes\n" {
			confirm = true
		}
	}

	if !confirm {
		return "", fmt.Errorf("Operation not confirmed. Aborting")
	}

	err = client.InfrastructureOperationCancel(retInfra.InfrastructureID)

	return "", err
}

func infrastructureGetCmd(c *Command, client MetalCloudClient) (string, error) {

	retInfra, err := getInfrastructureFromCommand(c, client)
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
			FieldName: "OBJECT_TYPE",
			FieldType: TypeString,
			FieldSize: 15,
		},
		SchemaField{
			FieldName: "LABEL",
			FieldType: TypeString,
			FieldSize: 25,
		},
		SchemaField{
			FieldName: "DETAILS",
			FieldType: TypeString,
			FieldSize: 75,
		},
		SchemaField{
			FieldName: "STATUS",
			FieldType: TypeString,
			FieldSize: 10,
		},
	}

	data := [][]interface{}{}

	iaList, err := client.InstanceArrays(retInfra.InfrastructureID)
	if err != nil {
		return "", err
	}

	for _, ia := range *iaList {
		status := ia.InstanceArrayServiceStatus
		if ia.InstanceArrayServiceStatus != "ordered" && ia.InstanceArrayOperation.InstanceArrayDeployType == "edit" && ia.InstanceArrayOperation.InstanceArrayDeployStatus == "not_started" {
			status = "edited"
		}

		volumeTemplateName := ""
		if ia.InstanceArrayOperation.VolumeTemplateID != 0 {
			vt, err := client.VolumeTemplateGet(ia.InstanceArrayOperation.VolumeTemplateID)
			if err != nil {
				return "", err
			}
			volumeTemplateName = fmt.Sprintf("%s [#%d] ", vt.VolumeTemplateDisplayName, vt.VolumeTemplateID)
		}

		fwMgmtDisabled := ""
		if !ia.InstanceArrayFirewallManaged {
			fwMgmtDisabled = " fw mgmt disabled"
		}
		details := fmt.Sprintf("%d instances (%d RAM, %d cores, %d disks %s %s%s)",
			ia.InstanceArrayOperation.InstanceArrayInstanceCount,
			ia.InstanceArrayOperation.InstanceArrayRAMGbytes,
			ia.InstanceArrayOperation.InstanceArrayProcessorCount*ia.InstanceArrayProcessorCoreCount,
			ia.InstanceArrayOperation.InstanceArrayDiskCount,
			ia.InstanceArrayOperation.InstanceArrayBootMethod,
			volumeTemplateName,
			fwMgmtDisabled,
		)

		data = append(data, []interface{}{
			ia.InstanceArrayID,
			"InstanceArray",
			ia.InstanceArrayOperation.InstanceArrayLabel,
			details,
			status,
		})

	}

	daList, err := client.DriveArrays(retInfra.InfrastructureID)
	if err != nil {
		return "", err
	}

	for _, da := range *daList {
		status := da.DriveArrayServiceStatus
		if da.DriveArrayServiceStatus != "ordered" && da.DriveArrayOperation.DriveArrayDeployType == "edit" && da.DriveArrayOperation.DriveArrayDeployStatus == "not_started" {
			status = "edited"
		}

		volumeTemplateName := ""
		if da.DriveArrayOperation.VolumeTemplateID != 0 {
			vt, err := client.VolumeTemplateGet(da.DriveArrayOperation.VolumeTemplateID)
			if err != nil {
				return "", err
			}
			volumeTemplateName = fmt.Sprintf("%s [#%d]", vt.VolumeTemplateDisplayName, vt.VolumeTemplateID)
		}

		attachedToInstanceArrayStr := ""
		for _, ia := range *iaList {
			if ia.InstanceArrayID == da.DriveArrayOperation.InstanceArrayID {
				attachedToInstanceArrayStr = fmt.Sprintf("%s [#%d]", ia.InstanceArrayLabel, ia.InstanceArrayID)
				break
			}
		}

		data = append(data, []interface{}{
			da.DriveArrayID,
			"DriveArray",
			da.DriveArrayOperation.DriveArrayLabel,
			fmt.Sprintf("%d drives - %.1f GB %s %s attached to: %s",
				da.DriveArrayOperation.DriveArrayCount,
				float64(da.DriveArrayOperation.DriveSizeMBytesDefault/1024),
				da.DriveArrayOperation.DriveArrayStorageType,
				volumeTemplateName,
				attachedToInstanceArrayStr),
			status,
		})

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

		sb.WriteString(fmt.Sprintf("Infrastructure %s (%d) - datacenter %s owner %s\n",
			retInfra.InfrastructureLabel,
			retInfra.InfrastructureID,
			retInfra.DatacenterName,
			retInfra.UserEmailOwner))
		sb.WriteString(GetTableAsString(data, schema))

		sb.WriteString(fmt.Sprintf("Total: %d elements\n\n", len(data)))
	}

	return sb.String(), nil
}

func getInfrastructureFromCommand(c *Command, client MetalCloudClient) (*metalcloud.Infrastructure, error) {

	if c.Arguments["infrastructure_id_or_label"] == nil {
		return nil, fmt.Errorf("Either an infrastructure ID or an infrastructure label must be provided")
	}

	switch v := c.Arguments["infrastructure_id_or_label"].(type) {

	case *int:
		if *v != _nilDefaultInt {
			return client.InfrastructureGet(*v)
		}

	case *string:
		infrastructureID, err := strconv.Atoi(*v)
		if err == nil {
			return client.InfrastructureGet(infrastructureID)
		}
		if *v == _nilDefaultStr {
			return nil, fmt.Errorf("Either an infrastructure ID or an infrastructure label must be provided")
		}
	default:
		return nil, fmt.Errorf("format not supported")

	}

	var infrastructure *metalcloud.Infrastructure

	ret, err := client.Infrastructures()
	if err != nil {
		return nil, err
	}

	for k, i := range *ret {
		if i.InfrastructureOperation.InfrastructureLabel == *c.Arguments["infrastructure_id_or_label"].(*string) {

			if infrastructure != nil {
				//if we found this infrastructure label, with the same name again, we throw an error
				return nil, fmt.Errorf("Infrastructures %d and %d both have the same label %s", infrastructure.InfrastructureID, i.InfrastructureID, *c.Arguments["infrastructure_id_or_label"].(*string))
			}

			infr := (*ret)[k]
			infrastructure = &infr
			//we let the search go on to check for ambiguous situationss
		}
	}

	if infrastructure == nil {
		return nil, fmt.Errorf("Could not find infrastructure with label %s", *c.Arguments["infrastructure_id_or_label"].(*string))
	}

	return infrastructure, nil
}
