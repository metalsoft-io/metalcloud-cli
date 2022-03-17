package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
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
				"format":              c.FlagSet.String("format", _nilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
				"filter":              c.FlagSet.String("filter", "*", "filter to use when searching for servers. Check the documentation for examples. Defaults to '*'"),
				"show_credentials":    c.FlagSet.Bool("show-credentials", false, green("(Flag)")+" If set returns the servers' IPMI credentials. (Slow for large queries)"),
				"show_rack_info":      c.FlagSet.Bool("show-rack-info", false, green("(Flag)")+" If set returns the servers' rack metadata"),
				"show_hardware":       c.FlagSet.Bool("show-hardware", false, green("(Flag)")+" If set returns the servers' hardware configuration"),
				"show_decommissioned": c.FlagSet.Bool("show-decommissioned", false, green("(Flag)")+" If set returns decommissioned servers which are normally hidden"),
			}
		},
		ExecuteFunc: serversListCmd,
		Endpoint:    DeveloperEndpoint,
		Example: `
metalcloud-cli server list --filter "available used" # to show all available and used servers. One of: [available|unavailable|used|cleaning|registering]
metalcloud-cli server list --show-credentials # to retrieve a list of credentials. Note: this will take a longer time.
		`,
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
				"show_credentials":  c.FlagSet.Bool("show-credentials", false, green("(Flag)")+" If set returns the servers' IPMI credentials"),
				"raw":               c.FlagSet.Bool("raw", false, green("(Flag)")+" If set returns the servers' raw object serialized using specified format"),
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
				"read_config_from_file": c.FlagSet.String("raw-config", _nilDefaultStr, red("(Required)")+" Read raw object from file"),
				"read_config_from_pipe": c.FlagSet.Bool("pipe", false, green("(Flag)")+" If set, read raw object from pipe instead of from a file. Either this flag or the --raw-config option must be used."),
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
				"read_config_from_file": c.FlagSet.String("raw-config", _nilDefaultStr, red("(Required)")+" Read raw object from file"),
				"read_config_from_pipe": c.FlagSet.Bool("pipe", false, green("(Flag)")+" If set, read raw object from pipe instead of from a file. Either this flag or the --raw-config option must be used."),
			}
		},
		ExecuteFunc: serverEditCmd,
		Endpoint:    DeveloperEndpoint,
	},
	{
		Description:  "Edit server's IPMI",
		Subject:      "server",
		AltSubject:   "srv",
		Predicate:    "edit-ipmi",
		AltPredicate: "update-ipmi",
		FlagSet:      flag.NewFlagSet("edit server IPMI", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"server_id_or_uuid": c.FlagSet.String("id", _nilDefaultStr, "Server's ID or UUID"),
				"ipmi_hostname":     c.FlagSet.String("ipmi-host", _nilDefaultStr, "The new IPMI hostname of the server. This command cannot be used in conjunction with config or pipe commands."),
				"ipmi_username":     c.FlagSet.String("ipmi-user", _nilDefaultStr, "The new IPMI username of the server. This command cannot be used in conjunction with config or pipe commands."),
				"ipmi_password":     c.FlagSet.String("ipmi-pass", _nilDefaultStr, "The new IPMI password of the server. This command cannot be used in conjunction with config or pipe commands."),
			}
		},
		ExecuteFunc: serverEditIPMICmd,
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
				"server_id":   c.FlagSet.Int("id", _nilDefaultInt, red("(Required)")+" Server's id."),
				"operation":   c.FlagSet.String("operation", _nilDefaultStr, red("(Required)")+" Power control operation, one of: on, off, reset, soft."),
				"autoconfirm": c.FlagSet.Bool("autoconfirm", false, green("(Flag)")+" If set it will assume action is confirmed"),
			}
		},
		ExecuteFunc: serverPowerControlCmd,
		Endpoint:    DeveloperEndpoint,
	},
	{
		Description:  "Change server status",
		Subject:      "server",
		AltSubject:   "srv",
		Predicate:    "status-set",
		AltPredicate: "status",
		FlagSet:      flag.NewFlagSet("", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"server_id":   c.FlagSet.Int("id", _nilDefaultInt, red("(Required)")+" Server's id."),
				"status":      c.FlagSet.String("status", _nilDefaultStr, red("(Required)")+" New server status. One of: 'available','decommissioned','removed_from_rack'"),
				"autoconfirm": c.FlagSet.Bool("autoconfirm", false, green("(Flag)")+" If set it will assume action is confirmed"),
			}
		},
		ExecuteFunc: serverStatusSetCmd,
		Endpoint:    DeveloperEndpoint,
	},
	{
		Description:  "Reregister server",
		Subject:      "server",
		AltSubject:   "srv",
		Predicate:    "reregister",
		AltPredicate: "re-register",
		FlagSet:      flag.NewFlagSet("", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"server_id":   c.FlagSet.Int("id", _nilDefaultInt, red("(Required)")+" Server's id."),
				"skip_ipmi":   c.FlagSet.Bool("do-not-set-ipmi", false, "If set, the system will not change the IPMI credentials."),
				"autoconfirm": c.FlagSet.Bool("autoconfirm", false, green("(Flag)")+" If set it will assume action is confirmed"),
			}
		},
		ExecuteFunc: serverReregisterCmd,
		Endpoint:    DeveloperEndpoint,
	},
	{
		Description:  "Change server server type",
		Subject:      "server",
		AltSubject:   "srv",
		Predicate:    "server-type-set",
		AltPredicate: "server-type",
		FlagSet:      flag.NewFlagSet("", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"server_id":   c.FlagSet.Int("id", _nilDefaultInt, red("(Required)")+" Server's id."),
				"server_type": c.FlagSet.String("server-type", _nilDefaultStr, red("(Required)")+" New server type. Can be an ID or label"),
				"autoconfirm": c.FlagSet.Bool("autoconfirm", false, green("(Flag)")+" If set it will assume action is confirmed"),
			}
		},
		ExecuteFunc: serverServerTypeSetCmd,
		Endpoint:    DeveloperEndpoint,
	},
	{
		Description:  "Change server rack information",
		Subject:      "server",
		AltSubject:   "srv",
		Predicate:    "rack-info-set",
		AltPredicate: "rack-info",
		FlagSet:      flag.NewFlagSet("", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"server_id":    c.FlagSet.Int("id", _nilDefaultInt, red("(Required)")+" Server's id."),
				"inventory_id": c.FlagSet.String("inventory-id", _nilDefaultStr, " New inventory id"),
				"rack_name":    c.FlagSet.String("rack-name", _nilDefaultStr, red("(Required)")+" New rack name."),
				"lower_u":      c.FlagSet.Int("lower-u", _nilDefaultInt, red("(Required)")+" Lower U of the equipment"),
				"upper_u":      c.FlagSet.Int("upper-u", _nilDefaultInt, red("(Required)")+" Upper U of the equipment"),
				"autoconfirm":  c.FlagSet.Bool("autoconfirm", false, green("(Flag)")+" If set it will assume action is confirmed"),
			}
		},
		ExecuteFunc: serverRackInfoSetCmd,
		Endpoint:    DeveloperEndpoint,
	},
	{
		Description:  "Lists server interfaces.",
		Subject:      "server",
		AltSubject:   "srv",
		Predicate:    "interfaces",
		AltPredicate: "intf",
		FlagSet:      flag.NewFlagSet("list server interfaces", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"format":            c.FlagSet.String("format", _nilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
				"server_id_or_uuid": c.FlagSet.Int("id", _nilDefaultInt, red("(Required)")+" Server's id."),
				"raw":               c.FlagSet.Bool("raw", false, green("(Flag)")+" When set the return will be a full dump of the object. This is useful when copying configurations. Only works with json and yaml formats."),
			}
		},
		ExecuteFunc: serverInterfacesListCmd,
		Endpoint:    DeveloperEndpoint,
	},
}

func serverPowerControlCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {
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

func serverStatusSetCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {
	serverID, ok := getIntParamOk(c.Arguments["server_id"])
	if !ok {
		return "", fmt.Errorf("-id is required")
	}

	newStatus, ok := getStringParamOk(c.Arguments["status"])
	if !ok {
		return "", fmt.Errorf("-status is required (one of: on, off, reset, soft)")
	}

	var server metalcloud.Server

	if !getBoolParam(c.Arguments["autoconfirm"]) {
		serverPtr, err := client.ServerGet(serverID, false)
		if err != nil {
			return "", err
		}
		server = *serverPtr
	}

	confirm, err := confirmCommand(c, func() string {

		confirmationMessage := ""

		if !getBoolParam(c.Arguments["autoconfirm"]) {

			confirmationMessage = fmt.Sprintf("Server #%s (%s) of datacenter %s. Current status: %s new status: %s  Are you sure? Type \"yes\" to continue:",
				blue(fmt.Sprintf("%d", server.ServerID)),
				yellow(server.ServerSerialNumber),
				server.DatacenterName,
				colorizeServerStatus(server.ServerStatus),
				colorizeServerStatus(newStatus),
			)
		}

		if strings.HasSuffix(os.Args[0], ".test") {
			confirmationMessage = ""
		}

		return confirmationMessage
	})

	if err != nil {
		return "", err
	}

	if confirm {
		err = client.ServerStatusUpdate(serverID, newStatus)
	}

	return "", err
}

func serverServerTypeSetCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {
	serverID, ok := getIntParamOk(c.Arguments["server_id"])
	if !ok {
		return "", fmt.Errorf("-id is required")
	}

	serverTypeStr, ok := getStringParamOk(c.Arguments["server_type"])
	if !ok {
		return "", fmt.Errorf("-server-type is required")
	}

	serverTypeID, _, isID := idOrLabel(serverTypeStr)
	var newServerType metalcloud.ServerType
	if !isID {
		fmt.Printf("%s", serverTypeStr)
		st, err := client.ServerTypeGetByLabel(serverTypeStr)
		if err != nil {
			return "", err
		}
		fmt.Printf("here: %v", newServerType)
		newServerType = *st

	} else {
		st, err := client.ServerTypeGet(serverTypeID)
		if err != nil {
			return "", err
		}
		newServerType = *st
	}

	var server metalcloud.Server

	if !getBoolParam(c.Arguments["autoconfirm"]) {
		serverPtr, err := client.ServerGet(serverID, false)
		if err != nil {
			return "", err
		}
		server = *serverPtr
	}

	confirm, err := confirmCommand(c, func() string {

		confirmationMessage := ""

		if !getBoolParam(c.Arguments["autoconfirm"]) {

			oldServerType, err := client.ServerTypeGet(server.ServerTypeID)
			if err != nil {
				return err.Error()
			}

			confirmationMessage = fmt.Sprintf("Server #%s (%s) of datacenter %s. Current server type: %s (#%s) new server type: %s (#%s) Are you sure? Type \"yes\" to continue:",
				blue(fmt.Sprintf("%d", server.ServerID)),
				yellow(server.ServerSerialNumber),
				server.DatacenterName,
				red(oldServerType.ServerTypeName),
				red(oldServerType.ServerTypeID),
				green(newServerType.ServerTypeName),
				green(newServerType.ServerTypeID),
			)
		}

		if strings.HasSuffix(os.Args[0], ".test") {
			confirmationMessage = ""
		}

		return confirmationMessage
	})

	if err != nil {
		return "", err
	}

	if confirm {
		err = client.ServerEditProperty(serverID, "server_type_id", newServerType.ServerTypeID)
	}

	return "", err
}

func serverRackInfoSetCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {
	serverID, ok := getIntParamOk(c.Arguments["server_id"])
	if !ok {
		return "", fmt.Errorf("-id is required")
	}

	serverRackName, ok := getStringParamOk(c.Arguments["rack_name"])
	if !ok {
		return "", fmt.Errorf("-rack-name is required")
	}

	serverRackLowerU, ok := getIntParamOk(c.Arguments["lower_u"])
	if !ok {
		return "", fmt.Errorf("-lower-u is required")
	}

	serverRackUpperU, ok := getIntParamOk(c.Arguments["upper_u"])
	if !ok {
		return "", fmt.Errorf("-upper-u is required")
	}

	var server metalcloud.Server

	serverPtr, err := client.ServerGet(serverID, false)
	if err != nil {
		return "", err
	}
	server = *serverPtr

	confirm, err := confirmCommand(c, func() string {

		confirmationMessage := ""

		if !getBoolParam(c.Arguments["autoconfirm"]) {

			oldRackInfo := getRackInfoSafe(server)

			oldServerRackInfo := fmt.Sprintf("InvID:%s Rack:%s U:%s-%s", oldRackInfo.InventoryID, oldRackInfo.RackName, oldRackInfo.LowerU, oldRackInfo.UpperU)

			serverInventoryIDStr := ""
			serverInventoryID, ok := getStringParamOk(c.Arguments["inventory_id"])
			if ok {
				serverInventoryIDStr = serverInventoryID
			}
			newServerRackInfo := fmt.Sprintf("InvID:%s Rack:%s U:%d-%d", serverInventoryIDStr, serverRackName, serverRackLowerU, serverRackUpperU)

			confirmationMessage = fmt.Sprintf("Server #%s (%s) of datacenter %s. Current server rack info %s new rack info: %s. Are you sure? Type \"yes\" to continue:",
				blue(fmt.Sprintf("%d", server.ServerID)),
				yellow(server.ServerSerialNumber),
				server.DatacenterName,
				red(oldServerRackInfo),
				green(newServerRackInfo),
			)
		}

		if strings.HasSuffix(os.Args[0], ".test") {
			confirmationMessage = ""
		}

		return confirmationMessage
	})

	if err != nil {
		return "", err
	}

	if confirm {
		server.ServerRackName = &serverRackName

		lowerUStr := fmt.Sprintf("%d", serverRackLowerU)
		server.ServerRackPositionLowerUnit = &lowerUStr

		upperUStr := fmt.Sprintf("%d", serverRackUpperU)
		server.ServerRackPositionUpperUnit = &upperUStr

		serverInventoryID, ok := getStringParamOk(c.Arguments["inventory_id"])
		if ok {
			server.ServerInventoryId = &serverInventoryID
		}

		_, err = client.ServerEdit(serverID, "complete", server)
	}

	return "", err
}

func serverReregisterCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {
	serverID, ok := getIntParamOk(c.Arguments["server_id"])
	if !ok {
		return "", fmt.Errorf("-id is required")
	}

	skipIpmi := getBoolParam(c.Arguments["skip_ipmi"])

	var server metalcloud.Server

	if !getBoolParam(c.Arguments["autoconfirm"]) {
		serverPtr, err := client.ServerGet(serverID, false)
		if err != nil {
			return "", err
		}
		server = *serverPtr
	}

	confirm, err := confirmCommand(c, func() string {

		confirmationMessage := ""

		if !getBoolParam(c.Arguments["autoconfirm"]) {

			confirmationMessage = fmt.Sprintf("Server #%s (%s) BMC IP:%s of datacenter %s. Are you sure? Type \"yes\" to continue:",
				blue(fmt.Sprintf("%d", server.ServerID)),
				yellow(server.ServerSerialNumber),
				server.ServerIPMIHost,
				server.DatacenterName,
			)
		}

		if strings.HasSuffix(os.Args[0], ".test") {
			confirmationMessage = ""
		}

		return confirmationMessage
	})

	if err != nil {
		return "", err
	}

	if confirm {
		err = client.ServerReregister(serverID, skipIpmi, false)
	}

	return "", err
}

func colorizeServerStatus(status string) string {

	switch status {
	case "available":
		return blue(status)
	case "used":
		return green(status)
	case "unavailable":
		return magenta(status)
	}
	return yellow(status)

}

func serversListCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	filter := getStringParam(c.Arguments["filter"])

	list, err := client.ServersSearch(convertToSearchFieldFormat(filter))
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
			FieldName: "STATUS",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},

		{
			FieldName: "SERVER_TYPE",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},

		{
			FieldName: "SERIAL_NUMBER",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "IPMI_HOST",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "ALLOCATED_TO",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "DATACENTER_NAME",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
	}

	if getBoolParam(c.Arguments["show_rack_info"]) {

		extraFields := []tableformatter.SchemaField{
			{
				FieldName: "TAGS",
				FieldType: tableformatter.TypeString,
				FieldSize: 4,
			},
			{
				FieldName: "INV_ID",
				FieldType: tableformatter.TypeString,
				FieldSize: 4,
			},
			{
				FieldName: "RACK",
				FieldType: tableformatter.TypeString,
				FieldSize: 4,
			},
			{
				FieldName: "RU_D",
				FieldType: tableformatter.TypeString,
				FieldSize: 4,
			},
			{
				FieldName: "RU_U",
				FieldType: tableformatter.TypeString,
				FieldSize: 4,
			},
		}

		schema = append(schema, extraFields...)

	}

	serverInterfaces := map[int][]metalcloud.SwitchInterfaceSearchResult{}

	if getBoolParam(c.Arguments["show_hardware"]) {
		extraFields := []tableformatter.SchemaField{
			{
				FieldName: "CONFIG.",
				FieldType: tableformatter.TypeString,
				FieldSize: 5,
			},
		}
		schema = append(schema, extraFields...)

		//retrieve interface information, it will help us show a more detailed data on
		//NICs.
		serverInterfacesList, err := client.SwitchInterfaceSearch("*")

		if err != nil {
			return "", err
		}

		//We save it in a map indexed by server id for quicker retrieval later
		for _, serverInterface := range *serverInterfacesList {
			serverInterfaces[serverInterface.ServerID] = append(serverInterfaces[serverInterface.ServerID], serverInterface)
		}
	}

	if getBoolParam(c.Arguments["show_credentials"]) {

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

		schema = append(schema, tableformatter.SchemaField{
			FieldName: "IPMI_SNMP_COMMUNITY",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		})

	}

	data := [][]interface{}{}

	statusCounts := map[string]int{
		"available":      0,
		"cleaning":       0,
		"registering":    0,
		"used":           0,
		"decommissioned": 0,
	}

	for _, s := range *list {

		if s.ServerStatus == "decommissioned" && !getBoolParam(c.Arguments["show_decommissioned"]) {
			continue
		}

		statusCounts[s.ServerStatus] = statusCounts[s.ServerStatus] + 1

		allocation := ""
		if s.ServerStatus == "used" || s.ServerStatus == "used_registering" {
			users := strings.Join(s.UserEmail[0], ",")

			allocation = fmt.Sprintf("%s %s (#%d) IA:#%d Infra:#%d",
				users,
				s.InstanceLabel[0],
				s.InstanceID[0],
				s.InstanceArrayID[0],
				s.InfrastructureID[0])
			if len(allocation) > 30 {
				allocation = truncateString(allocation, 10)
			}
		}

		credentialsUser := ""
		credentialsPass := ""
		//snmpCommunity := ""

		if getBoolParam(c.Arguments["show_credentials"]) {

			server, err := client.ServerGet(s.ServerID, true)

			if err != nil {
				return "", err
			}

			credentialsUser = fmt.Sprintf("%s", server.ServerIPMInternalUsername)
			credentialsPass = fmt.Sprintf("%s", server.ServerIPMInternalPassword)
			//snmpCommunity = fmt.Sprintf("%s", server.ServerMgmtSNMPCommunityPassword)

		}

		diskDescription := ""
		if s.ServerDiskCount > 0 {
			diskDescription = fmt.Sprintf("%s x %s GB %s",
				yellow(s.ServerDiskCount),
				yellow(s.ServerDiskSizeMbytes/1000),
				yellow(s.ServerDiskType))
		}

		//we index by capacity
		interfacesByCapacity := map[int][]metalcloud.SwitchInterfaceSearchResult{}

		for _, serverInterface := range serverInterfaces[s.ServerID] {
			interfacesByCapacity[serverInterface.ServerInterfaceCapacityMBPs] = append(interfacesByCapacity[serverInterface.ServerInterfaceCapacityMBPs], serverInterface)
		}

		interfaceDescription := ""

		for capacity, serverInterfaces := range interfacesByCapacity {
			interfaceDescription = interfaceDescription +
				fmt.Sprintf("%s x %s Gbps NICs",
					magenta(len(serverInterfaces)),
					magenta(capacity/1000))
		}

		hardwareConfig := ""
		if s.ServerProcessorCount > 0 {
			hardwareConfig = fmt.Sprintf("%s x %s cores(ht), %s GB RAM, %s, %s",
				blue(s.ServerProcessorCount),
				blue(s.ServerProcessorThreads*s.ServerProcessorCoreCount),
				red(s.ServerRAMGbytes),
				diskDescription,
				interfaceDescription,
			)
		}

		status := s.ServerStatus

		switch status {
		case "available":
			status = blue(status)
		case "used":
			status = green(status)
		case "unavailable":
			status = magenta(status)
		case "defective":
			status = red(status)
		default:
			status = yellow(status)

		}

		row := []interface{}{
			s.ServerID,
			status,
			s.ServerTypeName,
			s.ServerSerialNumber,
			s.ServerIPMIHost,
			allocation,
			s.DatacenterName,
		}

		if getBoolParam(c.Arguments["show_rack_info"]) {

			row = append(row, []interface{}{
				strings.Join(s.ServerTags, ","),
				s.ServerInventoryId,
				s.ServerRackName,
				s.ServerRackPositionLowerUnit,
				s.ServerRackPositionUpperUnit,
			}...)
		}

		if getBoolParam(c.Arguments["show_hardware"]) {
			row = append(row, []interface{}{
				hardwareConfig,
			}...)
		}

		if getBoolParam(c.Arguments["show_credentials"]) {
			row = append(row, []interface{}{
				credentialsUser,
				credentialsPass,
				//snmpCommunity,
			}...)
		}

		data = append(data, row)

	}

	table := tableformatter.Table{
		Data:   data,
		Schema: schema,
	}

	title := fmt.Sprintf("Servers: %d available %d used %d cleaning %d registering %d unavailable",
		statusCounts["available"],
		statusCounts["used"],
		statusCounts["cleaning"],
		statusCounts["registering"],
		statusCounts["unavailable"])

	if getBoolParam(c.Arguments["show_decommissioned"]) {
		title = title + fmt.Sprintf(" %d decommissioned", statusCounts["decommissioned"])
	}

	return table.RenderTable(title, "", getStringParam(c.Arguments["format"]))
}

func serverGetCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

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
			FieldName: "SERIAL NUMBER",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "DATACENTER_NAME",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "INVENTORY_ID",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "RACK_NAME",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "RACK_POSITION_LOWER_UNIT",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "RACK_POSITION_UPPER_UNIT",
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
	snmpCommunity := ""

	if showCredentials {

		schema = append(schema, tableformatter.SchemaField{
			FieldName: "CREDENTIALS",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		})

		credentials = fmt.Sprintf("User: %s Pass: %s", server.ServerIPMInternalUsername, server.ServerIPMInternalPassword)

		schema = append(schema, tableformatter.SchemaField{
			FieldName: "SNMP_COMMUNITY",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		})

		//snmpCommunity = server.ServerMgmtSNMPCommunityPassword

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

	serverInventoryID := ""
	if server.ServerInventoryId != nil {
		serverInventoryID = *server.ServerInventoryId
	}

	serverRackName := ""
	if server.ServerRackName != nil {
		serverRackName = *server.ServerRackName
	}

	serverRackPositionLowerUnit := ""
	if server.ServerRackPositionLowerUnit != nil {
		serverRackPositionLowerUnit = *server.ServerRackPositionLowerUnit
	}
	serverRackPositionUpperUnit := ""
	if server.ServerRackPositionUpperUnit != nil {
		serverRackPositionUpperUnit = *server.ServerRackPositionUpperUnit
	}

	data = append(data, []interface{}{
		server.ServerID,
		server.ServerSerialNumber,
		server.DatacenterName,
		serverInventoryID,
		serverRackName,
		serverRackPositionLowerUnit,
		serverRackPositionUpperUnit,
		serverTypeName,
		server.ServerStatus,
		server.ServerVendor,
		productName,
		configuration,
		disks,
		strings.Join(server.ServerTags, ","),
		server.ServerIPMIHost,
		allocation,
		credentials,
		snmpCommunity,
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
			table := tableformatter.Table{
				Data:   data,
				Schema: schema,
			}
			ret, err := table.RenderTransposedTable("server details", "", format)
			if err != nil {
				return "", err
			}

			sb.WriteString(ret)
		}
	}

	return sb.String(), nil
}

func serverCreateCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

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

func serverEditCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

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

func serverEditIPMICmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	server, err := getServerFromCommand("id", c, client, false)
	if err != nil {
		return "", err
	}

	newIPMIHostname, setIPMIHostname := getStringParamOk(c.Arguments["ipmi_hostname"])
	newIPMIUsername, setIPMIUsername := getStringParamOk(c.Arguments["ipmi_username"])
	newIPMIPassword, setIPMIPassword := getStringParamOk(c.Arguments["ipmi_password"])

	newServer := *server

	if setIPMIHostname {
		newServer.ServerIPMIHost = newIPMIHostname
	}

	if setIPMIUsername {
		newServer.ServerIPMInternalUsername = newIPMIUsername
	}

	if setIPMIPassword {
		newServer.ServerIPMInternalPassword = newIPMIPassword
	}

	_, err = client.ServerEditIPMI(server.ServerID, newServer)

	return "", err
}

func getServerFromCommand(paramName string, c *Command, client metalcloud.MetalCloudClient, decryptPassword bool) (*metalcloud.Server, error) {

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

func getServerTypeFromCommand(paramName string, c *Command, client metalcloud.MetalCloudClient) (*metalcloud.ServerType, error) {

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

func serverInterfacesListCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	server, err := getServerFromCommand("id", c, client, false)
	if err != nil {
		return "", err
	}

	list, err := client.SwitchInterfaceSearch(fmt.Sprintf("server_id:%d", server.ServerID))

	if err != nil {
		return "", err
	}

	schema := []tableformatter.SchemaField{
		{
			FieldName: "IDX",
			FieldType: tableformatter.TypeInt,
			FieldSize: 5,
		},
		{
			FieldName: "TYPE",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "SERVER INTERFACE",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "SWITCH INTERFACE",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "SWITCH",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "SWITCH MGMT",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "SWITCH INTERFACE MAC",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "CAPACITY",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},

		{
			FieldName: "IP",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
	}

	data := [][]interface{}{}
	for _, s := range *list {

		ips := flattenAndJoinStrings(s.IP)
		networkType := strings.Join(s.NetworkType, ",")

		switch_info := fmt.Sprintf("%s (#%d)",
			s.NetworkEquipmentIdentifierString,
			s.NetworkEquipmentID,
		)

		capacity := fmt.Sprintf("%d Gbps", int(s.ServerInterfaceCapacityMBPs/1000))

		data = append(data, []interface{}{
			s.ServerInterfaceIndex,
			networkType,
			s.ServerInterfaceMACAddress,
			s.NetworkEquipmentInterfaceIdentifierString,
			switch_info,
			s.NetworkEquipmentManagementAddress,
			s.NetworkEquipmentInterfaceMACAddress,
			capacity,
			ips,
		})

	}

	tableformatter.TableSorter(schema).OrderBy(schema[0].FieldName).Sort(data)

	var sb strings.Builder

	format := getStringParam(c.Arguments["format"])

	if getBoolParam(c.Arguments["raw"]) {
		ret, err := tableformatter.RenderRawObject(*list, format, "Server interfaces")
		if err != nil {
			return "", err
		}
		sb.WriteString(ret)
	} else {
		table := tableformatter.Table{
			Data:   data,
			Schema: schema,
		}
		ret, err := table.RenderTable(fmt.Sprintf("Server interfaces of server #%d %s", server.ServerID, server.ServerSerialNumber), "", format)
		if err != nil {
			return "", err
		}
		sb.WriteString(ret)
	}

	return sb.String(), nil
}

type rackInfo struct {
	InventoryID string
	RackName    string
	LowerU      string
	UpperU      string
}

func getRackInfoSafe(server metalcloud.Server) rackInfo {
	return rackInfo{
		InventoryID: getStringFromStringOrEmpty(server.ServerInventoryId),
		RackName:    getStringFromStringOrEmpty(server.ServerInventoryId),
		LowerU:      getStringFromStringOrEmpty(server.ServerRackPositionLowerUnit),
		UpperU:      getStringFromStringOrEmpty(server.ServerRackPositionUpperUnit),
	}
}

func getStringFromStringOrEmpty(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}
