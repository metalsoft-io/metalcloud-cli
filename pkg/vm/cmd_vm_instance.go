package vm

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"

	metalcloud2 "github.com/metalsoft-io/metal-cloud-sdk2-go"
	"github.com/metalsoft-io/metalcloud-cli/internal/colors"
	"github.com/metalsoft-io/metalcloud-cli/internal/command"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/metalsoft-io/tableformatter"
	"gopkg.in/yaml.v3"
)

var VmInstancesCmds = []command.Command{
	{
		Description:  "Get VM instance.",
		Subject:      "vm-instance",
		AltSubject:   "vm",
		Predicate:    "get",
		AltPredicate: "show",
		FlagSet:      flag.NewFlagSet("get VM instance", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_id": c.FlagSet.String("infra", command.NilDefaultStr, colors.Red("(Required)")+" Infrastructure's Id."),
				"vm_instance_id":    c.FlagSet.Int("id", 0, "VM instance Id"),
				"format":            c.FlagSet.String("format", command.NilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
			}
		},
		ExecuteFunc2:        vmInstanceGetCmd,
		PermissionsRequired: []string{command.VMS_READ},
	},
	{
		Description:  "Create VM instance.",
		Subject:      "vm-instance",
		AltSubject:   "vm",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("create VM instance", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_id":     c.FlagSet.String("infra", command.NilDefaultStr, colors.Red("(Required)")+" Infrastructure's Id."),
				"format":                c.FlagSet.String("format", "json", "The input format. Supported values are 'json','yaml'. The default format is json."),
				"read_config_from_file": c.FlagSet.String("raw-config", command.NilDefaultStr, colors.Red("(Required)")+" Read raw object from file"),
				"read_config_from_pipe": c.FlagSet.Bool("pipe", false, colors.Green("(Flag)")+" If set, read raw object from pipe instead of from a file. Either this flag or the --raw-config option must be used."),
				"return_id":             c.FlagSet.Bool("return-id", false, "Will print the ID of the created object. Useful for automating tasks."),
			}
		},
		ExecuteFunc2:        vmInstanceCreateCmd,
		PermissionsRequired: []string{command.VMS_WRITE},
	},
	{
		Description:  "Update VM instance.",
		Subject:      "vm-instance",
		AltSubject:   "vm",
		Predicate:    "edit",
		AltPredicate: "update",
		FlagSet:      flag.NewFlagSet("update VM instance", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_id":     c.FlagSet.String("infra", command.NilDefaultStr, colors.Red("(Required)")+" Infrastructure's Id."),
				"vm_instance_id":        c.FlagSet.Int("id", 0, "VM instance Id"),
				"format":                c.FlagSet.String("format", "json", "The input format. Supported values are 'json','yaml'. The default format is json."),
				"read_config_from_file": c.FlagSet.String("raw-config", command.NilDefaultStr, colors.Red("(Required)")+" Read raw object from file"),
				"read_config_from_pipe": c.FlagSet.Bool("pipe", false, colors.Green("(Flag)")+" If set, read raw object from pipe instead of from a file. Either this flag or the --raw-config option must be used."),
			}
		},
		ExecuteFunc2:        vmInstanceUpdateCmd,
		PermissionsRequired: []string{command.VMS_WRITE},
	},
	{
		Description:  "Delete VM instance.",
		Subject:      "vm-instance",
		AltSubject:   "vm",
		Predicate:    "delete",
		AltPredicate: "rm",
		FlagSet:      flag.NewFlagSet("delete VM instance", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_id": c.FlagSet.String("infra", command.NilDefaultStr, colors.Red("(Required)")+" Infrastructure's Id."),
				"vm_instance_id":    c.FlagSet.Int("id", 0, "VM instance Id"),
				"autoconfirm":       c.FlagSet.Bool("autoconfirm", false, colors.Green("(Flag)")+" If set it will assume action is confirmed"),
			}
		},
		ExecuteFunc2:        vmInstanceDeleteCmd,
		PermissionsRequired: []string{command.VMS_WRITE},
	},
	{
		Description:  "Change VM instance type.",
		Subject:      "vm-instance",
		AltSubject:   "vm",
		Predicate:    "change-type",
		AltPredicate: "ct",
		FlagSet:      flag.NewFlagSet("change VM instance type", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_id": c.FlagSet.String("infra", command.NilDefaultStr, colors.Red("(Required)")+" Infrastructure's Id."),
				"vm_instance_id":    c.FlagSet.Int("id", 0, "VM instance Id"),
				"vm_type_id":        c.FlagSet.Int("type-id", 0, "VM instance type Id"),
			}
		},
		ExecuteFunc2:        vmInstanceTypeChangeCmd,
		PermissionsRequired: []string{command.VMS_WRITE},
	},
	{
		Description:  "Get VM instance power status.",
		Subject:      "vm-instance",
		AltSubject:   "vm",
		Predicate:    "power-status",
		AltPredicate: "power",
		FlagSet:      flag.NewFlagSet("get VM instance power status", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_id": c.FlagSet.String("infra", command.NilDefaultStr, colors.Red("(Required)")+" Infrastructure's Id."),
				"vm_instance_id":    c.FlagSet.Int("id", 0, "VM instance Id"),
			}
		},
		ExecuteFunc2:        vmInstancePowerStatusGetCmd,
		PermissionsRequired: []string{command.VMS_READ},
	},
	{
		Description:  "Start VM instance.",
		Subject:      "vm-instance",
		AltSubject:   "vm",
		Predicate:    "start",
		AltPredicate: "run",
		FlagSet:      flag.NewFlagSet("start VM instance", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_id": c.FlagSet.String("infra", command.NilDefaultStr, colors.Red("(Required)")+" Infrastructure's Id."),
				"vm_instance_id":    c.FlagSet.Int("id", 0, "VM instance Id"),
				"autoconfirm":       c.FlagSet.Bool("autoconfirm", false, colors.Green("(Flag)")+" If set it will assume action is confirmed"),
			}
		},
		ExecuteFunc2:        vmInstanceStartCmd,
		PermissionsRequired: []string{command.VMS_WRITE},
	},
	{
		Description:  "Reboot VM instance.",
		Subject:      "vm-instance",
		AltSubject:   "vm",
		Predicate:    "reboot",
		AltPredicate: "reboot",
		FlagSet:      flag.NewFlagSet("reboot VM instance", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_id": c.FlagSet.String("infra", command.NilDefaultStr, colors.Red("(Required)")+" Infrastructure's Id."),
				"vm_instance_id":    c.FlagSet.Int("id", 0, "VM instance Id"),
				"autoconfirm":       c.FlagSet.Bool("autoconfirm", false, colors.Green("(Flag)")+" If set it will assume action is confirmed"),
			}
		},
		ExecuteFunc2:        vmInstanceRebootCmd,
		PermissionsRequired: []string{command.VMS_WRITE},
	},
	{
		Description:  "Shutdown VM instance.",
		Subject:      "vm-instance",
		AltSubject:   "vm",
		Predicate:    "shutdown",
		AltPredicate: "stop",
		FlagSet:      flag.NewFlagSet("shutdown VM instance", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_id": c.FlagSet.String("infra", command.NilDefaultStr, colors.Red("(Required)")+" Infrastructure's Id."),
				"vm_instance_id":    c.FlagSet.Int("id", 0, "VM instance Id"),
				"autoconfirm":       c.FlagSet.Bool("autoconfirm", false, colors.Green("(Flag)")+" If set it will assume action is confirmed"),
			}
		},
		ExecuteFunc2:        vmInstanceShutdownCmd,
		PermissionsRequired: []string{command.VMS_WRITE},
	},
}

