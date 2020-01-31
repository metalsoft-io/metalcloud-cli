package main

import (
	"flag"
	"fmt"
	"strings"

	interfaces "github.com/bigstepinc/metalcloud-cli/interfaces"
)

var serversCmds = []Command{

	Command{
		Description:  "Lists available servers",
		Subject:      "servers",
		AltSubject:   "srv",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list servers", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"format": c.FlagSet.String("format", _nilDefaultStr, "The output format. Supported values are 'json','csv'. The default format is human readable."),
				"filter": c.FlagSet.String("filter", "*", "filter to use when searching for servers. Check the documentation for examples. Defaults to '*'"),
			}
		},
		ExecuteFunc: serversListCmd,
	},

	Command{
		Description:  "Get server details",
		Subject:      "server",
		AltSubject:   "srv",
		Predicate:    "get",
		AltPredicate: "show",
		FlagSet:      flag.NewFlagSet("get server", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"id":               c.FlagSet.Int("id", _nilDefaultInt, "Server's ID"),
				"format":           c.FlagSet.String("format", _nilDefaultStr, "The output format. Supported values are 'json','csv'. The default format is human readable."),
				"show_credentials": c.FlagSet.Bool("show_credentials", false, "(Flag) If set returns the servers' IPMI credentials"),
			}
		},
		ExecuteFunc: serverGetCmd,
	},
}

func serversListCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	filter := *c.Arguments["filter"].(*string)
	if filter == _nilDefaultStr {
		filter = "*"
	}

	list, err := client.ServersSearch(filter)

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
			FieldName: "DATACENTER_NAME",
			FieldType: TypeString,
			FieldSize: 6,
		},
		SchemaField{
			FieldName: "SERVER_TYPE",
			FieldType: TypeString,
			FieldSize: 6,
		},
		SchemaField{
			FieldName: "STATUS",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "VENDOR",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "PRODUCT_NAME",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "SERIAL_NUMBER",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "CONFIG.",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "DISKS",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "TAGS",
			FieldType: TypeString,
			FieldSize: 4,
		},
		SchemaField{
			FieldName: "IPMI_HOST",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "ALLOCATED_TO.",
			FieldType: TypeString,
			FieldSize: 5,
		},
	}

	data := [][]interface{}{}
	for _, s := range *list {

		allocation := ""
		if s.ServerStatus == "used" || s.ServerStatus == "used_registering" {
			allocation = fmt.Sprintf("%s (#%d) IA:#%d Infra:#%d",
				s.InstanceLabel,
				s.InstanceID,
				s.InstanceArrayID,
				s.InfrastructureID)
		}
		productName := s.ServerProductName
		if len(s.ServerProductName) > 21 {
			productName = truncateString(s.ServerProductName, 18)
		}

		data = append(data, []interface{}{
			s.ServerID,
			s.DatacenterName,
			s.ServerTypeName,
			s.ServerStatus,
			s.ServerVendor,
			productName,
			s.ServerSerialNumber,
			fmt.Sprintf("%d GB RAM %d x %s (%d cores) ",
				s.ServerRAMGbytes,
				s.ServerProcessorCount,
				s.ServerProcessorName,
				s.ServerProcessorCoreCount),
			fmt.Sprintf("%d x %d GB [%s]",
				s.ServerDiskCount,
				s.ServerDiskSizeMbytes/1000,
				s.ServerDiskType),
			strings.Join(s.ServerTags, ","),
			s.ServerIPMIHost,
			allocation,
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
		sb.WriteString(fmt.Sprintf("Matching Servers for filter %s\n", filter))

		TableSorter(schema).OrderBy(
			schema[0].FieldName,
			schema[1].FieldName).Sort(data)

		AdjustFieldSizes(data, &schema)

		sb.WriteString(GetTableAsString(data, schema))

		sb.WriteString(fmt.Sprintf("Total: %d servers\n\n", len(*list)))
	}

	return sb.String(), nil
}

func serverGetCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	serverID := *c.Arguments["id"].(*int)
	if serverID == _nilDefaultInt {
		return "", fmt.Errorf("id is required")
	}

	showCredentials := false
	if c.Arguments["show_credentials"] != nil && *c.Arguments["show_credentials"].(*bool) {
		showCredentials = true
	}

	server, err := client.ServerGet(serverID, showCredentials)

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
			FieldName: "DATACENTER_NAME",
			FieldType: TypeString,
			FieldSize: 6,
		},
		SchemaField{
			FieldName: "SERVER_TYPE",
			FieldType: TypeString,
			FieldSize: 6,
		},
		SchemaField{
			FieldName: "STATUS",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "VENDOR",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "PRODUCT_NAME",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "SERIAL_NUMBER",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "CONFIG.",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "DISKS",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "TAGS",
			FieldType: TypeString,
			FieldSize: 4,
		},
		SchemaField{
			FieldName: "IPMI_HOST",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "ALLOCATED_TO.",
			FieldType: TypeString,
			FieldSize: 5,
		},
	}

	data := [][]interface{}{}

	serverType, err := client.ServerTypeGet(server.ServerTypeID)
	if err != nil {
		return "", err
	}

	allocation := ""

	if server.ServerStatus == "used" || server.ServerStatus == "used_registering" {
		searchRes, err := client.ServersSearch(fmt.Sprintf("+server_id:%d", server.ServerID))
		if err != nil {
			return "", err
		}

		if len(*searchRes) < 1 {
			return "", fmt.Errorf("Server not found by search function")
		}

		allocation = fmt.Sprintf("%s (#%d) IA:#%d Infra:#%d",
			(*searchRes)[0].InstanceLabel,
			(*searchRes)[0].InstanceID,
			(*searchRes)[0].InstanceArrayID,
			(*searchRes)[0].InfrastructureID)
	}

	productName := server.ServerProductName
	if len(server.ServerProductName) > 21 {
		productName = truncateString(server.ServerProductName, 18)
	}
	credentials := ""

	if showCredentials {

		schema = append(schema, SchemaField{
			FieldName: "CREDENTIALS",
			FieldType: TypeString,
			FieldSize: 5,
		})

		credentials = fmt.Sprintf("User: %s Pass: %s", server.ServerIPMInternalUsername, server.ServerIPMInternalPassword)

	}

	configuration := fmt.Sprintf("%d GB RAM %d x %s (%d cores) ",
		server.ServerRAMGbytes,
		server.ServerProcessorCount,
		server.ServerProcessorName,
		server.ServerProcessorCoreCount)

	disks := fmt.Sprintf("%d x %d GB [%s]",
		server.ServerDiskCount,
		server.ServerDiskSizeMbytes/1000,
		server.ServerDiskType)

	data = append(data, []interface{}{
		server.ServerID,
		server.DatacenterName,
		(*serverType).ServerTypeName,
		server.ServerStatus,
		server.ServerVendor,
		productName,
		server.ServerSerialNumber,
		configuration,
		disks,
		strings.Join(server.ServerTags, ","),
		server.ServerIPMIHost,
		allocation,
		credentials,
	})

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
		sb.WriteString("SERVER OVERVIEW\n")
		sb.WriteString("---------------\n")
		sb.WriteString(fmt.Sprintf("#%d %s %s\n",
			server.ServerID,
			serverType.ServerTypeName,
			server.DatacenterName))

		sb.WriteString(fmt.Sprintf("%s %s\n%s %s\n\n",
			server.ServerVendor,
			server.ServerProductName,
			server.ServerSerialNumber,
			server.ServerUUID))

		sb.WriteString("CONFIGURATION\n")
		sb.WriteString("------------\n")
		sb.WriteString(fmt.Sprintf("%s\n", configuration))
		sb.WriteString(fmt.Sprintf("%s\n\n", disks))

		sb.WriteString("ALLOCATION\n")
		sb.WriteString("----------\n")
		sb.WriteString(fmt.Sprintf("server_status: %s\nallocated to: %s\n\n", server.ServerStatus, allocation))

		if showCredentials {
			sb.WriteString("CREDENTIALS\n")
			sb.WriteString("-----------\n")
			sb.WriteString(fmt.Sprintf("%s\n", credentials))
		}

	}

	return sb.String(), nil
}
