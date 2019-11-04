package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	metalcloud "github.com/bigstepinc/metal-cloud-sdk-go"
)

//InfrastructureCmds commands affecting infrastructures
var InfrastructureCmds = []Command{

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
				"datacenter":           c.FlagSet.String("dc", "", "(Required) Infrastructure datacenter"),
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
				"format": c.FlagSet.String("format", "", "The output format. Supproted values are 'json','csv'. The default format is human readable."),
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
				"infrastructure_id": c.FlagSet.Int("id", 0, "(Required) Infrastructure's id"),
				"autoconfirm":       c.FlagSet.Bool("autoconfirm", false, "If true it does not ask for confirmation anymore"),
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
				"infrastructure_id":             c.FlagSet.Int("id", 0, "(Required) Infrastructure's id"),
				"hard_shutdown_after_timeout":   c.FlagSet.Bool("hard_shutdown_after_timeout", true, "(Optional, default true) Force a hard power off after timeout expired and the server is not powered off."),
				"attempt_soft_shutdown":         c.FlagSet.Bool("attempt_soft_shutdown", true, "(Optional, default true) If needed, atempt a soft (ACPI) power off of all the servers in the infrastructure before the deploy"),
				"soft_shutdown_timeout_seconds": c.FlagSet.Int("soft_shutdown_timeout_seconds", 180, "(Optional, default 180) Timeout to wait if hard_shutdown_after_timeout is set."),
				"allow_data_loss":               c.FlagSet.Bool("allow_data_loss", false, "(Optional, default false) If false, deploy will throw error if data loss is expected."),
				"skip_ansible":                  c.FlagSet.Bool("skip_ansible", false, "(Optional, default false) If true, some automatic provisioning steps will be skipped. This parameter should generally be ignored."),
				"autoconfirm":                   c.FlagSet.Bool("autoconfirm", false, "If true it does not ask for confirmation anymore"),
			}
		},
		ExecuteFunc: infrastructureDeleteCmd,
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
				"infrastructure_id": c.FlagSet.Int("id", 0, "(Required) Infrastructure's id"),
				"format":            c.FlagSet.String("format", "", "The output format. Supproted values are 'json','csv'. The default format is human readable."),
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
				"infrastructure_id": c.FlagSet.Int("id", 0, "(Required) Infrastructure's id"),
				"autoconfirm":       c.FlagSet.Bool("autoconfirm", false, "If true it does not ask for confirmation anymore"),
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

	if datacenter == nil || *datacenter.(*string) == "" {
		//	return "", fmt.Errorf("-dc <datacenter> is required")
		datacenter = GetDatacenter()
	}

	ia := metalcloud.Infrastructure{
		InfrastructureLabel: *infrastructureLabel.(*string),
		DatacenterName:      *datacenter.(*string),
	}
	_, err := client.InfrastructureCreate(ia)
	if err != nil {
		return "", err
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
			FieldSize: 30,
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
			FieldName: "CREATED",
			FieldType: TypeString,
			FieldSize: 20,
		},
		SchemaField{
			FieldName: "UPDATED",
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
			i.InfrastructureLabel,
			i.UserEmailOwner,
			relation,
			i.InfrastructureServiceStatus,
			i.InfrastructureCreatedTimestamp,
			i.InfrastructureUpdatedTimestamp,
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
		sb.WriteString(fmt.Sprintf("Infrastructures I have access to (as %s)\n", user))

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

	infrastructureID := c.Arguments["infrastructure_id"]

	if infrastructureID == nil || *infrastructureID.(*int) == 0 {
		return "", fmt.Errorf("-id <infrastructure_id> is required")
	}

	ret, err2 := client.InfrastructureGet(*infrastructureID.(*int))
	if err2 != nil {
		return "", err2
	}

	autoConfirm := c.Arguments["autoconfirm"]

	confirm := false

	if autoConfirm == nil || autoConfirm == false {
		fmt.Printf("Deleting infrastructure %s (%d). Are you sure? Type \"yes\" to continue:", ret.InfrastructureLabel, ret.InfrastructureID)
		reader := bufio.NewReader(os.Stdin)
		yes, _ := reader.ReadString('\n')

		if yes == "yes\n" {
			confirm = true
		}

	} else {
		confirm = true
	}

	if confirm {
		err := client.InfrastructureDelete(*infrastructureID.(*int))
		if err != nil {
			return "", err
		}
	}

	return "", nil
}

