package vm

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"

	metalcloud2 "github.com/metalsoft-io/metal-cloud-sdk2-go"
	"github.com/metalsoft-io/metalcloud-cli/internal/colors"
	"github.com/metalsoft-io/metalcloud-cli/internal/command"
	"github.com/metalsoft-io/tableformatter"
	"gopkg.in/yaml.v3"
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
				"show_credentials": c.FlagSet.Bool("show-credentials", false, colors.Green("(Flag)")+" If set returns the VM pool credentials."),
			}
		},
		ExecuteFunc2:        vmPoolListCmd,
		PermissionsRequired: []string{command.VM_POOLS_READ},
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
				"show_credentials": c.FlagSet.Bool("show-credentials", false, colors.Green("(Flag)")+" If set returns the VM pool credentials."),
			}
		},
		ExecuteFunc2:        vmPoolGetCmd,
		PermissionsRequired: []string{command.VM_POOLS_READ},
	},
	{
		Description:  "Create VM pool.",
		Subject:      "vm-pool",
		AltSubject:   "vm-pool",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("Create VM pool", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"format":                c.FlagSet.String("format", "json", "The input format. Supported values are 'json','yaml'. The default format is json."),
				"read_config_from_file": c.FlagSet.String("raw-config", command.NilDefaultStr, colors.Red("(Required)")+" Read raw object from file"),
				"read_config_from_pipe": c.FlagSet.Bool("pipe", false, colors.Green("(Flag)")+" If set, read raw object from pipe instead of from a file. Either this flag or the --raw-config option must be used."),
				"return_id":             c.FlagSet.Bool("return-id", false, "Will print the ID of the created object. Useful for automating tasks."),
			}
		},
		ExecuteFunc2:        vmPoolCreateCmd,
		PermissionsRequired: []string{command.VM_POOLS_WRITE},
	},
}

func vmPoolListCmd(ctx context.Context, c *command.Command, client *metalcloud2.APIClient) (string, error) {
	vmPools, response, err := client.VMPoolsApi.GetVMPools(ctx)
	if err != nil {
		return "", err
	}
	if response.StatusCode >= 400 {
		return "", fmt.Errorf("HTTP error: %s", response.Status)
	}

	showCredentials := command.GetBoolParam(c.Arguments["show_credentials"])
	format := command.GetStringParam(c.Arguments["format"])

	rawData := []metalcloud2.VmPool{}
	data := [][]interface{}{}

	statusCounts := map[string]int{
		"active":      0,
		"maintenance": 0,
	}

	for _, vmPool := range vmPools.Data {
		vmPool.Links = nil
		if !showCredentials {
			vmPool.Certificate = ""
		}

		rawData = append(rawData, vmPool)

		statusCounts[vmPool.Status] = statusCounts[vmPool.Status] + 1
		if vmPool.InMaintenance == 1.0 {
			statusCounts["maintenance"]++
		}

		capacityRam := formattedCapacity(
			float64(vmPool.UsedRamGB)/float64(vmPool.TotalRamGB),
			fmt.Sprintf("%.2f GB RAM used out of %.2f GB total", float64(vmPool.UsedRamGB), float64(vmPool.TotalRamGB)),
			format)

		capacityStorage := formattedCapacity(
			float64(vmPool.UsedRamGB)/float64(vmPool.TotalRamGB),
			fmt.Sprintf("%.2f GB RAM used out of %.2f GB total", float64(vmPool.UsedRamGB), float64(vmPool.TotalRamGB)),
			format)

		row := []interface{}{
			vmPool.Id,
			formattedStatus(vmPool.Status, format),
			vmPool.Name,
			vmPool.ManagementHost,
			capacityRam,
			capacityStorage,
			vmPool.DatacenterName,
		}
		if showCredentials {
			row = append(row, []interface{}{
				vmPool.Certificate,
			}...)
		}

		data = append(data, row)
	}

	switch format {
	case "json", "JSON":
		result, err := json.MarshalIndent(rawData, "", "  ")
		if err != nil {
			return "", err
		}
		return string(result), nil

	case "yaml", "YAML":
		result, err := yaml.Marshal(rawData)
		if err != nil {
			return "", err
		}
		return string(result), nil

	default:
		table := tableformatter.Table{
			Data:   data,
			Schema: fieldSchema(showCredentials),
		}

		title := fmt.Sprintf("VM pools: %d active %d maintenance",
			statusCounts["active"],
			statusCounts["maintenance"])

		return table.RenderTable(title, "", format)
	}
}

