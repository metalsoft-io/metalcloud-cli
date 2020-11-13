package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	metalcloud "github.com/bigstepinc/metal-cloud-sdk-go"
	interfaces "github.com/bigstepinc/metalcloud-cli/interfaces"
	"github.com/metalsoft-io/tableformatter"
)

var serversCmds = []Command{

	{
		Description:  "Lists all servers.",
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
		Description:  "Get server details.",
		Subject:      "server",
		AltSubject:   "srv",
		Predicate:    "get",
		AltPredicate: "show",
		FlagSet:      flag.NewFlagSet("get server", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"server_id_or_uuid": c.FlagSet.String("id", _nilDefaultStr, "Server's ID or UUID"),
				"format":            c.FlagSet.String("format", _nilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
				"show_credentials":  c.FlagSet.Bool("show-credentials", false, "(Flag) If set returns the servers' IPMI credentials"),
				"raw":               c.FlagSet.Bool("raw", false, "(Flag) If set returns the servers' raw object serialized using specified format"),
			}
		},
		ExecuteFunc: serverGetCmd,
		Endpoint:    DeveloperEndpoint,
	},

	{
		Description:  "Create server.",
		Subject:      "server",
		AltSubject:   "srv",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("create server", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"format":                c.FlagSet.String("format", "json", "The input format. Supported values are 'json','yaml'. The default format is json."),
				"read_config_from_file": c.FlagSet.String("raw-config", _nilDefaultStr, "(Required) Read raw object from file"),
				"read_config_from_pipe": c.FlagSet.Bool("pipe", false, "(Flag) If set, read raw object from pipe instead of from a file. Either this flag or the --raw-config option must be used."),
				"return_id":             c.FlagSet.Bool("return-id", false, "Will print the ID of the created object. Useful for automating tasks."),
			}
		},
		ExecuteFunc: serverCreateCmd,
		Endpoint:    DeveloperEndpoint,
	},

	{
		Description:  "Edit server.",
		Subject:      "server",
		AltSubject:   "srv",
		Predicate:    "edit",
		AltPredicate: "update",
		FlagSet:      flag.NewFlagSet("edit server", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"server_id_or_uuid":     c.FlagSet.String("id", _nilDefaultStr, "Server's ID or UUID"),
				"status":                c.FlagSet.String("status", _nilDefaultStr, "The new status of the server. Supported values are 'available','unavailable'. This command cannot be used in conjunction with config or pipe commands."),
				"ipmi_hostname":         c.FlagSet.String("ipmi-host", _nilDefaultStr, "The new IPMI hostname of the server. This command cannot be used in conjunction with config or pipe commands."),
				"ipmi_username":         c.FlagSet.String("ipmi-user", _nilDefaultStr, "The new IPMI username of the server. This command cannot be used in conjunction with config or pipe commands."),
				"ipmi_password":         c.FlagSet.String("ipmi-pass", _nilDefaultStr, "The new IPMI password of the server. This command cannot be used in conjunction with config or pipe commands."),
				"server_type":           c.FlagSet.String("server-type", _nilDefaultStr, "The new server type (id or label) of the server. This command cannot be used in conjunction with config or pipe commands."),
				"server_class":          c.FlagSet.String("server-class", _nilDefaultStr, "The new class of the server. This command cannot be used in conjunction with config or pipe commands."),
				"format":                c.FlagSet.String("format", "json", "The input format used when config or pipe commands are used. Supported values are 'json','yaml'. The default format is json."),
				"read_config_from_file": c.FlagSet.String("raw-config", _nilDefaultStr, "(Required) Read raw object from file"),
				"read_config_from_pipe": c.FlagSet.Bool("pipe", false, "(Flag) If set, read raw object from pipe instead of from a file. Either this flag or the --raw-config option must be used."),
			}
		},
		ExecuteFunc: serverEditCmd,
		Endpoint:    DeveloperEndpoint,
	},
	{
		Description:  "Change server power status",
		Subject:      "server",
		AltSubject:   "srv",
		Predicate:    "power-control",
		AltPredicate: "pwr",
		FlagSet:      flag.NewFlagSet("", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"server_id":   c.FlagSet.Int("id", _nilDefaultInt, "(Required) Server's id."),
				"operation":   c.FlagSet.String("operation", _nilDefaultStr, "(Required) Power control operation, one of: on, off, reset, soft."),
				"autoconfirm": c.FlagSet.Bool("autoconfirm", false, "If true it does not ask for confirmation anymore"),
			}
		},
		ExecuteFunc: serverPowerControlCmd,
	},
}

