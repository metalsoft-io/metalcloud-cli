package server

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
	"github.com/metalsoft-io/metalcloud-cli/internal/colors"
	"github.com/metalsoft-io/metalcloud-cli/internal/command"
	"github.com/metalsoft-io/metalcloud-cli/internal/configuration"
	"github.com/metalsoft-io/metalcloud-cli/internal/filtering"
	"github.com/metalsoft-io/metalcloud-cli/internal/stringutils"
	"github.com/metalsoft-io/tableformatter"
	"gopkg.in/yaml.v2"
)

var ServersCmds = []command.Command{
	{
		Description:  "Lists all servers.",
		Subject:      "server",
		AltSubject:   "srv",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list servers", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"format":              c.FlagSet.String("format", command.NilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
				"filter":              c.FlagSet.String("filter", "*", "filter to use when searching for servers. Check the documentation for examples. Defaults to '*'"),
				"show_credentials":    c.FlagSet.Bool("show-credentials", false, colors.Green("(Flag)")+" If set returns the servers' IPMI credentials. (Slow for large queries)"),
				"show_rack_info":      c.FlagSet.Bool("show-rack-info", false, colors.Green("(Flag)")+" If set returns the servers' rack metadata"),
				"show_hardware":       c.FlagSet.Bool("show-hardware", false, colors.Green("(Flag)")+" If set returns the servers' hardware configuration"),
				"show_decommissioned": c.FlagSet.Bool("show-decommissioned", false, colors.Green("(Flag)")+" If set returns decommissioned servers which are normally hidden"),
			}
		},
		ExecuteFunc: serversListCmd,
		Endpoint:    configuration.DeveloperEndpoint,
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
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"server_id_or_uuid": c.FlagSet.String("id", command.NilDefaultStr, "Server's ID or UUID"),
				"format":            c.FlagSet.String("format", command.NilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
				"show_credentials":  c.FlagSet.Bool("show-credentials", false, colors.Green("(Flag)")+" If set returns the servers' IPMI credentials"),
				"raw":               c.FlagSet.Bool("raw", false, colors.Green("(Flag)")+" If set returns the servers' raw object serialized using specified format"),
			}
		},
		ExecuteFunc: serverGetCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},

	{
		Description:  "Create server.",
		Subject:      "server",
		AltSubject:   "srv",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("create server", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"format":                c.FlagSet.String("format", "json", "The input format. Supported values are 'json','yaml'. The default format is json."),
				"read_config_from_file": c.FlagSet.String("raw-config", command.NilDefaultStr, colors.Red("(Required)")+" Read raw object from file"),
				"read_config_from_pipe": c.FlagSet.Bool("pipe", false, colors.Green("(Flag)")+" If set, read raw object from pipe instead of from a file. Either this flag or the --raw-config option must be used."),
				"return_id":             c.FlagSet.Bool("return-id", false, "Will print the ID of the created object. Useful for automating tasks."),
			}
		},
		ExecuteFunc: serverCreateCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},

	{
		Description:  "Register a server.",
		Subject:      "server",
		AltSubject:   "srv",
		Predicate:    "register",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("register server", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"datacenter":    c.FlagSet.String("datacenter", command.NilDefaultStr, colors.Red("(Required)")+" The datacenter in which this server is to be registered."),
				"server_vendor": c.FlagSet.String("server-vendor", command.NilDefaultStr, colors.Red("(Required)")+" Server vendor (driver) to use when interacting with the server. One of: `dell`,'hpe_legacy','hpe'."),
				"mgmt_address":  c.FlagSet.String("mgmt-address", command.NilDefaultStr, colors.Red("(Required)")+" IP or DNS record for the server's management interface (BMC)."),
				"mgmt_user":     c.FlagSet.String("mgmt-user", command.NilDefaultStr, colors.Red("(Required)")+" Server' BMC username."),
				"mgmt_pass":     c.FlagSet.String("mgmt-pass", command.NilDefaultStr, colors.Red("(Required)")+" Server' BMC password."),
				"return_id":     c.FlagSet.Bool("return-id", false, "Will print the ID of the created object. Useful for automating tasks."),
			}
		},
		ExecuteFunc: serverRegisterCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},

	{
		Description:  "Import server.",
		Subject:      "server",
		AltSubject:   "srv",
		Predicate:    "import",
		AltPredicate: "import-unmanaged",
		FlagSet:      flag.NewFlagSet("import server", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"format":                c.FlagSet.String("format", "json", "The input format. Supported values are 'json','yaml'. The default format is json."),
				"read_config_from_file": c.FlagSet.String("file", command.NilDefaultStr, colors.Red("(Required)")+" Read raw object from file"),
				"read_config_from_pipe": c.FlagSet.Bool("pipe", false, colors.Green("(Flag)")+" If set, read raw object from pipe instead of from a file. Either this flag or the --file option must be used."),
				"add_to_infra":          c.FlagSet.String("add-to-infra", command.NilDefaultStr, colors.Green("(Optional)")+" The infrastructure to use to add this server to. If set to 'auto' will use the settings in the file instead."),
				"return_id":             c.FlagSet.Bool("return-id", false, "Will print the ID of the created object. Useful for automating tasks."),
			}
		},
		ExecuteFunc: serverImportCmd,
		Endpoint:    configuration.DeveloperEndpoint,
		Example: `
The following fields are required:

#Example1
#========
datacenter: sonic-qts
serialNumber: NNAACC2
#serverType: M.15.15.1
serverTypeID: 9
interfaces:
   - mac: 00:B0:D0:63:C2:26
     switch: leaf-124
     switchInterface: Ethernet216
   - mac: aa:bb:cc:dd:02:ff
     switch: leaf-124
     switchInterface: Ethernet217

#Example2
#========
datacenter: sonic-qts
serialNumber: NNAACC2
serverType: M.15.15.1
interfaces:
   - mac: 00:B0:D0:63:C2:26
     switch: leaf-124
     switchInterface: Ethernet216
   - mac: aa:bb:cc:dd:02:ff
     switch: leaf-124
     switchInterface: Ethernet217

#Example3
#========
datacenter: sonic-qts
serialNumber: NNAACC2
serverType: M.15.15.1
label: testserv
infrastructure: myinfra
userEmail: alex@test.io
interfaces:
- mac: 00:B0:D0:63:C2:26
	switch: leaf-124
	switchInterface: Ethernet216
- mac: aa:bb:cc:dd:02:ff
	switch: leaf-124
	switchInterface: Ethernet217

$ metalcloud-cli server import -format yaml -file ./input.yaml
`,
	},

	{
		Description:  "Add a server to an infrastructure as an instance array",
		Subject:      "server",
		AltSubject:   "srv",
		Predicate:    "add-to-infra",
		AltPredicate: "add-to-infrastructure",
		FlagSet:      flag.NewFlagSet("Add server to infrastructure", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"server_id":                  c.FlagSet.Int("server-id", command.NilDefaultInt, colors.Red("(Required)")+" The server id"),
				"infrastructure_id_or_label": c.FlagSet.String("infra", command.NilDefaultStr, colors.Red("(Required)")+" The infrastructure ID or Label. Must exist"),
				"return_id":                  c.FlagSet.Bool("return-id", false, "Will print the ID of the created instance array. Useful for automating tasks."),
			}
		},
		ExecuteFunc: serverAddToInfraCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},

	{
		Description:  "Import server batch.",
		Subject:      "server",
		AltSubject:   "srv",
		Predicate:    "import-batch",
		AltPredicate: "import-batch",
		FlagSet:      flag.NewFlagSet("import server", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"format":                c.FlagSet.String("format", "json", "The input format. Supported values are 'json','yaml'. The only supported format is yaml."),
				"read_config_from_file": c.FlagSet.String("file", command.NilDefaultStr, colors.Red("(Required)")+" Read raw object from file"),
				"add_to_infra":          c.FlagSet.String("add-to-infra", command.NilDefaultStr, colors.Green("(Optional)")+" The infrastructure to use to add this server to. If set to 'auto' will use the settings in the file instead."),
				"return_id":             c.FlagSet.Bool("return-id", false, "(Optional) Will print the ID of the created object. Useful for automating tasks."),
			}
		},
		ExecuteFunc: serverImportBatchCmd,
		Endpoint:    configuration.DeveloperEndpoint,
		Example: `
This command is the batch version of server import. The file format uses "---" separator between records. For example:


datacenter: sonic-qts
serialNumber: NNAACC2
serverType: M.15.15.1
label: testserv
infrastructure: myinfra
userEmail: alex@test.io
interfaces:
- mac: 00:B0:D0:63:C2:26
	switch: leaf-124
	switchInterface: Ethernet216
- mac: aa:bb:cc:dd:02:ff
	switch: leaf-124
	switchInterface: Ethernet217
---
datacenter: sonic-qts
serialNumber: NNAACC3
serverType: M.15.15.1
label: testserv
infrastructure: myinfra
userEmail: alex@test.io
interfaces:
- mac: 00:B0:D0:63:C2:26
	switch: leaf-124
	switchInterface: Ethernet218
- mac: aa:bb:cc:dd:02:ff
	switch: leaf-124
	switchInterface: Ethernet219
---
datacenter: sonic-qts
serialNumber: NNAACC3
serverType: M.15.15.1
label: testserv
infrastructure: myinfra
userEmail: alex@test.io
interfaces:
- mac: 00:B0:D0:63:C2:26
	switch: leaf-124
	switchInterface: Ethernet220
- mac: aa:bb:cc:dd:02:ff
	switch: leaf-124
	switchInterface: Ethernet221
`,
	},

	{
		Description:  "Edit server.",
		Subject:      "server",
		AltSubject:   "srv",
		Predicate:    "edit",
		AltPredicate: "update",
		FlagSet:      flag.NewFlagSet("edit server", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"server_id_or_uuid":     c.FlagSet.String("id", command.NilDefaultStr, "Server's ID or UUID"),
				"status":                c.FlagSet.String("status", command.NilDefaultStr, "The new status of the server. Supported values are 'available','unavailable'. This command cannot be used in conjunction with config or pipe commands."),
				"ipmi_hostname":         c.FlagSet.String("ipmi-host", command.NilDefaultStr, "The new IPMI hostname of the server. This command cannot be used in conjunction with config or pipe commands."),
				"ipmi_username":         c.FlagSet.String("ipmi-user", command.NilDefaultStr, "The new IPMI username of the server. This command cannot be used in conjunction with config or pipe commands."),
				"ipmi_password":         c.FlagSet.String("ipmi-pass", command.NilDefaultStr, "The new IPMI password of the server. This command cannot be used in conjunction with config or pipe commands."),
				"server_type":           c.FlagSet.String("server-type", command.NilDefaultStr, "The new server type (id or label) of the server. This command cannot be used in conjunction with config or pipe commands."),
				"server_class":          c.FlagSet.String("server-class", command.NilDefaultStr, "The new class of the server. This command cannot be used in conjunction with config or pipe commands."),
				"format":                c.FlagSet.String("format", "json", "The input format used when config or pipe commands are used. Supported values are 'json','yaml'. The default format is json."),
				"read_config_from_file": c.FlagSet.String("raw-config", command.NilDefaultStr, colors.Red("(Required)")+" Read raw object from file"),
				"read_config_from_pipe": c.FlagSet.Bool("pipe", false, colors.Green("(Flag)")+" If set, read raw object from pipe instead of from a file. Either this flag or the --raw-config option must be used."),
			}
		},
		ExecuteFunc: serverEditCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
	{
		Description:  "Edit server's IPMI",
		Subject:      "server",
		AltSubject:   "srv",
		Predicate:    "edit-ipmi",
		AltPredicate: "update-ipmi",
		FlagSet:      flag.NewFlagSet("edit server IPMI", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"server_id_or_uuid":  c.FlagSet.String("id", command.NilDefaultStr, "Server's ID or UUID"),
				"ipmi_hostname":      c.FlagSet.String("ipmi-host", command.NilDefaultStr, "The new IPMI hostname of the server."),
				"ipmi_username":      c.FlagSet.String("ipmi-user", command.NilDefaultStr, "The new IPMI username of the server."),
				"ipmi_password":      c.FlagSet.String("ipmi-pass", command.NilDefaultStr, "The new IPMI password of the server."),
				"ipmi_update_in_bmc": c.FlagSet.Bool("update-credentials-on-bmc", false, "If set, the server's BMC credentials on the actual server will also be updated."),
			}
		},
		ExecuteFunc: serverEditIPMICmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
	{
		Description:  "Change server power status",
		Subject:      "server",
		AltSubject:   "srv",
		Predicate:    "power-control",
		AltPredicate: "pwr",
		FlagSet:      flag.NewFlagSet("", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"server_id":   c.FlagSet.Int("id", command.NilDefaultInt, colors.Red("(Required)")+" Server's id."),
				"operation":   c.FlagSet.String("operation", command.NilDefaultStr, colors.Red("(Required)")+" Power control operation, one of: on, off, reset, soft."),
				"autoconfirm": c.FlagSet.Bool("autoconfirm", false, colors.Green("(Flag)")+" If set it will assume action is confirmed"),
			}
		},
		ExecuteFunc: serverPowerControlCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
	{
		Description:  "Change server status",
		Subject:      "server",
		AltSubject:   "srv",
		Predicate:    "status-set",
		AltPredicate: "status",
		FlagSet:      flag.NewFlagSet("", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"server_id":   c.FlagSet.Int("id", command.NilDefaultInt, colors.Red("(Required)")+" Server's id."),
				"status":      c.FlagSet.String("status", command.NilDefaultStr, colors.Red("(Required)")+" New server status. One of: 'available','decommissioned','removed_from_rack'"),
				"autoconfirm": c.FlagSet.Bool("autoconfirm", false, colors.Green("(Flag)")+" If set it will assume action is confirmed"),
			}
		},
		ExecuteFunc: serverStatusSetCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
	{
		Description:  "Reregister server",
		Subject:      "server",
		AltSubject:   "srv",
		Predicate:    "reregister",
		AltPredicate: "re-register",
		FlagSet:      flag.NewFlagSet("", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"server_id":   c.FlagSet.Int("id", command.NilDefaultInt, colors.Red("(Required)")+" Server's id."),
				"skip_ipmi":   c.FlagSet.Bool("do-not-set-ipmi", false, "If set, the system will not change the IPMI credentials."),
				"autoconfirm": c.FlagSet.Bool("autoconfirm", false, colors.Green("(Flag)")+" If set it will assume action is confirmed"),
			}
		},
		ExecuteFunc: serverReregisterCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
	{
		Description:  "Change server server type",
		Subject:      "server",
		AltSubject:   "srv",
		Predicate:    "server-type-set",
		AltPredicate: "server-type",
		FlagSet:      flag.NewFlagSet("", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"server_id":   c.FlagSet.Int("id", command.NilDefaultInt, colors.Red("(Required)")+" Server's id."),
				"server_type": c.FlagSet.String("server-type", command.NilDefaultStr, colors.Red("(Required)")+" New server type. Can be an ID or label"),
				"autoconfirm": c.FlagSet.Bool("autoconfirm", false, colors.Green("(Flag)")+" If set it will assume action is confirmed"),
			}
		},
		ExecuteFunc: serverServerTypeSetCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
	{
		Description:  "Change server rack information",
		Subject:      "server",
		AltSubject:   "srv",
		Predicate:    "rack-info-set",
		AltPredicate: "rack-info",
		FlagSet:      flag.NewFlagSet("", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"server_id":   c.FlagSet.Int("id", command.NilDefaultInt, colors.Red("(Required)")+" Server's id."),
				"rack_name":   c.FlagSet.String("rack-name", command.NilDefaultStr, colors.Red("(Required)")+" New rack name."),
				"lower_u":     c.FlagSet.Int("lower-u", command.NilDefaultInt, colors.Red("(Required)")+" Lower U of the equipment"),
				"upper_u":     c.FlagSet.Int("upper-u", command.NilDefaultInt, colors.Red("(Required)")+" Upper U of the equipment"),
				"autoconfirm": c.FlagSet.Bool("autoconfirm", false, colors.Green("(Flag)")+" If set it will assume action is confirmed"),
			}
		},
		ExecuteFunc: serverRackInfoSetCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
	{
		Description:  "Change server inventory information",
		Subject:      "server",
		AltSubject:   "srv",
		Predicate:    "inventory-info-set",
		AltPredicate: "inventory-info",
		FlagSet:      flag.NewFlagSet("", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"server_id":    c.FlagSet.Int("id", command.NilDefaultInt, colors.Red("(Required)")+" Server's id."),
				"inventory_id": c.FlagSet.String("inventory-id", command.NilDefaultStr, colors.Red("(Required)")+" New inventory id"),
				"autoconfirm":  c.FlagSet.Bool("autoconfirm", false, colors.Green("(Flag)")+" If set it will assume action is confirmed"),
			}
		},
		ExecuteFunc: serverInventoryInfoSetCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
	{
		Description:  "Lists server interfaces.",
		Subject:      "server",
		AltSubject:   "srv",
		Predicate:    "interfaces",
		AltPredicate: "intf",
		FlagSet:      flag.NewFlagSet("list server interfaces", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"format":            c.FlagSet.String("format", command.NilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
				"server_id_or_uuid": c.FlagSet.Int("id", command.NilDefaultInt, colors.Red("(Required)")+" Server's id."),
				"raw":               c.FlagSet.Bool("raw", false, colors.Green("(Flag)")+" When set the return will be a full dump of the object. This is useful when copying configurations. Only works with json and yaml formats."),
			}
		},
		ExecuteFunc: serverInterfacesListCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
}

func serverPowerControlCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	serverID, ok := command.GetIntParamOk(c.Arguments["server_id"])
	if !ok {
		return "", fmt.Errorf("-id is required")
	}
	operation, ok := command.GetStringParamOk(c.Arguments["operation"])
	if !ok {
		return "", fmt.Errorf("-operation is required (one of: on, off, reset, soft)")
	}

	server, err := client.ServerGet(serverID, false)
	if err != nil {
		return "", err
	}

	confirm, err := command.ConfirmCommand(c, func() string {
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

func serverStatusSetCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	serverID, ok := command.GetIntParamOk(c.Arguments["server_id"])
	if !ok {
		return "", fmt.Errorf("-id is required")
	}

	newStatus, ok := command.GetStringParamOk(c.Arguments["status"])
	if !ok {
		return "", fmt.Errorf("-status is required (one of: on, off, reset, soft)")
	}

	var server metalcloud.Server

	if !command.GetBoolParam(c.Arguments["autoconfirm"]) {
		serverPtr, err := client.ServerGet(serverID, false)
		if err != nil {
			return "", err
		}
		server = *serverPtr
	}

	confirm, err := command.ConfirmCommand(c, func() string {

		confirmationMessage := ""

		if !command.GetBoolParam(c.Arguments["autoconfirm"]) {

			confirmationMessage = fmt.Sprintf("Server #%s (%s) of datacenter %s. Current status: %s new status: %s  Are you sure? Type \"yes\" to continue:",
				colors.Blue(fmt.Sprintf("%d", server.ServerID)),
				colors.Yellow(server.ServerSerialNumber),
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

func serverServerTypeSetCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	serverID, ok := command.GetIntParamOk(c.Arguments["server_id"])
	if !ok {
		return "", fmt.Errorf("-id is required")
	}

	serverTypeStr, ok := command.GetStringParamOk(c.Arguments["server_type"])
	if !ok {
		return "", fmt.Errorf("-server-type is required")
	}

	serverTypeID, _, isID := command.IdOrLabel(serverTypeStr)
	var newServerType metalcloud.ServerType
	if !isID {
		st, err := client.ServerTypeGetByLabel(serverTypeStr)
		if err != nil {
			return "", err
		}
		newServerType = *st
	} else {
		st, err := client.ServerTypeGet(serverTypeID)
		if err != nil {
			return "", err
		}
		newServerType = *st
	}

	var server metalcloud.Server

	if !command.GetBoolParam(c.Arguments["autoconfirm"]) {
		serverPtr, err := client.ServerGet(serverID, false)
		if err != nil {
			return "", err
		}
		server = *serverPtr
	}

	confirm, err := command.ConfirmCommand(c, func() string {

		confirmationMessage := ""

		if !command.GetBoolParam(c.Arguments["autoconfirm"]) {

			oldServerType := metalcloud.ServerType{
				ServerTypeName: "none",
				ServerTypeID:   0,
			}

			if server.ServerTypeID != 0 {
				st, err := client.ServerTypeGet(server.ServerTypeID)
				if err != nil {
					return err.Error()
				}
				oldServerType = *st
			}

			confirmationMessage = fmt.Sprintf("Server #%s (%s) of datacenter %s. Current server type: %s (#%s) new server type: %s (#%s) Are you sure? Type \"yes\" to continue:",
				colors.Blue(fmt.Sprintf("%d", server.ServerID)),
				colors.Yellow(server.ServerSerialNumber),
				server.DatacenterName,
				colors.Red(oldServerType.ServerTypeName),
				colors.Red(oldServerType.ServerTypeID),
				colors.Green(newServerType.ServerTypeName),
				colors.Green(newServerType.ServerTypeID),
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

func serverRackInfoSetCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	serverID, ok := command.GetIntParamOk(c.Arguments["server_id"])
	if !ok {
		return "", fmt.Errorf("-id is required")
	}

	serverRackName, ok := command.GetStringParamOk(c.Arguments["rack_name"])
	if !ok {
		return "", fmt.Errorf("-rack-name is required")
	}

	serverRackLowerU, ok := command.GetIntParamOk(c.Arguments["lower_u"])
	if !ok {
		return "", fmt.Errorf("-lower-u is required")
	}

	serverRackUpperU, ok := command.GetIntParamOk(c.Arguments["upper_u"])
	if !ok {
		return "", fmt.Errorf("-upper-u is required")
	}

	var server metalcloud.Server

	serverPtr, err := client.ServerGet(serverID, false)
	if err != nil {
		return "", err
	}
	server = *serverPtr

	confirm, err := command.ConfirmCommand(c, func() string {

		confirmationMessage := ""

		if !command.GetBoolParam(c.Arguments["autoconfirm"]) {

			oldRackInfo := getRackInfoSafe(server)

			oldServerRackInfo := fmt.Sprintf("Rack:%s U:%s-%s", oldRackInfo.RackName, oldRackInfo.LowerU, oldRackInfo.UpperU)

			newServerRackInfo := fmt.Sprintf("Rack:%s U:%d-%d", serverRackName, serverRackLowerU, serverRackUpperU)

			confirmationMessage = fmt.Sprintf("Server #%s (%s) of datacenter %s. Current server rack info %s new rack info: %s. Are you sure? Type \"yes\" to continue:",
				colors.Blue(fmt.Sprintf("%d", server.ServerID)),
				colors.Yellow(server.ServerSerialNumber),
				server.DatacenterName,
				colors.Red(oldServerRackInfo),
				colors.Green(newServerRackInfo),
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

		lowerUStr := fmt.Sprintf("%d", serverRackLowerU)
		upperUStr := fmt.Sprintf("%d", serverRackUpperU)

		serverRackEdit := metalcloud.ServerEditRack{
			ServerRackName:              &serverRackName,
			ServerRackPositionLowerUnit: &lowerUStr,
			ServerRackPositionUpperUnit: &upperUStr,
		}

		_, err = client.ServerEditRack(serverID, serverRackEdit)
	}

	return "", err
}

func serverInventoryInfoSetCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	serverID, ok := command.GetIntParamOk(c.Arguments["server_id"])
	if !ok {
		return "", fmt.Errorf("-id is required")
	}

	inventoryID, ok := command.GetStringParamOk(c.Arguments["inventory_id"])
	if !ok {
		return "", fmt.Errorf("inventory-id is required")
	}

	var server metalcloud.Server

	serverPtr, err := client.ServerGet(serverID, false)
	if err != nil {
		return "", err
	}
	server = *serverPtr

	confirm, err := command.ConfirmCommand(c, func() string {

		confirmationMessage := ""

		if !command.GetBoolParam(c.Arguments["autoconfirm"]) {

			oldInventoryID := getStringFromStringOrEmpty(server.ServerInventoryId)

			confirmationMessage = fmt.Sprintf("Server #%s (%s) of datacenter %s. Current inventory id: %s new inventory id: %s. Are you sure? Type \"yes\" to continue:",
				colors.Blue(fmt.Sprintf("%d", server.ServerID)),
				colors.Yellow(server.ServerSerialNumber),
				server.DatacenterName,
				colors.Red(oldInventoryID),
				colors.Green(inventoryID),
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

		serverEditInventory := metalcloud.ServerEditInventory{
			ServerInventoryId: &inventoryID,
		}

		_, err = client.ServerEditInventory(serverID, serverEditInventory)
	}

	return "", err
}

func serverReregisterCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	serverID, ok := command.GetIntParamOk(c.Arguments["server_id"])
	if !ok {
		return "", fmt.Errorf("-id is required")
	}

	skipIpmi := command.GetBoolParam(c.Arguments["skip_ipmi"])

	var server metalcloud.Server

	if !command.GetBoolParam(c.Arguments["autoconfirm"]) {
		serverPtr, err := client.ServerGet(serverID, false)
		if err != nil {
			return "", err
		}
		server = *serverPtr
	}

	confirm, err := command.ConfirmCommand(c, func() string {

		confirmationMessage := ""

		if !command.GetBoolParam(c.Arguments["autoconfirm"]) {

			confirmationMessage = fmt.Sprintf("Server #%s (%s) BMC IP:%s of datacenter %s. Are you sure? Type \"yes\" to continue:",
				colors.Blue(fmt.Sprintf("%d", server.ServerID)),
				colors.Yellow(server.ServerSerialNumber),
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
		return colors.Blue(status)
	case "used":
		return colors.Green(status)
	case "unavailable":
		return colors.Magenta(status)
	}
	return colors.Yellow(status)

}

func serversListCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	filter := command.GetStringParam(c.Arguments["filter"])

	list, err := client.ServersSearch(filtering.ConvertToSearchFieldFormat(filter))
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

	if command.GetBoolParam(c.Arguments["show_rack_info"]) {

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

	if command.GetBoolParam(c.Arguments["show_hardware"]) {
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

	if command.GetBoolParam(c.Arguments["show_credentials"]) {

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

	statusCounts := map[string]int{
		"available":      0,
		"cleaning":       0,
		"registering":    0,
		"used":           0,
		"decommissioned": 0,
	}

	for _, s := range *list {

		if s.ServerStatus == "decommissioned" && !command.GetBoolParam(c.Arguments["show_decommissioned"]) {
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
				allocation = stringutils.TruncateString(allocation, 10)
			}
		}

		credentialsUser := ""
		credentialsPass := ""
		//snmpCommunity := ""

		if command.GetBoolParam(c.Arguments["show_credentials"]) {

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
				colors.Yellow(s.ServerDiskCount),
				colors.Yellow(s.ServerDiskSizeMbytes/1000),
				colors.Yellow(s.ServerDiskType))
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
					colors.Magenta(len(serverInterfaces)),
					colors.Magenta(capacity/1000))
		}

		hardwareConfig := ""
		if s.ServerProcessorCount > 0 {
			hardwareConfig = fmt.Sprintf("%s x %s cores(ht), %s GB RAM, %s, %s",
				colors.Blue(s.ServerProcessorCount),
				colors.Blue(s.ServerProcessorCoreCount),
				colors.Red(s.ServerRAMGbytes),
				diskDescription,
				interfaceDescription,
			)
		}

		status := s.ServerStatus

		switch status {
		case "available":
			status = colors.Blue(status)
		case "used":
			status = colors.Green(status)
		case "unavailable":
			status = colors.Magenta(status)
		case "defective":
			status = colors.Red(status)
		default:
			status = colors.Yellow(status)

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

		if command.GetBoolParam(c.Arguments["show_rack_info"]) {

			row = append(row, []interface{}{
				strings.Join(s.ServerTags, ","),
				s.ServerInventoryId,
				s.ServerRackName,
				s.ServerRackPositionLowerUnit,
				s.ServerRackPositionUpperUnit,
			}...)
		}

		if command.GetBoolParam(c.Arguments["show_hardware"]) {
			row = append(row, []interface{}{
				hardwareConfig,
			}...)
		}

		if command.GetBoolParam(c.Arguments["show_credentials"]) {
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

	if command.GetBoolParam(c.Arguments["show_decommissioned"]) {
		title = title + fmt.Sprintf(" %d decommissioned", statusCounts["decommissioned"])
	}

	return table.RenderTable(title, "", command.GetStringParam(c.Arguments["format"]))
}

func serverGetCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	showCredentials := command.GetBoolParam(c.Arguments["show_credentials"])

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
		productName = stringutils.TruncateString(server.ServerProductName, 18)
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

	format := command.GetStringParam(c.Arguments["format"])

	if command.GetBoolParam(c.Arguments["raw"]) {
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

func serverCreateCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	var obj metalcloud.Server

	err := command.GetRawObjectFromCommand(c, &obj)

	ret, err := client.ServerCreate(obj, false)
	if err != nil {
		return "", err
	}

	if command.GetBoolParam(c.Arguments["return_id"]) {
		return fmt.Sprintf("%d", ret), nil
	}

	return "", err
}

func serverRegisterCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	datacenter, ok := command.GetStringParamOk(c.Arguments["datacenter"])
	if !ok {
		return "", fmt.Errorf("-datacenter is required")
	}
	server_vendor, ok := command.GetStringParamOk(c.Arguments["server_vendor"])
	if !ok {
		return "", fmt.Errorf("-server_vendor is required")
	}
	mgmt_address, ok := command.GetStringParamOk(c.Arguments["mgmt_address"])
	if !ok {
		return "", fmt.Errorf("-mgmt_address is required")
	}
	mgmt_user, ok := command.GetStringParamOk(c.Arguments["mgmt_user"])
	if !ok {
		return "", fmt.Errorf("-mgmt_user is required")
	}
	mgmt_pass, ok := command.GetStringParamOk(c.Arguments["mgmt_pass"])
	if !ok {
		return "", fmt.Errorf("-mgmt_pass is required")
	}

	obj := metalcloud.ServerCreateAndRegister{
		DatacenterName:           datacenter,
		ServerVendor:             server_vendor,
		ServerManagementAddress:  mgmt_address,
		ServerManagementUser:     mgmt_user,
		ServerManagementPassword: mgmt_pass,
	}

	ret, err := client.ServerCreateAndRegister(obj)
	if err != nil {
		return "", err
	}

	if command.GetBoolParam(c.Arguments["return_id"]) {
		return fmt.Sprintf("%d", ret), nil
	}

	return "", err
}

type ServerCreateUnmanagedInternal struct {
	ServerCreateUnmanaged metalcloud.ServerCreateUnmanaged `json:",inline" yaml:",inline"`

	//not used serverside but used by the CLI
	InstanceArrayLabel  *string `json:"instance_array_label,omitempty" yaml:"label,omitempty"`
	ServerTypeLabel     *string `json:"server_type_label,omitempty" yaml:"serverType,omitempty"`
	InfrastructureLabel *string `json:"infrastructure_label,omitempty" yaml:"infrastructure,omitempty"`
	InfrastructureID    *int    `json:"infrastructure_id,omitempty" yaml:"infrastructureID,omitempty"`
	UserEmail           *string `json:"user_email,omitempty" yaml:"userEmail,omitempty"`
	UserID              *int    `json:"user_id,omitempty" yaml:"userID,omitempty"`
}

func serverImportCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	var obj ServerCreateUnmanagedInternal

	err := command.GetRawObjectFromCommand(c, &obj)
	if err != nil {
		return "", err
	}

	if obj.ServerCreateUnmanaged.ServerSerialNumber == "" {
		return "", fmt.Errorf("server serial number cannot be empty")
	}

	if obj.ServerTypeLabel != nil {
		serverType, err := client.ServerTypeGetByLabel(*obj.ServerTypeLabel)
		if err != nil {
			return "", err
		}
		obj.ServerCreateUnmanaged.ServerTypeID = serverType.ServerTypeID
	}

	createdServer, err := client.ServerUnmanagedImport(obj.ServerCreateUnmanaged)
	if err != nil {
		return "", err
	}

	if v, ok := command.GetStringParamOk(c.Arguments["add_to_infra"]); ok {
		_, err := addServerToInfrastructure(createdServer.ServerID, &v, &obj, client)
		if err != nil {
			return "", err
		}
	}

	if command.GetBoolParam(c.Arguments["return_id"]) {
		return fmt.Sprintf("%d", createdServer.ServerID), nil
	}

	return "", err
}

// returns an infrastructure id as defined in the infrastructureID, InfrastructureLabel, UserEmail, UserID
// fields in the obj param
// if the infrastructure does not exist it will be created
// if no name is provided for the infrastructure a default one will be used
// if no user is specified the system will fallback to the user that is logged in (to whom the api key belongs)
func createInfrastructureIfDoesNotExist(obj *ServerCreateUnmanagedInternal, client metalcloud.MetalCloudClient) (int, error) {
	if obj.UserID != nil {
		user, err := client.UserGet(*obj.UserID)
		if err != nil {
			return 0, err
		}
		obj.UserEmail = &user.UserEmail
	}

	if obj.InfrastructureID != nil {
		infra, err := client.InfrastructureGet(*obj.InfrastructureID)
		if err != nil {
			return 0, err
		}
		obj.InfrastructureLabel = &infra.InfrastructureLabel
	}

	//if user not provided we use the current user
	if obj.UserEmail == nil {
		currentUserEmail := client.GetUserEmail()
		obj.UserEmail = &currentUserEmail
	}

	//if infrastructure not provided we use a hardcoded one
	if obj.InfrastructureLabel == nil {
		defaultInfrastructureName := "imported"
		obj.InfrastructureLabel = &defaultInfrastructureName
	}

	if obj.InstanceArrayLabel == nil {
		defaultInstanceArrayLabel := obj.ServerCreateUnmanaged.ServerSerialNumber
		obj.InstanceArrayLabel = &defaultInstanceArrayLabel
	}

	user, err := client.UserGetByEmail(*obj.UserEmail)
	if err != nil {
		return 0, err
	}

	// search for infrastructures with the given name that belong to the user
	foundInfrastructures, err := client.InfrastructureSearch(fmt.Sprintf("user_email:%s infrastructure_label:%s", *obj.UserEmail, *obj.InfrastructureLabel))
	if err != nil {
		return 0, err
	}

	var infrastructureID int

	if len(*foundInfrastructures) == 0 {
		// infra with name does not exist for user create it
		infrastructure := metalcloud.Infrastructure{
			InfrastructureLabel: *obj.InfrastructureLabel,
			UserIDowner:         user.UserID,
			DatacenterName:      obj.ServerCreateUnmanaged.DatacenterName,
		}
		createdInfra, err := client.InfrastructureCreate(infrastructure)
		if err != nil {
			return 0, err
		}
		infrastructureID = createdInfra.InfrastructureID
	} else {
		infrastructureID = (*foundInfrastructures)[0].InfrastructureID
	}

	return infrastructureID, nil
}

// addServerToInfrastructure adds server to an infrastructure by creating an instance array, adding it to an infrastructure
// the instance array will be called like the serial number unless overwritten by the InstanceArrayLabel entry
// provide an infrastructure id or label to the infrastructureIDOrLabel to use that infrastructure
// provide null or "auto" to automatically create the infrastructure based on the object details provided in the object
func addServerToInfrastructure(serverID int, infrastructureIDOrLabel *string, obj *ServerCreateUnmanagedInternal, client metalcloud.MetalCloudClient) (*metalcloud.InstanceArray, error) {

	var err error
	var infrastructureID int

	if infrastructureIDOrLabel != nil && *infrastructureIDOrLabel != "auto" {
		id, label, isID := command.IdOrLabelString(*infrastructureIDOrLabel)
		if !isID {
			log.Printf("infra label: %v", label)
			infra, err := client.InfrastructureGetByLabel(label)
			if err != nil {
				return nil, err
			}
			infrastructureID = infra.InfrastructureID
		} else {
			infrastructureID = id
		}

	} else {
		infrastructureID, err = createInfrastructureIfDoesNotExist(obj, client)
		if err != nil {
			return nil, err
		}
	}

	server, err := client.ServerGet(serverID, false)
	if err != nil {
		return nil, err
	}

	if server.ServerStatus == "decommissioned" {
		return nil, fmt.Errorf("The specified server is decomissioned")
	}

	serverTypeID := server.ServerTypeID

	instanceArrayLabel := server.ServerSerialNumber

	if obj != nil && *obj.InstanceArrayLabel != "" {
		instanceArrayLabel = *obj.InstanceArrayLabel
	}

	createdIA, err := createInstanceArrayWithOptions(infrastructureID, instanceArrayLabel, serverTypeID, serverID, 1, client)
	if err != nil {
		return nil, err
	}

	return createdIA, nil
}

func createInstanceArrayWithOptions(infrastructureID int, instanceArrayLabel string, serverTypeID int, serverID int, instanceCount int, client metalcloud.MetalCloudClient) (*metalcloud.InstanceArray, error) {
	// create instance array for the server that we just imported
	ia := metalcloud.InstanceArray{
		InstanceArrayLabel:         instanceArrayLabel,
		InstanceArrayInstanceCount: instanceCount,
	}

	createdIA, err := client.InstanceArrayCreate(infrastructureID, ia)
	if err != nil {
		return nil, err
	}

	instances, err := client.InstanceArrayInstances(createdIA.InstanceArrayID)
	if err != nil {
		return nil, err
	}

	for _, i := range *instances {

		i.InstanceOperation.ServerTypeID = serverTypeID
		i.InstanceOperation.PreferredServerIDsJSON = fmt.Sprintf("[%d]", serverID)

		_, err := client.InstanceEdit(i.InstanceID, i.InstanceOperation)

		if err != nil {
			return nil, err
		}
	}

	return createdIA, nil
}

func serverAddToInfraCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	serverID, ok := command.GetIntParamOk(c.Arguments["server_id"])
	if !ok {
		return "", fmt.Errorf("-server-id is required")
	}

	infrastructureIDOrLabel, ok := command.GetStringParamOk(c.Arguments["infrastructure_id_or_label"])
	if !ok {
		return "", fmt.Errorf("-infra is required")
	}

	retIA, err := addServerToInfrastructure(serverID, &infrastructureIDOrLabel, nil, client)
	if err != nil {
		return "", err
	}

	if command.GetBoolParam(c.Arguments["return_id"]) {
		return fmt.Sprintf("%d", retIA.InstanceArrayID), nil
	}

	return "", err
}

func getMultipleServerCreateUnmanagedInternalFromYamlFile(filePath string) ([]ServerCreateUnmanagedInternal, error) {

	file, err := os.Open(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []ServerCreateUnmanagedInternal{}, nil
		} else {
			return []ServerCreateUnmanagedInternal{}, err
		}
	}

	decoder := yaml.NewDecoder(file)

	records := []ServerCreateUnmanagedInternal{}

	for true {

		var record ServerCreateUnmanagedInternal

		err = decoder.Decode(&record)
		if err == nil {
			records = append(records, record)
		}

		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			} else {
				return nil, fmt.Errorf("Error while reading %s: %v", filePath, err)
			}
		}
	}

	file.Close()

	return records, nil
}

func serverImportBatchCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	filePath, ok := command.GetStringParamOk(c.Arguments["read_config_from_file"])
	if !ok {
		return "", fmt.Errorf("-file is required")
	}

	records, err := getMultipleServerCreateUnmanagedInternalFromYamlFile(filePath)
	if err != nil {
		return "", err
	}

	//set server type ids if they are set as labels
	for i, r := range records {

		if r.ServerTypeLabel != nil {
			serverType, err := client.ServerTypeGetByLabel(*r.ServerTypeLabel)
			if err != nil {
				return "", err
			}
			records[i].ServerCreateUnmanaged.ServerTypeID = serverType.ServerTypeID
		}
	}

	//perform a batch update. This helps perform interface swaps in one go
	embeddedObjects := []metalcloud.ServerCreateUnmanaged{}
	for _, o := range records {
		embeddedObjects = append(embeddedObjects, o.ServerCreateUnmanaged)
	}

	createdServerRecords, err := client.ServerUnmanagedImportBatch(embeddedObjects)
	if err != nil {
		return "", err
	}

	if v, ok := command.GetStringParamOk(c.Arguments["add_to_infra"]); ok {
		for _, record := range records {

			//because the order might have changed
			//find the server creation object in records for the
			//returned object for the same serial number
			serverID := 0
			for _, cr := range *createdServerRecords {
				if strings.ToLower(cr.ServerSerialNumber) == strings.ToLower(record.ServerCreateUnmanaged.ServerSerialNumber) {
					serverID = cr.ServerID
					break
				}
			}
			if serverID != 0 {
				_, err := addServerToInfrastructure(serverID, &v, &record, client)
				if err != nil {
					return "", err
				}
			}
		}
	}

	var s strings.Builder

	if command.GetBoolParam(c.Arguments["return_id"]) {
		for _, r := range *createdServerRecords {
			s.WriteString(fmt.Sprintf("%d\n", r.ServerID))
		}

		return s.String(), nil
	}

	return "", err
}

func serverEditCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	server, err := getServerFromCommand("id", c, client, false)
	if err != nil {
		return "", err
	}

	newStatus, setStatus := command.GetStringParamOk(c.Arguments["status"])
	newIPMIHostname, setIPMIHostname := command.GetStringParamOk(c.Arguments["ipmi_hostname"])
	newIPMIUsername, setIPMIUsername := command.GetStringParamOk(c.Arguments["ipmi_username"])
	newIPMIPassword, setIPMIPassword := command.GetStringParamOk(c.Arguments["ipmi_password"])

	_, setServerType := command.GetStringParamOk(c.Arguments["server_type"])
	newServerClass, setServerClass := command.GetStringParamOk(c.Arguments["server_class"])

	_, readFromFile := command.GetStringParamOk(c.Arguments["read_config_from_file"])
	readFromPipe := command.GetBoolParam(c.Arguments["read_config_from_pipe"])

	if (readFromFile || readFromPipe) && (setStatus || setIPMIHostname || setIPMIUsername || setIPMIPassword) {
		return "", fmt.Errorf("Cannot use --config or --pipe with --status or --ipmi-host or --ipmi-user or --ipmi-pass")
	}

	newServer := *server

	if readFromFile || readFromPipe {

		err = command.GetRawObjectFromCommand(c, &newServer)
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

func serverEditIPMICmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	server, err := getServerFromCommand("id", c, client, false)
	if err != nil {
		return "", err
	}

	newIPMIHostname, setIPMIHostname := command.GetStringParamOk(c.Arguments["ipmi_hostname"])
	newIPMIUsername, setIPMIUsername := command.GetStringParamOk(c.Arguments["ipmi_username"])
	newIPMIPassword, setIPMIPassword := command.GetStringParamOk(c.Arguments["ipmi_password"])
	IPMIUpdateInBMC := command.GetBoolParam(c.Arguments["ipmi_update_in_bmc"])

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

	_, err = client.ServerEditIPMI(server.ServerID, newServer, IPMIUpdateInBMC)

	return "", err
}

func getServerFromCommand(paramName string, c *command.Command, client metalcloud.MetalCloudClient, decryptPassword bool) (*metalcloud.Server, error) {

	m, err := command.GetParam(c, "server_id_or_uuid", paramName)
	if err != nil {
		return nil, err
	}

	id, uuid, isID := command.IdOrLabel(m)

	if isID {
		return client.ServerGet(id, decryptPassword)
	}

	return client.ServerGetByUUID(uuid, decryptPassword)
}

func getServerTypeFromCommand(paramName string, c *command.Command, client metalcloud.MetalCloudClient) (*metalcloud.ServerType, error) {

	m, err := command.GetParam(c, "server_type", paramName)
	if err != nil {
		return nil, err
	}

	id, label, isID := command.IdOrLabel(m)

	if isID {
		return client.ServerTypeGet(id)
	}

	return client.ServerTypeGetByLabel(label)
}

func serverInterfacesListCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

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
			FieldName: "SRV. ID",
			FieldType: tableformatter.TypeInt,
			FieldSize: 5,
		},
		{
			FieldName: "INTF. IDX",
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

		ips := stringutils.FlattenAndJoinStrings(s.IP)
		networkType := strings.Join(s.NetworkType, ",")

		switch_info := fmt.Sprintf("%s (#%d)",
			s.NetworkEquipmentIdentifierString,
			s.NetworkEquipmentID,
		)

		capacity := fmt.Sprintf("%d Gbps", int(s.ServerInterfaceCapacityMBPs/1000))

		data = append(data, []interface{}{
			server.ServerID,
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

	format := command.GetStringParam(c.Arguments["format"])

	if command.GetBoolParam(c.Arguments["raw"]) {
		for _, s := range *list {
			ret, err := tableformatter.RenderRawObject(s, format, "Server interfaces")
			if err != nil {
				return "", err
			}
			sb.WriteString(ret)
		}
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
		RackName:    getStringFromStringOrEmpty(server.ServerRackName),
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