func vmInstanceGetCmd(ctx context.Context, c *command.Command, client *metalcloud2.APIClient) (string, error) {
	infrastructureId, ok := command.GetIntParamOk(c.Arguments["infrastructure_id"])
	if !ok {
		return "", fmt.Errorf("-infra is required")
	}

	id, ok := command.GetIntParamOk(c.Arguments["vm_instance_id"])
	if !ok {
		return "", fmt.Errorf("-id is required")
	}

	vm, response, err := client.VMInstanceApi.GetVMInstance(ctx, float64(infrastructureId), float64(id))
	if err != nil {
		return "", err
	}
	if response.StatusCode >= 400 {
		return "", fmt.Errorf("HTTP error: %s", response.Status)
	}

	format := command.GetStringParam(c.Arguments["format"])

	vm.Links = nil

	switch format {
	case "json", "JSON":
		result, err := json.MarshalIndent(vm, "", "\t")
		if err != nil {
			return "", err
		}

		return string(result), nil

	case "yaml", "YAML":
		result, err := yaml.Marshal(vm)
		if err != nil {
			return "", err
		}

		return string(result), nil

	default:
		formattedData := [][]interface{}{}
		formattedData = append(formattedData, formattedVmInstanceRecord(vm))

		table := tableformatter.Table{
			Data:   formattedData,
			Schema: vmInstanceFieldSchema(),
		}

		return table.RenderTable("VM", "", format)
	}
}

