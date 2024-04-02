package variable

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
	"github.com/metalsoft-io/metalcloud-cli/internal/colors"
	"github.com/metalsoft-io/metalcloud-cli/internal/command"
	"github.com/metalsoft-io/metalcloud-cli/internal/configuration"
	"github.com/metalsoft-io/tableformatter"
)

var VariablesCmds = []command.Command{
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
		ExecuteFunc: variablesListCmd,
		Endpoint:    configuration.ExtendedEndpoint,
	},
	{
		Description:  "Create a variable.",
		Subject:      "variable",
		AltSubject:   "var",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("create variable", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"name":                   c.FlagSet.String("name", command.NilDefaultStr, "Variable's name"),
				"usage":                  c.FlagSet.String("usage", command.NilDefaultStr, "Variable's usage"),
				"read_content_from_pipe": c.FlagSet.Bool("pipe", false, "Read variable's content read from pipe instead of terminal input"),
				"return_id":              c.FlagSet.Bool("return-id", false, colors.Green("(Flag)")+" If set will print the ID of the created infrastructure. Useful for automating tasks."),
			}
		},
		ExecuteFunc: variableCreateCmd,
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

func variablesListCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

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

func variableCreateCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	variable := metalcloud.Variable{}

	if v, ok := command.GetStringParamOk(c.Arguments["name"]); ok {
		variable.VariableName = v
	} else {
		return "", fmt.Errorf("name is required")
	}

	variable.VariableUsage = command.GetStringParam(c.Arguments["usage"])

	var err error
	content := []byte{}

	if command.GetBoolParam(c.Arguments["read_content_from_pipe"]) {
		content, err = configuration.ReadInputFromPipe()
	} else {
		content, err = command.RequestInput("Variable content:")
	}

	if err != nil {
		return "", err
	}

	if len(content) == 0 {
		return "", fmt.Errorf("Content cannot be empty")
	}

	cleanedContent := strings.Trim(string(content), "\"\r\n")

	b, err := json.Marshal(cleanedContent)

	variable.VariableJSON = string(b)

	ret, err := client.VariableCreate(variable)
	if err != nil {
		return "", err
	}

	if c.Arguments["return_id"] != nil && *c.Arguments["return_id"].(*bool) {
		return fmt.Sprintf("%d", ret.VariableID), nil
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
