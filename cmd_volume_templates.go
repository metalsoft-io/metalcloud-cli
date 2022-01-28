package main

import (
	"flag"
	"fmt"
	"strings"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
	"github.com/metalsoft-io/tableformatter"
)

var volumeTemplateCmds = []Command{

	{
		Description:  "Lists available volume templates.",
		Subject:      "volume-template",
		AltSubject:   "vt",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list volume templates", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"format":     c.FlagSet.String("format", _nilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
				"local_only": c.FlagSet.Bool("local-only", false, "Show only templates that support local install"),
				"pxe_only":   c.FlagSet.Bool("pxe-only", false, "Show only templates that support pxe booting"),
			}
		},
		ExecuteFunc: volumeTemplatesListCmd,
		Endpoint:    ExtendedEndpoint,
	},
	{
		Description:  "Create volume templates.",
		Subject:      "volume-template",
		AltSubject:   "vt",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("create volume templates", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"drive_id":                   c.FlagSet.Int("id", _nilDefaultInt, red("(Required)") + " The id of the drive to create the volume template from"),
				"label":                      c.FlagSet.String("label", _nilDefaultStr, red("(Required)") + " The label of the volume template"),
				"description":                c.FlagSet.String("description", _nilDefaultStr, red("(Required)") + " The description of the volume template"),
				"display_name":               c.FlagSet.String("name", _nilDefaultStr, red("(Required)") + " The display name of the volume template"),
				"boot_type":                  c.FlagSet.String("boot-type", _nilDefaultStr, "The boot_type of the volume template. Possible values: 'uefi_only','legacy_only' "),
				"boot_methods_supported":     c.FlagSet.String("boot-methods-supported", _nilDefaultStr, "The boot_methods_supported of the volume template. Defaults to 'pxe_iscsi'."),
				"deprecation_status":         c.FlagSet.String("deprecation-status", _nilDefaultStr, "Deprecation status. Possible values: not_deprecated,deprecated_deny_provision,deprecated_allow_expand. Defaults to 'not_deprecated'."),
				"tags":                       c.FlagSet.String("tags", _nilDefaultStr, "The tags of the volume template, comma separated."),
				"os_bootstrap_function_name": c.FlagSet.String("os-bootstrap-function-name", _nilDefaultStr, "Optional property that selects the cloudinit configuration function. Can be one of: provisioner_os_cloudinit_prepare_centos, provisioner_os_cloudinit_prepare_rhel, provisioner_os_cloudinit_prepare_ubuntu, provisioner_os_cloudinit_prepare_windows."),
				"os_type":                    c.FlagSet.String("os-type", _nilDefaultStr, "Template operating system type. For example, Ubuntu or CentOS. If set, os-version and os-architecture flags are required as well."),
				"os_version":                 c.FlagSet.String("os-version", _nilDefaultStr, "Template operating system version. If set, os-type and os-architecture flags are required as well."),
				"os_architecture":            c.FlagSet.String("os-architecture", _nilDefaultStr, "Template operating system architecture.Possible values: none, unknown, x86, x86_64. If set, os-version and os-type flags are required as well."),
				"version":                    c.FlagSet.String("version", _nilDefaultStr, "Template version. Default value is 0.0.0"),
				"os_ready_method":            c.FlagSet.String("os-ready-method", _nilDefaultStr, "Possible values: 'wait_for_ssh', 'wait_for_signal_from_os'. Default value: 'wait_for_ssh'."),
				"return_id":                  c.FlagSet.Bool("return-id", false, "(Optional) Will print the ID of the created Volume Template. Useful for automating tasks."),
			}
		},
		ExecuteFunc: volumeTemplateCreateFromDriveCmd,
		Endpoint:    ExtendedEndpoint,
	},
	{
		Description:  "Allow other users of the platform to use the template.",
		Subject:      "volume-template",
		AltSubject:   "vt",
		Predicate:    "make-public",
		AltPredicate: "public",
		FlagSet:      flag.NewFlagSet("make volume template public", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"template_id_or_name":        c.FlagSet.String("id", _nilDefaultStr, "Volume template id or name"),
				"os_bootstrap_function_name": c.FlagSet.String("os-bootstrap-function-name", _nilDefaultStr, red("(Required)") + " Selects the cloudinit configuration function. Can be one of: provisioner_os_cloudinit_prepare_centos, provisioner_os_cloudinit_prepare_rhel, provisioner_os_cloudinit_prepare_ubuntu, provisioner_os_cloudinit_prepare_windows."),
			}
		},
		ExecuteFunc: volumeTemplateMakePublicCmd,
		Endpoint:    DeveloperEndpoint,
	},
	{
		Description:  "Stop other users of the platform from being able to use the template by allocating a specific owner.",
		Subject:      "volume-template",
		AltSubject:   "vt",
		Predicate:    "make-private",
		AltPredicate: "private",
		FlagSet:      flag.NewFlagSet("make volume template private", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"template_id_or_name": c.FlagSet.String("id", _nilDefaultStr, "Volume template id or name"),
				"user_id":             c.FlagSet.String("user-id", _nilDefaultStr, "New owner user id or email."),
			}
		},
		ExecuteFunc: volumeTemplateMakePrivateCmd,
		Endpoint:    DeveloperEndpoint,
	},
}

func volumeTemplatesListCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	vList, err := client.VolumeTemplates()
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
			FieldSize: 20,
		},
		{
			FieldName: "NAME",
			FieldType: tableformatter.TypeString,
			FieldSize: 20,
		},
		{
			FieldName: "SIZE",
			FieldType: tableformatter.TypeInt,
			FieldSize: 6,
		},
		{
			FieldName: "STATUS",
			FieldType: tableformatter.TypeString,
			FieldSize: 20,
		},
		{
			FieldName: "FLAGS",
			FieldType: tableformatter.TypeString,
			FieldSize: 10,
		},
	}

	localOnly := c.Arguments["local_only"] != nil && *c.Arguments["local_only"].(*bool)
	pxeOnly := c.Arguments["pxe_only"] != nil && *c.Arguments["pxe_only"].(*bool)

	data := [][]interface{}{}
	for _, v := range *vList {

		if localOnly && !v.VolumeTemplateLocalDiskSupported {
			continue
		}

		if pxeOnly && !strings.Contains(v.VolumeTemplateBootMethodsSupported, "pxe_iscsi") {
			continue
		}

		flags := []string{}

		flags = append(flags, strings.Split(v.VolumeTemplateBootMethodsSupported, ",")...)

		if v.VolumeTemplateLocalDiskSupported {
			flags = append(flags, "local")
		}

		data = append(data, []interface{}{
			v.VolumeTemplateID,
			v.VolumeTemplateLabel,
			v.VolumeTemplateDisplayName,
			v.VolumeTemplateSizeMBytes,
			v.VolumeTemplateDeprecationStatus,
			strings.Join(flags, ","),
		})

	}

	tableformatter.TableSorter(schema).OrderBy(schema[0].FieldName).Sort(data)

	table := tableformatter.Table{
		Data:   data,
		Schema: schema,
	}
	return table.RenderTable("Volume templates", "", getStringParam(c.Arguments["format"]))
}

func volumeTemplateCreateFromDriveCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {
	objVolumeTemplate := metalcloud.VolumeTemplate{
		VolumeTemplateDeprecationStatus:       getStringParam(c.Arguments["deprecation_status"]),
		VolumeTemplateBootMethodsSupported:    getStringParam(c.Arguments["boot_methods_supported"]),
		VolumeTemplateTags:                    strings.Split(getStringParam(c.Arguments["tags"]), ","),
		VolumeTemplateBootType:                getStringParam(c.Arguments["boot_type"]),
		VolumeTemplateOsBootstrapFunctionName: getStringParam(c.Arguments["os_bootstrap_function_name"]),
		VolumeTemplateVersion:                 getStringParam(c.Arguments["version"]),
		VolumeTemplateOSReadyMethod:           getStringParam(c.Arguments["os_ready_method"]),
	}

	driveID, ok := getIntParamOk(c.Arguments["drive_id"])
	if !ok {
		return "", fmt.Errorf("-id is required (drive id)")
	}

	if label, ok := getStringParamOk(c.Arguments["label"]); ok {
		objVolumeTemplate.VolumeTemplateLabel = label
	} else {
		return "", fmt.Errorf("-label is required")
	}

	if description, ok := getStringParamOk(c.Arguments["description"]); ok {
		objVolumeTemplate.VolumeTemplateDescription = description
	} else {
		return "", fmt.Errorf("-description is required")
	}

	if name, ok := getStringParamOk(c.Arguments["display_name"]); ok {
		objVolumeTemplate.VolumeTemplateDisplayName = name
	} else {
		return "", fmt.Errorf("-name is required")
	}

	os, err := getOperatingSystemFromCommand(c)

	if err != nil {
		return "", err
	}
	objVolumeTemplate.VolumeTemplateOperatingSystem = *os

	ret, err := client.VolumeTemplateCreateFromDrive(driveID, objVolumeTemplate)
	if err != nil {
		return "", err
	}

	if getBoolParam(c.Arguments["return_id"]) {
		return fmt.Sprintf("%d", ret.VolumeTemplateID), nil
	}

	return "", nil
}

