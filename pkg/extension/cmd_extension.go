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

var ExtensionCmds = []command.Command{
	{
		Description:  "Lists all extensions.",
		Subject:      "extension",
		AltSubject:   "extensions",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list extensions", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"format": c.FlagSet.String("format", command.NilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
			}
		},
		ExecuteFunc2:        extensionListCmd,
		PermissionsRequired: []string{command.EXTENSIONS_READ},
	},
	{
		Description:  "Get extension.",
		Subject:      "extension",
		AltSubject:   "extension",
		Predicate:    "get",
		AltPredicate: "show",
		FlagSet:      flag.NewFlagSet("get extension", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"extension_id": c.FlagSet.Int("id", 0, "extension id"),
				"format":       c.FlagSet.String("format", command.NilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
			}
		},
		ExecuteFunc2:        extensionGetCmd,
		PermissionsRequired: []string{command.EXTENSIONS_READ},
	},
	{
		Description:  "Create extension.",
		Subject:      "extension",
		AltSubject:   "extension",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("create extension", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"format":                c.FlagSet.String("format", "json", "The input format. Supported values are 'json','yaml'. The default format is json."),
				"read_config_from_file": c.FlagSet.String("raw-config", command.NilDefaultStr, colors.Red("(Required)")+" Read raw object from file"),
				"read_config_from_pipe": c.FlagSet.Bool("pipe", false, colors.Green("(Flag)")+" If set, read raw object from pipe instead of from a file. Either this flag or the --raw-config option must be used."),
				"return_id":             c.FlagSet.Bool("return-id", false, "Will print the ID of the created object. Useful for automating tasks."),
			}
		},
		ExecuteFunc2:        extensionCreateCmd,
		PermissionsRequired: []string{command.EXTENSIONS_WRITE},
	},
	{
		Description:  "Update extension.",
		Subject:      "extension",
		AltSubject:   "extension",
		Predicate:    "edit",
		AltPredicate: "update",
		FlagSet:      flag.NewFlagSet("update extension", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"extension_id":          c.FlagSet.Int("id", 0, "extension id"),
				"format":                c.FlagSet.String("format", "json", "The input format. Supported values are 'json','yaml'. The default format is json."),
				"read_config_from_file": c.FlagSet.String("raw-config", command.NilDefaultStr, colors.Red("(Required)")+" Read raw object from file"),
				"read_config_from_pipe": c.FlagSet.Bool("pipe", false, colors.Green("(Flag)")+" If set, read raw object from pipe instead of from a file. Either this flag or the --raw-config option must be used."),
			}
		},
		ExecuteFunc2:        extensionUpdateCmd,
		PermissionsRequired: []string{command.EXTENSIONS_WRITE},
	},
	{
		Description:  "Archive extension.",
		Subject:      "extension",
		AltSubject:   "extension",
		Predicate:    "archive",
		AltPredicate: "arc",
		FlagSet:      flag.NewFlagSet("archive extension", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"extension_id": c.FlagSet.Int("id", 0, "extension id"),
				"autoconfirm":  c.FlagSet.Bool("autoconfirm", false, colors.Green("(Flag)")+" If set it will assume action is confirmed"),
			}
		},
		ExecuteFunc2:        extensionArchiveCmd,
		PermissionsRequired: []string{command.EXTENSIONS_WRITE},
	},
	{
		Description:  "Publish extension.",
		Subject:      "extension",
		AltSubject:   "extension",
		Predicate:    "publish",
		AltPredicate: "pub",
		FlagSet:      flag.NewFlagSet("publish extension", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"extension_id": c.FlagSet.Int("id", 0, "extension id"),
				"autoconfirm":  c.FlagSet.Bool("autoconfirm", false, colors.Green("(Flag)")+" If set it will assume action is confirmed"),
			}
		},
		ExecuteFunc2:        extensionPublishCmd,
		PermissionsRequired: []string{command.EXTENSIONS_WRITE},
	},
}

func extensionListCmd(ctx context.Context, c *command.Command, client *metalcloud2.APIClient) (string, error) {
	extensionsList, response, err := client.ExtensionApi.GetExtensions(ctx, nil)
	if err != nil {
		return "", err
	}
	if response.StatusCode >= 400 {
		return "", fmt.Errorf("HTTP error: %s", response.Status)
	}

	format := command.GetStringParam(c.Arguments["format"])

	rawData := []metalcloud2.ExtensionInfoDto{}
	formattedData := [][]interface{}{}

	for _, extension := range extensionsList.Extensions {
		rawData = append(rawData, extension)
		formattedData = append(formattedData, formattedExtensionInfoRecord(extension))
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
			Schema: extensionFieldSchema(),
		}

		return table.RenderTable("Extensions", "", format)
	}
}

