package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	metalcloud "github.com/bigstepinc/metal-cloud-sdk-go"
	interfaces "github.com/bigstepinc/metalcloud-cli/interfaces"
)

var secretsCmds = []Command{

	Command{
		Description:  "Lists available secrets",
		Subject:      "secrets",
		AltSubject:   "sec",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list secrets", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"format": c.FlagSet.String("format", _nilDefaultStr, "The output format. Supported values are 'json','csv'. The default format is human readable."),
				"usage":  c.FlagSet.String("usage", _nilDefaultStr, "Secret's usage"),
			}
		},
		ExecuteFunc: secretsListCmd,
	},
	Command{
		Description:  "Create secret",
		Subject:      "secret",
		AltSubject:   "sec",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("create secret", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"name":                   c.FlagSet.String("name", _nilDefaultStr, "Secret's name"),
				"usage":                  c.FlagSet.String("usage", _nilDefaultStr, "Secret's usage"),
				"read_content_from_pipe": c.FlagSet.Bool("pipe", false, "Read secret's content read from pipe instead of terminal input"),
			}
		},
		ExecuteFunc: secretCreateCmd,
	},
	Command{
		Description:  "Delete secret",
		Subject:      "secret",
		AltSubject:   "sec",
		Predicate:    "delete",
		AltPredicate: "rm",
		FlagSet:      flag.NewFlagSet("delete secret", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"secret_id_or_name": c.FlagSet.String("id", _nilDefaultStr, "Secret's id or name"),
				"autoconfirm":       c.FlagSet.Bool("autoconfirm", false, "If true it does not ask for confirmation anymore"),
			}
		},
		ExecuteFunc: secretDeleteCmd,
	},
}

func secretsListCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	usage := *c.Arguments["usage"].(*string)
	if usage == _nilDefaultStr {
		usage = ""
	}

	list, err := client.Secrets(usage)

	if err != nil {
		return "", err
	}

	schema := []SchemaField{
		SchemaField{
			FieldName: "ID",
			FieldType: TypeInt,
			FieldSize: 6,
		},
		SchemaField{
			FieldName: "NAME",
			FieldType: TypeString,
			FieldSize: 20,
		},
		SchemaField{
			FieldName: "USAGE",
			FieldType: TypeString,
			FieldSize: 20,
		},
		SchemaField{
			FieldName: "CREATED",
			FieldType: TypeString,
			FieldSize: 20,
		},
		SchemaField{
			FieldName: "UPDATED",
			FieldType: TypeString,
			FieldSize: 20,
		},
	}

	user := GetUserEmail()

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

	var sb strings.Builder

	format := c.Arguments["format"]
	if format == nil {
		var f string
		f = ""
		format = &f
	}

	switch *format.(*string) {
	case "json", "JSON":
		ret, err := GetTableAsJSONString(data, schema)
		if err != nil {
			return "", err
		}
		sb.WriteString(ret)
	case "csv", "CSV":
		ret, err := GetTableAsCSVString(data, schema)
		if err != nil {
			return "", err
		}
		sb.WriteString(ret)

	default:
		sb.WriteString(fmt.Sprintf("Secrets I have access to (as %s)\n", user))

		TableSorter(schema).OrderBy(
			schema[0].FieldName,
			schema[1].FieldName).Sort(data)

		AdjustFieldSizes(data, &schema)

		sb.WriteString(GetTableAsString(data, schema))

		sb.WriteString(fmt.Sprintf("Total: %d secrets\n\n", len(*list)))
	}

	return sb.String(), nil
}

func secretCreateCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {
	secret := metalcloud.Secret{}

	if v := c.Arguments["name"]; v != nil && *v.(*string) != _nilDefaultStr {
		secret.SecretName = *v.(*string)
	}

	if v := c.Arguments["usage"]; v != nil && *v.(*string) != _nilDefaultStr {
		secret.SecretUsage = *v.(*string)
	}

	content := []byte{}
	if v := c.Arguments["read_content_from_pipe"]; *v.(*bool) {
		content = readInputFromPipe()
	} else {
		if runtime.GOOS == "windows" {
			content = requestInput("Secret content:")
		} else {
			content = requestInputSilent("Secret content:")
		}

	}

	secret.SecretBase64 = base64.StdEncoding.EncodeToString([]byte(content))

	_, err := client.SecretCreate(secret)

	return "", err
}

func secretDeleteCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	retS, err := getSecretFromCommand("id", c, client)
	if err != nil {
		return "", err
	}
	confirm := false

	if c.Arguments["autoconfirm"] != nil && *c.Arguments["autoconfirm"].(*bool) == true {
		confirm = true
	} else {

		confirmationMessage := fmt.Sprintf("Deleting secret %s (%d).  Are you sure? Type \"yes\" to continue:",
			retS.SecretName,
			retS.SecretID)

		//this is simply so that we don't output a text on the command line under go test
		if strings.HasSuffix(os.Args[0], ".test") {
			confirmationMessage = ""
		}

		confirm = requestConfirmation(confirmationMessage)
	}

	if !confirm {
		return "", fmt.Errorf("Operation not confirmed. Aborting")
	}

	err = client.SecretDelete(retS.SecretID)

	return "", err
}

func getSecretFromCommand(paramName string, c *Command, client interfaces.MetalCloudClient) (*metalcloud.Secret, error) {

	v, err := getParam(c, "secret_id_or_name", paramName)
	if err != nil {
		return nil, err
	}

	id, label, isID := idOrLabel(v)

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
