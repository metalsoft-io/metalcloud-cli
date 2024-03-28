package stagedefinition

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
	"github.com/metalsoft-io/tableformatter"
)

var StageDefinitionsCmds = []command.Command{
	{
		Description:  "Lists all stage definitions.",
		Subject:      "stage-definition",
		AltSubject:   "stagedef",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list stage definitions", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"format": c.FlagSet.String("format", command.NilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
			}
		},
		ExecuteFunc:   stageDefinitionsListCmd,
		Endpoint:      configuration.ExtendedEndpoint,
		AdminEndpoint: configuration.DeveloperEndpoint,
	},
	{
		Description:  "Create a stage definition.",
		Subject:      "stage-definition",
		AltSubject:   "stagedef",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("create stage definition", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"label":       c.FlagSet.String("label", command.NilDefaultStr, "Stage Definitions's label"),
				"icon":        c.FlagSet.String("icon", command.NilDefaultStr, "Icon image file in data URI format like this: data:image/png;base64,iVBOR="),
				"title":       c.FlagSet.String("title", command.NilDefaultStr, "Stage Definitions's title"),
				"description": c.FlagSet.String("description", command.NilDefaultStr, "Stage Definitions's description"),
				"type":        c.FlagSet.String("type", command.NilDefaultStr, "Stage Definitions's type. Possible values: HTTPRequest, AnsibleBundle, WorkflowReference"),
				"vars":        c.FlagSet.String("vars", command.NilDefaultStr, "Stage Definitions's variables. These must be available in the execution context, otherwise the stage cannot run."),

				"ansible_bundle_filename": c.FlagSet.String("ansible_bundle_filename", command.NilDefaultStr, "Ansible bundle's file path to load the bundle from. Must be a zip file. Required when type=AnsibleBundle"),

				"http_request_url":                  c.FlagSet.String("http-request-url", command.NilDefaultStr, "HTTP Requests's URL. Required when using type=HTTPRequest"),
				"http_request_method":               c.FlagSet.String("http-request-method", command.NilDefaultStr, "HTTP Requests's method. Required when using type=HTTPRequest"),
				"http_request_body_filename":        c.FlagSet.String("http-request-body-filename", command.NilDefaultStr, "HTTP Requests's content is read from this file. Can only be used when type=HTTPRequest"),
				"http_request_body_from_pipe":       c.FlagSet.Bool("http-request-body-from-pipe", false, "HTTP Requests's content is read from stdin. Can only be used when type=HTTPRequest"),
				"http_request_header_accept":        c.FlagSet.String("http-request-header-accept", command.NilDefaultStr, "HTTP Requests's Accept header. Can only be used when type=HTTPRequest"),
				"http_request_header_authorization": c.FlagSet.String("http-request-header-authorization", command.NilDefaultStr, "HTTP Requests's Authorization header. Can only be used when type=HTTPRequest"),
				"http_request_header_cookie":        c.FlagSet.String("http-request-header-cookie", command.NilDefaultStr, "HTTP Requests's Cookie header. Can only be used when type=HTTPRequest"),
				"http_request_header_user_agent":    c.FlagSet.String("http-request-header-user-agent", command.NilDefaultStr, "HTTP Requests's User-Agent header. Can only be used when type=HTTPRequest"),
				"http_request_redirect":             c.FlagSet.String("http-request-redirect", command.NilDefaultStr, "HTTP Requests's method. Can only be used when type=HTTPRequest"),
				"http_request_follow":               c.FlagSet.Int("http-request-follow", command.NilDefaultInt, "HTTP Requests's follow. Can only be used when type=HTTPRequest"),
				"http_request_no_compress":          c.FlagSet.Bool("http-request-no-compress", false, "HTTP Requests's compress disabled if set. Can only be used when type=HTTPRequest"),
				"http_request_timeout":              c.FlagSet.Int("http-request-timeout", command.NilDefaultInt, "HTTP Requests's timeout. Can only be used when type=HTTPRequest"),
				"http_request_size":                 c.FlagSet.Int("http-request-size", command.NilDefaultInt, "HTTP Requests's size. Can only be used when type=HTTPRequest"),

				"workflow_id_or_label": c.FlagSet.String("workflow", command.NilDefaultStr, "workflow to reference. Can only be used when type=WorkflowReference"),

				"return_id": c.FlagSet.Bool("return-id", false, colors.Green("(Flag)")+" If set will print the ID of the created object. Useful for automating tasks."),
			}
		},
		ExecuteFunc:   stageDefinitionCreateCmd,
		Endpoint:      configuration.ExtendedEndpoint,
		AdminEndpoint: configuration.DeveloperEndpoint,
	},
	{
		Description:  "Delete a stage definition.",
		Subject:      "stage-definition",
		AltSubject:   "stagedef",
		Predicate:    "delete",
		AltPredicate: "rm",
		FlagSet:      flag.NewFlagSet("delete stage", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"stage_id_or_name": c.FlagSet.String("id", command.NilDefaultStr, "stage's id or name"),
				"autoconfirm":      c.FlagSet.Bool("autoconfirm", false, colors.Green("(Flag)")+" If set it will assume action is confirmed"),
			}
		},
		ExecuteFunc:   stageDefinitionDeleteCmd,
		Endpoint:      configuration.ExtendedEndpoint,
		AdminEndpoint: configuration.DeveloperEndpoint,
	},
	{
		Description:  "Add a stage to an infrastructure pre or post deploy workflows.",
		Subject:      "stage-definition",
		AltSubject:   "stagedef",
		Predicate:    "add-to-infrastructure",
		AltPredicate: "addtoinfra",
		FlagSet:      flag.NewFlagSet("add stage", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"stage_id_or_name":           c.FlagSet.String("id", command.NilDefaultStr, "stage's id or name"),
				"infrastructure_id_or_label": c.FlagSet.String("infra", command.NilDefaultStr, "The infrastructure's id"),
				"runlevel":                   c.FlagSet.Int("runlevel", command.NilDefaultInt, "The runlevel"),
				"group":                      c.FlagSet.String("group", command.NilDefaultStr, "When to run the stage"),
			}
		},
		ExecuteFunc:   stageDefinitionAddToInfrastructureCmd,
		Endpoint:      configuration.ExtendedEndpoint,
		AdminEndpoint: configuration.DeveloperEndpoint,
	},
	{
		Description:  "Add stage into workflow.",
		Subject:      "stage-definition",
		AltSubject:   "stagedef",
		Predicate:    "add-to-workflow",
		AltPredicate: "addtowf",
		FlagSet:      flag.NewFlagSet("add stage", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"stage_id_or_name":     c.FlagSet.String("id", command.NilDefaultStr, "stage's id or name"),
				"workflow_id_or_label": c.FlagSet.String("workflow", command.NilDefaultStr, "The workflow's id"),
				"runlevel":             c.FlagSet.Int("runlevel", command.NilDefaultInt, "The runlevel"),
			}
		},
		ExecuteFunc:   stageDefinitionAddToWorkflowCmd,
		Endpoint:      configuration.ExtendedEndpoint,
		AdminEndpoint: configuration.DeveloperEndpoint,
	},
}

func stageDefinitionsListCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	list, err := client.StageDefinitions()

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
			FieldSize: 6,
		},
		{
			FieldName: "TITLE",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "DESCRIPTION",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "TYPE",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "VARS_REQUIRED",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "DEF.",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "CREATED",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "UPDATED",
			FieldType: tableformatter.TypeString,
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
		case "SSHExec":
			ssh := s.StageDefinition.(metalcloud.SSHExec)
			stageDef = fmt.Sprintf("SSHExec command: %s", ssh.Command)
		case "WorkflowReference":
			wr := s.StageDefinition.(metalcloud.WorkflowReference)
			stageDef = fmt.Sprintf("Workflow Reference Workflow id: %d", wr.WorkflowID)
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

	tableformatter.TableSorter(schema).OrderBy(schema[1].FieldName).Sort(data)

	table := tableformatter.Table{
		Data:   data,
		Schema: schema,
	}
	return table.RenderTable("Stage Definitions", "", command.GetStringParam(c.Arguments["format"]))
}

func stageDefinitionCreateCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	stage := metalcloud.StageDefinition{}

	if v := c.Arguments["label"]; v != nil && *v.(*string) != command.NilDefaultStr {
		stage.StageDefinitionLabel = *v.(*string)
	} else {
		return "", fmt.Errorf("label is required")
	}

	if v := c.Arguments["icon"]; v != nil && *v.(*string) != command.NilDefaultStr {
		stage.IconAssetDataURI = *v.(*string)
	}

	if v := c.Arguments["title"]; v != nil && *v.(*string) != command.NilDefaultStr {
		stage.StageDefinitionTitle = *v.(*string)
	} else {
		return "", fmt.Errorf("title is required")
	}

	if v := c.Arguments["description"]; v != nil && *v.(*string) != command.NilDefaultStr {
		stage.StageDefinitionDescription = *v.(*string)
	}

	if v := c.Arguments["type"]; v != nil && *v.(*string) != command.NilDefaultStr {
		stage.StageDefinitionType = *v.(*string)
	} else {
		return "", fmt.Errorf("type is required")
	}

	if v := c.Arguments["vars"]; v != nil && *v.(*string) != command.NilDefaultStr {
		vars := *v.(*string)
		stage.StageDefinitionVariablesNamesRequired = strings.Split(vars, ",")
	}

	switch stage.StageDefinitionType {
	case "AnsibleBundle":
		if v := c.Arguments["ansible_bundle_filename"]; v != nil && *v.(*string) != command.NilDefaultStr {
			ab := metalcloud.AnsibleBundle{}

			ab.AnsibleBundleArchiveFilename = *v.(*string)

			content, err := configuration.ReadInputFromFile(ab.AnsibleBundleArchiveFilename)
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
		var err error
		if command.GetBoolParam(c.Arguments["http_request_body_from_pipe"]) {
			content, err = configuration.ReadInputFromPipe()
			if err != nil {
				return "", err
			}

		} else {
			if filename, ok := command.GetStringParamOk(c.Arguments["http_request_body_filename"]); ok {

				c, err := configuration.ReadInputFromFile(filename)
				if err != nil {
					return "", err
				}

				content = c
			}
		}

		req.Options.BodyBufferBase64 = base64.StdEncoding.EncodeToString(content)

		if v := c.Arguments["http_request_url"]; v != nil && *v.(*string) != command.NilDefaultStr {
			req.URL = *v.(*string)
		} else {
			return "", fmt.Errorf("http_request_url is required if using HTTPRequest")
		}

		if v := c.Arguments["http_request_method"]; v != nil && *v.(*string) != command.NilDefaultStr {
			req.Options.Method = *v.(*string)
		} else {
			return "", fmt.Errorf("http_request_method is required if using HTTPRequest")
		}

		if v := c.Arguments["http_request_redirect"]; v != nil && *v.(*string) != command.NilDefaultStr {
			req.Options.Redirect = *v.(*string)
		}

		if v := c.Arguments["http_request_follow"]; v != nil && *v.(*int) != command.NilDefaultInt {
			req.Options.Follow = *v.(*int)
		}

		if v := c.Arguments["http_request_no_compress"]; v != nil && *v.(*bool) {
			req.Options.Compress = *v.(*bool)
		}

		if v := c.Arguments["http_request_timeout"]; v != nil && *v.(*int) != command.NilDefaultInt {
			req.Options.Timeout = *v.(*int)
		}

		if v := c.Arguments["http_request_size"]; v != nil && *v.(*int) != command.NilDefaultInt {
			req.Options.Size = *v.(*int)
		}

		req.Options.Headers = metalcloud.WebFetchAPIRequestHeaders{}

		if v := c.Arguments["http_request_header_accept"]; v != nil && *v.(*string) != command.NilDefaultStr {
			req.Options.Headers.Accept = *v.(*string)
		}

		if v := c.Arguments["http_request_header_authorization"]; v != nil && *v.(*string) != command.NilDefaultStr {
			req.Options.Headers.Authorization = *v.(*string)
		}

		if v := c.Arguments["http_request_header_cookie"]; v != nil && *v.(*string) != command.NilDefaultStr {
			req.Options.Headers.Cookie = *v.(*string)
		}

		if v := c.Arguments["http_request_header_user_agent"]; v != nil && *v.(*string) != command.NilDefaultStr {
			req.Options.Headers.UserAgent = *v.(*string)
		}

		stage.StageDefinition = req

	case "WorkflowReference":

		wf, err := command.GetWorkflowFromCommand("workflow", c, client)
		if err != nil {
			return "", err
		}

		wr := metalcloud.WorkflowReference{
			WorkflowID: wf.WorkflowID,
			Type:       "WorkflowReference",
		}

		stage.StageDefinition = wr
	default:
		return "", fmt.Errorf("Unknown stage definition type %s", stage.StageDefinitionType)
	}

	ret, err := client.StageDefinitionCreate(stage)

	if command.GetBoolParam(c.Arguments["return_id"]) {
		return fmt.Sprintf("%d", ret.StageDefinitionID), nil
	}

	return "", err
}

func stageDefinitionDeleteCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	retS, err := getStageDefinitionFromCommand("id", c, client)
	if err != nil {
		return "", err
	}
	confirm := false

	if command.GetBoolParam(c.Arguments["autoconfirm"]) {
		confirm = true
	} else {

		confirmationMessage := fmt.Sprintf("Deleting stage definition %s (%d).  Are you sure? Type \"yes\" to continue:",
			retS.StageDefinitionLabel,
			retS.StageDefinitionID)

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

	err = client.StageDefinitionDelete(retS.StageDefinitionID)

	return "", err
}

func stageDefinitionAddToInfrastructureCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	stage, err := getStageDefinitionFromCommand("id", c, client)
	if err != nil {
		return "", err
	}

	infra, err := command.GetInfrastructureFromCommand("infra_id", c, client)
	if err != nil {
		return "", err
	}

	runlevel := 0
	if v := c.Arguments["runlevel"]; v != nil && *v.(*int) != command.NilDefaultInt {
		runlevel = *v.(*int)
	}

	runmoment := "post_deploy"
	if v := c.Arguments["group"]; v != nil && *v.(*string) != command.NilDefaultStr {
		runmoment = *v.(*string)
	}

	err = client.InfrastructureDeployCustomStageAddIntoRunlevel(infra.InfrastructureID, stage.StageDefinitionID, runlevel, runmoment)

	return "", err
}

func stageDefinitionAddToWorkflowCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	stage, err := getStageDefinitionFromCommand("id", c, client)
	if err != nil {
		return "", err
	}

	w, err := command.GetWorkflowFromCommand("workflow", c, client)
	if err != nil {
		return "", err
	}

	runlevel := command.GetIntParam(c.Arguments["runlevel"])

	stages, err := client.WorkflowStages(w.WorkflowID)

	for _, s := range *stages {
		if s.WorkflowStageRunLevel == runlevel {
			err = client.WorkflowStageAddIntoRunLevel(w.WorkflowID, stage.StageDefinitionID, runlevel)
			return "", err
		}
	}

	err = client.WorkflowStageAddAsNewRunLevel(w.WorkflowID, stage.StageDefinitionID, runlevel)

	return "", err
}

func getStageDefinitionFromCommand(paramName string, c *command.Command, client metalcloud.MetalCloudClient) (*metalcloud.StageDefinition, error) {

	v, err := command.GetParam(c, "stage_id_or_name", paramName)
	if err != nil {
		return nil, err
	}

	id, label, isID := command.IdOrLabel(v)

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

	if isID {
		return nil, fmt.Errorf("Stage definition %d not found", id)
	}

	return nil, fmt.Errorf("Stage definition %s not found", label)
}
