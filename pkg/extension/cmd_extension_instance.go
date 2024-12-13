package extension

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

var ExtensionInstanceCmds = []command.Command{
	{
		Description:  "Lists all extension-instances.",
		Subject:      "extension-instance",
		AltSubject:   "extension-instances",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list extension-instances", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"format": c.FlagSet.String("format", command.NilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
			}
		},
		ExecuteFunc2:        extensionInstancesListCmd,
		PermissionsRequired: []string{command.EXTENSIONS_READ},
		MinApiVersion:       "v6.4",
	},
	{
		Description:  "Get extension-instance.",
		Subject:      "extension-instance",
		AltSubject:   "extension-instance",
		Predicate:    "get",
		AltPredicate: "show",
		FlagSet:      flag.NewFlagSet("get extension-instance", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"extension_instance_id": c.FlagSet.Int("id", 0, "extension-instance id"),
				"format":                c.FlagSet.String("format", command.NilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
			}
		},
		ExecuteFunc2:        extensionInstanceGetCmd,
		PermissionsRequired: []string{command.EXTENSIONS_READ},
		MinApiVersion:       "v6.4",
	},
	{
		Description:  "Create extension-instance.",
		Subject:      "extension-instance",
		AltSubject:   "extension-instance",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("create extension-instance", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_id":     c.FlagSet.Int("id", 0, "infrastructure id"),
				"format":                c.FlagSet.String("format", "json", "The input format. Supported values are 'json','yaml'. The default format is json."),
				"read_config_from_file": c.FlagSet.String("raw-config", command.NilDefaultStr, colors.Red("(Required)")+" Read raw object from file"),
				"read_config_from_pipe": c.FlagSet.Bool("pipe", false, colors.Green("(Flag)")+" If set, read raw object from pipe instead of from a file. Either this flag or the --raw-config option must be used."),
				"return_id":             c.FlagSet.Bool("return-id", false, "Will print the ID of the created object. Useful for automating tasks."),
			}
		},
		ExecuteFunc2:        extensionInstanceCreateCmd,
		PermissionsRequired: []string{command.EXTENSIONS_WRITE},
		MinApiVersion:       "v6.4",
	},
	{
		Description:  "Update extension-instance.",
		Subject:      "extension-instance",
		AltSubject:   "extension-instance",
		Predicate:    "edit",
		AltPredicate: "update",
		FlagSet:      flag.NewFlagSet("update extension-instance", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"extension_instance_id": c.FlagSet.Int("id", 0, "extension-instance id"),
				"format":                c.FlagSet.String("format", "json", "The input format. Supported values are 'json','yaml'. The default format is json."),
				"read_config_from_file": c.FlagSet.String("raw-config", command.NilDefaultStr, colors.Red("(Required)")+" Read raw object from file"),
				"read_config_from_pipe": c.FlagSet.Bool("pipe", false, colors.Green("(Flag)")+" If set, read raw object from pipe instead of from a file. Either this flag or the --raw-config option must be used."),
			}
		},
		ExecuteFunc2:        extensionInstanceUpdateCmd,
		PermissionsRequired: []string{command.EXTENSIONS_WRITE},
		MinApiVersion:       "v6.4",
	},
	{
		Description:  "Delete extension-instance.",
		Subject:      "extension-instance",
		AltSubject:   "extension-instance",
		Predicate:    "delete",
		AltPredicate: "del",
		FlagSet:      flag.NewFlagSet("delete extension-instance", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"extension_instance_id": c.FlagSet.Int("id", 0, "extension-instance id"),
				"autoconfirm":           c.FlagSet.Bool("autoconfirm", false, colors.Green("(Flag)")+" If set it will assume action is confirmed"),
			}
		},
		ExecuteFunc2:        extensionInstanceDeleteCmd,
		PermissionsRequired: []string{command.EXTENSIONS_WRITE},
		MinApiVersion:       "v6.4",
	},
}

func extensionInstancesListCmd(ctx context.Context, c *command.Command, client *metalcloud2.APIClient) (string, error) {
	extensionInstancesList, response, err := client.ExtensionInstanceApi.GetExtensionInstances(ctx, nil)
	if err != nil {
		return "", err
	}
	if response.StatusCode >= 400 {
		return "", fmt.Errorf("HTTP error: %s", response.Status)
	}

	format := command.GetStringParam(c.Arguments["format"])

	rawData := []metalcloud2.ExtensionInstanceDto{}
	formattedData := [][]interface{}{}

	for _, extensionInstance := range extensionInstancesList.ExtensionInstances {
		rawData = append(rawData, extensionInstance)
		formattedData = append(formattedData, formattedExtensionInstanceRecord(extensionInstance))
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
			Schema: extensionInstanceFieldSchema(),
		}

		return table.RenderTable("Extension instances", "", format)
	}
}

