package vm

import (
	"flag"
	"fmt"

	metalcloud2 "github.com/metalsoft-io/metal-cloud-sdk2-go"
	"github.com/metalsoft-io/metalcloud-cli/internal/colors"
	"github.com/metalsoft-io/metalcloud-cli/internal/command"
	"github.com/metalsoft-io/tableformatter"
)

var VmPoolsCmds = []command.Command{
	{
		Description:  "Lists all VM pools.",
		Subject:      "vm-pool",
		AltSubject:   "vm-pools",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("List VM pools", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"format":           c.FlagSet.String("format", command.NilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
				"show_credentials": c.FlagSet.Bool("show-credentials", false, colors.Green("(Flag)")+" If set returns the VM pool credentials. (Slow for large queries)"),
			}
		},
		ExecuteFunc2:        vmPoolListCmd,
		PermissionsRequired: []string{command.STORAGE_READ},
	},
	{
		Description:  "Get VM pool.",
		Subject:      "vm-pool",
		AltSubject:   "vm-pool",
		Predicate:    "get",
		AltPredicate: "show",
		FlagSet:      flag.NewFlagSet("Get VM pool", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"vm_pool_id":       c.FlagSet.Int("id", 0, "VM pool id"),
				"format":           c.FlagSet.String("format", command.NilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
				"show_credentials": c.FlagSet.Bool("show-credentials", false, colors.Green("(Flag)")+" If set returns the VM pool credentials. (Slow for large queries)"),
			}
		},
		ExecuteFunc2:        vmPoolGetCmd,
		PermissionsRequired: []string{command.STORAGE_READ},
	},
}

func vmPoolListCmd(c *command.Command, client *metalcloud2.APIClient) (string, error) {
	vmPools, response, err := client.VMPoolsApi.InventoryController1GetVMPools(nil)
	if err != nil {
		return "", err
	}
	if response.StatusCode >= 400 {
		return "", fmt.Errorf("HTTP error: %s", response.Status)
	}

	schema := []tableformatter.SchemaField{
		{
			FieldName: "ID",
			FieldType: tableformatter.TypeFloat,
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
			FieldName: "HOST",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},

		{
			FieldName: "RAM",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},

		{
			FieldName: "STORAGE",
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
			FieldName: "CERT",
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

	for _, vmPool := range vmPools.Data {
		if vmPool.Status == "decommissioned" && !command.GetBoolParam(c.Arguments["show_decommissioned"]) {
			continue
		}

		statusCounts[vmPool.Status] = statusCounts[vmPool.Status] + 1
		if vmPool.InMaintenance == 1.0 {
			statusCounts["maintenance"]++
		}

		status := vmPool.Status
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

		usedRamPercentage := float64(vmPool.UsedRamGB) / float64(vmPool.TotalRamGB)

		capacityRam := fmt.Sprintf("%.2f GB RAM used out of %.2f GB total",
			float64(vmPool.UsedRamGB),
			float64(vmPool.TotalRamGB),
		)

		if usedRamPercentage >= 0.8 {
			capacityRam = colors.Red(capacityRam)
		} else if usedRamPercentage >= 0.5 {
			capacityRam = colors.Red(capacityRam)
		} else {
			capacityRam = colors.Green(capacityRam)
		}

		usedStoragePercentage := float64(vmPool.UsedRamGB) / float64(vmPool.TotalRamGB)

		capacityStorage := fmt.Sprintf("%.2f GB RAM used out of %.2f GB total",
			float64(vmPool.UsedRamGB),
			float64(vmPool.TotalRamGB),
		)

		if usedStoragePercentage >= 0.8 {
			capacityStorage = colors.Red(capacityStorage)
		} else if usedStoragePercentage >= 0.5 {
			capacityStorage = colors.Red(capacityStorage)
		} else {
			capacityStorage = colors.Green(capacityStorage)
		}

		row := []interface{}{
			vmPool.Id,
			status,
			vmPool.Name,
			vmPool.ManagementHost,
			capacityRam,
			capacityStorage,
			vmPool.DatacenterName,
		}
		if command.GetBoolParam(c.Arguments["show_credentials"]) {
			row = append(row, []interface{}{
				vmPool.Certificate,
			}...)
		}

		data = append(data, row)
	}

	table := tableformatter.Table{
		Data:   data,
		Schema: schema,
	}

	title := fmt.Sprintf("VM pools: %d active %d maintenance",
		statusCounts["active"],
		statusCounts["maintenance"])

	if command.GetBoolParam(c.Arguments["show_decommissioned"]) {
		title = title + fmt.Sprintf(" %d decommissioned", statusCounts["decommissioned"])
	}

	return table.RenderTable(title, "", command.GetStringParam(c.Arguments["format"]))
}

func vmPoolGetCmd(c *command.Command, client *metalcloud2.APIClient) (string, error) {
	vmPoolId, ok := command.GetIntParamOk(c.Arguments["vm_pool_id"])
	if !ok {
		return "", fmt.Errorf("-id is required")
	}

	vmPool, response, err := client.VMPoolsApi.InventoryController1GetVMPool(nil, float64(vmPoolId))
	if err != nil {
		return "", err
	}
	if response.StatusCode >= 400 {
		return "", fmt.Errorf("HTTP error: %s", response.Status)
	}

	schema := []tableformatter.SchemaField{
		{
			FieldName: "ID",
			FieldType: tableformatter.TypeFloat,
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
			FieldName: "HOST",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},

		{
			FieldName: "RAM",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},

		{
			FieldName: "STORAGE",
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
			FieldName: "CERT",
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

	status := vmPool.Status
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

	usedRamPercentage := float64(vmPool.UsedRamGB) / float64(vmPool.TotalRamGB)

	capacityRam := fmt.Sprintf("%.2f GB RAM used out of %.2f GB total",
		float64(vmPool.UsedRamGB),
		float64(vmPool.TotalRamGB),
	)

	if usedRamPercentage >= 0.8 {
		capacityRam = colors.Red(capacityRam)
	} else if usedRamPercentage >= 0.5 {
		capacityRam = colors.Red(capacityRam)
	} else {
		capacityRam = colors.Green(capacityRam)
	}

	usedStoragePercentage := float64(vmPool.UsedRamGB) / float64(vmPool.TotalRamGB)

	capacityStorage := fmt.Sprintf("%.2f GB RAM used out of %.2f GB total",
		float64(vmPool.UsedRamGB),
		float64(vmPool.TotalRamGB),
	)

	if usedStoragePercentage >= 0.8 {
		capacityStorage = colors.Red(capacityStorage)
	} else if usedStoragePercentage >= 0.5 {
		capacityStorage = colors.Red(capacityStorage)
	} else {
		capacityStorage = colors.Green(capacityStorage)
	}

	row := []interface{}{
		vmPool.Id,
		status,
		vmPool.Name,
		vmPool.ManagementHost,
		capacityRam,
		capacityStorage,
		vmPool.DatacenterName,
	}
	if command.GetBoolParam(c.Arguments["show_credentials"]) {
		row = append(row, []interface{}{
			vmPool.Certificate,
		}...)
	}

	data = append(data, row)

	table := tableformatter.Table{
		Data:   data,
		Schema: schema,
	}

	title := fmt.Sprintf("VM pools: %d active %d maintenance",
		statusCounts["active"],
		statusCounts["maintenance"])

	if command.GetBoolParam(c.Arguments["show_decommissioned"]) {
		title = title + fmt.Sprintf(" %d decommissioned", statusCounts["decommissioned"])
	}

	return table.RenderTable(title, "", command.GetStringParam(c.Arguments["format"]))
}
