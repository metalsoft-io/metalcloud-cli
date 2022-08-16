package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
	"github.com/metalsoft-io/tableformatter"
)

//infrastructureCmds commands affecting infrastructures
var infrastructureCmds = []Command{

	{
		Description:  "Creates an infrastructure.",
		Subject:      "infrastructure",
		AltSubject:   "infra",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("create infrastructure", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_label": c.FlagSet.String("label", "", red("(Required)")+" Infrastructure's label"),
				"datacenter":           c.FlagSet.String("datacenter", _nilDefaultStr, red("(Required)")+" Infrastructure datacenter"),
				"return_id":            c.FlagSet.Bool("return-id", false, green("(Flag)")+" If set will print the ID of the created infrastructure. Useful for automating tasks."),
			}
		},
		ExecuteFunc: infrastructureCreateCmd,
	},
	{
		Description:  "Lists all infrastructures.",
		Subject:      "infrastructure",
		AltSubject:   "infra",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list infrastructure", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"format":       c.FlagSet.String("format", "", "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
				"filter":       c.FlagSet.String("filter", "*", "filter to use when searching for servers. Check the documentation for examples. Defaults to '*'"),
				"show_ordered": c.FlagSet.Bool("show-ordered", false, green("(Flag)")+" If set will also return ordered (created but not deployed) infrastructures. Default is false."),
				"show_deleted": c.FlagSet.Bool("show-deleted", false, green("(Flag)")+" If set will also return deleted infrastructures. Default is false."),
			}
		},
		ExecuteFunc: infrastructureListCmd,
		Endpoint:    DeveloperEndpoint,
	},
	{
		Description:  "Delete an infrastructure.",
		Subject:      "infrastructure",
		AltSubject:   "infra",
		Predicate:    "delete",
		AltPredicate: "rm",
		FlagSet:      flag.NewFlagSet("delete infrastructure", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_id_or_label": c.FlagSet.String("id", _nilDefaultStr, red("(Required)")+" Infrastructure's id or label. Note that using the 'label' might be ambiguous in certain situations."),
				"autoconfirm":                c.FlagSet.Bool("autoconfirm", false, green("(Flag)")+" If set it will assume action is confirmed"),
			}
		},
		ExecuteFunc: infrastructureDeleteCmd,
		Endpoint:    DeveloperEndpoint,
	},
	{
		Description:  "Deploy an infrastructure.",
		Subject:      "infrastructure",
		AltSubject:   "infra",
		Predicate:    "deploy",
		AltPredicate: "apply",
		FlagSet:      flag.NewFlagSet("deploy infrastructure", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_id_or_label":     c.FlagSet.String("id", _nilDefaultStr, red("(Required)")+" Infrastructure's id or label. Note that using the 'label' might be ambiguous in certain situations."),
				"no_hard_shutdown_after_timeout": c.FlagSet.Bool("no-hard-shutdown-after-timeout", false, green("(Flag)")+" If set do not force a hard power off after timeout expired and the server is not powered off."),
				"no_attempt_soft_shutdown":       c.FlagSet.Bool("no-attempt-soft-shutdown", false, green("(Flag)")+" If set,do not atempt a soft (ACPI) power off of all the servers in the infrastructure before the deploy"),
				"soft_shutdown_timeout_seconds":  c.FlagSet.Int("soft-shutdown-timeout-seconds", 180, "(Optional, default 180) Timeout to wait if hard_shutdown_after_timeout is set."),
				"allow_data_loss":                c.FlagSet.Bool("allow-data-loss", false, green("(Flag)")+" If set, deploy will not throw error if data loss is expected."),
				"skip_ansible":                   c.FlagSet.Bool("skip-ansible", false, green("(Flag)")+" If set, some automatic provisioning steps will be skipped. This parameter should generally be ignored."),
				"block_until_deployed":           c.FlagSet.Bool("blocking", false, green("(Flag)")+" If set, the operation will wait until deployment finishes."),
				"block_timeout":                  c.FlagSet.Int("block-timeout", 180*60, "Block timeout in seconds. After this timeout the application will return an error. Defaults to 180 minutes."),
				"block_check_interval":           c.FlagSet.Int("block-check-interval", 10, "Check interval for when blocking. Defaults to 10 seconds."),
				"autoconfirm":                    c.FlagSet.Bool("autoconfirm", false, green("(Flag)")+" If set it will assume action is confirmed"),
			}
		},
		ExecuteFunc: infrastructureDeployCmd,
		Endpoint:    DeveloperEndpoint,
	},
	{
		Description:  "Get infrastructure details.",
		Subject:      "infrastructure",
		AltSubject:   "infra",
		Predicate:    "get",
		AltPredicate: "show",
		FlagSet:      flag.NewFlagSet("get infrastructure", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_id_or_label": c.FlagSet.String("id", _nilDefaultStr, red("(Required)")+" Infrastructure's id or label. Note that using the 'label' might be ambiguous in certain situations."),
				"format":                     c.FlagSet.String("format", "", "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
			}
		},
		ExecuteFunc: infrastructureGetCmd,
		Endpoint:    DeveloperEndpoint,
	},
	{
		Description:  "Revert all changes of an infrastructure.",
		Subject:      "infrastructure",
		AltSubject:   "infra",
		Predicate:    "revert",
		AltPredicate: "undo",
		FlagSet:      flag.NewFlagSet("deploy infrastructure", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_id_or_label": c.FlagSet.String("id", _nilDefaultStr, red("(Required)")+" Infrastructure's id or label. Note that using the 'label' might be ambiguous in certain situations."),
				"autoconfirm":                c.FlagSet.Bool("autoconfirm", false, green("(Flag)")+" If set it will assume action is confirmed"),
			}
		},
		ExecuteFunc: infrastructureRevertCmd,
		Endpoint:    DeveloperEndpoint,
	},
	{
		Description:  "List stages of a workflow.",
		Subject:      "infrastructure",
		AltSubject:   "infra",
		Predicate:    "workflow-stages",
		AltPredicate: "workflow-stages",
		FlagSet:      flag.NewFlagSet("List stages of a workflow.", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_id_or_label": c.FlagSet.String("id", _nilDefaultStr, "The infrastructure's id"),
				"type":                       c.FlagSet.String("type", _nilDefaultStr, "stage definition type. possible values: pre_deploy, post_deploy"),
			}
		},
		ExecuteFunc: listWorkflowStagesCmd,
		Endpoint:    DeveloperEndpoint,
	},
}

func infrastructureCreateCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	infrastructureLabel := c.Arguments["infrastructure_label"]

	if infrastructureLabel == nil || *infrastructureLabel.(*string) == "" {
		return "", fmt.Errorf("-label is required")
	}

	datacenter := c.Arguments["datacenter"]

	if datacenter == nil || *datacenter.(*string) == "" {
		return "", fmt.Errorf("-datacenter is required")
	}

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

func infrastructureListCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	filter := getStringParam(c.Arguments["filter"])
	iList, err := client.InfrastructureSearch(filter)
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
			FieldSize: 5,
		},
		{
			FieldName: "OWNER",
			FieldType: tableformatter.TypeString,
			FieldSize: 20,
		},
		{
			FieldName: "DATACENTER",
			FieldType: tableformatter.TypeString,
			FieldSize: 10,
		},
		{
			FieldName: "CREATED",
			FieldType: tableformatter.TypeString,
			FieldSize: 10,
		},
		{
			FieldName: "UPDATED",
			FieldType: tableformatter.TypeString,
			FieldSize: 10,
		},
	}

	data := [][]interface{}{}
	for _, i := range *iList {

		if i.InfrastructureServiceStatus == "ordered" && !getBoolParam(c.Arguments["show_ordered"]) {
			continue
		}

		if i.InfrastructureServiceStatus == "deleted" && !getBoolParam(c.Arguments["show_deleted"]) {
			continue
		}

		status := ""

		if i.InfrastructureServiceStatus == "active" && i.AFCExecutedSuccess == i.AFCTotal {
			status = green("Deployed")
		}

		if i.InfrastructureServiceStatus == "ordered" && i.AFCTotal == 0 {
			status = blue("Ordered (deploy not started)")
		}

		if i.AFCExecutedSuccess < i.AFCTotal {

			if i.AFCThrownError == 0 {
				status = yellow(fmt.Sprintf("Deploy ongoing - %d/%d", i.AFCExecutedSuccess, i.AFCTotal))
			} else {
				status = red(fmt.Sprintf("Deploy ongoing - Thrown error at %d/%d", i.AFCExecutedSuccess, i.AFCTotal))
			}
		}

		if i.InfrastructureServiceStatus == "deleted" {
			status = magenta("Deleted")
		}

		userEmail := ""
		if len(i.UserEmail) > 0 {
			userEmail = i.UserEmail[0]
		}

		data = append(data, []interface{}{
			i.InfrastructureID,
			i.InfrastructureLabel,
			status,
			userEmail,
			i.DatacenterName,
			i.InfrastructureCreatedTimestamp,
			i.InfrastructureUpdatedTimestamp,
		})

	}

	tableformatter.TableSorter(schema).OrderBy(
		schema[3].FieldName,
		schema[0].FieldName,
		schema[1].FieldName).Sort(data)

	topLine := fmt.Sprintf("Infrastructures")

	table := tableformatter.Table{
		Data:   data,
		Schema: schema,
	}
	return table.RenderTable("Infrastructures", topLine, getStringParam(c.Arguments["format"]))
}

