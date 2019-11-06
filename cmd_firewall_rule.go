package main

import (
	"flag"
	"fmt"
	"strings"
)

var firewallRuleCmds = []Command{

	Command{
		Description:  "Lists instance array firewall rules",
		Subject:      "firewall_rules",
		AltSubject:   "fw",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list firewall rules", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"instance_array_id": c.FlagSet.Int("ia", _nilDefaultInt, "(Required) The instance array id"),
				"format":            c.FlagSet.String("format", _nilDefaultStr, "The output format. Supported values are 'json','csv'. The default format is human readable."),
			}
		},
		ExecuteFunc: firewallRulesListCmd,
	},
}

func firewallRulesListCmd(c *Command, client MetalCloudClient) (string, error) {

	instanceArrayID := c.Arguments["instance_array_id"]

	if instanceArrayID == nil || *instanceArrayID.(*int) == 0 {
		return "", fmt.Errorf("-ia <instance_array_id> is required")
	}

	retIA, err := client.InstanceArrayGet(*instanceArrayID.(*int))
	if err != nil {
		return "", err
	}

	schema := []SchemaField{
		SchemaField{
			FieldName: "INDEX",
			FieldType: TypeInt,
			FieldSize: 6,
		},
		SchemaField{
			FieldName: "PROTOCOL",
			FieldType: TypeString,
			FieldSize: 10,
		},
		SchemaField{
			FieldName: "PORT",
			FieldType: TypeString,
			FieldSize: 10,
		},
		SchemaField{
			FieldName: "SOURCE",
			FieldType: TypeString,
			FieldSize: 20,
		},

		SchemaField{
			FieldName: "ENABLED",
			FieldType: TypeBool,
			FieldSize: 10,
		},
		SchemaField{
			FieldName: "DESC.",
			FieldType: TypeString,
			FieldSize: 50,
		},
	}

	status := retIA.InstanceArrayServiceStatus
	if retIA.InstanceArrayServiceStatus != "ordered" && retIA.InstanceArrayOperation.InstanceArrayDeployType == "edit" && retIA.InstanceArrayOperation.InstanceArrayDeployStatus == "not_started" {
		status = "edited"
	}

	list := retIA.InstanceArrayOperation.InstanceArrayFirewallRules
	data := [][]interface{}{}
	idx := 0
	for _, fw := range list {

		portRange := "any"

		if fw.FirewallRulePortRangeStart != 0 {
			portRange = fmt.Sprintf("%d", fw.FirewallRulePortRangeStart)
		}

		if fw.FirewallRulePortRangeStart != fw.FirewallRulePortRangeEnd {
			portRange += fmt.Sprintf("-%d", fw.FirewallRulePortRangeEnd)
		}

		sourceIPRange := "any"

		if fw.FirewallRuleSourceIPAddressRangeStart != "" {
			sourceIPRange = fw.FirewallRuleSourceIPAddressRangeStart
		}

		if fw.FirewallRuleSourceIPAddressRangeStart != fw.FirewallRuleSourceIPAddressRangeEnd {
			sourceIPRange += fmt.Sprintf("-%s", fw.FirewallRuleSourceIPAddressRangeEnd)
		}

		data = append(data, []interface{}{
			idx,
			fw.FirewallRuleProtocol,
			portRange,
			sourceIPRange,
			fw.FirewallRuleEnabled,
			fw.FirewallRuleDescription,
		})

		idx++

	}

	var sb strings.Builder

	format := c.Arguments["format"]
	if format == nil {
		var f string
		f = ""
		format = &f
	}

	switch *format.(*string) {
	case "json", "JSON":
		ret, err := GetTableAsJSONString(data, schema)
		if err != nil {
			return "", err
		}
		sb.WriteString(ret)
	case "csv", "CSV":
		ret, err := GetTableAsCSVString(data, schema)
		if err != nil {
			return "", err
		}
		sb.WriteString(ret)

	default:
		sb.WriteString(fmt.Sprintf("Instance Array %s (%d) [%s] has the following firewall rules:\n", retIA.InstanceArrayLabel, retIA.InstanceArrayID, status))

		sb.WriteString(GetTableAsString(data, schema))

		sb.WriteString(fmt.Sprintf("Total: %d firewall rules\n\n", len(list)))
	}

	return sb.String(), nil
}
