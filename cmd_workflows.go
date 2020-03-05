package main

import (
	"flag"
	"fmt"
	"strings"

	metalcloud "github.com/bigstepinc/metal-cloud-sdk-go"
	interfaces "github.com/bigstepinc/metalcloud-cli/interfaces"
)

var workflowCmds = []Command{

	Command{
		Description:  "Lists available workflows",
		Subject:      "workflow",
		AltSubject:   "workflows",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list workflows", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"usage":  c.FlagSet.String("usage", _nilDefaultStr, "Workflow usage. One of infrastructure, network_equipment, server, free_standing, storage_pool, user, os_template"),
				"format": c.FlagSet.String("format", _nilDefaultStr, "The output format. Supported values are 'json','csv'. The default format is human readable."),
			}
		},
		ExecuteFunc: workflowsListCmd,
		Endpoint:    DeveloperEndpoint,
	},
	Command{
		Description:  "Lists available workflows",
		Subject:      "workflow",
		AltSubject:   "workflows",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list workflows", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"usage":  c.FlagSet.String("usage", _nilDefaultStr, "Workflow usage. One of infrastructure, network_equipment, server, free_standing, storage_pool, user, os_template."),
				"format": c.FlagSet.String("format", _nilDefaultStr, "The output format. Supported values are 'json','csv'. The default format is human readable."),
			}
		},
		ExecuteFunc: workflowsListCmd,
		Endpoint:    DeveloperEndpoint,
	},
	Command{
		Description:  "Get workflow details",
		Subject:      "workflow",
		AltSubject:   "wf",
		Predicate:    "get",
		AltPredicate: "show",
		FlagSet:      flag.NewFlagSet("list workflows", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"workflow_id_or_label": c.FlagSet.String("id", _nilDefaultStr, "Workflow's id or label."),
				"format":               c.FlagSet.String("format", _nilDefaultStr, "The output format. Supported values are 'json','csv'. The default format is human readable."),
			}
		},
		ExecuteFunc: workflowGetCmd,
		Endpoint:    DeveloperEndpoint,
	},
	Command{
		Description:  "Create workflow ",
		Subject:      "workflow",
		AltSubject:   "wf",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("create workflow", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"label":               c.FlagSet.String("label", _nilDefaultStr, "Workflow's label."),
				"title":               c.FlagSet.String("title", _nilDefaultStr, "Workflow's title."),
				"usage":               c.FlagSet.String("usage", _nilDefaultStr, "Workflow's usage, one of:  infrastructure, network_equipment, server, free_standing, storage_pool, user, os_template."),
				"description":         c.FlagSet.String("description", _nilDefaultStr, "Workflow's description"),
				"deprecated":          c.FlagSet.Bool("deprecated", false, "Flag. Workflow's deprecation status. Default false"),
				"icon_asset_data_uri": c.FlagSet.String("icon", _nilDefaultStr, "Workflow's icon data"),
				"return_id":           c.FlagSet.Bool("return_id", false, "(Flag) If set will print the ID of the created workflow. Useful for automating tasks."),
			}
		},
		ExecuteFunc: workflowCreateCmd,
		Endpoint:    DeveloperEndpoint,
	},
}

func workflowsListCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	usage := getStringParam(c.Arguments["usage"])

	list, err := client.WorkflowsWithUsage(usage)

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
			FieldName: "USAGE",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "DESCRIPTION",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "TITLE",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "OWNER",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "DEPRECATED",
			FieldType: TypeBool,
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
	for _, w := range *list {

		user, err := client.UserGet(w.UserIDOwner)
		if err != nil {
			return "", err
		}

		data = append(data, []interface{}{
			w.WorkflowID,
			w.WorkflowLabel,
			w.WorkflowUsage,
			w.WorkflowDescription,
			w.WorkflowTitle,
			user.UserEmail,
			w.WorkflowIsDeprecated,
			w.WorkflowCreatedTimestamp,
			w.WorkflowUpdatedTimestamp,
		})

	}

	return renderTable("Workflows", "", getStringParam(c.Arguments["format"]), data, schema)
}

func workflowGetCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	wf, err := getWorkflowFromCommand("id", c, client)
	if err != nil {
		return "", err
	}

	schema := []SchemaField{
		SchemaField{
			FieldName: "RUNLEVEL",
			FieldType: TypeInt,
			FieldSize: 6,
		},
		SchemaField{
			FieldName: "STAGES",
			FieldType: TypeString,
			FieldSize: 6,
		},
	}

	list, err := client.WorkflowStages(wf.WorkflowID)
	if err != nil {
		return "", err
	}

	runlevels := map[int][]*metalcloud.StageDefinition{}

	for _, s := range *list {
		stageDef, err := client.StageDefinitionGet(s.StageDefinitionID)
		if err != nil {
			return "", err
		}

		runlevels[s.WorkflowStageRunLevel] = append(runlevels[s.WorkflowStageRunLevel], stageDef)
	}

	data := [][]interface{}{}
	for k, stageDefs := range runlevels {

		var stageDescriptions []string

		for _, s := range stageDefs {
			stageDescription := fmt.Sprintf("%s(%d)", s.StageDefinitionTitle, s.StageDefinitionID)
			stageDescriptions = append(stageDescriptions, stageDescription)
		}

		data = append(data, []interface{}{
			k,
			strings.Join(stageDescriptions, " "),
		})

	}

	topLine := fmt.Sprintf("Workflow %s (%d) has the following stages:", wf.WorkflowLabel, wf.WorkflowID)
	return renderTable("Stages", topLine, getStringParam(c.Arguments["format"]), data, schema)
}

func workflowCreateCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	label, ok := getStringParamOk(c.Arguments["label"])
	if !ok {
		return "", fmt.Errorf("-label is required")
	}

	usage, ok := getStringParamOk(c.Arguments["usage"])
	if !ok {
		return "", fmt.Errorf("-usage is required. It must be one of infrastructure, network_equipment, server, free_standing, storage_pool, user, os_template")
	}

	wf := metalcloud.Workflow{
		WorkflowLabel:        label,
		WorkflowTitle:        getStringParam(c.Arguments["title"]),
		WorkflowUsage:        usage,
		WorkflowDescription:  getStringParam(c.Arguments["description"]),
		WorkflowIsDeprecated: getBoolParam(c.Arguments["deprecated"]),
		IconAssetDataURI:     getStringParam(c.Arguments["icon"]),
	}

	ret, err := client.WorkflowCreate(wf)
	if err != nil {
		return "", err
	}
	if getBoolParam(c.Arguments["return_id"]) {
		return fmt.Sprintf("%d", ret.WorkflowID), nil
	}

	return "", nil

}

func getWorkflowFromCommand(paramName string, c *Command, client interfaces.MetalCloudClient) (*metalcloud.Workflow, error) {

	v, err := getParam(c, "workflow_id_or_label", paramName)
	if err != nil {
		return nil, err
	}

	id, label, isID := idOrLabel(v)

	if isID {
		return client.WorkflowGet(id)
	}

	list, err := client.Workflows()
	if err != nil {
		return nil, err
	}

	for _, s := range *list {
		if s.WorkflowLabel == label {
			return &s, nil
		}
	}

	if isID {
		return nil, fmt.Errorf("workflow %d not found", id)
	}

	return nil, fmt.Errorf("workflow %s not found", label)

}