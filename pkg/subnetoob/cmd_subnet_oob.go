package subnetoob

import (
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

var SubnetOOBCmds = []command.Command{
	{
		Description:  "Lists OOB subnet ",
		Subject:      "subnet-oob",
		AltSubject:   "oob-subnet",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list oob subnets", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"format":     c.FlagSet.String("format", command.NilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
				"filter":     c.FlagSet.String("filter", "*", "Filter to restrict the results. Defaults to '*'"),
				"datacenter": c.FlagSet.String("datacenter", command.NilDefaultStr, "Quick filter to restrict the results to show only the subnets of a datacenter."),
			}
		},
		ExecuteFunc: subnetOOBListCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
	{
		Description:  "Get a subnet OOB.",
		Subject:      "subnet-oob",
		AltSubject:   "oob-subnet",
		Predicate:    "get",
		AltPredicate: "show",
		FlagSet:      flag.NewFlagSet("Get an OOB subnet details", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"subnet_oob_id": c.FlagSet.Int("id", command.NilDefaultInt, colors.Red("(Required)")+" Subnet oob's id"),
				"format":        c.FlagSet.String("format", "", "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
				"raw":           c.FlagSet.Bool("raw", false, colors.Green("(Flag)")+" When set the return will be a full dump of the object. This is useful when copying configurations. Only works with json and yaml formats."),
			}
		},
		ExecuteFunc: subnetOOBGetCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
	{
		Description:  "Create an oob subnet.",
		Subject:      "subnet-oob",
		AltSubject:   "oob-subnet",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("Create an oob subnet", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"format":                c.FlagSet.String("format", "json", "The input format. Supported values are 'json','yaml'. The default format is json."),
				"read_config_from_file": c.FlagSet.String("f", command.NilDefaultStr, colors.Red("(Required)")+" Read configuration from file"),
				"read_config_from_pipe": c.FlagSet.Bool("pipe", false, colors.Green("(Flag)")+" If set, read configuration from pipe instead of from a file. Either this flag or the -raw-config option must be used."),
				"return_id":             c.FlagSet.Bool("return-id", false, "Will print the ID of the created object. Useful for automating tasks."),
			}
		},
		ExecuteFunc: subnetOOBCreateCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
	{
		Description:  "Delete an OOB subnet.",
		Subject:      "subnet-oob",
		AltSubject:   "oob-subnet",
		Predicate:    "delete",
		AltPredicate: "rm",
		FlagSet:      flag.NewFlagSet("delete an oob subnet", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"subnet_oob_id": c.FlagSet.Int("id", command.NilDefaultInt, colors.Red("(Required)")+" Subnet's's id"),
				"autoconfirm":   c.FlagSet.Bool("autoconfirm", false, colors.Green("(Flag)")+" If set it will assume action is confirmed"),
			}
		},
		ExecuteFunc: subnetOOBDeleteCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
}

func subnetOOBListCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	filter := command.GetStringParam(c.Arguments["filter"])
	if datacenter, ok := command.GetStringParamOk(c.Arguments["datacenter"]); ok {
		filter = fmt.Sprintf("datacenter_name: %s %s", datacenter, filter)
	}

	list, err := client.SubnetOOBSearch(filter)

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
			FieldName: "DATACENTER",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "PREFIX",
			FieldType: tableformatter.TypeInt,
			FieldSize: 2,
		},
		{
			FieldName: "FOR",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "AUTOMATIC_ALLOC_ENABLED",
			FieldType: tableformatter.TypeBool,
			FieldSize: 8,
		},
		{
			FieldName: "RANGE_START",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "RANGE_END",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
	}

	data := [][]interface{}{}
	for _, s := range *list {

		data = append(data, []interface{}{

			s.SubnetOOBID,
			s.SubnetOOBLabel,
			s.DatacenterName,
			s.SubnetOOBPrefixSize,
			s.SubnetOOBAllocateForResourceType,
			s.SubnetOOBUseForAutoAllocation,
			s.SubnetOOBRangeStartHumanReadable,
			s.SubnetOOBRangeEndHumanReadable,
		})

	}

	tableformatter.TableSorter(schema).OrderBy(
		schema[0].FieldName,
		schema[1].FieldName,
		schema[2].FieldName).Sort(data)

	table := tableformatter.Table{
		Data:   data,
		Schema: schema,
	}
	return table.RenderTable("OOB Subnets", "", command.GetStringParam(c.Arguments["format"]))
}

func subnetOOBGetCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	id, ok := command.GetIntParamOk(c.Arguments["subnet_oob_id"])
	if !ok {
		return "", fmt.Errorf("-id is required")
	}

	s, err := client.SubnetOOBGet(id)
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
			FieldName: "DATACENTER",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "PREFIX",
			FieldType: tableformatter.TypeInt,
			FieldSize: 2,
		},
		{
			FieldName: "FOR",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "AUTOMATIC_ALLOC_ENABLED",
			FieldType: tableformatter.TypeBool,
			FieldSize: 8,
		},
		{
			FieldName: "RANGE_START",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "RANGE_END",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
	}

	data := [][]interface{}{{
		s.SubnetOOBID,
		s.SubnetOOBLabel,
		s.DatacenterName,
		s.SubnetOOBPrefixSize,
		s.SubnetOOBAllocateForResourceType,
		s.SubnetOOBUseForAutoAllocation,
		s.SubnetOOBRangeStartHumanReadable,
		s.SubnetOOBRangeEndHumanReadable,
	}}

	var sb strings.Builder

	format := command.GetStringParam(c.Arguments["format"])

	if command.GetBoolParam(c.Arguments["raw"]) {
		ret, err := objects.RenderRawObject(*s, format, "SubnetOOB")
		if err != nil {
			return "", err
		}
		sb.WriteString(ret)
	} else {
		table := tableformatter.Table{
			Data:   data,
			Schema: schema,
		}
		ret, err := table.RenderTransposedTable("subnet oob", "", format)
		if err != nil {
			return "", err
		}
		sb.WriteString(ret)
	}

	return sb.String(), nil
}

func subnetOOBCreateCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	var sn metalcloud.SubnetOOB

	err := command.GetRawObjectFromCommand(c, &sn)
	if err != nil {
		return "", err
	}

	ret, err := client.SubnetOOBCreate(sn)
	if err != nil {
		return "", err
	}

	if command.GetBoolParam(c.Arguments["return_id"]) {
		return fmt.Sprintf("%d", ret.SubnetOOBID), nil
	}

	return "", err
}

func subnetOOBDeleteCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	id, ok := command.GetIntParamOk(c.Arguments["subnet_oob_id"])
	if !ok {
		return "", fmt.Errorf("-id is required")
	}
	confirm := false

	obj, err := client.SubnetOOBGet(id)
	if err != nil {
		return "", err
	}

	if command.GetBoolParam(c.Arguments["autoconfirm"]) {
		confirm = true
	} else {

		confirmationMessage := fmt.Sprintf("Deleting oob subnet #%d (%s-%s).  Are you sure? Type \"yes\" to continue:",
			obj.SubnetOOBID,
			obj.SubnetOOBRangeStartHumanReadable,
			obj.SubnetOOBRangeEndHumanReadable)

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

	err = client.SubnetOOBDelete(obj.SubnetOOBID)

	return "", err
}
