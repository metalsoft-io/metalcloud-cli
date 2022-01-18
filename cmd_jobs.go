package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
	"github.com/metalsoft-io/tableformatter"
)

var jobsCmds = []Command{

	{
		Description:  "Lists all jobs.",
		Subject:      "job",
		AltSubject:   "jobs",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list jobs", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"format": c.FlagSet.String("format", _nilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
				"filter": c.FlagSet.String("filter", "*", "filter to use when searching for jobs. Check the documentation for examples. Defaults to '*'"),
				"limit":  c.FlagSet.Int("limit", 20, "how many jobs to show. Latest jobs first."),
				"watch":  c.FlagSet.String("watch", _nilDefaultStr, "If set to a human readable interval such as '4s', '1m' will print the job status until interrupted."),
			}
		},
		ExecuteFunc: jobListCmdWithWatch,
		Endpoint:    DeveloperEndpoint,
	},
	{
		Description:  "Get job details.",
		Subject:      "job",
		AltSubject:   "afc",
		Predicate:    "get",
		AltPredicate: "show",
		FlagSet:      flag.NewFlagSet("get job", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"job_id": c.FlagSet.String("id", _nilDefaultStr, "JOB ID"),
				"format": c.FlagSet.String("format", _nilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
				"watch":  c.FlagSet.String("watch", _nilDefaultStr, "If set to a human readable interval such as '4s', '1m' will print the job status until interrupted."),
			}
		},
		ExecuteFunc: jobGetCmdWithWatch,
		Endpoint:    DeveloperEndpoint,
	},
	{
		Description:  "Retry job.",
		Subject:      "job",
		AltSubject:   "afc",
		Predicate:    "retry",
		AltPredicate: "ret",
		FlagSet:      flag.NewFlagSet("retry job", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"job_id":      c.FlagSet.String("id", _nilDefaultStr, "JOB ID"),
				"autoconfirm": c.FlagSet.Bool("autoconfirm", false, "If true it does not ask for confirmation anymore"),
			}
		},
		ExecuteFunc: jobRetryCmd,
		Endpoint:    DeveloperEndpoint,
	},
	{
		Description:  "Skip job.",
		Subject:      "job",
		AltSubject:   "afc",
		Predicate:    "skip",
		AltPredicate: "skip",
		FlagSet:      flag.NewFlagSet("Skip job", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"job_id":      c.FlagSet.String("id", _nilDefaultStr, "JOB ID"),
				"autoconfirm": c.FlagSet.Bool("autoconfirm", false, "If true it does not ask for confirmation anymore"),
			}
		},
		ExecuteFunc: jobSkipCmd,
		Endpoint:    DeveloperEndpoint,
	},
	{
		Description:  "Delete job.",
		Subject:      "job",
		AltSubject:   "afc",
		Predicate:    "delete",
		AltPredicate: "rm",
		FlagSet:      flag.NewFlagSet("Delete job", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"job_id":      c.FlagSet.String("id", _nilDefaultStr, "JOB ID"),
				"autoconfirm": c.FlagSet.Bool("autoconfirm", false, "If true it does not ask for confirmation anymore"),
			}
		},
		ExecuteFunc: jobDeleteCmd,
		Endpoint:    DeveloperEndpoint,
	},
	{
		Description:  "Kill job.",
		Subject:      "job",
		AltSubject:   "afc",
		Predicate:    "kill",
		AltPredicate: "kill",
		FlagSet:      flag.NewFlagSet("Kill job", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"job_id":      c.FlagSet.String("id", _nilDefaultStr, "JOB ID"),
				"mark":        c.FlagSet.String("mark", "kill", "One of 'kill','stop_retrying','kill_and_stop_retrying','kill_and_stop_retrying','keep_alive'"),
				"autoconfirm": c.FlagSet.Bool("autoconfirm", false, "If true it does not ask for confirmation anymore"),
			}
		},
		ExecuteFunc: jobKillCmd,
		Endpoint:    DeveloperEndpoint,
	},
}

func jobListCmdWithWatch(c *Command, client metalcloud.MetalCloudClient) (string, error) {
	return funcWithWatch(c, client, jobListCmd)
}

func jobListCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	filter := getStringParam(c.Arguments["filter"])
	limit := getIntParam(c.Arguments["limit"])

	list, err := client.AFCSearch(filter, 0, limit)
	if err != nil {
		return "", err
	}

	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	schema := []tableformatter.SchemaField{
		{
			FieldName: "ID",
			FieldType: tableformatter.TypeInt,
			FieldSize: 6,
		},
		{
			FieldName: "STATUS",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},

		{
			FieldName: "AFFECTS",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "RETRIES",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "REQUEST",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "RESPONSE",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
	}

	data := [][]interface{}{}

	statusCounts := map[string]int{
		"thrown_error":                0,
		"thrown_error_while_retrying": 0,
		"returned_success":            0,
		"running":                     0,
	}

	for _, s := range *list {

		statusCounts[s.AFCStatus] = statusCounts[s.AFCStatus] + 1

		status := s.AFCStatus

		switch status {
		case "thrown_error":
			status = red(status)
		case "thrown_error_while_retrying":
			status = magenta(status)
		case "running":
			status = yellow(status)
		case "returned_success":
			status = green(status)
		default:
			status = yellow(status)
		}

		affects := ""
		if s.ServerID != 0 {
			affects = affects + fmt.Sprintf("Server: #%d ", s.ServerID)
		}
		if s.InstanceID != 0 {
			affects = affects + fmt.Sprintf("Inst: #%d ", s.InstanceID)
		}
		if s.InfrastructureID != 0 {
			affects = affects + fmt.Sprintf("Infra: #%d ", s.InfrastructureID)
		}

		if s.InfrastructureID == 0 && s.AFCGroupID != 0 {
			affects = affects + fmt.Sprintf("Group: #%d ", s.AFCGroupID)
		}

		retries := fmt.Sprintf("%d/%d", s.AFCRetryCount, s.AFCRetryMax)
		if s.AFCRetryCount >= s.AFCRetryMax {
			retries = red(retries)
		} else if s.AFCRetryCount < s.AFCRetryMax && s.AFCRetryCount > 1 {
			retries = yellow(retries)
		} else {
			retries = green(retries)
		}

		request := ""
		switch s.AFCFunctionName {
		case "infrastructure_provision":
			var paramsArr []interface{}
			err := json.Unmarshal([]byte(s.AFCParamsJSON), &paramsArr)
			if err != nil {
				return "", err
			}
			funcName := ""
			if len(paramsArr) >= 2 {
				funcName = paramsArr[1].(string)
			}
			var actualParams interface{}
			if len(paramsArr) >= 3 {
				actualParams = paramsArr[2]
			}
			request = fmt.Sprintf("%s(%+v)", funcName, actualParams)
		default:
			request = fmt.Sprintf("%s(%+v)", s.AFCFunctionName, s.AFCParamsJSON)
		}

		if len(request) > 40 {
			request = truncateString(request, 40)
		}
		var respObj map[string]string

		response := ""
		if s.AFCExceptionJSON != "" {
			err := json.Unmarshal([]byte(s.AFCExceptionJSON), &respObj)
			if err != nil {
				return "", err
			}

			response = fmt.Sprintf("%+v", respObj["message"])
		}

		if len(response) > 80 {
			response = truncateString(response, 80)
		}

		row := []interface{}{
			s.AFCID,
			status,
			affects,
			retries,
			request,
			response,
		}

		data = append(data, row)

	}

	table := tableformatter.Table{
		Data:   data,
		Schema: schema,
	}

	title := fmt.Sprintf("Jobs: %d thrown error %d thrown error retrying  %d running %d returned success",
		statusCounts["thrown_error"],
		statusCounts["thrown_error_while_retrying"],
		statusCounts["running"],
		statusCounts["returned_success"],
	)

	return table.RenderTable(title, "", getStringParam(c.Arguments["format"]))

}

func jobGetCmdWithWatch(c *Command, client metalcloud.MetalCloudClient) (string, error) {
	return funcWithWatch(c, client, jobGetCmd)
}

func jobGetCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	afc_id_s, ok := getStringParamOk(c.Arguments["job_id"])
	if !ok {
		return "", fmt.Errorf("-id required")
	}

	afc_id, err := strconv.Atoi(afc_id_s)
	if err != nil {
		return "", err
	}

	s, err := client.AFCGet(afc_id)
	if err != nil {
		return "", err
	}

	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	schema := []tableformatter.SchemaField{
		{
			FieldName: "ID",
			FieldType: tableformatter.TypeInt,
			FieldSize: 6,
		},
		{
			FieldName: "STATUS",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},

		{
			FieldName: "AFFECTS",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "RETRIES",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "REQUEST",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "RESPONSE",
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
			FieldSize: 5,
		},
	}

	data := [][]interface{}{}

	status := s.AFCStatus

	switch status {
	case "thrown_error":
		status = red(status)
	case "thrown_error_while_retrying":
		status = magenta(status)
	case "running":
		status = yellow(status)
	case "returned_success":
		status = green(status)
	default:
		status = yellow(status)
	}

	affects := ""
	if s.ServerID != 0 {
		affects = affects + fmt.Sprintf("Server: #%d ", s.ServerID)
	}
	if s.AFCGroupID != 0 {
		affects = affects + fmt.Sprintf("Group: #%d ", s.AFCGroupID)
	}
	if s.InstanceID != 0 {
		affects = affects + fmt.Sprintf("Instance: #%d ", s.InstanceID)
	}
	if s.InfrastructureID != 0 {
		affects = affects + fmt.Sprintf("Infrastructure: #%d ", s.InfrastructureID)
	}

	retries := fmt.Sprintf("%d/%d", s.AFCRetryCount, s.AFCRetryMax)
	if s.AFCRetryCount >= s.AFCRetryMax {
		retries = red(retries)
	} else if s.AFCRetryCount < s.AFCRetryMax && s.AFCRetryCount > 1 {
		retries = yellow(retries)
	} else {
		retries = green(retries)
	}

	request := ""
	switch s.AFCFunctionName {
	case "infrastructure_provision":
		var paramsArr []interface{}
		err := json.Unmarshal([]byte(s.AFCParamsJSON), &paramsArr)
		if err != nil {
			return "", err
		}
		funcName := ""
		if len(paramsArr) >= 2 {
			funcName = paramsArr[1].(string)
		}
		var actualParams interface{}
		if len(paramsArr) >= 3 {
			actualParams = paramsArr[2]
		}
		request = fmt.Sprintf("%s(%+v)", funcName, actualParams)
	default:
		request = fmt.Sprintf("%s(%+v)", s.AFCFunctionName, s.AFCParamsJSON)
	}

	if len(request) > 100 {
		request = truncateString(request, 100)
	}
	var respObj map[string]string

	response := ""
	if s.AFCExceptionJSON != "" {
		err := json.Unmarshal([]byte(s.AFCExceptionJSON), &respObj)
		if err != nil {
			return "", err
		}

		response = fmt.Sprintf("%+v", respObj["message"])
	}

	if len(response) > 100 {
		response = wrapToLength(response, 100)
	}

	row := []interface{}{
		s.AFCID,
		status,
		affects,
		retries,
		request,
		response,
		s.AFCCreatedTimestamp,
		s.AFCUpdatedTimestamp,
	}

	data = append(data, row)

	table := tableformatter.Table{
		Data:   data,
		Schema: schema,
	}

	title := fmt.Sprintf("Job details")

	interval, ok := getStringParamOk(c.Arguments["watch"])
	if ok {

		watch(func() (string, error) {
			return table.RenderTransposedTable(title, "", getStringParam(c.Arguments["format"]))
		},
			interval)
	}

	return table.RenderTransposedTable(title, "", getStringParam(c.Arguments["format"]))

}

func jobRetryCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	afc_id_s, ok := getStringParamOk(c.Arguments["job_id"])
	if !ok {
		return "", fmt.Errorf("-id required")
	}

	afc_id, err := strconv.Atoi(afc_id_s)
	if err != nil {
		return "", err
	}

	confirm, err := confirmCommand(c, func() string {

		confirmationMessage := fmt.Sprintf("Retrying Job #%d.  Are you sure? Type \"yes\" to continue:",
			afc_id,
		)

		if strings.HasSuffix(os.Args[0], ".test") {
			confirmationMessage = ""
		}

		return confirmationMessage
	})

	if err != nil {
		return "", err
	}

	if confirm {
		err = client.AFCRetryCall(afc_id)
	}

	return jobGetCmd(c, client)
}

func jobSkipCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	afc_id_s, ok := getStringParamOk(c.Arguments["job_id"])
	if !ok {
		return "", fmt.Errorf("-id required")
	}

	afc_id, err := strconv.Atoi(afc_id_s)
	if err != nil {
		return "", err
	}

	confirm, err := confirmCommand(c, func() string {

		confirmationMessage := fmt.Sprintf("Skipping Job #%d.  Are you sure? Type \"yes\" to continue:",
			afc_id,
		)

		if strings.HasSuffix(os.Args[0], ".test") {
			confirmationMessage = ""
		}

		return confirmationMessage
	})

	if err != nil {
		return "", err
	}

	if confirm {
		err = client.AFCSkip(afc_id)
	}

	return jobGetCmd(c, client)
}

func jobDeleteCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	afc_id_s, ok := getStringParamOk(c.Arguments["job_id"])
	if !ok {
		return "", fmt.Errorf("-id required")
	}

	afc_id, err := strconv.Atoi(afc_id_s)
	if err != nil {
		return "", err
	}

	confirm, err := confirmCommand(c, func() string {

		confirmationMessage := fmt.Sprintf("Deleting Job #%d.  Are you sure? Type \"yes\" to continue:",
			afc_id,
		)

		if strings.HasSuffix(os.Args[0], ".test") {
			confirmationMessage = ""
		}

		return confirmationMessage
	})

	if err != nil {
		return "", err
	}

	if confirm {
		err = client.AFCDelete(afc_id)
	}

	return jobGetCmd(c, client)
}

func jobKillCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	afc_id_s, ok := getStringParamOk(c.Arguments["job_id"])
	if !ok {
		return "", fmt.Errorf("-id required")
	}

	afc_id, err := strconv.Atoi(afc_id_s)
	if err != nil {
		return "", err
	}

	mark := getStringParam(c.Arguments["mark"])

	confirm, err := confirmCommand(c, func() string {

		confirmationMessage := fmt.Sprintf("Killing Job #%d.  Are you sure? Type \"yes\" to continue:",
			afc_id,
		)

		if strings.HasSuffix(os.Args[0], ".test") {
			confirmationMessage = ""
		}

		return confirmationMessage
	})

	if err != nil {
		return "", err
	}

	if confirm {
		err = client.AFCMarkForDeath(afc_id, mark)
	}

	return jobGetCmd(c, client)
}
