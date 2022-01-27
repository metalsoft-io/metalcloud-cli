package main

import (
	"flag"
	"fmt"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
	"github.com/metalsoft-io/tableformatter"
)

var storageCmds = []Command{

	{
		Description:  "Lists all storage.",
		Subject:      "storage",
		AltSubject:   "storages",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list storage", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"format":              c.FlagSet.String("format", _nilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
				"filter":              c.FlagSet.String("filter", "*", "filter to use when searching for servers. Check the documentation for examples. Defaults to '*'"),
				"show_credentials":    c.FlagSet.Bool("show-credentials", false, "(Flag) If set returns the servers' IPMI credentials. (Slow for large queries)"),
				"show_decommissioned": c.FlagSet.Bool("show-decommissioned", false, "(Flag) If set returns decommissioned servers which are normally hidden"),
			}
		},
		ExecuteFunc: storageListCmd,
		Endpoint:    DeveloperEndpoint,
	},
}

func storageListCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	filter := getStringParam(c.Arguments["filter"])

	list, err := client.StoragePoolSearch(filter)
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
			FieldName: "NAME",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},

		{
			FieldName: "ENDPOINT",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "CAPACITY",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "DATACENTER",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
	}

	if getBoolParam(c.Arguments["show_credentials"]) {

		schema = append(schema, tableformatter.SchemaField{
			FieldName: "USER",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		})

		schema = append(schema, tableformatter.SchemaField{
			FieldName: "PASS",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		})

	}

	data := [][]interface{}{}

	statusCounts := map[string]int{
		"available":      0,
		"maintenance":    0,
		"decommissioned": 0,
	}

	for _, s := range *list {

		if s.StoragePoolStatus == "decommissioned" && !getBoolParam(c.Arguments["show_decommissioned"]) {
			continue
		}

		statusCounts[s.StoragePoolStatus] = statusCounts[s.StoragePoolStatus] + 1

		credentialsUser := ""
		credentialsPass := ""

		if getBoolParam(c.Arguments["show_credentials"]) {

			storage, err := client.StoragePoolGet(s.StoragePoolID, true)

			if err != nil {
				return "", err
			}

			credentialsUser = fmt.Sprintf("%s", storage.StoragePoolUsername)
			credentialsPass = fmt.Sprintf("%s", storage.StoragePoolPassword)

		}

		status := s.StoragePoolStatus

		switch status {
		case "active":
			status = blue(status)
		case "maintenance":
			status = green(status)
		case "":
			status = green(status)

		default:
			status = yellow(status)

		}

		usedPercentage := float64(s.StoragePoolCapacityTotalCachedRealMbytes-s.StoragePoolCapacityFreeCachedRealMbytes) / float64(s.StoragePoolCapacityTotalCachedRealMbytes)

		capacity := fmt.Sprintf("%.2f TB physically used out of %.2f TB total, %0.2f TB virtually allocated",
			float64(s.StoragePoolCapacityTotalCachedRealMbytes-s.StoragePoolCapacityFreeCachedRealMbytes)/(1024*1024),
			float64(s.StoragePoolCapacityTotalCachedRealMbytes)/(1024*1024),
			float64(s.StoragePoolCapacityUsedCachedVirtualMbytes)/(1024*1024),
		)

		if usedPercentage >= 0.8 {
			capacity = red(capacity)
		} else if usedPercentage >= 0.5 {
			capacity = red(capacity)
		} else {
			capacity = green(capacity)
		}

		row := []interface{}{
			s.StoragePoolID,
			status,
			s.StoragePoolName,
			s.StoragePoolEndpoint,
			capacity,
			s.DatacenterName,
		}
		if getBoolParam(c.Arguments["show_credentials"]) {
			row = append(row, []interface{}{
				credentialsUser,
				credentialsPass,
			}...)
		}

		data = append(data, row)

	}

	table := tableformatter.Table{
		Data:   data,
		Schema: schema,
	}

	title := fmt.Sprintf("Storage pools: %d active %d maintenance",
		statusCounts["available"],
		statusCounts["maintenance"])

	if getBoolParam(c.Arguments["show_decommissioned"]) {
		title = title + fmt.Sprintf(" %d decommissioned", statusCounts["decommissioned"])
	}

	return table.RenderTable(title, "", getStringParam(c.Arguments["format"]))
}
