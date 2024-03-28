package storage

import (
	"flag"
	"fmt"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
	"github.com/metalsoft-io/metalcloud-cli/internal/colors"
	"github.com/metalsoft-io/metalcloud-cli/internal/command"
	"github.com/metalsoft-io/metalcloud-cli/internal/configuration"
	"github.com/metalsoft-io/metalcloud-cli/internal/filtering"
	"github.com/metalsoft-io/tableformatter"
)

var StorageCmds = []command.Command{
	{
		Description:  "Lists all storage.",
		Subject:      "storage",
		AltSubject:   "storages",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list storage", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"format":              c.FlagSet.String("format", command.NilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
				"filter":              c.FlagSet.String("filter", "*", "filter to use when searching for servers. Check the documentation for examples. Defaults to '*'"),
				"show_credentials":    c.FlagSet.Bool("show-credentials", false, colors.Green("(Flag)")+" If set returns the servers' IPMI credentials. (Slow for large queries)"),
				"show_decommissioned": c.FlagSet.Bool("show-decommissioned", false, colors.Green("(Flag)")+" If set returns decommissioned servers which are normally hidden"),
			}
		},
		ExecuteFunc: storageListCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
}

func storageListCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {

	filter := command.GetStringParam(c.Arguments["filter"])

	list, err := client.StoragePoolSearch(filtering.ConvertToSearchFieldFormat(filter))
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

	if command.GetBoolParam(c.Arguments["show_credentials"]) {

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
		"active":         0,
		"maintenance":    0,
		"decommissioned": 0,
	}

	for _, s := range *list {

		if s.StoragePoolStatus == "decommissioned" && !command.GetBoolParam(c.Arguments["show_decommissioned"]) {
			continue
		}

		statusCounts[s.StoragePoolStatus] = statusCounts[s.StoragePoolStatus] + 1

		if s.StoragePoolInMaintenance == true {
			statusCounts["maintenance"]++
		}
		credentialsUser := ""
		credentialsPass := ""

		if command.GetBoolParam(c.Arguments["show_credentials"]) {

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
			status = colors.Blue(status)
		case "maintenance":
			status = colors.Green(status)
		case "":
			status = colors.Green(status)

		default:
			status = colors.Yellow(status)

		}

		usedPercentage := float64(s.StoragePoolCapacityTotalCachedRealMbytes-s.StoragePoolCapacityFreeCachedRealMbytes) / float64(s.StoragePoolCapacityTotalCachedRealMbytes)

		capacity := fmt.Sprintf("%.2f TB physically used out of %.2f TB total, %0.2f TB virtually allocated",
			float64(s.StoragePoolCapacityTotalCachedRealMbytes-s.StoragePoolCapacityFreeCachedRealMbytes)/(1024*1024),
			float64(s.StoragePoolCapacityTotalCachedRealMbytes)/(1024*1024),
			float64(s.StoragePoolCapacityUsedCachedVirtualMbytes)/(1024*1024),
		)

		if usedPercentage >= 0.8 {
			capacity = colors.Red(capacity)
		} else if usedPercentage >= 0.5 {
			capacity = colors.Red(capacity)
		} else {
			capacity = colors.Green(capacity)
		}

		row := []interface{}{
			s.StoragePoolID,
			status,
			s.StoragePoolName,
			s.StoragePoolEndpoint,
			capacity,
			s.DatacenterName,
		}
		if command.GetBoolParam(c.Arguments["show_credentials"]) {
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
		statusCounts["active"],
		statusCounts["maintenance"])

	if command.GetBoolParam(c.Arguments["show_decommissioned"]) {
		title = title + fmt.Sprintf(" %d decommissioned", statusCounts["decommissioned"])
	}

	return table.RenderTable(title, "", command.GetStringParam(c.Arguments["format"]))
}
