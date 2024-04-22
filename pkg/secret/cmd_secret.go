package secret

import (
	"encoding/base64"
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

var SecretsCmds = []command.Command{
	{
		Description:  "Create a secret.",
		Subject:      "secret",
		AltSubject:   "sec",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("create secret", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"read_config_from_file": c.FlagSet.String("f", command.NilDefaultStr, colors.Red("(Required)")+" Read configuration from file in the format specified with --format."),
				"format":                c.FlagSet.String("format", "yaml", "The input format. Supported values are 'json','yaml'. The default format is json."),
				"return_id":             c.FlagSet.Bool("return-id", false, colors.Green("(Flag)")+" If set will print the ID of the created infrastructure. Useful for automating tasks."),
			}
		},
		ExecuteFunc: secretCreateCmd,
		Endpoint:    configuration.ExtendedEndpoint,
	},
	{
		Description:  "Lists available secrets.",
		Subject:      "secret",
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
		ExecuteFunc: secretListCmd,
		Endpoint:    configuration.ExtendedEndpoint,
	},
	{
		Description:  "Get a secret.",
		Subject:      "secret",
		AltSubject:   "sec",
		Predicate:    "get",
		AltPredicate: "show",
		FlagSet:      flag.NewFlagSet("get secret", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"secret_id": c.FlagSet.Int("id", command.NilDefaultInt, colors.Red("(Required)")+" ID of the secret"),
				"format":    c.FlagSet.String("format", "yaml", "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
				"usage":     c.FlagSet.String("usage", command.NilDefaultStr, "Secret's usage"),
			}
		},
		ExecuteFunc: secretGetCmd,
		Endpoint:    configuration.ExtendedEndpoint,
	},
	{
		Description:  "Update a secret.",
		Subject:      "secret",
		AltSubject:   "sec",
		Predicate:    "update",
		AltPredicate: "edit",
		FlagSet:      flag.NewFlagSet("update secret", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"read_config_from_file": c.FlagSet.String("f", command.NilDefaultStr, colors.Red("(Required)")+" Read configuration from file in the format specified with --format."),
				"format":                c.FlagSet.String("format", "yaml", "The input format. Supported values are 'json','yaml'. The default format is json."),
			}
		},
		ExecuteFunc: secretUpdateCmd,
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

func secretCreateCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	obj, err := objects.ReadSingleObjectFromCommand(c, client)
	if err != nil {
		return "", err
	}
	secret := (*obj).(metalcloud.Secret)

	secret.SecretBase64 = base64.StdEncoding.EncodeToString([]byte(secret.SecretBase64))
	createdSecret, err := client.SecretCreate(secret)
	if err != nil {
		return "", err
	}

	if command.GetBoolParam(c.Arguments["return_id"]) {
		return fmt.Sprintf("%d", createdSecret.SecretID), nil
	}

	return "", err
}

func secretListCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
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

func secretGetCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	secretID, ok := command.GetIntParamOk(c.Arguments["secret_id"])
	if !ok {
		return "", fmt.Errorf("-id required")
	}

	secret, err := client.SecretGet(secretID)
	if err != nil {
		return "", err
	}

	format := command.GetStringParam(c.Arguments["format"])
	ret, err := objects.RenderRawObject(*secret, format, "Secret")
	if err != nil {
		return "", err
	}

	return ret, nil
}

func secretUpdateCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	obj, err := objects.ReadSingleObjectFromCommand(c, client)
	if err != nil {
		return "", err
	}
	secret := (*obj).(metalcloud.Secret)

	secret.SecretBase64 = base64.StdEncoding.EncodeToString([]byte(secret.SecretBase64))
	_, err = client.SecretUpdate(secret.SecretID, secret)
	if err != nil {
		return "", err
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
