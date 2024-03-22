package custom_isos

import (
	"flag"
	"fmt"
	"os"
	"strings"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
	"github.com/metalsoft-io/metalcloud-cli/internal/colors"
	"github.com/metalsoft-io/metalcloud-cli/internal/command"
	"github.com/metalsoft-io/tableformatter"
)

var CustomISOCmds = []command.Command{
	{
		Description:  "List custom ISOs.",
		Subject:      "custom-iso",
		AltSubject:   "iso",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("custom-iso", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"user_id": c.FlagSet.Int("user-id", command.NilDefaultInt, "The user ID for which to list the custom iso. Defaults to the current user."),
			}
		},
		ExecuteFunc: customISOListCmd,
	},
	{
		Description:  "Creates a custom iso.",
		Subject:      "custom-iso",
		AltSubject:   "iso",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("custom-iso", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"label":        c.FlagSet.String("label", command.NilDefaultStr, colors.Red("(Required)")+" The custom iso's label"),
				"url":          c.FlagSet.String("url", command.NilDefaultStr, colors.Red("(Required)")+" The custom iso's location (http/https URL)"),
				"display_name": c.FlagSet.String("display-name", command.NilDefaultStr, "The custom iso's display name"),
				"username":     c.FlagSet.String("username", command.NilDefaultStr, "Username to authenticate to the http repository"),
				"password":     c.FlagSet.String("password", command.NilDefaultStr, "Password to authenticate to the http repository"),
				"return_id":    c.FlagSet.Bool("return-id", false, "(Optional) Will print the ID of the created object. Useful for automating tasks."),
			}
		},
		ExecuteFunc: customISOCreateCmd,
	},
	{
		Description:  "Update a custom iso.",
		Subject:      "custom-iso",
		AltSubject:   "iso",
		Predicate:    "update",
		AltPredicate: "edit",
		FlagSet:      flag.NewFlagSet("custom-iso", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"custom_iso_id_or_label": c.FlagSet.String("id", command.NilDefaultStr, colors.Red("(Required)")+" The custom iso's id or label"),
				"url":                    c.FlagSet.String("url", command.NilDefaultStr, colors.Red("(Required)")+" The custom iso's location (http/https URL)"),
				"label":                  c.FlagSet.String("label", command.NilDefaultStr, "The custom iso's label"),
				"display_name":           c.FlagSet.String("display-name", command.NilDefaultStr, "The custom iso's display name"),
				"username":               c.FlagSet.String("username", command.NilDefaultStr, "Username to authenticate to the http repository"),
				"password":               c.FlagSet.String("password", command.NilDefaultStr, "Password to authenticate to the http repository"),
				"return_id":              c.FlagSet.Bool("return-id", false, "(Optional) Will print the ID of the created Object. Useful for automating tasks."),
			}
		},
		ExecuteFunc: customISOUpdateCmd,
	},
	{
		Description:  "Delete a custom iso.",
		Subject:      "custom-iso",
		AltSubject:   "iso",
		Predicate:    "delete",
		AltPredicate: "rm",
		FlagSet:      flag.NewFlagSet("custom-iso", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"custom_iso_id_or_label": c.FlagSet.String("id", command.NilDefaultStr, colors.Red("(Required)")+" The custom iso's id or label"),
				"autoconfirm":            c.FlagSet.Bool("autoconfirm", false, colors.Green("(Flag)")+" If set it will assume action is confirmed"),
			}
		},
		ExecuteFunc: customISODeleteCmd,
	},
	{
		Description:  "Boot a custom iso on a server.",
		Subject:      "custom-iso",
		AltSubject:   "iso",
		Predicate:    "boot-on-server",
		AltPredicate: "boot",
		FlagSet:      flag.NewFlagSet("custom-iso", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"custom_iso_id_or_label": c.FlagSet.String("id", command.NilDefaultStr, colors.Red("(Required)")+" The custom iso's id or label"),
				"server_id":              c.FlagSet.Int("server-id", command.NilDefaultInt, colors.Red("(Required)")+" The server id"),
				"autoconfirm":            c.FlagSet.Bool("autoconfirm", false, colors.Green("(Flag)")+" If set it will assume action is confirmed"),
				"return_id":              c.FlagSet.Bool("return-id", false, "(Optional) Will print the ID of the created Object. Useful for automating tasks."),
			}
		},
		ExecuteFunc: customISOBootIntoServerCmd,
	},
}

type CustomISOCommandConfig struct {
	RequireLabel bool
}

func getCustomISOFromCommand(c *command.Command, config CustomISOCommandConfig) (*metalcloud.CustomISO, error) {
	var label, url, displayName, username, password string
	var ok bool

	if config.RequireLabel {
		label, ok = command.GetStringParamOk(c.Arguments["label"])
		if !ok {
			return nil, fmt.Errorf("-label is required")
		}
	} else {
		label, _ = command.GetStringParamOk(c.Arguments["label"]) // Not required, ignore ok
	}

	url, ok = command.GetStringParamOk(c.Arguments["url"])
	if !ok {
		return nil, fmt.Errorf("-url is required")
	}

	displayName, _ = command.GetStringParamOk(c.Arguments["display_name"]) // Optional, ignore ok
	username, _ = command.GetStringParamOk(c.Arguments["username"])        // Optional, ignore ok
	password, _ = command.GetStringParamOk(c.Arguments["password"])        // Optional, ignore ok

	return &metalcloud.CustomISO{
		CustomISOName:           label,
		CustomISOAccessURL:      url,
		CustomISODisplayName:    displayName,
		CustomISOAccessUsername: username,
		CustomISOAccessPassword: password,
	}, nil
}

func customISOCreateCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	config := CustomISOCommandConfig{RequireLabel: true}
	customISO, err := getCustomISOFromCommand(c, config)
	if err != nil {
		return "", err
	}

	ret, err := client.CustomISOCreate(*customISO)
	if err != nil {
		return "", err
	}

	if command.GetBoolParam(c.Arguments["return_id"]) {
		return fmt.Sprintf("%d", ret.CustomISOID), nil
	}

	return "", err
}

func customISOUpdateCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	custom_iso_id, label, isLabel := command.IdOrLabel(c.Arguments["custom_iso_id_or_label"])
	if isLabel {
		ciList, err := client.CustomISOs(client.GetUserID())
		if err != nil {
			return "", err
		}
		for _, ci := range *ciList {
			if ci.CustomISOName == label {
				custom_iso_id = ci.CustomISOID
			}
		}
	}

	config := CustomISOCommandConfig{RequireLabel: false}
	customISO, err := getCustomISOFromCommand(c, config)
	if err != nil {
		return "", err
	}

	_, err = client.CustomISOUpdate(custom_iso_id, *customISO)
	if err != nil {
		return "", err
	}

	return "", nil
}

func customISODeleteCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	custom_iso_id, label, isLabel := command.IdOrLabel(c.Arguments["custom_iso_id_or_label"])
	if isLabel {
		ciList, err := client.CustomISOs(client.GetUserID())
		if err != nil {
			return "", err
		}
		for _, ci := range *ciList {
			if ci.CustomISOName == label {
				custom_iso_id = ci.CustomISOID
			}
		}
	}

	customISO, err := client.CustomISOGet(custom_iso_id)
	if err != nil {
		return "", err
	}

	confirm := false

	if command.GetBoolParam(c.Arguments["autoconfirm"]) {
		confirm = true
	} else {

		confirmationMessage := fmt.Sprintf("Deleting Custom ISO %s (%d).  Are you sure? Type \"yes\" to continue:", customISO.CustomISOName, customISO.CustomISOID)

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

	err = client.CustomISODelete(customISO.CustomISOID)
	return "", err
}

func customISOBootIntoServerCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	customIsoId, label, isLabel := command.IdOrLabel(c.Arguments["custom_iso_id_or_label"])
	if isLabel {
		ciList, err := client.CustomISOs(client.GetUserID())
		if err != nil {
			return "", err
		}
		for _, ci := range *ciList {
			if ci.CustomISOName == label {
				customIsoId = ci.CustomISOID
			}
		}
	}

	serverId, ok := command.GetIntParamOk(c.Arguments["server_id"])

	if !ok {
		return "", fmt.Errorf("-server-id is required")
	}

	server, err := client.ServerGet(serverId, false)
	if err != nil {
		return "", err
	}

	customISO, err := client.CustomISOGet(customIsoId)
	if err != nil {
		return "", err
	}

	confirm := false

	if command.GetBoolParam(c.Arguments["autoconfirm"]) {
		confirm = true
	} else {

		confirmationMessage := fmt.Sprintf("Booting Custom ISO #%d (%s) on server #%d (%s).  Are you sure? Type \"yes\" to continue:", customISO.CustomISOID, customISO.CustomISOName, server.ServerID, server.ServerSerialNumber)

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

	AFCID, err := client.CustomISOBootIntoServer(customIsoId, serverId)
	if err != nil {
		return "", err
	}

	if command.GetBoolParam(c.Arguments["return_id"]) {
		return fmt.Sprintf("%d", AFCID), nil
	}

	return "", err
}

func customISOListCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	user_id, ok := command.GetIntParamOk(c.Arguments["user-id"])
	if !ok {
		user_id = client.GetUserID()
	}

	customISOList, err := client.CustomISOs(user_id)
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
			FieldName: "LABEL",
			FieldType: tableformatter.TypeString,
			FieldSize: 10,
		},
		{
			FieldName: "Display Name",
			FieldType: tableformatter.TypeString,
			FieldSize: 15,
		},
		{
			FieldName: "URL",
			FieldType: tableformatter.TypeString,
			FieldSize: 10,
		},
	}

	data := [][]interface{}{}
	for _, c := range *customISOList {

		data = append(data, []interface{}{
			c.CustomISOID,
			c.CustomISOName,
			c.CustomISODisplayName,
			c.CustomISOAccessURL,
		})
	}

	tableformatter.TableSorter(schema).OrderBy(schema[0].FieldName).Sort(data)

	table := tableformatter.Table{
		Data:   data,
		Schema: schema,
	}
	return table.RenderTable("Custom ISOs", "", command.GetStringParam(c.Arguments["format"]))
}