func vmInstanceCreateCmd(ctx context.Context, c *command.Command, client *metalcloud2.APIClient) (string, error) {
	infrastructureId, ok := command.GetIntParamOk(c.Arguments["infrastructure_id"])
	if !ok {
		return "", fmt.Errorf("-infra is required")
	}

	var obj metalcloud2.CreateVmInstance

	err := command.GetRawObjectFromCommand(c, &obj)
	if err != nil {
		return "", err
	}

	vm, _, err := client.VMInstanceApi.CreateVMInstance(ctx, obj, float64(infrastructureId))
	if err != nil {
		return "", err
	}

	if command.GetBoolParam(c.Arguments["return_id"]) {
		return fmt.Sprintf("%v", vm.Id), nil
	}

	return "", err
}

func vmInstanceUpdateCmd(ctx context.Context, c *command.Command, client *metalcloud2.APIClient) (string, error) {
	infrastructureId, ok := command.GetIntParamOk(c.Arguments["infrastructure_id"])
	if !ok {
		return "", fmt.Errorf("-infra is required")
	}

	id, ok := command.GetIntParamOk(c.Arguments["vm_instance_id"])
	if !ok {
		return "", fmt.Errorf("-id is required")
	}

	var obj metalcloud2.UpdateVmInstance

	err := command.GetRawObjectFromCommand(c, &obj)
	if err != nil {
		return "", err
	}

	_, _, err = client.VMInstanceApi.UpdateVMInstance(ctx, obj, float64(infrastructureId), float64(id))
	if err != nil {
		return "", err
	}

	return "", err
}

func vmInstanceDeleteCmd(ctx context.Context, c *command.Command, client *metalcloud2.APIClient) (string, error) {
	infrastructureId, ok := command.GetIntParamOk(c.Arguments["infrastructure_id"])
	if !ok {
		return "", fmt.Errorf("-infra is required")
	}

	id, ok := command.GetIntParamOk(c.Arguments["vm_instance_id"])
	if !ok {
		return "", fmt.Errorf("-id is required")
	}

	confirmed, err := utils.GetConfirmation(command.GetBoolParam(c.Arguments["autoconfirm"]), fmt.Sprintf("Deleting VM instance #%d.", id))
	if err != nil {
		return "", err
	}

	if !confirmed {
		return "", fmt.Errorf("Operation not confirmed. Aborting")
	}

	_, err = client.VMInstanceApi.DeleteVMInstance(ctx, float64(infrastructureId), float64(id))
	if err != nil {
		return "", err
	}

	return "", err
}

func vmInstanceTypeChangeCmd(ctx context.Context, c *command.Command, client *metalcloud2.APIClient) (string, error) {
	infrastructureId, ok := command.GetIntParamOk(c.Arguments["infrastructure_id"])
	if !ok {
		return "", fmt.Errorf("-infra is required")
	}

	id, ok := command.GetIntParamOk(c.Arguments["vm_instance_id"])
	if !ok {
		return "", fmt.Errorf("-id is required")
	}

	typeId, ok := command.GetIntParamOk(c.Arguments["vm_type_id"])
	if !ok {
		return "", fmt.Errorf("-type-id is required")
	}

	_, _, err := client.VMInstanceApi.ApplyVMTypeOnVMInstance(ctx, float64(infrastructureId), float64(id), float64(typeId))
	if err != nil {
		return "", err
	}

	return "", err
}