func infrastructureDeployCmd(c *Command, client MetalCloudClient) (string, error) {

	infrastructureID := c.Arguments["infrastructure_id"]

	if infrastructureID == nil || *infrastructureID.(*int) == 0 {
		return "", fmt.Errorf("-id <infrastructure_id> is required")
	}

	ret, err2 := client.InfrastructureGet(*infrastructureID.(*int))
	if err2 != nil {
		return "", err2
	}

	autoConfirm := c.Arguments["autoconfirm"]

	confirm := false

	if autoConfirm == nil || autoConfirm == false {
		fmt.Printf("Deploying infrastructure %s (%d). Are you sure? Type \"yes\" to continue:", ret.InfrastructureLabel, ret.InfrastructureID)
		reader := bufio.NewReader(os.Stdin)
		yes, _ := reader.ReadString('\n')

		if yes == "yes\n" {
			confirm = true
		}

	} else {
		confirm = true
	}

	timeout := 180
	if c.Arguments["soft_shutdown_timeout_seconds"] != nil {
		timeout = *c.Arguments["soft_shutdown_timeout_seconds"].(*int)
	}

	shutDownOptions := metalcloud.ShutdownOptions{
		HardShutdownAfterTimeout:   c.Arguments["hard_shutdown_after_timeout"] != nil && *c.Arguments["hard_shutdown_after_timeout"].(*bool),
		AttemptSoftShutdown:        c.Arguments["attempt_soft_shutdown"] != nil && *c.Arguments["attempt_soft_shutdown"].(*bool),
		SoftShutdownTimeoutSeconds: timeout,
	}

	if confirm {
		err := client.InfrastructureDeploy(*infrastructureID.(*int),
			shutDownOptions,
			c.Arguments["allow_data_loss"] != nil && *c.Arguments["allow_data_loss"].(*bool),
			c.Arguments["skip_ansible"] != nil && *c.Arguments["skip_ansible"].(*bool),
		)
		if err != nil {
			return "", err
		}
	}

	return "", nil
}

func infrastructureRevertCmd(c *Command, client MetalCloudClient) (string, error) {

	infrastructureID := c.Arguments["infrastructure_id"]

	if infrastructureID == nil || *infrastructureID.(*int) == 0 {
		return "", fmt.Errorf("-id <infrastructure_id> is required")
	}

	ret, err2 := client.InfrastructureGet(*infrastructureID.(*int))
	if err2 != nil {
		return "", err2
	}

	autoConfirm := c.Arguments["autoconfirm"]

	confirm := false

	if autoConfirm == nil || autoConfirm == false {
		fmt.Printf("Reverting infrastructure %s (%d) to the deployed state. Are you sure? Type \"yes\" to continue:", ret.InfrastructureLabel, ret.InfrastructureID)
		reader := bufio.NewReader(os.Stdin)
		yes, _ := reader.ReadString('\n')

		if yes == "yes\n" {
			confirm = true
		}

	} else {
		confirm = true
	}

	if confirm {
		err := client.InfrastructureOperationCancel(*infrastructureID.(*int))
		if err != nil {
			return "", err
		}
	}

	return "", nil
}

func infrastructureGetCmd(c *Command, client MetalCloudClient) (string, error) {

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
			FieldSize: 30,
		},
		SchemaField{
			FieldName: "DETAILS",
			FieldType: TypeString,
			FieldSize: 70,
		},
		SchemaField{
			FieldName: "STATUS",
			FieldType: TypeString,
			FieldSize: 10,
		},
	}

	infrastructureID := c.Arguments["infrastructure_id"]

	if infrastructureID == nil || *infrastructureID.(*int) == 0 {
		return "", fmt.Errorf("-id <infrastructure_id> is required")
	}

	data := [][]interface{}{}

	iaList, err := client.InstanceArrays(*infrastructureID.(*int))
	if err != nil {
		return "", err
	}

	for _, ia := range *iaList {
		status := ia.InstanceArrayServiceStatus
		if ia.InstanceArrayServiceStatus != "ordered" && ia.InstanceArrayOperation.InstanceArrayDeployType == "edit" && ia.InstanceArrayOperation.InstanceArrayDeployStatus == "not_started" {
			status = "edited"
		}
		data = append(data, []interface{}{
			ia.InstanceArrayID,
			"InstanceArray",
			ia.InstanceArrayOperation.InstanceArrayLabel,
			fmt.Sprintf("%d instances (%d RAM, %d cores, %d disks)",
				ia.InstanceArrayOperation.InstanceArrayInstanceCount,
				ia.InstanceArrayOperation.InstanceArrayRAMGbytes,
				ia.InstanceArrayOperation.InstanceArrayProcessorCount*ia.InstanceArrayProcessorCoreCount,
				ia.InstanceArrayOperation.InstanceArrayDiskCount),
			status,
		})

	}

	daList, err := client.DriveArrays(*infrastructureID.(*int))
	if err != nil {
		return "", err
	}

	for _, da := range *daList {
		status := da.DriveArrayServiceStatus
		if da.DriveArrayServiceStatus != "ordered" && da.DriveArrayOperation.DriveArrayDeployType == "edit" && da.DriveArrayOperation.DriveArrayDeployStatus == "not_started" {
			status = "edited"
		}
		data = append(data, []interface{}{
			da.DriveArrayID,
			"DriveArray",
			da.DriveArrayOperation.DriveArrayLabel,
			fmt.Sprintf("%d drives - %.1f GB %s (volume_template:%d) attached to: %d",
				da.DriveArrayOperation.DriveArrayCount,
				float64(da.DriveArrayOperation.DriveSizeMBytesDefault/1024),
				da.DriveArrayOperation.DriveArrayStorageType,
				da.DriveArrayOperation.VolumeTemplateID,
				da.DriveArrayOperation.InstanceArrayID),
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

		sb.WriteString(fmt.Sprintf("Infrastructures I have access to (as %s)\n", GetUserEmail()))
		/*
			TableSorter(schema).OrderBy(
				schema[3].FieldName,
				schema[0].FieldName,
				schema[1].FieldName).Sort(data)
		*/
		sb.WriteString(GetTableAsString(data, schema))

		sb.WriteString(fmt.Sprintf("Total: %d elements\n\n", len(data)))
	}

	return sb.String(), nil
}