type infrastructureConfirmAndDoFunc func(infraID int, c *Command, client metalcloud.MetalCloudClient) (string, error)

//infrastructureConfirmAndDo asks for confirmation and executes the given function
func infrastructureConfirmAndDo(operation string, c *Command, client metalcloud.MetalCloudClient, f infrastructureConfirmAndDoFunc) (string, error) {

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

	if getBoolParam(c.Arguments["autoconfirm"]) {
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

func infrastructureDeleteCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {
	return infrastructureConfirmAndDo("Delete", c, client,
		func(infraID int, c *Command, client metalcloud.MetalCloudClient) (string, error) {
			return "", client.InfrastructureDelete(infraID)
		})
}

func infrastructureDeployCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	return infrastructureConfirmAndDo("Deploy", c, client,
		func(infraID int, c *Command, client metalcloud.MetalCloudClient) (string, error) {

			shutDownOptions := metalcloud.ShutdownOptions{
				HardShutdownAfterTimeout:   !getBoolParam(c.Arguments["no_hard_shutdown_after_timeout"]),
				AttemptSoftShutdown:        !getBoolParam(c.Arguments["no_attempt_soft_shutdown"]),
				SoftShutdownTimeoutSeconds: getIntParam(c.Arguments["soft_shutdown_timeout_seconds"]),
			}

			err := client.InfrastructureDeploy(
				infraID,
				shutDownOptions,
				getBoolParam(c.Arguments["allow_data_loss"]),
				getBoolParam(c.Arguments["skip_ansible"]),
			)
			if err != nil {
				return "", err
			}

			if getBoolParam(c.Arguments["block_until_deployed"]) {

				time.Sleep(time.Duration(getIntParam(c.Arguments["block_check_interval"])) * time.Second) //wait until the system picks up the afc

				err := loopUntilInfraReady(infraID, getIntParam(c.Arguments["block_timeout"]), getIntParam(c.Arguments["block_check_interval"]), client)

				if err != nil && strings.HasPrefix(err.Error(), "timeout after") {
					return "", err
				} //else we ignore errors as they might be infrastrucure not found due to infrastructure being deleted
			}

			return "", nil
		})
}

func infrastructureRevertCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	return infrastructureConfirmAndDo("Revert", c, client,
		func(infraID int, c *Command, client metalcloud.MetalCloudClient) (string, error) {
			return "", client.InfrastructureOperationCancel(infraID)
		})
}

func infrastructureGetCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	retInfra, err := getInfrastructureFromCommand("id", c, client)
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
			FieldName: "OBJECT_TYPE",
			FieldType: tableformatter.TypeString,
			FieldSize: 15,
		},
		{
			FieldName: "LABEL",
			FieldType: tableformatter.TypeString,
			FieldSize: 10,
		},
		{
			FieldName: "DETAILS",
			FieldType: tableformatter.TypeString,
			FieldSize: 50,
		},
		{
			FieldName: "STATUS",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
	}

	data := [][]interface{}{}

	iaList, err := client.InstanceArrays(retInfra.InfrastructureID)
	if err != nil {
		return "", err
	}

	for _, ia := range *iaList {
		status := green(ia.InstanceArrayServiceStatus)
		if ia.InstanceArrayServiceStatus != "ordered" && ia.InstanceArrayOperation.InstanceArrayDeployType == "edit" && ia.InstanceArrayOperation.InstanceArrayDeployStatus == "not_started" {
			status = blue("edited")
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
			green("InstanceArray"),
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
		status := green(da.DriveArrayServiceStatus)
		if da.DriveArrayServiceStatus != "ordered" && da.DriveArrayOperation.DriveArrayDeployType == "edit" && da.DriveArrayOperation.DriveArrayDeployStatus == "not_started" {
			status = blue("edited")
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

		details := fmt.Sprintf("%d drives - %.1f GB %s %s attached to: %s, storage pool: #%d",
			da.DriveArrayOperation.DriveArrayCount,
			float64(da.DriveArrayOperation.DriveSizeMBytesDefault/1024),
			da.DriveArrayOperation.DriveArrayStorageType,
			volumeTemplateName,
			attachedToInstanceArrayStr,
			da.StoragePoolID,
		)

		data = append(data, []interface{}{
			da.DriveArrayID,
			blue("DriveArray"),
			da.DriveArrayOperation.DriveArrayLabel,
			details,
			status,
		})
	}

	sdaList, err := client.SharedDrives(retInfra.InfrastructureID)
	if err != nil {
		return "", err
	}

	for _, sda := range *sdaList {
		status := green(sda.SharedDriveServiceStatus)
		if sda.SharedDriveServiceStatus != "ordered" && sda.SharedDriveOperation.SharedDriveDeployType == "edit" && sda.SharedDriveOperation.SharedDriveDeployStatus == "not_started" {
			status = blue("edited")
		}

		details := fmt.Sprintf("%d GB size, type: %s, i/o limit policy: %s, WWW: %s, storage pool: #%d",
			int(sda.SharedDriveSizeMbytes/1024),
			sda.SharedDriveStorageType,
			sda.SharedDriveIOLimitPolicy,
			sda.SharedDriveWWN,
			sda.StoragePoolID,
		)

		data = append(data, []interface{}{
			sda.SharedDriveID,
			magenta("SharedDrive"),
			sda.SharedDriveLabel,
			details,
			status,
		})

	}

	topLine := fmt.Sprintf("Infrastructure %s (%d) - datacenter %s owner %s",
		retInfra.InfrastructureLabel,
		retInfra.InfrastructureID,
		retInfra.DatacenterName,
		retInfra.UserEmailOwner)

	table := tableformatter.Table{
		Data:   data,
		Schema: schema,
	}
	return table.RenderTable("Infrastructures", topLine, getStringParam(c.Arguments["format"]))
}

func listWorkflowStagesCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

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

	schema := []tableformatter.SchemaField{
		{
			FieldName: "ID",
			FieldType: tableformatter.TypeInt,
			FieldSize: 6,
		},
		{
			FieldName: "INFRASTRUCTRE",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "STAGE",
			FieldType: tableformatter.TypeString,
			FieldSize: 4,
		},
		{
			FieldName: "TYPE",
			FieldType: tableformatter.TypeString,
			FieldSize: 4,
		},
		{
			FieldName: "RUNLEVEL",
			FieldType: tableformatter.TypeInt,
			FieldSize: 5,
		},
		{
			FieldName: "OUTPUT",
			FieldType: tableformatter.TypeString,
			FieldSize: 20,
		},
	}

	data := [][]interface{}{}
	for _, s := range *list {

		stage, err := client.StageDefinitionGet(s.StageDefinitionID)
		if err != nil {
			return "", err
		}

		infrastructureLabel := ""

		if stage.StageDefinitionContext != "global" {
			infra, err := client.InfrastructureGet(s.InfrastructureID)
			if err != nil {
				return "", err
			}
			infrastructureLabel = infra.InfrastructureLabel
		}

		data = append(data, []interface{}{
			s.InfrastructureDeployCustomStageID,
			infrastructureLabel,
			stage.StageDefinitionLabel,
			s.InfrastructureDeployCustomStageType,
			s.InfrastructureDeployCustomStageRunLevel,
			s.InfrastructureDeployCustomStageExecOutputJSON,
		})

	}

	table := tableformatter.Table{
		Data:   data,
		Schema: schema,
	}
	return table.RenderTable("Workflow Stages", "", getStringParam(c.Arguments["format"]))
}

//getInfrastructureFromCommand returns an Infrastructure object using the infrastructure_id_or_label argument
func getInfrastructureFromCommand(paramName string, c *Command, client metalcloud.MetalCloudClient) (*metalcloud.Infrastructure, error) {

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

//loop until infra is ready
func loopUntilInfraReady(infraID int, timeoutSeconds int, checkIntervalSeconds int, client metalcloud.MetalCloudClient) error {
	c := make(chan error, 1)

	go func() {
		for {
			infra, err := client.InfrastructureGet(infraID)

			if err != nil {
				c <- err
				break
			}
			if infra.InfrastructureOperation.InfrastructureDeployStatus == "ongoing" {
				time.Sleep(time.Duration(checkIntervalSeconds) * time.Second)
			} else {
				break
			}
		}
		c <- nil
	}()

	select {
	case err := <-c:
		return err
	case <-time.After(time.Duration(timeoutSeconds) * time.Second):
		return fmt.Errorf("timeout after %d seconds while waiting for infrastructure to finish deploying", timeoutSeconds)
	}
}
