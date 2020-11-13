package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"strings"

	metalcloud "github.com/bigstepinc/metal-cloud-sdk-go"
	interfaces "github.com/bigstepinc/metalcloud-cli/interfaces"
	"github.com/metalsoft-io/tableformatter"
)

var osAssetsCmds = []Command{

	{
		Description:  "Lists all Assets.",
		Subject:      "asset",
		AltSubject:   "asset",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list assets", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"format": c.FlagSet.String("format", _nilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
				"usage":  c.FlagSet.String("usage", _nilDefaultStr, "Asset's usage"),
			}
		},
		ExecuteFunc: assetsListCmd,
		Endpoint:    ExtendedEndpoint,
	},
	{
		Description:  "Create asset.",
		Subject:      "asset",
		AltSubject:   "asset",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("create asset", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"filename":                c.FlagSet.String("filename", _nilDefaultStr, "Asset's filename"),
				"usage":                   c.FlagSet.String("usage", _nilDefaultStr, "Asset's usage. Possible values: \"bootloader\""),
				"mime":                    c.FlagSet.String("mime", _nilDefaultStr, "Asset's mime type. Possible values: \"text/plain\",\"application/octet-stream\""),
				"url":                     c.FlagSet.String("url", _nilDefaultStr, "Asset's source url. If present it will not read content anymore"),
				"variable_names_required": c.FlagSet.String("variable-names-required", _nilDefaultStr, "The names of the variables and secrets that are used in this asset, comma separated."),
				"read_content_from_pipe":  c.FlagSet.Bool("pipe", false, "Read assets's content read from pipe instead of terminal input"),
				"template_id_or_name":     c.FlagSet.String("template-id", _nilDefaultStr, "Template's id or name to associate. "),
				"path":                    c.FlagSet.String("path", _nilDefaultStr, "Path to associate asset to."),
				"variables_json":          c.FlagSet.String("variables-json", _nilDefaultStr, "JSON encoded variables object"),
				"delete_if_exists":        c.FlagSet.Bool("delete-if-exists", true, "Automatically delete the existing asset associated with the current template."),
				"return_id":               c.FlagSet.Bool("return-id", false, "(Flag) If set will print the ID of the created infrastructure. Useful for automating tasks."),
			}
		},
		ExecuteFunc: assetCreateCmd,
		Endpoint:    ExtendedEndpoint,
	},
	{
		Description:  "Delete asset.",
		Subject:      "asset",
		AltSubject:   "asset",
		Predicate:    "delete",
		AltPredicate: "rm",
		FlagSet:      flag.NewFlagSet("delete asset", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"asset_id_or_name": c.FlagSet.String("id", _nilDefaultStr, "Asset's id or name"),
				"autoconfirm":      c.FlagSet.Bool("autoconfirm", false, "If true it does not ask for confirmation anymore"),
			}
		},
		ExecuteFunc: assetDeleteCmd,
		Endpoint:    ExtendedEndpoint,
	},
	{
		Description:  "Add (associate) asset to template.",
		Subject:      "asset",
		AltSubject:   "asset",
		Predicate:    "associate",
		AltPredicate: "assign",
		FlagSet:      flag.NewFlagSet("associate template to asset", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"asset_id_or_name":    c.FlagSet.String("id", _nilDefaultStr, "Asset's id or filename"),
				"template_id_or_name": c.FlagSet.String("template-id", _nilDefaultStr, "Template's id or name"),
				"path":                c.FlagSet.String("path", _nilDefaultStr, "Path to associate asset to"),
				"variables_json":      c.FlagSet.String("variables-json", _nilDefaultStr, "JSON encoded variables object"),
			}
		},
		ExecuteFunc: associateAssetCmd,
		Endpoint:    ExtendedEndpoint,
	},
	{
		Description:  "Remove (unassign) asset from template.",
		Subject:      "asset",
		AltSubject:   "asset",
		Predicate:    "disassociate",
		AltPredicate: "unassign",
		FlagSet:      flag.NewFlagSet("disassociate asset from template", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"asset_id_or_name":    c.FlagSet.String("id", _nilDefaultStr, "Asset's id or filename"),
				"template_id_or_name": c.FlagSet.String("template-id", _nilDefaultStr, "Template's id or name"),
			}
		},
		ExecuteFunc: disassociateAssetCmd,
		Endpoint:    ExtendedEndpoint,
	},
	{
		Description:  "List associated assets.",
		Subject:      "asset",
		AltSubject:   "asset",
		Predicate:    "list-associated",
		AltPredicate: "assoc",
		FlagSet:      flag.NewFlagSet("associated assets", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"template_id_or_name": c.FlagSet.String("id", _nilDefaultStr, "Template's id or name"),
			}
		},
		ExecuteFunc: templateListAssociatedAssetsCmd,
		Endpoint:    ExtendedEndpoint,
	},
	{
		Description:  "Edit asset.",
		Subject:      "asset",
		AltSubject:   "asset",
		Predicate:    "edit",
		AltPredicate: "update",
		FlagSet:      flag.NewFlagSet("edit asset", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"asset_id_or_name":        c.FlagSet.String("id", _nilDefaultStr, "Asset's id or filename"),
				"filename":                c.FlagSet.String("filename", _nilDefaultStr, "Asset's filename"),
				"usage":                   c.FlagSet.String("usage", _nilDefaultStr, "Asset's usage. Possible values: \"bootloader\""),
				"mime":                    c.FlagSet.String("mime", _nilDefaultStr, "Required. Asset's mime type. Possible values: \"text/plain\",\"application/octet-stream\""),
				"url":                     c.FlagSet.String("url", _nilDefaultStr, "Asset's source url. If present it will not read content anymore"),
				"variable_names_required": c.FlagSet.String("variable-names-required", _nilDefaultStr, "The names of the variables and secrets that are used in this asset, comma separated."),
				"read_content_from_pipe":  c.FlagSet.Bool("pipe", false, "Read assets's content read from pipe instead of terminal input"),
				"template_id_or_name":     c.FlagSet.String("template-id", _nilDefaultStr, "Template's id or name to associate. "),
				"path":                    c.FlagSet.String("path", _nilDefaultStr, "Path to associate asset to."),
				"variables_json":          c.FlagSet.String("variables-json", _nilDefaultStr, "JSON encoded variables object"),
				"return_id":               c.FlagSet.Bool("return-id", false, "(Flag) If set will print the ID of the created infrastructure. Useful for automating tasks."),
			}
		},
		ExecuteFunc: assetEditCmd,
		Endpoint:    ExtendedEndpoint,
	},
}

func assetsListCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	list, err := client.OSAssets()

	if err != nil {
		return "", err
	}

	schema := []tableformatter.SchemaField{
		{
			FieldName: "ID",
			FieldType: tableformatter.TypeInt,
			FieldSize: 2,
		},
		{
			FieldName: "FILENAME",
			FieldType: tableformatter.TypeString,
			FieldSize: 20,
		},
		{
			FieldName: "FILE_SIZE_BYTES",
			FieldType: tableformatter.TypeInt,
			FieldSize: 4,
		},
		{
			FieldName: "FILE_MIME",
			FieldType: tableformatter.TypeString,
			FieldSize: 20,
		},
		{
			FieldName: "USAGE",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "SOURCE_URL",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "VARIABLE_NAMES_REQUIRED",
			FieldType: tableformatter.TypeString,
			FieldSize: 10,
		},
	}

	data := [][]interface{}{}
	for _, s := range *list {

		data = append(data, []interface{}{
			s.OSAssetID,
			s.OSAssetFileName,
			s.OSAssetFileSizeBytes,
			s.OSAssetFileMime,
			s.OSAssetUsage,
			s.OSAssetSourceURL,
			strings.Join(s.OSAssetVariableNamesRequired, ","),
		})

	}

	tableformatter.TableSorter(schema).OrderBy(schema[0].FieldName).Sort(data)

	table := tableformatter.Table{
		Data:   data,
		Schema: schema,
	}
	return table.RenderTable("Assets", "", getStringParam(c.Arguments["format"]))
}

func assetCreateCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {
	newObj := metalcloud.OSAsset{}
	updatedObj, err := updateAssetFromCommand(newObj, c, client, true)

	ret, err := client.OSAssetCreate(*updatedObj)

	if err != nil {
		return "", err
	}

	err = associateAssetFromCommand(ret.OSAssetID, ret.OSAssetFileName, c, client)

	if err != nil {
		return "", err
	}

	if getBoolParam(c.Arguments["return_id"]) {
		return fmt.Sprintf("%d", ret.OSAssetID), nil
	}

	return "", err
}

func associateAssetFromCommand(assetID int, assetFileName string, c *Command, client interfaces.MetalCloudClient) error {
	variablesJSON := "[]"
	if _, error := getParam(c, "template_id_or_name", "template-id"); error == nil {
		template, err := getOSTemplateFromCommand("template-id", c, client, false)
		if err != nil {
			return err
		}
		path, ok := getStringParamOk(c.Arguments["path"])
		if !ok {
			return fmt.Errorf("-path is required")
		}

		if v, ok := getStringParamOk(c.Arguments["variables_json"]); ok {
			variablesJSON = v
		}

		if del := getBoolParam(c.Arguments["delete_if_exists"]); del {
			list, err := client.OSTemplateOSAssets(template.VolumeTemplateID)

			if err != nil {
				return err
			}

			for _, a := range *list {
				if a.OSAsset.OSAssetFileName == assetFileName {
					err = client.OSTemplateRemoveOSAsset(template.VolumeTemplateID, a.OSAsset.OSAssetID)

					if err != nil {
						return err
					}
				}
			}
		}

		err = client.OSTemplateAddOSAsset(template.VolumeTemplateID, assetID, path, variablesJSON)

		if err != nil {
			return err
		}
	}
	return nil
}

func assetDeleteCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	retS, err := getOSAssetFromCommand("id", "asset_id_or_name", c, client)
	if err != nil {
		return "", err
	}

	confirm, err := confirmCommand(c, func() string {

		confirmationMessage := fmt.Sprintf("Deleting asset  %s (%d).  Are you sure? Type \"yes\" to continue:",
			retS.OSAssetFileName,
			retS.OSAssetID)

		//this is simply so that we don't output a text on the command line under go test
		if strings.HasSuffix(os.Args[0], ".test") {
			confirmationMessage = ""
		}

		return confirmationMessage
	})

	if err != nil {
		return "", err
	}

	if confirm {
		err = client.OSAssetDelete(retS.OSAssetID)
	}

	return "", err
}

//asset_id_or_name
func getOSAssetFromCommand(paramName string, internalParamName string, c *Command, client interfaces.MetalCloudClient) (*metalcloud.OSAsset, error) {

	v, err := getParam(c, internalParamName, paramName)
	if err != nil {
		return nil, err
	}

	id, name, isID := idOrLabel(v)

	if isID {
		return client.OSAssetGet(id)
	}

	list, err := client.OSAssets()
	if err != nil {
		return nil, err
	}

	for _, s := range *list {
		if s.OSAssetFileName == name {
			return &s, nil
		}
	}

	return nil, fmt.Errorf("Could not locate asset with file name '%s'", name)
}

func updateAssetFromCommand(obj metalcloud.OSAsset, c *Command, client interfaces.MetalCloudClient, checkRequired bool) (*metalcloud.OSAsset, error) {
	if v, ok := getStringParamOk(c.Arguments["filename"]); ok {
		obj.OSAssetFileName = v
	} else {
		if checkRequired {
			return nil, fmt.Errorf("-filename is required")
		}
	}

	if v, ok := getStringParamOk(c.Arguments["usage"]); ok {
		obj.OSAssetUsage = v
	}

	if v, ok := getStringParamOk(c.Arguments["mime"]); ok {
		obj.OSAssetFileMime = v
	}

	content := []byte{}

	if v, ok := getStringParamOk(c.Arguments["url"]); ok {
		obj.OSAssetSourceURL = v
	} else {
		if getBoolParam(c.Arguments["read_content_from_pipe"]) {
			_content, err := readInputFromPipe()
			if err != nil {
				return nil, err
			}
			content = _content
		} else {
			_content, err := requestInputSilent("Asset content:")
			if err != nil {
				return nil, err
			}
			content = _content
		}

		obj.OSAssetContentsBase64 = base64.StdEncoding.EncodeToString([]byte(content))

		if v, ok := getStringParamOk(c.Arguments["variable_names_required"]); ok {
			obj.OSAssetVariableNamesRequired = strings.Split(v, ",")
		}
	}

	return &obj, nil
}

func associateAssetCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	asset, err := getOSAssetFromCommand("id", "asset_id_or_name", c, client)
	if err != nil {
		return "", err
	}

	template, err := getOSTemplateFromCommand("template-id", c, client, false)
	if err != nil {
		return "", err
	}

	path, ok := getStringParamOk(c.Arguments["path"])
	if !ok {
		return "", fmt.Errorf("-path is required")
	}

	variablesJSON := getStringParam(c.Arguments["variables_json"])

	return "", client.OSTemplateAddOSAsset(template.VolumeTemplateID, asset.OSAssetID, path, variablesJSON)
}

func disassociateAssetCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	asset, err := getOSAssetFromCommand("id", "asset_id_or_name", c, client)
	if err != nil {
		return "", err
	}

	template, err := getOSTemplateFromCommand("template-id", c, client, false)
	if err != nil {
		return "", err
	}

	return "", client.OSTemplateRemoveOSAsset(template.VolumeTemplateID, asset.OSAssetID)
}

func templateListAssociatedAssetsCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	ret, err := getOSTemplateFromCommand("id", c, client, false)
	if err != nil {
		return "", err
	}

	list, err := client.OSTemplateOSAssets(ret.VolumeTemplateID)

	if err != nil {
		return "", err
	}

	schema := []tableformatter.SchemaField{
		{
			FieldName: "PATH",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "ID",
			FieldType: tableformatter.TypeInt,
			FieldSize: 2,
		},
		{
			FieldName: "FILENAME",
			FieldType: tableformatter.TypeString,
			FieldSize: 20,
		},
		{
			FieldName: "FILE_SIZE_BYTES",
			FieldType: tableformatter.TypeInt,
			FieldSize: 4,
		},
		{
			FieldName: "FILE_MIME",
			FieldType: tableformatter.TypeString,
			FieldSize: 20,
		},
		{
			FieldName: "USAGE",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "SOURCE_URL",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "VARIABLES_JSON",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
	}

	data := [][]interface{}{}
	for path, s := range *list {

		data = append(data, []interface{}{
			path,
			s.OSAsset.OSAssetID,
			s.OSAsset.OSAssetFileName,
			s.OSAsset.OSAssetFileSizeBytes,
			s.OSAsset.OSAssetFileMime,
			s.OSAsset.OSAssetUsage,
			s.OSAsset.OSAssetSourceURL,
			s.OSTemplateOSAssetVariablesJSON,
		})

	}

	tableformatter.TableSorter(schema).OrderBy(
		schema[0].FieldName,
		schema[1].FieldName).Sort(data)

	table := tableformatter.Table{
		Data:   data,
		Schema: schema,
	}
	return table.RenderTable("Associated assets", "", getStringParam(c.Arguments["format"]))
}

func assetEditCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {
	asset, err := getOSAssetFromCommand("id", "asset_id_or_name", c, client)
	if err != nil {
		fmt.Println(err)
	}

	if err != nil {
		return "", err
	}

	newObj := metalcloud.OSAsset{}

	updatedObj, err := updateAssetFromCommand(newObj, c, client, true)

	if err != nil {
		return "", err
	}

	ret, err := client.OSAssetUpdate(asset.OSAssetID, *updatedObj)

	if err != nil {
		return "", err
	}

	err = associateAssetFromCommand(ret.OSAssetID, ret.OSAssetFileName, c, client)

	if err != nil {
		return "", err
	}

	if getBoolParam(c.Arguments["return_id"]) {
		return fmt.Sprintf("%d", ret.OSAssetID), nil
	}

	return "", err
}
