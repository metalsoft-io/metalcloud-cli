package instance

import (
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

var InstanceCmds = []command.Command{
	{
		Description:  "Control power for an instance",
		Subject:      "instance",
		AltSubject:   "instance",
		Predicate:    "power-control",
		AltPredicate: "pwr",
		FlagSet:      flag.NewFlagSet("instance_array", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"instance_id": c.FlagSet.Int("id", command.NilDefaultInt, colors.Red("(Required)")+" Instances's id . Note that the 'label' this be ambiguous in certain situations."),
				"operation":   c.FlagSet.String("operation", command.NilDefaultStr, colors.Red("(Required)")+" Power control operation, one of: on, off, reset, soft"),
				"autoconfirm": c.FlagSet.Bool("autoconfirm", false, colors.Green("(Flag)")+" If set it will assume action is confirmed"),
			}
		},
		ExecuteFunc: instancePowerControlCmd,
	},

	{
		Description:  "Show an instance details",
		Subject:      "instance",
		AltSubject:   "instance",
		Predicate:    "get",
		AltPredicate: "show",
		FlagSet:      flag.NewFlagSet("instance get", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"instance_id":           c.FlagSet.Int("id", command.NilDefaultInt, colors.Red("(Required)")+" Instances's id . Note that the 'label' this be ambiguous in certain situations."),
				"show_credentials":      c.FlagSet.Bool("show-credentials", false, colors.Green("(Flag)")+" If set returns the instance's credentials"),
				"show_custom_variables": c.FlagSet.Bool("show-custom-variables", false, colors.Green("(Flag)")+" If set returns the instance's custom variables"),
				"format":                c.FlagSet.String("format", "yaml", "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
			}
		},
		ExecuteFunc: instanceGetCmd,
	},
	{
		Description:  "Replace an instance's associated server. (EXPERIMENTAL)",
		Subject:      "instance",
		AltSubject:   "instance",
		Predicate:    "server-replace",
		AltPredicate: "server-change",
		FlagSet:      flag.NewFlagSet("", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"instance_id":   c.FlagSet.Int("id", command.NilDefaultInt, colors.Red("(Required)")+" Instance's id."),
				"server_id":     c.FlagSet.Int("new-server-id", command.NilDefaultInt, colors.Red("(Required)")+" New server's id."),
				"autoconfirm":   c.FlagSet.Bool("autoconfirm", false, colors.Green("(Flag)")+" If set it will assume action is confirmed"),
				"return_afc_id": c.FlagSet.Bool("return-afc-id", false, colors.Green("(Flag)")+" If set it will return the AFC id of the operation."),
			}
		},
		ExecuteFunc: instanceServerReplaceCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
	{
		Description:  "Edit instance",
		Subject:      "instance",
		AltSubject:   "instance",
		Predicate:    "edit",
		AltPredicate: "update",
		FlagSet:      flag.NewFlagSet("", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"instance_id":      c.FlagSet.Int("id", command.NilDefaultInt, colors.Red("(Required)")+" Instance's id."),
				"custom_variables": c.FlagSet.String("custom-variables", command.NilDefaultStr, "Comma separated list of custom variables such as 'var1=value,var2=value'. If special characters need to be set use urlencode and pass the encoded string"),
			}
		},
		ExecuteFunc: instanceEditCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
}

func instancePowerControlCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	instanceID, ok := command.GetIntParamOk(c.Arguments["instance_id"])
	if !ok {
		return "", fmt.Errorf("-id is required (drive id)")
	}
	operation, ok := command.GetStringParamOk(c.Arguments["operation"])
	if !ok {
		return "", fmt.Errorf("-operation is required (one of: on, off, reset, soft)")
	}

	instance, err := client.InstanceGet(instanceID)
	if err != nil {
		return "", err
	}

	ia, err := client.InstanceArrayGet(instance.InstanceArrayID)
	if err != nil {
		return "", err
	}

	infra, err := client.InfrastructureGet(ia.InfrastructureID)
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

		confirmationMessage := fmt.Sprintf("%s instance %s (%d) of instance array %s (#%d) infrastructure %s (#%d).  Are you sure? Type \"yes\" to continue:",
			op,
			instance.InstanceLabel,
			instance.InstanceID,
			ia.InstanceArrayLabel,
			ia.InstanceArrayID,
			infra.InfrastructureLabel,
			infra.InfrastructureID,
		)

		//this is simply so that we don't output a text on the command line under go test
		if strings.HasSuffix(os.Args[0], ".test") {
			confirmationMessage = ""
		}

		return confirmationMessage

	})

	if err != nil {
		return "", err
	}

	if confirm {
		err = client.InstanceServerPowerSet(instanceID, operation)
	}

	return "", err
}

func instanceEditCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	instanceID, ok := command.GetIntParamOk(c.Arguments["instance_id"])
	if !ok {
		return "", fmt.Errorf("-id is required (drive id)")
	}

	instance, err := client.InstanceGet(instanceID)
	if err != nil {
		return "", err
	}

	if v, ok := command.GetStringParamOk(c.Arguments["custom_variables"]); ok {
		m, err := command.GetKeyValueMapFromString(v)
		if err != nil {
			return "", err
		}
		instance.InstanceOperation.InstanceCustomVariables = m
	}

	client.InstanceEdit(instanceID, instance.InstanceOperation)

	return "", err
}

func instanceGetCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	instanceID, ok := command.GetIntParamOk(c.Arguments["instance_id"])
	if !ok {
		return "", fmt.Errorf("-id is required (instance id)")
	}

	instance, err := client.InstanceGet(instanceID)
	if err != nil {
		return "", err
	}

	ia, err := client.InstanceArrayGet(instance.InstanceArrayID)
	if err != nil {
		return "", err
	}

	infra, err := client.InfrastructureGet(ia.InfrastructureID)
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
			FieldName: "SUBDOMAIN",
			FieldType: tableformatter.TypeString,
			FieldSize: 10,
		},
		{
			FieldName: "INSTANCE_ARRAY",
			FieldType: tableformatter.TypeString,
			FieldSize: 10,
		},
		{
			FieldName: "INFRASTRUCTURE",
			FieldType: tableformatter.TypeString,
			FieldSize: 10,
		},
		{
			FieldName: "PUBLIC_IPs",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "PRIVATE_IPs",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
	}

	publicIPS := getIPsAsStringArray(instance.InstanceCredentials.IPAddressesPublic)
	privateIPS := getIPsAsStringArray(instance.InstanceCredentials.IPAddressesPrivate)

	dataRow := []interface{}{
		instance.InstanceID,
		instance.InstanceSubdomainPermanent,
		ia.InstanceArrayLabel,
		infra.InfrastructureLabel,
		strings.Join(publicIPS, " "),
		strings.Join(privateIPS, " "),
	}

	if command.GetBoolParam(c.Arguments["show_credentials"]) {

		if v := instance.InstanceCredentials.SSH; v != nil {

			newFields := []tableformatter.SchemaField{
				{
					FieldName: "SSH_USERNAME",
					FieldType: tableformatter.TypeString,
					FieldSize: 10,
				},
				{
					FieldName: "SSH_PASSWORD",
					FieldType: tableformatter.TypeString,
					FieldSize: 10,
				},
				{
					FieldName: "SSH_PORT",
					FieldType: tableformatter.TypeInt,
					FieldSize: 10,
				},
			}

			schema = append(schema, newFields...)

			newData := []interface{}{
				v.Username,
				v.InitialPassword,
				v.Port,
			}
			dataRow = append(dataRow, newData...)
		}

		if v := instance.InstanceCredentials.RDP; v != nil {

			newFields := []tableformatter.SchemaField{
				{
					FieldName: "RDP_USERNAME",
					FieldType: tableformatter.TypeString,
					FieldSize: 5,
				},
				{
					FieldName: "RDP_PASSWORD",
					FieldType: tableformatter.TypeString,
					FieldSize: 5,
				},
				{
					FieldName: "RDP_PORT",
					FieldType: tableformatter.TypeInt,
					FieldSize: 5,
				},
			}

			schema = append(schema, newFields...)
			newData := []interface{}{
				v.Username,
				v.InitialPassword,
				v.Port,
			}
			dataRow = append(dataRow, newData...)
		}

		if v := instance.InstanceCredentials.ISCSI; v != nil {

			newFields := []tableformatter.SchemaField{
				{
					FieldName: "INITIATOR_IQN",
					FieldType: tableformatter.TypeString,
					FieldSize: 5,
				},
				{
					FieldName: "ISCSI_USERNAME",
					FieldType: tableformatter.TypeString,
					FieldSize: 5,
				},
				{
					FieldName: "ISCSI_PASSWORD",
					FieldType: tableformatter.TypeString,
					FieldSize: 5,
				},
			}

			schema = append(schema, newFields...)
			newData := []interface{}{
				v.InitiatorIQN,
				v.Username,
				v.Password,
			}
			dataRow = append(dataRow, newData...)
		}

		if v := instance.InstanceCredentials.SharedDrives; v != nil {

			for k, sd := range v {
				newFields := []tableformatter.SchemaField{
					{
						FieldName: fmt.Sprintf("SHARED_DRIVE_%s_TARGET_IP_ADDRESS", k),
						FieldType: tableformatter.TypeString,
						FieldSize: 5,
					},
					{
						FieldName: fmt.Sprintf("SHARED_DRIVE_%s_TARGET_PORT", k),
						FieldType: tableformatter.TypeInt,
						FieldSize: 5,
					},
					{
						FieldName: fmt.Sprintf("SHARED_DRIVE_%s_TARGET_IQN", k),
						FieldType: tableformatter.TypeString,
						FieldSize: 5,
					},
					{
						FieldName: fmt.Sprintf("SHARED_DRIVE_%s_LUN_ID", k),
						FieldType: tableformatter.TypeString,
						FieldSize: 5,
					},
				}

				schema = append(schema, newFields...)
				newData := []interface{}{
					sd.StorageIPAddress,
					sd.StoragePort,
					sd.TargetIQN,
					sd.LunID,
				}
				dataRow = append(dataRow, newData...)
			}
		}
	}

	if command.GetBoolParam(c.Arguments["show_custom_variables"]) {
		schema = append(schema, tableformatter.SchemaField{
			FieldName: "CUSTOM_VARS",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		})

		custom_vars := command.GetKeyValueStringFromMap(instance.InstanceCustomVariables)
		dataRow = append(dataRow, custom_vars)

	}

	data := [][]interface{}{dataRow}

	topRow := fmt.Sprintf("Instance %s",
		instance.InstanceSubdomainPermanent,
	)
	table := tableformatter.Table{
		Data:   data,
		Schema: schema,
	}
	return table.RenderTransposedTable("Records", topRow, command.GetStringParam(c.Arguments["format"]))
}