func getOperatingSystemFromCommand(c *Command) (*metalcloud.OperatingSystem, error) {
	var operatingSystem = metalcloud.OperatingSystem{}
	present := false

	if osType, ok := getStringParamOk(c.Arguments["os_type"]); ok {
		present = true
		operatingSystem.OperatingSystemType = osType
	}

	if osVersion, ok := getStringParamOk(c.Arguments["os_version"]); ok {
		if !present {
			return nil, fmt.Errorf("some of the operating system flags are missing")
		}
		operatingSystem.OperatingSystemVersion = osVersion
	} else if present {
		return nil, fmt.Errorf("os-version is required")
	}

	if osArchitecture, ok := getStringParamOk(c.Arguments["os_architecture"]); ok {
		if !present {
			return nil, fmt.Errorf("some of the operating system flags are missing")
		}
		operatingSystem.OperatingSystemArchitecture = osArchitecture
	} else if present {
		return nil, fmt.Errorf("os-architecture is required")
	}

	return &operatingSystem, nil
}

func getNetworkOperatingSystemFromCommand(c *Command) (*metalcloud.NetworkOperatingSystem, error) {
	var operatingSystem = metalcloud.NetworkOperatingSystem{}
	present := false

	if osType, ok := getStringParamOk(c.Arguments["network_os_type"]); ok {
		present = true
		operatingSystem.OperatingSystemType = osType
	}

	if osVersion, ok := getStringParamOk(c.Arguments["network_os_version"]); ok {
		if !present {
			return nil, fmt.Errorf("some of the network operating system flags are missing")
		}
		operatingSystem.OperatingSystemVersion = osVersion
	} else if present {
		return nil, fmt.Errorf("network-os-version is required")
	}

	if osArchitecture, ok := getStringParamOk(c.Arguments["network_os_architecture"]); ok {
		if !present {
			return nil, fmt.Errorf("some of the network operating system flags are missing")
		}
		operatingSystem.OperatingSystemArchitecture = osArchitecture
	} else if present {
		return nil, fmt.Errorf("network-os-architecture is required")
	}

	if osVendor, ok := getStringParamOk(c.Arguments["network_os_vendor"]); ok {
		if !present {
			return nil, fmt.Errorf("some of the network operating system flags are missing")
		}
		operatingSystem.OperatingSystemVendor = osVendor
	} else if present {
		return nil, fmt.Errorf("network-os-vendor is required")
	}

	if osMachine, ok := getStringParamOk(c.Arguments["network_os_machine"]); ok {
		if !present {
			return nil, fmt.Errorf("some of the network operating system flags are missing")
		}
		operatingSystem.OperatingSystemMachine = osMachine
	} else if present {
		return nil, fmt.Errorf("network-os-machine is required")
	}
	if osMachineRevision, ok := getStringParamOk(c.Arguments["network_os_machine_revision"]); ok {
		if !present {
			return nil, fmt.Errorf("some of the network operating system flags are missing")
		}
		operatingSystem.OperatingSystemMachineRevision = osMachineRevision
	} else if present {
		return nil, fmt.Errorf("network-os-machine-revision is required")
	}
	return &operatingSystem, nil
}

func getVolumeTemplateFromCommand(paramName string, c *Command, client metalcloud.MetalCloudClient) (*metalcloud.VolumeTemplate, error) {

	v, err := getParam(c, "template_id_or_name", paramName)
	if err != nil {
		return nil, err
	}

	id, label, isID := idOrLabel(v)

	if isID {
		return client.VolumeTemplateGet(id)
	}

	list, err := client.VolumeTemplates()
	if err != nil {
		return nil, err
	}

	for _, s := range *list {
		if s.VolumeTemplateLabel == label {
			return &s, nil
		}
	}

	if isID {
		return nil, fmt.Errorf("volume template %d not found", id)
	}

	return nil, fmt.Errorf("volume template %s not found", label)
}

func volumeTemplateMakePublicCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {
	template, err := getVolumeTemplateFromCommand("id", c, client)

	if err != nil {
		return "", err
	}

	osBootstrapFunctionName, ok := getStringParamOk(c.Arguments["os_bootstrap_function_name"])
	if !ok {
		return "", fmt.Errorf("-os-bootstrap-function-name is required")
	}

	err = client.VolumeTemplateMakePublic(template.VolumeTemplateID, osBootstrapFunctionName)

	if err != nil {
		return "", err
	}

	return "", nil
}

func volumeTemplateMakePrivateCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {
	template, err := getVolumeTemplateFromCommand("id", c, client)

	if err != nil {
		return "", err
	}

	user, err := getUserFromCommand("user-id", c, client)
	if err != nil {
		return "", err
	}

	if err = client.VolumeTemplateMakePrivate(template.VolumeTemplateID, user.UserID); err != nil {
		return "", err
	}

	return "", nil
}
