package main

import (
	"flag"
	"fmt"
	"strings"

	interfaces "github.com/bigstepinc/metalcloud-cli/interfaces"
)

var serversCmds = []Command{

	{
		Description:  "Lists available servers",
		Subject:      "server",
		AltSubject:   "srv",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list servers", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"format":           c.FlagSet.String("format", _nilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
				"filter":           c.FlagSet.String("filter", "*", "filter to use when searching for servers. Check the documentation for examples. Defaults to '*'"),
				"show_credentials": c.FlagSet.Bool("show-credentials", false, "(Flag) If set returns the servers' IPMI credentials. (Slow for large queries)"),
			}
		},
		ExecuteFunc: serversListCmd,
		Endpoint:    DeveloperEndpoint,
	},

	{
		Description:  "Get server details",
		Subject:      "server",
		AltSubject:   "srv",
		Predicate:    "get",
		AltPredicate: "show",
		FlagSet:      flag.NewFlagSet("get server", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"id":               c.FlagSet.Int("id", _nilDefaultInt, "Server's ID"),
				"format":           c.FlagSet.String("format", _nilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
				"show_credentials": c.FlagSet.Bool("show-credentials", false, "(Flag) If set returns the servers' IPMI credentials"),
			}
		},
		ExecuteFunc: serverGetCmd,
		Endpoint:    DeveloperEndpoint,
	},
}

func serversListCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	filter := getStringParam(c.Arguments["filter"])

	list, err := client.ServersSearch(filter)

	if err != nil {
		return "", err
	}

	schema := []SchemaField{
		{
			FieldName: "ID",
			FieldType: TypeInt,
			FieldSize: 6,
		},
		{
			FieldName: "DATACENTER_NAME",
			FieldType: TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "SERVER_TYPE",
			FieldType: TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "STATUS",
			FieldType: TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "VENDOR",
			FieldType: TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "PRODUCT_NAME",
			FieldType: TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "SERIAL_NUMBER",
			FieldType: TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "CONFIG.",
			FieldType: TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "TAGS",
			FieldType: TypeString,
			FieldSize: 4,
		},
		{
			FieldName: "IPMI_HOST",
			FieldType: TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "ALLOCATED_TO.",
			FieldType: TypeString,
			FieldSize: 5,
		},
	}

	showCredentials := false
	if c.Arguments["show_credentials"] != nil && *c.Arguments["show_credentials"].(*bool) {
		showCredentials = true

		schema = append(schema, SchemaField{
			FieldName: "IPMI_USER",
			FieldType: TypeString,
			FieldSize: 5,
		})

		schema = append(schema, SchemaField{
			FieldName: "IPMI_PASS",
			FieldType: TypeString,
			FieldSize: 5,
		})

	}

	data := [][]interface{}{}
	for _, s := range *list {

		allocation := ""
		if s.ServerStatus == "used" || s.ServerStatus == "used_registering" {
			users := strings.Join(s.UserEmail[0], ",")
			if len(users) > 30 {
				users = truncateString(users, 27)
			}
			allocation = fmt.Sprintf("%s %s (#%d) IA:#%d Infra:#%d",
				users,
				s.InstanceLabel[0],
				s.InstanceID[0],
				s.InstanceArrayID[0],
				s.InfrastructureID[0])
		}
		productName := s.ServerProductName
		if len(s.ServerProductName) > 21 {
			productName = truncateString(s.ServerProductName, 18)
		}
		diskDescription := ""
		if s.ServerDiskCount > 0 {
			diskDescription = fmt.Sprintf(" %d x %d GB %s", s.ServerDiskCount,
				s.ServerDiskSizeMbytes/1000,
				s.ServerDiskType)
		}

		credentialsUser := ""
		credentialsPass := ""

		if showCredentials {

			server, err := client.ServerGet(s.ServerID, showCredentials)

			if err != nil {
				return "", err
			}

			credentialsUser = fmt.Sprintf("%s", server.ServerIPMInternalUsername)
			credentialsPass = fmt.Sprintf("%s", server.ServerIPMInternalPassword)

		}
		data = append(data, []interface{}{
			s.ServerID,
			s.DatacenterName,
			s.ServerTypeName,
			s.ServerStatus,
			s.ServerVendor,
			productName,
			s.ServerSerialNumber,
			fmt.Sprintf("%d GB RAM %d x %s (%d cores)%s",
				s.ServerRAMGbytes,
				s.ServerProcessorCount,
				s.ServerProcessorName,
				s.ServerProcessorCoreCount,
				diskDescription,
			),
			strings.Join(s.ServerTags, ","),
			s.ServerIPMIHost,
			allocation,
			credentialsUser,
			credentialsPass,
		})

	}

	return renderTable("Servers", "", getStringParam(c.Arguments["format"]), data, schema)
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
		{
			FieldName: "ID",
			FieldType: TypeInt,
			FieldSize: 6,
		},
		{
			FieldName: "DATACENTER_NAME",
			FieldType: TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "SERVER_TYPE",
			FieldType: TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "STATUS",
			FieldType: TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "VENDOR",
			FieldType: TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "PRODUCT_NAME",
			FieldType: TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "SERIAL_NUMBER",
			FieldType: TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "CONFIG.",
			FieldType: TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "DISKS",
			FieldType: TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "TAGS",
			FieldType: TypeString,
			FieldSize: 4,
		},
		{
			FieldName: "IPMI_HOST",
			FieldType: TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "ALLOCATED_TO.",
			FieldType: TypeString,
			FieldSize: 5,
		},
	}

	data := [][]interface{}{}

	serverTypeName := "<no_server_type>"
	if server.ServerTypeID != 0 {
		serverType, err := client.ServerTypeGet(server.ServerTypeID)
		if err != nil {
			return "", err
		}
		serverTypeName = serverType.ServerTypeDisplayName
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
		serverTypeName,
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
			serverTypeName,
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
