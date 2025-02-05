package vm

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"

	metalcloud2 "github.com/metalsoft-io/metal-cloud-sdk2-go"
	"github.com/metalsoft-io/metalcloud-cli/internal/colors"
	"github.com/metalsoft-io/metalcloud-cli/internal/command"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/metalsoft-io/tableformatter"
	"gopkg.in/yaml.v3"
)

var VmInstanceGroupsCmds = []command.Command{
	{
		Description:  "Lists all VM instance groups of an infrastructure.",
		Subject:      "vm-instance-group",
		AltSubject:   "vm-ig",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list VM instance groups", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_id": c.FlagSet.Int("infra", command.NilDefaultInt, colors.Red("(Required)")+" Infrastructure's Id."),
				"format":            c.FlagSet.String("format", command.NilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
			}
		},
		ExecuteFunc2:        vmInstanceGroupListCmd,
		PermissionsRequired: []string{command.VMS_READ},
		MinApiVersion:       "v6.3",
	},
	{
		Description:  "Get VM instance group details.",
		Subject:      "vm-instance-group",
		AltSubject:   "vm-ig",
		Predicate:    "get",
		AltPredicate: "show",
		FlagSet:      flag.NewFlagSet("get VM instance group", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_id":    c.FlagSet.Int("infra", command.NilDefaultInt, colors.Red("(Required)")+" Infrastructure's Id."),
				"vm_instance_group_id": c.FlagSet.Int("id", 0, "VM instance group Id"),
				"format":               c.FlagSet.String("format", command.NilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
			}
		},
		ExecuteFunc2:        vmInstanceGroupGetCmd,
		PermissionsRequired: []string{command.VMS_READ},
		MinApiVersion:       "v6.3",
	},
	{
		Description:  "List VM instance group VMs.",
		Subject:      "vm-instance-group",
		AltSubject:   "vm-ig",
		Predicate:    "vms-list",
		AltPredicate: "vms-ls",
		FlagSet:      flag.NewFlagSet("list VM instance group VMs", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_id":    c.FlagSet.Int("infra", command.NilDefaultInt, colors.Red("(Required)")+" Infrastructure's Id."),
				"vm_instance_group_id": c.FlagSet.Int("id", 0, "VM instance group Id"),
				"format":               c.FlagSet.String("format", command.NilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
			}
		},
		ExecuteFunc2:        vmInstanceGroupVmsListCmd,
		PermissionsRequired: []string{command.VMS_READ},
		MinApiVersion:       "v6.3",
	},
	{
		Description:  "Create VM instance group.",
		Subject:      "vm-instance-group",
		AltSubject:   "vm-ig",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("create VM instance group", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_id":     c.FlagSet.Int("infra", command.NilDefaultInt, colors.Red("(Required)")+" Infrastructure's Id."),
				"format":                c.FlagSet.String("format", "json", "The input format. Supported values are 'json','yaml'. The default format is json."),
				"read_config_from_file": c.FlagSet.String("raw-config", command.NilDefaultStr, colors.Red("(Required)")+" Read raw object from file"),
				"read_config_from_pipe": c.FlagSet.Bool("pipe", false, colors.Green("(Flag)")+" If set, read raw object from pipe instead of from a file. Either this flag or the --raw-config option must be used."),
				"return_id":             c.FlagSet.Bool("return-id", false, "Will print the Id of the created object. Useful for automating tasks."),
			}
		},
		ExecuteFunc2:        vmInstanceGroupCreateCmd,
		PermissionsRequired: []string{command.VMS_WRITE},
		MinApiVersion:       "v6.3",
	},
	{
		Description:  "Update VM instance group.",
		Subject:      "vm-instance-group",
		AltSubject:   "vm-ig",
		Predicate:    "edit",
		AltPredicate: "update",
		FlagSet:      flag.NewFlagSet("update VM instance group", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_id":     c.FlagSet.Int("infra", command.NilDefaultInt, colors.Red("(Required)")+" Infrastructure's Id."),
				"vm_instance_group_id":  c.FlagSet.Int("id", 0, "VM instance group Id"),
				"format":                c.FlagSet.String("format", "json", "The input format. Supported values are 'json','yaml'. The default format is json."),
				"read_config_from_file": c.FlagSet.String("raw-config", command.NilDefaultStr, colors.Red("(Required)")+" Read raw object from file"),
				"read_config_from_pipe": c.FlagSet.Bool("pipe", false, colors.Green("(Flag)")+" If set, read raw object from pipe instead of from a file. Either this flag or the --raw-config option must be used."),
			}
		},
		ExecuteFunc2:        vmInstanceGroupUpdateCmd,
		PermissionsRequired: []string{command.VMS_WRITE},
		MinApiVersion:       "v6.3",
	},
	{
		Description:  "Delete VM instance group.",
		Subject:      "vm-instance-group",
		AltSubject:   "vm-ig",
		Predicate:    "delete",
		AltPredicate: "rm",
		FlagSet:      flag.NewFlagSet("delete VM instance group", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_id":    c.FlagSet.Int("infra", command.NilDefaultInt, colors.Red("(Required)")+" Infrastructure's Id."),
				"vm_instance_group_id": c.FlagSet.Int("id", 0, "VM instance group Id"),
				"autoconfirm":          c.FlagSet.Bool("autoconfirm", false, colors.Green("(Flag)")+" If set it will assume action is confirmed"),
			}
		},
		ExecuteFunc2:        vmInstanceGroupDeleteCmd,
		PermissionsRequired: []string{command.VMS_WRITE},
		MinApiVersion:       "v6.3",
	},
}

