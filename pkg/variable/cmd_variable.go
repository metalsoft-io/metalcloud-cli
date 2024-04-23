package variable

import (
	"flag"
	"fmt"
	"os"
	"strings"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
	"github.com/metalsoft-io/metalcloud-cli/internal/colors"
	"github.com/metalsoft-io/metalcloud-cli/internal/command"
	"github.com/metalsoft-io/metalcloud-cli/internal/configuration"
	"github.com/metalsoft-io/metalcloud-cli/internal/objects"
	"github.com/metalsoft-io/tableformatter"
)

var VariablesCmds = []command.Command{
	{
		Description:  "Create a variable.",
		Subject:      "variable",
		AltSubject:   "var",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("create variable", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"read_config_from_file": c.FlagSet.String("f", command.NilDefaultStr, colors.Red("(Required)")+" Read configuration from file in the format specified with --format."),
				"format":                c.FlagSet.String("format", "yaml", "The input format. Supported values are 'json','yaml'. The default format is json."),
				"return_id":             c.FlagSet.Bool("return-id", false, colors.Green("(Flag)")+" If set will print the ID of the created variable. Useful for automating tasks."),
			}
		},
		ExecuteFunc: variableCreateCmd,
		Endpoint:    configuration.ExtendedEndpoint,
	},
	{
		Description:  "Lists all variables.",
		Subject:      "variable",
		AltSubject:   "var",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list variables", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"format": c.FlagSet.String("format", command.NilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
				"usage":  c.FlagSet.String("usage", command.NilDefaultStr, "Variable's usage"),
			}
		},
		ExecuteFunc: variableListCmd,
		Endpoint:    configuration.ExtendedEndpoint,
	},
	{
		Description:  "Get a variable.",
		Subject:      "variable",
		AltSubject:   "var",
		Predicate:    "get",
		AltPredicate: "show",
		FlagSet:      flag.NewFlagSet("get variable", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"variable_id": c.FlagSet.Int("id", command.NilDefaultInt, colors.Red("(Required)")+" Variable ID."),
				"format":      c.FlagSet.String("format", "yaml", "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
			}
		},
		ExecuteFunc: variableGetCmd,
		Endpoint:    configuration.ExtendedEndpoint,
	},
	{
		Description:  "Update a variable.",
		Subject:      "variable",
		AltSubject:   "var",
		Predicate:    "update",
		AltPredicate: "edit",
		FlagSet:      flag.NewFlagSet("update variable", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"read_config_from_file": c.FlagSet.String("f", command.NilDefaultStr, colors.Red("(Required)")+" Read configuration from file in the format specified with --format."),
				"format":                c.FlagSet.String("format", "yaml", "The input format. Supported values are 'json','yaml'. The default format is json."),
			}
		},
		ExecuteFunc: variableUpdateCmd,
		Endpoint:    configuration.ExtendedEndpoint,
	},
	{
		Description:  "Delete a variable.",
		Subject:      "variable",
		AltSubject:   "var",
		Predicate:    "delete",
		AltPredicate: "rm",
		FlagSet:      flag.NewFlagSet("delete variable", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"variable_id_or_name": c.FlagSet.String("id", command.NilDefaultStr, "Variable's id or name"),
				"autoconfirm":         c.FlagSet.Bool("autoconfirm", false, colors.Green("(Flag)")+" If set it will assume action is confirmed"),
			}
		},
		ExecuteFunc: variableDeleteCmd,
		Endpoint:    configuration.ExtendedEndpoint,
	},
}

func variableCreateCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	obj, err := objects.ReadSingleObjectFromCommand(c, client)
	if err != nil {
		return "", err
	}
	variable := (*obj).(metalcloud.Variable)

	createdVariable, err := client.VariableCreate(variable)
	if err != nil {
		return "", err
	}

	if command.GetBoolParam(c.Arguments["return_id"]) {
		return fmt.Sprintf("%d", createdVariable.VariableID), nil
	}

	return "", err
}

func variableListCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	usage := *c.Arguments["usage"].(*string)
	if usage == command.NilDefaultStr {
		usage = ""
	}

	list, err := client.Variables(usage)

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
			FieldName: "NAME",
			FieldType: tableformatter.TypeString,
			FieldSize: 20,
		},
		{
			FieldName: "USAGE",
			FieldType: tableformatter.TypeString,
			FieldSize: 20,
		},
		{
			FieldName: "CREATED",
			FieldType: tableformatter.TypeString,
			FieldSize: 20,
		},
		{
			FieldName: "UPDATED",
			FieldType: tableformatter.TypeString,
			FieldSize: 20,
		},
	}

	data := [][]interface{}{}
	for _, s := range *list {

		data = append(data, []interface{}{
			s.VariableID,
			s.VariableName,
			s.VariableUsage,
			s.VariableCreatedTimestamp,
			s.VariableUpdatedTimestamp,
		})

	}

	tableformatter.TableSorter(schema).OrderBy(schema[0].FieldName).Sort(data)

	table := tableformatter.Table{
		Data:   data,
		Schema: schema,
	}
	return table.RenderTable("Variables", "", command.GetStringParam(c.Arguments["format"]))
}

func variableGetCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	variableID, ok := command.GetIntParamOk(c.Arguments["variable_id"])
	if !ok {
		return "", fmt.Errorf("-id required")
	}

	variable, err := client.VariableGet(variableID)
	if err != nil {
		return "", err
	}

	format := command.GetStringParam(c.Arguments["format"])
	ret, err := objects.RenderRawObject(*variable, format, "Variable")
	if err != nil {
		return "", err
	}

	return ret, nil
}

func variableUpdateCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	obj, err := objects.ReadSingleObjectFromCommand(c, client)
	if err != nil {
		return "", err
	}
	variable := (*obj).(metalcloud.Variable)

	_, err = client.VariableUpdate(variable.VariableID, variable)
	if err != nil {
		return "", err
	}

	return "", err
}

func variableDeleteCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	retS, err := getVariableFromCommand("id", c, client)
	if err != nil {
		return "", err
	}
	confirm := false

	if command.GetBoolParam(c.Arguments["autoconfirm"]) {
		confirm = true
	} else {

		confirmationMessage := fmt.Sprintf("Deleting variable %s (%d).  Are you sure? Type \"yes\" to continue:",
			retS.VariableName,
			retS.VariableID)

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

	err = client.VariableDelete(retS.VariableID)

	return "", err
}

func getVariableFromCommand(paramName string, c *command.Command, client metalcloud.MetalCloudClient) (*metalcloud.Variable, error) {

	v, err := command.GetParam(c, "variable_id_or_name", paramName)
	if err != nil {
		return nil, err
	}

	id, label, isID := command.IdOrLabel(v)

	if isID {
		return client.VariableGet(id)
	}

	variables, err := client.Variables("")
	if err != nil {
		return nil, err
	}

	for _, s := range *variables {
		if s.VariableName == label {
			return &s, nil
		}
	}

	return nil, fmt.Errorf("Could not locate variable with id/name %v", *v.(*interface{}))
}