func getIPsAsStringArray(ips []metalcloud.IP) []string {
	sList := []string{}
	for _, ip := range ips {
		sList = append(sList, ip.IPHumanReadable)
	}
	return sList
}

func instanceServerReplaceCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	instanceID, ok := command.GetIntParamOk(c.Arguments["instance_id"])
	if !ok {
		return "", fmt.Errorf("-id is required")
	}

	instance, err := client.InstanceGet(instanceID)
	if err != nil {
		return "", err
	}

	instanceArray, err := client.InstanceArrayGet(instance.InstanceArrayID)
	if err != nil {
		return "", err
	}

	infrastructure, err := client.InfrastructureGet(instanceArray.InfrastructureID)
	if err != nil {
		return "", err
	}

	newServerID, ok := command.GetIntParamOk(c.Arguments["server_id"])
	if !ok {
		return "", fmt.Errorf("-new-server-id is required")
	}

	server, err := client.ServerGet(newServerID, false)
	if err != nil {
		return "", err
	}

	confirm, err := command.ConfirmCommand(c, func() string {

		confirmationMessage := ""

		if !command.GetBoolParam(c.Arguments["autoconfirm"]) {

			confirmationMessage = fmt.Sprintf("%s\nInstance #%s of instance array #%s, infrastructure #%s belonging to user %s will "+
				"have the associated server replaced with the server #%s (SN:%s) MGMT IP:%s on datacenter %s. \nAre you sure? Type \"yes\" to continue:",
				colors.Red("WARNING: This feature is experimental."),
				colors.Red(fmt.Sprintf("%d", instance.InstanceID)),
				fmt.Sprintf("%d", instanceArray.InstanceArrayID),
				fmt.Sprintf("%d", infrastructure.InfrastructureID),
				colors.Yellow(fmt.Sprintf("%s", infrastructure.UserEmailOwner)),
				colors.Red(fmt.Sprintf("%d", server.ServerID)),
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

	afc := 0
	if confirm {
		afc, err = client.InstanceServerReplace(instanceID, newServerID)
	}

	if command.GetBoolParam(c.Arguments["return_afc_id"]) {
		return fmt.Sprintf("%d", afc), nil
	}

	return "", err
}