func extensionGetCmd(ctx context.Context, c *command.Command, client *metalcloud2.APIClient) (string, error) {
	extensionId, ok := command.GetIntParamOk(c.Arguments["extension_id"])
	if !ok {
		return "", fmt.Errorf("-id is required")
	}

	extension, response, err := client.ExtensionApi.GetExtension(ctx, float64(extensionId))
	if err != nil {
		return "", err
	}
	if response.StatusCode >= 400 {
		return "", fmt.Errorf("HTTP error: %s", response.Status)
	}

	format := command.GetStringParam(c.Arguments["format"])

	switch format {
	case "json", "JSON":
		result, err := json.MarshalIndent(extension, "", "\t")
		if err != nil {
			return "", err
		}

		return string(result), nil

	case "yaml", "YAML":
		result, err := yaml.Marshal(extension)
		if err != nil {
			return "", err
		}

		return string(result), nil

	default:
		formattedData := [][]interface{}{}
		formattedData = append(formattedData, formattedExtensionRecord(extension))

		table := tableformatter.Table{
			Data:   formattedData,
			Schema: extensionFieldSchema(),
		}

		return table.RenderTable("Extension", "", format)
	}
}

func extensionCreateCmd(ctx context.Context, c *command.Command, client *metalcloud2.APIClient) (string, error) {
	var obj metalcloud2.CreateExtensionDto

	err := command.GetRawObjectFromCommand(c, &obj)
	if err != nil {
		return "", err
	}

	extension, _, err := client.ExtensionApi.CreateExtension(ctx, obj)
	if err != nil {
		return "", err
	}

	if command.GetBoolParam(c.Arguments["return_id"]) {
		return fmt.Sprintf("%v", extension.Id), nil
	}

	return "", err
}

func extensionUpdateCmd(ctx context.Context, c *command.Command, client *metalcloud2.APIClient) (string, error) {
	extensionId, ok := command.GetIntParamOk(c.Arguments["extension_id"])
	if !ok {
		return "", fmt.Errorf("-id is required")
	}

	var obj metalcloud2.UpdateExtensionDto

	err := command.GetRawObjectFromCommand(c, &obj)
	if err != nil {
		return "", err
	}

	_, _, err = client.ExtensionApi.UpdateExtension(ctx, obj, float64(extensionId))
	if err != nil {
		return "", err
	}

	return "", err
}

func extensionArchiveCmd(ctx context.Context, c *command.Command, client *metalcloud2.APIClient) (string, error) {
	extensionId, ok := command.GetIntParamOk(c.Arguments["extension_id"])
	if !ok {
		return "", fmt.Errorf("-id is required")
	}

	confirmed, err := utils.GetConfirmation(command.GetBoolParam(c.Arguments["autoconfirm"]), fmt.Sprintf("Archiving Extension #%d.", extensionId))
	if err != nil {
		return "", err
	}

	if !confirmed {
		return "", fmt.Errorf("operation not confirmed, aborting")
	}

	_, err = client.ExtensionApi.ArchiveExtension(ctx, float64(extensionId))
	if err != nil {
		return "", err
	}

	return "", err
}

func extensionPublishCmd(ctx context.Context, c *command.Command, client *metalcloud2.APIClient) (string, error) {
	extensionId, ok := command.GetIntParamOk(c.Arguments["extension_id"])
	if !ok {
		return "", fmt.Errorf("-id is required")
	}

	confirmed, err := utils.GetConfirmation(command.GetBoolParam(c.Arguments["autoconfirm"]), fmt.Sprintf("Publishing Extension #%d.", extensionId))
	if err != nil {
		return "", err
	}

	if !confirmed {
		return "", fmt.Errorf("operation not confirmed, aborting")
	}

	_, err = client.ExtensionApi.PublishExtension(ctx, float64(extensionId))
	if err != nil {
		return "", err
	}

	return "", err
}

func extensionFieldSchema() []tableformatter.SchemaField {
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
			FieldName: "DESCRIPTION",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},

		{
			FieldName: "STATUS",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
	}

	return schema
}

func formattedExtensionInfoRecord(extension metalcloud2.ExtensionInfoDto) []interface{} {
	formattedRecord := []interface{}{
		extension.Id,
		extension.Name,
		extension.Label,
		extension.Description,
		extension.Status,
	}

	return formattedRecord
}

func formattedExtensionRecord(extension metalcloud2.ExtensionDto) []interface{} {
	formattedRecord := []interface{}{
		extension.Id,
		extension.Name,
		extension.Label,
		extension.Description,
		extension.Status,
	}

	return formattedRecord
}