func vmPoolGetCmd(ctx context.Context, c *command.Command, client *metalcloud2.APIClient) (string, error) {
	vmPoolId, ok := command.GetIntParamOk(c.Arguments["vm_pool_id"])
	if !ok {
		return "", fmt.Errorf("-id is required")
	}

	vmPool, response, err := client.VMPoolsApi.GetVMPool(ctx, float64(vmPoolId))
	if err != nil {
		return "", err
	}
	if response.StatusCode >= 400 {
		return "", fmt.Errorf("HTTP error: %s", response.Status)
	}

	showCredentials := command.GetBoolParam(c.Arguments["show_credentials"])
	format := command.GetStringParam(c.Arguments["format"])

	vmPool.Links = nil
	if !showCredentials {
		vmPool.Certificate = ""
	}

	switch format {
	case "json", "JSON":
		result, err := json.MarshalIndent(vmPool, "", "  ")
		if err != nil {
			return "", err
		}
		return string(result), nil

	case "yaml", "YAML":
		result, err := yaml.Marshal(vmPool)
		if err != nil {
			return "", err
		}
		return string(result), nil

	default:
		data := [][]interface{}{}

		capacityRam := formattedCapacity(
			float64(vmPool.UsedRamGB)/float64(vmPool.TotalRamGB),
			fmt.Sprintf("%.2f GB RAM used out of %.2f GB total", float64(vmPool.UsedRamGB), float64(vmPool.TotalRamGB)),
			format)

		capacityStorage := formattedCapacity(
			float64(vmPool.UsedRamGB)/float64(vmPool.TotalRamGB),
			fmt.Sprintf("%.2f GB RAM used out of %.2f GB total", float64(vmPool.UsedRamGB), float64(vmPool.TotalRamGB)),
			format)

		row := []interface{}{
			vmPool.Id,
			formattedStatus(vmPool.Status, format),
			vmPool.Name,
			vmPool.ManagementHost,
			capacityRam,
			capacityStorage,
			vmPool.DatacenterName,
		}
		if showCredentials {
			row = append(row, []interface{}{
				vmPool.Certificate,
			}...)
		}

		data = append(data, row)

		table := tableformatter.Table{
			Data:   data,
			Schema: fieldSchema(showCredentials),
		}

		return table.RenderTable("VM pool", "", format)
	}
}

func vmPoolCreateCmd(ctx context.Context, c *command.Command, client *metalcloud2.APIClient) (string, error) {
	var obj metalcloud2.CreateVmPool

	err := command.GetRawObjectFromCommand(c, &obj)
	if err != nil {
		return "", err
	}

	vmPool, _, err := client.VMPoolsApi.CreateVMPool(ctx, obj)
	if err != nil {
		return "", err
	}

	if command.GetBoolParam(c.Arguments["return_id"]) {
		return fmt.Sprintf("%v", vmPool.Id), nil
	}

	return "", err
}

func fieldSchema(showCredentials bool) []tableformatter.SchemaField {
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

	if showCredentials {
		schema = append(schema, tableformatter.SchemaField{
			FieldName: "CERT",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		})
	}

	return schema
}

func formattedStatus(status string, format string) string {
	if format != "" {
		return status
	}

	switch status {
	case "active":
		return colors.Blue(status)
	case "maintenance":
		return colors.Green(status)
	case "":
		return colors.Green(status)
	default:
		return colors.Yellow(status)
	}
}

func formattedCapacity(usedPercentage float64, capacity string, format string) string {
	if format != "" {
		return capacity
	}

	if usedPercentage >= 0.8 {
		capacity = colors.Red(capacity)
	} else if usedPercentage >= 0.5 {
		capacity = colors.Red(capacity)
	} else {
		capacity = colors.Green(capacity)
	}

	return capacity
}