func vmInstancePowerStatusGetCmd(ctx context.Context, c *command.Command, client *metalcloud2.APIClient) (string, error) {
	infrastructureId, ok := command.GetIntParamOk(c.Arguments["infrastructure_id"])
	if !ok {
		return "", fmt.Errorf("-infra is required")
	}

	id, ok := command.GetIntParamOk(c.Arguments["vm_instance_id"])
	if !ok {
		return "", fmt.Errorf("-id is required")
	}

	response, err := client.VMInstanceApi.GetVMInstancePowerStatus(ctx, float64(infrastructureId), float64(id))
	if err != nil {
		return "", err
	}
	if response.StatusCode >= 400 {
		return "", fmt.Errorf("HTTP error: %s", response.Status)
	}

	powerStatus, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(powerStatus), nil
}

func vmInstanceStartCmd(ctx context.Context, c *command.Command, client *metalcloud2.APIClient) (string, error) {
	infrastructureId, ok := command.GetIntParamOk(c.Arguments["infrastructure_id"])
	if !ok {
		return "", fmt.Errorf("-infra is required")
	}

	id, ok := command.GetIntParamOk(c.Arguments["vm_instance_id"])
	if !ok {
		return "", fmt.Errorf("-id is required")
	}

	confirmed, err := utils.GetConfirmation(command.GetBoolParam(c.Arguments["autoconfirm"]), fmt.Sprintf("Starting VM instance #%d.", id))
	if err != nil {
		return "", err
	}

	if !confirmed {
		return "", fmt.Errorf("Operation not confirmed. Aborting")
	}

	_, err = client.VMInstanceApi.StartVMInstance(ctx, float64(infrastructureId), float64(id))
	if err != nil {
		return "", err
	}

	return "", err
}

func vmInstanceRebootCmd(ctx context.Context, c *command.Command, client *metalcloud2.APIClient) (string, error) {
	infrastructureId, ok := command.GetIntParamOk(c.Arguments["infrastructure_id"])
	if !ok {
		return "", fmt.Errorf("-infra is required")
	}

	id, ok := command.GetIntParamOk(c.Arguments["vm_instance_id"])
	if !ok {
		return "", fmt.Errorf("-id is required")
	}

	confirmed, err := utils.GetConfirmation(command.GetBoolParam(c.Arguments["autoconfirm"]), fmt.Sprintf("Rebooting VM instance #%d.", id))
	if err != nil {
		return "", err
	}

	if !confirmed {
		return "", fmt.Errorf("Operation not confirmed. Aborting")
	}

	_, err = client.VMInstanceApi.RebootVMInstance(ctx, float64(infrastructureId), float64(id))
	if err != nil {
		return "", err
	}

	return "", err
}

func vmInstanceShutdownCmd(ctx context.Context, c *command.Command, client *metalcloud2.APIClient) (string, error) {
	infrastructureId, ok := command.GetIntParamOk(c.Arguments["infrastructure_id"])
	if !ok {
		return "", fmt.Errorf("-infra is required")
	}

	id, ok := command.GetIntParamOk(c.Arguments["vm_instance_id"])
	if !ok {
		return "", fmt.Errorf("-id is required")
	}

	confirmed, err := utils.GetConfirmation(command.GetBoolParam(c.Arguments["autoconfirm"]), fmt.Sprintf("Shutting down VM instance #%d.", id))
	if err != nil {
		return "", err
	}

	if !confirmed {
		return "", fmt.Errorf("Operation not confirmed. Aborting")
	}

	_, err = client.VMInstanceApi.ShutdownVMInstance(ctx, float64(infrastructureId), float64(id))
	if err != nil {
		return "", err
	}

	return "", err
}

func vmInstanceFieldSchema() []tableformatter.SchemaField {
	schema := []tableformatter.SchemaField{
		{
			FieldName: "ID",
			FieldType: tableformatter.TypeFloat,
			FieldSize: 6,
		},

		{
			FieldName: "LABEL",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},

		{
			FieldName: "TYPE",
			FieldType: tableformatter.TypeInt,
			FieldSize: 5,
		},

		{
			FieldName: "DISK",
			FieldType: tableformatter.TypeFloat,
			FieldSize: 5,
		},
	}

	return schema
}

func formattedVmInstanceRecord(vm metalcloud2.VmInstance) []interface{} {
	formattedRecord := []interface{}{
		vm.Id,
		vm.Label,
		int(vm.TypeId),
		vm.DiskSizeGB,
	}

	return formattedRecord
}
