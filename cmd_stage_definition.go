package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"strings"

	metalcloud "github.com/bigstepinc/metal-cloud-sdk-go"
	interfaces "github.com/bigstepinc/metalcloud-cli/interfaces"
)

var stageDefinitionsCmds = []Command{

	Command{
		Description:  "Lists available stage definitions",
		Subject:      "stages",
		AltSubject:   "stage_definitions",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list stage definitions", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"format": c.FlagSet.String("format", _nilDefaultStr, "The output format. Supported values are 'json','csv'. The default format is human readable."),
			}
		},
		ExecuteFunc: stageDefinitionsListCmd,
	},
	Command{
		Description:  "Create stage definition",
		Subject:      "stage",
		AltSubject:   "stage_definition",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("create stage definition", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"label":                             c.FlagSet.String("label", _nilDefaultStr, "Stage Definitions's label"),
				"icon":                              c.FlagSet.String("icon", _nilDefaultStr, "Icon image file in data URI format like this: data:image/png;base64,iVBOR="),
				"title":                             c.FlagSet.String("title", _nilDefaultStr, "Stage Definitions's title"),
				"description":                       c.FlagSet.String("description", _nilDefaultStr, "Stage Definitions's description"),
				"type":                              c.FlagSet.String("type", _nilDefaultStr, "Stage Definitions's type. Possible values: AnsibleBundle, HTTPRequest"),
				"vars":                              c.FlagSet.String("vars", _nilDefaultStr, "Stage Definitions's variables. These must be available in the execution context, otherwise the stage cannot run."),
				"ansible_bundle_filename":           c.FlagSet.String("ansible_bundle_filename", _nilDefaultStr, "Ansible bundle's file path to load the bundle from. Must be a zip file. Required when type=AnsibleBundle"),
				"http_request_url":                  c.FlagSet.String("http_request_url", _nilDefaultStr, "HTTP Requests's URL. Required when using type=HTTPRequest"),
				"http_request_method":               c.FlagSet.String("http_request_method", _nilDefaultStr, "HTTP Requests's method. Required when using type=HTTPRequest"),
				"http_request_body_filename":        c.FlagSet.String("http_request_body_filename", _nilDefaultStr, "HTTP Requests's content is read from this file. Can only be used when type=HTTPRequest"),
				"http_request_body_from_pipe":       c.FlagSet.Bool("http_request_body_from_pipe", false, "HTTP Requests's content is read from stdin. Can only be used when type=HTTPRequest"),
				"http_request_header_accept":        c.FlagSet.String("http_request_header_accept", _nilDefaultStr, "HTTP Requests's Accept header. Can only be used when type=HTTPRequest"),
				"http_request_header_authorization": c.FlagSet.String("http_request_header_authorization", _nilDefaultStr, "HTTP Requests's Authorization header. Can only be used when type=HTTPRequest"),
				"http_request_header_cookie":        c.FlagSet.String("http_request_header_cookie", _nilDefaultStr, "HTTP Requests's Cookie header. Can only be used when type=HTTPRequest"),
				"http_request_header_user_agent":    c.FlagSet.String("http_request_header_user_agent", _nilDefaultStr, "HTTP Requests's User-Agent header. Can only be used when type=HTTPRequest"),
				"http_request_redirect":             c.FlagSet.String("http_request_redirect", _nilDefaultStr, "HTTP Requests's method. Can only be used when type=HTTPRequest"),
				"http_request_follow":               c.FlagSet.Int("http_request_follow", _nilDefaultInt, "HTTP Requests's follow. Can only be used when type=HTTPRequest"),
				"http_request_no_compress":          c.FlagSet.Bool("http_request_no_compress", false, "HTTP Requests's compress disabled if set. Can only be used when type=HTTPRequest"),
				"http_request_timeout":              c.FlagSet.Int("http_request_timeout", _nilDefaultInt, "HTTP Requests's timeout. Can only be used when type=HTTPRequest"),
				"http_request_size":                 c.FlagSet.Int("http_request_size", _nilDefaultInt, "HTTP Requests's size. Can only be used when type=HTTPRequest"),
			}
		},
		ExecuteFunc: stageDefinitionCreateCmd,
	},
	Command{
		Description:  "Delete stage definition",
		Subject:      "stage",
		AltSubject:   "stage_definition",
		Predicate:    "delete",
		AltPredicate: "rm",
		FlagSet:      flag.NewFlagSet("delete stage", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"stage_id_or_name": c.FlagSet.String("id", _nilDefaultStr, "stage's id or name"),
				"autoconfirm":      c.FlagSet.Bool("autoconfirm", false, "If true it does not ask for confirmation anymore"),
			}
		},
		ExecuteFunc: stageDefinitionDeleteCmd,
	},
	Command{
		Description:  "Add stage into workflow",
		Subject:      "stage",
		AltSubject:   "stage_definition",
		Predicate:    "add",
		AltPredicate: "add_to_infrastructure",
		FlagSet:      flag.NewFlagSet("add stage", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"stage_id_or_name":           c.FlagSet.String("id", _nilDefaultStr, "stage's id or name"),
				"infrastructure_id_or_label": c.FlagSet.String("infra_id", _nilDefaultStr, "The infrastructure's id"),
				"runlevel":                   c.FlagSet.Int("runlevel", _nilDefaultInt, "The runlevel"),
				"runmoment":                  c.FlagSet.String("moment", _nilDefaultStr, "When to run the stage"),
			}
		},
		ExecuteFunc: stageDefinitionAddToWorkflowCmd,
	},
	Command{
		Description:  "Add stage into workflow",
		Subject:      "stage",
		AltSubject:   "stage_definition",
		Predicate:    "add",
		AltPredicate: "add_to_infrastructure",
		FlagSet:      flag.NewFlagSet("add stage", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"stage_id_or_name":           c.FlagSet.String("id", _nilDefaultStr, "stage's id or name"),
				"infrastructure_id_or_label": c.FlagSet.String("infra_id", _nilDefaultStr, "The infrastructure's id"),
				"runlevel":                   c.FlagSet.Int("runlevel", _nilDefaultInt, "The runlevel"),
				"group":                      c.FlagSet.String("group", _nilDefaultStr, "When to run the stage"),
			}
		},
		ExecuteFunc: stageDefinitionAddToWorkflowCmd,
	},
}

func stageDefinitionsListCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	list, err := client.StageDefinitions()

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
			FieldName: "LABEL",
			FieldType: TypeString,
			FieldSize: 6,
		},
		SchemaField{
			FieldName: "TITLE",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "DESCRIPTION",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "TYPE",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "VARS_REQUIRED",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "DEF.",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "CREATED",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "UPDATED",
			FieldType: TypeString,
			FieldSize: 4,
		},
	}

	data := [][]interface{}{}
	for _, s := range *list {

		stageDef := ""
		switch s.StageDefinitionType {
		case "AnsibleBundle":
			bundle := s.StageDefinition.(metalcloud.AnsibleBundle)
			stageDef = fmt.Sprintf("Ansible Bundle Filename: %s", bundle.AnsibleBundleArchiveFilename)
		case "HTTPRequest":
			req := s.StageDefinition.(metalcloud.HTTPRequest)
			stageDef = fmt.Sprintf("HTTP Request URI: %s", req.URL)
		}

		data = append(data, []interface{}{
			s.StageDefinitionID,
			s.StageDefinitionLabel,
			s.StageDefinitionTitle,
			s.StageDefinitionDescription,
			s.StageDefinitionType,
			strings.Join(s.StageDefinitionVariablesNamesRequired, ","),
			stageDef,
			s.StageDefinitionCreatedTimestamp,
			s.StageDefinitionUpdatedTimestamp,
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
		sb.WriteString(fmt.Sprintf("Stage Definitions:"))

		TableSorter(schema).OrderBy(
			schema[0].FieldName,
			schema[1].FieldName).Sort(data)

		AdjustFieldSizes(data, &schema)

		sb.WriteString(GetTableAsString(data, schema))

		sb.WriteString(fmt.Sprintf("Total: %d stage definitions\n\n", len(*list)))
	}

	return sb.String(), nil
}

func stageDefinitionCreateCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	stage := metalcloud.StageDefinition{}

	if v := c.Arguments["label"]; v != nil && *v.(*string) != _nilDefaultStr {
		stage.StageDefinitionLabel = *v.(*string)
	} else {
		return "", fmt.Errorf("label is required")
	}

	if v := c.Arguments["icon"]; v != nil && *v.(*string) != _nilDefaultStr {
		stage.IconAssetDataURI = *v.(*string)
	}

	if v := c.Arguments["title"]; v != nil && *v.(*string) != _nilDefaultStr {
		stage.StageDefinitionTitle = *v.(*string)
	} else {
		return "", fmt.Errorf("title is required")
	}

	if v := c.Arguments["description"]; v != nil && *v.(*string) != _nilDefaultStr {
		stage.StageDefinitionDescription = *v.(*string)
	}

	if v := c.Arguments["type"]; v != nil && *v.(*string) != _nilDefaultStr {
		stage.StageDefinitionType = *v.(*string)
	} else {
		return "", fmt.Errorf("type is required")
	}

	if v := c.Arguments["vars"]; v != nil && *v.(*string) != _nilDefaultStr {
		vars := *v.(*string)
		stage.StageDefinitionVariablesNamesRequired = strings.Split(vars, ",")
	}

	switch stage.StageDefinitionType {
	case "AnsibleBundle":
		if v := c.Arguments["ansible_bundle_filename"]; v != nil && *v.(*string) != _nilDefaultStr {
			ab := metalcloud.AnsibleBundle{}

			ab.AnsibleBundleArchiveFilename = *v.(*string)

			content, err := requestInputFromFile(ab.AnsibleBundleArchiveFilename)
			if err != nil {
				return "", err
			}

			ab.AnsibleBundleArchiveContentsBase64 = base64.StdEncoding.EncodeToString(content)
			ab.Type = "AnsibleBundle"
			stage.StageDefinition = ab
		}
	case "HTTPRequest":

		req := metalcloud.HTTPRequest{}
		req.Type = "HTTPRequest"
		req.Options = metalcloud.WebFetchAAPIOptions{}

		content := []byte{}
		if v := c.Arguments["http_request_body_from_pipe"]; *v.(*bool) {
			content = readInputFromPipe()
		} else {
			v := c.Arguments["http_request_body_filename"]
			if *v.(*string) != _nilDefaultStr {

				filename := *v.(*string)
				c, err := requestInputFromFile(filename)
				if err != nil {
					return "", err
				}

				content = c
			}
		}

		req.Options.BodyBufferBase64 = base64.StdEncoding.EncodeToString(content)

		if v := c.Arguments["http_request_url"]; v != nil && *v.(*string) != _nilDefaultStr {
			req.URL = *v.(*string)
		} else {
			return "", fmt.Errorf("http_request_url is required if using HTTPRequest")
		}

		if v := c.Arguments["http_request_method"]; v != nil && *v.(*string) != _nilDefaultStr {
			req.Options.Method = *v.(*string)
		} else {
			return "", fmt.Errorf("http_request_method is required if using HTTPRequest")
		}

		if v := c.Arguments["http_request_redirect"]; v != nil && *v.(*string) != _nilDefaultStr {
			req.Options.Redirect = *v.(*string)
		}

		if v := c.Arguments["http_request_follow"]; v != nil && *v.(*int) != _nilDefaultInt {
			req.Options.Follow = *v.(*int)
		}

		if v := c.Arguments["http_request_no_compress"]; v != nil && *v.(*bool) {
			req.Options.Compress = *v.(*bool)
		}

		if v := c.Arguments["http_request_timeout"]; v != nil && *v.(*int) != _nilDefaultInt {
			req.Options.Timeout = *v.(*int)
		}

		if v := c.Arguments["http_request_size"]; v != nil && *v.(*int) != _nilDefaultInt {
			req.Options.Size = *v.(*int)
		}

		req.Options.Headers = metalcloud.WebFetchAPIRequestHeaders{}

		if v := c.Arguments["http_request_header_accept"]; v != nil && *v.(*string) != _nilDefaultStr {
			req.Options.Headers.Accept = *v.(*string)
		}

		if v := c.Arguments["http_request_header_authorization"]; v != nil && *v.(*string) != _nilDefaultStr {
			req.Options.Headers.Authorization = *v.(*string)
		}

		if v := c.Arguments["http_request_header_cookie"]; v != nil && *v.(*string) != _nilDefaultStr {
			req.Options.Headers.Cookie = *v.(*string)
		}

		if v := c.Arguments["http_request_header_user_agent"]; v != nil && *v.(*string) != _nilDefaultStr {
			req.Options.Headers.UserAgent = *v.(*string)
		}

		stage.StageDefinition = req
	}

	_, err := client.StageDefinitionCreate(stage)

	return "", err
}

func stageDefinitionDeleteCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	retS, err := getStageDefinitionFromCommand("id", c, client)
	if err != nil {
		return "", err
	}
	confirm := false

	if c.Arguments["autoconfirm"] != nil && *c.Arguments["autoconfirm"].(*bool) == true {
		confirm = true
	} else {

		confirmationMessage := fmt.Sprintf("Deleting stage definition %s (%d).  Are you sure? Type \"yes\" to continue:",
			retS.StageDefinitionLabel,
			retS.StageDefinitionID)

		//this is simply so that we don't output a text on the command line under go test
		if strings.HasSuffix(os.Args[0], ".test") {
			confirmationMessage = ""
		}

		confirm = requestConfirmation(confirmationMessage)
	}

	if !confirm {
		return "", fmt.Errorf("Operation not confirmed. Aborting")
	}

	err = client.StageDefinitionDelete(retS.StageDefinitionID)

	return "", err
}

func stageDefinitionAddToWorkflowCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	stage, err := getStageDefinitionFromCommand("id", c, client)
	if err != nil {
		return "", err
	}

	infra, err := getInfrastructureFromCommand("infra_id", c, client)
	if err != nil {
		return "", err
	}

	runlevel := 0
	if v := c.Arguments["runlevel"]; v != nil && *v.(*int) != _nilDefaultInt {
		runlevel = *v.(*int)
	}

	runmoment := "post_deploy"
	if v := c.Arguments["group"]; v != nil && *v.(*string) != _nilDefaultStr {
		runmoment = *v.(*string)
	}

	err = client.InfrastructureDeployCustomStageAddIntoRunlevel(infra.InfrastructureID, stage.StageDefinitionID, runlevel, runmoment)

	return "", err
}

func getStageDefinitionFromCommand(paramName string, c *Command, client interfaces.MetalCloudClient) (*metalcloud.StageDefinition, error) {

	v, err := getParam(c, "stage_id_or_name", paramName)
	if err != nil {
		return nil, err
	}

	id, label, isID := idOrLabel(v)

	if isID {
		return client.StageDefinitionGet(id)
	}

	secrets, err := client.StageDefinitions()
	if err != nil {
		return nil, err
	}

	for _, s := range *secrets {
		if s.StageDefinitionLabel == label {
			return &s, nil
		}
	}

	return nil, fmt.Errorf("Could not locate stage definition with id/name %v", *v.(*interface{}))
}
