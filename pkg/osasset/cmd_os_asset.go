package osasset

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

var OsAssetsCmds = []command.Command{
	{
		Description:  "Lists all Assets.",
		Subject:      "asset",
		AltSubject:   "asset",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list assets", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"format": c.FlagSet.String("format", command.NilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
				"usage":  c.FlagSet.String("usage", command.NilDefaultStr, "Asset's usage"),
			}
		},
		ExecuteFunc: assetsListCmd,
		Endpoint:    configuration.ExtendedEndpoint,
	},
	{
		Description:  "Create asset.",
		Subject:      "asset",
		AltSubject:   "asset",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("create asset", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"filename":               c.FlagSet.String("filename", command.NilDefaultStr, "Asset's filename"),
				"usage":                  c.FlagSet.String("usage", command.NilDefaultStr, "Asset's usage. Possible values: \"bootloader\""),
				"mime":                   c.FlagSet.String("mime", command.NilDefaultStr, "Asset's mime type. Possible values: \"text/plain\",\"application/octet-stream\""),
				"template_type":          c.FlagSet.String("template-type", command.NilDefaultStr, "Asset's template type. Possible values: \"none\",\"simple\",\"advanced\""),
				"url":                    c.FlagSet.String("url", command.NilDefaultStr, "Asset's source url. If present it will not read content anymore"),
				"read_content_from_pipe": c.FlagSet.Bool("pipe", false, "Read assets's content read from pipe instead of terminal input"),
				"template_id_or_name":    c.FlagSet.String("template-id", command.NilDefaultStr, "Template's id or name to associate. "),
				"path":                   c.FlagSet.String("path", command.NilDefaultStr, "Path to associate asset to."),
				"variables_json":         c.FlagSet.String("variables-json", command.NilDefaultStr, "JSON encoded variables object"),
				"delete_if_exists":       c.FlagSet.Bool("delete-if-exists", false, "Automatically delete the existing asset associated with the current template."),
				"return_id":              c.FlagSet.Bool("return-id", false, colors.Green("(Flag)")+" If set will print the ID of the created infrastructure. Useful for automating tasks."),
			}
		},
		ExecuteFunc: AssetCreateCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
	{
		Description:  "Get asset contents",
		Subject:      "asset",
		AltSubject:   "asset",
		Predicate:    "get-contents",
		AltPredicate: "contents",
		FlagSet:      flag.NewFlagSet("get asset content", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"asset_id_or_name": c.FlagSet.String("id", command.NilDefaultStr, "Asset's id or name"),
			}
		},
		ExecuteFunc: assetGetCmd,
		Endpoint:    configuration.ExtendedEndpoint,
	},
	{
		Description:  "Delete asset.",
		Subject:      "asset",
		AltSubject:   "asset",
		Predicate:    "delete",
		AltPredicate: "rm",
		FlagSet:      flag.NewFlagSet("delete asset", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"asset_id_or_name": c.FlagSet.String("id", command.NilDefaultStr, "Asset's id or name"),
				"autoconfirm":      c.FlagSet.Bool("autoconfirm", false, colors.Green("(Flag)")+" If set it will assume action is confirmed"),
			}
		},
		ExecuteFunc: assetDeleteCmd,
		Endpoint:    configuration.ExtendedEndpoint,
	},
	{
		Description:  "Add (associate) asset to template.",
		Subject:      "asset",
		AltSubject:   "asset",
		Predicate:    "associate",
		AltPredicate: "assign",
		FlagSet:      flag.NewFlagSet("associate template to asset", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"asset_id_or_name":    c.FlagSet.String("id", command.NilDefaultStr, "Asset's id or filename"),
				"template_id_or_name": c.FlagSet.String("template-id", command.NilDefaultStr, "Template's id or name"),
				"path":                c.FlagSet.String("path", command.NilDefaultStr, "Path to associate asset to"),
				"variables_json":      c.FlagSet.String("variables-json", command.NilDefaultStr, "JSON encoded variables object"),
			}
		},
		ExecuteFunc: associateAssetCmd,
		Endpoint:    configuration.ExtendedEndpoint,
	},
	{
		Description:  "Remove (unassign) asset from template.",
		Subject:      "asset",
		AltSubject:   "asset",
		Predicate:    "disassociate",
		AltPredicate: "unassign",
		FlagSet:      flag.NewFlagSet("disassociate asset from template", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"asset_id_or_name":    c.FlagSet.String("id", command.NilDefaultStr, "Asset's id or filename"),
				"template_id_or_name": c.FlagSet.String("template-id", command.NilDefaultStr, "Template's id or name"),
			}
		},
		ExecuteFunc: disassociateAssetCmd,
		Endpoint:    configuration.ExtendedEndpoint,
	},
	{
		Description:  "Edit asset.",
		Subject:      "asset",
		AltSubject:   "asset",
		Predicate:    "edit",
		AltPredicate: "update",
		FlagSet:      flag.NewFlagSet("edit asset", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"asset_id_or_name":       c.FlagSet.String("id", command.NilDefaultStr, "Asset's id or filename"),
				"filename":               c.FlagSet.String("filename", command.NilDefaultStr, "Asset's filename"),
				"usage":                  c.FlagSet.String("usage", command.NilDefaultStr, "Asset's usage. Possible values: \"bootloader\""),
				"mime":                   c.FlagSet.String("mime", command.NilDefaultStr, "Required. Asset's mime type. Possible values: \"text/plain\",\"application/octet-stream\""),
				"template_type":          c.FlagSet.String("template-type", command.NilDefaultStr, "Asset's template type. Possible values: \"none\",\"simple\",\"advanced\""),
				"url":                    c.FlagSet.String("url", command.NilDefaultStr, "Asset's source url. If present it will not read content anymore"),
				"read_content_from_pipe": c.FlagSet.Bool("pipe", false, "Read assets's content read from pipe instead of terminal input"),
				"template_id_or_name":    c.FlagSet.String("template-id", command.NilDefaultStr, "Template's id or name to associate. "),
				"path":                   c.FlagSet.String("path", command.NilDefaultStr, "Path to associate asset to."),
				"variables_json":         c.FlagSet.String("variables-json", command.NilDefaultStr, "JSON encoded variables object"),
				"return_id":              c.FlagSet.Bool("return-id", false, colors.Green("(Flag)")+" If set will print the ID of the created infrastructure. Useful for automating tasks."),
			}
		},
		ExecuteFunc: assetEditCmd,
		Endpoint:    configuration.ExtendedEndpoint,
	},
	{
		Description:  "Allow other users of the platform to use the asset.",
		Subject:      "asset",
		AltSubject:   "asset",
		Predicate:    "make-public",
		AltPredicate: "public",
		FlagSet:      flag.NewFlagSet("make asset public", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"asset_id_or_name": c.FlagSet.String("id", command.NilDefaultStr, "Asset id or name"),
			}
		},
		ExecuteFunc: assetMakePublicCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
	{
		Description:  "Stop other users of the platform from being able to use the asset by allocating a specific owner.",
		Subject:      "asset",
		AltSubject:   "asset",
		Predicate:    "make-private",
		AltPredicate: "private",
		FlagSet:      flag.NewFlagSet("make asset private", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"asset_id_or_name": c.FlagSet.String("id", command.NilDefaultStr, "Asset id or name"),
				"user_id":          c.FlagSet.String("user-id", command.NilDefaultStr, "New owner user id or email."),
			}
		},
		ExecuteFunc: assetMakePrivateCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
}

func assetsListCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

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
			FieldName: "TEMPLATE_TYPE",
			FieldType: tableformatter.TypeString,
			FieldSize: 14,
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
			s.OSAssetTemplateType,
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
	return table.RenderTable("Assets", "", command.GetStringParam(c.Arguments["format"]))
}

func AssetCreateCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	return assetCreate(c, client, []byte{})
}

func AssetCreateWithContentCmd(c *command.Command, client metalcloud.MetalCloudClient, assetContent []byte) (string, error) {
	return assetCreate(c, client, assetContent)
}

func assetCreate(c *command.Command, client metalcloud.MetalCloudClient, assetContent []byte) (string, error) {
	newObj := metalcloud.OSAsset{}
	updatedObj, err := updateAssetFromCommand(newObj, c, client, true, assetContent)
	if err != nil {
		return "", err
	}

	ret, err := client.OSAssetCreate(*updatedObj)

	if err != nil {
		return "", err
	}

	err = associateAssetFromCommand(ret.OSAssetID, ret.OSAssetFileName, c, client)

	if err != nil {
		return "", err
	}

	if command.GetBoolParam(c.Arguments["return_id"]) {
		return fmt.Sprintf("%d", ret.OSAssetID), nil
	}

	return "", err
}

func associateAssetFromCommand(assetID int, assetFileName string, c *command.Command, client metalcloud.MetalCloudClient) error {
	variablesJSON := "[]"
	templateIsPublic := false
	if _, error := command.GetParam(c, "template_id_or_name", "template-id"); error == nil {
		template, err := command.GetOSTemplateFromCommand("template-id", c, client, false)
		if err != nil {
			return err
		}

		templateIsPublic = template.UserID == 0
		path, ok := command.GetStringParamOk(c.Arguments["path"])
		if !ok {
			return fmt.Errorf("-path is required")
		}

		if v, ok := command.GetStringParamOk(c.Arguments["variables_json"]); ok {
			variablesJSON = v
		}

		if del := command.GetBoolParam(c.Arguments["delete_if_exists"]); del {
			list, err := client.OSTemplateOSAssets(template.VolumeTemplateID)

			if err != nil {
				return err
			}

			for _, a := range *list {
				if a.OSAsset.OSAssetFileName == assetFileName || a.OSAssetFilePath == path {
					bUpdateTemplate := false
					if template.OSAssetBootloaderLocalInstall == a.OSAsset.OSAssetID {
						template.OSAssetBootloaderLocalInstall = 0
						bUpdateTemplate = true
					}
					if template.OSAssetBootloaderOSBoot == a.OSAsset.OSAssetID {
						template.OSAssetBootloaderOSBoot = 0
						bUpdateTemplate = true
					}
					if bUpdateTemplate {
						template, err = client.OSTemplateUpdate(template.VolumeTemplateID, *template)

						if err != nil {
							return err
						}
					}
					err = client.OSTemplateRemoveOSAsset(template.VolumeTemplateID, a.OSAsset.OSAssetID)

					if err != nil {
						return err
					}
				}
			}
		}

		if templateIsPublic {
			_, err = client.OSAssetMakePublic(assetID)

			if err != nil {
				return err
			}
		}

		err = client.OSTemplateAddOSAsset(template.VolumeTemplateID, assetID, path, variablesJSON)

		if err != nil {
			return err
		}
	}
	return nil
}

func assetGetCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	retS, err := command.GetOSAssetFromCommand("id", "asset_id_or_name", c, client)
	if err != nil {
		return "", err
	}

	if retS.OSAssetSourceURL != "" {
		return "", fmt.Errorf("No stored content. This command can only be used for assets that have content stored in the database. This asset is being pulled from '%s'.", retS.OSAssetSourceURL)
	}

	content, err := client.OSAssetGetStoredContent(retS.OSAssetID)
	if err != nil {
		return "", err
	}

	sDec, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return "", err
	}

	fmt.Print(string(sDec))

	return "", err
}

func assetDeleteCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	retS, err := command.GetOSAssetFromCommand("id", "asset_id_or_name", c, client)
	if err != nil {
		return "", err
	}

	confirm, err := command.ConfirmCommand(c, func() string {

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

func updateAssetFromCommand(obj metalcloud.OSAsset, c *command.Command, client metalcloud.MetalCloudClient, checkRequired bool, assetContent []byte) (*metalcloud.OSAsset, error) {
	if v, ok := command.GetStringParamOk(c.Arguments["filename"]); ok {
		obj.OSAssetFileName = v
	} else {
		if checkRequired {
			return nil, fmt.Errorf("-filename is required")
		}
	}

	if v, ok := command.GetStringParamOk(c.Arguments["template_type"]); ok {
		obj.OSAssetTemplateType = v
	} else {
		if checkRequired {
			return nil, fmt.Errorf("--template-type is required")
		}
	}

	if v, ok := command.GetStringParamOk(c.Arguments["usage"]); ok {
		obj.OSAssetUsage = v
	}

	if v, ok := command.GetStringParamOk(c.Arguments["mime"]); ok {
		obj.OSAssetFileMime = v
	}

	content := assetContent
	var err error

	if v, ok := command.GetStringParamOk(c.Arguments["url"]); ok {
		obj.OSAssetSourceURL = v
	} else {
		if len(content) == 0 {
			if command.GetBoolParam(c.Arguments["read_content_from_pipe"]) {
				_content, err := configuration.ReadInputFromPipe()
				if err != nil {
					return nil, err
				}
				content = _content
			} else {
				if runtime.GOOS == "windows" {
					content, err = command.RequestInput("Asset content:")

				} else {
					content, err = command.RequestInputSilent("Asset content:")
				}

				if err != nil {
					return nil, err
				}
			}
		}

		obj.OSAssetContentsBase64 = base64.StdEncoding.EncodeToString([]byte(content))
	}

	return &obj, nil
}

func associateAssetCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	asset, err := command.GetOSAssetFromCommand("id", "asset_id_or_name", c, client)
	if err != nil {
		return "", err
	}

	template, err := command.GetOSTemplateFromCommand("template-id", c, client, false)
	if err != nil {
		return "", err
	}

	path, ok := command.GetStringParamOk(c.Arguments["path"])
	if !ok {
		return "", fmt.Errorf("-path is required")
	}

	variablesJSON := "[]"
	if v, ok := command.GetStringParamOk(c.Arguments["variables_json"]); ok {
		variablesJSON = v
	}

	return "", client.OSTemplateAddOSAsset(template.VolumeTemplateID, asset.OSAssetID, path, variablesJSON)
}

func disassociateAssetCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	asset, err := command.GetOSAssetFromCommand("id", "asset_id_or_name", c, client)
	if err != nil {
		return "", err
	}

	template, err := command.GetOSTemplateFromCommand("template-id", c, client, false)
	if err != nil {
		return "", err
	}

	return "", client.OSTemplateRemoveOSAsset(template.VolumeTemplateID, asset.OSAssetID)
}

func assetEditCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	asset, err := command.GetOSAssetFromCommand("id", "asset_id_or_name", c, client)
	if err != nil {
		return "", err
	}

	if err != nil {
		return "", err
	}

	newObj := metalcloud.OSAsset{}

	updatedObj, err := updateAssetFromCommand(newObj, c, client, true, []byte{})

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

	if command.GetBoolParam(c.Arguments["return_id"]) {
		return fmt.Sprintf("%d", ret.OSAssetID), nil
	}

	return "", err
}

func assetMakePublicCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	asset, err := command.GetOSAssetFromCommand("id", "asset_id_or_name", c, client)
	if err != nil {
		return "", err
	}

	asset, err = client.OSAssetMakePublic(asset.OSAssetID)

	if err != nil {
		return "", err
	}

	return "", nil
}

func assetMakePrivateCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	asset, err := command.GetOSAssetFromCommand("id", "asset_id_or_name", c, client)
	if err != nil {
		return "", err
	}

	user, err := command.GetUserFromCommand("user-id", c, client)

	if err != nil {
		return "", err
	}

	asset, err = client.OSAssetMakePrivate(asset.OSAssetID, user.UserID)

	if err != nil {
		return "", err
	}

	return "", nil
}
