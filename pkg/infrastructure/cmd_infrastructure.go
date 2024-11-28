package infrastructure

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
	"github.com/metalsoft-io/metalcloud-cli/internal/colors"
	"github.com/metalsoft-io/metalcloud-cli/internal/command"
	"github.com/metalsoft-io/metalcloud-cli/internal/configuration"
	"github.com/metalsoft-io/tableformatter"
)

// infrastructureCmds commands affecting infrastructures
var InfrastructureCmds = []command.Command{
	{
		Description:  "Creates an infrastructure.",
		Subject:      "infrastructure",
		AltSubject:   "infra",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("create infrastructure", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_label": c.FlagSet.String("label", "", colors.Red("(Required)")+" Infrastructure's label"),
				"datacenter":           c.FlagSet.String("datacenter", command.NilDefaultStr, colors.Red("(Required)")+" Infrastructure datacenter"),
				"return_id":            c.FlagSet.Bool("return-id", false, colors.Green("(Flag)")+" If set will print the ID of the created infrastructure. Useful for automating tasks."),
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
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"format":       c.FlagSet.String("format", "", "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
				"filter":       c.FlagSet.String("filter", "*", "filter to use when searching for servers. Check the documentation for examples. Defaults to '*'"),
				"show_ordered": c.FlagSet.Bool("show-ordered", false, colors.Green("(Flag)")+" If set will also return ordered (created but not deployed) infrastructures. Default is false."),
				"show_deleted": c.FlagSet.Bool("show-deleted", false, colors.Green("(Flag)")+" If set will also return deleted infrastructures. Default is false."),
			}
		},
		ExecuteFunc: infrastructureListAdminCmd,
		Endpoint:    configuration.DeveloperEndpoint,
		AdminOnly:   true,
	},
	{ // This is a second version of the list command for users. It uses a different function that needs a different set
		//of permissions and returns only the user's infrastructures instead of all user's infrastructures.
		Description:  "Lists all infrastructures.",
		Subject:      "infrastructure",
		AltSubject:   "infra",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list infrastructure", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"format":       c.FlagSet.String("format", "", "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
				"show_deleted": c.FlagSet.Bool("show-deleted", false, colors.Green("(Flag)")+" If set will also return deleted infrastructures. Default is false."),
			}
		},
		ExecuteFunc: infrastructureListUserCmd,
		Endpoint:    configuration.UserEndpoint,
		UserOnly:    true, //notice this
	},
	{
		Description:  "Delete an infrastructure.",
		Subject:      "infrastructure",
		AltSubject:   "infra",
		Predicate:    "delete",
		AltPredicate: "rm",
		FlagSet:      flag.NewFlagSet("delete infrastructure", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_id_or_label": c.FlagSet.String("id", command.NilDefaultStr, colors.Red("(Required)")+" Infrastructure's id or label. Note that using the 'label' might be ambiguous in certain situations."),
				"autoconfirm":                c.FlagSet.Bool("autoconfirm", false, colors.Green("(Flag)")+" If set it will assume action is confirmed"),
			}
		},
		ExecuteFunc:   infrastructureDeleteCmd,
		Endpoint:      configuration.UserEndpoint,
		AdminEndpoint: configuration.DeveloperEndpoint,
	},
	{
		Description:  "Deploy an infrastructure.",
		Subject:      "infrastructure",
		AltSubject:   "infra",
		Predicate:    "deploy",
		AltPredicate: "apply",
		FlagSet:      flag.NewFlagSet("deploy infrastructure", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_id_or_label":     c.FlagSet.String("id", command.NilDefaultStr, colors.Red("(Required)")+" Infrastructure's id or label. Note that using the 'label' might be ambiguous in certain situations."),
				"no_hard_shutdown_after_timeout": c.FlagSet.Bool("no-hard-shutdown-after-timeout", false, colors.Green("(Flag)")+" If set do not force a hard power off after timeout expired and the server is not powered off."),
				"no_attempt_soft_shutdown":       c.FlagSet.Bool("no-attempt-soft-shutdown", false, colors.Green("(Flag)")+" If set,do not atempt a soft (ACPI) power off of all the servers in the infrastructure before the deploy"),
				"soft_shutdown_timeout_seconds":  c.FlagSet.Int("soft-shutdown-timeout-seconds", 180, "(Optional, default 180) Timeout to wait if hard_shutdown_after_timeout is set."),
				"allow_data_loss":                c.FlagSet.Bool("allow-data-loss", false, colors.Green("(Flag)")+" If set, deploy will not throw error if data loss is expected."),
				"skip_ansible":                   c.FlagSet.Bool("skip-ansible", false, colors.Green("(Flag)")+" If set, some automatic provisioning steps will be skipped. This parameter should generally be ignored."),
				"block_until_deployed":           c.FlagSet.Bool("blocking", false, colors.Green("(Flag)")+" If set, the operation will wait until deployment finishes."),
				"block_timeout":                  c.FlagSet.Int("block-timeout", 180*60, "Block timeout in seconds. After this timeout the application will return an error. Defaults to 180 minutes."),
				"block_check_interval":           c.FlagSet.Int("block-check-interval", 10, "Check interval for when blocking. Defaults to 10 seconds."),
				"autoconfirm":                    c.FlagSet.Bool("autoconfirm", false, colors.Green("(Flag)")+" If set it will assume action is confirmed"),
			}
		},
		ExecuteFunc:   infrastructureDeployCmd,
		Endpoint:      configuration.UserEndpoint,
		AdminEndpoint: configuration.DeveloperEndpoint,
	},
	{ //User version
		Description:  "Get infrastructure details.",
		Subject:      "infrastructure",
		AltSubject:   "infra",
		Predicate:    "get",
		AltPredicate: "show",
		FlagSet:      flag.NewFlagSet("get infrastructure", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_id_or_label": c.FlagSet.String("id", command.NilDefaultStr, colors.Red("(Required)")+" Infrastructure's id or label. Note that using the 'label' might be ambiguous in certain situations."),
				"format":                     c.FlagSet.String("format", "", "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
			}
		},
		ExecuteFunc:   infrastructureGetCmd,
		Endpoint:      configuration.UserEndpoint,
		AdminEndpoint: configuration.DeveloperEndpoint,
	},
	{
		Description:  "Revert all changes of an infrastructure.",
		Subject:      "infrastructure",
		AltSubject:   "infra",
		Predicate:    "revert",
		AltPredicate: "undo",
		FlagSet:      flag.NewFlagSet("deploy infrastructure", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_id_or_label": c.FlagSet.String("id", command.NilDefaultStr, colors.Red("(Required)")+" Infrastructure's id or label. Note that using the 'label' might be ambiguous in certain situations."),
				"autoconfirm":                c.FlagSet.Bool("autoconfirm", false, colors.Green("(Flag)")+" If set it will assume action is confirmed"),
			}
		},
		ExecuteFunc:   infrastructureRevertCmd,
		Endpoint:      configuration.UserEndpoint,
		AdminEndpoint: configuration.DeveloperEndpoint,
	},
	{
		Description:  "List stages of a workflow.",
		Subject:      "infrastructure",
		AltSubject:   "infra",
		Predicate:    "workflow-stages",
		AltPredicate: "workflow-stages",
		FlagSet:      flag.NewFlagSet("List stages of a workflow.", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_id_or_label": c.FlagSet.String("id", command.NilDefaultStr, "The infrastructure's id"),
				"type":                       c.FlagSet.String("type", command.NilDefaultStr, "stage definition type. possible values: pre_deploy, post_deploy"),
			}
		},
		ExecuteFunc:   listWorkflowStagesCmd,
		Endpoint:      configuration.DeveloperEndpoint,
		AdminEndpoint: configuration.DeveloperEndpoint,
	},
}

func infrastructureCreateCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

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

func infrastructureListAdminCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	filter := command.GetStringParam(c.Arguments["filter"])
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

		if i.InfrastructureServiceStatus == "ordered" && !command.GetBoolParam(c.Arguments["show_ordered"]) {
			continue
		}

		if i.InfrastructureServiceStatus == "deleted" && !command.GetBoolParam(c.Arguments["show_deleted"]) {
			continue
		}

		status := ""

		if i.InfrastructureServiceStatus == "active" && i.AFCExecutedSuccess == i.AFCTotal {
			status = colors.Green("Deployed")
		}

		if i.InfrastructureServiceStatus == "ordered" && i.AFCTotal == 0 {
			status = colors.Blue("Ordered (deploy not started)")
		}

		if i.AFCExecutedSuccess < i.AFCTotal {

			if i.AFCThrownError == 0 {
				status = colors.Yellow(fmt.Sprintf("Deploy ongoing - %d/%d", i.AFCExecutedSuccess, i.AFCTotal))
			} else {
				status = colors.Red(fmt.Sprintf("Deploy ongoing - Thrown error at %d/%d", i.AFCExecutedSuccess, i.AFCTotal))
			}
		}

		if i.InfrastructureServiceStatus == "deleted" {
			status = colors.Magenta("Deleted")
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
	return table.RenderTable("Infrastructures", topLine, command.GetStringParam(c.Arguments["format"]))
}

func infrastructureListUserCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	iList, err := client.Infrastructures()
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

		if i.InfrastructureServiceStatus == "deleted" && !command.GetBoolParam(c.Arguments["show_deleted"]) {
			continue
		}

		status := ""

		if i.InfrastructureServiceStatus == "active" {
			status = colors.Green("Deployed")
		}

		if i.InfrastructureServiceStatus == "ordered" {
			status = colors.Blue("Ordered (deploy not started)")
		}

		if i.InfrastructureServiceStatus == "deleted" {
			status = colors.Magenta("Deleted")
		}

		data = append(data, []interface{}{
			i.InfrastructureID,
			i.InfrastructureLabel,
			status,
			i.UserEmailOwner,
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
	return table.RenderTable("Infrastructures", topLine, command.GetStringParam(c.Arguments["format"]))
}

type infrastructureConfirmAndDoFunc func(infraID int, c *command.Command, client metalcloud.MetalCloudClient) (string, error)

// infrastructureConfirmAndDo asks for confirmation and executes the given function
func infrastructureConfirmAndDo(operation string, c *command.Command, client metalcloud.MetalCloudClient, f infrastructureConfirmAndDoFunc) (string, error) {

	val, err := command.GetParam(c, "infrastructure_id_or_label", "id")
	if err != nil {
		return "", err
	}

	infraID, err := command.GetIDOrDo(*val.(*string), func(label string) (int, error) {
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

	if command.GetBoolParam(c.Arguments["autoconfirm"]) {
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

		confirm, err = command.RequestConfirmation(confirmationMessage)
		if err != nil {
			return "", err
		}
	}

	if !confirm {
		return "", fmt.Errorf("Operation not confirmed. Aborting")
	}

	return f(infraID, c, client)
}

func infrastructureDeleteCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	return infrastructureConfirmAndDo("Delete", c, client,
		func(infraID int, c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
			return "", client.InfrastructureDelete(infraID)
		})
}

func infrastructureDeployCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	return infrastructureConfirmAndDo("Deploy", c, client,
		func(infraID int, c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

			shutDownOptions := metalcloud.ShutdownOptions{
				HardShutdownAfterTimeout:   !command.GetBoolParam(c.Arguments["no_hard_shutdown_after_timeout"]),
				AttemptSoftShutdown:        !command.GetBoolParam(c.Arguments["no_attempt_soft_shutdown"]),
				SoftShutdownTimeoutSeconds: command.GetIntParam(c.Arguments["soft_shutdown_timeout_seconds"]),
			}

			err := client.InfrastructureDeploy(
				infraID,
				shutDownOptions,
				command.GetBoolParam(c.Arguments["allow_data_loss"]),
				command.GetBoolParam(c.Arguments["skip_ansible"]),
			)
			if err != nil {
				return "", err
			}

			if command.GetBoolParam(c.Arguments["block_until_deployed"]) {

				time.Sleep(time.Duration(command.GetIntParam(c.Arguments["block_check_interval"])) * time.Second) //wait until the system picks up the afc

				err := loopUntilInfraReady(infraID, command.GetIntParam(c.Arguments["block_timeout"]), command.GetIntParam(c.Arguments["block_check_interval"]), client)

				if err != nil && strings.HasPrefix(err.Error(), "timeout after") {
					return "", err
				} //else we ignore errors as they might be infrastrucure not found due to infrastructure being deleted
			}

			return "", nil
		})
}

func infrastructureRevertCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	return infrastructureConfirmAndDo("Revert", c, client,
		func(infraID int, c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
			return "", client.InfrastructureOperationCancel(infraID)
		})
}

func infrastructureGetCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	retInfra, err := command.GetInfrastructureFromCommand("id", c, client)
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
		status := colors.Green(ia.InstanceArrayServiceStatus)
		if ia.InstanceArrayServiceStatus != "ordered" && ia.InstanceArrayOperation.InstanceArrayDeployType == "edit" && ia.InstanceArrayOperation.InstanceArrayDeployStatus == "not_started" {
			status = colors.Blue("edited")
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
			colors.Green("InstanceArray"),
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
		status := colors.Green(da.DriveArrayServiceStatus)
		if da.DriveArrayServiceStatus != "ordered" && da.DriveArrayOperation.DriveArrayDeployType == "edit" && da.DriveArrayOperation.DriveArrayDeployStatus == "not_started" {
			status = colors.Blue("edited")
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
			colors.Blue("DriveArray"),
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
		status := colors.Green(sda.SharedDriveServiceStatus)
		if sda.SharedDriveServiceStatus != "ordered" && sda.SharedDriveOperation.SharedDriveDeployType == "edit" && sda.SharedDriveOperation.SharedDriveDeployStatus == "not_started" {
			status = colors.Blue("edited")
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
			colors.Magenta("SharedDrive"),
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
	return table.RenderTable("resources", topLine, command.GetStringParam(c.Arguments["format"]))
}

func listWorkflowStagesCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	t := *c.Arguments["type"].(*string)
	if t == command.NilDefaultStr {
		t = "post_deploy"
	}

	retInfra, err := command.GetInfrastructureFromCommand("id", c, client)
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
	return table.RenderTable("Workflow Stages", "", command.GetStringParam(c.Arguments["format"]))
}

// loop until infra is ready
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
