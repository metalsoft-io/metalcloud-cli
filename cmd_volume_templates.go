package main

import (
	"flag"
	"strings"

	interfaces "github.com/bigstepinc/metalcloud-cli/interfaces"
)

var volumeTemplateyCmds = []Command{

	Command{
		Description:  "Lists available volume templates",
		Subject:      "volume_templates",
		AltSubject:   "vt",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list volume templates", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"format":     c.FlagSet.String("format", _nilDefaultStr, "The output format. Supported values are 'json','csv'. The default format is human readable."),
				"local_only": c.FlagSet.Bool("local_only", false, "Show only templates that support local install"),
				"pxe_only":   c.FlagSet.Bool("pxe_only", false, "Show only templates that support pxe booting"),
			}
		},
		ExecuteFunc: volumeTemplatesListCmd,
	},
}

func volumeTemplatesListCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	vList, err := client.VolumeTemplates()
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
			FieldName: "LABEL",
			FieldType: TypeString,
			FieldSize: 20,
		},
		SchemaField{
			FieldName: "NAME",
			FieldType: TypeString,
			FieldSize: 20,
		},
		SchemaField{
			FieldName: "SIZE",
			FieldType: TypeInt,
			FieldSize: 6,
		},
		SchemaField{
			FieldName: "STATUS",
			FieldType: TypeString,
			FieldSize: 20,
		},
		SchemaField{
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

	return renderTable("Volume templates", "", getStringParam(c.Arguments["format"]), data, schema)
}