func truncateString(s string, length int) string {
	str := s
	if len(str) > 0 {
		return str[:length] + "..."
	}
	return ""
}

func serverPowerControlCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {
	serverID, ok := getIntParamOk(c.Arguments["server_id"])
	if !ok {
		return "", fmt.Errorf("-id is required")
	}
	operation, ok := getStringParamOk(c.Arguments["operation"])
	if !ok {
		return "", fmt.Errorf("-operation is required (one of: on, off, reset, soft)")
	}

	server, err := client.ServerGet(serverID, false)
	if err != nil {
		return "", err
	}

	confirm, err := confirmCommand(c, func() string {
		op := ""
		switch operation {
		case "on":
			op = "Turning on"
		case "off":
			op = "Turning off (hard)"
		case "reset":
			op = "Rebooting"
		case "soft":
			op = "Shutting down"
		}

		confirmationMessage := fmt.Sprintf("%s server (%d) of datacenter %s.  Are you sure? Type \"yes\" to continue:",
			op,
			server.ServerID,
			server.DatacenterName,
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
		err = client.ServerPowerSet(serverID, operation)
	}

	return "", err
}

func serversListCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	filter := getStringParam(c.Arguments["filter"])

	list, err := client.ServersSearch(filter)

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
			FieldName: "DATACENTER_NAME",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "SERVER_TYPE",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "STATUS",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "VENDOR",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "PRODUCT_NAME",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "SERIAL_NUMBER",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "CONFIG.",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "TAGS",
			FieldType: tableformatter.TypeString,
			FieldSize: 4,
		},
		{
			FieldName: "IPMI_HOST",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "ALLOCATED_TO.",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
	}

	showCredentials := false
	if c.Arguments["show_credentials"] != nil && *c.Arguments["show_credentials"].(*bool) {
		showCredentials = true

		schema = append(schema, tableformatter.SchemaField{
			FieldName: "IPMI_USER",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		})

		schema = append(schema, tableformatter.SchemaField{
			FieldName: "IPMI_PASS",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		})

	}

	data := [][]interface{}{}
	for _, s := range *list {

		allocation := ""
		if s.ServerStatus == "used" || s.ServerStatus == "used_registering" {
			users := strings.Join(s.UserEmail[0], ",")
			// if len(users) > 30 {
			// 	users = truncateString(users, 27)
			// }
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

	table := tableformatter.Table{
		Data:   data,
		Schema: schema,
	}
	return table.RenderTable("Servers", "", getStringParam(c.Arguments["format"]))
}

func serverGetCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	showCredentials := getBoolParam(c.Arguments["show_credentials"])

	server, err := getServerFromCommand("id", c, client, showCredentials)
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
			FieldName: "DATACENTER_NAME",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "SERVER_TYPE",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "STATUS",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "VENDOR",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "PRODUCT_NAME",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "SERIAL_NUMBER",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "CONFIG.",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "DISKS",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "TAGS",
			FieldType: tableformatter.TypeString,
			FieldSize: 4,
		},
		{
			FieldName: "IPMI_HOST",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "ALLOCATED_TO.",
			FieldType: tableformatter.TypeString,
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

		schema = append(schema, tableformatter.SchemaField{
			FieldName: "CREDENTIALS",
			FieldType: tableformatter.TypeString,
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

	format := getStringParam(c.Arguments["format"])

	if getBoolParam(c.Arguments["raw"]) {
		ret, err := tableformatter.RenderRawObject(*server, format, "Server")
		if err != nil {
			return "", err
		}
		sb.WriteString(ret)
	} else {

		switch format {
		case "json", "JSON":
			table := tableformatter.Table{
				Data:   data,
				Schema: schema,
			}
			ret, err := table.RenderTableAsJSON()
			if err != nil {
				return "", err
			}
			sb.WriteString(ret)
		case "csv", "CSV":
			table := tableformatter.Table{
				Data:   data,
				Schema: schema,
			}
			ret, err := table.RenderTableAsCSV()
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
	}

	return sb.String(), nil
}

func serverCreateCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	var obj metalcloud.Server

	err := getRawObjectFromCommand(c, &obj)

	ret, err := client.ServerCreate(obj, false)
	if err != nil {
		return "", err
	}

	if getBoolParam(c.Arguments["return_id"]) {
		return fmt.Sprintf("%d", ret), nil
	}

	return "", err
}

func serverEditCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	server, err := getServerFromCommand("id", c, client, false)
	if err != nil {
		return "", err
	}

	newStatus, setStatus := getStringParamOk(c.Arguments["status"])
	newIPMIHostname, setIPMIHostname := getStringParamOk(c.Arguments["ipmi_hostname"])
	newIPMIUsername, setIPMIUsername := getStringParamOk(c.Arguments["ipmi_username"])
	newIPMIPassword, setIPMIPassword := getStringParamOk(c.Arguments["ipmi_password"])

	_, setServerType := getStringParamOk(c.Arguments["server_type"])
	newServerClass, setServerClass := getStringParamOk(c.Arguments["server_class"])

	_, readFromFile := getStringParamOk(c.Arguments["read_config_from_file"])
	readFromPipe := getBoolParam(c.Arguments["read_config_from_pipe"])

	if (readFromFile || readFromPipe) && (setStatus || setIPMIHostname || setIPMIUsername || setIPMIPassword) {
		return "", fmt.Errorf("Cannot use --config or --pipe with --status or --ipmi-host or --ipmi-user or --ipmi-pass")
	}

	newServer := *server

	if readFromFile || readFromPipe {

		err = getRawObjectFromCommand(c, &newServer)
		if err != nil {
			return "", err
		}
	}

	if setStatus {
		newServer.ServerStatus = newStatus
	}

	if setIPMIHostname {
		newServer.ServerIPMIHost = newIPMIHostname
	}

	if setIPMIUsername {
		newServer.ServerIPMInternalUsername = newIPMIUsername
	}

	if setIPMIPassword {
		newServer.ServerIPMInternalPassword = newIPMIPassword
	}

	if setServerClass {
		newServer.ServerClass = newServerClass
	}

	if setServerType {
		serverType, err := getServerTypeFromCommand("server-type", c, client)
		if err != nil {
			return "", err
		}
		newServer.ServerTypeID = serverType.ServerTypeID

	}
	_, err = client.ServerEditComplete(server.ServerID, newServer)

	return "", err
}

func getServerFromCommand(paramName string, c *Command, client interfaces.MetalCloudClient, decryptPassword bool) (*metalcloud.Server, error) {

	m, err := getParam(c, "server_id_or_uuid", paramName)
	if err != nil {
		return nil, err
	}

	id, uuid, isID := idOrLabel(m)

	if isID {
		return client.ServerGet(id, decryptPassword)
	}

	return client.ServerGetByUUID(uuid, decryptPassword)
}

func getServerTypeFromCommand(paramName string, c *Command, client interfaces.MetalCloudClient) (*metalcloud.ServerType, error) {

	m, err := getParam(c, "server_type", paramName)
	if err != nil {
		return nil, err
	}

	id, label, isID := idOrLabel(m)

	if isID {
		return client.ServerTypeGet(id)
	}

	return client.ServerTypeGetByLabel(label)
}
