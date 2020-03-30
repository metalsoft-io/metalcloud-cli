package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	metalcloud "github.com/bigstepinc/metal-cloud-sdk-go"
	interfaces "github.com/bigstepinc/metalcloud-cli/interfaces"
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
				"datacenter":           c.FlagSet.String("datacenter", GetDatacenter(), "(Required) Infrastructure datacenter"),
				"return_id":            c.FlagSet.Bool("return-id", false, "(Flag) If set will print the ID of the created infrastructure. Useful for automating tasks."),
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
				"no_hard_shutdown_after_timeout": c.FlagSet.Bool("no-hard-shutdown-after-timeout", false, "(Flag) If set do not force a hard power off after timeout expired and the server is not powered off."),
				"no_attempt_soft_shutdown":       c.FlagSet.Bool("no-attempt-soft-shutdown", false, "(Flag) If set,do not atempt a soft (ACPI) power off of all the servers in the infrastructure before the deploy"),
				"soft_shutdown_timeout_seconds":  c.FlagSet.Int("soft-shutdown-timeout-seconds", 180, "(Optional, default 180) Timeout to wait if hard_shutdown_after_timeout is set."),
				"allow_data_loss":                c.FlagSet.Bool("allow-data-loss", false, "(Flag) If set, deploy will throw error if data loss is expected."),
				"skip_ansible":                   c.FlagSet.Bool("skip-ansible", false, "(Flag) If set, some automatic provisioning steps will be skipped. This parameter should generally be ignored."),
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
	Command{
		Description:  "list stages of a workflow",
		Subject:      "infrastructure",
		AltSubject:   "infra",
		Predicate:    "workflow-stages",
		AltPredicate: "workflow-stages",
		FlagSet:      flag.NewFlagSet("list stages of a workflow", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_id_or_label": c.FlagSet.String("id", _nilDefaultStr, "The infrastructure's id"),
				"type":                       c.FlagSet.String("type", _nilDefaultStr, "stage definition type. possible values: pre_deploy, post_deploy"),
			}
		},
		ExecuteFunc: listWorkflowStagesCmd,
	},
}

func infrastructureCreateCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

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

func infrastructureListCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

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
			FieldSize: 15,
		},
		SchemaField{
			FieldName: "OWNER",
			FieldType: TypeString,
			FieldSize: 20,
		},
		SchemaField{
			FieldName: "REL.",
			FieldType: TypeString,
			FieldSize: 10,
		},
		SchemaField{
			FieldName: "STATUS",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "DATACENTER",
			FieldType: TypeString,
			FieldSize: 10,
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

	TableSorter(schema).OrderBy(
		schema[3].FieldName,
		schema[0].FieldName,
		schema[1].FieldName).Sort(data)

	topLine := fmt.Sprintf("Infrastructures I have access to (as %s) in datacenter %s\n", user, GetDatacenter())
	return renderTable("Infrastructures", topLine, getStringParam(c.Arguments["format"]), data, schema)
}

type infrastructureConfirmAndDoFunc func(infraID int, c *Command, client interfaces.MetalCloudClient) (string, error)

//infrastructureConfirmAndDo asks for confirmation and executes the given function
func infrastructureConfirmAndDo(operation string, c *Command, client interfaces.MetalCloudClient, f infrastructureConfirmAndDoFunc) (string, error) {

	val, err := getParam(c, "infrastructure_id_or_label", "infra")
	if err != nil {
		return "", err
	}

	infraID, err := getIDOrDo(*val.(*string), func(label string) (int, error) {
		ia, err := client.InfrastructureGetByLabel(label)
		if err != nil {
			return 0, err
		}
		return ia.InfrastructureID, nil
	})
	if err != nil {
		return "", err
	}

	confirm := false

	if c.Arguments["autoconfirm"] != nil && *c.Arguments["autoconfirm"].(*bool) == true {
		confirm = true
	} else {

		retInfra, err := client.InfrastructureGet(infraID)
		if err != nil {
			return "", err
		}

		confirmationMessage := fmt.Sprintf("%s infrastructure %s (%d). Are you sure? Type \"yes\" to continue:", operation, retInfra.InfrastructureLabel, retInfra.InfrastructureID)

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

	return f(infraID, c, client)
}

func infrastructureDeleteCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {
	return infrastructureConfirmAndDo("Delete", c, client,
		func(infraID int, c *Command, client interfaces.MetalCloudClient) (string, error) {
			return "", client.InfrastructureDelete(infraID)
		})
}

func infrastructureDeployCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	return infrastructureConfirmAndDo("Deploy", c, client,
		func(infraID int, c *Command, client interfaces.MetalCloudClient) (string, error) {

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

			return "", client.InfrastructureDeploy(
				infraID,
				shutDownOptions,
				c.Arguments["allow_data_loss"] != nil && *c.Arguments["allow_data_loss"].(*bool),
				c.Arguments["skip_ansible"] != nil && *c.Arguments["skip_ansible"].(*bool),
			)
		})
}

func infrastructureRevertCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	return infrastructureConfirmAndDo("Revert", c, client,
		func(infraID int, c *Command, client interfaces.MetalCloudClient) (string, error) {
			return "", client.InfrastructureOperationCancel(infraID)
		})
}

func infrastructureGetCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	retInfra, err := getInfrastructureFromCommand("id", c, client)
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
			FieldSize: 10,
		},
		SchemaField{
			FieldName: "DETAILS",
			FieldType: TypeString,
			FieldSize: 50,
		},
		SchemaField{
			FieldName: "STATUS",
			FieldType: TypeString,
			FieldSize: 5,
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

		AdjustFieldSizes(data, &schema)

		sb.WriteString(GetTableAsString(data, schema))

		sb.WriteString(fmt.Sprintf("Total: %d elements\n\n", len(data)))
	}

	return sb.String(), nil
}

func listWorkflowStagesCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	t := *c.Arguments["type"].(*string)
	if t == _nilDefaultStr {
		t = "post_deploy"
	}

	retInfra, err := getInfrastructureFromCommand("id", c, client)
	if err != nil {
		return "", err
	}

	list, err := client.InfrastructureDeployCustomStages(retInfra.InfrastructureID, t)

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
			FieldName: "INFRASTRUCTRE",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "STAGE",
			FieldType: TypeString,
			FieldSize: 4,
		},
		SchemaField{
			FieldName: "TYPE",
			FieldType: TypeString,
			FieldSize: 4,
		},
		SchemaField{
			FieldName: "RUNLEVEL",
			FieldType: TypeInt,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "OUTPUT",
			FieldType: TypeString,
			FieldSize: 20,
		},
	}

	data := [][]interface{}{}
	for _, s := range *list {

		infra, err := client.InfrastructureGet(s.InfrastructureID)
		if err != nil {
			return "", err
		}

		stage, err := client.StageDefinitionGet(s.StageDefinitionID)
		if err != nil {
			return "", err
		}

		data = append(data, []interface{}{
			s.InfrastructureDeployCustomStageID,
			infra.InfrastructureLabel,
			stage.StageDefinitionLabel,
			s.InfrastructureDeployCustomStageType,
			s.InfrastructureDeployCustomStageRunLevel,
			s.InfrastructureDeployCustomStageExecOutputJSON,
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
		sb.WriteString(fmt.Sprintf("Stage Definitions:\n"))

		TableSorter(schema).OrderBy(
			schema[0].FieldName,
			schema[1].FieldName).Sort(data)

		AdjustFieldSizes(data, &schema)

		sb.WriteString(GetTableAsString(data, schema))

		sb.WriteString(fmt.Sprintf("Total: %d \n\n", len(*list)))
	}

	return sb.String(), nil
}

//getInfrastructureFromCommand returns an Infrastructure object using the infrastructure_id_or_label argument
func getInfrastructureFromCommand(paramName string, c *Command, client interfaces.MetalCloudClient) (*metalcloud.Infrastructure, error) {

	m, err := getParam(c, "infrastructure_id_or_label", paramName)
	if err != nil {
		return nil, err
	}

	id, label, isID := idOrLabel(m)

	if isID {
		return client.InfrastructureGet(id)
	}

	return client.InfrastructureGetByLabel(label)
}