func extensionInstanceGetCmd(ctx context.Context, c *command.Command, client *metalcloud2.APIClient) (string, error) {
	extensionInstanceId, ok := command.GetIntParamOk(c.Arguments["extension_instance_id"])
	if !ok {
		return "", fmt.Errorf("-id is required")
	}

	extensionInstance, response, err := client.ExtensionInstanceApi.GetExtensionInstance(ctx, float64(extensionInstanceId))
	if err != nil {
		return "", err
	}
	if response.StatusCode >= 400 {
		return "", fmt.Errorf("HTTP error: %s", response.Status)
	}

	format := command.GetStringParam(c.Arguments["format"])

	switch format {
	case "json", "JSON":
		result, err := json.MarshalIndent(extensionInstance, "", "\t")
		if err != nil {
			return "", err
		}

		return string(result), nil

	case "yaml", "YAML":
		result, err := yaml.Marshal(extensionInstance)
		if err != nil {
			return "", err
		}

		return string(result), nil

	default:
		formattedData := [][]interface{}{}
		formattedData = append(formattedData, formattedExtensionInstanceRecord(extensionInstance))

		table := tableformatter.Table{
			Data:   formattedData,
			Schema: extensionFieldSchema(),
		}

		return table.RenderTable("Extension instance", "", format)
	}
}

func extensionInstanceCreateCmd(ctx context.Context, c *command.Command, client *metalcloud2.APIClient) (string, error) {
	infrastructureId, ok := command.GetIntParamOk(c.Arguments["infrastructure_id"])
	if !ok {
		return "", fmt.Errorf("-id is required")
	}

	var obj metalcloud2.CreateExtensionInstanceDto

	err := command.GetRawObjectFromCommand(c, &obj)
	if err != nil {
		return "", err
	}

	extensionInstance, _, err := client.ExtensionInstanceApi.CreateExtensionInstance(ctx, obj, float64(infrastructureId))
	if err != nil {
		return "", err
	}

	if command.GetBoolParam(c.Arguments["return_id"]) {
		return fmt.Sprintf("%v", extensionInstance.Id), nil
	}

	return "", err
}

func extensionInstanceUpdateCmd(ctx context.Context, c *command.Command, client *metalcloud2.APIClient) (string, error) {
	extensionInstanceId, ok := command.GetIntParamOk(c.Arguments["extension_instance_id"])
	if !ok {
		return "", fmt.Errorf("-id is required")
	}

	var obj metalcloud2.UpdateExtensionInstanceDto

	err := command.GetRawObjectFromCommand(c, &obj)
	if err != nil {
		return "", err
	}

	_, _, err = client.ExtensionInstanceApi.UpdateExtensionInstance(ctx, obj, float64(extensionInstanceId))
	if err != nil {
		return "", err
	}

	return "", err
}

func extensionInstanceDeleteCmd(ctx context.Context, c *command.Command, client *metalcloud2.APIClient) (string, error) {
	extensionInstanceId, ok := command.GetIntParamOk(c.Arguments["extension_instance_id"])
	if !ok {
		return "", fmt.Errorf("-id is required")
	}

	confirmed, err := utils.GetConfirmation(command.GetBoolParam(c.Arguments["autoconfirm"]), fmt.Sprintf("Deleting Extension Instance #%d.", extensionInstanceId))
	if err != nil {
		return "", err
	}

	if !confirmed {
		return "", fmt.Errorf("operation not confirmed, aborting")
	}

	_, err = client.ExtensionInstanceApi.DeleteExtensionInstance(ctx, float64(extensionInstanceId))
	if err != nil {
		return "", err
	}

	return "", err
}

func extensionInstanceFieldSchema() []tableformatter.SchemaField {
	schema := []tableformatter.SchemaField{
		{
			FieldName: "ID",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},

		{
			FieldName: "INFRASTRUCTURE ID",
			FieldType: tableformatter.TypeFloat,
			FieldSize: 6,
		},

		{
			FieldName: "LABEL",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
	}

	return schema
}

func formattedExtensionInstanceRecord(extensionInstance metalcloud2.ExtensionInstanceDto) []interface{} {
	formattedRecord := []interface{}{
		fmt.Sprintf("%.0f", extensionInstance.Id),
		extensionInstance.InfrastructureId,
		extensionInstance.Label,
	}

	return formattedRecord
}
