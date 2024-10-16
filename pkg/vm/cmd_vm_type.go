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

var VmTypesCmds = []command.Command{
	{
		Description:  "Lists all VM types.",
		Subject:      "vm-type",
		AltSubject:   "vm-types",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list VM types", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"format": c.FlagSet.String("format", command.NilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
			}
		},
		ExecuteFunc2:        vmTypeListCmd,
		PermissionsRequired: []string{command.VM_TYPES_READ},
	},
	{
		Description:  "Get VM type.",
		Subject:      "vm-type",
		AltSubject:   "vm-type",
		Predicate:    "get",
		AltPredicate: "show",
		FlagSet:      flag.NewFlagSet("get VM type", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"vm_type_id": c.FlagSet.Int("id", 0, "VM type id"),
				"format":     c.FlagSet.String("format", command.NilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
			}
		},
		ExecuteFunc2:        vmTypeGetCmd,
		PermissionsRequired: []string{command.VM_TYPES_READ},
	},
	{
		Description:  "Create VM type.",
		Subject:      "vm-type",
		AltSubject:   "vm-type",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("create VM type", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"format":                c.FlagSet.String("format", "json", "The input format. Supported values are 'json','yaml'. The default format is json."),
				"read_config_from_file": c.FlagSet.String("raw-config", command.NilDefaultStr, colors.Red("(Required)")+" Read raw object from file"),
				"read_config_from_pipe": c.FlagSet.Bool("pipe", false, colors.Green("(Flag)")+" If set, read raw object from pipe instead of from a file. Either this flag or the --raw-config option must be used."),
				"return_id":             c.FlagSet.Bool("return-id", false, "Will print the ID of the created object. Useful for automating tasks."),
			}
		},
		ExecuteFunc2:        vmTypeCreateCmd,
		PermissionsRequired: []string{command.VM_TYPES_WRITE},
	},
	{
		Description:  "Update VM type.",
		Subject:      "vm-type",
		AltSubject:   "vm-type",
		Predicate:    "edit",
		AltPredicate: "update",
		FlagSet:      flag.NewFlagSet("update VM type", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"vm_type_id":            c.FlagSet.Int("id", 0, "VM type id"),
				"format":                c.FlagSet.String("format", "json", "The input format. Supported values are 'json','yaml'. The default format is json."),
				"read_config_from_file": c.FlagSet.String("raw-config", command.NilDefaultStr, colors.Red("(Required)")+" Read raw object from file"),
				"read_config_from_pipe": c.FlagSet.Bool("pipe", false, colors.Green("(Flag)")+" If set, read raw object from pipe instead of from a file. Either this flag or the --raw-config option must be used."),
			}
		},
		ExecuteFunc2:        vmTypeUpdateCmd,
		PermissionsRequired: []string{command.VM_TYPES_WRITE},
	},
	{
		Description:  "Delete VM type.",
		Subject:      "vm-type",
		AltSubject:   "vm-type",
		Predicate:    "delete",
		AltPredicate: "rm",
		FlagSet:      flag.NewFlagSet("delete VM type", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"vm_type_id":  c.FlagSet.Int("id", 0, "VM type id"),
				"autoconfirm": c.FlagSet.Bool("autoconfirm", false, colors.Green("(Flag)")+" If set it will assume action is confirmed"),
			}
		},
		ExecuteFunc2:        vmTypeDeleteCmd,
		PermissionsRequired: []string{command.VM_TYPES_WRITE},
	},
}

func vmTypeListCmd(ctx context.Context, c *command.Command, client *metalcloud2.APIClient) (string, error) {
	vmTypes, response, err := client.VMTypesApi.GetVMTypes(ctx)
	if err != nil {
		return "", err
	}
	if response.StatusCode >= 400 {
		return "", fmt.Errorf("HTTP error: %s", response.Status)
	}

	format := command.GetStringParam(c.Arguments["format"])

	rawData := []metalcloud2.VmType{}
	formattedData := [][]interface{}{}

	for _, vmType := range vmTypes.Data {
		vmType.Links = nil

		rawData = append(rawData, vmType)

		formattedData = append(formattedData, formattedVMTypeRecord(vmType))
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
			Schema: vmTypeFieldSchema(),
		}

		return table.RenderTable("VM Types", "", format)
	}
}

