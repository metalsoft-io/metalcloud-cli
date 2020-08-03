package main

import (
	"flag"
	"fmt"
	"strings"

	interfaces "github.com/bigstepinc/metalcloud-cli/interfaces"
)

var volumeTemplateyCmds = []Command{

	{
		Description:  "Lists available volume templates",
		Subject:      "volume-template",
		AltSubject:   "vt",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list volume templates", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"format":     c.FlagSet.String("format", _nilDefaultStr, "The output format. Supported values are 'json','csv'. The default format is human readable."),
				"local_only": c.FlagSet.Bool("local-only", false, "Show only templates that support local install"),
				"pxe_only":   c.FlagSet.Bool("pxe-only", false, "Show only templates that support pxe booting"),
			}
		},
		ExecuteFunc: volumeTemplatesListCmd,
		Endpoint:    ExtendedEndpoint,
	},
	{
		Description:  "Create volume templates",
		Subject:      "volume-template",
		AltSubject:   "vt",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("create volume templates", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"drive_id":               c.FlagSet.Int("id", _nilDefaultInt, "(Required) The id of the drive to create the volume template from"),
				"label":                  c.FlagSet.String("label", _nilDefaultStr, "(Required) The label of the volume template"),
				"description":            c.FlagSet.String("description", _nilDefaultStr, "(Required) The description of the volume template"),
				"display_name":           c.FlagSet.String("name", _nilDefaultStr, "(Required) The display name of the volume template"),
				"boot_type":              c.FlagSet.String("boot-type", _nilDefaultStr, "(Required) The boot_type of the volume template. Possible values: 'uefi_only','legacy_only','hybrid' "),
				"boot_methods_supported": c.FlagSet.String("boot-methods-supported", _nilDefaultStr, "The boot_methods_supported of the volume template. Defaults to 'pxe_iscsi'."),
				"deprecation_status":     c.FlagSet.String("deprecation-status", _nilDefaultStr, "Deprecation status. Possible values: not_deprecated,deprecated_deny_provision,deprecated_allow_expand. Defaults to 'not_deprecated'."),
				"tags":                   c.FlagSet.String("tags", _nilDefaultStr, "The tags of the volume template, comma separated."),
				"return_id":              c.FlagSet.Bool("return-id", false, "(Optional) Will print the ID of the created Drive Array. Useful for automating tasks."),
			}
		},
		ExecuteFunc: volumeTemplateCreateCmd,
		Endpoint:    ExtendedEndpoint,
	},
}

func volumeTemplatesListCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	vList, err := client.VolumeTemplates()
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
			FieldName: "LABEL",
			FieldType: TypeString,
			FieldSize: 20,
		},
		{
			FieldName: "NAME",
			FieldType: TypeString,
			FieldSize: 20,
		},
		{
			FieldName: "SIZE",
			FieldType: TypeInt,
			FieldSize: 6,
		},
		{
			FieldName: "STATUS",
			FieldType: TypeString,
			FieldSize: 20,
		},
		{
			FieldName: "FLAGS",
			FieldType: TypeString,
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

	TableSorter(schema).OrderBy(schema[0].FieldName).Sort(data)

	return renderTable("Volume templates", "", getStringParam(c.Arguments["format"]), data, schema)
}

func volumeTemplateCreateCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	driveID, ok := getIntParamOk(c.Arguments["drive_id"])
	if !ok {
		return "", fmt.Errorf("-id is required (drive id)")
	}

	label, ok := getStringParamOk(c.Arguments["label"])
	if !ok {
		return "", fmt.Errorf("-label is required ")
	}

	description := getStringParam(c.Arguments["label"])
	name := getStringParam(c.Arguments["display_name"])
	bootType := getStringParam(c.Arguments["boot_type"])
	deprecationStatus := getStringParam(c.Arguments["deprecation_status"])
	bootMethodSupported := getStringParam(c.Arguments["boot_methods_supported"])
	tags := strings.Split(getStringParam(c.Arguments["tags"]), ",")

	if deprecationStatus == "" {
		deprecationStatus = "not_deprecated"
	}

	if bootMethodSupported == "" {
		bootMethodSupported = "pxe_iscsi"
	}

	ret, err := client.VolumeTemplateCreate(driveID, label, description, name, bootType, deprecationStatus, bootMethodSupported, tags)
	if err != nil {
		return "", err
	}

	if getBoolParam(c.Arguments["return_id"]) {
		return fmt.Sprintf("%d", ret.VolumeTemplateID), nil
	}

	return "", nil
}