func vmInstanceGroupListCmd(ctx context.Context, c *command.Command, client *metalcloud2.APIClient) (string, error) {
	infrastructureId, ok := command.GetIntParamOk(c.Arguments["infrastructure_id"])
	if !ok {
		return "", fmt.Errorf("-infra is required")
	}

	vmGroups, response, err := client.VMInstanceGroupApi.GetVMInstanceGroups(ctx, float64(infrastructureId))
	if err != nil {
		return "", err
	}
	if response.StatusCode >= 400 {
		return "", fmt.Errorf("HTTP error: %s", response.Status)
	}

	format := command.GetStringParam(c.Arguments["format"])

	rawData := []metalcloud2.VmInstanceGroup{}
	formattedData := [][]interface{}{}

	for _, vmGroup := range vmGroups {
		vmGroup.Links = nil

		rawData = append(rawData, vmGroup)

		formattedData = append(formattedData, formattedVmInstanceGroupRecord(vmGroup))
	}

	switch format {
	case "json", "JSON":
		result, err := json.MarshalIndent(rawData, "", "\t")
		if err != nil {
			return "", err
		}

		return string(result), nil

	case "yaml", "YAML":
		result, err := yaml.Marshal(rawData)
		if err != nil {
			return "", err
		}

		return string(result), nil

	default:
		table := tableformatter.Table{
			Data:   formattedData,
			Schema: vmInstanceGroupFieldSchema(),
		}

		return table.RenderTable("VM Instance Groups", "", format)
	}
}

func vmInstanceGroupGetCmd(ctx context.Context, c *command.Command, client *metalcloud2.APIClient) (string, error) {
	infrastructureId, ok := command.GetIntParamOk(c.Arguments["infrastructure_id"])
	if !ok {
		return "", fmt.Errorf("-infra is required")
	}

	id, ok := command.GetIntParamOk(c.Arguments["vm_instance_group_id"])
	if !ok {
		return "", fmt.Errorf("-id is required")
	}

	vmGroup, response, err := client.VMInstanceGroupApi.GetVMInstanceGroup(ctx, float64(infrastructureId), float64(id))
	if err != nil {
		return "", err
	}
	if response.StatusCode >= 400 {
		return "", fmt.Errorf("HTTP error: %s", response.Status)
	}

	format := command.GetStringParam(c.Arguments["format"])

	vmGroup.Links = nil

	switch format {
	case "json", "JSON":
		result, err := json.MarshalIndent(vmGroup, "", "\t")
		if err != nil {
			return "", err
		}

		return string(result), nil

	case "yaml", "YAML":
		result, err := yaml.Marshal(vmGroup)
		if err != nil {
			return "", err
		}

		return string(result), nil

	default:
		formattedData := [][]interface{}{}
		formattedData = append(formattedData, formattedVmInstanceGroupRecord(vmGroup))

		table := tableformatter.Table{
			Data:   formattedData,
			Schema: vmInstanceGroupFieldSchema(),
		}

		return table.RenderTable("VM", "", format)
	}
}

