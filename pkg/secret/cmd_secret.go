package secret

import (
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
	"github.com/metalsoft-io/metalcloud-cli/internal/colors"
	"github.com/metalsoft-io/metalcloud-cli/internal/command"
	"github.com/metalsoft-io/metalcloud-cli/internal/configuration"
	"github.com/metalsoft-io/tableformatter"
)

var SecretsCmds = []command.Command{
	{
		Description:  "Lists available secrets.",
		Subject:      "secrets",
		AltSubject:   "sec",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list secrets", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"format": c.FlagSet.String("format", command.NilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
				"usage":  c.FlagSet.String("usage", command.NilDefaultStr, "Secret's usage"),
			}
		},
		ExecuteFunc: secretsListCmd,
		Endpoint:    configuration.ExtendedEndpoint,
	},
	{
		Description:  "Create a secret.",
		Subject:      "secret",
		AltSubject:   "sec",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("create secret", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"name":                   c.FlagSet.String("name", command.NilDefaultStr, colors.Red("(Required)")+" Secret's name"),
				"usage":                  c.FlagSet.String("usage", command.NilDefaultStr, "Secret's usage"),
				"read_content_from_pipe": c.FlagSet.Bool("pipe", false, "Read secret's content read from pipe instead of terminal input"),
				"return_id":              c.FlagSet.Bool("return-id", false, colors.Green("(Flag)")+" If set will print the ID of the created infrastructure. Useful for automating tasks."),
			}
		},
		ExecuteFunc: secretCreateCmd,
		Endpoint:    configuration.ExtendedEndpoint,
	},
	{
		Description:  "Delete a secret.",
		Subject:      "secret",
		AltSubject:   "sec",
		Predicate:    "delete",
		AltPredicate: "rm",
		FlagSet:      flag.NewFlagSet("delete secret", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"secret_id_or_name": c.FlagSet.String("id", command.NilDefaultStr, "Secret's id or name"),
				"autoconfirm":       c.FlagSet.Bool("autoconfirm", false, colors.Green("(Flag)")+" If set it will assume action is confirmed"),
			}
		},
		ExecuteFunc: secretDeleteCmd,
		Endpoint:    configuration.ExtendedEndpoint,
	},
}

func secretsListCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	usage := *c.Arguments["usage"].(*string)
	if usage == command.NilDefaultStr {
		usage = ""
	}

	list, err := client.Secrets(usage)

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
			s.SecretID,
			s.SecretName,
			s.SecretUsage,
			s.SecretCreatedTimestamp,
			s.SecretUpdatedTimestamp,
		})

	}

	tableformatter.TableSorter(schema).OrderBy(schema[0].FieldName).Sort(data)

	table := tableformatter.Table{
		Data:   data,
		Schema: schema,
	}
	return table.RenderTable("Secrets", "", command.GetStringParam(c.Arguments["format"]))
}

func secretCreateCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	secret := metalcloud.Secret{}

	secretName, ok := command.GetStringParamOk(c.Arguments["name"])
	if !ok {
		return "", fmt.Errorf("name is required")
	} else {
		secret.SecretName = secretName
	}

	if v, ok := command.GetStringParamOk(c.Arguments["usage"]); ok {
		secret.SecretUsage = v
	}

	content := []byte{}
	var err error
	if v := c.Arguments["read_content_from_pipe"]; *v.(*bool) {
		content, err = configuration.ReadInputFromPipe()
	} else {
		if runtime.GOOS == "windows" {
			content, err = command.RequestInput("Secret content:")
		} else {
			content, err = command.RequestInputSilent("Secret content:")
		}
	}

	if err != nil {
		return "", err
	}

	if len(content) == 0 {
		return "", fmt.Errorf("Content cannot be empty")
	}

	secret.SecretBase64 = base64.StdEncoding.EncodeToString([]byte(content))

	ret, err := client.SecretCreate(secret)
	if err != nil {
		return "", err
	}

	if c.Arguments["return_id"] != nil && *c.Arguments["return_id"].(*bool) {
		return fmt.Sprintf("%d", ret.SecretID), nil
	}

	return "", err
}

func secretDeleteCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	retS, err := getSecretFromCommand("id", c, client)
	if err != nil {
		return "", err
	}
	confirm := false

	if command.GetBoolParam(c.Arguments["autoconfirm"]) {
		confirm = true
	} else {

		confirmationMessage := fmt.Sprintf("Deleting secret %s (%d).  Are you sure? Type \"yes\" to continue:",
			retS.SecretName,
			retS.SecretID)

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

	err = client.SecretDelete(retS.SecretID)

	return "", err
}

func getSecretFromCommand(paramName string, c *command.Command, client metalcloud.MetalCloudClient) (*metalcloud.Secret, error) {

	v, err := command.GetParam(c, "secret_id_or_name", paramName)
	if err != nil {
		return nil, err
	}

	id, label, isID := command.IdOrLabel(v)

	if isID {
		return client.SecretGet(id)
	}

	secrets, err := client.Secrets("")
	if err != nil {
		return nil, err
	}

	for _, s := range *secrets {
		if s.SecretName == label {
			return &s, nil
		}
	}

	return nil, fmt.Errorf("Could not locate secret with id/name %v", *v.(*interface{}))
}