func vmTypeGetCmd(ctx context.Context, c *command.Command, client *metalcloud2.APIClient) (string, error) {
	vmTypeId, ok := command.GetIntParamOk(c.Arguments["vm_type_id"])
	if !ok {
		return "", fmt.Errorf("-id is required")
	}

	vmType, response, err := client.VMTypesApi.GetVMType(ctx, float64(vmTypeId))
	if err != nil {
		return "", err
	}
	if response.StatusCode >= 400 {
		return "", fmt.Errorf("HTTP error: %s", response.Status)
	}

	format := command.GetStringParam(c.Arguments["format"])

	vmType.Links = nil

	switch format {
	case "json", "JSON":
		result, err := json.MarshalIndent(vmType, "", "\t")
		if err != nil {
			return "", err
		}

		return string(result), nil

	case "yaml", "YAML":
		result, err := yaml.Marshal(vmType)
		if err != nil {
			return "", err
		}

		return string(result), nil

	default:
		formattedData := [][]interface{}{}
		formattedData = append(formattedData, formattedVMTypeRecord(vmType))

		table := tableformatter.Table{
			Data:   formattedData,
			Schema: vmTypeFieldSchema(),
		}

		return table.RenderTable("VM Type", "", format)
	}
}

func vmTypeCreateCmd(ctx context.Context, c *command.Command, client *metalcloud2.APIClient) (string, error) {
	var obj metalcloud2.CreateVmType

	err := command.GetRawObjectFromCommand(c, &obj)
	if err != nil {
		return "", err
	}

	vmType, _, err := client.VMTypesApi.CreateVMType(ctx, obj)
	if err != nil {
		return "", err
	}

	if command.GetBoolParam(c.Arguments["return_id"]) {
		return fmt.Sprintf("%v", vmType.Id), nil
	}

	return "", err
}

func vmTypeUpdateCmd(ctx context.Context, c *command.Command, client *metalcloud2.APIClient) (string, error) {
	vmTypeId, ok := command.GetIntParamOk(c.Arguments["vm_type_id"])
	if !ok {
		return "", fmt.Errorf("-id is required")
	}

	var obj metalcloud2.UpdateVmType

	err := command.GetRawObjectFromCommand(c, &obj)
	if err != nil {
		return "", err
	}

	_, _, err = client.VMTypesApi.UpdateVMType(ctx, obj, float64(vmTypeId))
	if err != nil {
		return "", err
	}

	return "", err
}

func vmTypeDeleteCmd(ctx context.Context, c *command.Command, client *metalcloud2.APIClient) (string, error) {
	vmTypeId, ok := command.GetIntParamOk(c.Arguments["vm_type_id"])
	if !ok {
		return "", fmt.Errorf("-id is required")
	}

	confirmed, err := utils.GetConfirmation(command.GetBoolParam(c.Arguments["autoconfirm"]), fmt.Sprintf("Deleting VM type #%d.", vmTypeId))
	if err != nil {
		return "", err
	}

	if !confirmed {
		return "", fmt.Errorf("Operation not confirmed. Aborting")
	}

	_, err = client.VMTypesApi.DeleteVMType(ctx, float64(vmTypeId))
	if err != nil {
		return "", err
	}

	return "", err
}

func vmTypeFieldSchema() []tableformatter.SchemaField {
	schema := []tableformatter.SchemaField{
		{
			FieldName: "ID",
			FieldType: tableformatter.TypeFloat,
			FieldSize: 6,
		},

		{
			FieldName: "NAME",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},

		{
			FieldName: "LABEL",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},

		{
			FieldName: "CPU",
			FieldType: tableformatter.TypeInt,
			FieldSize: 5,
		},

		{
			FieldName: "RAM",
			FieldType: tableformatter.TypeFloat,
			FieldSize: 5,
		},
	}

	return schema
}

func formattedVMTypeRecord(vmType metalcloud2.VmType) []interface{} {
	formattedRecord := []interface{}{
		vmType.Id,
		vmType.Name,
		vmType.Label,
		int(vmType.CpuCores),
		vmType.RamGB,
	}

	return formattedRecord
}