func vmInstanceGroupVmsListCmd(ctx context.Context, c *command.Command, client *metalcloud2.APIClient) (string, error) {
	infrastructureId, ok := command.GetIntParamOk(c.Arguments["infrastructure_id"])
	if !ok {
		return "", fmt.Errorf("-infra is required")
	}

	id, ok := command.GetIntParamOk(c.Arguments["vm_instance_group_id"])
	if !ok {
		return "", fmt.Errorf("-id is required")
	}

	vms, response, err := client.VMInstanceGroupApi.GetVMInstanceGroupVMInstances(ctx, float64(infrastructureId), float64(id))
	if err != nil {
		return "", err
	}
	if response.StatusCode >= 400 {
		return "", fmt.Errorf("HTTP error: %s", response.Status)
	}

	format := command.GetStringParam(c.Arguments["format"])

	rawData := []metalcloud2.VmInstance{}
	formattedData := [][]interface{}{}

	for _, vm := range vms {
		vm.Links = nil

		rawData = append(rawData, vm)

		formattedData = append(formattedData, formattedVmInstanceRecord(vm))
	}

	switch format {
	case "json", "JSON":
		result, err := json.MarshalIndent(rawData, "", "\t")
		if err != nil {
			return "", err
		}

		return string(result), nil

	case "yaml", "YAML":
		result, err := yaml.Marshal(rawData)
		if err != nil {
			return "", err
		}

		return string(result), nil

	default:
		table := tableformatter.Table{
			Data:   formattedData,
			Schema: vmInstanceFieldSchema(),
		}

		return table.RenderTable("VM Instance Group VMs", "", format)
	}
}

func vmInstanceGroupCreateCmd(ctx context.Context, c *command.Command, client *metalcloud2.APIClient) (string, error) {
	infrastructureId, ok := command.GetIntParamOk(c.Arguments["infrastructure_id"])
	if !ok {
		return "", fmt.Errorf("-infra is required")
	}

	var obj metalcloud2.CreateVmInstanceGroup

	err := command.GetRawObjectFromCommand(c, &obj)
	if err != nil {
		return "", err
	}

	vmGroup, _, err := client.VMInstanceGroupApi.CreateVMInstanceGroup(ctx, obj, float64(infrastructureId))
	if err != nil {
		return "", err
	}

	if command.GetBoolParam(c.Arguments["return_id"]) {
		return fmt.Sprintf("%v", vmGroup.Id), nil
	}

	return "", err
}

func vmInstanceGroupUpdateCmd(ctx context.Context, c *command.Command, client *metalcloud2.APIClient) (string, error) {
	infrastructureId, ok := command.GetIntParamOk(c.Arguments["infrastructure_id"])
	if !ok {
		return "", fmt.Errorf("-infra is required")
	}

	id, ok := command.GetIntParamOk(c.Arguments["vm_instance_group_id"])
	if !ok {
		return "", fmt.Errorf("-id is required")
	}

	var obj metalcloud2.UpdateVmInstanceGroup

	err := command.GetRawObjectFromCommand(c, &obj)
	if err != nil {
		return "", err
	}

	_, _, err = client.VMInstanceGroupApi.UpdateVMInstanceGroup(ctx, obj, float64(infrastructureId), float64(id))
	if err != nil {
		return "", err
	}

	return "", err
}

func vmInstanceGroupDeleteCmd(ctx context.Context, c *command.Command, client *metalcloud2.APIClient) (string, error) {
	infrastructureId, ok := command.GetIntParamOk(c.Arguments["infrastructure_id"])
	if !ok {
		return "", fmt.Errorf("-infra is required")
	}

	id, ok := command.GetIntParamOk(c.Arguments["vm_instance_group_id"])
	if !ok {
		return "", fmt.Errorf("-id is required")
	}

	confirmed, err := utils.GetConfirmation(command.GetBoolParam(c.Arguments["autoconfirm"]), fmt.Sprintf("Deleting VM instance group #%d.", id))
	if err != nil {
		return "", err
	}

	if !confirmed {
		return "", fmt.Errorf("Operation not confirmed. Aborting")
	}

	_, err = client.VMInstanceGroupApi.DeleteVMInstanceGroup(ctx, float64(infrastructureId), float64(id))
	if err != nil {
		return "", err
	}

	return "", err
}

func vmInstanceGroupFieldSchema() []tableformatter.SchemaField {
	schema := []tableformatter.SchemaField{
		{
			FieldName: "ID",
			FieldType: tableformatter.TypeFloat,
			FieldSize: 6,
		},

		{
			FieldName: "INFRA",
			FieldType: tableformatter.TypeFloat,
			FieldSize: 6,
		},

		{
			FieldName: "COUNT",
			FieldType: tableformatter.TypeInt,
			FieldSize: 5,
		},

		{
			FieldName: "STATUS",
			FieldType: tableformatter.TypeFloat,
			FieldSize: 5,
		},
	}

	return schema
}

func formattedVmInstanceGroupRecord(vmGroup metalcloud2.VmInstanceGroup) []interface{} {
	formattedRecord := []interface{}{
		vmGroup.Id,
		vmGroup.InfrastructureId,
		int(vmGroup.InstanceCount),
		vmGroup.ServiceStatus,
	}

	return formattedRecord
}
